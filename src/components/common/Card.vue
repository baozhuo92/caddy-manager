<script setup>
defineProps({
  title: String,
  padding: {
    type: Boolean,
    default: true
  }
})
</script>

<template>
  <div class="card" :class="{ 'no-padding': !padding }">
    <div v-if="title" class="card-header">
      <h3 class="card-title">{{ title }}</h3>
      <slot name="actions" />
    </div>
    <div class="card-body">
      <slot />
    </div>
  </div>
</template>

<style scoped>
.card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  overflow: hidden;
}

.card.no-padding .card-body {
  padding: 0;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid var(--color-border);
}

.card-title {
  font-family: 'Space Grotesk', sans-serif;
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text);
  margin: 0;
}

.card-body {
  padding: 20px;
}
</style>