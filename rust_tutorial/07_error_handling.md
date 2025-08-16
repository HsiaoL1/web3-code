# Rust 错误处理

## 1. 错误处理概述

Rust 将错误分为两大类：

- **可恢复的错误**：使用 `Result<T, E>` 类型
- **不可恢复的错误**：使用 `panic!` 宏

Rust 没有异常系统，而是通过类型系统来处理错误，这使得错误处理更加明确和安全。

## 2. panic! 和不可恢复的错误

### 2.1 基本 panic!

```rust
fn main() {
    panic!("程序崩溃了！");
}
```

### 2.2 常见的 panic! 情况

```rust
fn main() {
    // 数组越界访问
    let v = vec![1, 2, 3];
    // v[99]; // 这会导致 panic!

    // 除零操作（对于整数）
    // let result = 10 / 0; // 这会导致 panic!

    // unwrap() 在 None 上调用
    let x: Option<i32> = None;
    // x.unwrap(); // 这会导致 panic!

    // expect() 提供自定义错误信息
    let x: Option<i32> = None;
    // x.expect("x 应该有一个值"); // 这会导致 panic! 并显示自定义消息
}
```

### 2.3 控制 panic! 行为

```bash
# 设置环境变量以获得更详细的回溯信息
RUST_BACKTRACE=1 cargo run
RUST_BACKTRACE=full cargo run
```

```rust
// 在程序中设置 panic hook
use std::panic;

fn main() {
    panic::set_hook(Box::new(|panic_info| {
        println!("自定义 panic 处理: {:?}", panic_info);
    }));

    panic!("测试自定义 panic 处理");
}
```

## 3. Result<T, E> 和可恢复的错误

### 3.1 Result 的定义

```rust
enum Result<T, E> {
    Ok(T),
    Err(E),
}
```

### 3.2 基本 Result 使用

```rust
use std::fs::File;

fn main() {
    let f = File::open("hello.txt");

    let f = match f {
        Ok(file) => file,
        Err(error) => {
            panic!("打开文件时出现问题: {:?}", error);
        }
    };
}
```

### 3.3 匹配不同的错误

```rust
use std::fs::File;
use std::io::ErrorKind;

fn main() {
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
}
```

### 3.4 unwrap 和 expect

```rust
use std::fs::File;

fn main() {
    // unwrap: 如果 Result 是 Ok，返回值；如果是 Err，调用 panic!
    // let f = File::open("hello.txt").unwrap();

    // expect: 类似 unwrap，但允许自定义错误消息
    let f = File::open("hello.txt").expect("无法打开 hello.txt");
}
```

## 4. 错误传播

### 4.1 手动传播错误

```rust
use std::fs::File;
use std::io::{self, Read};

fn read_username_from_file() -> Result<String, io::Error> {
    let f = File::open("hello.txt");

    let mut f = match f {
        Ok(file) => file,
        Err(e) => return Err(e),
    };

    let mut s = String::new();

    match f.read_to_string(&mut s) {
        Ok(_) => Ok(s),
        Err(e) => Err(e),
    }
}

fn main() {
    match read_username_from_file() {
        Ok(username) => println!("用户名: {}", username),
        Err(e) => println!("读取用户名时出错: {}", e),
    }
}
```

### 4.2 使用 ? 运算符

```rust
use std::fs::File;
use std::io::{self, Read};

fn read_username_from_file() -> Result<String, io::Error> {
    let mut f = File::open("hello.txt")?;
    let mut s = String::new();
    f.read_to_string(&mut s)?;
    Ok(s)
}

// 进一步简化
fn read_username_from_file_v2() -> Result<String, io::Error> {
    let mut s = String::new();
    File::open("hello.txt")?.read_to_string(&mut s)?;
    Ok(s)
}

// 最简化版本
fn read_username_from_file_v3() -> Result<String, io::Error> {
    std::fs::read_to_string("hello.txt")
}

fn main() {
    match read_username_from_file() {
        Ok(username) => println!("用户名: {}", username),
        Err(e) => println!("读取用户名时出错: {}", e),
    }
}
```

### 4.3 ? 在 main 函数中的使用

```rust
use std::error::Error;
use std::fs::File;

fn main() -> Result<(), Box<dyn Error>> {
    let f = File::open("hello.txt")?;
    Ok(())
}
```

### 4.4 ? 与 Option

```rust
fn last_char_of_first_line(text: &str) -> Option<char> {
    text.lines().next()?.chars().last()
}

fn main() {
    let text = "Hello\nWorld";
    match last_char_of_first_line(text) {
        Some(ch) => println!("最后一个字符: {}", ch),
        None => println!("没有找到字符"),
    }
}
```

## 5. 自定义错误类型

### 5.1 简单的自定义错误

```rust
#[derive(Debug)]
enum CalculatorError {
    DivisionByZero,
    NegativeSquareRoot,
    Overflow,
}

fn divide(a: f64, b: f64) -> Result<f64, CalculatorError> {
    if b == 0.0 {
        Err(CalculatorError::DivisionByZero)
    } else {
        Ok(a / b)
    }
}

fn square_root(x: f64) -> Result<f64, CalculatorError> {
    if x < 0.0 {
        Err(CalculatorError::NegativeSquareRoot)
    } else {
        Ok(x.sqrt())
    }
}

fn main() {
    match divide(10.0, 0.0) {
        Ok(result) => println!("结果: {}", result),
        Err(CalculatorError::DivisionByZero) => println!("错误: 除零"),
        Err(e) => println!("其他错误: {:?}", e),
    }

    match square_root(-4.0) {
        Ok(result) => println!("平方根: {}", result),
        Err(CalculatorError::NegativeSquareRoot) => println!("错误: 负数的平方根"),
        Err(e) => println!("其他错误: {:?}", e),
    }
}
```

### 5.2 实现 Error trait

```rust
use std::fmt;
use std::error::Error;

#[derive(Debug)]
enum MyError {
    Io(std::io::Error),
    Parse(std::num::ParseIntError),
    Custom(String),
}

impl fmt::Display for MyError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            MyError::Io(err) => write!(f, "IO 错误: {}", err),
            MyError::Parse(err) => write!(f, "解析错误: {}", err),
            MyError::Custom(msg) => write!(f, "自定义错误: {}", msg),
        }
    }
}

impl Error for MyError {
    fn source(&self) -> Option<&(dyn Error + 'static)> {
        match self {
            MyError::Io(err) => Some(err),
            MyError::Parse(err) => Some(err),
            MyError::Custom(_) => None,
        }
    }
}

impl From<std::io::Error> for MyError {
    fn from(error: std::io::Error) -> Self {
        MyError::Io(error)
    }
}

impl From<std::num::ParseIntError> for MyError {
    fn from(error: std::num::ParseIntError) -> Self {
        MyError::Parse(error)
    }
}

fn read_and_parse_file(filename: &str) -> Result<i32, MyError> {
    let content = std::fs::read_to_string(filename)?;
    let number = content.trim().parse::<i32>()?;
    Ok(number)
}

fn main() {
    match read_and_parse_file("number.txt") {
        Ok(number) => println!("读取的数字: {}", number),
        Err(e) => {
            println!("错误: {}", e);
            if let Some(source) = e.source() {
                println!("原因: {}", source);
            }
        }
    }
}
```

### 5.3 使用 thiserror crate

```toml
# Cargo.toml
[dependencies]
thiserror = "1.0"
```

```rust
use thiserror::Error;

#[derive(Error, Debug)]
pub enum DataStoreError {
    #[error("数据序列化失败")]
    Serialization(#[from] serde_json::Error),

    #[error("IO 操作失败")]
    Io(#[from] std::io::Error),

    #[error("数据验证失败: {message}")]
    Validation { message: String },

    #[error("未知的数据存储错误")]
    Unknown,
}

fn validate_data(data: &str) -> Result<(), DataStoreError> {
    if data.is_empty() {
        return Err(DataStoreError::Validation {
            message: "数据不能为空".to_string(),
        });
    }
    Ok(())
}

fn main() {
    match validate_data("") {
        Ok(_) => println!("数据验证成功"),
        Err(e) => println!("错误: {}", e),
    }
}
```

## 6. 错误处理的最佳实践

### 6.1 组合错误类型

```rust
use std::error::Error;
use std::fmt;

// 定义应用程序的顶级错误类型
#[derive(Debug)]
pub enum AppError {
    Database(DatabaseError),
    Network(NetworkError),
    Validation(ValidationError),
}

#[derive(Debug)]
pub struct DatabaseError {
    pub message: String,
}

#[derive(Debug)]
pub struct NetworkError {
    pub status_code: u16,
}

#[derive(Debug)]
pub struct ValidationError {
    pub field: String,
    pub message: String,
}

impl fmt::Display for AppError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            AppError::Database(err) => write!(f, "数据库错误: {}", err.message),
            AppError::Network(err) => write!(f, "网络错误: 状态码 {}", err.status_code),
            AppError::Validation(err) => write!(f, "验证错误: 字段 '{}' - {}", err.field, err.message),
        }
    }
}

impl Error for AppError {}

// 具体的业务函数
fn save_user(name: &str, email: &str) -> Result<(), AppError> {
    // 验证输入
    if name.is_empty() {
        return Err(AppError::Validation(ValidationError {
            field: "name".to_string(),
            message: "名称不能为空".to_string(),
        }));
    }

    if !email.contains('@') {
        return Err(AppError::Validation(ValidationError {
            field: "email".to_string(),
            message: "邮箱格式无效".to_string(),
        }));
    }

    // 模拟数据库操作
    if name == "forbidden" {
        return Err(AppError::Database(DatabaseError {
            message: "用户名被禁止".to_string(),
        }));
    }

    println!("用户 {} ({}) 保存成功", name, email);
    Ok(())
}

fn main() {
    let test_cases = vec![
        ("", "test@example.com"),
        ("John", "invalid-email"),
        ("forbidden", "test@example.com"),
        ("Alice", "alice@example.com"),
    ];

    for (name, email) in test_cases {
        match save_user(name, email) {
            Ok(_) => println!("✓ 用户保存成功"),
            Err(e) => println!("✗ {}", e),
        }
    }
}
```

### 6.2 Result 的链式操作

```rust
fn main() {
    let result = "42"
        .parse::<i32>()
        .map(|n| n * 2)
        .and_then(|n| {
            if n > 50 {
                Ok(n)
            } else {
                Err("数字太小".into())
            }
        })
        .map_err(|e| format!("处理失败: {}", e));

    match result {
        Ok(n) => println!("最终结果: {}", n),
        Err(e) => println!("错误: {}", e),
    }
}
```

### 6.3 使用 anyhow 进行快速错误处理

```toml
# Cargo.toml
[dependencies]
anyhow = "1.0"
```

```rust
use anyhow::{Result, Context};
use std::fs;

fn read_config() -> Result<String> {
    let content = fs::read_to_string("config.toml")
        .context("无法读取配置文件")?;

    if content.is_empty() {
        anyhow::bail!("配置文件为空");
    }

    Ok(content)
}

fn main() -> Result<()> {
    let config = read_config()
        .context("初始化配置失败")?;

    println!("配置: {}", config);
    Ok(())
}
```

## 7. 实践项目：文件处理器

```rust
use std::fs;
use std::io::{self, BufRead, BufReader, Write};
use std::path::Path;
use std::error::Error;
use std::fmt;

#[derive(Debug)]
enum FileProcessorError {
    Io(io::Error),
    Parse(String),
    Validation(String),
}

impl fmt::Display for FileProcessorError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            FileProcessorError::Io(err) => write!(f, "IO 错误: {}", err),
            FileProcessorError::Parse(msg) => write!(f, "解析错误: {}", msg),
            FileProcessorError::Validation(msg) => write!(f, "验证错误: {}", msg),
        }
    }
}

impl Error for FileProcessorError {}

impl From<io::Error> for FileProcessorError {
    fn from(error: io::Error) -> Self {
        FileProcessorError::Io(error)
    }
}

struct FileProcessor;

impl FileProcessor {
    fn new() -> Self {
        Self
    }

    fn validate_file_path(&self, path: &str) -> Result<(), FileProcessorError> {
        if path.is_empty() {
            return Err(FileProcessorError::Validation(
                "文件路径不能为空".to_string()
            ));
        }

        if !Path::new(path).exists() {
            return Err(FileProcessorError::Validation(
                format!("文件不存在: {}", path)
            ));
        }

        Ok(())
    }

    fn count_lines(&self, file_path: &str) -> Result<usize, FileProcessorError> {
        self.validate_file_path(file_path)?;

        let file = fs::File::open(file_path)?;
        let reader = BufReader::new(file);

        let count = reader.lines().count();
        Ok(count)
    }

    fn count_words(&self, file_path: &str) -> Result<usize, FileProcessorError> {
        self.validate_file_path(file_path)?;

        let content = fs::read_to_string(file_path)?;
        let word_count = content
            .split_whitespace()
            .count();

        Ok(word_count)
    }

    fn find_and_replace(&self,
                       file_path: &str,
                       find: &str,
                       replace: &str) -> Result<usize, FileProcessorError> {
        self.validate_file_path(file_path)?;

        if find.is_empty() {
            return Err(FileProcessorError::Validation(
                "搜索字符串不能为空".to_string()
            ));
        }

        let content = fs::read_to_string(file_path)?;
        let new_content = content.replace(find, replace);

        // 计算替换次数
        let replace_count = content.matches(find).count();

        // 写回文件
        fs::write(file_path, new_content)?;

        Ok(replace_count)
    }

    fn backup_file(&self, file_path: &str) -> Result<String, FileProcessorError> {
        self.validate_file_path(file_path)?;

        let backup_path = format!("{}.backup", file_path);
        fs::copy(file_path, &backup_path)?;

        Ok(backup_path)
    }

    fn process_csv(&self, file_path: &str) -> Result<Vec<Vec<String>>, FileProcessorError> {
        self.validate_file_path(file_path)?;

        let content = fs::read_to_string(file_path)?;
        let mut rows = Vec::new();

        for (line_num, line) in content.lines().enumerate() {
            if line.trim().is_empty() {
                continue;
            }

            let fields: Vec<String> = line
                .split(',')
                .map(|field| field.trim().to_string())
                .collect();

            if fields.is_empty() {
                return Err(FileProcessorError::Parse(
                    format!("第 {} 行为空或格式不正确", line_num + 1)
                ));
            }

            rows.push(fields);
        }

        if rows.is_empty() {
            return Err(FileProcessorError::Parse(
                "CSV 文件没有有效数据".to_string()
            ));
        }

        Ok(rows)
    }
}

// 创建测试文件的辅助函数
fn create_test_files() -> Result<(), Box<dyn Error>> {
    // 创建测试文本文件
    fs::write("test.txt", "Hello world\nThis is a test file\nWith multiple lines")?;

    // 创建测试 CSV 文件
    fs::write("test.csv", "Name,Age,City\nAlice,30,New York\nBob,25,Los Angeles\nCharlie,35,Chicago")?;

    Ok(())
}

fn main() -> Result<(), Box<dyn Error>> {
    // 创建测试文件
    create_test_files()?;

    let processor = FileProcessor::new();

    println!("=== 文件处理器演示 ===\n");

    // 测试行数统计
    match processor.count_lines("test.txt") {
        Ok(count) => println!("test.txt 有 {} 行", count),
        Err(e) => println!("统计行数失败: {}", e),
    }

    // 测试单词统计
    match processor.count_words("test.txt") {
        Ok(count) => println!("test.txt 有 {} 个单词", count),
        Err(e) => println!("统计单词失败: {}", e),
    }

    // 测试备份
    match processor.backup_file("test.txt") {
        Ok(backup_path) => println!("文件已备份到: {}", backup_path),
        Err(e) => println!("备份失败: {}", e),
    }

    // 测试查找替换
    match processor.find_and_replace("test.txt", "test", "example") {
        Ok(count) => println!("替换了 {} 处 'test' 为 'example'", count),
        Err(e) => println!("查找替换失败: {}", e),
    }

    // 测试 CSV 处理
    match processor.process_csv("test.csv") {
        Ok(rows) => {
            println!("\nCSV 数据:");
            for (i, row) in rows.iter().enumerate() {
                println!("行 {}: {:?}", i + 1, row);
            }
        },
        Err(e) => println!("处理 CSV 失败: {}", e),
    }

    // 测试错误情况
    println!("\n=== 错误处理测试 ===");

    match processor.count_lines("nonexistent.txt") {
        Ok(_) => println!("不应该到达这里"),
        Err(e) => println!("预期的错误: {}", e),
    }

    match processor.find_and_replace("test.txt", "", "replacement") {
        Ok(_) => println!("不应该到达这里"),
        Err(e) => println!("预期的错误: {}", e),
    }

    // 清理测试文件
    let _ = fs::remove_file("test.txt");
    let _ = fs::remove_file("test.txt.backup");
    let _ = fs::remove_file("test.csv");

    Ok(())
}
```

## 8. 错误处理总结

### 8.1 何时使用 panic! vs Result

- **使用 panic!**：

  - 程序遇到不可恢复的错误
  - 程序逻辑错误（如数组越界）
  - 测试中的断言失败

- **使用 Result**：
  - 可能失败的操作（如文件 IO、网络请求）
  - 用户输入验证
  - 外部依赖可能失败的情况

### 8.2 错误处理模式

```rust
// 1. 立即处理错误
fn handle_immediately() {
    match std::fs::read_to_string("file.txt") {
        Ok(content) => println!("内容: {}", content),
        Err(e) => println!("错误: {}", e),
    }
}

// 2. 传播错误
fn propagate_error() -> Result<String, std::io::Error> {
    let content = std::fs::read_to_string("file.txt")?;
    Ok(content.to_uppercase())
}

// 3. 转换错误类型
fn convert_error() -> Result<i32, Box<dyn std::error::Error>> {
    let content = std::fs::read_to_string("file.txt")?;
    let number = content.trim().parse::<i32>()?;
    Ok(number)
}

// 4. 提供默认值
fn with_default() -> String {
    std::fs::read_to_string("file.txt")
        .unwrap_or_else(|_| "默认内容".to_string())
}
```

错误处理是 Rust 编程的重要组成部分。通过掌握 `Result` 类型、`?` 操作符和自定义错误类型，你将能够编写出更加健壮和可维护的程序。
