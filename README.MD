# 一个简单的聊天室

![CI](https://travis-ci.com/alienyan/Chat-room-backend.svg?branch=master)
![issues](https://img.shields.io/github/issues/alienyan/Chat-room-backend)
![license](https://img.shields.io/github/license/alienyan/Chat-room-backend)

# 简介
一个可以前后端分离的简单聊天室，如果有什么好的想法可以提交issues/PR。

[前端](https://github.com/mr-chen01/talk)

# 使用方法
## 修改配置
修改config文件夹下的config.toml配置文件，如果需要用到密码找回就需要配置你的smtp（目前只支持邮箱找回）

## GF工具链
使用gf工具链可以更方便的编译

1. 安装GF-CLI，方法见：[https://github.com/gogf/gf-cli](https://github.com/gogf/gf-cli)

2. 修改config/config.toml下的compiler，删除其他的平台（比如想编译linux的话就删掉windows），修改编译出来的目录和二进制文件名称。（这步可以不做）

3. 执行gf build gf-app，得到所需的二进制文件

4. 执行

（或者第三步直接使用gf run main.go）

## 原生编译器
直接使用go build gf-app或者go run gf-app，或者gf-app改成main.go也可以。

# 部署方法
## docker部署
使用以下指令可以静态编译Linux平台amd64架构的可执行文件：
```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go
```
然后生成我们需要的docker镜像：
```
docker build -t chat-room-back .
```
然后可以推送到远程的registry服务器（可以查百度），最后在服务器上拉取到这个镜像，通过以下代码就可以运行本服务：
```
docker run chat-room-back
```

## 独立部署
### 1. `screen`
Screen 是一款由 GNU 计划开发的用于命令行终端切换的自由软件。用户可以通过该软件同时连接多个本地或远程的命令行会话，并在其间自由切换。这个是我最常用的一个工具。使用方法如下：
1. `screen -S chatroom` 创建一个叫name的session
2. 在新终端输入 `./main` 即可
3. 执行 `ctrl-a, ctrl-d` 暂时离开当前session;
4. 当你离开后想再次进入，可以输入：`screen -r chatroom` 进入
5. 执行 `screen -X -S chatroom quit` 结束程序;

### 2. `nohup`

我们可以使用简单的nohup命令来运行应用程序，使其作为后台守护进程运行，即使远程连接的SSH断开也不会影响程序的执行。在流行的Linux发行版中往往都默认安装好了nohup命令工具。 命令如下：
```
nohup ./main &
```
### 3. `systemctl`
Systemd 是 Linux 系统工具，用来启动守护进程，已成为大多数发行版的标准配置。
而 systemctl 是 Systemd 的主命令，用于管理系统。可以参考阮一峰对于 Systemd 的解读，文章的第四、五章节：Unit
其实我们大部分服务都有使用 systemctl 管理，比如 MySQL、Nginx 等等。
常见配置如下：
```
[Unit]
# 单元描述
Description=Chat Room Backend
# 在什么服务启动之后再执行本程序
After=mysql.service

[Service]
Type=simple
# 程序执行的目录
WorkingDirectory=/data/server/chatroomback/
# 启动的脚本命令
ExecStart=/data/server/chatroomback/main
# 重启条件
Restart=alway
# 几秒后重启
RestartSec=5

[Install]
WantedBy=multi-user.target
```
**使用方法**
1. 创建应用配置文件 /etc/systemd/system/chatroom.service, 内容如上;
2. 使用 systemctl daemon-reload 重新加载服务;
3. 执行 systemctl start chatroom 来启动服务;
4. 最后执行 systemctl status chatroom 来查看服务运行的状态信息;
5. 执行 systemctl enable chatroom 将服务添加到开机启动项;
6. 注意：执行的 chatroom 是使用文件名作为服务名;
7. 常见的命令有: start(启动), stop(停止), restart(重启), status(查看运行状态), enable(添加到开机启动项), disable(将程序从开机启动中移除)