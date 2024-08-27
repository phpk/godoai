<script lang="ts" setup>
import { useCreationStore } from "@/stores/creation";
import { t } from "@/i18n/index";
const creationStore:any = useCreationStore();
function changeSystemPrompt(val:any) {
  const data = creationStore.config.systemList.find((item:any) => item.name == val);
  creationStore.config.systemPrompt = data.prompt;
}
function changeLeaderPrompt(val:any) {
  const data = creationStore.config.leaderList.find((item:any) => item.name == val);
  creationStore.config.leaderPrompt = data.prompt;
}
function changeBuilderPrompt(val:any) {
  const data = creationStore.config.builderList.find((item:any) => item.name == val);
  creationStore.config.builderPrompt = data.prompt;
}
</script>
<template>
  <el-form label-width="150px" style="margin-top: 12px">
    <el-form-item :label="t('creation.model')">
      <el-select v-model="creationStore.config.model">
        <el-option
          v-for="(item, key) in creationStore.config.modelList"
          :key="key"
          :label="item.model"
          :value="item.model"
        />
      </el-select>
    </el-form-item>
    <!-- <el-form-item :label="t('chat.refknow')">
      <el-select
        v-model="creationStore.config.refKnowledge"
        :clearable="true"
        :filterable="true"
      >
        <el-option
          v-for="(item, key) in creationStore.config.knowList"
          :key="key"
          :label="item.name"
          :value="item.uuid"
        />
      </el-select>
    </el-form-item> -->
    <el-form-item :label="t('creation.sysPrompt')">
      <el-select
        v-model="creationStore.config.systemPromptTitle"
        @change="changeSystemPrompt"
      >
        <el-option
          v-for="(item, key) in creationStore.config.systemList"
          :key="key"
          :label="item.name"
          :value="item.name"
        />
      </el-select>
    </el-form-item>
    <el-form-item :label="t('creation.sysPromptContent')">
      <el-input
        type="textarea"
        v-model="creationStore.config.systemPrompt"
        :rows="2"
      ></el-input>
    </el-form-item>

    <el-form-item :label="t('creation.outlinePrompt')">
      <el-select
        v-model="creationStore.config.leaderPromptTitle"
        @change="changeLeaderPrompt"
      >
        <el-option
          v-for="(item, key) in creationStore.config.leaderList"
          :key="key"
          :label="item.name"
          :value="item.name"
        />
      </el-select>
    </el-form-item>
    <el-form-item :label="t('creation.outlinePromptContent')">
      <el-input
        type="textarea"
        v-model="creationStore.config.leaderPrompt"
        :rows="2"
      ></el-input>
    </el-form-item>

    <el-form-item :label="t('creation.contentPrompt')">
      <el-select
        v-model="creationStore.config.builderPromptTitle"
        @change="changeBuilderPrompt"
      >
        <el-option
          v-for="(item, key) in creationStore.config.builderList"
          :key="key"
          :label="item.name"
          :value="item.name"
        />
      </el-select>
    </el-form-item>
    <el-form-item :label="t('creation.contentPromptContent')">
      <el-input
        type="textarea"
        v-model="creationStore.config.builderPrompt"
        :rows="2"
      ></el-input>
    </el-form-item>

  </el-form>
  </template>