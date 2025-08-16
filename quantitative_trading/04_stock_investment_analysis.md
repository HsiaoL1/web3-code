# 第四章：股票投资与分析（上篇）

## 📖 章节概述

本章将全面介绍股票投资的理论基础和分析方法，包括股票市场基础知识、基本面分析、技术分析以及不同的投资策略。这是量化交易中股票策略开发的重要基础。

## 🎯 学习目标

- 掌握股票市场的基础知识和运作机制
- 学会基本面分析的方法和指标
- 熟练运用技术分析工具
- 理解价值投资和成长投资的差异
- 掌握股票估值的基本方法
- 为量化选股策略开发打下基础

## 4.1 股票市场基础

### 4.1.1 股票的基本概念

**股票的定义**
股票是股份公司发行的所有权凭证，代表持有者对公司的所有权份额。持有股票意味着：

- 享有公司盈利分配权（股息）
- 参与公司重大决策的投票权
- 公司清算时的剩余资产分配权
- 股票价格波动带来的资本利得

**股票的分类**

```python
import pandas as pd
import numpy as np
import matplotlib.pyplot as plt

# 股票分类体系
stock_classifications = {
    "按股东权利": {
        "普通股": {
            "特点": ["有投票权", "享受分红", "承担风险"],
            "占比": "95%以上"
        },
        "优先股": {
            "特点": ["优先分红", "固定股息", "一般无投票权"],
            "占比": "较少"
        }
    },
    "按发行地区": {
        "A股": {
            "定义": "境内发行，人民币交易",
            "交易所": ["上海证券交易所", "深圳证券交易所"],
            "投资者": "主要为境内投资者"
        },
        "B股": {
            "定义": "境内发行，外币交易",
            "特点": "对外开放，现已式微"
        },
        "H股": {
            "定义": "境外发行（香港）",
            "特点": "港币交易，国际化程度高"
        }
    },
    "按公司规模": {
        "大盘股": {
            "市值范围": "> 500亿元",
            "特点": ["稳定性强", "流动性好", "业绩相对稳定"]
        },
        "中盘股": {
            "市值范围": "50-500亿元",
            "特点": ["成长性较好", "波动性适中"]
        },
        "小盘股": {
            "市值范围": "< 50亿元",
            "特点": ["成长潜力大", "波动性高", "流动性相对较差"]
        }
    }
}

def display_stock_classification():
    """展示股票分类"""
    for category, subcategories in stock_classifications.items():
        print(f"\n=== {category} ===")
        for subcat, details in subcategories.items():
            print(f"\n{subcat}:")
            for key, value in details.items():
                if isinstance(value, list):
                    print(f"  {key}: {', '.join(value)}")
                else:
                    print(f"  {key}: {value}")

display_stock_classification()
```

### 4.1.2 股票市场结构

**一级市场（发行市场）**

```python
class IPOProcess:
    """IPO流程模拟"""

    def __init__(self, company_name, share_count, issue_price):
        self.company_name = company_name
        self.share_count = share_count
        self.issue_price = issue_price
        self.stages = []

    def add_stage(self, stage_name, description, duration_days):
        """添加IPO阶段"""
        self.stages.append({
            "stage": stage_name,
            "description": description,
            "duration": duration_days
        })

    def calculate_fundraising(self):
        """计算募资金额"""
        return self.share_count * self.issue_price

    def display_ipo_timeline(self):
        """显示IPO时间线"""
        print(f"\n{self.company_name} IPO流程:")
        print(f"发行股数: {self.share_count:,} 股")
        print(f"发行价格: {self.issue_price} 元")
        print(f"募资金额: {self.calculate_fundraising():,.0f} 元")

        print("\nIPO流程时间线:")
        cumulative_days = 0
        for stage in self.stages:
            print(f"第{cumulative_days+1}-{cumulative_days+stage['duration']}天: "
                 f"{stage['stage']} - {stage['description']}")
            cumulative_days += stage['duration']

# IPO流程示例
ipo_example = IPOProcess("科技创新公司", 100_000_000, 25.50)
ipo_example.add_stage("辅导期", "券商辅导，规范公司治理", 180)
ipo_example.add_stage("申报期", "提交申请材料", 30)
ipo_example.add_stage("审核期", "监管机构审核", 120)
ipo_example.add_stage("路演期", "向投资者推介", 15)
ipo_example.add_stage("定价期", "确定发行价格", 5)
ipo_example.add_stage("缴款期", "投资者缴款认购", 3)
ipo_example.add_stage("上市交易", "正式在交易所挂牌", 1)

ipo_example.display_ipo_timeline()
```

**二级市场（交易市场）**

```python
class StockExchange:
    """证券交易所模拟"""

    def __init__(self, name, trading_hours):
        self.name = name
        self.trading_hours = trading_hours
        self.listed_stocks = {}
        self.daily_stats = {
            "total_volume": 0,
            "total_turnover": 0,
            "advancing_stocks": 0,
            "declining_stocks": 0,
            "unchanged_stocks": 0
        }

    def list_stock(self, stock_code, company_name, listing_price, share_count):
        """股票上市"""
        self.listed_stocks[stock_code] = {
            "company_name": company_name,
            "listing_price": listing_price,
            "current_price": listing_price,
            "share_count": share_count,
            "market_cap": listing_price * share_count,
            "listing_date": pd.Timestamp.now()
        }
        print(f"✅ {company_name}({stock_code}) 成功在{self.name}上市")

    def update_stock_price(self, stock_code, new_price, volume):
        """更新股票价格"""
        if stock_code in self.listed_stocks:
            stock = self.listed_stocks[stock_code]
            old_price = stock["current_price"]

            stock["current_price"] = new_price
            stock["market_cap"] = new_price * stock["share_count"]

            # 更新市场统计
            self.daily_stats["total_volume"] += volume
            self.daily_stats["total_turnover"] += new_price * volume

            if new_price > old_price:
                self.daily_stats["advancing_stocks"] += 1
            elif new_price < old_price:
                self.daily_stats["declining_stocks"] += 1
            else:
                self.daily_stats["unchanged_stocks"] += 1

            change_pct = (new_price - old_price) / old_price * 100
            print(f"{stock['company_name']}({stock_code}): "
                 f"{old_price:.2f} → {new_price:.2f} ({change_pct:+.2f}%)")

    def get_market_summary(self):
        """获取市场概况"""
        total_stocks = len(self.listed_stocks)
        total_market_cap = sum(stock["market_cap"] for stock in self.listed_stocks.values())

        return {
            "exchange_name": self.name,
            "listed_companies": total_stocks,
            "total_market_cap": total_market_cap,
            "daily_volume": self.daily_stats["total_volume"],
            "daily_turnover": self.daily_stats["total_turnover"],
            "advancing_stocks": self.daily_stats["advancing_stocks"],
            "declining_stocks": self.daily_stats["declining_stocks"]
        }

# 模拟交易所运作
shanghai_exchange = StockExchange("上海证券交易所", "09:30-15:00")

# 股票上市
shanghai_exchange.list_stock("600519", "贵州茅台", 1000.00, 1_256_000_000)
shanghai_exchange.list_stock("000001", "平安银行", 15.50, 19_400_000_000)
shanghai_exchange.list_stock("000858", "五粮液", 200.00, 3_860_000_000)

# 模拟价格变动
shanghai_exchange.update_stock_price("600519", 1050.00, 2_000_000)
shanghai_exchange.update_stock_price("000001", 15.80, 50_000_000)
shanghai_exchange.update_stock_price("000858", 195.00, 8_000_000)

# 显示市场概况
market_summary = shanghai_exchange.get_market_summary()
print(f"\n=== {market_summary['exchange_name']} 市场概况 ===")
print(f"上市公司数: {market_summary['listed_companies']}")
print(f"总市值: {market_summary['total_market_cap']:,.0f} 元")
print(f"日成交量: {market_summary['daily_volume']:,} 股")
print(f"日成交额: {market_summary['daily_turnover']:,.0f} 元")
print(f"上涨股票: {market_summary['advancing_stocks']}")
print(f"下跌股票: {market_summary['declining_stocks']}")
```

### 4.1.3 股票交易机制

```python
class OrderBook:
    """股票订单簿"""

    def __init__(self, stock_code):
        self.stock_code = stock_code
        self.bid_orders = []  # 买单（按价格降序）
        self.ask_orders = []  # 卖单（按价格升序）
        self.trade_history = []

    def add_buy_order(self, price, volume, order_id):
        """添加买单"""
        order = {"price": price, "volume": volume, "order_id": order_id, "type": "buy"}
        self.bid_orders.append(order)
        self.bid_orders.sort(key=lambda x: x["price"], reverse=True)  # 按价格降序
        print(f"买单添加: {order_id} - {volume}股 @ {price}元")

    def add_sell_order(self, price, volume, order_id):
        """添加卖单"""
        order = {"price": price, "volume": volume, "order_id": order_id, "type": "sell"}
        self.ask_orders.append(order)
        self.ask_orders.sort(key=lambda x: x["price"])  # 按价格升序
        print(f"卖单添加: {order_id} - {volume}股 @ {price}元")

    def match_orders(self):
        """撮合交易"""
        while self.bid_orders and self.ask_orders:
            best_bid = self.bid_orders[0]
            best_ask = self.ask_orders[0]

            # 如果买价 >= 卖价，可以成交
            if best_bid["price"] >= best_ask["price"]:
                trade_price = best_ask["price"]  # 以卖价成交
                trade_volume = min(best_bid["volume"], best_ask["volume"])

                # 记录交易
                trade = {
                    "price": trade_price,
                    "volume": trade_volume,
                    "buy_order": best_bid["order_id"],
                    "sell_order": best_ask["order_id"],
                    "timestamp": pd.Timestamp.now()
                }
                self.trade_history.append(trade)

                print(f"✅ 成交: {trade_volume}股 @ {trade_price}元 "
                     f"(买单:{best_bid['order_id']} - 卖单:{best_ask['order_id']})")

                # 更新订单量
                best_bid["volume"] -= trade_volume
                best_ask["volume"] -= trade_volume

                # 移除已完全成交的订单
                if best_bid["volume"] == 0:
                    self.bid_orders.pop(0)
                if best_ask["volume"] == 0:
                    self.ask_orders.pop(0)
            else:
                break  # 无法继续撮合

    def get_best_bid_ask(self):
        """获取最优买卖价"""
        best_bid = self.bid_orders[0]["price"] if self.bid_orders else None
        best_ask = self.ask_orders[0]["price"] if self.ask_orders else None
        return best_bid, best_ask

    def get_market_depth(self, levels=5):
        """获取市场深度"""
        depth = {
            "bids": self.bid_orders[:levels],
            "asks": self.ask_orders[:levels]
        }
        return depth

    def display_order_book(self):
        """显示订单簿"""
        print(f"\n=== {self.stock_code} 订单簿 ===")

        # 显示卖单（价格从高到低）
        print("卖单 (ASK):")
        for i, order in enumerate(reversed(self.ask_orders[:5])):
            print(f"  {order['price']:.2f} × {order['volume']:,}")

        print("-" * 20)

        # 显示买单（价格从高到低）
        print("买单 (BID):")
        for order in self.bid_orders[:5]:
            print(f"  {order['price']:.2f} × {order['volume']:,}")

        # 显示最新成交
        if self.trade_history:
            latest_trade = self.trade_history[-1]
            print(f"\n最新成交: {latest_trade['volume']}股 @ {latest_trade['price']:.2f}元")

# 模拟订单簿交易
order_book = OrderBook("000001")

# 添加买单
order_book.add_buy_order(15.80, 10000, "BUY_001")
order_book.add_buy_order(15.75, 20000, "BUY_002")
order_book.add_buy_order(15.78, 15000, "BUY_003")

# 添加卖单
order_book.add_sell_order(15.85, 8000, "SELL_001")
order_book.add_sell_order(15.82, 12000, "SELL_002")
order_book.add_sell_order(15.90, 25000, "SELL_003")

# 显示订单簿
order_book.display_order_book()

# 添加能够成交的订单
print("\n添加市价买单...")
order_book.add_buy_order(15.85, 5000, "BUY_004")
order_book.match_orders()

# 再次显示订单簿
order_book.display_order_book()
```

## 4.2 基本面分析

### 4.2.1 财务报表分析

```python
class FinancialAnalysis:
    """财务分析工具"""

    def __init__(self, company_name):
        self.company_name = company_name
        self.income_statement = {}
        self.balance_sheet = {}
        self.cash_flow = {}

    def load_financial_data(self, year, income_data, balance_data, cashflow_data):
        """加载财务数据"""
        self.income_statement[year] = income_data
        self.balance_sheet[year] = balance_data
        self.cash_flow[year] = cashflow_data

    def calculate_profitability_ratios(self, year):
        """计算盈利能力指标"""
        income = self.income_statement[year]
        balance = self.balance_sheet[year]

        ratios = {}

        # 毛利率
        if income.get('revenue') and income.get('cost_of_sales'):
            ratios['gross_margin'] = (income['revenue'] - income['cost_of_sales']) / income['revenue']

        # 净利率
        if income.get('revenue') and income.get('net_income'):
            ratios['net_margin'] = income['net_income'] / income['revenue']

        # ROA (资产收益率)
        if income.get('net_income') and balance.get('total_assets'):
            ratios['roa'] = income['net_income'] / balance['total_assets']

        # ROE (净资产收益率)
        if income.get('net_income') and balance.get('shareholders_equity'):
            ratios['roe'] = income['net_income'] / balance['shareholders_equity']

        return ratios

    def calculate_liquidity_ratios(self, year):
        """计算流动性指标"""
        balance = self.balance_sheet[year]

        ratios = {}

        # 流动比率
        if balance.get('current_assets') and balance.get('current_liabilities'):
            ratios['current_ratio'] = balance['current_assets'] / balance['current_liabilities']

        # 速动比率
        if (balance.get('current_assets') and balance.get('inventory') and
            balance.get('current_liabilities')):
            quick_assets = balance['current_assets'] - balance['inventory']
            ratios['quick_ratio'] = quick_assets / balance['current_liabilities']

        # 现金比率
        if balance.get('cash_and_equivalents') and balance.get('current_liabilities'):
            ratios['cash_ratio'] = balance['cash_and_equivalents'] / balance['current_liabilities']

        return ratios

    def calculate_leverage_ratios(self, year):
        """计算杠杆指标"""
        balance = self.balance_sheet[year]

        ratios = {}

        # 资产负债率
        if balance.get('total_liabilities') and balance.get('total_assets'):
            ratios['debt_to_assets'] = balance['total_liabilities'] / balance['total_assets']

        # 股债比
        if balance.get('total_liabilities') and balance.get('shareholders_equity'):
            ratios['debt_to_equity'] = balance['total_liabilities'] / balance['shareholders_equity']

        return ratios

    def calculate_efficiency_ratios(self, year):
        """计算效率指标"""
        income = self.income_statement[year]
        balance = self.balance_sheet[year]

        ratios = {}

        # 资产周转率
        if income.get('revenue') and balance.get('total_assets'):
            ratios['asset_turnover'] = income['revenue'] / balance['total_assets']

        # 存货周转率
        if income.get('cost_of_sales') and balance.get('inventory'):
            ratios['inventory_turnover'] = income['cost_of_sales'] / balance['inventory']

        # 应收账款周转率
        if income.get('revenue') and balance.get('accounts_receivable'):
            ratios['receivables_turnover'] = income['revenue'] / balance['accounts_receivable']

        return ratios

    def comprehensive_analysis(self, year):
        """综合财务分析"""
        profitability = self.calculate_profitability_ratios(year)
        liquidity = self.calculate_liquidity_ratios(year)
        leverage = self.calculate_leverage_ratios(year)
        efficiency = self.calculate_efficiency_ratios(year)

        analysis = {
            "company": self.company_name,
            "year": year,
            "profitability": profitability,
            "liquidity": liquidity,
            "leverage": leverage,
            "efficiency": efficiency
        }

        return analysis

    def display_analysis(self, year):
        """显示财务分析结果"""
        analysis = self.comprehensive_analysis(year)

        print(f"\n=== {self.company_name} {year}年财务分析 ===")

        print("\n📈 盈利能力指标:")
        for ratio, value in analysis["profitability"].items():
            print(f"  {ratio}: {value:.2%}")

        print("\n💧 流动性指标:")
        for ratio, value in analysis["liquidity"].items():
            print(f"  {ratio}: {value:.2f}")

        print("\n🏦 杠杆指标:")
        for ratio, value in analysis["leverage"].items():
            print(f"  {ratio}: {value:.2%}")

        print("\n⚡ 效率指标:")
        for ratio, value in analysis["efficiency"].items():
            print(f"  {ratio}: {value:.2f}")

# 示例：财务分析
company_analysis = FinancialAnalysis("科技股份有限公司")

# 2023年财务数据（单位：亿元）
income_2023 = {
    'revenue': 1000,
    'cost_of_sales': 600,
    'gross_profit': 400,
    'operating_expenses': 250,
    'operating_income': 150,
    'net_income': 120
}

balance_2023 = {
    'total_assets': 2000,
    'current_assets': 800,
    'cash_and_equivalents': 200,
    'accounts_receivable': 150,
    'inventory': 100,
    'total_liabilities': 1200,
    'current_liabilities': 400,
    'shareholders_equity': 800
}

cashflow_2023 = {
    'operating_cash_flow': 180,
    'investing_cash_flow': -100,
    'financing_cash_flow': -50,
    'net_cash_flow': 30
}

company_analysis.load_financial_data(2023, income_2023, balance_2023, cashflow_2023)
company_analysis.display_analysis(2023)
```

### 4.2.2 估值方法

```python
class StockValuation:
    """股票估值工具"""

    def __init__(self, stock_code, company_name):
        self.stock_code = stock_code
        self.company_name = company_name
        self.current_price = None
        self.shares_outstanding = None

    def set_market_data(self, current_price, shares_outstanding):
        """设置市场数据"""
        self.current_price = current_price
        self.shares_outstanding = shares_outstanding

    def pe_valuation(self, eps, industry_pe_ratio):
        """市盈率估值法"""
        estimated_value = eps * industry_pe_ratio

        if self.current_price:
            discount_premium = (estimated_value - self.current_price) / self.current_price

        return {
            "method": "PE估值法",
            "eps": eps,
            "industry_pe": industry_pe_ratio,
            "estimated_value": estimated_value,
            "current_price": self.current_price,
            "discount_premium": discount_premium if self.current_price else None
        }

    def pb_valuation(self, book_value_per_share, industry_pb_ratio):
        """市净率估值法"""
        estimated_value = book_value_per_share * industry_pb_ratio

        discount_premium = None
        if self.current_price:
            discount_premium = (estimated_value - self.current_price) / self.current_price

        return {
            "method": "PB估值法",
            "book_value_per_share": book_value_per_share,
            "industry_pb": industry_pb_ratio,
            "estimated_value": estimated_value,
            "current_price": self.current_price,
            "discount_premium": discount_premium
        }

    def dcf_valuation(self, free_cash_flows, terminal_growth_rate, discount_rate):
        """现金流贴现估值法"""
        # 计算未来现金流现值
        present_values = []
        for i, fcf in enumerate(free_cash_flows):
            pv = fcf / ((1 + discount_rate) ** (i + 1))
            present_values.append(pv)

        # 计算终值
        terminal_fcf = free_cash_flows[-1] * (1 + terminal_growth_rate)
        terminal_value = terminal_fcf / (discount_rate - terminal_growth_rate)
        terminal_pv = terminal_value / ((1 + discount_rate) ** len(free_cash_flows))

        # 企业价值
        enterprise_value = sum(present_values) + terminal_pv

        # 每股价值（假设无净债务）
        value_per_share = enterprise_value / self.shares_outstanding if self.shares_outstanding else None

        discount_premium = None
        if self.current_price and value_per_share:
            discount_premium = (value_per_share - self.current_price) / self.current_price

        return {
            "method": "DCF估值法",
            "free_cash_flows": free_cash_flows,
            "discount_rate": discount_rate,
            "terminal_growth_rate": terminal_growth_rate,
            "enterprise_value": enterprise_value,
            "value_per_share": value_per_share,
            "current_price": self.current_price,
            "discount_premium": discount_premium
        }

    def dividend_discount_model(self, dividends, growth_rate, required_return):
        """股息贴现模型"""
        if growth_rate >= required_return:
            return {"error": "增长率不能大于等于要求回报率"}

        # Gordon增长模型
        next_dividend = dividends[-1] * (1 + growth_rate)
        intrinsic_value = next_dividend / (required_return - growth_rate)

        discount_premium = None
        if self.current_price:
            discount_premium = (intrinsic_value - self.current_price) / self.current_price

        return {
            "method": "股息贴现模型",
            "last_dividend": dividends[-1],
            "growth_rate": growth_rate,
            "required_return": required_return,
            "intrinsic_value": intrinsic_value,
            "current_price": self.current_price,
            "discount_premium": discount_premium
        }

    def comprehensive_valuation(self, valuation_data):
        """综合估值"""
        valuations = []

        # PE估值
        if 'eps' in valuation_data and 'industry_pe' in valuation_data:
            pe_val = self.pe_valuation(valuation_data['eps'], valuation_data['industry_pe'])
            valuations.append(pe_val)

        # PB估值
        if 'book_value_per_share' in valuation_data and 'industry_pb' in valuation_data:
            pb_val = self.pb_valuation(valuation_data['book_value_per_share'],
                                     valuation_data['industry_pb'])
            valuations.append(pb_val)

        # DCF估值
        if all(key in valuation_data for key in ['free_cash_flows', 'terminal_growth_rate', 'discount_rate']):
            dcf_val = self.dcf_valuation(valuation_data['free_cash_flows'],
                                       valuation_data['terminal_growth_rate'],
                                       valuation_data['discount_rate'])
            valuations.append(dcf_val)

        # 股息贴现模型
        if all(key in valuation_data for key in ['dividends', 'dividend_growth_rate', 'required_return']):
            ddm_val = self.dividend_discount_model(valuation_data['dividends'],
                                                 valuation_data['dividend_growth_rate'],
                                                 valuation_data['required_return'])
            if 'error' not in ddm_val:
                valuations.append(ddm_val)

        return valuations

    def display_valuation_results(self, valuations):
        """显示估值结果"""
        print(f"\n=== {self.company_name}({self.stock_code}) 估值分析 ===")
        print(f"当前股价: {self.current_price:.2f}元" if self.current_price else "当前股价: 未设置")

        estimated_values = []

        for val in valuations:
            print(f"\n📊 {val['method']}:")

            if val['method'] == "PE估值法":
                print(f"  每股收益(EPS): {val['eps']:.2f}元")
                print(f"  行业平均PE: {val['industry_pe']:.1f}倍")
                print(f"  估值价格: {val['estimated_value']:.2f}元")
                estimated_values.append(val['estimated_value'])

            elif val['method'] == "PB估值法":
                print(f"  每股净资产: {val['book_value_per_share']:.2f}元")
                print(f"  行业平均PB: {val['industry_pb']:.1f}倍")
                print(f"  估值价格: {val['estimated_value']:.2f}元")
                estimated_values.append(val['estimated_value'])

            elif val['method'] == "DCF估值法":
                print(f"  企业价值: {val['enterprise_value']:.0f}万元")
                print(f"  每股价值: {val['value_per_share']:.2f}元")
                if val['value_per_share']:
                    estimated_values.append(val['value_per_share'])

            elif val['method'] == "股息贴现模型":
                print(f"  最新股息: {val['last_dividend']:.2f}元")
                print(f"  内在价值: {val['intrinsic_value']:.2f}元")
                estimated_values.append(val['intrinsic_value'])

            if val.get('discount_premium') is not None:
                status = "高估" if val['discount_premium'] < 0 else "低估"
                print(f"  相对当前价格: {status} {abs(val['discount_premium']):.1%}")

        # 计算平均估值
        if estimated_values:
            avg_value = np.mean(estimated_values)
            print(f"\n💡 平均估值: {avg_value:.2f}元")

            if self.current_price:
                avg_discount_premium = (avg_value - self.current_price) / self.current_price
                status = "高估" if avg_discount_premium < 0 else "低估"
                print(f"综合判断: {status} {abs(avg_discount_premium):.1%}")

# 估值示例
valuation = StockValuation("000001", "平安银行")
valuation.set_market_data(current_price=15.80, shares_outstanding=19_400_000_000)

# 估值参数
valuation_params = {
    'eps': 1.75,  # 每股收益
    'industry_pe': 6.5,  # 银行业平均PE
    'book_value_per_share': 18.50,  # 每股净资产
    'industry_pb': 0.85,  # 银行业平均PB
    'free_cash_flows': [120, 135, 150, 165, 180],  # 未来5年自由现金流（亿元）
    'terminal_growth_rate': 0.03,  # 永续增长率
    'discount_rate': 0.10,  # 折现率
    'dividends': [0.80, 0.85, 0.90],  # 历史股息
    'dividend_growth_rate': 0.06,  # 股息增长率
    'required_return': 0.12  # 要求回报率
}

# 执行估值
results = valuation.comprehensive_valuation(valuation_params)
valuation.display_valuation_results(results)
```

## 4.3 技术分析进阶

### 4.3.1 K 线形态分析

```python
class CandlestickPatterns:
    """K线形态识别"""

    def __init__(self):
        self.patterns = {}

    def identify_doji(self, open_price, close_price, high_price, low_price):
        """识别十字星"""
        body = abs(close_price - open_price)
        total_range = high_price - low_price

        # 实体小于总波幅的5%被认为是十字星
        if total_range > 0 and body / total_range < 0.05:
            return True, "十字星 - 市场犹豫不决的信号"
        return False, None

    def identify_hammer(self, open_price, close_price, high_price, low_price):
        """识别锤子线"""
        body = abs(close_price - open_price)
        upper_shadow = high_price - max(open_price, close_price)
        lower_shadow = min(open_price, close_price) - low_price

        # 锤子线特征：下影线长，上影线短，实体小
        if (lower_shadow >= 2 * body and
            upper_shadow <= body * 0.5 and
            body > 0):
            pattern_type = "锤子线" if close_price > open_price else "上吊线"
            return True, f"{pattern_type} - 潜在反转信号"
        return False, None

    def identify_engulfing(self, prev_ohlc, curr_ohlc):
        """识别吞没形态"""
        prev_o, prev_h, prev_l, prev_c = prev_ohlc
        curr_o, curr_h, curr_l, curr_c = curr_ohlc

        prev_body = abs(prev_c - prev_o)
        curr_body = abs(curr_c - curr_o)

        # 看涨吞没
        if (prev_c < prev_o and  # 前一天阴线
            curr_c > curr_o and  # 当天阳线
            curr_o < prev_c and  # 当天开盘价低于前一天收盘价
            curr_c > prev_o):    # 当天收盘价高于前一天开盘价
            return True, "看涨吞没 - 强烈的看涨信号"

        # 看跌吞没
        if (prev_c > prev_o and  # 前一天阳线
            curr_c < curr_o and  # 当天阴线
            curr_o > prev_c and  # 当天开盘价高于前一天收盘价
            curr_c < prev_o):    # 当天收盘价低于前一天开盘价
            return True, "看跌吞没 - 强烈的看跌信号"

        return False, None

    def identify_morning_star(self, day1_ohlc, day2_ohlc, day3_ohlc):
        """识别启明星形态"""
        d1_o, d1_h, d1_l, d1_c = day1_ohlc
        d2_o, d2_h, d2_l, d2_c = day2_ohlc
        d3_o, d3_h, d3_l, d3_c = day3_ohlc

        # 第一天：阴线
        # 第二天：十字星或小实体，向下跳空
        # 第三天：阳线，向上跳空，收盘价进入第一天实体

        day1_bearish = d1_c < d1_o
        day3_bullish = d3_c > d3_o

        day2_small_body = abs(d2_c - d2_o) < abs(d1_c - d1_o) * 0.3
        gap_down = d2_h < d1_c
        gap_up = d3_o > d2_h
        penetration = d3_c > (d1_o + d1_c) / 2

        if (day1_bearish and day3_bullish and day2_small_body and
            gap_down and gap_up and penetration):
            return True, "启明星 - 强烈的底部反转信号"

        return False, None

    def analyze_candlestick_data(self, ohlc_data):
        """分析K线数据中的形态"""
        patterns_found = []

        for i in range(len(ohlc_data)):
            o, h, l, c = ohlc_data[i]

            # 单根K线形态
            is_doji, doji_msg = self.identify_doji(o, c, h, l)
            if is_doji:
                patterns_found.append({"day": i, "pattern": doji_msg})

            is_hammer, hammer_msg = self.identify_hammer(o, c, h, l)
            if is_hammer:
                patterns_found.append({"day": i, "pattern": hammer_msg})

            # 双根K线形态
            if i > 0:
                is_engulfing, engulfing_msg = self.identify_engulfing(
                    ohlc_data[i-1], ohlc_data[i]
                )
                if is_engulfing:
                    patterns_found.append({"day": i, "pattern": engulfing_msg})

            # 三根K线形态
            if i > 1:
                is_morning_star, morning_star_msg = self.identify_morning_star(
                    ohlc_data[i-2], ohlc_data[i-1], ohlc_data[i]
                )
                if is_morning_star:
                    patterns_found.append({"day": i, "pattern": morning_star_msg})

        return patterns_found

# K线形态分析示例
pattern_analyzer = CandlestickPatterns()

# 模拟K线数据 (开盘价, 最高价, 最低价, 收盘价)
sample_ohlc = [
    (100, 105, 98, 99),   # 第0天：小阴线
    (99, 101, 95, 96),    # 第1天：阴线
    (96, 97, 94, 96.5),   # 第2天：十字星
    (97, 103, 96, 102),   # 第3天：阳线
    (102, 108, 101, 107), # 第4天：阳线
    (107, 108, 103, 103), # 第5天：上影线长的阴线
]

print("K线形态分析结果:")
patterns = pattern_analyzer.analyze_candlestick_data(sample_ohlc)

for pattern in patterns:
    print(f"第{pattern['day']}天: {pattern['pattern']}")
```

### 4.3.2 支撑阻力分析

```python
class SupportResistanceAnalyzer:
    """支撑阻力分析器"""

    def __init__(self, price_data, volume_data=None):
        self.price_data = np.array(price_data)
        self.volume_data = np.array(volume_data) if volume_data else None
        self.support_levels = []
        self.resistance_levels = []

    def find_local_extremes(self, window=5):
        """寻找局部极值点"""
        highs = []
        lows = []

        for i in range(window, len(self.price_data) - window):
            # 局部最高点
            if all(self.price_data[i] >= self.price_data[i-window:i]) and \
               all(self.price_data[i] >= self.price_data[i+1:i+window+1]):
                highs.append((i, self.price_data[i]))

            # 局部最低点
            if all(self.price_data[i] <= self.price_data[i-window:i]) and \
               all(self.price_data[i] <= self.price_data[i+1:i+window+1]):
                lows.append((i, self.price_data[i]))

        return highs, lows

    def cluster_levels(self, levels, tolerance=0.02):
        """聚类相近的价格水平"""
        if not levels:
            return []

        levels_sorted = sorted(levels, key=lambda x: x[1])
        clusters = []
        current_cluster = [levels_sorted[0]]

        for i in range(1, len(levels_sorted)):
            current_price = levels_sorted[i][1]
            cluster_avg = np.mean([level[1] for level in current_cluster])

            if abs(current_price - cluster_avg) / cluster_avg <= tolerance:
                current_cluster.append(levels_sorted[i])
            else:
                clusters.append(current_cluster)
                current_cluster = [levels_sorted[i]]

        clusters.append(current_cluster)
        return clusters

    def calculate_level_strength(self, cluster):
        """计算支撑/阻力水平的强度"""
        # 触及次数
        touch_count = len(cluster)

        # 平均价格
        avg_price = np.mean([level[1] for level in cluster])

        # 时间跨度
        indices = [level[0] for level in cluster]
        time_span = max(indices) - min(indices) if len(indices) > 1 else 0

        # 成交量权重（如果有成交量数据）
        volume_weight = 1.0
        if self.volume_data is not None:
            volumes = [self.volume_data[level[0]] for level in cluster]
            volume_weight = np.mean(volumes) / np.mean(self.volume_data)

        # 综合强度评分
        strength = touch_count * 0.4 + (time_span / len(self.price_data)) * 0.3 + volume_weight * 0.3

        return {
            "price": avg_price,
            "strength": strength,
            "touch_count": touch_count,
            "time_span": time_span,
            "volume_weight": volume_weight
        }

    def identify_support_resistance(self, min_touches=2, min_strength=0.5):
        """识别支撑和阻力位"""
        highs, lows = self.find_local_extremes()

        # 聚类阻力位（高点）
        resistance_clusters = self.cluster_levels(highs)
        for cluster in resistance_clusters:
            level_info = self.calculate_level_strength(cluster)
            if (level_info["touch_count"] >= min_touches and
                level_info["strength"] >= min_strength):
                self.resistance_levels.append(level_info)

        # 聚类支撑位（低点）
        support_clusters = self.cluster_levels(lows)
        for cluster in support_clusters:
            level_info = self.calculate_level_strength(cluster)
            if (level_info["touch_count"] >= min_touches and
                level_info["strength"] >= min_strength):
                self.support_levels.append(level_info)

        # 按强度排序
        self.resistance_levels.sort(key=lambda x: x["strength"], reverse=True)
        self.support_levels.sort(key=lambda x: x["strength"], reverse=True)

    def get_current_levels(self, current_price):
        """获取当前价格附近的支撑阻力位"""
        nearby_resistance = [level for level in self.resistance_levels
                           if level["price"] > current_price]
        nearby_support = [level for level in self.support_levels
                        if level["price"] < current_price]

        # 取最近的几个重要位置
        next_resistance = sorted(nearby_resistance, key=lambda x: x["price"])[:3]
        next_support = sorted(nearby_support, key=lambda x: x["price"], reverse=True)[:3]

        return next_support, next_resistance

    def display_analysis(self, current_price):
        """显示支撑阻力分析结果"""
        print(f"=== 支撑阻力分析 ===")
        print(f"当前价格: {current_price:.2f}")

        support_levels, resistance_levels = self.get_current_levels(current_price)

        print(f"\n📈 主要阻力位:")
        for i, level in enumerate(resistance_levels):
            distance = (level["price"] - current_price) / current_price * 100
            print(f"  {i+1}. {level['price']:.2f} (距离: +{distance:.1f}%, "
                 f"强度: {level['strength']:.2f}, 触及次数: {level['touch_count']})")

        print(f"\n📉 主要支撑位:")
        for i, level in enumerate(support_levels):
            distance = (current_price - level["price"]) / current_price * 100
            print(f"  {i+1}. {level['price']:.2f} (距离: -{distance:.1f}%, "
                 f"强度: {level['strength']:.2f}, 触及次数: {level['touch_count']})")

# 支撑阻力分析示例
np.random.seed(42)
price_trend = np.cumsum(np.random.randn(200) * 0.5) + 100
volume_data = np.random.randint(1000, 5000, 200)

sr_analyzer = SupportResistanceAnalyzer(price_trend, volume_data)
sr_analyzer.identify_support_resistance(min_touches=2, min_strength=0.3)

current_price = price_trend[-1]
sr_analyzer.display_analysis(current_price)
```

## 📚 本章小结

在本章上篇中，我们深入学习了股票投资的基础知识，包括：

1. **股票市场基础**：理解了股票的基本概念、市场结构和交易机制
2. **基本面分析**：掌握了财务报表分析和多种估值方法
3. **技术分析进阶**：学习了 K 线形态识别和支撑阻力分析

这些知识为开发量化选股策略和理解市场运作机制提供了重要基础。

## 🔗 下一章预告

在[第四章：股票投资与分析（下篇）](04_stock_investment_analysis_part2.md)中，我们将学习量化选股策略、多因子模型、事件驱动策略以及 A 股市场的特色分析。

## 💡 实践练习

1. **财务分析实践**：选择一只真实股票，获取其财务数据并进行全面分析

2. **估值练习**：使用多种估值方法为同一只股票估值，比较结果差异

3. **K 线形态识别**：在真实股票数据中识别各种 K 线形态

4. **支撑阻力分析**：分析某只股票的历史支撑阻力位，验证其有效性

5. **综合分析报告**：结合基本面和技术面，写一份完整的股票分析报告

---

**建议学习时间：** 5-6 天  
**前置章节：** [第三章：Python 金融编程（上篇）](03_python_finance_programming.md)  
**下一章：** [股票投资与分析（下篇）](04_stock_investment_analysis_part2.md)
