# 系统设置

## - `/system/setting`

### 描述
该 API 用于更新系统的配置设置，支持批量设置多个配置项。当前支持的配置项包括 `dataDir` 和 `ipList`。

### 请求方法
- **POST**

### 请求 URL
- **/system/setting**

### 请求参数
- **Content-Type**: `application/json`
- **Body**:
  - **dataDir** (可选): 数据目录的路径。类型为字符串。
  - **ipList** (可选): 允许访问的 IP 地址列表。类型为字符串数组。

### 响应
- **200 OK**: 成功设置配置。
- **400 Bad Request**: 参数错误或无效。

### 示例
#### JavaScript 使用示例
```javascript
const request = {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify([
    { name: 'dataDir', value: '/path/to/new/data/directory' },
    { name: 'ipList', value: ['192.168.1.100', '192.168.1.101'] }
  ])
};

fetch('/system/setting', request)
  .then(response => response.json())
  .then(data => console.log(data))
  .catch(error => console.error('Error:', error));
```
### 参数说明
- dataDir: 设置数据目录的路径。如果指定的新路径不存在，则会尝试创建该目录。
- ipList: 设置允许访问的 IP 地址列表。列表中的 IP 地址必须是有效的 CIDR 格式。
### 返回值
- status: 状态消息。
- message: 操作成功或失败的信息。
### 注意事项
- 如果请求体中包含多个配置项，则这些配置项将被依次处理。
- 对于 dataDir 配置项，如果新路径与现有路径不同，则会尝试创建新路径，并确保其具有正确的权限。
- 对于 ipList 配置项，所有 IP 地址必须是有效的 CIDR 格式。
- 成功设置配置后，系统会重启相关服务以应用新的配置。




## 下载模型

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
    url: ['https://example.com/model.bin'],//可为空，如果是外部源则必填
    model: 'my-model',
    type: 'ollama',
    from: 'ollama' // 可选，默认为 'ollama'
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

## 下载服务器模型文件

### 描述
此 API 用于从服务器下载指定的模型文件。客户端可以通过 GET 请求并提供文件路径来触发文件下载。

### 请求 URL
- `/model/outserver`

### 请求方法
- `GET`

### 请求参数
- **`path`**: `string` - 查询字符串中的文件路径。

### 响应
- **200 OK**: 文件下载成功。
- **400 Bad Request**: 请求的文件路径无效。
- **404 Not Found**: 请求的文件不存在。
- **500 Internal Server Error**: 服务器内部错误。

### 响应头
- **`Content-Disposition`**: `attachment; filename="filename"` - 指示浏览器以附件形式下载文件。
- **`Content-Type`**: `application/octet-stream` - 文件类型。
- **`Content-Length`**: `integer` - 文件大小。

### 示例
#### 请求
```bash
GET /model/outserver?path=/path/to/file.txt
```

## 示例响应

### 成功
- **HTTP 状态码**: 200
- **响应头**:
  - **`Content-Disposition`**: `attachment; filename="file.txt"`
  - **`Content-Type`**: `application/octet-stream`
  - **`Content-Length`**: `12345` (假设文件大小为 12345 字节)

### 错误
- **HTTP 状态码**: 400
- **响应体**:
  ```json
  {
    "error": "Invalid file path"
  }
  ```
- **HTTP 状态码**: 404

- **响应体**:

```json
{
  "error": "File not found"
}
```
- **HTTP 状态码**: 500

- **响应体**:

```json
{
  "error": "Internal Server Error"
}
```

## 注意事项
- 请求的文件路径必须是有效的服务器文件路径。
- 如果文件不存在或路径无效，将返回相应的错误状态码。
- 服务器会设置适当的响应头以确保文件被正确地作为附件下载。
## 示例代码
```javascript
const url = 'http://example.com/model/outserver?path=/path/to/file.txt';

fetch(url)
  .then(response => {
    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }
    return response.blob();
  })
  .then(blob => {
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'file.txt';
    a.click();
  })
  .catch(error => console.error('Error:', error))
  ```

  