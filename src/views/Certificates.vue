<script setup>
import { onMounted } from 'vue'
import { useCertStore } from '../stores/certificates'
import Card from '../components/common/Card.vue'
import Btn from '../components/common/Btn.vue'
import Table from '../components/common/Table.vue'

const store = useCertStore()

const columns = [
  { key: 'property', label: '属性' },
  { key: 'value', label: '值' }
]

onMounted(() => {
  store.fetchCAInfo()
  store.fetchCertificates()
})

const certColumns = [
  { key: 'cert', label: '证书' }
]

const formatCertData = (certs) => {
  if (!certs) return []
  const lines = certs.split('-----BEGIN CERTIFICATE-----')
  return lines.filter(l => l.trim()).map(l => ({ cert: '-----BEGIN CERTIFICATE-----' + l.trim() }))
}
</script>

<template>
  <div class="cert-page">
    <Card title="CA信息" :padding="false">
      <template #actions>
        <Btn size="sm" @click="store.fetchCAInfo()">刷新</Btn>
      </template>

      <div v-if="store.loading" class="loading">加载中...</div>

      <div v-else-if="store.caInfo" class="ca-info">
        <Table :columns="columns" :data="Object.entries(store.caInfo).filter(([k]) => k !== 'root_certificate' && k !== 'intermediate_certificate').map(([k, v]) => ({ property: k, value: v }))" />

        <div class="cert-section">
          <h4>根证书</h4>
          <pre class="cert-content">{{ store.caInfo.root_certificate }}</pre>
        </div>

        <div class="cert-section">
          <h4>中间证书</h4>
          <pre class="cert-content">{{ store.caInfo.intermediate_certificate }}</pre>
        </div>
      </div>

      <div v-else class="empty">暂无CA信息</div>
    </Card>

    <Card title="证书链" :padding="false">
      <template #actions>
        <Btn size="sm" @click="store.fetchCertificates()">刷新</Btn>
      </template>

      <div v-if="store.loading" class="loading">加载中...</div>

      <div v-else-if="store.certificates" class="cert-chain">
        <pre class="cert-content">{{ store.certificates }}</pre>
      </div>

      <div v-else class="empty">暂无证书</div>
    </Card>
  </div>
</template>

<style scoped>
.cert-page {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.ca-info {
  padding: 0;
}

.cert-section {
  padding: 16px 20px;
  border-top: 1px solid var(--color-border);
}

.cert-section h4 {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 12px;
}

.cert-content {
  background: var(--color-background);
  padding: 12px;
  border-radius: 6px;
  font-family: 'Fira Code', monospace;
  font-size: 12px;
  white-space: pre-wrap;
  word-break: break-all;
  max-height: 200px;
  overflow-y: auto;
}

.loading,
.empty {
  text-align: center;
  padding: 40px;
  color: #94a3b8;
}
</style>