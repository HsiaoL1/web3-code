# 第三章：函数与控制流

## 本章学习目标

- 掌握函数的定义、参数和返回值
- 理解函数的可见性和状态可变性
- 学会使用条件语句和循环结构
- 掌握错误处理和异常机制
- 了解函数重载和回退函数
- 学习 Gas 优化技巧

## 3.1 函数基础

函数是智能合约的核心组件，用于实现业务逻辑和状态变更。

### 3.1.1 函数语法结构

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract FunctionBasics {
    uint256 public counter = 0;

    // 基本函数语法
    // function <函数名>(<参数列表>) <可见性> <状态可变性> [returns (<返回类型>)] { ... }

    // 无参数无返回值的函数
    function increment() public {
        counter++;
    }

    // 有参数的函数
    function add(uint256 _value) public {
        counter += _value;
    }

    // 有返回值的函数
    function getCounter() public view returns (uint256) {
        return counter;
    }

    // 多个参数和多个返回值
    function calculate(uint256 a, uint256 b)
        public
        pure
        returns (uint256 sum, uint256 product, uint256 difference)
    {
        sum = a + b;
        product = a * b;
        difference = a > b ? a - b : b - a;

        // 或者使用 return 语句
        // return (a + b, a * b, a > b ? a - b : b - a);
    }

    // 命名返回值
    function namedReturns(uint256 _input) public pure returns (uint256 doubled, bool isEven) {
        doubled = _input * 2;
        isEven = _input % 2 == 0;
        // 命名返回值可以直接赋值，不需要 return 语句
    }
}
```

### 3.1.2 函数参数和返回值

```solidity
contract ParametersAndReturns {
    struct Person {
        string name;
        uint256 age;
    }

    mapping(uint256 => Person) public people;
    uint256 public personCount;

    // 基本类型参数
    function createPerson(string memory _name, uint256 _age) public returns (uint256 id) {
        id = personCount++;
        people[id] = Person(_name, _age);
    }

    // 结构体参数
    function addPerson(Person memory _person) public returns (uint256) {
        uint256 id = personCount++;
        people[id] = _person;
        return id;
    }

    // 数组参数
    function processNumbers(uint256[] memory _numbers) public pure returns (uint256 total, uint256 average) {
        require(_numbers.length > 0, "Array cannot be empty");

        for (uint256 i = 0; i < _numbers.length; i++) {
            total += _numbers[i];
        }
        average = total / _numbers.length;
    }

    // 可变长度参数（使用 calldata 优化 Gas）
    function sumNumbers(uint256[] calldata _numbers) external pure returns (uint256) {
        uint256 sum = 0;
        for (uint256 i = 0; i < _numbers.length; i++) {
            sum += _numbers[i];
        }
        return sum;
    }

    // 返回复杂数据类型
    function getPerson(uint256 _id) public view returns (Person memory) {
        require(_id < personCount, "Person does not exist");
        return people[_id];
    }

    // 返回多个值
    function getPersonInfo(uint256 _id)
        public
        view
        returns (string memory name, uint256 age, bool exists)
    {
        if (_id < personCount) {
            Person memory person = people[_id];
            return (person.name, person.age, true);
        } else {
            return ("", 0, false);
        }
    }
}
```

## 3.2 函数修饰符详解

### 3.2.1 可见性修饰符

```solidity
contract VisibilityModifiers {
    uint256 private secretNumber = 42;

    // public: 可以从任何地方调用
    function publicFunction() public pure returns (string memory) {
        return "Anyone can call this function";
    }

    // external: 只能从外部调用（不能从合约内部调用）
    function externalFunction() external pure returns (string memory) {
        return "Only external calls allowed";
    }

    // internal: 只能从合约内部或继承合约调用
    function internalFunction() internal pure returns (string memory) {
        return "Internal use only";
    }

    // private: 只能从当前合约调用
    function privateFunction() private view returns (uint256) {
        return secretNumber;
    }

    // 内部调用示例
    function testInternalCalls() public view returns (string memory, uint256) {
        return (internalFunction(), privateFunction());
    }

    // 外部调用示例
    function testExternalCall() public view returns (string memory) {
        return this.externalFunction();  // 通过 this 调用 external 函数
    }
}
```

### 3.2.2 状态可变性修饰符

```solidity
contract StateMutability {
    uint256 public stateVariable = 100;
    uint256 public constant CONSTANT_VALUE = 42;

    // pure: 不读取也不修改状态
    function pureFunction(uint256 a, uint256 b) public pure returns (uint256) {
        return a + b + 10;  // 只使用参数和字面量
    }

    // view: 只读取状态，不修改
    function viewFunction() public view returns (uint256) {
        return stateVariable + CONSTANT_VALUE;  // 读取状态变量
    }

    // 默认: 可以修改状态
    function defaultFunction(uint256 _newValue) public {
        stateVariable = _newValue;  // 修改状态变量
    }

    // payable: 可以接收以太币
    function payableFunction() public payable {
        require(msg.value > 0, "Must send ether");
        // 处理接收到的以太币
    }

    // 状态可变性限制示例
    function demonstrateRestrictions() public view returns (uint256) {
        // view 函数中不能做的事情：
        // stateVariable = 200;  // 错误：不能修改状态
        // payableFunction();    // 错误：不能调用 payable 函数

        // view 函数中可以做的事情：
        uint256 localVar = 50;
        return stateVariable + localVar + CONSTANT_VALUE;
    }
}
```

### 3.2.3 函数重载

```solidity
contract FunctionOverloading {
    // 函数重载：相同名称，不同参数
    function process(uint256 _value) public pure returns (uint256) {
        return _value * 2;
    }

    function process(uint256 _value, uint256 _multiplier) public pure returns (uint256) {
        return _value * _multiplier;
    }

    function process(string memory _text) public pure returns (string memory) {
        return string(abi.encodePacked("Processed: ", _text));
    }

    // 调用重载函数的示例
    function testOverloading() public pure returns (uint256, uint256, string memory) {
        return (
            process(10),           // 调用第一个版本
            process(10, 5),        // 调用第二个版本
            process("hello")       // 调用第三个版本
        );
    }
}
```

## 3.3 控制流语句

### 3.3.1 条件语句

```solidity
contract ConditionalStatements {
    enum Status { Inactive, Active, Suspended, Terminated }

    function simpleIfElse(uint256 _value) public pure returns (string memory) {
        if (_value > 100) {
            return "High value";
        } else if (_value > 50) {
            return "Medium value";
        } else {
            return "Low value";
        }
    }

    function ternaryOperator(uint256 _value) public pure returns (string memory) {
        return _value > 50 ? "Greater than 50" : "50 or less";
    }

    function complexConditions(uint256 _value, bool _isActive)
        public
        pure
        returns (string memory)
    {
        if (_value > 100 && _isActive) {
            return "High value and active";
        } else if (_value > 50 || _isActive) {
            return "Medium value or active";
        } else {
            return "Low value and inactive";
        }
    }

    function switchLikeFunction(Status _status) public pure returns (string memory) {
        // Solidity 没有 switch 语句，使用 if-else 链
        if (_status == Status.Inactive) {
            return "Status is inactive";
        } else if (_status == Status.Active) {
            return "Status is active";
        } else if (_status == Status.Suspended) {
            return "Status is suspended";
        } else if (_status == Status.Terminated) {
            return "Status is terminated";
        } else {
            return "Unknown status";
        }
    }
}
```

### 3.3.2 循环语句

```solidity
contract LoopStatements {
    uint256[] public numbers;

    function forLoopExample() public {
        // 清空数组
        delete numbers;

        // for 循环添加数字
        for (uint256 i = 1; i <= 10; i++) {
            numbers.push(i * 2);
        }
    }

    function whileLoopExample() public {
        delete numbers;

        uint256 i = 1;
        while (i <= 5) {
            numbers.push(i * i);  // 添加平方数
            i++;
        }
    }

    function doWhileExample() public {
        delete numbers;

        uint256 i = 1;
        do {
            numbers.push(i * 3);
            i++;
        } while (i <= 3);
    }

    function findFirstEven(uint256[] calldata _numbers)
        external
        pure
        returns (uint256 index, bool found)
    {
        for (uint256 i = 0; i < _numbers.length; i++) {
            if (_numbers[i] % 2 == 0) {
                return (i, true);
            }
        }
        return (0, false);
    }

    function sumUntilLimit(uint256[] calldata _numbers, uint256 _limit)
        external
        pure
        returns (uint256 sum, uint256 count)
    {
        for (uint256 i = 0; i < _numbers.length; i++) {
            if (sum + _numbers[i] > _limit) {
                break;  // 中断循环
            }
            sum += _numbers[i];
            count++;
        }
    }

    function processEvenNumbers(uint256[] calldata _numbers)
        external
        pure
        returns (uint256[] memory evenNumbers)
    {
        // 首先计算偶数的数量
        uint256 evenCount = 0;
        for (uint256 i = 0; i < _numbers.length; i++) {
            if (_numbers[i] % 2 == 0) {
                evenCount++;
            }
        }

        // 创建结果数组
        evenNumbers = new uint256[](evenCount);
        uint256 index = 0;

        for (uint256 i = 0; i < _numbers.length; i++) {
            if (_numbers[i] % 2 == 0) {
                evenNumbers[index] = _numbers[i];
                index++;
            }
        }
    }

    // 注意：避免无限循环和过多的 Gas 消耗
    function gasEfficientLoop(uint256 _iterations) public pure returns (uint256) {
        require(_iterations <= 1000, "Too many iterations");  // 限制循环次数

        uint256 result = 0;
        for (uint256 i = 0; i < _iterations; i++) {
            result += i;
        }
        return result;
    }
}
```

## 3.4 错误处理机制

### 3.4.1 require、revert 和 assert

```solidity
contract ErrorHandling {
    mapping(address => uint256) public balances;
    address public owner;
    bool public contractActive = true;

    constructor() {
        owner = msg.sender;
        balances[msg.sender] = 1000;
    }

    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }

    modifier whenActive() {
        require(contractActive, "Contract is not active");
        _;
    }

    // require: 用于验证条件，适用于用户输入和外部条件
    function transfer(address _to, uint256 _amount) public whenActive {
        require(_to != address(0), "Cannot transfer to zero address");
        require(_amount > 0, "Transfer amount must be positive");
        require(balances[msg.sender] >= _amount, "Insufficient balance");

        balances[msg.sender] -= _amount;
        balances[_to] += _amount;
    }

    // revert: 用于复杂的错误条件
    function complexTransfer(address _to, uint256 _amount) public {
        if (_to == address(0)) {
            revert("Invalid recipient address");
        }

        if (_amount == 0) {
            revert("Amount cannot be zero");
        }

        if (balances[msg.sender] < _amount) {
            revert("Insufficient balance for transfer");
        }

        balances[msg.sender] -= _amount;
        balances[_to] += _amount;
    }

    // assert: 用于检查内部错误和不变量
    function assertExample(uint256 _value) public pure returns (uint256) {
        uint256 result = _value * 2;

        // assert 通常用于检查不应该失败的条件
        assert(result >= _value);  // 除非溢出，否则应该总是成立

        return result;
    }

    // 自定义错误（Solidity 0.8.4+）
    error InsufficientBalance(uint256 available, uint256 required);
    error InvalidAddress(address provided);

    function customErrorExample(address _to, uint256 _amount) public {
        if (_to == address(0)) {
            revert InvalidAddress(_to);
        }

        uint256 available = balances[msg.sender];
        if (available < _amount) {
            revert InsufficientBalance(available, _amount);
        }

        balances[msg.sender] -= _amount;
        balances[_to] += _amount;
    }

    // try-catch 用法（用于外部调用）
    function safeExternalCall(address _contract) public returns (bool success, bytes memory data) {
        try this.riskyFunction() returns (uint256 value) {
            return (true, abi.encode(value));
        } catch Error(string memory reason) {
            return (false, abi.encode(reason));
        } catch (bytes memory lowLevelData) {
            return (false, lowLevelData);
        }
    }

    function riskyFunction() public pure returns (uint256) {
        require(false, "This function always fails");
        return 42;
    }
}
```

### 3.4.2 错误处理最佳实践

```solidity
contract ErrorHandlingBestPractices {
    mapping(address => uint256) public balances;

    // 好的做法：清晰的错误消息
    function goodTransfer(address _to, uint256 _amount) public {
        require(_to != address(0), "Transfer: recipient is zero address");
        require(_amount > 0, "Transfer: amount must be positive");
        require(balances[msg.sender] >= _amount, "Transfer: insufficient balance");

        balances[msg.sender] -= _amount;
        balances[_to] += _amount;
    }

    // 更好的做法：使用自定义错误
    error ZeroAddress();
    error InsufficientBalance(uint256 balance, uint256 amount);
    error InvalidAmount();

    function betterTransfer(address _to, uint256 _amount) public {
        if (_to == address(0)) revert ZeroAddress();
        if (_amount == 0) revert InvalidAmount();
        if (balances[msg.sender] < _amount) {
            revert InsufficientBalance(balances[msg.sender], _amount);
        }

        balances[msg.sender] -= _amount;
        balances[_to] += _amount;
    }

    // 检查-效果-交互模式
    function withdrawPattern(uint256 _amount) public {
        // 1. 检查
        require(balances[msg.sender] >= _amount, "Insufficient balance");

        // 2. 效果
        balances[msg.sender] -= _amount;

        // 3. 交互
        payable(msg.sender).transfer(_amount);
    }
}
```

## 3.5 特殊函数

### 3.5.1 构造函数

```solidity
contract ConstructorExamples {
    address public owner;
    uint256 public creationTime;
    string public name;

    // 基本构造函数
    constructor(string memory _name) {
        owner = msg.sender;
        creationTime = block.timestamp;
        name = _name;
    }
}

// 继承中的构造函数
contract Parent {
    uint256 public parentValue;

    constructor(uint256 _value) {
        parentValue = _value;
    }
}

contract Child is Parent {
    uint256 public childValue;

    // 调用父合约构造函数
    constructor(uint256 _parentValue, uint256 _childValue) Parent(_parentValue) {
        childValue = _childValue;
    }
}
```

### 3.5.2 回退函数和接收函数

```solidity
contract FallbackAndReceive {
    mapping(address => uint256) public deposits;

    event Received(address sender, uint256 amount);
    event FallbackCalled(address sender, bytes data);

    // receive: 处理纯以太币转账
    receive() external payable {
        deposits[msg.sender] += msg.value;
        emit Received(msg.sender, msg.value);
    }

    // fallback: 处理不存在的函数调用
    fallback() external payable {
        deposits[msg.sender] += msg.value;
        emit FallbackCalled(msg.sender, msg.data);
    }

    function getBalance() public view returns (uint256) {
        return address(this).balance;
    }

    function withdraw() public {
        uint256 amount = deposits[msg.sender];
        require(amount > 0, "No deposits found");

        deposits[msg.sender] = 0;
        payable(msg.sender).transfer(amount);
    }
}
```

## 3.6 实战练习

### 练习 1：银行系统

创建一个简单的银行系统：

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract SimpleBank {
    // TODO: 实现存款功能
    // TODO: 实现取款功能
    // TODO: 实现转账功能
    // TODO: 添加利息计算
    // TODO: 实现账户冻结功能
}
```

**要求**：

1. 用户可以存款和取款
2. 实现用户间转账
3. 添加余额查询功能
4. 实现简单的利息计算
5. 管理员可以冻结账户

### 练习 2：投票系统

实现一个去中心化投票系统：

```solidity
contract VotingSystem {
    // TODO: 定义候选人结构
    // TODO: 实现投票功能
    // TODO: 添加投票限制
    // TODO: 实现结果统计
}
```

### 练习 3：拍卖合约

创建一个英式拍卖合约：

```solidity
contract Auction {
    // TODO: 实现出价功能
    // TODO: 添加时间限制
    // TODO: 实现退款机制
    // TODO: 处理拍卖结束
}
```

## 3.7 Gas 优化技巧

### 3.7.1 循环优化

```solidity
contract GasOptimization {
    uint256[] public data;

    // 低效的循环
    function inefficientLoop() public view returns (uint256) {
        uint256 sum = 0;
        for (uint256 i = 0; i < data.length; i++) {  // 每次读取 data.length
            sum += data[i];
        }
        return sum;
    }

    // 优化的循环
    function efficientLoop() public view returns (uint256) {
        uint256 sum = 0;
        uint256 length = data.length;  // 缓存长度
        for (uint256 i = 0; i < length; i++) {
            sum += data[i];
        }
        return sum;
    }

    // 更好的优化
    function veryEfficientLoop() public view returns (uint256) {
        uint256 sum = 0;
        uint256 length = data.length;
        for (uint256 i = 0; i < length;) {
            sum += data[i];
            unchecked { ++i; }  // 避免溢出检查
        }
        return sum;
    }
}
```

### 3.7.2 函数优化

```solidity
contract FunctionOptimization {
    // 使用 external 而不是 public（如果只从外部调用）
    function externalFunction(uint256[] calldata _data) external pure returns (uint256) {
        // calldata 比 memory 更省 Gas
        return _data.length;
    }

    // 短路求值优化
    function shortCircuit(uint256 _value) public pure returns (bool) {
        return _value > 0 && _value < 100;  // 如果第一个条件为 false，不会检查第二个
    }

    // 使用位运算
    function powerOfTwo(uint256 _value) public pure returns (bool) {
        return _value != 0 && (_value & (_value - 1)) == 0;
    }
}
```

## 3.8 章节总结

在本章中，我们全面学习了：

1. **函数基础**：语法结构、参数、返回值和重载
2. **函数修饰符**：可见性、状态可变性和 payable
3. **控制流**：条件语句、循环和流程控制
4. **错误处理**：require、revert、assert 和自定义错误
5. **特殊函数**：构造函数、回退函数和接收函数

### 关键知识点

- 正确选择函数修饰符可以提高安全性和效率
- 错误处理是智能合约安全的重要组成部分
- 循环操作需要注意 Gas 消耗限制
- 函数重载提供了代码的灵活性

### 下一章预告

在下一章《复杂数据结构》中，我们将学习：

- 数组的详细操作
- 映射的高级用法
- 结构体的设计模式
- 枚举的最佳实践

继续您的 Solidity 学习之旅：[第四章 - 复杂数据结构](04_complex_data_structures.md)
