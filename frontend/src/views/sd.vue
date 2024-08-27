<script setup lang="ts">
import { computed, onMounted, ref, nextTick } from "vue";
import {
  getSystemKey,
  parseJson,
} from "@/stores/config.ts";
import { useSdStore } from "@/stores/sd";
import { t } from "@/i18n/index";
import { notifyError, notifySuccess } from "@/util/msg";
import { ElScrollbar } from "element-plus";
const messageContainerRef = ref<InstanceType<typeof ElScrollbar>>();
const messageInnerRef = ref<HTMLDivElement>();

const sdStore = useSdStore();
onMounted(async () => {
  try {
    await sdStore.initPage();
    if (!sdStore.currentModel) {
      notifyError(t('sd.noImageModel'));
      return;
    }
  } catch (e:any) {
    notifyError(e.message);
  }
});

const imgIofo: any = ref({});

const messages: any = ref([]);
// const apiUrl = systemStore.getImageUrl();
const scrollToBottom = () => {
  nextTick(() => {
    if (messageContainerRef && messageContainerRef.value) {
      messageContainerRef.value!.setScrollTop(messageInnerRef.value!.clientHeight);
    }
  });
};
async function saveImg(data: any) {
  //console.log(data);
  const modelData:any = sdStore.getModel();
  console.log(modelData);
  if (!modelData) {
    notifyError(t('sd.modelNotExist'));
    return;
  }
  data.model = modelData.model;
  data.file_name = modelData.file_name;
  imgIofo.value.status = "creating";
  try {
    const completion = await fetch(getSystemKey("modelUrl") + "/image", {
      method: "POST",
      body: JSON.stringify(data),
    });
    //console.log(completion);
    if (!completion.ok) {
      notifyError(completion.statusText);
      imgIofo.value.status = "success";
      //await reStartApp("sd");
      return;
    }
    const reader: any = completion.body?.getReader();
    if (!reader) {
      notifyError(t("common.cantStream"));
      imgIofo.value.status = "success";
    }
    const decoder = new TextDecoder("utf-8");

    const read = async (): Promise<any> => {
      const { done, value } = await reader.read();
      if (done) {
        messages.value = [];
        imgIofo.value.status = "success";
        sdStore.addImgList(imgIofo.value.imgList);
        return reader.releaseLock();
      }
      const chunk = decoder.decode(value, { stream: true });
      //console.log(chunk);
      const res = parseJson(chunk);
      if (res && res.message) {
        imgIofo.value.imgList = res.image_list;
        messages.value = [...messages.value, ...res.message];
        scrollToBottom();
      }
      return read();
    };
    await read();
  } catch (e) {
    console.log(e);
    imgIofo.value.status = "success";
  }
}

const activeName = ref("txt2img");
async function deleteImage(img: string) {
  await sdStore.delImgList(img);
  const completion = await fetch(getSystemKey("sdUrl") + "/delete", {
    method: "POST",
    body: JSON.stringify({ path: img }),
  });
  //console.log(completion);
  if (!completion.ok) {
    notifyError(completion.statusText);
    return;
  }
  const res = await completion.json();
  if (res.code != 0) {
    notifyError("delete image failed");
    return;
  }

  notifySuccess("delete image success");
}
const imgList = computed(() => {
  const list: any = [];
  sdStore.imageList.forEach((img) => {
    list.push(getSystemKey("sdUrl") + "/view?path=" + img);
  });
  return list;
});
</script>

<template>
  <el-row justify="space-around" class="main" style="border-bottom: 0;">
    <el-col :span="12">
      <el-scrollbar class="scrollbar-height">
        <div style="margin: 12px; height: 88vh">
          <el-tabs v-model="activeName" class="sd-tabs">
            <el-tab-pane :label="t('sd.textToImage')" name="txt2img">
              <text-to-image @saveFn="saveImg" :info="imgIofo" />
            </el-tab-pane>
            <el-tab-pane :label="t('sd.imageToImage')" name="img2img">
              <image-to-image @saveFn="saveImg" :info="imgIofo" />
            </el-tab-pane>
            <el-tab-pane :label="t('sd.models')" name="models">
              <el-card>
                <div class="card-content">
                  <el-form label-width="100px">
                    <el-form-item :label="t('setting.defModel')">
                      <el-select
                        v-model="sdStore.currentModel"
                        @change="(val:any) => sdStore.updateModel(val)"
                      >
                        <el-option
                          v-for="(el, key) in sdStore.modelList"
                          :key="key"
                          :label="el.model"
                          :value="el.model"
                        />
                      </el-select>
                    </el-form-item>
                  </el-form>
                </div>
              </el-card>
            </el-tab-pane>
          </el-tabs>
          <div class="command-line" v-if="messages.length > 0">
            <el-scrollbar ref="messageContainerRef">
              <div class="terminal-window" ref="messageInnerRef">
                <div class="output" ref="output">
                  <div v-for="(line, index) in messages" :key="index" class="line">
                    {{ line }}
                  </div>
                </div>
              </div>
            </el-scrollbar>
          </div>
        </div>
      </el-scrollbar>
    </el-col>
    <el-col :span="12">
      <el-card style="min-height: 536px;">
        <template #header>
          <div class="card-header">{{ t('sd.imgList') }}</div>
        </template>
        <div class="image-grid">
          <div
            v-for="(image, index) in sdStore.imageList"
            :key="index"
            class="image-item"
          >
            <el-image
              :src="getSystemKey('sdUrl') + '/view?path=' + image"
              alt="Image"
              :zoom-rate="1.2"
              :max-scale="7"
              :min-scale="0.2"
              :preview-src-list="imgList"
              :initial-index="index"
              fit="cover"
              class="image"
            />
            <el-button
              icon="Delete"
              size="small"
              class="delete-button"
              @click.stop="deleteImage(image)"
              circle
            ></el-button>
          </div>
        </div>
        <el-row justify="center" style="margin-top: 15px" v-if="sdStore.page.pages > 1">
          <el-pagination
            background
            layout="prev, pager, next"
            v-model:current-page="sdStore.page.current"
            v-model:page-size="sdStore.page.size"
            :total="sdStore.page.total"
            @current-change="(val:any) => sdStore.pageClick(val)"
          />
        </el-row>
      </el-card>
    </el-col>
  </el-row>
</template>

<style scoped lang="scss">
.main {
  width: 96%;
  margin: 2%;
}
.scrollbar-height {
  height:calc(100vh - 180px);
}
.sd-tabs > .el-tabs__content {
  padding: 0px;
  color: #6b778c;
  font-size: 32px;
  font-weight: 600;
}
.image-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
}

.image-item {
  position: relative;
  border: 1px solid white; /* 白边效果 */
  overflow: hidden;
}

.image {
  width: 100%;
  height: auto;
}

.delete-button {
  position: absolute;
  top: 0;
  right: 0;
  padding: 5px 10px;
  cursor: pointer;
  opacity: 0; /* 默认隐藏 */
  transition: opacity 0.3s ease;
}

.image-item:hover .delete-button {
  opacity: 1; /* 鼠标悬停时显示 */
}

.command-line {
  width: 100%;
  height: 200px;
  margin-top: 20px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.terminal-window {
  width: 90%;
  background-color: #000;
  color: lime;
  padding: 1rem;
  border-radius: 5px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.5);
  overflow-y: auto;
}

.output {
  margin-top: 1rem;
}

.line {
  white-space: pre-wrap;
  font-size: 12px;
  text-align: left;
}
</style>
