package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

type Block struct {
	Hash 		[]byte  `json:"hash"`
	PrevHash 	[]byte	`json:"prev_hash"`
	Data 		[]byte  `json:"data"`
}

type BlockChain struct {
	Blocks []*Block
}

func (b *Block) DeriveHash(){
	info:=bytes.Join([][]byte{b.Data,b.PrevHash},[]byte{})
	fmt.Println(">>>>>>>>>>",&info)
	hash:=sha256.Sum256(info)
	b.Hash=hash[:]
}

func CreateBlock(data string,prevHash []byte) *Block {
	block:=&Block{
		Hash:     []byte{},
		PrevHash: prevHash,
		Data:     []byte(data),
	}
	block.DeriveHash()
	return block
}

func (chain *BlockChain) AddBlock(data string){
	prevBlock:=chain.Blocks[len(chain.Blocks)-1]
	block :=CreateBlock(data,prevBlock.PrevHash)
	chain.Blocks=append(chain.Blocks, block)
}

func Genesis() *Block{
	return CreateBlock("genesis",[]byte{})
}

func InitBlockChain() *BlockChain {
	return &BlockChain{Blocks:[]*Block{Genesis()}}
}

func main(){
	chain:=InitBlockChain()

	chain.AddBlock("first block after genesis")
	//chain.AddBlock("second block after genesis")
	//chain.AddBlock("third block after genesis")

	for _,block:=range chain.Blocks{
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("data: %s\n", block.Data)
		fmt.Printf("hash: %x\n", block.Hash)
		fmt.Println("\n")
	}
}