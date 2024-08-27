<script lang="ts" setup>
import { useCreationStore } from "@/stores/creation";
import { t } from "@/i18n/index";
import { notifyError, notifySuccess } from "@/util/msg";
import Vditor from "vditor";
import "vditor/dist/index.css";
import { getEditorOption, getEditorTitleOption } from "@/util/vditor";
import { onMounted, ref } from "vue";
import { asBlob } from "html-docx-js-extends";
import { saveAs } from "file-saver";
import { getSystemKey, getChatConfig } from "@/stores/config";
const creationStore = useCreationStore();

const help_title = ref();
const help_cate = ref();
const help_create = ref();
const help_save = ref();
const help_showList = ref();
const createSelected = async () => {
  if (creationStore.isLoading) {
    notifyError(t("creation.loading"));
    return;
  }
  if (creationStore.baseTitle === "") {
    notifyError(t("creation.emptyTitle"));
    return;
  }
  if (creationStore.baseEditer.getValue() === "") {
    notifyError(t("creation.emptyParagraph"));
    return;
  }
  const selected = creationStore.baseEditer.getSelection();
  //console.log(selected);
  if (!selected) {
    notifyError(t("creation.selectText"));
    return;
  }
  creationStore.isSelectLoading = true;
  const message = [
    {
      role: "system",
      content: creationStore.config.systemPrompt,
    },
    {
      role: "user",
      content: creationStore.getBuilderPrompt(creationStore.baseTitle, selected),
    },
  ];
  await postMsg(message, true);
  creationStore.isSelectLoading = false;
};

const createRow = async () => {
  if (creationStore.baseTitle === "") {
    notifyError(t("creation.emptyTitle"));
    return;
  }
  const message = [
    {
      role: "system",
      content: creationStore.config.systemPrompt,
    },
    {
      role: "user",
      content: creationStore.getLeaderPrompt(
        creationStore.currentCate,
        creationStore.baseTitle
      ),
    },
  ];
  creationStore.isLoading = true;
  await postMsg(message, false);
  creationStore.isLoading = false;
};
const systemPrompt = {
  role: "system",
  content: creationStore.config.systemPrompt,
};
const contentMsg: any = [systemPrompt];

const createContent = async () => {
  if (creationStore.isLoading) {
    notifyError(t("creation.loading"));
    return;
  }
  if (creationStore.baseTitle == "") {
    notifyError(t("creation.emptyTitle"));
    return;
  }
  if (creationStore.baseEditer.getValue() == "") {
    notifyError(t("creation.emptyParagraph"));
    return;
  }

  const resp = creationStore.baseEditer
    .getValue()
    .replace(/\n- /g, "\n1.") //兼容不同格式
    .split("\n");
  if (
    resp[0].startsWith(`## ${creationStore.baseTitle}`) ||
    resp[0].startsWith(`**${creationStore.currentCate}` + t("creation.outline"))
  ) {
    resp.shift();
  }
  creationStore.isContentLoading = true;
  for (const i in resp) {
    let line: any = resp[i];
    if (line == "") continue;
    //s.push(line.trim());
    const message = {
      role: "user",
      content: creationStore.getBuilderPrompt(creationStore.baseTitle, line.trim()),
    };

    contentMsg.push(message);
    // console.log(contentMsg)
    if (contentMsg.length < 11) {
      await postMsg(contentMsg, true);
    } else {
      const msg = [systemPrompt, ...contentMsg.slice(-9)];
      await postMsg(msg, true);
    }
  }
  creationStore.isContentLoading = false;
};

const postMsg = async (message: any, isContent: any) => {
  if (creationStore.config.model == "") {
    notifyError(t("creation.noModel"));
    return;
  }
  const chatConfig = getChatConfig("creation");
  const completion = await fetch(getSystemKey("chatApi"), {
    method: "POST",
    body: JSON.stringify({
      messages: message,
      options: chatConfig,
      model: creationStore.config.model,
      stream: false,
    }),
  });
  if (!completion.ok || completion.status !== 200) {
    const errorData = await completion.json();
    notifyError(errorData.error.message);
    return;
  }
  const res = await completion.json();
  //console.log(res)
  if (res && res.choices && res.choices.length > 0) {
    if (res.choices[0].message.content) {
      let msg = res.choices[0].message.content;
      if (isContent) {
        creationStore.targetContent += msg;
        creationStore.targetEditor.setValue(creationStore.targetContent);
      } else {
        creationStore.baseContent += msg;
        creationStore.baseEditer.setValue(creationStore.baseContent);
      }
      if (contentMsg.length > 1) {
        msg = "---\n" + msg;
      }
      contentMsg.push({
        role: "assistant",
        content: msg,
      });
    }
  }
  creationStore.isLoading = false;
  creationStore.isContentLoading = false;
};
function setBaseVal(val: any) {
  creationStore.baseEditer.setValue(val);
  creationStore.baseContent = val;
}
function getBaseVal() {
  return creationStore.baseContent;
}
function getTargetVal() {
  return creationStore.targetContent;
}
function setTargetVal(val: any) {
  creationStore.targetEditor.setValue(val);
  creationStore.targetContent = val;
}
const download = async () => {
  if (creationStore.baseTitle === "") {
    notifyError(t("creation.emptyTitle"));
    return;
  }
  if (creationStore.targetEditor.getValue() === "") {
    notifyError(t("creation.emptyContent"));
    return;
  }
  asBlob(creationStore.targetEditor.getHTML()).then((data: any) => {
    saveAs(data, creationStore.baseTitle + ".docx"); // save as docx file
  });
};
const saveData = async () => {
  if (creationStore.baseTitle === "") {
    notifyError(t("creation.emptyTitle"));
    return;
  }
  if (creationStore.baseEditer.getValue() === "") {
    notifyError(t("creation.emptyParagraph"));
    return;
  }
  await creationStore.addArticle();
  notifySuccess(t("creation.saved"));
};
const newCreate = () => {
  creationStore.baseTitle = "";
  creationStore.currentCate = "";
  creationStore.activeId = 0;
  creationStore.showLeft = false;
};
onMounted(async () => {
  const editorOptionsTitle: any = getEditorTitleOption();
  const editorOptions: any = getEditorOption();
  editorOptionsTitle.input = (val: any) => {
    creationStore.baseContent = val;
  };
  editorOptions.input = (val: any) => {
    creationStore.targetContent = val;
  };
  creationStore.baseEditer = new Vditor("baseContent", editorOptionsTitle);
  creationStore.targetEditor = new Vditor("targetContent", editorOptions);
  await creationStore.initConfig();
});
</script>
<template>
  <div>
    <el-dialog v-model="creationStore.showConfig" width="600">
      <creation-setting />
    </el-dialog>
    <el-drawer
      v-model="creationStore.showLeft"
      direction="ltr"
      style="height: 100vh"
      :show-close="false"
      :with-header="false"
    >
      <creation-left
        @setBaseFn="setBaseVal"
        @setTargetFn="setTargetVal"
        @newFn="newCreate"
      />
    </el-drawer>
    <el-page-header icon="null">
      <template #title>
        <div></div>
      </template>
      <template #content>
        <el-space>
          <el-button
            @click.stop="creationStore.showLeft = true"
            icon="Menu"
            circle
            ref="help_showList"
          />
          <el-button @click.stop="creationStore.handlerArea()" icon="Switch" circle />
          <el-button @click.stop="creationStore.showConfig = true" icon="Tools" circle />
          <el-button @click.stop="newCreate" icon="Plus" circle />
        </el-space>
      </template>
      <template #extra>
        <el-space class="mr-10" ref="help_cate">
          <el-select-v2
            v-model="creationStore.currentCate"
            :options="creationStore.creationCate"
            :placeholder="t('common.cate')"
            style="width: 240px; vertical-align: middle"
            allow-create
            filterable
            clearable
          />
        </el-space>
      </template>
    </el-page-header>
    <el-row style="padding: 12px" justify="space-around">
      <el-col :span="creationStore.areaSet.left">
        <el-card>
          <template #header>
            <div class="card-header">
              <el-input
                v-model="creationStore.baseTitle"
                :placeholder="t('creation.titlePlaceholder')"
                ref="help_title"
              >
              </el-input>
            </div>
          </template>
          <div id="baseContent"></div>
          <template #footer>
            <el-row justify="end" style="border-bottom: none;margin:0;padding:0">
              <el-button icon="Delete" @click="setBaseVal('')" size="small" circle />
              <copy-btn :text="getBaseVal()" />
            </el-row>
          </template>
        </el-card>
      </el-col>
      <el-col :span="creationStore.areaSet.right">
        <el-card>
          <template #header>
            <div class="card-header">
              <el-row type="flex" justify="space-between" style="border-bottom: none;padding-bottom: 0;margin:0;">
                <div>
                  <el-tooltip
                    class="box-item"
                    effect="dark"
                    :content="t('creation.generateParagraph')"
                    placement="top-start"
                  >
                    <el-button
                      :loading="creationStore.isLoading"
                      :disabled="creationStore.isLoading"
                      type="primary"
                      circle
                      icon="Memo"
                      ref="help_create"
                      @click="createRow"
                    />
                  </el-tooltip>
                  <el-tooltip
                    class="box-item"
                    effect="dark"
                    :content="t('creation.generateSelected')"
                    placement="top-start"
                  >
                    <el-button
                      :loading="creationStore.isSelectLoading"
                      :disabled="creationStore.isSelectLoading"
                      @click="createSelected"
                      type="success"
                      circle
                      icon="Tickets"
                    />
                  </el-tooltip>

                  <el-tooltip
                    class="box-item"
                    effect="dark"
                    :content="t('creation.generateAll')"
                    placement="top-start"
                  >
                    <el-button
                      :loading="creationStore.isContentLoading"
                      :disabled="creationStore.isContentLoading"
                      type="warning"
                      circle
                      icon="Collection"
                      @click="createContent"
                    />
                  </el-tooltip>
                </div>
                <el-tooltip
                  class="box-item"
                  effect="dark"
                  :content="t('creation.save')"
                  placement="top-start"
                >
                  <el-button
                    type="danger"
                    circle
                    icon="Finished"
                    ref="help_save"
                    @click="saveData"
                  />
                </el-tooltip>
              </el-row>
            </div>
          </template>
          <div id="targetContent"></div>
          <template #footer>
            <el-row justify="end" style="border-bottom: none;margin:0;padding:0">
              <el-button icon="Delete" @click="setTargetVal('')" size="small" circle />
              <copy-btn :text="getTargetVal()" />
              <el-button icon="Download" @click="download()" size="small" circle />
            </el-row>
          </template>
        </el-card>
      </el-col>
    </el-row>
    <el-tour v-model="creationStore.openHelp">
      <el-tour-step
        :target="help_title?.$el"
        :title="t('creation.titleInput')"
        :description="t('creation.descriptionInputTitle')"
      />
      <el-tour-step
        :target="help_cate?.$el"
        :title="t('creation.chooseCategory')"
        :description="t('creation.descriptionChooseCategory')"
      />
      <el-tour-step
        :target="help_create?.$el"
        :title="t('creation.operationButtons')"
        :description="t('creation.descriptionOperationButtons')"
      />
      <el-tour-step
        :target="help_save?.$el"
        :title="t('creation.saveData')"
        :description="t('creation.descriptionSaveData')"
      />
      <el-tour-step
        :target="help_showList?.$el"
        :title="t('creation.viewList')"
        :description="t('creation.descriptionViewList')"
      />
    </el-tour>
  </div>
</template>
<style lang="scss">
#targetContent,
#baseContent {
  text-align: left;
}
</style>
