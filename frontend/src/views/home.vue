<script setup lang="ts">
import { Vue3Lottie } from "vue3-lottie";
import { useRouter } from "vue-router";
import { t } from "@/i18n";
const router = useRouter();
</script>
<template>
  <div class="home-container">
    <h1>GodoAI</h1>
    <Vue3Lottie animationLink="/bot/home.json" :height="360" :width="360" />
  </div>
  <el-space wrap>
    <el-button plain type="info" @click="router.push('/model')">{{ t("common.model") }}</el-button>
    <el-button plain type="success" @click="router.push('/knowledge')">{{ t('model.knowledge') }}</el-button>
    <el-button plain type="danger" @click="router.push('/creation')">{{ t('model.creation') }}</el-button>
  </el-space>
</template>
<style lang="scss">
.home-container {
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
  padding-top: 100px;

  h1 {
    font-size: 2rem;
    font-weight: 500;
  }
}
</style>
