// SPDX-License-Identifier: MIT
// 这是一个简单的 ERC20 代币合约示例
pragma solidity ^0.8.20;


// SimpleToken 继承自 OpenZeppelin 的 ERC20 合约
// 这意味着 SimpleToken 将具有 ERC20 代币的所有标准功能
contract SimpleToken {
    // -- 状态变量
    // name 是代币的名称
    string public name;
    // symbol 是代币的符号
    string public symbol;
    // decimals 是代币的小数位数
    uint8 public decimals = 18; // 默认值为 18
    // totalSupply 是代币的总供应量
    uint256 public totalSupply;
    // balances 映射存储每个地址的代币余额
    mapping(address => uint256) public balances;

    // allowances 映射存储每个地址对其他地址的代币授权
    mapping(address => mapping(address => uint256)) public allowances;

    // -- 事件
    // Transfer 事件在代币转移时触发
    event Transfer(address indexed from, address indexed to, uint256 value);
    // Approval 事件在代币授权时触发
    event Approval(address indexed owner, address indexed spender, uint256 value);

    // -- 构造函数
    // 构造函数在合约部署时执行一次
    constructor(string memory _name, string memory _symbol, uint256 _initialSupply) {
        name = _name;
        symbol = _symbol;
        totalSupply = _initialSupply * (10 ** uint256(decimals)); // 考虑小数位
        balances[msg.sender] = totalSupply; // 将初始供应量分配给合约部署者
    }

    // -- 函数
    // transfer 函数用于将代币从一个地址转移到另一个地址
    function transfer(address to, uint256 value) public returns (bool) {
        require(to != address(0), "Invalid address");
        require(balances[msg.sender] >= value, "Insufficient balance");

        balances[msg.sender] -= value;
        balances[to] += value;

        // emit 是什么意思？
        // emit 关键字用于触发事件
        // 事件可以在区块链上记录重要的状态变化
        emit Transfer(msg.sender, to, value);
        return true;
    }

    // approval 函数用于授权其他地址可以花费一定数量的代币
    function approve(address spender, uint256 value) public returns (bool) {
        require(spender != address(0), "Invalid address");

        allowances[msg.sender][spender] = value;

        emit Approval(msg.sender, spender, value);
        return true;
    }

    // transferFrom 函数用于从一个地址转移代币到另一个地址
    function transferFrom(address from, address to, uint256 value) public returns (bool) {
        require(from != address(0), "Invalid from address");
        require(to != address(0), "Invalid to address");
        require(balances[from] >= value, "Insufficient balance");
        require(allowances[from][msg.sender] >= value, "Allowance exceeded");

        balances[from] -= value;
        balances[to] += value;
        allowances[from][msg.sender] -= value;

        emit Transfer(from, to, value);
        return true;
    }
}

