package main

import (
	pb "blockchain/proto"
)


func ConvertOutputToProto(output *Output) *pb.Output {
    return &pb.Output{
        Value:  int32(output.Value),
        Index:  int32(output.Index),
        Script: output.Script,
    }
}

// Convert Go Transaction struct to protobuf Transaction message
func ConvertTransactionToProto(tx *Transaction) *pb.Transaction {
    protoOutputs := make([]*pb.Output, len(tx.ListOfOutputs))
    for i, output := range tx.ListOfOutputs {
        protoOutputs[i] = ConvertOutputToProto(output)
    }

    return &pb.Transaction{
        VersionNumber:   int32(tx.VersionNumber),
        InCounter:       int32(tx.InCounter),
        ListOfInputs:    tx.ListOfInputs,
        OutCounter:      int32(tx.OutCounter),
        ListOfOutputs:   protoOutputs,
        OutputsString:   tx.OutputsString,
        TransactionHash: tx.TransactionHash,
    }
}

func ConvertHeaderToProto(header *Header) *pb.Header {
    return &pb.Header{
        Version:        int32(header.Version),
        HashPrevBlock:  header.hashPrevBlock,
        HashMerkleRoot: header.hashMerkleRoot,
        Timestamp:      header.Timestamp,
        Bits:           int32(header.Bits),
        Nonce:          header.Nonce,
    }
}

func ConvertBlockToProto(block *Block) *pb.Block {
    // Convert the Header
    protoHeader := ConvertHeaderToProto(block.BlockHeader)

    // Convert the Transactions
    protoTransactions := make([]*pb.Transaction, len(block.Transactions))
    for i, tx := range block.Transactions {
        protoTransactions[i] = ConvertTransactionToProto(tx)
    }

    return &pb.Block{
        MagicNumber:        int32(block.MagicNumber),
        Blocksize:          int32(block.Blocksize),
        BlockHeader:        protoHeader,
        TransactionCounter: int32(block.TransactionCounter),
        Transactions:       protoTransactions,
        Blockhash:          block.Blockhash,
    }
}

// Convert protobuf Output message to Go Output struct
func ConvertProtoToOutput(protoOutput *pb.Output) *Output {
    return &Output{
        Value:  int(protoOutput.Value),
        Index:  int(protoOutput.Index),
        Script: protoOutput.Script,
    }
}


// Convert protobuf Transaction message to Go Transaction struct
func ConvertProtoToTransaction(protoTx *pb.Transaction) *Transaction {
    outputs := make([]*Output, len(protoTx.ListOfOutputs))
    for i, protoOutput := range protoTx.ListOfOutputs {
        outputs[i] = ConvertProtoToOutput(protoOutput)
    }

    return &Transaction{
        VersionNumber:   int(protoTx.VersionNumber),
        InCounter:       int(protoTx.InCounter),
        ListOfInputs:    protoTx.ListOfInputs,
        OutCounter:      int(protoTx.OutCounter),
        ListOfOutputs:   outputs,
        OutputsString:   protoTx.OutputsString,
        TransactionHash: protoTx.TransactionHash,
    }
}

// Convert protobuf Header message to Go Header struct
func ConvertProtoToHeader(protoHeader *pb.Header) *Header {
    return &Header{
        Version:        int(protoHeader.Version),
        hashPrevBlock:  protoHeader.HashPrevBlock,
        hashMerkleRoot: protoHeader.HashMerkleRoot,
        Timestamp:      protoHeader.Timestamp,
        Bits:           int(protoHeader.Bits),
        Nonce:          protoHeader.Nonce,
    }
}

// Convert protobuf Block message to Go Block struct
func ConvertProtoToBlock(protoBlock *pb.Block) *Block {
    header := ConvertProtoToHeader(protoBlock.BlockHeader)

    transactions := make([]*Transaction, len(protoBlock.Transactions))
    for i, protoTx := range protoBlock.Transactions {
        transactions[i] = ConvertProtoToTransaction(protoTx)
    }

    return &Block{
        MagicNumber:        int(protoBlock.MagicNumber),
        Blocksize:          int(protoBlock.Blocksize),
        BlockHeader:        header,
        TransactionCounter: int(protoBlock.TransactionCounter),
        Transactions:       transactions,
        Blockhash:          protoBlock.Blockhash,
    }
}
