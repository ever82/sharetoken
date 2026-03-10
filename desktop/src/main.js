const { app, BrowserWindow, ipcMain, dialog, shell } = require('electron');
const path = require('path');
const { spawn } = require('child_process');
const fs = require('fs');

// Wallet manager
const walletManager = require('./wallet');

// 保持全局窗口对象，防止被垃圾回收
let mainWindow;
let sharetokenProcess = null;
let walletPaths = null;

// 确定平台
const platform = process.platform;
const isDev = process.env.NODE_ENV === 'development';

// 获取 sharetokend 二进制路径
function getSharetokenPath() {
  if (isDev) {
    const devPath = path.join(__dirname, '../../bin/sharetokend');
    const devPathWin = devPath + '.exe';
    return fs.existsSync(devPathWin) ? devPathWin : devPath;
  }

  // 打包后的路径 - 尝试多个可能的路径
  const binDir = path.join(process.resourcesPath, 'bin');
  const candidates = [
    path.join(binDir, 'sharetokend.exe'),  // Windows
    path.join(binDir, 'sharetokend'),      // Linux/macOS
  ];

  for (const candidate of candidates) {
    if (fs.existsSync(candidate)) {
      return candidate;
    }
  }

  // 默认返回平台对应的路径
  const exeName = platform === 'win32' ? 'sharetokend.exe' : 'sharetokend';
  return path.join(binDir, exeName);
}

// 获取前端资源路径
function getFrontendPath() {
  if (isDev) {
    // 开发模式：使用本地前端服务器
    return 'http://localhost:8080';
  }

  // 打包后的路径
  return path.join(process.resourcesPath, 'frontend', 'index.html');
}

// 启动 sharetoken 节点
function startSharetokenNode() {
  const sharetokenPath = getSharetokenPath();

  if (!fs.existsSync(sharetokenPath)) {
    console.error('ShareToken binary not found:', sharetokenPath);
    return Promise.reject(new Error('找不到 ShareToken 节点程序: ' + sharetokenPath));
  }

  const dataDir = path.join(app.getPath('userData'), 'sharetoken-data');
  const nodeDataTemplate = path.join(process.resourcesPath, 'node-data');

  // 如果数据目录不存在且模板存在，复制模板
  if (!fs.existsSync(dataDir) && fs.existsSync(nodeDataTemplate)) {
    console.log('Copying node data template to:', dataDir);
    copyDir(nodeDataTemplate, dataDir);
  }

  // 确保数据目录存在
  if (!fs.existsSync(dataDir)) {
    fs.mkdirSync(dataDir, { recursive: true });
  }

  const args = [
    'start',
    '--home', dataDir,
    '--rpc.laddr', 'tcp://127.0.0.1:26657',
    '--api.address', 'tcp://127.0.0.1:1317'
  ];

  console.log('Starting ShareToken node:', sharetokenPath, args);

  sharetokenProcess = spawn(sharetokenPath, args, {
    detached: false,
    stdio: ['ignore', 'pipe', 'pipe']
  });

  sharetokenProcess.stdout.on('data', (data) => {
    console.log(`[ShareToken] ${data}`);
  });

  sharetokenProcess.stderr.on('data', (data) => {
    console.error(`[ShareToken Error] ${data}`);
  });

  sharetokenProcess.on('close', (code) => {
    console.log(`ShareToken process exited with code ${code}`);
    sharetokenProcess = null;
  });

  // 等待节点启动
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve();
    }, 3000);
  });
}

// 复制目录
function copyDir(src, dest) {
  if (!fs.existsSync(dest)) {
    fs.mkdirSync(dest, { recursive: true });
  }
  const entries = fs.readdirSync(src, { withFileTypes: true });
  for (const entry of entries) {
    const srcPath = path.join(src, entry.name);
    const destPath = path.join(dest, entry.name);
    if (entry.isDirectory()) {
      copyDir(srcPath, destPath);
    } else {
      fs.copyFileSync(srcPath, destPath);
    }
  }
}

// 停止 sharetoken 节点
function stopSharetokenNode() {
  if (sharetokenProcess) {
    console.log('Stopping ShareToken node...');
    sharetokenProcess.kill();
    sharetokenProcess = null;
  }
}

// 创建主窗口
function createWindow() {
  mainWindow = new BrowserWindow({
    width: 1400,
    height: 900,
    minWidth: 1200,
    minHeight: 700,
    webPreferences: {
      nodeIntegration: false,
      contextIsolation: true,
      preload: path.join(__dirname, 'preload.js'),
      webSecurity: false // 允许加载本地文件
    },
    titleBarStyle: 'default',
    show: false,
    icon: path.join(__dirname, '../build/icon.png')
  });

  // 加载前端
  const frontendPath = getFrontendPath();

  if (frontendPath.startsWith('http')) {
    mainWindow.loadURL(frontendPath);
  } else {
    mainWindow.loadFile(frontendPath);
  }

  // 开发模式下打开开发者工具
  if (isDev) {
    mainWindow.webContents.openDevTools();
  }

  // 窗口准备好后显示
  mainWindow.once('ready-to-show', () => {
    mainWindow.show();

    // 初始化钱包
    initializeWallet().then((result) => {
      console.log('Wallet initialization result:', result);
      // 通知前端钱包状态
      if (mainWindow && !mainWindow.isDestroyed()) {
        mainWindow.webContents.send('wallet-initialized', result);
      }
    }).catch((err) => {
      console.error('Failed to initialize wallet:', err);
    });

    // 启动本地节点
    startSharetokenNode().then(() => {
      console.log('ShareToken node started');
    }).catch((err) => {
      console.error('Failed to start ShareToken node:', err);
      dialog.showErrorBox('启动错误', err.message);
    });
  });

  // 窗口关闭事件
  mainWindow.on('closed', () => {
    mainWindow = null;
    stopSharetokenNode();
  });

  // 处理新窗口请求（外部链接用系统浏览器打开）
  mainWindow.webContents.setWindowOpenHandler(({ url }) => {
    shell.openExternal(url);
    return { action: 'deny' };
  });
}

// 应用就绪
app.whenReady().then(() => {
  createWindow();

  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow();
    }
  });
});

// 所有窗口关闭时退出
app.on('window-all-closed', () => {
  stopSharetokenNode();

  if (platform !== 'darwin') {
    app.quit();
  }
});

// 应用退出前
app.on('before-quit', () => {
  stopSharetokenNode();
});

// 初始化钱包（首次启动自动创建）
async function initializeWallet() {
  walletPaths = walletManager.getWalletPaths(app.getPath('userData'));

  if (!walletManager.walletExists(walletPaths)) {
    console.log('First launch - creating new wallet...');
    try {
      const wallet = await walletManager.createWallet();
      walletManager.saveEncryptedMnemonic(wallet.mnemonic, walletPaths);
      console.log('Wallet created:', wallet.address);
      return { created: true, address: wallet.address, needsBackup: true };
    } catch (error) {
      console.error('Failed to create wallet:', error);
      return { created: false, error: error.message };
    }
  } else {
    console.log('Wallet already exists');
    const address = await walletManager.getWalletAddress(walletPaths);
    const needsBackup = walletManager.needsBackup(walletPaths);
    return { created: false, address, needsBackup };
  }
}

// IPC 通信处理

// 初始化钱包
ipcMain.handle('wallet-init', async () => {
  return await initializeWallet();
});

// 获取钱包地址
ipcMain.handle('wallet-get-address', async () => {
  if (!walletPaths) {
    walletPaths = walletManager.getWalletPaths(app.getPath('userData'));
  }
  return await walletManager.getWalletAddress(walletPaths);
});

// 获取钱包状态
ipcMain.handle('wallet-get-status', async () => {
  if (!walletPaths) {
    walletPaths = walletManager.getWalletPaths(app.getPath('userData'));
  }
  const exists = walletManager.walletExists(walletPaths);
  const needsBackup = walletManager.needsBackup(walletPaths);
  const address = exists ? await walletManager.getWalletAddress(walletPaths) : null;

  return {
    exists,
    address,
    needsBackup
  };
});

// 获取余额
ipcMain.handle('wallet-get-balance', async (event, address) => {
  if (!walletPaths) {
    walletPaths = walletManager.getWalletPaths(app.getPath('userData'));
  }

  try {
    const addr = address || await walletManager.getWalletAddress(walletPaths);
    if (!addr) {
      return { success: false, error: 'No wallet found' };
    }

    const balances = await walletManager.getBalance(addr);
    return { success: true, balances };
  } catch (error) {
    console.error('Failed to get balance:', error);
    return { success: false, error: error.message };
  }
});

// 发送交易
ipcMain.handle('wallet-send', async (event, { recipient, amount, denom, memo }) => {
  if (!walletPaths) {
    walletPaths = walletManager.getWalletPaths(app.getPath('userData'));
  }

  try {
    const result = await walletManager.sendTokens(
      recipient,
      amount,
      denom || 'stake',
      memo || '',
      walletPaths
    );
    return { success: true, result };
  } catch (error) {
    console.error('Failed to send tokens:', error);
    return { success: false, error: error.message };
  }
});

// 导出助记词
ipcMain.handle('wallet-export', async (event, { password }) => {
  if (!walletPaths) {
    walletPaths = walletManager.getWalletPaths(app.getPath('userData'));
  }

  try {
    const data = await walletManager.exportWalletData(walletPaths, password);
    walletManager.markBackupComplete(walletPaths);
    return { success: true, mnemonic: data.mnemonic };
  } catch (error) {
    console.error('Failed to export wallet:', error);
    return { success: false, error: error.message };
  }
});

// 标记备份完成
ipcMain.handle('wallet-mark-backup', async () => {
  if (!walletPaths) {
    walletPaths = walletManager.getWalletPaths(app.getPath('userData'));
  }
  walletManager.markBackupComplete(walletPaths);
  return { success: true };
});

// 从助记词恢复钱包
ipcMain.handle('wallet-restore', async (event, { mnemonic }) => {
  if (!walletPaths) {
    walletPaths = walletManager.getWalletPaths(app.getPath('userData'));
  }

  try {
    if (!walletManager.isValidMnemonic(mnemonic)) {
      return { success: false, error: 'Invalid mnemonic format' };
    }

    const wallet = await walletManager.restoreWallet(mnemonic);
    walletManager.saveEncryptedMnemonic(mnemonic, walletPaths);
    walletManager.markBackupComplete(walletPaths);

    return { success: true, address: wallet.address };
  } catch (error) {
    console.error('Failed to restore wallet:', error);
    return { success: false, error: error.message };
  }
});

// 删除钱包（用于切换）
ipcMain.handle('wallet-delete', async () => {
  if (!walletPaths) {
    walletPaths = walletManager.getWalletPaths(app.getPath('userData'));
  }

  try {
    walletManager.deleteWallet(walletPaths);
    return { success: true };
  } catch (error) {
    console.error('Failed to delete wallet:', error);
    return { success: false, error: error.message };
  }
});

// IPC 通信处理

// 获取应用版本
ipcMain.handle('get-app-version', () => {
  return app.getVersion();
});

// 获取平台信息
ipcMain.handle('get-platform', () => {
  return platform;
});

// 打开外部链接
ipcMain.handle('open-external', async (event, url) => {
  await shell.openExternal(url);
});

// 显示保存对话框
ipcMain.handle('show-save-dialog', async (event, options) => {
  const result = await dialog.showSaveDialog(mainWindow, options);
  return result;
});

// 显示打开对话框
ipcMain.handle('show-open-dialog', async (event, options) => {
  const result = await dialog.showOpenDialog(mainWindow, options);
  return result;
});

// 读取文件
ipcMain.handle('read-file', async (event, filePath) => {
  try {
    const data = fs.readFileSync(filePath, 'utf8');
    return { success: true, data };
  } catch (error) {
    return { success: false, error: error.message };
  }
});

// 写入文件
ipcMain.handle('write-file', async (event, filePath, data) => {
  try {
    fs.writeFileSync(filePath, data, 'utf8');
    return { success: true };
  } catch (error) {
    return { success: false, error: error.message };
  }
});

// 获取用户数据目录
ipcMain.handle('get-user-data-path', () => {
  return app.getPath('userData');
});

// 检查节点状态
ipcMain.handle('check-node-status', async () => {
  return {
    running: sharetokenProcess !== null,
    pid: sharetokenProcess ? sharetokenProcess.pid : null
  };
});

console.log('ShareToken Desktop starting...');
console.log('Platform:', platform);
console.log('UserData:', app.getPath('userData'));
