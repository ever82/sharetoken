package types

// QueryRequest and QueryResponse types for taskmarket module

// QueryTasksRequest is the request type for querying tasks
type QueryTasksRequest struct {
	Status   string `json:"status,omitempty"`
	Category string `json:"category,omitempty"`
	RequesterID string `json:"requester_id,omitempty"`
	WorkerID    string `json:"worker_id,omitempty"`
	Limit       int    `json:"limit,omitempty"`
	Offset      int    `json:"offset,omitempty"`
}

// QueryTasksResponse is the response type for querying tasks
type QueryTasksResponse struct {
	Tasks      []Task `json:"tasks"`
	TotalCount int    `json:"total_count"`
}

// QueryTaskRequest is the request type for querying a single task
type QueryTaskRequest struct {
	TaskID string `json:"task_id"`
}

// QueryTaskResponse is the response type for querying a single task
type QueryTaskResponse struct {
	Task Task `json:"task"`
}

// QueryApplicationsRequest is the request type for querying applications
type QueryApplicationsRequest struct {
	TaskID   string `json:"task_id,omitempty"`
	WorkerID string `json:"worker_id,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Offset   int    `json:"offset,omitempty"`
}

// QueryApplicationsResponse is the response type for querying applications
type QueryApplicationsResponse struct {
	Applications []Application `json:"applications"`
	TotalCount   int           `json:"total_count"`
}

// QueryAuctionRequest is the request type for querying an auction
type QueryAuctionRequest struct {
	TaskID string `json:"task_id"`
}

// QueryAuctionResponse is the response type for querying an auction
type QueryAuctionResponse struct {
	Auction Auction `json:"auction"`
}

// QueryBidsRequest is the request type for querying bids
type QueryBidsRequest struct {
	TaskID   string `json:"task_id,omitempty"`
	WorkerID string `json:"worker_id,omitempty"`
}

// QueryBidsResponse is the response type for querying bids
type QueryBidsResponse struct {
	Bids []Bid `json:"bids"`
}

// QueryReputationRequest is the request type for querying reputation
type QueryReputationRequest struct {
	UserID string `json:"user_id"`
}

// QueryReputationResponse is the response type for querying reputation
type QueryReputationResponse struct {
	Reputation Reputation `json:"reputation"`
}

// QueryRatingsRequest is the request type for querying ratings
type QueryRatingsRequest struct {
	UserID string `json:"user_id,omitempty"` // User being rated
	TaskID string `json:"task_id,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
}

// QueryRatingsResponse is the response type for querying ratings
type QueryRatingsResponse struct {
	Ratings    []Rating `json:"ratings"`
	TotalCount int      `json:"total_count"`
}

// QueryStatisticsRequest is the request type for querying statistics
type QueryStatisticsRequest struct{}

// QueryStatisticsResponse is the response type for querying statistics
type QueryStatisticsResponse struct {
	TotalTasks        int `json:"total_tasks"`
	OpenTasks         int `json:"open_tasks"`
	AssignedTasks     int `json:"assigned_tasks"`
	InProgressTasks   int `json:"in_progress_tasks"`
	CompletedTasks    int `json:"completed_tasks"`
	CancelledTasks    int `json:"cancelled_tasks"`
	TotalApplications int `json:"total_applications"`
	TotalBids         int `json:"total_bids"`
	TotalRatings      int `json:"total_ratings"`
}
