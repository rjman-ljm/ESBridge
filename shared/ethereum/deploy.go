// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package utils

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChainSafe/chainbridge-utils/keystore"
	bridge "github.com/rjman-self/Platdot/bindings/Bridge"
	erc20Handler "github.com/rjman-self/Platdot/bindings/ERC20Handler"
	erc721Handler "github.com/rjman-self/Platdot/bindings/ERC721Handler"
)

var (
	RelayerAddresses = []common.Address{
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.AliceKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.BobKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.CharlieKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.DaveKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.EveKey].Address()),
	}

	ZeroAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")
)

type DeployedContracts struct {
	BridgeAddress         common.Address
	ERC20HandlerAddress   common.Address
}

// DeployContracts deploys Bridge, Relayer, ERC20Handler, ERC721Handler and CentrifugeAssetHandler and returns the addresses
func DeployContracts(client *Client, chainID uint8, initialRelayerThreshold *big.Int) (*DeployedContracts, error) {
	bridgeAddr, err := deployBridge(client, chainID, RelayerAddresses, initialRelayerThreshold)
	if err != nil {
		return nil, err
	}

	erc20HandlerAddr, err := deployERC20Handler(client, bridgeAddr)
	if err != nil {
		return nil, err
	}

	deployedContracts := DeployedContracts{bridgeAddr, erc20HandlerAddr}

	return &deployedContracts, nil

}

func UpdateNonce(client *Client) error {
	newNonce, err := client.Client.PendingNonceAt(context.Background(), client.CallOpts.From)
	if err != nil {
		return err
	}

	client.Opts.Nonce = big.NewInt(int64(newNonce))

	return nil
}

func deployBridge(client *Client, chainID uint8, relayerAddrs []common.Address, initialRelayerThreshold *big.Int) (common.Address, error) {
	err := client.LockNonceAndUpdate()
	if err != nil {
		return ZeroAddress, err
	}

	bridgeAddr, tx, _, err := bridge.DeployBridge(client.Opts, client.Client, chainID, relayerAddrs, initialRelayerThreshold, big.NewInt(0), big.NewInt(100))
	if err != nil {
		return ZeroAddress, err
	}

	err = WaitForTx(client, tx)
	if err != nil {
		return ZeroAddress, err
	}

	client.UnlockNonce()

	return bridgeAddr, nil

}

func deployERC20Handler(client *Client, bridgeAddress common.Address) (common.Address, error) {
	err := client.LockNonceAndUpdate()
	if err != nil {
		return ZeroAddress, err
	}

	erc20HandlerAddr, tx, _, err := erc20Handler.DeployERC20Handler(client.Opts, client.Client, bridgeAddress, [][32]byte{}, []common.Address{}, []common.Address{})
	if err != nil {
		return ZeroAddress, err
	}

	err = WaitForTx(client, tx)
	if err != nil {
		return ZeroAddress, err
	}

	client.UnlockNonce()

	return erc20HandlerAddr, nil
}

func deployERC721Handler(client *Client, bridgeAddress common.Address) (common.Address, error) {
	err := client.LockNonceAndUpdate()
	if err != nil {
		return ZeroAddress, err
	}

	erc721HandlerAddr, tx, _, err := erc721Handler.DeployERC721Handler(client.Opts, client.Client, bridgeAddress, [][32]byte{}, []common.Address{}, []common.Address{})
	if err != nil {
		return ZeroAddress, err
	}
	err = WaitForTx(client, tx)
	if err != nil {
		return ZeroAddress, err
	}

	client.UnlockNonce()

	return erc721HandlerAddr, nil
}

