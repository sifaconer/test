package usecase

import (
	"api-test/src/common"
	"api-test/src/config"
	"encoding/base64"
	"testing"
)

func Test_encryption_GenerateMasterDemoKey(t *testing.T) {
	e := &encryption{
		log:    common.NewLogger(),
		config: config.NewConfig(),
	}
	key, _ := e.GenerateRandomKey()

	keyBase64 := base64.StdEncoding.EncodeToString(key)
	t.Logf("Nueva clave maestra (Base64): %s\n", keyBase64)
}
