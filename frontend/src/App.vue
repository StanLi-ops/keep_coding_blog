<script setup>
import Header from './components/Header.vue'
</script>

<template>
  <div class="app">
    <Header />
    <router-view></router-view>
  </div>
</template>

<style>
.app {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}
</style>
