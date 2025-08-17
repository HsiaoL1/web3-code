# 第四章：股票投资与分析（下篇）

## 📖 章节概述

本章下篇将深入探讨股票投资的高级分析方法，包括量化选股策略、多因子模型构建、事件驱动策略以及 A 股市场的特色分析。这些内容将帮助您从传统投资分析进阶到系统化的量化投资方法。

## 🎯 学习目标

- 掌握量化选股的核心方法和策略
- 学会构建和应用多因子模型
- 理解事件驱动投资策略的原理和实践
- 深入了解 A 股市场的特色和投资机会
- 掌握股票投资组合的构建和优化
- 建立完整的股票量化投资体系

## 4.6 量化选股策略

### 4.6.1 量化选股基础

量化选股是运用数学、统计学、计算机技术等方法，通过建立数量化模型来指导股票投资决策的过程。

```python
import pandas as pd
import numpy as np
from sklearn.ensemble import RandomForestRegressor
from sklearn.preprocessing import StandardScaler
from sklearn.model_selection import train_test_split
import warnings
warnings.filterwarnings('ignore')

class QuantitativeStockSelection:
    """量化选股策略框架"""

    def __init__(self):
        self.factors = {}
        self.models = {}
        self.scores = {}

    def calculate_financial_factors(self, financial_data):
        """计算财务因子"""
        factors = pd.DataFrame(index=financial_data.index)

        # 盈利能力因子
        factors['roe'] = financial_data['净利润'] / financial_data['股东权益']
        factors['roa'] = financial_data['净利润'] / financial_data['总资产']
        factors['gross_margin'] = financial_data['毛利润'] / financial_data['营业收入']
        factors['net_margin'] = financial_data['净利润'] / financial_data['营业收入']

        # 成长性因子
        factors['revenue_growth'] = financial_data['营业收入'].pct_change(4)  # 同比增长
        factors['profit_growth'] = financial_data['净利润'].pct_change(4)
        factors['eps_growth'] = financial_data['每股收益'].pct_change(4)

        # 估值因子
        factors['pe'] = financial_data['股价'] / financial_data['每股收益']
        factors['pb'] = financial_data['股价'] / financial_data['每股净资产']
        factors['ps'] = financial_data['市值'] / financial_data['营业收入']

        # 质量因子
        factors['debt_ratio'] = financial_data['总负债'] / financial_data['总资产']
        factors['current_ratio'] = financial_data['流动资产'] / financial_data['流动负债']
        factors['asset_turnover'] = financial_data['营业收入'] / financial_data['总资产']

        return factors

    def calculate_technical_factors(self, price_data):
        """计算技术因子"""
        factors = pd.DataFrame(index=price_data.index)

        # 动量因子
        factors['return_1m'] = price_data['close'].pct_change(20)
        factors['return_3m'] = price_data['close'].pct_change(60)
        factors['return_6m'] = price_data['close'].pct_change(120)
        factors['return_12m'] = price_data['close'].pct_change(240)

        # 反转因子
        factors['return_1w'] = price_data['close'].pct_change(5)
        factors['return_1d'] = price_data['close'].pct_change(1)

        # 波动率因子
        factors['volatility_20d'] = price_data['close'].pct_change().rolling(20).std()
        factors['volatility_60d'] = price_data['close'].pct_change().rolling(60).std()

        # 流动性因子
        factors['turnover_20d'] = price_data['volume'].rolling(20).mean()
        factors['amount_20d'] = (price_data['close'] * price_data['volume']).rolling(20).mean()

        return factors

    def factor_preprocessing(self, factors):
        """因子预处理"""
        # 去极值
        def winsorize(series, lower=0.05, upper=0.95):
            return series.clip(
                lower=series.quantile(lower),
                upper=series.quantile(upper)
            )

        # 标准化
        processed_factors = factors.copy()
        for col in processed_factors.columns:
            processed_factors[col] = winsorize(processed_factors[col])
            processed_factors[col] = (processed_factors[col] - processed_factors[col].mean()) / processed_factors[col].std()

        return processed_factors

    def build_selection_model(self, factors, returns, method='random_forest'):
        """构建选股模型"""
        # 准备数据
        X = factors.dropna()
        y = returns.loc[X.index]

        X_train, X_test, y_train, y_test = train_test_split(
            X, y, test_size=0.3, random_state=42
        )

        if method == 'random_forest':
            model = RandomForestRegressor(
                n_estimators=100,
                max_depth=10,
                random_state=42
            )

        model.fit(X_train, y_train)

        # 计算因子重要性
        feature_importance = pd.Series(
            model.feature_importances_,
            index=X.columns
        ).sort_values(ascending=False)

        return model, feature_importance

    def select_stocks(self, factors, model, top_n=50):
        """选择股票"""
        # 预测收益率
        predicted_returns = model.predict(factors.dropna())

        # 创建股票评分
        scores = pd.Series(predicted_returns, index=factors.dropna().index)

        # 选择前N只股票
        selected_stocks = scores.nlargest(top_n)

        return selected_stocks

# 使用示例
selector = QuantitativeStockSelection()

# 假设有财务数据和价格数据
# financial_factors = selector.calculate_financial_factors(financial_data)
# technical_factors = selector.calculate_technical_factors(price_data)
# all_factors = pd.concat([financial_factors, technical_factors], axis=1)
# processed_factors = selector.factor_preprocessing(all_factors)
# model, importance = selector.build_selection_model(processed_factors, future_returns)
# selected_stocks = selector.select_stocks(processed_factors, model)
```

### 4.6.2 多因子选股模型

```python
class MultiFactorStockSelection:
    """多因子选股模型"""

    def __init__(self):
        self.factor_categories = {
            'value': ['pe', 'pb', 'ps', 'pcf'],
            'growth': ['revenue_growth', 'profit_growth', 'eps_growth'],
            'quality': ['roe', 'roa', 'debt_ratio', 'current_ratio'],
            'momentum': ['return_1m', 'return_3m', 'return_6m'],
            'volatility': ['volatility_20d', 'volatility_60d']
        }

    def calculate_factor_scores(self, factors):
        """计算因子得分"""
        factor_scores = pd.DataFrame(index=factors.index)

        for category, factor_list in self.factor_categories.items():
            available_factors = [f for f in factor_list if f in factors.columns]
            if available_factors:
                # 计算每个类别的综合得分
                category_data = factors[available_factors]

                # 对于估值因子，数值越小越好（负号处理）
                if category == 'value':
                    category_data = -category_data
                elif category == 'volatility':
                    category_data = -category_data  # 波动率越小越好

                # 计算综合得分（等权重平均）
                factor_scores[f'{category}_score'] = category_data.mean(axis=1)

        return factor_scores

    def composite_scoring(self, factor_scores, weights=None):
        """复合打分"""
        if weights is None:
            weights = {
                'value_score': 0.25,
                'growth_score': 0.25,
                'quality_score': 0.25,
                'momentum_score': 0.15,
                'volatility_score': 0.10
            }

        composite_score = pd.Series(0, index=factor_scores.index)

        for factor, weight in weights.items():
            if factor in factor_scores.columns:
                composite_score += factor_scores[factor] * weight

        return composite_score

    def sector_neutral_selection(self, scores, sector_info, stocks_per_sector=5):
        """行业中性选股"""
        selected_stocks = []

        for sector in sector_info['sector'].unique():
            sector_stocks = sector_info[sector_info['sector'] == sector].index
            sector_scores = scores.loc[sector_stocks].dropna()

            if len(sector_scores) >= stocks_per_sector:
                top_stocks = sector_scores.nlargest(stocks_per_sector)
                selected_stocks.extend(top_stocks.index.tolist())

        return selected_stocks

# 使用示例
multi_factor_selector = MultiFactorStockSelection()
# factor_scores = multi_factor_selector.calculate_factor_scores(all_factors)
# composite_scores = multi_factor_selector.composite_scoring(factor_scores)
# selected_stocks = multi_factor_selector.sector_neutral_selection(composite_scores, sector_data)
```

## 4.7 多因子模型

### 4.7.1 Fama-French 三因子模型

```python
import numpy as np
import pandas as pd
from sklearn.linear_model import LinearRegression
import matplotlib.pyplot as plt
import seaborn as sns

class FamaFrenchModel:
    """Fama-French三因子模型"""

    def __init__(self):
        self.factors = {}
        self.portfolios = {}

    def construct_size_factor(self, market_cap, returns):
        """构建规模因子SMB (Small Minus Big)"""
        # 按市值排序，分为大盘股和小盘股
        size_median = market_cap.median()

        small_cap_stocks = market_cap[market_cap <= size_median].index
        big_cap_stocks = market_cap[market_cap > size_median].index

        # 计算小盘股和大盘股的平均收益率
        small_returns = returns.loc[small_cap_stocks].mean()
        big_returns = returns.loc[big_cap_stocks].mean()

        smb = small_returns - big_returns

        return smb

    def construct_value_factor(self, book_to_market, returns):
        """构建价值因子HML (High Minus Low)"""
        # 按账面市值比排序，分为高价值股和低价值股
        bm_terciles = book_to_market.quantile([1/3, 2/3])

        high_bm_stocks = book_to_market[book_to_market >= bm_terciles.iloc[1]].index
        low_bm_stocks = book_to_market[book_to_market <= bm_terciles.iloc[0]].index

        # 计算高价值股和低价值股的平均收益率
        high_returns = returns.loc[high_bm_stocks].mean()
        low_returns = returns.loc[low_bm_stocks].mean()

        hml = high_returns - low_returns

        return hml

    def construct_factors(self, returns, market_cap, book_to_market, market_return, risk_free_rate):
        """构建三因子"""
        # 市场因子
        market_premium = market_return - risk_free_rate

        # 规模因子
        smb = self.construct_size_factor(market_cap, returns)

        # 价值因子
        hml = self.construct_value_factor(book_to_market, returns)

        factors = pd.DataFrame({
            'MKT': market_premium,
            'SMB': smb,
            'HML': hml
        })

        return factors

    def run_regression(self, stock_returns, factors, risk_free_rate):
        """运行三因子回归"""
        # 计算超额收益
        excess_returns = stock_returns - risk_free_rate

        results = {}

        for stock in excess_returns.columns:
            y = excess_returns[stock].dropna()
            X = factors.loc[y.index]

            if len(y) > 0 and len(X) > 0:
                # 线性回归
                model = LinearRegression()
                model.fit(X, y)

                # 计算统计量
                alpha = model.intercept_
                betas = model.coef_
                r_squared = model.score(X, y)

                results[stock] = {
                    'alpha': alpha,
                    'beta_market': betas[0],
                    'beta_smb': betas[1],
                    'beta_hml': betas[2],
                    'r_squared': r_squared
                }

        return pd.DataFrame(results).T

    def factor_attribution(self, results, factors):
        """因子归因分析"""
        attribution = pd.DataFrame(index=results.index)

        # 计算各因子对收益的贡献
        attribution['market_contrib'] = results['beta_market'] * factors['MKT'].mean()
        attribution['size_contrib'] = results['beta_smb'] * factors['SMB'].mean()
        attribution['value_contrib'] = results['beta_hml'] * factors['HML'].mean()
        attribution['alpha_contrib'] = results['alpha']

        attribution['total_return'] = attribution.sum(axis=1)

        return attribution

# 使用示例
ff_model = FamaFrenchModel()
# factors = ff_model.construct_factors(returns, market_cap, book_to_market, market_return, rf_rate)
# regression_results = ff_model.run_regression(stock_returns, factors, rf_rate)
# attribution = ff_model.factor_attribution(regression_results, factors)
```

### 4.7.2 多因子模型优化

```python
class AdvancedMultiFactorModel:
    """高级多因子模型"""

    def __init__(self):
        self.factor_weights = {}
        self.factor_loadings = {}

    def dynamic_factor_selection(self, factors, returns, window=252):
        """动态因子选择"""
        selected_factors = {}

        for i in range(window, len(returns)):
            period_returns = returns.iloc[i-window:i]
            period_factors = factors.iloc[i-window:i]

            # 计算因子IC（信息系数）
            ic_scores = {}
            for factor in period_factors.columns:
                ic = period_factors[factor].corr(period_returns.mean(axis=1))
                ic_scores[factor] = abs(ic) if not np.isnan(ic) else 0

            # 选择前N个因子
            top_factors = sorted(ic_scores.items(), key=lambda x: x[1], reverse=True)[:10]
            selected_factors[returns.index[i]] = [f[0] for f in top_factors]

        return selected_factors

    def calculate_factor_ic(self, factors, returns, periods=[1, 5, 10, 20]):
        """计算因子信息系数"""
        ic_results = {}

        for period in periods:
            future_returns = returns.shift(-period)
            ic_values = {}

            for factor in factors.columns:
                ic_series = []
                for date in factors.index[:-period]:
                    factor_values = factors.loc[date, :]
                    future_ret = future_returns.loc[date + pd.Timedelta(days=period), :]

                    valid_data = pd.concat([factor_values, future_ret], axis=1).dropna()
                    if len(valid_data) > 10:
                        ic = valid_data.iloc[:, 0].corr(valid_data.iloc[:, 1])
                        ic_series.append(ic)

                ic_values[factor] = {
                    'mean_ic': np.mean(ic_series),
                    'ic_std': np.std(ic_series),
                    'ic_ir': np.mean(ic_series) / np.std(ic_series) if np.std(ic_series) > 0 else 0
                }

            ic_results[f'{period}d'] = ic_values

        return ic_results

    def factor_decay_analysis(self, factors, returns, max_periods=60):
        """因子衰减分析"""
        decay_results = {}

        for factor in factors.columns:
            ic_decay = []

            for period in range(1, max_periods + 1):
                future_returns = returns.shift(-period)
                ic_values = []

                for date in factors.index[:-period]:
                    factor_vals = factors.loc[date, factor]
                    future_ret = future_returns.loc[date, :]

                    valid_data = pd.concat([factor_vals, future_ret], axis=1).dropna()
                    if len(valid_data) > 10:
                        ic = valid_data.iloc[:, 0].corr(valid_data.iloc[:, 1])
                        ic_values.append(ic)

                ic_decay.append(np.mean(ic_values) if ic_values else 0)

            decay_results[factor] = ic_decay

        return pd.DataFrame(decay_results, index=range(1, max_periods + 1))

    def construct_factor_portfolio(self, factor_scores, weights_method='equal_weight'):
        """构建因子投资组合"""
        if weights_method == 'equal_weight':
            weights = pd.Series(1/len(factor_scores), index=factor_scores.index)

        elif weights_method == 'score_weight':
            # 根据因子得分分配权重
            positive_scores = factor_scores[factor_scores > 0]
            weights = positive_scores / positive_scores.sum()

        elif weights_method == 'long_short':
            # 多空组合
            top_quantile = factor_scores.quantile(0.8)
            bottom_quantile = factor_scores.quantile(0.2)

            long_stocks = factor_scores[factor_scores >= top_quantile]
            short_stocks = factor_scores[factor_scores <= bottom_quantile]

            weights = pd.Series(0, index=factor_scores.index)
            weights.loc[long_stocks.index] = 0.5 / len(long_stocks)
            weights.loc[short_stocks.index] = -0.5 / len(short_stocks)

        return weights

# 使用示例
advanced_model = AdvancedMultiFactorModel()
# selected_factors = advanced_model.dynamic_factor_selection(factors, returns)
# ic_analysis = advanced_model.calculate_factor_ic(factors, returns)
# decay_analysis = advanced_model.factor_decay_analysis(factors, returns)
```

## 4.8 事件驱动策略

### 4.8.1 盈利公告策略

```python
class EarningsAnnouncementStrategy:
    """盈利公告事件驱动策略"""

    def __init__(self):
        self.positions = {}
        self.results = {}

    def calculate_earnings_surprise(self, actual_eps, consensus_eps):
        """计算盈利超预期程度"""
        surprise = (actual_eps - consensus_eps) / abs(consensus_eps)
        return surprise

    def identify_earnings_events(self, earnings_data, surprise_threshold=0.05):
        """识别盈利事件"""
        events = {}

        for date, row in earnings_data.iterrows():
            surprise = self.calculate_earnings_surprise(row['actual_eps'], row['consensus_eps'])

            if abs(surprise) > surprise_threshold:
                events[date] = {
                    'stock': row['stock'],
                    'surprise': surprise,
                    'direction': 'positive' if surprise > 0 else 'negative'
                }

        return events

    def post_earnings_drift(self, price_data, earnings_events, holding_periods=[5, 10, 20]):
        """盈利公告后漂移分析"""
        drift_results = {}

        for event_date, event_info in earnings_events.items():
            stock = event_info['stock']
            surprise = event_info['surprise']

            if stock in price_data.columns:
                event_price = price_data.loc[event_date, stock]

                for period in holding_periods:
                    try:
                        future_date = event_date + pd.Timedelta(days=period)
                        future_price = price_data.loc[future_date, stock]

                        return_period = (future_price - event_price) / event_price

                        drift_results[f'{event_date}_{stock}_{period}d'] = {
                            'surprise': surprise,
                            'return': return_period,
                            'period': period
                        }
                    except:
                        continue

        return pd.DataFrame(drift_results).T

    def earnings_momentum_strategy(self, earnings_events, price_data, holding_period=10):
        """盈利动量策略"""
        trades = []

        for event_date, event_info in earnings_events.items():
            stock = event_info['stock']
            surprise = event_info['surprise']

            # 正向超预期做多，负向超预期做空
            if surprise > 0.1:  # 大幅超预期
                action = 'buy'
                target_return = 0.05  # 目标收益5%
            elif surprise < -0.1:  # 大幅低于预期
                action = 'sell'
                target_return = -0.05
            else:
                continue

            try:
                entry_price = price_data.loc[event_date, stock]
                exit_date = event_date + pd.Timedelta(days=holding_period)
                exit_price = price_data.loc[exit_date, stock]

                actual_return = (exit_price - entry_price) / entry_price
                if action == 'sell':
                    actual_return = -actual_return

                trades.append({
                    'event_date': event_date,
                    'stock': stock,
                    'action': action,
                    'surprise': surprise,
                    'entry_price': entry_price,
                    'exit_price': exit_price,
                    'return': actual_return,
                    'target_met': (actual_return >= target_return) if action == 'buy' else (actual_return >= -target_return)
                })
            except:
                continue

        return pd.DataFrame(trades)

# 使用示例
earnings_strategy = EarningsAnnouncementStrategy()
# earnings_events = earnings_strategy.identify_earnings_events(earnings_data)
# drift_analysis = earnings_strategy.post_earnings_drift(price_data, earnings_events)
# strategy_results = earnings_strategy.earnings_momentum_strategy(earnings_events, price_data)
```

### 4.8.2 并购重组策略

```python
class MergerArbitrageStrategy:
    """并购套利策略"""

    def __init__(self):
        self.active_deals = {}
        self.completed_deals = {}

    def identify_merger_opportunities(self, announcement_data):
        """识别并购机会"""
        opportunities = {}

        for idx, deal in announcement_data.iterrows():
            target = deal['target_stock']
            acquirer = deal['acquirer_stock']
            offer_price = deal['offer_price']
            deal_type = deal['deal_type']  # cash, stock, mixed
            expected_close = deal['expected_close_date']

            opportunities[deal['deal_id']] = {
                'target': target,
                'acquirer': acquirer,
                'offer_price': offer_price,
                'deal_type': deal_type,
                'announcement_date': deal['announcement_date'],
                'expected_close': expected_close,
                'probability': self.estimate_completion_probability(deal)
            }

        return opportunities

    def estimate_completion_probability(self, deal_info):
        """估计交易完成概率"""
        base_prob = 0.8  # 基础完成概率

        # 根据各种因素调整概率
        adjustments = 0

        # 监管风险
        if deal_info['industry'] in ['technology', 'healthcare', 'financial']:
            adjustments -= 0.1

        # 溢价水平
        premium = deal_info['offer_premium']
        if premium > 0.5:  # 溢价过高
            adjustments -= 0.15
        elif premium < 0.1:  # 溢价过低
            adjustments -= 0.1

        # 敌意收购
        if deal_info['hostile']:
            adjustments -= 0.2

        # 融资条件
        if deal_info['financing_contingent']:
            adjustments -= 0.1

        final_prob = max(0.1, min(0.95, base_prob + adjustments))
        return final_prob

    def calculate_arbitrage_spread(self, current_price, offer_price, completion_prob, time_to_close):
        """计算套利价差"""
        expected_return = completion_prob * (offer_price - current_price) / current_price
        annualized_return = expected_return * (365 / time_to_close)

        risk_adjusted_return = expected_return - (1 - completion_prob) * 0.2  # 假设失败损失20%

        return {
            'spread': (offer_price - current_price) / current_price,
            'expected_return': expected_return,
            'annualized_return': annualized_return,
            'risk_adjusted_return': risk_adjusted_return
        }

    def execute_arbitrage_strategy(self, opportunities, price_data, min_return=0.05):
        """执行套利策略"""
        positions = {}

        for deal_id, deal_info in opportunities.items():
            target = deal_info['target']
            current_price = price_data.loc[deal_info['announcement_date'], target]

            days_to_close = (deal_info['expected_close'] - deal_info['announcement_date']).days

            arb_metrics = self.calculate_arbitrage_spread(
                current_price, deal_info['offer_price'],
                deal_info['probability'], days_to_close
            )

            if arb_metrics['risk_adjusted_return'] > min_return:
                positions[deal_id] = {
                    'action': 'buy_target',
                    'target_stock': target,
                    'entry_price': current_price,
                    'target_price': deal_info['offer_price'],
                    'position_size': self.calculate_position_size(arb_metrics),
                    'expected_return': arb_metrics['expected_return']
                }

                # 对于股票交易，可能需要做空收购方
                if deal_info['deal_type'] == 'stock':
                    exchange_ratio = deal_info['exchange_ratio']
                    acquirer_price = price_data.loc[deal_info['announcement_date'], deal_info['acquirer']]

                    positions[f"{deal_id}_hedge"] = {
                        'action': 'sell_acquirer',
                        'acquirer_stock': deal_info['acquirer'],
                        'entry_price': acquirer_price,
                        'position_size': exchange_ratio * positions[deal_id]['position_size']
                    }

        return positions

    def calculate_position_size(self, arb_metrics, max_position_size=0.05):
        """计算仓位大小"""
        # 基于风险调整收益确定仓位
        risk_score = arb_metrics['risk_adjusted_return']
        position_size = min(max_position_size, risk_score * 0.1)
        return max(0.01, position_size)  # 最小仓位1%

# 使用示例
merger_strategy = MergerArbitrageStrategy()
# opportunities = merger_strategy.identify_merger_opportunities(merger_data)
# positions = merger_strategy.execute_arbitrage_strategy(opportunities, price_data)
```

## 4.9 A 股市场特色分析

### 4.9.1 A 股市场特征

```python
class AShareMarketAnalysis:
    """A股市场特色分析"""

    def __init__(self):
        self.market_features = {}

    def analyze_trading_patterns(self, price_data, volume_data):
        """分析A股交易模式"""
        patterns = {}

        # 计算日内收益率分布
        returns = price_data.pct_change().dropna()

        patterns['return_distribution'] = {
            'mean': returns.mean().mean(),
            'std': returns.std().mean(),
            'skewness': returns.skew().mean(),
            'kurtosis': returns.kurtosis().mean()
        }

        # 分析涨跌停现象
        patterns['limit_moves'] = self.analyze_limit_moves(returns)

        # 分析T+1交易影响
        patterns['t_plus_1_effect'] = self.analyze_t_plus_1_effect(returns)

        # 分析午间效应
        patterns['lunch_break_effect'] = self.analyze_lunch_break_effect(price_data)

        return patterns

    def analyze_limit_moves(self, returns, limit_threshold=0.098):
        """分析涨跌停现象"""
        limit_up = (returns >= limit_threshold).sum().sum()
        limit_down = (returns <= -limit_threshold).sum().sum()
        total_observations = len(returns) * len(returns.columns)

        return {
            'limit_up_frequency': limit_up / total_observations,
            'limit_down_frequency': limit_down / total_observations,
            'limit_up_count': limit_up,
            'limit_down_count': limit_down
        }

    def analyze_t_plus_1_effect(self, returns):
        """分析T+1交易制度影响"""
        # 分析大涨大跌后的反转效应
        extreme_up = returns > returns.quantile(0.95)
        extreme_down = returns < returns.quantile(0.05)

        next_day_returns_after_up = returns.shift(-1)[extreme_up].mean()
        next_day_returns_after_down = returns.shift(-1)[extreme_down].mean()

        return {
            'reversal_after_extreme_up': next_day_returns_after_up.mean(),
            'reversal_after_extreme_down': next_day_returns_after_down.mean()
        }

    def analyze_lunch_break_effect(self, price_data):
        """分析午间休市效应"""
        # 需要分钟级数据来分析上午收盘和下午开盘的关系
        # 这里简化处理，分析日内波动模式

        daily_range = (price_data.max(axis=1) - price_data.min(axis=1)) / price_data.mean(axis=1)

        return {
            'average_daily_range': daily_range.mean(),
            'range_volatility': daily_range.std()
        }

    def sector_rotation_analysis(self, sector_data, price_data):
        """行业轮动分析"""
        sector_returns = {}

        for sector, stocks in sector_data.items():
            sector_prices = price_data[stocks].mean(axis=1)
            sector_returns[sector] = sector_prices.pct_change()

        sector_returns_df = pd.DataFrame(sector_returns)

        # 计算行业相关性
        correlation_matrix = sector_returns_df.corr()

        # 分析行业动量
        momentum_periods = [20, 60, 120]
        sector_momentum = {}

        for period in momentum_periods:
            momentum = sector_returns_df.rolling(period).sum()
            sector_momentum[f'{period}d_momentum'] = momentum.iloc[-1].sort_values(ascending=False)

        return {
            'correlation_matrix': correlation_matrix,
            'sector_momentum': sector_momentum
        }

    def policy_impact_analysis(self, policy_dates, price_data):
        """政策影响分析"""
        policy_effects = {}

        for policy_date, policy_info in policy_dates.items():
            # 分析政策前后的市场表现
            pre_period = price_data.loc[policy_date - pd.Timedelta(days=5):policy_date - pd.Timedelta(days=1)]
            post_period = price_data.loc[policy_date:policy_date + pd.Timedelta(days=5)]

            pre_returns = pre_period.pct_change().mean()
            post_returns = post_period.pct_change().mean()

            policy_effects[policy_date] = {
                'policy_type': policy_info['type'],
                'pre_announcement_return': pre_returns.mean(),
                'post_announcement_return': post_returns.mean(),
                'impact_magnitude': post_returns.mean() - pre_returns.mean()
            }

        return policy_effects

# 使用示例
a_share_analysis = AShareMarketAnalysis()
# trading_patterns = a_share_analysis.analyze_trading_patterns(price_data, volume_data)
# sector_rotation = a_share_analysis.sector_rotation_analysis(sector_data, price_data)
# policy_impacts = a_share_analysis.policy_impact_analysis(policy_dates, price_data)
```

### 4.9.2 A 股特色投资策略

```python
class AShareSpecificStrategies:
    """A股特色投资策略"""

    def __init__(self):
        self.strategies = {}

    def new_stock_ipo_strategy(self, ipo_data, price_data):
        """新股申购和上市策略"""
        ipo_performance = {}

        for ipo_date, ipo_info in ipo_data.items():
            stock = ipo_info['stock_code']
            issue_price = ipo_info['issue_price']

            # 分析首日表现
            first_day_close = price_data.loc[ipo_date, stock]
            first_day_return = (first_day_close - issue_price) / issue_price

            # 分析后续表现
            returns_1w = (price_data.loc[ipo_date + pd.Timedelta(days=7), stock] - first_day_close) / first_day_close
            returns_1m = (price_data.loc[ipo_date + pd.Timedelta(days=30), stock] - first_day_close) / first_day_close

            ipo_performance[stock] = {
                'first_day_return': first_day_return,
                '1_week_return': returns_1w,
                '1_month_return': returns_1m,
                'industry': ipo_info['industry']
            }

        return pd.DataFrame(ipo_performance).T

    def st_stock_strategy(self, st_stocks, financial_data, price_data):
        """ST股票投资策略"""
        st_analysis = {}

        for stock in st_stocks:
            # 分析财务状况改善情况
            recent_financials = financial_data[financial_data['stock'] == stock].tail(4)

            if len(recent_financials) >= 2:
                # 检查盈利改善
                profit_improvement = recent_financials['净利润'].iloc[-1] > recent_financials['净利润'].iloc[-2]
                revenue_growth = recent_financials['营业收入'].pct_change().iloc[-1]

                # 分析摘帽概率
                delisting_risk = self.assess_delisting_risk(recent_financials)

                st_analysis[stock] = {
                    'profit_improvement': profit_improvement,
                    'revenue_growth': revenue_growth,
                    'delisting_risk': delisting_risk,
                    'investment_score': self.calculate_st_score(profit_improvement, revenue_growth, delisting_risk)
                }

        return st_analysis

    def assess_delisting_risk(self, financials):
        """评估退市风险"""
        latest_data = financials.iloc[-1]

        risk_factors = 0

        # 连续亏损
        if financials['净利润'].tail(2).sum() < 0:
            risk_factors += 1

        # 资不抵债
        if latest_data['净资产'] < 0:
            risk_factors += 2

        # 营业收入过低
        if latest_data['营业收入'] < 100000000:  # 1亿
            risk_factors += 1

        return min(risk_factors / 4, 1.0)  # 标准化为0-1

    def calculate_st_score(self, profit_improvement, revenue_growth, delisting_risk):
        """计算ST股票投资评分"""
        score = 0

        if profit_improvement:
            score += 0.4

        if revenue_growth > 0:
            score += 0.3

        score -= delisting_risk * 0.5

        return max(0, score)

    def theme_investment_strategy(self, theme_data, price_data):
        """主题投资策略"""
        theme_performance = {}

        for theme, stocks in theme_data.items():
            theme_prices = price_data[stocks]
            theme_returns = theme_prices.pct_change()

            # 计算主题指数
            theme_index = theme_prices.mean(axis=1)
            theme_momentum = theme_index.pct_change(20)  # 20日动量

            # 分析主题内股票表现分化
            latest_returns = theme_returns.iloc[-20:].mean()  # 近20日平均收益
            performance_dispersion = latest_returns.std()

            theme_performance[theme] = {
                'momentum_20d': theme_momentum.iloc[-1],
                'performance_dispersion': performance_dispersion,
                'best_performer': latest_returns.idxmax(),
                'worst_performer': latest_returns.idxmin(),
                'theme_strength': self.calculate_theme_strength(theme_momentum, performance_dispersion)
            }

        return theme_performance

    def calculate_theme_strength(self, momentum, dispersion):
        """计算主题强度"""
        # 动量强且分化小的主题得分高
        momentum_score = max(0, momentum) * 0.7
        concentration_score = (1 - min(dispersion, 0.1) / 0.1) * 0.3

        return momentum_score + concentration_score

    def quarterly_report_strategy(self, earnings_calendar, price_data):
        """季报披露策略"""
        earnings_trades = []

        for date, stocks_reporting in earnings_calendar.items():
            for stock in stocks_reporting:
                # 分析季报前后的价格走势
                pre_earnings = price_data.loc[date - pd.Timedelta(days=5):date - pd.Timedelta(days=1), stock]
                post_earnings = price_data.loc[date:date + pd.Timedelta(days=5), stock]

                if len(pre_earnings) > 0 and len(post_earnings) > 0:
                    pre_trend = (pre_earnings.iloc[-1] - pre_earnings.iloc[0]) / pre_earnings.iloc[0]
                    post_performance = (post_earnings.iloc[-1] - post_earnings.iloc[0]) / post_earnings.iloc[0]

                    earnings_trades.append({
                        'date': date,
                        'stock': stock,
                        'pre_earnings_trend': pre_trend,
                        'post_earnings_performance': post_performance,
                        'strategy_signal': self.generate_earnings_signal(pre_trend)
                    })

        return pd.DataFrame(earnings_trades)

    def generate_earnings_signal(self, pre_trend):
        """生成季报交易信号"""
        if pre_trend > 0.05:  # 季报前上涨超过5%
            return 'sell'  # 获利了结
        elif pre_trend < -0.05:  # 季报前下跌超过5%
            return 'buy'   # 逢低买入
        else:
            return 'hold'

# 使用示例
a_share_strategies = AShareSpecificStrategies()
# ipo_analysis = a_share_strategies.new_stock_ipo_strategy(ipo_data, price_data)
# st_analysis = a_share_strategies.st_stock_strategy(st_stocks, financial_data, price_data)
# theme_analysis = a_share_strategies.theme_investment_strategy(theme_data, price_data)
```

## 4.10 股票投资组合构建

### 4.10.1 组合优化理论

```python
import cvxpy as cp
from scipy.optimize import minimize

class PortfolioOptimization:
    """投资组合优化"""

    def __init__(self):
        self.expected_returns = None
        self.covariance_matrix = None

    def calculate_expected_returns(self, returns, method='historical_mean'):
        """计算预期收益率"""
        if method == 'historical_mean':
            return returns.mean()
        elif method == 'capm':
            # 使用CAPM模型估计
            pass
        elif method == 'black_litterman':
            # Black-Litterman模型
            pass

    def mean_variance_optimization(self, expected_returns, cov_matrix, risk_aversion=1):
        """均值-方差优化"""
        n_assets = len(expected_returns)

        # 定义变量
        weights = cp.Variable(n_assets)

        # 目标函数：最大化效用 = 预期收益 - 风险厌恶系数 * 风险
        portfolio_return = expected_returns.values @ weights
        portfolio_risk = cp.quad_form(weights, cov_matrix.values)
        utility = portfolio_return - 0.5 * risk_aversion * portfolio_risk

        # 约束条件
        constraints = [
            cp.sum(weights) == 1,  # 权重和为1
            weights >= 0           # 不允许做空
        ]

        # 求解
        problem = cp.Problem(cp.Maximize(utility), constraints)
        problem.solve()

        if weights.value is not None:
            optimal_weights = pd.Series(weights.value, index=expected_returns.index)
            return optimal_weights
        else:
            return None

    def min_variance_portfolio(self, cov_matrix):
        """最小方差组合"""
        n_assets = len(cov_matrix)

        weights = cp.Variable(n_assets)
        objective = cp.Minimize(cp.quad_form(weights, cov_matrix.values))

        constraints = [
            cp.sum(weights) == 1,
            weights >= 0
        ]

        problem = cp.Problem(objective, constraints)
        problem.solve()

        if weights.value is not None:
            return pd.Series(weights.value, index=cov_matrix.index)
        else:
            return None

    def risk_parity_portfolio(self, cov_matrix):
        """风险平价组合"""
        def risk_budget_objective(weights, cov_matrix):
            portfolio_variance = np.dot(weights, np.dot(cov_matrix, weights))
            marginal_contrib = np.dot(cov_matrix, weights)
            contrib = weights * marginal_contrib

            # 风险贡献应该相等
            risk_diffs = []
            for i in range(len(weights)):
                for j in range(i+1, len(weights)):
                    risk_diffs.append((contrib[i] - contrib[j])**2)

            return sum(risk_diffs)

        n_assets = len(cov_matrix)
        initial_weights = np.array([1/n_assets] * n_assets)

        constraints = [
            {'type': 'eq', 'fun': lambda x: np.sum(x) - 1},  # 权重和为1
        ]

        bounds = [(0, 1) for _ in range(n_assets)]  # 权重在0-1之间

        result = minimize(
            risk_budget_objective,
            initial_weights,
            args=(cov_matrix.values,),
            method='SLSQP',
            bounds=bounds,
            constraints=constraints
        )

        if result.success:
            return pd.Series(result.x, index=cov_matrix.index)
        else:
            return None

    def black_litterman_optimization(self, market_caps, cov_matrix, views=None, tau=0.025):
        """Black-Litterman模型"""
        # 市场均衡收益率
        market_weights = market_caps / market_caps.sum()
        pi = tau * np.dot(cov_matrix, market_weights)

        if views is None:
            # 没有观点时，返回市场权重
            return market_weights

        # 有观点时，结合观点和市场均衡
        P = views['P']  # 观点矩阵
        Q = views['Q']  # 观点收益率
        Omega = views['Omega']  # 观点不确定性

        # Black-Litterman公式
        M1 = np.linalg.inv(tau * cov_matrix)
        M2 = np.dot(P.T, np.dot(np.linalg.inv(Omega), P))
        M3 = np.dot(np.linalg.inv(tau * cov_matrix), pi)
        M4 = np.dot(P.T, np.dot(np.linalg.inv(Omega), Q))

        mu_bl = np.dot(np.linalg.inv(M1 + M2), M3 + M4)
        cov_bl = np.linalg.inv(M1 + M2)

        # 使用新的期望收益和协方差矩阵进行优化
        return self.mean_variance_optimization(
            pd.Series(mu_bl, index=cov_matrix.index),
            pd.DataFrame(cov_bl, index=cov_matrix.index, columns=cov_matrix.columns)
        )

# 使用示例
portfolio_optimizer = PortfolioOptimization()
# expected_returns = portfolio_optimizer.calculate_expected_returns(stock_returns)
# optimal_weights = portfolio_optimizer.mean_variance_optimization(expected_returns, cov_matrix)
# min_var_weights = portfolio_optimizer.min_variance_portfolio(cov_matrix)
# risk_parity_weights = portfolio_optimizer.risk_parity_portfolio(cov_matrix)
```

## 4.11 本章总结

在本章下篇中，我们深入学习了股票投资的高级分析方法：

### 🎯 核心要点回顾

1. **量化选股策略**

   - 财务因子和技术因子的构建
   - 多因子选股模型的实现
   - 行业中性选股方法

2. **多因子模型**

   - Fama-French 三因子模型
   - 动态因子选择和 IC 分析
   - 因子衰减和组合构建

3. **事件驱动策略**

   - 盈利公告策略和盈利后漂移
   - 并购套利策略
   - 事件概率评估方法

4. **A 股市场特色**

   - A 股独特的交易制度和市场特征
   - 涨跌停、T+1、主题投资等特色策略
   - 政策影响和行业轮动分析

5. **投资组合优化**
   - 均值-方差优化
   - 风险平价组合
   - Black-Litterman 模型

### 📈 实践应用指南

1. **策略开发流程**

   - 因子研究 → 模型构建 → 回测验证 → 实盘应用

2. **风险控制要点**

   - 因子暴露控制
   - 行业和风格中性
   - 仓位管理和止损

3. **A 股投资建议**
   - 关注政策导向和主题机会
   - 重视季报和年报时间窗口
   - 利用市场情绪和资金流向

## 📚 延伸学习

在[第五章：期货与衍生品交易（上篇）](05_futures_derivatives_trading.md)中，我们将学习期货合约基础、期货定价理论、套期保值策略以及期货量化策略。

## 💡 实践练习

1. **多因子模型实践**：构建一个包含价值、成长、质量、动量四大类因子的选股模型

2. **事件驱动策略**：开发一个基于盈利公告的交易策略并进行回测

3. **A 股特色策略**：分析某个热门主题的投资机会和风险

4. **组合优化**：使用不同优化方法构建投资组合并比较表现

5. **因子研究**：研究某个新因子的有效性和稳定性

---

**建议学习时间：** 5-6 天  
**前置章节：** [第四章：股票投资与分析（上篇）](04_stock_investment_analysis.md)  
**下一章：** [期货与衍生品交易（上篇）](05_futures_derivatives_trading.md)
