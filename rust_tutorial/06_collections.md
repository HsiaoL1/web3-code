# Rust 集合类型

## 1. 集合概述

Rust 标准库包含许多非常有用的数据结构，称为集合。大部分其他数据类型都代表一个特定的值，不过集合可以包含多个值。集合指向的数据是储存在堆上的，这意味着数据的数量不必在编译时就已知，并且还可以随着程序的运行增长或缩小。

## 2. 动态数组 Vec<T>

### 2.1 创建 Vector

```rust
fn main() {
    // 创建新的空 vector
    let v: Vec<i32> = Vec::new();

    // 使用宏创建包含初始值的 vector
    let v = vec![1, 2, 3];

    // 创建可变 vector
    let mut v = Vec::new();
    v.push(5);
    v.push(6);
    v.push(7);
    v.push(8);

    println!("v = {:?}", v);
}
```

### 2.2 读取 Vector 的元素

```rust
fn main() {
    let v = vec![1, 2, 3, 4, 5];

    // 方法1: 使用索引语法
    let third: &i32 = &v[2];
    println!("第三个元素是 {}", third);

    // 方法2: 使用 get 方法
    match v.get(2) {
        Some(third) => println!("第三个元素是 {}", third),
        None => println!("没有第三个元素。"),
    }

    // 处理无效索引
    // let does_not_exist = &v[100]; // 这会导致程序崩溃!

    match v.get(100) {
        Some(element) => println!("元素: {}", element),
        None => println!("索引超出范围"),
    }
}
```

### 2.3 遍历 Vector

```rust
fn main() {
    let v = vec![100, 32, 57];

    // 不可变引用遍历
    for i in &v {
        println!("{}", i);
    }

    // 可变引用遍历
    let mut v = vec![100, 32, 57];
    for i in &mut v {
        *i += 50; // 解引用并修改值
    }

    println!("修改后: {:?}", v);

    // 获取所有权遍历
    let v = vec![1, 2, 3];
    for i in v {
        println!("{}", i);
    }
    // println!("{:?}", v); // 编译错误! v 的所有权已被移动
}
```

### 2.4 Vector 的常用方法

```rust
fn main() {
    let mut v = vec![1, 2, 3, 4, 5];

    // 长度
    println!("长度: {}", v.len());

    // 容量
    println!("容量: {}", v.capacity());

    // 添加元素
    v.push(6);
    println!("添加后: {:?}", v);

    // 删除最后一个元素
    if let Some(last) = v.pop() {
        println!("删除的元素: {}", last);
    }
    println!("删除后: {:?}", v);

    // 插入元素
    v.insert(2, 99);
    println!("插入后: {:?}", v);

    // 删除指定位置元素
    let removed = v.remove(2);
    println!("删除的元素: {}, 结果: {:?}", removed, v);

    // 清空
    v.clear();
    println!("清空后: {:?}", v);

    // 检查是否为空
    println!("是否为空: {}", v.is_empty());
}
```

### 2.5 使用枚举储存多种类型

```rust
#[derive(Debug)]
enum SpreadsheetCell {
    Int(i32),
    Float(f64),
    Text(String),
}

fn main() {
    let row = vec![
        SpreadsheetCell::Int(3),
        SpreadsheetCell::Text(String::from("blue")),
        SpreadsheetCell::Float(10.12),
    ];

    for cell in &row {
        match cell {
            SpreadsheetCell::Int(i) => println!("整数: {}", i),
            SpreadsheetCell::Float(f) => println!("浮点数: {}", f),
            SpreadsheetCell::Text(s) => println!("文本: {}", s),
        }
    }

    println!("完整行: {:?}", row);
}
```

## 3. 字符串 String

### 3.1 创建字符串

```rust
fn main() {
    // 创建新的空字符串
    let mut s = String::new();

    // 从字符串字面值创建
    let data = "initial contents";
    let s = data.to_string();
    let s = "initial contents".to_string();
    let s = String::from("initial contents");

    // UTF-8 编码
    let hello = String::from("السلام عليكم");
    let hello = String::from("Dobrý den");
    let hello = String::from("Hello");
    let hello = String::from("שָׁלוֹם");
    let hello = String::from("नमस्ते");
    let hello = String::from("こんにちは");
    let hello = String::from("안녕하세요");
    let hello = String::from("你好");
    let hello = String::from("Olá");
    let hello = String::from("Здравствуйте");
    let hello = String::from("Hola");

    println!("各种语言的你好: {}", hello);
}
```

### 3.2 更新字符串

```rust
fn main() {
    // push_str 方法
    let mut s = String::from("foo");
    s.push_str("bar");
    println!("s = {}", s);

    // push 方法添加单个字符
    let mut s = String::from("lo");
    s.push('l');
    println!("s = {}", s);

    // 使用 + 运算符连接字符串
    let s1 = String::from("Hello, ");
    let s2 = String::from("world!");
    let s3 = s1 + &s2; // s1 被移动了，不能继续使用
    println!("s3 = {}", s3);
    // println!("s1 = {}", s1); // 编译错误!

    // 复杂的字符串连接
    let s1 = String::from("tic");
    let s2 = String::from("tac");
    let s3 = String::from("toe");
    let s = s1 + "-" + &s2 + "-" + &s3;
    println!("s = {}", s);

    // 使用 format! 宏
    let s1 = String::from("tic");
    let s2 = String::from("tac");
    let s3 = String::from("toe");
    let s = format!("{}-{}-{}", s1, s2, s3);
    println!("s = {}", s);
    println!("s1 仍然有效: {}", s1); // s1 仍然有效
}
```

### 3.3 字符串切片和索引

```rust
fn main() {
    let hello = "Здравствуйте";
    let s = &hello[0..4]; // 注意：这里是字节索引，不是字符索引
    println!("切片: {}", s);

    // 遍历字符串的方法

    // 按字符遍历
    for c in "नमस्ते".chars() {
        println!("{}", c);
    }

    // 按字节遍历
    for b in "नमस्ते".bytes() {
        println!("{}", b);
    }

    // 字符串长度
    let s = String::from("hello");
    println!("字节长度: {}", s.len());
    println!("字符数量: {}", s.chars().count());

    let s = String::from("नमस्ते");
    println!("字节长度: {}", s.len());
    println!("字符数量: {}", s.chars().count());
}
```

### 3.4 字符串的常用方法

```rust
fn main() {
    let s = String::from("  Hello, World!  ");

    // 去除空白
    println!("原始: '{}'", s);
    println!("trim: '{}'", s.trim());
    println!("trim_start: '{}'", s.trim_start());
    println!("trim_end: '{}'", s.trim_end());

    // 替换
    let s = "Hello, World!";
    println!("replace: {}", s.replace("World", "Rust"));
    println!("replace_n: {}", s.replacen("l", "L", 2));

    // 分割
    let s = "apple,banana,cherry";
    let fruits: Vec<&str> = s.split(',').collect();
    println!("分割结果: {:?}", fruits);

    // 检查包含
    let s = "Hello, World!";
    println!("contains 'World': {}", s.contains("World"));
    println!("starts_with 'Hello': {}", s.starts_with("Hello"));
    println!("ends_with '!': {}", s.ends_with("!"));

    // 大小写转换
    let s = "Hello, World!";
    println!("to_lowercase: {}", s.to_lowercase());
    println!("to_uppercase: {}", s.to_uppercase());

    // 重复
    let s = "na";
    println!("repeat: {}", s.repeat(3));
}
```

## 4. 哈希映射 HashMap<K, V>

### 4.1 创建 HashMap

```rust
use std::collections::HashMap;

fn main() {
    // 创建新的哈希映射
    let mut scores = HashMap::new();

    scores.insert(String::from("Blue"), 10);
    scores.insert(String::from("Yellow"), 50);

    println!("scores: {:?}", scores);

    // 从 vector 创建
    let teams = vec![String::from("Blue"), String::from("Yellow")];
    let initial_scores = vec![10, 50];

    let mut scores: HashMap<_, _> = teams.into_iter()
        .zip(initial_scores.into_iter())
        .collect();

    println!("从 vector 创建: {:?}", scores);
}
```

### 4.2 访问 HashMap 中的值

```rust
use std::collections::HashMap;

fn main() {
    let mut scores = HashMap::new();

    scores.insert(String::from("Blue"), 10);
    scores.insert(String::from("Yellow"), 50);

    // 使用 get 方法
    let team_name = String::from("Blue");
    let score = scores.get(&team_name);

    match score {
        Some(s) => println!("{} 队的分数: {}", team_name, s),
        None => println!("队伍不存在"),
    }

    // 遍历哈希映射
    for (key, value) in &scores {
        println!("{}: {}", key, value);
    }

    // 检查键是否存在
    if scores.contains_key("Blue") {
        println!("Blue 队存在");
    }
}
```

### 4.3 更新 HashMap

```rust
use std::collections::HashMap;

fn main() {
    let mut scores = HashMap::new();

    // 覆盖值
    scores.insert(String::from("Blue"), 10);
    scores.insert(String::from("Blue"), 25);
    println!("覆盖后: {:?}", scores);

    // 只在键没有值时插入
    scores.entry(String::from("Yellow")).or_insert(50);
    scores.entry(String::from("Blue")).or_insert(50);
    println!("条件插入后: {:?}", scores);

    // 根据旧值更新值
    let text = "hello world wonderful world";
    let mut map = HashMap::new();

    for word in text.split_whitespace() {
        let count = map.entry(word).or_insert(0);
        *count += 1;
    }

    println!("单词计数: {:?}", map);
}
```

### 4.4 HashMap 的常用方法

```rust
use std::collections::HashMap;

fn main() {
    let mut map = HashMap::new();
    map.insert("a", 1);
    map.insert("b", 2);
    map.insert("c", 3);

    // 长度
    println!("长度: {}", map.len());

    // 是否为空
    println!("是否为空: {}", map.is_empty());

    // 删除元素
    if let Some(value) = map.remove("b") {
        println!("删除的值: {}", value);
    }
    println!("删除后: {:?}", map);

    // 清空
    map.clear();
    println!("清空后: {:?}", map);

    // 容量相关
    let mut map = HashMap::with_capacity(10);
    println!("预分配容量的映射: {:?}", map);
}
```

## 5. 其他集合类型

### 5.1 VecDeque（双端队列）

```rust
use std::collections::VecDeque;

fn main() {
    let mut deque = VecDeque::new();

    // 前端添加
    deque.push_front(1);
    deque.push_front(2);

    // 后端添加
    deque.push_back(3);
    deque.push_back(4);

    println!("deque: {:?}", deque);

    // 前端和后端删除
    println!("前端删除: {:?}", deque.pop_front());
    println!("后端删除: {:?}", deque.pop_back());

    println!("最终 deque: {:?}", deque);
}
```

### 5.2 HashSet（哈希集合）

```rust
use std::collections::HashSet;

fn main() {
    let mut books = HashSet::new();

    books.insert("A Song of Ice and Fire");
    books.insert("The Hobbit");
    books.insert("Harry Potter");
    books.insert("The Hobbit"); // 重复插入，不会有效果

    println!("书籍集合: {:?}", books);

    // 检查是否包含
    if books.contains("The Hobbit") {
        println!("找到了《霍比特人》");
    }

    // 集合操作
    let fiction: HashSet<&str> = ["Harry Potter", "The Hobbit"].iter().cloned().collect();
    let fantasy: HashSet<&str> = ["The Hobbit", "Lord of the Rings"].iter().cloned().collect();

    // 交集
    let intersection: HashSet<_> = fiction.intersection(&fantasy).collect();
    println!("交集: {:?}", intersection);

    // 并集
    let union: HashSet<_> = fiction.union(&fantasy).collect();
    println!("并集: {:?}", union);

    // 差集
    let difference: HashSet<_> = fiction.difference(&fantasy).collect();
    println!("差集: {:?}", difference);
}
```

### 5.3 BTreeMap（有序映射）

```rust
use std::collections::BTreeMap;

fn main() {
    let mut map = BTreeMap::new();

    map.insert(3, "three");
    map.insert(1, "one");
    map.insert(2, "two");
    map.insert(5, "five");
    map.insert(4, "four");

    // BTreeMap 会保持键的排序
    for (key, value) in &map {
        println!("{}: {}", key, value);
    }

    // 范围查询
    println!("2 到 4 的范围:");
    for (key, value) in map.range(2..=4) {
        println!("{}: {}", key, value);
    }
}
```

## 6. 实践项目：学生成绩管理系统

```rust
use std::collections::HashMap;

#[derive(Debug, Clone)]
struct Student {
    id: u32,
    name: String,
    grades: Vec<f64>,
}

impl Student {
    fn new(id: u32, name: String) -> Self {
        Self {
            id,
            name,
            grades: Vec::new(),
        }
    }

    fn add_grade(&mut self, grade: f64) {
        if grade >= 0.0 && grade <= 100.0 {
            self.grades.push(grade);
        }
    }

    fn average_grade(&self) -> Option<f64> {
        if self.grades.is_empty() {
            None
        } else {
            let sum: f64 = self.grades.iter().sum();
            Some(sum / self.grades.len() as f64)
        }
    }

    fn highest_grade(&self) -> Option<f64> {
        self.grades.iter().cloned().fold(None, |max, grade| {
            Some(max.map_or(grade, |m| m.max(grade)))
        })
    }

    fn lowest_grade(&self) -> Option<f64> {
        self.grades.iter().cloned().fold(None, |min, grade| {
            Some(min.map_or(grade, |m| m.min(grade)))
        })
    }
}

struct GradeManager {
    students: HashMap<u32, Student>,
    subjects: HashMap<String, Vec<(u32, f64)>>, // 科目 -> (学生ID, 成绩)
}

impl GradeManager {
    fn new() -> Self {
        Self {
            students: HashMap::new(),
            subjects: HashMap::new(),
        }
    }

    fn add_student(&mut self, student: Student) {
        self.students.insert(student.id, student);
    }

    fn add_grade(&mut self, student_id: u32, subject: String, grade: f64) -> Result<(), String> {
        let student = self.students.get_mut(&student_id)
            .ok_or("学生不存在")?;

        student.add_grade(grade);

        // 记录科目成绩
        self.subjects.entry(subject).or_insert_with(Vec::new).push((student_id, grade));

        Ok(())
    }

    fn get_student_average(&self, student_id: u32) -> Result<Option<f64>, String> {
        let student = self.students.get(&student_id)
            .ok_or("学生不存在")?;
        Ok(student.average_grade())
    }

    fn get_subject_average(&self, subject: &str) -> Option<f64> {
        if let Some(grades) = self.subjects.get(subject) {
            if grades.is_empty() {
                None
            } else {
                let sum: f64 = grades.iter().map(|(_, grade)| grade).sum();
                Some(sum / grades.len() as f64)
            }
        } else {
            None
        }
    }

    fn get_top_students(&self, n: usize) -> Vec<(u32, String, f64)> {
        let mut student_averages: Vec<_> = self.students
            .iter()
            .filter_map(|(id, student)| {
                student.average_grade().map(|avg| (*id, student.name.clone(), avg))
            })
            .collect();

        student_averages.sort_by(|a, b| b.2.partial_cmp(&a.2).unwrap());
        student_averages.into_iter().take(n).collect()
    }

    fn list_all_students(&self) -> Vec<&Student> {
        self.students.values().collect()
    }

    fn list_subjects(&self) -> Vec<&String> {
        self.subjects.keys().collect()
    }
}

fn main() {
    let mut manager = GradeManager::new();

    // 添加学生
    manager.add_student(Student::new(1, "张三".to_string()));
    manager.add_student(Student::new(2, "李四".to_string()));
    manager.add_student(Student::new(3, "王五".to_string()));

    // 添加成绩
    let grades_data = vec![
        (1, "数学", 85.0),
        (1, "英语", 92.0),
        (1, "物理", 88.0),
        (2, "数学", 78.0),
        (2, "英语", 95.0),
        (2, "物理", 82.0),
        (3, "数学", 92.0),
        (3, "英语", 89.0),
        (3, "物理", 95.0),
    ];

    for (student_id, subject, grade) in grades_data {
        if let Err(e) = manager.add_grade(student_id, subject.to_string(), grade) {
            println!("添加成绩错误: {}", e);
        }
    }

    // 显示所有学生信息
    println!("=== 所有学生信息 ===");
    for student in manager.list_all_students() {
        println!("学生: {} (ID: {})", student.name, student.id);
        println!("  所有成绩: {:?}", student.grades);
        if let Some(avg) = student.average_grade() {
            println!("  平均成绩: {:.2}", avg);
        }
        if let Some(highest) = student.highest_grade() {
            println!("  最高成绩: {:.2}", highest);
        }
        if let Some(lowest) = student.lowest_grade() {
            println!("  最低成绩: {:.2}", lowest);
        }
        println!();
    }

    // 显示各科目平均成绩
    println!("=== 各科目平均成绩 ===");
    for subject in manager.list_subjects() {
        if let Some(avg) = manager.get_subject_average(subject) {
            println!("{}: {:.2}", subject, avg);
        }
    }

    // 显示前3名学生
    println!("\n=== 前3名学生 ===");
    for (rank, (id, name, avg)) in manager.get_top_students(3).iter().enumerate() {
        println!("第{}名: {} (ID: {}) - 平均分: {:.2}", rank + 1, name, id, avg);
    }
}
```

## 7. 集合性能比较

```rust
use std::collections::{HashMap, BTreeMap, HashSet, BTreeSet};
use std::time::Instant;

fn main() {
    let data: Vec<i32> = (0..100000).collect();

    // HashMap vs BTreeMap 插入性能
    let start = Instant::now();
    let mut hash_map = HashMap::new();
    for &i in &data {
        hash_map.insert(i, i * 2);
    }
    println!("HashMap 插入耗时: {:?}", start.elapsed());

    let start = Instant::now();
    let mut btree_map = BTreeMap::new();
    for &i in &data {
        btree_map.insert(i, i * 2);
    }
    println!("BTreeMap 插入耗时: {:?}", start.elapsed());

    // 查找性能比较
    let start = Instant::now();
    for i in 0..1000 {
        hash_map.get(&i);
    }
    println!("HashMap 查找耗时: {:?}", start.elapsed());

    let start = Instant::now();
    for i in 0..1000 {
        btree_map.get(&i);
    }
    println!("BTreeMap 查找耗时: {:?}", start.elapsed());
}
```

## 8. 练习题

1. **词频统计器**：编写一个程序，统计文本中每个单词的出现次数。

2. **购物车系统**：使用 HashMap 实现一个购物车，支持添加商品、删除商品、计算总价等功能。

3. **学生选课系统**：使用多种集合类型实现一个学生选课系统，包括学生信息、课程信息和选课关系。

集合是 Rust 编程中最常用的数据结构之一。通过熟练掌握 Vec、String、HashMap 等集合类型，你将能够处理更复杂的数据结构和算法问题。
