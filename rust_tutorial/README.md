
# Rust 语言教程

本教程旨在提供对Rust语言的全面介绍，重点介绍其核心概念和语法。

## 什么是 Rust？

Rust 是一种多范式、通用的编程语言，专注于性能、类型安全和并发性。它在语法上类似于 C++，但旨在提供更好的内存安全性，同时保持高性能。Rust 由 Mozilla Research 开发，并于 2015 年首次发布。

## 核心特性

*   **内存安全**: Rust 的所有权系统和借用检查器在编译时强制执行内存安全保证。这可以防止出现空指针、悬垂指针和数据竞争等常见错误。
*   **性能**: Rust 速度快且内存效率高。由于没有运行时或垃圾收集器，它可以为性能关键型服务提供支持，在嵌入式设备上运行，并且可以轻松地与其他语言集成。
*   **并发**: Rust 的所有权和类型系统是其处理并发的强大工具。该语言鼓励使用消息传递来进行线程之间的通信，从而更容易编写正确的并发代码。
*   **可靠性**: Rust 丰富的类型系统和所有权模型保证了内存安全和线程安全，并且它还使你能够在编译时消除许多类别的错误。

## 语法基础

### 变量和可变性

在 Rust 中，变量默认是不可变的。你可以使用 `mut` 关键字使变量可变。

```rust
let x = 5; // 不可变
let mut y = 10; // 可变
y = 11;
```

### 数据类型

Rust 是一种静态类型语言，这意味着它必须在编译时知道所有变量的类型。一些基本数据类型包括：

*   **标量类型**:
    *   整数 (`i8`, `i16`, `i32`, `i64`, `i128`, `isize`, `u8`, `u16`, `u32`, `u64`, `u128`, `usize`)
    *   浮点数 (`f32`, `f64`)
    *   布尔值 (`bool`)
    *   字符 (`char`)
*   **复合类型**:
    *   元组 (`(i32, f64, u8)`)
    *   数组 (`[i32; 5]`)

### 函数

函数使用 `fn` 关键字声明。

```rust
fn add(a: i32, b: i32) -> i32 {
    a + b
}
```

### 控制流

*   **`if` 表达式**:

```rust
let number = 3;

if number < 5 {
    println!("条件为 true");
} else {
    println!("条件为 false");
}
```

*   **循环**: Rust 有三种循环：`loop`、`while` 和 `for`。

```rust
// loop
loop {
    println!("again!");
    break;
}

// while
let mut number = 3;
while number != 0 {
    println!("{}!", number);
    number -= 1;
}

// for
let a = [10, 20, 30, 40, 50];
for element in a.iter() {
    println!("值为：{}", element);
}
```

### 所有权

所有权是 Rust 最独特的特性。它使 Rust 能够在没有垃圾收集器的情况下保证内存安全。

*   **规则**:
    1.  Rust 中的每个值都有一个变量，称为其所有者。
    2.  一次只能有一个所有者。
    3.  当所有者超出作用域时，该值将被删除。

### 结构体

结构体是一种将多个相关值组合在一起的自定义数据类型。

```rust
struct User {
    username: String,
    email: String,
    sign_in_count: u64,
    active: bool,
}
```

### 枚举

枚举允许你通过列举其可能的成员来定义一个类型。

```rust
enum IpAddrKind {
    V4,
    V6,
}
```

### 错误处理

Rust 使用 `Result<T, E>` 枚举进行可恢复的错误处理。

```rust
use std::fs::File;

fn main() {
    let f = File::open("hello.txt");

    let f = match f {
        Ok(file) => file,
        Err(error) => {
            panic!("打开文件时出现问题：{:?}", error)
        },
    };
}
```

## 工具

*   **Cargo**: Rust 的构建工具和包管理器。
*   **Rustfmt**: 一个用于格式化 Rust 代码的工具。
*   **Clippy**: 一个用于捕捉常见错误和改进 Rust 代码的 linter。

## 结论

Rust 是一种功能强大的语言，为现代系统编程提供了独特的特性组合。其对安全性、性能和并发性的关注使其成为各种应用程序的绝佳选择。
