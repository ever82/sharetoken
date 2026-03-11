package types

import (
	"fmt"
	"strings"
	"time"
)

// RatingDimension represents a rating dimension
type RatingDimension string

const (
	DimensionQuality         RatingDimension = "quality"         // Quality of work
	DimensionCommunication   RatingDimension = "communication"   // Communication
	DimensionTimeliness      RatingDimension = "timeliness"      // Timeliness
	DimensionProfessionalism RatingDimension = "professionalism" // Professionalism
)

// GetAllDimensions returns all rating dimensions
func GetAllDimensions() []RatingDimension {
	return []RatingDimension{
		DimensionQuality,
		DimensionCommunication,
		DimensionTimeliness,
		DimensionProfessionalism,
	}
}

// Rating represents a multi-dimensional rating
type Rating struct {
	ID        string                  `json:"id"`
	TaskID    string                  `json:"task_id"`
	RaterID   string                  `json:"rater_id"` // Who gave the rating
	RatedID   string                  `json:"rated_id"` // Who was rated
	Ratings   map[RatingDimension]int `json:"ratings"`  // 1-5 for each dimension
	Comment   string                  `json:"comment"`
	CreatedAt int64                   `json:"created_at"`
}

// NewRating creates a new rating
func NewRating(id, taskID, raterID, ratedID string) *Rating {
	return &Rating{
		ID:        id,
		TaskID:    taskID,
		RaterID:   raterID,
		RatedID:   ratedID,
		Ratings:   make(map[RatingDimension]int),
		CreatedAt: time.Now().Unix(),
	}
}

// SetRating sets a rating for a dimension
func (r *Rating) SetRating(dimension RatingDimension, value int) error {
	if value < 1 || value > 5 {
		return fmt.Errorf("rating must be between 1 and 5")
	}
	r.Ratings[dimension] = value
	return nil
}

// GetRating gets a rating for a dimension
func (r *Rating) GetRating(dimension RatingDimension) int {
	if val, ok := r.Ratings[dimension]; ok {
		return val
	}
	return 0
}

// GetAverage returns the average rating
func (r *Rating) GetAverage() float64 {
	if len(r.Ratings) == 0 {
		return 0
	}

	var sum int
	for _, v := range r.Ratings {
		sum += v
	}

	return float64(sum) / float64(len(r.Ratings))
}

// Validate validates the rating
func (r *Rating) Validate() error {
	if r.TaskID == "" {
		return fmt.Errorf("task ID cannot be empty")
	}
	if r.RaterID == "" {
		return fmt.Errorf("rater ID cannot be empty")
	}
	if r.RatedID == "" {
		return fmt.Errorf("rated ID cannot be empty")
	}
	if len(r.Ratings) == 0 {
		return fmt.Errorf("at least one rating dimension required")
	}

	for dim, val := range r.Ratings {
		if val < 1 || val > 5 {
			return fmt.Errorf("rating for %s must be between 1 and 5", dim)
		}
	}

	return nil
}

// IsComplete checks if all dimensions are rated
func (r *Rating) IsComplete() bool {
	for _, dim := range GetAllDimensions() {
		if _, ok := r.Ratings[dim]; !ok {
			return false
		}
	}
	return true
}

// Reputation represents a user's reputation
type Reputation struct {
	UserID         string                      `json:"user_id"`
	TotalRatings   int                         `json:"total_ratings"`
	AverageRating  float64                     `json:"average_rating"`
	RatingsByDim   map[RatingDimension]float64 `json:"ratings_by_dimension"`
	CompletedTasks int                         `json:"completed_tasks"`
	DisputeRate    float64                     `json:"dispute_rate"`
	OnTimeDelivery float64                     `json:"on_time_delivery"` // percentage
	CreatedAt      int64                       `json:"created_at"`
	UpdatedAt      int64                       `json:"updated_at"`
}

// NewReputation creates a new reputation
func NewReputation(userID string) *Reputation {
	now := time.Now().Unix()
	return &Reputation{
		UserID:       userID,
		RatingsByDim: make(map[RatingDimension]float64),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// AddRating adds a rating to the reputation
func (r *Reputation) AddRating(rating *Rating) {
	// Update dimension averages
	for dim, val := range rating.Ratings {
		if existing, ok := r.RatingsByDim[dim]; ok {
			// Weighted average
			r.RatingsByDim[dim] = (existing*float64(r.TotalRatings) + float64(val)) / float64(r.TotalRatings+1)
		} else {
			r.RatingsByDim[dim] = float64(val)
		}
	}

	r.TotalRatings++
	r.calculateAverage()
	r.UpdatedAt = time.Now().Unix()
}

// calculateAverage calculates the overall average
func (r *Reputation) calculateAverage() {
	if len(r.RatingsByDim) == 0 {
		r.AverageRating = 0
		return
	}

	var sum float64
	for _, v := range r.RatingsByDim {
		sum += v
	}

	r.AverageRating = sum / float64(len(r.RatingsByDim))
}

// GetRatingForDimension gets the rating for a specific dimension
func (r *Reputation) GetRatingForDimension(dim RatingDimension) float64 {
	if val, ok := r.RatingsByDim[dim]; ok {
		return val
	}
	return 0
}

// UpdateTaskCompletion updates task completion stats
func (r *Reputation) UpdateTaskCompleted(onTime bool) {
	r.CompletedTasks++
	// Update on-time delivery rate
	if r.CompletedTasks == 1 {
		if onTime {
			r.OnTimeDelivery = 100.0
		} else {
			r.OnTimeDelivery = 0.0
		}
	} else {
		previousOnTime := int(r.OnTimeDelivery * float64(r.CompletedTasks-1) / 100)
		if onTime {
			r.OnTimeDelivery = float64(previousOnTime+1) / float64(r.CompletedTasks) * 100
		} else {
			r.OnTimeDelivery = float64(previousOnTime) / float64(r.CompletedTasks) * 100
		}
	}
	r.UpdatedAt = time.Now().Unix()
}

// UpdateDisputeRate updates dispute rate
func (r *Reputation) UpdateDisputeRate(totalTasks int, disputedTasks int) {
	if totalTasks > 0 {
		r.DisputeRate = float64(disputedTasks) / float64(totalTasks) * 100
	}
	r.UpdatedAt = time.Now().Unix()
}

// GetStars returns a star rating representation
func (r *Reputation) GetStars() string {
	stars := int(r.AverageRating + 0.5)
	// Performance: 预分配容量，5个Unicode字符 = 15字节
	var result strings.Builder
	result.Grow(15)
	for i := 0; i < 5; i++ {
		if i < stars {
			result.WriteString("★")
		} else {
			result.WriteString("☆")
		}
	}
	return result.String()
}

// IsTrusted checks if user has trusted reputation (> 4.0 average, > 5 ratings)
func (r *Reputation) IsTrusted() bool {
	return r.AverageRating >= 4.0 && r.TotalRatings >= 5 && r.DisputeRate < 10.0
}

// IsNew checks if user is new (< 3 ratings)
func (r *Reputation) IsNew() bool {
	return r.TotalRatings < 3
}

// GetTrustLevel returns trust level
func (r *Reputation) GetTrustLevel() string {
	if r.IsTrusted() {
		return "trusted"
	}
	if r.IsNew() {
		return "new"
	}
	if r.AverageRating >= 3.0 {
		return "established"
	}
	return "caution"
}

// RatingSummary provides a summary of ratings
type RatingSummary struct {
	UserID        string                      `json:"user_id"`
	Overall       float64                     `json:"overall"`
	TotalReviews  int                         `json:"total_reviews"`
	Dimensions    map[RatingDimension]float64 `json:"dimensions"`
	RecentReviews []Rating                    `json:"recent_reviews"`
}

// NewRatingSummary creates a rating summary from reputation
func NewRatingSummary(reputation *Reputation) *RatingSummary {
	return &RatingSummary{
		UserID:       reputation.UserID,
		Overall:      reputation.AverageRating,
		TotalReviews: reputation.TotalRatings,
		Dimensions:   reputation.RatingsByDim,
	}
}
