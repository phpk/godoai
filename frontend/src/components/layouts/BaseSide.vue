<template>
  <el-menu
    default-active="0"
    class="el-menu-vertical"
    :collapse="isCollapse"
    @open="handleOpen"
    @close="handleClose"
  >
    <el-menu-item :index="String(key)" v-for="(item,key) in menuList" @click="toRoute(item.path, String(key))">
      <el-icon><component class="password-icon" :is="item.icon" /></el-icon>
      <template #title>{{ item.name }}</template>
    </el-menu-item>
  </el-menu>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { useRouter } from 'vue-router'
const router = useRouter()
const isCollapse = ref(true);
const activeId = ref('0');
import {t} from "@/i18n";
const menuList = [
  {
    name: t('common.home'),
    icon: 'HomeFilled',
    path: '/'
  },
  {
    name: t('common.chat'),
    icon: 'ChatDotRound',
    path: '/chat'
  },
  {
    name: t('model.knowledge'),
    icon: 'Document',
    path: '/knowledge'
  },
  {
    name: t('common.draw'),
    icon: 'PictureFilled',
    path: '/sd'
  },
  {
    name: t('model.creation'),
    icon: 'Reading',
    path: '/creation'
  },
  {
    name: t('model.translation'),
    icon: 'SetUp',
    path: '/trans'
  },
  {
    name: t('model.spoken'),
    icon: 'Microphone',
    path: '/spoken'
  },
  {
    name: t('model.recording'),
    icon: 'Headset',
    path: '/recording'
  },
  {
    name: t('common.model'),
    icon: 'Collection',
    path: '/model'
  },
  {
    name: t('model.assistant'),
    icon: 'Compass',
    path: '/assistant'
  },
  {
    name: t('common.setting'),
    icon: 'Setting',
    path: '/setting'
  },
  {
    name: t('common.help'),
    icon: 'QuestionFilled',
    path: '/help'
  }

]
function toRoute(path: string, id: string){
  activeId.value = id;
  router.push(path)
}
const handleOpen = (key: string, keyPath: string[]) => {
  console.log(key, keyPath);
};
const handleClose = (key: string, keyPath: string[]) => {
  console.log(key, keyPath);
};
</script>
<style lang="scss" scoped>
.el-menu-vertical:not(.el-menu--collapse) {
  width: 65px;
}
.el-menu-item {
  height: 50px;
}
</style>
