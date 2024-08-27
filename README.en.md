<p align="center">
    <img src="./build/appicon.png" width="120" height="120">
</p>

<h1 align="center">GodoAI</h1>
GodoAI is an innovative office product that combines cutting-edge AI technology with efficient office functions. By deeply integrating model downloads, model dialogues, knowledge base management, image generation, speech to text, text to speech, and comprehensive image processing capabilities, GodoAI achieves full coverage of office scenes, allowing every user to enjoy the convenience and efficiency brought by intelligence.
💝 A conscientious product with unique craftsmanship, fully localized deployment of core functions, and easy control with only 8GB of memory.

<div align="center">

[![license][license-image]][license-url] 

[English](README.en.md) | [简体中文](README.md)

[API](./frontend/public/help/zh-cn/api.md)

</div>


### 📥 Download(v1.0.0)

1. 💻 **Windows**:
   
- Windows (AMD64) [**Web**](https://gitee.com/ruitao_admin/godoai/releases/download/v1.0.0/godoai_windows_amd64.zip) [**Desktop**](https://gitee.com/ruitao_admin/godoai/releases/download/v1.0.0/godoai-windows-amd64-installer.zip)
- Windows (ARM64) [**Web**](https://gitee.com/ruitao_admin/godoai/releases/download/v1.0.0/godoai_windows_arm64.zip) [**Desktop**](https://gitee.com/ruitao_admin/godoai/releases/download/v1.0.0/godoai-windows-arm64-installer.zip)


2. 💼 **MacOS**:

- MacOS (AMD64) [**Web**](https://gitee.com/ruitao_admin/godoai/releases/download/v1.0.0/godoai_darwin_amd64.zip)
- MacOS (ARM64) [**Web**](https://gitee.com/ruitao_admin/godoai/releases/download/v1.0.0/godoai_darwin_arm64.zip)

3. 💽 **Linux**:

- Linux (AMD64) [**Web**](https://gitee.com/ruitao_admin/godoai/releases/download/v1.0.0/godoai_linux_amd64.zip)
- Linux (ARM64) [**Web**](https://gitee.com/ruitao_admin/godoai/releases/download/v1.0.0/godoai_linux_arm64.zip)

- Note: Start the server after downloading the web version. The access address is：http://localhost:56710/。


<!-- 图标和链接 -->
[license-image]: https://gitee.com/ruitao_admin/godoos/raw/master/docs/img/license_%20MIT.svg
[license-url]: https://spdx.org/licenses/MIT.html

## 💝 Highlights 
- ***Built in rich AI model library, users can download the required models with just one click according to their needs*** 
- ***Zero configuration, no registration required, download and use immediately*** 
- ***Pure Golang development backend, low resource consumption and high performance*** 
- ***Cross platform, Windows, Linux, MacOS can all be installed without Docker deployment***

## 💖 Open source address
- [Gitee](https://gitee.com/ruitao_admin/godoai)




## ⚡ Function Description and Preview

### 一、Model
- Built in rich AI model library, users can download the required models with just one click according to their needs
- It covers almost all basic models on the market, including chat models, text to image, image to text, image to image, text to speech, speech to text, embeddings and other models
<img src="https://godoos.com/static/images/ai/model.jpg" width="600" />

### 二、knowledge base
- Supports uploading most types of files and can crawl remote web pages
- Support local vector database
<img src="https://godoos.com/static/images/ai/konw.jpg" width="600" />

### 三、a literary creation
- Comprehensive functionality, including full process management from title setting to content editing, saving, and exporting
- Support associated knowledge base, allowing for unlimited length and creativity, unleashing infinite space
- Support paragraph generation, selected generation, and all generation
- After the creation is completed, it can be directly exported as Word
- Support viewing creation lists, configuring prompt advanced settings, etc
<img src="https://godoos.com/static/images/ai/create.jpg" width="600" />

### 四、generate images
- Can generate images without the need for GPU graphics card
- Using StableDiffusion, quickly generate high-quality images to meet various scenarios such as design, advertising, and marketing needs.
- Support all models in SD and customize parameters in SD
- Support tag based generation of prompt prompts
- The generated images will be automatically saved locally and support real-time viewing of backend generation logs
<img src="https://godoos.com/static/images/ai/draw.jpg" width="600" />

### 五、chat
- Support multimodal model dialogue
- Support saving chat content
- Support custom chat parameters
- Support associated knowledge base
<img src="https://godoos.com/static/images/ai/chat.jpg" width="600" />

### 六、translate
- Using large models to translate multiple languages
- Support custom translation parameters 
<img src="https://godoos.com/static/images/ai/trans.jpg" width="600" />

### 七、Language Chat
- It is recommended to choose the Telespeech model for voice to text conversion, which supports recognition of multiple dialects. Different models have different speech recognition languages.
- Supports voice to text and text to speech functions, allowing for quick communication with large models without the need for typing. 
<img src="https://godoos.com/static/images/ai/spoken.jpg" width="600" />

### 八、Listen to the recording
- Support voice to text conversion function, easy to handle meeting minutes and interview notes.
- Support saving meeting minutes for easy management and viewing in the future.
<img src="https://godoos.com/static/images/ai/recording.jpg" width="600" />


## Development
### Build
- Front end construction
```bash
cd frontend
pnpm i
pnpm build
```
- Backend construction
```bash
# go install github.com/wailsapp/wails/v2/cmd/wails@latest
wails build
# wails build -nsis -upx //you need install nsis and upx
```
## thank
- [element-plus](http://element-plus.org/)
- [vue3](https://v3.cn.vuejs.org/)
- [wails](https://wails.io/)
- [ollama](https://ollama.com/)

## 💕 Associated Projects
- [godoos](https://gitee.com/ruitao_admin/godoos)
- [godooa](https://gitee.com/ruitao_admin/gdoa)
- [gododb](https://gitee.com/ruitao_admin/gododb)

## weChat group
<img src="https://gitee.com/ruitao_admin/gdoa/raw/master/docs/wx.png" width="150" />
