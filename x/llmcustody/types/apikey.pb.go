// Code generated manually for protobuf support. DO NOT EDIT.
// source: sharetoken/llmcustody/apikey.proto

package types

import (
	"github.com/cosmos/gogoproto/proto"
)

// ProviderProto represents the provider type in protobuf
type ProviderProto int32

const (
	ProviderProto_PROVIDER_UNSPECIFIED ProviderProto = 0
	ProviderProto_PROVIDER_OPENAI      ProviderProto = 1
	ProviderProto_PROVIDER_ANTHROPIC   ProviderProto = 2
)

// AccessRuleProto is the protobuf message for AccessRule
type AccessRuleProto struct {
	ServiceId   string `protobuf:"bytes,1,opt,name=service_id,json=serviceId,proto3" json:"service_id,omitempty"`
	Allowed     bool   `protobuf:"varint,2,opt,name=allowed,proto3" json:"allowed,omitempty"`
	RateLimit   int64  `protobuf:"varint,3,opt,name=rate_limit,json=rateLimit,proto3" json:"rate_limit,omitempty"`
	MaxRequests int64  `protobuf:"varint,4,opt,name=max_requests,json=maxRequests,proto3" json:"max_requests,omitempty"`
	PricePerReq int64  `protobuf:"varint,5,opt,name=price_per_req,json=pricePerReq,proto3" json:"price_per_req,omitempty"`
}

func (m *AccessRuleProto) Reset()         { *m = AccessRuleProto{} }
func (m *AccessRuleProto) String() string { return proto.CompactTextString(m) }
func (*AccessRuleProto) ProtoMessage()    {}

// APIKeyProto is the protobuf message for APIKey
type APIKeyProto struct {
	Id           string             `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Provider     string             `protobuf:"bytes,2,opt,name=provider,proto3" json:"provider,omitempty"`
	EncryptedKey []byte             `protobuf:"bytes,3,opt,name=encrypted_key,json=encryptedKey,proto3" json:"encrypted_key,omitempty"`
	Hash         string             `protobuf:"bytes,4,opt,name=hash,proto3" json:"hash,omitempty"`
	Owner        string             `protobuf:"bytes,5,opt,name=owner,proto3" json:"owner,omitempty"`
	AccessRules  []*AccessRuleProto `protobuf:"bytes,6,rep,name=access_rules,json=accessRules,proto3" json:"access_rules,omitempty"`
	CreatedAt    int64              `protobuf:"varint,7,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	LastUsedAt   int64              `protobuf:"varint,8,opt,name=last_used_at,json=lastUsedAt,proto3" json:"last_used_at,omitempty"`
	UsageCount   int64              `protobuf:"varint,9,opt,name=usage_count,json=usageCount,proto3" json:"usage_count,omitempty"`
	Active       bool               `protobuf:"varint,10,opt,name=active,proto3" json:"active,omitempty"`
}

func (m *APIKeyProto) Reset()         { *m = APIKeyProto{} }
func (m *APIKeyProto) String() string { return proto.CompactTextString(m) }
func (*APIKeyProto) ProtoMessage()    {}

func (m *APIKeyProto) GetId() string           { return m.Id }
func (m *APIKeyProto) GetProvider() string     { return m.Provider }
func (m *APIKeyProto) GetEncryptedKey() []byte { return m.EncryptedKey }
func (m *APIKeyProto) GetHash() string         { return m.Hash }
func (m *APIKeyProto) GetOwner() string        { return m.Owner }
func (m *APIKeyProto) GetAccessRules() []*AccessRuleProto { return m.AccessRules }
func (m *APIKeyProto) GetCreatedAt() int64     { return m.CreatedAt }
func (m *APIKeyProto) GetLastUsedAt() int64    { return m.LastUsedAt }
func (m *APIKeyProto) GetUsageCount() int64    { return m.UsageCount }
func (m *APIKeyProto) GetActive() bool         { return m.Active }

// GenesisStateProto is the protobuf message for GenesisState
type GenesisStateProto struct {
	ApiKeys       []*APIKeyProto `protobuf:"bytes,1,rep,name=api_keys,json=apiKeys,proto3" json:"api_keys,omitempty"`
	EncryptionKey []byte         `protobuf:"bytes,2,opt,name=encryption_key,json=encryptionKey,proto3" json:"encryption_key,omitempty"`
}

func (m *GenesisStateProto) Reset()         { *m = GenesisStateProto{} }
func (m *GenesisStateProto) String() string { return proto.CompactTextString(m) }
func (*GenesisStateProto) ProtoMessage()    {}

func (m *GenesisStateProto) GetApiKeys() []*APIKeyProto { return m.ApiKeys }
func (m *GenesisStateProto) GetEncryptionKey() []byte   { return m.EncryptionKey }

func init() {
	proto.RegisterType((*AccessRuleProto)(nil), "sharetoken.llmcustody.AccessRule")
	proto.RegisterType((*APIKeyProto)(nil), "sharetoken.llmcustody.APIKey")
	proto.RegisterType((*GenesisStateProto)(nil), "sharetoken.llmcustody.GenesisState")
}
