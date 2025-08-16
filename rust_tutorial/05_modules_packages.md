# Rust 模块系统与包管理

## 1. 模块系统概述

Rust 的模块系统包括：

- **包（Packages）**：Cargo 的特性，让你构建、测试和分享 crate
- **Crates**：一个模块的树形结构，形成了库或二进制项目
- **模块（Modules）**：让你控制作用域和路径的私有性
- **路径（Paths）**：一个命名项的方式

## 2. 包和 Crate

### 2.1 创建新包

```bash
# 创建新的二进制包
cargo new my_project

# 创建新的库包
cargo new my_library --lib

# 查看项目结构
cd my_project
tree
```

### 2.2 包的结构

```
my_project/
├── Cargo.toml
├── src/
│   ├── main.rs      # 二进制crate根
│   ├── lib.rs       # 库crate根
│   └── bin/
│       ├── named-executable.rs
│       ├── another-executable.rs
│       └── multi-file-executable/
│           ├── main.rs
│           └── some_module.rs
```

### 2.3 Cargo.toml 配置

```toml
[package]
name = "my_project"
version = "0.1.0"
edition = "2021"
authors = ["Your Name <you@example.com>"]
description = "A sample Rust project"
license = "MIT"

[dependencies]
serde = { version = "1.0", features = ["derive"] }
tokio = { version = "1", features = ["full"] }
clap = "4.0"

[dev-dependencies]
assert_cmd = "2.0"

[build-dependencies]
cc = "1.0"

[[bin]]
name = "my-executable"
path = "src/bin/my-executable.rs"
```

## 3. 模块定义和使用

### 3.1 基本模块定义

```rust
// src/lib.rs
mod front_of_house {
    mod hosting {
        fn add_to_waitlist() {}
        fn seat_at_table() {}
    }

    mod serving {
        fn take_order() {}
        fn serve_order() {}
        fn take_payment() {}
    }
}

mod back_of_house {
    fn fix_incorrect_order() {
        cook_order();
        super::serve_order(); // 使用 super 访问父模块
    }

    fn cook_order() {}
}

fn serve_order() {}
```

### 3.2 公有和私有

```rust
mod front_of_house {
    pub mod hosting {
        pub fn add_to_waitlist() {
            println!("添加到等待列表");
        }

        fn seat_at_table() {
            println!("安排就座");
        }
    }
}

// 公有函数
pub fn eat_at_restaurant() {
    // 绝对路径
    crate::front_of_house::hosting::add_to_waitlist();

    // 相对路径
    front_of_house::hosting::add_to_waitlist();
}

// 私有函数（默认）
fn private_function() {
    println!("这是一个私有函数");
}
```

### 3.3 结构体和枚举的公有性

```rust
mod back_of_house {
    pub struct Breakfast {
        pub toast: String,
        seasonal_fruit: String, // 私有字段
    }

    impl Breakfast {
        pub fn summer(toast: &str) -> Breakfast {
            Breakfast {
                toast: String::from(toast),
                seasonal_fruit: String::from("peaches"),
            }
        }
    }

    pub enum Appetizer {
        Soup,     // 公有成员
        Salad,    // 公有成员
    }
}

pub fn eat_at_restaurant() {
    let mut meal = back_of_house::Breakfast::summer("Rye");
    meal.toast = String::from("Wheat");
    println!("我要 {} 吐司", meal.toast);

    // meal.seasonal_fruit = String::from("blueberries"); // 编译错误！私有字段

    let order1 = back_of_house::Appetizer::Soup;
    let order2 = back_of_house::Appetizer::Salad;
}
```

## 4. use 关键字

### 4.1 基本 use 用法

```rust
mod front_of_house {
    pub mod hosting {
        pub fn add_to_waitlist() {}
    }
}

use crate::front_of_house::hosting;

pub fn eat_at_restaurant() {
    hosting::add_to_waitlist();
}
```

### 4.2 use 的惯用方式

```rust
use std::collections::HashMap;
use std::fmt::Result;
use std::io::Result as IoResult; // 使用 as 关键字提供新名称

fn function1() -> Result {
    // --snip--
    Ok(())
}

fn function2() -> IoResult<()> {
    // --snip--
    Ok(())
}

// 函数的惯用方式：引入函数的父模块
use crate::front_of_house::hosting;

pub fn eat_at_restaurant() {
    hosting::add_to_waitlist();
}

// 结构体、枚举和其他项的惯用方式：指定完整路径
use std::collections::HashMap;

fn main() {
    let mut map = HashMap::new();
    map.insert(1, 2);
}
```

### 4.3 pub use 重导出

```rust
mod front_of_house {
    pub mod hosting {
        pub fn add_to_waitlist() {}
    }
}

pub use crate::front_of_house::hosting;

pub fn eat_at_restaurant() {
    hosting::add_to_waitlist();
}
```

### 4.4 使用嵌套路径清理大量 use 语句

```rust
// 之前
use std::cmp::Ordering;
use std::io;

// 现在
use std::{cmp::Ordering, io};

// 之前
use std::io;
use std::io::Write;

// 现在
use std::io::{self, Write};

// 引入所有公有定义
use std::collections::*;
```

## 5. 模块分离到不同文件

### 5.1 将模块分离到文件

```rust
// src/lib.rs
mod front_of_house;

pub use crate::front_of_house::hosting;

pub fn eat_at_restaurant() {
    hosting::add_to_waitlist();
}
```

```rust
// src/front_of_house.rs
pub mod hosting {
    pub fn add_to_waitlist() {
        println!("添加到等待列表");
    }
}
```

### 5.2 进一步分离子模块

```rust
// src/front_of_house.rs
pub mod hosting;
```

```rust
// src/front_of_house/hosting.rs
pub fn add_to_waitlist() {
    println!("添加到等待列表");
}

pub fn seat_at_table() {
    println!("安排就座");
}
```

## 6. 使用外部包

### 6.1 添加依赖

```toml
# Cargo.toml
[dependencies]
rand = "0.8.5"
```

### 6.2 使用外部 crate

```rust
use rand::Rng;

fn main() {
    let secret_number = rand::thread_rng().gen_range(1..=100);
    println!("秘密数字是: {}", secret_number);
}
```

## 7. 工作空间 (Workspaces)

### 7.1 创建工作空间

```toml
# Cargo.toml (根目录)
[workspace]
members = [
    "adder",
    "add-one",
]
```

```
add/
├── Cargo.toml
├── Cargo.lock
├── adder/
│   ├── Cargo.toml
│   └── src/
│       └── main.rs
└── add-one/
    ├── Cargo.toml
    └── src/
        └── lib.rs
```

### 7.2 工作空间中的包

```toml
# add-one/Cargo.toml
[package]
name = "add-one"
version = "0.1.0"
edition = "2021"

[dependencies]
```

```rust
// add-one/src/lib.rs
pub fn add_one(x: i32) -> i32 {
    x + 1
}
```

```toml
# adder/Cargo.toml
[package]
name = "adder"
version = "0.1.0"
edition = "2021"

[dependencies]
add-one = { path = "../add-one" }
```

```rust
// adder/src/main.rs
use add_one;

fn main() {
    let num = 10;
    println!("{} 加一等于 {}", num, add_one::add_one(num));
}
```

## 8. 实践项目：图书馆管理系统

### 8.1 项目结构

```
library_management/
├── Cargo.toml
├── src/
│   ├── main.rs
│   ├── lib.rs
│   ├── models/
│   │   ├── mod.rs
│   │   ├── book.rs
│   │   ├── user.rs
│   │   └── transaction.rs
│   ├── services/
│   │   ├── mod.rs
│   │   ├── library_service.rs
│   │   └── user_service.rs
│   └── utils/
│       ├── mod.rs
│       └── date_utils.rs
```

### 8.2 实现代码

```toml
# Cargo.toml
[package]
name = "library_management"
version = "0.1.0"
edition = "2021"

[dependencies]
chrono = { version = "0.4", features = ["serde"] }
serde = { version = "1.0", features = ["derive"] }
uuid = { version = "1.0", features = ["v4"] }
```

```rust
// src/lib.rs
pub mod models;
pub mod services;
pub mod utils;

pub use models::{Book, User, Transaction};
pub use services::{LibraryService, UserService};
```

```rust
// src/models/mod.rs
pub mod book;
pub mod user;
pub mod transaction;

pub use book::Book;
pub use user::User;
pub use transaction::Transaction;
```

```rust
// src/models/book.rs
use serde::{Deserialize, Serialize};
use uuid::Uuid;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Book {
    pub id: Uuid,
    pub title: String,
    pub author: String,
    pub isbn: String,
    pub available: bool,
}

impl Book {
    pub fn new(title: String, author: String, isbn: String) -> Self {
        Self {
            id: Uuid::new_v4(),
            title,
            author,
            isbn,
            available: true,
        }
    }

    pub fn borrow_book(&mut self) -> Result<(), String> {
        if self.available {
            self.available = false;
            Ok(())
        } else {
            Err("书籍已被借出".to_string())
        }
    }

    pub fn return_book(&mut self) {
        self.available = true;
    }
}
```

```rust
// src/models/user.rs
use serde::{Deserialize, Serialize};
use uuid::Uuid;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct User {
    pub id: Uuid,
    pub name: String,
    pub email: String,
    pub borrowed_books: Vec<Uuid>,
}

impl User {
    pub fn new(name: String, email: String) -> Self {
        Self {
            id: Uuid::new_v4(),
            name,
            email,
            borrowed_books: Vec::new(),
        }
    }

    pub fn borrow_book(&mut self, book_id: Uuid) {
        self.borrowed_books.push(book_id);
    }

    pub fn return_book(&mut self, book_id: Uuid) {
        self.borrowed_books.retain(|&id| id != book_id);
    }
}
```

```rust
// src/models/transaction.rs
use serde::{Deserialize, Serialize};
use uuid::Uuid;
use chrono::{DateTime, Utc};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum TransactionType {
    Borrow,
    Return,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Transaction {
    pub id: Uuid,
    pub user_id: Uuid,
    pub book_id: Uuid,
    pub transaction_type: TransactionType,
    pub timestamp: DateTime<Utc>,
}

impl Transaction {
    pub fn new(user_id: Uuid, book_id: Uuid, transaction_type: TransactionType) -> Self {
        Self {
            id: Uuid::new_v4(),
            user_id,
            book_id,
            transaction_type,
            timestamp: Utc::now(),
        }
    }
}
```

```rust
// src/services/mod.rs
pub mod library_service;
pub mod user_service;

pub use library_service::LibraryService;
pub use user_service::UserService;
```

```rust
// src/services/library_service.rs
use crate::models::{Book, User, Transaction, TransactionType};
use std::collections::HashMap;
use uuid::Uuid;

pub struct LibraryService {
    books: HashMap<Uuid, Book>,
    users: HashMap<Uuid, User>,
    transactions: Vec<Transaction>,
}

impl LibraryService {
    pub fn new() -> Self {
        Self {
            books: HashMap::new(),
            users: HashMap::new(),
            transactions: Vec::new(),
        }
    }

    pub fn add_book(&mut self, book: Book) {
        self.books.insert(book.id, book);
    }

    pub fn add_user(&mut self, user: User) {
        self.users.insert(user.id, user);
    }

    pub fn borrow_book(&mut self, user_id: Uuid, book_id: Uuid) -> Result<(), String> {
        let user = self.users.get_mut(&user_id)
            .ok_or("用户不存在")?;
        let book = self.books.get_mut(&book_id)
            .ok_or("书籍不存在")?;

        book.borrow_book()?;
        user.borrow_book(book_id);

        let transaction = Transaction::new(user_id, book_id, TransactionType::Borrow);
        self.transactions.push(transaction);

        Ok(())
    }

    pub fn return_book(&mut self, user_id: Uuid, book_id: Uuid) -> Result<(), String> {
        let user = self.users.get_mut(&user_id)
            .ok_or("用户不存在")?;
        let book = self.books.get_mut(&book_id)
            .ok_or("书籍不存在")?;

        book.return_book();
        user.return_book(book_id);

        let transaction = Transaction::new(user_id, book_id, TransactionType::Return);
        self.transactions.push(transaction);

        Ok(())
    }

    pub fn list_available_books(&self) -> Vec<&Book> {
        self.books.values().filter(|book| book.available).collect()
    }

    pub fn get_user_books(&self, user_id: Uuid) -> Result<Vec<&Book>, String> {
        let user = self.users.get(&user_id)
            .ok_or("用户不存在")?;

        let books: Vec<&Book> = user.borrowed_books
            .iter()
            .filter_map(|book_id| self.books.get(book_id))
            .collect();

        Ok(books)
    }
}
```

```rust
// src/services/user_service.rs
use crate::models::User;
use std::collections::HashMap;
use uuid::Uuid;

pub struct UserService {
    users: HashMap<Uuid, User>,
}

impl UserService {
    pub fn new() -> Self {
        Self {
            users: HashMap::new(),
        }
    }

    pub fn create_user(&mut self, name: String, email: String) -> Uuid {
        let user = User::new(name, email);
        let user_id = user.id;
        self.users.insert(user_id, user);
        user_id
    }

    pub fn get_user(&self, user_id: Uuid) -> Option<&User> {
        self.users.get(&user_id)
    }

    pub fn update_user_email(&mut self, user_id: Uuid, new_email: String) -> Result<(), String> {
        let user = self.users.get_mut(&user_id)
            .ok_or("用户不存在")?;
        user.email = new_email;
        Ok(())
    }

    pub fn list_all_users(&self) -> Vec<&User> {
        self.users.values().collect()
    }
}
```

```rust
// src/utils/mod.rs
pub mod date_utils;

pub use date_utils::*;
```

```rust
// src/utils/date_utils.rs
use chrono::{DateTime, Utc, Local};

pub fn format_utc_time(time: DateTime<Utc>) -> String {
    time.format("%Y-%m-%d %H:%M:%S UTC").to_string()
}

pub fn format_local_time(time: DateTime<Utc>) -> String {
    let local_time: DateTime<Local> = time.into();
    local_time.format("%Y-%m-%d %H:%M:%S").to_string()
}

pub fn time_since(time: DateTime<Utc>) -> String {
    let now = Utc::now();
    let duration = now.signed_duration_since(time);

    if duration.num_days() > 0 {
        format!("{} 天前", duration.num_days())
    } else if duration.num_hours() > 0 {
        format!("{} 小时前", duration.num_hours())
    } else if duration.num_minutes() > 0 {
        format!("{} 分钟前", duration.num_minutes())
    } else {
        "刚刚".to_string()
    }
}
```

```rust
// src/main.rs
use library_management::{Book, User, LibraryService, UserService};

fn main() {
    let mut library = LibraryService::new();
    let mut user_service = UserService::new();

    // 创建用户
    let user_id = user_service.create_user(
        "张三".to_string(),
        "zhangsan@example.com".to_string()
    );

    // 添加书籍
    let book1 = Book::new(
        "Rust程序设计语言".to_string(),
        "Steve Klabnik".to_string(),
        "978-1718500440".to_string()
    );
    let book1_id = book1.id;
    library.add_book(book1);

    let book2 = Book::new(
        "计算机程序的构造和解释".to_string(),
        "Harold Abelson".to_string(),
        "978-0262510875".to_string()
    );
    library.add_book(book2);

    // 借书
    match library.borrow_book(user_id, book1_id) {
        Ok(_) => println!("借书成功!"),
        Err(e) => println!("借书失败: {}", e),
    }

    // 列出可借书籍
    println!("\n当前可借书籍:");
    for book in library.list_available_books() {
        println!("- {} (作者: {})", book.title, book.author);
    }

    // 列出用户借阅的书籍
    match library.get_user_books(user_id) {
        Ok(books) => {
            println!("\n用户借阅的书籍:");
            for book in books {
                println!("- {} (作者: {})", book.title, book.author);
            }
        },
        Err(e) => println!("获取用户书籍失败: {}", e),
    }
}
```

## 9. 高级模块概念

### 9.1 条件编译

```rust
// src/lib.rs
#[cfg(feature = "encryption")]
pub mod encryption;

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        assert_eq!(2 + 2, 4);
    }
}

#[cfg(target_os = "windows")]
fn windows_specific_function() {
    println!("这个函数只在 Windows 上编译");
}

#[cfg(target_os = "linux")]
fn linux_specific_function() {
    println!("这个函数只在 Linux 上编译");
}
```

### 9.2 特性标志

```toml
# Cargo.toml
[features]
default = ["std"]
std = []
serde_support = ["serde"]
encryption = ["aes"]

[dependencies]
serde = { version = "1.0", optional = true }
aes = { version = "0.8", optional = true }
```

模块系统和包管理是 Rust 项目组织的核心。通过合理地使用模块、包和工作空间，你可以构建出结构清晰、易于维护的大型项目。
