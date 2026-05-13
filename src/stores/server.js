import { defineStore } from 'pinia'
import { ref } from 'vue'
import { caddyApi } from '../services/caddyApi'

export const useServerStore = defineStore('server', () => {
  const status = ref('unknown')
  const loading = ref(false)
  const error = ref(null)

  const checkStatus = async () => {
    loading.value = true
    try {
      await caddyApi.getConfig('/')
      status.value = 'running'
    } catch {
      status.value = 'stopped'
    } finally {
      loading.value = false
    }
  }

  const stopServer = async () => {
    loading.value = true
    error.value = null
    try {
      await caddyApi.stopServer()
      status.value = 'stopped'
      return true
    } catch (e) {
      error.value = e.message || '停止服务器失败'
      return false
    } finally {
      loading.value = false
    }
  }

  const loadConfig = async (config, contentType = 'application/json') => {
    loading.value = true
    error.value = null
    try {
      await caddyApi.loadConfig(config, contentType)
      status.value = 'running'
      return true
    } catch (e) {
      error.value = e.response?.data?.error || '加载配置失败'
      return false
    } finally {
      loading.value = false
    }
  }

  const adaptConfig = async (config, contentType = 'text/caddyfile') => {
    loading.value = true
    error.value = null
    try {
      const response = await caddyApi.adaptConfig(config, contentType)
      return { success: true, data: response.data }
    } catch (e) {
      error.value = e.response?.data?.error || '适配配置失败'
      return { success: false, error: error.value }
    } finally {
      loading.value = false
    }
  }

  return {
    status,
    loading,
    error,
    checkStatus,
    stopServer,
    loadConfig,
    adaptConfig
  }
})