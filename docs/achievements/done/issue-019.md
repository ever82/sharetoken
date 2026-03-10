# ACH-DEV-019: Node Role System

## Summary
Implemented a flexible node role system for ShareToken blockchain supporting four distinct node types with different capabilities and resource requirements.

## Node Roles

### 1. Light Node
- **Purpose**: Minimal resource usage for end users
- **Capabilities**:
  - Can run GenieBot plugin
  - Cannot validate (no consensus participation)
  - Cannot query full state (relies on trusted full nodes)
  - Cannot query history
- **Resource Requirements**: 10GB storage, 2GB memory
- **Pruning Strategy**: Aggressive pruning (keep minimal state)

### 2. Full Node
- **Purpose**: Standard validator node
- **Capabilities**:
  - Can validate (participates in consensus)
  - Can query full state
  - Can serve light clients
  - Cannot run plugins
- **Resource Requirements**: 100GB storage, 4GB memory
- **Pruning Strategy**: Custom (keep recent 100 blocks, every 10000th block)

### 3. Service Node
- **Purpose**: Execute AI services (LLM, Agent, Workflow)
- **Capabilities**:
  - Can run all service plugins
  - Cannot validate (focused on execution)
  - Uses WASM sandboxing for security
- **Resource Requirements**: 50GB storage, 8GB memory
- **Plugin Support**:
  - LLM API Key Custody
  - Agent Executor (OpenFang)
  - Workflow Executor (OpenFang Hands)

### 4. Archive Node
- **Purpose**: Complete historical data with indexing
- **Capabilities**:
  - Can validate
  - Can query full state and history
  - Can serve light clients
  - Can index blocks for explorer
  - Block explorer API support
- **Resource Requirements**: 1000GB storage, 16GB memory
- **Features**:
  - Transaction indexer
  - Event indexer
  - Block compression
  - Long-term retention

## Role Switching

### Hot Switch (No Restart Required)
- Light ↔ Service: Can switch without restart (just plugin load/unload)

### Cold Switch (Restart Required)
- Light → Full: Different pruning strategy
- Light → Archive: Different storage backend
- Full ↔ Archive: Different storage backend
- Any ↔ Archive: Requires full sync

## Implementation

### Types (`x/node/types/role.go`)
- `NodeRole` enum: Light, Full, Service, Archive
- `RoleConfig`: Role-specific configuration structs
- `NodeCapabilities`: Capability matrix for each role
- `CanSwitchTo()`: Determines switch feasibility

### Keeper (`x/node/keeper/keeper.go`)
- `NodeKeeper`: Manages current role and configuration
- `InitializeRole()`: Initialize node with a role
- `SwitchRoleHot()`: Perform hot role switch
- `RequestRoleSwitchForRestart()`: Schedule role change for next restart
- Plugin management integration

### CLI (`x/node/client/cli/tx.go`)
```bash
# Query current role
sharetokend query node role

# Query capabilities
sharetokend query node capabilities [role]

# Query node state
sharetokend query node state

# Initialize role
sharetokend tx node init-role light
sharetokend tx node init-role full
sharetokend tx node init-role service
sharetokend tx node init-role archive

# Switch role (detects if restart needed)
sharetokend tx node switch-role service --force

# Update configuration
sharetokend tx node update-config --config-file custom-config.json
```

## Configuration Files

### Light Node Config
```json
{
  "role": "light",
  "light_config": {
    "max_peers": 10,
    "pruning_strategy": "everything",
    "enable_genie_bot": true,
    "trusted_full_nodes": ["http://full-node-1:26657"]
  }
}
```

### Service Node Config
```json
{
  "role": "service",
  "service_config": {
    "max_peers": 20,
    "enable_llm_plugin": true,
    "enable_agent_plugin": true,
    "enable_workflow_plugin": true,
    "wasm_max_memory": 512,
    "wasm_max_cpu_time": 30000,
    "trusted_full_nodes": []
  }
}
```

### Archive Node Config
```json
{
  "role": "archive",
  "archive_config": {
    "max_peers": 50,
    "enable_block_explorer": true,
    "block_explorer_port": 8080,
    "index_all_transactions": true,
    "index_events": true,
    "compression_enabled": true,
    "archive_retention_years": 0
  }
}
```

## Capability Matrix

| Capability | Light | Full | Service | Archive |
|------------|-------|------|---------|---------|
| Validate | ❌ | ✅ | ❌ | ✅ |
| Query State | ❌ | ✅ | ❌ | ✅ |
| Query History | ❌ | ❌ | ❌ | ✅ |
| Serve Light Clients | ❌ | ✅ | ❌ | ✅ |
| Run Plugins | ✅ | ❌ | ✅ | ❌ |
| Index Blocks | ❌ | ❌ | ❌ | ✅ |
| Storage (GB) | 10 | 100 | 50 | 1000 |
| Memory (GB) | 2 | 4 | 8 | 16 |

## Testing

All tests pass:
- Role initialization
- Invalid role handling
- Double initialization prevention
- Capabilities validation
- Role switching logic
- Config persistence
- Start/Stop lifecycle

```bash
go test ./x/node/... -v
# PASS (11 tests)
```

## Files Created

```
x/node/
├── types/
│   └── role.go              # Role definitions and config
├── keeper/
│   ├── keeper.go            # Node keeper implementation
│   └── keeper_test.go       # Unit tests
└── client/cli/
    └── tx.go                # CLI commands
```

## Production Considerations

1. **Security**
   - Service nodes use WASM sandboxing for plugin isolation
   - Memory and CPU limits enforced per plugin
   - Trusted full nodes should use TLS

2. **Networking**
   - Light nodes connect to trusted full nodes for state queries
   - Service nodes focus on execution, minimal validation
   - Archive nodes handle high query load

3. **Upgrades**
   - Role changes requiring restart preserve old config
   - Config validation before applying changes
   - Graceful shutdown of plugins during role switch

## Next Steps

1. Integrate with actual plugin manager
2. Implement gRPC endpoints for role queries
3. Add metrics for role-specific operations
4. Create deployment guides for each role type
