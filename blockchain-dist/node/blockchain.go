package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// IF YOU WANT TO INCLUDE TRANSACTION OUTPUTS IN THE HASH THEN ADJUST THE HASHING FUNCTION
type Transaction struct {
    VersionNumber   int
    InCounter       int
    ListOfInputs    []string
    OutCounter      int
    ListOfOutputs   []*Output
	OutputsString   string
    TransactionHash string
}

func NewTransaction(versionNumber, inCounter int, inputs []string, outCounter int, outputs []*Output) *Transaction {
    tx := &Transaction{
        VersionNumber: versionNumber,
        InCounter:     inCounter,
        ListOfInputs:  inputs,
        OutCounter:    outCounter,
        ListOfOutputs: outputs,
    }
	tx.OutputsString = outputsString(outputs)
    tx.TransactionHash = tx.calculateTransactionHash()
    return tx
}

func (tx *Transaction) calculateTransactionHash() string {
    data := strconv.Itoa(tx.VersionNumber) + strconv.Itoa(tx.InCounter) + strings.Join(tx.ListOfInputs, "") +
        strconv.Itoa(tx.OutCounter) + tx.OutputsString
    return SHA256DoubleHash(data)
}

func (tx *Transaction) printTransaction() {
    fmt.Printf("Transaction Hash: %s\n", tx.TransactionHash)
    fmt.Printf("Version: %d\n", tx.VersionNumber)
    fmt.Printf("Inputs: %v\n", tx.ListOfInputs)
    fmt.Printf("Outputs: %v\n", tx.ListOfOutputs)
}

type Header struct {
    Version         int
    hashPrevBlock   string
    hashMerkleRoot  string
    Timestamp       int64
    Bits            int
    Nonce           uint32
}

func NewHeader(version int, hashPrevBlock, hashMerkleRoot string, bits int, nonce uint32) *Header {
    return &Header{
        Version:        version,
        hashPrevBlock:  hashPrevBlock,
        hashMerkleRoot: hashMerkleRoot,
        Timestamp:      time.Now().Unix(),
        Bits:           bits,
        Nonce:          nonce,
    }
}

type Block struct {
    MagicNumber        int
    Blocksize          int
    BlockHeader        *Header
    TransactionCounter int
    Transactions       []*Transaction
    Blockhash          string
}

func NewBlock(magicNumber, blocksize int, header *Header, transactions []*Transaction) *Block {
    block := &Block{
        MagicNumber:        magicNumber,
        Blocksize:          blocksize,
        BlockHeader:        header,
        TransactionCounter: len(transactions),
        Transactions:       transactions,
    }
    block.Blockhash = block.calculateBlockHash()
    return block
}

func (b *Block) calculateBlockHash() string {
    data := strconv.FormatInt(b.BlockHeader.Timestamp, 10) + b.BlockHeader.hashMerkleRoot + strconv.Itoa(b.BlockHeader.Bits) +
        strconv.Itoa(int(b.BlockHeader.Nonce)) + b.BlockHeader.hashPrevBlock
    return SHA256DoubleHash(data)
}

func (b *Block) printBlock() {
    fmt.Printf("Block Hash: %s\n", b.Blockhash)
    fmt.Printf("Magic Number: %d\n", b.MagicNumber)
    fmt.Printf("Blocksize: %d\n", b.Blocksize)
    fmt.Printf("Transaction Count: %d\n", b.TransactionCounter)
    for _, tx := range b.Transactions {
        tx.printTransaction()
    }
}

type JCoin struct {
    Blocks   []*Block
	MAX_TXNS int
	Reward int
}

func NewJCoin() *JCoin {

	genesisOutput := NewOutput(1, 0, "genesis")
	genesisTransaction := NewTransaction(1, 1, []string{"genesis_input"}, 1, []*Output{genesisOutput})
	genesisMerkleRoot := buildMerkleTree([]*Transaction{genesisTransaction})
    genesisHeader := NewHeader(1, strings.Repeat("0", 64), genesisMerkleRoot, 0, 0)
    genesisBlock := NewBlock(0xD9B4BEF9, 0, genesisHeader, []*Transaction{genesisTransaction})

    return &JCoin{
        Blocks: []*Block{genesisBlock},
		MAX_TXNS: 10,
		Reward: 5447,
    }
}

func (bc *JCoin) AddBlock(newBlock *Block) int {
	// throws an error if num transactions > max transactions
	if len(newBlock.Transactions) > bc.MAX_TXNS {
		return 1
	}

	bc.Blocks = append(bc.Blocks, newBlock)
	return 0
}

func (bc *JCoin) GetBlockByHeight(height int) *Block {
    if height < len(bc.Blocks) {
        return bc.Blocks[height]
    }
    return nil
}

func (bc *JCoin) GetBlockByHash(hash string) *Block {
    for _, block := range bc.Blocks {
        if block.Blockhash == hash {
            return block
        }
    }
    return nil
}

func (bc *JCoin) GetTransactionByHash(hash string) *Transaction {
    for _, block := range bc.Blocks {
        for _, tx := range block.Transactions {
            if tx.TransactionHash == hash {
                return tx
            }
        }
    }
    return nil
}

// func main() {

// 	miner := NewMiner(transactionPool, JCoin, 0x7fffff, 0x20)

// 	numBlocks := 0
// 	for {
// 		miner.Mine()
// 		numBlocks++
// 		fmt.Println("\nBlock at height ", numBlocks, ":")
// 		block := JCoin.GetBlockByHeight(numBlocks)
// 		if block != nil {
// 			block.printBlock()
// 		}
// 	}
// }
