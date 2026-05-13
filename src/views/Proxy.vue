<script setup>
import { ref, onMounted } from 'vue'
import { useProxyStore } from '../stores/proxy'
import { caddyApi } from '../services/caddyApi'
import Card from '../components/common/Card.vue'
import Btn from '../components/common/Btn.vue'
import Message from '../components/common/Message.vue'
import Table from '../components/common/Table.vue'
import Badge from '../components/common/Badge.vue'
import ProxyConfig from '../components/proxy/ProxyConfig.vue'

const store = useProxyStore()
const messages = ref([])
const configSnapshot = ref(null)

const upstreamCols = [
  { key: 'address', label: '地址' },
  { key: 'num_requests', label: '活跃请求' },
  { key: 'fails', label: '失败次数' },
  { key: 'status', label: '状态' }
]

const showMessage = (type, msg) => {
  messages.value.push({ type, message: msg, id: Date.now() })
  setTimeout(() => messages.value.shift(), 3000)
}

const getStatus = (up) => {
  if (up.fails > 3) return { variant: 'error', text: '异常' }
  if (up.num_requests > 0) return { variant: 'success', text: '活跃' }
  return { variant: 'default', text: '空闲' }
}

const fetchConfigData = async () => {
  try {
    const response = await caddyApi.getConfig('/apps/http/servers')
    configSnapshot.value = { apps: { http: { servers: response.data } } }
  } catch (e) {
    if (e.response?.status === 404) {
      configSnapshot.value = null
    } else {
      showMessage('error', '获取配置失败: ' + (e.message || '未知错误'))
    }
  }
}

const handleSave = async (payload) => {
  try {
    if (payload.action === 'fullConfig') {
      await caddyApi.loadConfig(payload.data)
      showMessage('success', '全部配置已保存')
    } else if (payload.action === 'addHandler') {
      await caddyApi.setConfig(payload.path, payload.data)
      showMessage('success', '已添加 reverse_proxy 处理器')
    } else if (payload.action === 'add') {
      await caddyApi.setConfig(`${payload.path}/...`, [payload.data])
      showMessage('success', '上游已添加')
    } else if (payload.action === 'update') {
      await caddyApi.updateConfig(`${payload.path}/${payload.index}`, payload.data)
      showMessage('success', '上游已更新')
    } else if (payload.action === 'delete') {
      await caddyApi.deleteConfig(`${payload.path}/${payload.index}`)
      showMessage('success', '上游已删除')
    }
    await fetchConfigData()
    await store.fetchUpstreams()
  } catch (e) {
    showMessage('error', '操作失败: ' + (e.response?.data?.error || e.message))
  }
}

onMounted(() => {
  store.fetchUpstreams()
  fetchConfigData()
})
</script>

<template>
  <div class="proxy-page">
    <div v-for="(msg, idx) in messages" :key="msg.id">
      <Message :type="msg.type" :message="msg.message" @close="messages.splice(idx, 1)" />
    </div>

    <ProxyConfig
      v-if="configSnapshot"
      :config="configSnapshot"
      @save="handleSave"
      @refresh="fetchConfigData"
    />

    <Card v-else title="提示">
      <p>未获取到 HTTP 服务器配置。请确保 Caddy 已运行并配置了 HTTP 服务器。</p>
      <Btn size="sm" @click="fetchConfigData" class="retry-btn">重试</Btn>
    </Card>

    <Card title="上游实时状态" :padding="false">
      <template #actions>
        <Btn size="sm" @click="store.fetchUpstreams()">刷新</Btn>
      </template>

      <div v-if="store.loading" class="loading">加载中...</div>

      <Table
        v-else-if="store.upstreams.length > 0"
        :columns="upstreamCols"
        :data="store.upstreams"
      >
        <template #address="{ value }">
          <code class="upstream-address">{{ value }}</code>
        </template>
        <template #num_requests="{ value }">
          <span class="request-count">{{ value }}</span>
        </template>
        <template #fails="{ value }">
          <Badge :variant="value > 3 ? 'error' : value > 0 ? 'warning' : 'default'" size="md">
            {{ value }}
          </Badge>
        </template>
        <template #status="{ row }">
          <Badge :variant="getStatus(row).variant" size="md">
            {{ getStatus(row).text }}
          </Badge>
        </template>
      </Table>

      <div v-else class="empty">暂无上游数据</div>
    </Card>
  </div>
</template>

<style scoped>
.proxy-page {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.upstream-address {
  font-family: 'Fira Code', monospace;
  font-size: 13px;
  color: var(--color-primary);
}

.request-count {
  font-weight: 600;
  font-size: 14px;
}

.loading, .empty {
  text-align: center;
  padding: 40px;
  color: #94a3b8;
}

.retry-btn {
  margin-top: 12px;
}
</style>