# 简易区块链项目

这是一个使用Go语言实现的简单区块链项目，旨在演示区块链的基本原理。

## 功能

- 基本的区块结构和区块链实现
- 工作量证明（Proof of Work）挖矿算法
- 区块链验证机制
- HTTP API 接口
- 命令行界面

## 项目结构

```
blockchain/
  ├── block/       # 区块结构
  ├── chain/       # 区块链结构
  ├── pow/         # 工作量证明算法
  ├── api/         # HTTP API接口
  ├── utils/       # 工具函数
  └── main.go      # 主程序
```

## 运行方式

### 构建项目

```bash
go build -o blockchain blockchain/main.go
```

### 命令行模式

```bash
./blockchain --cli
```

### HTTP服务器模式

```bash
./blockchain --port 8080
```

## API 接口

### 获取区块链信息

```
GET /chain
```

### 获取特定区块

```
GET /block?index=0
```

### 添加新区块（挖矿）

```
POST /mine
Content-Type: application/json

{
  "data": "交易数据"
}
```

## 概念解释

### 区块（Block）

区块是区块链的基本单位，包含：
- 时间戳
- 数据（交易记录）
- 前一个区块的哈希
- 自身的哈希
- 工作量证明相关字段（nonce和难度）

### 工作量证明（Proof of Work）

工作量证明是一种共识算法，要求节点（矿工）证明他们已经完成了一定量的计算工作。在本项目中，工作量证明算法要求计算一个使得区块哈希满足特定难度的nonce值。

### 区块链（Blockchain）

区块链是一系列通过哈希链接在一起的区块，构成一个不可篡改的分布式账本。
