<script setup>
import { ref, onMounted } from 'vue'
import { useSitesStore } from '../stores/sites'
import { caddyApi } from '../services/caddyApi'
import Card from '../components/common/Card.vue'
import Btn from '../components/common/Btn.vue'
import Badge from '../components/common/Badge.vue'
import Message from '../components/common/Message.vue'
import Table from '../components/common/Table.vue'
import SiteForm from '../components/proxy/SiteForm.vue'

const sitesStore = useSitesStore()
const messages = ref([])
const showForm = ref(false)
const editingSite = ref(null)
const upstreams = ref([])
const upstreamsLoading = ref(false)

const showMessage = (type, msg) => {
  messages.value.push({ type, message: msg, id: Date.now() })
  setTimeout(() => messages.value.shift(), 3000)
}

const fetchUpstreams = async () => {
  upstreamsLoading.value = true
  try {
    const res = await caddyApi.getReverseProxyUpstreams()
    upstreams.value = res.data || []
  } catch {
    upstreams.value = []
  } finally {
    upstreamsLoading.value = false
  }
}

const handleAdd = () => {
  editingSite.value = null
  showForm.value = true
}

const handleEdit = (site) => {
  editingSite.value = site
  showForm.value = true
}

const handleDelete = async (site) => {
  if (!confirm(`确定删除站点 ${site.name} 吗？`)) return
  const ok = await sitesStore.deleteSite(site)
  if (ok) {
    showMessage('success', '站点已删除')
  } else {
    showMessage('error', sitesStore.error || '删除失败')
  }
}

const handleSave = async (data) => {
  let ok
  if (editingSite.value) {
    data.serverName = editingSite.value.serverName
    ok = await sitesStore.updateSite(data)
  } else {
    ok = await sitesStore.addSite(data)
  }
  if (ok) {
    showMessage('success', editingSite.value ? '站点已更新' : '站点已创建')
    showForm.value = false
    editingSite.value = null
  } else {
    showMessage('error', sitesStore.error || '保存失败')
  }
}

const handleCancel = () => {
  showForm.value = false
  editingSite.value = null
}

onMounted(async () => {
  await sitesStore.loadSites()
  fetchUpstreams()
})
</script>

<template>
  <div class="proxy-page">
    <div v-for="(msg, idx) in messages" :key="msg.id">
      <Message :type="msg.type" :message="msg.message" @close="messages.splice(idx, 1)" />
    </div>

    <SiteForm
      v-if="showForm"
      :site="editingSite"
      @save="handleSave"
      @cancel="handleCancel"
    />

    <Card title="反向代理站点">
      <template #actions>
        <Btn size="sm" @click="sitesStore.loadSites()">刷新</Btn>
        <Btn size="sm" variant="primary" @click="handleAdd">新建站点</Btn>
      </template>

      <div v-if="sitesStore.loading" class="state-msg">加载中...</div>

      <div v-else-if="sitesStore.sites.length === 0" class="state-msg">
        <p>暂无反向代理站点。点击"新建站点"添加。</p>
      </div>

      <div v-else class="site-list">
        <div v-for="site in sitesStore.sites" :key="site.id" class="site-card">
          <div class="site-header">
            <div class="site-domain">
              <Badge variant="primary" size="lg">{{ site.domain || '(无域名)' }}</Badge>
              <Badge v-if="site.scheme !== 'all'" variant="default" size="sm">
                {{ site.scheme }}
              </Badge>
            </div>
            <div class="site-actions">
              <button class="icon-btn" title="编辑" @click="handleEdit(site)">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                  <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                </svg>
              </button>
              <button class="icon-btn danger" title="删除" @click="handleDelete(site)">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="3 6 5 6 21 6"/>
                  <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                </svg>
              </button>
            </div>
          </div>
          <div class="site-body">
            <div class="site-info">
              <div class="info-chip">
                <span class="chip-label">服务器</span>
                <span class="chip-value"><code>{{ site.serverName }}</code></span>
              </div>
              <div class="info-chip">
                <span class="chip-label">监听</span>
                <span class="chip-value">
                  <code v-if="site.listen?.length">{{ site.listen.join(', ') }}</code>
                  <span v-else class="none">默认 :443</span>
                </span>
              </div>
            </div>
            <div class="site-routes">
              <div class="routes-title">路由规则</div>
              <div v-for="(r, i) in site.routeEntries" :key="i" class="route-entry">
                <div class="route-match">
                  <Badge v-if="r.path" variant="default" size="sm">{{ r.path }}</Badge>
                  <Badge v-else variant="default" size="sm">/ (全部)</Badge>
                  <Badge v-if="r.websocket" variant="info" size="sm">WS</Badge>
                </div>
                <div class="route-upstreams">
                  <code v-for="(u, j) in r.upstreams" :key="j" class="upstream-tag">{{ u }}</code>
                </div>
              </div>
            </div>
            <div class="site-features">
              <Badge v-if="site.cors?.enabled" variant="success" size="sm">CORS</Badge>
              <Badge v-if="site.basicauth?.enabled" variant="warning" size="sm">BasicAuth</Badge>
              <Badge v-if="site.customHeaders?.request?.length > 0 || site.customHeaders?.response?.length > 0" variant="default" size="sm">自定义头</Badge>
            </div>
          </div>
        </div>
      </div>
    </Card>

    <Card title="上游实时状态" :padding="false">
      <template #actions>
        <Btn size="sm" @click="fetchUpstreams">刷新</Btn>
      </template>
      <div v-if="upstreamsLoading" class="state-msg">加载中...</div>
      <Table v-else-if="upstreams.length > 0" :columns="[
        { key: 'address', label: '地址' },
        { key: 'num_requests', label: '活跃请求' },
        { key: 'fails', label: '失败次数' },
        { key: 'status', label: '状态' }
      ]" :data="upstreams">
        <template #address="{ value }">
          <code class="addr">{{ value }}</code>
        </template>
        <template #fails="{ value }">
          <Badge :variant="value > 3 ? 'error' : value > 0 ? 'warning' : 'default'" size="md">{{ value }}</Badge>
        </template>
        <template #status="{ row }">
          <Badge :variant="row.fails > 3 ? 'error' : row.num_requests > 0 ? 'success' : 'default'" size="md">
            {{ row.fails > 3 ? '异常' : row.num_requests > 0 ? '活跃' : '空闲' }}
          </Badge>
        </template>
      </Table>
      <div v-else class="state-msg">暂无上游数据</div>
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

.state-msg {
  text-align: center;
  padding: 40px;
  color: #94a3b8;
}

.site-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.site-card {
  border: 1px solid var(--color-border);
  border-radius: 8px;
  overflow: hidden;
  transition: border-color 0.2s;
}

.site-card:hover {
  border-color: var(--color-secondary);
}

.site-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 16px;
  background: var(--color-background);
  border-bottom: 1px solid var(--color-border);
}

.site-domain {
  display: flex;
  align-items: center;
  gap: 8px;
}

.site-actions {
  display: flex;
  gap: 4px;
}

.icon-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid transparent;
  border-radius: 6px;
  background: none;
  cursor: pointer;
  color: #64748b;
  transition: all 0.2s;
}

.icon-btn svg {
  width: 16px;
  height: 16px;
}

.icon-btn:hover {
  background: var(--color-surface);
  border-color: var(--color-border);
  color: var(--color-primary);
}

.icon-btn.danger:hover {
  color: var(--color-error);
  border-color: var(--color-error);
}

.site-body {
  padding: 14px 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.site-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.info-chip {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

.chip-label {
  color: #64748b;
  min-width: 48px;
}

.chip-value {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.upstream-tag {
  padding: 2px 8px;
  background: var(--color-background);
  border: 1px solid var(--color-border);
  border-radius: 4px;
  font-family: 'Fira Code', monospace;
  font-size: 12px;
}

.site-routes {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.routes-title {
  font-size: 12px;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.route-entry {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 10px;
  background: var(--color-background);
  border: 1px solid var(--color-border);
  border-radius: 6px;
}

.route-match {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
}

.route-upstreams {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.site-features {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.none {
  color: #94a3b8;
}

.addr {
  font-family: 'Fira Code', monospace;
  color: var(--color-primary);
}
</style>
