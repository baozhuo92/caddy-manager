<script setup>
import { ref, watch } from 'vue'
import Card from '../common/Card.vue'
import Btn from '../common/Btn.vue'
import Input from '../common/Input.vue'

const props = defineProps({
  upstream: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['save', 'cancel'])

const form = ref({
  dial: '',
  max_requests: 0,
  max_fails: 0,
  fail_duration: ''
})

watch(() => props.upstream, (val) => {
  if (val) {
    form.value = {
      dial: val.dial || '',
      max_requests: val.max_requests || 0,
      max_fails: val.max_fails || 0,
      fail_duration: val.fail_duration || ''
    }
  } else {
    form.value = { dial: '', max_requests: 0, max_fails: 0, fail_duration: '' }
  }
}, { immediate: true })

const handleSave = () => {
  const data = { dial: form.value.dial }
  if (form.value.max_requests > 0) data.max_requests = form.value.max_requests
  if (form.value.max_fails > 0) data.max_fails = form.value.max_fails
  if (form.value.fail_duration) data.fail_duration = form.value.fail_duration
  emit('save', data)
}
</script>

<template>
  <Card :title="upstream ? '编辑上游' : '添加上游'" class="form-card">
    <div class="form-body">
      <Input
        v-model="form.dial"
        label="目标地址"
        placeholder="例如: 10.0.1.1:80 或 localhost:3000"
      />
      <div class="form-row">
        <Input
          v-model="form.max_requests"
          label="最大请求数"
          type="number"
          placeholder="0 = 不限"
        />
        <Input
          v-model="form.max_fails"
          label="最大失败次数"
          type="number"
          placeholder="0 = 不限"
        />
      </div>
      <Input
        v-model="form.fail_duration"
        label="失败持续时间"
        placeholder="例如: 30s, 5m"
      />
    </div>
    <div class="form-actions">
      <Btn variant="ghost" @click="$emit('cancel')">取消</Btn>
      <Btn variant="primary" :disabled="!form.dial.trim()" @click="handleSave">
        {{ upstream ? '保存修改' : '添加上游' }}
      </Btn>
    </div>
  </Card>
</template>

<style scoped>
.form-card {
  max-width: 520px;
}

.form-body {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-row {
  display: flex;
  gap: 12px;
}

.form-row > * {
  flex: 1;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 20px;
}
</style>