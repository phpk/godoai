<p align="center">
    <img src="./build/appicon.png" width="120" height="120">
</p>

<h1 align="center">GodoAI</h1>
一款集前沿AI技术与高效办公功能于一体的创新型办公产品，通过深度整合模型下载、模型对话、知识库管理、图片生成、语音转文字、文字转语音、以及全面的图像处理能力，GodoAI实现了办公场景的全覆盖，让每一位用户都能享受到智能化带来的便捷与高效。
💝良心产品，匠心独运，核心功能全面实现本地化部署，仅需8G内存，即可轻松驾驭。 

<div align="center">

[![license][license-image]][license-url] 

[English](README.en.md) | 简体中文

[API](./frontend/public/help/zh-cn/api.md)

</div>


### 📥 下载安装(v1.0.0)

1. 💻 **Windows 用户**:
   
- Windows (AMD64) [**Web版**](https://gitee.com/ruitao_admin/godoai/releases/download/v1.0.0/godoai_windows_amd64.zip) [**桌面版**](https://gitee.com/ruitao_admin/godoai/releases/download/v1.0.0/godoai-windows-amd64-installer.zip)
- Windows (ARM64) [**Web版**](https://gitee.com/ruitao_admin/godoai/releases/download/v1.0.0/godoai_windows_arm64.zip) [**桌面版**](https://gitee.com/ruitao_admin/godoai/releases/download/v1.0.0/godoai-windows-arm64-installer.zip)


2. 💼 **MacOS 用户**:

- MacOS (AMD64) [**Web版**](https://gitee.com/ruitao_admin/godoai/releases/download/v1.0.0/godoai_darwin_amd64.zip)
- MacOS (ARM64) [**Web版**](https://gitee.com/ruitao_admin/godoai/releases/download/v1.0.0/godoai_darwin_arm64.zip)

3. 💽 **Linux 用户**:

- Linux (AMD64) [**Web版**](https://gitee.com/ruitao_admin/godoai/releases/download/v1.0.0/godoai_linux_amd64.zip)
- Linux (ARM64) [**Web版**](https://gitee.com/ruitao_admin/godoai/releases/download/v1.0.0/godoai_linux_arm64.zip)

- 备注：web版下载后启动服务端。访问地址为：http://localhost:56710/。


<!-- 图标和链接 -->
[license-image]: https://gitee.com/ruitao_admin/godoos/raw/master/docs/img/license_%20MIT.svg
[license-url]: https://spdx.org/licenses/MIT.html

## 💝 亮点
- ***内置丰富的AI模型库，用户可根据需求一键下载所需模型***
- ***零配置，无需注册，下载即用***
- ***纯golang开发后端，低资源消耗和高性能***
- ***跨平台，Windows、Linux、MacOS都可以安装，无需docker部署***

## 💖 开源地址
- [Github](https://github.com/phpk/godoai)
- [Gitee](https://gitee.com/ruitao_admin/godoai)




## ⚡ 功能说明和预览

### 一、模型
- 内置丰富的AI模型库，用户可根据需求一键下载所需模型
- 几乎涵盖了市面上所有的基础大模型，包含聊天模型、文字转图片、图片转文字、图片转图片、文字转语音、语音转文字、Embeddings等模型 
<img src="https://godoos.com/static/images/ai/model.jpg" width="600" />

### 二、知识库
- 支持绝大多数类型文件上传，可抓取远程网页
- 支持本地向量数据库
<img src="https://godoos.com/static/images/ai/konw.jpg" width="600" />

### 三、创作
- 功能全面，包括从标题设定到内容编辑、保存、导出的全过程管理
- 支持关联知识库，可以不限制长度任意创作，发挥无限空间
- 支持段落生成、选中生成和全部生成
- 创作完成后可直接导出为word
- 支持查看创作列表、配置prompt高级设置等 
<img src="https://godoos.com/static/images/ai/create.jpg" width="600" />

### 四、绘图
- 无需GPU显卡也能生成图片
- 利用StableDiffusion，快速生成高质量图片，满足设计、广告、营销等多种场景需求。
- 支持SD中的所有模型,支持自定义SD中的参数
- 支持标签式生成prompt提示
- 生成的图片会自动保存到本地，支持实时查看后端的生成日志 
<img src="https://godoos.com/static/images/ai/draw.jpg" width="600" />

### 五、聊天
- 支持多模态模型对话
- 支持保存聊天内容
- 支持自定义聊天参数
- 支持关联知识库 
<img src="https://godoos.com/static/images/ai/chat.jpg" width="600" />

### 六、翻译
- 借助大模型可翻译多国语言
- 支持自定义翻译参数 
<img src="https://godoos.com/static/images/ai/trans.jpg" width="600" />

### 七、语聊
- 声音转文字中建议选择telespeech模型，支持多地方言识别。不同的模型语音识别语言不同。
- 支持语音转文字、文字转语音功能，无需打字即可快速和大模型沟通。 
<img src="https://godoos.com/static/images/ai/spoken.jpg" width="600" />

### 八、听录
- 支持语音转文字功能，会议记录、访谈笔记轻松搞定。
- 支持保存会议记录，方便后续管理和查看。
<img src="https://godoos.com/static/images/ai/recording.jpg" width="600" />


## 🏆 开发
### 构建
- 前端构建
```bash
cd frontend
pnpm i
pnpm build
```
- 后端构建
```bash
# go install github.com/wailsapp/wails/v2/cmd/wails@latest
wails build
# wails build -nsis -upx //you need install nsis and upx
```



## ❤️ 感谢
- [element-plus](http://element-plus.org/)
- [vue3](https://v3.cn.vuejs.org/)
- [wails](https://wails.io/)
- [ollama](https://ollama.com/)

## 💕 关联项目
- [godoos](https://gitee.com/ruitao_admin/godoos)
- [godooa](https://gitee.com/ruitao_admin/gdoa)
- [gododb](https://gitee.com/ruitao_admin/gododb)

## 微信群
<img src="https://gitee.com/ruitao_admin/gdoa/raw/master/docs/wx.png" width="150" />

## 开源

- 承诺永久免费开源
- 允许企业/个人单独使用
- 如用于商业活动或二次开发后发售，请购买相关版权
- 不提供私下维护工作，如有bug请 [issures](https://gitee.com/ruitao_admin/godoai/issues) 提交
- 请尊重作者的劳动成果

## 💌 支持作者

如果觉得不错，或者已经在使用了，希望你可以去 
<a target="_blank" href="https://gitee.com/ruitao_admin/godoai">Gitee</a> 帮我点个 ⭐ Star，这将是对我极大的鼓励与支持。