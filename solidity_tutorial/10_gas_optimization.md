# 第十章：Gas 优化与性能

## 本章学习目标

- 深入理解以太坊 Gas 机制
- 掌握存储优化技术
- 学习函数和循环优化方法
- 了解数据结构选择对性能的影响
- 掌握高级 Gas 优化技巧
- 学会使用工具分析和优化 Gas 消耗

## 10.1 Gas 机制深入理解

### 10.1.1 Gas 基础概念

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * Gas 消耗说明：
 *
 * 基本操作的 Gas 成本：
 * - ADD, SUB, MUL: 3 gas
 * - DIV, MOD: 5 gas
 * - SSTORE (存储写入): 20,000 gas (新值) 或 5,000 gas (更新)
 * - SLOAD (存储读取): 2,100 gas (冷读取) 或 100 gas (热读取)
 * - CALL: 700 gas + 转账成本
 * - CREATE: 32,000 gas
 */

contract GasBasics {
    uint256 public storageVar;

    // 不同操作的 Gas 消耗示例
    function expensiveOperations() public {
        storageVar = 100;           // 大约 20,000 gas (新存储)
        storageVar = 200;           // 大约 5,000 gas (更新存储)

        uint256 temp = storageVar;  // 大约 2,100 gas (冷读取)
        temp = storageVar;          // 大约 100 gas (热读取)

        // 局部变量操作几乎不消耗 Gas
        uint256 local1 = 10;
        uint256 local2 = 20;
        uint256 local3 = local1 + local2;
    }

    // Memory vs Storage vs Calldata 的 Gas 对比
    function memoryExample(uint256[] memory data) public pure returns (uint256) {
        data[0] = 100; // Memory 操作，便宜
        return data[0];
    }

    function storageExample(uint256[] storage data) internal returns (uint256) {
        data[0] = 100; // Storage 操作，昂贵
        return data[0];
    }

    function calldataExample(uint256[] calldata data) external pure returns (uint256) {
        // Calldata 只读，最便宜
        return data[0];
    }
}
```

### 10.1.2 Gas 计算实例

```solidity
contract GasAnalysis {
    struct User {
        string name;
        uint256 age;
        bool active;
        uint256[] scores;
    }

    mapping(address => User) public users;
    address[] public userList;

    // 高 Gas 消耗的操作
    function expensiveUserCreation(
        string memory name,
        uint256 age,
        uint256[] memory scores
    ) public {
        User storage user = users[msg.sender];
        user.name = name;       // 昂贵：动态字符串存储
        user.age = age;         // 便宜：uint256 存储
        user.active = true;     // 便宜：bool 存储

        // 非常昂贵：动态数组存储
        for (uint256 i = 0; i < scores.length; i++) {
            user.scores.push(scores[i]);
        }

        userList.push(msg.sender); // 昂贵：数组扩展
    }

    // 优化后的用户创建
    function optimizedUserCreation(
        string calldata name,   // 使用 calldata 而不是 memory
        uint256 age,
        uint256[] calldata scores
    ) public {
        require(scores.length <= 10, "Too many scores"); // 限制大小

        User storage user = users[msg.sender];
        user.name = name;
        user.age = age;
        user.active = true;

        // 直接赋值而不是循环 push
        delete user.scores; // 清空现有数据
        for (uint256 i = 0; i < scores.length; ) {
            user.scores.push(scores[i]);
            unchecked { ++i; } // 避免溢出检查
        }

        userList.push(msg.sender);
    }

    // Gas 使用报告
    function gasReport() public view returns (
        uint256 gasLeft,
        uint256 gasPrice,
        uint256 blockGasLimit
    ) {
        gasLeft = gasleft();
        gasPrice = tx.gasprice;
        blockGasLimit = block.gaslimit;
    }
}
```

## 10.2 存储优化技术

### 10.2.1 变量打包

```solidity
contract StorageOptimization {
    // ❌ 低效的存储布局 - 使用 3 个存储槽
    struct InEfficientStruct {
        uint128 a;    // 槽 0: 16 字节
        bool b;       // 槽 1: 1 字节 (浪费 31 字节)
        uint128 c;    // 槽 2: 16 字节
    }

    // ✅ 高效的存储布局 - 使用 2 个存储槽
    struct EfficientStruct {
        uint128 a;    // 槽 0: 16 字节
        uint128 c;    // 槽 0: 16 字节 (总共 32 字节)
        bool b;       // 槽 1: 1 字节
    }

    // ❌ 更复杂的低效示例
    struct ComplexInefficient {
        bool isActive;        // 槽 0: 1 字节
        uint256 timestamp;    // 槽 1: 32 字节
        uint8 level;          // 槽 2: 1 字节
        address owner;        // 槽 3: 20 字节
        uint16 count;         // 槽 4: 2 字节
    }

    // ✅ 优化后的版本
    struct ComplexEfficient {
        uint256 timestamp;    // 槽 0: 32 字节
        address owner;        // 槽 1: 20 字节
        uint16 count;         // 槽 1: 2 字节
        uint8 level;          // 槽 1: 1 字节
        bool isActive;        // 槽 1: 1 字节
        // 槽 1 总计: 20 + 2 + 1 + 1 = 24 字节 (还有 8 字节可用)
    }

    // 实际使用示例
    mapping(address => EfficientStruct) public efficientUsers;
    mapping(address => ComplexEfficient) public complexUsers;

    function updateUser(uint128 a, uint128 c, bool b) public {
        EfficientStruct storage user = efficientUsers[msg.sender];

        // 一次性更新同一槽的所有变量可以节省 Gas
        user.a = a;
        user.c = c;  // 这两个在同一个槽，更新成本较低

        user.b = b;  // 这个在不同槽，需要额外的 SSTORE
    }

    // 批量更新策略
    function batchUpdateSameSlot(uint128 a, uint128 c) public {
        EfficientStruct storage user = efficientUsers[msg.sender];
        // 更新同一槽中的变量，只消耗一次 SSTORE
        user.a = a;
        user.c = c;
    }
}
```

### 10.2.2 位操作优化

```solidity
contract BitOperations {
    // 使用位字段压缩多个布尔值
    mapping(address => uint256) public userFlags;

    // 标志位定义
    uint256 constant IS_ACTIVE = 1;        // 0x001
    uint256 constant IS_PREMIUM = 2;       // 0x010
    uint256 constant IS_VERIFIED = 4;      // 0x100
    uint256 constant IS_ADMIN = 8;         // 0x1000

    // 级别存储在高位 (8-15 位)
    uint256 constant LEVEL_MASK = 0xFF00;
    uint256 constant LEVEL_SHIFT = 8;

    function setUserActive(address user, bool active) public {
        if (active) {
            userFlags[user] |= IS_ACTIVE;
        } else {
            userFlags[user] &= ~IS_ACTIVE;
        }
    }

    function setUserPremium(address user, bool premium) public {
        if (premium) {
            userFlags[user] |= IS_PREMIUM;
        } else {
            userFlags[user] &= ~IS_PREMIUM;
        }
    }

    function setUserLevel(address user, uint8 level) public {
        require(level <= 255, "Level too high");

        // 清除现有级别并设置新级别
        userFlags[user] = (userFlags[user] & ~LEVEL_MASK) |
                         (uint256(level) << LEVEL_SHIFT);
    }

    function getUserFlags(address user) public view returns (
        bool isActive,
        bool isPremium,
        bool isVerified,
        bool isAdmin,
        uint8 level
    ) {
        uint256 flags = userFlags[user];

        isActive = (flags & IS_ACTIVE) != 0;
        isPremium = (flags & IS_PREMIUM) != 0;
        isVerified = (flags & IS_VERIFIED) != 0;
        isAdmin = (flags & IS_ADMIN) != 0;
        level = uint8((flags & LEVEL_MASK) >> LEVEL_SHIFT);
    }

    // 批量设置多个标志
    function setMultipleFlags(
        address user,
        bool active,
        bool premium,
        bool verified,
        uint8 level
    ) public {
        uint256 flags = 0;

        if (active) flags |= IS_ACTIVE;
        if (premium) flags |= IS_PREMIUM;
        if (verified) flags |= IS_VERIFIED;

        flags |= (uint256(level) << LEVEL_SHIFT) & LEVEL_MASK;

        userFlags[user] = flags; // 只需要一次 SSTORE
    }

    // 位计数优化
    function countSetBits(uint256 value) public pure returns (uint256) {
        uint256 count = 0;

        while (value != 0) {
            count++;
            value &= value - 1; // 清除最低位的 1
        }

        return count;
    }

    // 检查是否为 2 的幂
    function isPowerOfTwo(uint256 value) public pure returns (bool) {
        return value != 0 && (value & (value - 1)) == 0;
    }
}
```

### 10.2.3 短路存储模式

```solidity
contract ShortCircuitStorage {
    mapping(address => uint256) public balances;
    mapping(address => bool) public hasBalance;

    // ❌ 总是写入存储
    function inefficientUpdate(uint256 newBalance) public {
        balances[msg.sender] = newBalance; // 总是消耗 Gas
    }

    // ✅ 短路写入
    function efficientUpdate(uint256 newBalance) public {
        if (balances[msg.sender] != newBalance) {
            balances[msg.sender] = newBalance; // 只在值改变时写入
        }
    }

    // ✅ 零值优化
    function updateWithZeroCheck(uint256 newBalance) public {
        uint256 currentBalance = balances[msg.sender];

        if (newBalance == 0 && currentBalance != 0) {
            // 删除映射条目可以获得 Gas 退款
            delete balances[msg.sender];
            hasBalance[msg.sender] = false;
        } else if (newBalance != 0) {
            if (currentBalance != newBalance) {
                balances[msg.sender] = newBalance;
            }
            if (!hasBalance[msg.sender]) {
                hasBalance[msg.sender] = true;
            }
        }
    }

    // 批量零值清理（获得 Gas 退款）
    address[] public usersWithBalance;
    mapping(address => uint256) public userIndex;

    function clearZeroBalances(address[] calldata users) external {
        for (uint256 i = 0; i < users.length; i++) {
            if (balances[users[i]] == 0 && hasBalance[users[i]]) {
                delete balances[users[i]];
                hasBalance[users[i]] = false;

                // 从活跃用户列表中移除
                removeFromUserList(users[i]);
            }
        }
    }

    function removeFromUserList(address user) internal {
        uint256 index = userIndex[user];
        uint256 lastIndex = usersWithBalance.length - 1;

        if (index != lastIndex) {
            address lastUser = usersWithBalance[lastIndex];
            usersWithBalance[index] = lastUser;
            userIndex[lastUser] = index;
        }

        usersWithBalance.pop();
        delete userIndex[user];
    }
}
```

## 10.3 函数优化技术

### 10.3.1 函数可见性优化

```solidity
contract FunctionOptimization {
    uint256[] public data;

    // ❌ 使用 public（生成 getter，即使不需要）
    function publicFunction(uint256[] memory _data) public pure returns (uint256) {
        return _data.length;
    }

    // ✅ 使用 external（只从外部调用）
    function externalFunction(uint256[] calldata _data) external pure returns (uint256) {
        return _data.length; // calldata 比 memory 便宜
    }

    // ✅ 内部函数优化
    function internalSum(uint256[] memory _data) internal pure returns (uint256) {
        uint256 sum = 0;
        uint256 length = _data.length; // 缓存长度

        for (uint256 i = 0; i < length;) {
            sum += _data[i];
            unchecked { ++i; } // 避免溢出检查
        }

        return sum;
    }

    // 函数修饰符优化
    modifier expensiveModifier() {
        require(msg.sender != address(0), "Invalid sender");
        require(data.length > 0, "No data");
        _;
    }

    // ❌ 多个修饰符增加 Gas 成本
    function multipleModifiers() public expensiveModifier view returns (uint256) {
        return data.length;
    }

    // ✅ 合并检查到函数内
    function combinedChecks() public view returns (uint256) {
        require(msg.sender != address(0), "Invalid sender");
        require(data.length > 0, "No data");
        return data.length;
    }

    // 返回值优化
    function multipleReturns() public view returns (uint256, uint256, uint256) {
        uint256 length = data.length;
        return (length, length * 2, length * 3);
    }

    // 使用结构体返回多个值
    struct ReturnData {
        uint256 length;
        uint256 doubled;
        uint256 tripled;
    }

    function structReturn() public view returns (ReturnData memory result) {
        uint256 length = data.length;
        result = ReturnData(length, length * 2, length * 3);
    }
}
```

### 10.3.2 循环优化

```solidity
contract LoopOptimization {
    uint256[] public numbers;
    mapping(uint256 => bool) public exists;

    // ❌ 低效的循环
    function inefficientLoop() public view returns (uint256) {
        uint256 sum = 0;

        for (uint256 i = 0; i < numbers.length; i++) { // 每次读取 length
            sum += numbers[i]; // 可能的溢出检查
        }

        return sum;
    }

    // ✅ 优化的循环
    function optimizedLoop() public view returns (uint256) {
        uint256 sum = 0;
        uint256 length = numbers.length; // 缓存长度

        for (uint256 i = 0; i < length;) { // 移除递增
            sum += numbers[i];
            unchecked { ++i; } // 避免溢出检查，使用前缀递增
        }

        return sum;
    }

    // ✅ 逆序循环（某些情况下更高效）
    function reverseLoop() public view returns (uint256) {
        uint256 sum = 0;
        uint256 length = numbers.length;

        if (length == 0) return 0;

        uint256 i = length;
        do {
            unchecked { --i; }
            sum += numbers[i];
        } while (i > 0);

        return sum;
    }

    // ✅ 分块处理大数组
    function processChunk(uint256 start, uint256 end) public view returns (uint256) {
        require(start < numbers.length, "Start out of bounds");
        require(end <= numbers.length, "End out of bounds");
        require(start < end, "Invalid range");

        uint256 sum = 0;

        for (uint256 i = start; i < end;) {
            sum += numbers[i];
            unchecked { ++i; }
        }

        return sum;
    }

    // ✅ 跳出循环优化
    function findAndBreak(uint256 target) public view returns (bool found, uint256 index) {
        uint256 length = numbers.length;

        for (uint256 i = 0; i < length;) {
            if (numbers[i] == target) {
                return (true, i); // 立即返回，避免继续循环
            }
            unchecked { ++i; }
        }

        return (false, 0);
    }

    // ✅ 使用映射替代线性搜索
    function addNumber(uint256 number) public {
        require(!exists[number], "Number already exists");

        numbers.push(number);
        exists[number] = true;
    }

    function hasNumber(uint256 number) public view returns (bool) {
        return exists[number]; // O(1) 而不是 O(n)
    }

    // 批量操作优化
    function batchAddNumbers(uint256[] calldata newNumbers) external {
        require(newNumbers.length <= 100, "Batch too large");

        uint256 length = newNumbers.length;

        for (uint256 i = 0; i < length;) {
            uint256 number = newNumbers[i];

            if (!exists[number]) {
                numbers.push(number);
                exists[number] = true;
            }

            unchecked { ++i; }
        }
    }
}
```

## 10.4 数据结构选择优化

### 10.4.1 映射 vs 数组

```solidity
contract DataStructureComparison {
    // 方案 1: 使用数组（适合小数据集和顺序访问）
    struct ArrayBased {
        address[] users;
        mapping(address => uint256) userIndex;
    }

    // 方案 2: 使用映射（适合大数据集和随机访问）
    struct MappingBased {
        mapping(address => bool) userExists;
        mapping(address => uint256) userData;
    }

    ArrayBased public arraySystem;
    MappingBased public mappingSystem;

    // 数组方案的操作
    function addUserArray(address user, uint256 data) public {
        require(arraySystem.userIndex[user] == 0, "User exists");

        arraySystem.users.push(user);
        arraySystem.userIndex[user] = arraySystem.users.length; // 1-based index
    }

    function removeUserArray(address user) public {
        uint256 index = arraySystem.userIndex[user];
        require(index > 0, "User not found");

        uint256 arrayIndex = index - 1;
        uint256 lastIndex = arraySystem.users.length - 1;

        if (arrayIndex != lastIndex) {
            address lastUser = arraySystem.users[lastIndex];
            arraySystem.users[arrayIndex] = lastUser;
            arraySystem.userIndex[lastUser] = index;
        }

        arraySystem.users.pop();
        delete arraySystem.userIndex[user];
    }

    function getUserCountArray() public view returns (uint256) {
        return arraySystem.users.length; // O(1)
    }

    function getAllUsersArray() public view returns (address[] memory) {
        return arraySystem.users; // 可以获取所有用户
    }

    // 映射方案的操作
    function addUserMapping(address user, uint256 data) public {
        require(!mappingSystem.userExists[user], "User exists");

        mappingSystem.userExists[user] = true;
        mappingSystem.userData[user] = data;
    }

    function removeUserMapping(address user) public {
        require(mappingSystem.userExists[user], "User not found");

        delete mappingSystem.userExists[user];
        delete mappingSystem.userData[user]; // Gas 退款
    }

    function hasUserMapping(address user) public view returns (bool) {
        return mappingSystem.userExists[user]; // O(1)
    }

    // 性能比较
    function performanceTest(address[] calldata users) external {
        // 数组方案：添加用户 O(1)，查找用户 O(1)，删除用户 O(1)
        // 映射方案：添加用户 O(1)，查找用户 O(1)，删除用户 O(1)

        uint256 gasStart = gasleft();

        for (uint256 i = 0; i < users.length && i < 50; i++) {
            addUserMapping(users[i], i);
        }

        uint256 gasUsed = gasStart - gasleft();
        // 记录 gas 使用情况...
    }
}
```

### 10.4.2 字符串优化

```solidity
contract StringOptimization {
    // ❌ 动态字符串存储昂贵
    mapping(address => string) public userNames;

    // ✅ 使用字节数组
    mapping(address => bytes32) public userNamesBytes32;

    // ✅ 字符串哈希映射
    mapping(address => bytes32) public userNameHashes;
    mapping(bytes32 => string) public hashToName;

    // 字符串比较优化
    function compareStrings(string memory a, string memory b) public pure returns (bool) {
        return keccak256(abi.encodePacked(a)) == keccak256(abi.encodePacked(b));
    }

    // 短字符串优化
    function setShortName(bytes32 name) public {
        userNamesBytes32[msg.sender] = name;
    }

    function getShortName(address user) public view returns (string memory) {
        bytes32 nameBytes = userNamesBytes32[user];

        // 转换 bytes32 到 string
        uint256 length = 0;
        while (length < 32 && nameBytes[length] != 0) {
            length++;
        }

        bytes memory nameArray = new bytes(length);
        for (uint256 i = 0; i < length; i++) {
            nameArray[i] = nameBytes[i];
        }

        return string(nameArray);
    }

    // 长字符串哈希存储
    function setLongName(string calldata name) public {
        bytes32 hash = keccak256(abi.encodePacked(name));
        userNameHashes[msg.sender] = hash;

        if (bytes(hashToName[hash]).length == 0) {
            hashToName[hash] = name; // 只存储一次
        }
    }

    function getLongName(address user) public view returns (string memory) {
        bytes32 hash = userNameHashes[user];
        return hashToName[hash];
    }

    // 字符串连接优化
    function concatenateOptimized(
        string calldata a,
        string calldata b
    ) public pure returns (string memory) {
        return string(abi.encodePacked(a, b)); // 比多次调用更高效
    }

    // 字符串长度检查
    function isValidLength(string calldata str, uint256 maxLength) public pure returns (bool) {
        return bytes(str).length <= maxLength;
    }
}
```

## 10.5 高级优化技巧

### 10.5.1 内联汇编优化

```solidity
contract AssemblyOptimization {
    // 高效的内存复制
    function efficientMemCopy(
        bytes memory src,
        bytes memory dst,
        uint256 len
    ) public pure {
        assembly {
            let srcPtr := add(src, 0x20)
            let dstPtr := add(dst, 0x20)

            for { let i := 0 } lt(i, len) { i := add(i, 0x20) } {
                mstore(add(dstPtr, i), mload(add(srcPtr, i)))
            }
        }
    }

    // 高效的数组求和
    function assemblySum(uint256[] memory arr) public pure returns (uint256 sum) {
        assembly {
            let len := mload(arr)
            let ptr := add(arr, 0x20)

            for { let i := 0 } lt(i, len) { i := add(i, 1) } {
                sum := add(sum, mload(add(ptr, mul(i, 0x20))))
            }
        }
    }

    // 高效的位操作
    function countBitsAssembly(uint256 value) public pure returns (uint256 count) {
        assembly {
            for { } value { value := and(value, sub(value, 1)) } {
                count := add(count, 1)
            }
        }
    }

    // 高效的地址检查
    function isContract(address addr) public view returns (bool result) {
        assembly {
            result := gt(extcodesize(addr), 0)
        }
    }

    // 高效的字节操作
    function getByteAt(bytes memory data, uint256 index) public pure returns (bytes1 result) {
        assembly {
            result := byte(0, mload(add(add(data, 0x20), index)))
        }
    }
}
```

### 10.5.2 预计算和查找表

```solidity
contract PrecomputedOptimization {
    // 预计算的平方表
    mapping(uint256 => uint256) public squares;

    constructor() {
        // 预计算常用值的平方
        for (uint256 i = 1; i <= 100; i++) {
            squares[i] = i * i;
        }
    }

    function getSquare(uint256 x) public view returns (uint256) {
        if (x <= 100) {
            return squares[x]; // O(1) 查找
        } else {
            return x * x; // 动态计算
        }
    }

    // 预计算的阶乘表
    uint256[21] public factorials = [
        1, 1, 2, 6, 24, 120, 720, 5040, 40320, 362880,
        3628800, 39916800, 479001600, 6227020800, 87178291200,
        1307674368000, 20922789888000, 355687428096000, 6402373705728000,
        121645100408832000, 2432902008176640000
    ];

    function getFactorial(uint256 n) public view returns (uint256) {
        require(n <= 20, "Factorial too large");
        return factorials[n];
    }

    // 预计算的幂表
    mapping(uint256 => mapping(uint256 => uint256)) public powers;

    function initializePowers() public {
        for (uint256 base = 2; base <= 10; base++) {
            uint256 power = 1;
            for (uint256 exp = 0; exp <= 10; exp++) {
                powers[base][exp] = power;
                power *= base;
            }
        }
    }

    function getPower(uint256 base, uint256 exp) public view returns (uint256) {
        if (base <= 10 && exp <= 10) {
            return powers[base][exp];
        } else {
            return base ** exp; // 动态计算
        }
    }

    // 位掩码查找表
    uint256[256] public bitCounts;

    function initializeBitCounts() public {
        for (uint256 i = 0; i < 256; i++) {
            bitCounts[i] = countBits(i);
        }
    }

    function countBits(uint256 value) internal pure returns (uint256 count) {
        while (value != 0) {
            count++;
            value &= value - 1;
        }
    }

    function fastBitCount(uint256 value) public view returns (uint256 count) {
        while (value > 0) {
            count += bitCounts[value & 0xFF];
            value >>= 8;
        }
    }
}
```

### 10.5.3 批量操作优化

```solidity
contract BatchOptimization {
    mapping(address => uint256) public balances;
    mapping(address => bool) public isUser;

    event Transfer(address indexed from, address indexed to, uint256 value);
    event BatchTransfer(address indexed from, uint256 totalAmount, uint256 recipientCount);

    // ❌ 单个转账（高 Gas）
    function singleTransfer(address to, uint256 amount) public {
        require(balances[msg.sender] >= amount, "Insufficient balance");

        balances[msg.sender] -= amount;
        balances[to] += amount;

        emit Transfer(msg.sender, to, amount);
    }

    // ✅ 批量转账优化
    function batchTransfer(
        address[] calldata recipients,
        uint256[] calldata amounts
    ) external {
        require(recipients.length == amounts.length, "Array length mismatch");
        require(recipients.length <= 200, "Batch too large");

        uint256 totalAmount = 0;
        uint256 length = recipients.length;

        // 第一遍：计算总额并验证
        for (uint256 i = 0; i < length;) {
            require(recipients[i] != address(0), "Invalid recipient");
            require(amounts[i] > 0, "Invalid amount");

            totalAmount += amounts[i];
            unchecked { ++i; }
        }

        require(balances[msg.sender] >= totalAmount, "Insufficient balance");

        // 第二遍：执行转账
        balances[msg.sender] -= totalAmount;

        for (uint256 i = 0; i < length;) {
            balances[recipients[i]] += amounts[i];
            unchecked { ++i; }
        }

        emit BatchTransfer(msg.sender, totalAmount, length);
    }

    // 相同金额的批量转账
    function batchTransferSameAmount(
        address[] calldata recipients,
        uint256 amount
    ) external {
        require(recipients.length <= 200, "Batch too large");
        require(amount > 0, "Invalid amount");

        uint256 length = recipients.length;
        uint256 totalAmount = amount * length;

        require(balances[msg.sender] >= totalAmount, "Insufficient balance");

        balances[msg.sender] -= totalAmount;

        for (uint256 i = 0; i < length;) {
            require(recipients[i] != address(0), "Invalid recipient");
            balances[recipients[i]] += amount;
            unchecked { ++i; }
        }

        emit BatchTransfer(msg.sender, totalAmount, length);
    }

    // 批量用户注册
    function batchRegisterUsers(address[] calldata users) external {
        require(users.length <= 100, "Batch too large");

        uint256 length = users.length;

        for (uint256 i = 0; i < length;) {
            address user = users[i];
            require(user != address(0), "Invalid user");
            require(!isUser[user], "User already registered");

            isUser[user] = true;
            unchecked { ++i; }
        }
    }

    // 打包多个操作
    struct Operation {
        uint8 opType; // 0: transfer, 1: register, 2: unregister
        address target;
        uint256 amount;
    }

    function batchOperations(Operation[] calldata ops) external {
        require(ops.length <= 50, "Too many operations");

        uint256 length = ops.length;

        for (uint256 i = 0; i < length;) {
            Operation calldata op = ops[i];

            if (op.opType == 0) {
                // Transfer
                require(balances[msg.sender] >= op.amount, "Insufficient balance");
                balances[msg.sender] -= op.amount;
                balances[op.target] += op.amount;
            } else if (op.opType == 1) {
                // Register
                require(!isUser[op.target], "Already registered");
                isUser[op.target] = true;
            } else if (op.opType == 2) {
                // Unregister
                require(isUser[op.target], "Not registered");
                isUser[op.target] = false;
            }

            unchecked { ++i; }
        }
    }
}
```

## 10.6 Gas 分析工具

### 10.6.1 Gas 报告合约

```solidity
contract GasReporter {
    struct GasReport {
        uint256 gasUsed;
        uint256 gasPrice;
        uint256 gasCost;
        uint256 timestamp;
    }

    mapping(string => GasReport[]) public gasReports;

    modifier measureGas(string memory operation) {
        uint256 gasStart = gasleft();
        _;
        uint256 gasUsed = gasStart - gasleft();

        gasReports[operation].push(GasReport({
            gasUsed: gasUsed,
            gasPrice: tx.gasprice,
            gasCost: gasUsed * tx.gasprice,
            timestamp: block.timestamp
        }));
    }

    function expensiveOperation() public measureGas("expensiveOperation") {
        // 模拟昂贵操作
        for (uint256 i = 0; i < 100; i++) {
            // 一些计算
        }
    }

    function getGasReport(string memory operation)
        public
        view
        returns (GasReport[] memory)
    {
        return gasReports[operation];
    }

    function getAverageGas(string memory operation)
        public
        view
        returns (uint256 avgGas, uint256 samples)
    {
        GasReport[] memory reports = gasReports[operation];
        samples = reports.length;

        if (samples == 0) return (0, 0);

        uint256 totalGas = 0;
        for (uint256 i = 0; i < samples; i++) {
            totalGas += reports[i].gasUsed;
        }

        avgGas = totalGas / samples;
    }

    // 比较不同实现的 Gas 消耗
    function compareImplementations() public {
        uint256 gasStart;
        uint256 gasUsed1;
        uint256 gasUsed2;

        // 实现 1
        gasStart = gasleft();
        inefficientImplementation();
        gasUsed1 = gasStart - gasleft();

        // 实现 2
        gasStart = gasleft();
        efficientImplementation();
        gasUsed2 = gasStart - gasleft();

        // 记录对比结果
        gasReports["inefficient"].push(GasReport({
            gasUsed: gasUsed1,
            gasPrice: tx.gasprice,
            gasCost: gasUsed1 * tx.gasprice,
            timestamp: block.timestamp
        }));

        gasReports["efficient"].push(GasReport({
            gasUsed: gasUsed2,
            gasPrice: tx.gasprice,
            gasCost: gasUsed2 * tx.gasprice,
            timestamp: block.timestamp
        }));
    }

    function inefficientImplementation() internal pure {
        uint256 sum = 0;
        for (uint256 i = 0; i < 50; i++) {
            sum += i * i;
        }
    }

    function efficientImplementation() internal pure {
        uint256 sum = 0;
        for (uint256 i = 0; i < 50;) {
            sum += i * i;
            unchecked { ++i; }
        }
    }
}
```

### 10.6.2 优化检查清单

```solidity
/**
 * Gas 优化检查清单
 *
 * 1. 存储优化
 *    - [ ] 变量打包到同一个存储槽
 *    - [ ] 使用适当的变量类型大小
 *    - [ ] 删除未使用的存储变量获得退款
 *    - [ ] 使用 constant 和 immutable
 *
 * 2. 函数优化
 *    - [ ] 使用 external 而不是 public（适当时）
 *    - [ ] 使用 calldata 而不是 memory
 *    - [ ] 避免不必要的修饰符
 *    - [ ] 合并相似的函数
 *
 * 3. 循环优化
 *    - [ ] 缓存数组长度
 *    - [ ] 使用 unchecked 避免溢出检查
 *    - [ ] 使用前缀递增/递减
 *    - [ ] 实现早期退出
 *
 * 4. 数据结构选择
 *    - [ ] 映射 vs 数组的权衡
 *    - [ ] 使用位操作压缩数据
 *    - [ ] 字符串优化策略
 *
 * 5. 批量操作
 *    - [ ] 实现批量函数
 *    - [ ] 限制批量大小
 *    - [ ] 优化批量验证
 *
 * 6. 高级技巧
 *    - [ ] 使用内联汇编（谨慎）
 *    - [ ] 预计算常用值
 *    - [ ] 短路存储写入
 */

contract OptimizationChecklist {
    // 示例：应用所有优化技术的合约

    // 1. 优化的存储布局
    struct User {
        uint128 balance;    // 16 字节
        uint128 lastAction; // 16 字节 (槽 0: 32 字节)

        address addr;       // 20 字节
        uint8 level;        // 1 字节
        bool isActive;      // 1 字节 (槽 1: 22 字节)
    }

    mapping(address => User) public users;
    mapping(address => uint256) public userFlags; // 位压缩的标志

    // 2. 优化的常量
    uint256 public constant MAX_LEVEL = 100;
    uint256 private constant ACTIVE_FLAG = 1;

    // 3. 优化的事件
    event UserUpdated(address indexed user, uint256 indexed flags);

    // 4. 优化的函数
    function updateUser(
        uint128 balance,
        uint8 level,
        bool isActive
    ) external {
        // 输入验证
        require(level <= MAX_LEVEL, "Level too high");

        User storage user = users[msg.sender];

        // 批量更新同一槽的变量
        user.balance = balance;
        user.lastAction = uint128(block.timestamp);

        // 更新另一个槽的变量
        user.addr = msg.sender;
        user.level = level;
        user.isActive = isActive;

        // 使用位操作设置标志
        if (isActive) {
            userFlags[msg.sender] |= ACTIVE_FLAG;
        } else {
            userFlags[msg.sender] &= ~ACTIVE_FLAG;
        }

        emit UserUpdated(msg.sender, userFlags[msg.sender]);
    }

    // 5. 优化的批量操作
    function batchUpdateUsers(
        address[] calldata addresses,
        uint128[] calldata balances,
        uint8[] calldata levels
    ) external {
        uint256 length = addresses.length;
        require(length == balances.length && length == levels.length, "Length mismatch");
        require(length <= 50, "Batch too large");

        for (uint256 i = 0; i < length;) {
            User storage user = users[addresses[i]];
            user.balance = balances[i];
            user.level = levels[i];
            user.lastAction = uint128(block.timestamp);

            unchecked { ++i; }
        }
    }

    // 6. 使用查找表优化
    mapping(uint256 => uint256) public precomputed;

    function initializePrecomputed() external {
        for (uint256 i = 1; i <= 20; i++) {
            precomputed[i] = i * i * i; // 立方数
        }
    }

    function getCube(uint256 x) external view returns (uint256) {
        if (x <= 20) {
            return precomputed[x];
        } else {
            return x * x * x;
        }
    }
}
```

## 10.7 实战练习

### 练习 1：优化存储布局

给定一个低效的结构体，重新设计以最小化存储槽使用：

```solidity
// 需要优化的结构体
struct IneffcientStruct {
    bool flag1;
    uint256 bigNumber;
    uint8 smallNumber;
    bool flag2;
    address user;
    uint16 mediumNumber;
}
```

### 练习 2：循环优化挑战

优化以下函数以减少 Gas 消耗：

```solidity
function inefficientSum(uint256[] memory numbers) public pure returns (uint256) {
    uint256 total = 0;
    for (uint256 i = 0; i < numbers.length; i++) {
        if (numbers[i] > 0) {
            total = total + numbers[i];
        }
    }
    return total;
}
```

### 练习 3：设计高效的投票系统

创建一个 Gas 高效的投票合约，要求：

- 支持多个提案
- 用户可以投票给多个提案
- 快速查询投票结果
- 最小化存储成本

## 10.8 章节总结

在本章中，我们深入学习了：

1. **Gas 机制**：理解 Gas 计算和成本结构
2. **存储优化**：变量打包、位操作、短路存储
3. **函数优化**：可见性、参数类型、修饰符优化
4. **循环优化**：缓存、unchecked、逆序循环
5. **数据结构选择**：映射 vs 数组、字符串处理
6. **高级技巧**：内联汇编、预计算、批量操作
7. **分析工具**：Gas 报告和优化检查

### 优化原则

1. **测量优先**：始终测量 Gas 消耗再优化
2. **权衡考虑**：平衡可读性和效率
3. **渐进优化**：从影响最大的优化开始
4. **验证结果**：确保优化不破坏功能

Gas 优化是一个持续的过程，需要在开发的每个阶段都保持关注。

继续您的 Solidity 学习之旅：[第十一章 - 测试与部署](11_testing_deployment.md)
