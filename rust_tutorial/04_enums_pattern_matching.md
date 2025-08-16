# Rust 枚举与模式匹配

## 1. 枚举基础

枚举允许你通过列举可能的成员来定义一个类型。枚举在 Rust 中非常强大，因为每个成员可以有不同的类型和数量的关联数据。

### 1.1 定义枚举

```rust
// 简单枚举
enum IpAddrKind {
    V4,
    V6,
}

// 带数据的枚举
enum IpAddr {
    V4(u8, u8, u8, u8),
    V6(String),
}

// 复杂枚举
enum Message {
    Quit,                       // 没有关联数据
    Move { x: i32, y: i32 },   // 命名字段
    Write(String),              // 单个字符串
    ChangeColor(i32, i32, i32), // 三个整数
}

fn main() {
    let four = IpAddrKind::V4;
    let six = IpAddrKind::V6;

    let home = IpAddr::V4(127, 0, 0, 1);
    let loopback = IpAddr::V6(String::from("::1"));

    let quit = Message::Quit;
    let move_msg = Message::Move { x: 10, y: 20 };
    let write = Message::Write(String::from("hello"));
    let change_color = Message::ChangeColor(255, 0, 0);
}
```

### 1.2 枚举方法

```rust
impl Message {
    fn call(&self) {
        match self {
            Message::Quit => println!("退出程序"),
            Message::Move { x, y } => println!("移动到坐标 ({}, {})", x, y),
            Message::Write(text) => println!("写入文本: {}", text),
            Message::ChangeColor(r, g, b) => println!("改变颜色为 RGB({}, {}, {})", r, g, b),
        }
    }

    fn message_type(&self) -> &str {
        match self {
            Message::Quit => "Quit",
            Message::Move { .. } => "Move",
            Message::Write(_) => "Write",
            Message::ChangeColor(_, _, _) => "ChangeColor",
        }
    }
}

fn main() {
    let msg = Message::Write(String::from("Hello, Rust!"));
    msg.call();
    println!("消息类型: {}", msg.message_type());
}
```

## 2. Option 枚举

`Option` 是 Rust 标准库中最重要的枚举之一，用于表示可能存在或不存在的值。

### 2.1 Option 的定义

```rust
// 标准库中 Option 的定义（简化）
enum Option<T> {
    Some(T),
    None,
}
```

### 2.2 使用 Option

```rust
fn main() {
    let some_number = Some(5);
    let some_string = Some("a string");
    let absent_number: Option<i32> = None;

    println!("{:?}", some_number);
    println!("{:?}", some_string);
    println!("{:?}", absent_number);

    // 使用 Option 的例子
    let numbers = vec![1, 2, 3, 4, 5];
    let index = 2;

    match get_element(&numbers, index) {
        Some(value) => println!("索引 {} 处的值是: {}", index, value),
        None => println!("索引 {} 超出范围", index),
    }
}

fn get_element(vec: &Vec<i32>, index: usize) -> Option<&i32> {
    if index < vec.len() {
        Some(&vec[index])
    } else {
        None
    }
}
```

### 2.3 Option 的常用方法

```rust
fn main() {
    let x: Option<&str> = Some("air");
    assert_eq!(x.is_some(), true);
    assert_eq!(x.is_none(), false);

    let y: Option<&str> = None;
    assert_eq!(y.is_some(), false);
    assert_eq!(y.is_none(), true);

    // unwrap - 获取值或 panic
    let some_value = Some(10);
    println!("值: {}", some_value.unwrap());

    // unwrap_or - 提供默认值
    let none_value: Option<i32> = None;
    println!("值或默认值: {}", none_value.unwrap_or(42));

    // map - 转换值
    let some_string = Some("hello");
    let some_len = some_string.map(|s| s.len());
    println!("字符串长度: {:?}", some_len);

    // and_then - 链式操作
    let result = Some("123")
        .and_then(|s| s.parse::<i32>().ok())
        .map(|n| n * 2);
    println!("结果: {:?}", result);
}
```

## 3. match 控制流

`match` 是 Rust 中强大的控制流结构，它允许你将一个值与一系列模式进行比较。

### 3.1 基本 match 语法

```rust
#[derive(Debug)]
enum Coin {
    Penny,
    Nickel,
    Dime,
    Quarter(UsState),
}

#[derive(Debug)]
enum UsState {
    Alabama,
    Alaska,
    California,
    // ... 其他州
}

fn value_in_cents(coin: Coin) -> u8 {
    match coin {
        Coin::Penny => {
            println!("幸运便士!");
            1
        },
        Coin::Nickel => 5,
        Coin::Dime => 10,
        Coin::Quarter(state) => {
            println!("来自 {:?} 州的25美分!", state);
            25
        },
    }
}

fn main() {
    let coin = Coin::Quarter(UsState::Alaska);
    println!("硬币价值: {} 美分", value_in_cents(coin));
}
```

### 3.2 匹配 Option<T>

```rust
fn plus_one(x: Option<i32>) -> Option<i32> {
    match x {
        None => None,
        Some(i) => Some(i + 1),
    }
}

fn main() {
    let five = Some(5);
    let six = plus_one(five);
    let none = plus_one(None);

    println!("five: {:?}", five);
    println!("six: {:?}", six);
    println!("none: {:?}", none);
}
```

### 3.3 穷尽性检查

```rust
fn describe_number(x: Option<i32>) {
    match x {
        Some(1) => println!("一"),
        Some(2) => println!("二"),
        Some(3) => println!("三"),
        Some(n) => println!("其他数字: {}", n),
        None => println!("没有数字"),
    }
}

// 使用 _ 通配符
fn handle_some_cases(x: Option<i32>) {
    match x {
        Some(1) => println!("一"),
        Some(2) => println!("二"),
        _ => println!("其他情况"), // 处理所有其他情况
    }
}

fn main() {
    describe_number(Some(1));
    describe_number(Some(42));
    describe_number(None);

    handle_some_cases(Some(1));
    handle_some_cases(Some(5));
    handle_some_cases(None);
}
```

## 4. if let 简洁控制流

当你只关心一种匹配情况时，`if let` 提供了比 `match` 更简洁的语法。

### 4.1 基本 if let 语法

```rust
fn main() {
    let some_value = Some(3);

    // 使用 match
    match some_value {
        Some(3) => println!("三"),
        _ => (),
    }

    // 使用 if let（更简洁）
    if let Some(3) = some_value {
        println!("三");
    }

    // 更复杂的例子
    let favorite_color: Option<&str> = None;
    let is_tuesday = false;
    let age: Result<u8, _> = "34".parse();

    if let Some(color) = favorite_color {
        println!("使用你最喜欢的颜色: {} 作为背景", color);
    } else if is_tuesday {
        println!("星期二是绿色的一天!");
    } else if let Ok(age) = age {
        if age > 30 {
            println!("使用紫色作为背景颜色");
        } else {
            println!("使用橙色作为背景颜色");
        }
    } else {
        println!("使用蓝色作为背景颜色");
    }
}
```

### 4.2 if let 与 else

```rust
#[derive(Debug)]
enum Message {
    Hello { id: i32 },
    Goodbye,
}

fn main() {
    let msg = Message::Hello { id: 5 };

    if let Message::Hello { id } = msg {
        println!("找到 id: {}", id);
    } else {
        println!("其他消息");
    }
}
```

## 5. while let 循环

```rust
fn main() {
    let mut stack = Vec::new();
    stack.push(1);
    stack.push(2);
    stack.push(3);

    // 使用 while let 处理栈
    while let Some(top) = stack.pop() {
        println!("弹出: {}", top);
    }

    println!("栈现在是空的");
}
```

## 6. 复杂模式匹配

### 6.1 解构结构体

```rust
struct Point {
    x: i32,
    y: i32,
}

fn main() {
    let p = Point { x: 0, y: 7 };

    match p {
        Point { x, y: 0 } => println!("在 x 轴上的点 x = {}", x),
        Point { x: 0, y } => println!("在 y 轴上的点 y = {}", y),
        Point { x, y } => println!("在其他位置的点 ({}, {})", x, y),
    }

    // 使用 if let 解构
    if let Point { x: 0, y } = p {
        println!("这个点在 y 轴上，y = {}", y);
    }
}
```

### 6.2 解构枚举

```rust
enum Color {
    Rgb(i32, i32, i32),
    Hsv(i32, i32, i32),
}

enum MessageType {
    Quit,
    Move { x: i32, y: i32 },
    Write(String),
    ChangeColor(Color),
}

fn main() {
    let msg = MessageType::ChangeColor(Color::Hsv(0, 160, 255));

    match msg {
        MessageType::Quit => println!("退出"),
        MessageType::Move { x, y } => println!("移动到 ({}, {})", x, y),
        MessageType::Write(text) => println!("文本消息: {}", text),
        MessageType::ChangeColor(Color::Rgb(r, g, b)) => {
            println!("改变颜色为红: {}, 绿: {}, 蓝: {}", r, g, b);
        },
        MessageType::ChangeColor(Color::Hsv(h, s, v)) => {
            println!("改变颜色为色调: {}, 饱和度: {}, 值: {}", h, s, v);
        },
    }
}
```

### 6.3 匹配守卫

```rust
fn main() {
    let num = Some(4);

    match num {
        Some(x) if x < 5 => println!("小于五: {}", x),
        Some(x) => println!("大于等于五: {}", x),
        None => (),
    }

    let x = 4;
    let y = false;

    match x {
        4 | 5 | 6 if y => println!("yes"),
        _ => println!("no"),
    }
}
```

## 7. Result 枚举和错误处理

```rust
use std::fs::File;
use std::io::ErrorKind;

fn main() {
    // 基本错误处理
    let f = File::open("hello.txt");

    let f = match f {
        Ok(file) => file,
        Err(error) => match error.kind() {
            ErrorKind::NotFound => match File::create("hello.txt") {
                Ok(fc) => fc,
                Err(e) => panic!("创建文件时出现问题: {:?}", e),
            },
            other_error => panic!("打开文件时出现问题: {:?}", other_error),
        },
    };

    // 使用 unwrap_or_else
    let f = File::open("hello.txt").unwrap_or_else(|error| {
        if error.kind() == ErrorKind::NotFound {
            File::create("hello.txt").unwrap_or_else(|error| {
                panic!("创建文件时出现问题: {:?}", error);
            })
        } else {
            panic!("打开文件时出现问题: {:?}", error);
        }
    });
}

// 传播错误
fn read_username_from_file() -> Result<String, std::io::Error> {
    use std::fs;
    fs::read_to_string("hello.txt")
}

// 使用 ? 操作符
fn read_username_from_file_short() -> Result<String, std::io::Error> {
    use std::fs;
    let mut file = File::open("hello.txt")?;
    let mut username = String::new();
    file.read_to_string(&mut username)?;
    Ok(username)
}
```

## 8. 实践项目：简单计算器

```rust
#[derive(Debug)]
enum Operation {
    Add,
    Subtract,
    Multiply,
    Divide,
}

#[derive(Debug)]
enum CalculatorError {
    DivisionByZero,
    InvalidOperation,
}

struct Calculator;

impl Calculator {
    fn new() -> Self {
        Calculator
    }

    fn calculate(&self, a: f64, b: f64, op: Operation) -> Result<f64, CalculatorError> {
        match op {
            Operation::Add => Ok(a + b),
            Operation::Subtract => Ok(a - b),
            Operation::Multiply => Ok(a * b),
            Operation::Divide => {
                if b == 0.0 {
                    Err(CalculatorError::DivisionByZero)
                } else {
                    Ok(a / b)
                }
            }
        }
    }

    fn parse_operation(op_str: &str) -> Option<Operation> {
        match op_str {
            "+" => Some(Operation::Add),
            "-" => Some(Operation::Subtract),
            "*" => Some(Operation::Multiply),
            "/" => Some(Operation::Divide),
            _ => None,
        }
    }
}

fn main() {
    let calc = Calculator::new();

    let operations = vec![
        (10.0, 5.0, "+"),
        (10.0, 3.0, "-"),
        (4.0, 6.0, "*"),
        (15.0, 3.0, "/"),
        (10.0, 0.0, "/"),
        (5.0, 2.0, "%"),
    ];

    for (a, b, op_str) in operations {
        match Calculator::parse_operation(op_str) {
            Some(operation) => {
                match calc.calculate(a, b, operation) {
                    Ok(result) => println!("{} {} {} = {}", a, op_str, b, result),
                    Err(CalculatorError::DivisionByZero) => {
                        println!("{} {} {} = 错误: 除零", a, op_str, b);
                    },
                    Err(CalculatorError::InvalidOperation) => {
                        println!("{} {} {} = 错误: 无效操作", a, op_str, b);
                    },
                }
            },
            None => println!("{} {} {} = 错误: 未知操作符", a, op_str, b),
        }
    }
}
```

## 9. 高级模式匹配技巧

### 9.1 @ 绑定

```rust
#[derive(Debug)]
enum Message {
    Hello { id: i32 },
}

fn main() {
    let msg = Message::Hello { id: 5 };

    match msg {
        Message::Hello { id: id_variable @ 3..=7 } => {
            println!("找到一个在范围内的 id: {}", id_variable);
        },
        Message::Hello { id: 10..=12 } => {
            println!("找到一个在另一个范围内的 id");
        },
        Message::Hello { id } => {
            println!("找到其他 id: {}", id);
        },
    }
}
```

### 9.2 忽略值

```rust
fn main() {
    let numbers = (2, 4, 8, 16, 32);

    match numbers {
        (first, _, third, _, fifth) => {
            println!("一些数字: {}, {}, {}", first, third, fifth);
        },
    }

    let s = Some(String::from("Hello!"));

    if let Some(_) = s {
        println!("找到一个字符串");
    }

    // s 在这里仍然有效，因为我们没有使用 Some(s)
    println!("{:?}", s);
}
```

枚举和模式匹配是 Rust 中极其强大的特性，它们让你能够以类型安全的方式处理不同的情况。通过掌握这些概念，你将能够编写更加健壮和表达力强的代码。
