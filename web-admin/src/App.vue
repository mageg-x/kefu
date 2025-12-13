<template>
  <div id="app" class="container">
    
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const data = ref('')

onMounted(() => {
  console.log('Component mounted.')
})
</script>

<style scoped>
.container {
  padding: 20px;
}
</style>