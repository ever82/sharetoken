#!/usr/bin/env python3
"""
GenieBot 本地体验客户端
简单交互式对话工具
"""

import requests
import json
import sys

def chat_with_genie(message, session_id="demo-session"):
    """与 GenieBot 对话"""
    url = "http://localhost:8080/mcp"
    payload = {
        "jsonrpc": "2.0",
        "method": "tools/call",
        "id": 1,
        "params": {
            "name": "chat_with_genie",
            "arguments": {
                "message": message,
                "session_id": session_id
            }
        }
    }

    try:
        resp = requests.post(url, json=payload, timeout=10)
        data = resp.json()

        if "result" in data and "content" in data["result"]:
            text = data["result"]["content"][0]["text"]
            result = json.loads(text)
            return result["content"], result["cost"]
        elif "error" in data:
            return f"Error: {data['error']['message']}", 0
        else:
            return "Unknown response format", 0
    except Exception as e:
        return f"Request failed: {e}", 0

def query_balance(address):
    """查询余额"""
    url = "http://localhost:8080/mcp"
    payload = {
        "jsonrpc": "2.0",
        "method": "tools/call",
        "id": 1,
        "params": {
            "name": "query_balance",
            "arguments": {"address": address}
        }
    }

    try:
        resp = requests.post(url, json=payload, timeout=5)
        data = resp.json()
        text = data["result"]["content"][0]["text"]
        result = json.loads(text)
        return result["balance"], result["denom"]
    except Exception as e:
        return None, str(e)

def create_task(description, budget="100stt"):
    """创建任务"""
    url = "http://localhost:8080/mcp"
    payload = {
        "jsonrpc": "2.0",
        "method": "tools/call",
        "id": 1,
        "params": {
            "name": "create_task",
            "arguments": {
                "description": description,
                "budget": budget
            }
        }
    }

    try:
        resp = requests.post(url, json=payload, timeout=5)
        data = resp.json()
        text = data["result"]["content"][0]["text"]
        result = json.loads(text)
        return result["task_id"], result["status"]
    except Exception as e:
        return None, str(e)

def main():
    print("🧞 欢迎来到 GenieBot 本地体验!")
    print("")

    # 检查服务是否运行
    try:
        resp = requests.get("http://localhost:8080/health", timeout=2)
        print("✅ 服务状态:", resp.json()["status"])
        print("")
    except:
        print("❌ 服务未启动，请先运行:")
        print("   ./bin/agent-gateway -transport=http -port=8080")
        print("")
        print("或者使用便捷脚本:")
        print("   ./scripts/demo_geniebot.sh")
        sys.exit(1)

    print("可用命令:")
    print("  /balance <地址>  - 查询账户余额")
    print("  /task <描述>     - 创建新任务")
    print("  /quit            - 退出")
    print("")
    print("直接输入文字即可与 GenieBot 对话")
    print("-" * 50)
    print("")

    session_id = "demo-session"
    total_cost = 0.0

    while True:
        try:
            user_input = input("💬 你: ").strip()
        except EOFError:
            break
        except KeyboardInterrupt:
            print("\n")
            break

        if not user_input:
            continue

        if user_input == "/quit" or user_input == "/exit":
            break

        if user_input.startswith("/balance "):
            address = user_input[9:].strip()
            if address:
                balance, denom = query_balance(address)
                if balance:
                    print(f"💰 余额: {balance} {denom}")
                else:
                    print(f"❌ 查询失败: {denom}")
            else:
                print("⚠️  请提供地址: /balance cosmos1...")
            continue

        if user_input.startswith("/task "):
            desc = user_input[6:].strip()
            if desc:
                task_id, status = create_task(desc)
                if task_id:
                    print(f"📋 任务已创建: {task_id} (状态: {status})")
                else:
                    print(f"❌ 创建失败: {status}")
            else:
                print("⚠️  请提供描述: /task 描述内容")
            continue

        # 与 GenieBot 对话
        print("🧞 GenieBot 思考中...")
        response, cost = chat_with_genie(user_input, session_id)
        total_cost += cost

        print(f"🧞 GenieBot: {response}")
        print(f"   💰 本次费用: {cost} STT | 累计: {total_cost:.4f} STT")
        print("")

    print("")
    print("-" * 50)
    print(f"👋 再见! 本次会话总费用: {total_cost:.4f} STT")

if __name__ == "__main__":
    main()
