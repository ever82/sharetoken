package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"sharetoken/test/helpers"
)

// TaskTrackingTestSuite 测试 ACH-USER-006: Task Progress Tracking
type TaskTrackingTestSuite struct {
	suite.Suite
	chain *helpers.ChainHelper
	ctx   context.Context
}

func (s *TaskTrackingTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.chain = helpers.NewChainHelper(s.T())
}

// TestTaskListPagination 测试任务列表分页
func (s *TaskTrackingTestSuite) TestTaskListPagination() {
	user := s.chain.CreateAccountWithBalance("task_user", "5000stt")

	// 查询任务列表
	tasks, err := s.chain.QueryUserTasks(user.Address, 1, 10)
	require.NoError(s.T(), err, "任务列表查询应成功")

	// 验证每页数量
	require.LessOrEqual(s.T(), len(tasks), 10, "每页应不超过10条")

	s.T().Logf("任务列表分页测试通过")
}

// TestTaskSorting 测试任务排序
func (s *TaskTrackingTestSuite) TestTaskSorting() {
	user := s.chain.CreateAccountWithBalance("sort_user", "5000stt")

	// 按时间排序
	tasks, err := s.chain.QueryUserTasksSorted(user.Address, "time", "desc")
	require.NoError(s.T(), err)
	require.Greater(s.T(), len(tasks), 0)

	// 按优先级排序
	tasks, err = s.chain.QueryUserTasksSorted(user.Address, "priority", "desc")
	require.NoError(s.T(), err)

	s.T().Log("任务排序测试通过")
}

// TestTaskStatusUpdate 测试任务状态更新 (< 10秒)
func (s *TaskTrackingTestSuite) TestTaskStatusUpdate() {
	user := s.chain.CreateAccountWithBalance("status_user", "5000stt")

	// 创建任务
	taskID, _ := s.chain.CreateTask(user, "test task")

	// 更新状态
	start := time.Now()
	err := s.chain.UpdateTaskStatus(taskID, "in_progress")
	require.NoError(s.T(), err)

	// 查询更新后的状态
	status, _ := s.chain.QueryTaskStatus(taskID)
	elapsed := time.Since(start)

	require.Equal(s.T(), "in_progress", status.State)
	require.Less(s.T(), elapsed.Seconds(), 10.0,
		"状态更新应在10秒内同步")

	s.T().Logf("状态更新用时: %.2f秒", elapsed.Seconds())
}

// TestTaskProgressPercentage 测试进度百分比精度 (1%)
func (s *TaskTrackingTestSuite) TestTaskProgressPercentage() {
	user := s.chain.CreateAccountWithBalance("progress_user", "5000stt")

	taskID, _ := s.chain.CreateTask(user, "progress task")

	// 设置不同进度
	progressValues := []int{0, 25, 50, 75, 100}
	for _, p := range progressValues {
		err := s.chain.UpdateTaskProgress(taskID, p)
		require.NoError(s.T(), err)

		status, _ := s.chain.QueryTaskStatus(taskID)
		require.Equal(s.T(), p, status.Progress)
	}

	s.T().Log("进度百分比精度测试通过")
}

// TestMilestoneTracking 测试里程碑完成情况
func (s *TaskTrackingTestSuite) TestMilestoneTracking() {
	user := s.chain.CreateAccountWithBalance("milestone_user", "5000stt")

	// 创建带里程碑的任务
	milestones := []string{"milestone1", "milestone2", "milestone3"}
	taskID, _ := s.chain.CreateTaskWithMilestones(user, "milestone task", milestones)

	// 查询里程碑
	ms, err := s.chain.QueryTaskMilestones(taskID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), len(milestones), len(ms))

	// 完成第一个里程碑
	err = s.chain.CompleteMilestone(taskID, "milestone1")
	require.NoError(s.T(), err)

	s.T().Log("里程碑追踪测试通过")
}

// TestNotificationDelivery 测试通知送达 (< 30秒)
func (s *TaskTrackingTestSuite) TestNotificationDelivery() {
	user := s.chain.CreateAccountWithBalance("notify_user", "5000stt")

	// 创建任务触发通知
	start := time.Now()
	_, _ = s.chain.CreateTask(user, "notification task")

	// 检查通知是否收到
	notifications, err := s.chain.QueryNotifications(user.Address, 30)
	elapsed := time.Since(start)

	require.NoError(s.T(), err)
	require.Greater(s.T(), len(notifications), 0, "应收到通知")
	require.Less(s.T(), elapsed.Seconds(), 30.0,
		"通知应在30秒内送达")

	s.T().Logf("通知送达用时: %.2f秒", elapsed.Seconds())
}

// TestTaskCommunication 测试任务沟通（留言）
func (s *TaskTrackingTestSuite) TestTaskCommunication() {
	user := s.chain.CreateAccountWithBalance("comm_user", "5000stt")
	_ = s.chain.CreateAccount("comm_provider")

	taskID, _ := s.chain.CreateTask(user, "communication task")

	// 发送留言
	err := s.chain.SendTaskMessage(taskID, user.Address, "请尽快开始")
	require.NoError(s.T(), err)

	// 查询留言
	messages, err := s.chain.QueryTaskMessages(taskID)
	require.NoError(s.T(), err)
	require.Greater(s.T(), len(messages), 0)

	s.T().Log("任务沟通功能测试通过")
}

func TestTaskTrackingSuite(t *testing.T) {
	suite.Run(t, new(TaskTrackingTestSuite))
}
