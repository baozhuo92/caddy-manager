import axios from 'axios'


const DEFAULT_BASE_URL = '/api/'

const api = axios.create({
  baseURL: DEFAULT_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

export const caddyApi = {
  getConfig(path = '/') {
    return api.get(`/config${path}`)
  },

  setConfig(path, data) {
    return api.post(`/config${path}`, data)
  },

  createConfig(path, data) {
    return api.put(`/config${path}`, data)
  },

  updateConfig(path, data) {
    return api.patch(`/config${path}`, data)
  },

  deleteConfig(path) {
    return api.delete(`/config${path}`)
  },

  loadConfig(data, contentType = 'application/json') {
    return api.post('/load', data, {
      headers: { 'Content-Type': contentType }
    })
  },

  stopServer() {
    return api.post('/stop')
  },

  adaptConfig(data, contentType = 'text/caddyfile') {
    return api.post('/adapt', data, {
      headers: { 'Content-Type': contentType }
    })
  },

  getPKICA(id = 'local') {
    return api.get(`/pki/ca/${id}`)
  },

  getPKICertificates(id = 'local') {
    return api.get(`/pki/ca/${id}/certificates`)
  },

  getReverseProxyUpstreams() {
    return api.get('/reverse_proxy/upstreams')
  }
}

export default caddyApi