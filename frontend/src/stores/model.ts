import { defineStore } from "pinia";
import { ref } from "vue";
import { db, ModelLabel } from "./db.ts"
import { getSystemKey, setCurrentModel, updateCurrentModels, parseJson, setSystemKey } from './config.ts'
import { getBaseModelList, getBaseModelInfo } from './api.ts'
import { aiokLabels } from "./models/aioklabel.ts"
import { promptsZh, promptsEn } from "./models/prompt.ts"

export const useModelStore = defineStore('modelStore', () => {

  const labelList: any = ref([])
  const cateList: any = ref([])
  const modelList: any = ref([])
  const downList: any = ref([])
  const openHelp = ref(true)
  function getConfigKey(key: string) {
    return getSystemKey(key)
  }
  async function getLabelCate(cateName: string) {
    const list = await getLabelList()
    labelList.value = list.filter((d: ModelLabel) => {
      if (cateName == 'all') {
        return true
      } else {
        return d.action.indexOf(cateName) > -1
      }
    })
  }

  async function getLabelSearch(keyword: string) {
    const list = await getLabelList()
    if (!keyword || keyword == "") {
      labelList.value = list
    }
    labelList.value = list.filter((d: ModelLabel) => d.name.toLowerCase().includes(keyword.toLowerCase()))
  }
  async function getLabelList() {
    return await db.getAll("modelslabel")
    //return await db.getByField("modelslabel", "chanel", getSystemKey("currentChanel"))
  }
  async function delLabel(id: number) {
    await db.delete("modelslabel", id)
    labelList.value = await getLabelList()
  }
  async function checkLabelData(data: any) {
    const labelData = await db.get("modelslabel", { name: data.label })
    if (!labelData) {
      return
    }
    if (labelData.models.find((d: any) => d.model == data.model)) {
      return
    }
    labelData.models.push(data)

    await db.update("modelslabel", labelData.id, labelData)

  }

  async function getModelList() {
    const mList = await db.getAll("modelslist")
    let list: any
    if (!mList || mList.length < 1) {
      list = await getBaseModelList()
      if (list && list.length > 0) {
        await db.addAll("modelslist", list)
        modelList.value = list
        updateCurrentModels(list)
      } else {
        modelList.value = []
      }
    } else {
      modelList.value = mList
    }
  }
  function getModelInfo(model: string) {
    return modelList.value.find((d: any) => d.model == model)
  }

  async function getList() {
    cateList.value = getSystemKey('modelCate')
    labelList.value = await getLabelList()
    await getModelList()
    downList.value.forEach((_: any, index: number) => {
      downList.value[index].isLoading = 0
    })
  }
  function parseMsg(str: string) {
    const nres = { status: "" }
    try {
      //console.log(str)
      if (str == 'has done!') {
        return { status: 'success' }
      }
      const raw: any = str.split("\n")
      if (raw.length < 1) return nres
      // deno-lint-ignore no-explicit-any
      const rt: any = raw.filter((d: any) => d.trim() != "")
      //console.log(rt)
      if (rt.length > 0) {
        let msg = parseJson(rt.pop())
        if (msg) {
          return msg
        } else {
          msg = parseJson(rt.pop())
          return msg
        }
        //return JSON.parse(rt.pop())
      } else {
        return nres
      }
    } catch (error) {
      console.log(error);
      return nres
    }
  }

  async function addDownList(data: any) {
    //modelList.value.unshift(data)
    const has = modelList.value.find((d: any) => d.model == data.model)
    //console.log(has)
    if (!has) {
      //data = toRaw(data)
      const save = await getBaseModelInfo(data.model)
      //console.log(save)
      if (save) {
        modelList.value.unshift(save)
        return await db.addOne("modelslist", save)
      } else {
        console.log("not get model" + data.model)
      }

    }
  }
  async function deleteModelList(data: any) {
    //console.log(data)
    if (!data || !data.model) return
    modelList.value.forEach((d: any, index: number) => {
      if (d.model == data.model) {
        modelList.value.splice(index, 1);
      }
    });
    setCurrentModel(data.action, "", false);
    await db.deleteByField("modelslist", "model", data.model)
    //await db.delete("modelslist", data.id)
    await getModelList()
    updateCurrentModels(modelList.value)
  }

  function checkDownload(name: string) {
    return modelList.value.find((d: any) => d.model === name);
  }
  function addDownload(data: any) {
    const has = downList.value.find((d: any) => d.model === data.model)
    if (!has) {
      downList.value.unshift(data)
    } else {
      updateDownload(data)
    }

    return data
  }
  function deleteDownload(model: string) {
    //console.log(model)
    downList.value.forEach((d: any, index: number) => {
      if (d.model == model) {
        downList.value.splice(index, 1);
      }
    });
  }
  async function updateDownload(modelData: any) {
    const index = downList.value.findIndex((d: any) => d.model === modelData.model);
    if (index !== -1) {
      // 或者使用splice方法替换对象
      downList.value.splice(index, 1, {
        ...downList.value[index],
        status: modelData.status,
        progress: modelData.progress,
        isLoading: modelData.isLoading ?? 0,
      });
      if (modelData.status === "success") {
        await addDownList(modelData);
        await checkLabelData(modelData);
      }
    }
  }
  async function initSystem() {
    if (getSystemKey('isFirstRun')) return;
    await db.addAll("modelslabel", aiokLabels);
    await initPrompt()
    await getModelList()
    setSystemKey('isFirstRun', true);
  }

  async function clearSystem() {
    localStorage.clear();
    await db.clearAll();
  }

  async function initPrompt() {
    promptsZh.forEach((d: any) => {
      d.lang = "zh-cn"
      if (!d.action) {
        d.action = "chat"
      }

      if (!d.isdef) {
        d.isdef = 0
      }
    })
    promptsEn.forEach((d: any) => {
      d.lang = "en"
      if (!d.action) {
        d.action = "chat"
      }
      if (!d.isdef) {
        d.isdef = 0
      }
    })

    const save = [...promptsZh, ...promptsEn]
    await db.addAll("prompts", save)
  }

  return {
    cateList,
    labelList,
    modelList,
    downList,
    openHelp,
    getConfigKey,
    getList,
    getModelList,
    getModelInfo,
    checkDownload,
    addDownload,
    deleteDownload,
    updateDownload,
    checkLabelData,
    getLabelCate,
    getLabelSearch,
    getLabelList,
    delLabel,
    addDownList,
    deleteModelList,
    parseMsg,
    initSystem,
    clearSystem
  }

}, {
  persist: {
    enabled: true,
    strategies: [
      {
        storage: localStorage,
        paths: [
          "downList",
          "openHelp"
        ]
      }, // name 字段用localstorage存储
    ],
  }
})
