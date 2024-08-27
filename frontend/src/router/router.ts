import { createWebHistory, createRouter } from "vue-router";

const routes = [
  {
    path: '/', // 首页路径为空
    name: 'home',
    component: () => import('../views/home.vue'), // 首页组件
  },
  {
    path: '/chat',
    name: 'chat',
    component: () => import('../views/chat.vue'),
  },
  {
    path: '/sd',
    name: 'sd',
    component: () => import('../views/sd.vue'),
  },
  {
    path: '/creation',
    name: 'creation',
    component: () => import('../views/creation.vue'),
  },
  {
    path: '/trans',
    name: 'trans',
    component: () => import('../views/trans.vue'),
  },
  {
    path: '/spoken',
    name: 'spoken',
    component: () => import('../views/spoken.vue'),
  },
  {
    path: '/recording',
    name: 'recording',
    component: () => import('../views/recording.vue'),
  },
  {
    path: '/model',
    name: 'model',
    component: () => import('../views/model.vue'),
  },
  {
    path: '/knowledge',
    name: 'knowledge',
    component: () => import('../views/knowledge.vue'),
  },
  {
    path: '/assistant',
    name: 'assistant',
    component: () => import('../views/assistant.vue'),
  },
  {
    path: '/setting',
    name: 'setting',
    component: () => import('../views/setting.vue'),
  },
  {
    path: '/help',
    name: 'help',
    component: () => import('../views/help.vue'),
  },
  // 添加一个通配符路由，处理未匹配到的路径
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => import('../components/common/NotFound.vue')
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

export default router;