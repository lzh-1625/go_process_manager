<div align="center">
  <br/>
  <img src="https://raw.githubusercontent.com/lzh-1625/go_process_manager/refs/heads/master/resources/templates/favicon.ico" alt="GPM logo">
  <br/>
</div>


# Go Process Manager

Go Process Manager is a process management tool developed based on Golang, offering process management capabilities similar to `screen`, and supports Web UI for visual operations. This tool integrates log management, permission control, process performance monitoring, status push, and more. It supports CPU and memory limits through `cgroup` and also provides task management and scheduled task functionalities.

---

## [中文](./README_CN.md)

## Features

### Core Features

- **Process Management**: Supports creating, starting, stopping, restarting, and deleting processes, providing multi-process management capabilities similar to `screen`.
- **Web UI for Visual Operations**: Offers a user-friendly web interface, allowing users to easily manage processes, view logs, and monitor performance through a browser.
- **Log Management**: Supports storing logs in Elasticsearch 8 or SQLite, providing efficient log query and analysis capabilities.
- **Permission Management**: Provides role-based access control with three roles (root, admin, user), ensuring system security and flexibility.
- **Process Performance Monitoring**: Real-time monitoring of CPU, memory, and other resource usage, helping users optimize system performance.
- **Process Status Push**: Supports real-time status push for processes, allowing users to receive the latest process status via Web UI or API.
- **Resource Limits**: Implements CPU and memory limits through `cgroup`, preventing processes from consuming excessive system resources.
- **Task Management**: Supports task management, scheduled tasks, and API calls, meeting the needs of automated operations.

---

### Terminal Types

- **pty (Pseudo Terminal)**: Based on pseudo terminals, supports ANSI characters and shortcut key operations, suitable for interactive command-line programs.
- **std (Standard Input/Output)**: Based on standard input/output pipes, suitable for non-interactive programs or scripts.

---

## Role Permissions

| Role  | Role Management | Process Creation | Process Control | Log Viewing | Task Management |
| ----- | --------------- | ---------------- | --------------- | ----------- | --------------- |
| root  | ✔              | ✔               | ✔              | ✔          | ✔              |
| admin | ×              | ×               | ✔              | ✔          | ✔              |
| user  | ×              | ×               | Custom          | Custom      | ×              |

- **root**: Has the highest permissions, can manage all processes, logs, and user roles.
- **admin**: Can control processes and view logs but cannot create processes or manage roles.
- **user**: Permissions can be customized, suitable for regular users.

---

## User Guide

### Starting the Process

#### Windows

1. Download the Windows version of the binary file.
2. Double-click to run and start the service.

#### Linux

1. Download the Linux version of the binary file.
2. Use the following commands to grant execution permissions and start:

   ```bash
   chmod 777 ./go_process_manager
   ./go_process_manager
   ```

### Web Interface

1. After starting the service, access `http://[ip]:8797`.
2. Log in with the default credentials `root/root`.

### Demo

Visit the [Demo](http://xcon.top:9787/process) to experience the features. Use the credentials `root/root` to log in.

---

## Interface Showcase

### Process Management

![Process Management Interface](https://github.com/lzh-1625/go_process_manager/assets/59822923/50f31b99-41d4-4d8c-88fe-20c978385155)

- **Process List**: Displays all running processes, including process ID, name, status, resource usage, etc.
- **Action Buttons**: Supports starting, stopping, restarting, and deleting processes.

### Terminal Operations

![Terminal Operations Interface](https://github.com/lzh-1625/go_process_manager/assets/59822923/63eb6bec-353f-4d12-a1d9-95d89fccdac3)

- **Terminal Emulation**: Supports ANSI characters and shortcut key operations, providing an experience similar to a local terminal.
- **Input/Output**: Real-time display of process standard input and output.

### Log Viewing

![Log Viewing Interface](https://github.com/lzh-1625/go_process_manager/assets/59822923/6af8e228-7709-45c5-aba8-4b61dc825026)

- **Log Query**: Supports filtering logs by time, process ID, operator, etc.

---

## Additional Notes

### Log Management

- **Elasticsearch 8**: Suitable for large-scale log storage and query, supports distributed deployment.
- **SQLite**: Lightweight log storage, suitable for single-machine or small-scale use.

### Permission Control

- **root User**: Has the highest permissions, can manage all processes, logs, and user roles.
- **admin User**: Can control processes and view logs, suitable for operations personnel.
- **user User**: Permissions can be customized, suitable for regular users or developers.

### Resource Limits

- **CPU Limits**: Set CPU usage limits for processes through `cgroup`.
- **Memory Limits**: Set memory usage limits for processes through `cgroup`.

### Task Management

- **Scheduled Tasks**: Supports Cron expressions, allowing users to create scheduled tasks.
- **API Calls**: Provides APIs to trigger tasks.
- **Task Flow**: Chain execution of tasks.
- **Trigger Events**: Trigger tasks through process stop, start, or exceptions.

---

## Development and Deployment

### Environment Requirements

- **Golang**: Version 1.18 or higher.
- **Elasticsearch 8** (optional): For log storage.
- **SQLite** (optional): For lightweight log storage.

### Compilation and Execution

1. Clone the project:

   ```bash
   git clone https://github.com/lzh-1625/go_process_manager.git
   cd go_process_manager
   ```
2. Compile the project:

   ```bash
   go build -o go_process_manager
   ```
3. Run the project:

   ```bash
   ./go_process_manager
   ```

---

## Contributions and Feedback

We welcome submitting Issues and Pull Requests to help improve Go Process Manager. If you have any questions or suggestions, please provide feedback via [GitHub Issues](https://github.com/lzh-1625/go_process_manager/issues).

---

## License

This project is licensed under the [MIT License](https://opensource.org/licenses/MIT). For more details, please refer to the [LICENSE](LICENSE) file.

---

## Contact Us

- **Author**: lzh-1625
- **GitHub**: [go_process_manager](https://github.com/lzh-1625/go_process_manager)

Thank you for using Go Process Manager! We hope this tool brings convenience to your process management.
