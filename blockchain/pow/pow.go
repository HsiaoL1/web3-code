package pow

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"

	"github.com/hsiaocz/web3-code/blockchain/block"
)

// ProofOfWork 表示工作量证明结构
type ProofOfWork struct {
	block      *block.Block
	target     *big.Int
}

// NewProofOfWork 创建新的工作量证明实例
func NewProofOfWork(b *block.Block) *ProofOfWork {
	target := big.NewInt(1)
	// 左移 (256 - difficulty) 位
	target.Lsh(target, uint(256-b.Difficulty))

	pow := &ProofOfWork{b, target}
	return pow
}

// prepareData 准备工作量证明的数据
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			intToHex(pow.block.Timestamp),
			intToHex(int64(pow.block.Difficulty)),
			intToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

// Run 执行工作量证明算法，挖矿
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	maxNonce := math.MaxInt64

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			fmt.Printf("\r%x", hash)
			break
		} else {
			nonce++
		}
	}
	fmt.Printf("\n\n")

	return nonce, hash[:]
}

// Validate 验证区块是否有效
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1
	return isValid
}

// 辅助函数：将int64转为字节数组
func intToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		return []byte{}
	}
	return buff.Bytes()
} 