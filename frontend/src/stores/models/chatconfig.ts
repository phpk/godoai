export const aiokChatConfig = {
    chat : {
        key:"chat",
        contextLength: 10,
        num_keep: 5, //保留多少个最有可能的预测结果。这与top_k一起使用，决定模型在生成下一个词时考虑的词汇范围。
        num_predict: 3, //生成多少个预测结果
        top_p: 0.95,
        top_k: 40, //影响生成的随机性。较高的top_k值将使模型考虑更多的词汇
        temperature: 0.7, //影响生成的随机性。较低的温度产生更保守的输出，较高的温度产生更随机的输出。
    },
    translation: {
        key:"translation",
        num_keep: 5,
        num_predict: 1, 
        top_k: 40, 
        top_p: 0.95,
        temperature: 0.2,
    },
    creation:{
        key:"creation",
        num_keep: 3,
        num_predict: 1, 
        top_k: 40, 
        top_p: 0.95,
        temperature: 0.2,
    },
    knowledge:{
        key:"knowledge",
        contextLength: 10,
        num_keep: 5,
        num_predict: 1, 
        top_k: 40, 
        top_p: 0.95,
        temperature: 0.2,
    },
    spoken:{
        key:"spoken",
        contextLength: 10,
        num_keep: 5,
        num_predict: 1, 
        top_k: 40, 
        top_p: 0.95,
        temperature: 0.2,
    }  
}