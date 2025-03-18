<template>
  <nav class="sidebar">
    <h2 class="logo">AliCloud<br>多账户资产管理</h2>
    <ul>
      <li
        v-for="type in resourceTypes"
        :key="type.value"
        :class="{ active: currentType === type.value }"
        @click="selectType(type.value)"
      >
        {{ type.label }}
      </li>
    </ul>
  </nav>
</template>

<script setup>
import { ref } from 'vue'

const resourceTypes = [
  { value: 'ecs', label: 'ECS' },
  { value: 'rds', label: 'RDS' },
  { value: 'slb', label: 'SLB' },
  { value: 'polardb', label: 'PolarDB' }
]

const props = defineProps({
  currentType: String
})

const emit = defineEmits(['select-type'])

const selectType = (type) => {
  emit('select-type', type)
}
</script>

<style scoped>
.sidebar {
  width: 240px;
  background: #2c3e50;
  color: white;
  height: 100vh;
  padding: 1rem;
}

.logo {
  margin: 1rem 0 2rem;
  text-align: center;
}

ul {
  list-style: none;
  padding: 0;
}

li {
  padding: 1rem;
  cursor: pointer;
  border-radius: 4px;
  transition: background 0.3s;
}

li:hover {
  background: #34495e;
}

li.active {
  background: #3498db;
}
</style>
