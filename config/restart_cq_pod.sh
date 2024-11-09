#!/bin/bash

# 获取环境变量中的标签信息
POD_NAME="${POD_NAME}"
NAMESPACE="${NAMESPACE}"
ENV="${ENV}"
DEPLOYMENT="${DEPLOYMENT}"

# 打印标签信息
echo "Pod Name: $POD_NAME"
echo "Namespace: $NAMESPACE"
echo "Environment: $ENV"
echo "Deployment: $DEPLOYMENT"

# 业务逻辑示例：删除指定的 Pod
kubectl delete pod "$POD_NAME" -n "$NAMESPACE" --ignore-not-found=true

# 检查命令执行结果
if [ $? -eq 0 ]; then
    echo "Pod $POD_NAME deleted successfully."
else
    echo "Failed to delete Pod $POD_NAME."
    exit 1
fi
