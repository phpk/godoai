import { ElMessage } from 'element-plus'
import '@/assets/code.scss'
import { t } from "@/i18n/index";
import ClipboardJS from 'clipboard'
import hljs from 'highlight.js'
import 'highlight.js/scss/github-dark.scss'
import MarkdownIt from 'markdown-it'
import markdownItMathjax3 from 'markdown-it-mathjax3'
// Encode text to Base64
export function textToBase64(text: string) {
  return btoa(encodeURIComponent(text))
}

// Decode Base64 to text
export function base64ToText(base64: string) {
  return decodeURIComponent(atob(base64))
}

// 代码复制
const clipboard = new ClipboardJS('.code-header-copy', {
  text: function (trigger:any) {
    const base64Str = trigger.getAttribute('data-clipboard-text-base64')
    return base64Str ? base64ToText(base64Str) : ''
  }
})
clipboard.on('success', () => {
  //notifySuccess(t('common.copySuccess'))
  ElMessage({
    type: 'success',
    message: t('common.copySuccess'),
  })
})

// MarkdownIt
const markdown = new MarkdownIt({
  highlight: (str: string, lang: string) => {
    if (!lang) {
      lang = 'text'
    }
    let codeHtml = `<code class="hljs language-${lang}">${hljs.highlight(str, { language: lang }).value
      }</code>`
    codeHtml =
      `<div class="code-header">
        <div>${lang}</div>
        <div class="code-header-copy" data-clipboard-text-base64='${textToBase64(str)}'>
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" class="icon-sm"><path fill-rule="evenodd" clip-rule="evenodd" d="M12 4C10.8954 4 10 4.89543 10 6H14C14 4.89543 13.1046 4 12 4ZM8.53513 4C9.22675 2.8044 10.5194 2 12 2C13.4806 2 14.7733 2.8044 15.4649 4H17C18.6569 4 20 5.34315 20 7V19C20 20.6569 18.6569 22 17 22H7C5.34315 22 4 20.6569 4 19V7C4 5.34315 5.34315 4 7 4H8.53513ZM8 6H7C6.44772 6 6 6.44772 6 7V19C6 19.5523 6.44772 20 7 20H17C17.5523 20 18 19.5523 18 19V7C18 6.44772 17.5523 6 17 6H16C16 7.10457 15.1046 8 14 8H10C8.89543 8 8 7.10457 8 6Z" fill="currentColor"></path></svg>
          <span>Copy</span>
        </div>
      </div>` + codeHtml
    return `<pre class="code-body-area">${codeHtml}</pre>`
  }
})

// 支持数学公式，svg渲染，无需引入额外样式
markdown.use(markdownItMathjax3)

// 渲染函数
export const renderMarkdown = (content: string, isLoading: boolean) => {
  if (!isLoading) {
    return markdown.render(content)
  }

  // 加载中，显示闪烁光标
  const endFlag = '【end】'
  let htmlCode = markdown.render(content + endFlag)
  // 找到结束标识
  const endFlagIndex = htmlCode.lastIndexOf(endFlag)
  // 插入光标元素
  htmlCode =
    htmlCode.substring(0, endFlagIndex) +
    `<span class="chat-message-loading">丨</span>` +
    htmlCode.substring(endFlagIndex + endFlag.length)
  return htmlCode
}
