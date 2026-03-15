package types_test

import (
	"testing"

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
				Creator:  "creator-1",
				Title:    "Test Task",
				Budget:   1000,
				TaskType: types.TaskTypeOpen,
			},
			wantErr: false,
		},
		{
			name: "invalid - empty creator",
			msg: types.MsgCreateTask{
				Creator:  "",
				Title:    "Test Task",
				Budget:   1000,
				TaskType: types.TaskTypeOpen,
			},
			wantErr: true,
		},
		{
			name: "invalid - empty title",
			msg: types.MsgCreateTask{
				Creator:  "creator-1",
				Title:    "",
				Budget:   1000,
				TaskType: types.TaskTypeOpen,
			},
			wantErr: true,
		},
		{
			name: "invalid - zero budget",
			msg: types.MsgCreateTask{
				Creator:  "creator-1",
				Title:    "Test Task",
				Budget:   0,
				TaskType: types.TaskTypeOpen,
			},
			wantErr: true,
		},
		{
			name: "invalid - invalid task type",
			msg: types.MsgCreateTask{
				Creator:  "creator-1",
				Title:    "Test Task",
				Budget:   1000,
				TaskType: types.TaskType_TASK_TYPE_UNSPECIFIED,
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
	msg := types.MsgCreateTask{
		Creator: "creator-1",
		Title:   "Test",
		Budget:  1000,
	}
	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, "creator-1", signers[0])
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
				TaskId:  "task-1",
			},
			wantErr: false,
		},
		{
			name: "invalid - empty creator",
			msg: types.MsgUpdateTask{
				Creator: "",
				TaskId:  "task-1",
			},
			wantErr: true,
		},
		{
			name: "invalid - empty task ID",
			msg: types.MsgUpdateTask{
				Creator: "creator-1",
				TaskId:  "",
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
	msg := types.MsgUpdateTask{
		Creator: "creator-1",
		TaskId:  "task-1",
	}
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
				TaskId:  "task-1",
			},
			wantErr: false,
		},
		{
			name: "invalid - empty creator",
			msg: types.MsgPublishTask{
				Creator: "",
				TaskId:  "task-1",
			},
			wantErr: true,
		},
		{
			name: "invalid - empty task ID",
			msg: types.MsgPublishTask{
				Creator: "creator-1",
				TaskId:  "",
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
				TaskId:  "task-1",
			},
			wantErr: false,
		},
		{
			name: "invalid - empty creator",
			msg: types.MsgCancelTask{
				Creator: "",
				TaskId:  "task-1",
			},
			wantErr: true,
		},
		{
			name: "invalid - empty task ID",
			msg: types.MsgCancelTask{
				Creator: "creator-1",
				TaskId:  "",
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
				WorkerId: "worker-1",
				TaskId:   "task-1",
			},
			wantErr: false,
		},
		{
			name: "invalid - empty worker ID",
			msg: types.MsgStartTask{
				WorkerId: "",
				TaskId:   "task-1",
			},
			wantErr: true,
		},
		{
			name: "invalid - empty task ID",
			msg: types.MsgStartTask{
				WorkerId: "worker-1",
				TaskId:   "",
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
				WorkerId:      "worker-1",
				TaskId:        "task-1",
				ProposedPrice: 800,
			},
			wantErr: false,
		},
		{
			name: "invalid - empty worker ID",
			msg: types.MsgSubmitApplication{
				WorkerId:      "",
				TaskId:        "task-1",
				ProposedPrice: 800,
			},
			wantErr: true,
		},
		{
			name: "invalid - empty task ID",
			msg: types.MsgSubmitApplication{
				WorkerId:      "worker-1",
				TaskId:        "",
				ProposedPrice: 800,
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
				RequesterId:   "requester-1",
				ApplicationId: "app-1",
			},
			wantErr: false,
		},
		{
			name: "invalid - empty requester ID",
			msg: types.MsgAcceptApplication{
				RequesterId:   "",
				ApplicationId: "app-1",
			},
			wantErr: true,
		},
		{
			name: "invalid - empty application ID",
			msg: types.MsgAcceptApplication{
				RequesterId:   "requester-1",
				ApplicationId: "",
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
				RequesterId:   "requester-1",
				ApplicationId: "app-1",
			},
			wantErr: false,
		},
		{
			name: "invalid - empty requester ID",
			msg: types.MsgRejectApplication{
				RequesterId:   "",
				ApplicationId: "app-1",
			},
			wantErr: true,
		},
		{
			name: "invalid - empty application ID",
			msg: types.MsgRejectApplication{
				RequesterId:   "requester-1",
				ApplicationId: "",
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
				WorkerId: "worker-1",
				TaskId:   "task-1",
				Amount:   500,
			},
			wantErr: false,
		},
		{
			name: "invalid - empty worker ID",
			msg: types.MsgSubmitBid{
				WorkerId: "",
				TaskId:   "task-1",
				Amount:   500,
			},
			wantErr: true,
		},
		{
			name: "invalid - empty task ID",
			msg: types.MsgSubmitBid{
				WorkerId: "worker-1",
				TaskId:   "",
				Amount:   500,
			},
			wantErr: true,
		},
		{
			name: "invalid - zero amount",
			msg: types.MsgSubmitBid{
				WorkerId: "worker-1",
				TaskId:   "task-1",
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
				RequesterId: "requester-1",
				TaskId:      "task-1",
			},
			wantErr: false,
		},
		{
			name: "invalid - empty requester ID",
			msg: types.MsgCloseAuction{
				RequesterId: "",
				TaskId:      "task-1",
			},
			wantErr: true,
		},
		{
			name: "invalid - empty task ID",
			msg: types.MsgCloseAuction{
				RequesterId: "requester-1",
				TaskId:      "",
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
