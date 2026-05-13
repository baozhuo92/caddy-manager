import { defineStore } from 'pinia'
import { ref } from 'vue'
import { caddyApi } from '../services/caddyApi'

export const useConfigStore = defineStore('config', () => {
  const config = ref(null)
  const loading = ref(false)
  const error = ref(null)
  const currentPath = ref('/')

  const fetchConfig = async (path = '/') => {
    loading.value = true
    error.value = null
    currentPath.value = path
    try {
      const response = await caddyApi.getConfig(path)
      config.value = response.data
    } catch (e) {
      error.value = e.message || '获取配置失败'
    } finally {
      loading.value = false
    }
  }

  const updateConfig = async (path, data) => {
    loading.value = true
    error.value = null
    try {
      await caddyApi.setConfig(path, data)
      await fetchConfig(currentPath.value)
      return true
    } catch (e) {
      error.value = e.response?.data?.error || '更新配置失败'
      return false
    } finally {
      loading.value = false
    }
  }

  const createConfig = async (path, data) => {
    loading.value = true
    error.value = null
    try {
      await caddyApi.createConfig(path, data)
      await fetchConfig(currentPath.value)
      return true
    } catch (e) {
      error.value = e.response?.data?.error || '创建配置失败'
      return false
    } finally {
      loading.value = false
    }
  }

  const patchConfig = async (path, data) => {
    loading.value = true
    error.value = null
    try {
      await caddyApi.updateConfig(path, data)
      await fetchConfig(currentPath.value)
      return true
    } catch (e) {
      error.value = e.response?.data?.error || '修改配置失败'
      return false
    } finally {
      loading.value = false
    }
  }

  const deleteConfig = async (path) => {
    loading.value = true
    error.value = null
    try {
      await caddyApi.deleteConfig(path)
      await fetchConfig(currentPath.value)
      return true
    } catch (e) {
      error.value = e.response?.data?.error || '删除配置失败'
      return false
    } finally {
      loading.value = false
    }
  }

  const loadConfig = async (data, contentType = 'application/json') => {
    loading.value = true
    error.value = null
    try {
      await caddyApi.loadConfig(data, contentType)
      await fetchConfig('/')
      return true
    } catch (e) {
      error.value = e.response?.data?.error || '加载配置失败'
      return false
    } finally {
      loading.value = false
    }
  }

  const adaptConfig = async (data, contentType = 'text/caddyfile') => {
    loading.value = true
    error.value = null
    try {
      const response = await caddyApi.adaptConfig(data, contentType)
      return { success: true, data: response.data }
    } catch (e) {
      error.value = e.response?.data?.error || '适配配置失败'
      return { success: false, error: error.value }
    } finally {
      loading.value = false
    }
  }

  return {
    config,
    loading,
    error,
    currentPath,
    fetchConfig,
    updateConfig,
    createConfig,
    patchConfig,
    deleteConfig,
    loadConfig,
    adaptConfig
  }
})