package crypto

import (
	"crypto/sha512"
	"encoding/hex"
	"strings"
)

//https://github.com/hyperledger/sawtooth-sdk-go/blob/master/examples/intkey_go/src/sawtooth_intkey/handler/handler.go
func Hexdigest(str string) string {
	hash := sha512.New()
	hash.Write([]byte(str))
	hashBytes := hash.Sum(nil)
	return strings.ToLower(hex.EncodeToString(hashBytes))
}

