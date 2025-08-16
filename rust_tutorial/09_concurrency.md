# Rust 并发编程

## 1. 并发概述

Rust 的所有权系统和类型系统为并发编程提供了强大的安全保障，可以在编译时防止数据竞争。

### 1.1 并发 vs 并行

- **并发（Concurrency）**：处理多个任务，但不一定同时执行
- **并行（Parallelism）**：同时执行多个任务

```rust
// Rust 并发编程的核心原则：
// 1. 所有权防止数据竞争
// 2. 类型系统确保线程安全
// 3. 编译时检查而非运行时
```

## 2. 线程 (Threads)

### 2.1 创建线程

```rust
use std::thread;
use std::time::Duration;

fn main() {
    // 创建新线程
    thread::spawn(|| {
        for i in 1..10 {
            println!("来自生成线程的数字 {}!", i);
            thread::sleep(Duration::from_millis(1));
        }
    });

    // 主线程
    for i in 1..5 {
        println!("来自主线程的数字 {}!", i);
        thread::sleep(Duration::from_millis(1));
    }
}
```

### 2.2 等待线程完成

```rust
use std::thread;
use std::time::Duration;

fn main() {
    let handle = thread::spawn(|| {
        for i in 1..10 {
            println!("来自生成线程的数字 {}!", i);
            thread::sleep(Duration::from_millis(1));
        }
        "线程完成了!"
    });

    for i in 1..5 {
        println!("来自主线程的数字 {}!", i);
        thread::sleep(Duration::from_millis(1));
    }

    // 等待线程完成并获取返回值
    let result = handle.join().unwrap();
    println!("线程返回: {}", result);
}
```

### 2.3 move 闭包

```rust
use std::thread;

fn main() {
    let v = vec![1, 2, 3];

    // 使用 move 将所有权转移到线程
    let handle = thread::spawn(move || {
        println!("这是向量: {:?}", v);
    });

    handle.join().unwrap();

    // println!("v = {:?}", v); // 编译错误! v 的所有权已经移动
}
```

### 2.4 线程配置

```rust
use std::thread;

fn main() {
    // 配置线程
    let builder = thread::Builder::new()
        .name("worker".into())
        .stack_size(32 * 1024); // 32KB 栈大小

    let handle = builder.spawn(|| {
        println!("当前线程: {:?}", thread::current().name());
        println!("线程 ID: {:?}", thread::current().id());
        42
    }).unwrap();

    let result = handle.join().unwrap();
    println!("工作线程返回: {}", result);
}
```

## 3. 消息传递

### 3.1 通道 (Channels)

```rust
use std::sync::mpsc;
use std::thread;

fn main() {
    // 创建通道
    let (tx, rx) = mpsc::channel();

    thread::spawn(move || {
        let val = String::from("hi");
        tx.send(val).unwrap();
        // println!("val = {}", val); // 编译错误! val 已被移动
    });

    // 接收消息
    let received = rx.recv().unwrap();
    println!("收到: {}", received);
}
```

### 3.2 发送多个值

```rust
use std::sync::mpsc;
use std::thread;
use std::time::Duration;

fn main() {
    let (tx, rx) = mpsc::channel();

    thread::spawn(move || {
        let vals = vec![
            String::from("hi"),
            String::from("from"),
            String::from("the"),
            String::from("thread"),
        ];

        for val in vals {
            tx.send(val).unwrap();
            thread::sleep(Duration::from_secs(1));
        }
    });

    // 接收多个消息
    for received in rx {
        println!("收到: {}", received);
    }
}
```

### 3.3 多个发送者

```rust
use std::sync::mpsc;
use std::thread;
use std::time::Duration;

fn main() {
    let (tx, rx) = mpsc::channel();
    let tx1 = tx.clone();

    thread::spawn(move || {
        let vals = vec![
            String::from("thread1: hi"),
            String::from("thread1: from"),
        ];

        for val in vals {
            tx1.send(val).unwrap();
            thread::sleep(Duration::from_secs(1));
        }
    });

    thread::spawn(move || {
        let vals = vec![
            String::from("thread2: more"),
            String::from("thread2: messages"),
        ];

        for val in vals {
            tx.send(val).unwrap();
            thread::sleep(Duration::from_secs(1));
        }
    });

    for received in rx {
        println!("收到: {}", received);
    }
}
```

### 3.4 同步通道

```rust
use std::sync::mpsc;
use std::thread;
use std::time::Duration;

fn main() {
    // 创建同步通道，容量为 0 (完全同步)
    let (tx, rx) = mpsc::sync_channel(0);

    let handle = thread::spawn(move || {
        println!("发送前...");
        tx.send("同步消息").unwrap();
        println!("发送后...");
    });

    thread::sleep(Duration::from_secs(2));
    println!("接收: {}", rx.recv().unwrap());

    handle.join().unwrap();
}
```

## 4. 共享状态并发

### 4.1 互斥器 (Mutex)

```rust
use std::sync::Mutex;

fn main() {
    let m = Mutex::new(5);

    {
        let mut num = m.lock().unwrap();
        *num = 6;
    } // 锁在这里释放

    println!("m = {:?}", m);
}
```

### 4.2 在多个线程间共享 Mutex

```rust
use std::rc::Rc;
use std::sync::{Arc, Mutex};
use std::thread;

fn main() {
    // Arc: Atomic Reference Counting
    let counter = Arc::new(Mutex::new(0));
    let mut handles = vec![];

    for _ in 0..10 {
        let counter = Arc::clone(&counter);
        let handle = thread::spawn(move || {
            let mut num = counter.lock().unwrap();
            *num += 1;
        });
        handles.push(handle);
    }

    for handle in handles {
        handle.join().unwrap();
    }

    println!("结果: {}", *counter.lock().unwrap());
}
```

### 4.3 读写锁 (RwLock)

```rust
use std::sync::{Arc, RwLock};
use std::thread;
use std::time::Duration;

fn main() {
    let data = Arc::new(RwLock::new(vec![1, 2, 3]));
    let mut handles = vec![];

    // 多个读者
    for i in 0..3 {
        let data = Arc::clone(&data);
        let handle = thread::spawn(move || {
            let r = data.read().unwrap();
            println!("读者 {} 看到: {:?}", i, *r);
            thread::sleep(Duration::from_millis(100));
        });
        handles.push(handle);
    }

    // 一个写者
    let data_clone = Arc::clone(&data);
    let writer = thread::spawn(move || {
        thread::sleep(Duration::from_millis(50));
        let mut w = data_clone.write().unwrap();
        w.push(4);
        println!("写者添加了 4");
    });
    handles.push(writer);

    for handle in handles {
        handle.join().unwrap();
    }

    println!("最终数据: {:?}", *data.read().unwrap());
}
```

## 5. 原子类型

### 5.1 基本原子操作

```rust
use std::sync::atomic::{AtomicUsize, Ordering};
use std::sync::Arc;
use std::thread;

fn main() {
    let counter = Arc::new(AtomicUsize::new(0));
    let mut handles = vec![];

    for _ in 0..10 {
        let counter = Arc::clone(&counter);
        let handle = thread::spawn(move || {
            for _ in 0..100 {
                counter.fetch_add(1, Ordering::SeqCst);
            }
        });
        handles.push(handle);
    }

    for handle in handles {
        handle.join().unwrap();
    }

    println!("结果: {}", counter.load(Ordering::SeqCst));
}
```

### 5.2 内存排序

```rust
use std::sync::atomic::{AtomicBool, AtomicI32, Ordering};
use std::sync::Arc;
use std::thread;
use std::time::Duration;

fn main() {
    let data = Arc::new(AtomicI32::new(0));
    let ready = Arc::new(AtomicBool::new(false));

    let data_clone = Arc::clone(&data);
    let ready_clone = Arc::clone(&ready);

    // 生产者线程
    let producer = thread::spawn(move || {
        data_clone.store(42, Ordering::Relaxed);
        ready_clone.store(true, Ordering::Release); // Release 语义
    });

    // 消费者线程
    let consumer = thread::spawn(move || {
        while !ready.load(Ordering::Acquire) { // Acquire 语义
            thread::sleep(Duration::from_millis(1));
        }
        let value = data.load(Ordering::Relaxed);
        println!("读取到的值: {}", value);
    });

    producer.join().unwrap();
    consumer.join().unwrap();
}
```

## 6. 条件变量和屏障

### 6.1 条件变量 (Condvar)

```rust
use std::sync::{Arc, Condvar, Mutex};
use std::thread;
use std::time::Duration;

fn main() {
    let pair = Arc::new((Mutex::new(false), Condvar::new()));
    let pair2 = Arc::clone(&pair);

    // 等待线程
    thread::spawn(move || {
        let (lock, cvar) = &*pair2;
        let mut started = lock.lock().unwrap();

        while !*started {
            println!("等待条件满足...");
            started = cvar.wait(started).unwrap();
        }

        println!("条件满足，继续执行!");
    });

    // 主线程等待一段时间后通知
    thread::sleep(Duration::from_millis(2000));

    let (lock, cvar) = &*pair;
    let mut started = lock.lock().unwrap();
    *started = true;
    cvar.notify_one();

    println!("已发送通知");

    thread::sleep(Duration::from_millis(1000));
}
```

### 6.2 屏障 (Barrier)

```rust
use std::sync::{Arc, Barrier};
use std::thread;

fn main() {
    let mut handles = Vec::with_capacity(10);
    let barrier = Arc::new(Barrier::new(10));

    for i in 0..10 {
        let c = Arc::clone(&barrier);
        handles.push(thread::spawn(move || {
            println!("线程 {} 到达屏障前", i);
            c.wait();
            println!("线程 {} 通过屏障后", i);
        }));
    }

    for handle in handles {
        handle.join().unwrap();
    }
}
```

## 7. 实践项目：并发下载器

```rust
use std::sync::{Arc, Mutex};
use std::thread;
use std::time::{Duration, Instant};
use std::collections::VecDeque;

// 模拟下载任务
#[derive(Debug, Clone)]
struct DownloadTask {
    id: usize,
    url: String,
    size: usize, // 字节数
}

// 下载结果
#[derive(Debug)]
struct DownloadResult {
    task_id: usize,
    success: bool,
    duration: Duration,
    error: Option<String>,
}

// 线程池
struct ThreadPool {
    workers: Vec<Worker>,
    sender: std::sync::mpsc::Sender<Job>,
}

type Job = Box<dyn FnOnce() + Send + 'static>;

impl ThreadPool {
    fn new(size: usize) -> ThreadPool {
        assert!(size > 0);

        let (sender, receiver) = std::sync::mpsc::channel();
        let receiver = Arc::new(Mutex::new(receiver));
        let mut workers = Vec::with_capacity(size);

        for id in 0..size {
            workers.push(Worker::new(id, Arc::clone(&receiver)));
        }

        ThreadPool { workers, sender }
    }

    fn execute<F>(&self, f: F)
    where
        F: FnOnce() + Send + 'static,
    {
        let job = Box::new(f);
        self.sender.send(job).unwrap();
    }
}

struct Worker {
    id: usize,
    thread: thread::JoinHandle<()>,
}

impl Worker {
    fn new(id: usize, receiver: Arc<Mutex<std::sync::mpsc::Receiver<Job>>>) -> Worker {
        let thread = thread::spawn(move || loop {
            let job = receiver.lock().unwrap().recv().unwrap();
            println!("工作线程 {} 获得任务; 执行中。", id);
            job();
        });

        Worker { id, thread }
    }
}

// 下载管理器
struct DownloadManager {
    thread_pool: ThreadPool,
    results: Arc<Mutex<Vec<DownloadResult>>>,
    pending_tasks: Arc<Mutex<VecDeque<DownloadTask>>>,
    active_downloads: Arc<Mutex<usize>>,
    max_concurrent: usize,
}

impl DownloadManager {
    fn new(thread_count: usize, max_concurrent: usize) -> Self {
        Self {
            thread_pool: ThreadPool::new(thread_count),
            results: Arc::new(Mutex::new(Vec::new())),
            pending_tasks: Arc::new(Mutex::new(VecDeque::new())),
            active_downloads: Arc::new(Mutex::new(0)),
            max_concurrent,
        }
    }

    fn add_task(&self, task: DownloadTask) {
        let mut pending = self.pending_tasks.lock().unwrap();
        pending.push_back(task);
    }

    fn start_downloads(&self) {
        loop {
            // 检查是否可以启动新的下载
            let active_count = *self.active_downloads.lock().unwrap();
            if active_count >= self.max_concurrent {
                thread::sleep(Duration::from_millis(100));
                continue;
            }

            // 获取待下载任务
            let task = {
                let mut pending = self.pending_tasks.lock().unwrap();
                pending.pop_front()
            };

            match task {
                Some(task) => {
                    // 增加活跃下载计数
                    {
                        let mut active = self.active_downloads.lock().unwrap();
                        *active += 1;
                    }

                    let results = Arc::clone(&self.results);
                    let active_downloads = Arc::clone(&self.active_downloads);

                    self.thread_pool.execute(move || {
                        let result = simulate_download(&task);

                        // 存储结果
                        {
                            let mut results = results.lock().unwrap();
                            results.push(result);
                        }

                        // 减少活跃下载计数
                        {
                            let mut active = active_downloads.lock().unwrap();
                            *active -= 1;
                        }
                    });
                }
                None => {
                    // 没有更多任务，检查是否所有下载都完成了
                    if *self.active_downloads.lock().unwrap() == 0 {
                        break;
                    }
                    thread::sleep(Duration::from_millis(100));
                }
            }
        }
    }

    fn get_results(&self) -> Vec<DownloadResult> {
        let results = self.results.lock().unwrap();
        results.clone()
    }

    fn get_stats(&self) -> (usize, usize, usize) {
        let results = self.results.lock().unwrap();
        let pending = self.pending_tasks.lock().unwrap();
        let active = *self.active_downloads.lock().unwrap();

        (results.len(), pending.len(), active)
    }
}

// 模拟下载函数
fn simulate_download(task: &DownloadTask) -> DownloadResult {
    let start = Instant::now();

    println!("开始下载任务 {}: {}", task.id, task.url);

    // 模拟下载时间（基于大小）
    let download_time = Duration::from_millis((task.size / 1000) as u64);
    thread::sleep(download_time);

    // 模拟 10% 的失败率
    let success = task.id % 10 != 0;

    let result = DownloadResult {
        task_id: task.id,
        success,
        duration: start.elapsed(),
        error: if !success {
            Some("网络错误".to_string())
        } else {
            None
        },
    };

    if success {
        println!("✓ 任务 {} 下载成功 ({:?})", task.id, result.duration);
    } else {
        println!("✗ 任务 {} 下载失败", task.id);
    }

    result
}

fn main() {
    println!("=== 并发下载器演示 ===\n");

    let manager = DownloadManager::new(4, 3); // 4个工作线程，最多3个并发下载

    // 添加下载任务
    let tasks = vec![
        DownloadTask { id: 1, url: "https://example.com/file1.zip".to_string(), size: 1024000 },
        DownloadTask { id: 2, url: "https://example.com/file2.zip".to_string(), size: 2048000 },
        DownloadTask { id: 3, url: "https://example.com/file3.zip".to_string(), size: 512000 },
        DownloadTask { id: 4, url: "https://example.com/file4.zip".to_string(), size: 3072000 },
        DownloadTask { id: 5, url: "https://example.com/file5.zip".to_string(), size: 1536000 },
        DownloadTask { id: 6, url: "https://example.com/file6.zip".to_string(), size: 768000 },
        DownloadTask { id: 7, url: "https://example.com/file7.zip".to_string(), size: 2560000 },
        DownloadTask { id: 8, url: "https://example.com/file8.zip".to_string(), size: 1280000 },
        DownloadTask { id: 9, url: "https://example.com/file9.zip".to_string(), size: 896000 },
        DownloadTask { id: 10, url: "https://example.com/file10.zip".to_string(), size: 1792000 },
    ];

    for task in tasks {
        manager.add_task(task);
    }

    // 创建监控线程
    let manager_clone = Arc::new(manager);
    let monitor_manager = Arc::clone(&manager_clone);

    let monitor = thread::spawn(move || {
        loop {
            let (completed, pending, active) = monitor_manager.get_stats();
            println!("状态: 已完成 {}, 待处理 {}, 进行中 {}", completed, pending, active);

            if pending == 0 && active == 0 {
                break;
            }

            thread::sleep(Duration::from_millis(500));
        }
    });

    // 开始下载
    let start_time = Instant::now();
    manager_clone.start_downloads();

    // 等待监控线程完成
    monitor.join().unwrap();

    let total_time = start_time.elapsed();

    // 显示结果
    println!("\n=== 下载完成 ===");
    println!("总耗时: {:?}", total_time);

    let results = manager_clone.get_results();
    let successful = results.iter().filter(|r| r.success).count();
    let failed = results.len() - successful;

    println!("成功: {}, 失败: {}", successful, failed);

    println!("\n详细结果:");
    for result in results {
        if result.success {
            println!("✓ 任务 {}: {:?}", result.task_id, result.duration);
        } else {
            println!("✗ 任务 {}: {:?}", result.task_id, result.error.unwrap_or_default());
        }
    }
}
```

## 8. 异步编程预览

```rust
// 注意：这需要添加 tokio 依赖
/*
[dependencies]
tokio = { version = "1", features = ["full"] }
*/

/*
use tokio::time::{sleep, Duration};

async fn say_hello() {
    println!("Hello");
    sleep(Duration::from_secs(1)).await;
    println!("World");
}

#[tokio::main]
async fn main() {
    say_hello().await;
}
*/
```

## 9. 并发编程最佳实践

### 9.1 避免死锁

```rust
use std::sync::{Arc, Mutex};
use std::thread;

fn main() {
    let data1 = Arc::new(Mutex::new(0));
    let data2 = Arc::new(Mutex::new(0));

    let data1_clone = Arc::clone(&data1);
    let data2_clone = Arc::clone(&data2);

    // 总是以相同的顺序获取锁
    let handle1 = thread::spawn(move || {
        let _lock1 = data1_clone.lock().unwrap();
        let _lock2 = data2_clone.lock().unwrap();
        println!("线程1获得了两个锁");
    });

    let handle2 = thread::spawn(move || {
        let _lock1 = data1.lock().unwrap(); // 相同的顺序
        let _lock2 = data2.lock().unwrap();
        println!("线程2获得了两个锁");
    });

    handle1.join().unwrap();
    handle2.join().unwrap();
}
```

### 9.2 选择合适的同步原语

- **Mutex**: 独占访问
- **RwLock**: 多读单写
- **Atomic**: 简单的原子操作
- **Channel**: 消息传递
- **Condvar**: 条件等待

## 10. 练习题

1. **生产者-消费者问题**：使用通道实现经典的生产者-消费者模式。

2. **并发计算器**：创建一个多线程程序来计算大数组的和。

3. **线程安全的计数器**：使用不同的同步机制实现线程安全的计数器。

4. **Web 爬虫**：实现一个简单的并发 Web 爬虫（模拟版本）。

Rust 的并发编程模型通过所有权系统提供了强大的安全保障。通过理解和实践这些概念，你将能够编写出安全、高效的并发程序。
