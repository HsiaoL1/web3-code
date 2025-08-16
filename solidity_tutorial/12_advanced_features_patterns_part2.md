# 第十二章：高级特性与模式（下篇）- 工厂模式与 Oracle 集成

## 本章学习目标

- 掌握工厂模式的设计和实现
- 学会使用 CREATE2 进行确定性部署
- 了解 Oracle 集成的原理和应用
- 掌握 Chainlink 价格馈送使用
- 学会跨链交互的基础概念
- 理解去中心化身份和多签钱包模式

## 12.6 工厂模式

### 12.6.1 基础工厂模式

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * 简单的代币工厂
 * 允许用户部署自己的 ERC20 代币
 */

// 可部署的代币合约模板
contract SimpleToken {
    string public name;
    string public symbol;
    uint8 public decimals;
    uint256 public totalSupply;
    address public owner;

    mapping(address => uint256) public balanceOf;
    mapping(address => mapping(address => uint256)) public allowance;

    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);

    constructor(
        string memory _name,
        string memory _symbol,
        uint8 _decimals,
        uint256 _totalSupply,
        address _owner
    ) {
        name = _name;
        symbol = _symbol;
        decimals = _decimals;
        totalSupply = _totalSupply;
        owner = _owner;
        balanceOf[_owner] = _totalSupply;
        emit Transfer(address(0), _owner, _totalSupply);
    }

    function transfer(address to, uint256 value) public returns (bool) {
        require(balanceOf[msg.sender] >= value, "Insufficient balance");
        balanceOf[msg.sender] -= value;
        balanceOf[to] += value;
        emit Transfer(msg.sender, to, value);
        return true;
    }

    function approve(address spender, uint256 value) public returns (bool) {
        allowance[msg.sender][spender] = value;
        emit Approval(msg.sender, spender, value);
        return true;
    }

    function transferFrom(address from, address to, uint256 value) public returns (bool) {
        require(balanceOf[from] >= value, "Insufficient balance");
        require(allowance[from][msg.sender] >= value, "Insufficient allowance");

        balanceOf[from] -= value;
        balanceOf[to] += value;
        allowance[from][msg.sender] -= value;

        emit Transfer(from, to, value);
        return true;
    }
}

/**
 * 代币工厂合约
 */
contract TokenFactory is Ownable {
    struct TokenInfo {
        address tokenAddress;
        string name;
        string symbol;
        address creator;
        uint256 createdAt;
        bool isActive;
    }

    TokenInfo[] public deployedTokens;
    mapping(address => TokenInfo[]) public userTokens;
    mapping(address => bool) public isFactoryToken;

    uint256 public creationFee;
    uint256 public totalTokensCreated;

    event TokenCreated(
        address indexed tokenAddress,
        address indexed creator,
        string name,
        string symbol,
        uint256 totalSupply
    );

    event CreationFeeUpdated(uint256 oldFee, uint256 newFee);

    constructor(uint256 _creationFee) Ownable(msg.sender) {
        creationFee = _creationFee;
    }

    function createToken(
        string memory name,
        string memory symbol,
        uint8 decimals,
        uint256 totalSupply
    ) public payable returns (address tokenAddress) {
        require(msg.value >= creationFee, "Insufficient creation fee");
        require(bytes(name).length > 0, "Name cannot be empty");
        require(bytes(symbol).length > 0, "Symbol cannot be empty");
        require(totalSupply > 0, "Total supply must be greater than 0");

        // 部署新的代币合约
        SimpleToken newToken = new SimpleToken(
            name,
            symbol,
            decimals,
            totalSupply,
            msg.sender
        );

        tokenAddress = address(newToken);

        // 记录代币信息
        TokenInfo memory tokenInfo = TokenInfo({
            tokenAddress: tokenAddress,
            name: name,
            symbol: symbol,
            creator: msg.sender,
            createdAt: block.timestamp,
            isActive: true
        });

        deployedTokens.push(tokenInfo);
        userTokens[msg.sender].push(tokenInfo);
        isFactoryToken[tokenAddress] = true;
        totalTokensCreated++;

        // 退还多余的费用
        if (msg.value > creationFee) {
            payable(msg.sender).transfer(msg.value - creationFee);
        }

        emit TokenCreated(tokenAddress, msg.sender, name, symbol, totalSupply);
    }

    function setCreationFee(uint256 newFee) public onlyOwner {
        uint256 oldFee = creationFee;
        creationFee = newFee;
        emit CreationFeeUpdated(oldFee, newFee);
    }

    function withdrawFees() public onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "No fees to withdraw");
        payable(owner()).transfer(balance);
    }

    function getDeployedTokensCount() public view returns (uint256) {
        return deployedTokens.length;
    }

    function getUserTokensCount(address user) public view returns (uint256) {
        return userTokens[user].length;
    }

    function getTokenInfo(uint256 index) public view returns (TokenInfo memory) {
        require(index < deployedTokens.length, "Index out of bounds");
        return deployedTokens[index];
    }

    function getUserTokenInfo(address user, uint256 index) public view returns (TokenInfo memory) {
        require(index < userTokens[user].length, "Index out of bounds");
        return userTokens[user][index];
    }

    // 批量获取代币信息
    function getTokensBatch(uint256 start, uint256 count)
        public
        view
        returns (TokenInfo[] memory tokens)
    {
        require(start < deployedTokens.length, "Start index out of bounds");

        uint256 end = start + count;
        if (end > deployedTokens.length) {
            end = deployedTokens.length;
        }

        tokens = new TokenInfo[](end - start);
        for (uint256 i = start; i < end; i++) {
            tokens[i - start] = deployedTokens[i];
        }
    }
}
```

### 12.6.2 CREATE2 确定性部署

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * 使用 CREATE2 的确定性工厂
 * 可以预测合约地址并创建最小代理
 */
contract CREATE2Factory {
    event ContractDeployed(address indexed contractAddress, bytes32 indexed salt);

    // 使用 CREATE2 部署合约
    function deploy(bytes memory bytecode, bytes32 salt) public returns (address) {
        address deployedAddress;

        assembly {
            deployedAddress := create2(
                0,                // value
                add(bytecode, 0x20), // bytecode data
                mload(bytecode),     // bytecode length
                salt                 // salt
            )
        }

        require(deployedAddress != address(0), "Deployment failed");

        emit ContractDeployed(deployedAddress, salt);
        return deployedAddress;
    }

    // 预测合约地址
    function predictAddress(bytes memory bytecode, bytes32 salt)
        public
        view
        returns (address)
    {
        bytes32 hash = keccak256(
            abi.encodePacked(
                bytes1(0xff),
                address(this),
                salt,
                keccak256(bytecode)
            )
        );

        return address(uint160(uint256(hash)));
    }

    // 检查地址是否已部署
    function isDeployed(address addr) public view returns (bool) {
        uint256 size;
        assembly {
            size := extcodesize(addr)
        }
        return size > 0;
    }
}

/**
 * 最小代理工厂（EIP-1167）
 * 节省 Gas 的代理部署方式
 */
contract MinimalProxyFactory {
    address public immutable implementation;

    event ProxyDeployed(address indexed proxy, bytes32 indexed salt);

    constructor(address _implementation) {
        implementation = _implementation;
    }

    function createProxy(bytes32 salt) public returns (address proxy) {
        bytes memory bytecode = getCreationCode();

        assembly {
            proxy := create2(0, add(bytecode, 0x20), mload(bytecode), salt)
        }

        require(proxy != address(0), "Proxy deployment failed");

        emit ProxyDeployed(proxy, salt);
    }

    function predictProxyAddress(bytes32 salt) public view returns (address) {
        bytes memory bytecode = getCreationCode();
        bytes32 hash = keccak256(
            abi.encodePacked(
                bytes1(0xff),
                address(this),
                salt,
                keccak256(bytecode)
            )
        );

        return address(uint160(uint256(hash)));
    }

    function getCreationCode() public view returns (bytes memory) {
        return abi.encodePacked(
            hex"3d602d80600a3d3981f3363d3d373d3d3d363d73",
            implementation,
            hex"5af43d82803e903d91602b57fd5bf3"
        );
    }
}

/**
 * 多合约部署工厂
 */
contract MultiContractFactory {
    struct DeploymentRecord {
        address contractAddress;
        bytes32 contractType;
        bytes32 salt;
        address deployer;
        uint256 timestamp;
    }

    DeploymentRecord[] public deployments;
    mapping(bytes32 => address[]) public contractsByType;
    mapping(address => DeploymentRecord[]) public deploymentsByUser;

    event MultiContractDeployed(
        address[] contractAddresses,
        bytes32[] contractTypes,
        address deployer
    );

    function deployMultiple(
        bytes[] memory bytecodes,
        bytes32[] memory contractTypes,
        bytes32[] memory salts
    ) public returns (address[] memory deployedAddresses) {
        require(
            bytecodes.length == contractTypes.length &&
            contractTypes.length == salts.length,
            "Array lengths mismatch"
        );

        deployedAddresses = new address[](bytecodes.length);

        for (uint256 i = 0; i < bytecodes.length; i++) {
            address deployed = _deploy(bytecodes[i], salts[i]);
            deployedAddresses[i] = deployed;

            // 记录部署信息
            DeploymentRecord memory record = DeploymentRecord({
                contractAddress: deployed,
                contractType: contractTypes[i],
                salt: salts[i],
                deployer: msg.sender,
                timestamp: block.timestamp
            });

            deployments.push(record);
            contractsByType[contractTypes[i]].push(deployed);
            deploymentsByUser[msg.sender].push(record);
        }

        emit MultiContractDeployed(deployedAddresses, contractTypes, msg.sender);
    }

    function _deploy(bytes memory bytecode, bytes32 salt) internal returns (address) {
        address deployedAddress;

        assembly {
            deployedAddress := create2(
                0,
                add(bytecode, 0x20),
                mload(bytecode),
                salt
            )
        }

        require(deployedAddress != address(0), "Deployment failed");
        return deployedAddress;
    }

    function getContractsByType(bytes32 contractType)
        public
        view
        returns (address[] memory)
    {
        return contractsByType[contractType];
    }

    function getUserDeployments(address user)
        public
        view
        returns (DeploymentRecord[] memory)
    {
        return deploymentsByUser[user];
    }

    function getTotalDeployments() public view returns (uint256) {
        return deployments.length;
    }
}
```

### 12.6.3 高级工厂模式

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/proxy/Clones.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";

/**
 * 基于模板的高级工厂
 */
contract AdvancedFactory is AccessControl {
    using Clones for address;

    bytes32 public constant TEMPLATE_MANAGER_ROLE = keccak256("TEMPLATE_MANAGER_ROLE");
    bytes32 public constant DEPLOYER_ROLE = keccak256("DEPLOYER_ROLE");

    struct Template {
        address implementation;
        string name;
        string version;
        bytes4[] supportedInterfaces;
        bool isActive;
        uint256 deploymentCount;
    }

    struct DeploymentConfig {
        bytes32 templateId;
        bytes initData;
        bytes32 salt;
        address deployer;
    }

    mapping(bytes32 => Template) public templates;
    mapping(address => bytes32) public instanceToTemplate;
    mapping(bytes32 => address[]) public templateInstances;

    bytes32[] public templateIds;

    event TemplateAdded(
        bytes32 indexed templateId,
        address implementation,
        string name,
        string version
    );

    event TemplateUpdated(bytes32 indexed templateId, bool isActive);

    event ContractDeployed(
        bytes32 indexed templateId,
        address indexed instance,
        address indexed deployer,
        bytes32 salt
    );

    constructor() {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(TEMPLATE_MANAGER_ROLE, msg.sender);
        _grantRole(DEPLOYER_ROLE, msg.sender);
    }

    function addTemplate(
        bytes32 templateId,
        address implementation,
        string memory name,
        string memory version,
        bytes4[] memory supportedInterfaces
    ) public onlyRole(TEMPLATE_MANAGER_ROLE) {
        require(implementation != address(0), "Invalid implementation");
        require(templates[templateId].implementation == address(0), "Template exists");

        templates[templateId] = Template({
            implementation: implementation,
            name: name,
            version: version,
            supportedInterfaces: supportedInterfaces,
            isActive: true,
            deploymentCount: 0
        });

        templateIds.push(templateId);

        emit TemplateAdded(templateId, implementation, name, version);
    }

    function updateTemplateStatus(bytes32 templateId, bool isActive)
        public
        onlyRole(TEMPLATE_MANAGER_ROLE)
    {
        require(templates[templateId].implementation != address(0), "Template not found");
        templates[templateId].isActive = isActive;
        emit TemplateUpdated(templateId, isActive);
    }

    function deployFromTemplate(
        bytes32 templateId,
        bytes memory initData,
        bytes32 salt
    ) public onlyRole(DEPLOYER_ROLE) returns (address instance) {
        Template storage template = templates[templateId];
        require(template.implementation != address(0), "Template not found");
        require(template.isActive, "Template not active");

        // 使用 CREATE2 部署克隆
        instance = Clones.cloneDeterministic(template.implementation, salt);

        // 初始化克隆实例
        if (initData.length > 0) {
            (bool success, ) = instance.call(initData);
            require(success, "Initialization failed");
        }

        // 记录部署信息
        instanceToTemplate[instance] = templateId;
        templateInstances[templateId].push(instance);
        template.deploymentCount++;

        emit ContractDeployed(templateId, instance, msg.sender, salt);
    }

    function predictInstanceAddress(
        bytes32 templateId,
        bytes32 salt
    ) public view returns (address) {
        Template memory template = templates[templateId];
        require(template.implementation != address(0), "Template not found");

        return Clones.predictDeterministicAddress(
            template.implementation,
            salt,
            address(this)
        );
    }

    function batchDeploy(
        DeploymentConfig[] memory configs
    ) public onlyRole(DEPLOYER_ROLE) returns (address[] memory instances) {
        instances = new address[](configs.length);

        for (uint256 i = 0; i < configs.length; i++) {
            instances[i] = deployFromTemplate(
                configs[i].templateId,
                configs[i].initData,
                configs[i].salt
            );
        }
    }

    function getTemplateInstances(bytes32 templateId)
        public
        view
        returns (address[] memory)
    {
        return templateInstances[templateId];
    }

    function getAllTemplates() public view returns (bytes32[] memory) {
        return templateIds;
    }

    function getTemplateInfo(bytes32 templateId)
        public
        view
        returns (
            address implementation,
            string memory name,
            string memory version,
            bytes4[] memory supportedInterfaces,
            bool isActive,
            uint256 deploymentCount
        )
    {
        Template memory template = templates[templateId];
        return (
            template.implementation,
            template.name,
            template.version,
            template.supportedInterfaces,
            template.isActive,
            template.deploymentCount
        );
    }

    // 验证合约是否支持特定接口
    function supportsInterface(bytes32 templateId, bytes4 interfaceId)
        public
        view
        returns (bool)
    {
        bytes4[] memory interfaces = templates[templateId].supportedInterfaces;
        for (uint256 i = 0; i < interfaces.length; i++) {
            if (interfaces[i] == interfaceId) {
                return true;
            }
        }
        return false;
    }
}
```

## 12.7 Oracle 集成

### 12.7.1 Chainlink 价格馈送

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@chainlink/contracts/src/v0.8/interfaces/AggregatorV3Interface.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * Chainlink 价格馈送集成
 */
contract PriceOracle is Ownable {
    struct PriceFeed {
        AggregatorV3Interface feed;
        string description;
        uint8 decimals;
        bool isActive;
    }

    mapping(string => PriceFeed) public priceFeeds;
    string[] public availableFeeds;

    uint256 public constant STALENESS_THRESHOLD = 3600; // 1小时

    event PriceFeedAdded(string symbol, address feedAddress);
    event PriceFeedUpdated(string symbol, address feedAddress);
    event PriceFeedRemoved(string symbol);

    constructor() Ownable(msg.sender) {}

    function addPriceFeed(
        string memory symbol,
        address feedAddress,
        string memory description
    ) public onlyOwner {
        require(feedAddress != address(0), "Invalid feed address");
        require(address(priceFeeds[symbol].feed) == address(0), "Feed already exists");

        AggregatorV3Interface feed = AggregatorV3Interface(feedAddress);

        priceFeeds[symbol] = PriceFeed({
            feed: feed,
            description: description,
            decimals: feed.decimals(),
            isActive: true
        });

        availableFeeds.push(symbol);

        emit PriceFeedAdded(symbol, feedAddress);
    }

    function updatePriceFeed(
        string memory symbol,
        address feedAddress
    ) public onlyOwner {
        require(address(priceFeeds[symbol].feed) != address(0), "Feed not found");
        require(feedAddress != address(0), "Invalid feed address");

        AggregatorV3Interface feed = AggregatorV3Interface(feedAddress);
        priceFeeds[symbol].feed = feed;
        priceFeeds[symbol].decimals = feed.decimals();

        emit PriceFeedUpdated(symbol, feedAddress);
    }

    function removePriceFeed(string memory symbol) public onlyOwner {
        require(address(priceFeeds[symbol].feed) != address(0), "Feed not found");

        delete priceFeeds[symbol];

        // 从数组中移除
        for (uint256 i = 0; i < availableFeeds.length; i++) {
            if (keccak256(bytes(availableFeeds[i])) == keccak256(bytes(symbol))) {
                availableFeeds[i] = availableFeeds[availableFeeds.length - 1];
                availableFeeds.pop();
                break;
            }
        }

        emit PriceFeedRemoved(symbol);
    }

    function getLatestPrice(string memory symbol)
        public
        view
        returns (int256 price, uint256 timestamp, uint8 decimals)
    {
        PriceFeed memory feed = priceFeeds[symbol];
        require(address(feed.feed) != address(0), "Feed not found");
        require(feed.isActive, "Feed not active");

        (, int256 latestPrice, , uint256 updatedAt, ) = feed.feed.latestRoundData();

        require(latestPrice > 0, "Invalid price");
        require(block.timestamp - updatedAt <= STALENESS_THRESHOLD, "Price data stale");

        return (latestPrice, updatedAt, feed.decimals);
    }

    function getHistoricalPrice(string memory symbol, uint80 roundId)
        public
        view
        returns (int256 price, uint256 timestamp, uint8 decimals)
    {
        PriceFeed memory feed = priceFeeds[symbol];
        require(address(feed.feed) != address(0), "Feed not found");

        (, int256 historicalPrice, , uint256 updatedAt, ) = feed.feed.getRoundData(roundId);

        require(historicalPrice > 0, "Invalid price");

        return (historicalPrice, updatedAt, feed.decimals);
    }

    function getPriceInUSD(string memory symbol, uint256 amount)
        public
        view
        returns (uint256 usdValue)
    {
        (int256 price, , uint8 decimals) = getLatestPrice(symbol);

        usdValue = (amount * uint256(price)) / (10 ** decimals);
    }

    function comparePrices(string memory symbol1, string memory symbol2)
        public
        view
        returns (int256 ratio, uint8 decimals)
    {
        (int256 price1, , uint8 decimals1) = getLatestPrice(symbol1);
        (int256 price2, , uint8 decimals2) = getLatestPrice(symbol2);

        require(price2 > 0, "Invalid price for comparison");

        // 标准化到相同的小数位数
        if (decimals1 > decimals2) {
            price2 = price2 * int256(10 ** (decimals1 - decimals2));
            decimals = decimals1;
        } else if (decimals2 > decimals1) {
            price1 = price1 * int256(10 ** (decimals2 - decimals1));
            decimals = decimals2;
        } else {
            decimals = decimals1;
        }

        ratio = (price1 * int256(10 ** decimals)) / price2;
    }

    function getMultiplePrices(string[] memory symbols)
        public
        view
        returns (
            int256[] memory prices,
            uint256[] memory timestamps,
            uint8[] memory decimalsArray
        )
    {
        prices = new int256[](symbols.length);
        timestamps = new uint256[](symbols.length);
        decimalsArray = new uint8[](symbols.length);

        for (uint256 i = 0; i < symbols.length; i++) {
            (prices[i], timestamps[i], decimalsArray[i]) = getLatestPrice(symbols[i]);
        }
    }

    function isPriceStale(string memory symbol) public view returns (bool) {
        PriceFeed memory feed = priceFeeds[symbol];
        require(address(feed.feed) != address(0), "Feed not found");

        (, , , uint256 updatedAt, ) = feed.feed.latestRoundData();

        return block.timestamp - updatedAt > STALENESS_THRESHOLD;
    }

    function getFeedInfo(string memory symbol)
        public
        view
        returns (
            address feedAddress,
            string memory description,
            uint8 decimals,
            bool isActive
        )
    {
        PriceFeed memory feed = priceFeeds[symbol];
        return (
            address(feed.feed),
            feed.description,
            feed.decimals,
            feed.isActive
        );
    }

    function getAllAvailableFeeds() public view returns (string[] memory) {
        return availableFeeds;
    }
}

/**
 * 基于 Oracle 的价格敏感合约示例
 */
contract PriceSensitiveContract {
    PriceOracle public priceOracle;

    struct Position {
        uint256 amount;
        uint256 entryPrice;
        uint256 timestamp;
        bool isLong;
    }

    mapping(address => mapping(string => Position)) public positions;
    mapping(string => uint256) public collateralRatios; // 抵押率（基点）

    uint256 public constant LIQUIDATION_THRESHOLD = 8000; // 80%

    event PositionOpened(
        address indexed user,
        string symbol,
        uint256 amount,
        uint256 entryPrice,
        bool isLong
    );

    event PositionClosed(
        address indexed user,
        string symbol,
        uint256 pnl,
        bool isProfit
    );

    event Liquidation(
        address indexed user,
        string symbol,
        uint256 liquidationPrice
    );

    constructor(address _priceOracle) {
        priceOracle = PriceOracle(_priceOracle);

        // 设置默认抵押率
        collateralRatios["ETH"] = 7500; // 75%
        collateralRatios["BTC"] = 7500; // 75%
        collateralRatios["LINK"] = 6000; // 60%
    }

    function openPosition(
        string memory symbol,
        uint256 amount,
        bool isLong
    ) public payable {
        require(amount > 0, "Invalid amount");
        require(msg.value > 0, "Collateral required");

        (int256 currentPrice, , ) = priceOracle.getLatestPrice(symbol);
        require(currentPrice > 0, "Invalid price");

        uint256 requiredCollateral = (amount * uint256(currentPrice) * collateralRatios[symbol]) / 10000;
        require(msg.value >= requiredCollateral, "Insufficient collateral");

        positions[msg.sender][symbol] = Position({
            amount: amount,
            entryPrice: uint256(currentPrice),
            timestamp: block.timestamp,
            isLong: isLong
        });

        emit PositionOpened(msg.sender, symbol, amount, uint256(currentPrice), isLong);
    }

    function closePosition(string memory symbol) public {
        Position memory position = positions[msg.sender][symbol];
        require(position.amount > 0, "No position found");

        (int256 currentPrice, , ) = priceOracle.getLatestPrice(symbol);
        require(currentPrice > 0, "Invalid price");

        (uint256 pnl, bool isProfit) = calculatePnL(
            position.amount,
            position.entryPrice,
            uint256(currentPrice),
            position.isLong
        );

        delete positions[msg.sender][symbol];

        // 处理盈亏
        if (isProfit) {
            payable(msg.sender).transfer(pnl);
        }

        emit PositionClosed(msg.sender, symbol, pnl, isProfit);
    }

    function liquidatePosition(address user, string memory symbol) public {
        Position memory position = positions[user][symbol];
        require(position.amount > 0, "No position found");

        (int256 currentPrice, , ) = priceOracle.getLatestPrice(symbol);
        require(currentPrice > 0, "Invalid price");

        uint256 currentValue = position.amount * uint256(currentPrice);
        uint256 requiredCollateral = (currentValue * collateralRatios[symbol]) / 10000;

        // 检查是否达到清算条件
        uint256 collateralValue = (currentValue * LIQUIDATION_THRESHOLD) / 10000;
        require(collateralValue <= requiredCollateral, "Position not liquidatable");

        delete positions[user][symbol];

        emit Liquidation(user, symbol, uint256(currentPrice));
    }

    function calculatePnL(
        uint256 amount,
        uint256 entryPrice,
        uint256 currentPrice,
        bool isLong
    ) public pure returns (uint256 pnl, bool isProfit) {
        if (isLong) {
            if (currentPrice > entryPrice) {
                pnl = amount * (currentPrice - entryPrice) / entryPrice;
                isProfit = true;
            } else {
                pnl = amount * (entryPrice - currentPrice) / entryPrice;
                isProfit = false;
            }
        } else {
            if (currentPrice < entryPrice) {
                pnl = amount * (entryPrice - currentPrice) / entryPrice;
                isProfit = true;
            } else {
                pnl = amount * (currentPrice - entryPrice) / entryPrice;
                isProfit = false;
            }
        }
    }

    function getPositionValue(address user, string memory symbol)
        public
        view
        returns (uint256 currentValue, uint256 pnl, bool isProfit)
    {
        Position memory position = positions[user][symbol];
        require(position.amount > 0, "No position found");

        (int256 currentPrice, , ) = priceOracle.getLatestPrice(symbol);
        require(currentPrice > 0, "Invalid price");

        currentValue = position.amount * uint256(currentPrice);
        (pnl, isProfit) = calculatePnL(
            position.amount,
            position.entryPrice,
            uint256(currentPrice),
            position.isLong
        );
    }

    function checkLiquidationRisk(address user, string memory symbol)
        public
        view
        returns (bool atRisk, uint256 healthFactor)
    {
        Position memory position = positions[user][symbol];
        if (position.amount == 0) {
            return (false, 0);
        }

        (int256 currentPrice, , ) = priceOracle.getLatestPrice(symbol);
        if (currentPrice <= 0) {
            return (true, 0);
        }

        uint256 currentValue = position.amount * uint256(currentPrice);
        uint256 requiredCollateral = (currentValue * collateralRatios[symbol]) / 10000;
        uint256 liquidationThreshold = (currentValue * LIQUIDATION_THRESHOLD) / 10000;

        healthFactor = (requiredCollateral * 10000) / liquidationThreshold;
        atRisk = healthFactor <= 10000; // 100%
    }
}
```

### 12.7.2 自定义 Oracle 实现

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";

/**
 * 自定义 Oracle 系统
 * 支持多数据源聚合和去中心化价格馈送
 */
contract CustomOracle is AccessControl, ReentrancyGuard {
    bytes32 public constant ORACLE_ROLE = keccak256("ORACLE_ROLE");
    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");

    struct PriceData {
        uint256 price;
        uint256 timestamp;
        uint256 roundId;
        address oracle;
    }

    struct AggregatedPrice {
        uint256 price;
        uint256 timestamp;
        uint256 confidence;
        uint8 decimals;
    }

    mapping(string => PriceData[]) public priceHistory;
    mapping(string => AggregatedPrice) public latestPrices;
    mapping(address => bool) public authorizedOracles;
    mapping(string => uint8) public assetDecimals;
    mapping(string => uint256) public minOracleCount;
    mapping(string => uint256) public maxPriceDeviation; // 基点

    uint256 public constant MAX_PRICE_STALENESS = 3600; // 1小时
    uint256 public constant DEFAULT_MIN_ORACLES = 3;
    uint256 public constant DEFAULT_MAX_DEVIATION = 500; // 5%

    event PriceSubmitted(
        string indexed asset,
        uint256 price,
        uint256 timestamp,
        address oracle,
        uint256 roundId
    );

    event PriceAggregated(
        string indexed asset,
        uint256 price,
        uint256 confidence,
        uint256 timestamp
    );

    event OracleAdded(address indexed oracle);
    event OracleRemoved(address indexed oracle);

    constructor() {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(ADMIN_ROLE, msg.sender);
    }

    function addOracle(address oracle) public onlyRole(ADMIN_ROLE) {
        require(oracle != address(0), "Invalid oracle address");
        require(!authorizedOracles[oracle], "Oracle already authorized");

        authorizedOracles[oracle] = true;
        _grantRole(ORACLE_ROLE, oracle);

        emit OracleAdded(oracle);
    }

    function removeOracle(address oracle) public onlyRole(ADMIN_ROLE) {
        require(authorizedOracles[oracle], "Oracle not authorized");

        authorizedOracles[oracle] = false;
        _revokeRole(ORACLE_ROLE, oracle);

        emit OracleRemoved(oracle);
    }

    function setAssetConfig(
        string memory asset,
        uint8 decimals,
        uint256 minOracles,
        uint256 maxDeviation
    ) public onlyRole(ADMIN_ROLE) {
        assetDecimals[asset] = decimals;
        minOracleCount[asset] = minOracles == 0 ? DEFAULT_MIN_ORACLES : minOracles;
        maxPriceDeviation[asset] = maxDeviation == 0 ? DEFAULT_MAX_DEVIATION : maxDeviation;
    }

    function submitPrice(
        string memory asset,
        uint256 price,
        uint256 roundId
    ) public onlyRole(ORACLE_ROLE) nonReentrant {
        require(price > 0, "Invalid price");
        require(bytes(asset).length > 0, "Invalid asset");

        // 检查价格偏差
        if (latestPrices[asset].timestamp > 0) {
            uint256 lastPrice = latestPrices[asset].price;
            uint256 deviation = price > lastPrice
                ? ((price - lastPrice) * 10000) / lastPrice
                : ((lastPrice - price) * 10000) / lastPrice;

            require(
                deviation <= maxPriceDeviation[asset],
                "Price deviation too high"
            );
        }

        PriceData memory newPrice = PriceData({
            price: price,
            timestamp: block.timestamp,
            roundId: roundId,
            oracle: msg.sender
        });

        priceHistory[asset].push(newPrice);

        emit PriceSubmitted(asset, price, block.timestamp, msg.sender, roundId);

        // 触发聚合
        _aggregatePrice(asset);
    }

    function _aggregatePrice(string memory asset) internal {
        uint256 recentCount = 0;
        uint256 totalPrice = 0;
        uint256 minTimestamp = block.timestamp - MAX_PRICE_STALENESS;

        // 收集最近的价格数据
        uint256[] memory recentPrices = new uint256[](priceHistory[asset].length);

        for (uint256 i = priceHistory[asset].length; i > 0; i--) {
            PriceData memory data = priceHistory[asset][i - 1];
            if (data.timestamp >= minTimestamp) {
                recentPrices[recentCount] = data.price;
                recentCount++;
                if (recentCount >= 50) break; // 限制最大样本数
            }
        }

        require(recentCount >= minOracleCount[asset], "Insufficient oracle data");

        // 排序价格数据
        _quickSort(recentPrices, 0, int256(recentCount - 1));

        // 计算中位数
        uint256 medianPrice;
        if (recentCount % 2 == 0) {
            medianPrice = (recentPrices[recentCount / 2 - 1] + recentPrices[recentCount / 2]) / 2;
        } else {
            medianPrice = recentPrices[recentCount / 2];
        }

        // 计算置信度（基于价格分散程度）
        uint256 confidence = _calculateConfidence(recentPrices, recentCount, medianPrice);

        latestPrices[asset] = AggregatedPrice({
            price: medianPrice,
            timestamp: block.timestamp,
            confidence: confidence,
            decimals: assetDecimals[asset]
        });

        emit PriceAggregated(asset, medianPrice, confidence, block.timestamp);
    }

    function _calculateConfidence(
        uint256[] memory prices,
        uint256 count,
        uint256 median
    ) internal pure returns (uint256) {
        if (count == 0) return 0;

        uint256 totalDeviation = 0;
        for (uint256 i = 0; i < count; i++) {
            uint256 deviation = prices[i] > median
                ? prices[i] - median
                : median - prices[i];
            totalDeviation += (deviation * 10000) / median;
        }

        uint256 avgDeviation = totalDeviation / count;

        // 置信度 = 100% - 平均偏差
        return avgDeviation >= 10000 ? 0 : 10000 - avgDeviation;
    }

    function _quickSort(uint256[] memory arr, int256 left, int256 right) internal pure {
        int256 i = left;
        int256 j = right;
        if (i == j) return;
        uint256 pivot = arr[uint256(left + (right - left) / 2)];

        while (i <= j) {
            while (arr[uint256(i)] < pivot) i++;
            while (pivot < arr[uint256(j)]) j--;
            if (i <= j) {
                (arr[uint256(i)], arr[uint256(j)]) = (arr[uint256(j)], arr[uint256(i)]);
                i++;
                j--;
            }
        }

        if (left < j) _quickSort(arr, left, j);
        if (i < right) _quickSort(arr, i, right);
    }

    function getLatestPrice(string memory asset)
        public
        view
        returns (uint256 price, uint256 timestamp, uint256 confidence, uint8 decimals)
    {
        AggregatedPrice memory latest = latestPrices[asset];
        require(latest.timestamp > 0, "No price data available");
        require(
            block.timestamp - latest.timestamp <= MAX_PRICE_STALENESS,
            "Price data stale"
        );

        return (latest.price, latest.timestamp, latest.confidence, latest.decimals);
    }

    function getPriceHistory(string memory asset, uint256 limit)
        public
        view
        returns (PriceData[] memory)
    {
        uint256 length = priceHistory[asset].length;
        uint256 returnLength = limit > length ? length : limit;

        PriceData[] memory result = new PriceData[](returnLength);

        for (uint256 i = 0; i < returnLength; i++) {
            result[i] = priceHistory[asset][length - 1 - i];
        }

        return result;
    }

    function getTWAP(string memory asset, uint256 timeWindow)
        public
        view
        returns (uint256 twap)
    {
        require(timeWindow > 0, "Invalid time window");

        uint256 minTimestamp = block.timestamp - timeWindow;
        uint256 weightedSum = 0;
        uint256 totalWeight = 0;

        for (uint256 i = priceHistory[asset].length; i > 0; i--) {
            PriceData memory data = priceHistory[asset][i - 1];
            if (data.timestamp < minTimestamp) break;

            uint256 weight = data.timestamp - minTimestamp;
            weightedSum += data.price * weight;
            totalWeight += weight;
        }

        require(totalWeight > 0, "No data in time window");
        twap = weightedSum / totalWeight;
    }

    function isHealthy(string memory asset) public view returns (bool) {
        AggregatedPrice memory latest = latestPrices[asset];

        if (latest.timestamp == 0) return false;
        if (block.timestamp - latest.timestamp > MAX_PRICE_STALENESS) return false;
        if (latest.confidence < 5000) return false; // 最小50%置信度

        return true;
    }
}
```

## 12.8 章节总结

在本章的下篇中，我们学习了：

1. **工厂模式**：基础工厂、CREATE2 确定性部署、最小代理模式
2. **高级工厂**：模板管理、批量部署、接口支持检查
3. **Oracle 集成**：Chainlink 价格馈送、自定义 Oracle 系统
4. **价格敏感合约**：基于 Oracle 的金融产品实现
5. **聚合算法**：多数据源聚合、置信度计算、异常检测

### 关键要点

1. **工厂模式优势**：代码复用、标准化部署、成本控制
2. **CREATE2 应用**：地址可预测性、跨链部署一致性
3. **Oracle 重要性**：连接链上链下数据、价格发现机制
4. **数据质量**：多源验证、异常检测、时效性保证
5. **安全考虑**：价格操纵防护、数据源多样性

### 设计模式总结

- **工厂模式**：集中管理合约部署
- **代理模式**：实现合约升级能力
- **Oracle 模式**：安全的外部数据集成
- **聚合模式**：多数据源整合
- **治理模式**：去中心化决策机制

这些高级模式和特性为构建复杂的 DeFi 应用、DAO 系统和跨链协议提供了坚实的技术基础。

恭喜您完成了完整的 Solidity 智能合约开发教程！这 12 个章节涵盖了从基础语法到高级模式的全面内容，为您的区块链开发之旅奠定了坚实的基础。
