# 第十二章：高级特性与模式（上篇）- 代理模式与可升级合约

## 本章学习目标

- 理解智能合约的不可变性挑战
- 掌握代理模式的原理和实现
- 学会使用 OpenZeppelin Upgrades
- 了解透明代理和 UUPS 代理
- 掌握存储布局和升级兼容性
- 学会安全的合约升级策略

## 12.1 合约升级的必要性

### 12.1.1 不可变性的挑战

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * 传统智能合约的限制：
 * 1. 一旦部署，代码无法修改
 * 2. 发现漏洞时无法修复
 * 3. 功能需求变更时无法更新
 * 4. 用户需要迁移到新合约
 */

contract NonUpgradeableContract {
    mapping(address => uint256) public balances;
    address public owner;

    // 假设这里有一个漏洞
    function withdraw(uint256 amount) public {
        // 错误：没有检查余额
        balances[msg.sender] -= amount;
        payable(msg.sender).transfer(amount);
    }

    // 无法修复这个漏洞，只能部署新合约
}

/**
 * 合约升级的需求场景：
 * 1. 安全漏洞修复
 * 2. 功能增强和优化
 * 3. 业务逻辑调整
 * 4. Gas 优化
 * 5. 标准更新适配
 */

contract UpgradeableContract {
    // 使用代理模式实现升级能力
    // 状态存储在代理合约中
    // 逻辑存储在实现合约中

    mapping(address => uint256) public balances;
    address public owner;

    function withdraw(uint256 amount) public {
        require(balances[msg.sender] >= amount, "Insufficient balance");
        balances[msg.sender] -= amount;
        payable(msg.sender).transfer(amount);
    }

    // V2 可以添加新功能
    function emergencyPause() public {
        require(msg.sender == owner, "Not owner");
        // 新增的紧急暂停功能
    }
}
```

### 12.1.2 代理模式概述

```solidity
/**
 * 代理模式原理：
 *
 * 用户 → 代理合约 → 实现合约
 *              ↓
 *          存储状态
 *
 * 特点：
 * 1. 代理合约：存储状态，永不更改
 * 2. 实现合约：存储逻辑，可以更换
 * 3. delegatecall：在代理合约的上下文中执行实现合约的代码
 * 4. 透明升级：用户无感知的合约升级
 */

contract SimpleProxy {
    address public implementation;
    address public admin;

    constructor(address _implementation) {
        implementation = _implementation;
        admin = msg.sender;
    }

    modifier onlyAdmin() {
        require(msg.sender == admin, "Not admin");
        _;
    }

    function upgrade(address newImplementation) external onlyAdmin {
        implementation = newImplementation;
    }

    fallback() external payable {
        address impl = implementation;

        assembly {
            // 复制调用数据
            calldatacopy(0, 0, calldatasize())

            // 使用 delegatecall 调用实现合约
            let result := delegatecall(gas(), impl, 0, calldatasize(), 0, 0)

            // 复制返回数据
            returndatacopy(0, 0, returndatasize())

            switch result
            case 0 { revert(0, returndatasize()) }
            default { return(0, returndatasize()) }
        }
    }

    receive() external payable {
        // 处理纯ETH转账
    }
}
```

## 12.2 OpenZeppelin Upgrades

### 12.2.1 透明代理模式

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

/**
 * 可升级代币合约 V1
 * 使用 OpenZeppelin Upgrades 框架
 */
contract UpgradeableTokenV1 is
    Initializable,
    ERC20Upgradeable,
    OwnableUpgradeable,
    UUPSUpgradeable
{
    uint256 public constant MAX_SUPPLY = 1000000 * 10**18;

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize(
        string memory name,
        string memory symbol,
        address owner
    ) public initializer {
        __ERC20_init(name, symbol);
        __Ownable_init(owner);
        __UUPSUpgradeable_init();

        // 初始铸造
        _mint(owner, 100000 * 10**18);
    }

    function mint(address to, uint256 amount) public onlyOwner {
        require(totalSupply() + amount <= MAX_SUPPLY, "Exceeds max supply");
        _mint(to, amount);
    }

    function burn(uint256 amount) public {
        _burn(msg.sender, amount);
    }

    // 版本信息
    function version() public pure virtual returns (string memory) {
        return "1.0.0";
    }

    // UUPS升级授权
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
}

/**
 * 可升级代币合约 V2 - 添加新功能
 */
contract UpgradeableTokenV2 is UpgradeableTokenV1 {
    // 新增状态变量（必须在现有变量之后）
    mapping(address => bool) public blacklisted;
    bool public transfersEnabled;
    uint256 public burnFeeRate; // 销毁费率（基点）

    event BlacklistUpdated(address indexed account, bool blacklisted);
    event TransfersToggled(bool enabled);
    event BurnFeeUpdated(uint256 newRate);

    // 新的初始化函数（仅用于V2的新功能）
    function initializeV2() public reinitializer(2) {
        transfersEnabled = true;
        burnFeeRate = 100; // 1%
    }

    // 重写转账函数，添加黑名单检查
    function _beforeTokenTransfer(
        address from,
        address to,
        uint256 amount
    ) internal virtual override {
        super._beforeTokenTransfer(from, to, amount);

        require(!blacklisted[from], "From address is blacklisted");
        require(!blacklisted[to], "To address is blacklisted");
        require(transfersEnabled || from == address(0) || to == address(0), "Transfers disabled");
    }

    // 新增功能：黑名单管理
    function setBlacklisted(address account, bool _blacklisted) public onlyOwner {
        blacklisted[account] = _blacklisted;
        emit BlacklistUpdated(account, _blacklisted);
    }

    function batchSetBlacklisted(
        address[] calldata accounts,
        bool _blacklisted
    ) public onlyOwner {
        for (uint256 i = 0; i < accounts.length; i++) {
            blacklisted[accounts[i]] = _blacklisted;
            emit BlacklistUpdated(accounts[i], _blacklisted);
        }
    }

    // 新增功能：转账开关
    function setTransfersEnabled(bool _enabled) public onlyOwner {
        transfersEnabled = _enabled;
        emit TransfersToggled(_enabled);
    }

    // 新增功能：带费用的销毁
    function burnWithFee(uint256 amount) public {
        uint256 feeAmount = (amount * burnFeeRate) / 10000;
        uint256 burnAmount = amount - feeAmount;

        _transfer(msg.sender, owner(), feeAmount);
        _burn(msg.sender, burnAmount);
    }

    // 新增功能：设置销毁费率
    function setBurnFeeRate(uint256 newRate) public onlyOwner {
        require(newRate <= 1000, "Fee rate too high"); // 最大10%
        burnFeeRate = newRate;
        emit BurnFeeUpdated(newRate);
    }

    // 新增功能：批量转账
    function batchTransfer(
        address[] calldata recipients,
        uint256[] calldata amounts
    ) public {
        require(recipients.length == amounts.length, "Arrays length mismatch");

        for (uint256 i = 0; i < recipients.length; i++) {
            transfer(recipients[i], amounts[i]);
        }
    }

    // 新增功能：空投
    function airdrop(address[] calldata recipients, uint256 amount) public onlyOwner {
        for (uint256 i = 0; i < recipients.length; i++) {
            _mint(recipients[i], amount);
        }
    }

    // 更新版本号
    function version() public pure override returns (string memory) {
        return "2.0.0";
    }
}

/**
 * 可升级代币合约 V3 - 添加更多高级功能
 */
contract UpgradeableTokenV3 is UpgradeableTokenV2 {
    // 新增状态变量
    mapping(address => uint256) public lastTransferTime;
    uint256 public transferCooldown;
    mapping(address => bool) public vipUsers;
    uint256 public maxTransferAmount;

    event TransferCooldownUpdated(uint256 newCooldown);
    event VIPStatusUpdated(address indexed user, bool isVIP);
    event MaxTransferAmountUpdated(uint256 newAmount);

    function initializeV3() public reinitializer(3) {
        transferCooldown = 300; // 5分钟冷却
        maxTransferAmount = 10000 * 10**18; // 最大转账额度
    }

    // 重写转账函数，添加冷却时间和转账限额
    function _beforeTokenTransfer(
        address from,
        address to,
        uint256 amount
    ) internal virtual override {
        super._beforeTokenTransfer(from, to, amount);

        if (from != address(0) && !vipUsers[from]) {
            // 检查冷却时间
            require(
                block.timestamp >= lastTransferTime[from] + transferCooldown,
                "Transfer cooldown not met"
            );

            // 检查转账限额
            require(amount <= maxTransferAmount, "Amount exceeds max transfer");
        }

        if (from != address(0)) {
            lastTransferTime[from] = block.timestamp;
        }
    }

    // VIP用户管理
    function setVIPStatus(address user, bool isVIP) public onlyOwner {
        vipUsers[user] = isVIP;
        emit VIPStatusUpdated(user, isVIP);
    }

    // 设置转账冷却时间
    function setTransferCooldown(uint256 newCooldown) public onlyOwner {
        require(newCooldown <= 86400, "Cooldown too long"); // 最大24小时
        transferCooldown = newCooldown;
        emit TransferCooldownUpdated(newCooldown);
    }

    // 设置最大转账金额
    function setMaxTransferAmount(uint256 newAmount) public onlyOwner {
        maxTransferAmount = newAmount;
        emit MaxTransferAmountUpdated(newAmount);
    }

    // 紧急提取（仅VIP用户）
    function emergencyTransfer(address to, uint256 amount) public {
        require(vipUsers[msg.sender], "Not VIP user");
        _transfer(msg.sender, to, amount);
    }

    function version() public pure override returns (string memory) {
        return "3.0.0";
    }
}
```

### 12.2.2 部署和升级脚本

```javascript
// scripts/deploy-upgradeable.js
const { ethers, upgrades } = require("hardhat");

async function main() {
  console.log("部署可升级代币合约...");

  const [deployer] = await ethers.getSigners();
  console.log("部署账户:", deployer.address);

  // 部署可升级合约 V1
  const UpgradeableTokenV1 = await ethers.getContractFactory(
    "UpgradeableTokenV1"
  );

  console.log("部署代理合约...");
  const proxy = await upgrades.deployProxy(
    UpgradeableTokenV1,
    [
      "Upgradeable Token", // name
      "UGT", // symbol
      deployer.address, // owner
    ],
    {
      kind: "uups",
      initializer: "initialize",
    }
  );

  await proxy.waitForDeployment();

  console.log("代理合约地址:", proxy.target);
  console.log(
    "实现合约地址:",
    await upgrades.erc1967.getImplementationAddress(proxy.target)
  );
  console.log(
    "管理员地址:",
    await upgrades.erc1967.getAdminAddress(proxy.target)
  );

  // 验证部署
  const token = await ethers.getContractAt("UpgradeableTokenV1", proxy.target);
  console.log("代币名称:", await token.name());
  console.log("代币符号:", await token.symbol());
  console.log("版本:", await token.version());
  console.log("总供应量:", ethers.formatEther(await token.totalSupply()));

  // 保存部署信息
  const deploymentInfo = {
    network: network.name,
    proxyAddress: proxy.target,
    implementationAddress: await upgrades.erc1967.getImplementationAddress(
      proxy.target
    ),
    adminAddress: await upgrades.erc1967.getAdminAddress(proxy.target),
    deployer: deployer.address,
    version: "1.0.0",
    timestamp: new Date().toISOString(),
  };

  const fs = require("fs");
  fs.writeFileSync(
    `./deployments/upgradeable-${network.name}.json`,
    JSON.stringify(deploymentInfo, null, 2)
  );

  console.log("部署完成！");
  return proxy;
}

async function upgradeToV2() {
  console.log("升级到 V2...");

  const deploymentInfo = JSON.parse(
    fs.readFileSync(`./deployments/upgradeable-${network.name}.json`)
  );

  const UpgradeableTokenV2 = await ethers.getContractFactory(
    "UpgradeableTokenV2"
  );

  console.log("升级合约...");
  const upgraded = await upgrades.upgradeProxy(
    deploymentInfo.proxyAddress,
    UpgradeableTokenV2
  );

  console.log("调用 initializeV2...");
  await upgraded.initializeV2();

  console.log("升级完成！");
  console.log("新版本:", await upgraded.version());

  // 更新部署信息
  deploymentInfo.implementationAddress =
    await upgrades.erc1967.getImplementationAddress(upgraded.target);
  deploymentInfo.version = "2.0.0";
  deploymentInfo.lastUpgrade = new Date().toISOString();

  fs.writeFileSync(
    `./deployments/upgradeable-${network.name}.json`,
    JSON.stringify(deploymentInfo, null, 2)
  );

  return upgraded;
}

async function upgradeToV3() {
  console.log("升级到 V3...");

  const deploymentInfo = JSON.parse(
    fs.readFileSync(`./deployments/upgradeable-${network.name}.json`)
  );

  const UpgradeableTokenV3 = await ethers.getContractFactory(
    "UpgradeableTokenV3"
  );

  console.log("升级合约...");
  const upgraded = await upgrades.upgradeProxy(
    deploymentInfo.proxyAddress,
    UpgradeableTokenV3
  );

  console.log("调用 initializeV3...");
  await upgraded.initializeV3();

  console.log("升级完成！");
  console.log("新版本:", await upgraded.version());

  // 更新部署信息
  deploymentInfo.implementationAddress =
    await upgrades.erc1967.getImplementationAddress(upgraded.target);
  deploymentInfo.version = "3.0.0";
  deploymentInfo.lastUpgrade = new Date().toISOString();

  fs.writeFileSync(
    `./deployments/upgradeable-${network.name}.json`,
    JSON.stringify(deploymentInfo, null, 2)
  );

  return upgraded;
}

// 如果直接运行此脚本
if (require.main === module) {
  const command = process.argv[2];

  if (command === "deploy") {
    main()
      .then(() => process.exit(0))
      .catch((error) => {
        console.error(error);
        process.exit(1);
      });
  } else if (command === "upgrade-v2") {
    upgradeToV2()
      .then(() => process.exit(0))
      .catch((error) => {
        console.error(error);
        process.exit(1);
      });
  } else if (command === "upgrade-v3") {
    upgradeToV3()
      .then(() => process.exit(0))
      .catch((error) => {
        console.error(error);
        process.exit(1);
      });
  } else {
    console.log(
      "用法: npx hardhat run scripts/deploy-upgradeable.js --network <network> -- <command>"
    );
    console.log("命令: deploy, upgrade-v2, upgrade-v3");
  }
}

module.exports = { main, upgradeToV2, upgradeToV3 };
```

### 12.2.3 升级测试

```javascript
// test/Upgradeable.test.js
const { expect } = require("chai");
const { ethers, upgrades } = require("hardhat");

describe("可升级代币合约测试", function () {
  let token;
  let owner, user1, user2;

  beforeEach(async function () {
    [owner, user1, user2] = await ethers.getSigners();

    // 部署 V1
    const UpgradeableTokenV1 = await ethers.getContractFactory(
      "UpgradeableTokenV1"
    );
    token = await upgrades.deployProxy(
      UpgradeableTokenV1,
      ["Test Token", "TEST", owner.address],
      { kind: "uups", initializer: "initialize" }
    );
  });

  describe("V1 功能测试", function () {
    it("应该正确初始化", async function () {
      expect(await token.name()).to.equal("Test Token");
      expect(await token.symbol()).to.equal("TEST");
      expect(await token.owner()).to.equal(owner.address);
      expect(await token.version()).to.equal("1.0.0");
      expect(await token.totalSupply()).to.equal(ethers.parseEther("100000"));
    });

    it("所有者应该能铸造代币", async function () {
      await token.mint(user1.address, ethers.parseEther("1000"));
      expect(await token.balanceOf(user1.address)).to.equal(
        ethers.parseEther("1000")
      );
    });

    it("应该支持基本的 ERC20 功能", async function () {
      await token.mint(user1.address, ethers.parseEther("1000"));
      await token
        .connect(user1)
        .transfer(user2.address, ethers.parseEther("500"));

      expect(await token.balanceOf(user1.address)).to.equal(
        ethers.parseEther("500")
      );
      expect(await token.balanceOf(user2.address)).to.equal(
        ethers.parseEther("500")
      );
    });
  });

  describe("升级到 V2", function () {
    beforeEach(async function () {
      // 升级到 V2
      const UpgradeableTokenV2 = await ethers.getContractFactory(
        "UpgradeableTokenV2"
      );
      token = await upgrades.upgradeProxy(token.target, UpgradeableTokenV2);
      await token.initializeV2();
    });

    it("应该保持原有状态", async function () {
      expect(await token.name()).to.equal("Test Token");
      expect(await token.symbol()).to.equal("TEST");
      expect(await token.totalSupply()).to.equal(ethers.parseEther("100000"));
      expect(await token.version()).to.equal("2.0.0");
    });

    it("应该支持新的初始化状态", async function () {
      expect(await token.transfersEnabled()).to.equal(true);
      expect(await token.burnFeeRate()).to.equal(100);
    });

    it("应该支持黑名单功能", async function () {
      await token.mint(user1.address, ethers.parseEther("1000"));

      // 添加到黑名单
      await token.setBlacklisted(user1.address, true);
      expect(await token.blacklisted(user1.address)).to.equal(true);

      // 黑名单用户不能转账
      await expect(
        token.connect(user1).transfer(user2.address, ethers.parseEther("100"))
      ).to.be.revertedWith("From address is blacklisted");
    });

    it("应该支持转账开关", async function () {
      await token.mint(user1.address, ethers.parseEther("1000"));

      // 禁用转账
      await token.setTransfersEnabled(false);

      // 非铸造/销毁转账应该失败
      await expect(
        token.connect(user1).transfer(user2.address, ethers.parseEther("100"))
      ).to.be.revertedWith("Transfers disabled");

      // 铸造应该仍然可以
      await token.mint(user2.address, ethers.parseEther("100"));
      expect(await token.balanceOf(user2.address)).to.equal(
        ethers.parseEther("100")
      );
    });

    it("应该支持带费用的销毁", async function () {
      await token.mint(user1.address, ethers.parseEther("1000"));

      const initialOwnerBalance = await token.balanceOf(owner.address);
      const burnAmount = ethers.parseEther("100");
      const expectedFee = (burnAmount * 100n) / 10000n; // 1%
      const expectedBurn = burnAmount - expectedFee;

      await token.connect(user1).burnWithFee(burnAmount);

      expect(await token.balanceOf(user1.address)).to.equal(
        ethers.parseEther("1000") - burnAmount
      );
      expect(await token.balanceOf(owner.address)).to.equal(
        initialOwnerBalance + expectedFee
      );
    });

    it("应该支持批量转账", async function () {
      await token.mint(user1.address, ethers.parseEther("1000"));

      const recipients = [user2.address, owner.address];
      const amounts = [ethers.parseEther("100"), ethers.parseEther("200")];

      await token.connect(user1).batchTransfer(recipients, amounts);

      expect(await token.balanceOf(user2.address)).to.equal(
        ethers.parseEther("100")
      );
      expect(await token.balanceOf(owner.address)).to.be.gte(
        ethers.parseEther("100200")
      );
    });

    it("应该支持空投", async function () {
      const recipients = [user1.address, user2.address];
      const airdropAmount = ethers.parseEther("50");

      await token.airdrop(recipients, airdropAmount);

      expect(await token.balanceOf(user1.address)).to.equal(airdropAmount);
      expect(await token.balanceOf(user2.address)).to.equal(airdropAmount);
    });
  });

  describe("升级到 V3", function () {
    beforeEach(async function () {
      // 先升级到 V2
      const UpgradeableTokenV2 = await ethers.getContractFactory(
        "UpgradeableTokenV2"
      );
      token = await upgrades.upgradeProxy(token.target, UpgradeableTokenV2);
      await token.initializeV2();

      // 再升级到 V3
      const UpgradeableTokenV3 = await ethers.getContractFactory(
        "UpgradeableTokenV3"
      );
      token = await upgrades.upgradeProxy(token.target, UpgradeableTokenV3);
      await token.initializeV3();
    });

    it("应该保持所有之前的状态", async function () {
      expect(await token.version()).to.equal("3.0.0");
      expect(await token.transfersEnabled()).to.equal(true);
      expect(await token.burnFeeRate()).to.equal(100);
    });

    it("应该支持新的初始化状态", async function () {
      expect(await token.transferCooldown()).to.equal(300);
      expect(await token.maxTransferAmount()).to.equal(
        ethers.parseEther("10000")
      );
    });

    it("应该支持转账冷却时间", async function () {
      await token.mint(user1.address, ethers.parseEther("20000"));

      // 第一次转账
      await token
        .connect(user1)
        .transfer(user2.address, ethers.parseEther("1000"));

      // 立即再次转账应该失败
      await expect(
        token.connect(user1).transfer(user2.address, ethers.parseEther("1000"))
      ).to.be.revertedWith("Transfer cooldown not met");

      // 增加时间后应该成功
      await ethers.provider.send("evm_increaseTime", [301]);
      await ethers.provider.send("evm_mine");

      await token
        .connect(user1)
        .transfer(user2.address, ethers.parseEther("1000"));
    });

    it("应该支持最大转账限额", async function () {
      await token.mint(user1.address, ethers.parseEther("20000"));

      // 超过限额的转账应该失败
      await expect(
        token.connect(user1).transfer(user2.address, ethers.parseEther("15000"))
      ).to.be.revertedWith("Amount exceeds max transfer");

      // 在限额内的转账应该成功
      await token
        .connect(user1)
        .transfer(user2.address, ethers.parseEther("5000"));
    });

    it("VIP用户应该不受限制", async function () {
      await token.mint(user1.address, ethers.parseEther("20000"));
      await token.setVIPStatus(user1.address, true);

      // VIP用户可以超过限额转账
      await token
        .connect(user1)
        .transfer(user2.address, ethers.parseEther("15000"));
      expect(await token.balanceOf(user2.address)).to.equal(
        ethers.parseEther("15000")
      );

      // VIP用户可以立即再次转账
      await token
        .connect(user1)
        .transfer(user2.address, ethers.parseEther("1000"));
      expect(await token.balanceOf(user2.address)).to.equal(
        ethers.parseEther("16000")
      );
    });

    it("应该支持紧急转账", async function () {
      await token.mint(user1.address, ethers.parseEther("1000"));
      await token.setVIPStatus(user1.address, true);

      // VIP用户可以使用紧急转账
      await token
        .connect(user1)
        .emergencyTransfer(user2.address, ethers.parseEther("500"));
      expect(await token.balanceOf(user2.address)).to.equal(
        ethers.parseEther("500")
      );

      // 非VIP用户不能使用紧急转账
      await expect(
        token
          .connect(user2)
          .emergencyTransfer(user1.address, ethers.parseEther("100"))
      ).to.be.revertedWith("Not VIP user");
    });
  });

  describe("升级安全性", function () {
    it("只有所有者可以升级合约", async function () {
      const UpgradeableTokenV2 = await ethers.getContractFactory(
        "UpgradeableTokenV2"
      );

      await expect(
        upgrades.upgradeProxy(token.target, UpgradeableTokenV2.connect(user1))
      ).to.be.reverted;
    });

    it("升级后存储布局应该保持一致", async function () {
      // 在升级前设置一些状态
      await token.mint(user1.address, ethers.parseEther("1000"));
      await token
        .connect(user1)
        .transfer(user2.address, ethers.parseEther("500"));

      const beforeBalances = {
        user1: await token.balanceOf(user1.address),
        user2: await token.balanceOf(user2.address),
        total: await token.totalSupply(),
      };

      // 升级到 V2
      const UpgradeableTokenV2 = await ethers.getContractFactory(
        "UpgradeableTokenV2"
      );
      token = await upgrades.upgradeProxy(token.target, UpgradeableTokenV2);
      await token.initializeV2();

      // 验证状态保持不变
      expect(await token.balanceOf(user1.address)).to.equal(
        beforeBalances.user1
      );
      expect(await token.balanceOf(user2.address)).to.equal(
        beforeBalances.user2
      );
      expect(await token.totalSupply()).to.equal(beforeBalances.total);
    });
  });
});
```

## 12.3 存储布局管理

### 12.3.1 存储槽冲突处理

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * 存储布局规则：
 * 1. 新变量只能添加在现有变量之后
 * 2. 不能删除或重新排序现有变量
 * 3. 不能更改现有变量的类型
 * 4. 映射和数组可以安全扩展
 */

contract StorageV1 {
    uint256 public value1;           // 槽 0
    address public owner;            // 槽 1
    mapping(address => uint256) public balances; // 槽 2
    string public name;              // 槽 3
    bool public paused;              // 槽 4 (打包在一起)
    uint8 public decimals;           // 槽 4 (打包在一起)
}

// ✅ 正确的升级 - 只添加新变量
contract StorageV2 {
    uint256 public value1;           // 槽 0 - 保持不变
    address public owner;            // 槽 1 - 保持不变
    mapping(address => uint256) public balances; // 槽 2 - 保持不变
    string public name;              // 槽 3 - 保持不变
    bool public paused;              // 槽 4 - 保持不变
    uint8 public decimals;           // 槽 4 - 保持不变

    // 新变量添加在最后
    uint256 public newValue;         // 槽 5
    mapping(address => bool) public whitelist; // 槽 6
    address[] public admins;         // 槽 7
}

// ❌ 错误的升级示例
contract StorageV2Wrong {
    uint256 public value1;           // 槽 0
    uint256 public newValue;         // ❌ 插入新变量会打乱布局
    address public owner;            // 槽 2 - 现在变成了槽 2！
    mapping(address => uint256) public balances; // 槽 3 - 位置改变！
    string public name;              // 槽 4
    bool public paused;              // 槽 5
    uint8 public decimals;           // 槽 5
}

/**
 * 存储布局验证工具
 */
contract StorageLayoutValidator {
    // 使用事件记录存储布局
    event StorageLayout(string name, uint256 slot, string dataType);

    function logStorageLayout() public {
        emit StorageLayout("value1", 0, "uint256");
        emit StorageLayout("owner", 1, "address");
        emit StorageLayout("balances", 2, "mapping(address => uint256)");
        emit StorageLayout("name", 3, "string");
        emit StorageLayout("paused", 4, "bool");
        emit StorageLayout("decimals", 4, "uint8");
    }

    // 检查关键存储槽的值
    function validateStorageIntegrity() public view returns (bool) {
        uint256 slot0;
        address slot1;

        assembly {
            slot0 := sload(0)  // value1
            slot1 := sload(1)  // owner
        }

        // 这里可以添加更多验证逻辑
        return slot1 != address(0); // 简单检查：owner 不应该是零地址
    }
}
```

### 12.3.2 存储间隙技术

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * 使用存储间隙预留升级空间
 */
contract StorageGapExample {
    uint256 public value1;
    address public owner;
    mapping(address => uint256) public balances;

    // 预留50个存储槽用于未来升级
    uint256[50] private __gap;

    // 在未来的版本中，可以使用这些间隙
}

contract StorageGapExampleV2 {
    uint256 public value1;
    address public owner;
    mapping(address => uint256) public balances;

    // 使用了2个间隙槽
    uint256 public newValue1;
    uint256 public newValue2;

    // 剩余48个间隙槽
    uint256[48] private __gap;
}

/**
 * 复杂存储布局示例
 */
contract ComplexStorageLayout {
    // 基础状态变量
    uint256 public totalSupply;         // 槽 0
    string public name;                 // 槽 1
    string public symbol;               // 槽 2
    uint8 public decimals;              // 槽 3
    address public owner;               // 槽 3 (打包)
    bool public paused;                 // 槽 3 (打包)

    // 映射
    mapping(address => uint256) public balanceOf;           // 槽 4
    mapping(address => mapping(address => uint256)) public allowance; // 槽 5
    mapping(address => bool) public blacklisted;           // 槽 6

    // 数组
    address[] public holders;           // 槽 7

    // 嵌套结构
    struct UserInfo {
        uint256 balance;
        uint256 lastAction;
        bool isVIP;
        uint8 level;
    }
    mapping(address => UserInfo) public userInfo; // 槽 8

    // 预留间隙
    uint256[42] private __gap;

    function getStorageSlot(bytes32 slot) public view returns (bytes32 value) {
        assembly {
            value := sload(slot)
        }
    }

    function setStorageSlot(bytes32 slot, bytes32 value) public {
        require(msg.sender == owner, "Not owner");
        assembly {
            sstore(slot, value)
        }
    }

    // 获取映射存储位置
    function getMappingSlot(address key, uint256 mapSlot) public pure returns (bytes32) {
        return keccak256(abi.encode(key, mapSlot));
    }

    // 获取数组元素存储位置
    function getArrayElementSlot(uint256 index, uint256 arraySlot) public pure returns (bytes32) {
        bytes32 arrayDataSlot = keccak256(abi.encode(arraySlot));
        return bytes32(uint256(arrayDataSlot) + index);
    }
}
```

## 12.4 升级安全考虑

### 12.4.1 初始化安全

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

/**
 * 安全的初始化模式
 */
contract SecureInitialization is Initializable {
    address public owner;
    uint256 public value;
    bool public initialized;

    // 版本控制的初始化器
    uint8 private _initializationVersion;

    modifier onlyOwner() {
        require(msg.sender == owner, "Not owner");
        _;
    }

    // V1 初始化
    function initialize(address _owner, uint256 _value) public initializer {
        require(_owner != address(0), "Invalid owner");
        owner = _owner;
        value = _value;
        initialized = true;
        _initializationVersion = 1;
    }

    // V2 升级初始化
    function initializeV2(uint256 _newValue) public reinitializer(2) {
        require(_initializationVersion == 1, "Invalid previous version");
        value = _newValue;
        _initializationVersion = 2;
    }

    // V3 升级初始化
    function initializeV3() public reinitializer(3) {
        require(_initializationVersion == 2, "Invalid previous version");
        // V3 特定的初始化逻辑
        _initializationVersion = 3;
    }

    // 防止意外的重复初始化
    function getInitializationVersion() public view returns (uint8) {
        return _initializationVersion;
    }

    // 紧急重置（仅用于极端情况）
    function emergencyReset() public onlyOwner {
        require(block.timestamp > 0, "Emergency only"); // 额外的安全检查
        // 重置关键状态
        _initializationVersion = 0;
    }
}

/**
 * 构造函数禁用模式
 */
contract ConstructorDisabled {
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function _disableInitializers() internal {
        // 这会阻止任何人调用 initialize 函数
        // 通过将初始化状态设置为最大值
    }
}
```

### 12.4.2 权限控制和治理

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

/**
 * 基于治理的升级控制
 */
contract GovernedUpgradeable is AccessControlUpgradeable, UUPSUpgradeable {
    bytes32 public constant UPGRADER_ROLE = keccak256("UPGRADER_ROLE");
    bytes32 public constant GOVERNANCE_ROLE = keccak256("GOVERNANCE_ROLE");

    uint256 public constant UPGRADE_DELAY = 2 days;

    struct UpgradeProposal {
        address newImplementation;
        bytes data;
        uint256 proposedAt;
        bool executed;
        uint256 votes;
        mapping(address => bool) voted;
    }

    mapping(bytes32 => UpgradeProposal) public upgradeProposals;
    bytes32[] public proposalIds;

    event UpgradeProposed(
        bytes32 indexed proposalId,
        address newImplementation,
        bytes data,
        address proposer
    );

    event UpgradeExecuted(
        bytes32 indexed proposalId,
        address newImplementation
    );

    event VoteCast(
        bytes32 indexed proposalId,
        address voter,
        uint256 weight
    );

    function initialize(address admin) public initializer {
        __AccessControl_init();
        __UUPSUpgradeable_init();

        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(GOVERNANCE_ROLE, admin);
        _grantRole(UPGRADER_ROLE, admin);
    }

    // 提议升级
    function proposeUpgrade(
        address newImplementation,
        bytes calldata data
    ) public onlyRole(UPGRADER_ROLE) returns (bytes32 proposalId) {
        require(newImplementation != address(0), "Invalid implementation");

        proposalId = keccak256(
            abi.encodePacked(newImplementation, data, block.timestamp)
        );

        UpgradeProposal storage proposal = upgradeProposals[proposalId];
        proposal.newImplementation = newImplementation;
        proposal.data = data;
        proposal.proposedAt = block.timestamp;

        proposalIds.push(proposalId);

        emit UpgradeProposed(proposalId, newImplementation, data, msg.sender);
    }

    // 投票支持升级
    function voteForUpgrade(bytes32 proposalId) public onlyRole(GOVERNANCE_ROLE) {
        UpgradeProposal storage proposal = upgradeProposals[proposalId];
        require(proposal.newImplementation != address(0), "Proposal not found");
        require(!proposal.executed, "Already executed");
        require(!proposal.voted[msg.sender], "Already voted");

        proposal.voted[msg.sender] = true;
        proposal.votes += getVotingWeight(msg.sender);

        emit VoteCast(proposalId, msg.sender, getVotingWeight(msg.sender));
    }

    // 执行升级
    function executeUpgrade(bytes32 proposalId) public onlyRole(UPGRADER_ROLE) {
        UpgradeProposal storage proposal = upgradeProposals[proposalId];
        require(proposal.newImplementation != address(0), "Proposal not found");
        require(!proposal.executed, "Already executed");
        require(
            block.timestamp >= proposal.proposedAt + UPGRADE_DELAY,
            "Upgrade delay not met"
        );
        require(proposal.votes >= getRequiredVotes(), "Insufficient votes");

        proposal.executed = true;

        // 执行升级
        _upgradeToAndCall(proposal.newImplementation, proposal.data, false);

        emit UpgradeExecuted(proposalId, proposal.newImplementation);
    }

    // 获取投票权重（可以基于代币持有量等）
    function getVotingWeight(address voter) public view virtual returns (uint256) {
        // 简单实现：每个治理角色成员权重为1
        return hasRole(GOVERNANCE_ROLE, voter) ? 1 : 0;
    }

    // 获取所需票数
    function getRequiredVotes() public view virtual returns (uint256) {
        uint256 governanceMembers = getRoleMemberCount(GOVERNANCE_ROLE);
        return (governanceMembers * 2) / 3 + 1; // 2/3 多数
    }

    // 紧急升级（多重签名）
    function emergencyUpgrade(
        address newImplementation,
        bytes calldata data
    ) public onlyRole(DEFAULT_ADMIN_ROLE) {
        require(newImplementation != address(0), "Invalid implementation");

        // 需要至少2个管理员确认
        // 这里简化实现，实际应该使用多签
        _upgradeToAndCall(newImplementation, data, false);
    }

    function _authorizeUpgrade(address newImplementation) internal override onlyRole(UPGRADER_ROLE) {
        // 额外的升级验证逻辑
    }

    // 获取提议详情
    function getProposal(bytes32 proposalId) public view returns (
        address newImplementation,
        bytes memory data,
        uint256 proposedAt,
        bool executed,
        uint256 votes
    ) {
        UpgradeProposal storage proposal = upgradeProposals[proposalId];
        return (
            proposal.newImplementation,
            proposal.data,
            proposal.proposedAt,
            proposal.executed,
            proposal.votes
        );
    }

    // 获取所有提议
    function getAllProposals() public view returns (bytes32[] memory) {
        return proposalIds;
    }
}
```

## 12.5 章节总结

在本章的上篇中，我们深入学习了：

1. **升级必要性**：智能合约不可变性的挑战和升级需求
2. **代理模式**：透明代理和 UUPS 代理的原理和实现
3. **OpenZeppelin Upgrades**：标准化的升级框架使用
4. **存储布局**：升级时的存储兼容性和间隙技术
5. **安全考虑**：初始化安全、权限控制、治理机制

### 关键要点

1. **存储兼容性**：升级时必须保持存储布局一致性
2. **初始化控制**：使用 reinitializer 确保安全的版本升级
3. **权限管理**：实施适当的升级权限控制机制
4. **测试充分**：每次升级都需要完整的测试验证
5. **治理机制**：重要升级应该通过社区治理决定

### 升级最佳实践

- **渐进升级**：分阶段实施复杂升级
- **向后兼容**：确保新版本与旧接口兼容
- **充分测试**：使用 Fork 测试验证升级效果
- **文档完整**：记录每次升级的变更内容
- **应急计划**：准备升级回滚和紧急修复方案

代理模式是智能合约发展的重要技术，正确使用可以让 dApp 具备持续演进的能力。

继续学习：[第十二章：高级特性与模式（下篇）- 工厂模式与 Oracle 集成](12_advanced_features_patterns_part2.md)
