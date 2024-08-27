<script lang="ts" setup>
import { ref } from "vue";
import { t } from "@/i18n/index";
const linkUrl = ref("")
const emit = defineEmits(["url-saved"]);
function saveData() {
  emit("url-saved", linkUrl.value)
}
</script>

<template>
  <el-form label-width="100px" style="margin-top: 12px">
    <el-form-item label="Url">
      <el-input 
      v-model="linkUrl" 
      placeholder="https://xxx.com" 
      autofocus
        clearable/>
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="saveData">{{ t('common.save') }}</el-button>
    </el-form-item>
  </el-form>

  
</template>
