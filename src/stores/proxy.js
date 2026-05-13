import { defineStore } from 'pinia'
import { ref } from 'vue'
import { caddyApi } from '../services/caddyApi'

export const useProxyStore = defineStore('proxy', () => {
  const upstreams = ref([])
  const config = ref(null)
  const loading = ref(false)
  const error = ref(null)
  const configPath = ref('/apps/http/servers')

  const fetchUpstreams = async () => {
    loading.value = true
    error.value = null
    try {
      const response = await caddyApi.getReverseProxyUpstreams()
      upstreams.value = response.data
    } catch (e) {
      error.value = e.message || '获取上游信息失败'
    } finally {
      loading.value = false
    }
  }

  const fetchConfig = async (path = configPath.value) => {
    loading.value = true
    error.value = null
    configPath.value = path
    try {
      const response = await caddyApi.getConfig(path)
      config.value = response.data
    } catch (e) {
      error.value = e.response?.data?.error || '获取配置失败'
    } finally {
      loading.value = false
    }
  }

  const addUpstream = async (upstreamPath, data) => {
    loading.value = true
    error.value = null
    try {
      if (Array.isArray(data)) {
        await caddyApi.setConfig(`${upstreamPath}/...`, data)
      } else {
        await caddyApi.setConfig(upstreamPath, data)
      }
      await fetchUpstreams()
      return true
    } catch (e) {
      error.value = e.response?.data?.error || '添加上游失败'
      return false
    } finally {
      loading.value = false
    }
  }

  const updateUpstream = async (upstreamPath, index, data) => {
    loading.value = true
    error.value = null
    try {
      await caddyApi.updateConfig(`${upstreamPath}/${index}`, data)
      await fetchUpstreams()
      return true
    } catch (e) {
      error.value = e.response?.data?.error || '更新上游失败'
      return false
    } finally {
      loading.value = false
    }
  }

  const deleteUpstream = async (upstreamPath, index) => {
    loading.value = true
    error.value = null
    try {
      await caddyApi.deleteConfig(`${upstreamPath}/${index}`)
      await fetchUpstreams()
      return true
    } catch (e) {
      error.value = e.response?.data?.error || '删除上游失败'
      return false
    } finally {
      loading.value = false
    }
  }

  const saveProxyConfig = async (path, config) => {
    loading.value = true
    error.value = null
    try {
      await caddyApi.setConfig(path, config)
      await fetchUpstreams()
      return true
    } catch (e) {
      error.value = e.response?.data?.error || '保存配置失败'
      return false
    } finally {
      loading.value = false
    }
  }

  return {
    upstreams, config, loading, error, configPath,
    fetchUpstreams, fetchConfig, addUpstream, updateUpstream, deleteUpstream, saveProxyConfig
  }
})