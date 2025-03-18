<template>
  <div class="pagination">
    <button 
      :disabled="currentPage === 1"
      @click="handlePageChange(currentPage - 1)"
    >
      上一页
    </button>
    
    <span class="page-info">
      第 {{ currentPage }} 页 / 共 {{ totalPages }} 页
    </span>
    
    <button
      :disabled="currentPage >= totalPages"
      @click="handlePageChange(currentPage + 1)"
    >
      下一页
    </button>
    
    <select 
      v-model="pageSize" 
      @change="handlePageSizeChange"
      class="page-size-select"
    >
      <option value="5">每页5条</option>
      <option value="10">每页10条</option>
      <option value="20">每页20条</option>
      <option value="50">每页50条</option>
    </select>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useResourcesStore } from '@/store/useResourcesStore'

const store = useResourcesStore()
const pageSize = ref(store.pagination.pageSize)

const currentPage = computed(() => store.pagination.currentPage)
const totalPages = computed(() => 
  Math.ceil(store.pagination.totalItems / store.pagination.pageSize)
)

const handlePageChange = (page) => {
  store.setPage(page)
}

const handlePageSizeChange = () => {
  store.setPageSize(Number(pageSize.value))
}
</script>

<style scoped>
.pagination {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-top: 1rem;
  padding: 1rem;
  background: white;
  border-radius: 8px;
}

button {
  padding: 0.5rem 1rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: #f8f9fa;
  cursor: pointer;
  transition: all 0.2s;
}

button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

button:hover:not(:disabled) {
  background: #3498db;
  color: white;
  border-color: #3498db;
}

.page-size-select {
  padding: 0.5rem;
  border-radius: 4px;
  margin-left: auto;
}
</style>
