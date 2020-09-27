## 介绍

前段时间看到了 facebook 正式开源了一个go 的orm，想着是终于有一个大厂做背书的orm了，
所以通过该项目进行一个使用，通过gin + entORM 实现一个博客系统，目前刚刚写完博客管理后端的代码。

项目目录规范参考于：https://github.com/golang-standards/project-layout

本地配置Commitizen，本地安装之后，commit可以通过 即可只用git cz 添加commit信息

## 运行

安装数据库并创建数据库goBlog，在项目目录下执行make 即可以运行项目
登陆 127.0.0.1:8989/admin/login