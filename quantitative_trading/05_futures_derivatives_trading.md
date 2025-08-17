# 第五章：期货与衍生品交易（上篇）

## 📖 章节概述

本章上篇将全面介绍期货与衍生品交易的基础知识，包括期货合约的基本概念、期货定价理论、套期保值策略以及基础的期货量化策略。这些内容将为您建立期货市场投资的理论基础和实践技能。

## 🎯 学习目标

- 深入理解期货合约的结构和特点
- 掌握期货定价理论和影响因素
- 学会设计和实施套期保值策略
- 了解期货市场的风险管理机制
- 掌握基础期货量化策略的开发
- 建立期货交易的风险管控体系

## 5.1 期货合约基础

### 5.1.1 期货合约概念与特征

```python
import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
import seaborn as sns
from datetime import datetime, timedelta
import warnings
warnings.filterwarnings('ignore')

class FuturesContract:
    """期货合约基础类"""

    def __init__(self, symbol, underlying_asset, contract_size,
                 expiration_date, delivery_month, margin_ratio=0.1):
        self.symbol = symbol  # 合约代码
        self.underlying_asset = underlying_asset  # 标的资产
        self.contract_size = contract_size  # 合约规模
        self.expiration_date = expiration_date  # 到期日
        self.delivery_month = delivery_month  # 交割月份
        self.margin_ratio = margin_ratio  # 保证金比例

    def calculate_margin_requirement(self, price, position_size):
        """计算保证金要求"""
        contract_value = price * self.contract_size * position_size
        initial_margin = contract_value * self.margin_ratio

        return {
            'contract_value': contract_value,
            'initial_margin': initial_margin,
            'margin_ratio': self.margin_ratio
        }

    def calculate_pnl(self, entry_price, current_price, position_size, position_type='long'):
        """计算盈亏"""
        price_diff = current_price - entry_price

        if position_type == 'short':
            price_diff = -price_diff

        pnl = price_diff * self.contract_size * position_size

        return {
            'price_difference': price_diff,
            'total_pnl': pnl,
            'pnl_per_contract': pnl / position_size if position_size != 0 else 0
        }

    def time_to_expiration(self, current_date=None):
        """计算到期时间"""
        if current_date is None:
            current_date = datetime.now()

        time_diff = self.expiration_date - current_date
        return time_diff.days

    def get_contract_info(self):
        """获取合约信息"""
        return {
            'symbol': self.symbol,
            'underlying_asset': self.underlying_asset,
            'contract_size': self.contract_size,
            'expiration_date': self.expiration_date,
            'delivery_month': self.delivery_month,
            'margin_ratio': self.margin_ratio
        }

class FuturesMarketData:
    """期货市场数据处理"""

    def __init__(self):
        self.contracts = {}
        self.price_data = {}

    def add_contract(self, contract):
        """添加期货合约"""
        self.contracts[contract.symbol] = contract

    def calculate_basis(self, futures_price, spot_price):
        """计算基差"""
        basis = futures_price - spot_price
        basis_percentage = (basis / spot_price) * 100

        return {
            'absolute_basis': basis,
            'relative_basis': basis_percentage
        }

    def analyze_contango_backwardation(self, price_curve):
        """分析正向市场和反向市场"""
        # price_curve: {maturity: price} 字典
        sorted_prices = sorted(price_curve.items(), key=lambda x: x[0])

        market_structure = []
        for i in range(len(sorted_prices) - 1):
            current_price = sorted_prices[i][1]
            next_price = sorted_prices[i + 1][1]

            if next_price > current_price:
                structure = 'contango'  # 正向市场
            elif next_price < current_price:
                structure = 'backwardation'  # 反向市场
            else:
                structure = 'flat'

            market_structure.append({
                'period': f"{sorted_prices[i][0]} to {sorted_prices[i+1][0]}",
                'structure': structure,
                'price_diff': next_price - current_price,
                'slope': (next_price - current_price) / current_price
            })

        return market_structure

    def calculate_roll_yield(self, front_month_price, back_month_price, days_to_roll):
        """计算展期收益率"""
        price_ratio = front_month_price / back_month_price
        roll_yield = (price_ratio - 1) * (365 / days_to_roll)

        return {
            'price_ratio': price_ratio,
            'annualized_roll_yield': roll_yield,
            'market_structure': 'contango' if price_ratio < 1 else 'backwardation'
        }

# 使用示例
# 创建期货合约
cu_contract = FuturesContract(
    symbol='CU2312',
    underlying_asset='铜',
    contract_size=5,  # 5吨/手
    expiration_date=datetime(2023, 12, 15),
    delivery_month='2023-12',
    margin_ratio=0.08
)

market_data = FuturesMarketData()
market_data.add_contract(cu_contract)

# 计算保证金
# margin_info = cu_contract.calculate_margin_requirement(price=75000, position_size=10)
# pnl_info = cu_contract.calculate_pnl(entry_price=75000, current_price=76000, position_size=10)
```

### 5.1.2 期货交易机制

```python
class FuturesTradingSystem:
    """期货交易系统"""

    def __init__(self):
        self.positions = {}
        self.margin_account = 0
        self.daily_settlements = []

    def open_position(self, contract_symbol, price, quantity, direction='long'):
        """开仓"""
        if contract_symbol not in self.positions:
            self.positions[contract_symbol] = []

        position = {
            'entry_date': datetime.now(),
            'entry_price': price,
            'quantity': quantity,
            'direction': direction,
            'status': 'open'
        }

        self.positions[contract_symbol].append(position)

        # 计算保证金要求
        contract = self.get_contract_by_symbol(contract_symbol)
        margin_required = contract.calculate_margin_requirement(price, quantity)

        return {
            'position_id': len(self.positions[contract_symbol]) - 1,
            'margin_required': margin_required['initial_margin'],
            'position_info': position
        }

    def close_position(self, contract_symbol, position_id, exit_price):
        """平仓"""
        if contract_symbol in self.positions and position_id < len(self.positions[contract_symbol]):
            position = self.positions[contract_symbol][position_id]

            if position['status'] == 'open':
                position['exit_date'] = datetime.now()
                position['exit_price'] = exit_price
                position['status'] = 'closed'

                # 计算盈亏
                contract = self.get_contract_by_symbol(contract_symbol)
                pnl = contract.calculate_pnl(
                    position['entry_price'],
                    exit_price,
                    position['quantity'],
                    position['direction']
                )

                position['pnl'] = pnl['total_pnl']

                return {
                    'success': True,
                    'pnl': pnl['total_pnl'],
                    'position': position
                }

        return {'success': False, 'error': 'Position not found or already closed'}

    def daily_settlement(self, contract_symbol, settlement_price):
        """每日结算"""
        total_pnl = 0

        if contract_symbol in self.positions:
            for position in self.positions[contract_symbol]:
                if position['status'] == 'open':
                    contract = self.get_contract_by_symbol(contract_symbol)
                    pnl = contract.calculate_pnl(
                        position['entry_price'],
                        settlement_price,
                        position['quantity'],
                        position['direction']
                    )
                    total_pnl += pnl['total_pnl']

        settlement_record = {
            'date': datetime.now().date(),
            'contract': contract_symbol,
            'settlement_price': settlement_price,
            'daily_pnl': total_pnl,
            'margin_balance': self.margin_account + total_pnl
        }

        self.daily_settlements.append(settlement_record)
        return settlement_record

    def calculate_margin_call(self, maintenance_margin_ratio=0.05):
        """计算追保要求"""
        margin_calls = []

        for contract_symbol, positions in self.positions.items():
            total_margin_required = 0

            for position in positions:
                if position['status'] == 'open':
                    contract = self.get_contract_by_symbol(contract_symbol)
                    current_value = position['quantity'] * contract.contract_size * position['entry_price']
                    maintenance_margin = current_value * maintenance_margin_ratio
                    total_margin_required += maintenance_margin

            if total_margin_required > self.margin_account:
                margin_calls.append({
                    'contract': contract_symbol,
                    'required_margin': total_margin_required,
                    'current_balance': self.margin_account,
                    'margin_call_amount': total_margin_required - self.margin_account
                })

        return margin_calls

    def get_contract_by_symbol(self, symbol):
        """根据符号获取合约（简化实现）"""
        # 实际实现中应该从市场数据中获取
        return FuturesContract(symbol, 'commodity', 1, datetime.now() + timedelta(days=30), '2024-01')

# 使用示例
trading_system = FuturesTradingSystem()
trading_system.margin_account = 100000  # 初始保证金账户

# 开仓
# position_result = trading_system.open_position('CU2312', 75000, 10, 'long')
# settlement = trading_system.daily_settlement('CU2312', 76000)
```

## 5.2 期货定价理论

### 5.2.1 无套利定价模型

```python
import numpy as np
from scipy.optimize import fsolve
import math

class FuturesPricingModel:
    """期货定价模型"""

    def __init__(self):
        self.models = {}

    def cost_of_carry_model(self, spot_price, risk_free_rate, time_to_maturity,
                           dividend_yield=0, storage_cost=0, convenience_yield=0):
        """持有成本模型"""
        # F = S * e^((r - d + s - c) * T)
        # F: 期货价格, S: 现货价格, r: 无风险利率, d: 股息率
        # s: 储存成本, c: 便利收益率, T: 到期时间

        cost_of_carry = risk_free_rate - dividend_yield + storage_cost - convenience_yield
        futures_price = spot_price * np.exp(cost_of_carry * time_to_maturity)

        return {
            'theoretical_futures_price': futures_price,
            'cost_of_carry': cost_of_carry,
            'components': {
                'risk_free_rate': risk_free_rate,
                'dividend_yield': dividend_yield,
                'storage_cost': storage_cost,
                'convenience_yield': convenience_yield
            }
        }

    def commodity_futures_pricing(self, spot_price, risk_free_rate, time_to_maturity,
                                 storage_cost_rate, convenience_yield_rate):
        """商品期货定价"""
        net_cost = risk_free_rate + storage_cost_rate - convenience_yield_rate
        futures_price = spot_price * np.exp(net_cost * time_to_maturity)

        return {
            'futures_price': futures_price,
            'net_cost_of_carry': net_cost,
            'storage_cost_component': storage_cost_rate * time_to_maturity,
            'convenience_yield_component': convenience_yield_rate * time_to_maturity
        }

    def stock_index_futures_pricing(self, index_level, risk_free_rate,
                                   dividend_yield, time_to_maturity):
        """股指期货定价"""
        futures_price = index_level * np.exp((risk_free_rate - dividend_yield) * time_to_maturity)

        return {
            'theoretical_price': futures_price,
            'basis': futures_price - index_level,
            'implied_financing_cost': (risk_free_rate - dividend_yield) * time_to_maturity
        }

    def bond_futures_pricing(self, bond_price, repo_rate, time_to_maturity, accrued_interest=0):
        """国债期货定价"""
        # 考虑应计利息的国债期货定价
        clean_price = bond_price - accrued_interest
        futures_price = clean_price * np.exp(repo_rate * time_to_maturity)

        return {
            'clean_futures_price': futures_price,
            'dirty_futures_price': futures_price + accrued_interest,
            'repo_carry': clean_price * (np.exp(repo_rate * time_to_maturity) - 1)
        }

    def calculate_implied_convenience_yield(self, futures_price, spot_price,
                                          risk_free_rate, storage_cost, time_to_maturity):
        """计算隐含便利收益率"""
        def objective_function(convenience_yield):
            theoretical_price = self.cost_of_carry_model(
                spot_price, risk_free_rate, time_to_maturity,
                0, storage_cost, convenience_yield
            )['theoretical_futures_price']
            return theoretical_price - futures_price

        implied_convenience_yield = fsolve(objective_function, 0.05)[0]

        return {
            'implied_convenience_yield': implied_convenience_yield,
            'annual_rate': implied_convenience_yield,
            'market_signal': 'tight supply' if implied_convenience_yield > 0.03 else 'adequate supply'
        }

class ArbitrageStrategy:
    """套利策略"""

    def __init__(self):
        self.positions = []

    def identify_arbitrage_opportunity(self, market_futures_price, theoretical_futures_price,
                                     transaction_cost=0.001):
        """识别套利机会"""
        price_difference = market_futures_price - theoretical_futures_price
        relative_difference = abs(price_difference) / theoretical_futures_price

        if relative_difference > transaction_cost:
            if price_difference > 0:
                strategy = 'sell_futures_buy_underlying'  # 期货高估
                expected_profit = price_difference - (transaction_cost * theoretical_futures_price)
            else:
                strategy = 'buy_futures_sell_underlying'  # 期货低估
                expected_profit = abs(price_difference) - (transaction_cost * theoretical_futures_price)

            return {
                'arbitrage_opportunity': True,
                'strategy': strategy,
                'price_difference': price_difference,
                'expected_profit': expected_profit,
                'profit_margin': expected_profit / theoretical_futures_price
            }

        return {
            'arbitrage_opportunity': False,
            'price_difference': price_difference,
            'threshold_not_met': True
        }

    def execute_cash_carry_arbitrage(self, spot_price, futures_price, theoretical_price,
                                   position_size, financing_rate):
        """执行现货套利"""
        if futures_price > theoretical_price:
            # 卖出期货，买入现货
            strategy = {
                'action': 'cash_and_carry',
                'futures_position': 'short',
                'spot_position': 'long',
                'initial_investment': spot_price * position_size,
                'financing_cost': spot_price * position_size * financing_rate,
                'expected_profit': (futures_price - theoretical_price) * position_size
            }
        else:
            # 买入期货，卖出现货
            strategy = {
                'action': 'reverse_cash_and_carry',
                'futures_position': 'long',
                'spot_position': 'short',
                'cash_received': spot_price * position_size,
                'investment_income': spot_price * position_size * financing_rate,
                'expected_profit': (theoretical_price - futures_price) * position_size
            }

        return strategy

# 使用示例
pricing_model = FuturesPricingModel()
arbitrage_strategy = ArbitrageStrategy()

# 股指期货定价示例
# index_futures_price = pricing_model.stock_index_futures_pricing(
#     index_level=4500,
#     risk_free_rate=0.03,
#     dividend_yield=0.02,
#     time_to_maturity=0.25
# )

# 套利机会识别
# arbitrage_opportunity = arbitrage_strategy.identify_arbitrage_opportunity(
#     market_futures_price=4520,
#     theoretical_futures_price=index_futures_price['theoretical_price']
# )
```

### 5.2.2 期货价格影响因素分析

```python
class FuturesPriceAnalysis:
    """期货价格影响因素分析"""

    def __init__(self):
        self.factors = {}

    def interest_rate_sensitivity(self, spot_price, base_rate, time_to_maturity,
                                 rate_changes=np.linspace(-0.02, 0.02, 21)):
        """利率敏感性分析"""
        pricing_model = FuturesPricingModel()
        sensitivities = []

        for rate_change in rate_changes:
            new_rate = base_rate + rate_change
            futures_price = pricing_model.cost_of_carry_model(
                spot_price, new_rate, time_to_maturity
            )['theoretical_futures_price']

            sensitivities.append({
                'rate_change': rate_change,
                'new_rate': new_rate,
                'futures_price': futures_price,
                'price_change': futures_price - pricing_model.cost_of_carry_model(
                    spot_price, base_rate, time_to_maturity
                )['theoretical_futures_price']
            })

        return pd.DataFrame(sensitivities)

    def time_decay_analysis(self, spot_price, risk_free_rate, dividend_yield=0,
                           max_time=1.0, time_steps=50):
        """时间价值衰减分析"""
        pricing_model = FuturesPricingModel()
        time_points = np.linspace(max_time, 0.01, time_steps)

        time_analysis = []
        for time_to_maturity in time_points:
            futures_price = pricing_model.stock_index_futures_pricing(
                spot_price, risk_free_rate, dividend_yield, time_to_maturity
            )['theoretical_price']

            basis = futures_price - spot_price

            time_analysis.append({
                'time_to_maturity': time_to_maturity,
                'futures_price': futures_price,
                'basis': basis,
                'time_premium': basis / spot_price if spot_price != 0 else 0
            })

        return pd.DataFrame(time_analysis)

    def volatility_impact_analysis(self, price_data, window=30):
        """波动率对期货价格的影响"""
        returns = price_data.pct_change().dropna()
        rolling_volatility = returns.rolling(window).std() * np.sqrt(252)

        # 分析波动率与基差的关系
        if 'futures_price' in price_data.columns and 'spot_price' in price_data.columns:
            basis = price_data['futures_price'] - price_data['spot_price']

            volatility_analysis = pd.DataFrame({
                'volatility': rolling_volatility,
                'basis': basis,
                'basis_percentage': basis / price_data['spot_price']
            }).dropna()

            correlation = volatility_analysis['volatility'].corr(volatility_analysis['basis_percentage'])

            return {
                'volatility_data': volatility_analysis,
                'vol_basis_correlation': correlation,
                'analysis': 'Higher volatility typically increases time value' if correlation > 0 else 'Inverse relationship observed'
            }

        return {'volatility_data': rolling_volatility}

    def seasonal_pattern_analysis(self, price_data, commodity_type='agricultural'):
        """季节性模式分析"""
        price_data['month'] = price_data.index.month
        price_data['year'] = price_data.index.year

        monthly_returns = price_data.groupby('month')['price'].pct_change().mean()
        monthly_volatility = price_data.groupby('month')['price'].pct_change().std()

        seasonal_patterns = pd.DataFrame({
            'average_return': monthly_returns,
            'volatility': monthly_volatility
        })

        # 识别季节性特征
        if commodity_type == 'agricultural':
            harvest_months = [9, 10, 11]  # 秋季收获期
            planting_months = [3, 4, 5]   # 春季种植期

            seasonal_patterns['period_type'] = seasonal_patterns.index.map(
                lambda x: 'harvest' if x in harvest_months
                         else 'planting' if x in planting_months
                         else 'normal'
            )

        return seasonal_patterns

    def supply_demand_analysis(self, inventory_data, consumption_data, production_data):
        """供需分析"""
        # 计算供需平衡
        supply_demand_balance = production_data + inventory_data.shift(1) - consumption_data

        # 计算库存消费比
        inventory_consumption_ratio = inventory_data / consumption_data

        # 分析价格与基本面的关系
        fundamental_indicators = pd.DataFrame({
            'supply_demand_balance': supply_demand_balance,
            'inventory_ratio': inventory_consumption_ratio,
            'production_growth': production_data.pct_change(),
            'consumption_growth': consumption_data.pct_change()
        })

        return fundamental_indicators

# 使用示例
price_analysis = FuturesPriceAnalysis()

# 利率敏感性分析
# rate_sensitivity = price_analysis.interest_rate_sensitivity(
#     spot_price=4500, base_rate=0.03, time_to_maturity=0.25
# )

# 时间价值衰减分析
# time_decay = price_analysis.time_decay_analysis(
#     spot_price=4500, risk_free_rate=0.03, dividend_yield=0.02
# )
```

## 5.3 套期保值策略

### 5.3.1 套期保值基础理论

```python
class HedgingStrategy:
    """套期保值策略"""

    def __init__(self):
        self.hedge_positions = []

    def calculate_hedge_ratio(self, spot_returns, futures_returns, method='ols'):
        """计算套期保值比率"""
        if method == 'ols':
            # 最小二乘法
            correlation = np.corrcoef(spot_returns, futures_returns)[0, 1]
            spot_std = np.std(spot_returns)
            futures_std = np.std(futures_returns)

            hedge_ratio = correlation * (spot_std / futures_std)

        elif method == 'minimum_variance':
            # 最小方差套期保值比率
            covariance = np.cov(spot_returns, futures_returns)[0, 1]
            futures_variance = np.var(futures_returns)

            hedge_ratio = covariance / futures_variance

        return {
            'hedge_ratio': hedge_ratio,
            'correlation': np.corrcoef(spot_returns, futures_returns)[0, 1],
            'r_squared': np.corrcoef(spot_returns, futures_returns)[0, 1] ** 2,
            'method': method
        }

    def static_hedge(self, spot_position, spot_price, futures_price, hedge_ratio):
        """静态套期保值"""
        futures_position = -hedge_ratio * spot_position  # 负号表示反向

        hedge_portfolio = {
            'spot_position': spot_position,
            'futures_position': futures_position,
            'initial_spot_value': spot_position * spot_price,
            'initial_futures_value': futures_position * futures_price,
            'hedge_ratio': hedge_ratio,
            'hedge_effectiveness': None  # 需要后续计算
        }

        return hedge_portfolio

    def dynamic_hedge(self, spot_position, price_data, hedge_ratios, rebalance_frequency='weekly'):
        """动态套期保值"""
        hedge_history = []
        current_futures_position = 0

        rebalance_dates = self._get_rebalance_dates(price_data.index, rebalance_frequency)

        for date in rebalance_dates:
            if date in price_data.index and date in hedge_ratios.index:
                target_hedge_ratio = hedge_ratios.loc[date]
                target_futures_position = -target_hedge_ratio * spot_position

                position_change = target_futures_position - current_futures_position

                hedge_record = {
                    'date': date,
                    'hedge_ratio': target_hedge_ratio,
                    'target_position': target_futures_position,
                    'position_change': position_change,
                    'current_position': target_futures_position
                }

                hedge_history.append(hedge_record)
                current_futures_position = target_futures_position

        return pd.DataFrame(hedge_history)

    def _get_rebalance_dates(self, date_index, frequency):
        """获取再平衡日期"""
        if frequency == 'daily':
            return date_index
        elif frequency == 'weekly':
            return date_index[date_index.dayofweek == 0]  # 周一
        elif frequency == 'monthly':
            return date_index[date_index.day == 1]  # 月初
        else:
            return date_index

    def calculate_hedge_effectiveness(self, hedged_portfolio_returns, unhedged_returns):
        """计算套期保值有效性"""
        hedged_variance = np.var(hedged_portfolio_returns)
        unhedged_variance = np.var(unhedged_returns)

        variance_reduction = (unhedged_variance - hedged_variance) / unhedged_variance

        effectiveness_metrics = {
            'variance_reduction': variance_reduction,
            'hedge_effectiveness': variance_reduction,
            'hedged_volatility': np.std(hedged_portfolio_returns),
            'unhedged_volatility': np.std(unhedged_returns),
            'volatility_reduction': 1 - (np.std(hedged_portfolio_returns) / np.std(unhedged_returns))
        }

        return effectiveness_metrics

class PortfolioHedging:
    """投资组合套期保值"""

    def __init__(self):
        self.portfolio = {}
        self.hedge_instruments = {}

    def beta_hedging(self, portfolio_beta, index_futures_price, portfolio_value):
        """Beta套期保值"""
        # 计算需要的期货合约数量来对冲系统性风险
        futures_contracts_needed = -(portfolio_beta * portfolio_value) / index_futures_price

        return {
            'contracts_needed': futures_contracts_needed,
            'hedge_value': futures_contracts_needed * index_futures_price,
            'residual_beta': 0,  # 理论上完全对冲后Beta为0
            'systematic_risk_hedged': True
        }

    def duration_hedging(self, bond_portfolio_duration, portfolio_value,
                        futures_duration, futures_price):
        """久期套期保值"""
        # 计算需要的国债期货合约数量
        duration_hedge_ratio = bond_portfolio_duration / futures_duration
        contracts_needed = -(duration_hedge_ratio * portfolio_value) / futures_price

        return {
            'duration_hedge_ratio': duration_hedge_ratio,
            'contracts_needed': contracts_needed,
            'portfolio_duration': bond_portfolio_duration,
            'futures_duration': futures_duration,
            'hedged_duration': 0  # 完全对冲后久期为0
        }

    def currency_hedging(self, foreign_asset_value, exchange_rate,
                        currency_futures_price, contract_size):
        """外汇套期保值"""
        # 计算外汇期货合约数量
        exposure_in_foreign_currency = foreign_asset_value / exchange_rate
        contracts_needed = -(exposure_in_foreign_currency / contract_size)

        return {
            'foreign_exposure': exposure_in_foreign_currency,
            'contracts_needed': contracts_needed,
            'hedge_value': contracts_needed * currency_futures_price * contract_size,
            'currency_risk_hedged': True
        }

# 使用示例
hedging_strategy = HedgingStrategy()
portfolio_hedging = PortfolioHedging()

# 计算套期保值比率示例
# spot_returns = np.random.normal(0, 0.02, 252)
# futures_returns = np.random.normal(0, 0.025, 252)
# hedge_ratio_result = hedging_strategy.calculate_hedge_ratio(spot_returns, futures_returns)

# Beta套期保值示例
# beta_hedge = portfolio_hedging.beta_hedging(
#     portfolio_beta=1.2,
#     index_futures_price=4500,
#     portfolio_value=1000000
# )
```

### 5.3.2 套期保值策略实施

```python
class HedgingImplementation:
    """套期保值策略实施"""

    def __init__(self):
        self.active_hedges = {}

    def basis_risk_analysis(self, spot_prices, futures_prices):
        """基差风险分析"""
        basis = futures_prices - spot_prices
        basis_changes = basis.diff()

        basis_analysis = {
            'average_basis': basis.mean(),
            'basis_volatility': basis.std(),
            'basis_change_volatility': basis_changes.std(),
            'basis_correlation_with_spot': basis.corr(spot_prices),
            'convergence_rate': self._calculate_convergence_rate(basis)
        }

        return basis_analysis

    def _calculate_convergence_rate(self, basis_series):
        """计算基差收敛速度"""
        # 简化计算：基差绝对值的平均下降速度
        abs_basis = abs(basis_series)
        convergence_rate = -abs_basis.diff().mean()
        return convergence_rate

    def cross_hedge_analysis(self, target_asset_returns, hedge_instrument_returns):
        """交叉套期保值分析"""
        correlation = np.corrcoef(target_asset_returns, hedge_instrument_returns)[0, 1]

        # 计算交叉套期保值比率
        target_std = np.std(target_asset_returns)
        hedge_std = np.std(hedge_instrument_returns)

        cross_hedge_ratio = correlation * (target_std / hedge_std)

        # 评估交叉套期保值效果
        hedge_effectiveness = correlation ** 2

        return {
            'cross_hedge_ratio': cross_hedge_ratio,
            'correlation': correlation,
            'hedge_effectiveness': hedge_effectiveness,
            'recommendation': 'suitable' if correlation > 0.7 else 'not recommended'
        }

    def hedge_accounting_framework(self, hedge_type, hedge_item_fair_value,
                                  hedging_instrument_fair_value):
        """套期保值会计框架"""
        # 套期保值有效性测试
        hedge_effectiveness_ratio = abs(hedging_instrument_fair_value / hedge_item_fair_value)

        is_effective = 0.8 <= hedge_effectiveness_ratio <= 1.25

        accounting_treatment = {
            'hedge_type': hedge_type,  # 'fair_value', 'cash_flow', 'net_investment'
            'effectiveness_ratio': hedge_effectiveness_ratio,
            'is_effective': is_effective,
            'accounting_method': self._determine_accounting_method(hedge_type, is_effective)
        }

        return accounting_treatment

    def _determine_accounting_method(self, hedge_type, is_effective):
        """确定会计处理方法"""
        if not is_effective:
            return 'hedge_accounting_discontinued'

        if hedge_type == 'fair_value':
            return 'fair_value_hedge_accounting'
        elif hedge_type == 'cash_flow':
            return 'cash_flow_hedge_accounting'
        elif hedge_type == 'net_investment':
            return 'net_investment_hedge_accounting'
        else:
            return 'standard_accounting'

    def dynamic_hedge_optimization(self, returns_data, transaction_costs=0.001,
                                  rebalance_threshold=0.1):
        """动态套期保值优化"""
        optimal_hedge_schedule = []
        current_hedge_ratio = 0

        for date, return_data in returns_data.iterrows():
            # 重新计算最优套期保值比率
            recent_data = returns_data.loc[:date].tail(60)  # 最近60个观测值

            if len(recent_data) >= 30:
                spot_returns = recent_data['spot_returns']
                futures_returns = recent_data['futures_returns']

                # 计算当期最优套期保值比率
                covariance = np.cov(spot_returns, futures_returns)[0, 1]
                futures_variance = np.var(futures_returns)
                optimal_ratio = covariance / futures_variance

                # 判断是否需要调整
                ratio_change = abs(optimal_ratio - current_hedge_ratio)

                if ratio_change > rebalance_threshold:
                    # 考虑交易成本
                    expected_benefit = self._calculate_hedge_benefit(
                        spot_returns, futures_returns, optimal_ratio, current_hedge_ratio
                    )
                    transaction_cost = transaction_costs * abs(optimal_ratio - current_hedge_ratio)

                    if expected_benefit > transaction_cost:
                        optimal_hedge_schedule.append({
                            'date': date,
                            'old_ratio': current_hedge_ratio,
                            'new_ratio': optimal_ratio,
                            'adjustment': optimal_ratio - current_hedge_ratio,
                            'expected_benefit': expected_benefit,
                            'transaction_cost': transaction_cost,
                            'net_benefit': expected_benefit - transaction_cost
                        })
                        current_hedge_ratio = optimal_ratio

        return pd.DataFrame(optimal_hedge_schedule)

    def _calculate_hedge_benefit(self, spot_returns, futures_returns,
                                new_ratio, old_ratio):
        """计算套期保值收益"""
        # 简化计算：方差减少的价值
        old_hedged_returns = spot_returns + old_ratio * futures_returns
        new_hedged_returns = spot_returns + new_ratio * futures_returns

        old_variance = np.var(old_hedged_returns)
        new_variance = np.var(new_hedged_returns)

        variance_reduction = old_variance - new_variance
        return variance_reduction  # 可以转换为货币单位

# 使用示例
hedge_implementation = HedgingImplementation()

# 基差风险分析示例
# basis_risk = hedge_implementation.basis_risk_analysis(spot_prices, futures_prices)

# 交叉套期保值分析示例
# cross_hedge = hedge_implementation.cross_hedge_analysis(target_returns, hedge_returns)
```

## 5.4 期货量化策略

### 5.4.1 趋势跟踪策略

```python
class FuturesTrendStrategy:
    """期货趋势策略"""

    def __init__(self):
        self.positions = {}
        self.signals = {}

    def moving_average_strategy(self, price_data, short_window=20, long_window=60):
        """移动平均线策略"""
        signals = pd.DataFrame(index=price_data.index)
        signals['price'] = price_data
        signals['short_ma'] = price_data.rolling(window=short_window).mean()
        signals['long_ma'] = price_data.rolling(window=long_window).mean()

        # 生成信号
        signals['signal'] = 0
        signals['signal'][short_window:] = np.where(
            signals['short_ma'][short_window:] > signals['long_ma'][short_window:], 1, 0
        )
        signals['positions'] = signals['signal'].diff()

        return signals

    def breakout_strategy(self, price_data, lookback_period=20, atr_multiplier=2):
        """突破策略"""
        signals = pd.DataFrame(index=price_data.index)
        signals['price'] = price_data

        # 计算高点和低点
        signals['highest'] = price_data.rolling(window=lookback_period).max()
        signals['lowest'] = price_data.rolling(window=lookback_period).min()

        # 计算ATR用于止损
        high_low = price_data.rolling(window=2).max() - price_data.rolling(window=2).min()
        signals['atr'] = high_low.rolling(window=14).mean()

        # 生成信号
        signals['long_signal'] = price_data > signals['highest'].shift(1)
        signals['short_signal'] = price_data < signals['lowest'].shift(1)

        signals['position'] = 0
        signals.loc[signals['long_signal'], 'position'] = 1
        signals.loc[signals['short_signal'], 'position'] = -1

        # 计算止损位
        signals['long_stop'] = price_data - atr_multiplier * signals['atr']
        signals['short_stop'] = price_data + atr_multiplier * signals['atr']

        return signals

    def momentum_strategy(self, price_data, returns_data, momentum_period=12,
                         volatility_period=20):
        """动量策略"""
        signals = pd.DataFrame(index=price_data.index)
        signals['price'] = price_data
        signals['returns'] = returns_data

        # 计算动量指标
        signals['momentum'] = returns_data.rolling(window=momentum_period).sum()
        signals['volatility'] = returns_data.rolling(window=volatility_period).std()

        # 动量强度调整
        signals['risk_adjusted_momentum'] = signals['momentum'] / signals['volatility']

        # 生成信号
        momentum_threshold = signals['risk_adjusted_momentum'].quantile(0.7)

        signals['signal'] = 0
        signals.loc[signals['risk_adjusted_momentum'] > momentum_threshold, 'signal'] = 1
        signals.loc[signals['risk_adjusted_momentum'] < -momentum_threshold, 'signal'] = -1

        return signals

    def multi_timeframe_strategy(self, price_data, short_term_period=5,
                               medium_term_period=20, long_term_period=60):
        """多时间框架策略"""
        signals = pd.DataFrame(index=price_data.index)
        signals['price'] = price_data

        # 短期趋势
        signals['short_trend'] = price_data.rolling(window=short_term_period).mean()
        signals['short_signal'] = np.where(price_data > signals['short_trend'], 1, -1)

        # 中期趋势
        signals['medium_trend'] = price_data.rolling(window=medium_term_period).mean()
        signals['medium_signal'] = np.where(price_data > signals['medium_trend'], 1, -1)

        # 长期趋势
        signals['long_trend'] = price_data.rolling(window=long_term_period).mean()
        signals['long_signal'] = np.where(price_data > signals['long_trend'], 1, -1)

        # 综合信号
        signals['combined_signal'] = (
            signals['short_signal'] * 0.5 +
            signals['medium_signal'] * 0.3 +
            signals['long_signal'] * 0.2
        )

        # 最终位置
        signals['position'] = 0
        signals.loc[signals['combined_signal'] > 0.3, 'position'] = 1
        signals.loc[signals['combined_signal'] < -0.3, 'position'] = -1

        return signals

class FuturesArbitrageStrategy:
    """期货套利策略"""

    def __init__(self):
        self.arbitrage_positions = {}

    def calendar_spread_strategy(self, near_month_prices, far_month_prices,
                               spread_threshold=2):
        """跨期套利策略"""
        spread = far_month_prices - near_month_prices
        spread_z_score = (spread - spread.rolling(60).mean()) / spread.rolling(60).std()

        signals = pd.DataFrame(index=spread.index)
        signals['spread'] = spread
        signals['z_score'] = spread_z_score

        # 套利信号
        signals['long_spread'] = spread_z_score < -spread_threshold  # 价差过小，买远卖近
        signals['short_spread'] = spread_z_score > spread_threshold   # 价差过大，卖远买近

        signals['position'] = 0
        signals.loc[signals['long_spread'], 'position'] = 1   # 买价差
        signals.loc[signals['short_spread'], 'position'] = -1  # 卖价差

        return signals

    def inter_commodity_spread(self, commodity1_prices, commodity2_prices,
                             hedge_ratio=1, spread_threshold=2):
        """跨品种套利"""
        # 计算调整后的价差
        adjusted_spread = commodity1_prices - hedge_ratio * commodity2_prices
        spread_z_score = (adjusted_spread - adjusted_spread.rolling(60).mean()) / adjusted_spread.rolling(60).std()

        signals = pd.DataFrame(index=adjusted_spread.index)
        signals['spread'] = adjusted_spread
        signals['z_score'] = spread_z_score

        # 套利信号
        signals['long_spread'] = spread_z_score < -spread_threshold
        signals['short_spread'] = spread_z_score > spread_threshold

        return signals

    def basis_trading_strategy(self, spot_prices, futures_prices,
                             normal_basis_range=(-50, 50)):
        """基差交易策略"""
        basis = futures_prices - spot_prices

        signals = pd.DataFrame(index=basis.index)
        signals['basis'] = basis
        signals['spot_price'] = spot_prices
        signals['futures_price'] = futures_prices

        # 基差异常信号
        signals['basis_too_high'] = basis > normal_basis_range[1]
        signals['basis_too_low'] = basis < normal_basis_range[0]

        # 交易信号
        signals['sell_futures_buy_spot'] = signals['basis_too_high']
        signals['buy_futures_sell_spot'] = signals['basis_too_low']

        return signals

# 使用示例
trend_strategy = FuturesTrendStrategy()
arbitrage_strategy = FuturesArbitrageStrategy()

# 移动平均策略示例
# ma_signals = trend_strategy.moving_average_strategy(price_data)

# 日历价差套利示例
# calendar_signals = arbitrage_strategy.calendar_spread_strategy(near_prices, far_prices)
```

## 5.5 本章总结

在本章上篇中，我们全面学习了期货与衍生品交易的基础知识：

### 🎯 核心要点回顾

1. **期货合约基础**

   - 期货合约的结构和特征
   - 保证金制度和每日结算
   - 期货交易机制

2. **期货定价理论**

   - 无套利定价模型
   - 持有成本模型
   - 价格影响因素分析

3. **套期保值策略**

   - 套期保值比率计算
   - 静态和动态套期保值
   - 组合套期保值方法

4. **期货量化策略**
   - 趋势跟踪策略
   - 套利策略
   - 多时间框架分析

### 📈 实践应用指南

1. **风险管理要点**

   - 保证金管理和资金控制
   - 基差风险和时间风险
   - 流动性风险考虑

2. **策略开发建议**
   - 回测验证的重要性
   - 交易成本的考虑
   - 风险调整收益评估

## 📚 延伸学习

在[第五章：期货与衍生品交易（下篇）](05_futures_derivatives_trading_part2.md)中，我们将学习期权基础与策略、量化套利策略、高频交易入门以及风险管理技术。

## 💡 实践练习

1. **期货定价实践**：使用持有成本模型计算不同商品期货的理论价格

2. **套期保值策略**：为投资组合设计完整的套期保值方案

3. **趋势策略开发**：开发并回测一个多因子期货趋势策略

4. **套利机会识别**：分析真实市场数据中的套利机会

5. **风险管理系统**：建立期货交易的风险控制框架

---

**建议学习时间：** 4-5 天  
**前置章节：** [第四章：股票投资与分析（下篇）](04_stock_investment_analysis_part2.md)  
**下一章：** [期货与衍生品交易（下篇）](05_futures_derivatives_trading_part2.md)
