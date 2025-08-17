# 第8章 机器学习在量化投资中的应用（上篇）

## 📚 本章概览

机器学习技术正在深刻改变量化投资的面貌，从传统的线性模型到复杂的非线性模型，从规则基础的策略到数据驱动的策略。本章将深入探讨监督学习、无监督学习和强化学习在量化投资中的具体应用，包括特征工程、模型选择、验证方法等关键技术。

### 🎯 学习目标

- 理解机器学习在金融中的应用场景和挑战
- 掌握监督学习在收益预测中的应用
- 学会无监督学习在市场聚类和异常检测中的应用
- 了解强化学习在交易决策中的应用
- 掌握金融数据的特征工程技巧

## 第一节：机器学习在金融中的基础

### 1.1 金融机器学习的特殊性

```python
import numpy as np
import pandas as pd
from sklearn.model_selection import TimeSeriesSplit, cross_val_score
from sklearn.metrics import accuracy_score, precision_score, recall_score, f1_score
from sklearn.preprocessing import StandardScaler, RobustScaler
import matplotlib.pyplot as plt
import seaborn as sns
from typing import Tuple, List, Dict, Any

class FinancialMLFramework:
    """
    金融机器学习框架基类
    """
    
    def __init__(self, lookback_window: int = 252):
        self.lookback_window = lookback_window
        self.scaler = None
        self.models = {}
        self.feature_importance = {}
        
    def prepare_financial_data(self, price_data: pd.DataFrame, 
                             volume_data: pd.DataFrame = None,
                             fundamental_data: pd.DataFrame = None) -> pd.DataFrame:
        """
        准备金融数据
        """
        
        features = pd.DataFrame(index=price_data.index)
        
        # 价格相关特征
        returns = price_data.pct_change()
        features['returns'] = returns.mean(axis=1) if len(price_data.columns) > 1 else returns.iloc[:, 0]
        
        # 技术指标特征
        features = self._add_technical_features(features, price_data, volume_data)
        
        # 基本面特征
        if fundamental_data is not None:
            features = self._add_fundamental_features(features, fundamental_data)
        
        # 宏观特征
        features = self._add_macro_features(features, price_data)
        
        return features.dropna()
    
    def _add_technical_features(self, features: pd.DataFrame, 
                              price_data: pd.DataFrame,
                              volume_data: pd.DataFrame = None) -> pd.DataFrame:
        """
        添加技术指标特征
        """
        
        if len(price_data.columns) == 1:
            price = price_data.iloc[:, 0]
        else:
            price = price_data.mean(axis=1)  # 如果多个资产，取平均
        
        returns = price.pct_change()
        
        # 移动平均
        for window in [5, 10, 20, 50]:
            features[f'sma_{window}'] = price.rolling(window).mean()
            features[f'price_vs_sma_{window}'] = price / features[f'sma_{window}'] - 1
        
        # 指数移动平均
        for span in [12, 26]:
            features[f'ema_{span}'] = price.ewm(span=span).mean()
        
        # MACD
        ema_12 = price.ewm(span=12).mean()
        ema_26 = price.ewm(span=26).mean()
        features['macd'] = ema_12 - ema_26
        features['macd_signal'] = features['macd'].ewm(span=9).mean()
        features['macd_histogram'] = features['macd'] - features['macd_signal']
        
        # RSI
        delta = price.diff()
        gain = (delta.where(delta > 0, 0)).rolling(window=14).mean()
        loss = (-delta.where(delta < 0, 0)).rolling(window=14).mean()
        rs = gain / loss
        features['rsi'] = 100 - (100 / (1 + rs))
        
        # 布林带
        sma_20 = price.rolling(20).mean()
        std_20 = price.rolling(20).std()
        features['bollinger_upper'] = sma_20 + (std_20 * 2)
        features['bollinger_lower'] = sma_20 - (std_20 * 2)
        features['bollinger_position'] = (price - features['bollinger_lower']) / (features['bollinger_upper'] - features['bollinger_lower'])
        
        # 波动率特征
        for window in [10, 20, 60]:
            features[f'volatility_{window}'] = returns.rolling(window).std() * np.sqrt(252)
        
        # 动量特征
        for window in [5, 10, 20, 60]:
            features[f'momentum_{window}'] = returns.rolling(window).sum()
        
        # 成交量特征（如果有成交量数据）
        if volume_data is not None:
            volume = volume_data.mean(axis=1) if len(volume_data.columns) > 1 else volume_data.iloc[:, 0]
            features['volume_sma_20'] = volume.rolling(20).mean()
            features['volume_ratio'] = volume / features['volume_sma_20']
            
            # 成交量价格趋势
            features['vpt'] = (volume * (price.pct_change())).cumsum()
            
            # 资金流量指数
            typical_price = price  # 简化处理
            money_flow = typical_price * volume
            features['money_flow_20'] = money_flow.rolling(20).sum()
        
        return features
    
    def _add_fundamental_features(self, features: pd.DataFrame,
                                 fundamental_data: pd.DataFrame) -> pd.DataFrame:
        """
        添加基本面特征
        """
        
        # 估值指标
        if 'pe_ratio' in fundamental_data.columns:
            features['pe_ratio'] = fundamental_data['pe_ratio']
            features['pe_percentile'] = fundamental_data['pe_ratio'].rolling(252).rank(pct=True)
        
        if 'pb_ratio' in fundamental_data.columns:
            features['pb_ratio'] = fundamental_data['pb_ratio']
            features['pb_percentile'] = fundamental_data['pb_ratio'].rolling(252).rank(pct=True)
        
        # 财务健康度
        if 'debt_to_equity' in fundamental_data.columns:
            features['debt_to_equity'] = fundamental_data['debt_to_equity']
        
        if 'roe' in fundamental_data.columns:
            features['roe'] = fundamental_data['roe']
            features['roe_trend'] = fundamental_data['roe'].diff(4)  # 年度变化
        
        return features
    
    def _add_macro_features(self, features: pd.DataFrame,
                           price_data: pd.DataFrame) -> pd.DataFrame:
        """
        添加宏观特征
        """
        
        # 市场整体趋势（如果有多个资产）
        if len(price_data.columns) > 1:
            market_returns = price_data.pct_change().mean(axis=1)
            features['market_momentum_20'] = market_returns.rolling(20).sum()
            features['market_volatility_20'] = market_returns.rolling(20).std() * np.sqrt(252)
            
            # 市场分散度
            individual_returns = price_data.pct_change()
            features['market_dispersion'] = individual_returns.std(axis=1)
            
            # 相关性特征
            correlation_matrix = individual_returns.rolling(60).corr()
            features['avg_correlation'] = correlation_matrix.mean().mean()
        
        # 时间特征
        features['month'] = features.index.month
        features['quarter'] = features.index.quarter
        features['day_of_week'] = features.index.dayofweek
        features['day_of_month'] = features.index.day
        
        # 季节性特征
        features['is_month_end'] = (features.index.day > 25).astype(int)
        features['is_quarter_end'] = features.index.is_quarter_end.astype(int)
        features['is_year_end'] = features.index.is_year_end.astype(int)
        
        return features
    
    def handle_financial_data_issues(self, data: pd.DataFrame) -> pd.DataFrame:
        """
        处理金融数据特有问题
        """
        
        # 1. 异常值处理
        data_cleaned = data.copy()
        
        # 使用分位数方法处理异常值
        for column in data_cleaned.select_dtypes(include=[np.number]).columns:
            Q1 = data_cleaned[column].quantile(0.01)
            Q3 = data_cleaned[column].quantile(0.99)
            data_cleaned[column] = data_cleaned[column].clip(lower=Q1, upper=Q3)
        
        # 2. 缺失值处理
        # 前向填充（适合金融时间序列）
        data_cleaned = data_cleaned.fillna(method='ffill')
        
        # 3. 数据标准化（使用RobustScaler，对异常值更鲁棒）
        numeric_columns = data_cleaned.select_dtypes(include=[np.number]).columns
        if self.scaler is None:
            self.scaler = RobustScaler()
            data_cleaned[numeric_columns] = self.scaler.fit_transform(data_cleaned[numeric_columns])
        else:
            data_cleaned[numeric_columns] = self.scaler.transform(data_cleaned[numeric_columns])
        
        return data_cleaned
    
    def create_labels(self, price_data: pd.DataFrame, method: str = 'return_direction',
                     threshold: float = 0.0, horizon: int = 1) -> pd.Series:
        """
        创建标签
        """
        
        if len(price_data.columns) == 1:
            price = price_data.iloc[:, 0]
        else:
            price = price_data.mean(axis=1)
        
        returns = price.pct_change()
        
        if method == 'return_direction':
            # 收益方向预测
            future_returns = returns.shift(-horizon)
            labels = (future_returns > threshold).astype(int)
        
        elif method == 'return_quantile':
            # 收益分位数预测
            future_returns = returns.shift(-horizon)
            labels = pd.qcut(future_returns, q=3, labels=[0, 1, 2])
        
        elif method == 'volatility_regime':
            # 波动率状态预测
            volatility = returns.rolling(20).std() * np.sqrt(252)
            future_vol = volatility.shift(-horizon)
            vol_median = volatility.rolling(252).median()
            labels = (future_vol > vol_median).astype(int)
        
        elif method == 'drawdown_prediction':
            # 回撤预测
            cumulative_returns = (1 + returns).cumprod()
            running_max = cumulative_returns.expanding().max()
            drawdown = (cumulative_returns - running_max) / running_max
            future_drawdown = drawdown.shift(-horizon)
            labels = (future_drawdown < -0.05).astype(int)  # 5%回撤阈值
        
        return labels.dropna()
```

### 1.2 时间序列交叉验证

```python
class FinancialTimeSeriesSplit:
    """
    金融时间序列专用交叉验证
    """
    
    def __init__(self, n_splits: int = 5, test_size: int = 252, gap: int = 0):
        self.n_splits = n_splits
        self.test_size = test_size
        self.gap = gap
    
    def split(self, X: pd.DataFrame, y: pd.Series = None):
        """
        时间序列分割
        """
        
        n_samples = len(X)
        
        # 计算每个fold的大小
        fold_size = (n_samples - self.test_size) // self.n_splits
        
        for i in range(self.n_splits):
            # 训练集结束位置
            train_end = fold_size * (i + 1)
            
            # 测试集开始位置（考虑gap）
            test_start = train_end + self.gap
            test_end = test_start + self.test_size
            
            if test_end > n_samples:
                break
            
            train_indices = list(range(0, train_end))
            test_indices = list(range(test_start, test_end))
            
            yield train_indices, test_indices
    
    def purged_cross_validation(self, X: pd.DataFrame, y: pd.Series, 
                               model, purge_window: int = 5):
        """
        净化交叉验证（避免数据泄露）
        """
        
        scores = []
        
        for train_idx, test_idx in self.split(X, y):
            # 净化：移除测试集前后purge_window天的训练数据
            purged_train_idx = [
                idx for idx in train_idx 
                if not any(abs(idx - test_i) <= purge_window for test_i in test_idx)
            ]
            
            X_train = X.iloc[purged_train_idx]
            y_train = y.iloc[purged_train_idx]
            X_test = X.iloc[test_idx]
            y_test = y.iloc[test_idx]
            
            # 训练和预测
            model.fit(X_train, y_train)
            y_pred = model.predict(X_test)
            
            # 计算评分
            score = accuracy_score(y_test, y_pred)
            scores.append(score)
        
        return np.array(scores)

def combinatorial_purged_cross_validation(X: pd.DataFrame, y: pd.Series, 
                                        model, n_splits: int = 5, 
                                        n_test_groups: int = 2):
    """
    组合净化交叉验证
    """
    
    from itertools import combinations
    
    # 将数据分组
    group_size = len(X) // n_splits
    groups = [list(range(i * group_size, (i + 1) * group_size)) for i in range(n_splits)]
    
    scores = []
    
    # 组合选择测试组
    for test_groups in combinations(range(n_splits), n_test_groups):
        test_indices = []
        for group_idx in test_groups:
            test_indices.extend(groups[group_idx])
        
        train_indices = []
        for group_idx in range(n_splits):
            if group_idx not in test_groups:
                # 净化：不包含测试组相邻的组
                adjacent_to_test = any(abs(group_idx - test_group) <= 1 for test_group in test_groups)
                if not adjacent_to_test:
                    train_indices.extend(groups[group_idx])
        
        if len(train_indices) > 0 and len(test_indices) > 0:
            X_train = X.iloc[train_indices]
            y_train = y.iloc[train_indices]
            X_test = X.iloc[test_indices]
            y_test = y.iloc[test_indices]
            
            model.fit(X_train, y_train)
            y_pred = model.predict(X_test)
            
            score = accuracy_score(y_test, y_pred)
            scores.append(score)
    
    return np.array(scores)
```

## 第二节：监督学习在收益预测中的应用

### 2.1 分类模型：方向预测

```python
from sklearn.ensemble import RandomForestClassifier, GradientBoostingClassifier
from sklearn.linear_model import LogisticRegression
from sklearn.svm import SVC
from sklearn.metrics import classification_report, confusion_matrix
import xgboost as xgb
import lightgbm as lgb

class ReturnDirectionPredictor:
    """
    收益方向预测器
    """
    
    def __init__(self, models_config: Dict[str, Any] = None):
        if models_config is None:
            self.models_config = {
                'logistic': LogisticRegression(random_state=42),
                'random_forest': RandomForestClassifier(n_estimators=100, random_state=42),
                'gradient_boosting': GradientBoostingClassifier(random_state=42),
                'xgboost': xgb.XGBClassifier(random_state=42),
                'lightgbm': lgb.LGBMClassifier(random_state=42)
            }
        else:
            self.models_config = models_config
        
        self.fitted_models = {}
        self.feature_importance = {}
        
    def train_models(self, X_train: pd.DataFrame, y_train: pd.Series):
        """
        训练多个模型
        """
        
        for model_name, model in self.models_config.items():
            print(f"训练 {model_name} 模型...")
            
            model.fit(X_train, y_train)
            self.fitted_models[model_name] = model
            
            # 提取特征重要性
            if hasattr(model, 'feature_importances_'):
                self.feature_importance[model_name] = pd.Series(
                    model.feature_importances_, 
                    index=X_train.columns
                ).sort_values(ascending=False)
    
    def evaluate_models(self, X_test: pd.DataFrame, y_test: pd.Series) -> Dict[str, Dict]:
        """
        评估模型性能
        """
        
        results = {}
        
        for model_name, model in self.fitted_models.items():
            y_pred = model.predict(X_test)
            y_pred_proba = model.predict_proba(X_test)[:, 1] if hasattr(model, 'predict_proba') else None
            
            # 基本指标
            accuracy = accuracy_score(y_test, y_pred)
            precision = precision_score(y_test, y_pred)
            recall = recall_score(y_test, y_pred)
            f1 = f1_score(y_test, y_pred)
            
            # 金融特有指标
            long_positions = y_pred == 1
            if np.sum(long_positions) > 0:
                hit_rate_long = np.sum((y_test == 1) & long_positions) / np.sum(long_positions)
            else:
                hit_rate_long = 0
            
            short_positions = y_pred == 0
            if np.sum(short_positions) > 0:
                hit_rate_short = np.sum((y_test == 0) & short_positions) / np.sum(short_positions)
            else:
                hit_rate_short = 0
            
            results[model_name] = {
                'accuracy': accuracy,
                'precision': precision,
                'recall': recall,
                'f1_score': f1,
                'hit_rate_long': hit_rate_long,
                'hit_rate_short': hit_rate_short,
                'predictions': y_pred,
                'probabilities': y_pred_proba
            }
        
        return results
    
    def analyze_feature_importance(self, top_n: int = 20):
        """
        分析特征重要性
        """
        
        fig, axes = plt.subplots(2, 3, figsize=(18, 12))
        axes = axes.flatten()
        
        for i, (model_name, importance) in enumerate(self.feature_importance.items()):
            if i < len(axes):
                importance.head(top_n).plot(kind='barh', ax=axes[i])
                axes[i].set_title(f'{model_name} - Top {top_n} Features')
                axes[i].set_xlabel('Importance')
        
        plt.tight_layout()
        plt.show()
        
        # 计算特征重要性的平均排名
        all_features = set()
        for importance in self.feature_importance.values():
            all_features.update(importance.index)
        
        feature_ranks = pd.DataFrame(index=list(all_features))
        
        for model_name, importance in self.feature_importance.items():
            ranks = importance.rank(ascending=False)
            feature_ranks[model_name] = ranks
        
        feature_ranks['mean_rank'] = feature_ranks.mean(axis=1)
        feature_ranks = feature_ranks.sort_values('mean_rank')
        
        return feature_ranks

def trading_performance_analysis(y_true: pd.Series, y_pred: np.ndarray, 
                               returns: pd.Series) -> Dict[str, float]:
    """
    交易绩效分析
    """
    
    # 计算策略收益
    strategy_returns = []
    positions = []
    
    for i in range(len(y_pred)):
        if y_pred[i] == 1:  # 预测上涨，做多
            position = 1
        else:  # 预测下跌，做空或持现金
            position = 0  # 这里简化为持现金
        
        positions.append(position)
        strategy_return = position * returns.iloc[i]
        strategy_returns.append(strategy_return)
    
    strategy_returns = pd.Series(strategy_returns, index=returns.index[:len(y_pred)])
    
    # 计算绩效指标
    total_return = (1 + strategy_returns).prod() - 1
    annual_return = (1 + strategy_returns.mean()) ** 252 - 1
    annual_volatility = strategy_returns.std() * np.sqrt(252)
    sharpe_ratio = annual_return / annual_volatility if annual_volatility > 0 else 0
    
    # 最大回撤
    cumulative_returns = (1 + strategy_returns).cumprod()
    running_max = cumulative_returns.expanding().max()
    drawdown = (cumulative_returns - running_max) / running_max
    max_drawdown = drawdown.min()
    
    # 胜率
    win_rate = (strategy_returns > 0).mean()
    
    # 盈亏比
    winning_trades = strategy_returns[strategy_returns > 0]
    losing_trades = strategy_returns[strategy_returns < 0]
    
    if len(losing_trades) > 0:
        profit_loss_ratio = winning_trades.mean() / abs(losing_trades.mean())
    else:
        profit_loss_ratio = np.inf
    
    return {
        'total_return': total_return,
        'annual_return': annual_return,
        'annual_volatility': annual_volatility,
        'sharpe_ratio': sharpe_ratio,
        'max_drawdown': max_drawdown,
        'win_rate': win_rate,
        'profit_loss_ratio': profit_loss_ratio,
        'strategy_returns': strategy_returns
    }
```

### 2.2 回归模型：收益预测

```python
from sklearn.ensemble import RandomForestRegressor, GradientBoostingRegressor
from sklearn.linear_model import LinearRegression, Ridge, Lasso, ElasticNet
from sklearn.svm import SVR
from sklearn.metrics import mean_squared_error, mean_absolute_error, r2_score

class ReturnPredictor:
    """
    收益预测器
    """
    
    def __init__(self, models_config: Dict[str, Any] = None):
        if models_config is None:
            self.models_config = {
                'linear': LinearRegression(),
                'ridge': Ridge(alpha=1.0),
                'lasso': Lasso(alpha=0.1),
                'elastic_net': ElasticNet(alpha=0.1, l1_ratio=0.5),
                'random_forest': RandomForestRegressor(n_estimators=100, random_state=42),
                'gradient_boosting': GradientBoostingRegressor(random_state=42),
                'xgboost': xgb.XGBRegressor(random_state=42),
                'lightgbm': lgb.LGBMRegressor(random_state=42)
            }
        else:
            self.models_config = models_config
        
        self.fitted_models = {}
        self.feature_importance = {}
    
    def train_models(self, X_train: pd.DataFrame, y_train: pd.Series):
        """
        训练回归模型
        """
        
        for model_name, model in self.models_config.items():
            print(f"训练 {model_name} 模型...")
            
            model.fit(X_train, y_train)
            self.fitted_models[model_name] = model
            
            # 提取特征重要性或系数
            if hasattr(model, 'feature_importances_'):
                self.feature_importance[model_name] = pd.Series(
                    model.feature_importances_, 
                    index=X_train.columns
                ).sort_values(ascending=False)
            elif hasattr(model, 'coef_'):
                self.feature_importance[model_name] = pd.Series(
                    np.abs(model.coef_), 
                    index=X_train.columns
                ).sort_values(ascending=False)
    
    def evaluate_models(self, X_test: pd.DataFrame, y_test: pd.Series) -> Dict[str, Dict]:
        """
        评估回归模型
        """
        
        results = {}
        
        for model_name, model in self.fitted_models.items():
            y_pred = model.predict(X_test)
            
            # 回归指标
            mse = mean_squared_error(y_test, y_pred)
            mae = mean_absolute_error(y_test, y_pred)
            r2 = r2_score(y_test, y_pred)
            rmse = np.sqrt(mse)
            
            # 方向准确率
            direction_accuracy = np.mean(np.sign(y_test) == np.sign(y_pred))
            
            # 信息系数（IC）
            ic = np.corrcoef(y_test, y_pred)[0, 1]
            
            # 分位数分析
            # 将预测值分为5个分位数，分析实际收益
            pred_quantiles = pd.qcut(y_pred, q=5, labels=['Q1', 'Q2', 'Q3', 'Q4', 'Q5'])
            quantile_analysis = pd.DataFrame({
                'predicted': y_pred,
                'actual': y_test,
                'quantile': pred_quantiles
            }).groupby('quantile')['actual'].agg(['mean', 'std', 'count'])
            
            results[model_name] = {
                'mse': mse,
                'mae': mae,
                'rmse': rmse,
                'r2': r2,
                'direction_accuracy': direction_accuracy,
                'information_coefficient': ic,
                'predictions': y_pred,
                'quantile_analysis': quantile_analysis
            }
        
        return results
    
    def long_short_strategy_analysis(self, X_test: pd.DataFrame, y_test: pd.Series,
                                   returns: pd.Series, top_pct: float = 0.2) -> Dict[str, Dict]:
        """
        多空策略分析
        """
        
        strategy_results = {}
        
        for model_name, model in self.fitted_models.items():
            y_pred = model.predict(X_test)
            
            # 根据预测值构建多空组合
            n_top = int(len(y_pred) * top_pct)
            n_bottom = int(len(y_pred) * top_pct)
            
            # 每期重新排序和选股
            strategy_returns = []
            
            # 这里简化处理，假设是单期预测
            sorted_indices = np.argsort(y_pred)
            
            # 做多前top_pct的股票，做空后top_pct的股票
            long_indices = sorted_indices[-n_top:]
            short_indices = sorted_indices[:n_bottom]
            
            long_return = returns.iloc[long_indices].mean() if len(long_indices) > 0 else 0
            short_return = returns.iloc[short_indices].mean() if len(short_indices) > 0 else 0
            
            # 多空组合收益
            long_short_return = long_return - short_return
            
            strategy_results[model_name] = {
                'long_return': long_return,
                'short_return': short_return,
                'long_short_return': long_short_return,
                'long_indices': long_indices,
                'short_indices': short_indices
            }
        
        return strategy_results

def ensemble_prediction(models_dict: Dict[str, Any], X_test: pd.DataFrame,
                       method: str = 'average') -> np.ndarray:
    """
    集成预测
    """
    
    predictions = []
    weights = []
    
    for model_name, model_info in models_dict.items():
        model = model_info.get('model')
        weight = model_info.get('weight', 1.0)
        
        if model is not None:
            pred = model.predict(X_test)
            predictions.append(pred)
            weights.append(weight)
    
    predictions = np.array(predictions)
    weights = np.array(weights)
    
    if method == 'average':
        # 加权平均
        ensemble_pred = np.average(predictions, axis=0, weights=weights)
    
    elif method == 'median':
        # 中位数
        ensemble_pred = np.median(predictions, axis=0)
    
    elif method == 'rank_average':
        # 排名平均
        ranked_predictions = []
        for pred in predictions:
            ranks = pd.Series(pred).rank().values
            ranked_predictions.append(ranks)
        
        ensemble_ranks = np.average(ranked_predictions, axis=0, weights=weights)
        ensemble_pred = ensemble_ranks
    
    return ensemble_pred
```

## 第三节：无监督学习在市场分析中的应用

### 3.1 聚类分析：市场状态识别

```python
from sklearn.cluster import KMeans, DBSCAN, AgglomerativeClustering
from sklearn.mixture import GaussianMixture
from sklearn.decomposition import PCA
from sklearn.manifold import TSNE
import plotly.express as px
import plotly.graph_objects as go

class MarketRegimeAnalysis:
    """
    市场状态分析
    """
    
    def __init__(self, n_regimes: int = 3):
        self.n_regimes = n_regimes
        self.regime_models = {}
        self.regime_labels = None
        self.regime_characteristics = None
    
    def prepare_regime_features(self, price_data: pd.DataFrame,
                              volume_data: pd.DataFrame = None) -> pd.DataFrame:
        """
        准备状态识别特征
        """
        
        returns = price_data.pct_change().dropna()
        
        # 市场特征
        features = pd.DataFrame(index=returns.index)
        
        # 收益率特征
        features['mean_return'] = returns.mean(axis=1)
        features['return_volatility'] = returns.std(axis=1)
        features['return_skewness'] = returns.skew(axis=1)
        features['return_kurtosis'] = returns.kurtosis(axis=1)
        
        # 滚动统计特征
        windows = [20, 60, 120]
        for window in windows:
            features[f'mean_return_{window}d'] = features['mean_return'].rolling(window).mean()
            features[f'volatility_{window}d'] = features['return_volatility'].rolling(window).mean()
            features[f'vol_of_vol_{window}d'] = features['return_volatility'].rolling(window).std()
        
        # 相关性特征
        if len(price_data.columns) > 1:
            rolling_corr = returns.rolling(60).corr()
            avg_corr = []
            for date in returns.index:
                try:
                    corr_matrix = rolling_corr.loc[date]
                    # 计算平均相关系数（排除对角线）
                    mask = ~np.eye(corr_matrix.shape[0], dtype=bool)
                    avg_corr.append(corr_matrix.values[mask].mean())
                except:
                    avg_corr.append(np.nan)
            
            features['avg_correlation'] = avg_corr
        
        # 趋势特征
        price_avg = price_data.mean(axis=1)
        for window in [20, 50, 200]:
            sma = price_avg.rolling(window).mean()
            features[f'price_vs_sma_{window}'] = (price_avg / sma - 1).fillna(0)
        
        # 成交量特征
        if volume_data is not None:
            volume_avg = volume_data.mean(axis=1)
            features['volume_trend'] = volume_avg.rolling(20).mean()
            features['volume_volatility'] = volume_avg.rolling(20).std()
        
        return features.dropna()
    
    def identify_regimes(self, features: pd.DataFrame, method: str = 'kmeans'):
        """
        识别市场状态
        """
        
        # 标准化特征
        scaler = StandardScaler()
        features_scaled = scaler.fit_transform(features)
        
        if method == 'kmeans':
            model = KMeans(n_clusters=self.n_regimes, random_state=42)
        elif method == 'gmm':
            model = GaussianMixture(n_components=self.n_regimes, random_state=42)
        elif method == 'hierarchical':
            model = AgglomerativeClustering(n_clusters=self.n_regimes)
        elif method == 'dbscan':
            model = DBSCAN(eps=0.5, min_samples=5)
        
        # 拟合模型
        if method == 'gmm':
            regime_labels = model.fit_predict(features_scaled)
            regime_probs = model.predict_proba(features_scaled)
        else:
            regime_labels = model.fit_predict(features_scaled)
            regime_probs = None
        
        self.regime_models[method] = model
        self.regime_labels = pd.Series(regime_labels, index=features.index)
        
        # 分析各状态特征
        self.regime_characteristics = self._analyze_regime_characteristics(features, regime_labels)
        
        return regime_labels, regime_probs
    
    def _analyze_regime_characteristics(self, features: pd.DataFrame, 
                                      regime_labels: np.ndarray) -> pd.DataFrame:
        """
        分析各状态特征
        """
        
        regime_data = features.copy()
        regime_data['regime'] = regime_labels
        
        characteristics = regime_data.groupby('regime').agg({
            'mean_return': ['mean', 'std'],
            'return_volatility': ['mean', 'std'],
            'return_skewness': 'mean',
            'return_kurtosis': 'mean'
        })
        
        # 添加状态持续时间分析
        regime_series = pd.Series(regime_labels, index=features.index)
        duration_stats = []
        
        for regime in np.unique(regime_labels):
            regime_periods = (regime_series == regime).astype(int)
            
            # 计算连续期间
            consecutive_periods = []
            current_duration = 0
            
            for period in regime_periods:
                if period == 1:
                    current_duration += 1
                else:
                    if current_duration > 0:
                        consecutive_periods.append(current_duration)
                    current_duration = 0
            
            if current_duration > 0:
                consecutive_periods.append(current_duration)
            
            if consecutive_periods:
                duration_stats.append({
                    'regime': regime,
                    'avg_duration': np.mean(consecutive_periods),
                    'max_duration': np.max(consecutive_periods),
                    'frequency': len(consecutive_periods)
                })
        
        duration_df = pd.DataFrame(duration_stats)
        
        return characteristics, duration_df
    
    def visualize_regimes(self, features: pd.DataFrame, regime_labels: np.ndarray):
        """
        可视化市场状态
        """
        
        # PCA降维可视化
        pca = PCA(n_components=2)
        features_pca = pca.fit_transform(StandardScaler().fit_transform(features))
        
        fig = px.scatter(
            x=features_pca[:, 0], 
            y=features_pca[:, 1],
            color=regime_labels,
            title='Market Regimes - PCA Visualization',
            labels={'x': f'PC1 ({pca.explained_variance_ratio_[0]:.2%})',
                   'y': f'PC2 ({pca.explained_variance_ratio_[1]:.2%})'}
        )
        fig.show()
        
        # t-SNE可视化
        tsne = TSNE(n_components=2, random_state=42)
        features_tsne = tsne.fit_transform(StandardScaler().fit_transform(features.iloc[::10]))  # 抽样以提高速度
        
        fig = px.scatter(
            x=features_tsne[:, 0], 
            y=features_tsne[:, 1],
            color=regime_labels[::10],
            title='Market Regimes - t-SNE Visualization'
        )
        fig.show()
        
        # 时间序列可视化
        regime_series = pd.Series(regime_labels, index=features.index)
        
        fig = go.Figure()
        
        for regime in np.unique(regime_labels):
            regime_dates = regime_series[regime_series == regime].index
            y_values = [regime] * len(regime_dates)
            
            fig.add_trace(go.Scatter(
                x=regime_dates,
                y=y_values,
                mode='markers',
                name=f'Regime {regime}',
                marker=dict(size=3)
            ))
        
        fig.update_layout(
            title='Market Regimes Over Time',
            xaxis_title='Date',
            yaxis_title='Regime',
            height=400
        )
        fig.show()

def regime_based_strategy(regime_labels: pd.Series, returns: pd.Series,
                         lookback_window: int = 20) -> pd.Series:
    """
    基于状态的策略
    """
    
    strategy_returns = []
    
    # 计算每个状态下的历史表现
    regime_performance = {}
    for regime in regime_labels.unique():
        regime_mask = regime_labels == regime
        regime_returns = returns[regime_mask]
        regime_performance[regime] = {
            'mean_return': regime_returns.mean(),
            'volatility': regime_returns.std(),
            'sharpe': regime_returns.mean() / regime_returns.std() if regime_returns.std() > 0 else 0
        }
    
    # 根据当前状态调整仓位
    for i, (date, current_regime) in enumerate(regime_labels.items()):
        if i < lookback_window:
            strategy_returns.append(0)
            continue
        
        # 根据当前状态决定仓位
        regime_perf = regime_performance[current_regime]
        
        if regime_perf['mean_return'] > 0 and regime_perf['sharpe'] > 0.5:
            # 好的状态，做多
            position = 1.0
        elif regime_perf['mean_return'] < 0 or regime_perf['sharpe'] < -0.5:
            # 坏的状态，减仓或做空
            position = 0.0  # 这里简化为持现金
        else:
            # 中性状态，部分仓位
            position = 0.5
        
        # 计算策略收益
        if i < len(returns):
            strategy_return = position * returns.iloc[i]
            strategy_returns.append(strategy_return)
    
    return pd.Series(strategy_returns, index=regime_labels.index[:len(strategy_returns)])
```

### 3.2 异常检测：风险预警

```python
from sklearn.ensemble import IsolationForest
from sklearn.svm import OneClassSVM
from sklearn.covariance import EllipticEnvelope
from scipy import stats

class AnomalyDetectionSystem:
    """
    异常检测系统
    """
    
    def __init__(self, contamination: float = 0.1):
        self.contamination = contamination
        self.detectors = {}
        self.anomaly_scores = {}
        
    def prepare_anomaly_features(self, price_data: pd.DataFrame,
                                volume_data: pd.DataFrame = None) -> pd.DataFrame:
        """
        准备异常检测特征
        """
        
        returns = price_data.pct_change().dropna()
        features = pd.DataFrame(index=returns.index)
        
        # 收益率异常特征
        features['daily_return'] = returns.mean(axis=1) if len(returns.columns) > 1 else returns.iloc[:, 0]
        features['return_volatility'] = returns.std(axis=1)
        features['extreme_return'] = np.abs(features['daily_return'])
        
        # 滚动窗口异常特征
        windows = [5, 20, 60]
        for window in windows:
            # Z-score
            rolling_mean = features['daily_return'].rolling(window).mean()
            rolling_std = features['daily_return'].rolling(window).std()
            features[f'zscore_{window}d'] = (features['daily_return'] - rolling_mean) / rolling_std
            
            # 分位数位置
            features[f'quantile_position_{window}d'] = features['daily_return'].rolling(window).rank(pct=True)
        
        # 技术指标异常
        price_avg = price_data.mean(axis=1) if len(price_data.columns) > 1 else price_data.iloc[:, 0]
        
        # RSI极值
        delta = price_avg.diff()
        gain = (delta.where(delta > 0, 0)).rolling(14).mean()
        loss = (-delta.where(delta < 0, 0)).rolling(14).mean()
        rs = gain / loss
        rsi = 100 - (100 / (1 + rs))
        features['rsi'] = rsi
        features['rsi_extreme'] = ((rsi > 80) | (rsi < 20)).astype(int)
        
        # 布林带位置
        sma_20 = price_avg.rolling(20).mean()
        std_20 = price_avg.rolling(20).std()
        bollinger_position = (price_avg - sma_20) / (2 * std_20)
        features['bollinger_position'] = bollinger_position
        features['bollinger_extreme'] = (np.abs(bollinger_position) > 1).astype(int)
        
        # 成交量异常
        if volume_data is not None:
            volume_avg = volume_data.mean(axis=1) if len(volume_data.columns) > 1 else volume_data.iloc[:, 0]
            volume_ma = volume_avg.rolling(20).mean()
            features['volume_ratio'] = volume_avg / volume_ma
            features['volume_spike'] = (features['volume_ratio'] > 2).astype(int)
        
        # VIX特征（如果有波动率数据）
        implied_vol = features['return_volatility'].rolling(20).mean() * np.sqrt(252) * 100
        features['implied_vol'] = implied_vol
        features['vol_spike'] = (implied_vol > implied_vol.rolling(252).quantile(0.9)).astype(int)
        
        return features.dropna()
    
    def detect_anomalies(self, features: pd.DataFrame) -> Dict[str, pd.Series]:
        """
        使用多种方法检测异常
        """
        
        # 标准化特征
        scaler = RobustScaler()
        features_scaled = scaler.fit_transform(features)
        
        anomaly_results = {}
        
        # 1. Isolation Forest
        iso_forest = IsolationForest(contamination=self.contamination, random_state=42)
        iso_anomalies = iso_forest.fit_predict(features_scaled)
        iso_scores = iso_forest.decision_function(features_scaled)
        
        anomaly_results['isolation_forest'] = pd.Series(iso_anomalies, index=features.index)
        self.anomaly_scores['isolation_forest'] = pd.Series(iso_scores, index=features.index)
        
        # 2. One-Class SVM
        oc_svm = OneClassSVM(gamma='scale', nu=self.contamination)
        svm_anomalies = oc_svm.fit_predict(features_scaled)
        svm_scores = oc_svm.decision_function(features_scaled)
        
        anomaly_results['one_class_svm'] = pd.Series(svm_anomalies, index=features.index)
        self.anomaly_scores['one_class_svm'] = pd.Series(svm_scores, index=features.index)
        
        # 3. Elliptic Envelope
        elliptic = EllipticEnvelope(contamination=self.contamination, random_state=42)
        elliptic_anomalies = elliptic.fit_predict(features_scaled)
        elliptic_scores = elliptic.decision_function(features_scaled)
        
        anomaly_results['elliptic_envelope'] = pd.Series(elliptic_anomalies, index=features.index)
        self.anomaly_scores['elliptic_envelope'] = pd.Series(elliptic_scores, index=features.index)
        
        # 4. 统计方法
        statistical_anomalies = self._statistical_anomaly_detection(features)
        anomaly_results['statistical'] = statistical_anomalies
        
        # 5. 集成方法
        ensemble_anomalies = self._ensemble_anomaly_detection(anomaly_results)
        anomaly_results['ensemble'] = ensemble_anomalies
        
        self.detectors = {
            'isolation_forest': iso_forest,
            'one_class_svm': oc_svm,
            'elliptic_envelope': elliptic
        }
        
        return anomaly_results
    
    def _statistical_anomaly_detection(self, features: pd.DataFrame) -> pd.Series:
        """
        基于统计的异常检测
        """
        
        anomalies = []
        
        for _, row in features.iterrows():
            # 多变量异常检测
            # 使用马哈拉诺比斯距离
            try:
                mean_vec = features.mean()
                cov_matrix = features.cov()
                inv_cov = np.linalg.pinv(cov_matrix)
                
                diff = row - mean_vec
                mahal_dist = np.sqrt(diff.T @ inv_cov @ diff)
                
                # 使用卡方分布的临界值
                threshold = stats.chi2.ppf(1 - self.contamination, df=len(features.columns))
                is_anomaly = mahal_dist > threshold
                
                anomalies.append(-1 if is_anomaly else 1)
            except:
                anomalies.append(1)  # 默认为正常
        
        return pd.Series(anomalies, index=features.index)
    
    def _ensemble_anomaly_detection(self, anomaly_results: Dict[str, pd.Series]) -> pd.Series:
        """
        集成异常检测结果
        """
        
        # 排除集成方法本身
        methods = [key for key in anomaly_results.keys() if key != 'ensemble']
        
        # 投票法
        votes = pd.DataFrame({method: anomaly_results[method] for method in methods})
        anomaly_votes = (votes == -1).sum(axis=1)
        
        # 如果超过一半的方法认为是异常，则判定为异常
        threshold = len(methods) / 2
        ensemble_result = (anomaly_votes > threshold).astype(int)
        ensemble_result = ensemble_result.replace({0: 1, 1: -1})  # 转换为sklearn格式
        
        return ensemble_result
    
    def analyze_anomaly_patterns(self, anomaly_results: Dict[str, pd.Series],
                                features: pd.DataFrame, returns: pd.Series):
        """
        分析异常模式
        """
        
        analysis_results = {}
        
        for method, anomalies in anomaly_results.items():
            anomaly_dates = anomalies[anomalies == -1].index
            
            if len(anomaly_dates) > 0:
                # 异常期间的特征分析
                anomaly_features = features.loc[anomaly_dates]
                normal_features = features.loc[anomalies[anomalies == 1].index]
                
                feature_comparison = pd.DataFrame({
                    'anomaly_mean': anomaly_features.mean(),
                    'normal_mean': normal_features.mean(),
                    'anomaly_std': anomaly_features.std(),
                    'normal_std': normal_features.std()
                })
                
                feature_comparison['difference'] = (
                    feature_comparison['anomaly_mean'] - feature_comparison['normal_mean']
                ) / feature_comparison['normal_std']
                
                # 异常期间的收益分析
                anomaly_returns = returns.loc[anomaly_dates]
                normal_returns = returns.loc[anomalies[anomalies == 1].index]
                
                return_analysis = {
                    'anomaly_return_mean': anomaly_returns.mean(),
                    'normal_return_mean': normal_returns.mean(),
                    'anomaly_return_std': anomaly_returns.std(),
                    'normal_return_std': normal_returns.std(),
                    'anomaly_return_min': anomaly_returns.min(),
                    'anomaly_return_max': anomaly_returns.max()
                }
                
                analysis_results[method] = {
                    'feature_comparison': feature_comparison,
                    'return_analysis': return_analysis,
                    'anomaly_count': len(anomaly_dates),
                    'anomaly_rate': len(anomaly_dates) / len(anomalies)
                }
        
        return analysis_results
    
    def create_early_warning_system(self, features: pd.DataFrame, 
                                   lookback_window: int = 20) -> pd.DataFrame:
        """
        创建早期预警系统
        """
        
        warning_signals = pd.DataFrame(index=features.index)
        
        # 滚动异常检测
        for i in range(lookback_window, len(features)):
            # 使用历史数据训练
            train_data = features.iloc[i-lookback_window:i]
            current_data = features.iloc[[i]]
            
            # 简化的早期预警
            scaler = RobustScaler()
            train_scaled = scaler.fit_transform(train_data)
            current_scaled = scaler.transform(current_data)
            
            # 使用Isolation Forest
            iso_forest = IsolationForest(contamination=0.1, random_state=42)
            iso_forest.fit(train_scaled)
            
            anomaly_score = iso_forest.decision_function(current_scaled)[0]
            is_anomaly = iso_forest.predict(current_scaled)[0] == -1
            
            warning_signals.loc[features.index[i], 'anomaly_score'] = anomaly_score
            warning_signals.loc[features.index[i], 'is_anomaly'] = is_anomaly
            
            # 风险等级
            if anomaly_score < -0.5:
                risk_level = 'High'
            elif anomaly_score < -0.2:
                risk_level = 'Medium'
            else:
                risk_level = 'Low'
            
            warning_signals.loc[features.index[i], 'risk_level'] = risk_level
        
        return warning_signals.dropna()
```

## 📊 实战案例：机器学习量化策略系统

```python
def comprehensive_ml_trading_system():
    """
    综合机器学习交易系统案例
    """
    
    print("=== 机器学习量化策略系统 ===")
    
    # 1. 数据准备
    np.random.seed(42)
    dates = pd.date_range('2018-01-01', '2023-12-31', freq='D')
    n_assets = 5
    
    # 模拟多资产价格数据
    price_data = pd.DataFrame()
    for i in range(n_assets):
        # 添加一些趋势和周期性
        trend = np.linspace(0, 0.5, len(dates))
        seasonality = 0.1 * np.sin(2 * np.pi * np.arange(len(dates)) / 252)
        noise = np.random.normal(0, 0.02, len(dates))
        
        price_series = 100 * np.exp(np.cumsum(0.0008 + trend/len(dates) + seasonality + noise))
        price_data[f'Asset_{i+1}'] = price_series
    
    price_data.index = dates
    
    # 模拟成交量数据
    volume_data = pd.DataFrame(
        np.random.lognormal(10, 0.5, (len(dates), n_assets)),
        index=dates,
        columns=price_data.columns
    )
    
    print(f"数据期间: {dates[0].date()} 到 {dates[-1].date()}")
    print(f"资产数量: {n_assets}")
    
    # 2. 特征工程
    print("\n=== 特征工程 ===")
    
    ml_framework = FinancialMLFramework()
    features = ml_framework.prepare_financial_data(price_data, volume_data)
    
    print(f"生成特征数量: {len(features.columns)}")
    print(f"特征样本数量: {len(features)}")
    
    # 3. 创建标签
    print("\n=== 创建标签 ===")
    
    # 创建多种标签
    labels_dict = {}
    
    # 收益方向预测（第一个资产）
    labels_dict['direction'] = ml_framework.create_labels(
        price_data[['Asset_1']], method='return_direction', horizon=1
    )
    
    # 收益分位数预测
    labels_dict['quantile'] = ml_framework.create_labels(
        price_data[['Asset_1']], method='return_quantile', horizon=5
    )
    
    # 波动率状态预测
    labels_dict['volatility'] = ml_framework.create_labels(
        price_data[['Asset_1']], method='volatility_regime', horizon=5
    )
    
    # 4. 数据分割和预处理
    print("\n=== 数据分割 ===")
    
    # 对齐特征和标签
    common_dates = features.index.intersection(labels_dict['direction'].index)
    features_aligned = features.loc[common_dates]
    
    # 处理数据问题
    features_processed = ml_framework.handle_financial_data_issues(features_aligned)
    
    # 时间序列分割
    ts_split = FinancialTimeSeriesSplit(n_splits=5, test_size=252, gap=5)
    
    # 5. 收益方向预测
    print("\n=== 收益方向预测 ===")
    
    direction_predictor = ReturnDirectionPredictor()
    
    # 使用最后一次分割进行演示
    splits = list(ts_split.split(features_processed))
    train_idx, test_idx = splits[-1]
    
    X_train = features_processed.iloc[train_idx]
    X_test = features_processed.iloc[test_idx]
    y_train = labels_dict['direction'].iloc[train_idx]
    y_test = labels_dict['direction'].iloc[test_idx]
    
    # 训练模型
    direction_predictor.train_models(X_train, y_train)
    
    # 评估模型
    direction_results = direction_predictor.evaluate_models(X_test, y_test)
    
    print("\n方向预测结果:")
    for model_name, metrics in direction_results.items():
        print(f"{model_name}:")
        print(f"  准确率: {metrics['accuracy']:.3f}")
        print(f"  精确率: {metrics['precision']:.3f}")
        print(f"  召回率: {metrics['recall']:.3f}")
        print(f"  F1分数: {metrics['f1_score']:.3f}")
    
    # 6. 收益预测
    print("\n=== 收益预测 ===")
    
    # 创建连续收益标签
    returns = price_data['Asset_1'].pct_change().shift(-1)  # 下一期收益
    returns_aligned = returns.loc[common_dates].dropna()
    
    # 重新分割数据
    common_dates_returns = features_processed.index.intersection(returns_aligned.index)
    X_train_reg = features_processed.loc[common_dates_returns].iloc[train_idx]
    X_test_reg = features_processed.loc[common_dates_returns].iloc[test_idx]
    y_train_reg = returns_aligned.loc[common_dates_returns].iloc[train_idx]
    y_test_reg = returns_aligned.loc[common_dates_returns].iloc[test_idx]
    
    return_predictor = ReturnPredictor()
    return_predictor.train_models(X_train_reg, y_train_reg)
    
    regression_results = return_predictor.evaluate_models(X_test_reg, y_test_reg)
    
    print("\n收益预测结果:")
    for model_name, metrics in regression_results.items():
        print(f"{model_name}:")
        print(f"  R²: {metrics['r2']:.3f}")
        print(f"  RMSE: {metrics['rmse']:.4f}")
        print(f"  方向准确率: {metrics['direction_accuracy']:.3f}")
        print(f"  信息系数: {metrics['information_coefficient']:.3f}")
    
    # 7. 市场状态分析
    print("\n=== 市场状态分析 ===")
    
    regime_analyzer = MarketRegimeAnalysis(n_regimes=3)
    regime_features = regime_analyzer.prepare_regime_features(price_data, volume_data)
    
    regime_labels, _ = regime_analyzer.identify_regimes(regime_features, method='kmeans')
    
    # 分析状态特征
    characteristics = regime_analyzer.regime_characteristics
    print(f"\n识别出 {len(np.unique(regime_labels))} 个市场状态")
    
    # 8. 异常检测
    print("\n=== 异常检测 ===")
    
    anomaly_detector = AnomalyDetectionSystem(contamination=0.05)
    anomaly_features = anomaly_detector.prepare_anomaly_features(price_data, volume_data)
    
    anomaly_results = anomaly_detector.detect_anomalies(anomaly_features)
    
    print("\n异常检测结果:")
    for method, anomalies in anomaly_results.items():
        anomaly_count = (anomalies == -1).sum()
        anomaly_rate = anomaly_count / len(anomalies)
        print(f"{method}: {anomaly_count} 个异常 ({anomaly_rate:.2%})")
    
    # 9. 策略集成和回测
    print("\n=== 策略集成 ===")
    
    # 选择最佳模型
    best_direction_model = max(direction_results.keys(), 
                              key=lambda x: direction_results[x]['f1_score'])
    
    best_return_model = max(regression_results.keys(),
                           key=lambda x: regression_results[x]['information_coefficient'])
    
    print(f"最佳方向预测模型: {best_direction_model}")
    print(f"最佳收益预测模型: {best_return_model}")
    
    # 构建集成策略
    direction_signals = direction_results[best_direction_model]['predictions']
    return_predictions = regression_results[best_return_model]['predictions']
    
    # 简单集成策略：方向和大小结合
    strategy_signals = []
    for i in range(len(direction_signals)):
        if direction_signals[i] == 1 and return_predictions[i] > 0:
            signal = min(abs(return_predictions[i]) * 10, 1.0)  # 根据预测强度调整仓位
        elif direction_signals[i] == 0 and return_predictions[i] < 0:
            signal = 0.0  # 预测下跌，持现金
        else:
            signal = 0.5  # 信号不一致，中性仓位
        
        strategy_signals.append(signal)
    
    # 计算策略收益
    test_returns = returns_aligned.iloc[test_idx]
    strategy_returns = []
    
    for i, (signal, actual_return) in enumerate(zip(strategy_signals, test_returns)):
        strategy_return = signal * actual_return
        strategy_returns.append(strategy_return)
    
    strategy_returns = pd.Series(strategy_returns, index=test_returns.index)
    
    # 绩效分析
    total_return = (1 + strategy_returns).prod() - 1
    annual_return = (1 + strategy_returns.mean()) ** 252 - 1
    annual_vol = strategy_returns.std() * np.sqrt(252)
    sharpe_ratio = annual_return / annual_vol if annual_vol > 0 else 0
    
    # 基准比较
    benchmark_returns = test_returns  # 买入持有策略
    benchmark_total = (1 + benchmark_returns).prod() - 1
    benchmark_annual = (1 + benchmark_returns.mean()) ** 252 - 1
    benchmark_vol = benchmark_returns.std() * np.sqrt(252)
    benchmark_sharpe = benchmark_annual / benchmark_vol if benchmark_vol > 0 else 0
    
    print(f"\n=== 策略绩效 ===")
    print(f"策略总收益: {total_return:.2%}")
    print(f"策略年化收益: {annual_return:.2%}")
    print(f"策略年化波动: {annual_vol:.2%}")
    print(f"策略夏普比率: {sharpe_ratio:.3f}")
    print(f"\n基准总收益: {benchmark_total:.2%}")
    print(f"基准年化收益: {benchmark_annual:.2%}")
    print(f"基准年化波动: {benchmark_vol:.2%}")
    print(f"基准夏普比率: {benchmark_sharpe:.3f}")
    print(f"\n超额收益: {annual_return - benchmark_annual:.2%}")
    
    return {
        'features': features_processed,
        'direction_results': direction_results,
        'regression_results': regression_results,
        'regime_labels': regime_labels,
        'anomaly_results': anomaly_results,
        'strategy_returns': strategy_returns,
        'benchmark_returns': benchmark_returns
    }

# 运行综合案例
if __name__ == "__main__":
    ml_results = comprehensive_ml_trading_system()
```

## 📚 本章小结

本章详细介绍了机器学习在量化投资中的应用：

### 🔑 核心要点

1. **金融ML特殊性**: 时间序列特性、非平稳性、噪声大、样本不平衡
2. **监督学习应用**: 收益方向预测、收益大小预测、多空策略构建
3. **无监督学习应用**: 市场状态识别、异常检测、风险预警
4. **特征工程**: 技术指标、基本面指标、宏观指标、时间特征

### 📈 实践要点

- 使用时间序列交叉验证避免数据泄露
- 特征工程是成功的关键，需要金融领域知识
- 模型集成通常比单一模型表现更好
- 异常检测对风险管理具有重要意义

### 🎯 下章预告

下一章将探讨**深度学习和自然语言处理在金融中的应用**，包括神经网络、LSTM、Transformer等先进技术在量化投资中的具体应用。

---

**风险提示**: 机器学习模型存在过拟合风险，历史表现不代表未来收益，需要持续监控和调整模型。