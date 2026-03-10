const puppeteer = require('puppeteer');

(async () => {
  const browser = await puppeteer.launch({
    headless: true,
    args: ['--no-sandbox', '--disable-setuid-sandbox']
  });

  console.log('Testing ShareToken Desktop App...');

  // Test 1: Homepage
  console.log('1. Screenshot homepage...');
  const page1 = await browser.newPage();
  await page1.setViewport({ width: 1400, height: 900 });
  await page1.goto('file:///tmp/mac-test/ShareToken.app/Contents/Resources/frontend/index.html', {
    waitUntil: 'networkidle0',
    timeout: 30000
  });
  await page1.waitForTimeout(2000);
  await page1.screenshot({ path: '/tmp/screenshots/homepage.png', fullPage: true });
  console.log('   ✅ Saved: /tmp/screenshots/homepage.png');
  await page1.close();

  // Test 2: Wallet Page
  console.log('2. Screenshot wallet page...');
  const page2 = await browser.newPage();
  await page2.setViewport({ width: 1400, height: 900 });
  await page2.goto('file:///tmp/mac-test/ShareToken.app/Contents/Resources/frontend/index.html#/wallet', {
    waitUntil: 'networkidle0',
    timeout: 30000
  });
  await page2.waitForTimeout(2000);
  await page2.screenshot({ path: '/tmp/screenshots/wallet.png', fullPage: true });
  console.log('   ✅ Saved: /tmp/screenshots/wallet.png');
  await page2.close();

  // Test 3: Market Page
  console.log('3. Screenshot market page...');
  const page3 = await browser.newPage();
  await page3.setViewport({ width: 1400, height: 900 });
  await page3.goto('file:///tmp/mac-test/ShareToken.app/Contents/Resources/frontend/index.html#/market', {
    waitUntil: 'networkidle0',
    timeout: 30000
  });
  await page3.waitForTimeout(2000);
  await page3.screenshot({ path: '/tmp/screenshots/market.png', fullPage: true });
  console.log('   ✅ Saved: /tmp/screenshots/market.png');
  await page3.close();

  await browser.close();
  console.log('\n✅ All screenshots captured!');
})();
