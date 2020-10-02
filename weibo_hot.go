// 微博热搜榜:  "https://s.weibo.com/top/summary"

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

const link = "https://s.weibo.com/top/summary"

type Item struct {
	Id string `json:"id"`
	Link string `json:"link"`
	Keyword string `json:"keyword"`
	Number string `json:"number"`
	Label string `json:"label"`
}

func main() {
	var tablePattern = regexp.MustCompile(`(?s)<div.*?id="pl_top_realtimehot">.*?<tbody>(.*?)</tbody>.*?</div>`)
	html, err := getHtml(link)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 表格，保存了热搜榜中的内容
	table := tablePattern.FindSubmatch(html)[1]
		// <td.*?>.*?<a href="(.*?)".*?>(.*?)</a><span>(.*?)</span>
	var trPattern = regexp.MustCompile(`(?s)<tr.*?>.*?<td class=".*?ranktop">(?P<id>.*?)</td>.*?<td.*?>.*?<a href="(?P<link>.*?)".*?>(?P<keyword>.*?)</a>.*?<span>(?P<number>.*?)</span>.*?</td>.*?<td.*?>(.*?)</td>.*?</tr>`)
	// td-03 中可能存在一些i标签，我们只需要其中的文字
	var iPattern = regexp.MustCompile(`<i .*?>(.*?)</i>`)
	var items []Item
	for _, tr := range trPattern.FindAllSubmatch(table, -1){
		label := iPattern.FindSubmatch(tr[5])
		var temp []byte
		if len(label) == 0{
			temp = []byte("")
		}else{
			temp = label[1]
		}
		l, err := url.PathUnescape("https://s.weibo.com"+ string(tr[2]))
		if err != nil {
			l = "https://s.weibo.com"+ string(tr[2])
		}
		items = append(items, Item{
			Id:      string(tr[1]),
			Link:    l,
			Keyword: string(tr[3]),

			Number:  string(tr[4]),
			Label:   string(temp),
		})
	}

	file, err := os.Create("weibo_hot.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("  ", "    ")
	if err := encoder.Encode(&items); err!= nil{
		fmt.Println(err)
	}
}

func getHtml(url string)([]byte, error){
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return data, nil
}
