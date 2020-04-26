package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"

	"github.com/cosmos/peggy/x/oracle"
)

// AccountKeeper defines the expected account keeper
type AccountKeeper interface {
	GetAccount(sdk.Context, sdk.AccAddress) authexported.Account
	// SetModuleAccount(ctx sdk.Context, macc exported.ModuleAccountI)
}

// BankKeeper defines the expected bank keeper
type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error
}

// OracleKeeper defines the expected oracle keeper
type OracleKeeper interface {
	ProcessClaim(ctx sdk.Context, claim oracle.Claim) (oracle.Status, error)
	GetProphecy(ctx sdk.Context, id string) (oracle.Prophecy, bool)
}
