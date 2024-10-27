package service

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/48club/rpc-watchdog/config"
	types2 "github.com/48club/rpc-watchdog/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Watch(c chan types2.Chan) {
	for _, v := range config.Config.RpcList {
		go checkLoop(c, v)
	}
}

func checkLoop(c chan types2.Chan, rpc string) {
	for {
		c <- types2.Chan{
			Rpc: rpc,
			Err: checkRpc(rpc),
		}

		<-time.After(config.Config.Interval * time.Second)
	}
}

func checkRpc(rpc string) error {
	ec, err := ethclient.Dial(rpc)
	if err != nil {
		return err
	}

	defer ec.Close()

	h, err := ec.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return err
	}

	if tt := time.Now().Unix() - int64(h.Time); tt >= config.Config.AllowSlow {
		return fmt.Errorf("rpc %s is slow, %d seconds slower", rpc, tt)
	}

	return fakeTx(ec)
}

var to = common.HexToAddress("0x2164D118329b03677710127d7878a57Db2b1edc5")

func fakeTx(ec *ethclient.Client) error {
	pk, err := crypto.GenerateKey()
	if err != nil {
		return nil // ignore error
	}
	publicKey := pk.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil // ignore error
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := ec.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}

	value := big.NewInt(1e18)
	var data []byte

	gasLimit, err := ec.EstimateGas(context.Background(), ethereum.CallMsg{To: &to, Data: data})
	if err != nil {
		return err
	}

	gasPrice, err := ec.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	tx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(56)), pk)
	if err != nil {
		return nil // ignore error
	}

	err = ec.SendTransaction(context.Background(), signedTx)
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "insufficient balance") {
		return nil // ignore error
	}
	return err
}
