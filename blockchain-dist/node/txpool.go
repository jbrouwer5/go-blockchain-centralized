package main

// stores the current unverified transactions
type TxnMemoryPool struct {
	Transactions []*Transaction
}

func NewTxnMemoryPool() *TxnMemoryPool {
	return &TxnMemoryPool{
		Transactions: []*Transaction{},
	}
}

func (txMemPool *TxnMemoryPool) AddTransactionToPool(tx *Transaction) {
	txMemPool.Transactions = append(txMemPool.Transactions, tx)
}
