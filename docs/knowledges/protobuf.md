# Protocol Buffers 知识文档

## 1. 什么是 Protocol Buffers

Protocol Buffers (protobuf) 是 Google 开发的结构化数据序列化机制，比 XML 更小、更快、更简单。

### Cosmos SDK 为什么使用 Protobuf

- **性能优越**: 二进制格式，序列化/反序列化速度快，数据体积小
- **跨语言支持**: 原生支持 Go、TypeScript、Java、Python 等多种语言
- **类型安全**: 强类型定义，编译时检查，减少运行时错误
- **向后兼容**: 字段编号机制允许平滑升级，新旧版本可互操作
- **gRPC 集成**: 原生支持 gRPC 服务定义，自动生成客户端/服务端代码
- **状态机序列化**: Cosmos SDK 使用 protobuf 序列化区块链状态

## 2. Proto 文件结构与语法基础

### 文件结构示例

```protobuf
syntax = "proto3";                                    // 指定 proto3 语法

package sharetokens.trust.v1;                         // 包名，防止命名冲突

option go_package = "github.com/sharetokens/x/trust/types";  // Go 代码生成路径

import "google/protobuf/timestamp.proto";             // 导入标准类型
import "gogoproto/gogo.proto";                        // 导入 Cosmos 扩展
import "cosmos/msg/v1/msg.proto";                     // 导入 Cosmos 消息定义
```

### 基本数据类型

| Proto 类型 | Go 类型 | 说明 |
|-----------|---------|------|
| `string` | `string` | UTF-8 字符串 |
| `uint64` | `uint64` | 无符号 64 位整数 |
| `int64` | `int64` | 有符号 64 位整数 |
| `int32` | `int32` | 有符号 32 位整数 |
| `bool` | `bool` | 布尔值 |
| `bytes` | `[]byte` | 字节数组 |
| `double` | `float64` | 双精度浮点数 |

### 消息定义 (Message)

```protobuf
// MQRecord 存储用户的 MQ 信息
message MQRecord {
  // 用户地址 (字段编号 1)
  string address = 1;

  // 当前 MQ 分数 (字段编号 2)
  uint64 mq = 2;

  // MQ 变更历史 (repeated = 数组)
  repeated MQHistoryEntry history = 3;

  // 统计信息 (嵌套消息)
  MQStats stats = 4;

  // 创建时间 (使用 Google 标准类型)
  google.protobuf.Timestamp created_at = 5;
}
```

### 枚举定义 (Enum)

```protobuf
// DisputeStatus 定义争议的可能状态
enum DisputeStatus {
  DISPUTE_STATUS_FILED = 0;      // 刚提交
  DISPUTE_STATUS_MEDIATING = 1;  // AI 仲裁中
  DISPUTE_STATUS_SETTLED = 2;    // 已和解
  DISPUTE_STATUS_JURIED = 3;     // 升级到陪审团
  DISPUTE_STATUS_RESOLVED = 5;   // 已解决
  DISPUTE_STATUS_FINAL = 7;      // 最终状态
}
```

**注意**: 第一个枚举值必须为 0，作为默认值。

### oneof 联合类型

```protobuf
message MediationEvent {
  uint64 id = 1;
  google.protobuf.Timestamp timestamp = 2;
  string actor = 3;
  EventType type = 4;

  // 只能设置其中一个字段
  oneof content {
    MessageContent message = 5;
    EvidenceContent evidence = 6;
    ProposalContent proposal = 7;
    VerdictContent verdict = 9;
  }
}
```

### Map 类型

```protobuf
message JuryMember {
  string address = 1;
  uint64 mq = 2;

  // key: proposal_id, value: score
  map<uint64, int32> scores = 3;
}
```

## 3. 代码生成

### Go 代码生成

```bash
# 使用 buf 生成 (推荐)
buf generate

# 或使用 protoc 直接生成
protoc --go_out=. --go-grpc_out=. \
  -I proto \
  proto/sharetokens/trust/v1/*.proto
```

生成的 Go 代码示例:
```go
// 自动生成的结构体
type MQRecord struct {
    Address   string           `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
    Mq        uint64           `protobuf:"varint,2,opt,name=mq,proto3" json:"mq,omitempty"`
    History   []*MQHistoryEntry `protobuf:"bytes,3,rep,name=history,proto3" json:"history,omitempty"`
}
```

### TypeScript 代码生成

```bash
# 生成 TypeScript 类型
protoc --ts_out=. \
  -I proto \
  proto/sharetokens/trust/v1/*.proto
```

## 4. Cosmos SDK Proto 约定

### 目录结构规范

```
proto/
└── sharetokens/
    └── trust/
        └── v1/
            ├── trust.proto      # 核心类型定义
            ├── dispute.proto    # 争议相关类型
            ├── tx.proto         # 交易消息 (Msg service)
            ├── query.proto      # 查询服务 (Query service)
            └── genesis.proto    # 创世状态
```

### 交易消息 (Msg) 服务

```protobuf
// tx.proto - 定义区块链交易
service Msg {
  // 创建争议
  rpc CreateDispute(MsgCreateDispute) returns (MsgCreateDisputeResponse);

  // 发送消息
  rpc SendMessage(MsgSendMessage) returns (MsgSendMessageResponse);

  // 提交证据
  rpc SubmitEvidence(MsgSubmitEvidence) returns (MsgSubmitEvidenceResponse);
}

// 交易消息必须指定签名者
message MsgCreateDispute {
  option (cosmos.msg.v1.signer) = "plaintiff";  // 指定签名者字段

  string plaintiff = 1;   // 原告地址
  uint64 order_id = 2;    // 关联订单 ID
  string defendant = 3;   // 被告地址
  string title = 4;       // 争议标题
  string description = 5; // 详细描述
}

// 响应消息
message MsgCreateDisputeResponse {
  uint64 dispute_id = 1;  // 返回创建的争议 ID
}
```

### 查询服务 (Query)

```protobuf
// query.proto - 定义只读查询
service Query {
  // 查询单个争议
  rpc Dispute(QueryDisputeRequest) returns (QueryDisputeResponse);

  // 查询所有争议 (带分页)
  rpc Disputes(QueryDisputesRequest) returns (QueryDisputesResponse);

  // 查询用户 MQ
  rpc MQ(QueryMQRequest) returns (QueryMQResponse);
}

// 请求消息
message QueryDisputeRequest {
  uint64 dispute_id = 1;
}

// 响应消息
message QueryDisputeResponse {
  Dispute dispute = 1;
}

// 分页请求
message QueryDisputesRequest {
  DisputeStatus status = 1;                                    // 过滤条件
  cosmos.base.query.v1beta1.PageRequest pagination = 2;        // 分页参数
}

// 分页响应
message QueryDisputesResponse {
  repeated Dispute disputes = 1;
  cosmos.base.query.v1beta1.PageResponse pagination = 2;
}
```

### 创世状态 (Genesis)

```protobuf
// genesis.proto - 定义模块初始化状态
message GenesisState {
  // 模块配置
  MQConfig config = 1 [(gogoproto.nullable) = false];

  // MQ 记录列表
  repeated MQRecord mq_records = 3 [(gogoproto.nullable) = false];

  // 争议列表
  repeated Dispute disputes = 4 [(gogoproto.nullable) = false];

  // 下一个 ID (用于自增)
  uint64 next_dispute_id = 9;
}
```

### Gogo Proto 扩展

```protobuf
import "gogoproto/gogo.proto";

message GenesisState {
  // nullable = false: 生成值类型而非指针
  MQConfig config = 1 [(gogoproto.nullable) = false];

  // moretags: 添加额外的 Go struct tag
  string address = 2 [(gogoproto.moretags) = "yaml:\"address\""];
}
```

## 5. gRPC 服务定义

### 服务定义格式

```protobuf
service Query {
  // 简单 RPC
  rpc Dispute(QueryDisputeRequest) returns (QueryDisputeResponse);

  // 服务端流式 (Cosmos 较少使用)
  // rpc StreamEvents(StreamRequest) returns (stream Event);
}
```

### 在 Go 中实现服务

```go
// keeper/grpc_query.go
func (k Keeper) Dispute(c context.Context, req *types.QueryDisputeRequest) (*types.QueryDisputeResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "empty request")
    }

    dispute, found := k.GetDispute(sdk.UnwrapSDKContext(c), req.DisputeId)
    if !found {
        return nil, status.Error(codes.NotFound, "dispute not found")
    }

    return &types.QueryDisputeResponse{Dispute: &dispute}, nil
}
```

## 6. Cosmos 项目中 Proto 开发工作流

### 添加新的 Proto 消息

1. **定义消息** (在对应的 .proto 文件中):
```protobuf
// proto/sharetokens/trust/v1/trust.proto
message NewFeature {
  string name = 1;
  repeated string options = 2;
}
```

2. **添加交易消息** (在 tx.proto 中):
```protobuf
service Msg {
  rpc DoNewFeature(MsgDoNewFeature) returns (MsgDoNewFeatureResponse);
}

message MsgDoNewFeature {
  option (cosmos.msg.v1.signer) = "creator";
  string creator = 1;
  NewFeature feature = 2;
}
```

3. **生成代码**:
```bash
buf generate
```

4. **实现 Keeper 方法**:
```go
// x/trust/keeper/msg_server.go
func (k msgServer) DoNewFeature(goCtx context.Context, msg *types.MsgDoNewFeature) (*types.MsgDoNewFeatureResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)
    // 实现业务逻辑
    return &types.MsgDoNewFeatureResponse{}, nil
}
```

### 字段编号规则

- **1-15**: 常用字段，单字节编码
- **16-2047**: 较少使用字段，两字节编码
- **不要复用编号**: 删除字段时保留编号，标记为 `reserved`
- **不要修改类型**: 可兼容的类型转换除外

```protobuf
message Example {
  reserved 2, 3;           // 保留已删除的字段编号
  reserved "old_field";    // 保留已删除的字段名

  string new_field = 4;    // 新字段使用新编号
}
```

### 导入路径

```protobuf
// 导入 Google 标准类型
import "google/protobuf/timestamp.proto";
import "google/protobuf/any.proto";

// 导入 Cosmos SDK 类型
import "cosmos/base/query/v1beta1/pagination.proto";
import "cosmos/msg/v1/msg.proto";

// 导入项目内部类型 (相对路径)
import "proto/sharetokens/trust/v1/trust.proto";
import "proto/sharetokens/trust/v1/dispute.proto";
```

### 常见问题

1. **循环依赖**: 避免两个 proto 文件互相导入
2. **字段重命名**: 使用 `json_name` 保持 JSON 兼容性
3. **大数字**: 金额类字段使用 `string` 而非 `uint64` 避免精度问题

```protobuf
message Payment {
  // 金额用 string 表示，避免大数问题
  string amount = 1;
  string denom = 2;
}
```

## 参考资料

- [Protocol Buffers 官方文档](https://protobuf.dev/)
- [Cosmos SDK Proto 规范](https://docs.cosmos.network/main/build/building-modules/protobuf-annotations)
- [Buf 工具文档](https://buf.build/docs/)
