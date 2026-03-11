// API测试脚本 - 测试 GenieBot 后端功能
const http = require('http');

const AGENT_GATEWAY_URL = 'http://localhost:18080';
const NODE_RPC_URL = 'http://localhost:26657';

// 简单的 HTTP 请求函数
async function httpRequest(url, method = 'GET', data = null) {
  return new Promise((resolve, reject) => {
    const options = {
      hostname: 'localhost',
      port: url.includes('26657') ? 26657 : 18080,
      path: url.replace(/^http:\/\/localhost:\d+/, ''),
      method: method,
      headers: {
        'Content-Type': 'application/json'
      }
    };

    const req = http.request(options, (res) => {
      let responseData = '';
      res.on('data', (chunk) => {
        responseData += chunk;
      });
      res.on('end', () => {
        try {
          resolve(JSON.parse(responseData));
        } catch {
          resolve(responseData);
        }
      });
    });

    req.on('error', reject);

    if (data) {
      req.write(JSON.stringify(data));
    }
    req.end();
  });
}

// 测试套件
async function runTests() {
  console.log('🧪 开始 GenieBot API 测试\n');
  console.log('=' .repeat(60));

  let passCount = 0;
  let failCount = 0;

  // Test 1: Agent Gateway Health
  console.log('\n✓ Test 1: Agent Gateway Health Check');
  try {
    const result = await httpRequest(`${AGENT_GATEWAY_URL}/health`);
    if (result.status === 'ok') {
      console.log('  ✅ Agent Gateway 运行正常');
      console.log('  📊 版本:', result.version);
      passCount++;
    } else {
      throw new Error('Health check failed');
    }
  } catch (e) {
    console.log('  ❌ 测试失败:', e.message);
    failCount++;
  }

  // Test 2: A2A Agent Card
  console.log('\n✓ Test 2: A2A Agent Card');
  try {
    const result = await httpRequest(`${AGENT_GATEWAY_URL}/.well-known/agent.json`);
    if (result.name === 'ShareToken Agent') {
      console.log('  ✅ Agent Card 正常');
      console.log('  📊 名称:', result.name);
      console.log('  📊 能力:', result.capabilities?.join(', '));
      passCount++;
    } else {
      throw new Error('Invalid agent card');
    }
  } catch (e) {
    console.log('  ❌ 测试失败:', e.message);
    failCount++;
  }

  // Test 3: MCP Tools List
  console.log('\n✓ Test 3: MCP Tools List');
  try {
    const result = await httpRequest(`${AGENT_GATEWAY_URL}/mcp`, 'POST', {
      jsonrpc: '2.0',
      method: 'tools/list',
      id: 1
    });
    if (result.result?.tools?.length > 0) {
      console.log('  ✅ MCP Tools 可用');
      console.log('  📊 可用工具:', result.result.tools.map(t => t.name).join(', '));
      passCount++;
    } else {
      throw new Error('No tools available');
    }
  } catch (e) {
    console.log('  ❌ 测试失败:', e.message);
    failCount++;
  }

  // Test 4: Query Balance
  console.log('\n✓ Test 4: Query Balance');
  try {
    const result = await httpRequest(`${AGENT_GATEWAY_URL}/mcp`, 'POST', {
      jsonrpc: '2.0',
      method: 'tools/call',
      id: 1,
      params: {
        name: 'query_balance',
        arguments: { address: 'cosmos1test' }
      }
    });
    if (result.result?.content?.[0]?.text) {
      const data = JSON.parse(result.result.content[0].text);
      console.log('  ✅ 余额查询成功');
      console.log('  📊 地址:', data.address);
      console.log('  📊 余额:', data.balance, data.denom);
      passCount++;
    } else {
      throw new Error('Balance query failed');
    }
  } catch (e) {
    console.log('  ❌ 测试失败:', e.message);
    failCount++;
  }

  // Test 5: Chat with GenieBot
  console.log('\n✓ Test 5: Chat with GenieBot');
  try {
    const result = await httpRequest(`${AGENT_GATEWAY_URL}/mcp`, 'POST', {
      jsonrpc: '2.0',
      method: 'tools/call',
      id: 1,
      params: {
        name: 'chat_with_genie',
        arguments: { message: '你好', session_id: 'test-session' }
      }
    });
    if (result.result?.content?.[0]?.text) {
      const data = JSON.parse(result.result.content[0].text);
      console.log('  ✅ GenieBot 对话成功');
      console.log('  📊 回复:', data.content.substring(0, 50) + '...');
      console.log('  📊 费用:', data.cost, 'STT');
      passCount++;
    } else {
      throw new Error('Chat failed');
    }
  } catch (e) {
    console.log('  ❌ 测试失败:', e.message);
    failCount++;
  }

  // Test 6: Create Task
  console.log('\n✓ Test 6: Create Task');
  try {
    const result = await httpRequest(`${AGENT_GATEWAY_URL}/mcp`, 'POST', {
      jsonrpc: '2.0',
      method: 'tools/call',
      id: 1,
      params: {
        name: 'create_task',
        arguments: { description: '测试任务', budget: '100stt' }
      }
    });
    if (result.result?.content?.[0]?.text) {
      const data = JSON.parse(result.result.content[0].text);
      console.log('  ✅ 任务创建成功');
      console.log('  📊 任务ID:', data.task_id);
      console.log('  📊 预算:', data.budget);
      passCount++;
    } else {
      throw new Error('Task creation failed');
    }
  } catch (e) {
    console.log('  ❌ 测试失败:', e.message);
    failCount++;
  }

  // Test 7: A2A Tasks List
  console.log('\n✓ Test 7: A2A Tasks List');
  try {
    const result = await httpRequest(`${AGENT_GATEWAY_URL}/a2a/tasks`);
    if (result.tasks && Array.isArray(result.tasks)) {
      console.log('  ✅ 任务列表获取成功');
      console.log('  📊 任务数量:', result.tasks.length);
      passCount++;
    } else {
      throw new Error('Tasks list failed');
    }
  } catch (e) {
    console.log('  ❌ 测试失败:', e.message);
    failCount++;
  }

  // Test 8: A2A Status
  console.log('\n✓ Test 8: A2A Status');
  try {
    const result = await httpRequest(`${AGENT_GATEWAY_URL}/a2a/status`);
    if (result.status === 'online') {
      console.log('  ✅ Agent Gateway 状态正常');
      console.log('  📊 状态:', result.status);
      console.log('  📊 版本:', result.version);
      passCount++;
    } else {
      throw new Error('Status check failed');
    }
  } catch (e) {
    console.log('  ❌ 测试失败:', e.message);
    failCount++;
  }

  // Summary
  console.log('\n' + '='.repeat(60));
  console.log('📊 测试结果汇总');
  console.log('='.repeat(60));
  console.log(`✅ 通过: ${passCount}`);
  console.log(`❌ 失败: ${failCount}`);
  console.log(`📈 成功率: ${Math.round((passCount / (passCount + failCount)) * 100)}%`);

  if (failCount === 0) {
    console.log('\n🎉 所有测试通过！GenieBot 运行正常。');
  } else {
    console.log('\n⚠️ 部分测试失败，请检查服务状态。');
  }
}

// Run tests
runTests().catch(console.error);
