package main

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

// jcoin
// denominated by .001

type Miner struct {
	txnPool *TxnMemoryPool
	JCoin *JCoin
	coefficient int 
	exponent int 
}

func NewMiner (txnPool *TxnMemoryPool, JCoin *JCoin, coefficient int, exponent int) *Miner {
	miner := &Miner{
		txnPool: 		txnPool,
		JCoin: 	JCoin,
		coefficient: 	coefficient,
		exponent: 		exponent,
	}
	return miner
}

func (m *Miner) Mine () {
	target := uint64(m.coefficient) * uint64(math.Pow(2, float64(0x8 * (m.exponent - 0x3))))
	bits := 0x207fffff
	var nonce uint32  
	nonce = 0 
	// wait for max transactions
	for len(m.txnPool.Transactions) < 9 {}

	// collect recent transactions
	transactions := m.txnPool.Transactions[:9]
	m.txnPool.Transactions = m.txnPool.Transactions[9:]

	// add coinbase transaction
	coinbaseTransaction := NewTransaction(0, 1, []string{"coinbase"}, 1, []*Output{NewOutput(m.JCoin.Reward, 0, "coinbase")})
	transactions = append(transactions, coinbaseTransaction)
	
	merkleHash := buildMerkleTree(transactions)
	for {
		header := NewHeader(0, m.JCoin.Blocks[len(m.JCoin.Blocks)-1].Blockhash, merkleHash, bits, nonce)
		block := NewBlock(0xD9B4BEF9, 10, header, transactions)
		blockHashBytes, _ := hex.DecodeString(block.Blockhash)
		blockHashVal := uint64(0)
		for _, b := range blockHashBytes {
			blockHashVal = (blockHashVal << 8) | uint64(b)
		}
		unsignedBlockHashVal := uint64(blockHashVal)
		if unsignedBlockHashVal < target {
			m.JCoin.AddBlock(block)
			return
		}
		nonce++
	}
}

func SimulateTransactions(txnPool *TxnMemoryPool) {
	transactionCount := 0
	for {
		sleepTime := rand.Intn(2) + 3
		time.Sleep(time.Duration(sleepTime) * time.Second) 
		tx := NewTransaction(0, 1, []string{strconv.FormatInt(rand.Int63n(1000), 10), 
			string(time.DateTime)}, 1, []*Output{NewOutput(rand.Intn(200), 0, "transaction")})
		txnPool.AddTransactionToPool(tx)
		// search a transaction by hashs
		transactionCount++
		fmt.Println("\nAdding Transaction", transactionCount, ":")
		tx.printTransaction()
	}
}