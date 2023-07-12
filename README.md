# A Light-Weight GRPC Implementation

本仓库是一个轻量级的 GRPC 实现，主要用于学习 GRPC 的实现原理。只用了 1.4k 行左右的代码，实现了 GRPC 的核心功能，包括：

1. 支持 Json 和 Gob 两种消息编码方式
2. 支持并发和读写分离的服务端, 未实现工作池（业务工作池的实现可以参考 [zinx](https://github.com/SakuraILU/zinx)）
3. 支持异步和并发调用的高性能客户端
4. 通过反射实现了服务的自动注册
5. 客户端链接和服务端 Handle 的超时控制
6. 支持 RANDOM 和 ROUND_ROBIN 两种负载均衡策略
7. 一个基于 HTTP 协议的注册中心，支持服务的注册、发现和心跳保活

Main 包中的代码是一个简单的示例，包括了服务端（提供 StringReverse 和 Sort 两种服务）和客户端（多 goroutine 并发 Call RPC 服务），以及一个注册中心。
