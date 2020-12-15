# PortScan

golang开发的简单的TCP端口扫描工具

+ 支持单个IP或CIDR来指定扫描目标
+ 支持单个/多个端口或端口范围的组合来指定扫描端口
+ 支持自定义goroutine的数量来控制扫描的速度
+ 支持将结果以日志的形式输出到文件

```
Author: Loong716
Github: https://github.com/loong716
Usage: PortScan [-h Help] [-t IP/CIDR] [-p Ports] [-n Threads] [-o Output]

Options:
  -h    Print the help page.
  -t    Targets, Single IP or CIDR. Ex: 192.168.2.1 or 192.168.2.0/24
  -p    Port ranges. Ex: 1-65535 or 80,443 or 80,443,100-110
  -n    The Number of Goroutine, too large value is not recommended. (Default is 20)
  -o    Output file path. Ex: /tmp/result.txt or C:\Windows\Temp\result.txt
```

## Usage

`-h`

打印帮助页面

`-t`

指定目标，支持单个IP或CIDR。例如: `-t 192.168.2.1` 或 `-t 192.168.2.0/24`

`-p`

指定端口，支持单个端口、多个端口（以`,`分割）、端口范围（以`-`表示范围）。例如: `-p 1-65535` 或 `-p 80,443` 或 `-p 80,443,100-110`

`-n`

创建goroutine的数量，数量越大速度越快，但不建议使用数量过大的值。（默认为20，测试扫描一个c段的机器，每台机器10个端口，耗时约20s）

`-o`

输出文件的路径。例如：`-o /tmp/result.txt` 或 `-o C:\Windows\Temp\result.txt`

## Build

Windows：

```
# x64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build PortScan.go

# x86
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build PortScan.go
```

Linux:

```
# x64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build PortScan.go

# x86
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build PortScan.go
```

Mac OS:

```
# x64
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build PortScan.go

# x86
CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build PortScan.go
```
