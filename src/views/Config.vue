<script setup>
import { ref, onMounted, computed } from 'vue'
import { useConfigStore } from '../stores/config'
import Card from '../components/common/Card.vue'
import Btn from '../components/common/Btn.vue'

const store = useConfigStore()
const pathInput = ref('/')

const pathSegments = computed(() => {
  if (!store.currentPath || store.currentPath === '/') return ['/']
  return store.currentPath.split('/').filter(Boolean)
})

const navigateTo = (path) => {
  pathInput.value = path || '/'
  store.fetchConfig(path || '/')
}

onMounted(() => {
  store.fetchConfig('/')
})
</script>

<template>
  <div class="config-page">
    <Card title="配置查看">
      <template #actions>
        <Btn size="sm" @click="store.fetchConfig(pathInput.value)">刷新</Btn>
      </template>

      <div v-if="store.loading" class="state-msg">加载中...</div>
      <div v-else-if="store.config !== null" class="json-display">
        <pre>{{ JSON.stringify(store.config, null, 2) }}</pre>
      </div>
      <div v-else class="state-msg">暂无配置数据</div>
    </Card>
  </div>
</template>

<style scoped>
.config-page {
  padding: 24px;
}

.path-navigator {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  flex-wrap: wrap;
  gap: 12px;
}

.breadcrumb {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
}

.breadcrumb-item a {
  color: var(--color-primary);
  text-decoration: none;
  font-family: 'Fira Code', monospace;
  font-size: 13px;
}

.breadcrumb-item a:hover {
  text-decoration: underline;
}

.sep {
  margin: 0 6px;
  color: #94a3b8;
}

.path-input-group {
  display: flex;
  gap: 8px;
}

.path-input {
  padding: 8px 12px;
  border: 1px solid var(--color-border);
  border-radius: 6px;
  font-family: 'Fira Code', monospace;
  font-size: 13px;
  width: 280px;
}

.path-input:focus {
  outline: none;
  border-color: var(--color-primary);
}

.json-display {
  background: var(--color-background);
  border-radius: 6px;
  padding: 16px;
  overflow-x: auto;
}

.json-display pre {
  font-family: 'Fira Code', monospace;
  font-size: 13px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
}

.state-msg {
  text-align: center;
  padding: 40px;
  color: #94a3b8;
}
</style>