<script lang="ts" setup>
import { useKnowStore } from "@/stores/know";
import { getSystemKey,getSystemConfig } from "@/stores/config";
import { usePageStore } from "@/stores/page.ts";
import { Vue3Lottie } from "vue3-lottie";
import { onMounted, ref, toRaw } from "vue";
import { notifySuccess, notifyError } from "@/util/msg";
import { ElLoading } from "element-plus";
import { t } from "@/i18n/index";

const fileInput: any = ref(null);
const knowStore:any = useKnowStore();
const pageStore = usePageStore();

function selectFile() {
  fileInput.value.click();
}
onMounted(async () => {
  knowStore.getKnowList();
});
const showFile = async (data: any) => {
  knowStore.isShowFile = !knowStore.isShowFile;
  knowStore.current = toRaw(data);
  await knowStore.getFileList(knowStore.current.id);
};
const uploadFile = async (event: any) => {
  const files = event.target.files;
  if (!files || files.length < 1) {
    return;
  }

  const formData: any = new FormData();
  const names: any = knowStore.getFileNames();
  const rehas: any = [];
  for (let i = 0; i < files.length; i++) {
    const filename: any = files[i].name;
    if (!names.includes(filename)) {
      formData.append("files", files[i]);
      names.push(filename);
    } else {
      rehas.push(filename);
    }
  }
  if (rehas.length > 0) {
    notifyError(t('knowledge.duplicateFilenames') + rehas.join(","));
    //return
  }
  if (formData.length < 1) {
    notifyError(t('knowledge.noFilesUploaded'));
    return;
  }
  const loadingInstance = ElLoading.service({
    lock: true,
    text: "UpLoading",
  });
  const response = await fetch(getSystemKey('knowledgeUrl') + "/upload", {
    method: "POST",
    body: formData,
  });

  loadingInstance.close();
  //console.log(json)
  if (!response.ok) {
    notifyError(response.statusText);
    return;
  }
  const json = await response.json();
  if (json.code !== 0) {
    notifyError(json.msg);
    return;
  }
  const addData = json.data;
  if (addData.length < 1) {
    notifyError(t('knowledge.uploadFailed'));
    return;
  }
  const saveFiles: any = [];
  addData.forEach((d: any) => {
    d.kid = knowStore.current.id;
    d.status = 0;
    saveFiles.push(d.save_path);
  });
  //console.log(addData)
  //return
  await knowStore.addFile(addData);
  await addUploadData(saveFiles);
  notifySuccess(t('knowledge.uploadSuccess'));
};
async function addUploadData(saveFiles:any) {
  if (saveFiles.length > 0) {
    const config = getSystemConfig();
    const postConv = {
      name: knowStore.current.uuid,
      model: knowStore.current.embedmodel,
      files: saveFiles,
      config: {
        type: knowStore.current.type,
        apiUrl: knowStore.current.apiUrl,
        embedding:{
          apiUrl: config.embedApi,
          apiType: config.embedType,
          apiKey:config.embedApiKey,
          contextLength : knowStore.current.contextLength,
        }
      },
    };
    await saveData(postConv);
  } else {
    notifyError(t('knowledge.uploadFailed'));
  }
}
async function saveData(postConv:any) {
  const loadingInstance = ElLoading.service({
    lock: true,
    text: "Transfering",
  });
  const response = await fetch(getSystemKey('knowledgeUrl') + "/add", {
    method: "POST",
    body: JSON.stringify(postConv),
  });
  loadingInstance.close();
  if (!response.ok) {
    notifyError(response.statusText);
    return;
  }
  //console.log(response);
  const json = await response.json();
  //console.log(json);
  if (json.code == 0) {
    await knowStore.updateFiles(postConv.files,knowStore.current.id);
    notifySuccess(t('knowledge.transferSuccess'));
  } else {
    notifyError(json.message);
  }
}
async function delKnow(item:any) {
  const data = toRaw(item);
  await knowStore.delKonw(data.id);
  const config = getSystemConfig();
  const postData = {
    name: data.uuid,
    config: {
      type: knowStore.current.type,
      apiUrl: knowStore.current.apiUrl,
      embedding:{
          apiUrl: config.embedApi,
          apiType: config.embedType,
          apiKey:config.embedApiKey
        }
    },
  };
  const response = await fetch(config.knowledgeUrl + "/delete", {
    method: "POST",
    body: JSON.stringify(postData),
  });
  const json = await response.json();
  if (!response.ok) {
    notifyError(json.error.message);
    return;
  }
  notifySuccess(t('knowledge.deleteSuccess'));
}
async function deleteFile(item:any) {
  const fileData = toRaw(item);
  const knowData = await knowStore.getKnow(fileData.kid);
  await knowStore.delFile(fileData.id, knowData.id);
  const config = getSystemConfig();
  const postData = {
    name: knowData.uuid,
    model: knowData.embedmodel,
    file: fileData.save_path,
    config: {
      type: knowStore.current.type,
      apiUrl: knowStore.current.apiUrl,
      embedding:{
          apiUrl: config.embedApi,
          apiType: config.embedType,
          apiKey:config.embedApiKey
        }
    },
  };
  const response = await fetch(config.knowledgeUrl + "/deleteFile", {
    method: "POST",
    body: JSON.stringify(postData),
  });
  const json = await response.json();
  //console.log(json)
  if (!response.ok) {
    notifyError(json.error.message);
    return;
  }
  notifySuccess(t('knowledge.deleteSuccess'));
}
function showChat(item:any) {
  knowStore.current = toRaw(item);
  knowStore.dialog = true;
}
function showEdit(item:any) {
  knowStore.current = toRaw(item);
  knowStore.isAdd = true;
  knowStore.isEdit = true;
}
function showAdd() {
  knowStore.isAdd = true;
  knowStore.isEdit = false;
  knowStore.current = {
    title: "",
    embedmodel: "",
    chatmodel: "",
    type: "chromem",
    apiUrl: getSystemKey('chromaUrl'),
  };
}
async function handleUrlSaved(url: string) {
  var regex =
    /^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/)?[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)?$/;
  if (!regex.test(url)) {
    notifyError(t('knowledge.invalidUrl'));
    return;
  }
  knowStore.knowLink = false;
  const loadingInstance = ElLoading.service({
    lock: true,
    text: "UpLoading",
  });
  const response = await fetch(getSystemKey('knowledgeUrl') + "/url", {
    method: "POST",
    body: JSON.stringify({ url }),
  });
  loadingInstance.close();
  //console.log(json)
  if (!response.ok) {
    notifyError(response.statusText);
    return;
  }
  const addData = await response.json();
  notifySuccess(t('knowledge.uploadSuccess'));

  //return
  if (!addData) {
    notifyError(t('knowledge.uploadFailed'));
    return;
  }
  addData.kid = knowStore.current.id;
  const saveFiles: any = [addData.save_path];
  await knowStore.addFile([addData]);
  await addUploadData(saveFiles);
}
</script>

<template>
  <div>
  <el-dialog v-model="knowStore.knowLink" width="600" :title="t('knowledge.addUrl')">
    <know-link @url-saved="handleUrlSaved" />
  </el-dialog>
  <el-dialog v-model="knowStore.isAdd" width="600" :title="t('knowledge.mangerKnowledge')">
    <know-add />
  </el-dialog>
  <el-dialog
    v-model="knowStore.dialog"
    :fullscreen="true"
    :title="knowStore.current.title"
  >
    <know-chat />
  </el-dialog>
  <el-drawer
    v-model="knowStore.isShowFile"
    direction="ltr"
    style="height: 100vh"
    width="300px"
    :show-close="false"
    :with-header="false"
  >
    <el-scrollbar>
      <el-row type="flex" justify="end" style="margin-bottom: 15px">
        <input
          type="file"
          ref="fileInput"
          accept="image/*,.doc,.docx,.xls,.xlsx,.csv,.pdf,.txt,.md,.ppt,.pptx,video/*,audio/*"
          style="display: none"
          multiple
          @change="uploadFile"
        />
        <el-button icon="Upload" circle @click.stop="selectFile"></el-button>
        <el-button
          icon="Link"
          circle
          @click.stop="knowStore.knowLink = true"
        ></el-button>
      </el-row>
      <div
        class="list-item"
        v-for="(item, index) in knowStore.fileList"
        :key="index"
      >
        <div class="list-title">
          <el-tooltip
            class="box-item"
            effect="dark"
            :content="item.name"
            placement="top-start"
            >{{ item.name }}</el-tooltip
          >
        </div>
        <el-button
          class="delete-btn"
          icon="VideoPlay"
          size="small"
          v-if="item.status < 1"
          @click.stop="addUploadData([toRaw(item.save_path)])"
          circle
        ></el-button>
        <el-button
          class="delete-btn"
          icon="Delete"
          size="small"
          @click.stop="deleteFile(item)"
          circle
        ></el-button>
      </div>
    </el-scrollbar>
  </el-drawer>
  <el-page-header icon="null">
    <template #title>
      <div></div>
    </template>
    <template #content>
      <div class="flex items-center">
        <el-button @click.stop="showAdd" icon="Plus" circle />
      </div>
    </template>
  </el-page-header>

  <el-scrollbar v-if="knowStore.knowList && knowStore.knowList.length > 0" style="height: 500px;">
    <div style="padding: 5px 15px 50px 15px">
      <el-card
        class="model-item"
        v-for="item in knowStore.knowList"
        :key="item.id"
        shadow="hover"
        style="margin: 8px 0"
      >
        <div class="card-header">
          <el-row type="flex" justify="space-between">
            <span>{{ item.title }}</span>
            <div class="card-header-right">
              <el-button icon="ChatLineRound" circle @click="showChat(item)" />
              <el-button icon="Link" circle @click="showFile(item)" />
              <el-button icon="Edit" circle @click="showEdit(item)" />
              <el-button icon="Delete" circle @click="delKnow(item)" />
            </div>
          </el-row>
        </div>
        <div style="display: flex; align-items: center; margin-top: 5px;">
          <el-tag size="small" style="margin-right: 5px;">{{ item.type }}</el-tag>
          <el-tag size="small" style="margin-right: 5px;">{{ item.embedmodel }}</el-tag>
          <el-tag size="small">{{ item.chatmodel }}</el-tag>
        </div>
      </el-card>
      <el-row
        justify="center"
        style="margin-top: 15px"
        v-if="pageStore.page.pages > 1"
      >
        <el-pagination
          background
          layout="prev, pager, next"
          v-model:current-page="pageStore.page.current"
          v-model:page-size="pageStore.page.size"
          :total="pageStore.page.total"
          @current-change="(val:any) => pageStore.pageClick(val, 'knowledge')"
        />
      </el-row>
    </div>
  </el-scrollbar>
  <div class="no-message-container" v-else>
    <Vue3Lottie animationLink="/bot/know.json" :height="420" :width="420" />
  </div>
</div>
</template>
<style scoped>
@import "@/assets/left.scss";
</style>
