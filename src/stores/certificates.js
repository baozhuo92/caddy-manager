import { defineStore } from 'pinia'
import { ref } from 'vue'
import { caddyApi } from '../services/caddyApi'

export const useCertStore = defineStore('certificates', () => {
  const caInfo = ref(null)
  const certificates = ref(null)
  const loading = ref(false)
  const error = ref(null)

  const fetchCAInfo = async (id = 'local') => {
    loading.value = true
    error.value = null
    try {
      const response = await caddyApi.getPKICA(id)
      caInfo.value = response.data
    } catch (e) {
      error.value = e.message || '获取CA信息失败'
    } finally {
      loading.value = false
    }
  }

  const fetchCertificates = async (id = 'local') => {
    loading.value = true
    error.value = null
    try {
      const response = await caddyApi.getPKICertificates(id)
      certificates.value = response.data
    } catch (e) {
      error.value = e.message || '获取证书失败'
    } finally {
      loading.value = false
    }
  }

  return {
    caInfo,
    certificates,
    loading,
    error,
    fetchCAInfo,
    fetchCertificates
  }
})