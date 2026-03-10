# ShareToken Desktop 桌面应用

> 开箱即用的 ShareToken 图形界面应用

## 功能特性

- 🚀 **开箱即用** - 下载解压，双击打开，无需任何配置
- 🔐 **内置钱包** - 创建/导入钱包，查看余额，发送转账
- 🤖 **AI 对话** - 内置 GenieBot，一键调用 AI 服务
- 🛒 **服务市场** - 浏览和调用平台上的服务
- 🖥️ **本地节点** - 内置轻量级区块链节点

## 安装使用

### Windows

```powershell
# 下载并解压
Invoke-WebRequest -Uri https://github.com/ever82/sharetoken/releases/latest/download/ShareToken-0.1.0-win-x64.zip -OutFile sharetoken.zip
Expand-Archive -Path sharetoken.zip -DestinationPath .\ShareToken

# 运行
.\ShareToken\ShareToken.exe
```

### macOS

```bash
# 下载
curl -LO https://github.com/ever82/sharetoken/releases/latest/download/ShareToken-0.1.0-mac-x64.dmg

# 打开 DMG 并拖动到 Applications
open ShareToken-0.1.0-mac-x64.dmg
```

### Linux

```bash
# 下载 AppImage
curl -LO https://github.com/ever82/sharetoken/releases/latest/download/ShareToken-0.1.0-linux-x86_64.AppImage
chmod +x ShareToken-0.1.0-linux-x86_64.AppImage
./ShareToken-0.1.0-linux-x86_64.AppImage
```

## 开发

```bash
# 进入 desktop 目录
cd desktop

# 安装依赖
npm install

# 开发模式运行
npm start

# 构建
npm run build

# 构建特定平台
npm run build:win
npm run build:mac
npm run build:linux
```

## 技术栈

- **Electron** - 桌面应用框架
- **Vue.js** - 前端界面（复用 frontend 代码）
- **Go** - 嵌入式区块链节点

## 目录结构

```
desktop/
├── src/
│   ├── main.js      # Electron 主进程
│   └── preload.js   # 预加载脚本
├── build/           # 构建资源（图标等）
├── dist/            # 构建输出
├── package.json     # 项目配置
└── README.md        # 本文件
```

## 注意事项

1. **首次启动**：应用会自动初始化并启动本地节点，可能需要 10-30 秒
2. **数据目录**：
   - Windows: `%APPDATA%/ShareToken/sharetoken-data`
   - macOS: `~/Library/Application Support/ShareToken/sharetoken-data`
   - Linux: `~/.config/ShareToken/sharetoken-data`
3. **端口占用**：默认使用 26657 (RPC) 和 1317 (API) 端口

## 许可证

MIT
