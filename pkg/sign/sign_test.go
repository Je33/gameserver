package sign_test

import (
	"server/pkg/sign"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifySignature_Success(t *testing.T) {
	err := sign.VerifySignature(
		"0xeF209Bee800Ef5c7d20A67F46E007a970EAf9935",
		"test",
		"0x499cf8ce848eac151a49d23f95ce3fbfc7bf9bac709445458ca34de7afe8d98f2b9e2a34b4f7cf447145efef707355da665189cdea83a7464c804bda4cec556400",
	)

	assert.NoError(t, err)
}

func TestVerifySignature_Fail(t *testing.T) {
	testCases := [4][3]string{
		{
			"0xeF1c8b8c7f478c0BE246735c06aE80BEA3675D75",
			"test",
			"0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		},
		{
			"0xeF1c8b8c7f478c0BE246735c06aE80BEA3675D75",
			"test",
			"0x499cf8ce848eac151a49d23f95ce3fbfc7bf9bac709445458ca34de7afe8d98f2b9e2a34b4f7cf447145efef707355da665189cdea83a7464c804bda4cec556400",
		},
		{
			"0xeF1c8b8c7f478c0BE246735c06aE80BEA3675D75",
			"test",
			"dasdasd",
		},
		{
			"0x28e582BA14CD679FB08E47dC50b565c715Ce0979",
			"test1",
			"0x51118c7e70dc23560323aef56eba52ab989bbca085383b055be666d63b34dca8057fb210665ce5114eb53bf6b6644b7eddeeeaaf27b033738c14fa0f98d6d61f00",
		},
	}

	for _, testCase := range testCases {
		err := sign.VerifySignature(
			testCase[0],
			testCase[1],
			testCase[2],
		)

		assert.Error(t, err)
	}
}
