package interchaintest

import (
	"context"
	"fmt"

	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
)

func IBCTransferWithCustomFee(c *cosmos.CosmosChain, ctx context.Context,
	channelID string,
	keyName string,
	amount ibc.WalletAmount,
	fee ibc.WalletAmount,
	options ibc.TransferOptions,
) (tx ibc.Tx, err error) {
	tn := c.FullNodes[0]
	command := []string{
		"ibc-transfer", "transfer", "transfer", channelID,
		amount.Address, fmt.Sprintf("%s%s", amount.Amount.String(), amount.Denom), "--fee", fmt.Sprintf("%s%s", fee.Amount.String(), fee.Denom),
	}

	if options.Timeout != nil {
		if options.Timeout.NanoSeconds > 0 {
			command = append(command, "--packet-timeout-timestamp", fmt.Sprint(options.Timeout.NanoSeconds))
		} else if options.Timeout.Height > 0 {
			command = append(command, "--packet-timeout-height", fmt.Sprintf("0-%d", options.Timeout.Height))
		}
	}
	if options.Memo != "" {
		command = append(command, "--memo", options.Memo)
	}
	txHash, err := tn.ExecTx(ctx, keyName, command...)
	if err != nil {
		return tx, fmt.Errorf("send ibc transfer: %w", err)
	}
	txResp, err := c.GetTransaction(txHash)
	if err != nil {
		return tx, fmt.Errorf("failed to get transaction %s: %w", txHash, err)
	}
	if txResp.Code != 0 {
		return tx, fmt.Errorf("error in transaction (code: %d): %s", txResp.Code, txResp.RawLog)
	}
	tx.Height = uint64(txResp.Height)
	tx.TxHash = txHash
	// In cosmos, user is charged for entire gas requested, not the actual gas used.
	tx.GasSpent = txResp.GasWanted

	return tx, nil
}
