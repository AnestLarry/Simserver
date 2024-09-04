package Libs

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

/////////////////////////////////////////////////////////////////////////////////////////////////////
//  20231
/////////////////////////////////////////////////////////////////////////////////////////////////////

func LibsXTp[T any](cond bool, t T, f T) T {
	if cond {
		return t
	} else {
		return f
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
//  20231
/////////////////////////////////////////////////////////////////////////////////////////////////////

func VecContains[T comparable](array []T, val T) (index int) {
	index = -1
	for i := 0; i < len(array); i++ {
		if array[i] == val {
			index = i
			return
		}
	}
	return
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
//  20203
/////////////////////////////////////////////////////////////////////////////////////////////////////

func SendHTTPRequest(method, aUrl string, headers, data map[string]string, timeOut uint16,
	proxy string) ([]byte, []*http.Cookie, *http.Response) {
	var aClient *http.Client
	if proxy != "" {
		proxyUrl, _ := url.Parse(proxy)
		transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
		aClient = &http.Client{
			Timeout:   time.Duration(timeOut) * time.Millisecond,
			Transport: transport,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}}
	} else {
		aClient = &http.Client{
			Timeout: time.Duration(timeOut) * time.Millisecond,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}}
	}

	payload := make(url.Values)
	if data != nil {
		for i, j := range data {
			payload.Add(i, j)
		}
	}
	body := strings.NewReader(payload.Encode())
	req, _ := http.NewRequest(
		method,
		aUrl,
		body,
	)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.117 Safari/537.36")
	if headers != nil {
		for i, v := range headers {
			req.Header.Add(i, v)
		}
	}
	if method == "POST" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	r, err := aClient.Do(req)
	if err != nil {
		panic(err)
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	cookie := req.Cookies()
	return b, cookie, r
}

func SendGet(aUrl string, headers map[string]string) ([]byte, []*http.Cookie, *http.Response) {
	return SendHTTPRequest("GET", aUrl, nil, headers, 60, "")
}
func SendGetSimple(aUrl string) ([]byte, []*http.Cookie, *http.Response) {
	return SendGet(aUrl, nil)
}
func SendPost(aUrl string, headers, data map[string]string) ([]byte, []*http.Cookie, *http.Response) {
	return SendHTTPRequest("POST", aUrl, data, headers, 60, "")
}
func SendPostSimple(aUrl string, data map[string]string) ([]byte, []*http.Cookie, *http.Response) {
	return SendPost(aUrl, nil, data)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
//  20201
/////////////////////////////////////////////////////////////////////////////////////////////////////

func LibsXClear__20201() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "darwin":
		fallthrough
	case "linux":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
//  History
/////////////////////////////////////////////////////////////////////////////////////////////////////

func LibsXRangeInt(args ...int) chan int {
	if l := len(args); l < 1 || l > 3 {
		fmt.Println("error args length, xRangeInt requires 1-3 int arguments")
	}
	var start, stop int
	var step int = 1
	switch len(args) {
	case 1:
		stop = args[0]
		start = 0
	case 2:
		start, stop = args[0], args[1]
	case 3:
		start, stop, step = args[0], args[1], args[2]
	}

	ch := make(chan int)
	go func() {
		if step > 0 {
			for start < stop {
				ch <- start
				start = start + step
			}
		} else {
			for start > stop {
				ch <- start
				start = start + step
			}
		}
		close(ch)
	}()

	return ch
}

func LibsXClear() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// 判断所给路径文件/文件夹是否存在
func LibsXExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件夹
func LibsXIsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func LibsXIsFile(path string) bool {
	return !LibsXIsDir(path)
}

func LibsXSha1File(filePath string) []byte {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return h.Sum(nil)
}

func LibsXSha1FileString(filePath string) string {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(h.Sum(nil))
}

func LibsXexecCommand(commandName string, params []string) bool {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command(commandName, params...)

	//显示运行的命令
	//fmt.Println(cmd.Args)
	//StdoutPipe方法返回一个在命令Start后与命令标准输出关联的管道。Wait方法获知命令结束后会关闭这个管道，一般不需要显式的关闭该管道。
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		return false
	}

	cmd.Start()
	//创建一个流来读取管道内内容，这里逻辑是通过一行一行的读取的
	reader := bufio.NewReader(stdout)

	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(line)
	}

	//阻塞直到该命令执行完成，该命令必须是被Start方法开始执行的
	cmd.Wait()
	return true
}

func LibsXExecShell(commandName string, params []string) {
	// stdout,stderr
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command(commandName, params...)

	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	w := bytes.NewBuffer(nil)
	cmd.Stderr = w
	cmd.Stdout = &out

	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	cmd.Run()

	outs, ws := out.String(), w.String()
	if ws != "" {
		fmt.Println(ws)
	} else {
		fmt.Println(outs)
	}
}
