package main

import (
	"crypto/sha256"
	"encoding/hex"
)

func SHA256DoubleHash(data string) string {
    hash := sha256.Sum256([]byte(data))
    hash = sha256.Sum256(hash[:])
    return hex.EncodeToString(hash[:])
}
