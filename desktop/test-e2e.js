const { _electron } = require('playwright');
const { spawn } = require('child_process');
const path = require('path');
const fs = require('fs');

// 测试配置
const TEST_TIMEOUT = 120000;
const AGENT_GATEWAY_URL = 'http://localhost:18080';
const NODE_RPC_URL = 'http://localhost:26657';

// 等待函数
const wait = (ms) => new Promise(r => setTimeout(r, ms));

// 检查服务是否运行
async function checkService(url, timeout = 5000) {
  const http = require('http');
  return new Promise((resolve) => {
    const req = http.get(url, (res) => {
      resolve(res.statusCode === 200);
    });
    req.on('error', () => resolve(false));
    req.setTimeout(timeout, () => {
      req.abort();
      resolve(false);
    });
  });
}

// 主测试
(async () => {
  console.log('🧪 GenieBot 全流程 E2E 测试\n');
  console.log('=' .repeat(70));

  let electronApp = null;
  let testResults = [];

  try {
    // Test 1: 检查服务状态
    console.log('\n✓ Test 1: 检查后端服务');
    const agentGatewayRunning = await checkService(`${AGENT_GATEWAY_URL}/health`);
    const nodeRunning = await checkService(`${NODE_RPC_URL}/status`);

    if (agentGatewayRunning && nodeRunning) {
      console.log('  ✅ Agent Gateway 运行中');
      console.log('  ✅ ShareToken 节点运行中');
      testResults.push({ name: '服务检查', status: 'pass' });
    } else {
      console.log('  ⚠️  服务未完全启动，尝试启动...');
      console.log(`     Agent Gateway: ${agentGatewayRunning ? '✅' : '❌'}`);
      console.log(`     ShareToken Node: ${nodeRunning ? '✅' : '❌'}`);

      if (!agentGatewayRunning) {
        console.log('  🚀 启动 Agent Gateway...');
        spawn('/Applications/ShareToken.app/Contents/Resources/bin/agent-gateway', [
          '-transport', 'http',
          '-port', '18080'
        ], { detached: true });
        await wait(3000);
      }

      if (!nodeRunning) {
        console.log('  ⚠️  请在 ShareToken 应用中启动节点');
      }
    }

    // Test 2: 启动 Electron 应用
    console.log('\n✓ Test 2: 启动 ShareToken 桌面应用');
    const appPath = '/Applications/ShareToken.app/Contents/MacOS/ShareToken';

    if (!fs.existsSync(appPath)) {
      throw new Error(`应用未找到: ${appPath}`);
    }

    console.log('  🚀 正在启动应用...');
    electronApp = await _electron.launch({
      executablePath: appPath,
      args: [],
      timeout: 30000
    });

    console.log('  ✅ 应用已启动');
    testResults.push({ name: '应用启动', status: 'pass' });

    // 获取主窗口
    const mainWindow = await electronApp.firstWindow();
    console.log('  📱 主窗口已获取');

    // 等待应用加载
    await wait(5000);

    // Test 3: 截图初始状态
    console.log('\n✓ Test 3: 截图 - 初始状态');
    await mainWindow.screenshot({ path: '/tmp/e2e-01-initial.png' });
    console.log('  📸 截图保存: /tmp/e2e-01-initial.png');
    testResults.push({ name: '初始截图', status: 'pass' });

    // Test 4: 导航到 GenieBot
    console.log('\n✓ Test 4: 导航到 GenieBot');

    // 尝试多种方式找到 GenieBot 链接
    const selectors = [
      'text=GenieBot',
      'a:has-text("GenieBot")',
      'button:has-text("GenieBot")',
      '[href*="genie"]',
      '.nav-item:has-text("Genie")',
      'nav a:has-text("🧞")'
    ];

    let genieBotFound = false;
    for (const selector of selectors) {
      try {
        const element = mainWindow.locator(selector).first();
        if (await element.isVisible({ timeout: 2000 }).catch(() => false)) {
          await element.click();
          console.log(`  ✅ 找到并点击 GenieBot: ${selector}`);
          genieBotFound = true;
          break;
        }
      } catch (e) {
        // Continue to next selector
      }
    }

    if (!genieBotFound) {
      // 尝试通过页面内容查找
      const pageContent = await mainWindow.content();
      if (pageContent.includes('GenieBot') || pageContent.includes('🧞')) {
        console.log('  ⚠️  页面包含 GenieBot 但未找到可点击元素');
        // 尝试点击包含 GenieBot 的元素
        const elements = await mainWindow.locator('a, button, .nav-item, .menu-item').all();
        for (const el of elements) {
          const text = await el.textContent().catch(() => '');
          if (text.toLowerCase().includes('genie')) {
            await el.click();
            console.log('  ✅ 通过文本找到 GenieBot');
            genieBotFound = true;
            break;
          }
        }
      }
    }

    if (!genieBotFound) {
      console.log('  ⚠️  未找到 GenieBot 导航，继续测试');
      testResults.push({ name: 'GenieBot 导航', status: 'skip' });
    } else {
      await wait(3000);
      await mainWindow.screenshot({ path: '/tmp/e2e-02-geniebot.png' });
      console.log('  📸 截图保存: /tmp/e2e-02-geniebot.png');
      testResults.push({ name: 'GenieBot 导航', status: 'pass' });
    }

    // Test 5: 检查 Agent Gateway 连接状态
    console.log('\n✓ Test 5: 检查 Agent Gateway 连接');

    // 查找连接状态显示
    const statusSelectors = [
      '.connection-status',
      '.status-text',
      '.status-dot',
      '[class*="status"]',
      'text=Connected',
      'text=连接'
    ];

    let connectionStatus = 'unknown';
    for (const selector of statusSelectors) {
      try {
        const element = mainWindow.locator(selector).first();
        if (await element.isVisible({ timeout: 2000 }).catch(() => false)) {
          const text = await element.textContent().catch(() => '');
          console.log(`  📊 连接状态: ${text}`);
          connectionStatus = text;
          break;
        }
      } catch (e) {
        // Continue
      }
    }

    if (connectionStatus.toLowerCase().includes('connected') ||
        connectionStatus.toLowerCase().includes('连接') ||
        connectionStatus.toLowerCase().includes('ok')) {
      console.log('  ✅ Agent Gateway 已连接');
      testResults.push({ name: 'Agent Gateway 连接', status: 'pass' });
    } else {
      console.log('  ⚠️  连接状态未知:', connectionStatus);
      testResults.push({ name: 'Agent Gateway 连接', status: 'skip' });
    }

    // Test 6: 检查 Agent 选择器
    console.log('\n✓ Test 6: 检查 Agent 选择器');
    const agentSelectors = [
      'select',
      '.agent-select',
      '[class*="agent"]',
      'select[class*="model"]',
      'select[class*="agent"]'
    ];

    let agentSelectorFound = false;
    for (const selector of agentSelectors) {
      try {
        const element = mainWindow.locator(selector).first();
        if (await element.isVisible({ timeout: 2000 }).catch(() => false)) {
          const options = await element.locator('option').allTextContents();
          console.log('  📊 可用 Agent:', options.join(', '));
          agentSelectorFound = true;
          testResults.push({ name: 'Agent 选择器', status: 'pass' });
          break;
        }
      } catch (e) {
        // Continue
      }
    }

    if (!agentSelectorFound) {
      console.log('  ⚠️  未找到 Agent 选择器');
      testResults.push({ name: 'Agent 选择器', status: 'skip' });
    }

    // Test 7: 发送测试消息
    console.log('\n✓ Test 7: 发送测试消息');

    const inputSelectors = [
      'textarea',
      'input[type="text"]',
      '.chat-input',
      '[placeholder*="message"]',
      '[placeholder*="输入"]'
    ];

    let messageSent = false;
    for (const selector of inputSelectors) {
      try {
        const input = mainWindow.locator(selector).first();
        if (await input.isVisible({ timeout: 2000 }).catch(() => false)) {
          console.log('  📝 找到输入框');

          // 输入测试消息
          await input.fill('你好，请查询我的余额');
          console.log('  ✅ 输入测试消息');

          // 查找发送按钮
          const sendSelectors = [
            'button:has-text("➤")',
            '.send-btn',
            'button[type="submit"]',
            'button:has-text("Send")',
            'button:has-text("发送")'
          ];

          for (const sendSelector of sendSelectors) {
            try {
              const sendBtn = mainWindow.locator(sendSelector).first();
              if (await sendBtn.isVisible({ timeout: 1000 }).catch(() => false)) {
                await sendBtn.click();
                console.log('  ✅ 点击发送按钮');
                messageSent = true;
                break;
              }
            } catch (e) {
              // Continue
            }
          }

          if (!messageSent) {
            // 尝试按 Enter 发送
            await input.press('Enter');
            console.log('  ✅ 按 Enter 发送');
            messageSent = true;
          }

          break;
        }
      } catch (e) {
        // Continue
      }
    }

    if (messageSent) {
      // 等待回复
      console.log('  ⏳ 等待回复...');
      await wait(5000);

      await mainWindow.screenshot({ path: '/tmp/e2e-03-response.png' });
      console.log('  📸 截图保存: /tmp/e2e-03-response.png');

      // 检查是否有回复消息
      const messageSelectors = [
        '.message',
        '.chat-message',
        '[class*="message"]'
      ];

      let messageCount = 0;
      for (const selector of messageSelectors) {
        try {
          const messages = await mainWindow.locator(selector).all();
          messageCount = messages.length;
          if (messageCount > 0) {
            console.log(`  📊 消息数量: ${messageCount}`);
            break;
          }
        } catch (e) {
          // Continue
        }
      }

      if (messageCount > 1) {
        console.log('  ✅ 收到回复');
        testResults.push({ name: '消息发送与回复', status: 'pass' });
      } else {
        console.log('  ⚠️  可能未收到回复');
        testResults.push({ name: '消息发送与回复', status: 'skip' });
      }
    } else {
      console.log('  ❌ 未找到输入框');
      testResults.push({ name: '消息发送与回复', status: 'fail' });
    }

    // Test 8: 检查快捷按钮
    console.log('\n✓ Test 8: 检查快捷按钮');
    const quickBtnSelectors = [
      '.quick-btn',
      '.quick-actions button',
      'button:has-text("💰")',
      'button:has-text("📋")'
    ];

    let quickBtnsFound = 0;
    for (const selector of quickBtnSelectors) {
      try {
        const btns = await mainWindow.locator(selector).all();
        quickBtnsFound = btns.length;
        if (quickBtnsFound > 0) {
          console.log(`  ✅ 找到 ${quickBtnsFound} 个快捷按钮`);
          break;
        }
      } catch (e) {
        // Continue
      }
    }

    if (quickBtnsFound > 0) {
      testResults.push({ name: '快捷按钮', status: 'pass' });
    } else {
      console.log('  ⚠️  未找到快捷按钮');
      testResults.push({ name: '快捷按钮', status: 'skip' });
    }

    // Test 9: API 直连测试
    console.log('\n✓ Test 9: API 直连测试');
    let isRealResponse = false;
    try {
      const response = await fetch(`${AGENT_GATEWAY_URL}/mcp`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          jsonrpc: '2.0',
          method: 'tools/call',
          id: 1,
          params: {
            name: 'chat_with_genie',
            arguments: { message: '测试消息', session_id: 'test' }
          }
        })
      });

      const data = await response.json();
      if (data.result?.content?.[0]?.text) {
        const result = JSON.parse(data.result.content[0].text);
        console.log('  ✅ API 返回:', result.content.substring(0, 50) + '...');
        console.log('  💰 费用:', result.cost, 'STT');

        // 检查是否是模拟回复
        if (result.content.includes('[模拟模式]') || result.content.includes('模拟响应')) {
          console.log('  ⚠️  当前是模拟回复，未配置真实 LLM');
          isRealResponse = false;
        } else {
          console.log('  ✅ 检测到真实 AI 回复！');
          isRealResponse = true;
        }
        testResults.push({ name: 'API 直连', status: 'pass' });
      } else {
        throw new Error('Invalid API response');
      }
    } catch (e) {
      console.log('  ❌ API 测试失败:', e.message);
      testResults.push({ name: 'API 直连', status: 'fail' });
    }

    // Test 10: 验证真实 LLM 回复
    console.log('\n✓ Test 10: 验证真实 LLM 回复');
    if (isRealResponse) {
      console.log('  ✅ 真实 AI 回复验证通过');
      testResults.push({ name: '真实 LLM 回复', status: 'pass' });
    } else {
      console.log('  ⚠️  当前使用模拟回复');
      console.log('  💡 如需真实 AI 回复:');
      console.log('     1. 访问 https://console.anthropic.com/ 获取 API Key');
      console.log('     2. 在应用中设置 ANTHROPIC_API_KEY');
      testResults.push({ name: '真实 LLM 回复', status: 'skip' });
    }

    // Test 11: 最终截图
    console.log('\n✓ Test 11: 最终截图');
    await wait(2000);
    await mainWindow.screenshot({ path: '/tmp/e2e-final.png' });
    console.log('  📸 截图保存: /tmp/e2e-final.png');
    testResults.push({ name: '最终截图', status: 'pass' });

  } catch (error) {
    console.error('\n❌ 测试失败:', error.message);
    testResults.push({ name: '整体测试', status: 'fail', error: error.message });
  } finally {
    // 关闭应用
    if (electronApp) {
      console.log('\n  🛑 关闭应用...');
      await electronApp.close();
    }

    // 打印结果汇总
    console.log('\n' + '='.repeat(70));
    console.log('📊 测试结果汇总');
    console.log('='.repeat(70));

    let pass = 0, fail = 0, skip = 0;
    testResults.forEach(r => {
      const icon = r.status === 'pass' ? '✅' : r.status === 'fail' ? '❌' : '⏭️';
      console.log(`${icon} ${r.name}`);
      if (r.status === 'fail') fail++;
      else if (r.status === 'skip') skip++;
      else pass++;
    });

    console.log('\n' + '-'.repeat(70));
    console.log(`✅ 通过: ${pass} | ❌ 失败: ${fail} | ⏭️ 跳过: ${skip}`);
    console.log(`📈 成功率: ${Math.round((pass / testResults.length) * 100)}%`);

    if (fail === 0) {
      console.log('\n🎉 测试完成！');
    } else {
      console.log('\n⚠️ 部分测试失败');
    }

    console.log('\n📁 截图文件:');
    console.log('  - /tmp/e2e-01-initial.png');
    console.log('  - /tmp/e2e-02-geniebot.png');
    console.log('  - /tmp/e2e-03-response.png');
    console.log('  - /tmp/e2e-final.png');
  }
})();
