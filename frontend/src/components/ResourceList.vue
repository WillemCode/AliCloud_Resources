<template>
  <div class="resource-list">
    <!-- 加载状态 -->
    <div v-if="isLoading" class="loading">加载中...</div>
    
    <!-- 数据表格 -->
    <table v-else class="resource-table">
      <thead>
        <tr>
          <th v-for="header in tableHeaders" :key="header.key">
            {{ header.label }}
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(item, index) in resources" :key="index">
          <td v-for="header in tableHeaders" :key="header.key">
            <template v-if="header.key === 'Status'">
              <span class="status-dot" :class="statusClass(item.Status)"></span>
              {{ item.Status }}
            </template>
            <template v-else-if="header.key === 'ConnectionString' && item.ConnectionString">
              <div class="connection-cell">
                <div 
                  v-for="(addr, i) in item.ConnectionString.split(',')" 
                  :key="i"
                  class="connection-item"
                >
                  {{ addr }}
                </div>
              </div>
            </template>
            <template v-else>
              {{ item[header.key] || '-' }}
            </template>
          </td>
        </tr>
      </tbody>
    </table>

    <!-- 空状态 -->
    <div v-if="!isLoading && resources.length === 0" class="empty">
      未找到相关资源
    </div>

    <!-- 分页组件 -->
    <Pagination v-if="resources.length > 0" />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import Pagination from './Pagination.vue'

const props = defineProps({
  resources: Array,
  isLoading: Boolean,
  currentType: String
})

// 动态表头配置
const tableHeaders = computed(() => {
  const baseHeaders = [
    { key: 'InstanceID', label: '实例ID' },
    { key: 'CloudName', label: '账户名称'}, 
    { key: 'Status', label: '状态' },
    { key: 'RegionID', label: '区域' }
  ]

  switch (props.currentType) {
    case 'ecs':
      return [
        ...baseHeaders,
        { key: 'InstanceType', label: '实例规格' },
        { key: 'CPU', label: 'CPU' },
        { key: 'Memory', label: '内存(MB)' },
        { key: 'PrivateIP', label: '内网IP'},
        { key: 'PublicIP', label: '公网IP' },
        { key: 'InstanceName', label: '实例名称' }, 
        { key: 'OSName', label: '操作系统' }
      ]
    case 'rds':
      return [
        ...baseHeaders,
        { key: 'Description', label: '数据库名称' },
        { key: 'Engine', label: '引擎类型' },
        { key: 'Memory', label: '内存(MB)' },
        { key: 'ConnectionString', label: '连接地址' }
      ]
    case 'slb':
      return [
        ...baseHeaders,
        { key: 'LoadBalancerName', label: '负载均衡名称' },
        { key: 'IPAddress', label: '服务地址' },
        { key: 'Bandwidth', label: '带宽' },
        { key: 'NetworkType', label: '网络类型' }
      ]
    case 'polardb':
      return [
        ...baseHeaders,
        { key: 'Description', label: '集群名称' },
        { key: 'Engine', label: '引擎类型' },
        { key: 'DBNodeCount', label: '节点数量' },
        { key: 'MemorySize', label: '内存' },
        { key: 'ConnectionString', label: '连接地址' }
      ]
    default:
      return baseHeaders
  }
})

// 状态颜色逻辑
const statusClass = (status) => {
  if (!status) return 'status-unknown'
  
  status = status.toLowerCase()
  if (status.includes('running') || status.includes('active')) {
    return 'status-running'
  } else if (status.includes('stopped') || status.includes('inactive')) {
    return 'status-stopped'
  } else if (status.includes('error') || status.includes('failed')) {
    return 'status-error'
  } else {
    return 'status-unknown'
  }
}
</script>

<style scoped>
.resource-table {
  width: 100%;
  border-collapse: collapse;
  background: white;
}

th, td {
  padding: 1rem;
  text-align: left;
  border-bottom: 1px solid #eee;
}

th {
  background: #f8f9fa;
  font-weight: 500;
}

.connection-cell {
  max-width: 500px;
}

.connection-item {
  padding: 0.25rem 0;
  /* white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis; */
  word-break: break-all;
}

.status-dot {
  display: inline-block;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  margin-right: 8px;
}

.status-running {
  background-color: #52c41a;
}

.status-stopped {
  background-color: #faad14;
}

.status-error {
  background-color: #f5222d;
}

.status-unknown {
  background-color: #d9d9d9;
}

/* 响应式处理 */
@media (max-width: 768px) {
  .resource-table {
    display: block;
    overflow-x: auto;
  }
  
  th, td {
    min-width: 120px;
  }
}
</style>
