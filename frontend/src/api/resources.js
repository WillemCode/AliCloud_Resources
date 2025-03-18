import axios from 'axios'

const api = axios.create({
  baseURL: '/api', // 通过代理访问
  timeout: 5000
})

export default {
  async getResources(type, page = 1, pageSize = 10) {
    try {
      const response = await api.get(`/${type.toLowerCase()}`, {
        params: { page, pageSize }
      })
      return response.data
    } catch (error) {
      console.error('API Error:', error)
      throw error
    }
  },

  async searchResources(keyword, type = 'all', page = 1, pageSize = 10) {
    try {
      const response = await api.get('/search', { 
        params: { 
          q: keyword,
          type: type,
          page: page,
          pageSize: pageSize
        } 
      })
      return response.data
    } catch (error) {
      console.error('Search Error:', error)
      throw error
    }
  }
}
