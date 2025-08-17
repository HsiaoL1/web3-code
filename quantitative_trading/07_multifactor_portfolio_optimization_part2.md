# 第7章 多因子模型与组合优化（下篇）

## 📚 本章概览

本章将深入探讨投资组合优化理论在多因子模型中的应用，包括现代投资组合理论、风险平价策略、Black-Litterman模型以及动态调仓策略。这些技术将帮助我们将多因子模型转化为实际可执行的投资策略。

### 🎯 学习目标

- 掌握现代投资组合理论的核心概念
- 理解风险平价策略的原理和实现
- 学会Black-Litterman模型的应用
- 掌握动态调仓和再平衡策略
- 了解组合优化中的约束条件设计

## 第一节：现代投资组合理论

### 1.1 马科维茨均值-方差优化

```python
import numpy as np
import pandas as pd
from scipy.optimize import minimize
import matplotlib.pyplot as plt
from scipy.stats import norm

def markowitz_optimization(expected_returns, cov_matrix, risk_aversion=1.0, constraints=None):
    """
    马科维茨均值-方差优化
    
    Parameters:
    expected_returns: 期望收益率向量
    cov_matrix: 协方差矩阵
    risk_aversion: 风险厌恶系数
    constraints: 约束条件字典
    """
    
    n_assets = len(expected_returns)
    
    def objective_function(weights):
        # 期望收益
        portfolio_return = np.sum(weights * expected_returns)
        # 组合方差
        portfolio_variance = np.dot(weights.T, np.dot(cov_matrix, weights))
        # 效用函数 = 期望收益 - 0.5 * 风险厌恶系数 * 方差
        return -(portfolio_return - 0.5 * risk_aversion * portfolio_variance)
    
    # 默认约束：权重和为1
    cons = [{'type': 'eq', 'fun': lambda x: np.sum(x) - 1}]
    
    # 添加自定义约束
    if constraints:
        if 'max_weight' in constraints:
            max_weight = constraints['max_weight']
            cons.append({'type': 'ineq', 'fun': lambda x: max_weight - np.max(x)})
        
        if 'min_weight' in constraints:
            min_weight = constraints['min_weight']
            cons.append({'type': 'ineq', 'fun': lambda x: np.min(x) - min_weight})
        
        if 'sector_limits' in constraints:
            sector_limits = constraints['sector_limits']
            sector_mapping = constraints['sector_mapping']
            
            for sector, limit in sector_limits.items():
                sector_indices = [i for i, s in enumerate(sector_mapping) if s == sector]
                cons.append({
                    'type': 'ineq', 
                    'fun': lambda x, indices=sector_indices: limit - np.sum(x[indices])
                })
    
    # 权重边界 [0, 1]（允许做空时可以设置为 [-1, 1]）
    bounds = tuple((0, 1) for _ in range(n_assets))
    
    # 初始权重
    x0 = np.array([1/n_assets] * n_assets)
    
    # 优化
    result = minimize(objective_function, x0, method='SLSQP', 
                     bounds=bounds, constraints=cons)
    
    if result.success:
        optimal_weights = result.x
        portfolio_return = np.sum(optimal_weights * expected_returns)
        portfolio_risk = np.sqrt(np.dot(optimal_weights.T, np.dot(cov_matrix, optimal_weights)))
        
        return {
            'weights': optimal_weights,
            'expected_return': portfolio_return,
            'volatility': portfolio_risk,
            'sharpe_ratio': portfolio_return / portfolio_risk,
            'success': True
        }
    else:
        return {'success': False, 'message': result.message}

def efficient_frontier(expected_returns, cov_matrix, n_portfolios=50):
    """
    计算有效前沿
    """
    
    n_assets = len(expected_returns)
    
    # 计算最小方差组合
    def min_variance_objective(weights):
        return np.dot(weights.T, np.dot(cov_matrix, weights))
    
    cons = [{'type': 'eq', 'fun': lambda x: np.sum(x) - 1}]
    bounds = tuple((0, 1) for _ in range(n_assets))
    x0 = np.array([1/n_assets] * n_assets)
    
    min_var_result = minimize(min_variance_objective, x0, method='SLSQP',
                             bounds=bounds, constraints=cons)
    
    min_var_return = np.sum(min_var_result.x * expected_returns)
    max_return = np.max(expected_returns)
    
    # 在不同目标收益率下优化
    target_returns = np.linspace(min_var_return, max_return, n_portfolios)
    efficient_portfolios = []
    
    for target_return in target_returns:
        # 添加收益率约束
        return_constraint = {'type': 'eq', 'fun': lambda x: np.sum(x * expected_returns) - target_return}
        cons_with_return = cons + [return_constraint]
        
        result = minimize(min_variance_objective, x0, method='SLSQP',
                         bounds=bounds, constraints=cons_with_return)
        
        if result.success:
            portfolio_risk = np.sqrt(result.fun)
            efficient_portfolios.append({
                'return': target_return,
                'risk': portfolio_risk,
                'weights': result.x,
                'sharpe_ratio': target_return / portfolio_risk
            })
    
    return pd.DataFrame(efficient_portfolios)
```

### 1.2 多目标优化

```python
def multi_objective_optimization(expected_returns, cov_matrix, objectives):
    """
    多目标优化
    
    Parameters:
    expected_returns: 期望收益率
    cov_matrix: 协方差矩阵
    objectives: 目标函数列表，每个目标包含权重和函数
    """
    
    n_assets = len(expected_returns)
    
    def combined_objective(weights):
        total_objective = 0
        
        for obj in objectives:
            obj_type = obj['type']
            obj_weight = obj['weight']
            
            if obj_type == 'return':
                # 最大化收益（转为最小化负收益）
                value = -np.sum(weights * expected_returns)
            
            elif obj_type == 'risk':
                # 最小化风险
                value = np.dot(weights.T, np.dot(cov_matrix, weights))
            
            elif obj_type == 'concentration':
                # 最小化集中度（最大化分散度）
                value = np.sum(weights**2)  # Herfindahl指数
            
            elif obj_type == 'turnover':
                # 最小化换手率（需要前期权重）
                if 'prev_weights' in obj:
                    value = np.sum(np.abs(weights - obj['prev_weights']))
                else:
                    value = 0
            
            elif obj_type == 'tracking_error':
                # 最小化跟踪误差
                if 'benchmark_weights' in obj:
                    active_weights = weights - obj['benchmark_weights']
                    value = np.dot(active_weights.T, np.dot(cov_matrix, active_weights))
                else:
                    value = 0
            
            total_objective += obj_weight * value
        
        return total_objective
    
    # 约束和边界
    cons = [{'type': 'eq', 'fun': lambda x: np.sum(x) - 1}]
    bounds = tuple((0, 1) for _ in range(n_assets))
    x0 = np.array([1/n_assets] * n_assets)
    
    # 优化
    result = minimize(combined_objective, x0, method='SLSQP',
                     bounds=bounds, constraints=cons)
    
    return result

# 使用示例
def demo_multi_objective():
    """
    多目标优化演示
    """
    # 模拟数据
    n_assets = 10
    expected_returns = np.random.uniform(0.05, 0.15, n_assets)
    correlation_matrix = np.random.uniform(0.1, 0.7, (n_assets, n_assets))
    np.fill_diagonal(correlation_matrix, 1.0)
    volatilities = np.random.uniform(0.1, 0.3, n_assets)
    cov_matrix = np.outer(volatilities, volatilities) * correlation_matrix
    
    # 定义多个目标
    objectives = [
        {'type': 'return', 'weight': 0.4},      # 40%权重给收益
        {'type': 'risk', 'weight': 0.4},        # 40%权重给风险
        {'type': 'concentration', 'weight': 0.2} # 20%权重给分散度
    ]
    
    result = multi_objective_optimization(expected_returns, cov_matrix, objectives)
    
    if result.success:
        weights = result.x
        portfolio_return = np.sum(weights * expected_returns)
        portfolio_risk = np.sqrt(np.dot(weights.T, np.dot(cov_matrix, weights)))
        concentration = np.sum(weights**2)
        
        print(f"组合期望收益: {portfolio_return:.3f}")
        print(f"组合风险: {portfolio_risk:.3f}")
        print(f"集中度指数: {concentration:.3f}")
        print(f"权重分布: {weights}")
    
    return result
```

## 第二节：风险平价策略

### 2.1 经典风险平价

```python
def risk_parity_optimization(cov_matrix, method='equal_risk_contribution'):
    """
    风险平价优化
    
    Parameters:
    cov_matrix: 协方差矩阵
    method: 'equal_risk_contribution' 或 'volatility_weighted'
    """
    
    n_assets = len(cov_matrix)
    
    if method == 'equal_risk_contribution':
        # 等风险贡献度
        def objective_function(weights):
            # 计算组合波动率
            portfolio_vol = np.sqrt(np.dot(weights.T, np.dot(cov_matrix, weights)))
            
            # 计算边际风险贡献
            marginal_contrib = np.dot(cov_matrix, weights) / portfolio_vol
            
            # 计算风险贡献
            risk_contrib = weights * marginal_contrib
            
            # 目标风险贡献（等权重）
            target_contrib = portfolio_vol / n_assets
            
            # 最小化风险贡献的偏离
            return np.sum((risk_contrib - target_contrib)**2)
    
    elif method == 'volatility_weighted':
        # 波动率倒数加权
        individual_vols = np.sqrt(np.diag(cov_matrix))
        inv_vol_weights = (1 / individual_vols) / np.sum(1 / individual_vols)
        return inv_vol_weights
    
    else:
        raise ValueError("Method must be 'equal_risk_contribution' or 'volatility_weighted'")
    
    if method == 'volatility_weighted':
        return inv_vol_weights
    
    # 约束条件
    constraints = [
        {'type': 'eq', 'fun': lambda x: np.sum(x) - 1}  # 权重和为1
    ]
    
    # 边界条件（权重非负）
    bounds = tuple((0, 1) for _ in range(n_assets))
    
    # 初始权重
    x0 = np.array([1/n_assets] * n_assets)
    
    # 优化
    result = minimize(objective_function, x0, method='SLSQP',
                     bounds=bounds, constraints=constraints)
    
    return result.x if result.success else None

def analyze_risk_contributions(weights, cov_matrix):
    """
    分析风险贡献度
    """
    
    # 组合波动率
    portfolio_vol = np.sqrt(np.dot(weights.T, np.dot(cov_matrix, weights)))
    
    # 边际风险贡献
    marginal_contrib = np.dot(cov_matrix, weights) / portfolio_vol
    
    # 风险贡献
    risk_contrib = weights * marginal_contrib
    
    # 风险贡献占比
    risk_contrib_pct = risk_contrib / portfolio_vol
    
    return {
        'portfolio_volatility': portfolio_vol,
        'marginal_contributions': marginal_contrib,
        'risk_contributions': risk_contrib,
        'risk_contribution_percentages': risk_contrib_pct
    }

def hierarchical_risk_parity(cov_matrix, returns_data=None):
    """
    分层风险平价（HRP）
    """
    
    from scipy.cluster.hierarchy import linkage, dendrogram, fcluster
    from scipy.spatial.distance import squareform
    
    # 计算相关系数矩阵
    corr_matrix = cov_matrix / np.outer(np.sqrt(np.diag(cov_matrix)), np.sqrt(np.diag(cov_matrix)))
    
    # 转换为距离矩阵
    distance_matrix = np.sqrt(0.5 * (1 - corr_matrix))
    
    # 层次聚类
    linkage_matrix = linkage(squareform(distance_matrix), method='ward')
    
    # 获取聚类排序
    def get_cluster_order(linkage_matrix):
        n = len(linkage_matrix) + 1
        order = []
        
        def traverse(node):
            if node < n:
                order.append(node)
            else:
                left_child = int(linkage_matrix[node - n, 0])
                right_child = int(linkage_matrix[node - n, 1])
                traverse(left_child)
                traverse(right_child)
        
        traverse(2 * n - 2)
        return order
    
    cluster_order = get_cluster_order(linkage_matrix)
    
    # 重新排列协方差矩阵
    reordered_cov = cov_matrix[np.ix_(cluster_order, cluster_order)]
    
    # 递归分配权重
    def recursive_bisection(cov_matrix):
        n = len(cov_matrix)
        if n == 1:
            return np.array([1.0])
        
        # 分割点
        mid = n // 2
        
        # 左右子组合的协方差矩阵
        left_cov = cov_matrix[:mid, :mid]
        right_cov = cov_matrix[mid:, mid:]
        
        # 计算子组合的方差
        left_weights = recursive_bisection(left_cov)
        right_weights = recursive_bisection(right_cov)
        
        left_var = np.dot(left_weights.T, np.dot(left_cov, left_weights))
        right_var = np.dot(right_weights.T, np.dot(right_cov, right_weights))
        
        # 根据方差倒数分配权重
        total_inv_var = 1/left_var + 1/right_var
        left_alloc = (1/left_var) / total_inv_var
        right_alloc = (1/right_var) / total_inv_var
        
        # 组合权重
        weights = np.zeros(n)
        weights[:mid] = left_alloc * left_weights
        weights[mid:] = right_alloc * right_weights
        
        return weights
    
    # 计算HRP权重
    hrp_weights = recursive_bisection(reordered_cov)
    
    # 还原到原始顺序
    original_weights = np.zeros(len(cov_matrix))
    for i, original_idx in enumerate(cluster_order):
        original_weights[original_idx] = hrp_weights[i]
    
    return {
        'weights': original_weights,
        'cluster_order': cluster_order,
        'linkage_matrix': linkage_matrix
    }
```

### 2.2 风险预算策略

```python
def risk_budgeting_optimization(cov_matrix, risk_budgets):
    """
    风险预算策略
    
    Parameters:
    cov_matrix: 协方差矩阵
    risk_budgets: 风险预算向量（应该和为1）
    """
    
    n_assets = len(cov_matrix)
    
    # 确保风险预算和为1
    risk_budgets = np.array(risk_budgets) / np.sum(risk_budgets)
    
    def objective_function(weights):
        # 组合波动率
        portfolio_vol = np.sqrt(np.dot(weights.T, np.dot(cov_matrix, weights)))
        
        # 边际风险贡献
        marginal_contrib = np.dot(cov_matrix, weights) / portfolio_vol
        
        # 实际风险贡献占比
        risk_contrib_pct = weights * marginal_contrib / portfolio_vol
        
        # 最小化与目标风险预算的偏离
        return np.sum((risk_contrib_pct - risk_budgets)**2)
    
    # 约束条件
    constraints = [
        {'type': 'eq', 'fun': lambda x: np.sum(x) - 1}
    ]
    
    # 边界条件
    bounds = tuple((0, 1) for _ in range(n_assets))
    
    # 初始权重
    x0 = risk_budgets.copy()  # 以风险预算作为初始权重
    
    # 优化
    result = minimize(objective_function, x0, method='SLSQP',
                     bounds=bounds, constraints=constraints)
    
    return result.x if result.success else None

def adaptive_risk_parity(returns_data, lookback_window=252, rebalance_freq=20):
    """
    自适应风险平价策略
    """
    
    n_assets = returns_data.shape[1]
    n_periods = len(returns_data)
    
    # 存储权重和组合收益
    weights_history = []
    portfolio_returns = []
    
    for t in range(lookback_window, n_periods, rebalance_freq):
        # 历史数据窗口
        historical_returns = returns_data.iloc[t-lookback_window:t]
        
        # 计算协方差矩阵（可以使用指数加权或其他方法）
        cov_matrix = historical_returns.cov().values
        
        # 计算风险平价权重
        weights = risk_parity_optimization(cov_matrix)
        
        if weights is not None:
            weights_history.append(weights)
            
            # 计算下一期组合收益
            next_period_end = min(t + rebalance_freq, n_periods)
            period_returns = returns_data.iloc[t:next_period_end]
            
            for _, daily_returns in period_returns.iterrows():
                portfolio_return = np.sum(weights * daily_returns.values)
                portfolio_returns.append(portfolio_return)
    
    return np.array(weights_history), np.array(portfolio_returns)
```

## 第三节：Black-Litterman模型

### 3.1 经典Black-Litterman模型

```python
def black_litterman_optimization(market_caps, cov_matrix, risk_aversion=3.0, 
                               views=None, view_uncertainty=None, tau=0.025):
    """
    Black-Litterman模型
    
    Parameters:
    market_caps: 市场市值权重
    cov_matrix: 历史协方差矩阵
    risk_aversion: 风险厌恶系数
    views: 主观观点矩阵P（每行表示一个观点）
    view_uncertainty: 观点不确定性矩阵Omega
    tau: 缩放因子
    """
    
    # 步骤1: 计算隐含收益率（市场均衡收益率）
    implied_returns = risk_aversion * np.dot(cov_matrix, market_caps)
    
    # 如果没有主观观点，直接返回隐含收益率
    if views is None:
        return {
            'expected_returns': implied_returns,
            'posterior_cov': cov_matrix,
            'optimal_weights': market_caps
        }
    
    # 步骤2: 整合主观观点
    P = np.array(views['P'])  # 观点矩阵
    Q = np.array(views['Q'])  # 观点收益率
    
    # 如果没有提供观点不确定性矩阵，使用默认方法计算
    if view_uncertainty is None:
        # 使用He and Litterman (1999)的方法
        view_uncertainty = tau * np.dot(P, np.dot(cov_matrix, P.T))
    
    Omega = np.array(view_uncertainty)
    
    # 步骤3: 计算后验期望收益率
    # tau * Σ
    tau_sigma = tau * cov_matrix
    
    # 中间计算矩阵
    M1 = np.linalg.inv(tau_sigma)
    M2 = np.dot(P.T, np.dot(np.linalg.inv(Omega), P))
    M3 = np.dot(np.linalg.inv(tau_sigma), implied_returns)
    M4 = np.dot(P.T, np.dot(np.linalg.inv(Omega), Q))
    
    # 后验协方差矩阵
    posterior_cov = np.linalg.inv(M1 + M2)
    
    # 后验期望收益率
    posterior_returns = np.dot(posterior_cov, M3 + M4)
    
    # 步骤4: 计算最优权重
    optimal_weights = np.dot(np.linalg.inv(risk_aversion * cov_matrix), posterior_returns)
    
    return {
        'implied_returns': implied_returns,
        'posterior_returns': posterior_returns,
        'posterior_cov': posterior_cov,
        'optimal_weights': optimal_weights,
        'views_impact': posterior_returns - implied_returns
    }

def create_views_matrix(n_assets, views_dict):
    """
    创建观点矩阵
    
    Parameters:
    n_assets: 资产数量
    views_dict: 观点字典，格式如下：
    {
        'view1': {'assets': [0, 1], 'weights': [1, -1], 'expected_return': 0.02},
        'view2': {'assets': [2], 'weights': [1], 'expected_return': 0.05}
    }
    """
    
    n_views = len(views_dict)
    P = np.zeros((n_views, n_assets))
    Q = np.zeros(n_views)
    
    for i, (view_name, view_data) in enumerate(views_dict.items()):
        for j, asset_idx in enumerate(view_data['assets']):
            P[i, asset_idx] = view_data['weights'][j]
        Q[i] = view_data['expected_return']
    
    return {'P': P, 'Q': Q}

# 使用示例
def demo_black_litterman():
    """
    Black-Litterman模型演示
    """
    
    # 模拟数据
    n_assets = 5
    market_caps = np.array([0.3, 0.25, 0.2, 0.15, 0.1])  # 市场权重
    
    # 模拟协方差矩阵
    np.random.seed(42)
    correlation_matrix = 0.3 * np.ones((n_assets, n_assets))
    np.fill_diagonal(correlation_matrix, 1.0)
    volatilities = np.array([0.15, 0.18, 0.20, 0.25, 0.30])
    cov_matrix = np.outer(volatilities, volatilities) * correlation_matrix
    
    # 创建主观观点
    views_dict = {
        'view1': {  # 认为资产0比资产1多2%收益
            'assets': [0, 1], 
            'weights': [1, -1], 
            'expected_return': 0.02
        },
        'view2': {  # 认为资产2绝对收益5%
            'assets': [2], 
            'weights': [1], 
            'expected_return': 0.05
        }
    }
    
    views = create_views_matrix(n_assets, views_dict)
    
    # 运行Black-Litterman模型
    bl_result = black_litterman_optimization(
        market_caps=market_caps,
        cov_matrix=cov_matrix,
        risk_aversion=3.0,
        views=views,
        tau=0.025
    )
    
    print("=== Black-Litterman结果 ===")
    print(f"隐含收益率: {bl_result['implied_returns']}")
    print(f"后验收益率: {bl_result['posterior_returns']}")
    print(f"市场权重: {market_caps}")
    print(f"最优权重: {bl_result['optimal_weights']}")
    print(f"观点影响: {bl_result['views_impact']}")
    
    return bl_result
```

### 3.2 扩展Black-Litterman模型

```python
def dynamic_black_litterman(returns_data, market_caps, views_function, 
                          lookback_window=252, rebalance_freq=20):
    """
    动态Black-Litterman模型
    
    Parameters:
    returns_data: 收益率数据
    market_caps: 市场权重（可以是固定值或时间序列）
    views_function: 生成观点的函数
    """
    
    n_assets = returns_data.shape[1]
    n_periods = len(returns_data)
    
    weights_history = []
    portfolio_returns = []
    
    for t in range(lookback_window, n_periods, rebalance_freq):
        # 历史数据
        historical_returns = returns_data.iloc[t-lookback_window:t]
        
        # 计算协方差矩阵
        cov_matrix = historical_returns.cov().values
        
        # 获取当前期的市场权重
        if hasattr(market_caps, 'iloc'):
            current_market_caps = market_caps.iloc[t].values
        else:
            current_market_caps = market_caps
        
        # 生成观点
        views = views_function(historical_returns, t)
        
        # Black-Litterman优化
        bl_result = black_litterman_optimization(
            market_caps=current_market_caps,
            cov_matrix=cov_matrix,
            views=views
        )
        
        weights = bl_result['optimal_weights']
        # 权重归一化和约束
        weights = np.maximum(weights, 0)  # 非负约束
        weights = weights / np.sum(weights)  # 归一化
        
        weights_history.append(weights)
        
        # 计算组合收益
        next_period_end = min(t + rebalance_freq, n_periods)
        period_returns = returns_data.iloc[t:next_period_end]
        
        for _, daily_returns in period_returns.iterrows():
            portfolio_return = np.sum(weights * daily_returns.values)
            portfolio_returns.append(portfolio_return)
    
    return np.array(weights_history), np.array(portfolio_returns)

def momentum_views_function(historical_returns, current_time):
    """
    基于动量的观点生成函数
    """
    
    # 计算过去3个月的累计收益
    momentum_returns = historical_returns.iloc[-60:].mean() * 252  # 年化
    
    n_assets = len(momentum_returns)
    
    # 选择表现最好和最差的资产
    sorted_assets = momentum_returns.argsort()
    best_asset = sorted_assets[-1]
    worst_asset = sorted_assets[0]
    
    # 创建观点：最好的资产比最差的资产未来表现好
    relative_momentum = momentum_returns.iloc[best_asset] - momentum_returns.iloc[worst_asset]
    
    views_dict = {
        'momentum_view': {
            'assets': [best_asset, worst_asset],
            'weights': [1, -1],
            'expected_return': relative_momentum * 0.5  # 保守估计未来差异
        }
    }
    
    return create_views_matrix(n_assets, views_dict)

def mean_reversion_views_function(historical_returns, current_time):
    """
    基于均值回归的观点生成函数
    """
    
    # 计算长期均值和当前位置
    long_term_mean = historical_returns.mean() * 252
    recent_performance = historical_returns.iloc[-20:].mean() * 252
    
    # 计算偏离度
    deviation = recent_performance - long_term_mean
    
    n_assets = len(deviation)
    
    # 对偏离度最大的资产创建反向观点
    abs_deviation = np.abs(deviation)
    most_deviated_asset = abs_deviation.idxmax()
    most_deviated_idx = deviation.index.get_loc(most_deviated_asset)
    
    # 如果当前表现超出长期均值，预期会回归
    expected_correction = -deviation.iloc[most_deviated_idx] * 0.3  # 预期30%的回归
    
    views_dict = {
        'mean_reversion_view': {
            'assets': [most_deviated_idx],
            'weights': [1],
            'expected_return': expected_correction
        }
    }
    
    return create_views_matrix(n_assets, views_dict)
```

## 第四节：动态调仓策略

### 4.1 时间驱动调仓

```python
def time_based_rebalancing(returns_data, target_weights, rebalance_frequency='monthly'):
    """
    时间驱动的再平衡策略
    
    Parameters:
    returns_data: 收益率数据
    target_weights: 目标权重
    rebalance_frequency: 再平衡频率 ('daily', 'weekly', 'monthly', 'quarterly')
    """
    
    # 设置再平衡间隔
    freq_map = {
        'daily': 1,
        'weekly': 5,
        'monthly': 20,
        'quarterly': 60
    }
    
    rebalance_interval = freq_map.get(rebalance_frequency, 20)
    
    n_periods = len(returns_data)
    n_assets = len(target_weights)
    
    # 初始化
    current_weights = target_weights.copy()
    portfolio_values = [100]  # 初始组合价值
    weights_history = [current_weights.copy()]
    turnover_history = []
    
    for t in range(1, n_periods):
        # 获取当期收益率
        daily_returns = returns_data.iloc[t].values
        
        # 计算当期组合收益
        portfolio_return = np.sum(current_weights * daily_returns)
        
        # 更新组合价值
        new_portfolio_value = portfolio_values[-1] * (1 + portfolio_return)
        portfolio_values.append(new_portfolio_value)
        
        # 更新权重（由于价格变动导致的漂移）
        current_weights = current_weights * (1 + daily_returns)
        current_weights = current_weights / np.sum(current_weights)
        
        # 检查是否需要再平衡
        if t % rebalance_interval == 0:
            # 计算换手率
            turnover = np.sum(np.abs(target_weights - current_weights)) / 2
            turnover_history.append(turnover)
            
            # 再平衡到目标权重
            current_weights = target_weights.copy()
        
        weights_history.append(current_weights.copy())
    
    portfolio_returns = np.diff(portfolio_values) / portfolio_values[:-1]
    
    return {
        'portfolio_returns': portfolio_returns,
        'portfolio_values': portfolio_values,
        'weights_history': np.array(weights_history),
        'turnover_history': turnover_history,
        'average_turnover': np.mean(turnover_history) if turnover_history else 0
    }

def threshold_based_rebalancing(returns_data, target_weights, threshold=0.05):
    """
    阈值驱动的再平衡策略
    
    Parameters:
    returns_data: 收益率数据
    target_weights: 目标权重
    threshold: 再平衡阈值（权重偏离超过此值时触发再平衡）
    """
    
    n_periods = len(returns_data)
    n_assets = len(target_weights)
    
    # 初始化
    current_weights = target_weights.copy()
    portfolio_values = [100]
    weights_history = [current_weights.copy()]
    rebalance_dates = []
    turnover_history = []
    
    for t in range(1, n_periods):
        # 获取当期收益率
        daily_returns = returns_data.iloc[t].values
        
        # 计算当期组合收益
        portfolio_return = np.sum(current_weights * daily_returns)
        
        # 更新组合价值
        new_portfolio_value = portfolio_values[-1] * (1 + portfolio_return)
        portfolio_values.append(new_portfolio_value)
        
        # 更新权重
        current_weights = current_weights * (1 + daily_returns)
        current_weights = current_weights / np.sum(current_weights)
        
        # 检查是否超过阈值
        weight_deviation = np.abs(current_weights - target_weights)
        max_deviation = np.max(weight_deviation)
        
        if max_deviation > threshold:
            # 记录再平衡
            rebalance_dates.append(returns_data.index[t])
            
            # 计算换手率
            turnover = np.sum(np.abs(target_weights - current_weights)) / 2
            turnover_history.append(turnover)
            
            # 再平衡
            current_weights = target_weights.copy()
        
        weights_history.append(current_weights.copy())
    
    portfolio_returns = np.diff(portfolio_values) / portfolio_values[:-1]
    
    return {
        'portfolio_returns': portfolio_returns,
        'portfolio_values': portfolio_values,
        'weights_history': np.array(weights_history),
        'rebalance_dates': rebalance_dates,
        'turnover_history': turnover_history,
        'average_turnover': np.mean(turnover_history) if turnover_history else 0,
        'rebalance_frequency': len(rebalance_dates) / (n_periods / 252)  # 年化频率
    }
```

### 4.2 成本感知调仓

```python
def cost_aware_rebalancing(returns_data, target_weights, transaction_costs=0.001, 
                          min_trade_size=0.001):
    """
    考虑交易成本的再平衡策略
    
    Parameters:
    returns_data: 收益率数据
    target_weights: 目标权重
    transaction_costs: 交易成本（单边）
    min_trade_size: 最小交易规模
    """
    
    n_periods = len(returns_data)
    n_assets = len(target_weights)
    
    # 初始化
    current_weights = target_weights.copy()
    portfolio_values = [100]
    weights_history = [current_weights.copy()]
    total_costs = 0
    trade_history = []
    
    for t in range(1, n_periods):
        # 获取当期收益率
        daily_returns = returns_data.iloc[t].values
        
        # 计算当期组合收益
        portfolio_return = np.sum(current_weights * daily_returns)
        
        # 更新组合价值
        new_portfolio_value = portfolio_values[-1] * (1 + portfolio_return)
        portfolio_values.append(new_portfolio_value)
        
        # 更新权重
        current_weights = current_weights * (1 + daily_returns)
        current_weights = current_weights / np.sum(current_weights)
        
        # 计算需要的交易量
        weight_diff = target_weights - current_weights
        
        # 只对超过最小交易规模的进行交易
        trade_mask = np.abs(weight_diff) > min_trade_size
        
        if np.any(trade_mask):
            # 计算交易成本
            trade_volume = np.sum(np.abs(weight_diff[trade_mask]))
            costs = trade_volume * transaction_costs * new_portfolio_value
            
            # 计算预期收益改善（简化模型）
            expected_benefit = estimate_rebalancing_benefit(
                current_weights, target_weights, returns_data.iloc[max(0, t-20):t]
            )
            
            # 如果预期收益大于成本，则进行交易
            if expected_benefit > costs:
                current_weights[trade_mask] = target_weights[trade_mask]
                current_weights = current_weights / np.sum(current_weights)
                
                total_costs += costs
                trade_history.append({
                    'date': returns_data.index[t],
                    'trade_volume': trade_volume,
                    'costs': costs,
                    'expected_benefit': expected_benefit
                })
        
        weights_history.append(current_weights.copy())
    
    portfolio_returns = np.diff(portfolio_values) / portfolio_values[:-1]
    
    # 扣除交易成本后的收益
    net_portfolio_returns = portfolio_returns.copy()
    for trade in trade_history:
        trade_date_idx = returns_data.index.get_loc(trade['date']) - 1
        if trade_date_idx < len(net_portfolio_returns):
            net_portfolio_returns[trade_date_idx] -= trade['costs'] / portfolio_values[trade_date_idx]
    
    return {
        'gross_returns': portfolio_returns,
        'net_returns': net_portfolio_returns,
        'total_costs': total_costs,
        'trade_history': trade_history,
        'weights_history': np.array(weights_history)
    }

def estimate_rebalancing_benefit(current_weights, target_weights, historical_returns):
    """
    估计再平衡的预期收益改善
    """
    
    if len(historical_returns) == 0:
        return 0
    
    # 计算历史协方差矩阵
    cov_matrix = historical_returns.cov().values
    
    # 当前组合的预期风险
    current_risk = np.sqrt(np.dot(current_weights.T, np.dot(cov_matrix, current_weights)))
    
    # 目标组合的预期风险
    target_risk = np.sqrt(np.dot(target_weights.T, np.dot(cov_matrix, target_weights)))
    
    # 风险改善转化为收益改善（简化假设）
    risk_improvement = current_risk - target_risk
    expected_benefit = risk_improvement * 100  # 假设风险降低1%对应收益改善100bp
    
    return max(0, expected_benefit)
```

### 4.3 信号驱动调仓

```python
def signal_based_rebalancing(returns_data, signal_function, base_weights, 
                           signal_strength=0.1, max_deviation=0.2):
    """
    基于信号的动态调仓策略
    
    Parameters:
    returns_data: 收益率数据
    signal_function: 信号生成函数
    base_weights: 基础权重
    signal_strength: 信号强度参数
    max_deviation: 最大偏离度
    """
    
    n_periods = len(returns_data)
    n_assets = len(base_weights)
    
    # 初始化
    current_weights = base_weights.copy()
    portfolio_values = [100]
    weights_history = [current_weights.copy()]
    signals_history = []
    
    for t in range(1, n_periods):
        # 生成信号
        signals = signal_function(returns_data.iloc[:t])
        signals_history.append(signals)
        
        # 根据信号调整权重
        if signals is not None and len(signals) == n_assets:
            # 标准化信号
            normalized_signals = (signals - np.mean(signals)) / (np.std(signals) + 1e-8)
            
            # 计算目标权重调整
            weight_adjustments = signal_strength * normalized_signals
            
            # 应用调整，但限制最大偏离
            adjusted_weights = base_weights + weight_adjustments
            
            # 确保权重在合理范围内
            adjusted_weights = np.maximum(adjusted_weights, 
                                        base_weights - max_deviation)
            adjusted_weights = np.minimum(adjusted_weights, 
                                        base_weights + max_deviation)
            
            # 归一化
            adjusted_weights = np.maximum(adjusted_weights, 0)
            adjusted_weights = adjusted_weights / np.sum(adjusted_weights)
            
            current_weights = adjusted_weights
        
        # 计算当期收益
        daily_returns = returns_data.iloc[t].values
        portfolio_return = np.sum(current_weights * daily_returns)
        
        # 更新组合价值
        new_portfolio_value = portfolio_values[-1] * (1 + portfolio_return)
        portfolio_values.append(new_portfolio_value)
        
        weights_history.append(current_weights.copy())
    
    portfolio_returns = np.diff(portfolio_values) / portfolio_values[:-1]
    
    return {
        'portfolio_returns': portfolio_returns,
        'weights_history': np.array(weights_history),
        'signals_history': signals_history
    }

def momentum_signal_function(returns_data, lookback_periods=[20, 60, 120]):
    """
    动量信号函数
    """
    
    if len(returns_data) < max(lookback_periods):
        return None
    
    # 计算不同期间的动量
    momentum_scores = []
    
    for period in lookback_periods:
        period_returns = returns_data.iloc[-period:].mean() * 252  # 年化收益率
        momentum_scores.append(period_returns.values)
    
    # 综合动量得分（等权重平均）
    combined_momentum = np.mean(momentum_scores, axis=0)
    
    return combined_momentum

def mean_reversion_signal_function(returns_data, lookback_period=252):
    """
    均值回归信号函数
    """
    
    if len(returns_data) < lookback_period:
        return None
    
    # 计算长期均值
    long_term_mean = returns_data.iloc[-lookback_period:].mean()
    
    # 计算短期表现
    short_term_mean = returns_data.iloc[-20:].mean()
    
    # 计算偏离度（负值表示当前表现低于长期均值，预期回归）
    deviation = short_term_mean - long_term_mean
    
    # 转换为信号（均值回归策略）
    reversion_signal = -deviation.values  # 负偏离变为正信号
    
    return reversion_signal

def volatility_signal_function(returns_data, target_vol=0.15):
    """
    波动率信号函数
    """
    
    if len(returns_data) < 60:
        return None
    
    # 计算当前波动率
    current_vol = returns_data.iloc[-60:].std() * np.sqrt(252)
    
    # 计算与目标波动率的偏离
    vol_deviation = current_vol - target_vol
    
    # 生成信号：高波动率资产给予负信号（减少配置）
    vol_signal = -vol_deviation.values
    
    return vol_signal
```

## 📊 综合实战案例：多策略组合管理系统

```python
def comprehensive_portfolio_management_system():
    """
    综合组合管理系统案例
    """
    
    print("=== 综合组合管理系统 ===")
    
    # 1. 数据准备
    np.random.seed(42)
    dates = pd.date_range('2020-01-01', '2023-12-31', freq='D')
    n_assets = 8
    
    # 模拟收益率数据
    returns_data = pd.DataFrame(
        np.random.multivariate_normal(
            mean=[0.0008] * n_assets,  # 日均收益
            cov=0.3 * np.ones((n_assets, n_assets)) + 0.7 * np.eye(n_assets),
            size=len(dates)
        ) * 0.02,  # 调整波动率
        index=dates,
        columns=[f'Asset_{i+1}' for i in range(n_assets)]
    )
    
    # 模拟市值数据
    market_caps = np.array([0.2, 0.18, 0.15, 0.12, 0.10, 0.08, 0.07, 0.10])
    
    print(f"数据期间: {dates[0].date()} 到 {dates[-1].date()}")
    print(f"资产数量: {n_assets}")
    
    # 2. 多种组合策略比较
    strategies = {}
    
    # 等权重策略
    equal_weights = np.array([1/n_assets] * n_assets)
    strategies['Equal Weight'] = time_based_rebalancing(
        returns_data, equal_weights, 'monthly'
    )
    
    # 市值加权策略
    strategies['Market Cap'] = time_based_rebalancing(
        returns_data, market_caps, 'monthly'
    )
    
    # 风险平价策略
    print("\n计算风险平价权重...")
    recent_returns = returns_data.iloc[-252:]  # 最近一年数据
    cov_matrix = recent_returns.cov().values
    rp_weights = risk_parity_optimization(cov_matrix)
    
    if rp_weights is not None:
        strategies['Risk Parity'] = time_based_rebalancing(
            returns_data, rp_weights, 'monthly'
        )
    
    # 分层风险平价策略
    print("计算分层风险平价权重...")
    hrp_result = hierarchical_risk_parity(cov_matrix)
    strategies['HRP'] = time_based_rebalancing(
        returns_data, hrp_result['weights'], 'monthly'
    )
    
    # Black-Litterman策略
    print("运行Black-Litterman模型...")
    
    def simple_momentum_views(historical_returns, t):
        """简化的动量观点"""
        if len(historical_returns) < 60:
            return None
        
        momentum = historical_returns.iloc[-60:].mean() * 252
        sorted_assets = momentum.argsort()
        
        best_asset = sorted_assets[-1]
        worst_asset = sorted_assets[0]
        
        views_dict = {
            'momentum': {
                'assets': [best_asset, worst_asset],
                'weights': [1, -1],
                'expected_return': (momentum.iloc[best_asset] - momentum.iloc[worst_asset]) * 0.3
            }
        }
        
        return create_views_matrix(len(momentum), views_dict)
    
    bl_weights_history, bl_returns = dynamic_black_litterman(
        returns_data, market_caps, simple_momentum_views,
        lookback_window=252, rebalance_freq=60
    )
    
    strategies['Black-Litterman'] = {
        'portfolio_returns': bl_returns,
        'weights_history': bl_weights_history
    }
    
    # 信号驱动策略
    print("运行信号驱动策略...")
    
    def combined_signal_function(returns_data):
        """组合信号函数"""
        momentum_sig = momentum_signal_function(returns_data)
        reversion_sig = mean_reversion_signal_function(returns_data)
        vol_sig = volatility_signal_function(returns_data)
        
        if all(sig is not None for sig in [momentum_sig, reversion_sig, vol_sig]):
            # 等权重组合三种信号
            combined = (momentum_sig + reversion_sig + vol_sig) / 3
            return combined
        return None
    
    signal_result = signal_based_rebalancing(
        returns_data, combined_signal_function, equal_weights
    )
    strategies['Signal-Based'] = signal_result
    
    # 3. 策略绩效评估
    print("\n=== 策略绩效比较 ===")
    
    performance_metrics = {}
    
    for strategy_name, strategy_result in strategies.items():
        if 'portfolio_returns' in strategy_result:
            returns = strategy_result['portfolio_returns']
            
            # 计算绩效指标
            annual_return = np.mean(returns) * 252
            annual_vol = np.std(returns) * np.sqrt(252)
            sharpe_ratio = annual_return / annual_vol if annual_vol > 0 else 0
            
            # 计算最大回撤
            cumulative_returns = (1 + returns).cumprod()
            running_max = np.maximum.accumulate(cumulative_returns)
            drawdown = (cumulative_returns - running_max) / running_max
            max_drawdown = np.min(drawdown)
            
            # 胜率
            win_rate = np.sum(returns > 0) / len(returns)
            
            # Calmar比率
            calmar_ratio = annual_return / abs(max_drawdown) if max_drawdown != 0 else 0
            
            performance_metrics[strategy_name] = {
                'Annual Return': annual_return,
                'Annual Volatility': annual_vol,
                'Sharpe Ratio': sharpe_ratio,
                'Max Drawdown': max_drawdown,
                'Calmar Ratio': calmar_ratio,
                'Win Rate': win_rate,
                'Total Return': cumulative_returns[-1] - 1 if len(cumulative_returns) > 0 else 0
            }
    
    # 创建绩效比较表
    performance_df = pd.DataFrame(performance_metrics).T
    
    print(performance_df.round(4))
    
    # 4. 风险分析
    print("\n=== 风险分析 ===")
    
    for strategy_name, strategy_result in strategies.items():
        if 'portfolio_returns' in strategy_result and 'weights_history' in strategy_result:
            returns = strategy_result['portfolio_returns']
            weights_history = strategy_result['weights_history']
            
            # 计算VaR和CVaR
            var_95 = np.percentile(returns, 5)
            cvar_95 = np.mean(returns[returns <= var_95])
            
            # 权重集中度
            if len(weights_history) > 0:
                avg_concentration = np.mean([np.sum(w**2) for w in weights_history])
                max_weight = np.mean([np.max(w) for w in weights_history])
                
                print(f"\n{strategy_name}:")
                print(f"  VaR (95%): {var_95:.4f}")
                print(f"  CVaR (95%): {cvar_95:.4f}")
                print(f"  平均集中度: {avg_concentration:.4f}")
                print(f"  平均最大权重: {max_weight:.4f}")
    
    # 5. 策略相关性分析
    print("\n=== 策略相关性分析 ===")
    
    strategy_returns_df = pd.DataFrame()
    for strategy_name, strategy_result in strategies.items():
        if 'portfolio_returns' in strategy_result:
            returns = strategy_result['portfolio_returns']
            # 对齐数据长度
            min_length = min(len(returns), len(strategy_returns_df) if not strategy_returns_df.empty else len(returns))
            if strategy_returns_df.empty:
                strategy_returns_df[strategy_name] = returns[:min_length]
            else:
                strategy_returns_df[strategy_name] = returns[:min_length]
    
    if not strategy_returns_df.empty:
        correlation_matrix = strategy_returns_df.corr()
        print("\n策略收益率相关性矩阵:")
        print(correlation_matrix.round(3))
    
    return {
        'strategies': strategies,
        'performance_metrics': performance_df,
        'correlation_matrix': correlation_matrix if not strategy_returns_df.empty else None,
        'returns_data': returns_data
    }

# 运行综合案例
if __name__ == "__main__":
    comprehensive_results = comprehensive_portfolio_management_system()
```

## 📚 本章小结

本章深入探讨了投资组合优化的高级技术和实战应用：

### 🔑 核心要点

1. **现代投资组合理论**: 马科维茨优化、有效前沿、多目标优化
2. **风险平价策略**: 等风险贡献、分层风险平价、风险预算
3. **Black-Litterman模型**: 贝叶斯方法整合主观观点和市场信息
4. **动态调仓策略**: 时间驱动、阈值驱动、成本感知、信号驱动

### 📈 实践要点

- 不同优化方法适用于不同的市场环境和投资目标
- 风险平价在分散化方面表现优异，但可能牺牲收益
- Black-Litterman模型能够有效整合定量和定性信息
- 动态调仓需要平衡收益改善和交易成本

### 🎯 下章预告

下一章将探讨**机器学习在量化投资中的应用**，包括监督学习、无监督学习、强化学习在投资策略中的具体应用和实现方法。

---

**风险提示**: 组合优化基于历史数据和模型假设，实际效果可能与预期不符，投资决策需结合多种因素综合考虑。