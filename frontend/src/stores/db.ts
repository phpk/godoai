import Dexie from 'dexie'

export type ChatRole = 'user' | 'assistant' | 'system'
export type ChatType = 'chat' | 'code' | 'embed' | 'creation' | 'consultant' | 'img2txt' | 'translation' | 'spoken' | 'image' | 'video'
export type ChatTable = 'chats' | 'messages' | 'prompts' | 'models' | 'modelscate' | 'modelslabel' | 'modelspre' | 'modelslist' | 'knowledge' | 'knowcate' | 'knowfile' | 'knowmsg' | 'recording' | 'articles' | 'chatuser' | 'chatmsg' | 'chatfile'
export type ChanelType = 'GodoAI' | 'Ollama' | 'ChatGPT'
export type ModelType = 'llama' | 'sd' | 'rwkv' | 'whisper'
export type ModelLoadStatus = 'unload' | 'loading' | 'loaded'
export interface Chat {
  id?: number
  title: string
  chatType: ChatType
  prompt: string
  promptName?:string
  kid?:string
  model: string
  createdAt: Date
}

export interface Message {
  id?: number
  chatId: number
  chatType: ChatType
  role: ChatRole
  content: string
  createdAt: number
  link?:any
}

export interface Prompt {
  id?: number
  lang: string
  action: string //action
  ext?:string //action ext type
  prompt: string //content
  name: string //name
  isdef: number //0 or 1
  createdAt: Date
}
export interface Article {
  id?: number
  title: string
  cate: string
  base: string
  content: string
  createdAt: Date
}
export interface Model {
  id?: number
  name: string
  chanelType?: ChanelType
  modelType?: ModelType
  url?: string
  family?: string
  loadStatus?: ModelLoadStatus
  parameter_size?: string
  quantization_level?: string,
  hasModel?: string
  desc?: string
  endesc?: string
  createdAt?: Date
}
export interface ModelCate {
  id?: number
  name: string
  zhdesc: string
  endesc: string
  current?: string
}
export interface ModelPre {
  id?: number
  name: string
  from: string
  label: string
}
export interface ModelLabel {
  id?: number
  name: string
  zhdesc: string
  endesc: string
  family?: string
  action:string
  hasModel: string
  //引擎
  type?: string
}
export interface ModelList {
  id?: number
  model: string
  label: string
  status: string
  progress: number
  isLoading?:number
  url?:string
  more?: Array<string>
  action: Array<string>
  engine:string
  chanel:string
  template:string
  // deno-lint-ignore no-explicit-any
  info?: any
}
export interface ModelDown {
  name: string
  status: number
  label: string
}
export interface Knowledge {
  id?: number
  name: string
  title: string
  model: string
  embed: string
  ranker?: string
  files: Array<string>
  isme:boolean
  createdAt?: Date
}
// export interface Knowcate{
//   id?:number
//   name:string
// }
// export interface KnowTrunk{
//   name:string
//   content:string
// }
export interface Knowfile{
  id?:number
  name:string
  path:string
  cid:number
  content:string
  trunkType:string
  isTrunk:boolean
  //trunkData?: Array<KnowTrunk>
  createdAt?:Date
}
export interface Knowmsg{
  id?:number
  content:string
  role:string
  knowId:string
  createdAt?: Date
}
export const dbInit:any = new Dexie('GodoChatDatabase');
dbInit.version(1).stores({
  chats: '++id,title,prompt,promptName,ext,kid,createdAt',
  messages: '++id,chatId,role,content,createdAt',
  prompts: '++id,lang,action,prompt,name,ext,isdef,createdAt,[action+lang]',
  articles: '++id,title,cate,base,content,createdAt',
  knowledge: '++id,uuid,title,chatmodel,embedmodel,contextLength,createdAt',
  knowfile:'++id,name,save_path,content,split,kid,status,created_at',
  knowmsg: '++id,knowId,role,content,createdAt',
  modelscate: '++id,name,zhdesc,endesc,current',
  modelslabel: '++id,name,zhdesc,endesc,family,chanel,models,action,engine',
  modelslist: '++id,model,label,status,progress,url,file_name,action,chanel,engine,info,options',
  recording: '++id,name,time,local,desc,content,sumary,createdAt',
  chatuser:'++id,ip,hostname,userName,avatar,mobile,nickName,isOnline,updatedAt,createdAt',
  chatmsg:'++id,targetId,targetIp,senderInfo,reciperInfo,content,type,status,isRead,isMe,readAt,createdAt',
  chatfile:'++id,msgId,fileName,fileSize,fileType,fileUrl,filePath,fileExt'
});
export const db = {

  async getMaxId(tableName: ChatTable) {
    const data = await dbInit[tableName].orderBy('id').reverse().first()
    if (!data) {
      return 0
    } else {
      return data.id
    }
  },
  async getInsertId(tableName: ChatTable) {
    const id: any = await this.getMaxId(tableName)
    return id + 1
  },
  async getPage(tableName: ChatTable, page?: number, size?: number) {
    page = (!page || page < 1) ? 1 : page
    size = size ? size : 10
    const offset = (page - 1) * size
    return dbInit[tableName]
      .orderBy("id")
      .reverse()
      .offset(offset)
      .limit(size)
      .toArray();
  },
  async getAll(tableName: ChatTable) {
    return dbInit[tableName].toArray()
  },
  async count(tableName: ChatTable) {
    return dbInit[tableName].count()
  },
  async countSearch(tableName: ChatTable, whereObj?: any) {
    return dbInit[tableName].where(whereObj).count()
  },
  async pageSearch(tableName: ChatTable, page?: number, size?: number, whereObj?: any) {
    page = (!page || page < 1) ? 1 : page
    size = size ? size : 10
    const offset = (page - 1) * size
    //console.log(whereObj)
    return dbInit[tableName]
      .where(whereObj)
      .reverse()
      .offset(offset)
      .limit(size)
      .toArray();
  },
  async filter(tableName: ChatTable, filterFunc : any) {
    return dbInit[tableName].filter(filterFunc).toArray()
  },
  table(tableName: ChatTable) {
    return dbInit[tableName]
  },
  async getOne(tableName: ChatTable, Id: number) {
    return dbInit[tableName].get(Id)
  },
  async getRow(tableName: ChatTable, fieldName: string, val: any){
    return dbInit[tableName].where(fieldName).equals(val).first()
  },
  async get(tableName: ChatTable, whereObj : any) {
    //console.log(whereObj)
    const data = await dbInit[tableName].where(whereObj).first()
    //console.log(data)
    return data? data : false
  },
  async rows(tableName: ChatTable, whereObj: any) {
    return dbInit[tableName].where(whereObj).toArray()
  },
  async field(tableName: ChatTable, whereObj: any, field: string) {
    const data = await this.get(tableName, whereObj)
    return data ? data[field] : false
  },
  async getValue(tableName: ChatTable, fieldName: string, val: any, fName : string) {
    const row = await this.getRow(tableName, fieldName, val);
    return row[fName]
  },
  async getByField(tableName: ChatTable, fieldName: string, val: any) {
    return dbInit[tableName].where(fieldName).equals(val).toArray()
  },
  async addOne(tableName: ChatTable, data: any) {
    return dbInit[tableName].add(data)
  },
  async addAll(tableName: ChatTable, data: any) {
    return dbInit[tableName].bulkAdd(data)
  },
  async update(tableName: ChatTable, Id?: number, updates?: any) {
    return dbInit[tableName].update(Id, updates)
  },
  async modify(tableName: ChatTable, fieldName: string, val: any, updates: any) {
    return dbInit[tableName].where(fieldName).equals(val).modify(updates)
  },
  async delete(tableName: ChatTable, Id?: number) {
    return dbInit[tableName].delete(Id)
  },
  async deleteByField(tableName: ChatTable, fieldName: string, val: any) {
    return dbInit[tableName].where(fieldName).equals(val).delete()
  },
  async clear(tableName: ChatTable) {
    return dbInit[tableName].clear()
  },
  async clearAll() {
    return dbInit.delete()
  }
}
