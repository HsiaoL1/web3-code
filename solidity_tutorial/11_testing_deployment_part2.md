# 第十一章：测试与部署（下篇）- 集成测试与部署策略

## 本章学习目标

- 掌握智能合约集成测试技巧
- 学会测试合约间交互
- 了解模拟和存根技术
- 掌握主网部署流程
- 学会合约验证和监控
- 理解部署后的运维策略

## 11.7 集成测试

### 11.7.1 多合约交互测试

```solidity
// contracts/DEX.sol - 简单的去中心化交易所
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract SimpleDEX is ReentrancyGuard, Ownable {
    struct Pool {
        address tokenA;
        address tokenB;
        uint256 reserveA;
        uint256 reserveB;
        uint256 totalLiquidity;
        mapping(address => uint256) liquidity;
        bool exists;
    }

    mapping(bytes32 => Pool) public pools;
    mapping(address => bytes32[]) public userPools;

    uint256 public constant MINIMUM_LIQUIDITY = 10**3;
    uint256 public feeRate = 30; // 0.3%

    event PoolCreated(
        address indexed tokenA,
        address indexed tokenB,
        bytes32 indexed poolId
    );

    event LiquidityAdded(
        bytes32 indexed poolId,
        address indexed provider,
        uint256 amountA,
        uint256 amountB,
        uint256 liquidity
    );

    event LiquidityRemoved(
        bytes32 indexed poolId,
        address indexed provider,
        uint256 amountA,
        uint256 amountB,
        uint256 liquidity
    );

    event TokenSwapped(
        bytes32 indexed poolId,
        address indexed user,
        address tokenIn,
        address tokenOut,
        uint256 amountIn,
        uint256 amountOut
    );

    constructor() Ownable(msg.sender) {}

    function createPool(
        address tokenA,
        address tokenB
    ) external returns (bytes32 poolId) {
        require(tokenA != tokenB, "Identical tokens");
        require(tokenA != address(0) && tokenB != address(0), "Zero address");

        // 确保 tokenA < tokenB（标准化）
        if (tokenA > tokenB) {
            (tokenA, tokenB) = (tokenB, tokenA);
        }

        poolId = keccak256(abi.encodePacked(tokenA, tokenB));
        require(!pools[poolId].exists, "Pool already exists");

        Pool storage pool = pools[poolId];
        pool.tokenA = tokenA;
        pool.tokenB = tokenB;
        pool.exists = true;

        emit PoolCreated(tokenA, tokenB, poolId);
    }

    function addLiquidity(
        bytes32 poolId,
        uint256 amountADesired,
        uint256 amountBDesired,
        uint256 amountAMin,
        uint256 amountBMin
    ) external nonReentrant returns (uint256 amountA, uint256 amountB, uint256 liquidity) {
        Pool storage pool = pools[poolId];
        require(pool.exists, "Pool does not exist");

        (amountA, amountB) = _calculateOptimalAmounts(
            pool,
            amountADesired,
            amountBDesired,
            amountAMin,
            amountBMin
        );

        // 转移代币
        IERC20(pool.tokenA).transferFrom(msg.sender, address(this), amountA);
        IERC20(pool.tokenB).transferFrom(msg.sender, address(this), amountB);

        // 计算流动性
        if (pool.totalLiquidity == 0) {
            liquidity = _sqrt(amountA * amountB) - MINIMUM_LIQUIDITY;
            pool.totalLiquidity = MINIMUM_LIQUIDITY; // 锁定最小流动性
        } else {
            liquidity = _min(
                (amountA * pool.totalLiquidity) / pool.reserveA,
                (amountB * pool.totalLiquidity) / pool.reserveB
            );
        }

        require(liquidity > 0, "Insufficient liquidity minted");

        pool.liquidity[msg.sender] += liquidity;
        pool.totalLiquidity += liquidity;
        pool.reserveA += amountA;
        pool.reserveB += amountB;

        userPools[msg.sender].push(poolId);

        emit LiquidityAdded(poolId, msg.sender, amountA, amountB, liquidity);
    }

    function removeLiquidity(
        bytes32 poolId,
        uint256 liquidity,
        uint256 amountAMin,
        uint256 amountBMin
    ) external nonReentrant returns (uint256 amountA, uint256 amountB) {
        Pool storage pool = pools[poolId];
        require(pool.exists, "Pool does not exist");
        require(pool.liquidity[msg.sender] >= liquidity, "Insufficient liquidity");

        amountA = (liquidity * pool.reserveA) / pool.totalLiquidity;
        amountB = (liquidity * pool.reserveB) / pool.totalLiquidity;

        require(amountA >= amountAMin, "Insufficient amount A");
        require(amountB >= amountBMin, "Insufficient amount B");

        pool.liquidity[msg.sender] -= liquidity;
        pool.totalLiquidity -= liquidity;
        pool.reserveA -= amountA;
        pool.reserveB -= amountB;

        IERC20(pool.tokenA).transfer(msg.sender, amountA);
        IERC20(pool.tokenB).transfer(msg.sender, amountB);

        emit LiquidityRemoved(poolId, msg.sender, amountA, amountB, liquidity);
    }

    function swapExactTokensForTokens(
        bytes32 poolId,
        address tokenIn,
        uint256 amountIn,
        uint256 amountOutMin
    ) external nonReentrant returns (uint256 amountOut) {
        Pool storage pool = pools[poolId];
        require(pool.exists, "Pool does not exist");
        require(
            tokenIn == pool.tokenA || tokenIn == pool.tokenB,
            "Invalid token"
        );

        bool isTokenA = tokenIn == pool.tokenA;
        address tokenOut = isTokenA ? pool.tokenB : pool.tokenA;

        amountOut = getAmountOut(
            amountIn,
            isTokenA ? pool.reserveA : pool.reserveB,
            isTokenA ? pool.reserveB : pool.reserveA
        );

        require(amountOut >= amountOutMin, "Insufficient output amount");

        IERC20(tokenIn).transferFrom(msg.sender, address(this), amountIn);

        if (isTokenA) {
            pool.reserveA += amountIn;
            pool.reserveB -= amountOut;
        } else {
            pool.reserveB += amountIn;
            pool.reserveA -= amountOut;
        }

        IERC20(tokenOut).transfer(msg.sender, amountOut);

        emit TokenSwapped(poolId, msg.sender, tokenIn, tokenOut, amountIn, amountOut);
    }

    function getAmountOut(
        uint256 amountIn,
        uint256 reserveIn,
        uint256 reserveOut
    ) public view returns (uint256 amountOut) {
        require(amountIn > 0, "Insufficient input amount");
        require(reserveIn > 0 && reserveOut > 0, "Insufficient liquidity");

        uint256 amountInWithFee = amountIn * (10000 - feeRate);
        uint256 numerator = amountInWithFee * reserveOut;
        uint256 denominator = (reserveIn * 10000) + amountInWithFee;
        amountOut = numerator / denominator;
    }

    function _calculateOptimalAmounts(
        Pool storage pool,
        uint256 amountADesired,
        uint256 amountBDesired,
        uint256 amountAMin,
        uint256 amountBMin
    ) internal view returns (uint256 amountA, uint256 amountB) {
        if (pool.reserveA == 0 && pool.reserveB == 0) {
            (amountA, amountB) = (amountADesired, amountBDesired);
        } else {
            uint256 amountBOptimal = (amountADesired * pool.reserveB) / pool.reserveA;
            if (amountBOptimal <= amountBDesired) {
                require(amountBOptimal >= amountBMin, "Insufficient B amount");
                (amountA, amountB) = (amountADesired, amountBOptimal);
            } else {
                uint256 amountAOptimal = (amountBDesired * pool.reserveA) / pool.reserveB;
                require(amountAOptimal <= amountADesired && amountAOptimal >= amountAMin, "Insufficient A amount");
                (amountA, amountB) = (amountAOptimal, amountBDesired);
            }
        }
    }

    function _sqrt(uint256 y) internal pure returns (uint256 z) {
        if (y > 3) {
            z = y;
            uint256 x = y / 2 + 1;
            while (x < z) {
                z = x;
                x = (y / x + x) / 2;
            }
        } else if (y != 0) {
            z = 1;
        }
    }

    function _min(uint256 x, uint256 y) internal pure returns (uint256) {
        return x < y ? x : y;
    }

    // 查询函数
    function getPoolInfo(bytes32 poolId) external view returns (
        address tokenA,
        address tokenB,
        uint256 reserveA,
        uint256 reserveB,
        uint256 totalLiquidity
    ) {
        Pool storage pool = pools[poolId];
        return (
            pool.tokenA,
            pool.tokenB,
            pool.reserveA,
            pool.reserveB,
            pool.totalLiquidity
        );
    }

    function getUserLiquidity(bytes32 poolId, address user) external view returns (uint256) {
        return pools[poolId].liquidity[user];
    }
}
```

### 11.7.2 集成测试用例

```javascript
// test/DEX.integration.test.js
const { expect } = require("chai");
const { ethers } = require("hardhat");
const { loadFixture } = require("@nomicfoundation/hardhat-network-helpers");

describe("SimpleDEX 集成测试", function () {
  async function deployDEXFixture() {
    const [owner, trader1, trader2, liquidityProvider] =
      await ethers.getSigners();

    // 部署代币
    const TestToken = await ethers.getContractFactory("TestToken");
    const tokenA = await TestToken.deploy(
      "Token A",
      "TKA",
      ethers.parseEther("10000")
    );
    const tokenB = await TestToken.deploy(
      "Token B",
      "TKB",
      ethers.parseEther("10000")
    );

    // 部署 DEX
    const SimpleDEX = await ethers.getContractFactory("SimpleDEX");
    const dex = await SimpleDEX.deploy();

    // 为测试账户分发代币
    await tokenA.transfer(trader1.address, ethers.parseEther("1000"));
    await tokenA.transfer(trader2.address, ethers.parseEther("1000"));
    await tokenA.transfer(liquidityProvider.address, ethers.parseEther("2000"));

    await tokenB.transfer(trader1.address, ethers.parseEther("1000"));
    await tokenB.transfer(trader2.address, ethers.parseEther("1000"));
    await tokenB.transfer(liquidityProvider.address, ethers.parseEther("2000"));

    return {
      dex,
      tokenA,
      tokenB,
      owner,
      trader1,
      trader2,
      liquidityProvider,
    };
  }

  describe("完整交易流程", function () {
    it("应该支持完整的流动性提供和交易流程", async function () {
      const { dex, tokenA, tokenB, liquidityProvider, trader1 } =
        await loadFixture(deployDEXFixture);

      // 1. 创建交易对
      const tx1 = await dex.createPool(tokenA.target, tokenB.target);
      const receipt1 = await tx1.wait();

      // 从事件中获取 poolId
      const poolCreatedEvent = receipt1.logs.find(
        (log) =>
          log.topics[0] === dex.interface.getEvent("PoolCreated").topicHash
      );
      const poolId = poolCreatedEvent.topics[3];

      // 2. 流动性提供者添加流动性
      const liquidityAmountA = ethers.parseEther("100");
      const liquidityAmountB = ethers.parseEther("200");

      await tokenA
        .connect(liquidityProvider)
        .approve(dex.target, liquidityAmountA);
      await tokenB
        .connect(liquidityProvider)
        .approve(dex.target, liquidityAmountB);

      await expect(
        dex
          .connect(liquidityProvider)
          .addLiquidity(
            poolId,
            liquidityAmountA,
            liquidityAmountB,
            liquidityAmountA,
            liquidityAmountB
          )
      ).to.emit(dex, "LiquidityAdded");

      // 验证池子状态
      const poolInfo = await dex.getPoolInfo(poolId);
      expect(poolInfo.reserveA).to.equal(liquidityAmountA);
      expect(poolInfo.reserveB).to.equal(liquidityAmountB);

      // 3. 交易者进行代币交换
      const swapAmountIn = ethers.parseEther("10");

      await tokenA.connect(trader1).approve(dex.target, swapAmountIn);

      const expectedAmountOut = await dex.getAmountOut(
        swapAmountIn,
        liquidityAmountA,
        liquidityAmountB
      );

      const trader1BalanceBBefore = await tokenB.balanceOf(trader1.address);

      await expect(
        dex
          .connect(trader1)
          .swapExactTokensForTokens(
            poolId,
            tokenA.target,
            swapAmountIn,
            expectedAmountOut
          )
      ).to.emit(dex, "TokenSwapped");

      const trader1BalanceBAfter = await tokenB.balanceOf(trader1.address);
      expect(trader1BalanceBAfter - trader1BalanceBBefore).to.equal(
        expectedAmountOut
      );

      // 4. 验证池子储备金变化
      const poolInfoAfter = await dex.getPoolInfo(poolId);
      expect(poolInfoAfter.reserveA).to.equal(liquidityAmountA + swapAmountIn);
      expect(poolInfoAfter.reserveB).to.equal(
        liquidityAmountB - expectedAmountOut
      );
    });

    it("应该正确计算多次交易后的价格影响", async function () {
      const { dex, tokenA, tokenB, liquidityProvider, trader1, trader2 } =
        await loadFixture(deployDEXFixture);

      // 创建池子并添加流动性
      const poolId = await createPoolAndAddLiquidity(
        dex,
        tokenA,
        tokenB,
        liquidityProvider,
        ethers.parseEther("1000"),
        ethers.parseEther("2000")
      );

      // 第一次交易
      const firstSwapAmount = ethers.parseEther("50");
      await tokenA.connect(trader1).approve(dex.target, firstSwapAmount);

      const expectedOut1 = await dex.getAmountOut(
        firstSwapAmount,
        ethers.parseEther("1000"),
        ethers.parseEther("2000")
      );

      await dex
        .connect(trader1)
        .swapExactTokensForTokens(
          poolId,
          tokenA.target,
          firstSwapAmount,
          expectedOut1
        );

      // 获取交易后的储备金
      const poolInfoAfter1 = await dex.getPoolInfo(poolId);

      // 第二次交易（价格已经改变）
      const secondSwapAmount = ethers.parseEther("50");
      await tokenA.connect(trader2).approve(dex.target, secondSwapAmount);

      const expectedOut2 = await dex.getAmountOut(
        secondSwapAmount,
        poolInfoAfter1.reserveA,
        poolInfoAfter1.reserveB
      );

      // 第二次交易的输出应该小于第一次（价格影响）
      expect(expectedOut2).to.be.lt(expectedOut1);

      await dex
        .connect(trader2)
        .swapExactTokensForTokens(
          poolId,
          tokenA.target,
          secondSwapAmount,
          expectedOut2
        );
    });
  });

  describe("流动性管理", function () {
    it("应该正确计算和分配流动性代币", async function () {
      const { dex, tokenA, tokenB, liquidityProvider, trader1 } =
        await loadFixture(deployDEXFixture);

      const poolId = await createPoolAndAddLiquidity(
        dex,
        tokenA,
        tokenB,
        liquidityProvider,
        ethers.parseEther("100"),
        ethers.parseEther("200")
      );

      // 获取流动性提供者的流动性代币数量
      const liquidity1 = await dex.getUserLiquidity(
        poolId,
        liquidityProvider.address
      );
      expect(liquidity1).to.be.gt(0);

      // 另一个用户添加流动性
      const additionalAmountA = ethers.parseEther("50");
      const additionalAmountB = ethers.parseEther("100");

      await tokenA.connect(trader1).approve(dex.target, additionalAmountA);
      await tokenB.connect(trader1).approve(dex.target, additionalAmountB);

      await dex
        .connect(trader1)
        .addLiquidity(
          poolId,
          additionalAmountA,
          additionalAmountB,
          additionalAmountA,
          additionalAmountB
        );

      const liquidity2 = await dex.getUserLiquidity(poolId, trader1.address);

      // 验证流动性代币分配比例
      // trader1 添加了 50% 的流动性，应该获得相应比例的流动性代币
      expect(liquidity2).to.be.approximately(
        liquidity1 / 2n,
        ethers.parseEther("0.1")
      );
    });

    it("应该支持移除流动性", async function () {
      const { dex, tokenA, tokenB, liquidityProvider } = await loadFixture(
        deployDEXFixture
      );

      const poolId = await createPoolAndAddLiquidity(
        dex,
        tokenA,
        tokenB,
        liquidityProvider,
        ethers.parseEther("200"),
        ethers.parseEther("400")
      );

      const initialBalanceA = await tokenA.balanceOf(liquidityProvider.address);
      const initialBalanceB = await tokenB.balanceOf(liquidityProvider.address);

      const userLiquidity = await dex.getUserLiquidity(
        poolId,
        liquidityProvider.address
      );
      const halfLiquidity = userLiquidity / 2n;

      // 移除一半流动性
      await expect(
        dex
          .connect(liquidityProvider)
          .removeLiquidity(poolId, halfLiquidity, 0, 0)
      ).to.emit(dex, "LiquidityRemoved");

      const finalBalanceA = await tokenA.balanceOf(liquidityProvider.address);
      const finalBalanceB = await tokenB.balanceOf(liquidityProvider.address);

      // 验证收到的代币数量
      expect(finalBalanceA).to.be.gt(initialBalanceA);
      expect(finalBalanceB).to.be.gt(initialBalanceB);

      // 验证剩余流动性
      const remainingLiquidity = await dex.getUserLiquidity(
        poolId,
        liquidityProvider.address
      );
      expect(remainingLiquidity).to.be.approximately(halfLiquidity, 1n);
    });
  });

  describe("错误情况处理", function () {
    it("应该在滑点过大时拒绝交易", async function () {
      const { dex, tokenA, tokenB, liquidityProvider, trader1 } =
        await loadFixture(deployDEXFixture);

      const poolId = await createPoolAndAddLiquidity(
        dex,
        tokenA,
        tokenB,
        liquidityProvider,
        ethers.parseEther("100"),
        ethers.parseEther("200")
      );

      const swapAmount = ethers.parseEther("10");
      const expectedOut = await dex.getAmountOut(
        swapAmount,
        ethers.parseEther("100"),
        ethers.parseEther("200")
      );

      await tokenA.connect(trader1).approve(dex.target, swapAmount);

      // 设置过高的最小输出期望
      const unreasonableMinOut = expectedOut * 2n;

      await expect(
        dex
          .connect(trader1)
          .swapExactTokensForTokens(
            poolId,
            tokenA.target,
            swapAmount,
            unreasonableMinOut
          )
      ).to.be.revertedWith("Insufficient output amount");
    });

    it("应该在流动性不足时拒绝大额交易", async function () {
      const { dex, tokenA, tokenB, liquidityProvider, trader1 } =
        await loadFixture(deployDEXFixture);

      const poolId = await createPoolAndAddLiquidity(
        dex,
        tokenA,
        tokenB,
        liquidityProvider,
        ethers.parseEther("100"),
        ethers.parseEther("200")
      );

      // 尝试交换超过池子储备的数量
      const excessiveAmount = ethers.parseEther("500");

      await tokenA.connect(trader1).approve(dex.target, excessiveAmount);

      await expect(
        dex
          .connect(trader1)
          .swapExactTokensForTokens(poolId, tokenA.target, excessiveAmount, 0)
      ).to.be.reverted; // 会因为计算问题或储备不足而失败
    });
  });

  // 辅助函数
  async function createPoolAndAddLiquidity(
    dex,
    tokenA,
    tokenB,
    provider,
    amountA,
    amountB
  ) {
    // 创建池子
    const tx = await dex.createPool(tokenA.target, tokenB.target);
    const receipt = await tx.wait();

    const poolCreatedEvent = receipt.logs.find(
      (log) => log.topics[0] === dex.interface.getEvent("PoolCreated").topicHash
    );
    const poolId = poolCreatedEvent.topics[3];

    // 添加流动性
    await tokenA.connect(provider).approve(dex.target, amountA);
    await tokenB.connect(provider).approve(dex.target, amountB);

    await dex
      .connect(provider)
      .addLiquidity(poolId, amountA, amountB, amountA, amountB);

    return poolId;
  }
});
```

### 11.7.3 Fork 测试

```javascript
// test/Fork.test.js - 使用主网分叉进行测试
const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("Fork 测试", function () {
  let uniswapRouter;
  let weth;
  let usdc;
  let whale; // 持有大量代币的账户

  before(async function () {
    // 这个测试需要在 hardhat.config.js 中配置主网分叉
    // networks: {
    //   hardhat: {
    //     forking: {
    //       url: "https://eth-mainnet.alchemyapi.io/v2/YOUR_KEY",
    //       blockNumber: 18000000 // 可选：固定区块高度
    //     }
    //   }
    // }

    // 主网合约地址
    const UNISWAP_V2_ROUTER = "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D";
    const WETH_ADDRESS = "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2";
    const USDC_ADDRESS = "0xA0b86a33E6417bFbB1eED4cD4d20D12d6b5f62c2";
    const WHALE_ADDRESS = "0x28C6c06298d514Db089934071355E5743bf21d60"; // Binance hot wallet

    // 获取合约实例
    uniswapRouter = await ethers.getContractAt(
      [
        "function swapExactTokensForTokens(uint amountIn, uint amountOutMin, address[] calldata path, address to, uint deadline) external returns (uint[] memory amounts)",
        "function getAmountsOut(uint amountIn, address[] calldata path) external view returns (uint[] memory amounts)",
      ],
      UNISWAP_V2_ROUTER
    );

    weth = await ethers.getContractAt("IERC20", WETH_ADDRESS);
    usdc = await ethers.getContractAt("IERC20", USDC_ADDRESS);

    // 模拟鲸鱼账户
    await network.provider.request({
      method: "hardhat_impersonateAccount",
      params: [WHALE_ADDRESS],
    });

    whale = await ethers.getSigner(WHALE_ADDRESS);
  });

  it("应该能够在 Uniswap 上进行真实交易", async function () {
    const swapAmount = ethers.parseEther("1"); // 1 WETH

    // 检查鲸鱼账户的 WETH 余额
    const wethBalance = await weth.balanceOf(whale.address);
    expect(wethBalance).to.be.gte(swapAmount);

    // 获取预期输出
    const path = [weth.target, usdc.target];
    const amounts = await uniswapRouter.getAmountsOut(swapAmount, path);
    const expectedUsdcOut = amounts[1];

    // 记录交易前余额
    const usdcBalanceBefore = await usdc.balanceOf(whale.address);

    // 批准和交换
    await weth.connect(whale).approve(uniswapRouter.target, swapAmount);

    await uniswapRouter.connect(whale).swapExactTokensForTokens(
      swapAmount,
      (expectedUsdcOut * 95n) / 100n, // 5% 滑点
      path,
      whale.address,
      Math.floor(Date.now() / 1000) + 300 // 5 分钟后过期
    );

    // 验证交易结果
    const usdcBalanceAfter = await usdc.balanceOf(whale.address);
    const usdcReceived = usdcBalanceAfter - usdcBalanceBefore;

    expect(usdcReceived).to.be.gte((expectedUsdcOut * 95n) / 100n);
    expect(usdcReceived).to.be.lte(expectedUsdcOut);
  });

  after(async function () {
    // 停止模拟账户
    await network.provider.request({
      method: "hardhat_stopImpersonatingAccount",
      params: [WHALE_ADDRESS],
    });
  });
});
```

## 11.8 部署策略

### 11.8.1 部署脚本

```javascript
// scripts/deploy.js
const { ethers, network } = require("hardhat");
const fs = require("fs");

async function main() {
  console.log("开始部署到网络:", network.name);

  const [deployer] = await ethers.getSigners();
  console.log("部署账户:", deployer.address);
  console.log(
    "账户余额:",
    ethers.formatEther(await deployer.provider.getBalance(deployer.address))
  );

  // 部署参数
  const deployParams = {
    name: "Production Token",
    symbol: "PROD",
    initialSupply: ethers.parseEther("1000000"),
  };

  console.log("部署参数:", deployParams);

  // 部署 Token 合约
  console.log("\n正在部署 TestToken...");
  const TestToken = await ethers.getContractFactory("TestToken");
  const token = await TestToken.deploy(
    deployParams.name,
    deployParams.symbol,
    deployParams.initialSupply
  );

  await token.waitForDeployment();
  console.log("TestToken 部署成功:", token.target);

  // 部署 DEX 合约
  console.log("\n正在部署 SimpleDEX...");
  const SimpleDEX = await ethers.getContractFactory("SimpleDEX");
  const dex = await SimpleDEX.deploy();

  await dex.waitForDeployment();
  console.log("SimpleDEX 部署成功:", dex.target);

  // 验证部署
  console.log("\n验证部署...");
  const tokenName = await token.name();
  const tokenSymbol = await token.symbol();
  const totalSupply = await token.totalSupply();
  const owner = await token.owner();

  console.log("Token 名称:", tokenName);
  console.log("Token 符号:", tokenSymbol);
  console.log("总供应量:", ethers.formatEther(totalSupply));
  console.log("所有者:", owner);

  // 保存部署信息
  const deploymentInfo = {
    network: network.name,
    timestamp: new Date().toISOString(),
    deployer: deployer.address,
    contracts: {
      TestToken: {
        address: token.target,
        constructorArgs: [
          deployParams.name,
          deployParams.symbol,
          deployParams.initialSupply,
        ],
      },
      SimpleDEX: {
        address: dex.target,
        constructorArgs: [],
      },
    },
    gasUsed: {
      // 这里可以记录实际的 gas 使用量
    },
  };

  // 保存到文件
  const deploymentsDir = "./deployments";
  if (!fs.existsSync(deploymentsDir)) {
    fs.mkdirSync(deploymentsDir);
  }

  fs.writeFileSync(
    `${deploymentsDir}/${network.name}.json`,
    JSON.stringify(deploymentInfo, null, 2)
  );

  console.log(`\n部署信息已保存到 ${deploymentsDir}/${network.name}.json`);

  // 如果是测试网，进行一些初始化操作
  if (network.name !== "mainnet") {
    console.log("\n执行初始化操作...");

    // 设置铸造价格
    await token.updateMintPrice(ethers.parseEther("0.002"));
    console.log("铸造价格已更新");

    // 添加一些初始铸造者
    const initialMinters = [
      // 添加一些地址
    ];

    for (const minter of initialMinters) {
      await token.addMinter(minter);
      console.log("已添加铸造者:", minter);
    }
  }

  console.log("\n部署完成!");

  // 返回合约实例供其他脚本使用
  return { token, dex, deploymentInfo };
}

// 如果直接运行此脚本
if (require.main === module) {
  main()
    .then(() => process.exit(0))
    .catch((error) => {
      console.error(error);
      process.exit(1);
    });
}

module.exports = main;
```

### 11.8.2 多网络部署配置

```javascript
// hardhat.config.js - 多网络配置
require("@nomicfoundation/hardhat-toolbox");
require("@nomicfoundation/hardhat-verify");
require("dotenv").config();

const PRIVATE_KEY = process.env.PRIVATE_KEY || "";
const INFURA_KEY = process.env.INFURA_KEY || "";
const ETHERSCAN_API_KEY = process.env.ETHERSCAN_API_KEY || "";

module.exports = {
  solidity: {
    version: "0.8.20",
    settings: {
      optimizer: {
        enabled: true,
        runs: 200,
      },
    },
  },

  networks: {
    // 本地开发网络
    localhost: {
      url: "http://127.0.0.1:8545",
      chainId: 31337,
      accounts: [PRIVATE_KEY],
    },

    // 测试网络
    sepolia: {
      url: `https://sepolia.infura.io/v3/${INFURA_KEY}`,
      chainId: 11155111,
      accounts: [PRIVATE_KEY],
      gas: "auto",
      gasPrice: "auto",
    },

    goerli: {
      url: `https://goerli.infura.io/v3/${INFURA_KEY}`,
      chainId: 5,
      accounts: [PRIVATE_KEY],
      gas: "auto",
      gasPrice: "auto",
    },

    // 主网
    mainnet: {
      url: `https://mainnet.infura.io/v3/${INFURA_KEY}`,
      chainId: 1,
      accounts: [PRIVATE_KEY],
      gas: "auto",
      gasPrice: "auto",
    },

    // 其他网络
    polygon: {
      url: "https://rpc-mainnet.maticvigil.com/",
      chainId: 137,
      accounts: [PRIVATE_KEY],
    },

    bsc: {
      url: "https://bsc-dataseed.binance.org/",
      chainId: 56,
      accounts: [PRIVATE_KEY],
    },
  },

  // 合约验证
  etherscan: {
    apiKey: {
      mainnet: ETHERSCAN_API_KEY,
      sepolia: ETHERSCAN_API_KEY,
      goerli: ETHERSCAN_API_KEY,
      polygon: process.env.POLYGONSCAN_API_KEY,
      bsc: process.env.BSCSCAN_API_KEY,
    },
  },

  // Gas 报告
  gasReporter: {
    enabled: process.env.REPORT_GAS === "true",
    currency: "USD",
    gasPrice: 20,
    token: "ETH",
    coinmarketcap: process.env.COINMARKETCAP_API_KEY,
  },
};
```

### 11.8.3 部署验证脚本

```javascript
// scripts/verify.js - 合约验证
const { ethers, run, network } = require("hardhat");
const fs = require("fs");

async function main() {
  console.log("开始验证合约，网络:", network.name);

  // 读取部署信息
  const deploymentPath = `./deployments/${network.name}.json`;
  if (!fs.existsSync(deploymentPath)) {
    throw new Error(`找不到部署文件: ${deploymentPath}`);
  }

  const deploymentInfo = JSON.parse(fs.readFileSync(deploymentPath, "utf8"));

  console.log("部署信息:", deploymentInfo);

  // 验证 TestToken
  console.log("\n验证 TestToken...");
  try {
    await run("verify:verify", {
      address: deploymentInfo.contracts.TestToken.address,
      constructorArguments: deploymentInfo.contracts.TestToken.constructorArgs,
    });
    console.log("TestToken 验证成功");
  } catch (error) {
    if (error.message.includes("Already Verified")) {
      console.log("TestToken 已经验证过了");
    } else {
      console.error("TestToken 验证失败:", error.message);
    }
  }

  // 验证 SimpleDEX
  console.log("\n验证 SimpleDEX...");
  try {
    await run("verify:verify", {
      address: deploymentInfo.contracts.SimpleDEX.address,
      constructorArguments: deploymentInfo.contracts.SimpleDEX.constructorArgs,
    });
    console.log("SimpleDEX 验证成功");
  } catch (error) {
    if (error.message.includes("Already Verified")) {
      console.log("SimpleDEX 已经验证过了");
    } else {
      console.error("SimpleDEX 验证失败:", error.message);
    }
  }

  console.log("\n验证完成!");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
```

### 11.8.4 升级和迁移脚本

```javascript
// scripts/upgrade.js - 合约升级脚本
const { ethers, upgrades } = require("hardhat");

async function main() {
  console.log("开始升级合约...");

  const [deployer] = await ethers.getSigners();
  console.log("部署账户:", deployer.address);

  // 现有代理合约地址（如果使用代理模式）
  const PROXY_ADDRESS = process.env.PROXY_ADDRESS;

  if (!PROXY_ADDRESS) {
    throw new Error("请设置 PROXY_ADDRESS 环境变量");
  }

  // 部署新的实现合约
  console.log("部署新的实现合约...");
  const TestTokenV2 = await ethers.getContractFactory("TestTokenV2");

  console.log("升级代理合约...");
  const upgraded = await upgrades.upgradeProxy(PROXY_ADDRESS, TestTokenV2);

  console.log("合约升级成功！");
  console.log("代理地址:", upgraded.target);

  // 验证升级
  const version = await upgraded.version();
  console.log("新版本:", version);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
```

## 11.9 部署后监控和运维

### 11.9.1 事件监听脚本

```javascript
// scripts/monitor.js - 合约事件监控
const { ethers } = require("hardhat");
const fs = require("fs");

async function main() {
  console.log("开始监控合约事件...");

  // 读取部署信息
  const deploymentInfo = JSON.parse(
    fs.readFileSync("./deployments/mainnet.json", "utf8")
  );

  // 连接到已部署的合约
  const token = await ethers.getContractAt(
    "TestToken",
    deploymentInfo.contracts.TestToken.address
  );

  const dex = await ethers.getContractAt(
    "SimpleDEX",
    deploymentInfo.contracts.SimpleDEX.address
  );

  console.log("监控 Token 合约:", token.target);
  console.log("监控 DEX 合约:", dex.target);

  // 监听 Token 事件
  token.on("Transfer", (from, to, value, event) => {
    console.log("\n=== Token Transfer ===");
    console.log("From:", from);
    console.log("To:", to);
    console.log("Value:", ethers.formatEther(value));
    console.log("Block:", event.blockNumber);
    console.log("Transaction:", event.transactionHash);

    // 记录到文件或数据库
    logEvent("Transfer", {
      contract: token.target,
      from,
      to,
      value: value.toString(),
      blockNumber: event.blockNumber,
      transactionHash: event.transactionHash,
      timestamp: new Date().toISOString(),
    });
  });

  token.on("MinterAdded", (minter, event) => {
    console.log("\n=== Minter Added ===");
    console.log("Minter:", minter);
    console.log("Block:", event.blockNumber);

    logEvent("MinterAdded", {
      contract: token.target,
      minter,
      blockNumber: event.blockNumber,
      transactionHash: event.transactionHash,
      timestamp: new Date().toISOString(),
    });
  });

  // 监听 DEX 事件
  dex.on("PoolCreated", (tokenA, tokenB, poolId, event) => {
    console.log("\n=== Pool Created ===");
    console.log("Token A:", tokenA);
    console.log("Token B:", tokenB);
    console.log("Pool ID:", poolId);
    console.log("Block:", event.blockNumber);

    logEvent("PoolCreated", {
      contract: dex.target,
      tokenA,
      tokenB,
      poolId,
      blockNumber: event.blockNumber,
      transactionHash: event.transactionHash,
      timestamp: new Date().toISOString(),
    });
  });

  dex.on(
    "TokenSwapped",
    (poolId, user, tokenIn, tokenOut, amountIn, amountOut, event) => {
      console.log("\n=== Token Swapped ===");
      console.log("Pool ID:", poolId);
      console.log("User:", user);
      console.log("Token In:", tokenIn);
      console.log("Token Out:", tokenOut);
      console.log("Amount In:", ethers.formatEther(amountIn));
      console.log("Amount Out:", ethers.formatEther(amountOut));
      console.log("Block:", event.blockNumber);

      logEvent("TokenSwapped", {
        contract: dex.target,
        poolId,
        user,
        tokenIn,
        tokenOut,
        amountIn: amountIn.toString(),
        amountOut: amountOut.toString(),
        blockNumber: event.blockNumber,
        transactionHash: event.transactionHash,
        timestamp: new Date().toISOString(),
      });
    }
  );

  console.log("事件监控已启动，按 Ctrl+C 停止...");

  // 保持进程运行
  process.stdin.resume();
}

function logEvent(eventName, data) {
  const logDir = "./logs";
  if (!fs.existsSync(logDir)) {
    fs.mkdirSync(logDir);
  }

  const logFile = `${logDir}/events-${
    new Date().toISOString().split("T")[0]
  }.log`;
  const logEntry = `${new Date().toISOString()} [${eventName}] ${JSON.stringify(
    data
  )}\n`;

  fs.appendFileSync(logFile, logEntry);
}

main().catch(console.error);
```

### 11.9.2 健康检查脚本

```javascript
// scripts/health-check.js - 合约健康检查
const { ethers } = require("hardhat");
const fs = require("fs");

async function main() {
  console.log("开始健康检查...");

  const deploymentInfo = JSON.parse(
    fs.readFileSync("./deployments/mainnet.json", "utf8")
  );

  const token = await ethers.getContractAt(
    "TestToken",
    deploymentInfo.contracts.TestToken.address
  );

  const dex = await ethers.getContractAt(
    "SimpleDEX",
    deploymentInfo.contracts.SimpleDEX.address
  );

  const results = {
    timestamp: new Date().toISOString(),
    token: await checkTokenHealth(token),
    dex: await checkDEXHealth(dex),
    network: await checkNetworkHealth(),
  };

  console.log("健康检查结果:", JSON.stringify(results, null, 2));

  // 保存结果
  fs.writeFileSync(
    `./logs/health-check-${Date.now()}.json`,
    JSON.stringify(results, null, 2)
  );

  // 检查是否有警告
  const warnings = [];
  if (results.token.totalSupply > results.token.maxSupply * 0.9) {
    warnings.push("Token 供应量接近最大值");
  }

  if (results.network.gasPrice > ethers.parseUnits("50", "gwei")) {
    warnings.push("网络 Gas 价格过高");
  }

  if (warnings.length > 0) {
    console.log("⚠️  发现警告:");
    warnings.forEach((warning) => console.log(`  - ${warning}`));

    // 这里可以发送警报通知
    // await sendAlert(warnings);
  } else {
    console.log("✅ 所有检查通过");
  }
}

async function checkTokenHealth(token) {
  try {
    const [
      name,
      symbol,
      totalSupply,
      maxSupply,
      mintPrice,
      mintingEnabled,
      owner,
    ] = await Promise.all([
      token.name(),
      token.symbol(),
      token.totalSupply(),
      token.MAX_SUPPLY(),
      token.mintPrice(),
      token.mintingEnabled(),
      token.owner(),
    ]);

    return {
      status: "healthy",
      name,
      symbol,
      totalSupply: totalSupply.toString(),
      maxSupply: maxSupply.toString(),
      mintPrice: mintPrice.toString(),
      mintingEnabled,
      owner,
      supplyUtilization:
        ((Number(totalSupply) / Number(maxSupply)) * 100).toFixed(2) + "%",
    };
  } catch (error) {
    return {
      status: "error",
      error: error.message,
    };
  }
}

async function checkDEXHealth(dex) {
  try {
    const owner = await dex.owner();
    const feeRate = await dex.feeRate();

    // 这里可以添加更多 DEX 特定的检查
    // 比如检查主要交易对的流动性等

    return {
      status: "healthy",
      owner,
      feeRate: feeRate.toString(),
    };
  } catch (error) {
    return {
      status: "error",
      error: error.message,
    };
  }
}

async function checkNetworkHealth() {
  try {
    const provider = ethers.provider;
    const [blockNumber, gasPrice, network] = await Promise.all([
      provider.getBlockNumber(),
      provider.getFeeData(),
      provider.getNetwork(),
    ]);

    return {
      status: "healthy",
      blockNumber,
      gasPrice: gasPrice.gasPrice?.toString() || "0",
      maxFeePerGas: gasPrice.maxFeePerGas?.toString() || "0",
      maxPriorityFeePerGas: gasPrice.maxPriorityFeePerGas?.toString() || "0",
      chainId: network.chainId.toString(),
    };
  } catch (error) {
    return {
      status: "error",
      error: error.message,
    };
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
```

## 11.10 章节总结

在本章的下篇中，我们学习了：

1. **集成测试**：多合约交互、完整业务流程测试
2. **Fork 测试**：使用主网分叉进行真实环境测试
3. **部署策略**：多网络部署、脚本管理、参数配置
4. **合约验证**：源码验证、自动化验证流程
5. **监控运维**：事件监听、健康检查、警报系统

### 部署最佳实践

1. **渐进部署**：测试网 → 主网，小规模 → 大规模
2. **验证充分**：所有合约都应进行源码验证
3. **监控到位**：部署后持续监控合约状态
4. **备份计划**：准备应急响应和升级方案
5. **文档完整**：维护详细的部署和运维文档

### 测试策略总结

- **单元测试**：快速反馈，高覆盖率
- **集成测试**：验证合约交互
- **Fork 测试**：真实环境验证
- **性能测试**：Gas 优化验证
- **安全测试**：漏洞检测

智能合约的测试和部署是一个系统性工程，需要在开发的各个阶段都保持严格的质量控制。

继续学习：[第十二章：高级特性与模式（上篇）- 代理模式与可升级合约](12_advanced_features_patterns_part1.md)
