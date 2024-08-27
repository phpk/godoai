export const chatgptLabels = [
    {
        name: "gpt-4-turbo",
        family: "chatgpt",
        engine: "openai",
        chanel: "chatgpt",
        models:[
            {
                model: "gpt-4-turbo",
                params : {
                    "max_tokens": 4096,
                    "temperature": 0.7,
                    "top_p": 0.95,
                    "stream": true,
                    "stop": [],
                    "frequency_penalty": 0,
                    "presence_penalty": 0
                }
            }
        ],
        action: ["chat", "code" , "translation"],
        zhdesc: "OpenAI GPT 4 Turbo模型表现出色",
        endesc: "OpenAI GPT 4 Turbo model is extremely good"
    },
    {
        name: "gpt-4-vision-preview",
        family: "chatgpt",
        engine: "openai",
        chanel: "chatgpt",
        models: [
            {
                model: "gpt-4-vision-preview",
                params: {
                    "max_tokens": 4096,
                    "temperature": 0.7,
                    "top_p": 0.95,
                    "stream": true
                }
            }
        ],
        action: ["img2txt"],
        zhdesc: "OpenAI GPT-4视觉模型具有视觉理解功能",
        endesc: "OpenAI GPT-4 Vision model features vision understanding capabilities"
    },
    {
        name: "gpt-3.5-turbo",
        family: "chatgpt",
        engine: "openai",
        chanel: "chatgpt",
        models: [
            {
                model: "gpt-3.5-turbo",
                params: {
                    "max_tokens": 4096,
                    "temperature": 0.7,
                    "top_p": 0.95,
                    "stream": true,
                    "stop": [],
                    "frequency_penalty": 0,
                    "presence_penalty": 0
                }
            }
        ],
        action: ["chat","translation"],
        zhdesc: "OpenAI GPT 3.5 Turbo模型速度极快",
        endesc: "OpenAI GPT 3.5 Turbo model is extremely fast"
    },
    {
        name: "gpt-4o",
        family: "chatgpt",
        engine: "openai",
        chanel: "chatgpt",
        models: [
            {
                model: "gpt-4o",
                params: {
                    "max_tokens": 4096,
                    "temperature": 0.7,
                    "top_p": 0.95,
                    "stream": true,
                    "stop": [],
                    "frequency_penalty": 0,
                    "presence_penalty": 0
                }
            }
        ],
        action: ["chat", "code", "translation"],
        zhdesc: "OpenAI GPT 4o是一款速度快、质量高的新旗舰机型",
        endesc: "OpenAI GPT 4o is a new flagship model with fast speed and high quality"
    },
    {
        name: "text-embedding-ada-002",
        family: "chatgpt",
        engine: "openai",
        chanel: "chatgpt",
        models: [
            {
                model: "text-embedding-ada-002",
                params: {}
            }
        ],
        action: ["embeddings"],
        zhdesc: "text-embedding-ada-002",
        endesc: "text-embedding-ada-002"
    },
    {
        name: "text-embedding-3-large",
        family: "chatgpt",
        engine: "openai",
        chanel: "chatgpt",
        models: [
            {
                model: "text-embedding-3-large",
                params: {}
            }
        ],
        action: ["embeddings"],
        zhdesc: "text-embedding-3-large",
        endesc: "text-embedding-3-large"
    },
    {
        name: "text-embedding-3-small",
        family: "chatgpt",
        engine: "openai",
        chanel: "chatgpt",
        models: [
            {
                model: "text-embedding-3-small",
                params: {}
            }
        ],
        action: ["embeddings"],
        zhdesc: "text-embedding-3-small",
        endesc: "text-embedding-3-small"
    }
]