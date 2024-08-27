import { getApiUrl, getSystemKey } from "./config"
const apiUrl = getApiUrl()
export const startAllApp = async () => {
  const flag = await fetch(apiUrl + "/startall")
  return flag.ok
}
export const llmType = getSystemKey("llmType")
export const startllm = async () => {
  const flag = await fetch(apiUrl + "/start/" + llmType)
  return flag.ok
}
export const loadModel = async (model:string) => {
  const modelUrl= getSystemKey("modelUrl")
  const flag = await fetch(modelUrl + "/load/" + llmType + "/" + model)
  return flag.ok
}
export const stopApp = async (name: string) => {
  const flag = await fetch(apiUrl + "/stop/" + name)
  return flag.ok
}
export const startApp = async (name: string) => {
  const flag = await fetch(apiUrl + "/start/" + name)
  return flag.ok
}
export const reStartApp = async (name: string) => {
  const flag = await fetch(apiUrl + "/restart/" + name)
  return flag.ok
}
