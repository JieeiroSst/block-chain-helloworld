package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

var (
	Difficulty=25
)


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

func (b *Block) DeriveHash(){
	info:=bytes.Join([][]byte{b.Data,b.PrevHash},[]byte{})
	hash:=sha256.Sum256(info)
	b.Hash=hash[:]
}

func CreateBlock(data string,prevHash []byte) *Block {
	block:=&Block{
		Hash:     []byte{},
		PrevHash: prevHash,
		Data:     []byte(data),
		Nonce:    0,
	}
	pow:=NewProofOfWord(block)
	nonce,hash:=pow.Run()
	block.Nonce=nonce
	block.Hash=hash[:]

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

func NewProofOfWord(b *Block) *ProofOfWord {
	target:=big.NewInt(1)
	target.Lsh(target,uint(256-Difficulty))
	pow:=&ProofOfWord{b,target}
	return pow
}

func ToHex(num int64) []byte{
	buff:=new(bytes.Buffer)
	err:=binary.Write(buff,binary.BigEndian,num)
	if err!=nil {
		panic(err)
	}
	return buff.Bytes()
}

func (pow *ProofOfWord) InitNonce(nonce int) []byte{
	data:=bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},[]byte{},
		)
	fmt.Println()
	return data
}

func (pow *ProofOfWord) Run() (int,[]byte) {
	var intHash big.Int
	var hash [32]byte
	nonce:=0
	for nonce < math.MaxInt64 {
		data:=pow.InitNonce(nonce)
		hash=sha256.Sum256(data)
		fmt.Printf("\r %X",hash)
		intHash.SetBytes(hash[:])
		if intHash.Cmp(pow.Target)== -1{
			break
		}else{
			nonce++
		}
	}
	return nonce,hash[:]
}

func (pow *ProofOfWord) Validate() bool{
	var intHash big.Int
	hash:=pow.InitNonce(pow.Block.Nonce)
	intHash.SetBytes(hash[:])
	return intHash.Cmp(pow.Target) == -1
}

func main(){
	chain:=InitBlockChain()

	chain.AddBlock("first block after genesis")
	chain.AddBlock("second block after genesis")
	chain.AddBlock("third block after genesis")

	for _,block:=range chain.Blocks{
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("data: %s\n", block.Data)
		fmt.Printf("hash: %x\n", block.Hash)
		fmt.Println("\n")
	}
}