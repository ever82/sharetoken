package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ServiceDiscoveryTest 测试服务发现功能
type ServiceDiscoveryTest struct {
	E2ETestSuite
}

func TestServiceDiscovery(t *testing.T) {
	suite.Run(t, new(ServiceDiscoveryTest))
}

// TestBrowseByCategory 测试按类别浏览服务
func (s *ServiceDiscoveryTest) TestBrowseByCategory() {
	// 创建不同类别的服务
	categories := []string{"llm", "agent", "workflow"}

	for _, category := range categories {
		s.CreateTestService(category, "Test Service")
	}

	// 测试按类别浏览
	for _, category := range categories {
		s.Run(category, func() {
			services, err := s.BrowseServicesByCategory(category)
			s.Require().NoError(err)

			// 验证返回了该类别的服务
			assert.NotEmpty(s.T(), services, "%s 类别应有服务", category)
			for _, svc := range services {
				assert.Equal(s.T(), category, svc.Category)
			}
		})
	}
}

// TestSearchByKeyword 测试关键词搜索
func (s *ServiceDiscoveryTest) TestSearchByKeyword() {
	// 创建带关键词的服务
	s.CreateTestService("llm", "GPT-4 Translation Service")
	s.CreateTestService("llm", "GPT-4 Code Assistant")
	s.CreateTestService("agent", "Data Analysis Agent")

	// 测试搜索
	searchCases := []struct {
		keyword      string
		expectedMin  int
	}{
		{"translation", 1},
		{"GPT", 2},
		{"assistant", 1},
		{"analysis", 1},
		{"nonexistent", 0},
	}

	for _, tc := range searchCases {
		s.Run(tc.keyword, func() {
			results, err := s.SearchServices(tc.keyword)
			s.Require().NoError(err)

			// 验证搜索结果
			assert.GreaterOrEqual(s.T(), len(results), tc.expectedMin, "搜索 '%s' 应返回至少 %d 个结果", tc.keyword, tc.expectedMin)
		})
	}
}

// TestSortByPrice 测试按价格排序
func (s *ServiceDiscoveryTest) TestSortByPrice() {
	// 创建不同价格的服务
	prices := []int64{1000000, 500000, 2000000} // 1000, 500, 2000 STT
	for i, price := range prices {
		s.CreateTestServiceWithPrice("llm", "Service", price, i)
	}

	// 测试按价格升序
	services, err := s.SortServices("price", "asc")
	s.Require().NoError(err)

	// 验证排序
	for i := 1; i < len(services); i++ {
		assert.LessOrEqual(s.T(), services[i-1].Price, services[i].Price, "价格应升序排列")
	}

	// 测试按价格降序
	services, err = s.SortServices("price", "desc")
	s.Require().NoError(err)

	// 验证排序
	for i := 1; i < len(services); i++ {
		assert.GreaterOrEqual(s.T(), services[i-1].Price, services[i].Price, "价格应降序排列")
	}
}

// TestSortByRating 测试按评分排序
func (s *ServiceDiscoveryTest) TestSortByRating() {
	// 创建不同评分的服务
	ratings := []float64{4.5, 3.0, 5.0}
	for i, rating := range ratings {
		s.CreateTestServiceWithRating("llm", "Service", rating, i)
	}

	// 按评分排序
	services, err := s.SortServices("rating", "desc")
	s.Require().NoError(err)

	// 验证排序
	for i := 1; i < len(services); i++ {
		assert.GreaterOrEqual(s.T(), services[i-1].Rating, services[i].Rating, "评分应降序排列")
	}
}

// TestSortByResponseTime 测试按响应时间排序
func (s *ServiceDiscoveryTest) TestSortByResponseTime() {
	// 创建不同响应时间的服务
	responseTimes := []int64{1000, 500, 2000} // ms
	for i, rt := range responseTimes {
		s.CreateTestServiceWithResponseTime("llm", "Service", rt, i)
	}

	// 按响应时间排序（升序：最快优先）
	services, err := s.SortServices("response_time", "asc")
	s.Require().NoError(err)

	// 验证排序
	for i := 1; i < len(services); i++ {
		assert.LessOrEqual(s.T(), services[i-1].ResponseTime, services[i].ResponseTime, "响应时间应升序排列")
	}
}

// TestViewServiceDetails 测试查看服务详情
func (s *ServiceDiscoveryTest) TestViewServiceDetails() {
	// 创建服务
	serviceID := s.CreateTestService("llm", "Detailed Service")

	// 查询详情
	details, err := s.GetServiceDetails(serviceID)
	s.Require().NoError(err)

	// 验证详情字段
	assert.NotEmpty(s.T(), details.ID)
	assert.NotEmpty(s.T(), details.Name)
	assert.NotEmpty(s.T(), details.Description)
	assert.NotEmpty(s.T(), details.Category)
	assert.Greater(s.T(), details.Price, int64(0))
	assert.GreaterOrEqual(s.T(), details.Rating, float64(0))
	assert.NotNil(s.T(), details.Cases)
}

// TestViewServiceReviews 测试查看服务评价
func (s *ServiceDiscoveryTest) TestViewServiceReviews() {
	// 创建服务和评价
	serviceID := s.CreateTestService("llm", "Reviewed Service")
	s.CreateTestReview(serviceID, "user1", 5, "Great service!")
	s.CreateTestReview(serviceID, "user2", 4, "Good but slow")

	// 查询评价
	reviews, err := s.GetServiceReviews(serviceID)
	s.Require().NoError(err)

	// 验证评价
	assert.GreaterOrEqual(s.T(), len(reviews), 2, "应至少有 2 条评价")
	for _, review := range reviews {
		assert.GreaterOrEqual(s.T(), review.Rating, 1)
		assert.LessOrEqual(s.T(), review.Rating, 5)
		assert.NotEmpty(s.T(), review.Comment)
	}
}

// TestFavoriteService 测试收藏服务
func (s *ServiceDiscoveryTest) TestFavoriteService() {
	// 创建用户和服务
	user := s.CreateTestUser()
	serviceID := s.CreateTestService("llm", "Favorite Service")

	// 收藏服务
	err := s.FavoriteService(user.Address, serviceID)
	s.Require().NoError(err)

	// 查询收藏列表
	favorites, err := s.GetFavoriteServices(user.Address)
	s.Require().NoError(err)

	// 验证收藏
	found := false
	for _, fav := range favorites {
		if fav.ID == serviceID {
			found = true
			break
		}
	}
	assert.True(s.T(), found, "服务应在收藏列表中")

	// 取消收藏
	err = s.UnfavoriteService(user.Address, serviceID)
	s.Require().NoError(err)

	// 验证已取消
	favorites, err = s.GetFavoriteServices(user.Address)
	s.Require().NoError(err)

	found = false
	for _, fav := range favorites {
		if fav.ID == serviceID {
			found = true
			break
		}
	}
	assert.False(s.T(), found, "服务不应再在收藏列表中")
}

// TestPagination 测试分页
func (s *ServiceDiscoveryTest) TestPagination() {
	// 创建多个服务
	for i := 0; i < 25; i++ {
		s.CreateTestService("llm", "Paged Service")
	}

	// 测试分页
	page1, err := s.ListServicesWithPagination(1, 10)
	s.Require().NoError(err)
	assert.Equal(s.T(), 10, len(page1), "第一页应有 10 条")

	page2, err := s.ListServicesWithPagination(2, 10)
	s.Require().NoError(err)
	assert.Equal(s.T(), 10, len(page2), "第二页应有 10 条")

	page3, err := s.ListServicesWithPagination(3, 10)
	s.Require().NoError(err)
	assert.Equal(s.T(), 5, len(page3), "第三页应有 5 条")
}
