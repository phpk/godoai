<script setup lang="ts">
import { t } from "@/i18n/index";
import { computed, onMounted, ref } from "vue";
import { notifyError } from "@/util/msg";
import { getSystemConfig, getCurrentModelName, getChatConfig } from "@/stores/config.ts";

const langs = [
  {
    code: "en",
    name: "English",
    label: "English",
  },
  {
    code: "zh-CN",
    name: "Chinese Simplified",
    label: "中文(简体)",
  },
  {
    code: "zh-TW",
    name: "Chinese Traditional",
    label: "中文(繁體)",
  },
  {
    code: "ja",
    name: "Japanese",
    label: "日本語",
  },
  {
    code: "ko",
    name: "Korean",
    label: "한국어",
  },
  {
    code: "fr",
    name: "French",
    label: "Français",
  },
  {
    code: "de",
    name: "German",
    label: "Deutsch",
  },
  {
    code: "es",
    name: "Spanish",
    label: "Español",
  },
];

const currentLang = ref({
  code: "en",
  name: "English",
  label: "English",
});

function changeLang(code: string) {
  const data: any = langs.find((item) => item.code == code);
  currentLang.value = data;
}
const baseContent = ref("");
const targetContent = ref("");

const prompt = computed(() => {
  return `Translate into ${currentLang.value.name}`;
});

const isLoading = ref(false);
onMounted(async () => {
  const hasCurrent = getCurrentModelName("translation");
  if (!hasCurrent) {
    notifyError(t("index.notFindChatModel"));
    return;
  }
});

const translate = async () => {
  const config = getSystemConfig();
  if (baseContent.value === "") {
    isBaseContentEmpty.value = true;
    return;
  }
  isLoading.value = true;
  const chatUrl = config.chatApi;
  const chatConfig = getChatConfig("translation");
  const postData: any = {
    method: "POST",
    body: JSON.stringify({
      messages: [
        { role: "system", content: prompt.value },
        { role: "user", content: baseContent.value },
      ],
      message: baseContent.value,
      prompt: prompt.value,
      //model: "qwen2-0.5b",
      model: getCurrentModelName("translation"),
      stream: false,
      options: chatConfig,
    }),
  };
  const completion = await fetch(chatUrl, postData);
  if (!completion.ok || completion.status !== 200) {
    notifyError(completion.statusText);
    return;
  }
  const res = await completion.json();
  //console.log(res)
  const msg = res.choices[0].message.content;
  targetContent.value = msg;

  isLoading.value = false;
};

const isBaseContentEmpty = ref(false);
</script>

<template>
  <el-row style="padding: 12px;border-bottom: none;" justify="space-around">
    <el-col :span="11">
      <el-card>
        <template #header>
          <el-button
            icon="Compass"
            :loading="isLoading"
            :disabled="isLoading"
            @click="translate"
            >{{ t("model.translation") }}</el-button
          >
        </template>
        <el-input
          type="textarea"
          :rows="20"
          placeholder=""
          v-model="baseContent"
          @focus="isBaseContentEmpty = false"
        >
        </el-input>
        <template #footer>
          <el-row justify="end" style="border-bottom: none;margin:0;padding:0">
            <copy-btn :text="baseContent" />
          </el-row>
        </template>
      </el-card>
    </el-col>
    <el-col :span="12">
      <el-card>
        <template #header>
          <div class="card-header">
            <el-select v-model="currentLang.code" @change="changeLang">
              <el-option
                v-for="(item, key) in langs"
                :key="key"
                :label="item.label"
                :value="item.code"
              />
            </el-select>
          </div>
        </template>
        <el-input
          type="textarea"
          :rows="20"
          v-model="targetContent"
          @focus="isBaseContentEmpty = false"
        >
        </el-input>
        <template #footer>
          <el-row justify="end" style="border-bottom: none;margin:0;padding:0">
            <copy-btn :text="targetContent" />
          </el-row>
        </template>
      </el-card>
    </el-col>
  </el-row>
</template>

<style scoped lang="scss"></style>
