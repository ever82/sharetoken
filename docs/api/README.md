# ShareToken API Documentation

This document provides a comprehensive reference for all Query and Msg (transaction) endpoints available in the ShareToken blockchain.

## Table of Contents

- [Overview](#overview)
- [Modules](#modules)
  - [sharetoken](#sharetoken-module)
  - [identity](#identity-module)
  - [taskmarket](#taskmarket-module)
  - [escrow](#escrow-module)
  - [dispute](#dispute-module)
  - [trust](#trust-module)
  - [marketplace](#marketplace-module)
  - [crowdfunding](#crowdfunding-module)
  - [oracle](#oracle-module)
  - [llmcustody](#llmcustody-module)

---

## Overview

ShareToken uses Cosmos SDK with gRPC for query operations and transaction messages for state changes. All endpoints are organized by module.

### Common Patterns

- **Query endpoints**: Read-only operations that fetch state from the blockchain
- **Msg endpoints**: Transaction operations that modify state (require signing)
- **REST API**: Available at `/cosmos/{module}/{version}/...`
- **gRPC**: Native gRPC endpoints for efficient querying

---

## Modules

### sharetoken Module

The core token module managing ShareToken parameters and basic operations.

#### Query Endpoints

| Method | gRPC Path | Description |
|--------|-----------|-------------|
| `Params` | `/sharetoken.sharetoken.Query/Params` | Query module parameters |

**Example Request:**
```json
{
  "query": "Params"
}
```

**Example Response:**
```json
{
  "params": {
    // Module-specific parameters
  }
}
```

#### Msg Types

This module currently has no transaction messages defined.

---

### identity Module

Manages user identity registration, verification, and limit configurations.

#### Query Endpoints

| Method | gRPC Path | HTTP Endpoint | Description |
|--------|-----------|---------------|-------------|
| `Params` | `/sharetoken.identity.Query/Params` | `/ShareToken/sharetoken/identity/params` | Query module parameters |
| `Identity` | `/sharetoken.identity.Query/Identity` | `/ShareToken/sharetoken/identity/{address}` | Query identity by address |
| `Identities` | `/sharetoken.identity.Query/Identities` | `/ShareToken/sharetoken/identity/identities` | Query all identities (paginated) |
| `LimitConfig` | `/sharetoken.identity.Query/LimitConfig` | `/ShareToken/sharetoken/identity/{address}/limits` | Query user limit configuration |
| `IsVerified` | `/sharetoken.identity.Query/IsVerified` | `/ShareToken/sharetoken/identity/{address}/verified` | Check if address is verified |

**QueryIdentity Request:**
```json
{
  "address": "sharetoken1..."
}
```

**QueryIdentity Response:**
```json
{
  "identity": {
    "address": "sharetoken1...",
    "did": "did:sharetoken:abc123",
    "registration_time": 1700000000,
    "is_verified": true,
    "verification_provider": "wechat",
    "merkle_root": "a1b2c3...",
    "metadata_hash": "d4e5f6..."
  }
}
```

**QueryLimitConfig Response:**
```json
{
  "limit_config": {
    "address": "sharetoken1...",
    "tx_limit": {
      "max_single": "1000000000ustt",
      "max_daily": "10000000000ustt",
      "max_monthly": "100000000000ustt",
      "daily_spent": "0ustt",
      "monthly_spent": "0ustt"
    },
    "withdrawal_limit": {
      "max_daily": "5000000000ustt",
      "cooldown_hours": 24
    },
    "dispute_limit": {
      "max_active_disputes": 5,
      "current_active": 0
    },
    "service_limit": {
      "max_concurrent": 10,
      "rate_limit_per_minute": 60
    }
  }
}
```

#### Msg Types

| Message | Route | Description |
|---------|-------|-------------|
| `MsgRegisterIdentity` | `identity` | Register a new identity |
| `MsgVerifyIdentity` | `identity` | Verify identity with third-party provider |
| `MsgUpdateLimitConfig` | `identity` | Update user's limit configuration (governance) |
| `MsgResetDailyLimits` | `identity` | Reset daily limits (end blocker) |
| `MsgUpdateParams` | `identity` | Update module parameters (governance) |

**MsgRegisterIdentity:**
```json
{
  "type": "identity/RegisterIdentity",
  "value": {
    "address": "sharetoken1...",
    "did": "did:sharetoken:abc123",
    "metadata_hash": "hash_of_encrypted_metadata"
  }
}
```

**Response:**
```json
{
  "merkle_root": "generated_merkle_root_hash"
}
```

**MsgVerifyIdentity:**
```json
{
  "type": "identity/VerifyIdentity",
  "value": {
    "address": "sharetoken1...",
    "provider": "wechat",  // or "github", "google"
    "verification_hash": "hash_of_verification_result",
    "proof": "verification_proof_from_provider"
  }
}
```

**Response:**
```json
{
  "is_verified": true,
  "updated_merkle_root": "updated_merkle_root_hash"
}
```

---

### taskmarket Module

Decentralized task marketplace for posting, bidding, and completing tasks.

#### Query Types

| Query | Description | Request | Response |
|-------|-------------|---------|----------|
| `QueryTasksRequest` | Query tasks with filters | `status`, `category`, `requester_id`, `worker_id`, `limit`, `offset` | `QueryTasksResponse` |
| `QueryTaskRequest` | Query single task | `task_id` | `QueryTaskResponse` |
| `QueryApplicationsRequest` | Query applications | `task_id`, `worker_id`, `limit`, `offset` | `QueryApplicationsResponse` |
| `QueryBidsRequest` | Query bids | `task_id`, `worker_id` | `QueryBidsResponse` |
| `QueryAuctionRequest` | Query auction | `task_id` | `QueryAuctionResponse` |
| `QueryReputationRequest` | Query user reputation | `user_id` | `QueryReputationResponse` |
| `QueryRatingsRequest` | Query ratings | `user_id`, `task_id`, `limit`, `offset` | `QueryRatingsResponse` |
| `QueryStatisticsRequest` | Query market statistics | - | `QueryStatisticsResponse` |

**QueryTasksRequest:**
```json
{
  "status": "open",           // optional: draft, open, assigned, in_progress, review, completed, cancelled
  "category": "development",  // optional: development, design, writing, research, marketing, consulting, other
  "requester_id": "sharetoken1...", // optional
  "worker_id": "sharetoken1...",    // optional
  "limit": 10,
  "offset": 0
}
```

**QueryTasksResponse:**
```json
{
  "tasks": [
    {
      "id": "task-001",
      "title": "Build a DeFi Dashboard",
      "description": "Create a React dashboard...",
      "requester_id": "sharetoken1abc...",
      "worker_id": "",
      "type": "open",        // "open" or "auction"
      "category": "development",
      "status": "open",
      "budget": 500000000,   // in ustt (micro-STT)
      "currency": "STT",
      "deadline": 1704067200,
      "skills": ["react", "typescript", "cosmos-sdk"],
      "subtasks": [],
      "milestones": [],
      "created_at": 1700000000,
      "updated_at": 1700000000,
      "view_count": 25,
      "application_count": 3,
      "bid_count": 0
    }
  ],
  "total_count": 100
}
```

**QueryReputationResponse:**
```json
{
  "reputation": {
    "user_id": "sharetoken1...",
    "total_ratings": 15,
    "average_rating": 4.5,
    "ratings_by_dimension": {
      "quality": 4.6,
      "communication": 4.4,
      "timeliness": 4.5,
      "professionalism": 4.5
    },
    "completed_tasks": 12,
    "dispute_rate": 0.0,
    "on_time_delivery": 95.0,
    "created_at": 1700000000,
    "updated_at": 1700000000
  }
}
```

#### Msg Types

##### Task Management

| Message | Type | Description |
|---------|------|-------------|
| `MsgCreateTask` | `create_task` | Create a new task |
| `MsgUpdateTask` | `update_task` | Update an existing task |
| `MsgPublishTask` | `publish_task` | Publish a draft task |
| `MsgCancelTask` | `cancel_task` | Cancel a task |
| `MsgStartTask` | `start_task` | Start working on a task |

**MsgCreateTask:**
```json
{
  "type": "taskmarket/CreateTask",
  "value": {
    "creator": "sharetoken1...",
    "title": "Build a DeFi Dashboard",
    "description": "Create a React dashboard for DeFi protocols",
    "task_type": "open",       // "open" or "auction"
    "category": "development",
    "budget": 500000000,
    "currency": "STT",
    "deadline": 1704067200,
    "skills": ["react", "typescript"],
    "subtasks": [],
    "milestones": []
  }
}
```

**Response:**
```json
{
  "task_id": "task-001"
}
```

##### Applications (for Open Tasks)

| Message | Type | Description |
|---------|------|-------------|
| `MsgSubmitApplication` | `submit_application` | Submit application for open task |
| `MsgAcceptApplication` | `accept_application` | Accept a worker's application |
| `MsgRejectApplication` | `reject_application` | Reject a worker's application |

**MsgSubmitApplication:**
```json
{
  "type": "taskmarket/SubmitApplication",
  "value": {
    "worker_id": "sharetoken1...",
    "task_id": "task-001",
    "proposed_price": 450000000,
    "cover_letter": "I have 5 years experience...",
    "relevant_experience": ["Built 3 DeFi dashboards"],
    "portfolio_links": ["https://github.com/..."],
    "estimated_duration": 14  // days
  }
}
```

##### Bids (for Auction Tasks)

| Message | Type | Description |
|---------|------|-------------|
| `MsgSubmitBid` | `submit_bid` | Submit a bid for auction task |
| `MsgCloseAuction` | `close_auction` | Close auction and select winner |

**MsgSubmitBid:**
```json
{
  "type": "taskmarket/SubmitBid",
  "value": {
    "worker_id": "sharetoken1...",
    "task_id": "task-001",
    "amount": 400000000,  // Lower is better in auctions
    "message": "I can complete this in 10 days",
    "portfolio": "https://portfolio.example.com"
  }
}
```

##### Milestones

| Message | Type | Description |
|---------|------|-------------|
| `MsgSubmitMilestone` | `submit_milestone` | Submit milestone for review |
| `MsgApproveMilestone` | `approve_milestone` | Approve a submitted milestone |
| `MsgRejectMilestone` | `reject_milestone` | Reject a submitted milestone |

**MsgSubmitMilestone:**
```json
{
  "type": "taskmarket/SubmitMilestone",
  "value": {
    "worker_id": "sharetoken1...",
    "task_id": "task-001",
    "milestone_id": "milestone-001",
    "deliverables": "Completed frontend components..."
  }
}
```

##### Ratings

| Message | Type | Description |
|---------|------|-------------|
| `MsgSubmitRating` | `submit_rating` | Submit a rating for completed work |

**MsgSubmitRating:**
```json
{
  "type": "taskmarket/SubmitRating",
  "value": {
    "task_id": "task-001",
    "rater_id": "sharetoken1...",
    "rated_id": "sharetoken1...",
    "ratings": {
      "quality": 5,
      "communication": 4,
      "timeliness": 5,
      "professionalism": 5
    },
    "comment": "Excellent work, delivered on time!"
  }
}
```

---

### escrow Module

Manages escrow agreements for secure transactions between parties.

#### Types

**Escrow:**
```go
type Escrow struct {
    ID              string       // Escrow ID
    Requester       string       // Requester address
    Provider        string       // Provider address
    Amount          sdk.Coins    // Escrowed amount
    Status          EscrowStatus // pending, completed, disputed, refunded
    CreatedAt       int64
    ExpiresAt       int64
    CompletedAt     int64
    CompletionProof string
    DisputeID       string
    RefundAddress   string
}
```

**EscrowStatus values:**
- `pending` - Escrow is active and waiting for completion
- `completed` - Funds released to provider
- `disputed` - Under dispute resolution
- `refunded` - Refunded to requester

#### Query Types

Queries are handled through the keeper's direct query methods (no proto-generated queries yet).

#### Msg Types

The module handles escrow lifecycle through internal keeper methods. Direct Msg types are not yet exposed.

---

### dispute Module

Manages dispute resolution between parties.

#### Types

**Dispute:**
```go
type Dispute struct {
    ID          string        // Dispute ID
    EscrowID    string        // Associated escrow ID
    Requester   string        // Requester address
    Provider    string        // Provider address
    Status      DisputeStatus // open, mediating, voting, resolved, cancelled
    Reason      string        // Dispute reason
    Evidence    []Evidence    // Submitted evidence
    Votes       []Vote        // Votes on dispute
    Result      VoteResult    // Voting result
    CreatedAt   int64
    CompletedAt int64
}
```

**DisputeStatus values:**
- `open` - Dispute is open
- `mediating` - In mediation phase
- `voting` - Community voting phase
- `resolved` - Dispute resolved
- `cancelled` - Dispute cancelled

**Evidence:**
```go
type Evidence struct {
    SubmittedBy string // Address of submitter
    Type        string // Evidence type
    Content     string // Evidence content
    Timestamp   int64
}
```

**Vote:**
```go
type Vote struct {
    Voter     string  // Voter address
    Weight    sdk.Dec // Vote weight (based on MQ score)
    Decision  string  // "requester", "provider", or "split"
    Timestamp int64
}
```

---

### trust Module

Manages Moral Quotient (MQ) scores for reputation-based voting.

#### Types

**MQScore:**
```go
type MQScore struct {
    Address   string // User address
    Score     int32  // MQ score (0-100)
    Disputes  uint64 // Number of disputes participated
    Consensus uint64 // Times voted with consensus
    UpdatedAt int64
}
```

**Key Concepts:**
- Initial MQ score: 100
- Score changes based on dispute voting behavior
- Voting with consensus: +0.5% to +1% reward
- Voting against consensus: -1% to -3% penalty
- Higher MQ = higher voting weight in disputes

#### Query Types

| Query | Description |
|-------|-------------|
| `QueryMQScore` | Query user's MQ score |
| `QueryMQStats` | Query MQ statistics |

**QueryMQScore Response:**
```json
{
  "mq_score": {
    "address": "sharetoken1...",
    "score": 85,
    "disputes": 10,
    "consensus": 9,
    "updated_at": 1700000000
  }
}
```

---

### marketplace Module

Manages service listings for LLM and agent services.

#### Types

**Service:**
```go
type Service struct {
    ID          string       // Service ID
    Provider    string       // Provider address
    Name        string       // Service name
    Description string       // Service description
    Level       ServiceLevel // 1=LLM API, 2=Agent, 3=Workflow
    PricingMode PricingMode  // fixed, dynamic, auction
    Price       sdk.Coins    // Service price
    Active      bool         // Is service active
    CreatedAt   int64
}
```

**ServiceLevel values:**
- `1` - LLM API
- `2` - Agent
- `3` - Workflow

**PricingMode values:**
- `fixed` - Fixed price
- `dynamic` - Dynamic pricing
- `auction` - Auction-based pricing

---

### crowdfunding Module

Manages crowdfunding campaigns for ideas and projects.

#### Types

**Campaign:**
```go
type Campaign struct {
    ID          string         // Campaign ID
    IdeaID      string         // Associated idea ID
    CreatorID   string         // Creator address
    Title       string
    Description string
    Type        CampaignType   // investment, lending, donation
    Status      CampaignStatus // draft, active, funded, expired, cancelled, closed
    GoalAmount  uint64         // Target amount
    RaisedAmount uint64        // Currently raised
    Currency    string         // "STT", "USDC"
    StartTime   int64
    EndTime     int64
    MinContribution uint64
    MaxContribution uint64
    // Investment-specific
    EquityOffered float64
    Valuation     uint64
    // Lending-specific
    InterestRate      float64
    LoanTerm          int64
    RepaymentSchedule string
    // Statistics
    BackerCount int
    UpdateCount int
    CreatedAt   int64
    UpdatedAt   int64
}
```

**CampaignType values:**
- `investment` - Equity investment
- `lending` - Loan with interest
- `donation` - No return

**CampaignStatus values:**
- `draft` - Draft state
- `active` - Accepting contributions
- `funded` - Goal reached
- `expired` - Time ran out
- `cancelled` - Cancelled by creator
- `closed` - Funds distributed

**Idea:**
```go
type Idea struct {
    ID                string     // Idea ID
    Title             string
    Description       string
    CreatorID         string
    Status            IdeaStatus // draft, active, funding, developing, completed, archived
    CurrentVersion    int
    Tags              []string
    Categories        []string
    ViewCount         int
    ContributionCount int
    TotalWeight       float64
    CampaignID        string
    CreatedAt         int64
    UpdatedAt         int64
    PublishedAt       int64
}
```

**Contribution:**
```go
type Contribution struct {
    ID            string               // Contribution ID
    IdeaID        string
    ContributorID string
    Category      ContributionCategory // code, design, docs, research, marketing, testing
    Description   string
    Weight        float64              // Calculated weight
    RawScore      float64              // Base score
    Evidence      string               // Link to work
    Status        ContributionStatus   // pending, approved, rejected
    ReviewedBy    string
    ReviewedAt    int64
    CreatedAt     int64
}
```

---

### oracle Module

Provides price feed and oracle services.

#### Types

**Price:**
```go
type Price struct {
    Symbol     string      // Token symbol (e.g., "STT", "USDC")
    Price      sdk.Dec     // Price value
    Timestamp  int64       // Unix timestamp
    Source     PriceSource // chainlink, manual
    Confidence int32       // 0-100
}
```

**LLMPrice:**
```go
type LLMPrice struct {
    Provider    string  // e.g., "openai", "anthropic"
    Model       string  // e.g., "gpt-4", "claude-3"
    InputPrice  sdk.Dec // per 1K tokens
    OutputPrice sdk.Dec // per 1K tokens
    Currency    string  // e.g., "USD"
}
```

---

### llmcustody Module

Manages API key custody for LLM service providers.

#### Types

**AccessRule:**
```go
type AccessRule struct {
    Resource    string   // API resource path
    Methods     []string // Allowed HTTP methods
    RateLimit   int      // Requests per minute
    AllowedIPs  []string // Whitelisted IPs
    DeniedIPs   []string // Blacklisted IPs
}
```

#### Msg Types

| Message | Type | Description |
|---------|------|-------------|
| `MsgRegisterAPIKey` | `register_api_key` | Register a new API key |
| `MsgUpdateAPIKey` | `update_api_key` | Update API key settings |
| `MsgRevokeAPIKey` | `revoke_api_key` | Revoke an API key |
| `MsgRecordUsage` | `record_usage` | Record API usage |

**MsgRegisterAPIKey:**
```json
{
  "type": "llmcustody/RegisterAPIKey",
  "value": {
    "owner": "sharetoken1...",
    "provider": "openai",
    "encrypted_key": [/* encrypted key bytes */],
    "access_rules": [
      {
        "resource": "/v1/chat/completions",
        "methods": ["POST"],
        "rate_limit": 60,
        "allowed_ips": [],
        "denied_ips": []
      }
    ]
  }
}
```

**Response:**
```json
{
  "api_key_id": "key-abc123"
}
```

**MsgUpdateAPIKey:**
```json
{
  "type": "llmcustody/UpdateAPIKey",
  "value": {
    "owner": "sharetoken1...",
    "api_key_id": "key-abc123",
    "access_rules": [...],
    "active": true
  }
}
```

**MsgRecordUsage:**
```json
{
  "type": "llmcustody/RecordUsage",
  "value": {
    "api_key_id": "key-abc123",
    "service_id": "service-001",
    "request_count": 100,
    "token_count": 5000,
    "cost": 100000  // in ustt
  }
}
```

---

## Transaction Structure

### Standard Cosmos SDK Transaction

```json
{
  "body": {
    "messages": [
      {
        "@type": "/cosmos.bank.v1beta1.MsgSend",
        "from_address": "sharetoken1...",
        "to_address": "sharetoken1...",
        "amount": [
          {
            "denom": "ustt",
            "amount": "1000000"
          }
        ]
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [...],
    "fee": {
      "amount": [...],
      "gas_limit": "200000"
    }
  },
  "signatures": [...]
}
```

### Fee Structure

- **Gas price**: Configurable based on network conditions
- **Minimum gas**: Varies by message type complexity
- **Fee denomination**: `ustt` (micro-STT)

---

## Error Codes

Common error codes returned by the API:

| Code | Description |
|------|-------------|
| `1` | Internal error |
| `2` | Tx parse error |
| `3` | Invalid sequence |
| `4` | Unauthorized |
| `5` | Insufficient funds |
| `6` | Unknown request |
| `7` | Invalid address |
| `8` | Invalid coins |
| `9` | Unknown address |
| `10` | Invalid pubKey |
| `12` | Insufficient fee |
| `13` | Out of gas |
| `19` | Tx already in mempool |

Module-specific errors:

| Module | Error | Description |
|--------|-------|-------------|
| identity | InvalidAddress | Address format invalid |
| identity | InvalidDID | DID format invalid |
| identity | InvalidProvider | Provider not allowed |
| taskmarket | TaskNotFound | Task ID not found |
| taskmarket | InvalidTaskStatus | Task in wrong status |
| taskmarket | Unauthorized | Not authorized for action |
| escrow | InvalidEscrowID | Escrow ID invalid |
| escrow | InvalidAmount | Amount invalid |
| escrow | InvalidStatus | Status transition invalid |
| escrow | Unauthorized | Not authorized for action |

---

## Pagination

List queries support pagination using offset-based pagination:

```json
{
  "pagination": {
    "offset": 0,
    "limit": 10,
    "count_total": true
  }
}
```

**Response pagination:**
```json
{
  "pagination": {
    "next_key": "...",
    "total": 100
  }
}
```

---

## Rate Limits

The identity module enforces rate limits per user:

| Limit Type | Default Value |
|------------|---------------|
| Max single transaction | 1,000 STT |
| Max daily transactions | 10,000 STT |
| Max monthly transactions | 100,000 STT |
| Max daily withdrawal | 5,000 STT |
| Withdrawal cooldown | 24 hours |
| Max concurrent services | 10 |
| Rate limit per minute | 60 requests |

---

## Version Information

| Component | Version |
|-----------|---------|
| Cosmos SDK | v0.47.x |
| Tendermint | v0.37.x |
| Go | 1.21+ |

---

## Additional Resources

- [Cosmos SDK Documentation](https://docs.cosmos.network/)
- [Protobuf Definitions](./proto/) - See `/proto/sharetoken/` directory
- [CLI Documentation](./docs/cli.md)

---

*This documentation was generated for ShareToken. Last updated: 2026-03-13*
