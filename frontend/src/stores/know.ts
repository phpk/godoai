import { defineStore } from "pinia";
import { db } from "./db.ts"
import { ref, watchEffect } from "vue";
import {usePageStore} from "./page.ts"
import { getSystemConfig } from "@/stores/config";
import { getLang } from "@/i18n/index.ts";
const currentLang = getLang()
export const useKnowStore = defineStore('knowStore', () => {
  const knowList = ref([])
  const fileList:any = ref([])
  const cateList = ref([])
  const isShowFile = ref(false)
  const isAdd = ref(false)
  const knowLink = ref(false)
  const current = ref({})
  const dialog = ref(false)
  const isEdit = ref(false)
  const configList = ref([{
      type : "chromem",
      apiUrl : ""
    },
    {
      type : "chroma",
      apiUrl : "http://localhost:8000"
    },
    ])
  const pageStore = usePageStore()

  const getKnowList = async () => {
    // const res = await db.getAll('knowledge')
    await pageStore.getDataList('knowledge')
    knowList.value = pageStore.dataList
    
  }
  watchEffect(() => {
    knowList.value = [...pageStore.dataList];
  });
  const changeConfig = (type:string, url:string) => {
    configList.value.forEach((item:any) => {
      if(item.type === type){
        item.apiUrl = url
      }
    })
  }
  const getKnowAll = async () => {
    return await db.getAll('knowledge')

  }
  const getFileList = async (kid :number) => {
    const res = await db.getByField("knowfile","kid", kid)
    fileList.value = res
  }
  const getFileNames = () => {
    const names:any = []
    fileList.value.forEach((d : any) => {
      names.push(d.name)
    });
    return names
  }
  const addFile = async (data: any) => {
    await db.addAll('knowfile', data)
    await getFileList(data[0]['kid'])
  }

  const addKnow = async (data: any) => {
    //console.log(data)
    const id = await db.addOne('knowledge', data)
    await getKnowList()
    await getFileList(id)
  }
  const updateKnow = async (data : any) => {
    await db.update('knowledge', data.id, data)
    current.value = data
  }
  const getKnow = async (id : number) => {
    console.log(id)
    return await db.getOne('knowledge', id)
  }
  const delKonw = async (id : number) => {
    await db.delete('knowledge', id)
    await db.deleteByField("knowfile", "kid", id)
    await getKnowList()
  }

  const hasKnow = async (title : string) => {
    const res = await db.get('knowledge', { title: title })
    return res
  }
  const delFile = async (id : number, kid:number) => {
    await db.delete('knowfile', id)
    await getFileList(kid)
  }
  const hasFile = (names : any) => {
    const has = fileList.value.filter((d : any) => names.includes(d.name))
    //console.log(has)
    return has && has.length > 0
  }
  const updateFiles = async (data : any, kid: number) => {
    //console.log(data)
    const allList = await db.getByField("knowfile", "kid", kid)
    //console.log(allList)
    //console.log(kid)
    if(allList && allList.length > 0){
      const updates:any = []
      allList.forEach((d : any) => {
        if(data.includes(d.save_path)){
          updates.push({
            key : d.id,
            changes : {
              status : 1
            }
          })
        }
      })
      //console.log(updates)
      if(updates.length > 0){
        await db.table('knowfile').bulkUpdate(updates)
      }
    }
    await getFileList(kid)
  }
  const getCateList = async () => {
    const res = await db.getAll('knowcate')
    cateList.value = res
  }
  const getModelList = async (act : string) => {
    return await db.filter("modelslist", (d : any) => {
      //console.log(d)
      return d.action.includes(act)
    })
  }

  const askKnow = async (msg : string, uuid : string) => {
    const knowData = await db.get('knowledge', { uuid: uuid })
    if(!knowData)return false;
    const config = getSystemConfig()
    const askData = {
      message: msg,
      model: knowData.embedmodel,
      name: knowData.uuid,
      config: {
        type: knowData.type,
        apiUrl: knowData.apiUrl,
        embedding: {
          apiUrl: config.embedApi,
          apiType: config.embedType,
          apiKey: config.embedApiKey,
          contextLength : knowData.contextLength,
        }
      },
    };
    console.log(askData)
    const response = await fetch(config.knowledgeUrl + "/ask", {
      method: "POST",
      body: JSON.stringify(askData),
    });
    if (!response.ok) {
      return false;
    }
    const json = await response.json();
    if (json.code == 0 && json.data) {
      let prompt = await db.field("prompts", { lang: currentLang, action: "knowledge", isdef : 1}, "prompt")
      if(!prompt){
        return false;
      }

      let context: string = "";
      let currentLength = 0;
      const maxLength = knowData.contextLength;

      for (let i = 0; i < json.data.length && currentLength + json.data[i].content.length + 3 <= maxLength; i++) { // +3 accounts for the length of "- \n"
        const item = json.data[i];
        context += "- " + item.content + "\n";
        currentLength += item.content.length + 2; // +2 accounts for the length of "- " without the newline
      }

      prompt = prompt.replace("{content}", context);
      return {prompt, link : json.data}
    }else{
      return false
    }

  }

  return {
    knowList,
    fileList,
    cateList,
    isShowFile,
    current,
    isAdd,
    knowLink,
    dialog,
    isEdit,
    configList,
    changeConfig,
    getKnowList,
    getKnowAll,
    addFile,
    delFile,
    hasFile,
    addKnow,
    getKnow,
    delKonw,
    hasKnow,
    askKnow,
    updateKnow,
    getFileList,
    updateFiles,
    getFileNames,
    getCateList,
    getModelList
  }
}, {
  persist: {
    enabled: true,
    strategies: [
      {
        storage: localStorage,
        paths: [
          "configList"
        ]
      }, // name 字段用localstorage存储
    ],
  }
})
//export const knowStore:any = useKnowStore()
