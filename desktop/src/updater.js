const { autoUpdater } = require('electron-updater');
const { dialog, BrowserWindow } = require('electron');
const log = require('electron-log');

// 配置日志
autoUpdater.logger = log;
autoUpdater.logger.transports.file.level = 'info';

// 更新状态
let updateAvailable = false;
let updateDownloaded = false;

// 初始化自动更新
function initAutoUpdater(mainWindow) {
  // 检查更新时显示进度
  autoUpdater.on('checking-for-update', () => {
    log.info('正在检查更新...');
    sendStatusToWindow(mainWindow, 'checking');
  });

  // 有可用更新
  autoUpdater.on('update-available', (info) => {
    log.info('发现可用更新:', info);
    updateAvailable = true;
    sendStatusToWindow(mainWindow, 'available', info);

    // 显示更新提示对话框
    dialog.showMessageBox(mainWindow, {
      type: 'info',
      title: '发现新版本',
      message: `ShareToken ${info.version} 已发布`,
      detail: `当前版本: ${autoUpdater.currentVersion}\n新版本: ${info.version}\n\n更新内容:\n${info.releaseNotes || '性能优化和问题修复'}`,
      buttons: ['立即下载', '稍后提醒'],
      defaultId: 0
    }).then(({ response }) => {
      if (response === 0) {
        // 用户选择下载
        autoUpdater.downloadUpdate();
      }
    });
  });

  // 没有可用更新
  autoUpdater.on('update-not-available', (info) => {
    log.info('当前已是最新版本:', info);
    sendStatusToWindow(mainWindow, 'not-available', info);
  });

  // 下载进度
  autoUpdater.on('download-progress', (progressObj) => {
    log.info('下载进度:', progressObj);
    sendStatusToWindow(mainWindow, 'progress', progressObj);

    // 更新窗口标题显示进度
    const percent = Math.round(progressObj.percent);
    mainWindow.setTitle(`ShareToken - 下载更新 ${percent}%`);
  });

  // 更新下载完成
  autoUpdater.on('update-downloaded', (info) => {
    log.info('更新下载完成:', info);
    updateDownloaded = true;
    sendStatusToWindow(mainWindow, 'downloaded', info);
    mainWindow.setTitle('ShareToken');

    // 询问用户是否立即安装
    dialog.showMessageBox(mainWindow, {
      type: 'info',
      title: '更新已就绪',
      message: `ShareToken ${info.version} 已下载完成`,
      detail: '是否立即安装更新？应用将自动重启。',
      buttons: ['立即安装', '稍后安装'],
      defaultId: 0
    }).then(({ response }) => {
      if (response === 0) {
        // 立即安装并重启
        autoUpdater.quitAndInstall(true, true);
      }
    });
  });

  // 更新错误
  autoUpdater.on('error', (err) => {
    log.error('更新错误:', err);
    sendStatusToWindow(mainWindow, 'error', { message: err.message });

    // 只在用户主动检查更新时显示错误对话框
    if (isManualCheck) {
      dialog.showErrorBox('更新错误', `检查更新时发生错误:\n${err.message}`);
    }
  });
}

// 发送状态到渲染进程
function sendStatusToWindow(mainWindow, status, data = {}) {
  if (mainWindow && mainWindow.webContents) {
    mainWindow.webContents.send('update-status', { status, data });
  }
}

// 是否手动检查标志
let isManualCheck = false;

// 检查更新
function checkForUpdates(mainWindow, manual = false) {
  isManualCheck = manual;

  // 只在生产环境检查更新
  if (process.env.NODE_ENV === 'development') {
    log.info('开发模式，跳过更新检查');
    if (manual) {
      dialog.showMessageBox(mainWindow, {
        type: 'info',
        title: '开发模式',
        message: '开发模式下不检查更新'
      });
    }
    return;
  }

  autoUpdater.checkForUpdates();
}

// 立即安装更新
function installUpdate() {
  if (updateDownloaded) {
    autoUpdater.quitAndInstall(true, true);
  }
}

// 获取更新状态
function getUpdateStatus() {
  return {
    available: updateAvailable,
    downloaded: updateDownloaded
  };
}

// 设置自动检查（每小时检查一次）
function startAutoCheck(mainWindow) {
  // 首次启动延迟检查（避免启动时网络拥堵）
  setTimeout(() => {
    checkForUpdates(mainWindow);
  }, 30000); // 30秒后

  // 每小时检查一次
  setInterval(() => {
    checkForUpdates(mainWindow);
  }, 3600000); // 1小时
}

module.exports = {
  initAutoUpdater,
  checkForUpdates,
  installUpdate,
  getUpdateStatus,
  startAutoCheck
};
