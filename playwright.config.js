// @ts-check
const { defineConfig, devices } = require('@playwright/test');

/**
 * Playwright configuration for ShareToken Desktop E2E testing
 *
 * @see https://playwright.dev/docs/test-configuration
 */
module.exports = defineConfig({
  testDir: './e2e',

  /* Run tests in files in parallel */
  fullyParallel: false, // Electron tests should run sequentially

  /* Fail the build on CI if you accidentally left test.only in the source code */
  forbidOnly: !!process.env.CI,

  /* Retry on CI only */
  retries: process.env.CI ? 2 : 0,

  /* Opt out of parallel tests on CI */
  workers: process.env.CI ? 1 : 1,

  /* Reporter to use */
  reporter: [['html', { open: 'never' }], ['list']],

  /* Shared settings for all the projects below */
  use: {
    headless: true,
    viewport: { width: 1400, height: 900 },
    screenshot: 'only-on-failure',
    trace: 'on-first-retry',
    video: 'retain-on-failure',
  },

  /* Configure projects for different test types */
  projects: [
    {
      name: 'chromium',
      testMatch: /desktop\.spec\.js$/, // Static HTML file tests
      use: { browserName: 'chromium' },
    },
    {
      name: 'electron',
      testMatch: /desktop-wallet\.spec\.js$/, // Electron app tests
      use: {},
    },
  ],

  /* Output directory for test results */
  outputDir: './e2e/results/',
});
