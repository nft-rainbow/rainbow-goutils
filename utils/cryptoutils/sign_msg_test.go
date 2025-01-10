package cryptoutils

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestRecoverAddressFromSign(t *testing.T) {
	message := "hello world"
	signature := "0x3afdd93eabeaa8bec851a0ed0a6c9703c3504092fea6b9f6ad96e97ee6d8ef887f310bc41a720b418043682cfe3a85bf6dd8b39160096833f8738a43433259d61b" // 这里需要替换为实际的签名

	actual, err := RecoverMessage(message, signature)
	if err != nil {
		t.Error(err)
	}

	expect := common.HexToAddress("0xBeD38c825459994002257DFBB88371E243204B6c")
	if actual != expect {
		t.Errorf("Message signature verification failed. Expected address: %s, Actual recovered address: %s", expect.Hex(), actual.Hex())
	}
}
