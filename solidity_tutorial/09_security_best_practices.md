# 第九章：安全性与最佳实践

## 本章学习目标

- 了解智能合约常见的安全漏洞
- 掌握防范攻击的技术手段
- 学习安全编程模式和最佳实践
- 了解安全审计的方法和工具
- 掌握代码质量保证策略

## 9.1 常见安全漏洞

### 9.1.1 重入攻击 (Reentrancy Attack)

重入攻击是最危险的智能合约漏洞之一，攻击者利用外部调用重新进入函数。

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

// 漏洞示例：易受重入攻击的合约
contract VulnerableBank {
    mapping(address => uint256) public balances;

    function deposit() public payable {
        balances[msg.sender] += msg.value;
    }

    // 危险：易受重入攻击
    function withdraw(uint256 amount) public {
        require(balances[msg.sender] >= amount, "Insufficient balance");

        // 先转账再更新状态 - 这是错误的顺序！
        (bool success, ) = msg.sender.call{value: amount}("");
        require(success, "Transfer failed");

        balances[msg.sender] -= amount; // 攻击者可以在这之前重新调用 withdraw
    }

    function getBalance() public view returns (uint256) {
        return address(this).balance;
    }
}

// 攻击合约
contract ReentrancyAttacker {
    VulnerableBank public bank;
    uint256 public attackAmount = 1 ether;

    constructor(address _bankAddress) {
        bank = VulnerableBank(_bankAddress);
    }

    function attack() public payable {
        require(msg.value >= attackAmount, "Need at least 1 ETH to attack");

        bank.deposit{value: attackAmount}();
        bank.withdraw(attackAmount);
    }

    // 重入点：当银行调用转账时，这个函数会被执行
    receive() external payable {
        if (address(bank).balance >= attackAmount) {
            bank.withdraw(attackAmount); // 重入攻击！
        }
    }

    function getBalance() public view returns (uint256) {
        return address(this).balance;
    }
}

// 安全的银行合约实现
contract SecureBank {
    mapping(address => uint256) public balances;
    bool private locked;

    modifier noReentrant() {
        require(!locked, "ReentrancyGuard: reentrant call");
        locked = true;
        _;
        locked = false;
    }

    function deposit() public payable {
        balances[msg.sender] += msg.value;
    }

    // 安全的提取函数
    function withdraw(uint256 amount) public noReentrant {
        require(balances[msg.sender] >= amount, "Insufficient balance");

        // 检查-效果-交互模式
        // 1. 检查已在 require 中完成
        // 2. 效果：先更新状态
        balances[msg.sender] -= amount;

        // 3. 交互：最后进行外部调用
        (bool success, ) = msg.sender.call{value: amount}("");
        require(success, "Transfer failed");
    }

    // 替代方案：使用 pull 模式
    mapping(address => uint256) public pendingWithdrawals;

    function initiateWithdraw(uint256 amount) public {
        require(balances[msg.sender] >= amount, "Insufficient balance");

        balances[msg.sender] -= amount;
        pendingWithdrawals[msg.sender] += amount;
    }

    function withdraw() public {
        uint256 amount = pendingWithdrawals[msg.sender];
        require(amount > 0, "No pending withdrawal");

        pendingWithdrawals[msg.sender] = 0;

        (bool success, ) = msg.sender.call{value: amount}("");
        require(success, "Transfer failed");
    }
}
```

### 9.1.2 整数溢出和下溢

在 Solidity 0.8.0 之前，整数溢出是一个严重问题。

```solidity
// Solidity 0.8.0 之前的漏洞示例
contract IntegerOverflowExample {
    mapping(address => uint256) public balances;

    // 在 0.8.0 之前，这个函数容易受到下溢攻击
    function transfer(address to, uint256 amount) public {
        // 如果 amount > balances[msg.sender]，这会导致下溢
        // 在 0.8.0 之前会回绕到一个很大的数
        balances[msg.sender] -= amount; // 潜在的下溢
        balances[to] += amount;
    }

    // Solidity 0.8.0+ 的安全版本
    function safeTransfer(address to, uint256 amount) public {
        require(balances[msg.sender] >= amount, "Insufficient balance");

        // 在 0.8.0+ 中，这些操作会自动检查溢出
        balances[msg.sender] -= amount;
        balances[to] += amount;
    }

    // 显式使用 unchecked 来禁用溢出检查（谨慎使用）
    function unsafeTransfer(address to, uint256 amount) public {
        require(balances[msg.sender] >= amount, "Insufficient balance");

        unchecked {
            balances[msg.sender] -= amount;
            balances[to] += amount;
        }
    }
}

// 使用 SafeMath 库（0.8.0 之前需要）
library SafeMath {
    function add(uint256 a, uint256 b) internal pure returns (uint256) {
        uint256 c = a + b;
        require(c >= a, "SafeMath: addition overflow");
        return c;
    }

    function sub(uint256 a, uint256 b) internal pure returns (uint256) {
        require(b <= a, "SafeMath: subtraction overflow");
        return a - b;
    }

    function mul(uint256 a, uint256 b) internal pure returns (uint256) {
        if (a == 0) return 0;
        uint256 c = a * b;
        require(c / a == b, "SafeMath: multiplication overflow");
        return c;
    }

    function div(uint256 a, uint256 b) internal pure returns (uint256) {
        require(b > 0, "SafeMath: division by zero");
        return a / b;
    }
}
```

### 9.1.3 访问控制漏洞

不当的访问控制可能导致未授权的函数调用。

```solidity
contract AccessControlVulnerabilities {
    address public owner;
    mapping(address => bool) public admins;

    constructor() {
        owner = msg.sender;
    }

    // 漏洞：忘记访问控制修饰符
    function dangerousFunction() public {
        // 任何人都可以调用这个函数！
        selfdestruct(payable(msg.sender));
    }

    // 漏洞：tx.origin vs msg.sender
    function vulnerableFunction() public {
        require(tx.origin == owner, "Not owner"); // 危险！使用 tx.origin
        // 可以通过钓鱼攻击绕过
    }

    // 安全的访问控制
    modifier onlyOwner() {
        require(msg.sender == owner, "Not the owner");
        _;
    }

    modifier onlyAdmin() {
        require(admins[msg.sender] || msg.sender == owner, "Not an admin");
        _;
    }

    function secureFunction() public onlyOwner {
        // 只有所有者可以调用
    }

    function addAdmin(address admin) public onlyOwner {
        admins[admin] = true;
    }

    function removeAdmin(address admin) public onlyOwner {
        admins[admin] = false;
    }

    // 更安全的所有权转移
    address public pendingOwner;

    function transferOwnership(address newOwner) public onlyOwner {
        require(newOwner != address(0), "New owner cannot be zero address");
        pendingOwner = newOwner;
    }

    function acceptOwnership() public {
        require(msg.sender == pendingOwner, "Not pending owner");
        owner = pendingOwner;
        pendingOwner = address(0);
    }
}

// 钓鱼攻击示例
contract PhishingAttack {
    AccessControlVulnerabilities target;

    constructor(address _target) {
        target = AccessControlVulnerabilities(_target);
    }

    function attack() public {
        // 如果目标合约的所有者调用这个函数
        // tx.origin 仍然是所有者，但 msg.sender 是这个合约
        target.vulnerableFunction(); // 这会成功，因为 tx.origin 是所有者
    }
}
```

### 9.1.4 前置交易攻击 (Front-running)

攻击者监视内存池中的交易并抢先执行。

```solidity
contract FrontRunningVulnerable {
    mapping(address => uint256) public balances;
    uint256 public prize = 10 ether;
    bytes32 public hashedAnswer;
    bool public solved = false;

    constructor(bytes32 _hashedAnswer) payable {
        hashedAnswer = _hashedAnswer;
    }

    // 漏洞：答案在交易中可见
    function submitAnswer(string memory answer) public {
        require(!solved, "Already solved");
        require(keccak256(abi.encodePacked(answer)) == hashedAnswer, "Wrong answer");

        solved = true;
        balances[msg.sender] += prize;
    }
}

// 防止前置攻击的解决方案
contract FrontRunningResistant {
    mapping(address => uint256) public balances;
    uint256 public prize = 10 ether;
    bytes32 public hashedAnswer;
    bool public solved = false;

    // 承诺-揭示模式
    struct Commitment {
        bytes32 commitment;
        uint256 timestamp;
        bool revealed;
    }

    mapping(address => Commitment) public commitments;
    uint256 public commitPeriod = 1 hours;
    uint256 public revealPeriod = 1 hours;
    uint256 public gameStartTime;

    constructor(bytes32 _hashedAnswer) payable {
        hashedAnswer = _hashedAnswer;
        gameStartTime = block.timestamp;
    }

    // 阶段 1：提交承诺
    function commitAnswer(bytes32 commitment) public {
        require(
            block.timestamp < gameStartTime + commitPeriod,
            "Commit period ended"
        );
        require(commitments[msg.sender].commitment == 0, "Already committed");

        commitments[msg.sender] = Commitment({
            commitment: commitment,
            timestamp: block.timestamp,
            revealed: false
        });
    }

    // 阶段 2：揭示答案
    function revealAnswer(string memory answer, uint256 nonce) public {
        require(
            block.timestamp >= gameStartTime + commitPeriod,
            "Commit period not ended"
        );
        require(
            block.timestamp < gameStartTime + commitPeriod + revealPeriod,
            "Reveal period ended"
        );
        require(!solved, "Already solved");

        Commitment storage commitment = commitments[msg.sender];
        require(commitment.commitment != 0, "No commitment found");
        require(!commitment.revealed, "Already revealed");

        // 验证承诺
        bytes32 hash = keccak256(abi.encodePacked(answer, nonce, msg.sender));
        require(hash == commitment.commitment, "Invalid commitment");

        commitment.revealed = true;

        // 检查答案
        if (keccak256(abi.encodePacked(answer)) == hashedAnswer) {
            solved = true;
            balances[msg.sender] += prize;
        }
    }
}
```

### 9.1.5 时间操控攻击

依赖 `block.timestamp` 或 `block.number` 进行关键逻辑判断可能被矿工操控。

```solidity
contract TimeManipulationVulnerable {
    uint256 public gameEndTime;
    address public winner;
    uint256 public prize;

    constructor() payable {
        gameEndTime = block.timestamp + 1 hours;
        prize = msg.value;
    }

    // 漏洞：依赖 block.timestamp
    function play() public payable {
        require(block.timestamp < gameEndTime, "Game ended");
        require(msg.value > 0, "Must send ETH");

        // 矿工可以稍微调整 block.timestamp
        if (block.timestamp % 2 == 0) {
            winner = msg.sender;
        }
    }

    function claimPrize() public {
        require(block.timestamp >= gameEndTime, "Game not ended");
        require(msg.sender == winner, "Not the winner");

        payable(winner).transfer(prize);
    }
}

// 更安全的时间处理
contract TimeManipulationResistant {
    uint256 public gameEndBlock;
    address public winner;
    uint256 public prize;
    uint256 public constant BLOCK_TIME_TOLERANCE = 15; // 15 秒容差

    constructor() payable {
        gameEndBlock = block.number + 240; // 约 1 小时（假设 15 秒一个块）
        prize = msg.value;
    }

    function play() public payable {
        require(block.number < gameEndBlock, "Game ended");
        require(msg.value > 0, "Must send ETH");

        // 使用更难操控的随机性来源
        uint256 randomness = uint256(keccak256(abi.encodePacked(
            block.difficulty,
            block.timestamp,
            msg.sender,
            block.coinbase
        )));

        if (randomness % 2 == 0) {
            winner = msg.sender;
        }
    }

    // 使用区块号而不是时间戳
    function claimPrize() public {
        require(block.number >= gameEndBlock, "Game not ended");
        require(msg.sender == winner, "Not the winner");

        payable(winner).transfer(prize);
    }
}
```

## 9.2 安全编程模式

### 9.2.1 检查-效果-交互模式

这是防止重入攻击的最重要模式。

```solidity
contract ChecksEffectsInteractions {
    mapping(address => uint256) public balances;
    mapping(address => bool) public hasVoted;
    uint256 public totalVotes;

    // 正确的模式示例
    function withdraw(uint256 amount) public {
        // 1. 检查 (Checks)
        require(balances[msg.sender] >= amount, "Insufficient balance");
        require(amount > 0, "Amount must be positive");

        // 2. 效果 (Effects) - 更新状态
        balances[msg.sender] -= amount;

        // 3. 交互 (Interactions) - 外部调用
        (bool success, ) = msg.sender.call{value: amount}("");
        require(success, "Transfer failed");
    }

    // 投票函数示例
    function vote() public {
        // 1. 检查
        require(!hasVoted[msg.sender], "Already voted");
        require(balances[msg.sender] > 0, "Must have balance to vote");

        // 2. 效果
        hasVoted[msg.sender] = true;
        totalVotes++;

        // 3. 交互（如果需要）
        // 外部调用应该放在最后
    }
}
```

### 9.2.2 拉取支付模式

避免主动推送支付，让用户主动拉取。

```solidity
contract PullPayment {
    mapping(address => uint256) public payments;

    event PaymentDeposited(address indexed payee, uint256 amount);
    event PaymentWithdrawn(address indexed payee, uint256 amount);

    // 存入支付
    function depositPayment(address payee) public payable {
        require(payee != address(0), "Invalid payee");
        require(msg.value > 0, "Must send value");

        payments[payee] += msg.value;
        emit PaymentDeposited(payee, msg.value);
    }

    // 用户主动提取支付
    function withdrawPayment() public {
        uint256 payment = payments[msg.sender];
        require(payment > 0, "No payment available");

        payments[msg.sender] = 0;

        (bool success, ) = msg.sender.call{value: payment}("");
        require(success, "Transfer failed");

        emit PaymentWithdrawn(msg.sender, payment);
    }

    // 检查待提取金额
    function getPayment(address payee) public view returns (uint256) {
        return payments[payee];
    }
}

// 用于拍卖的拉取支付示例
contract SecureAuction {
    address public highestBidder;
    uint256 public highestBid;
    mapping(address => uint256) public pendingReturns;
    bool public auctionEnded;

    event BidPlaced(address indexed bidder, uint256 amount);
    event AuctionEnded(address winner, uint256 amount);

    function bid() public payable {
        require(!auctionEnded, "Auction ended");
        require(msg.value > highestBid, "Bid too low");

        // 将之前的最高出价者的资金标记为可提取
        if (highestBidder != address(0)) {
            pendingReturns[highestBidder] += highestBid;
        }

        highestBidder = msg.sender;
        highestBid = msg.value;

        emit BidPlaced(msg.sender, msg.value);
    }

    // 失败的竞拍者提取资金
    function withdraw() public {
        uint256 amount = pendingReturns[msg.sender];
        require(amount > 0, "No funds to withdraw");

        pendingReturns[msg.sender] = 0;

        (bool success, ) = msg.sender.call{value: amount}("");
        require(success, "Transfer failed");
    }

    function endAuction() public {
        auctionEnded = true;
        emit AuctionEnded(highestBidder, highestBid);
    }
}
```

### 9.2.3 限制模式

实施各种限制来防止滥用。

```solidity
contract RateLimitedContract {
    mapping(address => uint256) public lastActionTime;
    mapping(address => uint256) public actionCount;
    mapping(address => uint256) public dailyActions;
    mapping(address => uint256) public lastResetDay;

    uint256 public constant MIN_DELAY = 1 minutes;
    uint256 public constant MAX_DAILY_ACTIONS = 10;
    uint256 public constant MAX_ACTIONS_PER_USER = 100;

    modifier rateLimited() {
        require(
            block.timestamp >= lastActionTime[msg.sender] + MIN_DELAY,
            "Action too frequent"
        );

        // 重置每日计数
        uint256 currentDay = block.timestamp / 1 days;
        if (lastResetDay[msg.sender] < currentDay) {
            dailyActions[msg.sender] = 0;
            lastResetDay[msg.sender] = currentDay;
        }

        require(
            dailyActions[msg.sender] < MAX_DAILY_ACTIONS,
            "Daily limit exceeded"
        );

        require(
            actionCount[msg.sender] < MAX_ACTIONS_PER_USER,
            "User action limit exceeded"
        );

        lastActionTime[msg.sender] = block.timestamp;
        dailyActions[msg.sender]++;
        actionCount[msg.sender]++;
        _;
    }

    function performAction() public rateLimited {
        // 执行受限制的操作
    }

    // 管理员重置用户限制
    function resetUserLimits(address user) public {
        // 添加适当的访问控制
        actionCount[user] = 0;
        dailyActions[user] = 0;
    }
}
```

### 9.2.4 紧急停止模式

为合约添加紧急停止功能。

```solidity
contract EmergencyStop {
    bool public stopped = false;
    address public owner;
    mapping(address => bool) public authorizedToStop;

    event EmergencyStopActivated(address indexed activator);
    event EmergencyStopDeactivated(address indexed deactivator);

    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner");
        _;
    }

    modifier onlyAuthorized() {
        require(
            msg.sender == owner || authorizedToStop[msg.sender],
            "Not authorized"
        );
        _;
    }

    modifier stopInEmergency() {
        require(!stopped, "Contract is stopped");
        _;
    }

    modifier onlyInEmergency() {
        require(stopped, "Contract is not stopped");
        _;
    }

    constructor() {
        owner = msg.sender;
    }

    function addAuthorizedStopper(address stopper) public onlyOwner {
        authorizedToStop[stopper] = true;
    }

    function removeAuthorizedStopper(address stopper) public onlyOwner {
        authorizedToStop[stopper] = false;
    }

    function emergencyStop() public onlyAuthorized {
        stopped = true;
        emit EmergencyStopActivated(msg.sender);
    }

    function resumeContract() public onlyOwner {
        stopped = false;
        emit EmergencyStopDeactivated(msg.sender);
    }

    // 正常操作函数
    function normalFunction() public stopInEmergency {
        // 正常业务逻辑
    }

    // 紧急情况下的函数
    function emergencyWithdraw() public onlyInEmergency {
        // 紧急提取逻辑
        (bool success, ) = owner.call{value: address(this).balance}("");
        require(success, "Emergency withdraw failed");
    }
}
```

## 9.3 输入验证和数据清理

### 9.3.1 全面的输入验证

```solidity
contract InputValidation {
    struct User {
        string name;
        uint256 age;
        address wallet;
        bool isActive;
    }

    mapping(address => User) public users;
    mapping(string => address) public nameToAddress;

    event UserRegistered(address indexed user, string name);

    function registerUser(
        string calldata name,
        uint256 age,
        address wallet
    ) external {
        // 验证字符串长度
        require(bytes(name).length > 0, "Name cannot be empty");
        require(bytes(name).length <= 50, "Name too long");

        // 验证字符串内容
        require(isValidName(name), "Invalid characters in name");

        // 验证数值范围
        require(age >= 18 && age <= 120, "Invalid age");

        // 验证地址
        require(wallet != address(0), "Invalid wallet address");
        require(wallet != msg.sender, "Wallet cannot be sender");

        // 验证唯一性
        require(nameToAddress[name] == address(0), "Name already taken");
        require(bytes(users[msg.sender].name).length == 0, "User already registered");

        // 存储用户信息
        users[msg.sender] = User({
            name: name,
            age: age,
            wallet: wallet,
            isActive: true
        });

        nameToAddress[name] = msg.sender;

        emit UserRegistered(msg.sender, name);
    }

    function isValidName(string calldata name) public pure returns (bool) {
        bytes memory nameBytes = bytes(name);

        for (uint256 i = 0; i < nameBytes.length; i++) {
            bytes1 char = nameBytes[i];

            // 只允许字母、数字、空格和某些特殊字符
            if (!(
                (char >= 0x30 && char <= 0x39) || // 0-9
                (char >= 0x41 && char <= 0x5A) || // A-Z
                (char >= 0x61 && char <= 0x7A) || // a-z
                char == 0x20 || // space
                char == 0x2D || // hyphen
                char == 0x5F    // underscore
            )) {
                return false;
            }
        }

        return true;
    }

    // 验证数组输入
    function batchOperation(
        address[] calldata addresses,
        uint256[] calldata amounts
    ) external {
        require(addresses.length > 0, "Empty arrays");
        require(addresses.length == amounts.length, "Array length mismatch");
        require(addresses.length <= 100, "Too many items"); // 防止 DoS

        for (uint256 i = 0; i < addresses.length; i++) {
            require(addresses[i] != address(0), "Invalid address in array");
            require(amounts[i] > 0, "Invalid amount in array");

            // 检查重复地址
            for (uint256 j = i + 1; j < addresses.length; j++) {
                require(addresses[i] != addresses[j], "Duplicate address");
            }
        }

        // 执行批量操作
        for (uint256 i = 0; i < addresses.length; i++) {
            // 批量操作逻辑
        }
    }
}
```

### 9.3.2 防止溢出和边界检查

```solidity
contract SafeCalculations {
    using SafeMath for uint256;

    uint256 public constant MAX_SUPPLY = 1000000 * 10**18;
    uint256 public constant MIN_TRANSFER = 1000; // 最小转账金额

    function safeAdd(uint256 a, uint256 b) public pure returns (uint256) {
        // Solidity 0.8.0+ 自动检查溢出
        return a + b;
    }

    function safeSub(uint256 a, uint256 b) public pure returns (uint256) {
        require(a >= b, "Subtraction underflow");
        return a - b;
    }

    function safeMul(uint256 a, uint256 b) public pure returns (uint256) {
        if (a == 0) return 0;
        uint256 c = a * b;
        require(c / a == b, "Multiplication overflow");
        return c;
    }

    function safeDiv(uint256 a, uint256 b) public pure returns (uint256) {
        require(b > 0, "Division by zero");
        return a / b;
    }

    // 安全的百分比计算
    function calculatePercentage(
        uint256 amount,
        uint256 percentage
    ) public pure returns (uint256) {
        require(percentage <= 10000, "Percentage too high"); // 最大 100%
        require(amount > 0, "Amount must be positive");

        return (amount * percentage) / 10000;
    }

    // 复合利息计算（防止溢出）
    function calculateCompoundInterest(
        uint256 principal,
        uint256 rate, // 基点（10000 = 100%）
        uint256 periods
    ) public pure returns (uint256) {
        require(principal > 0, "Principal must be positive");
        require(rate > 0 && rate <= 10000, "Invalid rate");
        require(periods > 0 && periods <= 100, "Invalid periods");

        uint256 result = principal;

        for (uint256 i = 0; i < periods; i++) {
            uint256 interest = (result * rate) / 10000;

            // 检查是否会溢出
            require(result <= type(uint256).max - interest, "Calculation overflow");

            result += interest;
        }

        return result;
    }
}

// 0.8.0 之前使用的 SafeMath 库
library SafeMath {
    function add(uint256 a, uint256 b) internal pure returns (uint256) {
        uint256 c = a + b;
        require(c >= a, "SafeMath: addition overflow");
        return c;
    }

    function sub(uint256 a, uint256 b) internal pure returns (uint256) {
        require(b <= a, "SafeMath: subtraction overflow");
        return a - b;
    }

    function mul(uint256 a, uint256 b) internal pure returns (uint256) {
        if (a == 0) return 0;
        uint256 c = a * b;
        require(c / a == b, "SafeMath: multiplication overflow");
        return c;
    }

    function div(uint256 a, uint256 b) internal pure returns (uint256) {
        require(b > 0, "SafeMath: division by zero");
        return a / b;
    }

    function mod(uint256 a, uint256 b) internal pure returns (uint256) {
        require(b > 0, "SafeMath: modulo by zero");
        return a % b;
    }
}
```

## 9.4 Gas 优化与 DoS 防护

### 9.4.1 防止 Gas 限制 DoS 攻击

```solidity
contract GasOptimizedContract {
    struct User {
        address addr;
        uint256 balance;
        bool isActive;
    }

    User[] public users;
    mapping(address => uint256) public userIndex;

    // 不好的实现：可能导致 Gas 限制 DoS
    function badGetTotalBalance() public view returns (uint256) {
        uint256 total = 0;

        // 如果用户数量很大，这个循环可能会超出 Gas 限制
        for (uint256 i = 0; i < users.length; i++) {
            if (users[i].isActive) {
                total += users[i].balance;
            }
        }

        return total;
    }

    // 好的实现：分页处理
    function getTotalBalance(
        uint256 start,
        uint256 end
    ) public view returns (uint256 total, bool hasMore) {
        require(start < users.length, "Start index out of bounds");

        uint256 actualEnd = end > users.length ? users.length : end;
        require(start < actualEnd, "Invalid range");

        for (uint256 i = start; i < actualEnd; i++) {
            if (users[i].isActive) {
                total += users[i].balance;
            }
        }

        hasMore = actualEnd < users.length;
    }

    // 使用映射缓存总计
    uint256 public totalActiveBalance;

    function updateBalance(address user, uint256 newBalance) public {
        uint256 index = userIndex[user];
        require(index < users.length && users[index].addr == user, "User not found");

        if (users[index].isActive) {
            totalActiveBalance = totalActiveBalance - users[index].balance + newBalance;
        }

        users[index].balance = newBalance;
    }

    // 批量操作优化
    function batchUpdateBalances(
        address[] calldata addresses,
        uint256[] calldata balances
    ) external {
        require(addresses.length == balances.length, "Array length mismatch");
        require(addresses.length <= 50, "Batch too large"); // 限制批量大小

        for (uint256 i = 0; i < addresses.length; i++) {
            updateBalance(addresses[i], balances[i]);
        }
    }
}
```

### 9.4.2 循环优化技巧

```solidity
contract LoopOptimization {
    uint256[] public data;
    mapping(uint256 => bool) public exists;

    // 不优化的循环
    function inefficientSearch(uint256 target) public view returns (bool found, uint256 index) {
        for (uint256 i = 0; i < data.length; i++) { // 每次读取 data.length
            if (data[i] == target) {
                return (true, i);
            }
        }
        return (false, 0);
    }

    // 优化的循环
    function efficientSearch(uint256 target) public view returns (bool found, uint256 index) {
        uint256 length = data.length; // 缓存长度

        for (uint256 i = 0; i < length;) { // 避免边界检查
            if (data[i] == target) {
                return (true, i);
            }

            unchecked { ++i; } // 使用 unchecked 避免溢出检查
        }

        return (false, 0);
    }

    // 使用映射替代数组搜索
    function addData(uint256 value) public {
        require(!exists[value], "Value already exists");

        data.push(value);
        exists[value] = true;
    }

    function checkExists(uint256 value) public view returns (bool) {
        return exists[value]; // O(1) 而不是 O(n)
    }

    // 逆序循环优化（在某些情况下更有效）
    function reverseIteration() public view returns (uint256 sum) {
        uint256 length = data.length;

        if (length == 0) return 0;

        uint256 i = length;
        do {
            unchecked { --i; }
            sum += data[i];
        } while (i > 0);
    }
}
```

## 9.5 随机数安全

### 9.5.1 安全的随机数生成

```solidity
// 不安全的随机数生成
contract InsecureRandomness {
    function badRandom() public view returns (uint256) {
        // 危险！可以被矿工操控
        return uint256(keccak256(abi.encodePacked(block.timestamp, block.difficulty)));
    }

    function anotherBadRandom() public view returns (uint256) {
        // 危险！可以被预测
        return uint256(keccak256(abi.encodePacked(blockhash(block.number - 1))));
    }
}

// 更安全的随机数生成
contract SecureRandomness {
    uint256 private nonce = 0;
    mapping(address => uint256) private userNonces;

    // 承诺-揭示方案
    struct Commitment {
        bytes32 commitment;
        uint256 blockNumber;
        bool revealed;
    }

    mapping(address => Commitment) public commitments;

    function commitRandomness(bytes32 commitment) public {
        commitments[msg.sender] = Commitment({
            commitment: commitment,
            blockNumber: block.number,
            revealed: false
        });
    }

    function revealRandomness(uint256 randomValue, uint256 salt) public returns (uint256) {
        Commitment storage commitment = commitments[msg.sender];

        require(commitment.blockNumber > 0, "No commitment found");
        require(!commitment.revealed, "Already revealed");
        require(
            block.number > commitment.blockNumber + 1,
            "Must wait at least one block"
        );

        bytes32 hash = keccak256(abi.encodePacked(randomValue, salt, msg.sender));
        require(hash == commitment.commitment, "Invalid reveal");

        commitment.revealed = true;

        // 结合多个难以操控的值
        uint256 combinedRandomness = uint256(keccak256(abi.encodePacked(
            randomValue,
            blockhash(commitment.blockNumber + 1),
            block.coinbase,
            nonce++
        )));

        return combinedRandomness;
    }

    // 使用外部随机性源（如 Chainlink VRF）
    interface VRFConsumerBase {
        function requestRandomness(bytes32 keyHash, uint256 fee) external returns (bytes32 requestId);
    }

    // 为简单用例的弱随机数（不用于关键决策）
    function weakRandom() internal returns (uint256) {
        nonce++;
        userNonces[msg.sender]++;

        return uint256(keccak256(abi.encodePacked(
            block.timestamp,
            block.difficulty,
            msg.sender,
            nonce,
            userNonces[msg.sender],
            blockhash(block.number - 1)
        )));
    }
}
```

## 9.6 审计清单

### 9.6.1 代码审计检查清单

```solidity
/**
 * 智能合约安全审计清单
 *
 * 1. 重入攻击防护
 *    - [ ] 使用检查-效果-交互模式
 *    - [ ] 使用重入锁
 *    - [ ] 避免在状态更新前进行外部调用
 *
 * 2. 整数溢出/下溢
 *    - [ ] 使用 Solidity 0.8.0+ 的内置检查
 *    - [ ] 或使用 SafeMath 库
 *    - [ ] 谨慎使用 unchecked 块
 *
 * 3. 访问控制
 *    - [ ] 所有敏感函数都有适当的修饰符
 *    - [ ] 使用 msg.sender 而不是 tx.origin
 *    - [ ] 实现多步骤所有权转移
 *
 * 4. 输入验证
 *    - [ ] 验证所有外部输入
 *    - [ ] 检查数组长度
 *    - [ ] 验证地址不为零
 *
 * 5. Gas 优化和 DoS 防护
 *    - [ ] 避免无界循环
 *    - [ ] 实现分页功能
 *    - [ ] 限制批量操作大小
 *
 * 6. 时间依赖
 *    - [ ] 不依赖 block.timestamp 进行关键逻辑
 *    - [ ] 使用区块号而不是时间戳
 *
 * 7. 随机数
 *    - [ ] 不使用可预测的随机数源
 *    - [ ] 实现承诺-揭示方案
 *    - [ ] 考虑使用外部随机性源
 *
 * 8. 外部调用
 *    - [ ] 检查外部调用的返回值
 *    - [ ] 使用 address.call 而不是废弃的方法
 *    - [ ] 实现拉取支付模式
 *
 * 9. 紧急处理
 *    - [ ] 实现紧急停止功能
 *    - [ ] 提供升级或迁移机制
 *
 * 10. 事件和日志
 *     - [ ] 为重要操作添加事件
 *     - [ ] 使用 indexed 参数便于查询
 */

contract AuditedContract {
    // 示例：实现了安全检查的合约

    address public owner;
    bool public paused = false;
    bool private locked = false;

    mapping(address => uint256) public balances;

    event Transfer(address indexed from, address indexed to, uint256 value);
    event Paused();
    event Unpaused();
    event OwnershipTransferred(address indexed previousOwner, address indexed newOwner);

    modifier onlyOwner() {
        require(msg.sender == owner, "Caller is not the owner");
        _;
    }

    modifier whenNotPaused() {
        require(!paused, "Contract is paused");
        _;
    }

    modifier noReentrant() {
        require(!locked, "ReentrancyGuard: reentrant call");
        locked = true;
        _;
        locked = false;
    }

    constructor() {
        owner = msg.sender;
    }

    function transfer(address to, uint256 amount)
        public
        whenNotPaused
        noReentrant
        returns (bool)
    {
        // 1. 检查
        require(to != address(0), "Transfer to zero address");
        require(amount > 0, "Transfer amount must be positive");
        require(balances[msg.sender] >= amount, "Insufficient balance");

        // 2. 效果
        balances[msg.sender] -= amount;
        balances[to] += amount;

        // 3. 交互
        emit Transfer(msg.sender, to, amount);

        return true;
    }

    function pause() public onlyOwner {
        paused = true;
        emit Paused();
    }

    function unpause() public onlyOwner {
        paused = false;
        emit Unpaused();
    }

    function transferOwnership(address newOwner) public onlyOwner {
        require(newOwner != address(0), "New owner is the zero address");
        emit OwnershipTransferred(owner, newOwner);
        owner = newOwner;
    }
}
```

## 9.7 实战练习

### 练习 1：修复漏洞合约

找出并修复以下合约的安全漏洞：

```solidity
contract VulnerableContract {
    mapping(address => uint256) public balances;
    address public owner;

    constructor() {
        owner = msg.sender;
    }

    function deposit() public payable {
        balances[msg.sender] += msg.value;
    }

    function withdraw(uint256 amount) public {
        require(balances[msg.sender] >= amount);
        msg.sender.call{value: amount}("");
        balances[msg.sender] -= amount;
    }

    function emergencyWithdraw() public {
        require(tx.origin == owner);
        payable(owner).transfer(address(this).balance);
    }
}
```

### 练习 2：实现安全的拍卖系统

创建一个防止各种攻击的拍卖合约。

### 练习 3：设计安全的投票系统

实现一个防止重入、前置交易等攻击的投票系统。

## 9.8 章节总结

在本章中，我们全面学习了：

1. **常见漏洞**：重入攻击、整数溢出、访问控制等
2. **安全模式**：检查-效果-交互、拉取支付等
3. **输入验证**：全面的数据验证和边界检查
4. **Gas 优化**：防止 DoS 攻击的优化技巧
5. **随机数安全**：安全的随机数生成方法
6. **审计实践**：系统性的安全检查方法

### 关键原则

- 安全第一，性能其次
- 假设外部输入都是恶意的
- 使用经过验证的安全模式
- 实施深度防御策略
- 定期进行安全审计

安全是一个持续的过程，需要在整个开发生命周期中保持警惕。

继续您的 Solidity 学习之旅：[第十章 - Gas 优化与性能](10_gas_optimization.md)
