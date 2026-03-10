const { test, expect } = require('@playwright/test');

test('capture homepage screenshot', async ({ page }) => {
  await page.goto('file:///tmp/mac-test/ShareToken.app/Contents/Resources/frontend/index.html');
  await page.waitForTimeout(3000);
  await page.screenshot({ path: '/tmp/screenshots/homepage.png', fullPage: true });
  console.log('Homepage screenshot saved');
});

test('capture wallet page', async ({ page }) => {
  await page.goto('file:///tmp/mac-test/ShareToken.app/Contents/Resources/frontend/index.html#/wallet');
  await page.waitForTimeout(3000);
  await page.screenshot({ path: '/tmp/screenshots/wallet.png', fullPage: true });
  console.log('Wallet page screenshot saved');
});

test('capture market page', async ({ page }) => {
  await page.goto('file:///tmp/mac-test/ShareToken.app/Contents/Resources/frontend/index.html#/market');
  await page.waitForTimeout(3000);
  await page.screenshot({ path: '/tmp/screenshots/market.png', fullPage: true });
  console.log('Market page screenshot saved');
});
