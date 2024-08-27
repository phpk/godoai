<script setup lang="ts">
import { ref, toRaw, watch } from "vue";
import { getSystemKey } from "@/stores/config.ts";
import { notifySuccess, notifyError } from "@/util/msg";
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
  action: "img2img",
  seed: 42,
  //enable_hr: false,
  strength: 0.75,
  sampling_method: "euler_a",
  width: 512,
  height: 512,
  cfg_scale: 7,
  negative_prompt:
    "out of frame, lowers, text, error, cropped, worst quality, low quality, jpeg artifacts, ugly, duplicate, morbid, mutilated, out of frame, extra fingers, mutated hands, poorly drawn hands, poorly drawn face, mutation, deformed, blurry, dehydrated, bad anatomy, bad proportions, extra limbs, cloned face, disfigured, gross proportions, malformed limbs, missing arms, missing legs, extra arms, extra legs, fused fingers, too many fingers, long neck, username, watermark, signature",
  url: "",
  prompt: "",
  steps: 10,
  num: 1,
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
const fileInput: any = ref(null);
function selectFile() {
  fileInput.value.click();
}
const uploadFile = async (event: any) => {
  const files = event.target.files;
  if (!files || files.length < 1) {
    return;
  }
  const formData: any = new FormData();
  formData.append("files", files[0]);
  const response = await fetch(getSystemKey("modelUrl") + "/uploadimage", {
    method: "POST",
    // headers:{
    //   "Content-Type" : "multipart/form-data"
    // },
    body: formData,
  });
  if (response.ok) {
    const data = await response.json();
    if (data.path) {
      params.value.url = data.path;
      notifySuccess(t('common.uploadSuccess'));
    }
    console.log(params.value);
  } else {
    notifyError(t('common.uploadFailed'));
  }
};
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
      <el-row>
        <el-col :span="24">
          <el-button icon="Upload" @click="selectFile">{{ t('sd.uploadImg') }}</el-button>
          <input
            type="file"
            ref="fileInput"
            accept="image/png"
            style="display: none"
            @change="uploadFile"
          />
        </el-col>
      </el-row>
      <el-divider></el-divider>

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
      <el-collapse accordion>
        <el-collapse-item :title="t('sd.moreSettingsTitle')">
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

          <!-- Steps -->
          <h3 class="title">{{ t('sd.stepsTitle') }}</h3>
          <el-slider
            v-model="params.steps"
            :max="60"
            :min="1"
            :step="1"
            class="my-3"
          ></el-slider>

          <!-- Num -->
          <h3 class="title">{{ t('sd.numImagesTitle') }}</h3>
          <el-slider
            v-model="params.num"
            :max="6"
            :min="1"
            :step="1"
            class="my-3"
          ></el-slider>

          <!-- CFG Scale -->
          <h3 class="title">{{ t('sd.cfgScaleTitle') }}</h3>
          <el-slider
            v-model="params.cfg_scale"
            :max="20"
            :min="1"
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
