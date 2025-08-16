# Rust 泛型、Trait 和生命周期

## 1. 泛型 (Generics)

泛型允许我们编写可以处理多种类型的代码，同时保持类型安全。

### 1.1 函数中的泛型

```rust
// 不使用泛型的重复代码
fn largest_i32(list: &[i32]) -> i32 {
    let mut largest = list[0];
    for &item in list {
        if item > largest {
            largest = item;
        }
    }
    largest
}

fn largest_char(list: &[char]) -> char {
    let mut largest = list[0];
    for &item in list {
        if item > largest {
            largest = item;
        }
    }
    largest
}

// 使用泛型
fn largest<T: PartialOrd + Copy>(list: &[T]) -> T {
    let mut largest = list[0];
    for &item in list {
        if item > largest {
            largest = item;
        }
    }
    largest
}

fn main() {
    let number_list = vec![34, 50, 25, 100, 65];
    let result = largest(&number_list);
    println!("最大的数字是 {}", result);

    let char_list = vec!['y', 'm', 'a', 'q'];
    let result = largest(&char_list);
    println!("最大的字符是 {}", result);
}
```

### 1.2 结构体中的泛型

```rust
struct Point<T> {
    x: T,
    y: T,
}

// 多个泛型参数
struct PointMixed<T, U> {
    x: T,
    y: U,
}

impl<T> Point<T> {
    fn x(&self) -> &T {
        &self.x
    }
}

// 为特定类型实现方法
impl Point<f32> {
    fn distance_from_origin(&self) -> f32 {
        (self.x.powi(2) + self.y.powi(2)).sqrt()
    }
}

// 方法中的泛型
impl<T, U> PointMixed<T, U> {
    fn mixup<V, W>(self, other: PointMixed<V, W>) -> PointMixed<T, W> {
        PointMixed {
            x: self.x,
            y: other.y,
        }
    }
}

fn main() {
    let integer = Point { x: 5, y: 10 };
    let float = Point { x: 1.0, y: 4.0 };
    let mixed = PointMixed { x: 5, y: 4.0 };

    println!("integer.x = {}", integer.x());
    println!("float distance = {}", float.distance_from_origin());

    let p1 = PointMixed { x: 5, y: 10.4 };
    let p2 = PointMixed { x: "Hello", y: 'c' };
    let p3 = p1.mixup(p2);
    println!("p3.x = {}, p3.y = {}", p3.x, p3.y);
}
```

### 1.3 枚举中的泛型

```rust
// 标准库中的例子
enum Option<T> {
    Some(T),
    None,
}

enum Result<T, E> {
    Ok(T),
    Err(E),
}

// 自定义泛型枚举
enum Container<T> {
    Empty,
    Single(T),
    Multiple(Vec<T>),
}

impl<T> Container<T> {
    fn new() -> Self {
        Container::Empty
    }

    fn add(self, item: T) -> Self {
        match self {
            Container::Empty => Container::Single(item),
            Container::Single(existing) => Container::Multiple(vec![existing, item]),
            Container::Multiple(mut vec) => {
                vec.push(item);
                Container::Multiple(vec)
            }
        }
    }

    fn len(&self) -> usize {
        match self {
            Container::Empty => 0,
            Container::Single(_) => 1,
            Container::Multiple(vec) => vec.len(),
        }
    }
}

fn main() {
    let container = Container::new()
        .add("first")
        .add("second")
        .add("third");

    println!("容器大小: {}", container.len());
}
```

## 2. Trait

Trait 定义了某个特定类型拥有可能与其他类型共享的功能。

### 2.1 定义 Trait

```rust
pub trait Summary {
    fn summarize(&self) -> String;

    // 默认实现
    fn summarize_author(&self) -> String {
        String::from("(作者未知)")
    }

    // 使用默认实现的方法
    fn summarize_with_author(&self) -> String {
        format!("(文章作者: {})", self.summarize_author())
    }
}

pub struct NewsArticle {
    pub headline: String,
    pub location: String,
    pub author: String,
    pub content: String,
}

impl Summary for NewsArticle {
    fn summarize(&self) -> String {
        format!("{}, by {} ({})", self.headline, self.author, self.location)
    }

    fn summarize_author(&self) -> String {
        format!("{}", self.author)
    }
}

pub struct Tweet {
    pub username: String,
    pub content: String,
    pub reply: bool,
    pub retweet: bool,
}

impl Summary for Tweet {
    fn summarize(&self) -> String {
        format!("{}: {}", self.username, self.content)
    }

    fn summarize_author(&self) -> String {
        format!("@{}", self.username)
    }
}

fn main() {
    let tweet = Tweet {
        username: String::from("horse_ebooks"),
        content: String::from("当然，正如你们大多数人所知道的，人都是马"),
        reply: false,
        retweet: false,
    };

    println!("1 条新推文: {}", tweet.summarize());
    println!("作者信息: {}", tweet.summarize_with_author());

    let article = NewsArticle {
        headline: String::from("企鹅队再次获得斯坦利杯冠军！"),
        location: String::from("匹兹堡, PA, USA"),
        author: String::from("Iceburgh"),
        content: String::from("匹兹堡企鹅队再次成为了最好的曲棍球队"),
    };

    println!("新文章发布: {}", article.summarize());
}
```

### 2.2 Trait 作为参数

```rust
// Trait 绑定语法
pub fn notify(item: &impl Summary) {
    println!("突发新闻! {}", item.summarize());
}

// 更复杂的 Trait 绑定
pub fn notify_verbose<T: Summary>(item: &T) {
    println!("突发新闻! {}", item.summarize());
}

// 多个 Trait 绑定
pub fn notify_multiple(item: &(impl Summary + std::fmt::Display)) {
    println!("突发新闻! {}", item.summarize());
}

// 使用 where 子句
pub fn some_function<T, U>(t: &T, u: &U) -> i32
where
    T: std::fmt::Display + Clone,
    U: Clone + std::fmt::Debug,
{
    42
}

// 返回实现了 Trait 的类型
fn returns_summarizable() -> impl Summary {
    Tweet {
        username: String::from("horse_ebooks"),
        content: String::from("当然，正如你们大多数人所知道的，人都是马"),
        reply: false,
        retweet: false,
    }
}

fn main() {
    let tweet = Tweet {
        username: String::from("user"),
        content: String::from("Hello, world!"),
        reply: false,
        retweet: false,
    };

    notify(&tweet);

    let returned_tweet = returns_summarizable();
    println!("返回的推文: {}", returned_tweet.summarize());
}
```

### 2.3 常用的标准库 Trait

```rust
use std::fmt;

// Debug Trait
#[derive(Debug)]
struct Point {
    x: i32,
    y: i32,
}

// Clone Trait
#[derive(Clone, Debug)]
struct Person {
    name: String,
    age: u32,
}

// PartialEq 和 Eq Trait
#[derive(PartialEq, Eq, Debug)]
struct Book {
    title: String,
    author: String,
}

// Display Trait
impl fmt::Display for Point {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "({}, {})", self.x, self.y)
    }
}

// From 和 Into Trait
struct Circle {
    radius: f64,
}

impl From<f64> for Circle {
    fn from(radius: f64) -> Self {
        Circle { radius }
    }
}

// Iterator Trait
struct Counter {
    current: usize,
    max: usize,
}

impl Counter {
    fn new(max: usize) -> Counter {
        Counter { current: 0, max }
    }
}

impl Iterator for Counter {
    type Item = usize;

    fn next(&mut self) -> Option<Self::Item> {
        if self.current < self.max {
            let current = self.current;
            self.current += 1;
            Some(current)
        } else {
            None
        }
    }
}

fn main() {
    let p1 = Point { x: 1, y: 2 };
    let p2 = Point { x: 3, y: 4 };

    println!("Debug: {:?}", p1);
    println!("Display: {}", p1);

    let person1 = Person {
        name: "Alice".to_string(),
        age: 30,
    };
    let person2 = person1.clone();
    println!("原始: {:?}", person1);
    println!("克隆: {:?}", person2);

    let book1 = Book {
        title: "The Rust Book".to_string(),
        author: "Steve Klabnik".to_string(),
    };
    let book2 = Book {
        title: "The Rust Book".to_string(),
        author: "Steve Klabnik".to_string(),
    };
    println!("书籍相等: {}", book1 == book2);

    // From/Into 示例
    let circle: Circle = 5.0.into();
    let circle2 = Circle::from(3.0);
    println!("圆的半径: {}, {}", circle.radius, circle2.radius);

    // Iterator 示例
    let counter = Counter::new(3);
    for num in counter {
        println!("计数: {}", num);
    }
}
```

### 2.4 高级 Trait 用法

```rust
// 关联类型
pub trait Iterator {
    type Item;
    fn next(&mut self) -> Option<Self::Item>;
}

// 泛型参数 vs 关联类型
trait Add<Rhs = Self> {
    type Output;
    fn add(self, rhs: Rhs) -> Self::Output;
}

// Trait 对象
trait Draw {
    fn draw(&self);
}

struct Circle {
    radius: f64,
}

struct Rectangle {
    width: f64,
    height: f64,
}

impl Draw for Circle {
    fn draw(&self) {
        println!("绘制圆形，半径: {}", self.radius);
    }
}

impl Draw for Rectangle {
    fn draw(&self) {
        println!("绘制矩形，宽: {}, 高: {}", self.width, self.height);
    }
}

// 使用 Trait 对象
fn draw_shapes(shapes: &[Box<dyn Draw>]) {
    for shape in shapes {
        shape.draw();
    }
}

// 超级 Trait
trait Pilot {
    fn fly(&self);
}

trait Wizard {
    fn fly(&self);
}

struct Human;

impl Pilot for Human {
    fn fly(&self) {
        println!("像飞行员一样飞行");
    }
}

impl Wizard for Human {
    fn fly(&self) {
        println!("像巫师一样飞行");
    }
}

impl Human {
    fn fly(&self) {
        println!("*疯狂挥动手臂*");
    }
}

fn main() {
    let shapes: Vec<Box<dyn Draw>> = vec![
        Box::new(Circle { radius: 5.0 }),
        Box::new(Rectangle { width: 3.0, height: 4.0 }),
    ];

    draw_shapes(&shapes);

    // 消除歧义
    let person = Human;
    person.fly(); // 调用 Human 的方法
    Pilot::fly(&person); // 调用 Pilot trait 的方法
    Wizard::fly(&person); // 调用 Wizard trait 的方法
}
```

## 3. 生命周期 (Lifetimes)

生命周期确保引用在需要它们的时候保持有效。

### 3.1 生命周期基础

```rust
// 悬垂引用问题（编译错误）
/*
fn main() {
    let r;
    {
        let x = 5;
        r = &x; // 错误：x 的生命周期不够长
    }
    println!("r: {}", r);
}
*/

// 正确的版本
fn main() {
    let x = 5;
    let r = &x;
    println!("r: {}", r);
}
```

### 3.2 函数中的生命周期

```rust
// 需要生命周期注解
fn longest<'a>(x: &'a str, y: &'a str) -> &'a str {
    if x.len() > y.len() {
        x
    } else {
        y
    }
}

// 生命周期注解不改变生命周期的长短
fn first_word<'a>(s: &'a str) -> &'a str {
    let bytes = s.as_bytes();

    for (i, &item) in bytes.iter().enumerate() {
        if item == b' ' {
            return &s[0..i];
        }
    }

    &s[..]
}

fn main() {
    let string1 = String::from("long string is long");

    {
        let string2 = String::from("xyz");
        let result = longest(string1.as_str(), string2.as_str());
        println!("最长的字符串是 {}", result);
    }

    let sentence = "Hello world wonderful day";
    let word = first_word(sentence);
    println!("第一个单词: {}", word);
}
```

### 3.3 结构体中的生命周期

```rust
// 结构体包含引用时需要生命周期注解
struct ImportantExcerpt<'a> {
    part: &'a str,
}

impl<'a> ImportantExcerpt<'a> {
    fn level(&self) -> i32 {
        3
    }

    // 生命周期省略规则
    fn announce_and_return_part(&self, announcement: &str) -> &str {
        println!("请注意: {}", announcement);
        self.part
    }
}

// 包含多个引用的结构体
struct Article<'a, 'b> {
    title: &'a str,
    content: &'b str,
}

fn main() {
    let novel = String::from("写小说。很久很久以前...");
    let first_sentence = novel.split('.').next().expect("找不到句号");
    let i = ImportantExcerpt {
        part: first_sentence,
    };

    println!("重要摘录: {}", i.part);
    println!("级别: {}", i.level());

    let article = Article {
        title: "Rust 生命周期",
        content: "生命周期是 Rust 的一个重要概念...",
    };

    println!("文章标题: {}", article.title);
}
```

### 3.4 生命周期省略规则

```rust
// 编译器自动推断生命周期的情况

// 规则1: 每个引用参数都有自己的生命周期
fn first_word(s: &str) -> &str { // 等价于 fn first_word<'a>(s: &'a str) -> &'a str
    let bytes = s.as_bytes();
    for (i, &item) in bytes.iter().enumerate() {
        if item == b' ' {
            return &s[0..i];
        }
    }
    &s[..]
}

// 规则2: 如果只有一个输入生命周期参数，它被赋给所有输出生命周期参数
fn get_slice(s: &str) -> &str { // 等价于 fn get_slice<'a>(s: &'a str) -> &'a str
    &s[1..3]
}

// 规则3: 如果方法有多个输入生命周期参数且其中一个是 &self 或 &mut self，
// 那么所有输出生命周期参数被赋予 self 的生命周期
impl<'a> ImportantExcerpt<'a> {
    fn return_part(&self, announcement: &str) -> &str {
        println!("请注意: {}", announcement);
        self.part
    }
}

fn main() {
    let s = "Hello world";
    let word = first_word(s);
    println!("第一个单词: {}", word);

    let slice = get_slice(s);
    println!("切片: {}", slice);
}
```

### 3.5 静态生命周期

```rust
// 'static 生命周期表示能够存活于整个程序期间
let s: &'static str = "我有静态生命周期";

// 字符串字面值默认具有 'static 生命周期
fn get_static_str() -> &'static str {
    "这个字符串存储在程序的二进制文件中"
}

// 需要静态生命周期的函数
fn takes_static(s: &'static str) {
    println!("接收到静态字符串: {}", s);
}

fn main() {
    let static_str = get_static_str();
    takes_static(static_str);
    takes_static("另一个静态字符串");
}
```

## 4. 综合示例：缓存系统

```rust
use std::collections::HashMap;
use std::hash::Hash;

// 定义缓存 trait
trait Cache<K, V> {
    fn get(&self, key: &K) -> Option<&V>;
    fn set(&mut self, key: K, value: V);
    fn remove(&mut self, key: &K) -> Option<V>;
    fn clear(&mut self);
    fn len(&self) -> usize;
}

// 简单的内存缓存实现
struct MemoryCache<K, V> {
    data: HashMap<K, V>,
    max_size: usize,
}

impl<K, V> MemoryCache<K, V>
where
    K: Eq + Hash,
{
    fn new(max_size: usize) -> Self {
        Self {
            data: HashMap::new(),
            max_size,
        }
    }

    fn is_full(&self) -> bool {
        self.data.len() >= self.max_size
    }
}

impl<K, V> Cache<K, V> for MemoryCache<K, V>
where
    K: Eq + Hash + Clone,
{
    fn get(&self, key: &K) -> Option<&V> {
        self.data.get(key)
    }

    fn set(&mut self, key: K, value: V) {
        if self.is_full() && !self.data.contains_key(&key) {
            // 简单的 FIFO 策略：移除第一个元素
            if let Some(first_key) = self.data.keys().next().cloned() {
                self.data.remove(&first_key);
            }
        }
        self.data.insert(key, value);
    }

    fn remove(&mut self, key: &K) -> Option<V> {
        self.data.remove(key)
    }

    fn clear(&mut self) {
        self.data.clear();
    }

    fn len(&self) -> usize {
        self.data.len()
    }
}

// LRU 缓存实现
struct LRUCache<K, V> {
    data: HashMap<K, V>,
    access_order: Vec<K>,
    max_size: usize,
}

impl<K, V> LRUCache<K, V>
where
    K: Eq + Hash + Clone,
{
    fn new(max_size: usize) -> Self {
        Self {
            data: HashMap::new(),
            access_order: Vec::new(),
            max_size,
        }
    }

    fn update_access(&mut self, key: &K) {
        // 移除旧的访问记录
        self.access_order.retain(|k| k != key);
        // 添加到末尾（最近访问）
        self.access_order.push(key.clone());
    }

    fn evict_lru(&mut self) {
        if let Some(lru_key) = self.access_order.first().cloned() {
            self.data.remove(&lru_key);
            self.access_order.remove(0);
        }
    }
}

impl<K, V> Cache<K, V> for LRUCache<K, V>
where
    K: Eq + Hash + Clone,
{
    fn get(&self, key: &K) -> Option<&V> {
        self.data.get(key)
    }

    fn set(&mut self, key: K, value: V) {
        if self.data.contains_key(&key) {
            self.data.insert(key.clone(), value);
            self.update_access(&key);
        } else {
            if self.data.len() >= self.max_size {
                self.evict_lru();
            }
            self.data.insert(key.clone(), value);
            self.access_order.push(key);
        }
    }

    fn remove(&mut self, key: &K) -> Option<V> {
        self.access_order.retain(|k| k != key);
        self.data.remove(key)
    }

    fn clear(&mut self) {
        self.data.clear();
        self.access_order.clear();
    }

    fn len(&self) -> usize {
        self.data.len()
    }
}

// 带生命周期的缓存统计
struct CacheStats<'a> {
    cache_name: &'a str,
    hits: u64,
    misses: u64,
}

impl<'a> CacheStats<'a> {
    fn new(cache_name: &'a str) -> Self {
        Self {
            cache_name,
            hits: 0,
            misses: 0,
        }
    }

    fn hit(&mut self) {
        self.hits += 1;
    }

    fn miss(&mut self) {
        self.misses += 1;
    }

    fn hit_rate(&self) -> f64 {
        let total = self.hits + self.misses;
        if total == 0 {
            0.0
        } else {
            self.hits as f64 / total as f64
        }
    }
}

impl<'a> std::fmt::Display for CacheStats<'a> {
    fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
        write!(f, "{}: 命中率 {:.2}% ({} 命中, {} 未命中)",
               self.cache_name,
               self.hit_rate() * 100.0,
               self.hits,
               self.misses)
    }
}

// 带统计的缓存包装器
struct StatefulCache<C, K, V> {
    cache: C,
    stats: CacheStats<'static>,
    _phantom: std::marker::PhantomData<(K, V)>,
}

impl<C, K, V> StatefulCache<C, K, V>
where
    C: Cache<K, V>,
{
    fn new(cache: C, name: &'static str) -> Self {
        Self {
            cache,
            stats: CacheStats::new(name),
            _phantom: std::marker::PhantomData,
        }
    }

    fn get(&mut self, key: &K) -> Option<&V> {
        match self.cache.get(key) {
            Some(value) => {
                self.stats.hit();
                Some(value)
            }
            None => {
                self.stats.miss();
                None
            }
        }
    }

    fn set(&mut self, key: K, value: V) {
        self.cache.set(key, value);
    }

    fn stats(&self) -> &CacheStats {
        &self.stats
    }
}

fn main() {
    println!("=== 缓存系统演示 ===\n");

    // 测试内存缓存
    let mut memory_cache = StatefulCache::new(
        MemoryCache::new(3),
        "内存缓存"
    );

    // 添加一些数据
    memory_cache.set("key1", "value1");
    memory_cache.set("key2", "value2");
    memory_cache.set("key3", "value3");

    // 测试命中和未命中
    println!("查找 key1: {:?}", memory_cache.get(&"key1"));
    println!("查找 key2: {:?}", memory_cache.get(&"key2"));
    println!("查找 key4: {:?}", memory_cache.get(&"key4"));

    println!("{}\n", memory_cache.stats());

    // 测试 LRU 缓存
    let mut lru_cache = StatefulCache::new(
        LRUCache::new(2),
        "LRU缓存"
    );

    lru_cache.set("a", 1);
    lru_cache.set("b", 2);

    println!("LRU - 查找 a: {:?}", lru_cache.get(&"a"));

    // 添加新项，应该淘汰 b
    lru_cache.set("c", 3);

    println!("LRU - 查找 b: {:?}", lru_cache.get(&"b")); // 应该未命中
    println!("LRU - 查找 a: {:?}", lru_cache.get(&"a")); // 应该命中
    println!("LRU - 查找 c: {:?}", lru_cache.get(&"c")); // 应该命中

    println!("{}", lru_cache.stats());
}
```

## 5. 练习题

1. **实现一个泛型栈**：创建一个泛型 `Stack<T>` 结构体，实现 `push`、`pop` 和 `peek` 方法。

2. **自定义迭代器**：为你的栈实现 `Iterator` trait。

3. **比较器 trait**：创建一个 `Comparable` trait，并为不同类型实现它。

4. **生命周期练习**：编写一个函数，返回两个字符串切片中较长的那个，并确保生命周期正确。

泛型、Trait 和生命周期是 Rust 的三大核心概念。掌握它们将让你能够编写出更加灵活、安全和高效的代码。
