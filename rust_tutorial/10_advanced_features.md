# Rust 高级特性

## 1. 不安全的 Rust (Unsafe Rust)

不安全的 Rust 允许你执行五种通常被禁止的操作。

### 1.1 不安全超能力

```rust
fn main() {
    let mut num = 5;

    // 创建裸指针
    let r1 = &num as *const i32;
    let r2 = &mut num as *mut i32;

    // 创建指向任意内存地址的裸指针
    let address = 0x012345usize;
    let r = address as *const i32;

    unsafe {
        // 1. 解引用裸指针
        println!("r1 is: {}", *r1);
        println!("r2 is: {}", *r2);

        // 2. 调用不安全的函数或方法
        dangerous();

        // 3. 访问或修改可变静态变量
        HELLO_WORLD = "Hello, unsafe world!";
        println!("name is: {}", HELLO_WORLD);

        // 4. 实现不安全 trait
        // 5. 访问 union 的字段
    }
}

unsafe fn dangerous() {
    println!("这是一个不安全的函数");
}

static mut HELLO_WORLD: &str = "Hello, world!";
```

### 1.2 创建不安全代码的安全抽象

```rust
use std::slice;

fn split_at_mut(slice: &mut [i32], mid: usize) -> (&mut [i32], &mut [i32]) {
    let len = slice.len();
    let ptr = slice.as_mut_ptr();

    assert!(mid <= len);

    unsafe {
        (
            slice::from_raw_parts_mut(ptr, mid),
            slice::from_raw_parts_mut(ptr.add(mid), len - mid),
        )
    }
}

fn main() {
    let mut vector = vec![1, 2, 3, 4, 5, 6];
    let (left, right) = split_at_mut(&mut vector, 3);

    println!("左半部分: {:?}", left);
    println!("右半部分: {:?}", right);
}
```

### 1.3 使用 extern 函数调用外部代码

```rust
extern "C" {
    fn abs(input: i32) -> i32;
}

fn main() {
    unsafe {
        println!("C语言中 -3 的绝对值是: {}", abs(-3));
    }
}

// 从其他语言调用 Rust 函数
#[no_mangle]
pub extern "C" fn call_from_c() {
    println!("从C语言调用了一个Rust函数!");
}
```

### 1.4 访问和修改可变静态变量

```rust
static mut COUNTER: usize = 0;

fn add_to_count(inc: usize) {
    unsafe {
        COUNTER += inc;
    }
}

fn main() {
    add_to_count(3);

    unsafe {
        println!("COUNTER: {}", COUNTER);
    }
}
```

### 1.5 实现不安全 trait

```rust
unsafe trait Foo {
    // trait 的方法定义
}

unsafe impl Foo for i32 {
    // trait 的实现
}

fn main() {}
```

## 2. 高级 Trait

### 2.1 关联类型

```rust
pub trait Iterator {
    type Item; // 关联类型

    fn next(&mut self) -> Option<Self::Item>;
}

// 实现具有关联类型的 trait
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
    let mut counter = Counter::new(5);

    while let Some(value) = counter.next() {
        println!("计数: {}", value);
    }
}
```

### 2.2 默认泛型类型参数和运算符重载

```rust
use std::ops::Add;

#[derive(Debug, Copy, Clone, PartialEq)]
struct Point {
    x: i32,
    y: i32,
}

impl Add for Point {
    type Output = Point;

    fn add(self, other: Point) -> Point {
        Point {
            x: self.x + other.x,
            y: self.y + other.y,
        }
    }
}

// 为不同类型实现加法
struct Millimeters(u32);
struct Meters(u32);

impl Add<Meters> for Millimeters {
    type Output = Millimeters;

    fn add(self, other: Meters) -> Millimeters {
        Millimeters(self.0 + (other.0 * 1000))
    }
}

fn main() {
    assert_eq!(
        Point { x: 1, y: 0 } + Point { x: 2, y: 3 },
        Point { x: 3, y: 3 }
    );

    println!("两点相加: {:?}", Point { x: 1, y: 2 } + Point { x: 3, y: 4 });
}
```

### 2.3 完全限定语法与消歧义

```rust
trait Pilot {
    fn fly(&self);
}

trait Wizard {
    fn fly(&self);
}

struct Human;

impl Pilot for Human {
    fn fly(&self) {
        println!("驾驶飞机飞行");
    }
}

impl Wizard for Human {
    fn fly(&self) {
        println!("用魔法飞行");
    }
}

impl Human {
    fn fly(&self) {
        println!("*疯狂摆动双臂*");
    }
}

// 关联函数的消歧义
trait Animal {
    fn baby_name() -> String;
}

struct Dog;

impl Dog {
    fn baby_name() -> String {
        String::from("Spot")
    }
}

impl Animal for Dog {
    fn baby_name() -> String {
        String::from("puppy")
    }
}

fn main() {
    let person = Human;
    person.fly(); // 调用 Human 的方法
    Pilot::fly(&person); // 调用 Pilot trait 的方法
    Wizard::fly(&person); // 调用 Wizard trait 的方法

    // 关联函数消歧义
    println!("小狗叫: {}", Dog::baby_name());
    println!("小动物叫: {}", <Dog as Animal>::baby_name());
}
```

### 2.4 超级 Trait

```rust
use std::fmt;

trait OutlinePrint: fmt::Display {
    fn outline_print(&self) {
        let output = self.to_string();
        let len = output.len();
        println!("{}", "*".repeat(len + 4));
        println!("*{}*", " ".repeat(len + 2));
        println!("* {} *", output);
        println!("*{}*", " ".repeat(len + 2));
        println!("{}", "*".repeat(len + 4));
    }
}

struct Point {
    x: i32,
    y: i32,
}

impl fmt::Display for Point {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "({}, {})", self.x, self.y)
    }
}

impl OutlinePrint for Point {}

fn main() {
    let point = Point { x: 1, y: 3 };
    point.outline_print();
}
```

### 2.5 newtype 模式

```rust
use std::fmt;

struct Wrapper(Vec<String>);

impl fmt::Display for Wrapper {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "[{}]", self.0.join(", "))
    }
}

fn main() {
    let w = Wrapper(vec![String::from("hello"), String::from("world")]);
    println!("w = {}", w);
}
```

## 3. 高级类型

### 3.1 类型别名

```rust
type Kilometers = i32;
type Result<T> = std::result::Result<T, std::io::Error>;

fn main() {
    let x: i32 = 5;
    let y: Kilometers = 5;

    println!("x + y = {}", x + y);
}

// 减少重复的长类型
type Thunk = Box<dyn Fn() + Send + 'static>;

fn takes_long_type(f: Thunk) {
    // --snip--
}

fn returns_long_type() -> Thunk {
    Box::new(|| println!("hi"))
}
```

### 3.2 Never 类型

```rust
fn bar() -> ! {
    panic!("这个函数永远不会返回");
}

fn main() {
    let guess = "3";

    let guess: u32 = match guess.trim().parse() {
        Ok(num) => num,
        Err(_) => panic!("这不是一个数字!"), // panic! 的类型是 !
    };

    println!("你猜的数字是: {}", guess);
}
```

### 3.3 动态大小类型和 Sized trait

```rust
// str 是动态大小类型
fn generic<T: ?Sized>(t: &T) {
    // --snip--
}

// 等同于
fn generic_explicit<T: Sized>(t: T) {
    // --snip--
}

fn main() {
    let s1: &str = "Hello there!";
    let s2: &str = "How's it going?";

    // generic(s1); // 如果没有 ?Sized，这会编译错误
}
```

## 4. 高级函数和闭包

### 4.1 函数指针

```rust
fn add_one(x: i32) -> i32 {
    x + 1
}

fn do_twice(f: fn(i32) -> i32, arg: i32) -> i32 {
    f(arg) + f(arg)
}

fn main() {
    let answer = do_twice(add_one, 5);
    println!("答案是: {}", answer);

    // 使用函数指针与闭包
    let list_of_numbers = vec![1, 2, 3];
    let list_of_strings: Vec<String> = list_of_numbers
        .iter()
        .map(|i| i.to_string())
        .collect();

    // 或者
    let list_of_strings: Vec<String> = list_of_numbers
        .iter()
        .map(ToString::to_string)
        .collect();

    println!("字符串列表: {:?}", list_of_strings);
}
```

### 4.2 返回闭包

```rust
fn returns_closure() -> Box<dyn Fn(i32) -> i32> {
    Box::new(|x| x + 1)
}

fn main() {
    let f = returns_closure();
    println!("结果: {}", f(3));
}
```

## 5. 宏

### 5.1 声明式宏

```rust
// 简单的宏
macro_rules! say_hello {
    () => {
        println!("Hello, macro!");
    };
}

// 带参数的宏
macro_rules! create_function {
    ($func_name:ident) => {
        fn $func_name() {
            println!("函数 {:?} 被调用", stringify!($func_name));
        }
    };
}

// 重载的宏
macro_rules! print_result {
    ($expression:expr) => {
        println!("{:?} = {:?}", stringify!($expression), $expression);
    };
}

// vec! 宏的简化实现
macro_rules! vec_simple {
    ( $( $x:expr ),* ) => {
        {
            let mut temp_vec = Vec::new();
            $(
                temp_vec.push($x);
            )*
            temp_vec
        }
    };
}

create_function!(foo);
create_function!(bar);

fn main() {
    say_hello!();

    foo();
    bar();

    print_result!(1u32 + 1);
    print_result!({
        let mut vec = Vec::new();
        vec.push(1);
        vec.push(2);
        vec
    });

    let v = vec_simple![1, 2, 3];
    println!("向量: {:?}", v);
}
```

### 5.2 过程宏

```rust
// 注意：这需要在独立的 crate 中实现
/*
use proc_macro::TokenStream;
use quote::quote;
use syn::{parse_macro_input, DeriveInput};

#[proc_macro_derive(HelloMacro)]
pub fn hello_macro_derive(input: TokenStream) -> TokenStream {
    let ast = parse_macro_input!(input as DeriveInput);
    impl_hello_macro(&ast)
}

fn impl_hello_macro(ast: &syn::DeriveInput) -> TokenStream {
    let name = &ast.ident;
    let gen = quote! {
        impl HelloMacro for #name {
            fn hello_macro() {
                println!("Hello, Macro! My name is {}!", stringify!(#name));
            }
        }
    };
    gen.into()
}
*/
```

## 6. 实践项目：内存池分配器

```rust
use std::alloc::{GlobalAlloc, Layout, System};
use std::ptr;
use std::sync::atomic::{AtomicUsize, Ordering};

// 统计分配器
struct CountingAllocator;

static ALLOCATED: AtomicUsize = AtomicUsize::new(0);
static DEALLOCATED: AtomicUsize = AtomicUsize::new(0);

unsafe impl GlobalAlloc for CountingAllocator {
    unsafe fn alloc(&self, layout: Layout) -> *mut u8 {
        let ret = System.alloc(layout);
        if !ret.is_null() {
            ALLOCATED.fetch_add(layout.size(), Ordering::SeqCst);
        }
        ret
    }

    unsafe fn dealloc(&self, ptr: *mut u8, layout: Layout) {
        System.dealloc(ptr, layout);
        DEALLOCATED.fetch_add(layout.size(), Ordering::SeqCst);
    }
}

//#[global_allocator]
//static GLOBAL: CountingAllocator = CountingAllocator;

// 简单的内存池
struct MemoryPool {
    memory: Vec<u8>,
    offset: AtomicUsize,
}

impl MemoryPool {
    fn new(size: usize) -> Self {
        Self {
            memory: vec![0; size],
            offset: AtomicUsize::new(0),
        }
    }

    fn allocate(&self, size: usize, align: usize) -> Option<*mut u8> {
        let current_offset = self.offset.load(Ordering::SeqCst);
        let aligned_offset = (current_offset + align - 1) & !(align - 1);
        let new_offset = aligned_offset + size;

        if new_offset > self.memory.len() {
            return None; // 内存不足
        }

        // 尝试更新偏移量
        match self.offset.compare_exchange(
            current_offset,
            new_offset,
            Ordering::SeqCst,
            Ordering::SeqCst,
        ) {
            Ok(_) => {
                // 成功分配
                unsafe {
                    Some(self.memory.as_ptr().add(aligned_offset) as *mut u8)
                }
            }
            Err(_) => {
                // 并发冲突，重试
                self.allocate(size, align)
            }
        }
    }

    fn reset(&self) {
        self.offset.store(0, Ordering::SeqCst);
    }

    fn usage(&self) -> (usize, usize) {
        let used = self.offset.load(Ordering::SeqCst);
        (used, self.memory.len())
    }
}

// 智能指针包装
struct PoolPtr<T> {
    ptr: *mut T,
    _marker: std::marker::PhantomData<T>,
}

impl<T> PoolPtr<T> {
    fn new(ptr: *mut T) -> Self {
        Self {
            ptr,
            _marker: std::marker::PhantomData,
        }
    }
}

impl<T> std::ops::Deref for PoolPtr<T> {
    type Target = T;

    fn deref(&self) -> &Self::Target {
        unsafe { &*self.ptr }
    }
}

impl<T> std::ops::DerefMut for PoolPtr<T> {
    fn deref_mut(&mut self) -> &mut Self::Target {
        unsafe { &mut *self.ptr }
    }
}

// 高级内存管理器
struct AdvancedMemoryManager {
    pools: Vec<MemoryPool>,
    current_pool: AtomicUsize,
}

impl AdvancedMemoryManager {
    fn new(pool_count: usize, pool_size: usize) -> Self {
        let mut pools = Vec::with_capacity(pool_count);
        for _ in 0..pool_count {
            pools.push(MemoryPool::new(pool_size));
        }

        Self {
            pools,
            current_pool: AtomicUsize::new(0),
        }
    }

    fn allocate<T>(&self) -> Option<PoolPtr<T>> {
        let size = std::mem::size_of::<T>();
        let align = std::mem::align_of::<T>();

        // 尝试从当前池分配
        let current = self.current_pool.load(Ordering::SeqCst);

        for i in 0..self.pools.len() {
            let pool_index = (current + i) % self.pools.len();

            if let Some(ptr) = self.pools[pool_index].allocate(size, align) {
                return Some(PoolPtr::new(ptr as *mut T));
            }
        }

        None // 所有池都满了
    }

    fn reset_all(&self) {
        for pool in &self.pools {
            pool.reset();
        }
    }

    fn get_stats(&self) -> Vec<(usize, usize)> {
        self.pools.iter().map(|pool| pool.usage()).collect()
    }
}

// 使用示例
#[derive(Debug)]
struct TestStruct {
    id: u32,
    data: [u8; 64],
}

impl TestStruct {
    fn new(id: u32) -> Self {
        Self {
            id,
            data: [id as u8; 64],
        }
    }
}

fn main() {
    println!("=== 高级内存管理演示 ===\n");

    // 创建内存管理器
    let manager = AdvancedMemoryManager::new(3, 1024); // 3个池，每个1KB

    // 分配一些对象
    let mut objects = Vec::new();

    for i in 0..10 {
        if let Some(mut obj) = manager.allocate::<TestStruct>() {
            unsafe {
                ptr::write(obj.ptr, TestStruct::new(i));
            }
            println!("分配对象 {}: {:?}", i, *obj);
            objects.push(obj);
        } else {
            println!("分配对象 {} 失败", i);
        }
    }

    // 显示内存使用情况
    println!("\n内存池使用情况:");
    for (i, (used, total)) in manager.get_stats().iter().enumerate() {
        println!("池 {}: {}/{} 字节 ({:.1}%)",
                 i, used, total, (*used as f32 / *total as f32) * 100.0);
    }

    // 演示不安全操作
    unsafe {
        println!("\n=== 不安全操作演示 ===");

        // 创建裸指针
        let x = 42;
        let raw = &x as *const i32;

        println!("通过裸指针访问: {}", *raw);

        // 调用外部函数
        println!("C语言abs(-42): {}", abs(-42));
    }

    // 显示全局分配统计（如果启用了自定义分配器）
    println!("\n全局分配统计:");
    println!("已分配: {} 字节", ALLOCATED.load(Ordering::SeqCst));
    println!("已释放: {} 字节", DEALLOCATED.load(Ordering::SeqCst));
}

extern "C" {
    fn abs(input: i32) -> i32;
}
```

## 7. 高级模式匹配

### 7.1 析构结构体和元组

```rust
struct Point {
    x: i32,
    y: i32,
}

fn main() {
    let p = Point { x: 0, y: 7 };

    match p {
        Point { x, y: 0 } => println!("在x轴上，x = {}", x),
        Point { x: 0, y } => println!("在y轴上，y = {}", y),
        Point { x, y } => println!("在其他位置：({}, {})", x, y),
    }

    // 嵌套析构
    let ((feet, inches), Point { x, y }) = ((3, 10), Point { x: 3, y: -10 });
    println!("脚: {}, 英寸: {}, 点: ({}, {})", feet, inches, x, y);
}
```

### 7.2 匹配守卫

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

### 7.3 @ 绑定

```rust
enum Message {
    Hello { id: i32 },
}

fn main() {
    let msg = Message::Hello { id: 5 };

    match msg {
        Message::Hello {
            id: id_variable @ 3..=7,
        } => println!("找到一个在范围内的id: {}", id_variable),
        Message::Hello { id: 10..=12 } => {
            println!("找到一个在另一个范围内的id")
        }
        Message::Hello { id } => println!("找到其他id: {}", id),
    }
}
```

## 8. 练习题

1. **实现自定义智能指针**：创建一个类似 `Box<T>` 的智能指针。

2. **宏练习**：编写一个宏来生成简单的 getter 和 setter 方法。

3. **内存管理**：实现一个简单的垃圾回收器概念验证。

4. **类型系统探索**：创建一个类型安全的状态机。

Rust 的高级特性为你提供了强大的工具来构建复杂、高性能的系统。虽然这些特性需要谨慎使用，但它们为解决特定问题提供了必要的灵活性。
