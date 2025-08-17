# 第7章 多因子模型与组合优化（上篇）

## 📚 本章概览

多因子模型是现代量化投资的核心理论基础，它通过识别和利用影响资产价格的多个因子来构建投资策略。本章将深入介绍因子投资理论、常见因子类型、因子有效性检验方法，以及多因子模型的构建技术。

### 🎯 学习目标

- 理解因子投资的理论基础和发展历程
- 掌握常见因子的分类和特征
- 学会因子有效性检验方法
- 掌握多因子模型的构建技术
- 了解因子选择和组合的实用方法

## 第一节：因子投资理论基础

### 1.1 因子投资的发展历程

#### 资本资产定价模型（CAPM）

```python
# CAPM模型实现
import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
from scipy import stats

def capm_analysis(stock_returns, market_returns, risk_free_rate=0.02):
    """
    CAPM模型分析
    
    Parameters:
    stock_returns: 股票收益率序列
    market_returns: 市场收益率序列
    risk_free_rate: 无风险利率
    """
    
    # 计算超额收益率
    excess_stock_returns = stock_returns - risk_free_rate/252
    excess_market_returns = market_returns - risk_free_rate/252
    
    # 回归分析
    slope, intercept, r_value, p_value, std_err = stats.linregress(
        excess_market_returns, excess_stock_returns
    )
    
    # Beta系数
    beta = slope
    alpha = intercept
    
    # 理论期望收益率
    expected_return = risk_free_rate + beta * (excess_market_returns.mean() * 252)
    
    results = {
        'beta': beta,
        'alpha': alpha,
        'r_squared': r_value**2,
        'p_value': p_value,
        'expected_return': expected_return,
        'tracking_error': std_err * np.sqrt(252)
    }
    
    return results

# 示例使用
# results = capm_analysis(stock_returns, market_returns)
# print(f"Beta: {results['beta']:.3f}")
# print(f"Alpha: {results['alpha']:.3f}")
# print(f"R²: {results['r_squared']:.3f}")
```

#### Fama-French三因子模型

```python
def fama_french_three_factor(returns, market_returns, smb, hml, risk_free_rate=0.02):
    """
    Fama-French三因子模型分析
    
    Parameters:
    returns: 股票收益率
    market_returns: 市场收益率
    smb: 小市值减大市值因子
    hml: 高账面市值比减低账面市值比因子
    risk_free_rate: 无风险利率
    """
    
    # 计算超额收益率
    excess_returns = returns - risk_free_rate/252
    excess_market = market_returns - risk_free_rate/252
    
    # 构建回归矩阵
    X = np.column_stack([excess_market, smb, hml])
    X = np.column_stack([np.ones(len(X)), X])  # 添加常数项
    
    # 多元线性回归
    coeffs = np.linalg.lstsq(X, excess_returns, rcond=None)[0]
    
    alpha = coeffs[0]
    beta_market = coeffs[1]
    beta_smb = coeffs[2]
    beta_hml = coeffs[3]
    
    # 拟合度计算
    y_pred = X @ coeffs
    r_squared = 1 - np.sum((excess_returns - y_pred)**2) / np.sum((excess_returns - excess_returns.mean())**2)
    
    return {
        'alpha': alpha,
        'beta_market': beta_market,
        'beta_smb': beta_smb,
        'beta_hml': beta_hml,
        'r_squared': r_squared
    }
```

### 1.2 因子投资的理论基础

#### 套利定价理论（APT）

套利定价理论认为资产收益率可以由多个因子的线性组合来解释：

$$R_i = \alpha_i + \beta_{i1}F_1 + \beta_{i2}F_2 + ... + \beta_{ik}F_k + \epsilon_i$$

其中：
- $R_i$ 是资产i的收益率
- $\alpha_i$ 是资产i的特异性收益
- $\beta_{ij}$ 是资产i对因子j的敞口
- $F_j$ 是因子j的收益率
- $\epsilon_i$ 是随机误差项

```python
def apt_model(returns, factors):
    """
    套利定价理论模型实现
    
    Parameters:
    returns: 资产收益率矩阵 (T x N)
    factors: 因子收益率矩阵 (T x K)
    """
    
    # 添加常数项
    X = np.column_stack([np.ones(len(factors)), factors])
    
    # 对每个资产进行回归
    betas = []
    alphas = []
    r_squareds = []
    
    for i in range(returns.shape[1]):
        y = returns.iloc[:, i]
        coeffs = np.linalg.lstsq(X, y, rcond=None)[0]
        
        alphas.append(coeffs[0])
        betas.append(coeffs[1:])
        
        # 计算R²
        y_pred = X @ coeffs
        r_squared = 1 - np.sum((y - y_pred)**2) / np.sum((y - y.mean())**2)
        r_squareds.append(r_squared)
    
    return {
        'alphas': np.array(alphas),
        'betas': np.array(betas),
        'r_squareds': np.array(r_squareds)
    }
```

## 第二节：常见因子介绍

### 2.1 基本面因子

#### 价值因子（Value Factors）

```python
def calculate_value_factors(price_data, fundamental_data):
    """
    计算价值因子
    
    Parameters:
    price_data: 价格数据
    fundamental_data: 基本面数据（包含总市值、净资产、营业收入、净利润等）
    """
    
    factors = pd.DataFrame(index=price_data.index)
    
    # 市净率倒数（BP）
    factors['BP'] = fundamental_data['book_value'] / fundamental_data['market_cap']
    
    # 市盈率倒数（EP）
    factors['EP'] = fundamental_data['earnings'] / fundamental_data['market_cap']
    
    # 市销率倒数（SP）
    factors['SP'] = fundamental_data['sales'] / fundamental_data['market_cap']
    
    # 市现率倒数（CFP）
    factors['CFP'] = fundamental_data['cash_flow'] / fundamental_data['market_cap']
    
    # EBITDA/EV
    ev = fundamental_data['market_cap'] + fundamental_data['debt'] - fundamental_data['cash']
    factors['EBITDA_EV'] = fundamental_data['ebitda'] / ev
    
    return factors

def value_factor_score(factors):
    """
    计算综合价值因子得分
    """
    # 标准化处理
    standardized = factors.apply(lambda x: (x - x.mean()) / x.std())
    
    # 等权重合成
    value_score = standardized.mean(axis=1)
    
    return value_score
```

#### 质量因子（Quality Factors）

```python
def calculate_quality_factors(fundamental_data):
    """
    计算质量因子
    """
    
    factors = pd.DataFrame(index=fundamental_data.index)
    
    # ROE（净资产收益率）
    factors['ROE'] = fundamental_data['net_income'] / fundamental_data['book_value']
    
    # ROA（总资产收益率）
    factors['ROA'] = fundamental_data['net_income'] / fundamental_data['total_assets']
    
    # ROIC（投入资本回报率）
    invested_capital = fundamental_data['total_assets'] - fundamental_data['cash'] - fundamental_data['current_liabilities']
    factors['ROIC'] = fundamental_data['operating_income'] / invested_capital
    
    # 毛利率
    factors['Gross_Margin'] = fundamental_data['gross_profit'] / fundamental_data['revenue']
    
    # 净利率
    factors['Net_Margin'] = fundamental_data['net_income'] / fundamental_data['revenue']
    
    # 资产负债率
    factors['Debt_Ratio'] = fundamental_data['total_debt'] / fundamental_data['total_assets']
    
    # 流动比率
    factors['Current_Ratio'] = fundamental_data['current_assets'] / fundamental_data['current_liabilities']
    
    # 利息覆盖倍数
    factors['Interest_Coverage'] = fundamental_data['ebit'] / fundamental_data['interest_expense']
    
    return factors
```

#### 成长因子（Growth Factors）

```python
def calculate_growth_factors(fundamental_data, periods=[1, 2, 3]):
    """
    计算成长因子
    
    Parameters:
    fundamental_data: 基本面数据
    periods: 计算增长率的期间（年）
    """
    
    factors = pd.DataFrame(index=fundamental_data.index)
    
    for period in periods:
        # 营业收入增长率
        factors[f'Sales_Growth_{period}Y'] = (
            fundamental_data['revenue'] / fundamental_data['revenue'].shift(period*4) - 1
        )
        
        # 净利润增长率
        factors[f'Earnings_Growth_{period}Y'] = (
            fundamental_data['net_income'] / fundamental_data['net_income'].shift(period*4) - 1
        )
        
        # 总资产增长率
        factors[f'Asset_Growth_{period}Y'] = (
            fundamental_data['total_assets'] / fundamental_data['total_assets'].shift(period*4) - 1
        )
        
        # 净资产增长率
        factors[f'Equity_Growth_{period}Y'] = (
            fundamental_data['book_value'] / fundamental_data['book_value'].shift(period*4) - 1
        )
    
    # 计算增长率稳定性
    for metric in ['Sales_Growth', 'Earnings_Growth']:
        growth_cols = [f'{metric}_{p}Y' for p in periods]
        factors[f'{metric}_Stability'] = -factors[growth_cols].std(axis=1)  # 负号表示稳定性越高越好
    
    return factors
```

### 2.2 技术因子

#### 动量因子（Momentum Factors）

```python
def calculate_momentum_factors(price_data):
    """
    计算动量因子
    """
    
    factors = pd.DataFrame(index=price_data.index)
    
    # 价格动量
    factors['Mom_1M'] = price_data.pct_change(20)  # 1个月动量
    factors['Mom_3M'] = price_data.pct_change(60)  # 3个月动量
    factors['Mom_6M'] = price_data.pct_change(120) # 6个月动量
    factors['Mom_12M'] = price_data.pct_change(240) # 12个月动量
    
    # 跳过最近一个月的动量（避免短期反转）
    factors['Mom_12_1'] = (price_data / price_data.shift(240) - 1) - (price_data / price_data.shift(20) - 1)
    
    # 风险调整动量
    returns = price_data.pct_change()
    rolling_vol = returns.rolling(60).std() * np.sqrt(252)
    factors['Risk_Adj_Mom'] = factors['Mom_3M'] / rolling_vol
    
    return factors

def calculate_reversal_factors(price_data, volume_data=None):
    """
    计算反转因子
    """
    
    factors = pd.DataFrame(index=price_data.index)
    returns = price_data.pct_change()
    
    # 短期反转
    factors['Reversal_1W'] = -returns.rolling(5).mean()
    factors['Reversal_1M'] = -returns.rolling(20).mean()
    
    # 长期反转
    factors['Reversal_36M'] = -price_data.pct_change(720)
    
    # 成交量加权反转
    if volume_data is not None:
        vwap_returns = (returns * volume_data).rolling(20).sum() / volume_data.rolling(20).sum()
        factors['Volume_Reversal'] = -vwap_returns
    
    return factors
```

#### 波动率因子

```python
def calculate_volatility_factors(price_data, volume_data=None):
    """
    计算波动率因子
    """
    
    factors = pd.DataFrame(index=price_data.index)
    returns = price_data.pct_change()
    
    # 历史波动率
    factors['Vol_20D'] = returns.rolling(20).std() * np.sqrt(252)
    factors['Vol_60D'] = returns.rolling(60).std() * np.sqrt(252)
    factors['Vol_252D'] = returns.rolling(252).std() * np.sqrt(252)
    
    # 波动率变化
    factors['Vol_Change'] = factors['Vol_20D'] / factors['Vol_60D'] - 1
    
    # 上下行波动率
    upside_returns = returns.where(returns > 0, 0)
    downside_returns = returns.where(returns < 0, 0)
    
    factors['Upside_Vol'] = upside_returns.rolling(60).std() * np.sqrt(252)
    factors['Downside_Vol'] = downside_returns.rolling(60).std() * np.sqrt(252)
    factors['Vol_Skew'] = factors['Upside_Vol'] / factors['Downside_Vol']
    
    # GARCH波动率
    factors['GARCH_Vol'] = calculate_garch_volatility(returns)
    
    return factors

def calculate_garch_volatility(returns, window=252):
    """
    计算GARCH模型波动率
    """
    from arch import arch_model
    
    garch_vol = pd.Series(index=returns.index, dtype=float)
    
    for i in range(window, len(returns)):
        try:
            data = returns.iloc[i-window:i] * 100  # 转换为百分比
            model = arch_model(data, vol='GARCH', p=1, q=1)
            res = model.fit(disp='off')
            garch_vol.iloc[i] = res.conditional_volatility.iloc[-1] / 100
        except:
            garch_vol.iloc[i] = returns.iloc[i-window:i].std()
    
    return garch_vol * np.sqrt(252)
```

### 2.3 市场结构因子

#### 流动性因子

```python
def calculate_liquidity_factors(price_data, volume_data, high_data, low_data):
    """
    计算流动性因子
    """
    
    factors = pd.DataFrame(index=price_data.index)
    returns = price_data.pct_change()
    
    # 换手率
    factors['Turnover'] = volume_data.rolling(20).mean()
    factors['Turnover_Volatility'] = volume_data.rolling(20).std()
    
    # Amihud非流动性指标
    factors['Amihud'] = np.abs(returns) / volume_data
    factors['Amihud_20D'] = factors['Amihud'].rolling(20).mean()
    
    # 价格冲击
    factors['Price_Impact'] = np.abs(returns) / np.log(volume_data + 1)
    
    # Bid-Ask价差代理
    factors['Bid_Ask_Spread'] = 2 * np.abs(returns) / (high_data + low_data)
    
    # 零收益率天数比例
    zero_return_days = (returns == 0).rolling(20).sum()
    factors['Zero_Return_Days'] = zero_return_days / 20
    
    # Roll价差估计
    factors['Roll_Spread'] = 2 * np.sqrt(-returns.rolling(2).cov().iloc[::2, 1::2].values.diagonal())
    
    return factors
```

#### 规模因子

```python
def calculate_size_factors(market_cap_data, price_data):
    """
    计算规模因子
    """
    
    factors = pd.DataFrame(index=market_cap_data.index)
    
    # 对数市值
    factors['Log_MarketCap'] = np.log(market_cap_data)
    
    # 自由流通市值
    factors['Log_FloatCap'] = np.log(market_cap_data * 0.7)  # 假设流通比例70%
    
    # 市值排名
    factors['MarketCap_Rank'] = market_cap_data.rank(pct=True)
    
    # 价格水平
    factors['Price_Level'] = np.log(price_data)
    
    return factors
```

## 第三节：因子有效性检验

### 3.1 单因子检验

```python
def single_factor_test(factor_data, return_data, n_quantiles=5):
    """
    单因子有效性检验
    
    Parameters:
    factor_data: 因子数据
    return_data: 收益率数据
    n_quantiles: 分位数数量
    """
    
    results = {}
    
    # 按因子值分组
    factor_quantiles = pd.qcut(factor_data, n_quantiles, labels=False)
    
    # 计算各分位数组合收益率
    quantile_returns = []
    for q in range(n_quantiles):
        mask = (factor_quantiles == q)
        if mask.sum() > 0:
            portfolio_return = return_data[mask].mean()
            quantile_returns.append(portfolio_return)
        else:
            quantile_returns.append(np.nan)
    
    quantile_returns = pd.Series(quantile_returns, index=range(n_quantiles))
    
    # 多空组合收益率
    long_short_return = quantile_returns.iloc[-1] - quantile_returns.iloc[0]
    
    # IC计算（信息系数）
    ic = factor_data.corr(return_data)
    
    # Rank IC计算
    rank_ic = factor_data.rank().corr(return_data.rank())
    
    # t统计量
    t_stat = ic * np.sqrt(len(factor_data) - 2) / np.sqrt(1 - ic**2)
    
    results = {
        'quantile_returns': quantile_returns,
        'long_short_return': long_short_return,
        'ic': ic,
        'rank_ic': rank_ic,
        't_statistic': t_stat,
        'factor_mean': factor_data.mean(),
        'factor_std': factor_data.std()
    }
    
    return results

def ic_analysis(factor_data, return_data, window=20):
    """
    IC分析（信息系数时间序列分析）
    """
    
    ic_series = []
    rank_ic_series = []
    dates = []
    
    for i in range(window, len(factor_data)):
        start_idx = i - window
        end_idx = i
        
        factor_slice = factor_data.iloc[start_idx:end_idx]
        return_slice = return_data.iloc[start_idx:end_idx]
        
        # 计算IC
        ic = factor_slice.corr(return_slice)
        rank_ic = factor_slice.rank().corr(return_slice.rank())
        
        ic_series.append(ic)
        rank_ic_series.append(rank_ic)
        dates.append(factor_data.index[end_idx-1])
    
    ic_df = pd.DataFrame({
        'IC': ic_series,
        'Rank_IC': rank_ic_series
    }, index=dates)
    
    # IC统计指标
    ic_stats = {
        'IC_mean': ic_df['IC'].mean(),
        'IC_std': ic_df['IC'].std(),
        'IC_IR': ic_df['IC'].mean() / ic_df['IC'].std(),  # IR = Information Ratio
        'IC_win_rate': (ic_df['IC'] > 0).mean(),
        'Rank_IC_mean': ic_df['Rank_IC'].mean(),
        'Rank_IC_std': ic_df['Rank_IC'].std(),
        'Rank_IC_IR': ic_df['Rank_IC'].mean() / ic_df['Rank_IC'].std(),
        'Rank_IC_win_rate': (ic_df['Rank_IC'] > 0).mean()
    }
    
    return ic_df, ic_stats
```

### 3.2 因子回归分析

```python
def factor_regression_analysis(factor_data, return_data, benchmark_return=None):
    """
    因子回归分析
    """
    
    from scipy import stats
    
    # 准备数据
    if benchmark_return is not None:
        excess_return = return_data - benchmark_return
    else:
        excess_return = return_data
    
    # 回归分析
    slope, intercept, r_value, p_value, std_err = stats.linregress(
        factor_data.dropna(), excess_return.dropna()
    )
    
    # 构建结果
    results = {
        'beta': slope,
        'alpha': intercept,
        'r_squared': r_value**2,
        'p_value': p_value,
        'std_error': std_err,
        't_statistic': slope / std_err if std_err != 0 else np.nan,
        'residuals': excess_return - (slope * factor_data + intercept)
    }
    
    return results

def factor_exposure_analysis(factor_data, return_data, periods=[20, 60, 120, 252]):
    """
    因子暴露度分析
    """
    
    exposures = {}
    
    for period in periods:
        rolling_corr = factor_data.rolling(period).corr(return_data.rolling(period).mean())
        exposures[f'{period}D'] = rolling_corr
    
    # 暴露度稳定性
    exposure_stability = pd.DataFrame()
    for period in periods:
        exposure_stability[f'{period}D_std'] = exposures[f'{period}D'].rolling(252).std()
    
    return exposures, exposure_stability
```

## 第四节：多因子模型构建

### 4.1 因子选择方法

```python
def factor_selection(factor_data, return_data, method='stepwise'):
    """
    因子选择
    
    Parameters:
    factor_data: 因子数据矩阵
    return_data: 收益率数据
    method: 选择方法 ('stepwise', 'lasso', 'elastic_net', 'pca')
    """
    
    from sklearn.feature_selection import f_regression, SelectKBest
    from sklearn.linear_model import LassoCV, ElasticNetCV
    from sklearn.decomposition import PCA
    
    # 去除缺失值
    combined_data = pd.concat([factor_data, return_data], axis=1).dropna()
    X = combined_data.iloc[:, :-1]
    y = combined_data.iloc[:, -1]
    
    if method == 'stepwise':
        # 逐步回归选择
        selected_factors = stepwise_selection(X, y)
        
    elif method == 'lasso':
        # Lasso回归选择
        lasso = LassoCV(cv=5, random_state=42)
        lasso.fit(X, y)
        
        # 选择非零系数的因子
        selected_mask = lasso.coef_ != 0
        selected_factors = X.columns[selected_mask].tolist()
        
    elif method == 'elastic_net':
        # 弹性网络选择
        elastic_net = ElasticNetCV(cv=5, random_state=42)
        elastic_net.fit(X, y)
        
        selected_mask = elastic_net.coef_ != 0
        selected_factors = X.columns[selected_mask].tolist()
        
    elif method == 'pca':
        # 主成分分析
        pca = PCA(n_components=0.95)  # 保留95%的方差
        X_pca = pca.fit_transform(X)
        
        # 返回主成分和权重
        components_df = pd.DataFrame(
            pca.components_.T,
            columns=[f'PC{i+1}' for i in range(pca.n_components_)],
            index=X.columns
        )
        
        return {
            'method': method,
            'components': components_df,
            'explained_variance_ratio': pca.explained_variance_ratio_,
            'transformed_data': X_pca
        }
    
    return {
        'method': method,
        'selected_factors': selected_factors,
        'n_factors': len(selected_factors)
    }

def stepwise_selection(X, y, threshold_in=0.01, threshold_out=0.05):
    """
    逐步回归因子选择
    """
    included = []
    
    while True:
        changed = False
        
        # Forward step
        excluded = list(set(X.columns) - set(included))
        new_pval = pd.Series(index=excluded, dtype=float)
        
        for new_column in excluded:
            if len(included) == 0:
                new_pval[new_column] = f_regression(X[[new_column]], y)[1][0]
            else:
                new_pval[new_column] = f_regression(X[included + [new_column]], y)[1][-1]
        
        best_pval = new_pval.min()
        if best_pval < threshold_in:
            best_feature = new_pval.idxmin()
            included.append(best_feature)
            changed = True
        
        # Backward step
        if len(included) > 1:
            pvals = f_regression(X[included], y)[1]
            worst_pval = pvals.max()
            
            if worst_pval > threshold_out:
                worst_feature = included[pvals.argmax()]
                included.remove(worst_feature)
                changed = True
        
        if not changed:
            break
    
    return included
```

### 4.2 因子正交化

```python
def factor_orthogonalization(factor_data, method='gram_schmidt'):
    """
    因子正交化
    
    Parameters:
    factor_data: 因子数据矩阵
    method: 正交化方法 ('gram_schmidt', 'symmetric', 'regression')
    """
    
    if method == 'gram_schmidt':
        # Gram-Schmidt正交化
        orthogonal_factors = gram_schmidt_orthogonalization(factor_data)
        
    elif method == 'symmetric':
        # 对称正交化
        correlation_matrix = factor_data.corr()
        eigenvals, eigenvecs = np.linalg.eigh(correlation_matrix)
        
        # 构建正交化矩阵
        sqrt_eigenvals = np.sqrt(np.maximum(eigenvals, 1e-8))
        orthogonal_matrix = eigenvecs @ np.diag(1/sqrt_eigenvals) @ eigenvecs.T
        
        # 标准化因子
        standardized_factors = (factor_data - factor_data.mean()) / factor_data.std()
        orthogonal_factors = standardized_factors @ orthogonal_matrix
        orthogonal_factors.columns = factor_data.columns
        
    elif method == 'regression':
        # 回归正交化
        orthogonal_factors = regression_orthogonalization(factor_data)
    
    return orthogonal_factors

def gram_schmidt_orthogonalization(factor_data):
    """
    Gram-Schmidt正交化过程
    """
    factors = factor_data.copy()
    orthogonal_factors = pd.DataFrame(index=factors.index)
    
    for i, column in enumerate(factors.columns):
        if i == 0:
            # 第一个因子直接标准化
            orthogonal_factors[column] = (factors[column] - factors[column].mean()) / factors[column].std()
        else:
            # 后续因子需要去除与前面因子的相关性
            current_factor = factors[column]
            
            # 投影到已有正交因子空间的分量
            for prev_column in orthogonal_factors.columns:
                proj_coef = current_factor.cov(orthogonal_factors[prev_column]) / orthogonal_factors[prev_column].var()
                current_factor = current_factor - proj_coef * orthogonal_factors[prev_column]
            
            # 标准化
            orthogonal_factors[column] = (current_factor - current_factor.mean()) / current_factor.std()
    
    return orthogonal_factors

def regression_orthogonalization(factor_data, order=None):
    """
    回归正交化
    """
    if order is None:
        order = factor_data.columns.tolist()
    
    orthogonal_factors = pd.DataFrame(index=factor_data.index)
    
    for i, factor_name in enumerate(order):
        if i == 0:
            # 第一个因子保持不变
            orthogonal_factors[factor_name] = factor_data[factor_name]
        else:
            # 对已有因子回归，取残差
            y = factor_data[factor_name]
            X = orthogonal_factors.iloc[:, :i]
            
            # 线性回归
            coeffs = np.linalg.lstsq(X.values, y.values, rcond=None)[0]
            predicted = X.values @ coeffs
            residual = y - predicted
            
            orthogonal_factors[factor_name] = residual
    
    # 标准化
    orthogonal_factors = (orthogonal_factors - orthogonal_factors.mean()) / orthogonal_factors.std()
    
    return orthogonal_factors
```

### 4.3 因子权重确定

```python
def factor_weighting(factor_returns, return_data, method='equal_weight'):
    """
    因子权重确定
    
    Parameters:
    factor_returns: 因子收益率矩阵
    return_data: 目标收益率
    method: 权重确定方法
    """
    
    if method == 'equal_weight':
        # 等权重
        weights = pd.Series(1/len(factor_returns.columns), index=factor_returns.columns)
        
    elif method == 'ic_weight':
        # IC加权
        ics = []
        for factor in factor_returns.columns:
            ic = factor_returns[factor].corr(return_data)
            ics.append(abs(ic))
        
        weights = pd.Series(ics, index=factor_returns.columns)
        weights = weights / weights.sum()
        
    elif method == 'regression_weight':
        # 回归权重
        from sklearn.linear_model import LinearRegression
        
        model = LinearRegression()
        model.fit(factor_returns, return_data)
        
        weights = pd.Series(model.coef_, index=factor_returns.columns)
        weights = weights / weights.sum()
        
    elif method == 'risk_parity':
        # 风险平价权重
        weights = risk_parity_weights(factor_returns)
        
    elif method == 'max_sharpe':
        # 最大夏普比率权重
        weights = max_sharpe_weights(factor_returns, return_data)
    
    return weights

def risk_parity_weights(factor_returns):
    """
    风险平价权重计算
    """
    from scipy.optimize import minimize
    
    cov_matrix = factor_returns.cov()
    n_factors = len(factor_returns.columns)
    
    def risk_parity_objective(weights):
        portfolio_vol = np.sqrt(weights.T @ cov_matrix @ weights)
        marginal_contrib = cov_matrix @ weights / portfolio_vol
        contrib = weights * marginal_contrib
        target_contrib = portfolio_vol / n_factors
        return np.sum((contrib - target_contrib)**2)
    
    # 约束条件
    constraints = {'type': 'eq', 'fun': lambda x: np.sum(x) - 1}
    bounds = tuple((0, 1) for _ in range(n_factors))
    
    # 初始权重
    x0 = np.array([1/n_factors] * n_factors)
    
    # 优化
    result = minimize(risk_parity_objective, x0, method='SLSQP', 
                     bounds=bounds, constraints=constraints)
    
    weights = pd.Series(result.x, index=factor_returns.columns)
    return weights

def max_sharpe_weights(factor_returns, target_return, risk_free_rate=0.02):
    """
    最大夏普比率权重计算
    """
    from scipy.optimize import minimize
    
    mean_returns = factor_returns.mean() * 252
    cov_matrix = factor_returns.cov() * 252
    n_factors = len(factor_returns.columns)
    
    def negative_sharpe_ratio(weights):
        portfolio_return = np.sum(mean_returns * weights)
        portfolio_vol = np.sqrt(np.dot(weights.T, np.dot(cov_matrix, weights)))
        return -(portfolio_return - risk_free_rate) / portfolio_vol
    
    # 约束条件
    constraints = {'type': 'eq', 'fun': lambda x: np.sum(x) - 1}
    bounds = tuple((0, 1) for _ in range(n_factors))
    
    # 初始权重
    x0 = np.array([1/n_factors] * n_factors)
    
    # 优化
    result = minimize(negative_sharpe_ratio, x0, method='SLSQP',
                     bounds=bounds, constraints=constraints)
    
    weights = pd.Series(result.x, index=factor_returns.columns)
    return weights
```

## 📊 实战案例：A股市场多因子模型构建

```python
def build_multifactor_model_case_study():
    """
    A股市场多因子模型构建案例
    """
    
    # 1. 数据准备（示例数据结构）
    print("=== A股多因子模型构建案例 ===")
    
    # 模拟数据
    np.random.seed(42)
    dates = pd.date_range('2020-01-01', '2023-12-31', freq='D')
    n_stocks = 100
    
    # 模拟价格数据
    price_data = pd.DataFrame(
        np.random.randn(len(dates), n_stocks).cumsum() + 100,
        index=dates,
        columns=[f'Stock_{i:03d}' for i in range(n_stocks)]
    )
    
    # 模拟基本面数据
    fundamental_data = pd.DataFrame({
        'market_cap': np.random.lognormal(10, 1, n_stocks),
        'book_value': np.random.lognormal(8, 1, n_stocks),
        'earnings': np.random.lognormal(6, 1, n_stocks),
        'revenue': np.random.lognormal(9, 1, n_stocks)
    }, index=price_data.columns)
    
    # 模拟成交量数据
    volume_data = pd.DataFrame(
        np.random.lognormal(10, 1, (len(dates), n_stocks)),
        index=dates,
        columns=price_data.columns
    )
    
    # 2. 因子计算
    print("\n1. 计算因子数据...")
    
    # 价值因子
    value_factors = calculate_value_factors(price_data, fundamental_data)
    
    # 动量因子
    momentum_factors = calculate_momentum_factors(price_data)
    
    # 波动率因子
    volatility_factors = calculate_volatility_factors(price_data, volume_data)
    
    # 规模因子
    size_factors = calculate_size_factors(fundamental_data['market_cap'], price_data)
    
    # 合并所有因子
    all_factors = pd.concat([
        value_factors.iloc[-252:],  # 最近一年数据
        momentum_factors.iloc[-252:],
        volatility_factors.iloc[-252:],
        size_factors.iloc[-252:]
    ], axis=1).dropna()
    
    print(f"共计算了 {len(all_factors.columns)} 个因子")
    
    # 3. 因子筛选
    print("\n2. 因子筛选...")
    
    returns = price_data.pct_change().iloc[-252:].dropna()
    
    # 对每只股票进行因子筛选
    selected_factors_dict = {}
    for stock in returns.columns[:5]:  # 选择前5只股票作为示例
        factor_selection_result = factor_selection(
            all_factors, returns[stock], method='lasso'
        )
        selected_factors_dict[stock] = factor_selection_result['selected_factors']
    
    # 统计最常用的因子
    all_selected = []
    for factors in selected_factors_dict.values():
        all_selected.extend(factors)
    
    factor_counts = pd.Series(all_selected).value_counts()
    top_factors = factor_counts.head(8).index.tolist()
    
    print(f"选择的核心因子: {top_factors}")
    
    # 4. 因子正交化
    print("\n3. 因子正交化...")
    
    core_factors = all_factors[top_factors]
    orthogonal_factors = factor_orthogonalization(core_factors, method='gram_schmidt')
    
    # 5. 模型构建和回测
    print("\n4. 模型构建和回测...")
    
    # 构建因子投资组合
    factor_scores = orthogonal_factors.mean(axis=1)  # 简单平均
    
    # 按因子得分选股
    top_stocks = factor_scores.nlargest(20).index  # 选择前20只股票
    
    # 计算投资组合收益
    portfolio_weights = pd.Series(1/len(top_stocks), index=top_stocks)
    portfolio_returns = (returns[top_stocks] * portfolio_weights).sum(axis=1)
    
    # 基准收益（等权重所有股票）
    benchmark_returns = returns.mean(axis=1)
    
    # 6. 绩效评估
    print("\n5. 绩效评估...")
    
    portfolio_cumret = (1 + portfolio_returns).cumprod() - 1
    benchmark_cumret = (1 + benchmark_returns).cumprod() - 1
    
    # 计算绩效指标
    portfolio_annual_return = portfolio_returns.mean() * 252
    portfolio_volatility = portfolio_returns.std() * np.sqrt(252)
    portfolio_sharpe = portfolio_annual_return / portfolio_volatility
    
    benchmark_annual_return = benchmark_returns.mean() * 252
    benchmark_volatility = benchmark_returns.std() * np.sqrt(252)
    benchmark_sharpe = benchmark_annual_return / benchmark_volatility
    
    max_drawdown = (portfolio_cumret / portfolio_cumret.cummax() - 1).min()
    
    print(f"\n=== 绩效报告 ===")
    print(f"因子组合年化收益率: {portfolio_annual_return:.2%}")
    print(f"基准组合年化收益率: {benchmark_annual_return:.2%}")
    print(f"超额年化收益率: {portfolio_annual_return - benchmark_annual_return:.2%}")
    print(f"因子组合夏普比率: {portfolio_sharpe:.3f}")
    print(f"基准组合夏普比率: {benchmark_sharpe:.3f}")
    print(f"最大回撤: {max_drawdown:.2%}")
    
    return {
        'factors': all_factors,
        'orthogonal_factors': orthogonal_factors,
        'portfolio_returns': portfolio_returns,
        'benchmark_returns': benchmark_returns,
        'selected_stocks': top_stocks
    }

# 运行案例
if __name__ == "__main__":
    case_results = build_multifactor_model_case_study()
```

## 📚 本章小结

本章详细介绍了多因子模型的理论基础和实践方法：

### 🔑 核心要点

1. **因子投资理论**: 从CAPM到APT，理解因子投资的理论演进
2. **因子分类**: 基本面因子、技术因子、市场结构因子的特征和计算
3. **因子检验**: 单因子检验、IC分析、回归分析等有效性验证方法
4. **模型构建**: 因子选择、正交化、权重确定的完整流程

### 📈 实践要点

- 因子计算需要高质量的数据支持
- 因子有效性会随时间变化，需要持续监控
- 因子间相关性需要通过正交化处理
- 权重确定方法直接影响策略效果

### 🎯 下章预告

下一章将深入探讨**投资组合优化理论、风险平价策略、动态调仓策略**等高级组合管理技术，以及如何将多因子模型应用到实际的投资组合构建中。

---

**风险提示**: 因子投资存在因子失效风险，历史表现不代表未来收益，请结合实际情况谨慎使用。