import {createApp} from 'vue'
import App from './App.vue'
import pinia from './stores/index.ts'
import router from './router/router'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import "./assets/index.scss";
import "uno.css";
import "element-plus/theme-chalk/src/message.scss";
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import {i18n} from './i18n/index.ts'
const app = createApp(App)
app.use(router)
app.use(ElementPlus)
app.use(pinia)
app.use(i18n)

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}
app.mount("#app");
