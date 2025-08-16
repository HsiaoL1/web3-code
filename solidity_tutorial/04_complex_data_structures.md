# 第四章：复杂数据结构

## 本章学习目标

- 掌握数组的创建、操作和优化技巧
- 深入理解映射的使用和嵌套映射
- 学会设计和使用结构体
- 了解枚举的高级应用
- 掌握数据结构的组合使用
- 学习存储优化策略

## 4.1 数组详解

数组是存储相同类型元素集合的数据结构，分为定长数组和动态数组。

### 4.1.1 数组的基本操作

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract ArrayBasics {
    // 动态数组声明
    uint256[] public dynamicArray;
    string[] public names;

    // 定长数组声明
    uint256[5] public fixedArray;
    bool[3] public flags = [true, false, true];

    // 二维数组
    uint256[][] public matrix;

    constructor() {
        // 初始化定长数组
        fixedArray = [1, 2, 3, 4, 5];

        // 初始化动态数组
        dynamicArray.push(10);
        dynamicArray.push(20);
        dynamicArray.push(30);
    }

    // 添加元素
    function addElement(uint256 _value) public {
        dynamicArray.push(_value);
    }

    // 删除最后一个元素
    function removeLastElement() public {
        require(dynamicArray.length > 0, "Array is empty");
        dynamicArray.pop();
    }

    // 获取数组长度
    function getArrayLength() public view returns (uint256) {
        return dynamicArray.length;
    }

    // 获取指定位置的元素
    function getElement(uint256 _index) public view returns (uint256) {
        require(_index < dynamicArray.length, "Index out of bounds");
        return dynamicArray[_index];
    }

    // 更新指定位置的元素
    function updateElement(uint256 _index, uint256 _value) public {
        require(_index < dynamicArray.length, "Index out of bounds");
        dynamicArray[_index] = _value;
    }

    // 清空整个数组
    function clearArray() public {
        delete dynamicArray;
    }

    // 删除指定位置的元素（保持顺序）
    function removeElement(uint256 _index) public {
        require(_index < dynamicArray.length, "Index out of bounds");

        // 将后面的元素向前移动
        for (uint256 i = _index; i < dynamicArray.length - 1; i++) {
            dynamicArray[i] = dynamicArray[i + 1];
        }
        dynamicArray.pop();
    }

    // 删除指定位置的元素（不保持顺序，更省 Gas）
    function removeElementUnordered(uint256 _index) public {
        require(_index < dynamicArray.length, "Index out of bounds");

        // 用最后一个元素替换要删除的元素
        dynamicArray[_index] = dynamicArray[dynamicArray.length - 1];
        dynamicArray.pop();
    }
}
```

### 4.1.2 数组的高级操作

```solidity
contract AdvancedArrays {
    uint256[] public numbers;
    address[] public users;

    // 批量添加元素
    function batchAdd(uint256[] calldata _values) external {
        for (uint256 i = 0; i < _values.length; i++) {
            numbers.push(_values[i]);
        }
    }

    // 查找元素
    function findElement(uint256 _value) public view returns (bool found, uint256 index) {
        for (uint256 i = 0; i < numbers.length; i++) {
            if (numbers[i] == _value) {
                return (true, i);
            }
        }
        return (false, 0);
    }

    // 数组排序（冒泡排序 - 仅适用于小数组）
    function bubbleSort() public {
        uint256 length = numbers.length;
        for (uint256 i = 0; i < length - 1; i++) {
            for (uint256 j = 0; j < length - i - 1; j++) {
                if (numbers[j] > numbers[j + 1]) {
                    // 交换元素
                    uint256 temp = numbers[j];
                    numbers[j] = numbers[j + 1];
                    numbers[j + 1] = temp;
                }
            }
        }
    }

    // 计算数组统计信息
    function getStatistics() public view returns (uint256 sum, uint256 average, uint256 max, uint256 min) {
        require(numbers.length > 0, "Array is empty");

        sum = numbers[0];
        max = numbers[0];
        min = numbers[0];

        for (uint256 i = 1; i < numbers.length; i++) {
            sum += numbers[i];
            if (numbers[i] > max) {
                max = numbers[i];
            }
            if (numbers[i] < min) {
                min = numbers[i];
            }
        }

        average = sum / numbers.length;
    }

    // 筛选元素
    function filterEvenNumbers() public view returns (uint256[] memory) {
        // 首先计算符合条件的元素数量
        uint256 count = 0;
        for (uint256 i = 0; i < numbers.length; i++) {
            if (numbers[i] % 2 == 0) {
                count++;
            }
        }

        // 创建结果数组并填充
        uint256[] memory evenNumbers = new uint256[](count);
        uint256 index = 0;
        for (uint256 i = 0; i < numbers.length; i++) {
            if (numbers[i] % 2 == 0) {
                evenNumbers[index] = numbers[i];
                index++;
            }
        }

        return evenNumbers;
    }

    // 数组去重
    function removeDuplicates() public {
        if (numbers.length <= 1) return;

        uint256 writeIndex = 1;
        for (uint256 readIndex = 1; readIndex < numbers.length; readIndex++) {
            bool isDuplicate = false;
            for (uint256 checkIndex = 0; checkIndex < writeIndex; checkIndex++) {
                if (numbers[readIndex] == numbers[checkIndex]) {
                    isDuplicate = true;
                    break;
                }
            }
            if (!isDuplicate) {
                numbers[writeIndex] = numbers[readIndex];
                writeIndex++;
            }
        }

        // 调整数组大小
        while (numbers.length > writeIndex) {
            numbers.pop();
        }
    }
}
```

## 4.2 映射详解

映射是 Solidity 中最重要的数据结构之一，类似于哈希表或字典。

### 4.2.1 基本映射操作

```solidity
contract MappingBasics {
    // 基本映射
    mapping(address => uint256) public balances;
    mapping(string => bool) public isValidName;
    mapping(uint256 => address) public idToAddress;

    // 嵌套映射
    mapping(address => mapping(address => uint256)) public allowances;
    mapping(address => mapping(string => bool)) public userPermissions;

    uint256 public nextId = 1;

    // 设置余额
    function setBalance(address _user, uint256 _amount) public {
        balances[_user] = _amount;
    }

    // 增加余额
    function addBalance(address _user, uint256 _amount) public {
        balances[_user] += _amount;
    }

    // 检查余额是否足够
    function hasEnoughBalance(address _user, uint256 _amount) public view returns (bool) {
        return balances[_user] >= _amount;
    }

    // 注册用户
    function registerUser(address _user, string memory _name) public {
        require(!isValidName[_name], "Name already taken");
        require(idToAddress[nextId] == address(0), "ID already assigned");

        idToAddress[nextId] = _user;
        isValidName[_name] = true;
        nextId++;
    }

    // 授权额度管理
    function approve(address _spender, uint256 _amount) public {
        allowances[msg.sender][_spender] = _amount;
    }

    // 获取授权额度
    function getAllowance(address _owner, address _spender) public view returns (uint256) {
        return allowances[_owner][_spender];
    }

    // 权限管理
    function grantPermission(address _user, string memory _permission) public {
        userPermissions[_user][_permission] = true;
    }

    function revokePermission(address _user, string memory _permission) public {
        userPermissions[_user][_permission] = false;
    }

    function hasPermission(address _user, string memory _permission) public view returns (bool) {
        return userPermissions[_user][_permission];
    }
}
```

### 4.2.2 映射的高级用法

```solidity
contract AdvancedMappings {
    struct User {
        string name;
        uint256 balance;
        bool isActive;
        uint256[] ownedTokens;
    }

    // 映射到结构体
    mapping(address => User) public users;

    // 映射到数组
    mapping(address => uint256[]) public userTransactions;
    mapping(string => address[]) public categoryUsers;

    // 复杂的嵌套映射
    mapping(address => mapping(uint256 => mapping(string => bool))) public complexMapping;

    event UserRegistered(address indexed user, string name);
    event TransactionAdded(address indexed user, uint256 amount);

    // 注册用户
    function registerUser(string memory _name) public {
        require(bytes(_name).length > 0, "Name cannot be empty");
        require(bytes(users[msg.sender].name).length == 0, "User already registered");

        users[msg.sender] = User({
            name: _name,
            balance: 0,
            isActive: true,
            ownedTokens: new uint256[](0)
        });

        emit UserRegistered(msg.sender, _name);
    }

    // 添加交易记录
    function addTransaction(uint256 _amount) public {
        require(users[msg.sender].isActive, "User is not active");

        userTransactions[msg.sender].push(_amount);
        users[msg.sender].balance += _amount;

        emit TransactionAdded(msg.sender, _amount);
    }

    // 获取用户的交易数量
    function getUserTransactionCount(address _user) public view returns (uint256) {
        return userTransactions[_user].length;
    }

    // 获取用户的所有交易
    function getUserTransactions(address _user) public view returns (uint256[] memory) {
        return userTransactions[_user];
    }

    // 添加代币到用户
    function addTokenToUser(address _user, uint256 _tokenId) public {
        users[_user].ownedTokens.push(_tokenId);
    }

    // 获取用户拥有的代币
    function getUserTokens(address _user) public view returns (uint256[] memory) {
        return users[_user].ownedTokens;
    }

    // 将用户添加到类别
    function addUserToCategory(address _user, string memory _category) public {
        categoryUsers[_category].push(_user);
    }

    // 获取某类别的所有用户
    function getCategoryUsers(string memory _category) public view returns (address[] memory) {
        return categoryUsers[_category];
    }

    // 复杂映射操作
    function setComplexData(uint256 _id, string memory _key, bool _value) public {
        complexMapping[msg.sender][_id][_key] = _value;
    }

    function getComplexData(address _user, uint256 _id, string memory _key) public view returns (bool) {
        return complexMapping[_user][_id][_key];
    }
}
```

## 4.3 结构体详解

结构体允许创建自定义的复合数据类型。

### 4.3.1 结构体的基本用法

```solidity
contract StructBasics {
    // 定义结构体
    struct Product {
        uint256 id;
        string name;
        uint256 price;
        bool isAvailable;
        address seller;
        string[] tags;
    }

    struct Order {
        uint256 orderId;
        address buyer;
        uint256 productId;
        uint256 quantity;
        uint256 totalPrice;
        OrderStatus status;
        uint256 timestamp;
    }

    enum OrderStatus {
        Pending,
        Confirmed,
        Shipped,
        Delivered,
        Cancelled
    }

    // 存储结构体
    Product[] public products;
    Order[] public orders;
    mapping(uint256 => Product) public productById;
    mapping(address => Order[]) public userOrders;

    uint256 public nextProductId = 1;
    uint256 public nextOrderId = 1;

    event ProductAdded(uint256 indexed productId, string name, uint256 price);
    event OrderCreated(uint256 indexed orderId, address indexed buyer, uint256 productId);

    // 添加产品
    function addProduct(
        string memory _name,
        uint256 _price,
        string[] memory _tags
    ) public returns (uint256) {
        uint256 productId = nextProductId++;

        Product memory newProduct = Product({
            id: productId,
            name: _name,
            price: _price,
            isAvailable: true,
            seller: msg.sender,
            tags: _tags
        });

        products.push(newProduct);
        productById[productId] = newProduct;

        emit ProductAdded(productId, _name, _price);
        return productId;
    }

    // 创建订单
    function createOrder(uint256 _productId, uint256 _quantity) public payable {
        require(productById[_productId].isAvailable, "Product not available");
        require(_quantity > 0, "Quantity must be positive");

        uint256 totalPrice = productById[_productId].price * _quantity;
        require(msg.value >= totalPrice, "Insufficient payment");

        uint256 orderId = nextOrderId++;
        Order memory newOrder = Order({
            orderId: orderId,
            buyer: msg.sender,
            productId: _productId,
            quantity: _quantity,
            totalPrice: totalPrice,
            status: OrderStatus.Pending,
            timestamp: block.timestamp
        });

        orders.push(newOrder);
        userOrders[msg.sender].push(newOrder);

        emit OrderCreated(orderId, msg.sender, _productId);
    }

    // 更新产品信息
    function updateProduct(
        uint256 _productId,
        string memory _name,
        uint256 _price,
        bool _isAvailable
    ) public {
        require(productById[_productId].seller == msg.sender, "Only seller can update");

        productById[_productId].name = _name;
        productById[_productId].price = _price;
        productById[_productId].isAvailable = _isAvailable;

        // 同时更新数组中的产品
        for (uint256 i = 0; i < products.length; i++) {
            if (products[i].id == _productId) {
                products[i] = productById[_productId];
                break;
            }
        }
    }

    // 获取产品信息
    function getProduct(uint256 _productId) public view returns (Product memory) {
        return productById[_productId];
    }

    // 获取用户的所有订单
    function getUserOrders(address _user) public view returns (Order[] memory) {
        return userOrders[_user];
    }

    // 获取产品数量
    function getProductCount() public view returns (uint256) {
        return products.length;
    }
}
```

### 4.3.2 结构体的高级应用

```solidity
contract AdvancedStructs {
    // 嵌套结构体
    struct Address {
        string street;
        string city;
        string country;
        string zipCode;
    }

    struct ContactInfo {
        string email;
        string phone;
        Address physicalAddress;
    }

    struct Company {
        string name;
        ContactInfo contact;
        address[] employees;
        mapping(address => bool) isEmployee;
        uint256 foundedYear;
    }

    // 结构体数组和映射
    Company[] public companies;
    mapping(string => uint256) public companyNameToId;
    mapping(address => uint256[]) public employeeCompanies;

    // 创建公司
    function createCompany(
        string memory _name,
        string memory _email,
        string memory _phone,
        string memory _street,
        string memory _city,
        string memory _country,
        string memory _zipCode,
        uint256 _foundedYear
    ) public returns (uint256) {
        require(companyNameToId[_name] == 0, "Company name already exists");

        uint256 companyId = companies.length;
        companies.push();

        Company storage newCompany = companies[companyId];
        newCompany.name = _name;
        newCompany.contact.email = _email;
        newCompany.contact.phone = _phone;
        newCompany.contact.physicalAddress = Address(_street, _city, _country, _zipCode);
        newCompany.foundedYear = _foundedYear;

        companyNameToId[_name] = companyId + 1; // +1 to avoid 0

        return companyId;
    }

    // 添加员工
    function addEmployee(uint256 _companyId, address _employee) public {
        require(_companyId < companies.length, "Company does not exist");
        require(!companies[_companyId].isEmployee[_employee], "Already an employee");

        companies[_companyId].employees.push(_employee);
        companies[_companyId].isEmployee[_employee] = true;
        employeeCompanies[_employee].push(_companyId);
    }

    // 获取公司信息
    function getCompanyInfo(uint256 _companyId)
        public
        view
        returns (
            string memory name,
            string memory email,
            string memory phone,
            Address memory addr,
            uint256 foundedYear,
            uint256 employeeCount
        )
    {
        require(_companyId < companies.length, "Company does not exist");

        Company storage company = companies[_companyId];
        return (
            company.name,
            company.contact.email,
            company.contact.phone,
            company.contact.physicalAddress,
            company.foundedYear,
            company.employees.length
        );
    }

    // 检查是否为员工
    function isEmployeeOf(uint256 _companyId, address _employee) public view returns (bool) {
        require(_companyId < companies.length, "Company does not exist");
        return companies[_companyId].isEmployee[_employee];
    }

    // 获取员工所属的公司
    function getEmployeeCompanies(address _employee) public view returns (uint256[] memory) {
        return employeeCompanies[_employee];
    }
}
```

## 4.4 枚举详解

枚举用于创建具有有限可能值的用户定义类型。

### 4.4.1 枚举的基本用法

```solidity
contract EnumBasics {
    // 定义枚举
    enum TaskStatus {
        Created,    // 0
        InProgress, // 1
        Completed,  // 2
        Cancelled   // 3
    }

    enum Priority {
        Low,
        Medium,
        High,
        Critical
    }

    struct Task {
        string title;
        string description;
        TaskStatus status;
        Priority priority;
        address assignee;
        uint256 deadline;
        uint256 createdAt;
    }

    Task[] public tasks;
    mapping(address => uint256[]) public userTasks;

    event TaskCreated(uint256 indexed taskId, address indexed assignee);
    event TaskStatusChanged(uint256 indexed taskId, TaskStatus oldStatus, TaskStatus newStatus);

    // 创建任务
    function createTask(
        string memory _title,
        string memory _description,
        Priority _priority,
        address _assignee,
        uint256 _deadline
    ) public returns (uint256) {
        uint256 taskId = tasks.length;

        tasks.push(Task({
            title: _title,
            description: _description,
            status: TaskStatus.Created,
            priority: _priority,
            assignee: _assignee,
            deadline: _deadline,
            createdAt: block.timestamp
        }));

        userTasks[_assignee].push(taskId);

        emit TaskCreated(taskId, _assignee);
        return taskId;
    }

    // 更新任务状态
    function updateTaskStatus(uint256 _taskId, TaskStatus _newStatus) public {
        require(_taskId < tasks.length, "Task does not exist");
        require(
            tasks[_taskId].assignee == msg.sender,
            "Only assignee can update status"
        );

        TaskStatus oldStatus = tasks[_taskId].status;
        tasks[_taskId].status = _newStatus;

        emit TaskStatusChanged(_taskId, oldStatus, _newStatus);
    }

    // 获取任务状态
    function getTaskStatus(uint256 _taskId) public view returns (TaskStatus) {
        require(_taskId < tasks.length, "Task does not exist");
        return tasks[_taskId].status;
    }

    // 开始任务
    function startTask(uint256 _taskId) public {
        require(_taskId < tasks.length, "Task does not exist");
        require(tasks[_taskId].assignee == msg.sender, "Only assignee can start task");
        require(tasks[_taskId].status == TaskStatus.Created, "Task not in created status");

        updateTaskStatus(_taskId, TaskStatus.InProgress);
    }

    // 完成任务
    function completeTask(uint256 _taskId) public {
        require(_taskId < tasks.length, "Task does not exist");
        require(tasks[_taskId].assignee == msg.sender, "Only assignee can complete task");
        require(tasks[_taskId].status == TaskStatus.InProgress, "Task not in progress");

        updateTaskStatus(_taskId, TaskStatus.Completed);
    }

    // 取消任务
    function cancelTask(uint256 _taskId) public {
        require(_taskId < tasks.length, "Task does not exist");
        require(tasks[_taskId].assignee == msg.sender, "Only assignee can cancel task");
        require(
            tasks[_taskId].status == TaskStatus.Created ||
            tasks[_taskId].status == TaskStatus.InProgress,
            "Cannot cancel completed task"
        );

        updateTaskStatus(_taskId, TaskStatus.Cancelled);
    }

    // 获取特定状态的任务数量
    function getTaskCountByStatus(TaskStatus _status) public view returns (uint256) {
        uint256 count = 0;
        for (uint256 i = 0; i < tasks.length; i++) {
            if (tasks[i].status == _status) {
                count++;
            }
        }
        return count;
    }

    // 获取用户的任务
    function getUserTasks(address _user) public view returns (uint256[] memory) {
        return userTasks[_user];
    }
}
```

### 4.4.2 枚举的高级应用

```solidity
contract AdvancedEnums {
    // 权限枚举
    enum Role {
        None,
        User,
        Moderator,
        Admin,
        SuperAdmin
    }

    // 投票选项枚举
    enum VoteOption {
        Abstain,
        For,
        Against
    }

    // 合约状态枚举
    enum ContractState {
        Inactive,
        Active,
        Paused,
        Terminated
    }

    struct Proposal {
        string title;
        string description;
        uint256 deadline;
        mapping(address => VoteOption) votes;
        uint256 forVotes;
        uint256 againstVotes;
        uint256 abstainVotes;
        bool executed;
    }

    mapping(address => Role) public userRoles;
    Proposal[] public proposals;
    ContractState public contractState = ContractState.Inactive;

    modifier onlyRole(Role _minRole) {
        require(userRoles[msg.sender] >= _minRole, "Insufficient permissions");
        _;
    }

    modifier inState(ContractState _state) {
        require(contractState == _state, "Invalid contract state");
        _;
    }

    // 设置用户角色
    function setUserRole(address _user, Role _role) public onlyRole(Role.Admin) {
        userRoles[_user] = _role;
    }

    // 激活合约
    function activateContract() public onlyRole(Role.SuperAdmin) {
        contractState = ContractState.Active;
    }

    // 暂停合约
    function pauseContract() public onlyRole(Role.Admin) inState(ContractState.Active) {
        contractState = ContractState.Paused;
    }

    // 恢复合约
    function resumeContract() public onlyRole(Role.Admin) inState(ContractState.Paused) {
        contractState = ContractState.Active;
    }

    // 创建提案
    function createProposal(
        string memory _title,
        string memory _description,
        uint256 _votingPeriod
    ) public onlyRole(Role.Moderator) inState(ContractState.Active) returns (uint256) {
        uint256 proposalId = proposals.length;
        proposals.push();

        Proposal storage newProposal = proposals[proposalId];
        newProposal.title = _title;
        newProposal.description = _description;
        newProposal.deadline = block.timestamp + _votingPeriod;
        newProposal.executed = false;

        return proposalId;
    }

    // 投票
    function vote(uint256 _proposalId, VoteOption _option)
        public
        onlyRole(Role.User)
        inState(ContractState.Active)
    {
        require(_proposalId < proposals.length, "Proposal does not exist");
        require(block.timestamp <= proposals[_proposalId].deadline, "Voting period ended");
        require(
            proposals[_proposalId].votes[msg.sender] == VoteOption.Abstain,
            "Already voted"
        );

        Proposal storage proposal = proposals[_proposalId];
        proposal.votes[msg.sender] = _option;

        if (_option == VoteOption.For) {
            proposal.forVotes++;
        } else if (_option == VoteOption.Against) {
            proposal.againstVotes++;
        } else {
            proposal.abstainVotes++;
        }
    }

    // 获取提案结果
    function getProposalResults(uint256 _proposalId)
        public
        view
        returns (uint256 forVotes, uint256 againstVotes, uint256 abstainVotes, bool passed)
    {
        require(_proposalId < proposals.length, "Proposal does not exist");

        Proposal storage proposal = proposals[_proposalId];
        forVotes = proposal.forVotes;
        againstVotes = proposal.againstVotes;
        abstainVotes = proposal.abstainVotes;

        // 简单多数决定
        passed = forVotes > againstVotes;
    }

    // 检查角色权限
    function hasPermission(address _user, Role _requiredRole) public view returns (bool) {
        return userRoles[_user] >= _requiredRole;
    }

    // 获取角色名称（用于前端显示）
    function getRoleName(Role _role) public pure returns (string memory) {
        if (_role == Role.None) return "None";
        if (_role == Role.User) return "User";
        if (_role == Role.Moderator) return "Moderator";
        if (_role == Role.Admin) return "Admin";
        if (_role == Role.SuperAdmin) return "SuperAdmin";
        return "Unknown";
    }
}
```

## 4.5 组合数据结构

在实际应用中，经常需要组合使用多种数据结构。

### 4.5.1 复杂系统设计

```solidity
contract MarketplaceSystem {
    // 枚举定义
    enum ProductCategory { Electronics, Clothing, Books, Home, Sports }
    enum OrderStatus { Pending, Paid, Shipped, Delivered, Cancelled, Refunded }
    enum UserType { Buyer, Seller, Admin }

    // 结构体定义
    struct User {
        string name;
        string email;
        UserType userType;
        bool isVerified;
        uint256 reputation;
        uint256 joinDate;
        address[] trustedBy;
    }

    struct Product {
        uint256 id;
        string name;
        string description;
        ProductCategory category;
        uint256 price;
        uint256 quantity;
        address seller;
        string[] images;
        uint256 rating;
        uint256 reviewCount;
        bool isActive;
    }

    struct Order {
        uint256 id;
        uint256 productId;
        address buyer;
        address seller;
        uint256 quantity;
        uint256 totalPrice;
        OrderStatus status;
        uint256 orderDate;
        uint256 deliveryDate;
        string shippingAddress;
    }

    struct Review {
        uint256 orderId;
        address reviewer;
        uint256 rating;
        string comment;
        uint256 timestamp;
    }

    // 状态变量
    mapping(address => User) public users;
    mapping(uint256 => Product) public products;
    mapping(uint256 => Order) public orders;
    mapping(uint256 => Review[]) public productReviews;

    // 分类映射
    mapping(ProductCategory => uint256[]) public productsByCategory;
    mapping(address => uint256[]) public userProducts;
    mapping(address => uint256[]) public buyerOrders;
    mapping(address => uint256[]) public sellerOrders;

    // 计数器
    uint256 public nextProductId = 1;
    uint256 public nextOrderId = 1;

    // 事件
    event UserRegistered(address indexed user, UserType userType);
    event ProductListed(uint256 indexed productId, address indexed seller);
    event OrderCreated(uint256 indexed orderId, address indexed buyer, address indexed seller);
    event OrderStatusChanged(uint256 indexed orderId, OrderStatus newStatus);
    event ReviewAdded(uint256 indexed productId, address indexed reviewer, uint256 rating);

    // 用户注册
    function registerUser(
        string memory _name,
        string memory _email,
        UserType _userType
    ) public {
        require(bytes(users[msg.sender].name).length == 0, "User already registered");

        users[msg.sender] = User({
            name: _name,
            email: _email,
            userType: _userType,
            isVerified: false,
            reputation: 100,
            joinDate: block.timestamp,
            trustedBy: new address[](0)
        });

        emit UserRegistered(msg.sender, _userType);
    }

    // 列出产品
    function listProduct(
        string memory _name,
        string memory _description,
        ProductCategory _category,
        uint256 _price,
        uint256 _quantity,
        string[] memory _images
    ) public returns (uint256) {
        require(users[msg.sender].userType == UserType.Seller, "Only sellers can list products");
        require(_price > 0, "Price must be positive");
        require(_quantity > 0, "Quantity must be positive");

        uint256 productId = nextProductId++;

        products[productId] = Product({
            id: productId,
            name: _name,
            description: _description,
            category: _category,
            price: _price,
            quantity: _quantity,
            seller: msg.sender,
            images: _images,
            rating: 0,
            reviewCount: 0,
            isActive: true
        });

        productsByCategory[_category].push(productId);
        userProducts[msg.sender].push(productId);

        emit ProductListed(productId, msg.sender);
        return productId;
    }

    // 创建订单
    function createOrder(
        uint256 _productId,
        uint256 _quantity,
        string memory _shippingAddress
    ) public payable returns (uint256) {
        require(products[_productId].isActive, "Product not active");
        require(products[_productId].quantity >= _quantity, "Insufficient quantity");
        require(_quantity > 0, "Quantity must be positive");

        uint256 totalPrice = products[_productId].price * _quantity;
        require(msg.value >= totalPrice, "Insufficient payment");

        uint256 orderId = nextOrderId++;

        orders[orderId] = Order({
            id: orderId,
            productId: _productId,
            buyer: msg.sender,
            seller: products[_productId].seller,
            quantity: _quantity,
            totalPrice: totalPrice,
            status: OrderStatus.Pending,
            orderDate: block.timestamp,
            deliveryDate: 0,
            shippingAddress: _shippingAddress
        });

        // 更新产品库存
        products[_productId].quantity -= _quantity;

        buyerOrders[msg.sender].push(orderId);
        sellerOrders[products[_productId].seller].push(orderId);

        emit OrderCreated(orderId, msg.sender, products[_productId].seller);
        return orderId;
    }

    // 更新订单状态
    function updateOrderStatus(uint256 _orderId, OrderStatus _newStatus) public {
        require(orders[_orderId].seller == msg.sender, "Only seller can update order");
        require(orders[_orderId].status != OrderStatus.Delivered, "Order already delivered");

        orders[_orderId].status = _newStatus;

        if (_newStatus == OrderStatus.Delivered) {
            orders[_orderId].deliveryDate = block.timestamp;
        }

        emit OrderStatusChanged(_orderId, _newStatus);
    }

    // 添加评价
    function addReview(
        uint256 _orderId,
        uint256 _rating,
        string memory _comment
    ) public {
        require(orders[_orderId].buyer == msg.sender, "Only buyer can review");
        require(orders[_orderId].status == OrderStatus.Delivered, "Order not delivered yet");
        require(_rating >= 1 && _rating <= 5, "Rating must be between 1 and 5");

        uint256 productId = orders[_orderId].productId;

        productReviews[productId].push(Review({
            orderId: _orderId,
            reviewer: msg.sender,
            rating: _rating,
            comment: _comment,
            timestamp: block.timestamp
        }));

        // 更新产品评分
        updateProductRating(productId);

        emit ReviewAdded(productId, msg.sender, _rating);
    }

    // 更新产品评分
    function updateProductRating(uint256 _productId) internal {
        Review[] storage reviews = productReviews[_productId];
        if (reviews.length == 0) return;

        uint256 totalRating = 0;
        for (uint256 i = 0; i < reviews.length; i++) {
            totalRating += reviews[i].rating;
        }

        products[_productId].rating = totalRating / reviews.length;
        products[_productId].reviewCount = reviews.length;
    }

    // 获取分类产品
    function getProductsByCategory(ProductCategory _category)
        public
        view
        returns (uint256[] memory)
    {
        return productsByCategory[_category];
    }

    // 获取用户订单
    function getUserOrders(address _user, bool _isBuyer)
        public
        view
        returns (uint256[] memory)
    {
        return _isBuyer ? buyerOrders[_user] : sellerOrders[_user];
    }

    // 获取产品评价
    function getProductReviews(uint256 _productId)
        public
        view
        returns (Review[] memory)
    {
        return productReviews[_productId];
    }
}
```

## 4.6 存储优化策略

### 4.6.1 Gas 优化技巧

```solidity
contract StorageOptimization {
    // 不好的做法：每个变量占用一个存储槽
    struct BadUser {
        bool isActive;      // 32 字节
        uint8 level;        // 32 字节
        uint256 balance;    // 32 字节
        bool isPremium;     // 32 字节
    }

    // 好的做法：打包变量到同一个存储槽
    struct GoodUser {
        uint256 balance;    // 32 字节
        bool isActive;      // 1 字节
        bool isPremium;     // 1 字节
        uint8 level;        // 1 字节
        // 剩余 29 字节未使用，但整体只占用两个存储槽
    }

    // 更好的做法：使用位字段
    struct OptimizedUser {
        uint256 balance;
        uint256 flags;  // 用位来存储多个布尔值和小整数
    }

    mapping(address => OptimizedUser) public users;

    // 位操作常量
    uint256 constant IS_ACTIVE_MASK = 1;
    uint256 constant IS_PREMIUM_MASK = 2;
    uint256 constant LEVEL_MASK = 0xFF00;  // 8 位用于等级
    uint256 constant LEVEL_SHIFT = 8;

    function setUserFlags(address _user, bool _isActive, bool _isPremium, uint8 _level) public {
        uint256 flags = 0;

        if (_isActive) {
            flags |= IS_ACTIVE_MASK;
        }

        if (_isPremium) {
            flags |= IS_PREMIUM_MASK;
        }

        flags |= (uint256(_level) << LEVEL_SHIFT) & LEVEL_MASK;

        users[_user].flags = flags;
    }

    function getUserFlags(address _user)
        public
        view
        returns (bool isActive, bool isPremium, uint8 level)
    {
        uint256 flags = users[_user].flags;

        isActive = (flags & IS_ACTIVE_MASK) != 0;
        isPremium = (flags & IS_PREMIUM_MASK) != 0;
        level = uint8((flags & LEVEL_MASK) >> LEVEL_SHIFT);
    }
}
```

## 4.7 实战练习

### 练习 1：图书馆管理系统

```solidity
// TODO: 设计一个图书馆管理系统
contract LibrarySystem {
    // TODO: 定义书籍、用户、借阅记录等结构体
    // TODO: 实现借书、还书功能
    // TODO: 添加书籍搜索功能
    // TODO: 实现罚金系统
}
```

### 练习 2：社交网络

```solidity
// TODO: 创建一个简单的社交网络
contract SocialNetwork {
    // TODO: 用户资料管理
    // TODO: 好友关系管理
    // TODO: 发布和点赞功能
    // TODO: 消息系统
}
```

### 练习 3：供应链追踪

```solidity
// TODO: 实现供应链追踪系统
contract SupplyChain {
    // TODO: 产品生命周期管理
    // TODO: 所有权转移记录
    // TODO: 质量检验系统
    // TODO: 溯源功能
}
```

## 4.8 章节总结

在本章中，我们深入学习了：

1. **数组操作**：动态数组、定长数组的创建和操作
2. **映射应用**：基本映射、嵌套映射和复杂映射
3. **结构体设计**：自定义数据类型和嵌套结构体
4. **枚举应用**：状态管理和权限控制
5. **组合使用**：复杂系统的数据结构设计
6. **存储优化**：Gas 效率和存储策略

### 关键知识点

- 选择合适的数据结构对合约性能至关重要
- 结构体打包可以显著减少 Gas 消耗
- 映射提供了高效的键值存储
- 枚举增强了代码的可读性和安全性

### 下一章预告

在下一章《合约交互与事件》中，我们将学习：

- 事件的定义和使用
- 合约间的调用和交互
- 接口和抽象合约
- 库的创建和使用

继续您的 Solidity 学习之旅：[第五章 - 合约交互与事件](05_contract_interaction_events.md)
