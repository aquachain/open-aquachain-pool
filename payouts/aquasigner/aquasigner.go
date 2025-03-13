package aquasigner

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"gitlab.com/aquachain/aquachain/common"
	"gitlab.com/aquachain/aquachain/core/types"
	"gitlab.com/aquachain/aquachain/crypto"
	"gitlab.com/aquachain/aquachain/rlp"
)

type Signer struct {
	k *btcec.PrivateKey
}

func Enabled() bool {
	return private_key != nil
}

func Account() string {
	if Enabled() {
		return crypto.PubkeyToAddress(private_key.PubKey()).Hex()
	}
	return ""
}

func SignTx(chainId *big.Int, tx *types.Transaction) (txhash string, rawtx string, err error) {
	if !Enabled() {
		return "", "", fmt.Errorf("signing is disabled")
	}
	return Signer{k: private_key}.SignTx(chainId, tx)
}

// SignTx without broadcasting, returns txhash and hex raw signed tx
func (s Signer) SignTx(chainId *big.Int, tx *types.Transaction) (string, string, error) {
	// confirm some stuff
	if chainId == nil || chainId.Sign() != 1 || tx == nil {
		return "", "", fmt.Errorf("invalid input")
	}
	to := tx.To()
	if to == nil {
		return "", "", fmt.Errorf("invalid recipient address")
	}
	if *to == (common.Address{}) {
		return "", "", fmt.Errorf("not sending to zero address")
	}
	// sign tx offline
	eip155signer := types.NewEIP155Signer(chainId)
	signedTx, err := types.SignTx(tx, eip155signer, s.k)
	if err != nil {
		return "", "", fmt.Errorf("could not sign transaction: %v", err)
	}
	txhash := signedTx.Hash().Hex()
	// hex raw signed
	raw, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return "", "", fmt.Errorf("could not encode signed transaction: %v", err)
	}
	hexraw := fmt.Sprintf("0x%02x", raw)
	log.Printf("signed tx to %s: /tx/%s?raw=%x\n", to.Hex(), txhash, hexraw)
	return txhash, hexraw, nil
}

// if file, it only needs to exist at boot time then can be shredded
var private_key = func() *btcec.PrivateKey {
	got := os.Getenv("PRIVATE_KEY_HEX")
	os.Setenv("PRIVATE_KEY_HEX", "") // clear it
	filename := os.Getenv("PRIVATE_KEY_HEXFILE")
	os.Setenv("PRIVATE_KEY_HEXFILE", "") // clear it
	if len(got) == 0 && len(filename) != 0 {
		got0, err := os.ReadFile(filename)
		if err != nil {
			log.Fatalf("could not read PRIVATE_KEY_HEXFILE: %v", err)
		}
		got = string(got0)
	}
	got = strings.TrimSpace(got) // removes newlines
	if len(got) != 0 {
		key, err := crypto.HexToBtcec(got)
		if err != nil {
			log.Fatalf("could not parse signing key: %v", err)
		}
		log.Println("local signer enabled")
		return key
	}
	return nil
}()
