package keeper

import (
	"sharetoken/x/taskmarket/types"
)

// Querier is a type alias for Keeper to separate query functions
type Querier struct {
	Keeper
}

// NewQuerier returns a new Querier instance
func NewQuerier(k Keeper) Querier {
	return Querier{Keeper: k}
}

// QueryServer implements the gRPC query server interface
type QueryServer struct {
	Keeper
}

// NewQueryServer creates a new QueryServer instance
func NewQueryServer(k Keeper) QueryServer {
	return QueryServer{Keeper: k}
}

// GetTask queries a task by ID
func (q Querier) GetTask(req *types.QueryTaskRequest) (*types.QueryTaskResponse, error) {
	task := q.Keeper.legacyKeeper.GetTask(req.TaskID)
	if task == nil {
		return nil, types.ErrTaskNotFound
	}
	return &types.QueryTaskResponse{Task: *task}, nil
}

// GetTasks queries tasks with filters
func (q Querier) GetTasks(req *types.QueryTasksRequest) (*types.QueryTasksResponse, error) {
	var tasks []types.Task

	if req.Status != "" {
		// Filter by status
		allTasks := q.Keeper.legacyKeeper.GetAllTasks()
		for _, task := range allTasks {
			if string(task.Status) == req.Status {
				tasks = append(tasks, *task)
			}
		}
	} else if req.RequesterID != "" {
		taskList := q.Keeper.legacyKeeper.GetTasksByRequester(req.RequesterID)
		for _, t := range taskList {
			tasks = append(tasks, *t)
		}
	} else if req.WorkerID != "" {
		taskList := q.Keeper.legacyKeeper.GetTasksByWorker(req.WorkerID)
		for _, t := range taskList {
			tasks = append(tasks, *t)
		}
	} else {
		allTasks := q.Keeper.legacyKeeper.GetAllTasks()
		for _, task := range allTasks {
			tasks = append(tasks, *task)
		}
	}

	// Apply pagination
	totalCount := len(tasks)
	if req.Offset < 0 {
		req.Offset = 0
	}
	if req.Limit <= 0 {
		req.Limit = 100
	}
	end := req.Offset + req.Limit
	if end > len(tasks) {
		end = len(tasks)
	}
	if req.Offset < len(tasks) {
		tasks = tasks[req.Offset:end]
	} else {
		tasks = []types.Task{}
	}

	return &types.QueryTasksResponse{
		Tasks:      tasks,
		TotalCount: totalCount,
	}, nil
}

// GetApplications queries applications
func (q Querier) GetApplications(req *types.QueryApplicationsRequest) (*types.QueryApplicationsResponse, error) {
	var apps []*types.Application

	if req.TaskID != "" {
		apps = q.Keeper.legacyKeeper.GetApplicationsByTask(req.TaskID)
	}

	// Convert to slice of values
	var result []types.Application
	for _, app := range apps {
		result = append(result, *app)
	}

	return &types.QueryApplicationsResponse{
		Applications: result,
		TotalCount:   len(result),
	}, nil
}

// GetAuction queries an auction
func (q Querier) GetAuction(req *types.QueryAuctionRequest) (*types.QueryAuctionResponse, error) {
	auction := q.Keeper.legacyKeeper.GetAuction(req.TaskID)
	if auction == nil {
		return nil, types.ErrAuctionNotFound
	}
	return &types.QueryAuctionResponse{Auction: *auction}, nil
}

// GetReputation queries a user's reputation
func (q Querier) GetReputation(req *types.QueryReputationRequest) (*types.QueryReputationResponse, error) {
	rep := q.Keeper.legacyKeeper.GetReputation(req.UserID)
	if rep == nil {
		rep = types.NewReputation(req.UserID)
	}
	return &types.QueryReputationResponse{Reputation: *rep}, nil
}

// GetStatistics returns marketplace statistics
func (q Querier) GetStatistics(req *types.QueryStatisticsRequest) (*types.QueryStatisticsResponse, error) {
	stats := q.Keeper.legacyKeeper.GetTaskStatistics()

	return &types.QueryStatisticsResponse{
		TotalTasks:        stats["total_tasks"].(int),
		OpenTasks:         stats["open_tasks"].(int),
		AssignedTasks:     stats["assigned_tasks"].(int),
		InProgressTasks:   stats["in_progress_tasks"].(int),
		CompletedTasks:    stats["completed_tasks"].(int),
		CancelledTasks:    stats["cancelled_tasks"].(int),
		TotalApplications: stats["total_applications"].(int),
		TotalBids:         stats["total_bids"].(int),
		TotalRatings:      stats["total_ratings"].(int),
	}, nil
}
