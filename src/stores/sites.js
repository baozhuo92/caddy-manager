import { defineStore } from 'pinia'
import { ref } from 'vue'
import { caddyApi } from '../services/caddyApi'

let nextId = 1

function nextServerName(servers) {
  const names = Object.keys(servers || {})
  let idx = 0
  while (names.includes(`srv${idx}`)) idx++
  return `srv${idx}`
}

function extractPort(listen) {
  if (!listen || listen.length !== 1) return ''
  const addr = listen[0]
  const m = addr.match(/:(\d+)$/)
  if (!m) return ''
  const port = m[1]
  return (port === '80' || port === '443') ? '' : port
}

function parseSites(servers) {
  const result = []
  for (const [sName, server] of Object.entries(servers || {})) {
    const rpRoutes = (server.routes || []).filter(r =>
      (r.handle || []).some(h => h.handler === 'reverse_proxy')
    )
    if (rpRoutes.length === 0) continue

    const firstRp = rpRoutes[0]
    const match = firstRp.match?.[0]
    const hosts = match?.host || []
    const handle = firstRp.handle || []
    const headersH = handle.find(h => h.handler === 'headers')
    const authH = handle.find(h => h.handler === 'authentication')

    const routeEntries = rpRoutes.map(r => {
      const h = r.handle || []
      const rp = h.find(h => h.handler === 'reverse_proxy')
      const m = r.match?.[0]
      return {
        path: m?.path?.[0] || '',
        upstreams: (rp?.upstreams || []).map(u => u.dial),
        websocket: rp?.transport?.protocol === 'http'
      }
    })

    result.push({
      id: `s${nextId++}`,
      name: hosts[0] || sName,
      domain: hosts[0] || '',
      port: extractPort(server.listen),
      scheme: match?.scheme || '',
      routeEntries,
      serverName: sName,
      serverPath: `/apps/http/servers/${sName}`,
      listen: server.listen || [],
      cors: parseCors(headersH),
      basicauth: parseAuth(authH),
      customHeaders: parseHeaders(headersH)
    })
  }
  return result
}

function parseCors(h) {
  if (!h) return { enabled: false, origins: [''], methods: ['*'], headers: ['*'], credentials: false }
  const set = h.response?.set || {}
  const o = set['Access-Control-Allow-Origin']
  return {
    enabled: !!o,
    origins: o || [''],
    methods: set['Access-Control-Allow-Methods'] || ['*'],
    headers: set['Access-Control-Allow-Headers'] || ['*'],
    credentials: set['Access-Control-Allow-Credentials'] === 'true'
  }
}

function parseAuth(h) {
  if (!h) return { enabled: false, users: [{ username: '', password: '' }] }
  const p = h.providers?.basicauth
  if (!p) return { enabled: false, users: [{ username: '', password: '' }] }
  const users = Object.entries(p.users || {}).map(([u, pwd]) => ({ username: u, password: pwd }))
  return { enabled: users.length > 0, users: users.length ? users : [{ username: '', password: '' }] }
}

function parseHeaders(h) {
  if (!h) return { request: [], response: [] }
  const corsKeys = [
    'Access-Control-Allow-Origin', 'Access-Control-Allow-Methods',
    'Access-Control-Allow-Headers', 'Access-Control-Allow-Credentials',
    'Access-Control-Expose-Headers', 'Access-Control-Max-Age'
  ]
  const req = Object.entries(h.request?.set || {}).map(([k, v]) => ({ key: k, value: fmt(v) }))
  const res = Object.entries(h.response?.set || {}).filter(([k]) => !corsKeys.includes(k)).map(([k, v]) => ({ key: k, value: fmt(v) }))
  return { request: req, response: res }
}

function fmt(v) { return Array.isArray(v) ? v.join(', ') : String(v) }
function parseVal(v) { return typeof v === 'string' && v.includes(',') ? v.split(',').map(s => s.trim()) : v }

function buildListen(d) {
  const domain = d.domain?.trim()
  const port = d.port?.trim()

  if (domain && port && port !== '80' && port !== '443') return [`${domain}:${port}`]
  if (port && port !== '80' && port !== '443') return [`:${port}`]
  if (domain) return []
  return []
}

function buildSiteHandlers(d) {
  const handle = []
  const hasCors = d.cors?.enabled && d.cors?.origins?.some(o => o.trim())
  const hasReqH = d.customHeaders?.request?.some(h => h.key)
  const hasResH = d.customHeaders?.response?.some(h => h.key)

  if (hasCors || hasReqH || hasResH) {
    const reqSet = {}
    const resSet = {}
    if (hasCors) {
      resSet['Access-Control-Allow-Origin'] = d.cors.origins.filter(o => o.trim())
      resSet['Access-Control-Allow-Methods'] = d.cors.methods || ['*']
      resSet['Access-Control-Allow-Headers'] = d.cors.headers || ['*']
      if (d.cors.credentials) resSet['Access-Control-Allow-Credentials'] = 'true'
    }
    for (const h of (d.customHeaders?.request || [])) { if (h.key) reqSet[h.key] = parseVal(h.value) }
    for (const h of (d.customHeaders?.response || [])) { if (h.key) resSet[h.key] = parseVal(h.value) }
    const hd = { handler: 'headers' }
    if (Object.keys(reqSet).length) hd.request = { set: reqSet }
    if (Object.keys(resSet).length) hd.response = { set: resSet }
    handle.push(hd)
  }

  if (d.basicauth?.enabled && d.basicauth.users?.some(u => u.username)) {
    const users = {}
    for (const u of d.basicauth.users) { if (u.username) users[u.username] = u.password || '' }
    handle.push({ handler: 'authentication', providers: { basicauth: { users } } })
  }
  return handle
}

function buildRoutes(d) {
  const siteHandlers = buildSiteHandlers(d)
  const domain = d.domain?.trim()
  const routes = []

  for (const entry of (d.routeEntries || [])) {
    const upstreams = (entry.upstreams || []).filter(u => u.trim())
    if (upstreams.length === 0) continue

    const rp = { handler: 'reverse_proxy', upstreams: upstreams.map(u => ({ dial: u })) }
    if (entry.websocket) rp.transport = { protocol: 'http' }

    const route = { handle: [...siteHandlers, rp], terminal: true }
    const match = {}
    if (domain) match.host = [domain]
    if (entry.path?.trim()) match.path = [entry.path.trim()]
    if (d.scheme && d.scheme !== 'all') match.scheme = d.scheme
    if (Object.keys(match).length) route.match = [match]

    routes.push(route)
  }
  return routes
}

export const useSitesStore = defineStore('sites', () => {
  const servers = ref({})
  const sites = ref([])
  const loading = ref(false)
  const error = ref(null)

  async function loadSites() {
    loading.value = true
    error.value = null
    try {
      const r = await caddyApi.getConfig('/apps/http/servers')
      servers.value = r.data || {}
      sites.value = parseSites(servers.value)
    } catch (e) {
      if (e.response?.status !== 404) error.value = e.message
      servers.value = {}
      sites.value = []
    } finally { loading.value = false }
  }

  async function addSite(d) {
    loading.value = true
    error.value = null
    try {
      const sName = nextServerName(servers.value)
      const serverCfg = { listen: buildListen(d), routes: buildRoutes(d) }
      await caddyApi.createConfig(`/apps/http/servers/${sName}`, serverCfg)
      await loadSites()
      return true
    } catch (e) {
      error.value = e.response?.data?.error || e.message || '添加失败'
      return false
    } finally { loading.value = false }
  }

  async function updateSite(d) {
    loading.value = true
    error.value = null
    const sName = d.serverName
    let backup = null
    try {
      const backupResp = await caddyApi.getConfig(`/apps/http/servers/${sName}`)
      backup = backupResp.data
    } catch (_) {}

    try {
      const serverCfg = { listen: buildListen(d), routes: buildRoutes(d) }
      await caddyApi.deleteConfig(`/apps/http/servers/${sName}`)
      await caddyApi.createConfig(`/apps/http/servers/${sName}`, serverCfg)
      await loadSites()
      return true
    } catch (e) {
      if (backup) {
        try { await caddyApi.createConfig(`/apps/http/servers/${sName}`, backup) } catch (_) {}
      }
      error.value = e.response?.data?.error || e.message || '更新失败'
      return false
    } finally { loading.value = false }
  }

  async function deleteSite(s) {
    loading.value = true
    error.value = null
    try {
      await caddyApi.deleteConfig(s.serverPath)
      await loadSites()
      return true
    } catch (e) {
      error.value = e.response?.data?.error || e.message || '删除失败'
      return false
    } finally { loading.value = false }
  }

  return { servers, sites, loading, error, loadSites, addSite, updateSite, deleteSite }
})
