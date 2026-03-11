# LLM API Key 密钥管理方案

> 本文档描述 ShareToken LLM Custody 模块的密钥管理架构和 KMS/HSM 集成方案

---

## 架构概述

```
┌─────────────────────────────────────────────────────────────────┐
│                         用户层                                   │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │  Provider    │  │   Service    │  │    User      │          │
│  │  Dashboard   │  │   Consumer   │  │   Wallet     │          │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘          │
└─────────┼─────────────────┼─────────────────┼──────────────────┘
          │                 │                 │
          ▼                 ▼                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                      API Gateway 层                              │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  Authentication │ Rate Limiting │ Request Routing       │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
          │
          ▼
┌─────────────────────────────────────────────────────────────────┐
│                     LLM Custody 模块                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │   Pricing   │  │   Access    │  │   Billing   │             │
│  │   Engine    │  │   Control   │  │   Service   │             │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘             │
│         │                │                │                     │
│         └────────────────┼────────────────┘                     │
│                          │                                      │
│                          ▼                                      │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │              Secure Key Management                       │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────────────┐      │   │
│  │  │ Key Hash │  │Encrypted │  │   WASM Sandbox   │      │   │
│  │  │ Storage  │  │   Key    │  │  (Decryption)    │      │   │
│  │  │ (Chain)  │  │ Storage  │  │                  │      │   │
│  │  └──────────┘  └──────────┘  └──────────────────┘      │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
          │
          ▼
┌─────────────────────────────────────────────────────────────────┐
│                      KMS/HSM 层                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │   AWS KMS    │  │  HashiCorp   │  │   YubiHSM    │          │
│  │              │  │    Vault     │  │              │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└─────────────────────────────────────────────────────────────────┘
          │
          ▼
┌─────────────────────────────────────────────────────────────────┐
│                     LLM Providers                                │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │    OpenAI    │  │  Anthropic   │  │   Others     │          │
│  │              │  │              │  │              │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└─────────────────────────────────────────────────────────────────┘
```

---

## 密钥层次结构

### 三层密钥架构

```
┌─────────────────────────────────────┐
│         Master Key (MK)             │  ← 由 KMS/HSM 管理
│    最高安全级别，永不离开 HSM        │
└─────────────┬───────────────────────┘
              │
              ▼ Encrypt/Decrypt
┌─────────────────────────────────────┐
│      Key Encryption Key (KEK)       │  ← 由 MK 加密存储
│     用于加密数据加密密钥 (DEK)       │
└─────────────┬───────────────────────┘
              │
              ▼ Encrypt/Decrypt
┌─────────────────────────────────────┐
│       Data Encryption Key (DEK)     │  ← 每个 API Key 一个
│      用于加密实际的 LLM API Key      │
└─────────────────────────────────────┘
```

### 密钥管理策略

| 密钥类型 | 存储位置 | 生命周期 | 轮换策略 |
|---------|---------|---------|---------|
| Master Key | HSM/KMS | 长期 | 每年或怀疑泄露时 |
| KEK | 加密后存储 | 中期 | 每季度 |
| DEK | 加密后存储 | 短期 | 每次 API Key 更新 |

---

## KMS/HSM 集成方案

### 方案一：AWS KMS (推荐用于 AWS 部署)

```go
package kms

import (
    "context"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/kms"
    "github.com/aws/aws-sdk-go-v2/service/kms/types"
)

// AWSKMSClient AWS KMS 客户端
type AWSKMSClient struct {
    client *kms.Client
    keyID  string
}

// NewAWSKMSClient 创建新的 AWS KMS 客户端
func NewAWSKMSClient(region, keyID string) (*AWSKMSClient, error) {
    cfg, err := config.LoadDefaultConfig(context.TODO(),
        config.WithRegion(region),
    )
    if err != nil {
        return nil, err
    }

    return &AWSKMSClient{
        client: kms.NewFromConfig(cfg),
        keyID:  keyID,
    }, nil
}

// GenerateDataKey 生成数据密钥
func (c *AWSKMSClient) GenerateDataKey() ([]byte, []byte, error) {
    result, err := c.client.GenerateDataKey(context.TODO(), &kms.GenerateDataKeyInput{
        KeyId:   aws.String(c.keyID),
        KeySpec: types.DataKeySpecAes256,
    })
    if err != nil {
        return nil, nil, err
    }

    // Plaintext: 明文 DEK（使用后立即擦除）
    // CiphertextBlob: 由 KMS 加密的 DEK（可以安全存储）
    return result.Plaintext, result.CiphertextBlob, nil
}

// Decrypt 解密数据密钥
func (c *AWSKMSClient) Decrypt(ciphertext []byte) ([]byte, error) {
    result, err := c.client.Decrypt(context.TODO(), &kms.DecryptInput{
        CiphertextBlob: ciphertext,
        KeyId:          aws.String(c.keyID),
    })
    if err != nil {
        return nil, err
    }
    return result.Plaintext, nil
}
```

**优点：**
- 完全托管，无需维护 HSM 硬件
- 与 AWS 生态系统集成良好
- 自动密钥轮换支持
- FIPS 140-2 Level 2 合规

**缺点：**
- 依赖 AWS 服务
- 有 API 调用费用

---

### 方案二：HashiCorp Vault (推荐用于自托管)

```go
package kms

import (
    "context"
    "github.com/hashicorp/vault/api"
)

// VaultClient HashiCorp Vault 客户端
type VaultClient struct {
    client *api.Client
    path   string
}

// NewVaultClient 创建新的 Vault 客户端
func NewVaultClient(address, token, path string) (*VaultClient, error) {
    config := &api.Config{
        Address: address,
    }

    client, err := api.NewClient(config)
    if err != nil {
        return nil, err
    }

    client.SetToken(token)

    return &VaultClient{
        client: client,
        path:   path,
    }, nil
}

// GenerateDataKey 使用 Vault Transit 生成数据密钥
func (c *VaultClient) GenerateDataKey(context string) (map[string]interface{}, error) {
    data := map[string]interface{}{
        "context": context, // 密钥派生上下文
    }

    secret, err := c.client.Logical().Write(
        c.path+"/datakey/plaintext/aes-256-gcm96",
        data,
    )
    if err != nil {
        return nil, err
    }

    return secret.Data, nil
}

// Encrypt 使用 Vault Transit 加密数据
func (c *VaultClient) Encrypt(plaintext []byte, context string) (string, error) {
    data := map[string]interface{}{
        "plaintext": base64.StdEncoding.EncodeToString(plaintext),
        "context":   context,
    }

    secret, err := c.client.Logical().Write(
        c.path+"/encrypt/my-key",
        data,
    )
    if err != nil {
        return "", err
    }

    return secret.Data["ciphertext"].(string), nil
}

// Decrypt 使用 Vault Transit 解密数据
func (c *VaultClient) Decrypt(ciphertext, context string) ([]byte, error) {
    data := map[string]interface{}{
        "ciphertext": ciphertext,
        "context":    context,
    }

    secret, err := c.client.Logical().Write(
        c.path+"/decrypt/my-key",
        data,
    )
    if err != nil {
        return nil, err
    }

    return base64.StdEncoding.DecodeString(
        secret.Data["plaintext"].(string),
    )
}
```

**优点：**
- 自托管，完全控制
- 支持多种后端（Consul, etcd, S3 等）
- 丰富的访问控制策略
- 支持动态密钥和 PKI

**缺点：**
- 需要维护 Vault 集群
- 需要配置高可用

---

### 方案三：YubiHSM 2 (推荐用于最高安全要求)

```go
package kms

// YubiHSMClient YubiHSM 2 客户端
// 注：实际实现需要使用 YubiHSM Go SDK

type YubiHSMClient struct {
    connector *yubihsm.Connector
    session   *yubihsm.Session
    authKeyID uint16
}

// NewYubiHSMClient 创建新的 YubiHSM 客户端
func NewYubiHSMClient(url string, password string) (*YubiHSMClient, error) {
    connector, err := yubihsm.NewConnector(url)
    if err != nil {
        return nil, err
    }

    // 使用默认认证密钥 (0x0001) 创建会话
    session, err := connector.CreateSession(0x0001, password)
    if err != nil {
        return nil, err
    }

    return &YubiHSMClient{
        connector: connector,
        session:   session,
        authKeyID: 0x0001,
    }, nil
}

// GenerateAESKey 在 HSM 中生成 AES 密钥
func (c *YubiHSMClient) GenerateAESKey(keyID uint16, label string) error {
    // 密钥永远存储在 HSM 中，无法导出
    _, err := c.session.GenerateAESKey(
        keyID,
        label,
        yubihsm.DomainAll,
        yubihsm.CapabilityEncryptAES | yubihsm.CapabilityDecryptAES,
        yubihsm.AlgorithmAES256,
    )
    return err
}

// Encrypt 使用 HSM 加密数据
func (c *YubiHSMClient) Encrypt(keyID uint16, plaintext []byte) ([]byte, error) {
    return c.session.EncryptAES(keyID, plaintext)
}

// Decrypt 使用 HSM 解密数据
func (c *YubiHSMClient) Decrypt(keyID uint16, ciphertext []byte) ([]byte, error) {
    return c.session.DecryptAES(keyID, ciphertext)
}

// Close 关闭会话
func (c *YubiHSMClient) Close() {
    c.session.Close()
    c.connector.Close()
}
```

**优点：**
- 硬件级安全，密钥永不离开 HSM
- FIPS 140-2 Level 3 合规
- 物理防篡改保护

**缺点：**
- 需要物理设备
- 成本较高
- 单点故障风险（需配置集群）

---

## 推荐部署方案

### 开发环境

```yaml
# 使用软件加密（当前实现）
kms:
  provider: "software"
  config:
    master_key: "${MASTER_KEY_ENV}"  # 从环境变量读取
```

### 测试环境

```yaml
# 使用 Vault 开发模式
kms:
  provider: "vault"
  config:
    address: "http://vault:8200"
    path: "transit"
    token: "${VAULT_TOKEN}"
```

### 生产环境

```yaml
# 使用 AWS KMS 或 Vault 集群
kms:
  provider: "aws_kms"  # 或 "vault"
  config:
    region: "us-east-1"
    key_id: "arn:aws:kms:..."
    # Vault 配置备选
    # address: "https://vault.production.local"
    # path: "transit/llmcustody"
```

---

## 安全最佳实践

### 1. 密钥零化 (Secret Zeroization)

```go
// 敏感数据使用后立即擦除
func ProcessAPIKey(encryptedKey []byte) {
    // 解密
    decrypted, err := decrypt(encryptedKey)
    if err != nil {
        return
    }

    // 确保零化
    defer func() {
        for i := range decrypted {
            decrypted[i] = 0
        }
    }()

    // 使用密钥...
    useKey(decrypted)
}
```

### 2. 内存保护

```go
// 使用 mlock 防止内存交换到磁盘
import "golang.org/x/sys/unix"

func SecureAlloc(size int) ([]byte, error) {
    data := make([]byte, size)
    if err := unix.Mlock(data); err != nil {
        return nil, err
    }
    return data, nil
}
```

### 3. 审计日志

```go
// 所有密钥操作记录审计日志
type AuditLog struct {
    Timestamp   time.Time `json:"timestamp"`
    Action      string    `json:"action"`      // encrypt, decrypt, rotate
    KeyID       string    `json:"key_id"`
    UserID      string    `json:"user_id"`
    Success     bool      `json:"success"`
    Error       string    `json:"error,omitempty"`
    IPAddress   string    `json:"ip_address"`
}
```

### 4. 密钥轮换

```go
// 定期轮换策略
type KeyRotationPolicy struct {
    Automatic   bool          `json:"automatic"`
    Interval    time.Duration `json:"interval"`     // 轮换间隔
    GracePeriod time.Duration `json:"grace_period"` // 旧密钥保留期
    NotifyBefore time.Duration `json:"notify_before"` // 提前通知时间
}
```

---

## 部署检查清单

- [ ] KMS/HSM 已配置并测试
- [ ] 主密钥已生成并备份
- [ ] 密钥轮换策略已配置
- [ ] 审计日志已启用
- [ ] 监控告警已配置
- [ ] 灾难恢复计划已测试
- [ ] 安全审计已通过

---

## 参考文档

- [AWS KMS Best Practices](https://docs.aws.amazon.com/kms/latest/developerguide/best-practices.html)
- [HashiCorp Vault Documentation](https://www.vaultproject.io/docs)
- [YubiHSM 2 Documentation](https://developers.yubico.com/YubiHSM2/)
- [NIST Key Management Guidelines](https://csrc.nist.gov/projects/key-management)
