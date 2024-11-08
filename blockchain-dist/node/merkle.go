package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// double hash function
func doubleSHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	hash = sha256.Sum256(hash[:])
	return hash[:]
}

// hashes two branches
func hashPair(left, right []byte) []byte {
	concatenated := append(left, right...)
	return doubleSHA256(concatenated)
}

// takes the transactions and builds a tree
func buildMerkleTree(transactions []*Transaction) string {
	var hashes [][]byte 

	for i := 0; i < len(transactions); i++ {
		hash := transactions[i].TransactionHash
		hashBytes, err := hex.DecodeString(hash)
		if err != nil {
			fmt.Println("Error decoding hash:", err)
			return ""
		}
		hashes = append(hashes, hashBytes)
	}

	for len(hashes) > 1 {
		var newLevel [][]byte
		for i := 0; i < len(hashes); i += 2 {
			if i+1 == len(hashes) {
				newLevel = append(newLevel, hashPair(hashes[i], hashes[i]))
			} else {
				newLevel = append(newLevel, hashPair(hashes[i], hashes[i+1]))
			}
		}
		hashes = newLevel
	}
	return hex.EncodeToString(hashes[0])
}
