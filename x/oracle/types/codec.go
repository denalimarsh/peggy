package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {}

// var (
// 	amino     = codec.New()
// 	ModuleCdc = codec.NewHybridCodec(amino)
// )

// func init() {
// 	RegisterCodec(amino)
// 	codec.RegisterCrypto(amino)
// 	amino.Seal() // TODO: required?
// }
