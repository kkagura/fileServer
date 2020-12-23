package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

var (
	startPort int
	silent    bool
	help      bool
)

func init() {
	flag.IntVar(&startPort, "p", 9001, "传入默认的端口号")
	flag.BoolVar(&silent, "s", false, "禁止自动打开浏览器")
	flag.BoolVar(&help, "h", false, "获取参数列表")
}

func main() {
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Print(err)
		return
	}

	httpDir := http.FileServer(http.Dir(dir))
	http.Handle("/", httpDir)

	ch := make(chan string)
	ch2 := make(chan string)
	go serve(startPort, ch, ch2)
	port := <-ch
	url := "http://127.0.0.1:" + port
	fmt.Print("\nserver running at ", url)
	if !silent {
		cmd := exec.Command("cmd", "/c", "start "+url)
		cmd.Run()
	}
	<-ch2
}

func serve(port int, ch chan string, ch2 chan string) {
	port = findPort(port)
	portStr := strconv.Itoa(port)
	ch <- portStr
	http.ListenAndServe(":"+portStr, nil)
	ch2 <- "1"
}

func findPort(portNumber int) int {
	if !isPortInUse(portNumber) {
		return portNumber
	}
	fmt.Print("\n端口", portNumber, "已被占用，正在查找下一可用端口")
	portNumber++
	return findPort(portNumber)
}

func isPortInUse(portNumber int) bool {
	con, _ := net.DialTimeout("tcp", net.JoinHostPort("", strconv.Itoa(portNumber)), time.Second)
	if con != nil {
		con.Close()
		return true
	}
	return false
}
