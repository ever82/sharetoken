package e2e

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// DataExportTest 测试数据导出功能
type DataExportTest struct {
	E2ETestSuite
}

func TestDataExport(t *testing.T) {
	suite.Run(t, new(DataExportTest))
}

// TestExportTransactionHistory 测试交易记录导出
func (s *DataExportTest) TestExportTransactionHistory() {
	// 创建用户并执行一些交易
	user := s.CreateVerifiedUser("github")
	s.FundAccount(user.Address, 100000000)

	// 执行转账
	recipient := s.CreateVerifiedUser("github")
	s.Transfer(user.Address, recipient.Address, 10000000)
	s.Transfer(user.Address, recipient.Address, 5000000)

	// 导出交易记录
	export, err := s.ExportUserData(user.Address, "transactions")
	s.Require().NoError(err)
	assert.NotEmpty(s.T(), export)

	// 验证 JSON 格式
	var transactions []map[string]interface{}
	err = json.Unmarshal(export, &transactions)
	s.Require().NoError(err)
	assert.GreaterOrEqual(s.T(), len(transactions), 2)
}

// TestExportTaskHistory 测试任务历史导出
func (s *DataExportTest) TestExportTaskHistory() {
	// 创建用户并创建任务
	user := s.CreateVerifiedUser("github")

	// 创建几个任务
	for i := 0; i < 3; i++ {
		task := WorkflowTask{
			Title:       "Test Task",
			Description: "Test Description",
			Budget:      1000000,
		}
		s.CreateWorkflowTask(user.Address, task)
	}

	// 导出任务历史
	export, err := s.ExportUserData(user.Address, "tasks")
	s.Require().NoError(err)

	// 验证
	var tasks []map[string]interface{}
	err = json.Unmarshal(export, &tasks)
	s.Require().NoError(err)
	assert.GreaterOrEqual(s.T(), len(tasks), 3)
}

// TestExportReviews 测试评价导出
func (s *DataExportTest) TestExportReviews() {
	// 创建用户并提交评价
	user := s.CreateVerifiedUser("github")
	provider := s.CreateVerifiedUser("github")

	// 完成任务并评价
	task := WorkflowTask{Title: "Review Task", Budget: 1000000}
	taskID, _ := s.CreateWorkflowTask(user.Address, task)
	s.FundEscrowForTask(user.Address, taskID, 1000000)
	s.AcceptWorkflowTask(provider.Address, taskID)
	s.CompleteMilestone(provider.Address, taskID, 0)
	s.ConfirmDelivery(user.Address, taskID)

	review := TaskReview{
		TaskID:  taskID,
		Rating:  5,
		Comment: "Great work!",
	}
	s.SubmitTaskReview(user.Address, review)

	// 导出评价
	export, err := s.ExportUserData(user.Address, "reviews")
	s.Require().NoError(err)

	var reviews []map[string]interface{}
	err = json.Unmarshal(export, &reviews)
	s.Require().NoError(err)
	assert.GreaterOrEqual(s.T(), len(reviews), 1)
}

// TestExportSettings 测试设置导出
func (s *DataExportTest) TestExportSettings() {
	// 创建用户
	user := s.CreateVerifiedUser("github")

	// 导出设置
	export, err := s.ExportUserData(user.Address, "settings")
	s.Require().NoError(err)

	// 验证 JSON 格式
	var settings map[string]interface{}
	err = json.Unmarshal(export, &settings)
	s.Require().NoError(err)
	assert.NotNil(s.T(), settings)
}

// TestFullDataExport 测试完整数据导出
func (s *DataExportTest) TestFullDataExport() {
	// 创建用户并生成数据
	user := s.CreateVerifiedUser("github")
	s.FundAccount(user.Address, 100000000)

	// 创建任务
	task := WorkflowTask{Title: "Full Export Test", Budget: 1000000}
	s.CreateWorkflowTask(user.Address, task)

	// 导出所有数据
	export, err := s.ExportAllUserData(user.Address)
	s.Require().NoError(err)

	// 验证结构
	var fullExport map[string]interface{}
	err = json.Unmarshal(export, &fullExport)
	s.Require().NoError(err)

	// 验证包含所有字段
	assert.Contains(s.T(), fullExport, "user_id")
	assert.Contains(s.T(), fullExport, "export_time")
	assert.Contains(s.T(), fullExport, "transactions")
	assert.Contains(s.T(), fullExport, "tasks")
	assert.Contains(s.T(), fullExport, "reviews")
	assert.Contains(s.T(), fullExport, "settings")
}

// TestExportFormatCSV 测试 CSV 格式导出
func (s *DataExportTest) TestExportFormatCSV() {
	// 创建用户
	user := s.CreateVerifiedUser("github")
	s.FundAccount(user.Address, 100000000)
	s.Transfer(user.Address, s.CreateVerifiedUser("github").Address, 1000000)

	// 导出为 CSV
	export, err := s.ExportUserDataFormat(user.Address, "transactions", "csv")
	s.Require().NoError(err)

	// 验证 CSV 格式（简单检查）
	csvContent := string(export)
	assert.Contains(s.T(), csvContent, "tx_hash") // 应包含表头
	assert.Contains(s.T(), csvContent, "from")
	assert.Contains(s.T(), csvContent, "to")
	assert.Contains(s.T(), csvContent, "amount")
}

// TestAccountDeletionRequest 测试账户删除申请
func (s *DataExportTest) TestAccountDeletionRequest() {
	// 创建用户
	user := s.CreateVerifiedUser("github")

	// 申请删除账户
	requestID, err := s.RequestAccountDeletion(user.Address)
	s.Require().NoError(err)
	assert.NotEmpty(s.T(), requestID)

	// 验证状态为 pending
	status, err := s.GetAccountDeletionStatus(requestID)
	s.Require().NoError(err)
	assert.Equal(s.T(), "pending", status)
}

// TestAccountDeletionCoolingPeriod 测试账户删除冷静期
func (s *DataExportTest) TestAccountDeletionCoolingPeriod() {
	// 创建用户
	user := s.CreateVerifiedUser("github")

	// 申请删除
	requestID, _ := s.RequestAccountDeletion(user.Address)

	// 验证冷静期信息
	coolingPeriod, err := s.GetDeletionCoolingPeriod(requestID)
	s.Require().NoError(err)
	assert.Equal(s.T(), 30, coolingPeriod.Days) // 30 天冷静期
	assert.False(s.T(), coolingPeriod.CanDelete)

	// 验证可以取消
	canCancel, err := s.CanCancelDeletion(requestID)
	s.Require().NoError(err)
	assert.True(s.T(), canCancel)
}

// TestCancelAccountDeletion 测试取消账户删除
func (s *DataExportTest) TestCancelAccountDeletion() {
	// 创建用户
	user := s.CreateVerifiedUser("github")

	// 申请删除
	requestID, _ := s.RequestAccountDeletion(user.Address)

	// 取消删除
	err := s.CancelAccountDeletion(user.Address, requestID)
	s.Require().NoError(err)

	// 验证状态已取消
	status, _ := s.GetAccountDeletionStatus(requestID)
	assert.Equal(s.T(), "cancelled", status)
}

// TestAccountDeletionDataHandling 测试删除后数据处理
func (s *DataExportTest) TestAccountDeletionDataHandling() {
	// 创建用户并生成数据
	user := s.CreateVerifiedUser("github")

	// 存储一些数据
	s.FundAccount(user.Address, 10000000)

	// 模拟删除完成（跳过冷静期）
	err := s.CompleteAccountDeletion(user.Address)
	s.Require().NoError(err)

	// 验证链上数据只保留哈希
	data, err := s.GetDeletedUserData(user.Address)
	s.Require().NoError(err)
	assert.Empty(s.T(), data.PlainData, "明文数据应被清除")
	assert.NotEmpty(s.T(), data.Hash, "应保留哈希")
	assert.Empty(s.T(), data.PersonalInfo, "个人信息应被清除")
}

// Helper types

type DeletionCoolingPeriod struct {
	Days       int
	CanDelete  bool
	ExpiresAt  int64
}

type DeletedUserData struct {
	Address      string
	Hash         string
	PlainData    map[string]interface{}
	PersonalInfo map[string]interface{}
	DeletedAt    int64
}
