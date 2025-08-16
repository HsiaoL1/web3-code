# 第三章：Python 金融编程（上篇）

## 📖 章节概述

本章将深入学习 Python 在金融领域的应用，包括金融数据的获取、处理、分析和可视化。这是量化交易的技术基础，将为后续的策略开发提供必要的编程技能。

## 🎯 学习目标

- 掌握 Python 金融库的使用方法
- 学会从各种数据源获取金融数据
- 熟练进行金融数据的清洗和处理
- 掌握金融数据的可视化技巧
- 了解基础的量化策略实现方法
- 建立完整的金融数据分析工作流

## 3.1 Python 金融生态系统

### 3.1.1 核心金融库介绍

**数据处理库**

- **pandas**：数据分析和操作的核心库
- **numpy**：数值计算基础
- **scipy**：科学计算扩展

**金融专业库**

- **yfinance**：Yahoo Finance 数据获取
- **pandas-datareader**：多数据源接口
- **tushare**：中国 A 股数据（免费版）
- **akshare**：中文财经数据接口

**技术分析库**

- **ta-lib**：技术分析指标库
- **pandas-ta**：基于 pandas 的技术分析
- **stockstats**：股票统计指标

**机器学习库**

- **scikit-learn**：机器学习算法
- **tensorflow/pytorch**：深度学习框架

**可视化库**

- **matplotlib**：基础绘图
- **seaborn**：统计可视化
- **plotly**：交互式图表
- **mplfinance**：专业金融图表

```python
# 安装必要的库
# pip install pandas numpy matplotlib seaborn yfinance pandas-datareader
# pip install ta-lib plotly akshare tushare scikit-learn

import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
import seaborn as sns
import warnings
warnings.filterwarnings('ignore')

# 设置中文字体和图表样式
plt.rcParams['font.sans-serif'] = ['SimHei']
plt.rcParams['axes.unicode_minus'] = False
sns.set_style("whitegrid")

print("Python金融分析环境配置完成！")
```

### 3.1.2 开发环境配置

**推荐开发环境：**

1. **Jupyter Notebook/Lab**：交互式开发
2. **PyCharm**：专业 IDE
3. **VS Code**：轻量级编辑器
4. **Spyder**：科学计算 IDE

**环境管理：**

```bash
# 创建虚拟环境
conda create -n quant_trading python=3.9
conda activate quant_trading

# 安装基础包
conda install pandas numpy matplotlib seaborn jupyter
pip install yfinance akshare tushare plotly ta-lib
```

## 3.2 金融数据获取

### 3.2.1 使用 yfinance 获取国外市场数据

```python
import yfinance as yf
from datetime import datetime, timedelta

def get_stock_data(symbol, period="1y", interval="1d"):
    """
    获取股票数据

    Parameters:
    symbol: 股票代码
    period: 时间周期 (1d, 5d, 1mo, 3mo, 6mo, 1y, 2y, 5y, 10y, ytd, max)
    interval: 数据间隔 (1m, 2m, 5m, 15m, 30m, 60m, 90m, 1h, 1d, 5d, 1wk, 1mo, 3mo)
    """
    try:
        ticker = yf.Ticker(symbol)
        data = ticker.history(period=period, interval=interval)

        # 添加股票代码列
        data['Symbol'] = symbol

        return data
    except Exception as e:
        print(f"获取{symbol}数据时出错: {e}")
        return None

# 获取单只股票数据
apple_data = get_stock_data("AAPL", "2y", "1d")
print("苹果股票数据:")
print(apple_data.head())
print(f"数据形状: {apple_data.shape}")

# 获取多只股票数据
symbols = ["AAPL", "GOOGL", "MSFT", "TSLA", "AMZN"]
stock_data = {}

for symbol in symbols:
    data = get_stock_data(symbol, "1y", "1d")
    if data is not None:
        stock_data[symbol] = data
        print(f"成功获取{symbol}数据")

# 合并多只股票的收盘价
close_prices = pd.DataFrame()
for symbol, data in stock_data.items():
    close_prices[symbol] = data['Close']

print("\n多只股票收盘价数据:")
print(close_prices.head())
```

### 3.2.2 使用 akshare 获取中国市场数据

```python
import akshare as ak

def get_china_stock_data(symbol, period="daily", adjust="qfq"):
    """
    获取中国股票数据

    Parameters:
    symbol: 股票代码（如"000001"）
    period: 周期（daily, weekly, monthly）
    adjust: 复权类型（qfq前复权, hfq后复权, ""不复权）
    """
    try:
        # 获取股票历史数据
        if period == "daily":
            data = ak.stock_zh_a_hist(symbol=symbol, adjust=adjust)
        else:
            data = ak.stock_zh_a_hist(symbol=symbol, period=period, adjust=adjust)

        # 重命名列以保持一致性
        data.columns = ['Date', 'Open', 'Close', 'High', 'Low',
                       'Volume', 'Amount', 'Amplitude', 'Pct_change',
                       'Change', 'Turnover']

        # 设置日期索引
        data['Date'] = pd.to_datetime(data['Date'])
        data.set_index('Date', inplace=True)

        return data
    except Exception as e:
        print(f"获取{symbol}数据时出错: {e}")
        return None

# 获取平安银行数据
ping_an_data = get_china_stock_data("000001")
if ping_an_data is not None:
    print("平安银行股票数据:")
    print(ping_an_data.head())

# 获取股票基本信息
def get_stock_info():
    """获取A股股票基本信息"""
    try:
        # 获取股票列表
        stock_list = ak.stock_info_a_code_name()
        return stock_list
    except Exception as e:
        print(f"获取股票列表出错: {e}")
        return None

stock_list = get_stock_info()
if stock_list is not None:
    print(f"\nA股股票总数: {len(stock_list)}")
    print("前10只股票:")
    print(stock_list.head(10))
```

### 3.2.3 获取宏观经济数据

```python
def get_macro_data():
    """获取宏观经济数据"""
    try:
        # 获取GDP数据
        gdp_data = ak.macro_china_gdp()

        # 获取CPI数据
        cpi_data = ak.macro_china_cpi()

        # 获取PMI数据
        pmi_data = ak.macro_china_pmi()

        return gdp_data, cpi_data, pmi_data
    except Exception as e:
        print(f"获取宏观数据出错: {e}")
        return None, None, None

# 获取指数数据
def get_index_data(symbol="sh000001", period="daily"):
    """
    获取指数数据
    symbol: 指数代码（sh000001上证指数, sz399001深证成指等）
    """
    try:
        data = ak.stock_zh_index_daily(symbol=symbol)
        data['date'] = pd.to_datetime(data['date'])
        data.set_index('date', inplace=True)
        return data
    except Exception as e:
        print(f"获取指数数据出错: {e}")
        return None

# 获取上证指数数据
sh_index = get_index_data("sh000001")
if sh_index is not None:
    print("上证指数数据:")
    print(sh_index.tail())
```

### 3.2.4 实时数据获取

```python
def get_realtime_data():
    """获取实时行情数据"""
    try:
        # 获取A股实时行情
        realtime_data = ak.stock_zh_a_spot_em()

        # 选择主要列
        columns = ['代码', '名称', '最新价', '涨跌幅', '涨跌额',
                  '成交量', '成交额', '振幅', '最高', '最低', '今开', '昨收']

        realtime_selected = realtime_data[columns]

        return realtime_selected
    except Exception as e:
        print(f"获取实时数据出错: {e}")
        return None

# 获取实时数据示例
realtime_stocks = get_realtime_data()
if realtime_stocks is not None:
    print("A股实时行情（前10只）:")
    print(realtime_stocks.head(10))
```

## 3.3 数据清洗与处理

### 3.3.1 数据质量检查

```python
def data_quality_check(data, name="数据"):
    """
    数据质量检查
    """
    print(f"\n=== {name}质量检查 ===")

    # 基本信息
    print(f"数据形状: {data.shape}")
    print(f"数据类型:\n{data.dtypes}")

    # 缺失值检查
    missing_count = data.isnull().sum()
    missing_pct = (missing_count / len(data)) * 100

    print(f"\n缺失值统计:")
    for col in data.columns:
        if missing_count[col] > 0:
            print(f"{col}: {missing_count[col]} ({missing_pct[col]:.2f}%)")

    # 重复值检查
    duplicates = data.duplicated().sum()
    print(f"\n重复行数: {duplicates}")

    # 数值型列的异常值检查
    numeric_cols = data.select_dtypes(include=[np.number]).columns
    print(f"\n数值型列统计:")
    print(data[numeric_cols].describe())

    # 检查零值和负值（如果适用）
    for col in numeric_cols:
        zero_count = (data[col] == 0).sum()
        negative_count = (data[col] < 0).sum()
        if zero_count > 0:
            print(f"{col}列零值数量: {zero_count}")
        if negative_count > 0:
            print(f"{col}列负值数量: {negative_count}")

# 对获取的数据进行质量检查
if apple_data is not None:
    data_quality_check(apple_data, "苹果股票数据")
```

### 3.3.2 数据清洗方法

```python
def clean_stock_data(data):
    """
    股票数据清洗
    """
    cleaned_data = data.copy()

    # 1. 删除重复行
    cleaned_data = cleaned_data.drop_duplicates()

    # 2. 处理缺失值
    # 对于价格数据，使用前向填充
    price_columns = ['Open', 'High', 'Low', 'Close']
    for col in price_columns:
        if col in cleaned_data.columns:
            cleaned_data[col] = cleaned_data[col].fillna(method='ffill')

    # 对于成交量，使用0填充
    if 'Volume' in cleaned_data.columns:
        cleaned_data['Volume'] = cleaned_data['Volume'].fillna(0)

    # 3. 异常值处理
    # 检查价格的逻辑关系：High >= max(Open, Close), Low <= min(Open, Close)
    if all(col in cleaned_data.columns for col in price_columns):
        # 修正异常的高低价
        cleaned_data['High'] = np.maximum(cleaned_data['High'],
                                        np.maximum(cleaned_data['Open'], cleaned_data['Close']))
        cleaned_data['Low'] = np.minimum(cleaned_data['Low'],
                                       np.minimum(cleaned_data['Open'], cleaned_data['Close']))

    # 4. 计算收益率
    if 'Close' in cleaned_data.columns:
        cleaned_data['Returns'] = cleaned_data['Close'].pct_change()
        cleaned_data['Log_Returns'] = np.log(cleaned_data['Close'] / cleaned_data['Close'].shift(1))

    # 5. 移除极端异常值（收益率超过3个标准差）
    if 'Returns' in cleaned_data.columns:
        returns_std = cleaned_data['Returns'].std()
        returns_mean = cleaned_data['Returns'].mean()

        # 标记异常值但不删除，而是用中位数替换
        outlier_mask = np.abs(cleaned_data['Returns'] - returns_mean) > 3 * returns_std
        outlier_count = outlier_mask.sum()

        if outlier_count > 0:
            print(f"发现{outlier_count}个异常收益率值，将进行处理")
            cleaned_data.loc[outlier_mask, 'Returns'] = cleaned_data['Returns'].median()

    print(f"数据清洗完成，原始数据: {len(data)}行，清洗后: {len(cleaned_data)}行")

    return cleaned_data

# 清洗苹果股票数据
if apple_data is not None:
    apple_clean = clean_stock_data(apple_data)
    print("\n清洗后的数据:")
    print(apple_clean.head())
```

### 3.3.3 数据标准化和归一化

```python
from sklearn.preprocessing import StandardScaler, MinMaxScaler

def normalize_data(data, method='standard'):
    """
    数据标准化/归一化

    method: 'standard' (标准化) 或 'minmax' (归一化)
    """
    numeric_columns = data.select_dtypes(include=[np.number]).columns

    if method == 'standard':
        scaler = StandardScaler()
        normalized_data = data.copy()
        normalized_data[numeric_columns] = scaler.fit_transform(data[numeric_columns])
        print("使用标准化处理数据")
    elif method == 'minmax':
        scaler = MinMaxScaler()
        normalized_data = data.copy()
        normalized_data[numeric_columns] = scaler.fit_transform(data[numeric_columns])
        print("使用归一化处理数据")
    else:
        print("不支持的标准化方法")
        return data, None

    return normalized_data, scaler

# 标准化处理示例
if apple_clean is not None:
    # 选择需要标准化的列
    features_to_normalize = ['Open', 'High', 'Low', 'Close', 'Volume']
    data_subset = apple_clean[features_to_normalize].dropna()

    normalized_data, scaler = normalize_data(data_subset, 'standard')

    print("标准化前后对比:")
    print("原始数据统计:")
    print(data_subset.describe())
    print("\n标准化后数据统计:")
    print(normalized_data.describe())
```

## 3.4 技术指标计算

### 3.4.1 趋势指标

```python
def calculate_moving_averages(data, short_window=20, long_window=50):
    """
    计算移动平均线
    """
    result = data.copy()

    # 简单移动平均
    result[f'SMA_{short_window}'] = data['Close'].rolling(window=short_window).mean()
    result[f'SMA_{long_window}'] = data['Close'].rolling(window=long_window).mean()

    # 指数移动平均
    result[f'EMA_{short_window}'] = data['Close'].ewm(span=short_window).mean()
    result[f'EMA_{long_window}'] = data['Close'].ewm(span=long_window).mean()

    # 移动平均收敛发散(MACD)
    ema_12 = data['Close'].ewm(span=12).mean()
    ema_26 = data['Close'].ewm(span=26).mean()
    result['MACD'] = ema_12 - ema_26
    result['MACD_Signal'] = result['MACD'].ewm(span=9).mean()
    result['MACD_Histogram'] = result['MACD'] - result['MACD_Signal']

    return result

def calculate_bollinger_bands(data, window=20, num_std=2):
    """
    计算布林带
    """
    result = data.copy()

    # 中轨（移动平均线）
    result['BB_Middle'] = data['Close'].rolling(window=window).mean()

    # 标准差
    rolling_std = data['Close'].rolling(window=window).std()

    # 上轨和下轨
    result['BB_Upper'] = result['BB_Middle'] + (rolling_std * num_std)
    result['BB_Lower'] = result['BB_Middle'] - (rolling_std * num_std)

    # 布林带宽度和位置
    result['BB_Width'] = (result['BB_Upper'] - result['BB_Lower']) / result['BB_Middle']
    result['BB_Position'] = (data['Close'] - result['BB_Lower']) / (result['BB_Upper'] - result['BB_Lower'])

    return result

# 计算技术指标
if apple_clean is not None:
    apple_with_indicators = calculate_moving_averages(apple_clean)
    apple_with_indicators = calculate_bollinger_bands(apple_with_indicators)

    print("添加技术指标后的数据:")
    print(apple_with_indicators[['Close', 'SMA_20', 'SMA_50', 'EMA_20', 'EMA_50',
                                'MACD', 'BB_Upper', 'BB_Lower']].tail())
```

### 3.4.2 动量指标

```python
def calculate_momentum_indicators(data, window=14):
    """
    计算动量指标
    """
    result = data.copy()

    # RSI (相对强弱指数)
    delta = data['Close'].diff()
    gain = (delta.where(delta > 0, 0)).rolling(window=window).mean()
    loss = (-delta.where(delta < 0, 0)).rolling(window=window).mean()
    rs = gain / loss
    result['RSI'] = 100 - (100 / (1 + rs))

    # 随机振荡器 (Stochastic Oscillator)
    lowest_low = data['Low'].rolling(window=window).min()
    highest_high = data['High'].rolling(window=window).max()
    result['%K'] = ((data['Close'] - lowest_low) / (highest_high - lowest_low)) * 100
    result['%D'] = result['%K'].rolling(window=3).mean()

    # 威廉指标 (%R)
    result['Williams_%R'] = ((highest_high - data['Close']) / (highest_high - lowest_low)) * -100

    # 动量指标 (Momentum)
    result['Momentum'] = data['Close'] / data['Close'].shift(window) * 100

    # 变化率 (Rate of Change)
    result['ROC'] = ((data['Close'] - data['Close'].shift(window)) / data['Close'].shift(window)) * 100

    return result

def calculate_volume_indicators(data, window=20):
    """
    计算成交量指标
    """
    result = data.copy()

    # 成交量移动平均
    result['Volume_MA'] = data['Volume'].rolling(window=window).mean()

    # 成交量比率
    result['Volume_Ratio'] = data['Volume'] / result['Volume_MA']

    # 价量确认指标 (OBV - On Balance Volume)
    price_change = data['Close'].diff()
    obv = np.where(price_change > 0, data['Volume'],
                   np.where(price_change < 0, -data['Volume'], 0))
    result['OBV'] = pd.Series(obv, index=data.index).cumsum()

    # 资金流量指标 (MFI - Money Flow Index)
    typical_price = (data['High'] + data['Low'] + data['Close']) / 3
    raw_money_flow = typical_price * data['Volume']

    positive_flow = raw_money_flow.where(typical_price.diff() > 0, 0).rolling(window=window).sum()
    negative_flow = raw_money_flow.where(typical_price.diff() < 0, 0).rolling(window=window).sum()

    money_ratio = positive_flow / negative_flow
    result['MFI'] = 100 - (100 / (1 + money_ratio))

    return result

# 计算动量和成交量指标
if apple_with_indicators is not None:
    apple_full_indicators = calculate_momentum_indicators(apple_with_indicators)
    apple_full_indicators = calculate_volume_indicators(apple_full_indicators)

    print("所有技术指标:")
    indicator_columns = ['Close', 'RSI', '%K', '%D', 'Williams_%R',
                        'Momentum', 'ROC', 'Volume_Ratio', 'MFI']
    print(apple_full_indicators[indicator_columns].tail())
```

### 3.4.3 波动率指标

```python
def calculate_volatility_indicators(data, window=20):
    """
    计算波动率指标
    """
    result = data.copy()

    # 真实波幅 (True Range)
    high_low = data['High'] - data['Low']
    high_close_prev = np.abs(data['High'] - data['Close'].shift(1))
    low_close_prev = np.abs(data['Low'] - data['Close'].shift(1))

    true_range = np.maximum(high_low, np.maximum(high_close_prev, low_close_prev))
    result['TR'] = true_range

    # 平均真实波幅 (ATR)
    result['ATR'] = result['TR'].rolling(window=window).mean()

    # 历史波动率
    returns = np.log(data['Close'] / data['Close'].shift(1))
    result['Historical_Volatility'] = returns.rolling(window=window).std() * np.sqrt(252)

    # 最高价最低价波动率
    hl_volatility = np.log(data['High'] / data['Low'])
    result['HL_Volatility'] = hl_volatility.rolling(window=window).mean()

    # 帕金森波动率估计器
    result['Parkinson_Volatility'] = np.sqrt(
        hl_volatility.rolling(window=window).mean() * 252 / (4 * np.log(2))
    )

    return result

# 计算波动率指标
if apple_full_indicators is not None:
    apple_complete = calculate_volatility_indicators(apple_full_indicators)

    print("波动率指标:")
    volatility_columns = ['Close', 'ATR', 'Historical_Volatility',
                         'HL_Volatility', 'Parkinson_Volatility']
    print(apple_complete[volatility_columns].tail())
```

## 3.5 金融数据可视化

### 3.5.1 基础价格图表

```python
def plot_price_chart(data, title="股票价格图表", figsize=(12, 8)):
    """
    绘制价格图表
    """
    fig, axes = plt.subplots(3, 1, figsize=figsize, sharex=True)

    # 价格图
    axes[0].plot(data.index, data['Close'], label='收盘价', linewidth=1.5)
    if 'SMA_20' in data.columns:
        axes[0].plot(data.index, data['SMA_20'], label='SMA(20)', alpha=0.7)
    if 'SMA_50' in data.columns:
        axes[0].plot(data.index, data['SMA_50'], label='SMA(50)', alpha=0.7)

    axes[0].set_title(title)
    axes[0].set_ylabel('价格')
    axes[0].legend()
    axes[0].grid(True, alpha=0.3)

    # 成交量图
    axes[1].bar(data.index, data['Volume'], alpha=0.6, color='gray')
    axes[1].set_ylabel('成交量')
    axes[1].grid(True, alpha=0.3)

    # MACD图
    if all(col in data.columns for col in ['MACD', 'MACD_Signal', 'MACD_Histogram']):
        axes[2].plot(data.index, data['MACD'], label='MACD', linewidth=1.5)
        axes[2].plot(data.index, data['MACD_Signal'], label='Signal', linewidth=1.5)
        axes[2].bar(data.index, data['MACD_Histogram'], label='Histogram', alpha=0.6)
        axes[2].axhline(y=0, color='black', linestyle='-', alpha=0.5)
        axes[2].set_ylabel('MACD')
        axes[2].legend()
        axes[2].grid(True, alpha=0.3)

    axes[2].set_xlabel('日期')
    plt.tight_layout()
    plt.show()

# 绘制苹果股票图表
if apple_complete is not None:
    plot_price_chart(apple_complete.tail(252), "苹果股票(AAPL) - 近1年")
```

### 3.5.2 K 线图（蜡烛图）

```python
import mplfinance as mpf

def plot_candlestick_chart(data, title="K线图", volume=True, mav=(20, 50), figsize=(12, 8)):
    """
    绘制K线图
    """
    # 准备数据，mplfinance需要特定的列名
    plot_data = data[['Open', 'High', 'Low', 'Close', 'Volume']].copy()

    # 设置样式
    mc = mpf.make_marketcolors(up='red', down='green',    # 中国习惯：红涨绿跌
                              edge='inherit',
                              volume='inherit')

    style = mpf.make_mpf_style(marketcolors=mc,
                              gridstyle='-',
                              gridcolor='lightgray')

    # 添加附加绘图
    addplot_list = []

    # 添加移动平均线
    if mav:
        for ma_period in mav:
            ma_col = f'SMA_{ma_period}'
            if ma_col in data.columns:
                addplot_list.append(
                    mpf.make_addplot(data[ma_col], color=f'C{len(addplot_list)}', width=1.5)
                )

    # 添加布林带
    if all(col in data.columns for col in ['BB_Upper', 'BB_Lower']):
        addplot_list.extend([
            mpf.make_addplot(data['BB_Upper'], color='blue', alpha=0.5, linestyle='--'),
            mpf.make_addplot(data['BB_Lower'], color='blue', alpha=0.5, linestyle='--')
        ])

    # 绘制图表
    mpf.plot(plot_data,
            type='candle',
            style=style,
            title=title,
            volume=volume,
            addplot=addplot_list if addplot_list else None,
            figsize=figsize,
            tight_layout=True,
            show_nontrading=False)

# 绘制K线图
if apple_complete is not None:
    # 获取最近3个月的数据
    recent_data = apple_complete.tail(60)
    plot_candlestick_chart(recent_data, "苹果股票(AAPL) K线图 - 近3个月")
```

### 3.5.3 技术指标可视化

```python
def plot_technical_indicators(data, figsize=(15, 12)):
    """
    绘制技术指标图表
    """
    fig, axes = plt.subplots(4, 2, figsize=figsize)

    # RSI
    if 'RSI' in data.columns:
        axes[0, 0].plot(data.index, data['RSI'])
        axes[0, 0].axhline(y=70, color='r', linestyle='--', alpha=0.7, label='超买线')
        axes[0, 0].axhline(y=30, color='g', linestyle='--', alpha=0.7, label='超卖线')
        axes[0, 0].axhline(y=50, color='black', linestyle='-', alpha=0.5)
        axes[0, 0].set_title('RSI (相对强弱指数)')
        axes[0, 0].set_ylim(0, 100)
        axes[0, 0].legend()
        axes[0, 0].grid(True, alpha=0.3)

    # 随机振荡器
    if all(col in data.columns for col in ['%K', '%D']):
        axes[0, 1].plot(data.index, data['%K'], label='%K')
        axes[0, 1].plot(data.index, data['%D'], label='%D')
        axes[0, 1].axhline(y=80, color='r', linestyle='--', alpha=0.7)
        axes[0, 1].axhline(y=20, color='g', linestyle='--', alpha=0.7)
        axes[0, 1].set_title('随机振荡器 (%K %D)')
        axes[0, 1].set_ylim(0, 100)
        axes[0, 1].legend()
        axes[0, 1].grid(True, alpha=0.3)

    # 布林带位置
    if all(col in data.columns for col in ['Close', 'BB_Upper', 'BB_Lower', 'BB_Position']):
        axes[1, 0].plot(data.index, data['Close'], label='收盘价')
        axes[1, 0].plot(data.index, data['BB_Upper'], label='上轨', alpha=0.7)
        axes[1, 0].plot(data.index, data['BB_Lower'], label='下轨', alpha=0.7)
        axes[1, 0].fill_between(data.index, data['BB_Upper'], data['BB_Lower'], alpha=0.1)
        axes[1, 0].set_title('布林带')
        axes[1, 0].legend()
        axes[1, 0].grid(True, alpha=0.3)

    # ATR
    if 'ATR' in data.columns:
        axes[1, 1].plot(data.index, data['ATR'])
        axes[1, 1].set_title('ATR (平均真实波幅)')
        axes[1, 1].grid(True, alpha=0.3)

    # 成交量指标
    if all(col in data.columns for col in ['Volume', 'Volume_MA']):
        axes[2, 0].bar(data.index, data['Volume'], alpha=0.6, label='成交量')
        axes[2, 0].plot(data.index, data['Volume_MA'], color='red', label='成交量均线')
        axes[2, 0].set_title('成交量')
        axes[2, 0].legend()
        axes[2, 0].grid(True, alpha=0.3)

    # 历史波动率
    if 'Historical_Volatility' in data.columns:
        axes[2, 1].plot(data.index, data['Historical_Volatility'])
        axes[2, 1].set_title('历史波动率')
        axes[2, 1].grid(True, alpha=0.3)

    # 收益率分布
    if 'Returns' in data.columns:
        returns_clean = data['Returns'].dropna()
        axes[3, 0].hist(returns_clean, bins=50, alpha=0.7, density=True)
        axes[3, 0].axvline(returns_clean.mean(), color='red', linestyle='--', label='均值')
        axes[3, 0].set_title('收益率分布')
        axes[3, 0].legend()
        axes[3, 0].grid(True, alpha=0.3)

    # QQ图
    if 'Returns' in data.columns:
        from scipy import stats
        returns_clean = data['Returns'].dropna()
        stats.probplot(returns_clean, dist="norm", plot=axes[3, 1])
        axes[3, 1].set_title('收益率QQ图')
        axes[3, 1].grid(True, alpha=0.3)

    plt.tight_layout()
    plt.show()

# 绘制技术指标图表
if apple_complete is not None:
    plot_technical_indicators(apple_complete.tail(252))
```

### 3.5.4 相关性分析可视化

```python
def plot_correlation_analysis(stock_data, figsize=(12, 10)):
    """
    绘制股票相关性分析图
    """
    # 计算收益率
    returns_data = pd.DataFrame()
    for symbol, data in stock_data.items():
        returns_data[symbol] = data['Close'].pct_change()

    returns_data = returns_data.dropna()

    fig, axes = plt.subplots(2, 2, figsize=figsize)

    # 相关性热力图
    correlation_matrix = returns_data.corr()
    sns.heatmap(correlation_matrix, annot=True, cmap='coolwarm', center=0,
                ax=axes[0, 0], fmt='.3f')
    axes[0, 0].set_title('股票收益率相关性矩阵')

    # 累计收益率对比
    cumulative_returns = (1 + returns_data).cumprod()
    cumulative_returns.plot(ax=axes[0, 1])
    axes[0, 1].set_title('累计收益率对比')
    axes[0, 1].legend()
    axes[0, 1].grid(True, alpha=0.3)

    # 收益率分布对比
    returns_data.plot.hist(bins=50, alpha=0.7, ax=axes[1, 0])
    axes[1, 0].set_title('收益率分布对比')
    axes[1, 0].legend()
    axes[1, 0].grid(True, alpha=0.3)

    # 滚动相关性（以第一只股票为基准）
    if len(returns_data.columns) > 1:
        base_stock = returns_data.columns[0]
        rolling_corr = pd.DataFrame()

        for col in returns_data.columns[1:]:
            rolling_corr[f'{base_stock} vs {col}'] = (
                returns_data[base_stock].rolling(60)
                .corr(returns_data[col])
            )

        rolling_corr.plot(ax=axes[1, 1])
        axes[1, 1].set_title(f'60日滚动相关性 (基准: {base_stock})')
        axes[1, 1].legend()
        axes[1, 1].grid(True, alpha=0.3)

    plt.tight_layout()
    plt.show()

    return returns_data, correlation_matrix

# 相关性分析
if stock_data:
    returns_df, corr_matrix = plot_correlation_analysis(stock_data)
    print("收益率统计:")
    print(returns_df.describe())
```

## 📚 本章小结

在本章上篇中，我们学习了 Python 金融编程的基础知识，包括：

1. **环境配置**：了解了 Python 金融生态系统和开发环境
2. **数据获取**：掌握了从多个数据源获取金融数据的方法
3. **数据清洗**：学会了数据质量检查和清洗技术
4. **技术指标**：实现了常用的技术分析指标
5. **数据可视化**：掌握了专业的金融图表绘制方法

这些技能为后续的策略开发和回测分析奠定了坚实的基础。

## 🔗 下一章预告

在[第三章：Python 金融编程（下篇）](03_python_finance_programming_part2.md)中，我们将学习更高级的内容，包括机器学习在金融中的应用、策略回测框架的构建，以及交易执行系统的设计。

## 💡 实践练习

1. **数据获取练习**：选择 5 只不同行业的股票，获取它们过去 2 年的数据

2. **技术指标计算**：为这些股票计算所有技术指标，并分析指标的有效性

3. **可视化实践**：创建一个综合分析报告，包含价格图表、技术指标图和相关性分析

4. **数据清洗实践**：识别和处理数据中的异常值，比较处理前后的差异

5. **自定义指标**：尝试实现一个自定义的技术指标（如自适应移动平均）

---

**建议学习时间：** 4-5 天  
**前置章节：** [第二章：金融数学与统计基础](02_financial_mathematics.md)  
**下一章：** [Python 金融编程（下篇）](03_python_finance_programming_part2.md)
