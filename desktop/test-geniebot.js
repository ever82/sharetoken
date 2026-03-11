const { webkit } = require('playwright');
const path = require('path');

(async () => {
  console.log('🧪 Starting GenieBot E2E Test...\n');

  // Launch browser
  const browser = await webkit.launch({
    headless: false,
    slowMo: 500
  });

  const context = await browser.newContext({
    viewport: { width: 1400, height: 900 }
  });

  const page = await context.newPage();

  try {
    // Test 1: Load the app
    console.log('Test 1: Loading ShareToken Desktop...');
    await page.goto('file:///Applications/ShareToken.app/Contents/Resources/frontend/index.html');
    await page.waitForLoadState('networkidle');
    console.log('✅ App loaded\n');

    // Wait for page to be ready
    await page.waitForTimeout(3000);

    // Take screenshot of initial state
    await page.screenshot({ path: '/tmp/test-01-initial.png' });
    console.log('📸 Screenshot saved: /tmp/test-01-initial.png\n');

    // Test 2: Check GenieBot navigation
    console.log('Test 2: Navigating to GenieBot...');

    // Look for GenieBot link/button
    const genieBotLink = await page.locator('text=GenieBot, button:has-text("GenieBot"), a:has-text("GenieBot")').first();
    if (await genieBotLink.isVisible().catch(() => false)) {
      await genieBotLink.click();
      console.log('✅ Clicked GenieBot link\n');
    } else {
      // Try to find by href
      const links = await page.locator('a, button').all();
      for (const link of links) {
        const text = await link.textContent().catch(() => '');
        if (text.toLowerCase().includes('genie') || text.includes('🧞')) {
          await link.click();
          console.log('✅ Found and clicked GenieBot element\n');
          break;
        }
      }
    }

    await page.waitForTimeout(2000);
    await page.screenshot({ path: '/tmp/test-02-geniebot.png' });
    console.log('📸 Screenshot saved: /tmp/test-02-geniebot.png\n');

    // Test 3: Check Agent Gateway connection
    console.log('Test 3: Checking Agent Gateway connection...');
    const connectionStatus = await page.locator('.connection-status, .status-text').textContent().catch(() => 'Not found');
    console.log('Connection status:', connectionStatus);

    const gatewayUrl = await page.locator('.gateway-url').textContent().catch(() => 'Not found');
    console.log('Gateway URL:', gatewayUrl);

    // Test 4: Check balance display
    console.log('\nTest 4: Checking balance...');
    const balance = await page.locator('.balance').textContent().catch(() => 'Not found');
    console.log('Balance:', balance);

    // Test 5: Send a message
    console.log('\nTest 5: Sending test message...');
    const input = await page.locator('textarea, .chat-input, [placeholder*="message"]').first();
    if (await input.isVisible().catch(() => false)) {
      await input.fill('你好，请查询我的余额');
      console.log('✅ Typed message\n');

      // Click send button
      const sendBtn = await page.locator('button:has-text("➤"), .send-btn, button[type="submit"]').first();
      if (await sendBtn.isVisible().catch(() => false)) {
        await sendBtn.click();
        console.log('✅ Clicked send button\n');

        // Wait for response
        await page.waitForTimeout(3000);
        await page.screenshot({ path: '/tmp/test-03-response.png' });
        console.log('📸 Screenshot saved: /tmp/test-03-response.png\n');

        // Check for response
        const messages = await page.locator('.message').all();
        console.log('Number of messages:', messages.length);

        if (messages.length > 1) {
          const lastMessage = await messages[messages.length - 1].textContent();
          console.log('Last message:', lastMessage.substring(0, 100));
        }
      } else {
        console.log('⚠️ Send button not found\n');
      }
    } else {
      console.log('⚠️ Input field not found\n');
    }

    // Test 6: Check Agent selector
    console.log('Test 6: Checking Agent selector...');
    const agentSelect = await page.locator('select.agent-select, .agent-selector select').first();
    if (await agentSelect.isVisible().catch(() => false)) {
      const options = await agentSelect.locator('option').allTextContents();
      console.log('Available agents:', options);
      console.log('✅ Agent selector found\n');
    } else {
      console.log('⚠️ Agent selector not found\n');
    }

    // Test 7: Check quick actions
    console.log('Test 7: Checking quick actions...');
    const quickBtns = await page.locator('.quick-btn, .quick-actions button').all();
    console.log('Quick action buttons:', quickBtns.length);
    for (const btn of quickBtns) {
      const text = await btn.textContent().catch(() => '');
      console.log(' -', text);
    }

    console.log('\n✅ All tests completed!');

    // Final screenshot
    await page.screenshot({ path: '/tmp/test-final.png' });
    console.log('📸 Final screenshot: /tmp/test-final.png');

  } catch (error) {
    console.error('❌ Test failed:', error.message);
    await page.screenshot({ path: '/tmp/test-error.png' });
    console.log('📸 Error screenshot: /tmp/test-error.png');
  } finally {
    await browser.close();
    console.log('\n🧪 Browser closed');
  }
})();
