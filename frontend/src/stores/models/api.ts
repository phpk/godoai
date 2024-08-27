export const modelFrom = ["hf-mirror.com", "huggingface.co", "ollama.com"];

export const chanelList = ["GodoAI", "Ollama", "ChatGPT"];


export interface ChannelApiInterface {
  [key: string]: string;
}
export const localOllamaUrl = "http://localhost:56715"
export const localFileUrl = "http://localhost:56713"
export const apiUrls: ChannelApiInterface = {
  GodoAI: "http://localhost:56710",
  Ollama: "http://localhost:11434",
  ChatGPT: "https://api.openai-proxy.com",
};

export const apiUrlRef = [
  {
    'name' : '进程',
    'key': 'apiUrl',
    'port' : '56710'
  },
  {
    'name': '模型管理',
    'key': 'modelUrl',
    'port': '56711'
  },
  {
    'name': '文字识别',
    'key': 'goconvUrl',
    'port': '56712'
  },
  {
    'name': '文件管理',
    'key': 'fileUrl',
    'port': '56713'
  },
  {
    'name': '知识库',
    'key': 'knowledgeUrl',
    'port': '56714'
  },
  {
    'name': '聊天模型',
    'key': 'ollamaUrl',
    'port': '56715'
  },
  {
    'name': '图像模型',
    'key': 'sdUrl',
    'port': '56716'
  },
  {
    'name': '声音模型',
    'key': 'voiceUrl',
    'port': '56717'
  },
  {
    'name': '内网聊天',
    'key': 'localchatUrl',
    'port': '56718'
  },

]

export const apiChat: ChannelApiInterface = {
  GodoAI: "/api/completions",
  Ollama: "/v1/chat/completions",
  ChatGPT: "/v1/chat/completions",
};

export const apiKeys: ChannelApiInterface = {
  GodoAI: "",
  Ollama: "",
  ChatGPT: "",
};

export const apiDown: ChannelApiInterface = {
  GodoAI: "/model/download",
  Ollama: "/api/pull",
};
export const apiTags: ChannelApiInterface = {
  GodoAI: "/api/tags",
  Ollama: "/api/tags",
};
export const apiDelete: ChannelApiInterface = {
  GodoAI: "/model/delete",
  Ollama: "/api/delete",
};
export const apiImage:string = "/image"
export const apiVoice:string = "/voice"

