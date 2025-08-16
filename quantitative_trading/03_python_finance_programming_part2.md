# 第三章：Python 金融编程（下篇）

## 📖 章节概述

本章下篇将深入学习 Python 在金融领域的高级应用，包括机器学习在金融中的应用、高级数据分析技术、策略回测框架的构建，以及交易执行系统的设计。这些内容将帮助您从基础的数据分析进阶到专业的量化交易系统开发。

## 🎯 学习目标

- 掌握机器学习在金融预测中的应用
- 学会构建专业的策略回测框架
- 了解交易执行系统的设计原理
- 掌握高级数据分析技术和另类数据处理
- 学会策略组合和风险管理的实现
- 建立完整的量化交易工作流

## 3.6 机器学习在金融中的应用

### 3.6.1 特征工程

```python
import pandas as pd
import numpy as np
from sklearn.preprocessing import StandardScaler, RobustScaler
from sklearn.feature_selection import SelectKBest, f_regression, mutual_info_regression
import warnings
warnings.filterwarnings('ignore')

class FinancialFeatureEngineering:
    """金融特征工程工具"""

    def __init__(self):
        self.scaler = None
        self.feature_selector = None
        self.feature_names = []

    def create_technical_features(self, df):
        """创建技术分析特征"""
        features = df.copy()

        # 价格相关特征
        features['price_change'] = features['Close'].pct_change()
        features['price_change_2d'] = features['Close'].pct_change(2)
        features['price_change_5d'] = features['Close'].pct_change(5)
        features['price_change_10d'] = features['Close'].pct_change(10)

        # 移动平均特征
        for window in [5, 10, 20, 50]:
            features[f'sma_{window}'] = features['Close'].rolling(window).mean()
            features[f'price_to_sma_{window}'] = features['Close'] / features[f'sma_{window}'] - 1
            features[f'sma_{window}_slope'] = features[f'sma_{window}'].diff(5) / features[f'sma_{window}'].shift(5)

        # 波动率特征
        for window in [5, 10, 20]:
            features[f'volatility_{window}'] = features['price_change'].rolling(window).std()
            features[f'volatility_{window}_rank'] = features[f'volatility_{window}'].rolling(60).rank(pct=True)

        # 成交量特征
        features['volume_change'] = features['Volume'].pct_change()
        features['volume_sma_20'] = features['Volume'].rolling(20).mean()
        features['volume_ratio'] = features['Volume'] / features['volume_sma_20']

        # 价量结合特征
        features['price_volume_trend'] = (features['Close'] - features['Close'].shift(1)) * features['Volume']
        features['money_flow'] = features['Close'] * features['Volume']
        features['money_flow_ratio'] = features['money_flow'] / features['money_flow'].rolling(20).mean()

        # 技术指标特征
        features = self._add_technical_indicators(features)

        return features

    def _add_technical_indicators(self, df):
        """添加技术指标"""
        # RSI
        delta = df['Close'].diff()
        gain = (delta.where(delta > 0, 0)).rolling(window=14).mean()
        loss = (-delta.where(delta < 0, 0)).rolling(window=14).mean()
        rs = gain / loss
        df['rsi'] = 100 - (100 / (1 + rs))
        df['rsi_oversold'] = (df['rsi'] < 30).astype(int)
        df['rsi_overbought'] = (df['rsi'] > 70).astype(int)

        # MACD
        exp1 = df['Close'].ewm(span=12).mean()
        exp2 = df['Close'].ewm(span=26).mean()
        df['macd'] = exp1 - exp2
        df['macd_signal'] = df['macd'].ewm(span=9).mean()
        df['macd_histogram'] = df['macd'] - df['macd_signal']

        # 布林带
        df['bb_middle'] = df['Close'].rolling(20).mean()
        bb_std = df['Close'].rolling(20).std()
        df['bb_upper'] = df['bb_middle'] + (bb_std * 2)
        df['bb_lower'] = df['bb_middle'] - (bb_std * 2)
        df['bb_width'] = (df['bb_upper'] - df['bb_lower']) / df['bb_middle']
        df['bb_position'] = (df['Close'] - df['bb_lower']) / (df['bb_upper'] - df['bb_lower'])

        return df

    def create_fundamental_features(self, price_df, fundamental_df):
        """创建基本面特征"""
        # 这里假设fundamental_df包含财务数据
        features = price_df.copy()

        if fundamental_df is not None:
            # 估值比率
            if 'PE' in fundamental_df.columns:
                features['pe_ratio'] = fundamental_df['PE']
                features['pe_percentile'] = fundamental_df['PE'].rolling(252).rank(pct=True)

            if 'PB' in fundamental_df.columns:
                features['pb_ratio'] = fundamental_df['PB']
                features['pb_percentile'] = fundamental_df['PB'].rolling(252).rank(pct=True)

            # 盈利能力
            if 'ROE' in fundamental_df.columns:
                features['roe'] = fundamental_df['ROE']
                features['roe_trend'] = fundamental_df['ROE'].diff(4)  # 季度变化

        return features

    def create_market_features(self, df, market_index):
        """创建市场相关特征"""
        features = df.copy()

        # 相对市场表现
        features['market_return'] = market_index.pct_change()
        features['relative_return'] = features['price_change'] - features['market_return']

        # Beta计算
        for window in [30, 60, 120]:
            correlation = features['price_change'].rolling(window).corr(features['market_return'])
            market_vol = features['market_return'].rolling(window).std()
            stock_vol = features['price_change'].rolling(window).std()
            features[f'beta_{window}'] = correlation * (stock_vol / market_vol)

        # 市场情绪指标
        features['market_volatility'] = features['market_return'].rolling(20).std()
        features['market_trend'] = features['market_return'].rolling(10).mean()

        return features

    def select_features(self, X, y, method='mutual_info', k=50):
        """特征选择"""
        if method == 'mutual_info':
            selector = SelectKBest(score_func=mutual_info_regression, k=k)
        elif method == 'f_regression':
            selector = SelectKBest(score_func=f_regression, k=k)
        else:
            raise ValueError("支持的方法: 'mutual_info', 'f_regression'")

        X_selected = selector.fit_transform(X, y)
        selected_features = X.columns[selector.get_support()]

        self.feature_selector = selector
        self.feature_names = selected_features

        print(f"从{len(X.columns)}个特征中选择了{len(selected_features)}个特征")
        print("选择的特征:")
        for i, feature in enumerate(selected_features[:10]):  # 显示前10个
            score = selector.scores_[X.columns.get_loc(feature)]
            print(f"  {i+1}. {feature}: {score:.4f}")

        return X_selected, selected_features

    def scale_features(self, X_train, X_test=None, method='robust'):
        """特征缩放"""
        if method == 'standard':
            scaler = StandardScaler()
        elif method == 'robust':
            scaler = RobustScaler()
        else:
            raise ValueError("支持的方法: 'standard', 'robust'")

        X_train_scaled = scaler.fit_transform(X_train)
        self.scaler = scaler

        if X_test is not None:
            X_test_scaled = scaler.transform(X_test)
            return X_train_scaled, X_test_scaled

        return X_train_scaled

# 示例使用
# 假设我们有股票数据
np.random.seed(42)
dates = pd.date_range('2020-01-01', periods=1000, freq='D')
stock_data = pd.DataFrame({
    'Close': 100 + np.cumsum(np.random.randn(1000) * 0.02),
    'High': lambda x: x['Close'] * (1 + np.random.uniform(0, 0.02, 1000)),
    'Low': lambda x: x['Close'] * (1 - np.random.uniform(0, 0.02, 1000)),
    'Volume': np.random.randint(1000000, 10000000, 1000)
}, index=dates)

stock_data['Open'] = stock_data['Close'].shift(1).fillna(stock_data['Close'].iloc[0])

# 特征工程
feature_engineer = FinancialFeatureEngineering()
enriched_data = feature_engineer.create_technical_features(stock_data)

print("特征工程结果:")
print(f"原始特征数: {len(stock_data.columns)}")
print(f"工程化后特征数: {len(enriched_data.columns)}")
print(f"新增特征: {list(set(enriched_data.columns) - set(stock_data.columns))[:10]}")
```

### 3.6.2 机器学习模型应用

```python
from sklearn.ensemble import RandomForestRegressor, GradientBoostingRegressor
from sklearn.linear_model import LinearRegression, Ridge, Lasso
from sklearn.svm import SVR
from sklearn.metrics import mean_squared_error, mean_absolute_error, r2_score
from sklearn.model_selection import TimeSeriesSplit, GridSearchCV
import xgboost as xgb
import lightgbm as lgb

class FinancialMLModels:
    """金融机器学习模型"""

    def __init__(self):
        self.models = {}
        self.predictions = {}
        self.feature_importance = {}

    def prepare_data(self, df, target_col='future_return', lookback=20):
        """准备机器学习数据"""
        # 创建目标变量（未来收益率）
        df[target_col] = df['Close'].shift(-1) / df['Close'] - 1

        # 移除包含NaN的行
        df_clean = df.dropna()

        # 分离特征和目标
        feature_cols = [col for col in df_clean.columns
                       if col not in ['Close', 'Open', 'High', 'Low', target_col]]

        X = df_clean[feature_cols]
        y = df_clean[target_col]

        return X, y, df_clean

    def time_series_split(self, X, y, n_splits=5):
        """时间序列交叉验证分割"""
        tscv = TimeSeriesSplit(n_splits=n_splits)

        splits = []
        for train_idx, test_idx in tscv.split(X):
            splits.append({
                'train_idx': train_idx,
                'test_idx': test_idx,
                'train_size': len(train_idx),
                'test_size': len(test_idx)
            })

        return splits

    def train_linear_models(self, X_train, y_train, X_test, y_test):
        """训练线性模型"""
        models = {
            'Linear Regression': LinearRegression(),
            'Ridge': Ridge(alpha=1.0),
            'Lasso': Lasso(alpha=0.1)
        }

        results = {}
        for name, model in models.items():
            # 训练模型
            model.fit(X_train, y_train)

            # 预测
            y_pred = model.predict(X_test)

            # 评估
            mse = mean_squared_error(y_test, y_pred)
            mae = mean_absolute_error(y_test, y_pred)
            r2 = r2_score(y_test, y_pred)

            results[name] = {
                'model': model,
                'predictions': y_pred,
                'mse': mse,
                'mae': mae,
                'r2': r2
            }

            self.models[name] = model

        return results

    def train_tree_models(self, X_train, y_train, X_test, y_test):
        """训练树模型"""
        models = {
            'Random Forest': RandomForestRegressor(n_estimators=100, random_state=42),
            'Gradient Boosting': GradientBoostingRegressor(n_estimators=100, random_state=42),
            'XGBoost': xgb.XGBRegressor(n_estimators=100, random_state=42),
            'LightGBM': lgb.LGBMRegressor(n_estimators=100, random_state=42)
        }

        results = {}
        for name, model in models.items():
            # 训练模型
            model.fit(X_train, y_train)

            # 预测
            y_pred = model.predict(X_test)

            # 评估
            mse = mean_squared_error(y_test, y_pred)
            mae = mean_absolute_error(y_test, y_pred)
            r2 = r2_score(y_test, y_pred)

            # 特征重要性
            if hasattr(model, 'feature_importances_'):
                importance = dict(zip(X_train.columns, model.feature_importances_))
                self.feature_importance[name] = importance

            results[name] = {
                'model': model,
                'predictions': y_pred,
                'mse': mse,
                'mae': mae,
                'r2': r2
            }

            self.models[name] = model

        return results

    def hyperparameter_tuning(self, X_train, y_train, model_type='xgboost'):
        """超参数调优"""
        if model_type == 'xgboost':
            model = xgb.XGBRegressor(random_state=42)
            param_grid = {
                'n_estimators': [50, 100, 200],
                'max_depth': [3, 6, 10],
                'learning_rate': [0.01, 0.1, 0.2],
                'subsample': [0.8, 1.0]
            }
        elif model_type == 'random_forest':
            model = RandomForestRegressor(random_state=42)
            param_grid = {
                'n_estimators': [50, 100, 200],
                'max_depth': [5, 10, None],
                'min_samples_split': [2, 5, 10],
                'min_samples_leaf': [1, 2, 4]
            }
        else:
            raise ValueError("支持的模型类型: 'xgboost', 'random_forest'")

        # 时间序列交叉验证
        tscv = TimeSeriesSplit(n_splits=3)

        # 网格搜索
        grid_search = GridSearchCV(
            model, param_grid, cv=tscv,
            scoring='neg_mean_squared_error',
            n_jobs=-1, verbose=1
        )

        grid_search.fit(X_train, y_train)

        print(f"最佳参数 ({model_type}):")
        for param, value in grid_search.best_params_.items():
            print(f"  {param}: {value}")

        print(f"最佳CV得分: {-grid_search.best_score_:.6f}")

        return grid_search.best_estimator_

    def ensemble_predictions(self, predictions_dict, method='average'):
        """集成预测"""
        if method == 'average':
            # 简单平均
            pred_values = list(predictions_dict.values())
            ensemble_pred = np.mean(pred_values, axis=0)
        elif method == 'weighted':
            # 基于R²的加权平均
            weights = []
            predictions = []

            for model_name, pred in predictions_dict.items():
                if model_name in self.models:
                    # 这里简化处理，实际应该用验证集计算权重
                    weight = 1.0  # 可以基于模型性能计算权重
                    weights.append(weight)
                    predictions.append(pred)

            weights = np.array(weights) / np.sum(weights)
            ensemble_pred = np.average(predictions, axis=0, weights=weights)

        return ensemble_pred

    def display_model_comparison(self, results):
        """显示模型比较结果"""
        print("\n=== 模型性能比较 ===")
        print(f"{'模型':<15} {'MSE':<12} {'MAE':<12} {'R²':<12}")
        print("-" * 55)

        for model_name, metrics in results.items():
            print(f"{model_name:<15} {metrics['mse']:<12.6f} "
                 f"{metrics['mae']:<12.6f} {metrics['r2']:<12.6f}")

    def display_feature_importance(self, model_name, top_k=10):
        """显示特征重要性"""
        if model_name not in self.feature_importance:
            print(f"模型 {model_name} 没有特征重要性信息")
            return

        importance = self.feature_importance[model_name]
        sorted_features = sorted(importance.items(), key=lambda x: x[1], reverse=True)

        print(f"\n=== {model_name} 特征重要性 (Top {top_k}) ===")
        for i, (feature, score) in enumerate(sorted_features[:top_k]):
            print(f"{i+1:2d}. {feature:<25} {score:.4f}")

# 示例使用
ml_models = FinancialMLModels()

# 准备数据
X, y, data_clean = ml_models.prepare_data(enriched_data)

# 时间序列分割
split_point = int(len(X) * 0.8)
X_train, X_test = X.iloc[:split_point], X.iloc[split_point:]
y_train, y_test = y.iloc[:split_point], y.iloc[split_point:]

print(f"训练集大小: {len(X_train)}")
print(f"测试集大小: {len(X_test)}")

# 特征选择和缩放
feature_engineer = FinancialFeatureEngineering()
X_train_selected, selected_features = feature_engineer.select_features(
    X_train, y_train, method='mutual_info', k=20
)
X_test_selected = X_test[selected_features]

X_train_scaled, X_test_scaled = feature_engineer.scale_features(
    X_train_selected, X_test_selected, method='robust'
)

# 转换回DataFrame以保持列名
X_train_final = pd.DataFrame(X_train_scaled, columns=selected_features, index=X_train.index)
X_test_final = pd.DataFrame(X_test_scaled, columns=selected_features, index=X_test.index)

# 训练模型
print("\n训练线性模型...")
linear_results = ml_models.train_linear_models(X_train_final, y_train, X_test_final, y_test)

print("\n训练树模型...")
tree_results = ml_models.train_tree_models(X_train_final, y_train, X_test_final, y_test)

# 合并结果
all_results = {**linear_results, **tree_results}

# 显示结果
ml_models.display_model_comparison(all_results)

# 显示特征重要性
ml_models.display_feature_importance('XGBoost')
```

## 3.7 策略回测框架

### 3.7.1 回测引擎设计

```python
import pandas as pd
import numpy as np
from datetime import datetime, timedelta
import matplotlib.pyplot as plt

class BacktestEngine:
    """策略回测引擎"""

    def __init__(self, initial_capital=1000000, commission=0.0003):
        self.initial_capital = initial_capital
        self.commission = commission
        self.current_capital = initial_capital
        self.positions = {}  # 持仓
        self.trades = []  # 交易记录
        self.daily_returns = []  # 日收益率
        self.portfolio_values = []  # 组合价值
        self.timestamps = []  # 时间戳

    def reset(self):
        """重置回测状态"""
        self.current_capital = self.initial_capital
        self.positions = {}
        self.trades = []
        self.daily_returns = []
        self.portfolio_values = []
        self.timestamps = []

    def add_order(self, timestamp, symbol, quantity, price, order_type='market'):
        """添加订单"""
        order = {
            'timestamp': timestamp,
            'symbol': symbol,
            'quantity': quantity,
            'price': price,
            'order_type': order_type,
            'commission': abs(quantity * price * self.commission)
        }
        return order

    def execute_order(self, order):
        """执行订单"""
        symbol = order['symbol']
        quantity = order['quantity']
        price = order['price']
        commission = order['commission']

        # 计算交易价值
        trade_value = quantity * price

        # 检查资金是否足够（买入时）
        if quantity > 0:  # 买入
            total_cost = trade_value + commission
            if total_cost > self.current_capital:
                print(f"资金不足，无法买入 {symbol}")
                return False

            self.current_capital -= total_cost
        else:  # 卖出
            if symbol not in self.positions:
                print(f"没有持仓 {symbol}，无法卖出")
                return False

            if abs(quantity) > self.positions[symbol]:
                print(f"持仓不足，无法卖出 {symbol}")
                return False

            self.current_capital += abs(trade_value) - commission

        # 更新持仓
        if symbol not in self.positions:
            self.positions[symbol] = 0

        self.positions[symbol] += quantity

        # 如果持仓为0，删除该股票
        if self.positions[symbol] == 0:
            del self.positions[symbol]

        # 记录交易
        trade = {
            'timestamp': order['timestamp'],
            'symbol': symbol,
            'quantity': quantity,
            'price': price,
            'trade_value': trade_value,
            'commission': commission,
            'cash_after_trade': self.current_capital
        }

        self.trades.append(trade)
        return True

    def update_portfolio_value(self, timestamp, price_data):
        """更新组合价值"""
        # 计算持仓市值
        position_value = 0
        for symbol, quantity in self.positions.items():
            if symbol in price_data:
                position_value += quantity * price_data[symbol]

        # 总资产 = 现金 + 持仓市值
        total_value = self.current_capital + position_value

        self.timestamps.append(timestamp)
        self.portfolio_values.append(total_value)

        # 计算日收益率
        if len(self.portfolio_values) > 1:
            daily_return = (total_value - self.portfolio_values[-2]) / self.portfolio_values[-2]
            self.daily_returns.append(daily_return)

        return total_value

    def get_performance_metrics(self):
        """计算绩效指标"""
        if len(self.portfolio_values) < 2:
            return {}

        # 转换为numpy数组方便计算
        returns = np.array(self.daily_returns)
        portfolio_values = np.array(self.portfolio_values)

        # 总收益率
        total_return = (portfolio_values[-1] - portfolio_values[0]) / portfolio_values[0]

        # 年化收益率
        trading_days = len(portfolio_values)
        annual_return = (portfolio_values[-1] / portfolio_values[0]) ** (252 / trading_days) - 1

        # 年化波动率
        annual_volatility = np.std(returns) * np.sqrt(252)

        # 夏普比率（假设无风险利率为3%）
        risk_free_rate = 0.03
        sharpe_ratio = (annual_return - risk_free_rate) / annual_volatility if annual_volatility > 0 else 0

        # 最大回撤
        cumulative_returns = portfolio_values / portfolio_values[0]
        running_max = np.maximum.accumulate(cumulative_returns)
        drawdowns = (cumulative_returns - running_max) / running_max
        max_drawdown = np.min(drawdowns)

        # 胜率
        winning_trades = len([t for t in self.trades if t['quantity'] * t['price'] > 0])  # 简化计算
        total_trades = len(self.trades)
        win_rate = winning_trades / total_trades if total_trades > 0 else 0

        # Calmar比率
        calmar_ratio = annual_return / abs(max_drawdown) if max_drawdown != 0 else 0

        return {
            'total_return': total_return,
            'annual_return': annual_return,
            'annual_volatility': annual_volatility,
            'sharpe_ratio': sharpe_ratio,
            'max_drawdown': max_drawdown,
            'calmar_ratio': calmar_ratio,
            'win_rate': win_rate,
            'total_trades': total_trades,
            'final_value': portfolio_values[-1],
            'start_value': portfolio_values[0]
        }

    def plot_performance(self):
        """绘制回测结果"""
        if len(self.portfolio_values) < 2:
            print("数据不足，无法绘图")
            return

        fig, axes = plt.subplots(3, 1, figsize=(12, 10))

        # 组合净值曲线
        cumulative_returns = np.array(self.portfolio_values) / self.portfolio_values[0]
        axes[0].plot(self.timestamps, cumulative_returns)
        axes[0].set_title('组合净值曲线')
        axes[0].set_ylabel('累计收益率')
        axes[0].grid(True, alpha=0.3)

        # 日收益率分布
        if len(self.daily_returns) > 0:
            axes[1].hist(self.daily_returns, bins=50, alpha=0.7)
            axes[1].set_title('日收益率分布')
            axes[1].set_xlabel('日收益率')
            axes[1].set_ylabel('频次')
            axes[1].grid(True, alpha=0.3)

        # 回撤曲线
        running_max = np.maximum.accumulate(cumulative_returns)
        drawdowns = (cumulative_returns - running_max) / running_max
        axes[2].fill_between(self.timestamps, drawdowns, 0, alpha=0.3, color='red')
        axes[2].set_title('回撤曲线')
        axes[2].set_ylabel('回撤')
        axes[2].set_xlabel('时间')
        axes[2].grid(True, alpha=0.3)

        plt.tight_layout()
        plt.show()

    def get_trade_analysis(self):
        """交易分析"""
        if not self.trades:
            return {}

        trades_df = pd.DataFrame(self.trades)

        analysis = {
            'total_trades': len(trades_df),
            'total_commission': trades_df['commission'].sum(),
            'average_trade_size': trades_df['trade_value'].abs().mean(),
            'largest_trade': trades_df['trade_value'].abs().max(),
            'trade_frequency': len(trades_df) / len(self.timestamps) if self.timestamps else 0
        }

        # 按股票分组的交易统计
        symbol_stats = trades_df.groupby('symbol').agg({
            'quantity': 'sum',
            'trade_value': 'sum',
            'commission': 'sum'
        }).round(2)

        analysis['symbol_statistics'] = symbol_stats

        return analysis

class Strategy:
    """策略基类"""

    def __init__(self, name):
        self.name = name

    def generate_signals(self, data):
        """生成交易信号"""
        raise NotImplementedError("子类必须实现generate_signals方法")

class MovingAverageCrossStrategy(Strategy):
    """移动平均交叉策略"""

    def __init__(self, short_window=20, long_window=50):
        super().__init__("移动平均交叉策略")
        self.short_window = short_window
        self.long_window = long_window

    def generate_signals(self, data):
        """生成交易信号"""
        signals = pd.DataFrame(index=data.index)
        signals['price'] = data['Close']

        # 计算移动平均
        signals['short_ma'] = data['Close'].rolling(window=self.short_window).mean()
        signals['long_ma'] = data['Close'].rolling(window=self.long_window).mean()

        # 生成信号
        signals['signal'] = 0
        signals['signal'][self.short_window:] = np.where(
            signals['short_ma'][self.short_window:] > signals['long_ma'][self.short_window:], 1, 0
        )

        # 交易信号（信号变化时）
        signals['positions'] = signals['signal'].diff()

        return signals

# 示例回测
def run_backtest_example():
    """运行回测示例"""
    # 创建模拟数据
    np.random.seed(42)
    dates = pd.date_range('2020-01-01', periods=500, freq='D')
    price_data = pd.DataFrame({
        'Close': 100 + np.cumsum(np.random.randn(500) * 0.02)
    }, index=dates)

    # 创建策略
    strategy = MovingAverageCrossStrategy(short_window=20, long_window=50)

    # 创建回测引擎
    backtest = BacktestEngine(initial_capital=1000000, commission=0.0003)

    # 生成信号
    signals = strategy.generate_signals(price_data)

    print(f"=== {strategy.name} 回测结果 ===")

    # 执行回测
    position = 0  # 当前持仓

    for date, row in signals.iterrows():
        price = row['price']
        signal_change = row['positions']

        # 更新组合价值
        current_prices = {'STOCK': price}
        backtest.update_portfolio_value(date, current_prices)

        # 处理交易信号
        if signal_change == 1:  # 买入信号
            if position == 0:  # 当前空仓，买入
                shares_to_buy = int(backtest.current_capital * 0.95 / price)  # 95%仓位
                if shares_to_buy > 0:
                    order = backtest.add_order(date, 'STOCK', shares_to_buy, price)
                    if backtest.execute_order(order):
                        position = shares_to_buy

        elif signal_change == -1:  # 卖出信号
            if position > 0:  # 当前有持仓，卖出
                order = backtest.add_order(date, 'STOCK', -position, price)
                if backtest.execute_order(order):
                    position = 0

    # 计算绩效指标
    metrics = backtest.get_performance_metrics()

    print(f"初始资金: {backtest.initial_capital:,.0f}")
    print(f"最终价值: {metrics.get('final_value', 0):,.0f}")
    print(f"总收益率: {metrics.get('total_return', 0):.2%}")
    print(f"年化收益率: {metrics.get('annual_return', 0):.2%}")
    print(f"年化波动率: {metrics.get('annual_volatility', 0):.2%}")
    print(f"夏普比率: {metrics.get('sharpe_ratio', 0):.2f}")
    print(f"最大回撤: {metrics.get('max_drawdown', 0):.2%}")
    print(f"Calmar比率: {metrics.get('calmar_ratio', 0):.2f}")
    print(f"交易次数: {metrics.get('total_trades', 0)}")

    # 交易分析
    trade_analysis = backtest.get_trade_analysis()
    print(f"\n交易分析:")
    print(f"总手续费: {trade_analysis.get('total_commission', 0):.2f}")
    print(f"平均交易规模: {trade_analysis.get('average_trade_size', 0):,.0f}")

    # 绘制结果
    backtest.plot_performance()

    return backtest, signals

# 运行示例
backtest_result, strategy_signals = run_backtest_example()
```

### 3.7.2 多策略回测框架

```python
class MultiStrategyBacktest:
    """多策略回测框架"""

    def __init__(self, initial_capital=1000000, commission=0.0003):
        self.initial_capital = initial_capital
        self.commission = commission
        self.strategies = {}
        self.results = {}

    def add_strategy(self, strategy, allocation_weight):
        """添加策略"""
        self.strategies[strategy.name] = {
            'strategy': strategy,
            'weight': allocation_weight,
            'engine': BacktestEngine(
                initial_capital=self.initial_capital * allocation_weight,
                commission=self.commission
            )
        }

    def run_backtest(self, data_dict):
        """运行多策略回测"""
        print("开始多策略回测...")

        for strategy_name, strategy_info in self.strategies.items():
            print(f"\n回测策略: {strategy_name}")

            strategy = strategy_info['strategy']
            engine = strategy_info['engine']

            # 为每个策略运行回测
            if hasattr(strategy, 'symbols'):
                # 多资产策略
                symbols = strategy.symbols
            else:
                # 单资产策略，使用第一个数据
                symbols = [list(data_dict.keys())[0]]

            # 生成信号并执行交易
            self._run_single_strategy(strategy, engine, data_dict, symbols)

            # 保存结果
            self.results[strategy_name] = {
                'engine': engine,
                'metrics': engine.get_performance_metrics(),
                'trades': engine.get_trade_analysis()
            }

    def _run_single_strategy(self, strategy, engine, data_dict, symbols):
        """运行单个策略的回测"""
        # 这里简化处理，实际应该根据策略类型调整
        main_symbol = symbols[0]
        data = data_dict[main_symbol]

        signals = strategy.generate_signals(data)
        position = 0

        for date, row in signals.iterrows():
            if main_symbol not in data_dict:
                continue

            price = data.loc[date, 'Close'] if date in data.index else None
            if price is None:
                continue

            # 更新组合价值
            current_prices = {main_symbol: price}
            engine.update_portfolio_value(date, current_prices)

            # 处理交易信号
            if 'positions' in row and not pd.isna(row['positions']):
                signal_change = row['positions']

                if signal_change == 1 and position == 0:
                    # 买入
                    shares_to_buy = int(engine.current_capital * 0.95 / price)
                    if shares_to_buy > 0:
                        order = engine.add_order(date, main_symbol, shares_to_buy, price)
                        if engine.execute_order(order):
                            position = shares_to_buy

                elif signal_change == -1 and position > 0:
                    # 卖出
                    order = engine.add_order(date, main_symbol, -position, price)
                    if engine.execute_order(order):
                        position = 0

    def compare_strategies(self):
        """比较策略性能"""
        if not self.results:
            print("没有回测结果可比较")
            return

        print("\n=== 策略性能比较 ===")

        comparison_data = []
        for strategy_name, result in self.results.items():
            metrics = result['metrics']
            comparison_data.append({
                '策略': strategy_name,
                '总收益率': metrics.get('total_return', 0),
                '年化收益率': metrics.get('annual_return', 0),
                '年化波动率': metrics.get('annual_volatility', 0),
                '夏普比率': metrics.get('sharpe_ratio', 0),
                '最大回撤': metrics.get('max_drawdown', 0),
                'Calmar比率': metrics.get('calmar_ratio', 0),
                '交易次数': metrics.get('total_trades', 0)
            })

        comparison_df = pd.DataFrame(comparison_data)
        comparison_df = comparison_df.set_index('策略')

        # 格式化显示
        pd.set_option('display.float_format', '{:.4f}'.format)
        print(comparison_df)
        pd.reset_option('display.float_format')

        return comparison_df

    def plot_strategy_comparison(self):
        """绘制策略比较图"""
        if not self.results:
            print("没有回测结果可绘制")
            return

        fig, axes = plt.subplots(2, 2, figsize=(15, 10))

        # 净值曲线比较
        for strategy_name, result in self.results.items():
            engine = result['engine']
            if len(engine.portfolio_values) > 0:
                cumulative_returns = np.array(engine.portfolio_values) / engine.portfolio_values[0]
                axes[0, 0].plot(engine.timestamps, cumulative_returns, label=strategy_name)

        axes[0, 0].set_title('策略净值曲线比较')
        axes[0, 0].set_ylabel('累计收益率')
        axes[0, 0].legend()
        axes[0, 0].grid(True, alpha=0.3)

        # 收益率对比
        returns = []
        names = []
        for strategy_name, result in self.results.items():
            metrics = result['metrics']
            returns.append(metrics.get('annual_return', 0))
            names.append(strategy_name)

        axes[0, 1].bar(names, returns)
        axes[0, 1].set_title('年化收益率对比')
        axes[0, 1].set_ylabel('年化收益率')
        axes[0, 1].tick_params(axis='x', rotation=45)

        # 夏普比率对比
        sharpe_ratios = []
        for strategy_name, result in self.results.items():
            metrics = result['metrics']
            sharpe_ratios.append(metrics.get('sharpe_ratio', 0))

        axes[1, 0].bar(names, sharpe_ratios)
        axes[1, 0].set_title('夏普比率对比')
        axes[1, 0].set_ylabel('夏普比率')
        axes[1, 0].tick_params(axis='x', rotation=45)

        # 最大回撤对比
        max_drawdowns = []
        for strategy_name, result in self.results.items():
            metrics = result['metrics']
            max_drawdowns.append(abs(metrics.get('max_drawdown', 0)))

        axes[1, 1].bar(names, max_drawdowns)
        axes[1, 1].set_title('最大回撤对比')
        axes[1, 1].set_ylabel('最大回撤')
        axes[1, 1].tick_params(axis='x', rotation=45)

        plt.tight_layout()
        plt.show()

# 示例：多策略回测
class RSIStrategy(Strategy):
    """RSI策略"""

    def __init__(self, period=14, oversold=30, overbought=70):
        super().__init__("RSI策略")
        self.period = period
        self.oversold = oversold
        self.overbought = overbought

    def generate_signals(self, data):
        """生成RSI交易信号"""
        signals = pd.DataFrame(index=data.index)
        signals['price'] = data['Close']

        # 计算RSI
        delta = data['Close'].diff()
        gain = (delta.where(delta > 0, 0)).rolling(window=self.period).mean()
        loss = (-delta.where(delta < 0, 0)).rolling(window=self.period).mean()
        rs = gain / loss
        signals['rsi'] = 100 - (100 / (1 + rs))

        # 生成信号
        signals['signal'] = 0
        signals.loc[signals['rsi'] < self.oversold, 'signal'] = 1  # 买入
        signals.loc[signals['rsi'] > self.overbought, 'signal'] = -1  # 卖出

        # 交易信号
        signals['positions'] = signals['signal'].diff()

        return signals

def run_multi_strategy_example():
    """运行多策略回测示例"""
    # 创建模拟数据
    np.random.seed(42)
    dates = pd.date_range('2020-01-01', periods=500, freq='D')
    stock_data = pd.DataFrame({
        'Close': 100 + np.cumsum(np.random.randn(500) * 0.02)
    }, index=dates)

    # 创建多策略回测框架
    multi_backtest = MultiStrategyBacktest(initial_capital=1000000)

    # 添加策略
    ma_strategy = MovingAverageCrossStrategy(short_window=10, long_window=30)
    rsi_strategy = RSIStrategy(period=14, oversold=30, overbought=70)

    multi_backtest.add_strategy(ma_strategy, allocation_weight=0.5)
    multi_backtest.add_strategy(rsi_strategy, allocation_weight=0.5)

    # 准备数据
    data_dict = {'STOCK': stock_data}

    # 运行回测
    multi_backtest.run_backtest(data_dict)

    # 比较结果
    comparison = multi_backtest.compare_strategies()

    # 绘制比较图
    multi_backtest.plot_strategy_comparison()

    return multi_backtest

# 运行多策略示例
multi_result = run_multi_strategy_example()
```

## 📚 本章小结

在本章下篇中，我们深入学习了 Python 金融编程的高级内容，包括：

1. **机器学习应用**：掌握了特征工程和机器学习模型在金融预测中的应用
2. **策略回测框架**：学会了构建专业的回测引擎和多策略比较系统

这些高级技能将帮助您从基础的数据分析进阶到专业的量化交易系统开发。

## 🔗 下一章预告

在[第四章：股票投资与分析（下篇）](04_stock_investment_analysis_part2.md)中，我们将学习量化选股策略、多因子模型、事件驱动策略以及 A 股市场的特色分析。

## 💡 实践练习

1. **特征工程实践**：为真实股票数据创建完整的技术和基本面特征集

2. **机器学习模型**：比较不同机器学习算法在股票收益预测中的表现

3. **回测框架**：实现一个完整的策略回测系统，支持多资产交易

4. **策略开发**：开发并回测一个自定义的量化交易策略

5. **多策略组合**：构建多策略投资组合并分析其风险收益特征

---

**建议学习时间：** 4-5 天  
**前置章节：** [第三章：Python 金融编程（上篇）](03_python_finance_programming.md)  
**下一章：** [股票投资与分析（下篇）](04_stock_investment_analysis_part2.md)
