<script lang="ts" setup>
import { onMounted, onUnmounted, ref } from "vue";
import { getWaveBlob } from "@/util/audio.ts";
import { useRecordingStore } from "@/stores/recording.ts";
import { notifyError } from "@/util/msg";
import { getSystemKey, getModel, getCurrentModelName} from "@/stores/config.ts"
const recordingStore = useRecordingStore();
const isRecording = ref(false);
const sendTime = 5000;
let audioWorker: Worker | null = null;

// 录音实例
const mediaRecorder = ref<any>();
let buffer: Blob[] = []; // 缓冲区，存储过去5秒内的数据
let lastSendTime: number | null = null; // 存储上一次发送的时间戳
let timeoutId: any = null; // 存储setTimeout的ID
const currentAudio = ref("");
let rec:any;
const isSpeaking = ref(false);
const isBackend = ref(false);
const isShowTime = ref(true);
function initSpeechRecognition() {
  if (
    !("SpeechRecognition" in window) &&
    !("webkitSpeechRecognition" in window)
  ) {
    console.error("当前浏览器不支持SpeechRecognition语音识别");
    return null;
  }
  rec =
    "SpeechRecognition" in window
    ? new (window as any).SpeechRecognition()
    : new (window as any).webkitSpeechRecognition();

  rec.onstart = () => {
    console.log("开始监听语音");
  };

  rec.onspeechstart = async () => {
    isSpeaking.value = true;
  };

  rec.onspeechend = () => {
    isSpeaking.value = false;
  };
  rec.onresult = async (event:any) => {
    const last = event.results.length - 1; // 获取最后一个识别结果
    const msg = event.results[last][0].transcript; // 获取识别到的文本
    //console.log("用户说话的结果:", msg); // 打印结果
    if(msg && msg.length > 0){
      recordingStore.content += getFormatTime() + msg + "\n";
    }
    
    //rec.start()
  };
  rec.onerror = (event:any) => {
    console.error("语音识别错误", event.error);
    rec.stop();
  };

  rec.onend = () => {
    if (!isBackend.value && !isSpeaking.value) {
      rec.start();
    }
  };
}
// function changeBackend() {
//   isBackend.value = !isBackend.value;
//   if (isBackend.value) {
//     if (isSpeaking.value) {
//       rec.stop();
//     }
//     rec = null;
//   } else {
//     initSpeechRecognition();
//     //rec.start()
//   }
// }
async function startRecord() {
  isRecording.value = true;
  if (!isBackend.value) {
    rec.start();
  } else {
    await startRecording();
  }
}
function stopRecord() {
  isRecording.value = false;
  isSpeaking.value = true;
  if (!isBackend.value) {
    rec.stop();
  } else {
    stopRecording();
  }
}
const isMediaRecorderSupported = () => {
  return !!(
    navigator.mediaDevices &&
    (navigator.mediaDevices as any).getUserMedia &&
    MediaRecorder
  );
};
// 开始录音
async function startRecording() {
  if (!isBackend.value) {
    return;
  }
  if (!currentAudio.value) {
    isRecording.value = false;
    notifyError("需要下载声音转文字模型");
    return;
  }
  if (!isMediaRecorderSupported()) {
    isRecording.value = false;
    notifyError("浏览器不支持MediaRecorder API");
    return;
  }
  //isRecording.value = true;
  const constraints = {
    audio: {
      echoCancellation: true,
      noiseSuppression: true,
      //autoGainControl: false, // 尝试关闭自动增益控制
      sampleRate: 16000, // 设置采样率为16kHz
      channelCount: 1, // 设置为单声道
    },
    video: false,
  };
  const stream = await navigator.mediaDevices.getUserMedia(constraints);

  mediaRecorder.value = new MediaRecorder(stream);

  mediaRecorder.value.addEventListener("dataavailable", handleDataAvailable);
  mediaRecorder.value.addEventListener("stop", handleStop);

  mediaRecorder.value.start();

  // 启动定时器
  lastSendTime = performance.now();
  //sendBufferedData();
  timeoutId = setTimeout(sendBufferedData, sendTime);
}

// 停止录音
function stopRecording() {
  if (!isBackend.value) {
    return;
  }

  if (mediaRecorder.value) {
    mediaRecorder.value.stop();
  }

  // 清理定时器
  clearTimeout(timeoutId);
  timeoutId = null;
}

// 当录音数据可用时，收集数据块
function handleDataAvailable(event: BlobEvent) {
  //console.log(event.data)
  console.log("Buffer updated:", buffer.length);
  if (event.data.size > 0) {
    buffer.push(event.data); // 添加到缓冲区
  }
}

// 发送缓冲区内的所有数据
async function sendBufferedData() {
  if (lastSendTime === null || performance.now() - lastSendTime < sendTime) {
    return; // 如果未满5秒，不发送
  }
  sendMsg(true);
  // 更新上次发送时间
  lastSendTime = performance.now();

  // 重新设置定时器
  timeoutId = setTimeout(sendBufferedData, sendTime); // 调整延迟以确保5秒间隔
}
async function sendMsg(flag:any) {
  try {
    mediaRecorder.value.stop();
    console.log(buffer);
    if (buffer.length < 1) {
      if (flag) {
        mediaRecorder.value.start();
      }
      return;
    }
    const blobData = new Blob(buffer, { type: "audio/wav" });
    buffer = [];
    if (flag) {
      mediaRecorder.value.start();
    }
    // For 16-bit audio
    const blob = await getWaveBlob(blobData, false);
    console.log(blob);
    // Convert audio blob to a regular blob and serialize other data
    const modelData = await getModel(currentAudio.value);
  
    const formData = new FormData();
    const file = new File([blob], "recording.wav", { type: "audio/wav" });
    formData.append("file", file);
    formData.append("model", currentAudio.value);
    formData.append("params", JSON.stringify(modelData.params));
   //console.log(apiUrl);

    const responseFromServer = await fetch(getSystemKey("modelUrl") + "/voice", {
      method: "POST",
      body: formData,
    });
    
    if (!responseFromServer.ok) {
      throw new Error(`Error sending audio: ${responseFromServer.statusText}`);
    }
    const data = await responseFromServer.json();
    if(data.data && data.data.txt.length > 0){
      recordingStore.content += getFormatTime() + data.data.txt + "\n";
    }
  } catch (error) {
    console.error("Error sending audio:", error);
  }
}
// 录音停止时，处理并上传所有剩余的录音数据
async function handleStop() {
  //const blobData = new Blob(buffer, { type: "audio/wav" });
  sendMsg(false);
  //buffer = []; // 清空缓冲区
}
function getFormatTime() {
  if (isShowTime.value) {
    const now = new Date();
    const month = now.getMonth() + 1;
    const formattedMonth = month < 10 ? "0" + month : month.toString();
    const minute = now.getMinutes();
    const formatMinute = minute < 10 ? "0" + minute : minute.toString();
    const milliseconds = now.getMilliseconds();
    const formattedMilliseconds =
      milliseconds < 10 ? "0" + milliseconds : milliseconds.toString();
    return `- ${now.getFullYear()}-${formattedMonth}-${now.getDate()} ${now.getHours()}:${formatMinute}:${formattedMilliseconds}\n`;
  } else {
    return "";
  }
}
onMounted(async () => {
  initSpeechRecognition();

  currentAudio.value = getCurrentModelName("audio")


  await recordingStore.getRecordingList();
});
onUnmounted(() => {
  clearTimeout(timeoutId);
  if (audioWorker) {
    audioWorker.terminate();
    audioWorker = null;
  }
});
async function saveRecoding(saveData:any) {
  saveData.content = recordingStore.content;
  //console.log(saveData)
  await recordingStore.saveRecordingData(saveData);
}
</script>
<template>
  <div>
  <el-dialog v-model="recordingStore.showAdd" width="600">
    <recording-add @saveFn="saveRecoding" />
  </el-dialog>
  <el-drawer
    v-model="recordingStore.showLeft"
    direction="ltr"
    style="height: 100vh"
    :show-close="false"
    :with-header="false"
  >
    <recording-left />
  </el-drawer>
  <el-page-header icon="null">
    <template #title>
      <div></div>
    </template>
    <template #content>
      <el-button
        @click.stop="recordingStore.showLeft = true"
        icon="Menu"
        circle
      />
      <el-button
        @click.stop="isShowTime = !isShowTime"
        :icon="isShowTime ? 'Watch' : 'QuartzWatch'"
        circle
      />
    </template>
    <template #extra>
      <el-space class="mr-10">
        <el-button
          @click.stop="startRecord"
          v-if="!isRecording"
          icon="CaretRight"
          circle
        />
        <el-button @click.stop="stopRecord" v-if="isRecording" circle>
          <el-icon class="is-loading">
            <Loading />
          </el-icon>
        </el-button>
        <el-button
          @click.stop="recordingStore.showAdd = true"
          icon="Check"
          circle
        />
      </el-space>
    </template>
  </el-page-header>
  <el-card style="margin: 20px">
    <el-input
      type="textarea"
      :rows="25"
      v-model="recordingStore.content"
    />
  </el-card>
</div>
</template>
