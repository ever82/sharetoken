const { contextBridge, ipcRenderer } = require('electron');

// 暴露安全的 API 给前端
contextBridge.exposeInMainWorld('electronAPI', {
  // 应用信息
  getAppVersion: () => ipcRenderer.invoke('get-app-version'),
  getPlatform: () => ipcRenderer.invoke('get-platform'),
  getUserDataPath: () => ipcRenderer.invoke('get-user-data-path'),

  // 外部链接
  openExternal: (url) => ipcRenderer.invoke('open-external', url),

  // 文件操作
  showSaveDialog: (options) => ipcRenderer.invoke('show-save-dialog', options),
  showOpenDialog: (options) => ipcRenderer.invoke('show-open-dialog', options),
  readFile: (filePath) => ipcRenderer.invoke('read-file', filePath),
  writeFile: (filePath, data) => ipcRenderer.invoke('write-file', filePath, data),

  // 节点状态
  checkNodeStatus: () => ipcRenderer.invoke('check-node-status'),

  // 是否是 Electron 环境
  isElectron: true
});

console.log('Electron preload script loaded');
