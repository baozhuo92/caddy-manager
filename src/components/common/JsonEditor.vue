<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  modelValue: [String, Object],
  readonly: {
    type: Boolean,
    default: false
  },
  height: {
    type: String,
    default: '300px'
  }
})

const emit = defineEmits(['update:modelValue'])

const content = ref('')

watch(() => props.modelValue, (val) => {
  content.value = typeof val === 'string' ? val : JSON.stringify(val, null, 2)
}, { immediate: true })

const handleInput = () => {
  try {
    const parsed = JSON.parse(content.value)
    emit('update:modelValue', parsed)
  } catch {
    emit('update:modelValue', content.value)
  }
}

const formatJson = () => {
  try {
    const parsed = JSON.parse(content.value)
    content.value = JSON.stringify(parsed, null, 2)
  } catch {}
}
</script>

<template>
  <div class="json-editor">
    <div class="editor-toolbar" v-if="!readonly">
      <button @click="formatJson" class="toolbar-btn">格式化</button>
    </div>
    <textarea
      v-model="content"
      class="editor-textarea"
      :style="{ height }"
      :readonly="readonly"
      @input="handleInput"
      spellcheck="false"
    ></textarea>
  </div>
</template>

<style scoped>
.json-editor {
  border: 1px solid var(--color-border);
  border-radius: 6px;
  overflow: hidden;
  background: var(--color-surface);
}

.editor-toolbar {
  padding: 8px 12px;
  border-bottom: 1px solid var(--color-border);
  background: var(--color-background);
}

.toolbar-btn {
  padding: 4px 12px;
  border: 1px solid var(--color-border);
  border-radius: 4px;
  background: var(--color-surface);
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.toolbar-btn:hover {
  background: var(--color-primary);
  color: white;
  border-color: var(--color-primary);
}

.editor-textarea {
  width: 100%;
  padding: 12px;
  border: none;
  font-family: 'Fira Code', monospace;
  font-size: 13px;
  line-height: 1.5;
  color: var(--color-text);
  background: var(--color-surface);
  resize: vertical;
}

.editor-textarea:focus {
  outline: none;
}

.editor-textarea[readonly] {
  background: var(--color-background);
}
</style>