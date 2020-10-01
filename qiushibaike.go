package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const PageNum = 13

type Jokes struct {
	Author  string `json:"author"`
	Link    string `json:"link"`
	Content string `json:"content"`
	Type    string `json:"type"`
}

func main() {
	// 每一页对应的url链接
	var url = "https://www.qiushibaike.com/text/page/%d/"
	var jokes []Jokes
	var divPattern = regexp.MustCompile(`(?s)<div class="article.*?typs_(.*?)".*?>.*?<div class="author.*?>.*?<a onclick.*?>.*?<h2>(.*?)</h2>.*?<a href="(.*?)".*? class="contentHerf".*?>.*?<div class="content">.*?<span>(.*?)</span>.*?</div>`)

	for i := 1; i < PageNum; i++ {
		resp, err := http.Get(fmt.Sprintf(url, i))
		if err != nil {
			fmt.Println("Get error: ", err)
			continue
		}

		data, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Println("ReadAll Err: ", err)
			continue
		}

		for _, div := range divPattern.FindAllSubmatch(data, -1) {
			jokes = append(jokes, Jokes{
				Author:  strings.TrimSpace(string(div[2])),
				Link:    "https://www.qiushibaike.com" + string(div[3]),
				Content: strings.TrimSpace(string(div[4])),
				Type:    string(div[1]),
			})
		}
	}
	file, err := os.Create("jokes.json")
	if err != nil {
		fmt.Println("Create file Error")
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("  ", "    ")
	if err := encoder.Encode(jokes); err != nil {
		fmt.Println("Encode Err: ", err)
	}
}
