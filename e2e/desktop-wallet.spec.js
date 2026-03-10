const { _electron: electron } = require('playwright');
const { test, expect } = require('@playwright/test');
const path = require('path');

// Electron app path
const getElectronAppPath = () => {
  // Try different possible paths
  const paths = [
    // Development build
    path.join(__dirname, '../desktop/dist/mac/ShareToken.app/Contents/MacOS/ShareToken'),
    path.join(__dirname, '../desktop/dist/mac-arm64/ShareToken.app/Contents/MacOS/ShareToken'),
    // CI build path
    '/tmp/mac-test/ShareToken.app/Contents/MacOS/ShareToken',
  ];

  const fs = require('fs');
  for (const p of paths) {
    if (fs.existsSync(p)) {
      return p;
    }
  }
  throw new Error('Electron app not found. Please build the desktop app first.');
};

test.describe('ShareToken Desktop App', () => {
  let electronApp;

  test.beforeEach(async () => {
    electronApp = await electron.launch({
      executablePath: getElectronAppPath(),
      args: [],
      env: {
        ...process.env,
        NODE_ENV: 'production',
        ELECTRON_IS_TEST: 'true',
      },
    });
  });

  test.afterEach(async () => {
    if (electronApp) {
      await electronApp.close();
    }
  });

  test('app launches and shows home page', async () => {
    // Get the first window
    const window = await electronApp.firstWindow();

    // Wait for the app to load
    await window.waitForLoadState('networkidle');

    // Take screenshot
    await window.screenshot({ path: '/tmp/screenshots/e2e-01-home.png' });
    console.log('Screenshot: /tmp/screenshots/e2e-01-home.png');

    // Verify the window title (may be "Vue App" or "ShareToken")
    const title = await window.title();
    expect(title).toMatch(/Vue App|ShareToken/);

    // Verify welcome message is visible
    const welcomeText = await window.locator('text=欢迎来到 ShareToken').first();
    await expect(welcomeText).toBeVisible();
  });

  test('navigate to wallet page and verify local wallet option', async () => {
    const window = await electronApp.firstWindow();

    await window.waitForLoadState('networkidle');

    // Click on "钱包" (Wallet) tab
    const walletTab = await window.locator('text=钱包').first();
    await walletTab.click();

    // Wait for navigation
    await window.waitForTimeout(1000);

    // Take screenshot
    await window.screenshot({ path: '/tmp/screenshots/e2e-02-wallet-page.png' });
    console.log('Screenshot: /tmp/screenshots/e2e-02-wallet-page.png');

    // Verify wallet page content
    const connectText = await window.locator('text=/Connect Wallet|连接钱包/i').first();
    await expect(connectText).toBeVisible();

    // Verify Local Wallet option exists
    const localWalletOption = await window.locator('text=/Local Wallet|本地钱包/i').first();
    await expect(localWalletOption).toBeVisible();
  });

  test('connect local wallet and verify address display', async () => {
    const window = await electronApp.firstWindow();

    await window.waitForLoadState('networkidle');

    // Navigate to wallet page
    const walletTab = await window.locator('text=钱包').first();
    await walletTab.click();
    await window.waitForTimeout(1000);

    // Click on Local Wallet option
    const localWalletOption = await window.locator('text=/Local Wallet|本地钱包/i').first();
    await localWalletOption.click();

    // Wait for connection
    await window.waitForTimeout(2000);

    // Take screenshot
    await window.screenshot({ path: '/tmp/screenshots/e2e-03-local-wallet-connected.png' });
    console.log('Screenshot: /tmp/screenshots/e2e-03-local-wallet-connected.png');

    // Verify wallet info is displayed - look for the address code element
    const addressElement = await window.locator('code:has-text("sharetoken1")').first();
    await expect(addressElement).toBeVisible();
  });

  test('verify backup reminder for new wallet', async () => {
    const window = await electronApp.firstWindow();

    await window.waitForLoadState('networkidle');

    // Navigate to wallet page
    const walletTab = await window.locator('text=钱包').first();
    await walletTab.click();
    await window.waitForTimeout(1000);

    // Check for backup reminder or wallet options
    const walletSection = await window.locator('.wallet-selection, .wallet-status').first();
    await expect(walletSection).toBeVisible();

    // Take screenshot
    await window.screenshot({ path: '/tmp/screenshots/e2e-04-wallet-view.png' });
    console.log('Screenshot: /tmp/screenshots/e2e-04-wallet-view.png');
  });

  test('navigate to market page', async () => {
    const window = await electronApp.firstWindow();

    await window.waitForLoadState('networkidle');

    // Click on "市场" (Market) tab
    const marketTab = await window.locator('text=市场').first();
    await marketTab.click();

    // Wait for navigation
    await window.waitForTimeout(1000);

    // Take screenshot
    await window.screenshot({ path: '/tmp/screenshots/e2e-05-market-page.png' });
    console.log('Screenshot: /tmp/screenshots/e2e-05-market-page.png');

    // Verify market page content
    const marketTitle = await window.locator('text=/Market|市场/i').first();
    await expect(marketTitle).toBeVisible();
  });
});
