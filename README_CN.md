<div align="center">
  <br/>
  <img src="https://raw.githubusercontent.com/lzh-1625/go_process_manager/refs/heads/master/resources/templates/favicon.ico" alt="GPM logo">
  <br/>
</div>

# Go Process Manager

Go Process Manager 是一个基于 Golang 开发的进程管理工具，提供了类似 `screen` 的进程管理功能，并支持 Web UI 可视化操作。该工具集成了日志管理、权限控制、进程性能监控、状态推送等功能，支持通过 `cgroup` 实现 CPU 和内存限制，同时还提供了任务管理和定时任务功能。

---

## 功能特性

### 核心功能

- **进程管理**支持进程的创建、启动、停止、重启和删除操作，提供类似 `screen` 的多进程管理能力。
- **Web UI 可视化操作**提供友好的 Web 界面，用户可以通过浏览器轻松管理进程、查看日志和监控性能。
- **日志管理**支持将日志存储到 Elasticsearch 8 或 SQLite 中，提供高效的日志查询和分析功能。
- **权限管理**提供三种角色（root、admin、user）的权限控制，确保系统安全性和灵活性。
- **进程性能监控**实时监控进程的 CPU、内存等资源使用情况，帮助用户优化系统性能。
- **进程状态推送**支持进程状态的实时推送，用户可以通过 Web UI 或 API 推送进程的最新状态。
- **资源限制**通过 `cgroup` 实现 CPU 和内存的限制，防止进程占用过多系统资源。
- **任务管理**
  支持任务管理、定时任务和 API 调用，满足自动化运维需求。

---

### 终端类型

- **pty（伪终端）**基于伪终端实现，支持 ANSI 字符和快捷键操作，适合交互式命令行程序。
- **std（标准输入输出）**
  基于标准输入输出管道实现，适合非交互式程序或脚本。

---

## 角色权限

| 角色  | 角色管理 | 进程创建 | 操控进程 | 日志查看 | 任务管理 |
| ----- | -------- | -------- | -------- | -------- | -------- |
| root  | ✔       | ✔       | ✔       | ✔       | ✔       |
| admin | ×       | ×       | ✔       | ✔       | ✔       |
| user  | ×       | ×       | 自定义   | 自定义   | ×       |

- **root**：拥有最高权限，可以管理所有进程、日志和用户角色。
- **admin**：可以操控进程和查看日志，但不能创建进程或管理角色。
- **user**：权限可自定义，适合普通用户使用。

---

## 使用指南

### 启动进程

#### Windows

1. 下载 Windows 版本的二进制文件。
2. 双击运行即可启动服务。

#### Linux

1. 下载 Linux 版本的二进制文件。
2. 使用以下命令赋予执行权限并启动：

   ```bash
   chmod 777 ./go_process_manager
   ./go_process_manager
   ```

### Web 界面

1. 启动服务后，访问 `http://[ip]:8797`。
2. 使用默认账号密码 `root/root` 登录。

### Demo 演示

访问 [Demo 演示](http://xcon.top:9787/process) 体验功能，使用账号 `root/root` 登录。

---

## 界面展示

### 进程管理

![进程管理界面](https://github.com/lzh-1625/go_process_manager/assets/59822923/50f31b99-41d4-4d8c-88fe-20c978385155)

- **进程列表**：显示所有运行的进程，包括进程 ID、名称、状态、资源使用情况等。
- **操作按钮**：支持启动、停止、重启和删除进程。

### 终端操作

![终端操作界面](https://github.com/lzh-1625/go_process_manager/assets/59822923/63eb6bec-353f-4d12-a1d9-95d89fccdac3)

- **终端模拟**：支持 ANSI 字符和快捷键操作，提供类似本地终端的体验。
- **输入输出**：实时显示进程的标准输入和输出。

### 日志查看

![日志查看界面](https://github.com/lzh-1625/go_process_manager/assets/59822923/6af8e228-7709-45c5-aba8-4b61dc825026)

- **日志查询**：支持按时间、进程 ID、操作用户等条件过滤日志。

---

## 补充说明

### 日志管理

- **Elasticsearch 8**：适合大规模日志存储和查询，支持分布式部署。
- **SQLite**：轻量级日志存储，适合单机或小规模使用。

### 权限控制

- **root 用户**：拥有最高权限，可以管理所有进程、日志和用户角色。
- **admin 用户**：可以操控进程和查看日志，适合运维人员使用。
- **user 用户**：权限可自定义，适合普通用户或开发人员使用。

### 资源限制

- **CPU 限制**：通过 `cgroup` 设置进程的 CPU 使用上限。
- **内存限制**：通过 `cgroup` 设置进程的内存使用上限。

### 任务管理

- **定时任务**：支持 Cron 表达式，用户可以创建定时任务。
- **API 调用**：提供 API触发任务。
- **任务流**：任务的链式执行。
- **触发事件**：通过进程的停止、启动、异常触发任务。

---

## 开发与部署

### 环境要求

- **Golang**：版本 1.18 或以上。
- **Elasticsearch 8**（可选）：用于日志存储。
- **SQLite**（可选）：用于轻量级日志存储。

### 编译与运行

1. 克隆项目：

   ```bash
   git clone https://github.com/lzh-1625/go_process_manager.git
   cd go_process_manager
   ```
2. 编译项目：

   ```bash
   go build -o go_process_manager
   ```
3. 运行项目：

   ```bash
   ./go_process_manager
   ```

---

## 贡献与反馈

欢迎提交 Issue 和 Pull Request，帮助我们改进 Go Process Manager。如果有任何问题或建议，请通过 [GitHub Issues](https://github.com/lzh-1625/go_process_manager/issues) 反馈。

---

## 许可证

本项目采用 [MIT 许可证](https://opensource.org/licenses/MIT)，详情请参阅 [LICENSE](LICENSE) 文件。

---

## 联系我们

- **作者**：lzh-1625
- **GitHub**：[go_process_manager](https://github.com/lzh-1625/go_process_manager)

感谢您使用 Go Process Manager！希望这个工具能为您的进程管理带来便利。
