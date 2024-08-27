

# 知识库管理
## 描述
- 用于管理知识库，包括创建、列出、删除知识库，向知识库添加文件、通过 URL 添加文件、查询知识库、上传文件、展示文件详情、显示图片以及转换文件。

## 请求 URL
- /knowledge/list
- /knowledge/create
- /knowledge/delete
- /knowledge/deleteFile
- /knowledge/ask
- /knowledge/add
- /knowledge/upload
- /knowledge/url
- /knowledge/filedetail
- /knowledge/showimage

1. 知识库列表
### 请求 URL
- /knowledge/list
### 请求方法
- GET
### 请求参数
- 无
### 示例
```
fetch('/knowledge/list', {
  method: 'GET',
})
  .then(response => response.json())
  .then(data => console.log(data))
  .catch(error => console.error('Error:', error));
```

2. 创建知识库
## 请求 URL
- /knowledge/create
## 请求方法
- POST
## 请求参数
- name: 知识库名称
- description: 知识库描述
## 示例
```
const request = {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    name: 'my-knowledge-base',
    model: 'gemma:2b',
    config: {
      type: 'local',
      apiUrl: 'https://api.example.com',
      apiKey: 'your-api-key',
      embedding: {
        apiUrl: 'https://embeddings.api.example.com',
        apiKey: 'your-embeddings-api-key',
        apiType: 'openai',
        contextLength: 1024
      }
    }
  })
};

fetch('/knowledge/create', request)
  .then(response => response.json())
  .then(data => console.log(data))
  .catch(error => console.error('Error:', error));
```

3. 删除知识库
## 请求 URL
- /knowledge/delete
## 请求方法
- POST
## 请求参数
- name: 要删除的知识库名称
## 示例
```
const request = {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    name: 'my-knowledge-base'
  })
};

fetch('/knowledge/delete', request)
  .then(response => response.json())
  .then(data => console.log(data))
  .catch(error => console.error('Error:', error));
```
4. 删除文件
## 请求 URL
- /knowledge/deleteFile
## 请求方法
- POST
## 请求参数
- name: 要删除的文件所在的知识库名称
- fileName: 要删除的文件名称
## 示例
```
const request = {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    name: 'my-knowledge-base',
    model: 'gpt-3.5-turbo',
    file: 'example.txt',
    config: {
      type: 'local',
      apiUrl: 'https://api.example.com',
      apiKey: 'your-api-key',
      embedding: {
        apiUrl: 'https://embeddings.api.example.com',
        apiKey: 'your-embeddings-api-key',
        apiType: 'openai',
        contextLength: 1024
      }
    }
  })
};

fetch('/knowledge/deleteFile', request)
  .then(response => response.json())
  .then(data => console.log(data))
  .catch(error => console.error('Error:', error));
```
5. 查询知识库
## 请求 URL
- /knowledge/ask
## 请求方法
- POST
## 请求参数
- name: 要查询的知识库名称
- question: 要查询的问题
## 示例
```
const request = {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    name: 'my-knowledge-base',
    model: 'gpt-3.5-turbo',
    message: 'What is the capital of France?',
    config: {
      type: 'local',
      apiUrl: 'https://api.example.com',
      apiKey: 'your-api-key',
      embedding: {
        apiUrl: 'https://embeddings.api.example.com',
        apiKey: 'your-embeddings-api-key',
        apiType: 'openai',
        contextLength: 1024
      }
    }
  })
};

fetch('/knowledge/ask', request)
  .then(response => response.json())
  .then(data => console.log(data))
  .catch(error => console.error('Error:', error));
```
6. 向知识库添加文件
## 请求 URL
- /knowledge/add
## 请求方法
- POST
## 请求参数
- name: 要添加文件的知识库名称
- file: 要添加的文件
## 示例
```
const request = {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    name: 'my-knowledge-base',
    model: 'gpt-3.5-turbo',
    files: ['example.txt', 'example.pdf'],
    config: {
      type: 'local',
      apiUrl: 'https://api.example.com',
      apiKey: 'your-api-key',
      embedding: {
        apiUrl: 'https://embeddings.api.example.com',
        apiKey: 'your-embeddings-api-key',
        apiType: 'openai',
        contextLength: 1024
      }
    }
  })
};

fetch('/knowledge/add', request)
  .then(response => response.json())
  .then(data => console.log(data))
  .catch(error => console.error('Error:', error));
```
7. 通过 URL 添加文件
## 请求 URL
- /knowledge/url
## 请求方法
- POST
## 请求参数
- name: 要添加文件的知识库名称
- url: 要添加文件的 URL
## 示例
```
const formData = new FormData();
formData.append('file', fileInput.files[0]);

fetch('/knowledge/upload', {
  method: 'POST',
  body: formData,
})
  .then(response => response.json())
  .then(data => console.log(data))
  .catch(error => console.error('Error:', error));
```
8. 获取文件详情
## 请求 URL
- /knowledge/filedetail
## 请求方法
- POST
## 请求参数
- name: 要获取文件详情的知识库名称
- fileName: 要获取文件详情的文件名称
## 示例
```
const request = {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    url: 'https://example.com/document.pdf'
  })
};

fetch('/knowledge/url', request)
  .then(response => response.json())
  .then(data => console.log(data))
  .catch(error => console.error('Error:', error));
```
9. 展示文件详情
## 请求 URL
- /knowledge/filedetail
## 请求方法
- GET
## 请求参数
- name: 要展示文件详情的知识库名称
- fileName: 要展示文件详情的文件名称
## 示例
```
fetch(`/knowledge/filedetail?filename=example.txt`)
  .then(response => response.json())
  .then(data => console.log(data))
  .catch(error => console.error('Error:', error));
```

10. 显示图片
## 请求 URL
- /knowledge/showimage
## 请求方法
- GET
## 请求参数
- name: 要展示图片的知识库名称
- fileName: 要展示图片的文件名称
## 示例
```
fetch(`/knowledge/showimage?filename=image.jpg`)
  .then(response => response.blob())
  .then(blob => {
    const imageUrl = URL.createObjectURL(blob);
    const imgElement = document.createElement('img');
    imgElement.src = imageUrl;
    document.body.appendChild(imgElement);
  })
  .catch(error => console.error('Error:', error));

```


