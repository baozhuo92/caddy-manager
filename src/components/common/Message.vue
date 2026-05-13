<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  type: {
    type: String,
    default: 'info'
  },
  message: String,
  duration: {
    type: Number,
    default: 3000
  }
})

const emit = defineEmits(['close'])
const visible = ref(true)

watch(() => props.message, () => {
  visible.value = true
  if (props.duration > 0) {
    setTimeout(() => {
      visible.value = false
      emit('close')
    }, props.duration)
  }
})
</script>

<template>
  <Transition name="fade">
    <div v-if="visible" class="message" :class="type">
      <span class="message-icon">
        {{ type === 'success' ? '✓' : type === 'error' ? '✕' : type === 'warning' ? '!' : 'i' }}
      </span>
      <span class="message-text">{{ message }}</span>
      <button class="message-close" @click="visible = false; $emit('close')">×</button>
    </div>
  </Transition>
</template>

<style scoped>
.message {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 6px;
  font-size: 14px;
  margin-bottom: 16px;
}

.message.success {
  background: #ecfdf5;
  border: 1px solid #a7f3d0;
  color: #065f46;
}

.message.error {
  background: #fef2f2;
  border: 1px solid #fecaca;
  color: #991b1b;
}

.message.warning {
  background: #fffbeb;
  border: 1px solid #fde68a;
  color: #92400e;
}

.message.info {
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  color: #1e40af;
}

.message-icon {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  font-size: 12px;
  font-weight: 600;
}

.message.success .message-icon { background: #10b981; color: white; }
.message.error .message-icon { background: #ef4444; color: white; }
.message.warning .message-icon { background: #f59e0b; color: white; }
.message.info .message-icon { background: #3b82f6; color: white; }

.message-text {
  flex: 1;
}

.message-close {
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  color: inherit;
  opacity: 0.6;
}

.message-close:hover {
  opacity: 1;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>