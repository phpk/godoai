## API 文档 - 下载模型

### 描述
此 API 用于下载指定的模型文件，并支持多种来源，如本地路径或远程 URL。它还支持跟踪下载进度。

### 请求 URL
- `/model/download`

### 请求方法
- `POST`

### 请求参数
- **Content-Type**: `application/json`
- **Body**:
  - **`Url`**: `[]string` - 模型文件的 URL 或本地路径列表。
  - **`Model`**: `string` - 模型名称。
  - **`Type`**: `string` - 模型类型，例如 `ollama` 或 `llm`。
  - **`From`**: `string` (可选) - 指定模型来源，默认为 `ollama`。

### 响应
- **200 OK**: 成功响应，包含下载状态和相关信息。
- **其他 HTTP 状态码**: 错误响应，包含错误信息。

### 示例
```javascript
const request = {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    Url: ['https://example.com/model.bin', '/path/to/local/model.bin'],
    Model: 'my-model',
    Type: 'ollama',
    From: 'ollama' // 可选，默认为 'ollama'
  })
};

fetch('/model/download', request)
  .then(response => response.json())
  .then(data => console.log(data))
  .catch(error => console.error('Error:', error));
```

## 注意事项
- 如果模型已经存在，则会直接返回成功消息而不进行下载。
- 如果 `From` 参数未提供，则默认值为 `ollama`。
- 下载过程中会跟踪进度，并在完成后更新模型信息。

## 返回数据结构
- **`Paths`**: `[]string` - 下载的文件路径列表。
- **`Status`**: `string` - 下载状态，例如 `"success"`。
- **`CreatedAt`**: `time.Time` - 创建时间。
- **`Info`**: `map[string]interface{}` - 其他信息，例如总大小 (`tsize`) 和人类可读的大小 (`size`)。

## 错误处理
- **400 Bad Request**: 请求体解析失败或请求参数无效。
- **500 Internal Server Error**: 内部服务器错误，例如配置加载失败、文件路径获取失败、下载失败等。

## 示例响应
```json
{
  "Paths": ["/path/to/downloaded/model.bin"],
  "Status": "success",
  "CreatedAt": "2023-04-01T12:00:00Z",
  "Info": {
    "tsize": 104857600,
    "size": "100 MB"
  }
}
```

