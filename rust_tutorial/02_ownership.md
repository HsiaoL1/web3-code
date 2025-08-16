# Rust 所有权系统

## 1. 所有权概念

所有权是 Rust 最独特和重要的特性，它使 Rust 能够在没有垃圾回收器的情况下保证内存安全。

### 1.1 所有权规则

1. Rust 中每个值都有一个变量，这个变量被称为值的**所有者**
2. 一个值在任何时刻只能有一个所有者
3. 当所有者离开作用域时，这个值将被丢弃

### 1.2 变量作用域

```rust
fn main() {
    {                      // s 在这里无效，它尚未声明
        let s = "hello";   // 从此处起，s 是有效的

        // 使用 s
        println!("{}", s);
    }                      // 此作用域已结束，s 不再有效
}
```

### 1.3 String 类型

```rust
fn main() {
    // 字符串字面值（不可变）
    let s1 = "hello";

    // String 类型（可变）
    let mut s2 = String::from("hello");
    s2.push_str(", world!");
    println!("{}", s2);
}
```

## 2. 内存与分配

### 2.1 栈与堆

```rust
fn main() {
    // 栈上的数据 - 固定大小
    let x = 5;
    let y = x; // 复制

    // 堆上的数据 - 动态大小
    let s1 = String::from("hello");
    let s2 = s1; // 移动（move），s1 不再有效

    // println!("{}", s1); // 编译错误！
    println!("{}", s2); // 正确
}
```

### 2.2 克隆数据

```rust
fn main() {
    let s1 = String::from("hello");
    let s2 = s1.clone(); // 深拷贝

    println!("s1 = {}, s2 = {}", s1, s2); // 两个都有效
}
```

### 2.3 Copy trait

```rust
fn main() {
    // 实现了 Copy trait 的类型可以直接复制
    let x = 5;
    let y = x;
    println!("x = {}, y = {}", x, y); // 都有效

    // 以下类型实现了 Copy：
    // - 所有整数类型
    // - 布尔类型
    // - 浮点类型
    // - 字符类型
    // - 元组（如果所有元素都是 Copy）
}
```

## 3. 所有权与函数

### 3.1 传递值给函数

```rust
fn main() {
    let s = String::from("hello");  // s 进入作用域

    takes_ownership(s);             // s 的值移动到函数里
                                    // s 到这里不再有效

    let x = 5;                      // x 进入作用域

    makes_copy(x);                  // x 应该移动函数里，
                                    // 但 i32 是 Copy 的，
                                    // 所以在后面可继续使用 x

} // 这里, x 先移出了作用域，然后是 s。但因为 s 的值已被移走，
  // 没有特殊之处

fn takes_ownership(some_string: String) { // some_string 进入作用域
    println!("{}", some_string);
} // 这里，some_string 移出作用域并调用 `drop` 方法。占用的内存被释放

fn makes_copy(some_integer: i32) { // some_integer 进入作用域
    println!("{}", some_integer);
} // 这里，some_integer 移出作用域。不会有特殊操作
```

### 3.2 返回值与作用域

```rust
fn main() {
    let s1 = gives_ownership();         // gives_ownership 将返回值
                                        // 移给 s1

    let s2 = String::from("hello");     // s2 进入作用域

    let s3 = takes_and_gives_back(s2);  // s2 被移动到
                                        // takes_and_gives_back 中,
                                        // 它也将返回值移给 s3
} // 这里, s3 移出作用域并被丢弃。s2 也移出作用域，但已被移走，
  // 所以什么也不会发生。s1 移出作用域并被丢弃

fn gives_ownership() -> String {             // gives_ownership 将返回值移动给
                                             // 调用它的函数

    let some_string = String::from("yours"); // some_string 进入作用域

    some_string                              // 返回 some_string 并移出给调用的函数
}

// takes_and_gives_back 将传入字符串并返回该值
fn takes_and_gives_back(a_string: String) -> String { // a_string 进入作用域

    a_string  // 返回 a_string 并移出给调用的函数
}
```

### 3.3 返回多个值

```rust
fn main() {
    let s1 = String::from("hello");

    let (s2, len) = calculate_length(s1);

    println!("'{}' 的长度是 {}。", s2, len);
}

fn calculate_length(s: String) -> (String, usize) {
    let length = s.len(); // len() 返回字符串的长度

    (s, length)
}
```

## 4. 引用与借用

### 4.1 引用

引用允许你使用值但不获取其所有权：

```rust
fn main() {
    let s1 = String::from("hello");

    let len = calculate_length(&s1); // &s1 创建一个指向 s1 值的引用

    println!("'{}' 的长度是 {}。", s1, len);
}

fn calculate_length(s: &String) -> usize { // s 是对 String 的引用
    s.len()
} // 这里，s 离开了作用域。但因为它并不拥有引用值的所有权，
  // 所以什么也不会发生
```

### 4.2 可变引用

```rust
fn main() {
    let mut s = String::from("hello");

    change(&mut s);

    println!("{}", s);
}

fn change(some_string: &mut String) {
    some_string.push_str(", world");
}
```

### 4.3 引用的规则

```rust
fn main() {
    let mut s = String::from("hello");

    // 规则1：在任意给定时间，要么只能有一个可变引用，要么只能有多个不可变引用
    let r1 = &s; // 没问题
    let r2 = &s; // 没问题
    println!("{} and {}", r1, r2);
    // 此位置之后 r1 和 r2 不再使用

    let r3 = &mut s; // 没问题
    println!("{}", r3);

    // 规则2：引用必须总是有效的（防止悬垂引用）
}
```

### 4.4 悬垂引用

```rust
fn main() {
    let reference_to_nothing = dangle();
}

// 这段代码不会编译！
// fn dangle() -> &String { // dangle 返回一个字符串的引用
//     let s = String::from("hello"); // s 是一个新字符串
//     &s // 返回字符串 s 的引用
// } // 这里 s 离开作用域并被丢弃。其内存被释放。危险！

// 正确的做法
fn no_dangle() -> String {
    let s = String::from("hello");
    s // 直接返回字符串
}
```

## 5. 切片类型

切片让你引用集合中一段连续的元素序列，而不用引用整个集合。

### 5.1 字符串切片

```rust
fn main() {
    let s = String::from("hello world");

    let hello = &s[0..5];  // 或 &s[..5]
    let world = &s[6..11]; // 或 &s[6..]
    let whole = &s[..];    // 整个字符串

    println!("{} {} {}", hello, world, whole);

    let word = first_word(&s);
    println!("第一个单词是: {}", word);
}

fn first_word(s: &String) -> &str {
    let bytes = s.as_bytes();

    for (i, &item) in bytes.iter().enumerate() {
        if item == b' ' {
            return &s[0..i];
        }
    }

    &s[..]
}
```

### 5.2 其他切片

```rust
fn main() {
    let a = [1, 2, 3, 4, 5];
    let slice = &a[1..3]; // 类型是 &[i32]

    for element in slice {
        println!("{}", element);
    }
}
```

## 6. 实践练习

### 练习 1：实现一个简单的单词计数器

```rust
fn main() {
    let text = "hello world hello rust";
    let word_count = count_words(text);
    println!("单词数量: {}", word_count);

    let first = get_first_word(text);
    println!("第一个单词: {}", first);
}

fn count_words(s: &str) -> usize {
    if s.trim().is_empty() {
        return 0;
    }

    s.split_whitespace().count()
}

fn get_first_word(s: &str) -> &str {
    let bytes = s.as_bytes();

    for (i, &item) in bytes.iter().enumerate() {
        if item == b' ' {
            return &s[..i];
        }
    }

    s
}
```

### 练习 2：实现字符串反转

```rust
fn main() {
    let s = String::from("hello");
    let reversed = reverse_string(&s);
    println!("原字符串: {}", s);
    println!("反转后: {}", reversed);
}

fn reverse_string(s: &str) -> String {
    s.chars().rev().collect()
}
```

所有权系统是 Rust 的核心特性，理解它对于编写安全高效的 Rust 代码至关重要。通过练习这些例子，你将更好地掌握所有权、借用和切片的概念。
