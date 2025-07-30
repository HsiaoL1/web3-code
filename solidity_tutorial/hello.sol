// SPDX-License-Identifier: MIT
// This is a simple Solidity contract that stores a greeting message.
pragma solidity ^0.8.0;

// contract 声明合约的意思
// 在 Solidity 中，合约类似于类的概念
// 合约可以包含状态变量、函数和事件等
// HelloWorld 合约是一个简单的合约示例 HelloWorld是合约的名称
// 合约的名称可以是任意的，但通常遵循 PascalCase 命名
// 状态变量是合约的存储数据
// 在这个例子中，greeting 是一个状态变量，用于存储问候语
// 构造函数在合约部署时执行一次
// 在这个例子中，构造函数将 greeting 初始化为 "Hello, World!"
contract HelloWorld {
    string public greeting;

    constructor() {
        greeting = "Hello, World!";
    }
    // 这是什么意思？
    // public 关键字表示该函数可以被外部调用
    // view 关键字表示该函数不会修改合约的状态
    // 返回值类型是 string memory，表示返回一个字符串
    // getGreeting 函数用于获取问候语
    // 这个函数可以被外部调用来获取当前的问候语
    // memory 关键字表示该字符串存储在内存中，而不是存储在区块链的存储中
    // Solidity 中的函数可以有不同的可见性修饰符
    function getGreeting() public view returns (string memory) {
        return greeting;
    }
}

// 如何运行
// 1. 安装 Solidity 编译器（solc）
// 2. 使用 solc 编译合约：solc --bin --abi hello.sol
// 3. 部署合约到以太坊网络（可以使用 Remix IDE 或 Truffle 等工具）
// 4. 调用 getGreeting 函数获取问候语
// 5. 可以使用 Web3.js 或 Ethers.js 等库与合约进行交互
// 6. 在 Remix IDE 中，可以直接编译和部署合约
// 7. 在 Remix 中，选择 Solidity 编译器版本 0.8.0 或更高版本
// 8. 点击 "Compile hello.sol" 按钮编译合约
// 9. 部署合约后，可以在 Remix 的 "Deployed Contracts" 部分找到合约实例
// 10. 点击 getGreeting 按钮调用 getGreeting 函数获取问候语
// 11. 可以在 Remix 的控制台中查看输出结果
// 12. 如果需要在本地环境中运行，可以使用 Truffle 或 Hardhat 等开发框架
// 13. 在 Truffle 中，创建一个新的项目并将 hello.sol 文件放入 contracts 目录

// 如何下载solidity编译器
// 1. 使用 npm 安装 solc：npm install -g solc`