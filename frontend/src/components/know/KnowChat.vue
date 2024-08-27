<script setup lang="ts">
import { Vue3Lottie } from "vue3-lottie";
import { t } from "@/i18n/index";
import { ElScrollbar } from "element-plus";
import { notifyError } from "@/util/msg";
import { useKnowStore } from "@/stores/know";
import { getSystemKey,getChatConfig } from "@/stores/config";
import { computed, nextTick, ref, watch } from "vue";

const knowStore: any = useKnowStore();
const messageContainerRef = ref<InstanceType<typeof ElScrollbar>>();
const messageInnerRef = ref<HTMLDivElement>();

// User Input Message
const userMessage = ref("");
const prompt = ref(t("knowledge.defPrompt"));

// Prompt Message
const promptMessage = computed(() => {
  return [
    {
      content: prompt.value,
      role: "system",
    },
  ];
});

// Message List
const messages:any = ref([]);

const requestMessages = computed(() => {
  //const contextLen = systemStore.modelConfig.chat.contextLength
  const chatConfig = getChatConfig('knowledge');
  const contextLen = chatConfig.contextLength;
  if (messages.value.length <= contextLen) {
    return [...promptMessage.value, ...messages.value];
  } else {
    // 截取最新的10条信息
    const slicedMessages = messages.value.slice(-contextLen);
    return [...promptMessage.value, ...slicedMessages];
  }
});

const isLoading = ref(false);

// Send Messsage
const sendMessage = async () => {
  if (isLoading.value) return;
  if (userMessage.value) {
    // Add the message to the list
    messages.value.push({
      content: userMessage.value,
      role: "user"
    });

    // Clear the input
    userMessage.value = "";

    // Create a completion
    await createCompletion();
  }
};

const createCompletion = async () => {
  isLoading.value = true;
  try {
    const usermsg = messages.value[messages.value.length - 1];
    //console.log(usermsg)
    const promptData = await knowStore.askKnow(usermsg.content, knowStore.current.uuid)
    const robMessage: any = {
      content: "",
      role: "assistant",
    };
    
    if(promptData) {
      prompt.value = promptData.prompt;
      robMessage.link = promptData.link;
    }
    const chatConfig = getChatConfig('knowledge');
    const completion = await fetch(
      getSystemKey("chatApi"),
      {
        method: "POST",
        body: JSON.stringify({
          messages: requestMessages.value,
          model: knowStore.current.chatmodel,
          stream: false,
          options: chatConfig,
        }),
      }
    );

    // Handle errors
    if (!completion.ok) {
      const errorData = await completion.json();
      notifyError(errorData.error.message);
      return;
    }
    const res = await completion.json();
    //console.log(res)
    if(res && res.choices && res.choices.length > 0){
      if(res.choices[0].message.content){
        const msg = res.choices[0].message.content;
        robMessage.content = msg;
        messages.value.push(robMessage);
      }
    }
    isLoading.value = false;
  } catch (error:any) {
    isLoading.value = false;
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
  () => messages.value,
  (_) => {
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


</script>

<template>
  <div class="chat-bot" style="height: 92vh">
  <div class="messsage-area">
    <el-scrollbar
      v-if="messages.length > 0"
      class="message-container"
      ref="messageContainerRef"
    >
      <div ref="messageInnerRef">
        <chat-message
          v-for="(message, key) in messages"
          :key="key"
          :content="message.content"
          :link="message.link"
          :role="message.role"
          :createdAt="Date.now()"
        />
      </div>
    </el-scrollbar>
    <div class="no-message-container" v-else>
      <Vue3Lottie animationLink="/bot/knowchat.json" :height="300" :width="300" />
    </div>
  </div>
  <div class="input-area" style="margin-top:10px">
      <div class="input-panel d-flex align-end pa-1">
        <el-row :gutter="24">
          <el-col :span="21">
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
              v-if="!isLoading"
              @click="sendMessage"
              icon="Promotion"
              type="danger"
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