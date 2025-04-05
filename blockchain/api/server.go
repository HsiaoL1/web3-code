package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hsiaocz/web3-code/blockchain/chain"
)

// BlockchainServer 区块链API服务器
type BlockchainServer struct {
	blockchain *chain.Blockchain
}

// BlockData 用于添加区块的请求数据结构
type BlockData struct {
	Data string `json:"data"`
}

// BlockResponse 区块查询响应结构
type BlockResponse struct {
	Index        int    `json:"index"`
	Timestamp    int64  `json:"timestamp"`
	Data         string `json:"data"`
	Hash         string `json:"hash"`
	PrevHash     string `json:"prev_hash"`
	Difficulty   int    `json:"difficulty"`
	Nonce        int    `json:"nonce"`
	IsValid      bool   `json:"is_valid"`
}

// ChainResponse 区块链响应结构
type ChainResponse struct {
	Length int             `json:"length"`
	Blocks []BlockResponse `json:"blocks"`
	IsValid bool           `json:"is_valid"`
}

// NewBlockchainServer 创建一个新的区块链服务器
func NewBlockchainServer() *BlockchainServer {
	return &BlockchainServer{
		blockchain: chain.NewBlockchain(),
	}
}

// Start 启动HTTP服务器
func (server *BlockchainServer) Start(port string) {
	// 获取区块链信息
	http.HandleFunc("/chain", server.handleGetChain)
	
	// 获取特定区块信息
	http.HandleFunc("/block", server.handleGetBlock)
	
	// 添加新区块
	http.HandleFunc("/mine", server.handleMineBlock)
	
	log.Printf("Starting blockchain server at :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// handleGetChain 处理获取整个区块链的请求
func (server *BlockchainServer) handleGetChain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var blocks []BlockResponse
	for i, block := range server.blockchain.Blocks {
		pow := server.blockchain.IsValid()
		blocks = append(blocks, BlockResponse{
			Index:      i,
			Timestamp:  block.Timestamp,
			Data:       string(block.Data),
			Hash:       string(block.Hash),
			PrevHash:   string(block.PrevBlockHash),
			Difficulty: block.Difficulty,
			Nonce:      block.Nonce,
			IsValid:    pow,
		})
	}

	response := ChainResponse{
		Length:  len(server.blockchain.Blocks),
		Blocks:  blocks,
		IsValid: server.blockchain.IsValid(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleGetBlock 处理获取特定区块的请求
func (server *BlockchainServer) handleGetBlock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	indexStr := r.URL.Query().Get("index")
	if indexStr == "" {
		http.Error(w, "Missing index parameter", http.StatusBadRequest)
		return
	}

	var index int
	_, err := fmt.Sscanf(indexStr, "%d", &index)
	if err != nil {
		http.Error(w, "Invalid index parameter", http.StatusBadRequest)
		return
	}

	block, err := server.blockchain.GetBlock(index)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := BlockResponse{
		Index:      index,
		Timestamp:  block.Timestamp,
		Data:       string(block.Data),
		Hash:       string(block.Hash),
		PrevHash:   string(block.PrevBlockHash),
		Difficulty: block.Difficulty,
		Nonce:      block.Nonce,
		IsValid:    true, // 此区块已经在区块链中，默认为有效
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleMineBlock 处理挖矿（添加新区块）的请求
func (server *BlockchainServer) handleMineBlock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var blockData BlockData
	err := json.NewDecoder(r.Body).Decode(&blockData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if blockData.Data == "" {
		http.Error(w, "Data field is required", http.StatusBadRequest)
		return
	}

	// 添加新区块
	server.blockchain.AddBlock(blockData.Data)
	
	// 获取最新添加的区块
	newBlock := server.blockchain.GetLatestBlock()
	index := len(server.blockchain.Blocks) - 1
	
	response := BlockResponse{
		Index:      index,
		Timestamp:  newBlock.Timestamp,
		Data:       string(newBlock.Data),
		Hash:       string(newBlock.Hash),
		PrevHash:   string(newBlock.PrevBlockHash),
		Difficulty: newBlock.Difficulty,
		Nonce:      newBlock.Nonce,
		IsValid:    true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
} 