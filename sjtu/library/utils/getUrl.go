package utils

import (
	"fmt"
	"html"
	"net/http"
	"regexp"
)

var NextPage = regexp.MustCompile(`(?s);<a href=(.*?) title=Next Page>`)

// 获取到url链接
// 每一次访问的url有可能不太一样
func GetUrl(client http.Client) (string, error) {
	const url = "http://opac.lib.sjtu.edu.cn/F"
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var buff [512]byte
	num, err := resp.Body.Read(buff[:])
	if err != nil {
		return "", err
	}
	pattern, err := regexp.Compile(`URL=(.*?)\?func`)
	if err != nil {
		return "", err
	}
	return string(pattern.FindSubmatch(buff[:num])[1]), nil
}

// 不断获取下一页的链接
func GetNextPage(client http.Client, url string, urlChan chan string) {
	urlChan <- url
	html_source, err := getHtml(client, url)
	if err != nil {
		fmt.Println("get html source err: ", err)
		return
	}
	next := NextPage.FindSubmatch(html_source)
	if next == nil {
		close(urlChan)
		return
	}

	GetNextPage(client, html.UnescapeString(string(next[1])), urlChan)
}



