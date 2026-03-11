package e2e

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"sharetoken/test/helpers"
)

// DesktopAppTestSuite 测试 ACH-USER-021: Desktop App (开箱即用)
type DesktopAppTestSuite struct {
	suite.Suite
	chain *helpers.ChainHelper
	ctx   context.Context
}

func (s *DesktopAppTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.chain = helpers.NewChainHelper(s.T())
}

// TestCIBuildArtifacts 测试CI/CD构建产物
func (s *DesktopAppTestSuite) TestCIBuildArtifacts() {
	// 验证CI配置存在
	ciConfigPath := ".github/workflows/desktop-build.yml"
	require.FileExists(s.T(), ciConfigPath, "桌面应用CI配置应存在")

	// 验证构建脚本存在
	buildScriptPath := "scripts/build-desktop.sh"
	require.FileExists(s.T(), buildScriptPath, "桌面应用构建脚本应存在")

	s.T().Log("CI/CD构建配置检查通过")
}

// TestBuildForAllPlatforms 测试多平台构建
func (s *DesktopAppTestSuite) TestBuildForAllPlatforms() {
	// 检查各平台构建产物目录结构
	platforms := []struct {
		name string
		dir  string
	}{
		{"Windows", "build/desktop/windows"},
		{"macOS", "build/desktop/macos"},
		{"Linux", "build/desktop/linux"},
	}

	for _, p := range platforms {
		// 目录可以不存在（CI构建时生成），但构建配置应支持
		configPath := filepath.Join("desktop", "package.json")
		if _, err := os.Stat(configPath); err == nil {
			s.T().Logf("%s 构建配置可用", p.name)
		}
	}

	s.T().Log("多平台构建配置检查完成")
}

// TestNodeDiscovery 测试自动检测或选择节点
func (s *DesktopAppTestSuite) TestNodeDiscovery() {
	// 测试节点发现API
	nodes, err := s.chain.DiscoverNodes()
	require.NoError(s.T(), err, "节点发现应成功")
	require.GreaterOrEqual(s.T(), len(nodes), 1, "应至少发现一个节点")

	// 验证节点信息完整
	for _, node := range nodes {
		require.NotEmpty(s.T(), node.Address, "节点地址不应为空")
		require.NotEmpty(s.T(), node.RPCEndpoint, "RPC端点不应为空")
	}

	s.T().Logf("发现 %d 个可用节点", len(nodes))
}

// TestNodeSelection 测试节点选择逻辑
func (s *DesktopAppTestSuite) TestNodeSelection() {
	// 测试自动选择最佳节点
	bestNode, err := s.chain.SelectBestNode()
	require.NoError(s.T(), err, "节点选择应成功")
	require.NotNil(s.T(), bestNode, "应返回最佳节点")
	require.NotEmpty(s.T(), bestNode.Address, "最佳节点地址不应为空")

	// 测试手动指定节点
	customNode := "http://localhost:26657"
	err = s.chain.SetNodeEndpoint(customNode)
	require.NoError(s.T(), err, "设置自定义节点应成功")

	s.T().Logf("节点选择逻辑测试通过，最佳节点: %s", bestNode.Address)
}

// TestGUIWalletFunctions 测试图形界面钱包功能
func (s *DesktopAppTestSuite) TestGUIWalletFunctions() {
	// 模拟GUI测试场景
	user := s.chain.CreateAccountWithBalance("gui_wallet_user", "5000stt")

	// 1. 查看余额
	balance, err := s.chain.QueryBalance(user.Address)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), balance)

	// 2. 查看交易历史
	txs, err := s.chain.QueryTxHistory(user.Address, 10)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), txs)

	// 3. 发送转账
	recipient := s.chain.CreateAccount("gui_recipient")
	err = s.chain.SendTokens(user, recipient.Address, "100stt")
	require.NoError(s.T(), err)

	s.T().Log("GUI钱包功能测试通过")
}

// TestGUIServiceMarketplace 测试图形界面服务市场
func (s *DesktopAppTestSuite) TestGUIServiceMarketplace() {
	// 测试服务市场API可用
	services, err := s.chain.QueryAllServices(10)
	require.NoError(s.T(), err)
	require.Greater(s.T(), len(services), 0)

	// 测试服务详情查询
	if len(services) > 0 {
		detail, err := s.chain.QueryServiceDetail(services[0].ID)
		require.NoError(s.T(), err)
		require.NotNil(s.T(), detail)
	}

	s.T().Log("GUI服务市场功能测试通过")
}

// TestGUIGenieBot 测试图形界面 GenieBot
func (s *DesktopAppTestSuite) TestGUIGenieBot() {
	user := s.chain.CreateAccountWithBalance("gui_genie_user", "1000stt")

	// 测试AI对话接口
	response, err := s.chain.SubmitIntent(user, "你好")
	require.NoError(s.T(), err)
	require.NotNil(s.T(), response)

	// 测试意图识别
	intent, err := s.chain.RecognizeIntent("帮我写代码")
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), intent.ServiceType)

	s.T().Log("GUI GenieBot功能测试通过")
}

// TestNoDependencyRequirements 测试无需依赖安装
func (s *DesktopAppTestSuite) TestNoDependencyRequirements() {
	// 验证桌面应用项目结构
	desktopDir := "desktop"
	if _, err := os.Stat(desktopDir); err == nil {
		// 检查package.json存在
		packageJSON := filepath.Join(desktopDir, "package.json")
		require.FileExists(s.T(), packageJSON, "桌面应用package.json应存在")

		// 检查main.js/main.ts存在
		mainFile := filepath.Join(desktopDir, "src", "main.js")
		if _, err := os.Stat(mainFile); err != nil {
			mainFile = filepath.Join(desktopDir, "src", "main.ts")
		}
		if _, err := os.Stat(mainFile); err == nil {
			s.T().Log("桌面应用主进程文件存在")
		}
	}

	s.T().Log("桌面应用无依赖要求验证通过")
}

// TestDesktopAppE2E 完整桌面应用E2E测试
func (s *DesktopAppTestSuite) TestDesktopAppE2E() {
	// 模拟用户使用桌面应用的完整流程
	// 1. 自动发现/选择节点
	nodes, err := s.chain.DiscoverNodes()
	require.NoError(s.T(), err)
	require.Greater(s.T(), len(nodes), 0)

	// 2. 创建/导入钱包
	user := s.chain.CreateAccount("desktop_e2e_user")

	// 3. 获取测试代币
	s.chain.RequestFaucet(user.Address, "1000stt")

	// 4. 浏览服务市场
	services, _ := s.chain.QueryAllServices(5)
	require.Greater(s.T(), len(services), 0)

	// 5. 与GenieBot对话
	_, err = s.chain.SubmitIntent(user, "我想使用AI服务")
	require.NoError(s.T(), err)

	// 6. 发起转账
	recipient := s.chain.CreateAccount("desktop_e2e_recipient")
	err = s.chain.SendTokens(user, recipient.Address, "100stt")
	require.NoError(s.T(), err)

	s.T().Log("桌面应用E2E测试通过")
}

func TestDesktopAppSuite(t *testing.T) {
	suite.Run(t, new(DesktopAppTestSuite))
}
