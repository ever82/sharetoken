# Quick Task Plan: Trust System Module Foundation

**Created:** 2026-03-03
**Mode:** quick
**Target Context:** ~30%

---

## Goal

Implement the foundational structure for Trust System Module (x/trust) containing:
1. Proto type definitions for MQ scoring and Dispute arbitration
2. Core Go types and Keeper structure
3. Basic MQ scoring implementation (zero-sum redistribution)

This is the **FIRST implementation task** - focus on foundation only.

---

## Context

From `.planning/models/09-dispute.md`:

**Trust System Architecture:**
- MQ (Moral Quotient): Zero-sum reputation scoring with convergence properties
- Dispute Arbitration: AI-mediated resolution with jury voting fallback

**Key Design Principles:**
- Zero-Sum: Total MQ is constant, loser's loss = winner's gain
- Convergence: Higher MQ = harder to increase (logarithmic rewards)
- Duty of Justice: Jury participation mandatory, absence = penalty
- Controlled Risk: Maximum 3% MQ loss per transaction
- Wisdom of Crowds: Random jury selection, multi-proposal scoring

**Module Dependencies:** auth, identity, escrow (Cosmos SDK modules)

---

## Tasks

### Task 1: Proto Definitions

**Files:**
- `proto/sharetokens/trust/v1/trust.proto` - Core MQ types
- `proto/sharetokens/trust/v1/dispute.proto` - Dispute types
- `proto/sharetokens/trust/v1/tx.proto` - Transaction messages
- `proto/sharetokens/trust/v1/query.proto` - Query services
- `proto/sharetokens/trust/v1/genesis.proto` - Genesis state

**Action:**
Create protobuf definitions following the design in `09-dispute.md`:

1. **trust.proto** - Define core types:
   ```protobuf
   // MQ Configuration
   message MQConfig {
     uint64 initial_mq = 1;           // 100
     string risk_rate = 2;            // 0.03 (3%)
     string lambda = 3;               // 1.5 (baseline deviation)
     string max_deviation = 4;        // 6.0
     JurySizeConfig jury_size = 5;
   }

   message JurySizeConfig {
     uint64 small = 1;   // 5
     uint64 medium = 2;  // 11
     uint64 large = 3;   // 21
     uint64 huge = 4;    // 31
   }

   // MQ Record
   message MQRecord {
     string address = 1;
     uint64 mq = 2;
     repeated MQHistoryEntry history = 3;
     MQStats stats = 4;
     google.protobuf.Timestamp created_at = 5;
     google.protobuf.Timestamp updated_at = 6;
   }

   message MQHistoryEntry {
     int64 change = 1;
     string reason = 2;
     uint64 related_id = 3;
     google.protobuf.Timestamp timestamp = 4;
   }

   message MQStats {
     uint64 total_disputes = 1;
     uint64 disputes_won = 2;
     uint64 disputes_lost = 3;
     uint64 jury_duties = 4;
     uint64 correct_votes = 5;
     uint64 wrong_votes = 6;
     int64 total_gained = 7;
     int64 total_lost = 8;
   }
   ```

2. **dispute.proto** - Define dispute lifecycle types:
   ```protobuf
   enum DisputeStatus {
     DISPUTE_STATUS_FILED = 0;
     DISPUTE_STATUS_MEDIATING = 1;
     DISPUTE_STATUS_SETTLED = 2;
     DISPUTE_STATUS_JURIED = 3;
     DISPUTE_STATUS_FORFEIT = 4;
     DISPUTE_STATUS_RESOLVED = 5;
     DISPUTE_STATUS_APPEALED = 6;
     DISPUTE_STATUS_FINAL = 7;
   }

   enum DisputeType {
     DISPUTE_TYPE_SERVICE_NOT_DELIVERED = 0;
     DISPUTE_TYPE_PAYMENT_NOT_RECEIVED = 1;
     DISPUTE_TYPE_QUALITY_ISSUE = 2;
     DISPUTE_TYPE_FRAUD = 3;
     DISPUTE_TYPE_OTHER = 4;
   }

   message Dispute {
     uint64 id = 1;
     uint64 order_id = 2;
     DisputeParty plaintiff = 3;
     DisputeParty defendant = 4;
     DisputeType type = 5;
     string title = 6;
     string description = 7;
     MediationSession mediation = 8;
     JuryInfo jury = 9;
     Resolution resolution = 10;
     DisputeStatus status = 11;
     uint64 appeal_count = 12;
     google.protobuf.Timestamp created_at = 13;
   }

   message DisputeParty {
     string address = 1;
     uint64 mq_stake = 2;
   }

   message MediationSession {
     uint64 dispute_id = 1;
     repeated MediationEvent events = 2;
     string status = 3;  // "active" | "settled" | "juried" | "forfeit"
     MediationOutcome outcome = 4;
     google.protobuf.Timestamp started_at = 5;
   }

   message MediationEvent {
     uint64 id = 1;
     google.protobuf.Timestamp timestamp = 2;
     string actor = 3;  // "A" | "B" | "AI"
     string type = 4;   // "message" | "evidence" | "proposal" | "score" | "verdict"
     bytes content = 5; // Serialized content based on type
   }
   ```

3. **tx.proto** - Transaction messages:
   ```protobuf
   service Msg {
     rpc InitializeMQ(MsgInitializeMQ) returns (MsgInitializeMQResponse);
     rpc CreateDispute(MsgCreateDispute) returns (MsgCreateDisputeResponse);
     rpc SendMessage(MsgSendMessage) returns (MsgSendMessageResponse);
     rpc SubmitEvidence(MsgSubmitEvidence) returns (MsgSubmitEvidenceResponse);
     rpc SubmitProposal(MsgSubmitProposal) returns (MsgSubmitProposalResponse);
     rpc SubmitScore(MsgSubmitScore) returns (MsgSubmitScoreResponse);
     rpc CastJuryVote(MsgCastJuryVote) returns (MsgCastJuryVoteResponse);
     rpc ResolveDispute(MsgResolveDispute) returns (MsgResolveDisputeResponse);
   }
   ```

4. **query.proto** - Query services:
   ```protobuf
   service Query {
     // MQ queries
     rpc MQ(QueryMQRequest) returns (QueryMQResponse);
     rpc MQHistory(QueryMQHistoryRequest) returns (QueryMQHistoryResponse);
     rpc MQLeaderboard(QueryMQLeaderboardRequest) returns (QueryMQLeaderboardResponse);

     // Dispute queries
     rpc Dispute(QueryDisputeRequest) returns (QueryDisputeResponse);
     rpc Disputes(QueryDisputesRequest) returns (QueryDisputesResponse);
     rpc MediationEvents(QueryMediationEventsRequest) returns (QueryMediationEventsResponse);
   }
   ```

**Verify:**
```bash
# Proto files exist
ls proto/sharetokens/trust/v1/*.proto

# Proto compiles successfully (requires protoc)
protoc --go_out=. --go-grpc_out=. proto/sharetokens/trust/v1/*.proto
```

**Done:**
- All 5 proto files created with complete type definitions
- Types match design document specifications
- Proto compilation succeeds without errors

---

### Task 2: Core Types and Keeper Structure

**Files:**
- `x/trust/types/keys.go` - Store keys
- `x/trust/types/errors.go` - Error definitions
- `x/trust/types/trust.go` - Go type definitions
- `x/trust/types/genesis.go` - Genesis state
- `x/trust/keeper/keeper.go` - Main keeper
- `x/trust/module.go` - Module definition

**Action:**
Implement the Go types and keeper structure:

1. **types/keys.go** - Define store keys:
   ```go
   package types

   const (
       ModuleName = "trust"
       StoreKey   = "trust"
       RouterKey  = "trust"

       // Store prefixes
       MQKeyPrefix       = "MQ/value/"
       MQHistoryPrefix   = "MQ/history/"
       DisputeKeyPrefix  = "Dispute/value/"
       EvidenceKeyPrefix = "Evidence/value/"
   )

   func MQKey(address sdk.AccAddress) []byte {
       return append([]byte(MQKeyPrefix), address.Bytes()...)
   }

   func DisputeKey(id uint64) []byte {
       return append([]byte(DisputeKeyPrefix), sdk.Uint64ToBigEndian(id)...)
   }
   ```

2. **types/errors.go** - Define errors:
   ```go
   package types

   var (
       ErrMQAlreadyInitialized = sdkerrors.Register(ModuleName, 1, "MQ already initialized")
       ErrMQNotFound           = sdkerrors.Register(ModuleName, 2, "MQ not found")
       ErrInvalidMQAmount      = sdkerrors.Register(ModuleName, 3, "invalid MQ amount")
       ErrDisputeNotFound      = sdkerrors.Register(ModuleName, 4, "dispute not found")
       ErrInvalidDisputeStatus = sdkerrors.Register(ModuleName, 5, "invalid dispute status")
       ErrNotAuthorized        = sdkerrors.Register(ModuleName, 6, "not authorized")
       ErrJuryAlreadyVoted     = sdkerrors.Register(ModuleName, 7, "juror already voted")
   )
   ```

3. **types/trust.go** - Go types matching proto:
   ```go
   package types

   type MQRecord struct {
       Address   sdk.AccAddress
       MQ        uint64
       History   []MQHistoryEntry
       Stats     MQStats
       CreatedAt time.Time
       UpdatedAt time.Time
   }

   type MQConfig struct {
       InitialMQ     uint64
       RiskRate      sdk.Dec  // 0.03
       Lambda        sdk.Dec  // 1.5
       MaxDeviation  sdk.Dec  // 6.0
       ChangeFactor  sdk.Dec  // 0.01
       JurySize      JurySizeConfig
   }

   type Dispute struct {
       Id          uint64
       OrderId     uint64
       Plaintiff   DisputeParty
       Defendant   DisputeParty
       Type        DisputeType
       Title       string
       Description string
       Mediation   MediationSession
       Jury        *JuryInfo
       Resolution  *Resolution
       Status      DisputeStatus
       AppealCount uint64
       CreatedAt   time.Time
   }
   ```

4. **keeper/keeper.go** - Main keeper structure:
   ```go
   package keeper

   type Keeper struct {
       storeKey     sdk.StoreKey
       cdc          codec.BinaryCodec
       bankKeeper   bankkeeper.Keeper
       escrowKeeper escrowkeeper.Keeper

       // Configuration
       config types.MQConfig
   }

   func NewKeeper(
       cdc codec.BinaryCodec,
       storeKey sdk.StoreKey,
       bankKeeper bankkeeper.Keeper,
       escrowKeeper escrowkeeper.Keeper,
   ) Keeper {
       return Keeper{
           storeKey:     storeKey,
           cdc:          cdc,
           bankKeeper:   bankKeeper,
           escrowKeeper: escrowKeeper,
           config:       DefaultMQConfig(),
       }
   }

   // MQ Operations (stub implementations)
   func (k Keeper) GetMQ(ctx sdk.Context, address sdk.AccAddress) uint64
   func (k Keeper) SetMQ(ctx sdk.Context, address sdk.AccAddress, mq uint64)
   func (k Keeper) HasMQ(ctx sdk.Context, address sdk.AccAddress) bool
   func (k Keeper) InitializeMQ(ctx sdk.Context, address sdk.AccAddress) error

   // Dispute Operations (stub implementations)
   func (k Keeper) CreateDispute(ctx sdk.Context, msg types.MsgCreateDispute) (uint64, error)
   func (k Keeper) GetDispute(ctx sdk.Context, id uint64) (types.Dispute, error)
   func (k Keeper) SetDispute(ctx sdk.Context, dispute types.Dispute)
   ```

**Verify:**
```bash
# Go files compile successfully
cd x/trust && go build ./...

# Unit tests pass (if any)
go test ./x/trust/types/... -v
```

**Done:**
- Store keys defined with proper prefixes
- Error codes registered (1-10 range)
- Go types match proto definitions
- Keeper structure with dependency injection
- Basic CRUD operations stubbed

---

### Task 3: Basic MQ Scoring Implementation

**Files:**
- `x/trust/keeper/mq.go` - MQ operations
- `x/trust/keeper/redistribution.go` - Zero-sum redistribution
- `x/trust/keeper/jury.go` - Jury selection

**Action:**
Implement the core MQ scoring logic:

1. **keeper/mq.go** - MQ initialization and CRUD:
   ```go
   package keeper

   // DefaultMQConfig returns default configuration
   func DefaultMQConfig() types.MQConfig {
       return types.MQConfig{
           InitialMQ:    100,
           RiskRate:     sdk.MustNewDecFromStr("0.03"),
           Lambda:       sdk.MustNewDecFromStr("1.5"),
           MaxDeviation: sdk.MustNewDecFromStr("6.0"),
           ChangeFactor: sdk.MustNewDecFromStr("0.01"),
           JurySize: types.JurySizeConfig{
               Small:  5,
               Medium: 11,
               Large:  21,
               Huge:   31,
           },
       }
   }

   // InitializeMQ creates initial MQ record for new user
   func (k Keeper) InitializeMQ(ctx sdk.Context, address sdk.AccAddress) error {
       if k.HasMQ(ctx, address) {
           return types.ErrMQAlreadyInitialized
       }

       record := types.MQRecord{
           Address:   address,
           MQ:        k.config.InitialMQ, // 100
           History: []types.MQHistoryEntry{{
               Change:    int64(k.config.InitialMQ),
               Reason:    "initial_verification",
               Timestamp: ctx.BlockTime(),
           }},
           CreatedAt: ctx.BlockTime(),
           UpdatedAt: ctx.BlockTime(),
       }

       k.SetMQRecord(ctx, record)
       return nil
   }

   // GetMQ retrieves current MQ for address
   func (k Keeper) GetMQ(ctx sdk.Context, address sdk.AccAddress) uint64 {
       record, found := k.GetMQRecord(ctx, address)
       if !found {
           return 0
       }
       return record.MQ
   }

   // SetMQ updates MQ value with history tracking
   func (k Keeper) SetMQ(ctx sdk.Context, address sdk.AccAddress, mq uint64, reason string, relatedId uint64) {
       record, found := k.GetMQRecord(ctx, address)
       if !found {
           return
       }

       change := int64(mq) - int64(record.MQ)

       record.History = append(record.History, types.MQHistoryEntry{
           Change:    change,
           Reason:    reason,
           RelatedId: relatedId,
           Timestamp: ctx.BlockTime(),
       })
       record.MQ = mq
       record.UpdatedAt = ctx.BlockTime()

       k.SetMQRecord(ctx, record)
   }
   ```

2. **keeper/redistribution.go** - Zero-sum redistribution algorithm:
   ```go
   package keeper

   // ScoringDeviation represents deviation from consensus
   type ScoringDeviation struct {
       Address      sdk.AccAddress
       MQ           uint64
       Deviations   []float64
       StdDeviation float64
       IsPenalized  bool  // d > lambda
       IsRewarded   bool  // d < lambda
   }

   // CalculateScoringDeviations calculates each juror's deviation from consensus
   func (k Keeper) CalculateScoringDeviations(
       ctx sdk.Context,
       proposals []types.ProposalScore,
       jurors map[string]uint64,
   ) []ScoringDeviation {
       results := make([]ScoringDeviation, 0)
       lambda := k.config.Lambda.MustFloat64()

       // Calculate MQ-weighted consensus for each proposal
       consensusScores := make([]float64, len(proposals))
       for i, proposal := range proposals {
           var weightedSum, totalWeight float64
           for addr, score := range proposal.Scores {
               mq := float64(jurors[addr])
               weightedSum += float64(score) * mq
               totalWeight += mq
           }
           if totalWeight > 0 {
               consensusScores[i] = weightedSum / totalWeight
           }
       }

       // Calculate deviation for each juror
       for addr, mq := range jurors {
           deviation := ScoringDeviation{
               Address:    sdk.AccAddress(addr),
               MQ:         mq,
               Deviations: make([]float64, len(proposals)),
           }

           var sumSquaredDiff float64
           for i, proposal := range proposals {
               score := float64(proposal.Scores[addr])
               diff := score - consensusScores[i]
               deviation.Deviations[i] = diff
               sumSquaredDiff += diff * diff
           }

           deviation.StdDeviation = math.Sqrt(sumSquaredDiff / float64(len(proposals)))
           deviation.IsPenalized = deviation.StdDeviation > lambda
           deviation.IsRewarded = deviation.StdDeviation < lambda

           results = append(results, deviation)
       }

       return results
   }

   // RedistributeMQ performs zero-sum redistribution
   func (k Keeper) RedistributeMQ(
       ctx sdk.Context,
       deviations []ScoringDeviation,
   ) map[string]int64 {
       changes := make(map[string]int64)
       lambda := k.config.Lambda.MustFloat64()
       maxD := k.config.MaxDeviation.MustFloat64()
       riskRate := k.config.RiskRate.MustFloat64()

       // Step 1: Calculate penalty pool
       var penaltyPool float64
       for _, d := range deviations {
           if d.IsPenalized {
               // Penalty rate = 3% * (d - lambda) / (max_d - lambda)
               penaltyRate := riskRate * (d.StdDeviation - lambda) / (maxD - lambda)
               penaltyRate = math.Min(penaltyRate, riskRate) // Cap at 3%

               loss := float64(d.MQ) * penaltyRate
               penaltyPool += loss
               changes[d.Address.String()] = -int64(math.Round(loss))
           }
       }

       // Step 2: Calculate reward distribution
       type rewardCandidate struct {
           addr    string
           mq      uint64
           contrib float64
       }

       var candidates []rewardCandidate
       var totalContrib float64

       for _, d := range deviations {
           if d.IsRewarded {
               // Contribution = (lambda - d) * log(MQ + 1) [convergence]
               contrib := (lambda - d.StdDeviation) * math.Log(float64(d.MQ)+1)
               candidates = append(candidates, rewardCandidate{
                   addr:    d.Address.String(),
                   mq:      d.MQ,
                   contrib: contrib,
               })
               totalContrib += contrib
           }
       }

       // Step 3: Distribute reward pool
       if totalContrib > 0 && penaltyPool > 0 {
           allocFactor := penaltyPool / totalContrib
           for _, c := range candidates {
               gain := c.contrib * allocFactor
               changes[c.addr] = changes[c.addr] + int64(math.Round(gain))
           }
       }

       // Step 4: Apply changes (ensure non-negative)
       for addrStr, change := range changes {
           addr, _ := sdk.AccAddressFromBech32(addrStr)
           currentMQ := k.GetMQ(ctx, addr)
           newMQ := int64(currentMQ) + change
           if newMQ < 0 {
               newMQ = 0 // MQ never negative
           }
           k.SetMQ(ctx, addr, uint64(newMQ), "redistribution", 0)
       }

       return changes
   }
   ```

3. **keeper/jury.go** - Jury selection based on MQ:
   ```go
   package keeper

   // DetermineJurySize determines jury size based on dispute amount
   func (k Keeper) DetermineJurySize(amount sdk.Coins) uint64 {
       // Simplified: use medium size for now
       // TODO: Implement proper threshold checking
       return k.config.JurySize.Medium // 11
   }

   // SelectJury selects jurors using MQ-weighted random selection
   func (k Keeper) SelectJury(
       ctx sdk.Context,
       size uint64,
       exclude map[string]bool,
   ) []sdk.AccAddress {
       // Get all eligible jurors (MQ >= 50)
       var candidates []types.JurorCandidate
       // TODO: Implement juror iteration

       var totalMQ uint64
       for _, c := range candidates {
           totalMQ += c.MQ
       }

       selected := make([]sdk.AccAddress, 0, size)
       selectedSet := make(map[string]bool)

       for uint64(len(selected)) < size {
           // Weighted random selection
           random := rand.Uint64() % totalMQ
           var cumulative uint64

           for _, candidate := range candidates {
               if selectedSet[candidate.Address.String()] {
                   continue
               }

               cumulative += candidate.MQ
               if random < cumulative {
                   selected = append(selected, candidate.Address)
                   selectedSet[candidate.Address.String()] = true
                   break
               }
           }
       }

       return selected
   }

   // GetVoteWeight returns voting weight (direct MQ value)
   func (k Keeper) GetVoteWeight(ctx sdk.Context, juror sdk.AccAddress) sdk.Dec {
       mq := k.GetMQ(ctx, juror)
       return sdk.NewDec(int64(mq))
   }
   ```

**Verify:**
```bash
# Unit tests for redistribution algorithm
go test ./x/trust/keeper/... -run TestRedistributeMQ -v

# Verify zero-sum property
# Total MQ before == Total MQ after redistribution
```

**Done:**
- MQ initialization with default 100 MQ
- MQ CRUD with history tracking
- Zero-sum redistribution algorithm implemented
- Convergence via logarithmic rewards
- Jury selection with MQ-weighted random
- Voting weight = MQ (no multipliers)
- MQ never goes negative (floor at 0)

---

## Dependencies

**External:**
- Cosmos SDK (auth, bank modules)
- Protobuf compiler (protoc)
- Go 1.19+

**Internal (will be mocked):**
- x/identity - Identity verification
- x/escrow - Escrow keeper

---

## Notes

1. **This is foundation only** - Full dispute mediation flow comes in follow-up tasks
2. **Use MockKeepers** for testing until x/identity and x/escrow are implemented
3. **Proto-first approach** - All types defined in proto, Go types generated
4. **Zero-sum verification** - Every redistribution should preserve total MQ
5. **Convergence test** - High MQ users should gain less proportionally

---

## Success Criteria

- [ ] All proto files compile successfully
- [ ] Go types compile and match proto definitions
- [ ] Keeper structure with dependency injection
- [ ] MQ initialization creates record with 100 MQ
- [ ] Zero-sum redistribution algorithm works correctly
- [ ] Jury selection uses MQ-weighted random
- [ ] Unit tests verify core algorithm properties

---

## Next Steps

After this plan:
1. Implement full dispute lifecycle (MsgServer)
2. Add AI mediation event handling
3. Implement evidence storage (hash only)
4. Add appeal mechanism
5. Wire up to x/escrow for dispute locking
