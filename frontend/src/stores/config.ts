import { aiokChatConfig } from "./models/chatconfig.ts"
import { db } from './db.ts'
export const configStoreType = localStorage.getItem('GodoOS-storeType') || 'browser';
export const systemCates = ['chat', 'translation', 'code', 'img2txt', 'image', 'tts', 'audio', 'embeddings']
/**
 * 获取系统配置信息。
 * 从本地存储中获取或初始化系统配置对象，并根据条件决定是否更新本地存储中的配置。
 * @param ifset 是否将配置信息更新回本地存储
 * @returns 当前系统配置对象
 */
export const getSystemConfig = (ifset = false) => {
    // 从本地存储中尝试获取配置信息，若不存在则使用默认空对象
    const configSetting = localStorage.getItem('GodoOS-config') || '{}';
    // 解析配置信息为JSON对象
    const config = JSON.parse(configSetting);

    // 初始化配置对象的各项属性，若本地存储中已存在则不进行覆盖
    if (!config.version) {
        config.version = '1.0.0';
    }
    if (!config.isFirstRun) {
        config.isFirstRun = false;
    }
    if (!config.lang) {
        config.lang = '';
    }
    // 初始化API相关URL，若本地存储中已存在则不进行覆盖
    if (!config.apiUrl) {
        config.apiUrl = 'http://localhost:56710';
    }
    if (!config.dataDir) {
        config.dataDir = '';
    }
    if (!config.chromaUrl) {
        config.chromaUrl = 'http://localhost:8000';
    }
    // 初始化当前频道的配置状态
    if (!config.aiokCurrent) {
        config.aiokCurrent = {}
        systemCates.forEach(d => {
            config.aiokCurrent[d] = ""
        })
    }
    if (!config.aiokChatConfig) {
        config.aiokChatConfig = aiokChatConfig
    }
    if (!config.embedApiKey) {
        config.embedApiKey = config.chatGptKey;
    }
    // if (!config.embedType) {
    //     config.embedType = config.currentChanel;
    // }
    // 初始化路由
    config.modelCate = systemCates
    config.knowledgeUrl = config.apiUrl + '/knowledge'
    config.modelUrl = config.apiUrl + '/model'
    config.chatApi = config.apiUrl + '/model/chat'
    config.embedApi = config.apiUrl + '/model/embeddings'
    config.modelDownApi = config.apiUrl + "/model/download"
    config.modelDeleteApi = config.apiUrl + "/model/delete"
    // 初始化系统相关信息，若本地存储中已存在则不进行覆盖
    //system
    if (!config.systemInfo) {
        config.systemInfo = {};
    }
    if(!config.dataDir){
        config.dataDir = ""
    }
    if (!config.theme) {
        config.theme = 'light';
    }
  
    
    // 根据参数决定是否更新本地存储中的配置信息
    if (ifset) {
        setSystemConfig(config)
    }
    // 返回配置对象
    return config;
};

export function getApiUrl() {
    return getSystemKey('apiUrl')
}

export async function getModelList(action: string = '') {
    const list = await db.getAll('modelslist')
    if (action == '') {
        return list
    } else {
        return list.filter((d: any) => d.action.includes(action))
    }
}
export async function getModel(model: string) {
    return db.get('modelslist', { model: model })
}
export async function getPrompt(action: string) {
    const config = getSystemConfig();
    const data = await db.get("prompts", {
        action,
        isdef: 1,
        lang: config.lang
    })
    if (data) {
        return data
    } else {
        return {
            name: 'defalut',
            prompt: ''
        }
    }
}
export async function getPromptList(action: any) {
    const config = getSystemConfig();
    if (action instanceof Array) {
        const list = await db.rows("prompts", {
            isdef: 1,
            lang: config.lang
        })
        return list.filter((d: any) => action.includes(d.action))
    } else {
        return await db.rows("prompts", {
            action: action,
            lang: config.lang
        })
    }

}
export function getCurrentModelName(cate: string = 'chat') {
    const config = getSystemConfig();
    return config.aiokCurrent[cate]
}
export function getCurrentModelList(list: any, cate: string = '') {
    if (!list || list.length < 1) return []
    const resItem = (item: any) => {
        if (!cate || cate == '') {
            return true
        } else {
            return item.action && item.action.includes(cate)
        }
    }
    list = JSON.parse(JSON.stringify(list))
    //console.log(list)
    return list.filter((item: any) => resItem(item))
}
export function getCurrents() {
    const config = getSystemConfig();
    return config.aiokCurrent
}

export function setCurrentModel(cate: any, val: string, ifsetDef = false) {
    const config = getSystemConfig();
    if (cate instanceof Array) {
        cate.forEach((d: any) => {
            if (ifsetDef && config.aiokCurrent[d] != '') {
                config.aiokCurrent[d] = val
            }
            if (!ifsetDef) {
                config.aiokCurrent[d] = val
            }
        })
    } else {
        if (ifsetDef && config.aiokCurrent[cate] != '') {
            config.aiokCurrent[cate] = val
        }
        if (!ifsetDef) {
            config.aiokCurrent[cate] = val
        }
    }
 
    setSystemConfig(config)
}

export function updateCurrentModels(models: any) {
    //console.log(models)
    if (!models || models.length < 1) return;
    //console.log(models)
    const config = getSystemConfig();
    systemCates.forEach((k: any) => {
        const has = models.find((m: any) => m.action.includes(k))
        if (has && config.aiokCurrent[k] == '') {
            config.aiokCurrent[k] = has.model
        }
    })
    setSystemConfig(config)
}
export function setChatConfig(chatConfig: any) {
    const config = getSystemConfig();
    config.aiokChatConfig = chatConfig
    setSystemConfig(config)
}
export function getChatConfig(action: string = '') {
    const config = getSystemConfig();
    if (action == '') {
        return config.aiokChatConfig
    } else {
        return config.aiokChatConfig[action]
    }
}
export function isWindowsOS() {
    return /win64|wow64|win32|win16|wow32/i.test(navigator.userAgent);
}
export function parseJson(str: string) {
    try {
        return JSON.parse(str);
    } catch (e) {
        return undefined;
    }
};


export const getSystemKey = (key: string, ifset = false) => {
    const config = getSystemConfig(ifset);
    if (key.indexOf('.') > -1) {
        const keys = key.split('.');
        return config[keys[0]][keys[1]];
    } else {
        return config[key];
    }
}

export const setSystemKey = (key: string, val: any) => {
    const config = getSystemConfig();
    config[key] = val;
    localStorage.setItem('GodoOS-config', JSON.stringify(config));
    if (key === 'storeType') {
        localStorage.setItem('GodoOS-storeType', val);
    }
};

export const setSystemConfig = (config: any) => {
    localStorage.setItem('GodoOS-config', JSON.stringify(config));
};

export const clearSystemConfig = () => {
    localStorage.clear()
};