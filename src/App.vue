<script setup>
import Sidebar from './components/common/Sidebar.vue'
import Header from './components/common/Header.vue'
</script>

<template>
  <div class="app-layout">
    <Sidebar />
    <div class="main-area">
      <Header />
      <main class="content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<style>
@import url('https://fonts.googleapis.com/css2?family=DM+Sans:wght@400;500;700&family=Space+Grotesk:wght@400;500;600;700&family=Fira+Code:wght@400;500;600;700&display=swap');

.app-layout {
  display: flex;
  min-height: 100vh;
}

.main-area {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.content {
  flex: 1;
  overflow-y: auto;
  background: var(--color-background);
}
</style>