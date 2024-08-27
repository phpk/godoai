<script lang="ts" setup>
import { onMounted, ref } from 'vue';
import { useDark, useToggle } from "@vueuse/core";
import { useModelStore } from "@/stores/model";
import { changeLang } from "@/i18n";
import { useI18n } from 'vue-i18n';
import { RestartApp } from "@/util/goutil";
import { useUpgradeStore } from '@/stores/upgrade';
import { t } from "@/i18n";
const { locale } = useI18n()
const isDark = useDark();
const toggleDark = useToggle(isDark);
const modelStore = useModelStore();
const progressModalRef: any = ref(null);
const updateStore = useUpgradeStore();
async function clearSystem() {
  await modelStore.clearSystem();
  RestartApp()
}
function changeLanguage() {
  const lang = changeLang()
  locale.value = lang
  //location.reload()
  RestartApp()
}
onMounted(async () => {
  await updateStore.checkUpdate()
})

</script>

<template>
  <ProgressModal :title="t('updatePrompt.title')" ref="progressModalRef" />
  <el-menu class="el-menu-demo" mode="horizontal" :ellipsis="false">
    <el-menu-item index="0">
      <img style="width: 30px" src="/logo.png" alt="AIOK" />
    </el-menu-item>
    <div class="flex-grow" />
    <el-menu-item @click="changeLanguage()">
      <el-tooltip class="box-item" effect="dark" :content="t('setting.switchLang')" placement="bottom-end">
        <el-icon>
          <Location />
        </el-icon>
      </el-tooltip>
    </el-menu-item>
    <el-menu-item @click="clearSystem()">
      <el-tooltip class="box-item" effect="dark" :content="t('setting.clearSystem')" placement="bottom-end">
        <el-icon>
          <Brush />
        </el-icon>
      </el-tooltip>
    </el-menu-item>
    <el-menu-item @click="toggleDark()">
      <el-tooltip class="box-item" effect="dark" :content="t('setting.switchLang')" placement="bottom-end">
        <el-icon v-if="isDark">
          <Moon />
        </el-icon>
        <el-icon v-else>
          <Sunny />
        </el-icon>
      </el-tooltip>
    </el-menu-item>
  </el-menu>
</template>
