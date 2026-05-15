<script setup>
import { watch, ref, nextTick } from 'vue'

const props = defineProps({
  visible: Boolean,
  title: String,
  width: { type: String, default: '560px' }
})

const emit = defineEmits(['close'])
const overlayRef = ref(null)

watch(() => props.visible, async (v) => {
  document.body.style.overflow = v ? 'hidden' : ''
  if (v) {
    await nextTick()
    overlayRef.value?.focus()
  }
})

function onOverlayClick(e) {
  if (e.target === e.currentTarget) emit('close')
}
</script>

<template>
  <Teleport to="body">
    <Transition name="drawer">
      <div v-if="visible" ref="overlayRef" class="drawer-overlay" tabindex="-1" @click="onOverlayClick" @keydown.esc="$emit('close')">
        <div class="drawer-panel" :style="{ width }">
          <div class="drawer-header">
            <h3 class="drawer-title">{{ title }}</h3>
            <button class="drawer-close" @click="$emit('close')" title="关闭">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="20" height="20">
                <path d="M18 6 6 18M6 6l12 12"/>
              </svg>
            </button>
          </div>
          <div class="drawer-body">
            <slot />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.drawer-overlay {
  position: fixed;
  inset: 0;
  z-index: 1000;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  justify-content: flex-end;
  outline: none;
}

.drawer-panel {
  height: 100%;
  background: var(--color-surface);
  box-shadow: -4px 0 24px rgba(0, 0, 0, 0.12);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.drawer-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 18px 24px;
  border-bottom: 1px solid var(--color-border);
  flex-shrink: 0;
}

.drawer-title {
  font-family: 'Space Grotesk', sans-serif;
  font-size: 17px;
  font-weight: 600;
  color: var(--color-text);
  margin: 0;
}

.drawer-close {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  border-radius: 6px;
  background: transparent;
  color: #64748b;
  cursor: pointer;
  transition: all 0.2s;
}

.drawer-close:hover {
  background: var(--color-border);
  color: var(--color-text);
}

.drawer-body {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.drawer-enter-active {
  animation: drawerFadeIn 0.2s ease;
}
.drawer-leave-active {
  animation: drawerFadeIn 0.2s ease reverse;
}
.drawer-enter-active .drawer-panel {
  animation: drawerSlideIn 0.25s ease;
}
.drawer-leave-active .drawer-panel {
  animation: drawerSlideIn 0.25s ease reverse;
}

@keyframes drawerFadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes drawerSlideIn {
  from { transform: translateX(100%); }
  to { transform: translateX(0); }
}
</style>
