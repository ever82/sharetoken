// API测试脚本 - 验证真实 LLM 回复
const http = require('http');

const AGENT_GATEWAY_URL = 'http://localhost:18080';

// 检查回复是否真实（不是模拟的）
function isRealResponse(text) {
  // 模拟回复的特征
  const mockPatterns = [
    "这是来自 GenieBot 的模拟响应",
    "Mock response",
    "模拟响应"
  ];

  for (const pattern of mockPatterns) {
    if (text.includes(pattern)) {
      return false;
    }
  }
  return true;
}

// HTTP 请求
async function httpRequest(url, method = 'GET', data = null) {
  return new Promise((resolve, reject) => {
    const options = {
      hostname: 'localhost',
      port: url.includes('26657') ? 26657 : 18080,
      path: url.replace(/^http:\/\/localhost:\d+/, ''),
      method: method,
      headers: { 'Content-Type': 'application/json' }
    };

    const req = http.request(options, (res) => {
      let responseData = '';
      res.on('data', (chunk) => { responseData += chunk; });
      res.on('end', () => {
        try { resolve(JSON.parse(responseData)); } catch { resolve(responseData); }
      });
    });

    req.on('error', reject);
    if (data) req.write(JSON.stringify(data));
    req.end();
  });
}

// 主测试
async function runTests() {
  console.log('🧪 GenieBot 真实 LLM 回复测试\n');
  console.log('='.repeat(70));

  // Test 1: 发送中文消息
  console.log('\n✓ Test 1: 发送中文消息 "你好"');
  try {
    const result = await httpRequest(`${AGENT_GATEWAY_URL}/mcp`, 'POST', {
      jsonrpc: '2.0',
      method: 'tools/call',
      id: 1,
      params: {
        name: 'chat_with_genie',
        arguments: { message: '你好', session_id: 'test-chinese' }
      }
    });

    if (result.result?.content?.[0]?.text) {
      const data = JSON.parse(result.result.content[0].text);
      console.log('  📨 用户: 你好');
      console.log('  🤖 GenieBot:', data.content);
      console.log('  💰 费用:', data.cost, 'STT');

      if (isRealResponse(data.content)) {
        console.log('  ✅ 检测到真实 LLM 回复！');
      } else {
        console.log('  ⚠️  当前是模拟回复（未配置 ANTHROPIC_API_KEY）');
        console.log('  💡 设置环境变量: export ANTHROPIC_API_KEY=your_key');
      }
    }
  } catch (e) {
    console.log('  ❌ 错误:', e.message);
  }

  // Test 2: 发送具体问题
  console.log('\n✓ Test 2: 发送具体问题');
  try {
    const result = await httpRequest(`${AGENT_GATEWAY_URL}/mcp`, 'POST', {
      jsonrpc: '2.0',
      method: 'tools/call',
      id: 2,
      params: {
        name: 'chat_with_genie',
        arguments: { message: '什么是区块链技术？', session_id: 'test-blockchain' }
      }
    });

    if (result.result?.content?.[0]?.text) {
      const data = JSON.parse(result.result.content[0].text);
      console.log('  📨 用户: 什么是区块链技术？');
      console.log('  🤖 GenieBot:', data.content.substring(0, 100) + '...');
      console.log('  💰 费用:', data.cost, 'STT');

      if (isRealResponse(data.content)) {
        console.log('  ✅ 真实 AI 回复！');
      } else {
        console.log('  ⚠️  模拟回复');
      }
    }
  } catch (e) {
    console.log('  ❌ 错误:', e.message);
  }

  // Test 3: 检查 API Key 配置
  console.log('\n✓ Test 3: 检查 API Key 配置');
  const apiKey = process.env.ANTHROPIC_API_KEY;
  if (apiKey) {
    console.log('  ✅ ANTHROPIC_API_KEY 已设置');
    console.log('  📊 Key 前缀:', apiKey.substring(0, 10) + '...');
  } else {
    console.log('  ⚠️  ANTHROPIC_API_KEY 未设置');
    console.log('  💡 如需真实 LLM 回复，请设置:');
    console.log('     export ANTHROPIC_API_KEY=sk-ant-api...');
  }

  console.log('\n' + '='.repeat(70));
  console.log('💡 使用真实 LLM 的方法:');
  console.log('  1. 获取 Claude API Key: https://console.anthropic.com/');
  console.log('  2. 设置环境变量: export ANTHROPIC_API_KEY=your_key');
  console.log('  3. 重启 Agent Gateway');
  console.log('  4. 再次测试');
}

runTests().catch(console.error);
