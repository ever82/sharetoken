# Issue #11: ACH-DEV-011 LLM API Key Custody Plugin

> LLM Provider API Key 的安全托管与代理服务

---

## 验收标准

### 核心功能
1. ⏳ API Key 加密存储（链上存加密后hash）
2. ⏳ WASM 沙箱内解密使用
3. ⏳ 使用后立即清除（Secret Zeroization）
4. ⏳ 访问控制与定价配置
5. ⏳ 支持 OpenAI / Anthropic API 代理
6. ⏳ 密钥管理方案文档化

---

## 技术实现细节

### 模块结构
```
x/llmcustody/
├── keeper/
│   ├── keeper.go
│   └── keeper_test.go
├── types/
│   ├── apikey.go
│   ├── encryption.go
│   └── wasm.go
└── wasm/
    └── runtime.go
```

### 核心类型
```go
type APIKey struct {
    ID           string    // Key ID
    Provider     string    // openai/anthropic
    EncryptedKey []byte    // 加密后的密钥
    Hash         string    // 密钥hash
    Owner        string    // 所有者地址
    AccessRules  []Rule    // 访问控制规则
}
```

---

## 状态

- ⏳ 进行中

---

## 关联 Issue
#11
