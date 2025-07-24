# Solidity 编程语言入门教程

欢迎学习 Solidity！本教程旨在为初学者提供一个全面而详细的 Solidity 学习路径，从基础语法到核心概念，再到编写一个完整的智能合约。

**学习前建议**：
*   已经理解区块链基础概念（交易、区块、Gas等）。
*   有一定其他编程语言（如JavaScript, Python, C++）的基础会很有帮助。

---

## 目录

1.  [Solidity 简介](#1-solidity-简介)
2.  [开发环境准备](#2-开发环境准备)
3.  [第一个 Solidity 合约：HelloWorld](#3-第一个-solidity-合约helloworld)
4.  [合约结构详解](#4-合约结构详解)
5.  [数据类型 (Data Types)](#5-数据类型-data-types)
6.  [变量 (Variables)](#6-变量-variables)
7.  [函数 (Functions)](#7-函数-functions)
8.  [流程控制 (Control Flow)](#8-流程控制-control-flow)
9.  [映射 (Mappings)](#9-映射-mappings)
10. [结构体 (Structs)](#10-结构体-structs)
11. [修饰器 (Modifiers)](#11-修饰器-modifiers)
12. [事件 (Events)](#12-事件-events)
13. [错误处理 (Error Handling)](#13-错误处理-error-handling)
14. [继承 (Inheritance)](#14-继承-inheritance)
15. [综合示例：创建你自己的代币 (SimpleToken)](#15-综合示例创建你自己的代币-simpletoken)

---

## 1. Solidity 简介

Solidity 是一种面向合约的、为实现智能合约而创建的高级编程语言。它受到了 C++、Python 和 JavaScript 的影响，被设计用来在以太坊虚拟机（EVM）上运行。

**核心特点**：
*   **静态类型**：变量的类型在编译时就需要确定。
*   **面向合约**：支持继承、库和复杂的用户定义类型等特性。
*   **运行在 EVM 上**：编译后的代码是以太坊虚拟机可以理解的字节码。

---

## 2. 开发环境准备

对于初学者来说，最简单的入门方式是使用在线的 IDE，无需在本地安装任何东西。

*   **Remix IDE**: [https://remix.ethereum.org/](https://remix.ethereum.org/)

Remix 是一个强大、开源的工具，可以帮助你直接在浏览器中编写、编译、部署和调试 Solidity 代码。本教程的所有示例都可以在 Remix 中直接运行。

---

## 3. 第一个 Solidity 合约：HelloWorld

让我们从一个经典的 "Hello, World!" 程序开始。

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

// 定义一个名为 HelloWorld 的合约
contract HelloWorld {
    // 定义一个 string 类型的状态变量 greet
    string public greet = "Hello, World!";
}
```

**代码解释**：
1.  `// SPDX-License-Identifier: MIT`：这是代码的许可证标识符，建议所有 Solidity 源文件都加上。它告诉大家你的代码是基于 MIT 许可证开源的。
2.  `pragma solidity ^0.8.20;`：这是编译器版本指令。它告诉编译器，这份代码应该使用 `0.8.20` 或更高但低于 `0.9.0` 版本的编译器来编译。`^` 符号表示向上兼容。
3.  `contract HelloWorld { ... }`：`contract` 关键字用于声明一个合约。`HelloWorld` 是合约的名称，合约的所有代码都包含在花括号 `{}` 内。
4.  `string public greet = "Hello, World!";`：这行代码声明了一个名为 `greet` 的**状态变量**。
    *   `string` 是变量的类型。
    *   `public` 是一个可见性修饰符，它会自动为这个变量创建一个名为 `greet()` 的 `getter` 函数，任何人都可以调用这个函数来读取变量的值。
    *   `"Hello, World!"` 是赋给这个变量的初始值。

**如何在 Remix 中运行**：
1.  打开 Remix IDE。
2.  在 `contracts` 目录下创建一个新文件，命名为 `HelloWorld.sol`。
3.  将上面的代码复制进去。
4.  在左侧的 "Solidity Compiler" 标签页中，点击 "Compile HelloWorld.sol" 按钮。
5.  编译成功后，切换到 "Deploy & Run Transactions" 标签页。
6.  点击 "Deploy" 按钮。
7.  部署成功后，你会在下方 "Deployed Contracts" 区域看到你的 `HelloWorld` 合约。展开它，你会看到一个蓝色的 `greet` 按钮，点击它就会返回 "Hello, World!"。

---

## 4. 合约结构详解

一个典型的 Solidity 合约包含以下几个部分：

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract MyContract {
    // 1. 状态变量 (State Variables)
    uint256 public myNumber;
    address public owner;

    // 2. 构造函数 (Constructor)
    constructor(uint256 _initialNumber) {
        myNumber = _initialNumber;
        owner = msg.sender; // msg.sender 是部署合约的账户地址
    }

    // 3. 函数 (Functions)
    function setNumber(uint256 _newNumber) public {
        // 只有合约的拥有者才能调用这个函数
        require(msg.sender == owner, "You are not the owner.");
        myNumber = _newNumber;
    }

    // ... 其他部分，如修饰器、事件等
}
```

*   **状态变量**：永久存储在合约存储空间中的变量，它们的值会写入区块链。
*   **构造函数**：一个可选的、特殊的函数，名为 `constructor`。它在合约被创建时仅执行一次，通常用于初始化状态变量。
*   **函数**：合约的可执行代码单元，用于读取或修改状态变量。

---

## 5. 数据类型 (Data Types)

Solidity 是静态类型语言，主要数据类型包括：

*   **布尔型 (bool)**：`true` 或 `false`。
*   **整型 (int / uint)**：
    *   `int`：有符号整数（可以为负）。
    *   `uint`：无符号整数（只能为正或零）。
    *   可以指定位数，如 `uint8`, `uint16`, `uint256`。`uint` 和 `int` 默认是 `uint256` 和 `int256`。`uint256` 是最常用的类型。
*   **地址 (address)**：用于存储以太坊地址（20字节）。它有两种类型：
    *   `address`：普通地址。
    *   `address payable`：可支付地址，可以接收以太币，拥有 `transfer` 和 `send` 方法。
*   **字节数组 (bytes)**：
    *   `bytes1`, `bytes2`, ..., `bytes32`：定长字节数组。
    *   `bytes`：动态大小的字节数组。
    *   `string`：动态大小的 UTF-8 编码字符串，本质上是特殊的 `bytes`。
*   **枚举 (enum)**：用于创建用户定义类型，通常用于表示一组有限的状态。
    ```solidity
    enum Status { Pending, Shipped, Delivered }
    Status public myStatus = Status.Pending;
    ```

---

## 6. 变量 (Variables)

Solidity 中有三种类型的变量：

1.  **状态变量 (State Variables)**：
    *   永久存储在区块链上。
    *   在所有函数之外声明。
    *   消耗大量的 Gas。
    *   示例：`uint256 public myNumber;`

2.  **局部变量 (Local Variables)**：
    *   仅在函数执行期间存在。
    *   在函数内部声明。
    *   不消耗 Gas（除非是写入存储的复杂操作）。
    *   示例：
        ```solidity
        function myFunction() public pure returns (uint256) {
            uint256 i = 10; // 这是一个局部变量
            return i;
        }
        ```

3.  **全局变量 (Global Variables)**：
    *   存在于全局命名空间中，提供了关于区块链和交易的信息。
    *   常用全局变量：
        *   `msg.sender` (address): 当前调用发起方的地址。
        *   `msg.value` (uint): 随调用发送的以太币数量 (以 wei 为单位)。
        *   `block.timestamp` (uint): 当前区块的时间戳。
        *   `tx.origin` (address): 交易的原始发送方地址。

---

## 7. 函数 (Functions)

函数是合约的核心逻辑所在。

```solidity
// 语法结构
// function <函数名>(<参数类型> <参数名>) <可见性> <状态可变性> [returns (<返回类型>)] { ... }

contract FunctionExamples {
    uint256 public number;

    // 这是一个可以修改状态变量的函数
    function setNumber(uint256 _newNumber) public {
        number = _newNumber;
    }

    // 这是一个只读函数
    function getNumber() public view returns (uint256) {
        return number;
    }

    // 这是一个纯函数，既不读取也不修改状态
    function add(uint256 a, uint256 b) public pure returns (uint256) {
        return a + b;
    }
}
```

**可见性 (Visibility)**：
*   `public`：任何账户（外部或其他合约）都可以调用。
*   `private`：只能在当前合约内部调用。
*   `internal`：只能在当前合约及继承它的子合约中调用。
*   `external`：只能从合约外部调用（不能在合约内部调用，除非用 `this.func()` 的方式）。通常用于接收外部数据。

**状态可变性 (State Mutability)**：
*   **无修饰符** (默认)：可以修改状态。
*   `view`：只读函数，承诺不修改状态。调用 `view` 函数不需要花费 Gas（除非被另一个会修改状态的函数调用）。
*   `pure`：纯函数，承诺既不读取也不修改状态。例如，只对输入参数进行计算。
*   `payable`：可支付函数，允许在调用时接收以太币。

---

## 8. 流程控制 (Control Flow)

Solidity 支持常见的流程控制语句，与 JavaScript 或 C++ 非常相似。

```solidity
function example(uint256 _input) public pure returns (string memory) {
    // if-else
    if (_input == 0) {
        return "Input is zero";
    } else if (_input < 10) {
        return "Input is less than 10";
    } else {
        return "Input is 10 or greater";
    }

    // for 循环
    uint256 sum = 0;
    for (uint256 i = 1; i <= _input; i++) {
        sum += i;
    }
    // 注意：在实际合约中要非常小心使用循环，
    // 因为循环次数过多会导致 Gas 耗尽。
}
```

---

## 9. 映射 (Mappings)

映射是 Solidity 中极其重要和常用的数据结构，可以理解为哈希表或字典。

**语法**：`mapping(KeyType => ValueType) public myMapping;`

*   `KeyType`：可以是任何基本类型，如 `uint`, `address`。不能是映射、结构体或数组。
*   `ValueType`：可以是任何类型，包括另一个映射或结构体。

```solidity
contract Bank {
    // 记录每个地址的余额
    mapping(address => uint256) public balances;

    // 存款
    function deposit() public payable {
        require(msg.value > 0, "Deposit amount must be greater than 0");
        balances[msg.sender] += msg.value;
    }

    // 查看余额
    function getBalance(address _account) public view returns (uint256) {
        return balances[_account];
    }
}
```

**关键特性**：
*   所有可能的 `Key` 都已经“存在”，并被初始化为对应 `ValueType` 的默认值（如 `uint` 的 `0`）。
*   你不能遍历一个映射，也无法获取它的大小。这是设计上的取舍，为了保证 Gas 效率。如果需要遍历，通常会额外用一个数组来存储所有的 `Key`。

---

## 10. 结构体 (Structs)

结构体允许你创建更复杂的数据类型，将多个变量组合在一起。

```solidity
contract TodoList {
    // 定义一个 Todo 项的结构体
    struct Todo {
        string text;
        bool completed;
    }

    // 创建一个存储 Todo 的数组
    Todo[] public todos;

    function create(string memory _text) public {
        // 创建一个新的 Todo 实例并添加到数组中
        todos.push(Todo({
            text: _text,
            completed: false
        }));
    }

    function get(uint256 _index) public view returns (string memory, bool) {
        Todo storage todo = todos[_index];
        return (todo.text, todo.completed);
    }

    function toggleCompleted(uint256 _index) public {
        Todo storage todo = todos[_index];
        todo.completed = !todo.completed;
    }
}
```

---

## 11. 修饰器 (Modifiers)

修饰器是一种可重用的代码，用于在函数执行前检查某些条件。它们非常适合用于权限控制。

```solidity
contract OwnerContract {
    address public owner;

    constructor() {
        owner = msg.sender;
    }

    // 定义一个名为 onlyOwner 的修饰器
    modifier onlyOwner() {
        require(msg.sender == owner, "Caller is not the owner");
        _; // 这个下划线是占位符，代表被修饰的函数体将在此处执行
    }

    // 使用 onlyOwner 修饰器来保护这个函数
    function changeOwner(address _newOwner) public onlyOwner {
        owner = _newOwner;
    }
}
```

当 `changeOwner` 函数被调用时，会先执行 `onlyOwner` 修饰器中的 `require` 检查。如果检查通过，则执行 `_` 所在位置的函数体；如果失败，则交易会回滚。

---

## 12. 事件 (Events)

事件是 EVM 日志功能的便捷接口。它们用于通知外部世界（如你的 App 前端）合约中发生了什么事情。

当事件被触发时，它会将参数存储在交易的日志中。这些日志与合约地址关联，并永久保存在区块链上。这比直接读取合约状态要便宜得多。

```solidity
contract EventExample {
    // 定义一个事件
    // `indexed` 关键字可以让外部应用更容易地过滤和搜索这些事件
    event Transfer(address indexed from, address indexed to, uint256 amount);

    mapping(address => uint256) public balances;

    function transfer(address _to, uint256 _amount) public {
        require(balances[msg.sender] >= _amount, "Insufficient balance");
        balances[msg.sender] -= _amount;
        balances[_to] += _amount;

        // 触发事件
        emit Transfer(msg.sender, _to, _amount);
    }
}
```

---

## 13. 错误处理 (Error Handling)

Solidity 提供了三种主要的方式来处理错误和回滚交易：

1.  **`require(bool condition, string memory message)`**:
    *   用于检查前置条件，如函数输入或外部合约状态。
    *   如果 `condition` 为 `false`，交易会回滚。
    *   会退还剩余的 Gas。
    *   最常用的错误处理方式。

2.  **`assert(bool condition)`**:
    *   用于检查内部错误或不变量。
    *   如果 `condition` 为 `false`，交易会回滚。
    *   **会消耗掉所有 Gas**。
    *   通常用于检查代码中理论上不应该发生的错误。

3.  **`revert()`**:
    *   直接触发回滚。
    *   可以与 `if/else` 结合使用，提供更复杂的逻辑。
    *   `revert("Error message");` 等同于 `require(false, "Error message");`

---

## 14. 继承 (Inheritance)

Solidity 支持多重继承，允许你构建模块化、可重用的合约。

```solidity
// 父合约
contract Ownable {
    address public owner;

    constructor() {
        owner = msg.sender;
    }

    modifier onlyOwner() {
        require(msg.sender == owner, "Not the owner");
        _;
    }
}

// 子合约继承了 Ownable
contract MyToken is Ownable {
    // 现在 MyToken 合约自动拥有了 owner 状态变量和 onlyOwner 修饰器
    string public name = "My Token";

    function renameToken(string memory _newName) public onlyOwner {
        name = _newName;
    }
}
```

---

## 15. 综合示例：创建你自己的代币 (SimpleToken)

现在，让我们把上面学到的大部分概念结合起来，创建一个简单的、符合部分 ERC-20 标准的代币合约。

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

// 引入 OpenZeppelin 的 ERC20 标准接口，这是一个很好的实践
// 但为了教学，我们在这里手动实现一个简化版
// import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract SimpleToken {
    // --- 状态变量 ---

    string public name;
    string public symbol;
    uint8 public decimals = 18; // 18 是以太坊中常见的精度
    uint256 public totalSupply;

    // 映射：存储每个地址的余额
    mapping(address => uint256) public balanceOf;

    // 映射：存储授权额度 (owner => spender => amount)
    mapping(address => mapping(address => uint256)) public allowance;

    // --- 事件 ---

    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);

    // --- 构造函数 ---

    constructor(string memory _name, string memory _symbol, uint256 _initialSupply) {
        name = _name;
        symbol = _symbol;
        // 总供应量要乘以精度
        totalSupply = _initialSupply * (10**decimals);
        // 将初始供应量全部分配给合约的创建者
        balanceOf[msg.sender] = totalSupply;
    }

    // --- 函数 ---

    /**
     * @dev 将 `_value` 数量的代币发送到 `_to` 地址
     */
    function transfer(address _to, uint256 _value) public returns (bool success) {
        require(_to != address(0), "ERC20: transfer to the zero address");
        require(balanceOf[msg.sender] >= _value, "ERC20: transfer amount exceeds balance");

        balanceOf[msg.sender] -= _value;
        balanceOf[_to] += _value;

        emit Transfer(msg.sender, _to, _value);
        return true;
    }

    /**
     * @dev 允许 `_spender` 从你的账户中提取最多 `_value` 数量的代币
     */
    function approve(address _spender, uint256 _value) public returns (bool success) {
        require(_spender != address(0), "ERC20: approve to the zero address");

        allowance[msg.sender][_spender] = _value;

        emit Approval(msg.sender, _spender, _value);
        return true;
    }

    /**
     * @dev 从 `_from` 地址向 `_to` 地址发送 `_value` 数量的代币。
     * 这个函数需要 `_from` 地址预先通过 `approve` 函数授权。
     */
    function transferFrom(address _from, address _to, uint256 _value) public returns (bool success) {
        require(_from != address(0), "ERC20: transfer from the zero address");
        require(_to != address(0), "ERC20: transfer to the zero address");
        require(balanceOf[_from] >= _value, "ERC20: transfer amount exceeds balance");
        require(allowance[_from][msg.sender] >= _value, "ERC20: transfer amount exceeds allowance");

        balanceOf[_from] -= _value;
        balanceOf[_to] += _value;
        allowance[_from][msg.sender] -= _value;

        emit Transfer(_from, _to, _value);
        return true;
    }
}
```

**如何使用这个合约**：
1.  在 Remix 中编译并部署，部署时需要提供 `_name` (例如 "My Simple Token"), `_symbol` (例如 "MST"), 和 `_initialSupply` (例如 10000)。
2.  部署后，你可以：
    *   在 `balanceOf` 中输入你自己的地址，查询你的余额，应该等于 `10000 * 10**18`。
    *   使用 `transfer` 函数向另一个地址发送一些代币。
    *   使用 `approve` 授权另一个地址可以从你的账户中花费一定额度的代币。
    *   让被授权的地址调用 `transferFrom` 来完成转账。

---

恭喜你！你已经完成了 Solidity 的入门学习。接下来，你可以深入研究更高级的主题，例如：
*   **OpenZeppelin 合约库**：学习如何使用经过安全审计的标准合约（如 ERC20, ERC721）来构建应用。
*   **DeFi (去中心化金融)**：了解 Uniswap, Aave 等协议的原理。
*   **Gas 优化**：学习如何编写更省 Gas 的代码。
*   **安全实践**：学习防止重入攻击等常见的安全漏洞。

继续探索，区块链的世界充满了无限可能！
