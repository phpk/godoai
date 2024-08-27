<script lang="ts" setup>
import DefaultLayout from '@/components/layouts/DefaultLayout.vue';
import SimpleLayout from '@/components/layouts/SimpleLayout.vue';
import { useModelStore } from "@/stores/model";
import { onMounted, ref, watch } from "vue";
import { useRouter } from 'vue-router';
const currentLayout = ref(DefaultLayout);
const store = useModelStore();
const router = useRouter();
onMounted(() => {
  store.initSystem();
});
watch(() => router.currentRoute.value, (newRoute) => {
  //console.log('New route:', newRoute);
  if (newRoute.name === 'ask') {
    currentLayout.value = SimpleLayout;
  } else {
    currentLayout.value = DefaultLayout;
  }
}, { immediate: true });
</script>
<template>
  <el-config-provider size="large">
    <!-- <component :is="currentLayout"></component> -->
      <component :is="currentLayout">
      </component>
  </el-config-provider>
</template>

<style>
#app {
  text-align: center;
  color: var(--el-text-color-primary);
}

</style>