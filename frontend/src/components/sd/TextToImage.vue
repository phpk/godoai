<script setup lang="ts">
import { ref, toRaw, watch } from "vue";
import { notifyError } from "@/util/msg";
import { t } from "@/i18n/index";
const props = defineProps({
  info: {
    type: Object,
    default: () => ({}),
  },
});
const showTips = ref(false)
const sampling_methods = [
  "euler",
  "euler_a",
  "heun",
  "dpm2",
  "dpm++2s_a",
  "dpm++2m",
  "dpm++2mv2",
  "lcm",
];
const params = ref({
  action: "txt2img",
  seed: 42,
  strength: 0.75,
  //enable_hr: false,
  width: 512,
  height: 512,
  sampling_method: "euler_a",
  cfg_scale: 7,
  negative_prompt:
    "out of frame, lowers, text, error, cropped, worst quality, low quality, jpeg artifacts, ugly, duplicate, morbid, mutilated, out of frame, extra fingers, mutated hands, poorly drawn hands, poorly drawn face, mutation, deformed, blurry, dehydrated, bad anatomy, bad proportions, extra limbs, cloned face, disfigured, gross proportions, malformed limbs, missing arms, missing legs, extra arms, extra legs, fused fingers, too many fingers, long neck, username, watermark, signature",
  prompt: "",
  num: 1,
  steps: 10,
});

const emit = defineEmits(["saveFn"]);
const imageGenerate = () => {
  const save = toRaw(params.value);
  if (save.prompt == "") {
    notifyError(t('invalidPrompt'));
    return;
  }
  emit("saveFn", save);
};
const isCreating = ref(false);
// 使用 watch 监听 props.info 的变化
watch(
  () => props.info,
  (newInfo, _) => {
    const newData = toRaw(newInfo);
    if (newData && newData.status != "success") {
      // 当 status 改变为 "success" 时，设置 isCreating 为 true
      isCreating.value = true;
    } else {
      // 其他情况下，设置 isCreating 为 false
      isCreating.value = false;
    }
  },
  { deep: true } // 使用 deep: true 以监听对象属性的变化
);
</script>

<template>
  <el-dialog
    v-model="showTips"
    title="Tips"
    width="88%"
  >
    <ImagePrompt />
  </el-dialog>
  <el-card>
    <div class="card-content">
      <!-- Prompt -->
      <h3 class="title">
        <el-button @click="showTips = true">{{ t('sd.promptTitle') }}</el-button>
      </h3>
      <el-input
        type="textarea"
        :row="3"
        v-model="params.prompt"
        :placeholder="t('sd.promptPlaceholder')"
        class="my-3"
        clearable
      ></el-input>

      <!-- Width -->
      <h3 class="title">{{ t('sd.widthTitle') }}</h3>
      <el-slider
        v-model="params.width"
        :max="1024"
        :min="8"
        :step="8"
        class="my-3"
      ></el-slider>

      <!-- Height -->
      <h3 class="title">{{ t('sd.heightTitle') }}</h3>
      <el-slider
        v-model="params.height"
        :max="1024"
        :min="8"
        :step="8"
        class="my-3"
      ></el-slider>

      <!-- More Settings -->
      <el-collapse accordion class="mb-5">
        <el-collapse-item :title="t('sd.moreSettingsTitle')">
          <div slot="content">
            <!-- Steps -->
            <h3 class="title">{{ t('sd.stepsTitle') }}</h3>
            <el-slider
              v-model="params.steps"
              :min="1"
              :max="60"
              :step="1"
              class="my-3"
            ></el-slider>

            <!-- Num -->
            <h3 class="title">{{ t('sd.numImagesTitle') }}</h3>
            <el-slider
              v-model="params.num"
              :min="1"
              :max="6"
              :step="1"
              class="my-3"
            ></el-slider>

            <!-- Cfg Scale -->
            <h3 class="title">{{ t('sd.cfgScaleTitle') }}</h3>
            <el-slider
              v-model="params.cfg_scale"
              :min="1"
              :max="20"
              :step="1"
              class="my-3"
            ></el-slider>
            <h3 class="title">{{ t('sd.strengthTitle') }}</h3>
            <el-slider
              v-model="params.strength"
              :max="1"
              :min="0"
              :step="0.01"
              class="my-3"
            ></el-slider>
            <h3 class="title">{{ t('sd.seedTitle') }}</h3>
            <el-slider
              v-model="params.seed"
              :max="100"
              :min="-1"
              :step="1"
              class="my-3"
            ></el-slider>
            <h3 class="title">{{ t('sd.samplerTitle') }}</h3>
          <el-select class="my-3" v-model="params.sampling_method">
            <el-option
              v-for="(el, key) in sampling_methods"
              :key="key"
              :label="el"
              :value="el"
            />
          </el-select>
            <!-- Negative Prompt -->
            <h3 class="title">{{ t('sd.negativePromptTitle') }}</h3>
            <el-input
              type="textarea"
              v-model="params.negative_prompt"
              :placeholder="t('sd.negativePromptPlaceholder')"
              class="my-3"
              rows="4"
              clearable
            ></el-input>
          </div>
        </el-collapse-item>
      </el-collapse>
      <el-row justify="center" style="margin-top: 10px">
        <el-button
          size="large"
          type="info"
          block
          :loading="isCreating"
          @click="imageGenerate"
          plain
        >
        {{ t('sd.generateButtonText') }}
        </el-button>
      </el-row>
    </div>
  </el-card>
</template>

<style scoped lang="scss">
.title {
  font-weight: 500;
  font-size: 14px;
  margin-top: 10px;
  margin-bottom: 10px;
  text-align: left;
  justify-content: left;
}
.my-3 {
  margin-bottom: 10px;
  width: 90%;
  margin-left: 20px;
}
</style>
