# auto-helper
自动签到助手

## v1.0 掘金社区自动签到，自动抽奖
```text
./build.sh
./main -h
usage: main [<flags>]

Flags:
  -h, --help              Show context-sensitive help (also try --help-long and --help-man).
      --server="tencent"  Set up server type, Default is tencent scf and main run, eg --server="tencent" or --server="run"
      --url="https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxxxxxxxxxx"  
                          Set up we conn bot URL, Default is eg
      --cookie="cookie"   Set up jue jin cookie
```
## Tencent scf使用
server 生成腾讯云函数执行二进制文件，部署至腾讯云函数即可。
```text
1.注册腾讯云并且安装Serverless Framework
2. 二进制安装：curl -o- -L https://slss.io/install | bash
3. 检查更新：serverless upgrade
4. 查看版本：serverless -v
5. 使用帮助：sls -h
6. 初始化新项目：sls init
7. 获取应用详情: sls info
8. 部署应用到云端：sls deploy
9. 移除应用：sls remove
```

### 部署实例

初始化scf项目后，目录如下：
```
project
│   README.md
│   README_EN.md 
│   main
│   serverless.yaml
```
将编译生成的scf二进制文件替换main，编辑serverless.yaml文件

```text
app: auto-helper-xxxxx
component: scf
name: scf-golang

inputs:
  name: ${app}-scf-golang
  src: ./
  handler: main
  runtime: Go1
  namespace: default
  region: ap-guangzhou
  description: '掘金社区自动签到与幸运抽奖'
  memorySize: 128
  timeout: 60
  events: # 触发器
    - timer: # 定时触发器
        name: #触发器名称，默认timer-${name}-${stage}
        parameters:
          cronExpression: '0 30 15 * * * *'
          enable: true
```

定义触发器为定时器触发即可。亦可定义不同类型的触发器。执行
```shell
sls deploy
```
将项目部署至腾讯云scf即可。

## todo jing dong 自动签到领京豆
