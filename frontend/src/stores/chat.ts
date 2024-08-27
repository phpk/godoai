import { defineStore } from "pinia"
import { db } from "./db.js"
import { t } from "@/i18n/index.ts"
import {
  getCurrentModelName,
  getPrompt
} from "@/stores/config";
import { ref } from "vue";
export const useChatStore = defineStore('chat', () => {
  //const modelsStore = useModelsStore()
  const activeId: any = ref(0)
  const currentMessage:any = ref({})
  // 聊天列表
  const chatList: any = ref([])
  const chatInfo:any = ref({})
  const messageList : any = ref([])
  const newChat = async() => {
    const currentModel = getCurrentModelName('chat')
    if (!currentModel) {
      return false
    }
    const promptData = await getPrompt('chat')
    return await addChat(t('chat.newchat'), currentModel, promptData, "")
  }
  const initChat = async () => {
    if (activeId.value === 0) {
      await newChat()
    }
  }
  const getActiveChat = async () => {
    chatInfo.value = await db.getOne('chats', activeId.value)
    messageList.value = await db.getByField('messages', 'chatId', activeId.value)
    chatList.value = await db.getAll('chats')
    return { chatInfo, messageList, chatList }
  }
  const getChatList = async () => {
    chatList.value = await db.getAll('chats')
    return chatList
  }
  // 添加聊天
  async function addChat(title: string, model: string, promptData: any, kid:string) {
    const newChat = {
      title,
      prompt: promptData.prompt,
      promptName: promptData.name,
      model,
      createdAt: Date.now(),
      kid
    }
    //console.log(newChat)
    activeId.value = await db.addOne('chats', newChat)
    return activeId.value

  }
  async function setActiveId(newId: number) {
    activeId.value = newId
    chatInfo.value = await db.getOne('chats', newId)
    messageList.value = await db.getByField('messages', 'chatId', newId)
  }
  // 删除单个聊天
  async function deleteChat(chatId: number) {
    await db.delete('chats', chatId)
    await db.deleteByField('messages','chatId', chatId)
    //如果删除的id是当前id
    let id;
    if (chatId == activeId.value) {
      //
      const list = await db.getAll('chats')
      if(list.length > 0) {
        id = list[0]['id']
        
      }else{
        id = await newChat()
      }
      setActiveId(id)
    }
    chatList.value = await db.getAll('chats')
  }

  // 更新聊天菜单标题
  async function updateTitle(chatId: number, title: string) {
    await db.update('chats', chatId, {title})
  }

  // 清空所有Chat
  async function clearChat() {
    await db.clear('chats')
    await db.clear('messages')
  }


  // 新增历史消息
  async function addMessages(chatId: number, message: any) {
    const currentChat = await db.getOne('chats', chatId)
    //console.log(currentChat)
    if (currentChat) {
      return await db.addOne('messages', message)
    }
  }

  async function getChat(chatId: number) {
    const chats = await db.getOne('chats', chatId)
    //console.log(chats)
    const messages = await db.getByField('messages', 'chatId', chatId)
    return { chats, messages }
  }

  // 获取指定id的聊天的历史记录

  async function getChatHistory(chatId: number) {
    return await db.getByField('messages', 'chatId', chatId)
  }

  // 删除指定id的聊天的历史记录
  async function clearChatHistory(chatId: number) {
    await db.deleteByField('messages', 'chatId', chatId)
  }

  // 更新聊天配置
  async function updateChat(config: any, chatId: number) {
    //console.log(config)
    return await db.update('chats',chatId, config)
  }
  return {
    activeId,
    chatList,
    messageList,
    chatInfo,
    currentMessage,
    initChat,
    setActiveId,
    getActiveChat,
    getChatList,
    addChat,
    updateTitle,
    deleteChat,
    clearChat,
    addMessages,
    getChat,
    getChatHistory,
    clearChatHistory,
    updateChat
  }

}, {
  persist: {
    enabled: true,
    strategies: [
      {
        storage: localStorage,
        paths: [
          "activeId"
        ]
      }, // name 字段用localstorage存储
    ],
  }
})
//export const chatStore = useChatStore()
