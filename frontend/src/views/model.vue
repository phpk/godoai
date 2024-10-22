<script setup lang="ts">
import { onMounted, ref, toRaw } from "vue";
import { t, getLang } from "@/i18n/index.ts";
import { useModelStore } from "@/stores/model.ts";
import { getSystemConfig, setCurrentModel } from "@/stores/config.ts";
import { notifySuccess, notifyError, buyVip } from "@/util/msg.ts";
import { Vue3Lottie } from "vue3-lottie";
import { Search } from "@element-plus/icons-vue";
import { useRouter } from 'vue-router'
const router = useRouter();
const currentLang = getLang()
const modelStore = useModelStore();
const config = getSystemConfig();
const downLeft = ref(false);
const downAdd = ref(false);
const labelEditor = ref(false);
const labelId = ref(0);
const searchKey = ref("");
const currentCate = ref("all");

const help_label = ref();
const help_adddown = ref();
const help_addlabel = ref();
const help_showdown = ref();
const showDetail = ref(false);
const detailModel = ref("")
let downloadAbort:any = {};
onMounted(async () => {
  await modelStore.getList();
});

async function showCate(name: any) {
  await modelStore.getLabelCate(name.paneName);
}
async function showSearch() {
  await modelStore.getLabelSearch(searchKey.value);
}

async function downLabel(modelData:any, labelData:any) {
  labelData = toRaw(labelData);
  const saveData = {
    model: modelData.model,
    label: labelData.name,
    action: labelData.action,
    engine: labelData.engine,
    url: modelData.url ?? [],
    from: modelData.from ? modelData.from : labelData.from,
    type: modelData.type ?? "",
    file_name: modelData.file_name ?? "",
    //options: modelData.options ?? {},
    params: modelData.params ?? {},
    info: modelData.info ?? {},
  };
  await download(saveData);
}

async function download(saveData:any) {
  saveData = toRaw(saveData);
  saveData.info = toRaw(saveData.info);
  saveData.url = toRaw(saveData.url);
  saveData.params = toRaw(saveData.params);
  downAdd.value = false;
  downLeft.value = true;
  const has = modelStore.checkDownload(saveData.model);
  if (has) {
    notifyError(t('model.labelDown'));
    return;
  }
  //console.log(saveData);
  const downUrl = config.modelDownApi;

  try {
    const completion = await fetch(downUrl, {
      method: "POST",
      body: JSON.stringify(saveData),
    });
    downloadAbort[saveData.model] = false;
    if (!completion.ok) {
      const errorData = await completion.json();
      notifyError(errorData.message);
      return;
    }

    saveData.status = "loading";
    saveData.progress = 0;
    modelStore.addDownload(saveData);
    await handleDown(saveData, completion);
  } catch (error:any) {
    notifyError(error.message);
  }
}
function cancelDownload(model: string) {
  modelStore.deleteDownload(model);
  if (!downloadAbort[model]) {
    downloadAbort[model] = true;
    notifySuccess(t('model.downChanel'));
  }
}

async function handleDown(modelData:any, completion:any) {
  const reader: any = completion.body?.getReader();
  if (!reader) {
    notifyError(t("common.cantStream"));
  }

  while (true) {
    try {
      const { done, value } = await reader.read();
      if (done) {
        //console.log("has done!");
        reader.releaseLock();
        break;
      }
      if (downloadAbort[modelData.model]) {
        break;
      }
      const rawjson = new TextDecoder().decode(value);
      //console.log(rawjson);
      const msg = modelStore.parseMsg(rawjson);
      if(msg.message && msg.code) {
        notifyError(msg.message);
        break;
      }
      if (msg.status == "") {
        continue;
      }
      modelData.status = msg.status;
      if (msg.total && msg.completed && msg.total > 0) {
        modelData.isLoading = 1;
        modelData.progress = Math.ceil((msg.completed / msg.total) * 100);
      } else {
        modelData.progress = 0;
      }
      if (msg.status == "success") {
        modelData.isLoading = 0;
        modelData.progress = 0;
      }
      //console.log(modelData);
      await modelStore.updateDownload(modelData);
      if (msg.status == "success") {
        modelStore.deleteDownload(modelData.model);
        setCurrentModel(toRaw(modelData.action), modelData.model, true);
      }
    } catch (error) {
      console.error("An error occurred:", error);
      break;
    }
  }
}

async function deleteModel(modelData: any) {
  modelData = toRaw(modelData);
  //console.log(modelData)
  try {
    await modelStore.deleteModelList(modelData);
    //return
    const postData = {
      method: "POST",
      body: JSON.stringify({
        url: modelData.url,
        model: modelData.model,
        engine: modelData.engine,
      }),
    };
    const delUrl = config.modelDeleteApi;
    const completion = await fetch(delUrl, postData);
    if (completion.status === 404) {
      notifyError(completion.statusText);
      return;
    }
    if (completion.status === 200) {
      notifySuccess("success!");
    }
  } catch (error:any) {
    console.log(error);
    notifyError(error.message);
  }
}
function labelShow(val:any) {
  labelId.value = val;
  labelEditor.value = true;
}
function closeLabel() {
  labelId.value = 0;
  labelEditor.value = false;
}
async function refreshList() {
  modelStore.labelList = await modelStore.getLabelList();
}
async function delLabel(id: number) {
  await modelStore.delLabel(id);
  notifySuccess("success!");
}
function getModelStatus(model: string) {
  let name = t('model.noDown');
  if (modelStore.modelList.find((item:any) => item.model === model)) {
    name = t('model.hasDown');
  }
  if (modelStore.downList.find((item:any) => item.model === model)) {
    name = t('model.downloading');
  }
  return name;
}
function showModel(model: string) {
  detailModel.value = model;
  showDetail.value = true;
}
function showbuyVip(){
  buyVip().then(()=>{
    router.push('/help')
  })
}
</script>
<template>
   <el-dialog v-model="showDetail" width="600" append-to-body>
    <DownModelInfo :model="detailModel" />
  </el-dialog>
  <div class="app-container">
    <el-drawer
      v-model="downLeft"
      direction="ltr"
      :show-close="false"
      :with-header="false"
      :size="300"
    >
      <div>
        <el-tag size="large" style="margin-bottom: 10px">{{ t('model.downloading') }}</el-tag>
        <div class="pa-2">
          <Vue3Lottie
            animationLink="/bot/search.json"
            :height="200"
            :width="200"
            v-if="modelStore.downList.length < 1"
          />
          <el-space direction="vertical" v-else>
            <el-card
              v-for="(val, key) in modelStore.downList"
              :key="key"
              class="box-card"
              style="width: 250px"
            >
              <div class="card-header">
                <span>{{ val.model }}</span>
              </div>
              <div class="text item" v-if="val.progress && val.progress > 0">
                <el-progress
                  :text-inside="true"
                  :stroke-width="15"
                  :percentage="val.progress"
                />
              </div>
              <div class="drawer-model-actions" style="margin-top: 10px">
                <el-tag size="small">{{ val.status }}</el-tag>
                <el-icon :size="18" color="red" @click="cancelDownload(val.model)">
                  <Delete />
                </el-icon>
                <el-icon
                  :size="18"
                  color="blue"
                  v-if="val.progress > 0 && val.isLoading < 1 && val.status != 'success'"
                  @click="download(toRaw(val))"
                >
                  <VideoPlay />
                </el-icon>
              </div>
            </el-card>
          </el-space>
        </div>
        <el-tag size="large" style="margin: 10px auto">{{ t('model.hasDown') }}</el-tag>
        <div class="pa-2">
          <div
            class="list-item"
            v-for="(item, index) in modelStore.modelList"
            :key="index"
          >
            <div class="list-title" @click="showModel(item.model)">
              {{ item.model }}
            </div>
            <el-button
              class="delete-btn"
              icon="Delete"
              size="small"
              @click.stop="deleteModel(item)"
              circle
            ></el-button>
          </div>
        </div>
      </div>
    </el-drawer>
    <el-dialog v-model="labelEditor" width="600" :title="t('model.modelLabel')">
      <down-labeleditor
        @closeFn="closeLabel"
        @refreshFn="refreshList"
        :labelId="labelId"
      />
    </el-dialog>
    
    <el-page-header icon="null">
      <template #title>
        <div></div>
      </template>
      <template #content>
        <el-button
          @click.stop="downLeft = !downLeft"
          icon="Menu"
          circle
          ref="help_showdown"
        />
        <el-button @click.stop="showbuyVip" icon="Plus" circle ref="help_adddown" />
        <el-button
          @click.stop="labelShow(0)"
          icon="CollectionTag"
          circle
          ref="help_addlabel"
        />
      </template>
      <template #extra>
        <el-space class="mr-10">
          <el-input
            :placeholder="t('model.search')"
            v-model="searchKey"
            v-on:keydown.enter="showSearch"
            style="width: 200px"
            :suffix-icon="Search"
          />
        </el-space>
      </template>
    </el-page-header>

    <div class="flex-fill ml-10 mr-10">
      <el-tabs v-model="currentCate" @tab-click="showCate" ref="help_label">
        <el-tab-pane :label="t('model.all')" name="all" />
        <el-tab-pane
          :label="t('model.' + item)"
          :name="item"
          v-for="(item, key) in modelStore.cateList"
          :key="key"
        />
      </el-tabs>
    </div>

    <el-scrollbar class="scrollbarHeightList">
      <div class="model-list">
        <div
          v-for="item in modelStore.labelList"
          :key="item.name"
          class="model-item flex align-center pa-5"
        >
          <div class="flex-fill mx-5">
            <div class="font-weight-bold">
              {{ item.name }}
            </div>
            <div class="desc">
              {{ currentLang == "zh-cn" ? item.zhdesc : item.endesc }}
            </div>
            <div></div>
          </div>
          <div class="drawer-model-actions">
            <el-popover placement="left" :width="300" trigger="click">
              <template #reference>
                <el-button icon="Download" circle />
              </template>
              <template #default>
                <div
                  v-for="(el, index) in item.models"
                  :key="index"
                  :value="el.model"
                  @click="downLabel(el, item)"
                  class="list-column"
                >
                  <div class="list-column-title">
                    {{ el.model }}
                    <el-tag size="small" type="info">{{
                      getModelStatus(el.model)
                    }}</el-tag>
                  </div>
                  <div class="list-footer">
                    <el-tag size="small" type="primary">size:{{ el.info.size }}</el-tag>
                    <el-tag size="small" type="success">cpu:{{ el.info.cpu }}</el-tag>
                    <el-tag size="small" type="danger">gpu:{{ el.info.gpu }}</el-tag>
                  </div>
                </div>
              </template>
            </el-popover>
            <el-button icon="Edit" circle @click="labelShow(item.id)" />
            <el-button
              @click.stop="delLabel(item.id)"
              icon="Delete"
              v-if="item.models.length === 0"
              circle
            />
          </div>
        </div>
      </div>
    </el-scrollbar>
    <el-tour v-model="modelStore.openHelp">
      <el-tour-step
        :target="help_label?.$el"
        :title="t('model.help_label')"
        :description="t('model.help_labelDesc')"
      />
      <el-tour-step
        :target="help_showdown?.$el"
        :title="t('model.help_showdown')"
        :description="t('model.help_showdownDesc')"
      />
      <el-tour-step
        :target="help_adddown?.$el"
        :title="t('model.help_adddown')"
        :description="t('model.help_adddownDesc')"
      />
      <el-tour-step
        :target="help_addlabel?.$el"
        :title="t('model.help_addlabel')"
        :description="t('model.help_addlabelDesc')"
      />
    </el-tour>
  </div>
</template>

<style scoped lang="scss">
@import "@/assets/list.scss";
@import "@/assets/left.scss";
</style>
