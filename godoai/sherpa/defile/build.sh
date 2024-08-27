#!/bin/bash

# 根据当前系统设置OS和SUFFIX
case $(uname) in
    Linux)
        OS="linux"
        SUFFIX=""
        ;;
    Darwin)
        OS="darwin"
        SUFFIX=""
        ;;
    CYGWIN*|MINGW32*|MSYS*|MINGW*)
        OS="windows"
        SUFFIX=".exe"
        ;;
    *)
        echo "Unsupported operating system."
        exit 1
        ;;
esac

# 定义输出路径
OUT_PATH="../../deps/${OS}/voice"

# 检查目录是否存在，不存在则创建
if [ ! -d "$OUT_PATH" ]; then
    mkdir -p "$OUT_PATH"
fi

# 构建输出文件名
OUTPUT_FILE="${OUT_PATH}/non-streaming-decode-files${SUFFIX}"

# 设置GOOS和GOARCH环境变量
export GOOS=$OS
export GOARCH=$ARCH

# 编译Go程序，并处理可能的错误
if ! go build -o "$OUTPUT_FILE" ./main.go; then
    echo "编译失败，请检查错误并尝试解决。"
    exit 1
else
    echo "编译成功，生成文件: $OUTPUT_FILE"
fi