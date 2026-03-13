// Code generated manually for protobuf support. DO NOT EDIT.
// source: sharetoken/marketplace/service.proto

package types

import (
	"github.com/cosmos/gogoproto/proto"
)

// CoinProto is the protobuf message for Coin
type CoinProto struct {
	Denom  string `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
	Amount string `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (m *CoinProto) Reset()         { *m = CoinProto{} }
func (m *CoinProto) String() string { return proto.CompactTextString(m) }
func (*CoinProto) ProtoMessage()    {}

func (m *CoinProto) GetDenom() string  { return m.Denom }
func (m *CoinProto) GetAmount() string { return m.Amount }

// ServiceProto is the protobuf message for Service
type ServiceProto struct {
	Id          string       `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Provider    string       `protobuf:"bytes,2,opt,name=provider,proto3" json:"provider,omitempty"`
	Name        string       `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Description string       `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	Level       int32        `protobuf:"varint,5,opt,name=level,proto3" json:"level,omitempty"`
	PricingMode int32        `protobuf:"varint,6,opt,name=pricing_mode,json=pricingMode,proto3" json:"pricing_mode,omitempty"`
	Price       []*CoinProto `protobuf:"bytes,7,rep,name=price,proto3" json:"price,omitempty"`
	Active      bool         `protobuf:"varint,8,opt,name=active,proto3" json:"active,omitempty"`
	CreatedAt   int64        `protobuf:"varint,9,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
}

func (m *ServiceProto) Reset()         { *m = ServiceProto{} }
func (m *ServiceProto) String() string { return proto.CompactTextString(m) }
func (*ServiceProto) ProtoMessage()    {}

func (m *ServiceProto) GetId() string          { return m.Id }
func (m *ServiceProto) GetProvider() string    { return m.Provider }
func (m *ServiceProto) GetName() string        { return m.Name }
func (m *ServiceProto) GetDescription() string { return m.Description }
func (m *ServiceProto) GetLevel() int32        { return m.Level }
func (m *ServiceProto) GetPricingMode() int32  { return m.PricingMode }
func (m *ServiceProto) GetPrice() []*CoinProto { return m.Price }
func (m *ServiceProto) GetActive() bool        { return m.Active }
func (m *ServiceProto) GetCreatedAt() int64    { return m.CreatedAt }

// GenesisStateProto is the protobuf message for GenesisState
type GenesisStateProto struct {
	Services []*ServiceProto `protobuf:"bytes,1,rep,name=services,proto3" json:"services,omitempty"`
}

func (m *GenesisStateProto) Reset()         { *m = GenesisStateProto{} }
func (m *GenesisStateProto) String() string { return proto.CompactTextString(m) }
func (*GenesisStateProto) ProtoMessage()    {}

func (m *GenesisStateProto) GetServices() []*ServiceProto { return m.Services }

func init() {
	proto.RegisterType((*CoinProto)(nil), "sharetoken.marketplace.Coin")
	proto.RegisterType((*ServiceProto)(nil), "sharetoken.marketplace.Service")
	proto.RegisterType((*GenesisStateProto)(nil), "sharetoken.marketplace.GenesisState")
}
