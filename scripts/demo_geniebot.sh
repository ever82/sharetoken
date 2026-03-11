#!/bin/bash
# GenieBot 本地体验脚本

echo "🧞 启动 GenieBot 本地服务..."
echo ""

# 启动后台服务
./bin/agent-gateway -transport=http -port=8080 &
PID=$!

# 等待服务启动
sleep 1

echo "✅ 服务已启动!"
echo ""
echo "📍 可用端点:"
echo "   - Health:    http://localhost:8080/health"
echo "   - Agent Card: http://localhost:8080/.well-known/agent.json"
echo "   - MCP API:   http://localhost:8080/mcp"
echo "   - A2A Tasks: http://localhost:8080/a2a/tasks"
echo ""
echo "💬 体验 GenieBot 对话:"
echo ""

# 测试对话
curl -s -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "id": 1,
    "params": {
      "name": "chat_with_genie",
      "arguments": {
        "message": "你好，请介绍一下 ShareToken"
      }
    }
  }' | python3 -m json.tool 2>/dev/null || cat

echo ""
echo ""
echo "🎮 交互式体验 (输入 'quit' 退出):"
echo ""

while true; do
    read -p "💬 你对 GenieBot 说: " message
    if [ "$message" = "quit" ] || [ "$message" = "exit" ]; then
        break
    fi

    # Escape special characters
    escaped_msg=$(echo "$message" | sed 's/"/\\"/g')

    curl -s -X POST http://localhost:8080/mcp \
      -H "Content-Type: application/json" \
      -d "{
        \"jsonrpc\": \"2.0\",
        \"method\": \"tools/call\",
        \"id\": 1,
        \"params\": {
          \"name\": \"chat_with_genie\",
          \"arguments\": {
            \"message\": \"$escaped_msg\"
          }
        }
      }" | python3 -c "import json,sys; d=json.load(sys.stdin); print('\n🧞 GenieBot:', json.loads(d['result']['content'][0]['text'])['content'], '\n💰 费用:', json.loads(d['result']['content'][0]['text'])['cost'], 'STT\n')" 2>/dev/null || echo "响应解析失败"
done

echo ""
echo "👋 关闭服务..."
kill $PID 2>/dev/null
echo "✅ 服务已停止"
