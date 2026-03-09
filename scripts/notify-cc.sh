#!/bin/bash
# notify-cc.sh - 通知CC（协调者）的脚本
# 用法: ./scripts/notify-cc.sh "消息内容"

MESSAGE="$1"
if [ -z "$MESSAGE" ]; then
    echo "用法: $0 '消息内容'"
    exit 1
fi

echo "[$(date '+%Y-%m-%d %H:%M:%S')] CC-NOTIFY: $MESSAGE"
echo "$MESSAGE" > /tmp/cc-notification.txt

# 如果clawdbot可用，发送通知
if command -v clawdbot &> /dev/null; then
    clawdbot notify "$MESSAGE" 2>/dev/null || true
fi
