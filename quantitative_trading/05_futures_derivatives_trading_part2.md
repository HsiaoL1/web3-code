# 第五章：期货与衍生品交易（下篇）

## 📖 章节概述

本章下篇将深入探讨期权基础与策略、高级量化套利策略、高频交易入门以及期货市场的风险管理技术。这些内容将帮助您从基础期货交易进阶到复杂衍生品策略和专业风险管理。

## 🎯 学习目标

- 掌握期权定价理论和交易策略
- 学会高级量化套利策略的设计与实施
- 了解高频交易的基本原理和技术
- 建立完善的风险管理体系
- 掌握衍生品组合的风险度量方法
- 学会期权波动率交易策略

## 5.6 期权基础与策略

### 5.6.1 期权定价模型

```python
import numpy as np
import pandas as pd
from scipy.stats import norm
from scipy.optimize import minimize_scalar
import matplotlib.pyplot as plt
import warnings
warnings.filterwarnings('ignore')

class OptionPricingModels:
    """期权定价模型"""

    def __init__(self):
        self.pricing_models = {}

    def black_scholes_call(self, S, K, T, r, sigma):
        """Black-Scholes看涨期权定价"""
        d1 = (np.log(S / K) + (r + 0.5 * sigma ** 2) * T) / (sigma * np.sqrt(T))
        d2 = d1 - sigma * np.sqrt(T)

        call_price = S * norm.cdf(d1) - K * np.exp(-r * T) * norm.cdf(d2)
        return call_price

    def black_scholes_put(self, S, K, T, r, sigma):
        """Black-Scholes看跌期权定价"""
        d1 = (np.log(S / K) + (r + 0.5 * sigma ** 2) * T) / (sigma * np.sqrt(T))
        d2 = d1 - sigma * np.sqrt(T)

        put_price = K * np.exp(-r * T) * norm.cdf(-d2) - S * norm.cdf(-d1)
        return put_price

    def calculate_greeks(self, S, K, T, r, sigma, option_type='call'):
        """计算期权希腊字母"""
        d1 = (np.log(S / K) + (r + 0.5 * sigma ** 2) * T) / (sigma * np.sqrt(T))
        d2 = d1 - sigma * np.sqrt(T)

        # Delta
        if option_type == 'call':
            delta = norm.cdf(d1)
        else:  # put
            delta = norm.cdf(d1) - 1

        # Gamma
        gamma = norm.pdf(d1) / (S * sigma * np.sqrt(T))

        # Theta
        if option_type == 'call':
            theta = (-S * norm.pdf(d1) * sigma / (2 * np.sqrt(T))
                    - r * K * np.exp(-r * T) * norm.cdf(d2)) / 365
        else:  # put
            theta = (-S * norm.pdf(d1) * sigma / (2 * np.sqrt(T))
                    + r * K * np.exp(-r * T) * norm.cdf(-d2)) / 365

        # Vega
        vega = S * norm.pdf(d1) * np.sqrt(T) / 100

        # Rho
        if option_type == 'call':
            rho = K * T * np.exp(-r * T) * norm.cdf(d2) / 100
        else:  # put
            rho = -K * T * np.exp(-r * T) * norm.cdf(-d2) / 100

        return {
            'delta': delta,
            'gamma': gamma,
            'theta': theta,
            'vega': vega,
            'rho': rho
        }

    def implied_volatility(self, market_price, S, K, T, r, option_type='call'):
        """计算隐含波动率"""
        def objective_function(sigma):
            if option_type == 'call':
                theoretical_price = self.black_scholes_call(S, K, T, r, sigma)
            else:
                theoretical_price = self.black_scholes_put(S, K, T, r, sigma)

            return abs(theoretical_price - market_price)

        result = minimize_scalar(objective_function, bounds=(0.01, 3.0), method='bounded')
        return result.x if result.success else None

    def binomial_option_pricing(self, S, K, T, r, sigma, n_steps=100, option_type='call', style='european'):
        """二叉树期权定价模型"""
        dt = T / n_steps
        u = np.exp(sigma * np.sqrt(dt))
        d = 1 / u
        p = (np.exp(r * dt) - d) / (u - d)

        # 构建价格树
        price_tree = np.zeros((n_steps + 1, n_steps + 1))
        price_tree[0, 0] = S

        for i in range(1, n_steps + 1):
            price_tree[i, 0] = price_tree[i-1, 0] * u
            for j in range(1, i + 1):
                price_tree[i, j] = price_tree[i-1, j-1] * d

        # 计算期权价值
        option_tree = np.zeros((n_steps + 1, n_steps + 1))

        # 到期日价值
        for j in range(n_steps + 1):
            if option_type == 'call':
                option_tree[n_steps, j] = max(0, price_tree[n_steps, j] - K)
            else:  # put
                option_tree[n_steps, j] = max(0, K - price_tree[n_steps, j])

        # 向前递推
        for i in range(n_steps - 1, -1, -1):
            for j in range(i + 1):
                option_value = np.exp(-r * dt) * (p * option_tree[i+1, j] + (1-p) * option_tree[i+1, j+1])

                if style == 'american':
                    if option_type == 'call':
                        early_exercise = max(0, price_tree[i, j] - K)
                    else:
                        early_exercise = max(0, K - price_tree[i, j])
                    option_tree[i, j] = max(option_value, early_exercise)
                else:
                    option_tree[i, j] = option_value

        return option_tree[0, 0]

class OptionTradingStrategies:
    """期权交易策略"""

    def __init__(self):
        self.strategies = {}
        self.pricing_model = OptionPricingModels()

    def covered_call_strategy(self, stock_price, strike_price, call_premium,
                             stock_quantity=100):
        """备兑看涨期权策略"""
        # 持有股票+卖出看涨期权
        strategy_cost = stock_price * stock_quantity - call_premium * stock_quantity

        # 计算不同股价下的损益
        price_range = np.linspace(stock_price * 0.8, stock_price * 1.3, 50)
        payoffs = []

        for price in price_range:
            stock_payoff = (price - stock_price) * stock_quantity
            option_payoff = call_premium * stock_quantity - max(0, price - strike_price) * stock_quantity
            total_payoff = stock_payoff + option_payoff
            payoffs.append(total_payoff)

        return {
            'strategy_cost': strategy_cost,
            'max_profit': (strike_price - stock_price + call_premium) * stock_quantity,
            'breakeven': stock_price - call_premium,
            'price_range': price_range,
            'payoffs': payoffs
        }

    def protective_put_strategy(self, stock_price, strike_price, put_premium,
                               stock_quantity=100):
        """保护性看跌期权策略"""
        strategy_cost = stock_price * stock_quantity + put_premium * stock_quantity

        price_range = np.linspace(stock_price * 0.6, stock_price * 1.4, 50)
        payoffs = []

        for price in price_range:
            stock_payoff = (price - stock_price) * stock_quantity
            option_payoff = max(0, strike_price - price) * stock_quantity - put_premium * stock_quantity
            total_payoff = stock_payoff + option_payoff
            payoffs.append(total_payoff)

        return {
            'strategy_cost': strategy_cost,
            'max_loss': (stock_price + put_premium - strike_price) * stock_quantity,
            'breakeven': stock_price + put_premium,
            'price_range': price_range,
            'payoffs': payoffs
        }

    def straddle_strategy(self, strike_price, call_premium, put_premium, position_type='long'):
        """跨式策略"""
        if position_type == 'long':
            strategy_cost = call_premium + put_premium
            multiplier = 1
        else:  # short straddle
            strategy_cost = -(call_premium + put_premium)
            multiplier = -1

        price_range = np.linspace(strike_price * 0.7, strike_price * 1.3, 50)
        payoffs = []

        for price in price_range:
            call_payoff = max(0, price - strike_price) - call_premium
            put_payoff = max(0, strike_price - price) - put_premium
            total_payoff = (call_payoff + put_payoff) * multiplier
            payoffs.append(total_payoff)

        breakeven_up = strike_price + call_premium + put_premium
        breakeven_down = strike_price - call_premium - put_premium

        return {
            'strategy_cost': strategy_cost,
            'breakeven_up': breakeven_up,
            'breakeven_down': breakeven_down,
            'price_range': price_range,
            'payoffs': payoffs
        }

    def iron_condor_strategy(self, low_put_strike, high_put_strike,
                           low_call_strike, high_call_strike,
                           put_spreads_premium, call_spreads_premium):
        """铁鹰策略"""
        # 卖出跨式+买入保护
        net_premium = put_spreads_premium + call_spreads_premium

        price_range = np.linspace(low_put_strike * 0.9, high_call_strike * 1.1, 100)
        payoffs = []

        for price in price_range:
            # Put spread payoff
            if price <= low_put_strike:
                put_payoff = -(high_put_strike - low_put_strike) + put_spreads_premium
            elif price <= high_put_strike:
                put_payoff = -(price - low_put_strike) + put_spreads_premium
            else:
                put_payoff = put_spreads_premium

            # Call spread payoff
            if price >= high_call_strike:
                call_payoff = -(high_call_strike - low_call_strike) + call_spreads_premium
            elif price >= low_call_strike:
                call_payoff = -(price - low_call_strike) + call_spreads_premium
            else:
                call_payoff = call_spreads_premium

            total_payoff = put_payoff + call_payoff
            payoffs.append(total_payoff)

        return {
            'net_premium': net_premium,
            'max_profit': net_premium,
            'max_loss': max(high_put_strike - low_put_strike, high_call_strike - low_call_strike) - net_premium,
            'price_range': price_range,
            'payoffs': payoffs
        }

# 使用示例
option_pricing = OptionPricingModels()
option_strategies = OptionTradingStrategies()

# Black-Scholes定价示例
# call_price = option_pricing.black_scholes_call(S=100, K=105, T=0.25, r=0.05, sigma=0.2)
# greeks = option_pricing.calculate_greeks(S=100, K=105, T=0.25, r=0.05, sigma=0.2)

# 备兑看涨策略示例
# covered_call = option_strategies.covered_call_strategy(
#     stock_price=100, strike_price=105, call_premium=3, stock_quantity=100
# )
```

### 5.6.2 波动率交易策略

```python
class VolatilityTradingStrategies:
    """波动率交易策略"""

    def __init__(self):
        self.vol_strategies = {}

    def calculate_historical_volatility(self, price_data, window=30):
        """计算历史波动率"""
        returns = price_data.pct_change().dropna()

        rolling_vol = returns.rolling(window=window).std() * np.sqrt(252)

        return {
            'daily_returns': returns,
            'historical_volatility': rolling_vol,
            'current_vol': rolling_vol.iloc[-1],
            'vol_percentile': (rolling_vol.iloc[-1] > rolling_vol).mean()
        }

    def volatility_cone_analysis(self, price_data, periods=[10, 20, 30, 60, 90]):
        """波动率锥分析"""
        returns = price_data.pct_change().dropna()
        vol_cone = {}

        for period in periods:
            rolling_vol = returns.rolling(window=period).std() * np.sqrt(252)

            vol_cone[f'{period}d'] = {
                'current': rolling_vol.iloc[-1],
                'min': rolling_vol.min(),
                'max': rolling_vol.max(),
                'median': rolling_vol.median(),
                'percentile_25': rolling_vol.quantile(0.25),
                'percentile_75': rolling_vol.quantile(0.75)
            }

        return vol_cone

    def volatility_surface_analysis(self, options_data):
        """波动率曲面分析"""
        # options_data应包含不同行权价和到期日的期权数据
        vol_surface = {}

        for expiry in options_data['expiry'].unique():
            expiry_data = options_data[options_data['expiry'] == expiry]

            vol_surface[expiry] = {}
            for strike in expiry_data['strike'].unique():
                option_data = expiry_data[expiry_data['strike'] == strike]

                if len(option_data) > 0:
                    vol_surface[expiry][strike] = option_data['implied_vol'].iloc[0]

        return vol_surface

    def volatility_arbitrage_strategy(self, historical_vol, implied_vol,
                                    vol_threshold=0.05):
        """波动率套利策略"""
        vol_difference = implied_vol - historical_vol
        relative_difference = vol_difference / historical_vol

        if abs(relative_difference) > vol_threshold:
            if vol_difference > 0:
                # 隐含波动率高估，卖出波动率
                strategy = {
                    'action': 'sell_volatility',
                    'recommendation': 'sell_straddle_or_strangle',
                    'expected_direction': 'volatility_decrease',
                    'confidence': min(abs(relative_difference) / vol_threshold, 2.0)
                }
            else:
                # 隐含波动率低估，买入波动率
                strategy = {
                    'action': 'buy_volatility',
                    'recommendation': 'buy_straddle_or_strangle',
                    'expected_direction': 'volatility_increase',
                    'confidence': min(abs(relative_difference) / vol_threshold, 2.0)
                }
        else:
            strategy = {
                'action': 'no_action',
                'recommendation': 'volatility_fairly_priced'
            }

        return strategy

    def delta_neutral_hedging(self, option_positions, current_stock_price,
                            stock_delta=1.0):
        """Delta中性对冲"""
        total_delta = 0

        for position in option_positions:
            option_delta = position['delta']
            quantity = position['quantity']
            total_delta += option_delta * quantity

        # 计算需要的股票数量来对冲Delta
        hedge_stock_quantity = -total_delta / stock_delta

        return {
            'total_portfolio_delta': total_delta,
            'hedge_stock_quantity': hedge_stock_quantity,
            'hedge_cost': hedge_stock_quantity * current_stock_price,
            'portfolio_delta_after_hedge': 0
        }

    def gamma_scalping_strategy(self, position_gamma, stock_price_change,
                              hedge_frequency='daily'):
        """Gamma套利策略"""
        # Gamma scalping利用Gamma效应获利
        gamma_pnl = 0.5 * position_gamma * (stock_price_change ** 2)

        rebalancing_frequency_map = {
            'intraday': 0.25,
            'daily': 1.0,
            'weekly': 7.0
        }

        frequency_factor = rebalancing_frequency_map.get(hedge_frequency, 1.0)

        return {
            'gamma_pnl': gamma_pnl,
            'optimal_rebalance_threshold': np.sqrt(frequency_factor * 0.01),  # 简化计算
            'strategy_effectiveness': gamma_pnl / abs(stock_price_change) if stock_price_change != 0 else 0
        }

# 使用示例
vol_strategies = VolatilityTradingStrategies()

# 历史波动率计算示例
# hist_vol = vol_strategies.calculate_historical_volatility(price_data)

# 波动率锥分析示例
# vol_cone = vol_strategies.volatility_cone_analysis(price_data)
```

## 5.7 量化套利策略

### 5.7.1 统计套利策略

```python
class StatisticalArbitrageStrategies:
    """统计套利策略"""

    def __init__(self):
        self.cointegration_models = {}

    def pairs_trading_strategy(self, price1, price2, lookback_window=60,
                              entry_threshold=2, exit_threshold=0.5):
        """配对交易策略"""
        # 计算价格比率
        ratio = price1 / price2
        rolling_mean = ratio.rolling(window=lookback_window).mean()
        rolling_std = ratio.rolling(window=lookback_window).std()

        # 计算Z-score
        z_score = (ratio - rolling_mean) / rolling_std

        # 生成交易信号
        signals = pd.DataFrame(index=ratio.index)
        signals['ratio'] = ratio
        signals['z_score'] = z_score
        signals['position'] = 0

        # 开仓信号
        signals.loc[z_score > entry_threshold, 'position'] = -1  # 卖出价差
        signals.loc[z_score < -entry_threshold, 'position'] = 1   # 买入价差

        # 平仓信号
        signals.loc[abs(z_score) < exit_threshold, 'position'] = 0

        # 计算持仓变化
        signals['position_change'] = signals['position'].diff()

        return signals

    def cointegration_analysis(self, price_series1, price_series2):
        """协整分析"""
        from statsmodels.tsa.stattools import coint
        from statsmodels.api import OLS

        # 协整检验
        coint_stat, p_value, critical_values = coint(price_series1, price_series2)

        # 计算协整关系
        ols_model = OLS(price_series1, price_series2).fit()
        hedge_ratio = ols_model.params[0]

        # 计算残差
        spread = price_series1 - hedge_ratio * price_series2

        return {
            'cointegration_statistic': coint_stat,
            'p_value': p_value,
            'critical_values': critical_values,
            'hedge_ratio': hedge_ratio,
            'spread': spread,
            'is_cointegrated': p_value < 0.05
        }

    def mean_reversion_strategy(self, spread_series, entry_z_score=2,
                               exit_z_score=0, stop_loss_z_score=3):
        """均值回归策略"""
        rolling_mean = spread_series.rolling(window=60).mean()
        rolling_std = spread_series.rolling(window=60).std()

        z_score = (spread_series - rolling_mean) / rolling_std

        signals = pd.DataFrame(index=spread_series.index)
        signals['spread'] = spread_series
        signals['z_score'] = z_score
        signals['position'] = 0

        # 均值回归信号
        signals.loc[z_score > entry_z_score, 'position'] = -1   # 卖出
        signals.loc[z_score < -entry_z_score, 'position'] = 1    # 买入
        signals.loc[abs(z_score) < exit_z_score, 'position'] = 0  # 平仓

        # 止损信号
        signals.loc[z_score > stop_loss_z_score, 'position'] = 0
        signals.loc[z_score < -stop_loss_z_score, 'position'] = 0

        return signals

    def basket_arbitrage_strategy(self, basket_weights, individual_prices,
                                 basket_price, arbitrage_threshold=0.02):
        """篮子套利策略"""
        # 计算理论篮子价格
        theoretical_basket_price = sum(basket_weights[i] * individual_prices.iloc[:, i]
                                     for i in range(len(basket_weights)))

        # 计算价差
        arbitrage_spread = basket_price - theoretical_basket_price
        relative_spread = arbitrage_spread / theoretical_basket_price

        signals = pd.DataFrame(index=basket_price.index)
        signals['basket_price'] = basket_price
        signals['theoretical_price'] = theoretical_basket_price
        signals['spread'] = arbitrage_spread
        signals['relative_spread'] = relative_spread
        signals['position'] = 0

        # 套利信号
        signals.loc[relative_spread > arbitrage_threshold, 'position'] = -1  # 卖篮子买个股
        signals.loc[relative_spread < -arbitrage_threshold, 'position'] = 1   # 买篮子卖个股

        return signals

class HighFrequencyStrategies:
    """高频交易策略"""

    def __init__(self):
        self.hft_strategies = {}

    def market_making_strategy(self, bid_ask_data, inventory_limit=1000,
                              spread_multiplier=1.2):
        """做市策略"""
        signals = pd.DataFrame(index=bid_ask_data.index)
        signals['bid'] = bid_ask_data['bid']
        signals['ask'] = bid_ask_data['ask']
        signals['spread'] = bid_ask_data['ask'] - bid_ask_data['bid']
        signals['mid_price'] = (bid_ask_data['bid'] + bid_ask_data['ask']) / 2

        # 计算目标买卖价
        target_spread = signals['spread'] * spread_multiplier
        signals['target_bid'] = signals['mid_price'] - target_spread / 2
        signals['target_ask'] = signals['mid_price'] + target_spread / 2

        # 库存管理
        signals['inventory'] = 0  # 简化处理，实际需要跟踪成交
        signals['skew'] = 0  # 根据库存调整报价偏移

        return signals

    def momentum_ignition_detection(self, price_data, volume_data,
                                   price_threshold=0.005, volume_threshold=2):
        """动量点火检测"""
        returns = price_data.pct_change()
        volume_ratio = volume_data / volume_data.rolling(20).mean()

        # 检测异常价格和成交量
        price_anomaly = abs(returns) > price_threshold
        volume_anomaly = volume_ratio > volume_threshold

        momentum_signals = price_anomaly & volume_anomaly

        return {
            'momentum_ignition_signals': momentum_signals,
            'signal_frequency': momentum_signals.sum() / len(momentum_signals),
            'average_return_on_signal': returns[momentum_signals].mean()
        }

    def latency_arbitrage_opportunity(self, exchange1_prices, exchange2_prices,
                                    latency_ms=10, transaction_cost=0.001):
        """延迟套利机会"""
        # 模拟延迟
        lagged_prices = exchange1_prices.shift(1)  # 简化处理延迟

        price_difference = exchange2_prices - lagged_prices
        arbitrage_opportunity = abs(price_difference) > transaction_cost

        potential_profit = abs(price_difference) - transaction_cost
        profitable_opportunities = potential_profit > 0

        return {
            'arbitrage_opportunities': arbitrage_opportunity,
            'profitable_opportunities': profitable_opportunities,
            'average_profit': potential_profit[profitable_opportunities].mean(),
            'opportunity_frequency': profitable_opportunities.sum() / len(profitable_opportunities)
        }

    def order_flow_imbalance_strategy(self, order_book_data):
        """订单流失衡策略"""
        # 计算订单流失衡
        bid_volume = order_book_data['bid_volume']
        ask_volume = order_book_data['ask_volume']

        order_imbalance = (bid_volume - ask_volume) / (bid_volume + ask_volume)

        # 根据订单流失衡预测价格方向
        signals = pd.DataFrame(index=order_book_data.index)
        signals['order_imbalance'] = order_imbalance
        signals['predicted_direction'] = np.sign(order_imbalance)

        # 生成交易信号
        imbalance_threshold = 0.3
        signals['trade_signal'] = 0
        signals.loc[order_imbalance > imbalance_threshold, 'trade_signal'] = 1
        signals.loc[order_imbalance < -imbalance_threshold, 'trade_signal'] = -1

        return signals

# 使用示例
stat_arb = StatisticalArbitrageStrategies()
hft_strategies = HighFrequencyStrategies()

# 配对交易示例
# pairs_signals = stat_arb.pairs_trading_strategy(price1, price2)

# 协整分析示例
# coint_result = stat_arb.cointegration_analysis(price_series1, price_series2)
```

### 5.7.2 跨市场套利策略

```python
class CrossMarketArbitrage:
    """跨市场套利策略"""

    def __init__(self):
        self.arbitrage_opportunities = {}

    def index_arbitrage_strategy(self, index_price, futures_price,
                               component_prices, component_weights,
                               transaction_cost=0.002):
        """指数套利策略"""
        # 计算理论指数价格
        theoretical_index = sum(component_weights[i] * component_prices.iloc[:, i]
                               for i in range(len(component_weights)))

        # 计算期现价差
        basis = futures_price - index_price
        theoretical_basis = futures_price - theoretical_index

        signals = pd.DataFrame(index=index_price.index)
        signals['index_price'] = index_price
        signals['futures_price'] = futures_price
        signals['theoretical_index'] = theoretical_index
        signals['basis'] = basis
        signals['theoretical_basis'] = theoretical_basis

        # 套利信号
        arbitrage_threshold = transaction_cost * 2

        # 正向套利：卖期货买现货
        signals['positive_arbitrage'] = theoretical_basis > arbitrage_threshold

        # 反向套利：买期货卖现货
        signals['negative_arbitrage'] = theoretical_basis < -arbitrage_threshold

        # 计算预期收益
        signals['expected_profit'] = abs(theoretical_basis) - transaction_cost
        signals['profitable'] = signals['expected_profit'] > 0

        return signals

    def currency_arbitrage_strategy(self, usd_cny_rate, usd_eur_rate,
                                  eur_cny_rate, transaction_cost=0.001):
        """货币套利策略"""
        # 计算理论汇率
        theoretical_eur_cny = usd_cny_rate / usd_eur_rate

        # 计算套利价差
        arbitrage_spread = eur_cny_rate - theoretical_eur_cny
        relative_spread = arbitrage_spread / theoretical_eur_cny

        signals = pd.DataFrame(index=usd_cny_rate.index)
        signals['usd_cny'] = usd_cny_rate
        signals['usd_eur'] = usd_eur_rate
        signals['eur_cny'] = eur_cny_rate
        signals['theoretical_eur_cny'] = theoretical_eur_cny
        signals['arbitrage_spread'] = arbitrage_spread
        signals['relative_spread'] = relative_spread

        # 套利机会
        min_profit_threshold = transaction_cost * 3  # 考虑三次交易的成本

        signals['arbitrage_opportunity'] = abs(relative_spread) > min_profit_threshold
        signals['arbitrage_direction'] = np.where(relative_spread > 0, 'sell_eur_cny', 'buy_eur_cny')

        return signals

    def etf_arbitrage_strategy(self, etf_price, nav_price, creation_redemption_cost=0.003):
        """ETF套利策略"""
        # 计算ETF折溢价
        premium_discount = (etf_price - nav_price) / nav_price

        signals = pd.DataFrame(index=etf_price.index)
        signals['etf_price'] = etf_price
        signals['nav_price'] = nav_price
        signals['premium_discount'] = premium_discount

        # 套利信号
        arbitrage_threshold = creation_redemption_cost

        # 申购套利：ETF溢价过高
        signals['creation_arbitrage'] = premium_discount > arbitrage_threshold

        # 赎回套利：ETF折价过高
        signals['redemption_arbitrage'] = premium_discount < -arbitrage_threshold

        # 计算预期收益
        signals['expected_profit'] = abs(premium_discount) - creation_redemption_cost
        signals['profitable'] = signals['expected_profit'] > 0

        return signals

    def commodity_arbitrage_strategy(self, spot_price, futures_price,
                                   storage_cost, financing_cost, convenience_yield):
        """商品套利策略"""
        # 计算理论期货价格
        net_cost_of_carry = financing_cost + storage_cost - convenience_yield
        theoretical_futures = spot_price * np.exp(net_cost_of_carry * 0.25)  # 假设3个月到期

        # 计算套利价差
        arbitrage_spread = futures_price - theoretical_futures

        signals = pd.DataFrame(index=spot_price.index)
        signals['spot_price'] = spot_price
        signals['futures_price'] = futures_price
        signals['theoretical_futures'] = theoretical_futures
        signals['arbitrage_spread'] = arbitrage_spread
        signals['relative_spread'] = arbitrage_spread / theoretical_futures

        # 套利策略
        transaction_cost = 0.005  # 商品交易成本较高

        # 现货套利：期货高估
        signals['cash_and_carry'] = arbitrage_spread > transaction_cost

        # 反向套利：期货低估
        signals['reverse_cash_carry'] = arbitrage_spread < -transaction_cost

        return signals

class RiskManagementSystem:
    """风险管理系统"""

    def __init__(self):
        self.risk_metrics = {}

    def calculate_var(self, returns, confidence_level=0.05, method='historical'):
        """计算风险价值"""
        if method == 'historical':
            var = np.percentile(returns, confidence_level * 100)
        elif method == 'parametric':
            mean_return = np.mean(returns)
            std_return = np.std(returns)
            var = mean_return + norm.ppf(confidence_level) * std_return

        return {
            'var': var,
            'confidence_level': confidence_level,
            'method': method
        }

    def calculate_cvar(self, returns, confidence_level=0.05):
        """计算条件风险价值"""
        var = np.percentile(returns, confidence_level * 100)
        cvar = returns[returns <= var].mean()

        return {
            'cvar': cvar,
            'var': var,
            'tail_expectation': cvar
        }

    def portfolio_risk_metrics(self, portfolio_returns):
        """投资组合风险指标"""
        metrics = {
            'volatility': np.std(portfolio_returns) * np.sqrt(252),
            'sharpe_ratio': np.mean(portfolio_returns) / np.std(portfolio_returns) * np.sqrt(252),
            'max_drawdown': self._calculate_max_drawdown(portfolio_returns),
            'var_5pct': self.calculate_var(portfolio_returns, 0.05)['var'],
            'cvar_5pct': self.calculate_cvar(portfolio_returns, 0.05)['cvar']
        }

        return metrics

    def _calculate_max_drawdown(self, returns):
        """计算最大回撤"""
        cumulative_returns = (1 + returns).cumprod()
        running_max = cumulative_returns.expanding().max()
        drawdown = (cumulative_returns - running_max) / running_max
        max_drawdown = drawdown.min()

        return max_drawdown

    def position_sizing(self, expected_return, volatility, kelly_fraction=0.25):
        """仓位管理"""
        # Kelly公式变种
        kelly_optimal = expected_return / (volatility ** 2)
        conservative_size = kelly_optimal * kelly_fraction

        return {
            'kelly_optimal': kelly_optimal,
            'conservative_size': conservative_size,
            'max_position': min(conservative_size, 0.1)  # 限制最大仓位10%
        }

    def stress_testing(self, portfolio_positions, stress_scenarios):
        """压力测试"""
        stress_results = {}

        for scenario_name, scenario_changes in stress_scenarios.items():
            portfolio_impact = 0

            for asset, position in portfolio_positions.items():
                if asset in scenario_changes:
                    asset_impact = position * scenario_changes[asset]
                    portfolio_impact += asset_impact

            stress_results[scenario_name] = {
                'portfolio_impact': portfolio_impact,
                'scenario_description': scenario_changes
            }

        return stress_results

# 使用示例
cross_market_arb = CrossMarketArbitrage()
risk_management = RiskManagementSystem()

# 指数套利示例
# index_arb_signals = cross_market_arb.index_arbitrage_strategy(
#     index_price, futures_price, component_prices, component_weights
# )

# 风险管理示例
# risk_metrics = risk_management.portfolio_risk_metrics(portfolio_returns)
```

## 5.8 本章总结

在本章下篇中，我们深入学习了期货与衍生品的高级交易策略：

### 🎯 核心要点回顾

1. **期权基础与策略**

   - Black-Scholes 定价模型和希腊字母
   - 各种期权交易策略组合
   - 波动率交易和套利策略

2. **量化套利策略**

   - 统计套利和配对交易
   - 跨市场套利机会识别
   - 高频交易策略基础

3. **风险管理技术**

   - VaR 和 CVaR 风险度量
   - 投资组合风险控制
   - 压力测试和情景分析

4. **高级交易技术**
   - 做市策略和订单流分析
   - 延迟套利和动量点火
   - Delta 中性对冲和 Gamma 套利

### 📈 实践应用指南

1. **策略实施要点**

   - 交易成本的精确计算
   - 流动性风险的评估
   - 执行滑点的控制

2. **风险控制建议**

   - 多层次风险管理体系
   - 实时监控和预警机制
   - 应急处置预案

3. **技术实现考虑**
   - 低延迟交易系统
   - 数据质量和完整性
   - 回测与实盘的差异

## 📚 延伸学习

在[第六章：外汇与商品交易](06_forex_commodity_trading.md)中，我们将学习外汇市场特点、汇率影响因素、商品期货分析以及跨市场套利策略。

## 💡 实践练习

1. **期权策略分析**：构建不同期权策略的损益图和风险收益分析

2. **波动率交易**：开发基于波动率预测的期权交易策略

3. **统计套利**：寻找协整股票对并实施配对交易策略

4. **风险管理**：建立完整的衍生品交易风险控制系统

5. **高频策略**：设计简单的做市策略并分析盈利能力

---

**建议学习时间：** 5-6 天  
**前置章节：** [第五章：期货与衍生品交易（上篇）](05_futures_derivatives_trading.md)  
**下一章：** [外汇与商品交易](06_forex_commodity_trading.md)
