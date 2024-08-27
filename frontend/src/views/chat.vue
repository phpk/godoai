<script setup lang="ts">
import { t } from "@/i18n/index";
import { useChatStore } from "@/stores/chat";
import {
  getCurrentModelName,
  getSystemKey,
  getChatConfig,
  getSystemConfig
} from "@/stores/config";

import { useKnowStore } from "@/stores/know";
import { computed, nextTick, onMounted, ref, toRaw, watch } from "vue";
import { Vue3Lottie } from "vue3-lottie";
import { notifySuccess, notifyInfo, notifyError, buyVip } from "@/util/msg";
import { ElScrollbar } from "element-plus";
//import { startApp, startllm, llmType } from "@/stores/app";
import { useRouter } from 'vue-router'
const router = useRouter();
const chatStore = useChatStore();
const knowStore = useKnowStore();

const isPadding = ref(false); //是否发送中
const drawerChatBox = ref(); //侧栏
const editInfo: any = ref({}); //编辑聊天信息

const chatDialogBox = ref(false);
const isEditor = ref(true);
const imageInput: any = ref(null);
let imageData = ref("");
const messageContainerRef = ref<InstanceType<typeof ElScrollbar>>();
const messageInnerRef = ref<HTMLDivElement>();
// User Input Message
const userMessage = ref("");
const konwAll: any = ref([]);
const llmType = getSystemKey("llmType")
let hasCurrent:any;
// Prompt Message
const promptMessage = computed(() => {
  return [
    {
      content: chatStore.chatInfo.prompt,
      chatType: "chat",
      chatId: chatStore.activeId,
      role: "system",
      id: Date.now(),
      createAt: Date.now(),
    },
  ];
});

onMounted(async () => {
  hasCurrent = getCurrentModelName("chat");
  if (!hasCurrent) {
    notifyError(t("index.notFindChatModel"));
    return;
  }
  await chatStore.initChat();
  await chatStore.getActiveChat();
  konwAll.value = await knowStore.getKnowAll();
  scrollToBottom();
});
const requestMessages = computed(() => {
  const chatConfig = getChatConfig('chat');
  const contextLen = chatConfig.contextLength;
  //console.log(contextLen)
  if (chatStore.messageList.length <= contextLen) {
    return [...promptMessage.value, ...chatStore.messageList];
  } else {
    // 截取最新的10条信息
    const slicedMessages = chatStore.messageList.slice(-contextLen);
    return [...promptMessage.value, ...slicedMessages];
  }
});

// Send Messsage
const sendMessage = async () => {
  if (!hasCurrent) {
    notifyError(t("index.notFindChatModel"));
    return;
  }
  if (userMessage.value) {
    // Add the message to the list
    if (isPadding.value === true) return;
    const saveMessage = {
      content: userMessage.value,
      chatType: "chat",
      chatId: chatStore.activeId,
      role: "user",
      id: Date.now(),
      createdAt: Date.now(),
    };
    chatStore.messageList.push(saveMessage);
    await chatStore.addMessages(chatStore.activeId, saveMessage);

    // Clear the input
    userMessage.value = "";

    // Create a completion
    isPadding.value = true;
    await createCompletion();
  }
};

const createCompletion = async () => {
  try {
    const config = getSystemConfig()
    const messageId = Date.now();
    const saveMessage: any = {
      content: "",
      role: "assistant",
      chatType: "chat",
      chatId: chatStore.activeId,
      id: messageId,
      createdAt: messageId,
    };
    if (chatStore.chatInfo.kid && chatStore.chatInfo.kid != "") {
      const usermsg = chatStore.messageList[chatStore.messageList.length - 1];
      const promptData = await knowStore.askKnow(
        usermsg.content,
        chatStore.chatInfo.kid
      );
      if (promptData) {
        chatStore.chatInfo.prompt = promptData.prompt;
        saveMessage.link = promptData.link;
      }
    } 
    const chatConfig = getChatConfig('chat');
    let postMsg: any = {
      messages: requestMessages.value,
      model: chatStore.chatInfo.model,
      stream: false,
      type: llmType,
      options: chatConfig,
    };
    if (imageData.value != "") {
      const img2txtModel = await getCurrentModelName("img2txt");
      const usermsg = chatStore.messageList[chatStore.messageList.length - 1];
      postMsg = {
        model: img2txtModel,
        //"prompt":userMessage.value,
        stream: false,
        type: llmType,
        //options: chatConfig,
        messages: [
          {
            role: "user",
            content: usermsg.content,
            images: [imageData.value],
          },
        ],
      };
    }
    const postData: any = {
      method: "POST",
      body: JSON.stringify(postMsg),
    };
    
    const completion = await fetch(config.chatApi, postData);
    imageData.value = "";
    if (!completion.ok) {
      const errorData = await completion.json();
      notifyError(errorData.error.message);
      isPadding.value = false;
      return;
    }
    const res = await completion.json();
    //console.log(res)
    if(res && res.choices && res.choices.length > 0){
      if(res.choices[0].message.content){
        const msg = res.choices[0].message.content;
        saveMessage.content = msg;
        chatStore.messageList.push(saveMessage);
        await chatStore.addMessages(chatStore.activeId, saveMessage);
      }
    }
    isPadding.value = false;
  } catch (error:any) {
    isPadding.value = false;
    notifyError(error.message);
  }
};
const scrollToBottom = () => {
  nextTick(() => {
    if (messageContainerRef && messageContainerRef.value) {
      messageContainerRef.value!.setScrollTop(
        messageInnerRef.value!.clientHeight
      );
    }
  });
};
watch(
  () => chatStore.messageList,
  () => {
    scrollToBottom();
  },
  {
    deep: true,
  }
);

const handleKeydown = (e:any) => {
  if (e.key === "Enter" && (e.altKey || e.shiftKey)) {
    // 当同时按下 alt或者shift 和 enter 时，插入一个换行符
    e.preventDefault();
    userMessage.value += "\n";
  } else if (e.key === "Enter") {
    // 当只按下 enter 时，发送消息
    e.preventDefault();
    sendMessage();
  }
};
const clearChat = () => {
  chatStore.clearChatHistory(chatStore.activeId);
  chatStore.messageList = [];
};
const changeRoom = async (id:any) => {
  chatStore.setActiveId(id);
};
const deleteChat = async (id:any) => {
  chatStore.activeId = await chatStore.deleteChat(id);
};
const showBox = (flag:any) => {
  isEditor.value = flag;
  if (flag === true) {
    //console.log(chatStore.chatInfo)
    editInfo.value = toRaw(chatStore.chatInfo);
    //console.log(editInfo.value);
  } else {
    editInfo.value = {
      title: "",
      model: "",
      prompt: "",
      promptName: "",
      kid: "",
    };
  }
  chatDialogBox.value = true;
};
const changeInfo = async (info: any) => {
  info = toRaw(info);

  if (!info.title) {
    notifyInfo(t("chat.inputTitle"));
    return;
  }
  if (!info.model) {
    notifyInfo(t("chat.selectModel"));
    return;
  }
  if (!info.prompt) {
    info.prompt = "";
  }
  if(info.kid && info.kid != ""){
    buyVip().then(()=>{
      router.push('/help')
    })
    return
  }
  delete info.id;

  if (isEditor.value) {
    console.log(info);
    await chatStore.updateChat(info, chatStore.activeId);
    chatStore.chatInfo = info;
    notifySuccess(t("chat.editsuccess"));
  } else {
    const promptData = {
      prompt: info.prompt,
      promptName: info.promptName,
    };
    await chatStore.addChat(info.title, info.model, promptData, info.kid);
    notifySuccess(t("chat.addsuccess"));
  }
  await chatStore.getActiveChat();
  chatDialogBox.value = false;
};

const selectImage = async () => {
  const img2txtModel = await getCurrentModelName("img2txt");
  if (!img2txtModel) {
    notifyError(t("chat.notEyeModel"));
    return;
  }

  imageInput.value.click();
};

const uploadImage = async (event: any) => {
  const file = event.target.files[0];
  if (!file) {
    return;
  }
  const reader = new FileReader();
  reader.onload = (e: any) => {
    imageData.value = e.target.result.split(",")[1];
  };

  reader.readAsDataURL(file);
};
</script>
<template>
  <el-dialog v-model="chatDialogBox" width="600" append-to-body>
    <chat-edit-info :dataInfo="editInfo" @saveFn="changeInfo" />
  </el-dialog>

  <el-drawer
    v-model="drawerChatBox"
    direction="ltr"
    style="height: 100vh"
    :show-close="false"
    :with-header="false"
  >
    <el-scrollbar>
      <el-space direction="vertical">
        <el-card
          v-for="(item, key) in chatStore.chatList"
          :key="key"
          class="box-card"
          style="width: 190px"
          :shadow="chatStore.activeId == item.id ? 'always' : 'hover'"
        >
          <el-row type="flex" justify="space-between">
            <el-button
              type="info"
              text
              @click.stop="changeRoom(item.id)"
              size="small"
            >
              {{ item.title }}
            </el-button>
            <el-button
              icon="Delete"
              size="small"
              @click.stop="deleteChat(item.id)"
              circle
            ></el-button>
          </el-row>
        </el-card>
      </el-space>
    </el-scrollbar>
  </el-drawer>
  <div class="chat-bot">
    <el-page-header icon="null">
      <template #title>
        <div></div>
      </template>
      <template #content>
        <el-space>
          <el-button
            @click.stop="drawerChatBox = !drawerChatBox"
            icon="Menu"
            circle
          />
          <el-button @click.stop="clearChat" icon="DeleteFilled" circle />
          <el-button @click.stop="showBox(true)" icon="Tools" circle />
          <el-button @click.stop="showBox(false)" icon="Plus" circle />
        </el-space>
      </template>
    </el-page-header>

    <div class="messsage-area">
      <el-scrollbar
        v-if="chatStore.messageList.length > 0"
        class="message-container"
        ref="messageContainerRef"
      >
        <div ref="messageInnerRef">
          <chat-message
            v-for="message in chatStore.messageList"
            :key="message.messageId"
            :content="message.content"
            :link="message.link"
            :role="message.role"
            :createdAt="message.createdAt"
          />
        </div>
      </el-scrollbar>
      <div class="no-message-container" v-else>
        <Vue3Lottie animationLink="/bot/chat.json" :height="420" :width="420" />
      </div>
    </div>
    <div class="input-area">
      <div class="input-panel d-flex align-end pa-1">
        <el-row :gutter="24" style="border-bottom: none;">
          <el-col :span="2">
            <el-button
              @click="selectImage"
              size="large"
              icon="Paperclip"
              circle
            />
            <input
              type="file"
              ref="imageInput"
              accept="image/*"
              style="display: none"
              @change="uploadImage"
            />
          </el-col>
          <el-col :span="19">
            <el-input
              v-model="userMessage"
              :placeholder="t('chat.askme')"
              size="large"
              clearable
              @keydown="handleKeydown"
              autofocus
            />
          </el-col>
          <el-col :span="2">
            <el-button
              v-if="!isPadding"
              @click="sendMessage"
              icon="Promotion"
              type="info"
              size="large"
              circle
            />

            <el-button
              type="primary"
              size="large"
              v-else
              loading-icon="Eleme"
              icon="Loading"
              loading
              circle
            />
          </el-col>
        </el-row>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
@import "@/assets/chat.scss";
</style>
