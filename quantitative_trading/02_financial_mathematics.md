# 第二章：金融数学与统计基础

## 📖 章节概述

本章将为您提供量化交易所需的数学和统计学基础知识。这些知识是理解和开发量化策略的基石，涵盖了概率论、统计学、时间序列分析等核心内容。

## 🎯 学习目标

- 掌握金融中常用的概率论和统计学概念
- 理解时间序列数据的特性和分析方法
- 学会计算和分析相关性、协方差等统计指标
- 了解金融风险的数学度量方法
- 建立数学模型思维，为后续策略开发打基础

## 2.1 概率论基础

### 2.1.1 基本概念

**概率（Probability）**
概率是衡量事件发生可能性大小的数值，取值范围为[0,1]。

- P(A) = 0：事件 A 不可能发生
- P(A) = 1：事件 A 必然发生
- 0 < P(A) < 1：事件 A 有一定可能发生

**条件概率（Conditional Probability）**
在事件 B 发生的条件下，事件 A 发生的概率：

```
P(A|B) = P(A ∩ B) / P(B)
```

**金融应用示例：**

- P(股价上涨|利好消息) = 在利好消息发布条件下股价上涨的概率
- P(违约|信用评级下调) = 在信用评级下调条件下发生违约的概率

### 2.1.2 概率分布

**正态分布（Normal Distribution）**
最重要的概率分布，在金融市场中广泛应用：

- **概率密度函数**：f(x) = (1/√(2πσ²)) × e^(-(x-μ)²/(2σ²))
- **参数**：μ（均值）、σ（标准差）
- **记号**：X ~ N(μ, σ²)

**金融应用：**

- 股票收益率通常假设服从正态分布
- 许多风险模型基于正态分布假设
- Black-Scholes 期权定价模型的基础

**对数正态分布（Log-Normal Distribution）**
如果 ln(X)服从正态分布，则 X 服从对数正态分布：

- **应用**：股票价格建模（价格不能为负）
- **特点**：右偏分布，有长尾特征

**t 分布（Student's t-Distribution）**
当样本量较小时，更适合描述金融收益：

- **特点**：比正态分布有更厚的尾部
- **应用**：小样本统计推断、风险模型

### 2.1.3 贝叶斯定理

**公式：**

```
P(A|B) = P(B|A) × P(A) / P(B)
```

**金融应用示例：**

```python
# 贝叶斯更新投资观点
# P(牛市|经济指标好) = P(经济指标好|牛市) × P(牛市) / P(经济指标好)

# 先验概率
prior_bull = 0.3  # 牛市的先验概率

# 似然函数
likelihood = 0.8  # 牛市情况下经济指标好的概率

# 边际概率
marginal = 0.5  # 经济指标好的总概率

# 后验概率
posterior_bull = (likelihood * prior_bull) / marginal
print(f"观察到经济指标好后，牛市的概率为: {posterior_bull:.2f}")
```

## 2.2 描述性统计

### 2.2.1 中心趋势度量

**算术平均数（Arithmetic Mean）**

```
μ = (x₁ + x₂ + ... + xₙ) / n = Σxᵢ / n
```

**几何平均数（Geometric Mean）**

```
几何平均 = ⁿ√(x₁ × x₂ × ... × xₙ) = (∏xᵢ)^(1/n)
```

**金融应用：**

- 算术平均：计算平均收益率
- 几何平均：计算复合年化收益率

```python
import numpy as np

# 示例：计算投资收益率
returns = np.array([0.10, -0.05, 0.15, 0.08, -0.02])

# 算术平均收益率
arithmetic_mean = np.mean(returns)
print(f"算术平均收益率: {arithmetic_mean:.4f}")

# 几何平均收益率（年化）
geometric_mean = np.prod(1 + returns) ** (1/len(returns)) - 1
print(f"几何平均收益率: {geometric_mean:.4f}")
```

**中位数（Median）**

- 将数据按大小排序后位于中间位置的值
- 对异常值不敏感
- 适用于偏态分布

### 2.2.2 离散程度度量

**方差（Variance）**

```
σ² = Σ(xᵢ - μ)² / n
```

**标准差（Standard Deviation）**

```
σ = √(σ²)
```

**变异系数（Coefficient of Variation）**

```
CV = σ / μ
```

**金融应用：**

- 标准差：衡量投资风险
- 变异系数：比较不同投资的风险收益比

```python
# 风险度量示例
stock_returns = np.array([0.12, -0.08, 0.15, 0.05, -0.03, 0.09, -0.02, 0.11])

# 计算风险指标
mean_return = np.mean(stock_returns)
volatility = np.std(stock_returns, ddof=1)  # 样本标准差
cv = volatility / mean_return

print(f"平均收益率: {mean_return:.4f}")
print(f"波动率: {volatility:.4f}")
print(f"变异系数: {cv:.4f}")
```

### 2.2.3 分布形状度量

**偏度（Skewness）**
衡量分布的对称性：

```
偏度 = E[(X-μ)³] / σ³
```

- 偏度 = 0：对称分布
- 偏度 > 0：右偏（正偏）
- 偏度 < 0：左偏（负偏）

**峰度（Kurtosis）**
衡量分布的尖锐程度：

```
峰度 = E[(X-μ)⁴] / σ⁴
```

- 峰度 = 3：正态分布
- 峰度 > 3：尖峰分布（厚尾）
- 峰度 < 3：扁平分布（薄尾）

```python
from scipy import stats

# 计算偏度和峰度
returns = np.random.normal(0, 0.02, 1000)  # 模拟收益率数据

skewness = stats.skew(returns)
kurt = stats.kurtosis(returns)

print(f"偏度: {skewness:.4f}")
print(f"峰度: {kurt:.4f}")
```

## 2.3 相关性与协方差分析

### 2.3.1 协方差（Covariance）

**定义：**
协方差衡量两个随机变量的线性关系强度：

```
Cov(X,Y) = E[(X-μₓ)(Y-μᵧ)]
```

**样本协方差：**

```
Cov(X,Y) = Σ(xᵢ-x̄)(yᵢ-ȳ) / (n-1)
```

**性质：**

- Cov(X,Y) > 0：正相关
- Cov(X,Y) < 0：负相关
- Cov(X,Y) = 0：不相关

### 2.3.2 相关系数（Correlation Coefficient）

**皮尔逊相关系数：**

```
ρ(X,Y) = Cov(X,Y) / (σₓ × σᵧ)
```

**特点：**

- 取值范围：[-1, 1]
- |ρ| = 1：完全线性相关
- ρ = 0：线性无关
- 标准化的协方差，便于比较

```python
import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns

# 示例：计算股票相关性
# 假设我们有几只股票的收益率数据
np.random.seed(42)
dates = pd.date_range('2023-01-01', periods=100)

# 生成模拟数据
stock_data = pd.DataFrame({
    'AAPL': np.random.normal(0.001, 0.02, 100),
    'GOOGL': np.random.normal(0.0008, 0.018, 100),
    'MSFT': np.random.normal(0.0012, 0.019, 100),
    'TSLA': np.random.normal(0.002, 0.035, 100)
}, index=dates)

# 计算相关系数矩阵
correlation_matrix = stock_data.corr()
print("相关系数矩阵:")
print(correlation_matrix)

# 可视化相关性热力图
plt.figure(figsize=(8, 6))
sns.heatmap(correlation_matrix, annot=True, cmap='coolwarm', center=0)
plt.title('股票收益率相关性热力图')
plt.show()
```

### 2.3.3 相关性的金融意义

**投资组合理论中的应用：**

1. **分散化效果**：

   - 低相关性资产组合可以降低整体风险
   - 负相关资产具有对冲效果

2. **风险度量**：

   - 投资组合方差 = Σᵢ Σⱼ wᵢwⱼσᵢⱼ
   - 其中 σᵢⱼ = ρᵢⱼσᵢσⱼ

3. **资产配置**：
   - 相关性矩阵是最优化配置的关键输入
   - 动态调整基于相关性变化

```python
# 投资组合风险计算示例
def portfolio_risk(weights, cov_matrix):
    """计算投资组合风险"""
    portfolio_variance = np.dot(weights.T, np.dot(cov_matrix, weights))
    portfolio_volatility = np.sqrt(portfolio_variance)
    return portfolio_volatility

# 计算协方差矩阵
cov_matrix = stock_data.cov() * 252  # 年化协方差矩阵

# 等权重投资组合
equal_weights = np.array([0.25, 0.25, 0.25, 0.25])

# 计算投资组合风险
portfolio_vol = portfolio_risk(equal_weights, cov_matrix)
print(f"等权重投资组合年化波动率: {portfolio_vol:.4f}")
```

## 2.4 时间序列分析基础

### 2.4.1 时间序列的基本概念

**时间序列（Time Series）**
按时间顺序排列的数据序列：X(t₁), X(t₂), ..., X(tₙ)

**平稳性（Stationarity）**

- **强平稳**：联合分布不随时间变化
- **弱平稳**：均值和方差不随时间变化，协方差只依赖于时间间隔

**平稳性检验：**

```python
from statsmodels.tsa.stattools import adfuller

def check_stationarity(timeseries, title):
    """ADF检验判断序列平稳性"""
    result = adfuller(timeseries, autolag='AIC')

    print(f'{title}:')
    print(f'ADF统计量: {result[0]:.6f}')
    print(f'p值: {result[1]:.6f}')
    print(f'临界值:')
    for key, value in result[4].items():
        print(f'\t{key}: {value:.3f}')

    if result[1] <= 0.05:
        print("序列是平稳的")
    else:
        print("序列是非平稳的")
    print("-" * 50)

# 示例：检验股价序列的平稳性
np.random.seed(42)
dates = pd.date_range('2023-01-01', periods=252)
stock_price = 100 * np.cumprod(1 + np.random.normal(0.0008, 0.02, 252))
stock_returns = np.diff(np.log(stock_price))

check_stationarity(stock_price, "股票价格序列")
check_stationarity(stock_returns, "股票收益率序列")
```

### 2.4.2 自相关和偏自相关

**自相关函数（ACF）**
序列与其滞后值的相关系数：

```
ρₖ = Corr(Xₜ, Xₜ₋ₖ)
```

**偏自相关函数（PACF）**
控制中间变量后的相关系数

```python
from statsmodels.tsa.stattools import acf, pacf
from statsmodels.graphics.tsaplots import plot_acf, plot_pacf

# 计算并绘制ACF和PACF
fig, axes = plt.subplots(2, 1, figsize=(12, 8))

# ACF图
plot_acf(stock_returns, ax=axes[0], lags=20, title='自相关函数(ACF)')

# PACF图
plot_pacf(stock_returns, ax=axes[1], lags=20, title='偏自相关函数(PACF)')

plt.tight_layout()
plt.show()
```

### 2.4.3 移动平均和指数平滑

**简单移动平均（SMA）**

```
SMAₜ(n) = (Xₜ + Xₜ₋₁ + ... + Xₜ₋ₙ₊₁) / n
```

**指数移动平均（EMA）**

```
EMAₜ = α × Xₜ + (1-α) × EMAₜ₋₁
```

其中 α 是平滑参数（0 < α < 1）

```python
def calculate_moving_averages(prices, short_window=20, long_window=50):
    """计算移动平均线"""
    data = pd.DataFrame({'Price': prices})

    # 简单移动平均
    data[f'SMA_{short_window}'] = data['Price'].rolling(window=short_window).mean()
    data[f'SMA_{long_window}'] = data['Price'].rolling(window=long_window).mean()

    # 指数移动平均
    data[f'EMA_{short_window}'] = data['Price'].ewm(span=short_window).mean()
    data[f'EMA_{long_window}'] = data['Price'].ewm(span=long_window).mean()

    return data

# 示例应用
ma_data = calculate_moving_averages(stock_price)

# 绘制图形
plt.figure(figsize=(12, 6))
plt.plot(ma_data.index, ma_data['Price'], label='股价', alpha=0.7)
plt.plot(ma_data.index, ma_data['SMA_20'], label='SMA(20)', alpha=0.8)
plt.plot(ma_data.index, ma_data['SMA_50'], label='SMA(50)', alpha=0.8)
plt.plot(ma_data.index, ma_data['EMA_20'], label='EMA(20)', alpha=0.8)
plt.plot(ma_data.index, ma_data['EMA_50'], label='EMA(50)', alpha=0.8)
plt.legend()
plt.title('股价与移动平均线')
plt.show()
```

## 2.5 风险度量方法

### 2.5.1 波动率模型

**历史波动率**
基于历史收益率计算：

```
σ = √(Σ(rᵢ - r̄)² / (n-1))
```

**GARCH 模型**
广义自回归条件异方差模型：

```
σₜ² = ω + Σαᵢεₜ₋ᵢ² + Σβⱼσₜ₋ⱼ²
```

```python
from arch import arch_model

# GARCH模型估计波动率
def estimate_garch_volatility(returns):
    """使用GARCH(1,1)模型估计波动率"""
    # 将收益率转换为百分比形式
    returns_pct = returns * 100

    # 建立GARCH(1,1)模型
    model = arch_model(returns_pct, vol='Garch', p=1, q=1)

    # 拟合模型
    fitted_model = model.fit(disp='off')

    # 提取条件波动率
    conditional_volatility = fitted_model.conditional_volatility / 100

    return conditional_volatility, fitted_model

# 应用示例
garch_vol, garch_model = estimate_garch_volatility(stock_returns)

print("GARCH模型参数:")
print(garch_model.summary())

# 绘制波动率图
plt.figure(figsize=(12, 6))
plt.plot(garch_vol.index, garch_vol, label='GARCH条件波动率')
plt.plot(garch_vol.index, pd.Series(stock_returns).rolling(30).std(),
         label='30日滚动波动率', alpha=0.7)
plt.legend()
plt.title('波动率比较')
plt.show()
```

### 2.5.2 VaR（风险价值）

**历史模拟法**

```python
def calculate_var_historical(returns, confidence_level=0.05):
    """历史模拟法计算VaR"""
    return np.percentile(returns, confidence_level * 100)

# 计算不同置信水平的VaR
returns_daily = stock_returns
var_95 = calculate_var_historical(returns_daily, 0.05)
var_99 = calculate_var_historical(returns_daily, 0.01)

print(f"95% VaR: {var_95:.4f}")
print(f"99% VaR: {var_99:.4f}")
```

**参数法（正态分布假设）**

```python
def calculate_var_parametric(returns, confidence_level=0.05):
    """参数法计算VaR"""
    mu = np.mean(returns)
    sigma = np.std(returns)
    z_score = stats.norm.ppf(confidence_level)
    var = mu + z_score * sigma
    return var

var_95_param = calculate_var_parametric(returns_daily, 0.05)
var_99_param = calculate_var_parametric(returns_daily, 0.01)

print(f"参数法95% VaR: {var_95_param:.4f}")
print(f"参数法99% VaR: {var_99_param:.4f}")
```

### 2.5.3 期望亏空（ES/CVaR）

```python
def calculate_expected_shortfall(returns, confidence_level=0.05):
    """计算期望亏空"""
    var = calculate_var_historical(returns, confidence_level)
    es = np.mean(returns[returns <= var])
    return es

es_95 = calculate_expected_shortfall(returns_daily, 0.05)
es_99 = calculate_expected_shortfall(returns_daily, 0.01)

print(f"95% Expected Shortfall: {es_95:.4f}")
print(f"99% Expected Shortfall: {es_99:.4f}")
```

## 2.6 回归分析

### 2.6.1 线性回归

**简单线性回归：**

```
Y = α + βX + ε
```

**多元线性回归：**

```
Y = α + β₁X₁ + β₂X₂ + ... + βₖXₖ + ε
```

```python
from sklearn.linear_model import LinearRegression
from sklearn.metrics import r2_score

# 示例：市场模型回归
def market_model_regression(stock_returns, market_returns):
    """市场模型回归分析"""
    # 创建回归模型
    model = LinearRegression()

    # 拟合模型
    X = market_returns.reshape(-1, 1)
    y = stock_returns
    model.fit(X, y)

    # 预测值
    y_pred = model.predict(X)

    # 计算统计指标
    alpha = model.intercept_
    beta = model.coef_[0]
    r_squared = r2_score(y, y_pred)

    return alpha, beta, r_squared, y_pred

# 生成市场收益率数据
market_returns = np.random.normal(0.001, 0.015, len(stock_returns))

# 执行回归
alpha, beta, r2, predictions = market_model_regression(stock_returns, market_returns)

print(f"Alpha (截距): {alpha:.6f}")
print(f"Beta (市场敏感度): {beta:.4f}")
print(f"R²: {r2:.4f}")

# 绘制回归结果
plt.figure(figsize=(10, 6))
plt.scatter(market_returns, stock_returns, alpha=0.6, label='实际数据')
plt.plot(market_returns, predictions, color='red', label='回归线')
plt.xlabel('市场收益率')
plt.ylabel('股票收益率')
plt.title(f'市场模型回归 (Beta={beta:.4f}, R²={r2:.4f})')
plt.legend()
plt.show()
```

### 2.6.2 多因子模型

```python
def multi_factor_model(returns, factors):
    """多因子模型回归"""
    from sklearn.linear_model import LinearRegression

    model = LinearRegression()
    model.fit(factors, returns)

    # 因子载荷（回归系数）
    factor_loadings = model.coef_
    alpha = model.intercept_
    r_squared = model.score(factors, returns)

    return alpha, factor_loadings, r_squared

# 示例：三因子模型
# 假设有市场、规模、价值三个因子
np.random.seed(42)
n_obs = len(stock_returns)

factor_data = pd.DataFrame({
    'Market': np.random.normal(0.001, 0.015, n_obs),
    'Size': np.random.normal(0.0005, 0.01, n_obs),
    'Value': np.random.normal(0.0003, 0.008, n_obs)
})

# 执行多因子回归
alpha_mf, loadings, r2_mf = multi_factor_model(stock_returns, factor_data.values)

print("多因子模型结果:")
print(f"Alpha: {alpha_mf:.6f}")
for i, factor in enumerate(factor_data.columns):
    print(f"{factor}因子载荷: {loadings[i]:.4f}")
print(f"R²: {r2_mf:.4f}")
```

## 2.7 假设检验

### 2.7.1 正态性检验

```python
from scipy.stats import shapiro, jarque_bera, normaltest

def test_normality(data, alpha=0.05):
    """正态性检验"""
    print("正态性检验结果:")

    # Shapiro-Wilk检验
    stat_sw, p_sw = shapiro(data)
    print(f"Shapiro-Wilk: 统计量={stat_sw:.4f}, p值={p_sw:.4f}")

    # Jarque-Bera检验
    stat_jb, p_jb = jarque_bera(data)
    print(f"Jarque-Bera: 统计量={stat_jb:.4f}, p值={p_jb:.4f}")

    # D'Agostino's检验
    stat_da, p_da = normaltest(data)
    print(f"D'Agostino: 统计量={stat_da:.4f}, p值={p_da:.4f}")

    # 结论
    tests = [p_sw, p_jb, p_da]
    if all(p > alpha for p in tests):
        print(f"结论: 在{alpha}显著性水平下，不能拒绝正态分布假设")
    else:
        print(f"结论: 在{alpha}显著性水平下，拒绝正态分布假设")

# 检验收益率的正态性
test_normality(stock_returns)
```

### 2.7.2 统计显著性检验

```python
from scipy.stats import ttest_1samp, ttest_ind

def test_mean_return(returns, expected_return=0, alpha=0.05):
    """检验平均收益率是否显著不为零"""
    t_stat, p_value = ttest_1samp(returns, expected_return)

    print(f"单样本t检验:")
    print(f"零假设: 平均收益率 = {expected_return}")
    print(f"t统计量: {t_stat:.4f}")
    print(f"p值: {p_value:.4f}")

    if p_value < alpha:
        print(f"在{alpha}显著性水平下，拒绝零假设")
    else:
        print(f"在{alpha}显著性水平下，不能拒绝零假设")

# 检验平均收益率是否显著不为零
test_mean_return(stock_returns)
```

## 📚 本章小结

本章为您介绍了量化交易所需的数学和统计学基础，包括：

1. **概率论基础**：掌握了概率、条件概率、贝叶斯定理等核心概念
2. **描述性统计**：学会了计算和解释各种统计指标
3. **相关性分析**：理解了协方差和相关系数在投资中的应用
4. **时间序列**：了解了平稳性、自相关等时间序列特性
5. **风险度量**：掌握了波动率、VaR、ES 等风险度量方法
6. **回归分析**：学会了线性回归和多因子模型的应用
7. **假设检验**：了解了统计推断的基本方法

这些数学和统计学工具将在后续的策略开发和风险管理中发挥重要作用。

## 🔗 下一章预告

在[第三章：Python 金融编程（上篇）](03_python_finance_programming.md)中，我们将学习如何使用 Python 进行金融数据分析，包括数据获取、处理、可视化等实践技能。

## 💡 实践练习

1. **统计分析**：选择一只股票，计算其收益率的基本统计指标（均值、标准差、偏度、峰度）

2. **相关性分析**：选择 3-5 只股票，计算它们之间的相关系数矩阵，分析投资组合的分散化效果

3. **波动率建模**：使用 GARCH 模型拟合某股票的波动率，比较与历史波动率的差异

4. **风险度量**：计算投资组合的 VaR 和 ES，分析不同置信水平下的风险

5. **回归分析**：构建市场模型，分析某股票的 beta 系数和超额收益

---

**建议学习时间：** 3-4 天  
**前置章节：** [第一章：量化交易入门](01_quantitative_trading_basics.md)  
**下一章：** [Python 金融编程（上篇）](03_python_finance_programming.md)
