基于golang的进程管理工具

# 基础功能
- 类似screen的进程管理
- 提供web ui可视化操作
- 基于Elasticsearch8或sqlite的日志管理
- 三种角色的权限管理
- 进程性能监控
- 进程状态推送

# 终端类型
## pty 
基于伪终端实现，支持ansi字符、快捷键，仅支持linux
## std
基于stdin、stdout管道实现，支持所有平台

# 角色

| 角色  | 角色管理 | 进程创建 | 操控进程 | 日志查看 |
| ----- | -------- | -------- | -------- | ---- |
| root  | ✔        | ✔        | ✔        | ✔    |
| admin | ×        | ×        | ✔        | ✔    |
| user  | ×        | ×        | 自定义  | 自定义    |


# 如何使用
## 启动进程
### windows
下载windows版本双击运行
### linux
下载linux版本
使用命令
```
chmod 777 ./xpm
./xpm
```
## web界面
访问http://[ip]:8797
默认账号密码 root/root

## demo演示
http://xcon.top:9787/process
root/root

# 界面
### 进程
![image](https://github.com/lzh-1625/x_process_manager/assets/59822923/50f31b99-41d4-4d8c-88fe-20c978385155)

### 终端
![image](https://github.com/lzh-1625/x_process_manager/assets/59822923/63eb6bec-353f-4d12-a1d9-95d89fccdac3)

### 日志
![image](https://github.com/lzh-1625/x_process_manager/assets/59822923/6af8e228-7709-45c5-aba8-4b61dc825026)

### 监控
cpu 内存 水位线

### cgroup

### 定时任务