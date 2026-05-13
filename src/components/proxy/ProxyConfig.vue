<script setup>
import { ref, computed, watch } from 'vue'
import Card from '../common/Card.vue'
import Btn from '../common/Btn.vue'
import Input from '../common/Input.vue'
import Message from '../common/Message.vue'
import Badge from '../common/Badge.vue'
import JsonEditor from '../common/JsonEditor.vue'
import UpstreamForm from './UpstreamForm.vue'

const props = defineProps({
  config: Object
})

const emit = defineEmits(['save', 'refresh'])

const messages = ref([])
const selectedServer = ref('')
const selectedRouteIdx = ref(-1)
const editingUpstream = ref(null)
const editingUpstreamIdx = ref(-1)
const showAddForm = ref(false)
const selectedTab = ref('upstreams')
const fullConfig = ref('')
const showFullEditor = ref(false)

const servers = computed(() => {
  if (!props.config?.apps?.http?.servers) return {}
  return props.config.apps.http.servers
})

const serverNames = computed(() => Object.keys(servers.value))

const routes = computed(() => {
  if (!selectedServer.value) return []
  return servers.value[selectedServer.value]?.routes || []
})

const proxyHandlers = computed(() => {
  const idx = selectedRouteIdx.value
  if (idx < 0 || idx >= routes.value.length) return []
  const handlers = routes.value[idx]?.handle || []
  return handlers.filter(h => h.handler === 'reverse_proxy')
})

const currentUpstreams = computed(() => {
  if (proxyHandlers.value.length === 0) return []
  return proxyHandlers.value[0]?.upstreams || []
})

const showMessage = (type, msg) => {
  messages.value.push({ type, message: msg, id: Date.now() })
}

const dismissMessage = (idx) => {
  messages.value.splice(idx, 1)
}

const updateFullConfig = () => {
  fullConfig.value = JSON.stringify(props.config, null, 2)
}

watch(() => props.config, () => {
  updateFullConfig()
}, { deep: true, immediate: true })

watch(serverNames, (names) => {
  if (names.length > 0 && !selectedServer.value) {
    selectedServer.value = names[0]
  }
})

const getUpstreamsPath = () => {
  const server = selectedServer.value
  const idx = selectedRouteIdx.value
  if (!server || idx < 0) return ''
  return `/apps/http/servers/${server}/routes/${idx}/handle/0/upstreams`
}

const handleAddUpstream = async (data) => {
  const path = getUpstreamsPath()
  if (!path) { showMessage('error', '请先选择服务器和路由'); return }
  emit('save', { path, action: 'add', data })
  showAddForm.value = false
}

const handleEditUpstream = (upstream, index) => {
  editingUpstream.value = { ...upstream }
  editingUpstreamIdx.value = index
}

const handleUpdateUpstream = async (data) => {
  const path = getUpstreamsPath()
  if (!path) return
  emit('save', { path, action: 'update', data, index: editingUpstreamIdx.value })
  editingUpstream.value = null
  editingUpstreamIdx.value = -1
}

const handleDeleteUpstream = (index) => {
  const path = getUpstreamsPath()
  if (!path) return
  emit('save', { path, action: 'delete', index })
}

const hasProxyHandler = computed(() => proxyHandlers.value.length > 0)

const addProxyHandler = () => {
  const server = selectedServer.value
  const idx = selectedRouteIdx.value
  if (!server || idx < 0) return
  const routePath = `/apps/http/servers/${server}/routes/${idx}/handle`
  const proxyHandler = {
    handler: 'reverse_proxy',
    upstreams: [{ dial: 'example.com:80' }]
  }
  emit('save', { path: routePath, action: 'addHandler', data: proxyHandler })
}

const toggleFullEditor = () => {
  showFullEditor.value = !showFullEditor.value
  if (showFullEditor.value) updateFullConfig()
}

const saveFullConfig = () => {
  try {
    const parsed = JSON.parse(fullConfig.value)
    emit('save', { action: 'fullConfig', data: parsed })
    showFullEditor.value = false
  } catch {
    showMessage('error', 'JSON 格式错误')
  }
}
</script>

<template>
  <div class="proxy-config">
    <div v-for="(msg, idx) in messages" :key="msg.id">
      <Message :type="msg.type" :message="msg.message" @close="dismissMessage(idx)" />
    </div>

    <Card title="选择服务器和路由">
      <div class="selector-row">
        <div class="selector-group">
          <label class="selector-label">服务器</label>
          <select v-model="selectedServer" class="selector">
            <option v-for="name in serverNames" :key="name" :value="name">{{ name }}</option>
          </select>
        </div>
        <div class="selector-group">
          <label class="selector-label">路由</label>
          <select v-model="selectedRouteIdx" class="selector">
            <option v-for="(route, idx) in routes" :key="idx" :value="idx">
              路由 {{ idx }}
              <template v-if="route.match?.[0]?.host"> - {{ route.match[0].host.join(', ') }}</template>
            </option>
          </select>
        </div>
        <Btn variant="ghost" size="sm" @click="$emit('refresh')">刷新</Btn>
      </div>
    </Card>

    <Card v-if="!selectedServer" title="提示">
      <p>未找到 HTTP 服务器配置。请在 Caddy 配置中添加 HTTP 服务器。</p>
    </Card>

    <Card v-else-if="selectedRouteIdx < 0" title="提示">
      <p>请选择一个路由来管理反向代理。</p>
    </Card>

    <Card v-else-if="!hasProxyHandler" title="反向代理">
      <div class="empty-proxy">
        <p>当前路由没有 reverse_proxy 处理器</p>
        <Btn size="sm" @click="addProxyHandler">添加 reverse_proxy</Btn>
      </div>
    </Card>

    <template v-else>
      <div class="tabs">
        <button class="tab" :class="{ active: selectedTab === 'upstreams' }" @click="selectedTab = 'upstreams'">
          上游服务器 {{ currentUpstreams.length }}
        </button>
        <button class="tab" :class="{ active: selectedTab === 'config' }" @click="selectedTab = 'config'">
          处理器配置
        </button>
      </div>

      <Card v-if="selectedTab === 'upstreams'" :padding="false">
        <template #actions>
          <Btn size="sm" @click="showAddForm = !showAddForm">
            {{ showAddForm ? '取消' : '添加上游' }}
          </Btn>
        </template>

        <UpstreamForm
          v-if="showAddForm"
          @save="handleAddUpstream"
          @cancel="showAddForm = false"
        />

        <div v-if="editingUpstream" class="edit-section">
          <UpstreamForm
            :upstream="editingUpstream"
            @save="handleUpdateUpstream"
            @cancel="editingUpstream = null"
          />
        </div>

        <table v-if="currentUpstreams.length > 0" class="upstream-table">
          <thead>
            <tr>
              <th>#</th>
              <th>目标地址</th>
              <th>最大请求</th>
              <th>最大失败</th>
              <th>失败时长</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(up, idx) in currentUpstreams" :key="idx">
              <td>{{ idx }}</td>
              <td><code class="dial-text">{{ up.dial }}</code></td>
              <td>{{ up.max_requests || '-' }}</td>
              <td>{{ up.max_fails || '-' }}</td>
              <td>{{ up.fail_duration || '-' }}</td>
              <td>
                <div class="row-actions">
                  <button class="action-btn edit" @click="handleEditUpstream(up, idx)">编辑</button>
                  <button class="action-btn delete" @click="handleDeleteUpstream(idx)">删除</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>

        <div v-else-if="!showAddForm" class="empty-state">
          <p>暂无上游服务器</p>
        </div>
      </Card>

      <Card v-if="selectedTab === 'config'">
        <div class="handler-config">
          <div class="handler-info">
            <h4>Handler 配置</h4>
            <pre>{{ JSON.stringify(proxyHandlers[0], null, 2) }}</pre>
          </div>
          <Btn variant="ghost" size="sm" @click="toggleFullEditor">
            {{ showFullEditor ? '隐藏编辑器' : '编辑完整配置' }}
          </Btn>
        </div>

        <div v-if="showFullEditor" class="full-editor">
          <textarea
            v-model="fullConfig"
            class="full-editor-textarea"
            spellcheck="false"
          ></textarea>
          <div class="editor-actions">
            <Btn variant="ghost" @click="showFullEditor = false">取消</Btn>
            <Btn variant="primary" @click="saveFullConfig">保存全部配置</Btn>
          </div>
        </div>
      </Card>
    </template>
  </div>
</template>

<style scoped>
.proxy-config {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.selector-row {
  display: flex;
  gap: 16px;
  align-items: flex-end;
  flex-wrap: wrap;
}

.selector-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.selector-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text);
}

.selector {
  padding: 8px 12px;
  border: 1px solid var(--color-border);
  border-radius: 6px;
  font-size: 14px;
  min-width: 180px;
  background: var(--color-surface);
}

.selector:focus {
  outline: none;
  border-color: var(--color-primary);
}

.tabs {
  display: flex;
  gap: 0;
  border-bottom: 2px solid var(--color-border);
}

.tab {
  padding: 10px 20px;
  border: none;
  background: none;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  color: #64748b;
  border-bottom: 2px solid transparent;
  margin-bottom: -2px;
  transition: all 0.2s;
}

.tab.active {
  color: var(--color-primary);
  border-bottom-color: var(--color-primary);
}

.upstream-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
}

.upstream-table th {
  text-align: left;
  padding: 12px 16px;
  background: var(--color-background);
  font-weight: 600;
  border-bottom: 1px solid var(--color-border);
}

.upstream-table td {
  padding: 12px 16px;
  border-bottom: 1px solid var(--color-border);
}

.upstream-table tr:hover td {
  background: var(--color-background);
}

.dial-text {
  font-family: 'Fira Code', monospace;
  color: var(--color-primary);
}

.row-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  padding: 4px 10px;
  border: 1px solid var(--color-border);
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
  background: var(--color-surface);
  transition: all 0.2s;
}

.action-btn.edit:hover {
  border-color: var(--color-primary);
  color: var(--color-primary);
}

.action-btn.delete:hover {
  border-color: var(--color-error);
  color: var(--color-error);
}

.edit-section {
  padding: 20px;
  border-bottom: 1px solid var(--color-border);
}

.empty-proxy {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.empty-state {
  text-align: center;
  padding: 40px;
  color: #94a3b8;
}

.handler-config {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.handler-config h4 {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 12px;
}

.handler-config pre {
  background: var(--color-background);
  padding: 12px;
  border-radius: 6px;
  font-family: 'Fira Code', monospace;
  font-size: 12px;
  white-space: pre-wrap;
  max-height: 200px;
  overflow-y: auto;
}

.full-editor {
  margin-top: 16px;
}

.full-editor-textarea {
  width: 100%;
  min-height: 300px;
  padding: 12px;
  border: 1px solid var(--color-border);
  border-radius: 6px;
  font-family: 'Fira Code', monospace;
  font-size: 13px;
  line-height: 1.5;
  resize: vertical;
}

.full-editor-textarea:focus {
  outline: none;
  border-color: var(--color-primary);
}

.editor-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 12px;
}
</style>