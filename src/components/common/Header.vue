<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useServerStore } from '../../stores/server'
import { caddyApi } from '../../services/caddyApi'

const route = useRoute()
const serverStore = useServerStore()

const serverUrl = ref('/api/')
const isConnected = ref(false)

const pageTitle = computed(() => {
  const titles = {
    '/': '仪表盘',
    '/config': '配置管理',
    '/servers': '服务器',
    '/certificates': '证书',
    '/proxy': '反向代理'
  }
  return titles[route.path] || 'Caddy Manager'
})

const checkConnection = async () => {
  try {
    await caddyApi.getConfig('/')
    isConnected.value = true
    serverStore.status = 'running'
  } catch {
    isConnected.value = false
    serverStore.status = 'stopped'
  }
}

onMounted(() => {
  checkConnection()
  setInterval(checkConnection, 10000)
})

watch(serverUrl, (newUrl) => {
  caddyApi.defaults.baseURL = newUrl
})
</script>

<template>
  <header class="header">
    <div class="header-left">
      <span class="page-title">{{ pageTitle }}</span>
    </div>
    <div class="header-right">
      <div class="status-indicator">
        <span class="status-dot" :class="isConnected ? 'connected' : 'disconnected'"></span>
        <span class="status-text">{{ isConnected ? '已连接' : '未连接' }}</span>
      </div>
      <div class="server-url">
        <input type="text" v-model="serverUrl" class="url-input" placeholder="Caddy API 地址" />
      </div>
    </div>
  </header>
</template>

<style scoped>
.header {
  height: 60px;
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
}

.header-left {
  display: flex;
  align-items: center;
}

.page-title {
  font-family: 'Space Grotesk', sans-serif;
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 24px;
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-dot.connected {
  background: var(--color-success);
}

.status-dot.disconnected {
  background: var(--color-error);
}

.status-text {
  font-size: 14px;
  color: var(--color-text);
}

.server-url {
  display: flex;
  align-items: center;
}

.url-input {
  padding: 8px 12px;
  border: 1px solid var(--color-border);
  border-radius: 6px;
  font-family: 'Fira Code', monospace;
  font-size: 13px;
  color: var(--color-text);
  background: var(--color-background);
  width: 200px;
}

.url-input:focus {
  outline: none;
  border-color: var(--color-primary);
}
</style>