import { defineStore } from 'pinia'
import {db} from './db.ts'
import { ref } from "vue"
import { t } from "@/i18n/index";
export const useRecordingStore = defineStore('recording', () => {
  const recordingList:any = ref([])
  const showLeft = ref(false)
  const showAdd = ref(false)
  const currentId = ref(0)
  const content = ref("")
  const title = ref(t('recording.newRecording'))
  const page =ref({
    current: 1,
    size: 10,
    total: 0,
    pages: 0,
    visible:5
  })
  const addRecording = async (recording: any) => {
    recordingList.value.push(recording)

    currentId.value = await db.addOne('recording', recording)
  }
  const saveRecordingData = async(saveData : any) => {
    saveData.createdAt = new Date()
    title.value = saveData.name
    if (currentId.value > 0) {
      await updateRecording(saveData);
    } else {
      await addRecording(saveData);
    }

    await getRecordingList();
    showLeft.value = true;
    showAdd.value = false;
  }
  const newRecording = () => {
    currentId.value = 0
    content.value = ""
    title.value = t('recording.newRecording')
  }
  const getRecording = async (id : number) => {
    return await db.getOne('recording', id)
  }
  const getRecordingList = async () => {
    //recordingList.value = await db.getAll('recording')
    recordingList.value = await db.getPage('recording',page.value.current, page.value.size)
    if (recordingList.value.length == 0){
      page.value.current = page.value.current > 1 ? page.value.current - 1 : 1
      recordingList.value = await db.getPage('recording', page.value.current, page.value.size)
    }
    await getPageCount()

  }
  const getPageCount = async () => {
    page.value.total = await db.count('recording')
    page.value.pages = Math.floor(page.value.total / 10)
    // 检查是否有余数
    if (page.value.total % 10 !== 0) {
      // 如果有余数，则加1
      page.value.pages++;
    }
    //console.log(pageCount.value)
    return page.value
  }
  const pageClick = async (pageNum: any) => {
    //console.log(pageNum)
    page.value.current = pageNum
    await getRecordingList()
  }
  const updateRecording = async (recording: any) => {
    //console.log(recording)
    await db.update('recording', currentId.value, recording)
  }
  const changeRecording = async (id : number) => {
    currentId.value = id
    const data = await db.getOne('recording', id)
    content.value = data.content
    title.value = data.name
    showLeft.value = false
  }
  const deleteRecording = async (id : number) => {
    await db.delete('recording', id)
    currentId.value = 0
    await getRecordingList()
  }
  return{
    title,
    content,
    currentId,
    showAdd,
    showLeft,
    recordingList,
    page,
    newRecording,
    getPageCount,
    pageClick,
    saveRecordingData,
    addRecording,
    getRecording,
    getRecordingList,
    updateRecording,
    changeRecording,
    deleteRecording
  }

})
