package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"sharetoken/x/taskmarket/types"
)

// MsgCreateTask Tests

func TestMsgCreateTask_ValidateBasic(t *testing.T) {
	tests := []struct {
		name    string
		msg     types.MsgCreateTask
		wantErr bool
	}{
		{
			name: "valid message",
			msg: types.MsgCreateTask{
				Creator:     "creator-1",
				Title:       "Test Task",
				Budget:      1000,
				TaskTypeVal: types.TaskTypeOpen,
			},
			wantErr: false,
		},
		{
			name: "invalid - empty creator",
			msg: types.MsgCreateTask{
				Creator:     "",
				Title:       "Test Task",
				Budget:      1000,
				TaskTypeVal: types.TaskTypeOpen,
			},
			wantErr: true,
		},
		{
			name: "invalid - empty title",
			msg: types.MsgCreateTask{
				Creator:     "creator-1",
				Title:       "",
				Budget:      1000,
				TaskTypeVal: types.TaskTypeOpen,
			},
			wantErr: true,
		},
		{
			name: "invalid - zero budget",
			msg: types.MsgCreateTask{
				Creator:     "creator-1",
				Title:       "Test Task",
				Budget:      0,
				TaskTypeVal: types.TaskTypeOpen,
			},
			wantErr: true,
		},
		{
			name: "invalid - invalid task type",
			msg: types.MsgCreateTask{
				Creator:     "creator-1",
				Title:       "Test Task",
				Budget:      1000,
				TaskTypeVal: "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgCreateTask_GetSigners(t *testing.T) {
	validAddress := sdk.AccAddress([]byte("test_address_1")).String()
	invalidAddress := "invalid"

	tests := []struct {
		name         string
		creator      string
		expectNil    bool
	}{
		{
			name:      "valid address",
			creator:   validAddress,
			expectNil: false,
		},
		{
			name:      "invalid address",
			creator:   invalidAddress,
			expectNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := types.NewMsgCreateTask(tt.creator, "Test", "Description", types.TaskTypeOpen, types.CategoryDevelopment, 1000)
			signers := msg.GetSigners()
			if tt.expectNil {
				require.Nil(t, signers)
			} else {
				require.Len(t, signers, 1)
			}
		})
	}
}

func TestMsgCreateTask_RouteAndType(t *testing.T) {
	msg := types.NewMsgCreateTask("creator-1", "Test", "Description", types.TaskTypeOpen, types.CategoryDevelopment, 1000)
	require.Equal(t, types.RouterKey, msg.Route())
	require.Equal(t, types.TypeMsgCreateTask, msg.Type())
}

// MsgUpdateTask Tests

func TestMsgUpdateTask_ValidateBasic(t *testing.T) {
	tests := []struct {
		name    string
		msg     types.MsgUpdateTask
		wantErr bool
	}{
		{
			name: "valid message",
			msg: types.MsgUpdateTask{
				Creator: "creator-1",
				TaskID:  "task-1",
			},
			wantErr: false,
		},
		{
			name: "invalid - empty creator",
			msg: types.MsgUpdateTask{
				Creator: "",
				TaskID:  "task-1",
			},
			wantErr: true,
		},
		{
			name: "invalid - empty task ID",
			msg: types.MsgUpdateTask{
				Creator: "creator-1",
				TaskID:  "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgUpdateTask_GetSigners(t *testing.T) {
	validAddress := sdk.AccAddress([]byte("test_address_1")).String()
	msg := types.NewMsgUpdateTask(validAddress, "task-1")
	signers := msg.GetSigners()
	require.Len(t, signers, 1)
}

// MsgPublishTask Tests

func TestMsgPublishTask_ValidateBasic(t *testing.T) {
	tests := []struct {
		name    string
		msg     types.MsgPublishTask
		wantErr bool
	}{
		{
			name: "valid message",
			msg: types.MsgPublishTask{
				Creator: "creator-1",
				TaskID:  "task-1",
			},
			wantErr: false,
		},
		{
			name: "invalid - empty creator",
			msg: types.MsgPublishTask{
				Creator: "",
				TaskID:  "task-1",
			},
			wantErr: true,
		},
		{
			name: "invalid - empty task ID",
			msg: types.MsgPublishTask{
				Creator: "creator-1",
				TaskID:  "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// MsgCancelTask Tests

func TestMsgCancelTask_ValidateBasic(t *testing.T) {
	tests := []struct {
		name    string
		msg     types.MsgCancelTask
		wantErr bool
	}{
		{
			name: "valid message",
			msg: types.MsgCancelTask{
				Creator: "creator-1",
				TaskID:  "task-1",
			},
			wantErr: false,
		},
		{
			name: "invalid - empty creator",
			msg: types.MsgCancelTask{
				Creator: "",
				TaskID:  "task-1",
			},
			wantErr: true,
		},
		{
			name: "invalid - empty task ID",
			msg: types.MsgCancelTask{
				Creator: "creator-1",
				TaskID:  "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// MsgStartTask Tests

func TestMsgStartTask_ValidateBasic(t *testing.T) {
	tests := []struct {
		name    string
		msg     types.MsgStartTask
		wantErr bool
	}{
		{
			name: "valid message",
			msg: types.MsgStartTask{
				WorkerID: "worker-1",
				TaskID:   "task-1",
			},
			wantErr: false,
		},
		{
			name: "invalid - empty worker ID",
			msg: types.MsgStartTask{
				WorkerID: "",
				TaskID:   "task-1",
			},
			wantErr: true,
		},
		{
			name: "invalid - empty task ID",
			msg: types.MsgStartTask{
				WorkerID: "worker-1",
				TaskID:   "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// MsgSubmitApplication Tests

func TestMsgSubmitApplication_ValidateBasic(t *testing.T) {
	tests := []struct {
		name    string
		msg     types.MsgSubmitApplication
		wantErr bool
	}{
		{
			name: "valid message",
			msg: types.MsgSubmitApplication{
				WorkerID:      "worker-1",
				TaskID:        "task-1",
				ProposedPrice: 800,
			},
			wantErr: false,
		},
		{
			name: "invalid - empty worker ID",
			msg: types.MsgSubmitApplication{
				WorkerID:      "",
				TaskID:        "task-1",
				ProposedPrice: 800,
			},
			wantErr: true,
		},
		{
			name: "invalid - empty task ID",
			msg: types.MsgSubmitApplication{
				WorkerID:      "worker-1",
				TaskID:        "",
				ProposedPrice: 800,
			},
			wantErr: true,
		},
		{
			name: "invalid - zero price",
			msg: types.MsgSubmitApplication{
				WorkerID:      "worker-1",
				TaskID:        "task-1",
				ProposedPrice: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// MsgAcceptApplication Tests

func TestMsgAcceptApplication_ValidateBasic(t *testing.T) {
	tests := []struct {
		name    string
		msg     types.MsgAcceptApplication
		wantErr bool
	}{
		{
			name: "valid message",
			msg: types.MsgAcceptApplication{
				RequesterID:   "requester-1",
				ApplicationID: "app-1",
			},
			wantErr: false,
		},
		{
			name: "invalid - empty requester ID",
			msg: types.MsgAcceptApplication{
				RequesterID:   "",
				ApplicationID: "app-1",
			},
			wantErr: true,
		},
		{
			name: "invalid - empty application ID",
			msg: types.MsgAcceptApplication{
				RequesterID:   "requester-1",
				ApplicationID: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// MsgRejectApplication Tests

func TestMsgRejectApplication_ValidateBasic(t *testing.T) {
	tests := []struct {
		name    string
		msg     types.MsgRejectApplication
		wantErr bool
	}{
		{
			name: "valid message",
			msg: types.MsgRejectApplication{
				RequesterID:   "requester-1",
				ApplicationID: "app-1",
			},
			wantErr: false,
		},
		{
			name: "invalid - empty requester ID",
			msg: types.MsgRejectApplication{
				RequesterID:   "",
				ApplicationID: "app-1",
			},
			wantErr: true,
		},
		{
			name: "invalid - empty application ID",
			msg: types.MsgRejectApplication{
				RequesterID:   "requester-1",
				ApplicationID: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// MsgSubmitBid Tests

func TestMsgSubmitBid_ValidateBasic(t *testing.T) {
	tests := []struct {
		name    string
		msg     types.MsgSubmitBid
		wantErr bool
	}{
		{
			name: "valid message",
			msg: types.MsgSubmitBid{
				WorkerID: "worker-1",
				TaskID:   "task-1",
				Amount:   500,
			},
			wantErr: false,
		},
		{
			name: "invalid - empty worker ID",
			msg: types.MsgSubmitBid{
				WorkerID: "",
				TaskID:   "task-1",
				Amount:   500,
			},
			wantErr: true,
		},
		{
			name: "invalid - empty task ID",
			msg: types.MsgSubmitBid{
				WorkerID: "worker-1",
				TaskID:   "",
				Amount:   500,
			},
			wantErr: true,
		},
		{
			name: "invalid - zero amount",
			msg: types.MsgSubmitBid{
				WorkerID: "worker-1",
				TaskID:   "task-1",
				Amount:   0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// MsgCloseAuction Tests

func TestMsgCloseAuction_ValidateBasic(t *testing.T) {
	tests := []struct {
		name    string
		msg     types.MsgCloseAuction
		wantErr bool
	}{
		{
			name: "valid message",
			msg: types.MsgCloseAuction{
				RequesterID: "requester-1",
				TaskID:      "task-1",
			},
			wantErr: false,
		},
		{
			name: "invalid - empty requester ID",
			msg: types.MsgCloseAuction{
				RequesterID: "",
				TaskID:      "task-1",
			},
			wantErr: true,
		},
		{
			name: "invalid - empty task ID",
			msg: types.MsgCloseAuction{
				RequesterID: "requester-1",
				TaskID:      "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// MsgSubmitMilestone Tests

func TestMsgSubmitMilestone_ValidateBasic(t *testing.T) {
	tests := []struct {
		name    string
		msg     types.MsgSubmitMilestone
		wantErr bool
	}{
		{
			name: "valid message",
			msg: types.MsgSubmitMilestone{
				WorkerID:    "worker-1",
				TaskID:      "task-1",
				MilestoneID: "milestone-1",
			},
			wantErr: false,
		},
		{
			name: "invalid - empty worker ID",
			msg: types.MsgSubmitMilestone{
				WorkerID:    "",
				TaskID:      "task-1",
				MilestoneID: "milestone-1",
			},
			wantErr: true,
		},
		{
			name: "invalid - empty task ID",
			msg: types.MsgSubmitMilestone{
				WorkerID:    "worker-1",
				TaskID:      "",
				MilestoneID: "milestone-1",
			},
			wantErr: true,
		},
		{
			name: "invalid - empty milestone ID",
			msg: types.MsgSubmitMilestone{
				WorkerID:    "worker-1",
				TaskID:      "task-1",
				MilestoneID: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// MsgApproveMilestone Tests

func TestMsgApproveMilestone_ValidateBasic(t *testing.T) {
	tests := []struct {
		name    string
		msg     types.MsgApproveMilestone
		wantErr bool
	}{
		{
			name: "valid message",
			msg: types.MsgApproveMilestone{
				RequesterID: "requester-1",
				TaskID:      "task-1",
				MilestoneID: "milestone-1",
			},
			wantErr: false,
		},
		{
			name: "invalid - empty requester ID",
			msg: types.MsgApproveMilestone{
				RequesterID: "",
				TaskID:      "task-1",
				MilestoneID: "milestone-1",
			},
			wantErr: true,
		},
		{
			name: "invalid - empty task ID",
			msg: types.MsgApproveMilestone{
				RequesterID: "requester-1",
				TaskID:      "",
				MilestoneID: "milestone-1",
			},
			wantErr: true,
		},
		{
			name: "invalid - empty milestone ID",
			msg: types.MsgApproveMilestone{
				RequesterID: "requester-1",
				TaskID:      "task-1",
				MilestoneID: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// MsgRejectMilestone Tests

func TestMsgRejectMilestone_ValidateBasic(t *testing.T) {
	tests := []struct {
		name    string
		msg     types.MsgRejectMilestone
		wantErr bool
	}{
		{
			name: "valid message",
			msg: types.MsgRejectMilestone{
				RequesterID: "requester-1",
				TaskID:      "task-1",
				MilestoneID: "milestone-1",
			},
			wantErr: false,
		},
		{
			name: "invalid - empty requester ID",
			msg: types.MsgRejectMilestone{
				RequesterID: "",
				TaskID:      "task-1",
				MilestoneID: "milestone-1",
			},
			wantErr: true,
		},
		{
			name: "invalid - empty task ID",
			msg: types.MsgRejectMilestone{
				RequesterID: "requester-1",
				TaskID:      "",
				MilestoneID: "milestone-1",
			},
			wantErr: true,
		},
		{
			name: "invalid - empty milestone ID",
			msg: types.MsgRejectMilestone{
				RequesterID: "requester-1",
				TaskID:      "task-1",
				MilestoneID: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// MsgSubmitRating Tests

func TestMsgSubmitRating_ValidateBasic(t *testing.T) {
	tests := []struct {
		name    string
		msg     types.MsgSubmitRating
		wantErr bool
	}{
		{
			name: "valid message",
			msg: types.MsgSubmitRating{
				TaskID:  "task-1",
				RaterID: "rater-1",
				RatedID: "rated-1",
				Ratings: map[string]int{"quality": 5},
			},
			wantErr: false,
		},
		{
			name: "invalid - empty task ID",
			msg: types.MsgSubmitRating{
				TaskID:  "",
				RaterID: "rater-1",
				RatedID: "rated-1",
				Ratings: map[string]int{"quality": 5},
			},
			wantErr: true,
		},
		{
			name: "invalid - empty rater ID",
			msg: types.MsgSubmitRating{
				TaskID:  "task-1",
				RaterID: "",
				RatedID: "rated-1",
				Ratings: map[string]int{"quality": 5},
			},
			wantErr: true,
		},
		{
			name: "invalid - empty rated ID",
			msg: types.MsgSubmitRating{
				TaskID:  "task-1",
				RaterID: "rater-1",
				RatedID: "",
				Ratings: map[string]int{"quality": 5},
			},
			wantErr: true,
		},
		{
			name: "invalid - no ratings",
			msg: types.MsgSubmitRating{
				TaskID:  "task-1",
				RaterID: "rater-1",
				RatedID: "rated-1",
				Ratings: map[string]int{},
			},
			wantErr: true,
		},
		{
			name: "invalid - nil ratings",
			msg: types.MsgSubmitRating{
				TaskID:  "task-1",
				RaterID: "rater-1",
				RatedID: "rated-1",
				Ratings: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
