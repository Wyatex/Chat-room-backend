# 一个简单的聊天室

![CI](https://travis-ci.com/alienyan/Chat-room-backend.svg?branch=master)
![issues](https://img.shields.io/github/issues/alienyan/Chat-room-backend)
![license](https://img.shields.io/github/license/alienyan/Chat-room-backend)

# 简介
一个可以前后端分离的简单聊天室，如果有什么好的想法可以提交issues/PR。

前端：（还未完成）

# 使用方法
## GF工具链
使用gf工具链可以更方便的编译

1. 安装GF-CLI，方法见：[https://github.com/gogf/gf-cli](https://github.com/gogf/gf-cli)

2. 修改config/config.toml下的compiler，删除其他的平台（比如想编译linux的话就删掉windows），修改编译出来的目录和二进制文件名称。（这步可以不做）

3. 执行gf build gf-app，得到所需的二进制文件

4. 执行

（或者第三步直接使用gf run main.go）

## 原生编译器
直接使用go build gf-app或者go run gf-app，或者gf-app改成main.go也可以。

