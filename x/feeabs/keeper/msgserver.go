package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/notional-labs/feeabstraction/v1/x/feeabs/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{
		Keeper: keeper,
	}
}

var _ types.MsgServer = msgServer{}

// Need to remove this
func (k Keeper) SendQuerySpotPrice(goCtx context.Context, msg *types.MsgSendQuerySpotPrice) (*types.MsgSendQuerySpotPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	denom := "ibc/E7D5E9D0E9BF8B7354929A817DD28D4D017E745F638954764AA88522A7A409EC"
	hostChainConfig, err := k.GetHostZoneConfig(ctx, denom)
	if err != nil {
		return &types.MsgSendQuerySpotPriceResponse{}, nil
	}
	_, err = sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, err
	}
	err = k.handleOsmosisIbcQuery(ctx, hostChainConfig)
	if err != nil {
		return nil, err
	}

	return &types.MsgSendQuerySpotPriceResponse{}, nil
}

// Need to remove this
func (k Keeper) SwapCrossChain(goCtx context.Context, msg *types.MsgSwapCrossChain) (*types.MsgSwapCrossChainResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	denom := "ibc/E7D5E9D0E9BF8B7354929A817DD28D4D017E745F638954764AA88522A7A409EC"
	hostChainConfig, err := k.GetHostZoneConfig(ctx, denom)
	if err != nil {
		return &types.MsgSwapCrossChainResponse{}, nil
	}
	_, err = sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, err
	}
	err = k.transferIBCTokenToOsmosisContract(ctx, hostChainConfig)
	if err != nil {
		return nil, err
	}

	return &types.MsgSwapCrossChainResponse{}, nil
}
