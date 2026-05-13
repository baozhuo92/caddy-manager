<script setup>
defineProps({
  variant: {
    type: String,
    default: 'default'
  },
  size: {
    type: String,
    default: 'md'
  }
})
</script>

<template>
  <span class="badge" :class="[`badge-${variant}`, `badge-${size}`]">
    <slot />
  </span>
</template>

<style scoped>
.badge {
  display: inline-flex;
  align-items: center;
  border-radius: 4px;
  font-weight: 500;
}

.badge-sm {
  padding: 2px 6px;
  font-size: 11px;
}

.badge-md {
  padding: 4px 10px;
  font-size: 12px;
}

.badge-lg {
  padding: 6px 12px;
  font-size: 14px;
}

.badge-default {
  background: #e2e8f0;
  color: #475569;
}

.badge-primary {
  background: #dbeafe;
  color: #1e40af;
}

.badge-success {
  background: #d1fae5;
  color: #065f46;
}

.badge-warning {
  background: #fef3c7;
  color: #92400e;
}

.badge-error {
  background: #fee2e2;
  color: #991b1b;
}

.badge-info {
  background: #e0f2fe;
  color: #075985;
}
</style>