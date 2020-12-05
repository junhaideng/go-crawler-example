package utils

import (
	"net/http"
	"sync"
	"fmt"
)

func Search(counter int, client http.Client, url string, data *[]BookInfo, m sync.Mutex, wg *sync.WaitGroup) {
	fmt.Printf("正在搜索 %d 页...\n", counter)
	defer wg.Done()
	defer func(){
		fmt.Printf("第 %d 页搜索完毕 √\n", counter)
	}()
	html_source, err := getHtml(client, url)
	if err != nil {
		return
	}
	m.Lock()
	defer m.Unlock()
	*data = append(*data, Extract(html_source)...)
}
