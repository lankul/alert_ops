#!/bin/bash

# 解析传递的 podname
for arg in "$@"
do
    case $arg in
        podname=*)
        POD_NAME="${arg#*=}"
        shift
        ;;
    esac
done

# 检查是否获取到了 podname
if [ -z "$POD_NAME" ]; then
  echo "Error: podname not provided."
  exit 1
fi

# 删除指定的 Pod
sshpass -p 'your_password' ssh -o StrictHostKeyChecking=no your_user@your_remote_host "kubectl delete pod ${POD_NAME} --ignore-not-found=true"

if [ $? -eq 0 ]; then
    echo "Pod ${POD_NAME} deleted successfully on remote host."
else
    echo "Failed to delete Pod ${POD_NAME} on remote host."
    exit 1
fi
