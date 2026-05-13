<script setup>
import { ref, onMounted } from 'vue'
import { useServerStore } from '../stores/server'
import Card from '../components/common/Card.vue'
import Btn from '../components/common/Btn.vue'
import JsonEditor from '../components/common/JsonEditor.vue'
import Message from '../components/common/Message.vue'
import Badge from '../components/common/Badge.vue'

const store = useServerStore()

const configInput = ref('')
const showEditor = ref(false)
const messages = ref([])
const adaptResult = ref(null)

onMounted(() => {
  store.checkStatus()
})

const showMessage = (type, msg) => {
  messages.value.push({ type, message: msg })
  setTimeout(() => messages.value.shift(), 3000)
}

const handleStop = async () => {
  if (!confirm('确定要停止Caddy服务器吗？')) return
  const success = await store.stopServer()
  if (success) {
    showMessage('success', '服务器已停止')
  } else {
    showMessage('error', store.error || '停止失败')
  }
}

const handleLoadConfig = async () => {
  if (!configInput.value.trim()) {
    showMessage('error', '请输入配置内容')
    return
  }
  const success = await store.loadConfig(configInput.value)
  if (success) {
    showMessage('success', '配置加载成功')
    configInput.value = ''
  } else {
    showMessage('error', store.error || '加载失败')
  }
}

const handleAdapt = async () => {
  if (!configInput.value.trim()) {
    showMessage('error', '请输入Caddyfile内容')
    return
  }
  const result = await store.adaptConfig(configInput.value, 'text/caddyfile')
  if (result.success) {
    adaptResult.value = result.data
    showMessage('success', '适配成功')
  } else {
    showMessage('error', result.error || '适配失败')
    adaptResult.value = null
  }
}
</script>

<template>
  <div class="server-page">
    <div v-for="(msg, idx) in messages" :key="idx">
      <Message :type="msg.type" :message="msg.message" />
    </div>

    <Card title="服务器状态">
      <div class="status-row">
        <div class="status-item">
          <span class="status-label">状态</span>
          <Badge :variant="store.status === 'running' ? 'success' : 'error'" size="md">
            {{ store.status === 'running' ? '运行中' : '已停止' }}
          </Badge>
        </div>
        <div class="status-item">
          <span class="status-label">操作</span>
          <Btn variant="danger" :loading="store.loading" @click="handleStop">
            停止服务器
          </Btn>
        </div>
      </div>
    </Card>

    <Card title="加载配置">
      <template #actions>
        <Btn size="sm" variant="ghost" @click="showEditor = !showEditor">
          {{ showEditor ? '隐藏编辑器' : '显示编辑器' }}
        </Btn>
      </template>

      <div class="config-section">
        <div class="config-type-selector">
          <label class="radio-label">
            <input type="radio" value="json" checked /> JSON配置
          </label>
          <label class="radio-label">
            <input type="radio" value="caddyfile" /> Caddyfile
          </label>
        </div>

        <div v-if="showEditor" class="editor-section">
          <textarea
            v-model="configInput"
            class="config-textarea"
            placeholder="输入JSON配置或Caddyfile内容..."
          ></textarea>
          <div class="editor-actions">
            <Btn variant="primary" @click="handleLoadConfig">加载配置</Btn>
            <Btn variant="cta" @click="handleAdapt">适配为JSON</Btn>
          </div>
        </div>

        <div v-if="adaptResult" class="adapt-result">
          <h4>适配结果：</h4>
          <JsonEditor :model-value="adaptResult" :readonly="true" :height="'300px'" />
        </div>
      </div>
    </Card>
  </div>
</template>

<style scoped>
.server-page {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.status-row {
  display: flex;
  gap: 40px;
  align-items: center;
}

.status-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.status-label {
  font-weight: 500;
  color: var(--color-text);
}

.config-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.config-type-selector {
  display: flex;
  gap: 20px;
}

.radio-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.config-textarea {
  width: 100%;
  min-height: 200px;
  padding: 12px;
  border: 1px solid var(--color-border);
  border-radius: 6px;
  font-family: 'Fira Code', monospace;
  font-size: 13px;
  line-height: 1.5;
  resize: vertical;
}

.config-textarea:focus {
  outline: none;
  border-color: var(--color-primary);
}

.editor-actions {
  display: flex;
  gap: 12px;
  margin-top: 12px;
}

.adapt-result {
  margin-top: 20px;
}

.adapt-result h4 {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 12px;
}
</style>