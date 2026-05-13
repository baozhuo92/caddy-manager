<script setup>
import { ref, onMounted, computed } from 'vue'
import { useConfigStore } from '../stores/config'
import Card from '../components/common/Card.vue'
import Btn from '../components/common/Btn.vue'
import JsonEditor from '../components/common/JsonEditor.vue'
import Message from '../components/common/Message.vue'
import Badge from '../components/common/Badge.vue'

const store = useConfigStore()

const pathInput = ref('/')
const editValue = ref('')
const showEditor = ref(false)
const messages = ref([])

const pathSegments = computed(() => {
  if (!store.currentPath || store.currentPath === '/') return ['/']
  return store.currentPath.split('/').filter(Boolean)
})

const navigateTo = (path) => {
  pathInput.value = path || '/'
  store.fetchConfig(path || '/')
}

const refreshConfig = () => {
  store.fetchConfig(pathInput.value)
}

onMounted(() => {
  store.fetchConfig('/')
})

const showMessage = (type, msg) => {
  messages.value.push({ type, message: msg })
  setTimeout(() => {
    messages.value.shift()
  }, 3000)
}

const handleUpdate = async () => {
  try {
    const data = JSON.parse(editValue.value)
    const success = await store.updateConfig(pathInput.value, data)
    if (success) {
      showMessage('success', '配置更新成功')
      showEditor.value = false
    } else {
      showMessage('error', store.error || '更新失败')
    }
  } catch (e) {
    showMessage('error', 'JSON格式错误')
  }
}

const handleDelete = async () => {
  if (!confirm('确定要删除当前配置吗？')) return
  const success = await store.deleteConfig(pathInput.value)
  if (success) {
    showMessage('success', '配置删除成功')
    pathInput.value = '/'
    store.fetchConfig('/')
  } else {
    showMessage('error', store.error || '删除失败')
  }
}
</script>

<template>
  <div class="config-page">
    <div v-for="(msg, idx) in messages" :key="idx">
      <Message :type="msg.type" :message="msg.message" />
    </div>

    <Card title="配置管理">
      <template #actions>
        <Btn size="sm" @click="refreshConfig">刷新</Btn>
      </template>

      <div class="path-navigator">
        <div class="breadcrumb">
          <span
            v-for="(segment, idx) in pathSegments"
            :key="idx"
            class="breadcrumb-item"
          >
            <span class="breadcrumb-sep">/</span>
            <a href="#" @click.prevent="navigateTo(pathSegments.slice(0, idx + 1).join('/'))">
              {{ segment || '根' }}
            </a>
          </span>
        </div>
        <div class="path-input-group">
          <input
            v-model="pathInput"
            type="text"
            class="path-input"
            placeholder="输入路径，如 /apps/http/servers"
          />
          <Btn size="sm" @click="navigateTo(pathInput)">跳转</Btn>
        </div>
      </div>

      <div v-if="store.loading" class="loading">加载中...</div>

      <div v-else-if="store.config !== null" class="config-content">
        <div class="config-actions">
          <Btn variant="primary" size="sm" @click="showEditor = true">编辑</Btn>
          <Btn variant="danger" size="sm" @click="handleDelete">删除</Btn>
        </div>

        <JsonEditor
          v-if="showEditor"
          v-model="editValue"
          :height="'400px'"
          class="editor-wrapper"
        />
        <div v-else class="json-display">
          <pre>{{ JSON.stringify(store.config, null, 2) }}</pre>
        </div>

        <div v-if="showEditor" class="editor-actions">
          <Btn @click="showEditor = false">取消</Btn>
          <Btn variant="primary" @click="handleUpdate">保存</Btn>
        </div>
      </div>

      <div v-else class="empty">
        <p>暂无配置数据</p>
      </div>
    </Card>

    <Card title="快捷路径" class="quick-paths">
      <div class="path-buttons">
        <Btn
          v-for="path in ['/', '/apps', '/apps/http', '/apps/http/servers']"
          :key="path"
          variant="ghost"
          size="sm"
          @click="navigateTo(path)"
        >
          {{ path }}
        </Btn>
      </div>
    </Card>
  </div>
</template>

<style scoped>
.config-page {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
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

.breadcrumb-sep {
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

.config-actions {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
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

.editor-wrapper {
  margin-bottom: 16px;
}

.editor-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 16px;
}

.quick-paths .path-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.loading,
.empty {
  text-align: center;
  padding: 40px;
  color: #94a3b8;
}
</style>