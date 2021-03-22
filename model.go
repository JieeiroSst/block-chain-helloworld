package main

import "math/big"

type Block struct {
	Hash 		[]byte  `json:"hash"`
	PrevHash 	[]byte	`json:"prev_hash"`
	Data 		[]byte  `json:"data"`
	Nonce    	int		`json:"nonce"`
}

type BlockChain struct {
	Blocks []*Block
}

type ProofOfWord struct {
	Block *Block
	Target *big.Int
}