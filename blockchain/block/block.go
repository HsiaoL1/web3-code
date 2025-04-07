package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"strconv"
	"time"
)

// Block 表示区块链中的一个区块
type Block struct {
	Timestamp     int64  // 时间戳
	Data          []byte // 区块数据
	PrevBlockHash []byte // 前一个区块的哈希
	Hash          []byte // 当前区块的哈希
	Nonce         int    // 工作量证明计数器
	Difficulty    int    // 挖矿难度
}

// NewBlock 创建并返回新的区块
func NewBlock(data string, prevBlockHash []byte, difficulty int) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
		Difficulty:    difficulty,
	}

	// 通过调用工作量证明来挖矿，设置nonce和hash
	// 这里先返回创建的区块，挖矿过程将在blockchain包中进行
	return block
}

// SetHash 计算并设置区块哈希
func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{
		b.PrevBlockHash,
		b.Data,
		timestamp,
		[]byte(strconv.FormatInt(int64(b.Nonce), 10)),
		[]byte(strconv.Itoa(b.Difficulty)),
	}, []byte{})

	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

// Serialize 将区块序列化为字节数组
func (b *Block) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

// DeserializeBlock 从字节数组反序列化为区块
func DeserializeBlock(d []byte) (*Block, error) {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		return nil, err
	}

	return &block, nil
}

// NewGenesisBlock 创建创世区块
func NewGenesisBlock(difficulty int) *Block {
	return NewBlock("Genesis Block", []byte{}, difficulty)
} 
