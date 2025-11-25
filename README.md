````markdown
# GoCannon: High-Concurrency Load Testing CLI 💣

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=flat&logo=go&logoColor=white)
![CLI](https://img.shields.io/badge/type-CLI-orange)
![Concurrency](https://img.shields.io/badge/pattern-Worker%20Pool-blueviolet)
![Performance](https://img.shields.io/badge/performance-High-success)

**GoCannon** is a lightweight, high-performance HTTP load testing tool built from scratch using **Golang**.

Unlike standard scripts that spawn threads uncontrollably, GoCannon implements the **Worker Pool Pattern** to manage concurrency efficiently, allowing it to simulate thousands of requests with minimal memory footprint.

---

## 🧠 Why I Built This?

As an aspiring **SDET / Backend Engineer**, I wanted to move beyond just *using* tools like k6 or JMeter. I wanted to understand **how they work under the hood**.

This project demonstrates my understanding of:
* **Golang Concurrency:** Using Goroutines and Channels instead of heavy OS threads.
* **Synchronization:** Using `sync.WaitGroup` to coordinate parallel execution.
* **Resource Management:** Implementing a Worker Pool to throttle traffic and prevent resource exhaustion.

---

## ⚙️ Architecture (The Worker Pool)

Instead of creating 1,000 threads for 1,000 requests, GoCannon creates a fixed number of "Workers" (e.g., 50) that process a queue of jobs.

```mermaid
graph TD
    A[User Input] -->|Generate Jobs| B(Jobs Channel)
    B --> C{Worker Pool}
    C -->|Worker 1| D[HTTP Request]
    C -->|Worker 2| D
    C -->|Worker N| D
    D -->|Response Data| E(Results Channel)
    E --> F[Aggregator / Main Thread]
    F --> G[Final Report]
````

-----

## 🛠 Tech Stack

  * **Language:** Go (Golang) 1.21
  * **Dependencies:** None (Zero dependencies). Uses pure Go Standard Library:
      * `net/http` for networking.
      * `sync` for WaitGroups.
      * `time` for benchmarking.
      * `flag` for CLI argument parsing.

-----

## 🚀 How to Use

### 1\. Build the Binary

Compiles the Go code into an executable file.

```bash
go build -o gocannon
```

### 2\. Run a Load Test

**Basic Syntax:**

```bash
./gocannon --url="<TARGET_URL>" --requests=<TOTAL> --concurrency=<WORKERS>
```

**Example:**
Simulate **1,000 requests** with **50 concurrent users** hitting a local server.

```bash
./gocannon --url="http://localhost:8080/" --requests=1000 --concurrency=50
```

-----

## 📊 Sample Output

```text
🚀 Attacking: http://localhost:8080/
📦 Total Requests: 1000 | 🔥 Concurrency: 50
---------------------------------------------------

🏁 FINISHED!
⏱️  Total Duration    : 150.4ms
✅ Success Requests  : 1000
❌ Failed Requests   : 0
📊 Average Latency   : 4.2ms
🚀 Throughput (RPS)  : 6648.21 req/sec
```

-----

## 📝 Key Features

  * **Concurrency Control:** Strictly limits the number of active connections using the Worker Pool pattern.
  * **Thread-Safe Communication:** Uses Go Channels to pass data between workers and the main thread without locks (Mutex).
  * **Detailed Reporting:** Calculates Requests Per Second (RPS), Average Latency, and Success/Error rates.

-----

## 👤 Author

**Bryan Chan**
