<template>
  <div class="search-bar">
    <input
      v-model="keyword"
      type="text"
      placeholder="Search resources..."
      @input="handleSearch"
    />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useResourcesStore } from '@/store/useResourcesStore'

const store = useResourcesStore()
const keyword = ref('')

const handleSearch = () => {
  store.searchKeyword = keyword.value
  store.searchResources()
}
</script>

<style scoped>
.search-bar {
  padding: 1rem;
  margin-bottom: 1rem;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

input {
  width: 100%;
  padding: 0.8rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  outline: none;
  font-size: 1rem;
}

input:focus {
  border-color: #3498db;
}
</style>
