package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 接收的参数
var (
	help bool
	targets string
	portRanges string
	numOfgoroutine int
	file string
)


// Usage
func usage() {
	fmt.Fprintf(
		os.Stderr, `Author: Loong716
Github: https://github.com/loong716
Usage: PortScan [-h Help] [-t IP/CIDR] [-p Ports] [-n Threads] [-o Output]

Options:
`)
	flagSet := flag.CommandLine
	order := []string{"h", "t", "p", "n", "o"}
	for _, name := range order {
		fl4g := flagSet.Lookup(name)
		fmt.Printf("  -%s  ", fl4g.Name)
		fmt.Printf("  %s\n", fl4g.Usage)
	}
}


// 初始化参数
func init() {
	flag.BoolVar(&help, "h", false, "Print the help page.")
	flag.StringVar(&file, "o", "", "Output file path. Ex: /tmp/result.txt or C:\\Windows\\Temp\\result.txt")
	flag.StringVar(&targets, "t", "", "Targets, Single IP or CIDR. Ex: 192.168.2.1 or 192.168.2.0/24")
	flag.StringVar(&portRanges, "p", "", "Port ranges. Ex: 1-65535 or 80,443 or 80,443,100-110")
	flag.IntVar(&numOfgoroutine, "n", 20, "The Number of Goroutine, too large value is not recommended. (Default is 20)")

	// 替换默认的Usage
	flag.Usage = usage
}


func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}


// 将CIDR转换为IP，若用户输入为单个IP直接返回
func CIDR2IP(cidr string) ([]string, error) {
	var hosts []string
	if !strings.ContainsAny(cidr, "/") {
		hosts = append(hosts, cidr)
		return hosts, nil
	}

	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		fmt.Println("Error: Invalid CIDR! Please use '-h' to see usage.")
		return nil, err
	}

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		hosts = append(hosts, ip.String())
	}
	return hosts[1 : len(hosts)-1], nil
}


// 解析用户输入的端口范围
func ParsePortRange(portList string) ([]int, error) {
	var ports []int
	// 先以','分割
	portList2 := strings.Split(portList, ",")

	for _, i := range portList2 {
		// 如果存在'-'，则解析范围
		if strings.ContainsAny(i, "-") {
			a := strings.Split(i, "-")
			startPort, err := strconv.Atoi(a[0])
			if err != nil {
				fmt.Println("Error: StartPort strconv error! Please use '-h' to see usage.")
				os.Exit(1)
			}
			endPort, err := strconv.Atoi(a[1])
			if err != nil {
				fmt.Println("Error: EndPort strconv error! Please use '-h' to see usage.")
				os.Exit(1)
			}
			for j := startPort; j <= endPort; j++ {
				ports = append(ports, j)
			}
		} else {
			// 不存在'-'，直接加入ports
			singlePort, err := strconv.Atoi(i)
			if err != nil {
				fmt.Println("Error: SinglePort strconv error! Please use '-h' to see usage.")
				os.Exit(1)
			}
			ports = append(ports, singlePort)
		}
	}
	return ports, nil
}


// 探测端口是否开放
func isOpen(target string) bool {
	conn, err := net.DialTimeout("tcp", target, time.Millisecond*time.Duration(200))
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}


func portScan(hosts []string, ports []int) int {
	wg := sync.WaitGroup{}
	targetsChan := make(chan string, 100)
	poolCount := numOfgoroutine  // goroutine的数量

	// 创建goroutine
	for i := 0; i <= poolCount; i++ {
		go func() {
			for j := range targetsChan {
				openFlag := isOpen(j)
				if openFlag {
					if file != ""{
						log.Printf("   %s is open!\r\n", j)
					} else {
						fmt.Printf("%s is open!\n", j)
					}
				}
				wg.Done()
			}
		}()
	}

	// 将ip与port拼接为target，通过channel传输给goroutine
	for _, m := range ports {
		portString := strconv.Itoa(m)
		for _, n := range hosts {
			target := n + ":" + portString
			targetsChan <- target
			wg.Add(1)
		}
	}

	wg.Wait()
	return 0
}


func main() {
	startTime := time.Now()
	flag.Parse()

	// 没有参数输入或输入'-h'时输出Usage
	if flag.NFlag() == 0 || help {
		flag.Usage()
		os.Exit(1)
	}

	hosts, err := CIDR2IP(targets)
	if err != nil {
		fmt.Println("Error: CIDR conversion failed! Please use '-h' to see usage.")
		os.Exit(1)
	}


	ports, err := ParsePortRange(portRanges)
	if err != nil {
		fmt.Println("Error: ParsePortRange failed! Please use '-h' to see usage.")
		os.Exit(1)
	}

	// 如果设置了'-o'参数，则使用log模块输出到文件
	if file != "" {
		logFile, err := os.OpenFile(file, os.O_RDWR | os.O_CREATE | os.O_APPEND, os.ModeAppend | os.ModePerm)
		if err != nil {
			fmt.Println("Error: Open output file failed!")
			os.Exit(1)
		}

		defer logFile.Close()
		out := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(out)
		fmt.Printf("[*] Port Scan Start...\n\n")
		portScan(hosts, ports)
		spendTime := time.Since(startTime).Seconds()
		fmt.Printf("\n[*] Finished. Take %.4fs.\n", spendTime)
	} else {
		fmt.Printf("[*] Port Scan Start...\n\n")
		portScan(hosts, ports)
		spendTime := time.Since(startTime).Seconds()
		fmt.Printf("\n[*] Finished. Take %.4fs.\n", spendTime)
	}

}