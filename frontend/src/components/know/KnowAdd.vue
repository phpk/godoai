<script lang="ts" setup>
import { onMounted, ref, toRaw } from "vue";
import { useKnowStore } from "@/stores/know";
import { t } from "@/i18n/index";
import { notifySuccess, notifyError, buyVip } from "@/util/msg";
import { getSystemConfig } from "@/stores/config";
const knowStore: any = useKnowStore();

const chatModelList:any = ref([]);
const embedModelList:any = ref([]);

onMounted(async () => {
  chatModelList.value = await knowStore.getModelList("chat");
  embedModelList.value = await knowStore.getModelList("embeddings");
});

async function saveData() {
  const save: any = toRaw(knowStore.current);
  if (!save.title || save.title == "") {
    notifyError(t('knowledge.pleaseInputKnowledgeBaseName'));
    return;
  }
  if (!save.embedmodel || save.embedmodel == "") {
    notifyError(t('knowledge.pleaseSelectEmbeddingModel'));
    return;
  }
  //save.num_ctx = modelData.params.num_ctx;
  if (save.contextLength < 50 || save.contextLength > getContextLength()) {
    notifyError(t('knowledge.pleaseInputContextLength'));
    return;
  }
  if (!save.chatmodel || save.chatmodel == "") {
    notifyError(t('knowledge.pleaseSelectChatModel'));
    return;
  }
  
  // 检测apiurl是否能fetch正常
  if (knowStore.current.type == "chroma") {
    buyVip()
    return;
  }

  if (knowStore.isEdit) {
    //save.id = knowStore.current.id;
    await knowStore.updateKnow(save);
    notifySuccess(t('common.saveSuccess'));
  } else {
    const has = await knowStore.hasKnow(save.title);
    if(has){
      notifyError(t('knowledge.knowledgeBaseNameExists'));
      return;
    }
    const config = getSystemConfig()
    const postData = {
      name: save.title,
      model: save.embedmodel,
      config: {
        type: save.type,
        apiUrl: save.apiUrl,
        embedding:{
          apiUrl: config.embedApi,
          apiType: config.embedType,
          apiKey:config.embedApiKey,
          contextLength : save.contextLength
        }
      },
    };
    const completion = await fetch(config.knowledgeUrl + "/create", {
      method: "POST",
      body: JSON.stringify(postData),
    });
    if (!completion.ok) {
      notifyError(completion.statusText);
      return;
    }
    const res = await completion.json();
    if (res.code !== 0 || !res.data || !res.data.Id) {
      notifyError(res.message);
      return;
    }
    save.uuid = res.data.Id;
    console.log(res);

    await knowStore.addKnow(save);
    notifySuccess(t('common.saveSuccess'));
  }

  knowStore.isAdd = false;
}
function setUrl(url: string) {
  if (knowStore.current.type != "chromem") {
    const list = toRaw(knowStore.configList);
    //console.log(list)
    knowStore.current.apiUrl = list.find((item:any) => item.type == url).apiUrl;
  }
}
function changeUrl(_: any) {
  knowStore.changeConfig(knowStore.current.type, knowStore.current.apiUrl);
}
function getContextLength(){
  const modelData:any = embedModelList.value.find((d:any) => d.model == knowStore.current.embedmodel)
  if(!modelData.info || !modelData.info.context_length){
    notifyError(t('knowledge.pleaseInputContextLength'));
    return;
  }
  return modelData.info.context_length;
}
function changeLength(){
  knowStore.current.contextLength = getContextLength();
}
</script>

<template>
  <el-form label-width="180px" style="margin-top: 12px">
    <el-form-item :label="t('knowledge.knowledgeBaseName')">
      <el-input
        v-model="knowStore.current.title"
        placeholder=""
        prefix-icon="Compass"
        :readonly="knowStore.isEdit"
        clearable
      ></el-input>
    </el-form-item>
    <el-form-item :label="t('knowledge.embeddingModel')" :placeholder="t('knowledge.pleaseSelectEmbeddingModel')">
      <el-select 
      v-model="knowStore.current.embedmodel" 
      :disabled="knowStore.isEdit"
      @change="changeLength"
      >
        <el-option
          v-for="(item, key) in embedModelList"
          :key="key"
          :label="item.model"
          :value="item.model"
        />
      </el-select>
    </el-form-item>
    <el-form-item :label="t('knowledge.chatModel')">
      <el-select v-model="knowStore.current.chatmodel" :placeholder="t('knowledge.pleaseSelectChatModel')">
        <el-option
          v-for="(item, key) in chatModelList"
          :key="key"
          :label="item.model"
          :value="item.model"
        />
      </el-select>
    </el-form-item>

    <el-form-item :label="t('knowledge.databaseEngine')" prop="type">
      <el-select
        v-model="knowStore.current.type"
        :placeholder="t('knowledge.databaseEngine')"
        :disabled="knowStore.isEdit"
      >
        <el-option
          v-for="item in knowStore.configList"
          :key="item.type"
          :label="item.type"
          :value="item.type"
          @change="setUrl"
        />
      </el-select>
    </el-form-item>
    <el-form-item v-if="knowStore.current.type !== 'chromem'" :label="t('knowledge.engineAddress')">
      <el-input
        v-model="knowStore.current.apiUrl"
        placeholder=""
        prefix-icon="Connection"
        clearable
        @blur="changeUrl"
        :readonly="knowStore.isEdit"
      ></el-input>
    </el-form-item>
    <el-form-item :label="t('knowledge.sliceLength')">
      <el-input
        v-model="knowStore.current.contextLength"
        placeholder=""
        prefix-icon="Compass"
        clearable
      ></el-input>
    </el-form-item>
    <el-form-item>
      <el-button type="primary" icon="CirclePlus" @click="saveData"
        >{{ t('common.save') }}</el-button
      >
    </el-form-item>
  </el-form>
  
</template>
