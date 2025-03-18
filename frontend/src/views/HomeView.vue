<template>
  <div class="home">
    <Sidebar
      :current-type="store.currentType"
      @select-type="handleTypeChange"
    />
    <main class="content">
      <SearchBar />
      <ResourceList
        :resources="store.resources[store.currentType]"
        :is-loading="store.isLoading"
        :current-type="store.currentType"
      />
    </main>
  </div>
</template>

<script setup>
import { watch } from 'vue'
import { useResourcesStore } from '@/store/useResourcesStore.js'
import Sidebar from '@/components/Sidebar.vue'
import SearchBar from '@/components/SearchBar.vue'
import ResourceList from '@/components/ResourceList.vue'

const store = useResourcesStore()

const handleTypeChange = (newType) => {
  store.currentType = newType
  store.pagination = {  // 切换类型时重置分页
    currentPage: 1,
    pageSize: 10,
    totalItems: 0
  }
  store.fetchResources()
}

watch(
  () => store.currentType,
  (newVal, oldVal) => {
    if (newVal !== oldVal && !store.searchKeyword) {
      store.fetchResources()
    }
  },
  { immediate: true }
)

// 初始化加载数据
store.fetchResources()

</script>

<style scoped>
.home {
  display: flex;
  min-height: 100vh;
}

.content {
  flex: 1;
  padding: 2rem;
  background: #f5f6fa;
}
</style>
