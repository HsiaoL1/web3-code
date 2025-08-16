# 第一章：Solidity 基础入门

## 本章学习目标

- 了解 Solidity 语言的特点和用途
- 搭建 Solidity 开发环境
- 编写第一个智能合约
- 理解智能合约的基本结构
- 掌握合约的编译和部署流程

## 1.1 什么是 Solidity？

Solidity 是一种面向合约的、为实现智能合约而创建的高级编程语言。它由以太坊团队开发，专门用于在以太坊虚拟机（EVM）上运行。

### 核心特点

- **静态类型语言**：变量类型在编译时确定，有助于提前发现错误
- **面向合约**：支持继承、库和复杂的用户定义类型
- **图灵完备**：可以实现任何计算逻辑（受 Gas 限制）
- **EVM 兼容**：编译为以太坊虚拟机字节码
- **语法类似 JavaScript**：易于 Web 开发者上手

### 语言影响

Solidity 受到以下语言的影响：

- **C++**：语法结构和类型系统
- **Python**：代码风格和可读性
- **JavaScript**：函数和对象的概念

### 应用场景

- **智能合约开发**：去中心化应用的核心逻辑
- **DeFi 协议**：去中心化金融产品
- **NFT 项目**：数字资产和收藏品
- **DAO 治理**：去中心化自治组织
- **供应链管理**：透明的物流追踪

## 1.2 开发环境搭建

### 1.2.1 在线开发环境 - Remix IDE

对于初学者，我们强烈推荐使用 Remix IDE，这是一个功能强大的在线开发环境。

**Remix IDE 地址**：https://remix.ethereum.org/

**Remix 的优势**：

- 无需安装，直接在浏览器中使用
- 集成编译器、调试器和部署工具
- 内置多个 Solidity 版本
- 支持多种测试网络
- 丰富的插件生态

**Remix 界面介绍**：

1. **文件浏览器**：管理合约文件
2. **编辑器**：编写和编辑代码
3. **Solidity 编译器**：编译合约
4. **部署和运行**：部署和测试合约
5. **调试器**：调试合约执行
6. **终端**：查看输出和日志

### 1.2.2 本地开发环境（可选）

如果您想要搭建本地开发环境，可以选择以下工具：

**Hardhat（推荐）**：

```bash
npm install --save-dev hardhat
npx hardhat
```

**Truffle**：

```bash
npm install -g truffle
truffle init
```

**Foundry**：

```bash
curl -L https://foundry.paradigm.xyz | bash
foundryup
```

### 1.2.3 准备工作

1. **安装浏览器钱包**（推荐 MetaMask）
2. **获取测试 ETH**：
   - Sepolia 测试网：https://sepoliafaucet.com/
   - Goerli 测试网：https://goerlifaucet.com/
3. **加入开发者社区**：
   - Ethereum StackExchange
   - Discord/Telegram 技术群组

## 1.3 第一个智能合约：Hello World

让我们从经典的 "Hello World" 程序开始 Solidity 之旅。

### 1.3.1 创建合约文件

在 Remix 中创建一个新文件 `HelloWorld.sol`：

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title HelloWorld
 * @dev 最简单的智能合约示例
 */
contract HelloWorld {
    // 状态变量：存储问候语
    string public greeting;

    // 构造函数：合约部署时执行
    constructor() {
        greeting = "Hello, Solidity World!";
    }

    // 公共函数：获取问候语
    function getGreeting() public view returns (string memory) {
        return greeting;
    }

    // 公共函数：设置新的问候语
    function setGreeting(string memory _newGreeting) public {
        greeting = _newGreeting;
    }
}
```

### 1.3.2 代码详解

**1. 许可证标识符**

```solidity
// SPDX-License-Identifier: MIT
```

- 用于指定代码的开源许可证
- MIT 是最常用的许可证类型
- 这行注释必须在文件开头

**2. 编译器版本指令**

```solidity
pragma solidity ^0.8.20;
```

- 指定 Solidity 编译器版本
- `^0.8.20` 表示兼容 0.8.20 及以上版本（但低于 0.9.0）
- 确保代码在指定版本上正常工作

**3. 合约声明**

```solidity
contract HelloWorld {
    // 合约内容
}
```

- `contract` 关键字定义智能合约
- `HelloWorld` 是合约名称（应遵循 PascalCase 命名规范）

**4. 状态变量**

```solidity
string public greeting;
```

- `string` 是数据类型，用于存储文本
- `public` 是可见性修饰符，自动生成 getter 函数
- `greeting` 是变量名

**5. 构造函数**

```solidity
constructor() {
    greeting = "Hello, Solidity World!";
}
```

- 合约部署时只执行一次
- 用于初始化状态变量
- 不能被显式调用

**6. 函数定义**

```solidity
function getGreeting() public view returns (string memory) {
    return greeting;
}
```

- `function` 关键字定义函数
- `public` 表示可以从外部调用
- `view` 表示只读函数，不修改状态
- `returns (string memory)` 指定返回类型

### 1.3.3 编译和部署

**步骤 1：编译合约**

1. 在 Remix 中选择 "Solidity Compiler" 标签
2. 选择合适的编译器版本（0.8.20 或更高）
3. 点击 "Compile HelloWorld.sol"
4. 确认编译成功，无错误和警告

**步骤 2：部署合约**

1. 切换到 "Deploy & Run Transactions" 标签
2. 选择环境（推荐使用 "Remix VM" 进行测试）
3. 确认合约选择为 "HelloWorld"
4. 点击 "Deploy" 按钮
5. 在底部确认部署成功

**步骤 3：与合约交互**

1. 在 "Deployed Contracts" 区域找到您的合约
2. 点击蓝色的 `greeting` 按钮查看当前问候语
3. 在 `setGreeting` 输入框中输入新的问候语
4. 点击橙色的 `setGreeting` 按钮执行交易
5. 再次点击 `greeting` 确认更新成功

## 1.4 智能合约的基本结构

一个完整的 Solidity 合约通常包含以下组件：

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

// 导入其他合约或库
import "./OtherContract.sol";

/**
 * @title ExampleContract
 * @dev 智能合约基本结构示例
 */
contract ExampleContract {
    // 1. 状态变量声明
    uint256 public totalSupply;
    address public owner;
    bool public isActive;

    // 2. 事件声明
    event Transfer(address indexed from, address indexed to, uint256 value);
    event OwnershipTransferred(address indexed previousOwner, address indexed newOwner);

    // 3. 修饰器声明
    modifier onlyOwner() {
        require(msg.sender == owner, "Not the contract owner");
        _;
    }

    modifier whenActive() {
        require(isActive, "Contract is not active");
        _;
    }

    // 4. 构造函数
    constructor(uint256 _initialSupply) {
        totalSupply = _initialSupply;
        owner = msg.sender;
        isActive = true;
    }

    // 5. 公共函数
    function transferOwnership(address _newOwner) public onlyOwner {
        require(_newOwner != address(0), "New owner cannot be zero address");
        emit OwnershipTransferred(owner, _newOwner);
        owner = _newOwner;
    }

    function toggleActive() public onlyOwner {
        isActive = !isActive;
    }

    // 6. 内部函数
    function _internalFunction() internal pure returns (bool) {
        return true;
    }

    // 7. 私有函数
    function _privateFunction() private pure returns (bool) {
        return true;
    }
}
```

### 结构组件说明

1. **状态变量**：永久存储在区块链上的数据
2. **事件**：记录合约状态变化的日志
3. **修饰器**：可重用的函数执行条件检查
4. **构造函数**：合约部署时的初始化逻辑
5. **函数**：合约的主要功能实现

## 1.5 实战练习

### 练习 1：个人信息合约

创建一个存储个人信息的智能合约：

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract PersonalInfo {
    // TODO: 定义状态变量存储姓名、年龄、邮箱

    // TODO: 实现构造函数

    // TODO: 实现获取信息的函数

    // TODO: 实现更新信息的函数
}
```

**要求**：

1. 存储姓名、年龄和邮箱地址
2. 只有合约部署者可以更新信息
3. 任何人都可以查看信息

### 练习 2：简单计数器

实现一个可以增加、减少和重置的计数器：

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract Counter {
    // TODO: 实现计数器功能
}
```

**要求**：

1. 初始计数为 0
2. 提供增加、减少、重置功能
3. 添加事件记录每次操作

## 1.6 章节总结

在本章中，我们学习了：

1. **Solidity 基础概念**：了解了 Solidity 的特点和应用场景
2. **开发环境搭建**：学会使用 Remix IDE 进行开发
3. **第一个合约**：编写、编译和部署了 Hello World 合约
4. **合约结构**：理解了智能合约的基本组成部分

### 关键知识点

- Solidity 是静态类型的面向合约编程语言
- 合约包含状态变量、函数、事件等组件
- `public` 可见性自动生成 getter 函数
- 构造函数在合约部署时执行一次
- Remix IDE 是优秀的在线开发环境

### 下一章预告

在下一章《数据类型与变量》中，我们将深入学习：

- Solidity 的各种数据类型
- 状态变量与局部变量的区别
- 可见性修饰符的详细用法
- 全局变量的使用

继续您的 Solidity 学习之旅：[第二章 - 数据类型与变量](02_data_types_variables.md)
