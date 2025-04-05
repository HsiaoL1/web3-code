package chain

import (
	"errors"
	"fmt"

	"github.com/hsiaocz/web3-code/blockchain/block"
	"github.com/hsiaocz/web3-code/blockchain/pow"
)

// MineDifficulty 挖矿难度
const MineDifficulty = 12

// Blockchain 表示区块链结构
type Blockchain struct {
	Blocks []*block.Block
}

// NewBlockchain 创建一个新的区块链
func NewBlockchain() *Blockchain {
	genesisBlock := block.NewGenesisBlock(MineDifficulty)
	// 为创世区块计算工作量证明
	powInstance := pow.NewProofOfWork(genesisBlock)
	nonce, hash := powInstance.Run()
	genesisBlock.Nonce = nonce
	genesisBlock.Hash = hash

	bc := &Blockchain{[]*block.Block{genesisBlock}}
	return bc
}

// AddBlock 向区块链中添加新区块
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := block.NewBlock(data, prevBlock.Hash, MineDifficulty)
	
	// 计算工作量证明
	powInstance := pow.NewProofOfWork(newBlock)
	nonce, hash := powInstance.Run()
	newBlock.Nonce = nonce
	newBlock.Hash = hash
	
	bc.Blocks = append(bc.Blocks, newBlock)
}

// GetBlock 根据索引获取区块
func (bc *Blockchain) GetBlock(index int) (*block.Block, error) {
	if index < 0 || index >= len(bc.Blocks) {
		return nil, errors.New("block index out of range")
	}
	return bc.Blocks[index], nil
}

// GetLatestBlock 获取最新的区块
func (bc *Blockchain) GetLatestBlock() *block.Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

// IsValid 验证区块链是否有效
func (bc *Blockchain) IsValid() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		previousBlock := bc.Blocks[i-1]

		// 验证哈希
		powInstance := pow.NewProofOfWork(currentBlock)
		if !powInstance.Validate() {
			return false
		}

		// 验证当前区块的前一个区块哈希是否指向了前一个区块的哈希
		if string(currentBlock.PrevBlockHash) != string(previousBlock.Hash) {
			return false
		}
	}
	return true
}

// PrintChain 打印区块链信息
func (bc *Blockchain) PrintChain() {
	for i, block := range bc.Blocks {
		fmt.Printf("Block %d:\n", i)
		fmt.Printf("  PrevHash: %x\n", block.PrevBlockHash)
		fmt.Printf("  Data: %s\n", block.Data)
		fmt.Printf("  Hash: %x\n", block.Hash)
		
		// 验证区块的工作量证明
		pow := pow.NewProofOfWork(block)
		fmt.Printf("  PoW: %v\n\n", pow.Validate())
	}
} 