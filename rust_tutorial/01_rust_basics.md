# Rust 基础篇

## 1. Rust 简介

Rust 是一种系统编程语言，专注于安全性、速度和并发性。它由 Mozilla Research 开发，于 2015 年首次稳定发布。

### 1.1 Rust 的核心特点

- **内存安全**：无需垃圾回收器即可防止内存泄漏和空指针
- **零成本抽象**：高级特性不会带来运行时开销
- **并发安全**：类型系统防止数据竞争
- **跨平台**：支持多种操作系统和架构

### 1.2 安装 Rust

```bash
# 通过 rustup 安装（推荐）
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# 验证安装
rustc --version
cargo --version
```

### 1.3 第一个 Rust 程序

```rust
fn main() {
    println!("Hello, Rust!");
}
```

编译并运行：

```bash
rustc hello.rs
./hello
```

## 2. 基本语法

### 2.1 变量和常量

```rust
fn main() {
    // 不可变变量
    let x = 5;

    // 可变变量
    let mut y = 10;
    y = 15;

    // 常量
    const MAX_POINTS: u32 = 100_000;

    // 变量遮蔽
    let x = x + 1;
    let x = x * 2;
    println!("x的值是: {}", x); // 12
}
```

### 2.2 数据类型

#### 标量类型

```rust
fn main() {
    // 整数类型
    let a: i32 = -42;
    let b: u32 = 42;
    let c: isize = 100; // 架构相关

    // 浮点类型
    let d: f64 = 3.14;
    let e: f32 = 2.718;

    // 布尔类型
    let is_active: bool = true;

    // 字符类型
    let emoji: char = '😎';
    let chinese: char = '中';
}
```

#### 复合类型

```rust
fn main() {
    // 元组
    let tuple: (i32, f64, u8) = (500, 6.4, 1);
    let (x, y, z) = tuple; // 解构
    let first = tuple.0;   // 索引访问

    // 数组
    let array: [i32; 5] = [1, 2, 3, 4, 5];
    let same_value = [3; 5]; // [3, 3, 3, 3, 3]
    let first_element = array[0];
}
```

### 2.3 函数

```rust
fn main() {
    let result = add(5, 3);
    println!("结果: {}", result);

    let (sum, product) = calculate(4, 6);
    println!("和: {}, 积: {}", sum, product);
}

// 基本函数
fn add(a: i32, b: i32) -> i32 {
    a + b // 表达式，无分号
}

// 返回多个值
fn calculate(x: i32, y: i32) -> (i32, i32) {
    (x + y, x * y)
}

// 无返回值函数
fn print_message(msg: &str) {
    println!("消息: {}", msg);
}
```

### 2.4 控制流

#### if 表达式

```rust
fn main() {
    let number = 6;

    if number % 4 == 0 {
        println!("数字能被4整除");
    } else if number % 3 == 0 {
        println!("数字能被3整除");
    } else {
        println!("数字不能被4或3整除");
    }

    // if 作为表达式
    let condition = true;
    let number = if condition { 5 } else { 6 };
    println!("数字的值是: {}", number);
}
```

#### 循环

```rust
fn main() {
    // loop 循环
    let mut counter = 0;
    let result = loop {
        counter += 1;
        if counter == 10 {
            break counter * 2; // 返回值
        }
    };
    println!("结果: {}", result);

    // while 循环
    let mut number = 3;
    while number != 0 {
        println!("{}!", number);
        number -= 1;
    }

    // for 循环
    let a = [10, 20, 30, 40, 50];
    for element in a.iter() {
        println!("值: {}", element);
    }

    // 范围循环
    for number in 1..4 {
        println!("{}!", number);
    }
}
```

## 3. 注释和文档

````rust
// 单行注释

/*
 * 多行注释
 */

/// 文档注释 - 描述函数功能
///
/// # Examples
///
/// ```
/// let result = add_two(5);
/// assert_eq!(result, 7);
/// ```
fn add_two(x: i32) -> i32 {
    x + 2
}

//! 这是模块级别的文档注释
````

## 4. 练习

创建一个程序，实现以下功能：

1. 定义一个函数计算摄氏度到华氏度的转换
2. 使用循环打印斐波那契数列的前 10 个数字
3. 判断一个数字是否为质数

```rust
fn main() {
    // 温度转换
    let celsius = 25.0;
    let fahrenheit = celsius_to_fahrenheit(celsius);
    println!("{}°C = {}°F", celsius, fahrenheit);

    // 斐波那契数列
    println!("斐波那契数列前10个数:");
    fibonacci(10);

    // 质数检查
    let num = 17;
    if is_prime(num) {
        println!("{} 是质数", num);
    } else {
        println!("{} 不是质数", num);
    }
}

fn celsius_to_fahrenheit(celsius: f64) -> f64 {
    celsius * 9.0 / 5.0 + 32.0
}

fn fibonacci(n: usize) {
    let mut a = 0;
    let mut b = 1;

    for i in 0..n {
        if i == 0 {
            print!("{} ", a);
        } else if i == 1 {
            print!("{} ", b);
        } else {
            let next = a + b;
            print!("{} ", next);
            a = b;
            b = next;
        }
    }
    println!();
}

fn is_prime(n: u32) -> bool {
    if n < 2 {
        return false;
    }
    for i in 2..=(n as f64).sqrt() as u32 {
        if n % i == 0 {
            return false;
        }
    }
    true
}
```
