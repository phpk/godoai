<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue';
import { getLang } from '@/i18n';
import MarkdownIt from 'markdown-it';
import { notifyError, notifySuccess } from "@/util/msg";
import {
    getApiUrl
} from "@/stores/config";
const searchQuery = ref('');
const lang = getLang();
// 定义渲染类型的枚举
enum RenderType {
    None = 'none',
    Markdown = 'markdown',
    Key = 'key',
    SearchResults = 'search-results',
    UserCenter = 'user-center',
}
const renderType = ref<RenderType>(RenderType.None); // 定义渲染类型
const mdContent = ref('');
const searchResults: any = ref([]);
const showUserCenter = ref(false);
const userCenterUrl = ref('https://godoos.com/home/login/');
const iframeRef = ref<HTMLIFrameElement>();
const applyCode = ref("");
const licenseCode = ref("");
const navItems = [
    { title: '帮助首页', link: 'home' },
    { title: '版本对比', link: 'version' },
];
const combinedList = [
    { text: '模型管理', link: 'model' },
    { text: '模型下载', link: 'down' },
    { text: 'AI对话', link: 'aichat' },
    { text: 'AI创作', link: 'creation' },
    { text: 'Prompt设置', link: 'prompt' },
    { text: 'AI绘图', link: 'sd' },
    { text: '系统设置', link: 'sysset' },
    { text: 'AI翻译', link: 'trans' },
    { text: 'AI语音', link: 'spoken' },
    { text: 'AI听录', link: 'recording' },
];



async function renderMd(link: string) {
    renderType.value = RenderType.None; // 重置渲染类型
    showUserCenter.value = false;
    searchResults.value = [];
    searchQuery.value = "";

    if (link === 'home') {
        return;
    }

    if (link === 'user-center') {
        renderType.value = RenderType.UserCenter;
        return;
    }
    if (link === 'key') {
        renderType.value = RenderType.Key;
        return;
    }

    const path = `/help/${lang}/${link}.md`;
    const response = await fetch(path);
    if (!response.ok) {
        return;
    }

    const markdown = new MarkdownIt();
    const content = await response.text();
    mdContent.value = markdown.render(content);

    renderType.value = RenderType.Markdown;
}

async function search() {
    searchResults.value = [];
    const fileList: any = [...navItems, ...combinedList];
    for (const item of fileList) {
        const path = `/help/${lang}/${item.link}.md`;
        const response = await fetch(path);
        if (!response.ok) {
            continue;
        }

        const content = await response.text();
        const markdown = new MarkdownIt();
        const htmlContent = markdown.render(content);
        const plainText = htmlContent.replace(/<[^>]*>/g, '');
        const query = searchQuery.value.toLowerCase();

        let matchIndex = plainText.toLowerCase().indexOf(query);
        while (matchIndex !== -1) {
            const contextStart = Math.max(matchIndex - 20, 0);
            const contextEnd = Math.min(matchIndex + query.length + 20, plainText.length);
            const highlightedContext = plainText.substring(contextStart, contextEnd)
                .replace(new RegExp(`(${query})`, 'gi'), '<span style="color:red">$1</span>');

            searchResults.value.push({
                text: item.text,
                link: item.link,
                context: highlightedContext,
            });

            matchIndex = plainText.toLowerCase().indexOf(query, matchIndex + query.length);
        }
    }

    if (searchResults.value.length > 0) {
        renderType.value = RenderType.SearchResults;
    } else {
        renderType.value = RenderType.None;
    }
}
function sendIframeMessage(message: any) {
    if (iframeRef.value && iframeRef.value.contentWindow) {
        iframeRef.value.contentWindow.postMessage(message, userCenterUrl.value);
    }
}
async function eventHandler(event: any) {
    //console.log("eventHandler", event)
    if (event.data.action === 'requestApplyCode') {
        sendIframeMessage({
            action: 'provideApplyCode',
            applyCode: applyCode.value
        })
    }
}
async function getApplyCode() {
    const apiUrl = getApiUrl()
    const completion = await fetch(apiUrl + "/getLicense")
    if (completion.ok) {
        const res = await completion.json()
        applyCode.value = res.data
    }
}
async function saveCode() {
    const apiUrl = getApiUrl()
    if (licenseCode.value == '' && licenseCode.value.length < 32) {
        notifyError("密钥不能为空")
        return
    }
    const completion = await fetch(apiUrl + "/setLicense", {
        method: "POST",
        body: JSON.stringify({
            licenseCode: licenseCode.value
        })
    })
    if (!completion.ok) {
        notifyError("服务器错误")
        return
    }
    const res = await completion.json()
    if (res.message == 'sucess') {
        notifySuccess("设置成功")
    } else {
        notifyError(res.message)
    }
}
onMounted(async () => {
    await getApplyCode()
    window.addEventListener("message", eventHandler);
});

onUnmounted(() => {
    window.removeEventListener("message", eventHandler);
});
</script>
<template>
    <div id="help-page">
        <el-menu mode="horizontal" :ellipsis="false">
            <el-menu-item v-for="(item, index) in navItems" :key="index" :index="item.link"
                @click="renderMd(item.link)">
                {{ item.title }}
            </el-menu-item>
            <div class="flex-grow" />
            <el-input v-model="searchQuery" style="margin:0;height: 45px;" placeholder="Search...">
                <template #append>
                    <el-button icon="Search" @click="search()" />
                </template>
            </el-input>
        </el-menu>

        <el-scrollbar class="help-content">
            <el-space v-show="renderType === 'none'" wrap class="list-container">
                <div v-for="(item, key) in combinedList" :key="key" class="list-item" @click="renderMd(item.link)">
                    {{ item.text }}
                </div>
            </el-space>
            <div v-show="renderType === 'search-results' && searchResults.length > 0">
                <h3>搜索结果:</h3>
                <div v-for="(result, key) in searchResults" :key="key" class="search-result-item">
                    <div v-html="result.context"></div>
                    <div class="result-link" @click="renderMd(result.link)">
                        查看全文
                    </div>
                </div>
            </div>
            <div v-show="renderType === 'markdown'">
                <div v-html="mdContent" class="markdown-content"></div>
            </div>
            <div v-show="renderType === 'user-center'">
                <iframe :src="userCenterUrl" ref="iframeRef" frameborder="0" class="user-center-iframe"></iframe>
            </div>
            <div v-show="renderType === 'key'">
                <el-space wrap>
                    <el-form style="margin-top: 12px;width:84vw">
                        <el-form-item>
                            <el-input type="textarea" placeholder="输入License" v-model="licenseCode"
                                :rows="6"></el-input>
                        </el-form-item>
                        <el-form-item>
                            <el-button type="primary" @click="saveCode">激活</el-button>
                        </el-form-item>
                    </el-form>
                    <el-card>
                        <template #header>申请码</template>
                        <div class="keyarea">{{ applyCode }}</div>
                        <template #footer>
                            <el-row justify="end" style="border-bottom: none;margin:0;padding:0">
                                <copy-btn :text="applyCode" style="margin-left: 30px;" />
                            </el-row>
                        </template>
                    </el-card>

                    
                </el-space>
            </div>
        </el-scrollbar>


    </div>
</template>
<style scoped>
#help-page {
    min-height: 100vh;
    /* 确保页面至少与视口高度相同 */
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 20px;
    /* 添加内边距 */
    color: var(--el-text-color-primary);
    background-color: var(--el-color-bg-base);
}

.help-content {
    width: 100%;
    height: calc(100vh - 170px);
    text-align: left;
    line-height: 30px;
}

/* 为 Element Plus 的组件添加样式 */
.el-menu {
    width: 100%;
    background-color: var(--el-color-bg-base);
    justify-content: center;
}

.el-input {
    width: 50%;
    /* 设置搜索框宽度为50% */
    margin-top: 20px;
    margin-bottom: 20px;
    /* 添加底部边距 */
    background-color: var(--el-fill-color-light);
    border-color: var(--el-border-color);
}

/* 列表容器和列表项样式 */
.list-container {
    width: 100%;
    max-width: 800px;
    display: flex;
    margin-top: 30px;
    flex-wrap: wrap;
    justify-content: center;
    /* 居中对齐列表 */
    gap: 20px;
    /* 列表项之间的间距 */
}

.list-item {
    width: 290px;
    /* 每个列表项占据一半宽度，减去gap的一半 */
    background-color: var(--el-color-bg-base);
    list-style-type: none;
    border: 1px solid var(--el-border-color);
    padding: 15px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    transition: all 0.3s ease;
    text-align: left;
}

.search-results {
    margin-top: 20px;
    width: 100%;
    max-width: 800px;
}

.search-result-item {
    margin-bottom: 20px;
    padding: 10px;
    background-color: var(--el-color-bg-base);
    border: 1px solid var(--el-border-color);
}

.user-center-iframe {
    width: 100%;
    height: 600px;
    border: none;
}

.result-link {
    cursor: pointer;
    color: var(--el-color-primary);
    text-decoration: underline;
}

/* 内容容器的样式 */
.markdown-content {
    padding: 20px;
}

.markdown-content :deep(h1, h2, h3, h4, h5, h6) {
    color: var(--el-text-color-primary);
    margin-bottom: 0.5em;
    line-height: 1.2;
    font-family: Arial, sans-serif;
}

/* 根据级别调整字体大小 */
.markdown-content :deep(h1) {
    font-size: 1.2em;
}

.markdown-content :deep(h2) {
    font-size: 1.1em;
}

.markdown-content :deep(h3) {
    font-size: 1em;
}

.markdown-content :deep(h4) {
    font-size: 0.9em;
}

.markdown-content :deep(h5) {
    font-size: 0.8em;
}

.markdown-content :deep(h6) {
    font-size: 0.7em;
}

/* 在原有的样式基础上添加 */
.markdown-content :deep(pre) {
    background-color: var(--el-fill-color-light);
    padding: 10px;
    overflow-x: auto;
    border-radius: 4px;
    margin: 10px 0;
    white-space: pre-wrap;
    word-wrap: break-word;
}

.markdown-content :deep(code) {
    font-family: 'Source Code Pro', monospace;
    font-size: 0.9em;
    color: var(--el-text-color-primary);
    padding: 2px 4px;
    background-color: var(--el-fill-color-light);
    border-radius: 4px;
}

.markdown-content :deep(table) {
    width: 80%;
    border-collapse: collapse;
    margin-bottom: 1rem;
}

.markdown-content :deep(th),
.markdown-content :deep(td) {
    border: 1px solid var(--el-border-color);
    padding: 8px;
    text-align: left;
}

.markdown-content :deep(th) {
    background-color: var(--el-fill-color-light);
    font-weight: bold;
}

.markdown-content :deep(tr:nth-child(even)) {
    background-color: var(--el-color-bg-base);
}

.markdown-content :deep(tr:hover) {
    background-color: var(--el-color-info-light-9);
}

.keyarea {
    border: 1px solid var(--el-border-color);
    margin-top: 20px;
    padding: 20px;
    text-align: left;
    width: calc(80vw - 40px);
    word-wrap: break-word;
}
</style>