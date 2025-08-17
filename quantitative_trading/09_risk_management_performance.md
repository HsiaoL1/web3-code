# 第9章 风险管理与绩效评估

## 📚 本章概览

风险管理是量化投资的核心组成部分，它不仅关乎资本的保护，更是可持续盈利的基石。本章将深入探讨现代风险管理理论与实践，包括风险度量方法、VaR与CVaR计算、回撤控制策略、绩效评估指标以及归因分析技术。

### 🎯 学习目标

- 掌握风险度量的核心概念和计算方法
- 理解VaR、CVaR等风险指标的应用
- 学会设计和实施回撤控制系统
- 掌握绩效评估的全面框架
- 了解绩效归因分析技术

## 第一节：风险度量指标

### 1.1 基础风险指标

```python
import numpy as np
import pandas as pd
from scipy import stats
from scipy.optimize import minimize
import matplotlib.pyplot as plt
import seaborn as sns
from typing import Union, Dict, List, Tuple
import warnings
warnings.filterwarnings('ignore')

class RiskMetrics:
    """
    风险度量指标计算器
    """
    
    def __init__(self, returns: pd.Series, benchmark_returns: pd.Series = None, 
                 risk_free_rate: float = 0.02):
        self.returns = returns.dropna()
        self.benchmark_returns = benchmark_returns
        self.risk_free_rate = risk_free_rate
        self.trading_days = 252
        
    def calculate_basic_metrics(self) -> Dict[str, float]:
        """
        计算基础风险指标
        """
        
        metrics = {}
        
        # 收益率统计
        metrics['total_return'] = (1 + self.returns).prod() - 1
        metrics['annualized_return'] = (1 + self.returns.mean()) ** self.trading_days - 1
        metrics['annualized_volatility'] = self.returns.std() * np.sqrt(self.trading_days)
        
        # 夏普比率
        excess_return = metrics['annualized_return'] - self.risk_free_rate
        metrics['sharpe_ratio'] = excess_return / metrics['annualized_volatility'] if metrics['annualized_volatility'] > 0 else 0
        
        # 偏度和峰度
        metrics['skewness'] = self.returns.skew()
        metrics['kurtosis'] = self.returns.kurtosis()
        metrics['excess_kurtosis'] = metrics['kurtosis'] - 3
        
        # 胜率
        metrics['win_rate'] = (self.returns > 0).mean()
        
        # 最佳和最差单日收益
        metrics['best_day'] = self.returns.max()
        metrics['worst_day'] = self.returns.min()
        
        # 收益分布统计
        metrics['positive_returns_mean'] = self.returns[self.returns > 0].mean()
        metrics['negative_returns_mean'] = self.returns[self.returns < 0].mean()
        
        # 盈亏比
        if not pd.isna(metrics['negative_returns_mean']) and metrics['negative_returns_mean'] != 0:
            metrics['profit_loss_ratio'] = metrics['positive_returns_mean'] / abs(metrics['negative_returns_mean'])
        else:
            metrics['profit_loss_ratio'] = np.inf
        
        return metrics
    
    def calculate_drawdown_metrics(self) -> Dict[str, float]:
        """
        计算回撤指标
        """
        
        # 累计收益曲线
        cumulative_returns = (1 + self.returns).cumprod()
        
        # 历史最高点
        running_max = cumulative_returns.expanding().max()
        
        # 回撤
        drawdown = (cumulative_returns - running_max) / running_max
        
        metrics = {}
        metrics['max_drawdown'] = drawdown.min()
        metrics['current_drawdown'] = drawdown.iloc[-1]
        
        # 回撤持续时间
        drawdown_periods = self._calculate_drawdown_periods(drawdown)
        
        if drawdown_periods:
            metrics['max_drawdown_duration'] = max(period['duration'] for period in drawdown_periods)
            metrics['avg_drawdown_duration'] = np.mean([period['duration'] for period in drawdown_periods])
            metrics['drawdown_count'] = len(drawdown_periods)
        else:
            metrics['max_drawdown_duration'] = 0
            metrics['avg_drawdown_duration'] = 0
            metrics['drawdown_count'] = 0
        
        # Calmar比率
        if abs(metrics['max_drawdown']) > 0:
            annualized_return = (1 + self.returns.mean()) ** self.trading_days - 1
            metrics['calmar_ratio'] = annualized_return / abs(metrics['max_drawdown'])
        else:
            metrics['calmar_ratio'] = np.inf
        
        # 回撤风险调整收益率
        if metrics['max_drawdown'] != 0:
            metrics['sterling_ratio'] = annualized_return / abs(metrics['max_drawdown'])
        else:
            metrics['sterling_ratio'] = np.inf
        
        return metrics, drawdown
    
    def _calculate_drawdown_periods(self, drawdown: pd.Series) -> List[Dict]:
        """
        计算回撤期间
        """
        
        periods = []
        in_drawdown = False
        start_date = None
        
        for date, dd in drawdown.items():
            if dd < -0.001 and not in_drawdown:  # 开始回撤（阈值1bp）
                in_drawdown = True
                start_date = date
                peak_value = (1 + self.returns.loc[:date]).prod()
                
            elif dd >= -0.001 and in_drawdown:  # 结束回撤
                in_drawdown = False
                end_date = date
                trough_value = (1 + self.returns.loc[:date]).prod()
                
                duration = (end_date - start_date).days
                max_dd = drawdown.loc[start_date:end_date].min()
                
                periods.append({
                    'start_date': start_date,
                    'end_date': end_date,
                    'duration': duration,
                    'max_drawdown': max_dd,
                    'peak_value': peak_value,
                    'trough_value': trough_value
                })
        
        # 处理未结束的回撤
        if in_drawdown:
            end_date = drawdown.index[-1]
            duration = (end_date - start_date).days
            max_dd = drawdown.loc[start_date:].min()
            trough_value = (1 + self.returns.loc[:end_date]).prod()
            
            periods.append({
                'start_date': start_date,
                'end_date': end_date,
                'duration': duration,
                'max_drawdown': max_dd,
                'peak_value': peak_value,
                'trough_value': trough_value
            })
        
        return periods
    
    def calculate_relative_metrics(self) -> Dict[str, float]:
        """
        计算相对指标（需要基准）
        """
        
        if self.benchmark_returns is None:
            return {}
        
        # 对齐数据
        aligned_data = pd.DataFrame({
            'portfolio': self.returns,
            'benchmark': self.benchmark_returns
        }).dropna()
        
        portfolio_returns = aligned_data['portfolio']
        benchmark_returns = aligned_data['benchmark']
        
        metrics = {}
        
        # 超额收益
        excess_returns = portfolio_returns - benchmark_returns
        metrics['excess_return_annualized'] = excess_returns.mean() * self.trading_days
        
        # 跟踪误差
        metrics['tracking_error'] = excess_returns.std() * np.sqrt(self.trading_days)
        
        # 信息比率
        if metrics['tracking_error'] > 0:
            metrics['information_ratio'] = metrics['excess_return_annualized'] / metrics['tracking_error']
        else:
            metrics['information_ratio'] = 0
        
        # Beta
        covariance = np.cov(portfolio_returns, benchmark_returns)[0, 1]
        benchmark_variance = benchmark_returns.var()
        
        if benchmark_variance > 0:
            metrics['beta'] = covariance / benchmark_variance
        else:
            metrics['beta'] = 1
        
        # Alpha
        portfolio_annual_return = (1 + portfolio_returns.mean()) ** self.trading_days - 1
        benchmark_annual_return = (1 + benchmark_returns.mean()) ** self.trading_days - 1
        
        metrics['alpha'] = portfolio_annual_return - (self.risk_free_rate + metrics['beta'] * (benchmark_annual_return - self.risk_free_rate))
        
        # 相关系数
        metrics['correlation'] = portfolio_returns.corr(benchmark_returns)
        
        # 上行/下行捕获率
        up_periods = benchmark_returns > 0
        down_periods = benchmark_returns < 0
        
        if up_periods.sum() > 0:
            metrics['upside_capture'] = portfolio_returns[up_periods].mean() / benchmark_returns[up_periods].mean()
        else:
            metrics['upside_capture'] = 1
            
        if down_periods.sum() > 0:
            metrics['downside_capture'] = portfolio_returns[down_periods].mean() / benchmark_returns[down_periods].mean()
        else:
            metrics['downside_capture'] = 1
        
        return metrics
    
    def calculate_downside_metrics(self, target_return: float = 0) -> Dict[str, float]:
        """
        计算下行风险指标
        """
        
        metrics = {}
        
        # 下行收益
        downside_returns = self.returns[self.returns < target_return]
        
        if len(downside_returns) > 0:
            # 下行偏差
            downside_deviation = np.sqrt(((downside_returns - target_return) ** 2).mean())
            metrics['downside_deviation'] = downside_deviation * np.sqrt(self.trading_days)
            
            # Sortino比率
            annualized_return = (1 + self.returns.mean()) ** self.trading_days - 1
            excess_return = annualized_return - self.risk_free_rate
            
            if metrics['downside_deviation'] > 0:
                metrics['sortino_ratio'] = excess_return / metrics['downside_deviation']
            else:
                metrics['sortino_ratio'] = np.inf
            
            # 下行频率
            metrics['downside_frequency'] = len(downside_returns) / len(self.returns)
            
            # 平均下行收益
            metrics['average_downside_return'] = downside_returns.mean()
        else:
            metrics['downside_deviation'] = 0
            metrics['sortino_ratio'] = np.inf
            metrics['downside_frequency'] = 0
            metrics['average_downside_return'] = 0
        
        return metrics
```

### 1.2 高级风险指标

```python
class AdvancedRiskMetrics:
    """
    高级风险度量指标
    """
    
    def __init__(self, returns: pd.Series, confidence_levels: List[float] = [0.95, 0.99]):
        self.returns = returns.dropna()
        self.confidence_levels = confidence_levels
        self.trading_days = 252
        
    def calculate_var_cvar(self, method: str = 'historical') -> Dict[str, Dict[str, float]]:
        """
        计算VaR和CVaR
        
        Parameters:
        method: 'historical', 'parametric', 'monte_carlo'
        """
        
        results = {}
        
        for confidence_level in self.confidence_levels:
            alpha = 1 - confidence_level
            
            if method == 'historical':
                # 历史模拟法
                var = self.returns.quantile(alpha)
                cvar = self.returns[self.returns <= var].mean()
                
            elif method == 'parametric':
                # 参数法（假设正态分布）
                mean_return = self.returns.mean()
                std_return = self.returns.std()
                var = mean_return + std_return * stats.norm.ppf(alpha)
                
                # CVaR计算
                standardized_var = (var - mean_return) / std_return
                cvar = mean_return - std_return * stats.norm.pdf(standardized_var) / alpha
                
            elif method == 'monte_carlo':
                # 蒙特卡洛模拟
                np.random.seed(42)
                simulated_returns = np.random.normal(
                    self.returns.mean(), 
                    self.returns.std(), 
                    size=10000
                )
                var = np.percentile(simulated_returns, alpha * 100)
                cvar = simulated_returns[simulated_returns <= var].mean()
            
            # 年化VaR和CVaR
            var_annual = var * np.sqrt(self.trading_days)
            cvar_annual = cvar * np.sqrt(self.trading_days)
            
            results[f'{confidence_level:.0%}'] = {
                'var_daily': var,
                'cvar_daily': cvar,
                'var_annual': var_annual,
                'cvar_annual': cvar_annual
            }
        
        return results
    
    def calculate_expected_shortfall(self, confidence_level: float = 0.95) -> float:
        """
        计算期望损失（Expected Shortfall）
        """
        
        alpha = 1 - confidence_level
        var = self.returns.quantile(alpha)
        
        # ES就是CVaR
        es = self.returns[self.returns <= var].mean()
        
        return es
    
    def calculate_coherent_risk_measures(self) -> Dict[str, float]:
        """
        计算一致性风险度量
        """
        
        measures = {}
        
        # 期望损失（多个置信水平的平均）
        es_values = []
        for cl in self.confidence_levels:
            es = self.calculate_expected_shortfall(cl)
            es_values.append(es)
        
        measures['average_expected_shortfall'] = np.mean(es_values)
        
        # 条件尾部期望（Conditional Tail Expectation）
        tail_threshold = self.returns.quantile(0.1)  # 最差10%
        measures['conditional_tail_expectation'] = self.returns[self.returns <= tail_threshold].mean()
        
        # 最大损失（Maximum Loss）
        measures['maximum_loss'] = self.returns.min()
        
        return measures
    
    def calculate_garch_metrics(self) -> Dict[str, any]:
        """
        计算GARCH模型风险指标
        """
        
        try:
            from arch import arch_model
            
            # 拟合GARCH(1,1)模型
            model = arch_model(self.returns * 100, vol='GARCH', p=1, q=1)
            fitted_model = model.fit(disp='off')
            
            # 条件波动率
            conditional_volatility = fitted_model.conditional_volatility / 100
            
            # GARCH预测
            forecast = fitted_model.forecast(horizon=1)
            forecasted_variance = forecast.variance.iloc[-1, 0] / 10000
            forecasted_volatility = np.sqrt(forecasted_variance)
            
            return {
                'conditional_volatility': conditional_volatility,
                'forecasted_volatility': forecasted_volatility,
                'garch_params': fitted_model.params,
                'aic': fitted_model.aic,
                'bic': fitted_model.bic
            }
            
        except ImportError:
            # 如果没有arch包，使用EWMA作为替代
            return self._calculate_ewma_volatility()
    
    def _calculate_ewma_volatility(self, lambda_param: float = 0.94) -> Dict[str, any]:
        """
        计算EWMA波动率
        """
        
        # 指数加权移动平均波动率
        ewma_var = self.returns.ewm(alpha=1-lambda_param).var()
        ewma_vol = np.sqrt(ewma_var)
        
        return {
            'ewma_volatility': ewma_vol,
            'current_ewma_volatility': ewma_vol.iloc[-1],
            'forecasted_volatility': ewma_vol.iloc[-1]  # 假设持续当前水平
        }
    
    def calculate_extreme_value_metrics(self) -> Dict[str, float]:
        """
        计算极值理论指标
        """
        
        # 选择阈值（经验法则：选择5%的极端值）
        threshold = self.returns.quantile(0.05)
        extreme_returns = self.returns[self.returns <= threshold]
        
        if len(extreme_returns) < 10:
            return {'error': 'Not enough extreme observations'}
        
        # 超额分布
        excesses = threshold - extreme_returns
        
        # 拟合广义帕累托分布
        try:
            # 使用矩估计
            mean_excess = excesses.mean()
            var_excess = excesses.var()
            
            # 形状参数和尺度参数的矩估计
            xi = 0.5 * (mean_excess**2 / var_excess - 1)  # 形状参数
            beta = 0.5 * mean_excess * (mean_excess**2 / var_excess + 1)  # 尺度参数
            
            # EVT-VaR计算
            n = len(self.returns)
            nu = len(extreme_returns)
            
            evt_var_95 = threshold - (beta / xi) * ((n/nu * 0.05)**(-xi) - 1)
            evt_var_99 = threshold - (beta / xi) * ((n/nu * 0.01)**(-xi) - 1)
            
            return {
                'threshold': threshold,
                'shape_parameter': xi,
                'scale_parameter': beta,
                'evt_var_95': evt_var_95,
                'evt_var_99': evt_var_99,
                'extreme_observations': len(extreme_returns)
            }
            
        except:
            return {'error': 'EVT calculation failed'}
    
    def calculate_tail_risk_measures(self) -> Dict[str, float]:
        """
        计算尾部风险度量
        """
        
        measures = {}
        
        # 左尾系数（负偏度）
        measures['left_tail_coefficient'] = -self.returns.skew()
        
        # 尾部比率
        q95 = self.returns.quantile(0.95)
        q05 = self.returns.quantile(0.05)
        median = self.returns.median()
        
        if median != q05:
            measures['tail_ratio'] = (q95 - median) / (median - q05)
        else:
            measures['tail_ratio'] = np.inf
        
        # 尾部指数
        sorted_returns = self.returns.sort_values()
        tail_size = int(len(sorted_returns) * 0.05)
        
        if tail_size > 1:
            tail_returns = sorted_returns.iloc[:tail_size]
            log_ratios = []
            
            for i in range(1, len(tail_returns)):
                if tail_returns.iloc[i] != 0 and tail_returns.iloc[i-1] != 0:
                    log_ratios.append(np.log(abs(tail_returns.iloc[i]) / abs(tail_returns.iloc[i-1])))
            
            if log_ratios:
                measures['tail_index'] = 1 / np.mean(log_ratios) if np.mean(log_ratios) != 0 else np.inf
            else:
                measures['tail_index'] = np.inf
        else:
            measures['tail_index'] = np.inf
        
        return measures
```

## 第二节：VaR与CVaR深入分析

### 2.1 多种VaR计算方法

```python
class VaRCalculator:
    """
    VaR计算器 - 多种方法实现
    """
    
    def __init__(self, returns: pd.Series, positions: pd.Series = None):
        self.returns = returns.dropna()
        self.positions = positions
        self.portfolio_value = 1000000  # 默认组合价值
        
    def historical_var(self, confidence_level: float = 0.95, 
                      window_size: int = None) -> pd.Series:
        """
        历史模拟VaR
        """
        
        if window_size is None:
            window_size = len(self.returns)
        
        var_series = []
        
        for i in range(window_size, len(self.returns) + 1):
            window_returns = self.returns.iloc[i-window_size:i]
            var = window_returns.quantile(1 - confidence_level)
            var_series.append(var)
        
        return pd.Series(var_series, index=self.returns.index[window_size-1:])
    
    def parametric_var(self, confidence_level: float = 0.95, 
                      window_size: int = 252) -> pd.Series:
        """
        参数法VaR（滚动窗口）
        """
        
        var_series = []
        
        for i in range(window_size, len(self.returns) + 1):
            window_returns = self.returns.iloc[i-window_size:i]
            
            mean_return = window_returns.mean()
            std_return = window_returns.std()
            
            # 假设正态分布
            var = mean_return - std_return * stats.norm.ppf(confidence_level)
            var_series.append(var)
        
        return pd.Series(var_series, index=self.returns.index[window_size-1:])
    
    def cornish_fisher_var(self, confidence_level: float = 0.95, 
                          window_size: int = 252) -> pd.Series:
        """
        Cornish-Fisher修正VaR（考虑偏度和峰度）
        """
        
        var_series = []
        
        for i in range(window_size, len(self.returns) + 1):
            window_returns = self.returns.iloc[i-window_size:i]
            
            mean_return = window_returns.mean()
            std_return = window_returns.std()
            skewness = window_returns.skew()
            kurtosis = window_returns.kurtosis()
            
            # 正态分布分位数
            z = stats.norm.ppf(confidence_level)
            
            # Cornish-Fisher修正
            z_cf = (z + 
                   (z**2 - 1) * skewness / 6 +
                   (z**3 - 3*z) * kurtosis / 24 -
                   (2*z**3 - 5*z) * skewness**2 / 36)
            
            var = mean_return - std_return * z_cf
            var_series.append(var)
        
        return pd.Series(var_series, index=self.returns.index[window_size-1:])
    
    def filtered_historical_var(self, confidence_level: float = 0.95,
                               garch_model: any = None) -> pd.Series:
        """
        过滤历史模拟VaR（结合GARCH）
        """
        
        if garch_model is None:
            # 使用简单的EWMA
            vol_model = self.returns.ewm(alpha=0.06).std()
        else:
            vol_model = garch_model
        
        # 标准化收益率
        standardized_returns = self.returns / vol_model
        
        # 使用标准化收益率的分位数
        var_multiplier = standardized_returns.quantile(1 - confidence_level)
        
        # 重新缩放到当前波动率
        filtered_var = var_multiplier * vol_model
        
        return filtered_var
    
    def monte_carlo_var(self, confidence_level: float = 0.95, 
                       n_simulations: int = 10000,
                       distribution: str = 'normal') -> float:
        """
        蒙特卡洛VaR
        """
        
        np.random.seed(42)
        
        mean_return = self.returns.mean()
        std_return = self.returns.std()
        
        if distribution == 'normal':
            simulated_returns = np.random.normal(mean_return, std_return, n_simulations)
        
        elif distribution == 't':
            # 拟合t分布
            params = stats.t.fit(self.returns)
            df, loc, scale = params
            simulated_returns = stats.t.rvs(df, loc, scale, size=n_simulations)
        
        elif distribution == 'skewed_t':
            # 偏t分布（简化实现）
            skewness = self.returns.skew()
            simulated_returns = np.random.normal(mean_return, std_return, n_simulations)
            
            # 添加偏度调整
            simulated_returns += skewness * std_return * 0.1 * np.random.randn(n_simulations)
        
        var = np.percentile(simulated_returns, (1 - confidence_level) * 100)
        
        return var
    
    def component_var(self, weights: np.ndarray, cov_matrix: np.ndarray,
                     confidence_level: float = 0.95) -> Dict[str, any]:\n        \"\"\"\n        成分VaR（组合VaR分解）\n        \"\"\"\n        \n        # 组合收益率和波动率\n        portfolio_return = np.sum(weights * self.returns.mean())\n        portfolio_variance = np.dot(weights.T, np.dot(cov_matrix, weights))\n        portfolio_volatility = np.sqrt(portfolio_variance)\n        \n        # 组合VaR\n        z_score = stats.norm.ppf(confidence_level)\n        portfolio_var = -(portfolio_return - z_score * portfolio_volatility)\n        \n        # 边际VaR\n        marginal_var = -z_score * np.dot(cov_matrix, weights) / portfolio_volatility\n        \n        # 成分VaR\n        component_var = weights * marginal_var\n        \n        # 百分比贡献\n        var_contributions = component_var / portfolio_var\n        \n        return {\n            'portfolio_var': portfolio_var,\n            'marginal_var': marginal_var,\n            'component_var': component_var,\n            'var_contributions': var_contributions,\n            'portfolio_volatility': portfolio_volatility\n        }\n    \n    def incremental_var(self, base_weights: np.ndarray, new_weights: np.ndarray,\n                       cov_matrix: np.ndarray, confidence_level: float = 0.95) -> float:\n        \"\"\"\n        增量VaR\n        \"\"\"\n        \n        # 基础组合VaR\n        base_var_result = self.component_var(base_weights, cov_matrix, confidence_level)\n        base_var = base_var_result['portfolio_var']\n        \n        # 新组合VaR\n        new_var_result = self.component_var(new_weights, cov_matrix, confidence_level)\n        new_var = new_var_result['portfolio_var']\n        \n        # 增量VaR\n        incremental_var = new_var - base_var\n        \n        return incremental_var\n    \n    def backtesting_var(self, var_estimates: pd.Series, \n                       actual_returns: pd.Series,\n                       confidence_level: float = 0.95) -> Dict[str, any]:\n        \"\"\"\n        VaR回测\n        \"\"\"\n        \n        # 对齐数据\n        aligned_data = pd.DataFrame({\n            'var': var_estimates,\n            'returns': actual_returns\n        }).dropna()\n        \n        # VaR违背次数\n        violations = aligned_data['returns'] < aligned_data['var']\n        n_violations = violations.sum()\n        n_observations = len(aligned_data)\n        \n        # 预期违背次数\n        expected_violations = n_observations * (1 - confidence_level)\n        \n        # 违背率\n        violation_rate = n_violations / n_observations\n        expected_rate = 1 - confidence_level\n        \n        # Kupiec POF测试\n        if n_violations > 0 and n_violations < n_observations:\n            likelihood_ratio = -2 * np.log(\n                (expected_rate ** n_violations * (1 - expected_rate) ** (n_observations - n_violations)) /\n                (violation_rate ** n_violations * (1 - violation_rate) ** (n_observations - n_violations))\n            )\n            pof_p_value = 1 - stats.chi2.cdf(likelihood_ratio, df=1)\n        else:\n            likelihood_ratio = np.inf\n            pof_p_value = 0\n        \n        # 独立性测试（违背的聚类）\n        violation_runs = self._calculate_runs(violations)\n        \n        return {\n            'n_violations': n_violations,\n            'expected_violations': expected_violations,\n            'violation_rate': violation_rate,\n            'expected_rate': expected_rate,\n            'kupiec_lr': likelihood_ratio,\n            'kupiec_p_value': pof_p_value,\n            'violation_runs': violation_runs,\n            'violations_dates': aligned_data.index[violations].tolist()\n        }\n    \n    def _calculate_runs(self, violations: pd.Series) -> Dict[str, int]:\n        \"\"\"\n        计算违背的游程\n        \"\"\"\n        \n        runs = []\n        current_run = 0\n        \n        for violation in violations:\n            if violation:\n                current_run += 1\n            else:\n                if current_run > 0:\n                    runs.append(current_run)\n                    current_run = 0\n        \n        if current_run > 0:\n            runs.append(current_run)\n        \n        return {\n            'total_runs': len(runs),\n            'max_run_length': max(runs) if runs else 0,\n            'avg_run_length': np.mean(runs) if runs else 0\n        }\n```\n\n### 2.2 条件VaR (CVaR) 高级应用\n\n```python\nclass CVaROptimizer:\n    \"\"\"\n    CVaR优化器\n    \"\"\"\n    \n    def __init__(self, returns_matrix: pd.DataFrame, confidence_level: float = 0.95):\n        self.returns_matrix = returns_matrix.dropna()\n        self.confidence_level = confidence_level\n        self.alpha = 1 - confidence_level\n        \n    def optimize_cvar_portfolio(self, target_return: float = None,\n                               max_weight: float = 1.0,\n                               min_weight: float = 0.0) -> Dict[str, any]:\n        \"\"\"\n        CVaR优化组合\n        \"\"\"\n        \n        from scipy.optimize import minimize\n        \n        n_assets = len(self.returns_matrix.columns)\n        n_scenarios = len(self.returns_matrix)\n        \n        # 定义优化变量: [权重, VaR, 辅助变量]\n        # 辅助变量用于CVaR计算\n        \n        def objective(x):\n            weights = x[:n_assets]\n            var = x[n_assets]\n            u = x[n_assets+1:n_assets+1+n_scenarios]\n            \n            # CVaR = VaR + (1/alpha) * E[max(0, -R - VaR)]\n            cvar = var + (1 / self.alpha) * np.mean(u)\n            return cvar\n        \n        def constraint_portfolio_return(x):\n            weights = x[:n_assets]\n            portfolio_returns = self.returns_matrix @ weights\n            return portfolio_returns.mean() - target_return if target_return else 0\n        \n        def constraint_weights_sum(x):\n            weights = x[:n_assets]\n            return np.sum(weights) - 1\n        \n        def constraint_auxiliary_variables(x):\n            weights = x[:n_assets]\n            var = x[n_assets]\n            u = x[n_assets+1:n_assets+1+n_scenarios]\n            \n            portfolio_returns = self.returns_matrix @ weights\n            # u >= max(0, -R - VaR)\n            return u - np.maximum(0, -portfolio_returns - var)\n        \n        # 约束条件\n        constraints = [\n            {'type': 'eq', 'fun': constraint_weights_sum}\n        ]\n        \n        if target_return is not None:\n            constraints.append({\n                'type': 'eq', 'fun': constraint_portfolio_return\n            })\n        \n        # 辅助变量约束\n        for i in range(n_scenarios):\n            constraints.append({\n                'type': 'ineq', \n                'fun': lambda x, idx=i: (x[n_assets+1+idx] - \n                                       max(0, -(self.returns_matrix.iloc[idx] @ x[:n_assets]) - x[n_assets]))\n            })\n        \n        # 边界条件\n        bounds = []\n        # 权重边界\n        for i in range(n_assets):\n            bounds.append((min_weight, max_weight))\n        \n        # VaR边界（无限制）\n        bounds.append((-np.inf, np.inf))\n        \n        # 辅助变量边界（非负）\n        for i in range(n_scenarios):\n            bounds.append((0, np.inf))\n        \n        # 初始值\n        x0 = np.concatenate([\n            np.ones(n_assets) / n_assets,  # 等权重\n            [0],  # VaR初始值\n            np.zeros(n_scenarios)  # 辅助变量初始值\n        ])\n        \n        # 优化\n        result = minimize(objective, x0, method='SLSQP', bounds=bounds, constraints=constraints)\n        \n        if result.success:\n            optimal_weights = result.x[:n_assets]\n            optimal_var = result.x[n_assets]\n            optimal_cvar = result.fun\n            \n            # 计算组合统计\n            portfolio_returns = self.returns_matrix @ optimal_weights\n            portfolio_mean = portfolio_returns.mean()\n            portfolio_std = portfolio_returns.std()\n            \n            return {\n                'success': True,\n                'weights': optimal_weights,\n                'var': optimal_var,\n                'cvar': optimal_cvar,\n                'expected_return': portfolio_mean,\n                'volatility': portfolio_std,\n                'portfolio_returns': portfolio_returns\n            }\n        else:\n            return {'success': False, 'message': result.message}\n    \n    def cvar_efficient_frontier(self, n_points: int = 20) -> pd.DataFrame:\n        \"\"\"\n        生成CVaR有效前沿\n        \"\"\"\n        \n        # 确定收益率范围\n        individual_returns = self.returns_matrix.mean()\n        min_return = individual_returns.min()\n        max_return = individual_returns.max()\n        \n        target_returns = np.linspace(min_return, max_return, n_points)\n        \n        efficient_portfolios = []\n        \n        for target_return in target_returns:\n            result = self.optimize_cvar_portfolio(target_return=target_return)\n            \n            if result['success']:\n                efficient_portfolios.append({\n                    'target_return': target_return,\n                    'expected_return': result['expected_return'],\n                    'cvar': result['cvar'],\n                    'var': result['var'],\n                    'volatility': result['volatility'],\n                    'weights': result['weights']\n                })\n        \n        return pd.DataFrame(efficient_portfolios)\n    \n    def compare_mean_variance_cvar(self, target_return: float) -> Dict[str, any]:\n        \"\"\"\n        比较均值-方差和CVaR优化\n        \"\"\"\n        \n        # CVaR优化\n        cvar_result = self.optimize_cvar_portfolio(target_return=target_return)\n        \n        # 均值-方差优化\n        mv_result = self._mean_variance_optimization(target_return)\n        \n        # 计算两种组合的风险指标\n        comparison = {\n            'cvar_portfolio': self._calculate_risk_metrics(cvar_result['weights']),\n            'mv_portfolio': self._calculate_risk_metrics(mv_result['weights'])\n        }\n        \n        return comparison\n    \n    def _mean_variance_optimization(self, target_return: float) -> Dict[str, any]:\n        \"\"\"\n        均值-方差优化（用于比较）\n        \"\"\"\n        \n        n_assets = len(self.returns_matrix.columns)\n        \n        # 期望收益和协方差矩阵\n        expected_returns = self.returns_matrix.mean().values\n        cov_matrix = self.returns_matrix.cov().values\n        \n        def objective(weights):\n            return 0.5 * np.dot(weights.T, np.dot(cov_matrix, weights))\n        \n        def constraint_return(weights):\n            return np.dot(weights, expected_returns) - target_return\n        \n        def constraint_weights(weights):\n            return np.sum(weights) - 1\n        \n        constraints = [\n            {'type': 'eq', 'fun': constraint_return},\n            {'type': 'eq', 'fun': constraint_weights}\n        ]\n        \n        bounds = [(0, 1) for _ in range(n_assets)]\n        x0 = np.ones(n_assets) / n_assets\n        \n        result = minimize(objective, x0, method='SLSQP', bounds=bounds, constraints=constraints)\n        \n        return {\n            'success': result.success,\n            'weights': result.x if result.success else None\n        }\n    \n    def _calculate_risk_metrics(self, weights: np.ndarray) -> Dict[str, float]:\n        \"\"\"\n        计算给定权重的风险指标\n        \"\"\"\n        \n        portfolio_returns = self.returns_matrix @ weights\n        \n        # 基础统计\n        metrics = {\n            'expected_return': portfolio_returns.mean(),\n            'volatility': portfolio_returns.std(),\n            'skewness': portfolio_returns.skew(),\n            'kurtosis': portfolio_returns.kurtosis()\n        }\n        \n        # VaR和CVaR\n        var_95 = portfolio_returns.quantile(0.05)\n        cvar_95 = portfolio_returns[portfolio_returns <= var_95].mean()\n        \n        metrics.update({\n            'var_95': var_95,\n            'cvar_95': cvar_95,\n            'max_loss': portfolio_returns.min()\n        })\n        \n        return metrics\n```\n\n## 第三节：回撤控制策略\n\n### 3.1 动态风险控制\n\n```python\nclass DrawdownController:\n    \"\"\"\n    回撤控制器\n    \"\"\"\n    \n    def __init__(self, max_drawdown_threshold: float = 0.1,\n                 lookback_window: int = 252):\n        self.max_drawdown_threshold = max_drawdown_threshold\n        self.lookback_window = lookback_window\n        self.position_history = []\n        self.equity_curve = []\n        \n    def calculate_position_size(self, signal_strength: float, \n                              current_equity: float,\n                              historical_returns: pd.Series,\n                              volatility_target: float = 0.15) -> float:\n        \"\"\"\n        基于风险的仓位规模计算\n        \"\"\"\n        \n        # 波动率调整仓位\n        current_vol = historical_returns.rolling(self.lookback_window).std().iloc[-1] * np.sqrt(252)\n        vol_adjustment = volatility_target / current_vol if current_vol > 0 else 1\n        \n        # 基础仓位\n        base_position = signal_strength * vol_adjustment\n        \n        # 回撤调整\n        drawdown_adjustment = self._calculate_drawdown_adjustment(\n            current_equity, historical_returns\n        )\n        \n        # 最终仓位\n        final_position = base_position * drawdown_adjustment\n        \n        # 限制仓位范围\n        final_position = np.clip(final_position, -1.0, 1.0)\n        \n        return final_position\n    \n    def _calculate_drawdown_adjustment(self, current_equity: float,\n                                     historical_returns: pd.Series) -> float:\n        \"\"\"\n        计算回撤调整因子\n        \"\"\"\n        \n        if len(self.equity_curve) < 10:\n            return 1.0\n        \n        # 计算当前回撤\n        equity_series = pd.Series(self.equity_curve + [current_equity])\n        running_max = equity_series.expanding().max()\n        current_drawdown = (equity_series.iloc[-1] - running_max.iloc[-1]) / running_max.iloc[-1]\n        \n        # 回撤调整函数\n        if current_drawdown >= 0:\n            # 没有回撤，正常仓位\n            return 1.0\n        elif abs(current_drawdown) < self.max_drawdown_threshold * 0.5:\n            # 小幅回撤，轻微减仓\n            return 0.8\n        elif abs(current_drawdown) < self.max_drawdown_threshold:\n            # 中等回撤，显著减仓\n            return 0.5\n        else:\n            # 大幅回撤，停止交易\n            return 0.0\n    \n    def implement_stop_loss(self, entry_price: float, current_price: float,\n                           position_type: str, stop_loss_pct: float = 0.05) -> bool:\n        \"\"\"\n        实施止损\n        \"\"\"\n        \n        if position_type == 'long':\n            loss = (entry_price - current_price) / entry_price\n            return loss >= stop_loss_pct\n        elif position_type == 'short':\n            loss = (current_price - entry_price) / entry_price\n            return loss >= stop_loss_pct\n        \n        return False\n    \n    def implement_time_stop(self, entry_date: pd.Timestamp, \n                           current_date: pd.Timestamp,\n                           max_holding_days: int = 30) -> bool:\n        \"\"\"\n        实施时间止损\n        \"\"\"\n        \n        holding_days = (current_date - entry_date).days\n        return holding_days >= max_holding_days\n    \n    def calculate_kelly_criterion(self, win_rate: float, avg_win: float, avg_loss: float) -> float:\n        \"\"\"\n        计算凯利公式仓位\n        \"\"\"\n        \n        if avg_loss == 0:\n            return 0\n        \n        # 凯利公式: f = (bp - q) / b\n        # 其中 b = avg_win/avg_loss, p = win_rate, q = 1 - win_rate\n        b = avg_win / abs(avg_loss)\n        p = win_rate\n        q = 1 - win_rate\n        \n        kelly_fraction = (b * p - q) / b\n        \n        # 限制凯利比例（通常使用1/4或1/2凯利）\n        conservative_kelly = kelly_fraction * 0.25\n        \n        return max(0, min(conservative_kelly, 0.2))  # 最大20%仓位\n    \n    def adaptive_volatility_targeting(self, returns: pd.Series, \n                                    target_vol: float = 0.15,\n                                    rebalance_freq: int = 5) -> pd.Series:\n        \"\"\"\n        自适应波动率目标\n        \"\"\"\n        \n        position_sizes = []\n        \n        for i in range(len(returns)):\n            if i < 20:  # 需要足够的历史数据\n                position_sizes.append(1.0)\n                continue\n            \n            # 计算历史波动率\n            historical_vol = returns.iloc[max(0, i-252):i].std() * np.sqrt(252)\n            \n            # 波动率调整\n            if historical_vol > 0:\n                vol_adjustment = target_vol / historical_vol\n                # 限制调整幅度\n                vol_adjustment = np.clip(vol_adjustment, 0.2, 2.0)\n            else:\n                vol_adjustment = 1.0\n            \n            position_sizes.append(vol_adjustment)\n        \n        return pd.Series(position_sizes, index=returns.index)\n    \n    def regime_aware_sizing(self, returns: pd.Series, regime_labels: pd.Series) -> pd.Series:\n        \"\"\"\n        状态感知仓位调整\n        \"\"\"\n        \n        # 不同状态下的仓位乘数\n        regime_multipliers = {\n            0: 0.5,  # 熊市/高波动状态\n            1: 1.0,  # 正常状态\n            2: 1.2   # 牛市/低波动状态\n        }\n        \n        position_multipliers = regime_labels.map(regime_multipliers).fillna(1.0)\n        \n        return position_multipliers\n\nclass RiskBudgetSystem:\n    \"\"\"\n    风险预算系统\n    \"\"\"\n    \n    def __init__(self, total_risk_budget: float = 0.02):\n        self.total_risk_budget = total_risk_budget  # 日度风险预算\n        self.strategy_allocations = {}\n        \n    def allocate_risk_budget(self, strategies: Dict[str, Dict]) -> Dict[str, float]:\n        \"\"\"\n        分配风险预算\n        \n        strategies: {\n            'strategy_name': {\n                'expected_return': float,\n                'volatility': float,\n                'sharpe_ratio': float,\n                'max_drawdown': float\n            }\n        }\n        \"\"\"\n        \n        # 基于夏普比率的风险预算分配\n        sharpe_ratios = {name: info['sharpe_ratio'] for name, info in strategies.items()}\n        \n        # 处理负夏普比率\n        min_sharpe = min(sharpe_ratios.values())\n        if min_sharpe < 0:\n            adjusted_sharpe = {name: sr - min_sharpe + 0.1 for name, sr in sharpe_ratios.items()}\n        else:\n            adjusted_sharpe = sharpe_ratios\n        \n        # 权重分配\n        total_adjusted_sharpe = sum(adjusted_sharpe.values())\n        \n        if total_adjusted_sharpe > 0:\n            risk_weights = {name: sr / total_adjusted_sharpe for name, sr in adjusted_sharpe.items()}\n        else:\n            # 等权重分配\n            risk_weights = {name: 1/len(strategies) for name in strategies.keys()}\n        \n        # 计算每个策略的风险预算\n        risk_budgets = {}\n        for name, weight in risk_weights.items():\n            strategy_vol = strategies[name]['volatility']\n            \n            # 每日风险预算\n            daily_vol_budget = self.total_risk_budget * weight\n            \n            # 转换为仓位限制\n            if strategy_vol > 0:\n                max_position = daily_vol_budget / (strategy_vol / np.sqrt(252))\n            else:\n                max_position = 0\n            \n            risk_budgets[name] = {\n                'risk_weight': weight,\n                'daily_vol_budget': daily_vol_budget,\n                'max_position': min(max_position, 1.0)  # 限制最大仓位\n            }\n        \n        return risk_budgets\n    \n    def monitor_risk_utilization(self, strategy_positions: Dict[str, float],\n                               strategy_returns: Dict[str, pd.Series]) -> Dict[str, float]:\n        \"\"\"\n        监控风险利用率\n        \"\"\"\n        \n        utilization = {}\n        \n        for strategy_name, position in strategy_positions.items():\n            if strategy_name in strategy_returns:\n                returns = strategy_returns[strategy_name]\n                \n                # 计算当前策略的风险贡献\n                strategy_vol = returns.std() * np.sqrt(252)\n                current_risk = abs(position) * strategy_vol / np.sqrt(252)\n                \n                # 计算利用率\n                if strategy_name in self.strategy_allocations:\n                    budget = self.strategy_allocations[strategy_name]['daily_vol_budget']\n                    utilization[strategy_name] = current_risk / budget if budget > 0 else 0\n                else:\n                    utilization[strategy_name] = 0\n        \n        return utilization\n    \n    def rebalance_risk_budget(self, current_utilization: Dict[str, float],\n                            target_utilization: float = 0.8) -> Dict[str, float]:\n        \"\"\"\n        重新平衡风险预算\n        \"\"\"\n        \n        adjustments = {}\n        \n        for strategy_name, utilization in current_utilization.items():\n            if utilization > target_utilization:\n                # 超出目标，减少风险分配\n                adjustment_factor = target_utilization / utilization\n            elif utilization < target_utilization * 0.5:\n                # 利用不足，可以增加风险分配\n                adjustment_factor = min(target_utilization / utilization, 1.2)\n            else:\n                # 在合理范围内，不调整\n                adjustment_factor = 1.0\n            \n            adjustments[strategy_name] = adjustment_factor\n        \n        return adjustments\n```\n\n### 3.2 动态对冲策略\n\n```python\nclass DynamicHedgingSystem:\n    \"\"\"\n    动态对冲系统\n    \"\"\"\n    \n    def __init__(self, hedge_threshold: float = 0.05):\n        self.hedge_threshold = hedge_threshold\n        self.hedge_positions = {}\n        \n    def calculate_portfolio_beta(self, portfolio_returns: pd.Series,\n                               market_returns: pd.Series,\n                               window: int = 252) -> pd.Series:\n        \"\"\"\n        计算滚动组合Beta\n        \"\"\"\n        \n        aligned_data = pd.DataFrame({\n            'portfolio': portfolio_returns,\n            'market': market_returns\n        }).dropna()\n        \n        rolling_beta = []\n        \n        for i in range(window, len(aligned_data)):\n            window_data = aligned_data.iloc[i-window:i]\n            \n            # 计算beta\n            covariance = window_data['portfolio'].cov(window_data['market'])\n            market_variance = window_data['market'].var()\n            \n            if market_variance > 0:\n                beta = covariance / market_variance\n            else:\n                beta = 1.0\n            \n            rolling_beta.append(beta)\n        \n        return pd.Series(rolling_beta, index=aligned_data.index[window:])\n    \n    def calculate_hedge_ratio(self, portfolio_value: float,\n                            current_beta: float,\n                            target_beta: float = 0.0) -> float:\n        \"\"\"\n        计算对冲比率\n        \"\"\"\n        \n        # 需要对冲的beta暴露\n        beta_exposure = current_beta - target_beta\n        \n        # 对冲比率（以组合价值的百分比表示）\n        hedge_ratio = -beta_exposure\n        \n        return hedge_ratio\n    \n    def implement_beta_hedge(self, portfolio_value: float,\n                           current_beta: float,\n                           market_price: float,\n                           futures_multiplier: float = 250) -> Dict[str, any]:\n        \"\"\"\n        实施Beta对冲\n        \"\"\"\n        \n        hedge_ratio = self.calculate_hedge_ratio(portfolio_value, current_beta)\n        \n        # 计算需要的期货合约数量\n        futures_value = market_price * futures_multiplier\n        contracts_needed = (hedge_ratio * portfolio_value) / futures_value\n        \n        # 四舍五入到整数合约\n        contracts_to_trade = round(contracts_needed)\n        \n        # 实际对冲比率\n        actual_hedge_ratio = (contracts_to_trade * futures_value) / portfolio_value\n        \n        return {\n            'target_hedge_ratio': hedge_ratio,\n            'contracts_needed': contracts_needed,\n            'contracts_to_trade': contracts_to_trade,\n            'actual_hedge_ratio': actual_hedge_ratio,\n            'hedge_value': contracts_to_trade * futures_value\n        }\n    \n    def pairs_hedge(self, primary_returns: pd.Series,\n                   hedge_returns: pd.Series,\n                   window: int = 60) -> pd.Series:\n        \"\"\"\n        配对对冲\n        \"\"\"\n        \n        aligned_data = pd.DataFrame({\n            'primary': primary_returns,\n            'hedge': hedge_returns\n        }).dropna()\n        \n        hedge_ratios = []\n        \n        for i in range(window, len(aligned_data)):\n            window_data = aligned_data.iloc[i-window:i]\n            \n            # 回归计算对冲比率\n            X = window_data['hedge'].values.reshape(-1, 1)\n            y = window_data['primary'].values\n            \n            # 简单线性回归\n            coef = np.cov(y, window_data['hedge'])[0, 1] / np.var(window_data['hedge'])\n            hedge_ratios.append(-coef)  # 负号表示对冲\n        \n        return pd.Series(hedge_ratios, index=aligned_data.index[window:])\n    \n    def volatility_hedge(self, portfolio_returns: pd.Series,\n                        vix_data: pd.Series,\n                        target_vol: float = 0.15) -> pd.Series:\n        \"\"\"\n        波动率对冲\n        \"\"\"\n        \n        # 计算滚动波动率\n        rolling_vol = portfolio_returns.rolling(20).std() * np.sqrt(252)\n        \n        # 波动率超出目标时的对冲信号\n        vol_excess = rolling_vol - target_vol\n        \n        # VIX对冲比率（简化模型）\n        hedge_signals = []\n        \n        for i, (date, vol_ex) in enumerate(vol_excess.items()):\n            if abs(vol_ex) > self.hedge_threshold:\n                # 需要对冲\n                if vol_ex > 0:  # 波动率过高\n                    hedge_signal = min(vol_ex / target_vol, 0.5)  # 买入VIX\n                else:  # 波动率过低\n                    hedge_signal = max(vol_ex / target_vol, -0.5)  # 卖出VIX\n            else:\n                hedge_signal = 0\n            \n            hedge_signals.append(hedge_signal)\n        \n        return pd.Series(hedge_signals, index=vol_excess.index)\n    \n    def tail_risk_hedge(self, portfolio_returns: pd.Series,\n                       put_prices: pd.Series,\n                       hedge_trigger: float = -0.02) -> pd.Series:\n        \"\"\"\n        尾部风险对冲\n        \"\"\"\n        \n        hedge_signals = []\n        \n        for i, (date, ret) in enumerate(portfolio_returns.items()):\n            # 计算最近的最大损失\n            recent_returns = portfolio_returns.iloc[max(0, i-20):i+1]\n            recent_min = recent_returns.min()\n            \n            # 如果接近触发阈值，增加对冲\n            if recent_min < hedge_trigger:\n                # 计算对冲强度\n                hedge_intensity = min(abs(recent_min) / abs(hedge_trigger), 1.0)\n                hedge_signals.append(hedge_intensity)\n            else:\n                hedge_signals.append(0)\n        \n        return pd.Series(hedge_signals, index=portfolio_returns.index)\n```\n\n## 第四节：绩效评估框架\n\n### 4.1 综合绩效指标\n\n```python\nclass PerformanceEvaluator:\n    \"\"\"\n    绩效评估器\n    \"\"\"\n    \n    def __init__(self, portfolio_returns: pd.Series, \n                 benchmark_returns: pd.Series = None,\n                 risk_free_rate: float = 0.02):\n        self.portfolio_returns = portfolio_returns.dropna()\n        self.benchmark_returns = benchmark_returns\n        self.risk_free_rate = risk_free_rate\n        self.trading_days = 252\n        \n    def calculate_comprehensive_metrics(self) -> Dict[str, float]:\n        \"\"\"\n        计算综合绩效指标\n        \"\"\"\n        \n        metrics = {}\n        \n        # 基础收益指标\n        total_return = (1 + self.portfolio_returns).prod() - 1\n        annualized_return = (1 + self.portfolio_returns.mean()) ** self.trading_days - 1\n        \n        metrics['total_return'] = total_return\n        metrics['annualized_return'] = annualized_return\n        \n        # 风险指标\n        volatility = self.portfolio_returns.std() * np.sqrt(self.trading_days)\n        metrics['volatility'] = volatility\n        \n        # 风险调整收益指标\n        if volatility > 0:\n            metrics['sharpe_ratio'] = (annualized_return - self.risk_free_rate) / volatility\n        else:\n            metrics['sharpe_ratio'] = 0\n        \n        # 回撤指标\n        cumulative_returns = (1 + self.portfolio_returns).cumprod()\n        running_max = cumulative_returns.expanding().max()\n        drawdown = (cumulative_returns - running_max) / running_max\n        \n        metrics['max_drawdown'] = drawdown.min()\n        \n        if abs(metrics['max_drawdown']) > 0:\n            metrics['calmar_ratio'] = annualized_return / abs(metrics['max_drawdown'])\n        else:\n            metrics['calmar_ratio'] = np.inf\n        \n        # 下行风险指标\n        downside_returns = self.portfolio_returns[self.portfolio_returns < 0]\n        if len(downside_returns) > 0:\n            downside_deviation = downside_returns.std() * np.sqrt(self.trading_days)\n            metrics['downside_deviation'] = downside_deviation\n            \n            if downside_deviation > 0:\n                metrics['sortino_ratio'] = (annualized_return - self.risk_free_rate) / downside_deviation\n            else:\n                metrics['sortino_ratio'] = np.inf\n        else:\n            metrics['downside_deviation'] = 0\n            metrics['sortino_ratio'] = np.inf\n        \n        # 分布指标\n        metrics['skewness'] = self.portfolio_returns.skew()\n        metrics['kurtosis'] = self.portfolio_returns.kurtosis()\n        \n        # 胜率指标\n        metrics['win_rate'] = (self.portfolio_returns > 0).mean()\n        \n        winning_returns = self.portfolio_returns[self.portfolio_returns > 0]\n        losing_returns = self.portfolio_returns[self.portfolio_returns < 0]\n        \n        if len(winning_returns) > 0 and len(losing_returns) > 0:\n            metrics['profit_loss_ratio'] = winning_returns.mean() / abs(losing_returns.mean())\n        else:\n            metrics['profit_loss_ratio'] = np.inf\n        \n        # VaR和CVaR\n        var_95 = self.portfolio_returns.quantile(0.05)\n        cvar_95 = self.portfolio_returns[self.portfolio_returns <= var_95].mean()\n        \n        metrics['var_95'] = var_95\n        metrics['cvar_95'] = cvar_95\n        \n        # 基准相关指标\n        if self.benchmark_returns is not None:\n            benchmark_metrics = self._calculate_benchmark_metrics()\n            metrics.update(benchmark_metrics)\n        \n        return metrics\n    \n    def _calculate_benchmark_metrics(self) -> Dict[str, float]:\n        \"\"\"\n        计算相对基准的指标\n        \"\"\"\n        \n        # 对齐数据\n        aligned_data = pd.DataFrame({\n            'portfolio': self.portfolio_returns,\n            'benchmark': self.benchmark_returns\n        }).dropna()\n        \n        if len(aligned_data) == 0:\n            return {}\n        \n        portfolio_ret = aligned_data['portfolio']\n        benchmark_ret = aligned_data['benchmark']\n        \n        metrics = {}\n        \n        # 超额收益\n        excess_returns = portfolio_ret - benchmark_ret\n        metrics['excess_return'] = excess_returns.mean() * self.trading_days\n        \n        # 跟踪误差\n        metrics['tracking_error'] = excess_returns.std() * np.sqrt(self.trading_days)\n        \n        # 信息比率\n        if metrics['tracking_error'] > 0:\n            metrics['information_ratio'] = metrics['excess_return'] / metrics['tracking_error']\n        else:\n            metrics['information_ratio'] = 0\n        \n        # Beta和Alpha\n        covariance = portfolio_ret.cov(benchmark_ret)\n        benchmark_variance = benchmark_ret.var()\n        \n        if benchmark_variance > 0:\n            metrics['beta'] = covariance / benchmark_variance\n        else:\n            metrics['beta'] = 1\n        \n        portfolio_annual = (1 + portfolio_ret.mean()) ** self.trading_days - 1\n        benchmark_annual = (1 + benchmark_ret.mean()) ** self.trading_days - 1\n        \n        metrics['alpha'] = portfolio_annual - (self.risk_free_rate + metrics['beta'] * (benchmark_annual - self.risk_free_rate))\n        \n        # 相关系数\n        metrics['correlation'] = portfolio_ret.corr(benchmark_ret)\n        \n        # 上行/下行捕获率\n        up_periods = benchmark_ret > 0\n        down_periods = benchmark_ret < 0\n        \n        if up_periods.sum() > 0:\n            metrics['upside_capture'] = portfolio_ret[up_periods].mean() / benchmark_ret[up_periods].mean()\n        else:\n            metrics['upside_capture'] = 1\n            \n        if down_periods.sum() > 0:\n            metrics['downside_capture'] = portfolio_ret[down_periods].mean() / benchmark_ret[down_periods].mean()\n        else:\n            metrics['downside_capture'] = 1\n        \n        return metrics\n    \n    def rolling_performance_analysis(self, window: int = 252) -> pd.DataFrame:\n        \"\"\"\n        滚动绩效分析\n        \"\"\"\n        \n        rolling_metrics = []\n        \n        for i in range(window, len(self.portfolio_returns) + 1):\n            window_returns = self.portfolio_returns.iloc[i-window:i]\n            \n            # 计算窗口内指标\n            total_ret = (1 + window_returns).prod() - 1\n            vol = window_returns.std() * np.sqrt(self.trading_days)\n            \n            if vol > 0:\n                sharpe = (window_returns.mean() * self.trading_days - self.risk_free_rate) / vol\n            else:\n                sharpe = 0\n            \n            # 回撤\n            cum_ret = (1 + window_returns).cumprod()\n            running_max = cum_ret.expanding().max()\n            dd = ((cum_ret - running_max) / running_max).min()\n            \n            rolling_metrics.append({\n                'date': self.portfolio_returns.index[i-1],\n                'total_return': total_ret,\n                'volatility': vol,\n                'sharpe_ratio': sharpe,\n                'max_drawdown': dd\n            })\n        \n        return pd.DataFrame(rolling_metrics).set_index('date')\n    \n    def generate_performance_report(self) -> str:\n        \"\"\"\n        生成绩效报告\n        \"\"\"\n        \n        metrics = self.calculate_comprehensive_metrics()\n        \n        report = \"\"\"\n=== 投资组合绩效报告 ===\n\n收益指标:\n  总收益率: {total_return:.2%}\n  年化收益率: {annualized_return:.2%}\n\n风险指标:\n  年化波动率: {volatility:.2%}\n  最大回撤: {max_drawdown:.2%}\n  下行偏差: {downside_deviation:.2%}\n\n风险调整收益:\n  夏普比率: {sharpe_ratio:.3f}\n  Calmar比率: {calmar_ratio:.3f}\n  Sortino比率: {sortino_ratio:.3f}\n\n分布特征:\n  偏度: {skewness:.3f}\n  峰度: {kurtosis:.3f}\n\n交易统计:\n  胜率: {win_rate:.2%}\n  盈亏比: {profit_loss_ratio:.3f}\n\n风险度量:\n  VaR(95%): {var_95:.2%}\n  CVaR(95%): {cvar_95:.2%}\n        \"\"\".format(**metrics)\n        \n        if self.benchmark_returns is not None:\n            benchmark_section = \"\"\"\n相对基准指标:\n  超额收益: {excess_return:.2%}\n  跟踪误差: {tracking_error:.2%}\n  信息比率: {information_ratio:.3f}\n  Beta: {beta:.3f}\n  Alpha: {alpha:.2%}\n  相关系数: {correlation:.3f}\n  上行捕获率: {upside_capture:.3f}\n  下行捕获率: {downside_capture:.3f}\n            \"\"\".format(**{k: v for k, v in metrics.items() if k in [\n                'excess_return', 'tracking_error', 'information_ratio', \n                'beta', 'alpha', 'correlation', 'upside_capture', 'downside_capture'\n            ]})\n            \n            report += benchmark_section\n        \n        return report\n```\n\n### 4.2 绩效归因分析\n\n```python\nclass PerformanceAttribution:\n    \"\"\"\n    绩效归因分析\n    \"\"\"\n    \n    def __init__(self, portfolio_returns: pd.Series,\n                 factor_returns: pd.DataFrame,\n                 weights: pd.DataFrame = None):\n        self.portfolio_returns = portfolio_returns\n        self.factor_returns = factor_returns\n        self.weights = weights\n        \n    def factor_attribution(self, method: str = 'regression') -> Dict[str, any]:\n        \"\"\"\n        因子归因分析\n        \"\"\"\n        \n        # 对齐数据\n        aligned_data = pd.concat([\n            self.portfolio_returns.to_frame('portfolio'),\n            self.factor_returns\n        ], axis=1).dropna()\n        \n        portfolio_ret = aligned_data['portfolio']\n        factor_ret = aligned_data[self.factor_returns.columns]\n        \n        if method == 'regression':\n            # 回归归因\n            from sklearn.linear_model import LinearRegression\n            \n            model = LinearRegression()\n            model.fit(factor_ret, portfolio_ret)\n            \n            # 因子暴露度\n            factor_exposures = pd.Series(model.coef_, index=factor_ret.columns)\n            \n            # 因子贡献\n            factor_contributions = factor_exposures * factor_ret.mean() * 252\n            \n            # 特异性收益\n            predicted_returns = model.predict(factor_ret)\n            residual_returns = portfolio_ret - predicted_returns\n            idiosyncratic_return = residual_returns.mean() * 252\n            \n            # 总归因\n            total_attributed = factor_contributions.sum() + idiosyncratic_return\n            actual_return = portfolio_ret.mean() * 252\n            \n            return {\n                'factor_exposures': factor_exposures,\n                'factor_contributions': factor_contributions,\n                'idiosyncratic_return': idiosyncratic_return,\n                'total_attributed': total_attributed,\n                'actual_return': actual_return,\n                'r_squared': model.score(factor_ret, portfolio_ret)\n            }\n        \n        elif method == 'decomposition':\n            # 方差分解归因\n            portfolio_var = portfolio_ret.var()\n            \n            variance_attribution = {}\n            \n            for factor in factor_ret.columns:\n                # 单因子回归\n                factor_data = factor_ret[factor].values.reshape(-1, 1)\n                model = LinearRegression()\n                model.fit(factor_data, portfolio_ret)\n                \n                # 该因子解释的方差\n                factor_r2 = model.score(factor_data, portfolio_ret)\n                variance_contribution = factor_r2 * portfolio_var\n                \n                variance_attribution[factor] = {\n                    'variance_contribution': variance_contribution,\n                    'volatility_contribution': np.sqrt(variance_contribution),\n                    'r_squared': factor_r2\n                }\n            \n            return variance_attribution\n    \n    def sector_attribution(self, sector_returns: pd.DataFrame,\n                          benchmark_weights: pd.DataFrame,\n                          portfolio_weights: pd.DataFrame) -> pd.DataFrame:\n        \"\"\"\n        行业归因分析（Brinson模型）\n        \"\"\"\n        \n        # 确保数据对齐\n        common_dates = sector_returns.index.intersection(\n            benchmark_weights.index\n        ).intersection(portfolio_weights.index)\n        \n        sector_ret = sector_returns.loc[common_dates]\n        bench_weights = benchmark_weights.loc[common_dates]\n        port_weights = portfolio_weights.loc[common_dates]\n        \n        attribution_results = []\n        \n        for date in common_dates:\n            date_sector_ret = sector_ret.loc[date]\n            date_bench_weights = bench_weights.loc[date]\n            date_port_weights = port_weights.loc[date]\n            \n            # 基准收益率\n            benchmark_return = (date_bench_weights * date_sector_ret).sum()\n            \n            # 组合收益率\n            portfolio_return = (date_port_weights * date_sector_ret).sum()\n            \n            # 归因分解\n            sector_attribution = []\n            \n            for sector in date_sector_ret.index:\n                # 配置效应\n                allocation_effect = (\n                    (date_port_weights[sector] - date_bench_weights[sector]) * \n                    (date_sector_ret[sector] - benchmark_return)\n                )\n                \n                # 选择效应（这里简化，假设没有个股选择）\n                selection_effect = 0\n                \n                # 交互效应\n                interaction_effect = (\n                    (date_port_weights[sector] - date_bench_weights[sector]) * \n                    (date_sector_ret[sector] - benchmark_return)\n                ) * 0  # 简化为0\n                \n                sector_attribution.append({\n                    'date': date,\n                    'sector': sector,\n                    'allocation_effect': allocation_effect,\n                    'selection_effect': selection_effect,\n                    'interaction_effect': interaction_effect,\n                    'total_effect': allocation_effect + selection_effect + interaction_effect\n                })\n            \n            attribution_results.extend(sector_attribution)\n        \n        return pd.DataFrame(attribution_results)\n    \n    def style_attribution(self, style_factors: pd.DataFrame) -> Dict[str, float]:\n        \"\"\"\n        风格归因分析\n        \"\"\"\n        \n        # 常见风格因子：价值、成长、质量、动量、波动率、规模等\n        style_attribution = {}\n        \n        for factor in style_factors.columns:\n            # 计算组合在该风格因子上的暴露度\n            factor_exposure = self.portfolio_returns.corr(style_factors[factor])\n            \n            # 因子收益贡献\n            factor_return = style_factors[factor].mean() * 252\n            factor_contribution = factor_exposure * factor_return\n            \n            style_attribution[factor] = {\n                'exposure': factor_exposure,\n                'factor_return': factor_return,\n                'contribution': factor_contribution\n            }\n        \n        return style_attribution\n    \n    def risk_attribution(self, risk_factors: pd.DataFrame) -> Dict[str, any]:\n        \"\"\"\n        风险归因分析\n        \"\"\"\n        \n        # 计算组合在各风险因子上的Beta\n        risk_betas = {}\n        \n        for factor in risk_factors.columns:\n            # 回归计算Beta\n            aligned_data = pd.DataFrame({\n                'portfolio': self.portfolio_returns,\n                'factor': risk_factors[factor]\n            }).dropna()\n            \n            if len(aligned_data) > 20:\n                covariance = aligned_data['portfolio'].cov(aligned_data['factor'])\n                factor_variance = aligned_data['factor'].var()\n                \n                if factor_variance > 0:\n                    beta = covariance / factor_variance\n                else:\n                    beta = 0\n                \n                risk_betas[factor] = beta\n        \n        # 风险贡献计算\n        portfolio_variance = self.portfolio_returns.var()\n        factor_contributions = {}\n        \n        for factor, beta in risk_betas.items():\n            factor_variance = risk_factors[factor].var()\n            \n            # 该因子对组合方差的贡献\n            variance_contribution = (beta ** 2) * factor_variance\n            \n            factor_contributions[factor] = {\n                'beta': beta,\n                'variance_contribution': variance_contribution,\n                'volatility_contribution': np.sqrt(variance_contribution),\n                'contribution_percentage': variance_contribution / portfolio_variance\n            }\n        \n        return factor_contributions\n    \n    def time_varying_attribution(self, window: int = 252) -> pd.DataFrame:\n        \"\"\"\n        时变归因分析\n        \"\"\"\n        \n        rolling_attribution = []\n        \n        for i in range(window, len(self.portfolio_returns)):\n            # 滚动窗口数据\n            window_portfolio = self.portfolio_returns.iloc[i-window:i]\n            window_factors = self.factor_returns.iloc[i-window:i]\n            \n            # 临时创建归因分析器\n            temp_attribution = PerformanceAttribution(\n                window_portfolio, window_factors\n            )\n            \n            # 计算因子归因\n            attribution_result = temp_attribution.factor_attribution()\n            \n            # 记录结果\n            result_row = {'date': self.portfolio_returns.index[i-1]}\n            result_row.update(attribution_result['factor_contributions'].to_dict())\n            result_row['idiosyncratic'] = attribution_result['idiosyncratic_return']\n            result_row['r_squared'] = attribution_result['r_squared']\n            \n            rolling_attribution.append(result_row)\n        \n        return pd.DataFrame(rolling_attribution).set_index('date')\n```\n\n## 📊 综合实战案例：完整的风险管理系统\n\n```python\ndef comprehensive_risk_management_system():\n    \"\"\"\n    综合风险管理系统案例\n    \"\"\"\n    \n    print(\"=== 综合风险管理系统 ===\\n\")\n    \n    # 1. 数据准备\n    np.random.seed(42)\n    dates = pd.date_range('2020-01-01', '2023-12-31', freq='D')\n    \n    # 模拟投资组合收益率\n    portfolio_returns = pd.Series(\n        np.random.normal(0.001, 0.015, len(dates)),  # 日均收益0.1%，波动率1.5%\n        index=dates\n    )\n    \n    # 添加一些真实的市场特征\n    # 偶尔的大幅下跌\n    crash_dates = np.random.choice(len(dates), size=5, replace=False)\n    for crash_idx in crash_dates:\n        portfolio_returns.iloc[crash_idx] = -0.05  # 5%的单日跌幅\n    \n    # 模拟基准收益率\n    benchmark_returns = pd.Series(\n        np.random.normal(0.0008, 0.012, len(dates)),\n        index=dates\n    )\n    \n    # 模拟因子收益率\n    factor_returns = pd.DataFrame({\n        'market': benchmark_returns,\n        'value': np.random.normal(0.0003, 0.008, len(dates)),\n        'growth': np.random.normal(0.0005, 0.010, len(dates)),\n        'momentum': np.random.normal(0.0002, 0.012, len(dates))\n    }, index=dates)\n    \n    print(f\"分析期间: {dates[0].date()} 到 {dates[-1].date()}\")\n    print(f\"样本数量: {len(portfolio_returns)} 个交易日\\n\")\n    \n    # 2. 基础风险指标计算\n    print(\"=== 基础风险指标 ===\")\n    \n    risk_metrics = RiskMetrics(portfolio_returns, benchmark_returns)\n    basic_metrics = risk_metrics.calculate_basic_metrics()\n    \n    print(f\"年化收益率: {basic_metrics['annualized_return']:.2%}\")\n    print(f\"年化波动率: {basic_metrics['annualized_volatility']:.2%}\")\n    print(f\"夏普比率: {basic_metrics['sharpe_ratio']:.3f}\")\n    print(f\"偏度: {basic_metrics['skewness']:.3f}\")\n    print(f\"峰度: {basic_metrics['kurtosis']:.3f}\")\n    print(f\"胜率: {basic_metrics['win_rate']:.2%}\")\n    \n    # 回撤分析\n    drawdown_metrics, drawdown_series = risk_metrics.calculate_drawdown_metrics()\n    print(f\"\\n最大回撤: {drawdown_metrics['max_drawdown']:.2%}\")\n    print(f\"Calmar比率: {drawdown_metrics['calmar_ratio']:.3f}\")\n    print(f\"回撤次数: {drawdown_metrics['drawdown_count']}\")\n    print(f\"平均回撤持续时间: {drawdown_metrics['avg_drawdown_duration']:.1f} 天\")\n    \n    # 相对指标\n    relative_metrics = risk_metrics.calculate_relative_metrics()\n    if relative_metrics:\n        print(f\"\\nBeta: {relative_metrics['beta']:.3f}\")\n        print(f\"Alpha: {relative_metrics['alpha']:.2%}\")\n        print(f\"信息比率: {relative_metrics['information_ratio']:.3f}\")\n        print(f\"跟踪误差: {relative_metrics['tracking_error']:.2%}\")\n    \n    # 3. VaR和CVaR分析\n    print(\"\\n=== VaR和CVaR分析 ===\")\n    \n    advanced_risk = AdvancedRiskMetrics(portfolio_returns)\n    \n    # 不同方法的VaR\n    var_methods = ['historical', 'parametric', 'monte_carlo']\n    \n    for method in var_methods:\n        var_results = advanced_risk.calculate_var_cvar(method=method)\n        print(f\"\\n{method.upper()} 方法:\")\n        for confidence, values in var_results.items():\n            print(f\"  {confidence} VaR: {values['var_daily']:.2%} (日), {values['var_annual']:.2%} (年)\")\n            print(f\"  {confidence} CVaR: {values['cvar_daily']:.2%} (日), {values['cvar_annual']:.2%} (年)\")\n    \n    # 极值理论分析\n    evt_results = advanced_risk.calculate_extreme_value_metrics()\n    if 'error' not in evt_results:\n        print(f\"\\n极值理论分析:\")\n        print(f\"  阈值: {evt_results['threshold']:.2%}\")\n        print(f\"  形状参数: {evt_results['shape_parameter']:.3f}\")\n        print(f\"  EVT-VaR(95%): {evt_results['evt_var_95']:.2%}\")\n        print(f\"  EVT-VaR(99%): {evt_results['evt_var_99']:.2%}\")\n    \n    # 4. 动态风险控制\n    print(\"\\n=== 动态风险控制 ===\")\n    \n    drawdown_controller = DrawdownController(max_drawdown_threshold=0.15)\n    \n    # 模拟交易信号和仓位控制\n    signals = np.random.uniform(-1, 1, 100)  # 模拟信号\n    current_equity = 1000000  # 当前权益\n    \n    position_sizes = []\n    for signal in signals[:10]:  # 只展示前10个\n        position_size = drawdown_controller.calculate_position_size(\n            signal_strength=signal,\n            current_equity=current_equity,\n            historical_returns=portfolio_returns.iloc[-252:],\n            volatility_target=0.15\n        )\n        position_sizes.append(position_size)\n    \n    print(f\"动态仓位控制示例 (前10个信号):\")\n    for i, (signal, position) in enumerate(zip(signals[:10], position_sizes)):\n        print(f\"  信号 {i+1}: {signal:.2f} → 仓位: {position:.2f}\")\n    \n    # Kelly公式仓位计算\n    win_rate = basic_metrics['win_rate']\n    avg_win = basic_metrics['positive_returns_mean']\n    avg_loss = basic_metrics['negative_returns_mean']\n    \n    kelly_position = drawdown_controller.calculate_kelly_criterion(win_rate, avg_win, avg_loss)\n    print(f\"\\nKelly公式建议仓位: {kelly_position:.2%}\")\n    \n    # 5. VaR回测\n    print(\"\\n=== VaR回测 ===\")\n    \n    var_calculator = VaRCalculator(portfolio_returns)\n    \n    # 计算历史VaR\n    historical_var = var_calculator.historical_var(confidence_level=0.95, window_size=252)\n    \n    # 回测VaR\n    if len(historical_var) > 100:\n        test_returns = portfolio_returns.loc[historical_var.index]\n        \n        backtest_results = var_calculator.backtesting_var(\n            historical_var, test_returns, confidence_level=0.95\n        )\n        \n        print(f\"VaR违背次数: {backtest_results['n_violations']}\")\n        print(f\"预期违背次数: {backtest_results['expected_violations']:.1f}\")\n        print(f\"违背率: {backtest_results['violation_rate']:.2%}\")\n        print(f\"Kupiec检验p值: {backtest_results['kupiec_p_value']:.4f}\")\n        \n        if backtest_results['kupiec_p_value'] > 0.05:\n            print(\"VaR模型通过Kupiec检验 ✓\")\n        else:\n            print(\"VaR模型未通过Kupiec检验 ✗\")\n    \n    # 6. 绩效评估\n    print(\"\\n=== 绩效评估 ===\")\n    \n    perf_evaluator = PerformanceEvaluator(portfolio_returns, benchmark_returns)\n    comprehensive_metrics = perf_evaluator.calculate_comprehensive_metrics()\n    \n    print(f\"总收益率: {comprehensive_metrics['total_return']:.2%}\")\n    print(f\"年化收益率: {comprehensive_metrics['annualized_return']:.2%}\")\n    print(f\"年化波动率: {comprehensive_metrics['volatility']:.2%}\")\n    print(f\"夏普比率: {comprehensive_metrics['sharpe_ratio']:.3f}\")\n    print(f\"Sortino比率: {comprehensive_metrics['sortino_ratio']:.3f}\")\n    print(f\"最大回撤: {comprehensive_metrics['max_drawdown']:.2%}\")\n    \n    if benchmark_returns is not None:\n        print(f\"\\n相对基准指标:\")\n        print(f\"超额收益: {comprehensive_metrics['excess_return']:.2%}\")\n        print(f\"信息比率: {comprehensive_metrics['information_ratio']:.3f}\")\n        print(f\"Beta: {comprehensive_metrics['beta']:.3f}\")\n        print(f\"Alpha: {comprehensive_metrics['alpha']:.2%}\")\n    \n    # 7. 绩效归因\n    print(\"\\n=== 绩效归因分析 ===\")\n    \n    attribution = PerformanceAttribution(portfolio_returns, factor_returns)\n    factor_attribution = attribution.factor_attribution(method='regression')\n    \n    print(f\"因子归因结果 (R² = {factor_attribution['r_squared']:.3f}):\")\n    \n    for factor, contribution in factor_attribution['factor_contributions'].items():\n        print(f\"  {factor}: {contribution:.2%}\")\n    \n    print(f\"  特异性收益: {factor_attribution['idiosyncratic_return']:.2%}\")\n    print(f\"  总归因收益: {factor_attribution['total_attributed']:.2%}\")\n    print(f\"  实际收益: {factor_attribution['actual_return']:.2%}\")\n    \n    # 8. 风险预算系统\n    print(\"\\n=== 风险预算系统 ===\")\n    \n    # 模拟多策略组合\n    strategies = {\n        'momentum': {\n            'expected_return': 0.12,\n            'volatility': 0.18,\n            'sharpe_ratio': 0.6,\n            'max_drawdown': -0.15\n        },\n        'mean_reversion': {\n            'expected_return': 0.08,\n            'volatility': 0.12,\n            'sharpe_ratio': 0.5,\n            'max_drawdown': -0.08\n        },\n        'trend_following': {\n            'expected_return': 0.10,\n            'volatility': 0.20,\n            'sharpe_ratio': 0.4,\n            'max_drawdown': -0.12\n        }\n    }\n    \n    risk_budget_system = RiskBudgetSystem(total_risk_budget=0.02)\n    risk_allocations = risk_budget_system.allocate_risk_budget(strategies)\n    \n    print(f\"风险预算分配:\")\n    for strategy, allocation in risk_allocations.items():\n        print(f\"  {strategy}:\")\n        print(f\"    风险权重: {allocation['risk_weight']:.1%}\")\n        print(f\"    日度风险预算: {allocation['daily_vol_budget']:.2%}\")\n        print(f\"    最大仓位: {allocation['max_position']:.1%}\")\n    \n    # 9. 压力测试\n    print(\"\\n=== 压力测试 ===\")\n    \n    # 模拟压力情景\n    stress_scenarios = {\n        '2008金融危机': portfolio_returns - 0.03,  # 整体下移3%\n        '新冠疫情': portfolio_returns * 1.5,  # 波动率增加50%\n        '黑天鹅事件': portfolio_returns.copy()\n    }\n    \n    # 为黑天鹅事件添加极端下跌\n    stress_scenarios['黑天鹅事件'].iloc[100] = -0.20  # 单日20%跌幅\n    \n    for scenario_name, scenario_returns in stress_scenarios.items():\n        scenario_metrics = RiskMetrics(scenario_returns)\n        scenario_basic = scenario_metrics.calculate_basic_metrics()\n        scenario_dd, _ = scenario_metrics.calculate_drawdown_metrics()\n        \n        print(f\"\\n{scenario_name}:\")\n        print(f\"  年化收益: {scenario_basic['annualized_return']:.2%}\")\n        print(f\"  年化波动: {scenario_basic['annualized_volatility']:.2%}\")\n        print(f\"  最大回撤: {scenario_dd['max_drawdown']:.2%}\")\n        print(f\"  夏普比率: {scenario_basic['sharpe_ratio']:.3f}\")\n    \n    print(\"\\n=== 风险管理系统总结 ===\")\n    print(\"✓ 完成基础风险指标计算\")\n    print(\"✓ 完成VaR/CVaR多方法验证\")\n    print(\"✓ 完成动态风险控制设计\")\n    print(\"✓ 完成VaR模型回测验证\")\n    print(\"✓ 完成综合绩效评估\")\n    print(\"✓ 完成绩效归因分析\")\n    print(\"✓ 完成风险预算分配\")\n    print(\"✓ 完成压力测试分析\")\n    \n    return {\n        'basic_metrics': basic_metrics,\n        'drawdown_metrics': drawdown_metrics,\n        'var_backtest': backtest_results if 'backtest_results' in locals() else None,\n        'performance_metrics': comprehensive_metrics,\n        'attribution_results': factor_attribution,\n        'risk_allocations': risk_allocations\n    }\n\n# 运行综合案例\nif __name__ == \"__main__\":\n    risk_management_results = comprehensive_risk_management_system()\n```\n\n## 📚 本章小结\n\n本章全面介绍了风险管理与绩效评估的核心技术：\n\n### 🔑 核心要点\n\n1. **风险度量**: 从基础指标到高级风险度量的完整体系\n2. **VaR/CVaR**: 多种计算方法和实际应用技巧\n3. **回撤控制**: 动态仓位调整和风险预算管理\n4. **绩效评估**: 绝对和相对绩效的全面评估框架\n5. **归因分析**: 因子归因、风格归因和风险归因技术\n\n### 📈 实践要点\n\n- 风险管理应该是动态和自适应的\n- 多种VaR方法结合使用提高准确性\n- 回撤控制需要平衡风险和收益\n- 绩效评估要考虑风险调整和基准比较\n- 归因分析有助于理解收益来源和风险暴露\n\n### 🎯 本章完成\n\n至此，第三部分"高级量化策略"的所有章节已经完成，涵盖了：\n- 多因子模型与组合优化\n- 机器学习在量化投资中的应用\n- 风险管理与绩效评估\n\n这些内容为构建完整的量化投资系统提供了坚实的理论基础和实践指导。\n\n---\n\n**风险提示**: 风险管理是量化投资的核心，任何模型和方法都有其局限性，需要结合实际市场情况灵活应用。