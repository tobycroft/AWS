# AWS
AWS is a simple transfer layer,It is a copy version of AWS swoole version

# AWS是什么
AWS是一个转发器，全称AWebsocket，这是AIM（1.0）的前置程序，通过AWS可以实现功能层热更新等功能

# 为什么立项
因为之前的Websocket中间件使用的是PHP+Swoole技术方案，
在（PHP）方案初期设计的时候，就是设计成热插拔的形式的，后端TP支撑，
那么最近在一个人多的程序里面，就经常出现这个中间件卡死的问题，
用户可以连入，但是就是无法完成鉴权等验证，必须重启swoole才可以，
重点是PHP程序里面没有trycatch（我写程序会尽量避免trycatch），
如果出现问题会直接报错，但是swoole那边运行依旧很稳定没有任何报错，也没有任何鉴权动作产生，
排除了问题后，于是这次用Go做一个AWS前置中间件

# 性能&稳定性
如果按照swoole官方的说法，Go程序性能不如swoole，这里不讨论了

稳定性方面主要是热更新程序的时候，所有连入端都不会掉，无闪断，后端程序可以平顺升级


# AWS vs AIM 架构设计

这次的架构方案还是沿用MChat（PHP版）v2的架构设计，因为Go特性，所以在使用协程后，大群效能可以提高1200%左右