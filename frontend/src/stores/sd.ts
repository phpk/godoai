import { defineStore } from 'pinia'
import { ref } from "vue"
import {getModelList,getCurrentModelName,setCurrentModel} from "./config.ts"
export const useSdStore = defineStore('sd', () => {
  const imgList:any = ref([])
  const imageList = ref([])
  const modelList:any =ref([])
  const currentModel = ref("")
  const page = ref({
    current: 1,
    size: 9,
    total: 0,
    pages: 0,
  })
  const getDataList = () => {
    page.value.total = imgList.value.length
    page.value.pages = Math.ceil(page.value.total / page.value.size)
    imageList.value = imgList.value.slice((page.value.current - 1) * page.value.size, page.value.current * page.value.size)
  }
  const pageClick = async (pageNum: number) => {
    if(pageNum < 1){
      pageNum = 1
    }
    page.value.current = pageNum
    getDataList()
}
  const updateImgList = (list: []) => {
    imgList.value = list
  }
  const addImgList = (img: string[]) => {
    //console.log(img)
    if(!img || img.length < 1) return
    img.forEach((item:any) => {
      if (!imgList.value.includes(item)) {
        imgList.value.unshift(item)
      }
    })
    getDataList()
  }
  const initPage = async () => {
    currentModel.value = await getCurrentModelName("image")
    modelList.value = await getModelList("image")
    getDataList()
  }
  const getModel = () => {
    return modelList.value.find((item : any) => item.model == currentModel.value)
  }
  const delImgList = (img: string) => {
    imgList.value = imgList.value.filter((item : string) => {
      return item !== img
    })
    getDataList()
  }
  const updateModel = (model : string) => {
    currentModel.value = model
    setCurrentModel("image",model)
  }
  return {
    imgList,
    imageList,
    modelList,
    currentModel,
    page,
    getModel,
    updateImgList,
    delImgList,
    addImgList,
    initPage,
    updateModel,
    pageClick
  }
}, {
  persist: {
    enabled: true,
    strategies: [
      { storage: localStorage, paths: ["imgList"] }, // name 字段用localstorage存储
    ],
  }
})
//export const sdStore = useSdStore()
