
# Golang-based Process Management Tool

## Basic Features
- Similar to screen's process management
- Web UI for visual operations
- Log management based on Elasticsearch
- Role-based permission management (three roles)
- Process performance monitoring
- Process state notifications

## Terminal Types
### pty
Implemented based on pseudo terminal, supporting ANSI characters and shortcuts. Only supports Linux.
### std
Implemented based on stdin and stdout pipes, supporting all platforms.

## Roles

| Role  | Role Management | Process Creation | Process Control | Log View |
| ----- | -------- | -------- | -------- | ---- |
| root  | ✔        | ✔        | ✔        | ✔    |
| admin | ×        | ×        | ✔        | ✔    |
| user  | ×        | ×        | Configurable | ×    |

## How to Use
### Starting a Process
#### Windows
Download the Windows version and double-click to run.
#### Linux
Download the Linux version.
Use the command:
```
chmod 777 ./xpm
./xpm
```
### Web Interface
Access http://[ip]:8797
Default username and password: root/root

## Interface
### Process
![image](https://github.com/lzh-1625/x_process_manager/assets/59822923/50f31b99-41d4-4d8c-88fe-20c978385155)

### Terminal
![image](https://github.com/lzh-1625/x_process_manager/assets/59822923/63eb6bec-353f-4d12-a1d9-95d89fccdac3)

### Log
![image](https://github.com/lzh-1625/x_process_manager/assets/59822923/6af8e228-7709-45c5-aba8-4b61dc825026)
