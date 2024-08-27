import { defineStore } from 'pinia'
import { db, Prompt } from "./db.ts"
import { getLang } from "@/i18n/index.ts"
import { ref } from "vue";
import { creationCate } from "./models/cate.ts"
import { getCurrentModelName } from "./config.ts"

export const useCreationStore = defineStore('creation', () => {
  const currentLang = getLang()
  const baseTitle = ref("");
  const activeId = ref(0);
  const baseContent: any = ref("");
  const targetContent: any = ref("");
  const isLoading = ref(false);
  const isContentLoading = ref(false);
  const isSelectLoading = ref(false);
  const currentCate = ref("");
  const showLeft = ref(false);
  const baseEditer: any = ref("");
  const targetEditor: any = ref("");
  const articleList: any = ref([])
  const showConfig = ref(false)
  const openHelp = ref(true)

  const config = ref({
    model: "",
    modelList: [],
    knowList:[],
    refKnowledge:"",
    systemPrompt: "",
    systemPromptTitle: "",
    systemList: [],
    leaderPrompt: "",
    leaderPromptTitle: "",
    leaderList: [],
    builderPrompt: "",
    builderPromptTitle: "",
    builderList: []
  })
  const page = ref({
    current: 1,
    size: 10,
    total: 0,
    pages: 0,
    visible: 5
  })
  const areaSet = ref({
    left : 10,
    right: 13
  })
  function handlerArea(){
    if(areaSet.value.left == 10) {
      areaSet.value.left = 0
      areaSet.value.right = 24
    }else{
      areaSet.value.left = 10
      areaSet.value.right = 13
    }
  }
  async function initModel() {
    config.value.modelList = await db.filter("modelslist", (d: any) => {
      return d.action.includes("chat")
    })
    if (config.value.model == "") {
      //config.value.model = await db.field("modelscate", { name: "chat" }, "current")
      config.value.model = getCurrentModelName("chat");
    }
  }
  async function initKnow() {
    config.value.knowList = await db.getAll("knowledge")
  }
  async function initPrompt() {
    const list = await db.getAll("prompts")
    config.value.systemList = list.filter((item: Prompt) => item.action === "creation_system" && item.lang == currentLang)
    config.value.leaderList = list.filter((item: Prompt) => item.action === "creation_leader" && item.lang == currentLang)
    config.value.builderList = list.filter((item: Prompt) => item.action === "creation_builder" && item.lang == currentLang)
    const builderPrompt = list.find((item: Prompt) => item.action === "creation_builder")
    if (builderPrompt && config.value.builderPrompt == ""){
      config.value.builderPrompt = builderPrompt.prompt
      config.value.builderPromptTitle = builderPrompt.name
    }
    const leaderPrompt = list.find((item: Prompt) => item.action === "creation_leader")
    if(leaderPrompt && config.value.leaderPrompt == ""){
      config.value.leaderPrompt = leaderPrompt.prompt
      config.value.leaderPromptTitle = leaderPrompt.name
    }
    const systemPrompt = list.find((item: Prompt) => item.action === "creation_system")
    if(systemPrompt && config.value.systemPrompt == ""){
      config.value.systemPrompt = systemPrompt.prompt
      config.value.systemPromptTitle = systemPrompt.name
    }
  }
  async function initConfig() {
    await initPrompt()
    await initModel()
    await initKnow()
    await getActicleList()
  }
  function getLeaderPrompt(cate: string, title: string) {
    let res = config.value.leaderPrompt.replace(/\{title\}/g, title)
    res = res.replace(/\{cate\}/g, cate)
    return res
  }
  function getBuilderPrompt(title: string, content: string) {
    let res = config.value.builderPrompt.replace(/\{title\}/g, title)
    res = res.replace(/\{content\}/g, content)
    return res
  }

  async function getActicleList() {
    articleList.value = await db.getPage('articles',
      page.value.current,
      page.value.size)
    if (articleList.value.length == 0) {
      page.value.current = page.value.current > 1 ? page.value.current - 1 : 1
      articleList.value = await db.getPage('articles', page.value.current, page.value.size)
    }
    await getPageCount()
  }

  const getPageCount = async () => {
    page.value.total = await db.count('articles')
    page.value.pages = Math.floor(page.value.total / 10)
    // 检查是否有余数
    if (page.value.total % 10 !== 0) {
      // 如果有余数，则加1
      page.value.pages++;
    }
    //console.log(pageCount.value)
    return page.value
  }

  async function getActicle(id: number) {
    return await db.getOne("articles", id)
  }

  async function delArticle(id: number) {
    await db.delete("articles", id)
    await getActicleList()
  }

  const addArticle = async () => {
    const save = {
      title: baseTitle.value,
      base: baseEditer.value.getValue(),
      content: targetEditor.value.getValue(),
      cate: currentCate.value,
      createdAt: Date.now()
    }
    if (activeId.value < 1) {
      activeId.value = await db.addOne("articles", save)
    } else {
      await db.update("articles", activeId.value, save)
    }
    await getActicleList()
  }
  const pageClick = async (pageNum: any) => {
    //console.log(pageNum)
    page.value.current = pageNum
    await getActicleList()
  }
  return {
    creationCate,
    baseTitle,
    baseContent,
    targetContent,
    activeId,
    isLoading,
    isContentLoading,
    isSelectLoading,
    currentCate,
    showLeft,
    showConfig,
    baseEditer,
    targetEditor,
    articleList,
    page,
    config,
    areaSet,
    openHelp,
    handlerArea,
    addArticle,
    getActicle,
    delArticle,
    pageClick,
    getLeaderPrompt,
    getBuilderPrompt,
    initConfig,
    getActicleList
  }

}, {
  persist: {
    enabled: true,
    strategies: [
      {
        storage: localStorage,
        paths: [
          "openHelp",
          "config"
        ]
      }, // name 字段用localstorage存储
    ],
  }
})
//export const creationStore = useCreationStore()
