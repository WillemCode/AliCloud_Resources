import { defineStore } from 'pinia'
import api from '@/api/resources.js'

export const useResourcesStore = defineStore('resources', {
  state: () => ({
    // 新增分页状态
    pagination: {
      currentPage: 1,
      pageSize: 10,
      totalItems: 0
    },
    currentType: 'ecs',
    searchKeyword: '',
    resources: {
      ecs: [],
      rds: [],
      slb: [],
      polardb: []
    },
    isLoading: false
  }),
  actions: {
    async fetchResources() {
      this.isLoading = true
      try {
        const { data, total, page, pageSize } = await api.getResources(
          this.currentType,
          this.pagination.currentPage,
          this.pagination.pageSize
        )
        this.resources[this.currentType] = data
        this.pagination = { 
          currentPage: page,
          pageSize: pageSize,
          totalItems: total
        }
      } catch (error) {
        console.error('获取资源失败:', error)
      } finally {
        this.isLoading = false
      }
    },
    // 新增分页操作方法
    setPage(page) {
      this.pagination.currentPage = page
      this.fetchResources()
    },
    setPageSize(size) {
      this.pagination.pageSize = size
      this.pagination.currentPage = 1
      this.fetchResources()
    },
    
    async initialize() {
      await this.fetchResources()
    },
    async searchResources() {
      if (!this.searchKeyword) {
        await this.fetchResources()
        return
      }
      
      this.isLoading = true
      try {
        const { data, total, page, pageSize } = await api.searchResources(this.searchKeyword, this.currentType)
        this.resources[this.currentType] = data
        this.pagination = { 
          currentPage: page,
          pageSize: pageSize,
          totalItems: total
        }
      } catch (error) {
        console.error('搜索失败:', error)
      } finally {
        this.isLoading = false
      }
    }
  }
})
