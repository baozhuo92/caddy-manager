<script setup>
import { ref, watch } from 'vue'
import Btn from '../common/Btn.vue'
import Input from '../common/Input.vue'

const props = defineProps({
  site: { type: Object, default: null }
})

const emit = defineEmits(['save', 'cancel'])

const form = ref(emptyForm())
const tab = ref('basic')

function emptyForm() {
  return {
    domain: '',
    port: '',
    scheme: '',
    routeEntries: [{ path: '', upstreams: [''], websocket: false }],
    cors: { enabled: false, origins: [''], methods: ['*'], headers: ['*'], credentials: false },
    basicauth: { enabled: false, users: [{ username: '', password: '' }] },
    customHeaders: { request: [], response: [] }
  }
}

watch(() => props.site, (val) => {
  if (val) {
    const entries = val.routeEntries?.length
      ? val.routeEntries.map(e => ({
          path: e.path || '',
          upstreams: e.upstreams?.length ? [...e.upstreams] : [''],
          websocket: e.websocket || false
        }))
      : [{ path: '', upstreams: val.upstreams?.length ? [...val.upstreams] : [''], websocket: val.websocket || false }]
    form.value = {
      domain: val.domain || '',
      port: val.port || '',
      scheme: val.scheme || '',
      routeEntries: entries,
      cors: {
        enabled: val.cors?.enabled || false,
        origins: val.cors?.origins?.length ? [...val.cors.origins] : [''],
        methods: val.cors?.methods?.length ? [...val.cors.methods] : ['*'],
        headers: val.cors?.headers?.length ? [...val.cors.headers] : ['*'],
        credentials: val.cors?.credentials || false
      },
      basicauth: {
        enabled: val.basicauth?.enabled || false,
        users: val.basicauth?.users?.length ? val.basicauth.users.map(u => ({ ...u })) : [{ username: '', password: '' }]
      },
      customHeaders: {
        request: val.customHeaders?.request?.length ? val.customHeaders.request.map(h => ({ ...h })) : [],
        response: val.customHeaders?.response?.length ? val.customHeaders.response.map(h => ({ ...h })) : []
      }
    }
  } else {
    form.value = emptyForm()
  }
}, { immediate: true })

function addRouteEntry() {
  form.value.routeEntries.push({ path: '', upstreams: [''], websocket: false })
}
function removeRouteEntry(idx) {
  form.value.routeEntries.splice(idx, 1)
  if (form.value.routeEntries.length === 0) addRouteEntry()
}
function addEntryUpstream(eIdx) {
  form.value.routeEntries[eIdx].upstreams.push('')
}
function removeEntryUpstream(eIdx, uIdx) {
  const u = form.value.routeEntries[eIdx].upstreams
  u.splice(uIdx, 1)
  if (u.length === 0) u.push('')
}
function addCorsOrigin() { form.value.cors.origins.push('') }
function removeCorsOrigin(idx) {
  form.value.cors.origins.splice(idx, 1)
  if (form.value.cors.origins.length === 0) form.value.cors.origins.push('')
}
function addUser() { form.value.basicauth.users.push({ username: '', password: '' }) }
function removeUser(idx) { form.value.basicauth.users.splice(idx, 1) }
function addReqHeader() { form.value.customHeaders.request.push({ key: '', value: '' }) }
function addResHeader() { form.value.customHeaders.response.push({ key: '', value: '' }) }
function removeHeader(arr, idx) { arr.splice(idx, 1) }

function handleSave() {
  const data = {
    domain: form.value.domain,
    port: form.value.port,
    scheme: form.value.scheme,
    routeEntries: form.value.routeEntries.map(e => ({
      path: e.path.trim(),
      upstreams: e.upstreams.filter(d => d.trim()),
      websocket: e.websocket
    })).filter(e => e.upstreams.length > 0),
    cors: form.value.cors,
    basicauth: form.value.basicauth,
    customHeaders: {
      request: form.value.customHeaders.request.filter(h => h.key),
      response: form.value.customHeaders.response.filter(h => h.key)
    }
  }

  if (props.site) {
    data.id = props.site.id
    data.serverName = props.site.serverName
  }

  emit('save', data)
}

const isValid = ref(true)
watch(form, () => {
  const hasDomain = form.value.domain.trim().length > 0
  const hasPort = form.value.port.trim().length > 0
  const hasAnyUpstream = form.value.routeEntries.some(e => e.upstreams.some(d => d.trim().length > 0))
  isValid.value = (hasDomain || hasPort) && hasAnyUpstream
}, { deep: true })
</script>

<template>
  <div class="site-form">
    <div class="tabs">
      <button class="tab" :class="{ active: tab === 'basic' }" @click="tab = 'basic'">基础</button>
      <button class="tab" :class="{ active: tab === 'routes' }" @click="tab = 'routes'">路由</button>
      <button class="tab" :class="{ active: tab === 'cors' }" @click="tab = 'cors'">CORS</button>
      <button class="tab" :class="{ active: tab === 'auth' }" @click="tab = 'auth'">认证</button>
      <button class="tab" :class="{ active: tab === 'headers' }" @click="tab = 'headers'">自定义头</button>
    </div>

    <div class="form-body">
      <div v-show="tab === 'basic'" class="form-section">
        <p class="required-hint">* 至少需要填写域名或端口</p>
        <div class="field-group">
          <label class="field-label">域名</label>
          <Input v-model="form.domain" placeholder="例如: example.com" />
          <span class="field-hint">留空则仅监听端口</span>
        </div>
        <div class="field-group">
          <label class="field-label">端口</label>
          <Input v-model="form.port" placeholder="例如: 19000" />
          <span class="field-hint">留空则默认 80/443，不设域名仅端口则监听 `:端口`</span>
        </div>
        <div class="field-group">
          <label class="field-label">协议</label>
          <select v-model="form.scheme" class="select">
            <option value="">HTTP / HTTPS (全部)</option>
            <option value="http">仅 HTTP</option>
            <option value="https">仅 HTTPS</option>
          </select>
        </div>
      </div>

      <div v-show="tab === 'routes'" class="form-section">
        <div class="section-header">
          <span class="section-title">路由规则</span>
          <Btn size="sm" @click="addRouteEntry">添加路由</Btn>
        </div>
        <div v-for="(entry, eIdx) in form.routeEntries" :key="eIdx" class="route-card">
          <div class="route-header">
            <span class="route-index">路由 {{ eIdx + 1 }}</span>
            <button class="remove-btn" @click="removeRouteEntry(eIdx)" v-if="form.routeEntries.length > 1">×</button>
          </div>
          <div class="route-body">
            <div class="field-group">
              <label class="field-label">匹配路径</label>
              <Input v-model="entry.path" placeholder='例如 /api/*（留空匹配全部）' />
            </div>
            <label class="checkbox-label">
              <input type="checkbox" v-model="entry.websocket" />
              <span>WebSocket 支持</span>
            </label>
            <div class="sub-section-header">
              <span class="sub-label">上游服务器</span>
              <Btn size="sm" variant="ghost" @click="addEntryUpstream(eIdx)">添加上游</Btn>
            </div>
            <div v-for="(_, uIdx) in entry.upstreams" :key="uIdx" class="array-item">
              <Input v-model="entry.upstreams[uIdx]" :placeholder="`例如: localhost:${3000 + uIdx}`" />
              <button class="remove-btn" @click="removeEntryUpstream(eIdx, uIdx)" v-if="entry.upstreams.length > 1">×</button>
            </div>
          </div>
        </div>
      </div>

      <div v-show="tab === 'cors'" class="form-section">
        <label class="checkbox-label">
          <input type="checkbox" v-model="form.cors.enabled" />
          <span>启用 CORS</span>
        </label>
        <template v-if="form.cors.enabled">
          <div class="section-header">
            <span class="section-title">允许的域名 (Origin)</span>
            <Btn size="sm" @click="addCorsOrigin">添加</Btn>
          </div>
          <div v-for="(_, idx) in form.cors.origins" :key="idx" class="array-item">
            <Input v-model="form.cors.origins[idx]" placeholder="* 或 https://example.com" />
            <button class="remove-btn" @click="removeCorsOrigin(idx)" v-if="form.cors.origins.length > 1">×</button>
          </div>
          <div class="field-group">
            <label class="field-label">允许的方法 (逗号分隔)</label>
            <Input v-model="form.cors.methods" placeholder="GET, POST, PUT, DELETE, OPTIONS" />
          </div>
          <div class="field-group">
            <label class="field-label">允许的请求头 (逗号分隔)</label>
            <Input v-model="form.cors.headers" placeholder="Content-Type, Authorization" />
          </div>
          <label class="checkbox-label">
            <input type="checkbox" v-model="form.cors.credentials" />
            <span>允许携带凭证 (Credentials)</span>
          </label>
        </template>
      </div>

      <div v-show="tab === 'auth'" class="form-section">
        <label class="checkbox-label">
          <input type="checkbox" v-model="form.basicauth.enabled" />
          <span>启用 BasicAuth 认证</span>
        </label>
        <template v-if="form.basicauth.enabled">
          <div class="section-header">
            <span class="section-title">用户</span>
            <Btn size="sm" @click="addUser">添加用户</Btn>
          </div>
          <div v-for="(_, idx) in form.basicauth.users" :key="idx" class="auth-user-row">
            <Input v-model="form.basicauth.users[idx].username" placeholder="用户名" />
            <Input v-model="form.basicauth.users[idx].password" type="password" placeholder="密码" />
            <button class="remove-btn" @click="removeUser(idx)">×</button>
          </div>
          <p class="hint">密码将发送到 Caddy API，由 Caddy 自动加密存储。</p>
        </template>
      </div>

      <div v-show="tab === 'headers'" class="form-section">
        <div class="sub-section">
          <div class="section-header">
            <span class="section-title">请求头 (Request)</span>
            <Btn size="sm" @click="addReqHeader">添加</Btn>
          </div>
          <div v-for="(_, idx) in form.customHeaders.request" :key="idx" class="header-row">
            <Input v-model="form.customHeaders.request[idx].key" placeholder="Header 名" />
            <Input v-model="form.customHeaders.request[idx].value" placeholder="值" />
            <button class="remove-btn" @click="removeHeader(form.customHeaders.request, idx)">×</button>
          </div>
        </div>
        <div class="sub-section">
          <div class="section-header">
            <span class="section-title">响应头 (Response)</span>
            <Btn size="sm" @click="addResHeader">添加</Btn>
          </div>
          <div v-for="(_, idx) in form.customHeaders.response" :key="idx" class="header-row">
            <Input v-model="form.customHeaders.response[idx].key" placeholder="Header 名" />
            <Input v-model="form.customHeaders.response[idx].value" placeholder="值" />
            <button class="remove-btn" @click="removeHeader(form.customHeaders.response, idx)">×</button>
          </div>
        </div>
      </div>
    </div>

    <div class="form-footer">
      <Btn variant="ghost" @click="$emit('cancel')">取消</Btn>
      <Btn variant="primary" :disabled="!isValid" @click="handleSave">
        {{ site ? '保存修改' : '创建站点' }}
      </Btn>
    </div>
  </div>
</template>

<style scoped>
.site-form {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.tabs {
  display: flex;
  border-bottom: 2px solid var(--color-border);
  margin: -24px -24px 20px;
  padding: 0 24px;
  flex-shrink: 0;
}

.tab {
  padding: 12px 16px;
  border: none;
  background: none;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  color: #64748b;
  border-bottom: 2px solid transparent;
  margin-bottom: -2px;
  transition: all 0.2s;
}

.tab.active {
  color: var(--color-primary);
  border-bottom-color: var(--color-primary);
}

.form-body {
  flex: 1;
  overflow-y: auto;
}

.form-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.field-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.field-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text);
}

.field-hint {
  font-size: 12px;
  color: #94a3b8;
}

.required-hint {
  font-size: 13px;
  color: var(--color-error);
  margin: 0;
  padding: 8px 12px;
  background: #fef2f2;
  border-radius: 6px;
  border: 1px solid #fecaca;
}

.select {
  padding: 10px 12px;
  border: 1px solid var(--color-border);
  border-radius: 6px;
  font-size: 14px;
  background: var(--color-surface);
  width: 100%;
}

.select:focus {
  outline: none;
  border-color: var(--color-primary);
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 14px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 8px;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text);
}

.array-item {
  display: flex;
  gap: 8px;
  align-items: flex-start;
}

.array-item > * {
  flex: 1;
}

.remove-btn {
  width: 32px;
  height: 38px;
  border: 1px solid transparent;
  border-radius: 6px;
  background: none;
  color: var(--color-error);
  font-size: 18px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.remove-btn:hover {
  background: #fef2f2;
}

.auth-user-row {
  display: flex;
  gap: 8px;
  align-items: flex-start;
}

.auth-user-row > * {
  flex: 1;
}

.header-row {
  display: flex;
  gap: 8px;
  align-items: flex-start;
}

.header-row > * {
  flex: 1;
}

.sub-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px 0;
}

.sub-section + .sub-section {
  border-top: 1px solid var(--color-border);
}

.hint {
  font-size: 12px;
  color: #64748b;
  margin: 0;
}

.form-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 16px;
  margin-top: 16px;
  border-top: 1px solid var(--color-border);
  flex-shrink: 0;
}

.route-card {
  border: 1px solid var(--color-border);
  border-radius: 8px;
  overflow: hidden;
}

.route-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: var(--color-background);
  border-bottom: 1px solid var(--color-border);
}

.route-index {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text);
}

.route-body {
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.sub-section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.sub-label {
  font-size: 13px;
  color: #64748b;
}
</style>
