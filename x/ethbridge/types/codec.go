package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

type Codec interface {
	codec.Marshaler
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateEthBridgeClaim{}, "ethbridge/MsgCreateEthBridgeClaim", nil)
	cdc.RegisterConcrete(MsgBurn{}, "ethbridge/MsgBurn", nil)
	cdc.RegisterConcrete(MsgLock{}, "ethbridge/MsgLock", nil)
}

var (
	amino     = codec.New()
	ModuleCdc = codec.NewHybridCodec(amino)
)

func init() {
	RegisterCodec(amino)
	codec.RegisterCrypto(amino)
	amino.Seal() // TODO: required?
}
