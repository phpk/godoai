import {getSystemKey} from './config.ts'
export async function getBaseModelList(){
    return await getAiokList()
}
export async function getAiokList() {
    try {
        const url = getSystemKey('modelUrl') + "/tags"
        const response: any = await fetch(url);
        if (!response.ok) {
            return []
        }
        const res = await response.json();
        //console.log(res)
        if (res !== null && !res.code) {
            return res
        } else {
            return []
        }
    } catch (e) {
       return []
    }
}
export async function getBaseModelInfo(model: string) {
    const url = getSystemKey('modelUrl') + "/show?model=" + model
    const response: any = await fetch(url);
    if (!response.ok) {
        return false
    }
    const res = await response.json();
    //console.log(res)
    if (res !== null) {
        return res
    } else {
        return false
    }
}
