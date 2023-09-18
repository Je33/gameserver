package sign

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func VerifySignature(wallet, message, sign string) error {

	hash := crypto.Keccak256Hash([]byte(message))

	signature, err := hexutil.Decode(sign)
	if err != nil {
		return err
	}

	sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), signature)
	if err != nil {
		return err
	}

	sigPublicKeyBytes := crypto.FromECDSAPub(sigPublicKeyECDSA)

	signatureNoRecoverID := signature[:len(signature)-1]
	verified := crypto.VerifySignature(sigPublicKeyBytes, hash.Bytes(), signatureNoRecoverID)
	if !verified {
		return fmt.Errorf("signature invalid")
	}

	sigAddress := crypto.PubkeyToAddress(*sigPublicKeyECDSA)
	if wallet != sigAddress.String() {
		return fmt.Errorf("signer invalid")
	}

	return nil
}
