<script setup lang="ts">
import {
  nextTick,
  onBeforeUnmount,
  onMounted,
  reactive,
  ref,
  watch,
} from "vue";
import { Vue3Lottie } from "vue3-lottie";
import { notifyError } from "@/util/msg";
import { getWaveBlob } from "@/util/audio.ts";
import { ElScrollbar } from "element-plus";
import { t } from "@/i18n/index";
import {
  getSystemKey,
  getSystemConfig,
  getChatConfig,
  getCurrentModelName,
  getModel,
} from "@/stores/config.ts";
//const pageLoading = ref(true);
const messageContainerRef = ref<InstanceType<typeof ElScrollbar>>();
const messageInnerRef = ref<HTMLDivElement>();


interface Message {
  content: string;
  role: "user" | "assistant" | "system";
}
// const apiUrl = systemStore.getVoiceUrl();
const currentChat = ref("");
const currentAudio = ref("");
const currentTts = ref("");

// Message List
const messages = ref<Message[]>([
  {
    content:t('spoken.systemWelcome'),
    role: "system",
  },
]);

// User Input Message
const userMessage = ref("");
const ifPlaying = ref(true);
const isLoading = ref(false);
let rec:any;
let audioElement:any;
const sexGender = ref(false);

async function initSpeechRecognition() {
  if (!("SpeechRecognition" in window) && !("webkitSpeechRecognition" in window)) {
    //notifyError("当前浏览器不支持语音识别");
    return null;
  }
  rec =
    "SpeechRecognition" in window
      ? new (window as any).SpeechRecognition()
      : new (window as any).webkitSpeechRecognition();

  rec.onstart = () => {
    console.log("开始监听语音");
  };

  rec.onspeechstart = () => {
    if(state.isResponse || state.isPlaying || state.isRecording)return;
    startRecord()
  };

  rec.onspeechend = () => {
    stopRecord()
  };
  rec.onresult = async () => {
  };
  rec.onerror = (event:any) => {
    console.error("语音识别错误", event.error);
    
  };

  rec.onend = () => {
    rec.start();
  };

  rec.start(); // 开始监听
  //return rec;
}
async function speechBackend(msg: string, apiUrl: string) {
  if(!ifPlaying.value)return;
  if(!currentTts.value){
    notifyError(t('spoken.downloadTTSFirst'));
    return;
  }
  const modelData = await getModel(currentTts.value);
  //console.log(modelData);
  const ttsData = {
    text: msg,
    sid: sexGender.value ? 167 : 99,
    model: currentTts.value,
    params: modelData.params,
  };
  state.isPlaying = true;
  try {
    const ttsRes = await fetch(apiUrl + "/tts", {
      method: "POST",
      body: JSON.stringify(ttsData),
    });

    if (!ttsRes.ok) {
      state.isPlaying = false;
      console.log(ttsRes.statusText)
      return
    }

    const ttsJson = await ttsRes.json();
    //console.log(ttsJson);
    if(ttsJson.data && ttsJson.data.txt != ""){
      playBase64Audio(ttsJson.data.txt);
    }else{
      state.isPlaying = false;
    }
 
  } catch (error) {
    state.isPlaying = false;
    console.log(error)
  }
  

}
function playBase64Audio(base64String: string) {
  // 将Base64字符串转换为byte数组
  const byteCharacters = atob(base64String);
  // 计算byte数组的长度
  const byteArrays:any = [];
  for (let offset = 0; offset < byteCharacters.length; offset += 1024) {
    const slice = byteCharacters.slice(offset, offset + 1024);
    const byteNumbers = new Array(slice.length);
    for (let i = 0; i < slice.length; i++) {
      byteNumbers[i] = slice.charCodeAt(i);
    }
    const byteArray = new Uint8Array(byteNumbers);
    byteArrays.push(byteArray);
  }

  // 将byte数组转换为Blob
  const blob = new Blob(byteArrays, { type: 'audio/wav' }); // 根据你的音频类型调整MIME类型，例如 'audio/mpeg' 对于MP3

  // 创建一个URL表示这个Blob
  const url = URL.createObjectURL(blob);

  // 创建或获取audio元素
  audioElement = document.getElementById('myAudio');
  if (!audioElement) {
    audioElement = document.createElement('audio');
    audioElement.id = 'myAudio';
    document.body.appendChild(audioElement); // 可能需要根据实际情况调整插入位置
  }

  // 设置音频的src为新创建的URL
  audioElement.src = url;

  // 添加事件监听器以确保音频加载完成后播放
  audioElement.addEventListener('canplaythrough', function () {
    audioElement.play();
  });

  // 添加事件监听器以侦听音频播放完成
  audioElement.addEventListener('ended', function () {
    //console.log('音频播放完成');
    state.isPlaying = false;
    // 在这里可以添加播放完成后的逻辑，例如调用其他函数
  });
}


// Send Messsage
const sendMessage = async () => {
  if (userMessage.value) {
    // Add the message to the list
    messages.value.push({
      content: userMessage.value,
      role: "user",
    });

    state.isRecording = false;
    state.isResponse = true;
    // Create a completion
    await createCompletion();

    // Clear the input
    userMessage.value = "";
  }
};
onMounted(async () => {
  currentChat.value = getCurrentModelName("chat");
  currentAudio.value = getCurrentModelName("audio");
  currentTts.value = getCurrentModelName("tts");
  if (!currentChat.value) {
    notifyError(t('spoken.downloadTextFirst'));
    return;
  }
  initSpeechRecognition();
  // await startApp("voice");
  // pageLoading.value = false;
});
onBeforeUnmount(() => {
  rec?.stop();
  rec = null;
});
const createCompletion = async () => {
  const config:any = getSystemConfig();
  const chatUrl = config.chatApi;
  const chatConfig = getChatConfig('translation');
  const postMsg: any = {
    messages: messages.value,
    model: currentChat.value,
    stream: false,
    options: chatConfig,
  };
  const postData: any = {
    method: "POST",
    body: JSON.stringify(postMsg),
  };
  const completion: any = await fetch(chatUrl, postData);
  //console.log(completion);
  isLoading.value = false;
  if (!completion.ok) {
    notifyError(completion.statusText);
    return;
  }
  const res = await completion.json();
  //console.log(res)
  const msg = res.choices[0].message;
  messages.value.push(msg);
  //console.log(messages.value)
  state.isResponse = false;
  await speechBackend(msg.content, config.modelUrl);
};

const scrollToBottom = () => {
  nextTick(() => {
    if (messageContainerRef && messageContainerRef.value) {
      messageContainerRef.value!.setScrollTop(messageInnerRef.value!.clientHeight);
    }
  });
};
watch(
  () => messages.value,
  () => {
    scrollToBottom();
  },
  {
    deep: true,
  }
);

const recorder = ref<any>();
async function startRecord() {
  try {
    await startRecording();
  } catch (error) {
    console.log(error);
  }
}
function stopRecord() {
  try {
    stopRecording();
  } catch (error) {
    console.log(error);
  }
}
const startRecording = async () => {

  if (!currentAudio.value) {
    notifyError(t('spoken.downloadVocieFirst'));
    return;
  }
  const constraints = {
    audio: {
      echoCancellation: true,
      noiseSuppression: true,
      sampleRate: 16000, // 设置采样率为16kHz
      channelCount: 1, // 设置为单声道
    },
    video: false,
  };
  try {
    const stream = await navigator.mediaDevices.getUserMedia(constraints);
    recorder.value = new MediaRecorder(stream);
    const audioChunks = <any>[];
    // 录音开始
    recorder.value.start();
    state.isRecording = true;

    // 录音数据
    recorder.value.ondataavailable = (e: any) => {
      audioChunks.push(e.data);
    };
    recorder.value.onstop = async () => {
      const blobData = new Blob(audioChunks, { type: "audio/wav" });
      // For 16-bit audio
      const blob = await getWaveBlob(blobData, false);
      //console.log(blob)
      const file = new File([blob], "recording.wav", {
        type: "audio/wav",
      });
      //console.log(file)
      const formData = new FormData();
      formData.append("file", file);

      //console.log(currentAudio.value);
      formData.append("model", currentAudio.value);
      const modelData = await getModel(currentAudio.value);
      //console.log(modelData);
      formData.append("params", JSON.stringify(modelData.params));
      //console.log(formData);
      const res = await fetch(getSystemKey("modelUrl") + "/voice", {
        body: formData,
        method: "POST",
      });
      //console.log(res);
      if (!res.ok) {
        notifyError(res.statusText);
        return;
      }
      const data = await res.json();
      //console.log(data);
      
      // 停止媒体流的所有轨道
      stream.getTracks().forEach((track) => track.stop());

      if (data.data.txt && data.data.txt != "") {
        userMessage.value = data.data.txt;
        sendMessage();
      } else {
        state.isRecording = false;
      }
    };
  } catch (err) {
    console.error(err);
  }
};

const stopRecording = () => {

  if (recorder.value) {
    recorder.value.stop();
    state.isRecording = false;
    state.isResponse = true;
  }
};
const changePlaying = () => {
  if(!currentTts.value){
    notifyError(t('spoken.downloadTTSFirst'));
    return;
  }
  if(audioElement && ifPlaying.value){
    audioElement.pause(); // 使用 pause 方法来停止播放
    state.isPlaying = false;
  }
  ifPlaying.value = !ifPlaying.value;
  

};
const state = reactive({
  isRecording: false,
  isResponse: false,
  isPlaying: false,
});

const clearMessages = () => {
  state.isRecording = false;
  state.isResponse = false;
  state.isPlaying = false;
};

</script>

<template>
  <div class="chat-bot">
    <el-page-header icon="null">
      <template #title>
        <div></div>
      </template>
      <template #content>
        <el-space>
          <el-button
            @click.stop="sexGender = !sexGender"
            :icon="sexGender ? 'UserFilled' : 'Avatar'"
            circle
          />
          <el-button
            @click.stop="changePlaying"
            :icon="ifPlaying ? 'Mic' : 'Mute'"
            circle
          />
          <el-button @click.stop="clearMessages" icon="DeleteFilled" circle />
        </el-space>
      </template>
    </el-page-header>
    <div class="messsage-area">
      <el-scrollbar
        v-if="messages.length > 1"
        class="message-container"
        ref="messageContainerRef"
      >
        <div ref="messageInnerRef" class="message-inner">
          <chat-message
            v-for="(message, key) in messages"
            :key="key"
            :content="message.content"
            :role="message.role"
            :createdAt="Date.now()"
          />
        </div>
      </el-scrollbar>
      <Vue3Lottie v-else 
      animationLink="/bot/spoken.json" 
      class="no-message-container" 
      :height="450" 
      :width="450" />
    </div>
    <div class="input-area">
      <Vue3Lottie
        animationLink="/bot/sbtn.json"
        :height="60"
        :width="60"
        v-if="state.isRecording"
        @click="stopRecord"
        class="animation-recording"
      />

      <span class="loader2" v-else-if="state.isResponse"> </span>

      <span class="loader1" v-else-if="state.isPlaying"> </span>

      <el-button
        v-else
        size="large"
        style="width: 60px; height: 60px"
        plain
        circle
        type="info"
        @click="startRecord"
        class="input-panel"
      >
        <el-icon :size="32">
          <Microphone />
        </el-icon>
      </el-button>
    </div>
  </div>
</template>

<style scoped lang="scss">
@import "@/assets/spoken.scss";
</style>
