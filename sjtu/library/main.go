package main

import (
	"os"
	"sync"
	"fmt"
	"net/http"
	"flag"
	"lib/utils"
	"strings"
	"encoding/json"
	"time"
)

// 需要搜索的内容
var request string
var filename string 

func init() {
	flag.StringVar(&request, "req", "", "搜索关键字")
	flag.StringVar(&filename, "filename", "lib.json", "保存文件名")
}

func main() {
	flag.Parse()
	if strings.TrimSpace(request) == "" {
		fmt.Println("请指定搜索关键字")
		return
	}
	start := time.Now()
	var client http.Client
	req_url, err := utils.GetUrl(client)
	if err != nil {
		fmt.Println("Err: ", err)
		return
	}
	var urlChan = make(chan string, 10)
	var data = make([]utils.BookInfo, 0)
	var m sync.Mutex
	var wg sync.WaitGroup
	var counter = 0
	// producer
	go utils.GetNextPage(client, req_url+fmt.Sprintf("?func=find-b&request=%s&find_code=WRD", request), urlChan)
	for url := range urlChan{
		wg.Add(1)
		// consumer
		counter ++
		go utils.Search(counter, client, url, &data, m, &wg)
	}
	wg.Wait()
	// 创建文件保存搜索内容
	file, err := os.Create(filename)
	if err != nil{
		fmt.Println("create file err: ", err)
		return 
	}
	encoder := json.NewEncoder(file)
	// 设置格式
	encoder.SetIndent(" ", "  ")
	// 不对html进行转义
	encoder.SetEscapeHTML(false)
	err = encoder.Encode(data)
	if err != nil{
		fmt.Println("save file error: ", err)
		return
	}
	fmt.Printf("搜索完成, 共用时: %.2f s", float64(time.Since(start))/float64(time.Second))
}
