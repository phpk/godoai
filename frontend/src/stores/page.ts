import { defineStore } from "pinia";
import { db } from "./db.ts"
import { ref } from "vue";
export const usePageStore = defineStore('pageStore', () => {
    const page = ref({
        current: 1,
        size: 10,
        total: 0,
        pages: 0,
    })
    const dataList = ref([])
    const getDataList = async (t : any) => {
        dataList.value = await db.getPage(t, page.value.current, page.value.size)
        //console.log(dataList.value)
        if (dataList.value.length == 0) {
            page.value.current = page.value.current > 1 ? page.value.current - 1 : 1
            dataList.value = await db.getPage(t, page.value.current, page.value.size)
        }
        await getPageCount(t)
    }
    const getPageCount = async (t: any) => {
        page.value.total = await db.count(t)
        page.value.pages = Math.ceil(page.value.total / page.value.size)
        return page.value
    }
    const pageClick = async (pageNum: any, t: any) => {
        console.log(pageNum)
        page.value.current = pageNum
        await getDataList(t)
    }
    return {
        page,
        dataList,
        getDataList,
        pageClick,
        getPageCount
    }
})