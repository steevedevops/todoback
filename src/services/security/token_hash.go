package security

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"

	"github.com/thanhpk/randstr"
)

func TokenHashmd5(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	theHash := hex.EncodeToString(hasher.Sum(nil))
	theToken := theHash
	return theToken
}

func TokenRandom() string {
	return randstr.String(40)
}

func TokenHashsha256(text string) string {
	hasher := sha256.New()
	hasher.Write([]byte(text))
	theHash := hex.EncodeToString(hasher.Sum(nil))
	theToken := theHash
	return theToken
}

func TokenHashsha512(text string) string {
	hasher := sha512.New()
	hasher.Write([]byte(text))
	theHash := hex.EncodeToString(hasher.Sum(nil))
	theToken := theHash
	return theToken
}
