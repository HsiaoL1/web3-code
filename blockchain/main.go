package main

import (
	"flag"
	"fmt"

	"github.com/hsiaocz/web3-code/blockchain/api"
	"github.com/hsiaocz/web3-code/blockchain/chain"
)

func main() {
	// 命令行参数解析
	cliMode := flag.Bool("cli", false, "Run in CLI mode instead of server mode")
	port := flag.String("port", "8080", "Port to run the HTTP server on")
	flag.Parse()

	if *cliMode {
		// CLI模式：创建区块链并进行简单交互
		runCLI()
	} else {
		// 服务器模式：启动HTTP API服务器
		server := api.NewBlockchainServer()
		server.Start(*port)
	}
}

func runCLI() {
	bc := chain.NewBlockchain()
	
	// 添加一些测试区块
	bc.AddBlock("Send 1 BTC to Alice")
	bc.AddBlock("Send 2 BTC to Bob")
	bc.AddBlock("Send 0.5 BTC to Charlie")
	
	// 打印区块链信息
	fmt.Println("\n区块链信息:")
	bc.PrintChain()
	
	// 验证区块链
	fmt.Printf("\n区块链有效性: %v\n", bc.IsValid())
} 