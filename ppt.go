package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	DOWNLOAD = "./download"
	NUM      = 50
)

func main() {
	start := time.Now()
	var urlChan = make(chan string, 10)
	var wg sync.WaitGroup
	go func() {
		for i := 1; i < NUM; i++ {
			wg.Add(1)
			urlChan <- "http://www.ypppt.com/p/d.php?aid=" + strconv.Itoa(i)
		}
		close(urlChan)
	}()

	for url := range urlChan {
		// 休眠一段时间
		// 不进行休眠的话一段时间里goroutine可能过大(上万)
		// 从而导致程序崩溃
		time.Sleep(100 * time.Millisecond)
		go GetDownloadLink(url, &wg)
	}

	wg.Wait()
	fmt.Println(time.Since(start))
}

func GetDownloadLink(url string, wg *sync.WaitGroup) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Err: ", err)
		wg.Done()
		return
	}
	statusCode := resp.StatusCode
	for statusCode != 200 {
		if statusCode == http.StatusNotFound {
			fmt.Printf("URL: %s 无效\n", url)
			wg.Done()
			return
		}
		resp, err := http.Get(url)
		resp.Body.Close()
		fmt.Println("重试: ", resp.StatusCode)
		if err != nil {
			fmt.Println("Err: ", err)
			wg.Done()
			return
		}
		statusCode = resp.StatusCode
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Err: ", err)
		wg.Done()
		return
	}
	resp.Body.Close()

	// 我们只需要取第一个url即可
	downloadLinkPattern := regexp.MustCompile(`(?s)<.*?class="down clear">.*?<li>.*?<a href="(.*?)".*?>`)

	namePattern := regexp.MustCompile(`(?s)<div class="de">.*?<h1>(.*?)-.*?下载页</h1>`)
	//fmt.Println(string(data))
	link := string(downloadLinkPattern.FindSubmatch(data)[1])
	name := string(namePattern.FindSubmatch(data)[1])
	if !strings.HasPrefix(link, "http") && len(link) > 8 {
		link = "http://www.youpinppt.com" + link[8:]
	}
	if len(link) < 8 {
		wg.Done()
		return
	}
	go downloadFile(link, strings.TrimSpace(name), wg)

}

func downloadFile(link, name string, wg *sync.WaitGroup) {
	resp, err := http.Get(link)
	if err != nil {
		fmt.Println("Err", err)
		wg.Done()
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Err: ", err)
		wg.Done()
		return
	}
	defer resp.Body.Close()
	fmt.Printf("开始下载：%s, link: %s \n", name, link)
	ext := filepath.Ext(link)
	if len(ext) == 0 {
		fmt.Println("暂不支持该文件类型下载，请自行下载: ", link)
		wg.Done()
		return
	}
	file, err := os.Create(filepath.Join(DOWNLOAD, name) + ext)
	if err != nil {
		fmt.Println("Err: ", err)
		wg.Done()
		return
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		fmt.Println("Err", err)
		wg.Done()
		return
	} else {
		fmt.Printf("下载 %s 完成\n", name)
		wg.Done()
	}
}
