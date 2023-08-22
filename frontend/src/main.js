// main.ts
import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs' //语言设为中文
import 'element-plus/dist/index.css'
import App from './App.vue'
import store from './store'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'




const app = createApp(App)
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}
app.use(ElementPlus, {
    locale: zhCn,
})
app.use(store)
app.mount('#app')