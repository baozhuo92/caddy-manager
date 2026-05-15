<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { caddyApi } from '../services/caddyApi'
import Card from '../components/common/Card.vue'
import Btn from '../components/common/Btn.vue'
import Badge from '../components/common/Badge.vue'

const router = useRouter()

const serverInfo = ref(null)
const isConnected = ref(false)
const loading = ref(true)

const quickActions = [
  { label: '配置管理', path: '/conf', icon: 'config', desc: '查看和编辑配置' },
  { label: '服务器控制', path: '/servers', icon: 'server', desc: '启动/停止服务' },
  { label: '证书管理', path: '/certificates', icon: 'cert', desc: '查看CA和证书' },
  { label: '反向代理', path: '/proxy', icon: 'proxy', desc: '监控上游状态' }
]

const fetchServerInfo = async () => {
  loading.value = true
  try {
    const response = await caddyApi.getConfig('/')
    serverInfo.value = response.data
    isConnected.value = true
  } catch {
    isConnected.value = false
    serverInfo.value = null
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchServerInfo()
})
</script>

<template>
  <div class="dashboard">
    <div class="welcome-section">
      <h1 class="welcome-title">欢迎使用 Caddy Manager</h1>
      <p class="welcome-desc">管理您的 Caddy Web 服务器</p>
    </div>

    <div class="status-section">
      <Card>
        <div class="status-content">
          <div class="status-info">
            <span class="status-label">服务器状态</span>
            <Badge :variant="isConnected ? 'success' : 'error'" size="lg">
              {{ isConnected ? '运行中' : '未连接' }}
            </Badge>
          </div>
          <Btn size="sm" @click="fetchServerInfo">刷新状态</Btn>
        </div>
      </Card>
    </div>

    <div class="quick-actions">
      <h2 class="section-title">快捷操作</h2>
      <div class="action-grid">
        <div
          v-for="action in quickActions"
          :key="action.path"
          class="action-card"
          @click="router.push(action.path)"
        >
          <div class="action-icon">
            <svg v-if="action.icon === 'config'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
              <polyline points="14 2 14 8 20 8"/>
            </svg>
            <svg v-else-if="action.icon === 'server'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="2" y="2" width="20" height="8" rx="2" ry="2"/>
              <rect x="2" y="14" width="20" height="8" rx="2" ry="2"/>
              <line x1="6" y1="6" x2="6.01" y2="6"/>
              <line x1="6" y1="18" x2="6.01" y2="18"/>
            </svg>
            <svg v-else-if="action.icon === 'cert'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <line x1="12" y1="8" x2="12" y2="12"/>
              <line x1="12" y1="16" x2="12.01" y2="16"/>
            </svg>
            <svg v-else-if="action.icon === 'proxy'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="16 18 22 12 16 6"/>
              <polyline points="8 6 2 12 8 18"/>
            </svg>
          </div>
          <div class="action-text">
            <span class="action-label">{{ action.label }}</span>
            <span class="action-desc">{{ action.desc }}</span>
          </div>
        </div>
      </div>
    </div>

    <div v-if="serverInfo" class="config-overview">
      <Card title="当前配置概览">
        <div class="config-info">
          <div class="info-item">
            <span class="info-label">应用数量</span>
            <span class="info-value">{{ Object.keys(serverInfo.apps || {}).length }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">HTTP服务器</span>
            <span class="info-value">{{ serverInfo.apps?.http?.servers ? Object.keys(serverInfo.apps.http.servers).length : 0 }}</span>
          </div>
        </div>
      </Card>
    </div>
  </div>
</template>

<style scoped>
.dashboard {
  padding: 24px;
}

.welcome-section {
  margin-bottom: 32px;
}

.welcome-title {
  font-family: 'Space Grotesk', sans-serif;
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text);
  margin: 0 0 8px 0;
}

.welcome-desc {
  font-size: 16px;
  color: #64748b;
  margin: 0;
}

.status-section {
  margin-bottom: 32px;
}

.status-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.status-info {
  display: flex;
  align-items: center;
  gap: 16px;
}

.status-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text);
}

.section-title {
  font-family: 'Space Grotesk', sans-serif;
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text);
  margin: 0 0 16px 0;
}

.quick-actions {
  margin-bottom: 32px;
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 16px;
}

.action-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.action-card:hover {
  border-color: var(--color-primary);
  box-shadow: 0 4px 12px rgba(30, 64, 175, 0.1);
}

.action-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-background);
  border-radius: 8px;
  color: var(--color-primary);
}

.action-icon svg {
  width: 20px;
  height: 20px;
}

.action-text {
  display: flex;
  flex-direction: column;
}

.action-label {
  font-size: 15px;
  font-weight: 600;
  color: var(--color-text);
}

.action-desc {
  font-size: 13px;
  color: #64748b;
}

.config-overview .config-info {
  display: flex;
  gap: 40px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-label {
  font-size: 13px;
  color: #64748b;
}

.info-value {
  font-size: 20px;
  font-weight: 600;
  color: var(--color-primary);
}
</style>