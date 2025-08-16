# 第二章：数据类型与变量

## 本章学习目标

- 掌握 Solidity 的基本数据类型
- 理解状态变量、局部变量和全局变量的区别
- 学会使用可见性修饰符
- 了解变量的存储位置和生命周期
- 掌握常量和不可变变量的使用

## 2.1 基本数据类型

Solidity 作为静态类型语言，提供了丰富的基本数据类型。每种类型都有其特定的用途和存储要求。

### 2.1.1 布尔类型 (bool)

布尔类型只有两个值：`true` 和 `false`。

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract BooleanExample {
    bool public isActive = true;
    bool public isCompleted = false;

    function toggleActive() public {
        isActive = !isActive;  // 取反操作
    }

    function checkCondition(uint256 _value) public pure returns (bool) {
        return _value > 100;  // 比较操作返回布尔值
    }

    function logicalOperations(bool a, bool b) public pure returns (bool, bool, bool) {
        return (
            a && b,  // 逻辑与
            a || b,  // 逻辑或
            a != b   // 不等于
        );
    }
}
```

**布尔运算符**：

- `!` (逻辑非)
- `&&` (逻辑与)
- `||` (逻辑或)
- `==` (等于)
- `!=` (不等于)

### 2.1.2 整数类型

Solidity 提供有符号和无符号整数类型，支持从 8 位到 256 位。

```solidity
contract IntegerExample {
    // 无符号整数（只能为正数或零）
    uint8 public smallNumber = 255;      // 0 到 255
    uint16 public mediumNumber = 65535;  // 0 到 65,535
    uint256 public largeNumber;          // 0 到 2^256 - 1
    uint public defaultUint;             // 默认是 uint256

    // 有符号整数（可以为负数）
    int8 public signedSmall = -128;      // -128 到 127
    int256 public signedLarge;           // -2^255 到 2^255 - 1
    int public defaultInt;               // 默认是 int256

    function integerOperations(uint256 a, uint256 b) public pure returns (uint256, uint256, uint256, uint256) {
        return (
            a + b,  // 加法
            a - b,  // 减法（注意：无符号整数不能为负）
            a * b,  // 乘法
            a / b   // 除法（整数除法，向下取整）
        );
    }

    function comparisonOperations(uint256 a, uint256 b) public pure returns (bool, bool, bool, bool) {
        return (
            a > b,   // 大于
            a < b,   // 小于
            a >= b,  // 大于等于
            a <= b   // 小于等于
        );
    }
}
```

**重要注意事项**：

- Solidity 0.8.0+ 默认启用溢出检查
- 无符号整数减法结果不能为负数
- 除零会导致交易回滚

### 2.1.3 地址类型 (address)

地址类型专门用于存储以太坊地址（20 字节）。

```solidity
contract AddressExample {
    address public owner;
    address payable public wallet;  // 可接收以太币的地址

    constructor() {
        owner = msg.sender;  // 设置合约部署者为所有者
        wallet = payable(msg.sender);  // 转换为 payable 地址
    }

    function getAddressInfo(address _addr) public view returns (uint256, bytes32) {
        return (
            _addr.balance,  // 地址的以太币余额（wei 单位）
            keccak256(abi.encodePacked(_addr))  // 地址的哈希值
        );
    }

    function sendEther(address payable _to) public payable {
        require(msg.value > 0, "Must send some ether");
        _to.transfer(msg.value);  // 发送以太币
    }

    function checkCode(address _addr) public view returns (bool) {
        uint256 size;
        assembly {
            size := extcodesize(_addr)
        }
        return size > 0;  // 检查是否为合约地址
    }
}
```

**地址类型成员**：

- `balance`：地址的以太币余额
- `transfer(uint256)`：发送以太币（失败时回滚）
- `send(uint256)`：发送以太币（返回布尔值）
- `call(bytes)`：低级调用函数

### 2.1.4 字节类型

Solidity 提供定长和动长字节数组。

```solidity
contract BytesExample {
    bytes1 public singleByte = 0xFF;
    bytes32 public hash;
    bytes public dynamicBytes;

    function bytesOperations() public {
        // 定长字节数组
        bytes32 data = keccak256("Hello, World!");
        hash = data;

        // 动态字节数组
        dynamicBytes = "Hello, Solidity!";
    }

    function getBytesInfo() public view returns (uint256, bytes1) {
        return (
            dynamicBytes.length,    // 动态字节数组长度
            hash[0]                 // 定长字节数组的第一个字节
        );
    }

    function compareBytes(bytes32 a, bytes32 b) public pure returns (bool) {
        return a == b;  // 字节数组比较
    }
}
```

### 2.1.5 字符串类型 (string)

字符串本质上是动态字节数组，用于存储 UTF-8 编码的文本。

```solidity
contract StringExample {
    string public name;
    string public description;

    constructor() {
        name = "Solidity Tutorial";
        description = "Learning Solidity step by step";
    }

    function setName(string memory _name) public {
        name = _name;
    }

    function concatenateStrings(string memory a, string memory b)
        public
        pure
        returns (string memory)
    {
        return string(abi.encodePacked(a, " ", b));
    }

    function compareStrings(string memory a, string memory b)
        public
        pure
        returns (bool)
    {
        return keccak256(abi.encodePacked(a)) == keccak256(abi.encodePacked(b));
    }

    function getStringLength(string memory _str) public pure returns (uint256) {
        return bytes(_str).length;  // 转换为字节数组获取长度
    }
}
```

### 2.1.6 枚举类型 (enum)

枚举用于创建用户定义的类型，表示一组有限的选项。

```solidity
contract EnumExample {
    // 定义订单状态枚举
    enum OrderStatus {
        Pending,    // 0
        Confirmed,  // 1
        Shipped,    // 2
        Delivered,  // 3
        Cancelled   // 4
    }

    OrderStatus public currentStatus;

    constructor() {
        currentStatus = OrderStatus.Pending;
    }

    function confirmOrder() public {
        require(currentStatus == OrderStatus.Pending, "Order not pending");
        currentStatus = OrderStatus.Confirmed;
    }

    function shipOrder() public {
        require(currentStatus == OrderStatus.Confirmed, "Order not confirmed");
        currentStatus = OrderStatus.Shipped;
    }

    function deliverOrder() public {
        require(currentStatus == OrderStatus.Shipped, "Order not shipped");
        currentStatus = OrderStatus.Delivered;
    }

    function cancelOrder() public {
        require(
            currentStatus == OrderStatus.Pending ||
            currentStatus == OrderStatus.Confirmed,
            "Cannot cancel order"
        );
        currentStatus = OrderStatus.Cancelled;
    }

    function getStatus() public view returns (OrderStatus) {
        return currentStatus;
    }

    function getStatusAsInt() public view returns (uint8) {
        return uint8(currentStatus);  // 枚举可以转换为整数
    }
}
```

## 2.2 变量类型详解

### 2.2.1 状态变量 (State Variables)

状态变量永久存储在区块链上，是合约的核心数据。

```solidity
contract StateVariableExample {
    // 状态变量声明
    uint256 public totalSupply = 1000000;  // 公共状态变量
    mapping(address => uint256) private balances;  // 私有状态变量
    address internal owner;  // 内部状态变量

    // 常量状态变量（编译时确定）
    uint256 public constant MAX_SUPPLY = 10000000;
    string public constant NAME = "Example Token";

    // 不可变状态变量（构造函数中设置）
    uint256 public immutable DECIMALS;
    address public immutable CREATOR;

    constructor(uint256 _decimals) {
        DECIMALS = _decimals;
        CREATOR = msg.sender;
        owner = msg.sender;
    }

    function updateTotalSupply(uint256 _newSupply) public {
        require(msg.sender == owner, "Only owner can update");
        require(_newSupply <= MAX_SUPPLY, "Exceeds maximum supply");
        totalSupply = _newSupply;
    }
}
```

**状态变量特点**：

- 存储在区块链上，永久保存
- 消耗存储 Gas（昂贵）
- 可以设置可见性修饰符
- 支持常量和不可变修饰符

### 2.2.2 局部变量 (Local Variables)

局部变量只在函数执行期间存在。

```solidity
contract LocalVariableExample {
    uint256 public stateVar = 100;

    function localVariableDemo(uint256 _input) public view returns (uint256, uint256) {
        // 局部变量声明
        uint256 localVar = 50;
        uint256 result;

        // 局部变量计算
        if (_input > stateVar) {
            uint256 difference = _input - stateVar;  // 块作用域变量
            result = localVar + difference;
        } else {
            result = localVar + stateVar;
        }

        return (localVar, result);
    }

    function scopeExample() public pure returns (uint256) {
        uint256 x = 10;

        {
            uint256 y = 20;  // 块作用域
            x = x + y;
        }
        // y 在这里不可访问

        return x;  // 返回 30
    }
}
```

**局部变量特点**：

- 存储在内存或栈中
- 函数执行结束后销毁
- Gas 消耗相对较低
- 作用域限制在声明的代码块内

### 2.2.3 全局变量 (Global Variables)

Solidity 提供了许多全局变量，用于访问区块链和交易信息。

```solidity
contract GlobalVariableExample {
    event GlobalInfo(
        address sender,
        uint256 value,
        uint256 timestamp,
        uint256 blockNumber
    );

    function demonstrateGlobalVars() public payable {
        // 消息相关全局变量
        address sender = msg.sender;      // 调用者地址
        uint256 value = msg.value;       // 发送的以太币数量
        bytes calldata data = msg.data;  // 调用数据
        bytes4 sig = msg.sig;           // 函数选择器

        // 交易相关全局变量
        address origin = tx.origin;      // 交易发起者
        uint256 gasPrice = tx.gasprice; // Gas 价格

        // 区块相关全局变量
        address coinbase = block.coinbase;    // 矿工地址
        uint256 difficulty = block.difficulty; // 区块难度
        uint256 gasLimit = block.gaslimit;    // 区块 Gas 限制
        uint256 number = block.number;        // 区块号
        uint256 timestamp = block.timestamp;  // 区块时间戳
        bytes32 blockHash = blockhash(number - 1); // 指定区块哈希

        emit GlobalInfo(sender, value, timestamp, number);
    }

    function timeBasedFunction() public view returns (bool) {
        // 使用时间戳进行逻辑判断
        return block.timestamp > 1640995200; // 2022年1月1日之后
    }

    function requireOrigin() public view {
        // 确保是直接调用，不是通过其他合约
        require(msg.sender == tx.origin, "Contract calls not allowed");
    }
}
```

**常用全局变量分类**：

1. **消息变量**：

   - `msg.sender`：当前调用发起方
   - `msg.value`：发送的以太币数量
   - `msg.data`：完整的调用数据
   - `msg.gas`：剩余 Gas（已废弃）

2. **交易变量**：

   - `tx.origin`：交易的原始发送方
   - `tx.gasprice`：交易的 Gas 价格

3. **区块变量**：
   - `block.coinbase`：当前区块矿工地址
   - `block.timestamp`：当前区块时间戳
   - `block.number`：当前区块号
   - `block.difficulty`：当前区块难度

## 2.3 可见性修饰符

可见性修饰符控制函数和状态变量的访问范围。

```solidity
contract VisibilityExample {
    uint256 public publicVar = 100;      // 任何人都可以访问
    uint256 internal internalVar = 200;  // 当前合约和子合约可以访问
    uint256 private privateVar = 300;    // 只有当前合约可以访问

    // public 函数：任何人都可以调用
    function publicFunction() public pure returns (string memory) {
        return "This is a public function";
    }

    // external 函数：只能从外部调用
    function externalFunction() external pure returns (string memory) {
        return "This is an external function";
    }

    // internal 函数：当前合约和子合约可以调用
    function internalFunction() internal pure returns (string memory) {
        return "This is an internal function";
    }

    // private 函数：只有当前合约可以调用
    function privateFunction() private pure returns (string memory) {
        return "This is a private function";
    }

    function testInternalCall() public view returns (uint256, string memory) {
        // 可以调用 internal 和 private 函数
        return (internalVar, internalFunction());
    }

    function testExternalCall() public view returns (string memory) {
        // 调用 external 函数需要使用 this
        return this.externalFunction();
    }
}

// 继承合约，测试可见性
contract ChildContract is VisibilityExample {
    function accessParentVars() public view returns (uint256, uint256) {
        return (
            publicVar,    // 可以访问
            internalVar   // 可以访问
            // privateVar  // 无法访问，会编译错误
        );
    }

    function callParentFunctions() public pure returns (string memory) {
        return internalFunction();  // 可以调用 internal 函数
        // return privateFunction();  // 无法调用，会编译错误
    }
}
```

**可见性总结表**：

| 修饰符   | 外部合约 | 继承合约 | 当前合约 | 自动 getter |
| -------- | -------- | -------- | -------- | ----------- |
| public   | ✓        | ✓        | ✓        | ✓           |
| external | ✓        | ✗        | ✗\*      | ✗           |
| internal | ✗        | ✓        | ✓        | ✗           |
| private  | ✗        | ✗        | ✓        | ✗           |

\*external 函数可以通过 `this.functionName()` 在内部调用

## 2.4 存储位置

Solidity 中的数据可以存储在不同的位置，每种位置有不同的 Gas 消耗和行为。

```solidity
contract StorageLocationExample {
    // 状态变量存储在 storage 中
    uint256[] public storageArray;
    mapping(uint256 => string) public storageMapping;

    struct Person {
        string name;
        uint256 age;
    }

    Person[] public people;

    function storageExample() public {
        // storage：指向状态变量的引用
        uint256[] storage localStorageArray = storageArray;
        localStorageArray.push(42);  // 直接修改状态变量

        Person storage newPerson = people.push();
        newPerson.name = "Alice";
        newPerson.age = 25;
    }

    function memoryExample() public view returns (uint256[] memory) {
        // memory：创建临时副本
        uint256[] memory localMemoryArray = new uint256[](3);
        localMemoryArray[0] = 1;
        localMemoryArray[1] = 2;
        localMemoryArray[2] = 3;

        return localMemoryArray;  // 返回内存中的数组
    }

    function calldataExample(uint256[] calldata _inputArray)
        external
        pure
        returns (uint256)
    {
        // calldata：只读的输入数据
        return _inputArray.length;  // 可以读取但不能修改
    }

    function stringLocationExample(string memory _str) public pure returns (string memory) {
        // 字符串参数和返回值通常使用 memory
        return string(abi.encodePacked("Hello, ", _str));
    }

    function compareStorageAndMemory() public {
        storageArray.push(100);

        // storage 引用：修改会影响原始数据
        uint256[] storage storageRef = storageArray;
        storageRef[0] = 999;  // storageArray[0] 现在是 999

        // memory 副本：修改不会影响原始数据
        uint256[] memory memoryArray = storageArray;
        memoryArray[0] = 888;  // storageArray[0] 仍然是 999
    }
}
```

**存储位置特点**：

1. **storage**：

   - 永久存储在区块链上
   - Gas 消耗最高
   - 状态变量默认位置
   - 引用类型会创建引用

2. **memory**：

   - 临时存储在内存中
   - 中等 Gas 消耗
   - 函数参数和返回值的默认位置
   - 总是创建副本

3. **calldata**：
   - 只读的外部数据存储
   - Gas 消耗最低
   - 只能用于 external 函数参数
   - 不能修改

## 2.5 实战练习

### 练习 1：学生成绩管理系统

创建一个管理学生成绩的智能合约：

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract StudentGrades {
    // TODO: 定义学生结构体
    // TODO: 实现添加学生功能
    // TODO: 实现更新成绩功能
    // TODO: 实现查询功能
    // TODO: 添加适当的访问控制
}
```

**要求**：

1. 定义学生结构体（姓名、学号、成绩）
2. 使用映射存储学生信息
3. 实现添加、更新、查询功能
4. 添加教师权限控制
5. 使用事件记录重要操作

### 练习 2：数据类型转换工具

实现各种数据类型之间的转换：

```solidity
contract TypeConverter {
    // TODO: 实现字符串到数字的转换
    // TODO: 实现地址到字符串的转换
    // TODO: 实现字节数组操作
    // TODO: 实现枚举转换
}
```

### 练习 3：时间锁合约

创建一个基于时间的锁定合约：

```solidity
contract TimeLock {
    // TODO: 实现时间锁定功能
    // TODO: 添加解锁条件检查
    // TODO: 使用全局变量获取时间
}
```

## 2.6 最佳实践

### 2.6.1 变量命名规范

```solidity
contract NamingConventions {
    // 状态变量：使用驼峰命名法
    uint256 public totalSupply;
    address private contractOwner;

    // 常量：使用大写字母和下划线
    uint256 public constant MAX_SUPPLY = 1000000;

    // 函数参数：使用下划线前缀
    function transfer(address _to, uint256 _amount) public {
        // 局部变量：使用驼峰命名法
        uint256 currentBalance = balanceOf[msg.sender];
        require(currentBalance >= _amount, "Insufficient balance");
    }
}
```

### 2.6.2 Gas 优化建议

```solidity
contract GasOptimization {
    // 使用 uint256 而不是较小的整数类型
    uint256 public optimizedVar;  // 好
    uint8 public unoptimizedVar;  // 在某些情况下可能消耗更多 Gas

    // 使用 constant 和 immutable
    uint256 public constant FIXED_VALUE = 100;  // 编译时常量
    address public immutable OWNER;             // 运行时常量

    constructor() {
        OWNER = msg.sender;
    }

    // 批量操作而不是单个操作
    function batchTransfer(address[] calldata _recipients, uint256[] calldata _amounts)
        external
    {
        require(_recipients.length == _amounts.length, "Array length mismatch");

        for (uint256 i = 0; i < _recipients.length; i++) {
            // 批量转账逻辑
        }
    }
}
```

### 2.6.3 安全考虑

```solidity
contract SecurityConsiderations {
    mapping(address => uint256) public balances;

    function secureTransfer(address _to, uint256 _amount) public {
        // 1. 检查输入参数
        require(_to != address(0), "Cannot transfer to zero address");
        require(_amount > 0, "Amount must be positive");

        // 2. 检查状态
        require(balances[msg.sender] >= _amount, "Insufficient balance");

        // 3. 更新状态
        balances[msg.sender] -= _amount;
        balances[_to] += _amount;

        // 4. 触发事件
        emit Transfer(msg.sender, _to, _amount);
    }

    event Transfer(address indexed from, address indexed to, uint256 value);
}
```

## 2.7 章节总结

在本章中，我们深入学习了：

1. **基本数据类型**：布尔、整数、地址、字节、字符串和枚举
2. **变量类型**：状态变量、局部变量和全局变量的特点和用法
3. **可见性修饰符**：public、external、internal、private 的区别和应用
4. **存储位置**：storage、memory、calldata 的特点和 Gas 消耗

### 关键知识点

- 选择合适的数据类型可以优化 Gas 消耗
- 状态变量的修改会产生 Gas 费用
- 可见性修饰符控制函数和变量的访问范围
- 理解存储位置对于优化合约性能至关重要

### 下一章预告

在下一章《函数与控制流》中，我们将学习：

- 函数的定义和调用
- 函数修饰符和状态可变性
- 条件语句和循环结构
- 错误处理和异常机制

继续您的 Solidity 学习之旅：[第三章 - 函数与控制流](03_functions_control_flow.md)
