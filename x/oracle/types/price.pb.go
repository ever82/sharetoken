// Code generated manually for protobuf support. DO NOT EDIT.
// source: sharetoken/oracle/price.proto

package types

import (
	"github.com/cosmos/gogoproto/proto"
)

// PriceProto is the protobuf message for Price
type PriceProto struct {
	Symbol     string `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Price      string `protobuf:"bytes,2,opt,name=price,proto3" json:"price,omitempty"`
	Timestamp  int64  `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Source     int32  `protobuf:"varint,4,opt,name=source,proto3" json:"source,omitempty"`
	Confidence int32  `protobuf:"varint,5,opt,name=confidence,proto3" json:"confidence,omitempty"`
}

func (m *PriceProto) Reset()         { *m = PriceProto{} }
func (m *PriceProto) String() string { return proto.CompactTextString(m) }
func (*PriceProto) ProtoMessage()    {}

func (m *PriceProto) GetSymbol() string     { return m.Symbol }
func (m *PriceProto) GetPrice() string      { return m.Price }
func (m *PriceProto) GetTimestamp() int64   { return m.Timestamp }
func (m *PriceProto) GetSource() int32      { return m.Source }
func (m *PriceProto) GetConfidence() int32  { return m.Confidence }

// LLMPriceProto is the protobuf message for LLMPrice
type LLMPriceProto struct {
	Provider    string `protobuf:"bytes,1,opt,name=provider,proto3" json:"provider,omitempty"`
	Model       string `protobuf:"bytes,2,opt,name=model,proto3" json:"model,omitempty"`
	InputPrice  string `protobuf:"bytes,3,opt,name=input_price,json=inputPrice,proto3" json:"input_price,omitempty"`
	OutputPrice string `protobuf:"bytes,4,opt,name=output_price,json=outputPrice,proto3" json:"output_price,omitempty"`
	Currency    string `protobuf:"bytes,5,opt,name=currency,proto3" json:"currency,omitempty"`
}

func (m *LLMPriceProto) Reset()         { *m = LLMPriceProto{} }
func (m *LLMPriceProto) String() string { return proto.CompactTextString(m) }
func (*LLMPriceProto) ProtoMessage()    {}

func (m *LLMPriceProto) GetProvider() string    { return m.Provider }
func (m *LLMPriceProto) GetModel() string       { return m.Model }
func (m *LLMPriceProto) GetInputPrice() string  { return m.InputPrice }
func (m *LLMPriceProto) GetOutputPrice() string { return m.OutputPrice }
func (m *LLMPriceProto) GetCurrency() string    { return m.Currency }

// GenesisStateProto is the protobuf message for GenesisState
type GenesisStateProto struct {
	Prices []*PriceProto `protobuf:"bytes,1,rep,name=prices,proto3" json:"prices,omitempty"`
}

func (m *GenesisStateProto) Reset()         { *m = GenesisStateProto{} }
func (m *GenesisStateProto) String() string { return proto.CompactTextString(m) }
func (*GenesisStateProto) ProtoMessage()    {}

func (m *GenesisStateProto) GetPrices() []*PriceProto { return m.Prices }

func init() {
	proto.RegisterType((*PriceProto)(nil), "sharetoken.oracle.Price")
	proto.RegisterType((*LLMPriceProto)(nil), "sharetoken.oracle.LLMPrice")
	proto.RegisterType((*GenesisStateProto)(nil), "sharetoken.oracle.GenesisState")
}
