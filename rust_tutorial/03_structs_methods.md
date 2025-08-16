# Rust 结构体与方法

## 1. 结构体基础

结构体是一种自定义数据类型，允许你将多个相关的值组合在一起形成一个有意义的组合。

### 1.1 定义和实例化结构体

```rust
// 定义结构体
struct User {
    username: String,
    email: String,
    sign_in_count: u64,
    active: bool,
}

fn main() {
    // 创建结构体实例
    let user1 = User {
        email: String::from("someone@example.com"),
        username: String::from("someusername123"),
        active: true,
        sign_in_count: 1,
    };

    // 访问结构体字段
    println!("用户名: {}", user1.username);
    println!("邮箱: {}", user1.email);
}
```

### 1.2 可变结构体

```rust
fn main() {
    let mut user1 = User {
        email: String::from("someone@example.com"),
        username: String::from("someusername123"),
        active: true,
        sign_in_count: 1,
    };

    // 修改字段值
    user1.email = String::from("anotheremail@example.com");
    user1.sign_in_count += 1;

    println!("新邮箱: {}", user1.email);
    println!("登录次数: {}", user1.sign_in_count);
}
```

### 1.3 结构体更新语法

```rust
fn main() {
    let user1 = User {
        email: String::from("someone@example.com"),
        username: String::from("someusername123"),
        active: true,
        sign_in_count: 1,
    };

    // 使用结构体更新语法
    let user2 = User {
        email: String::from("another@example.com"),
        username: String::from("anotherusername567"),
        ..user1  // 其余字段从 user1 获取值
    };

    // 注意：user1 的 String 字段已被移动，不能再使用
    // println!("{}", user1.username); // 编译错误！
    println!("user2 邮箱: {}", user2.email);
    println!("user2 是否活跃: {}", user2.active);
}
```

### 1.4 使用函数创建结构体

```rust
fn build_user(email: String, username: String) -> User {
    User {
        email,    // 字段初始化简写语法
        username, // 等同于 username: username
        active: true,
        sign_in_count: 1,
    }
}

fn main() {
    let user = build_user(
        String::from("test@example.com"),
        String::from("testuser")
    );

    println!("创建的用户: {}", user.username);
}
```

## 2. 元组结构体

```rust
// 元组结构体
struct Color(i32, i32, i32);
struct Point(i32, i32, i32);

fn main() {
    let black = Color(0, 0, 0);
    let origin = Point(0, 0, 0);

    // 访问元组结构体的字段
    println!("黑色的红色分量: {}", black.0);
    println!("原点的 x 坐标: {}", origin.0);

    // 解构元组结构体
    let Color(r, g, b) = black;
    println!("RGB: ({}, {}, {})", r, g, b);
}
```

## 3. 类单元结构体

```rust
// 类单元结构体，没有任何字段
struct Unit;

fn main() {
    let unit = Unit;
    // 主要用于实现 trait，在泛型编程中有用
}
```

## 4. 结构体示例程序

```rust
#[derive(Debug)] // 添加 Debug trait 以便打印
struct Rectangle {
    width: u32,
    height: u32,
}

fn main() {
    let rect1 = Rectangle { width: 30, height: 50 };
    let rect2 = Rectangle { width: 10, height: 40 };
    let rect3 = Rectangle { width: 35, height: 55 };

    println!("矩形1: {:?}", rect1);
    println!("矩形1的面积是 {} 平方像素。", area(&rect1));

    println!("rect1 能容纳 rect2 吗？{}", rect1.can_hold(&rect2));
    println!("rect1 能容纳 rect3 吗？{}", rect1.can_hold(&rect3));
}

// 计算面积的函数
fn area(rectangle: &Rectangle) -> u32 {
    rectangle.width * rectangle.height
}

// 我们将在下一节中用方法重写这个函数
```

## 5. 方法语法

方法与函数类似，但它们在结构体的上下文中被定义。

### 5.1 定义方法

```rust
#[derive(Debug)]
struct Rectangle {
    width: u32,
    height: u32,
}

impl Rectangle {
    // 方法的第一个参数总是 self
    fn area(&self) -> u32 {
        self.width * self.height
    }

    fn can_hold(&self, other: &Rectangle) -> bool {
        self.width > other.width && self.height > other.height
    }

    // 获取所有权的方法（较少使用）
    fn max_area(self, other: Rectangle) -> Rectangle {
        if self.area() > other.area() {
            self
        } else {
            other
        }
    }

    // 可变借用方法
    fn double_size(&mut self) {
        self.width *= 2;
        self.height *= 2;
    }
}

fn main() {
    let rect1 = Rectangle { width: 30, height: 50 };
    let rect2 = Rectangle { width: 10, height: 40 };

    println!("矩形1的面积是 {} 平方像素。", rect1.area());
    println!("rect1 能容纳 rect2 吗？{}", rect1.can_hold(&rect2));

    let mut rect3 = Rectangle { width: 20, height: 30 };
    println!("rect3 双倍前: {:?}", rect3);
    rect3.double_size();
    println!("rect3 双倍后: {:?}", rect3);
}
```

### 5.2 关联函数

```rust
impl Rectangle {
    // 关联函数（不以 self 作为参数）
    fn square(size: u32) -> Rectangle {
        Rectangle { width: size, height: size }
    }

    fn from_dimensions(width: u32, height: u32) -> Rectangle {
        Rectangle { width, height }
    }
}

fn main() {
    // 使用 :: 语法调用关联函数
    let square = Rectangle::square(20);
    let rect = Rectangle::from_dimensions(30, 40);

    println!("正方形: {:?}", square);
    println!("矩形: {:?}", rect);
}
```

### 5.3 多个 impl 块

```rust
impl Rectangle {
    fn area(&self) -> u32 {
        self.width * self.height
    }
}

impl Rectangle {
    fn perimeter(&self) -> u32 {
        2 * (self.width + self.height)
    }
}

fn main() {
    let rect = Rectangle { width: 30, height: 50 };
    println!("面积: {}", rect.area());
    println!("周长: {}", rect.perimeter());
}
```

## 6. 高级结构体概念

### 6.1 结构体字段的生命周期

```rust
// 这个结构体不能编译，因为引用需要生命周期参数
// struct User {
//     username: &str,  // 需要生命周期
//     email: &str,     // 需要生命周期
//     active: bool,
// }

// 正确的方式（我们将在生命周期章节详细讨论）
struct User<'a> {
    username: &'a str,
    email: &'a str,
    active: bool,
}
```

### 6.2 结构体的 Debug trait

```rust
// 手动实现 Debug
struct Point {
    x: i32,
    y: i32,
}

impl std::fmt::Debug for Point {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        f.debug_struct("Point")
            .field("x", &self.x)
            .field("y", &self.y)
            .finish()
    }
}

// 使用 derive 宏自动实现（推荐）
#[derive(Debug, Clone, PartialEq)]
struct Rectangle {
    width: u32,
    height: u32,
}

fn main() {
    let point = Point { x: 1, y: 2 };
    let rect = Rectangle { width: 30, height: 50 };

    println!("点: {:?}", point);
    println!("矩形: {:#?}", rect); // 美化输出
}
```

## 7. 实践项目：图书管理系统

```rust
#[derive(Debug, Clone)]
struct Book {
    title: String,
    author: String,
    pages: u32,
    available: bool,
}

#[derive(Debug)]
struct Library {
    name: String,
    books: Vec<Book>,
}

impl Book {
    fn new(title: String, author: String, pages: u32) -> Book {
        Book {
            title,
            author,
            pages,
            available: true,
        }
    }

    fn borrow_book(&mut self) -> Result<(), String> {
        if self.available {
            self.available = false;
            Ok(())
        } else {
            Err("书籍已被借出".to_string())
        }
    }

    fn return_book(&mut self) {
        self.available = true;
    }

    fn info(&self) -> String {
        format!("《{}》 - {} ({} 页) - {}",
            self.title,
            self.author,
            self.pages,
            if self.available { "可借" } else { "已借出" }
        )
    }
}

impl Library {
    fn new(name: String) -> Library {
        Library {
            name,
            books: Vec::new(),
        }
    }

    fn add_book(&mut self, book: Book) {
        self.books.push(book);
    }

    fn find_book(&mut self, title: &str) -> Option<&mut Book> {
        self.books.iter_mut().find(|book| book.title == title)
    }

    fn list_available_books(&self) -> Vec<&Book> {
        self.books.iter().filter(|book| book.available).collect()
    }

    fn borrow_book(&mut self, title: &str) -> Result<(), String> {
        match self.find_book(title) {
            Some(book) => book.borrow_book(),
            None => Err("未找到该书籍".to_string()),
        }
    }

    fn return_book(&mut self, title: &str) -> Result<(), String> {
        match self.find_book(title) {
            Some(book) => {
                book.return_book();
                Ok(())
            },
            None => Err("未找到该书籍".to_string()),
        }
    }
}

fn main() {
    let mut library = Library::new("中央图书馆".to_string());

    // 添加书籍
    library.add_book(Book::new("Rust程序设计".to_string(), "Steve Klabnik".to_string(), 500));
    library.add_book(Book::new("计算机程序的构造和解释".to_string(), "Harold Abelson".to_string(), 600));
    library.add_book(Book::new("算法导论".to_string(), "Thomas H. Cormen".to_string(), 1200));

    println!("=== {} ===", library.name);

    // 列出可借书籍
    println!("\n可借书籍:");
    for book in library.list_available_books() {
        println!("  {}", book.info());
    }

    // 借书
    println!("\n尝试借阅《Rust程序设计》...");
    match library.borrow_book("Rust程序设计") {
        Ok(_) => println!("借阅成功!"),
        Err(e) => println!("借阅失败: {}", e),
    }

    // 再次列出可借书籍
    println!("\n当前可借书籍:");
    for book in library.list_available_books() {
        println!("  {}", book.info());
    }

    // 归还书籍
    println!("\n归还《Rust程序设计》...");
    match library.return_book("Rust程序设计") {
        Ok(_) => println!("归还成功!"),
        Err(e) => println!("归还失败: {}", e),
    }

    // 最终状态
    println!("\n最终可借书籍:");
    for book in library.list_available_books() {
        println!("  {}", book.info());
    }
}
```

## 8. 练习题

1. **银行账户系统**：创建一个 `BankAccount` 结构体，包含账户号、余额和账户类型。实现存款、取款和查询余额的方法。

2. **学生成绩管理**：创建 `Student` 和 `Grade` 结构体，实现添加成绩、计算平均分和显示成绩单的功能。

3. **几何形状**：创建不同的几何形状结构体（圆形、三角形、正方形），为每个形状实现计算面积和周长的方法。

结构体和方法是 Rust 中组织代码的重要方式，它们帮助你创建更加模块化和可维护的程序。通过练习这些概念，你将能够设计出更好的数据结构和 API。
