package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"os"
	"path"
	"strings"
)

type Data struct {
	Link string `json:"link"`
	Description string `json:"description"`
	Language string `json:"language"`
	Stars string `json:"stars"`
	Fork string `json:"fork"`
	StarToday string `json:"star_today"`
}

const url = "https://github.com/trending"

func main(){
	collector := colly.NewCollector()
	var data []Data
	// 爬取网页中的所有article 元素
	collector.OnHTML("article", func(el *colly.HTMLElement){
		// 仓库链接
		link, _ := el.DOM.Find("h1").Find("a").Attr("href")
		// 仓库描述
		description := strings.TrimSpace(el.DOM.Find("p").Text())
		// 使用语言
		language := el.DOM.Find("span[itemprop=programmingLanguage]").Text()
		// star数目
		stars := strings.TrimSpace(el.DOM.Find(".muted-link").First().Text())
		// fork 数目
		fork := strings.TrimSpace(el.DOM.Find(".muted-link").Last().Text())
		// 今日star数目
		starToday := strings.TrimSpace(el.DOM.Find(".float-sm-right").Text())
		// 添加到data中
		data = append(data, Data{
			Link:        path.Join(url, link),
			Description: description,
			Language:    language,
			Stars:       stars,
			Fork:        fork,
			StarToday:   starToday,
		})
	})

	collector.Visit(url)
	// 创建文件
	file, err := os.Create("trending.json")
	if err != nil{
		fmt.Println("err: ", err)
		return
	}
	encoder := json.NewEncoder(file)
	// 不转义html文本
	encoder.SetEscapeHTML(false)
	// 设置缩进
	encoder.SetIndent(" ", "  ")
	encoder.Encode(data)
}