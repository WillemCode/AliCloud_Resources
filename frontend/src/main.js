import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'

const app = createApp(App)
const pinia = createPinia()
app.use(pinia)
app.use(router)

// 初始化store
const initializeStore = async () => {
  const store = useResourcesStore()
  await store.initialize()
}
initializeStore()

app.mount('#app')
