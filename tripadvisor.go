package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"
)

const PageNum = 1170 // 手动设置，或者首先执行一遍代码获取

type Comment struct {
	Author string `json:"author"`
	Text   string `json:"text"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var m sync.Mutex
	var wg sync.WaitGroup
	var urlTemplate = "https://www.tripadvisor.co.uk/Attraction_Review-g294212-d319086-Reviews-or%d-Forbidden_City_The_Palace_Museum-Beijing.html"
	var comments []Comment
	var dataPattern = regexp.MustCompile(`"mgmtResponse":.*?,"text":"(.*?)","username":"(.*?)",`)
	file, err := os.Create("tripadvisor.json")
	if err != nil {
		fmt.Println("Create file error: ", err)
		return
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("  ", "    ")

	// 该网站要控制一下同时请求的连接数
	// 如果连接数太多了，请求会失败
	var limit = make(chan struct{}, 10)

	var client http.Client
	start := time.Now()
	for i := 0; i < PageNum; i++ {
		wg.Add(1)
		limit <- struct{}{}
		go func(i int) {
			defer func() {
				fmt.Printf("已完成第%d页\n", i)
				<-limit
				wg.Done()
			}()
			fmt.Printf("正在处理第%d页\n", i)
			url := fmt.Sprintf(urlTemplate, i*5)
			resp, err := client.Get(url)
			if err != nil {
				fmt.Println("Get error: ", err)
				return
			}
			resource, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			for _, text := range dataPattern.FindAllSubmatch(resource, -1) {
				m.Lock()
				comments = append(comments, Comment{
					Author: string(text[2]),
					Text:   string(text[1]),
				})
				m.Unlock()
			}
		}(i)

	}
	wg.Wait()

	encoder.Encode(comments)
	fmt.Println("共花时: ", time.Since(start))

}
