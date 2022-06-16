/*
 * @Author: feng.pan
 * @Date: 2022-06-16 19:03:12
 * @LastEditors: feng.pan
 * @LastEditTime: 2022-06-16 19:09:46
 * @FilePath: /work/biququDownload/main.go
 * @Description: todo...
 */
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func delErr[T any](res T, err error) T {
	if err != nil {
		panic(err)
	}
	return res
}

func judgeCode(code int, msg string) {
	if code != 200 {
		panic(code)
		fmt.Println("错误：", code, msg)
		return
	}
}

func main() {
	var (
		url string
	)
	fmt.Print("请输入笔趣阁小说地址 \n")
	fmt.Scanf("%s", &url)
	BookLists := make([]string, 0)
	getAllList(url, BookLists)
}

func getAllList(url string, BookLists []string) {
	fmt.Println("即将下载改地址小说：", url)
	res := delErr(http.Get(url))
	defer res.Body.Close()
	judgeCode(res.StatusCode, "获取章节列表错误！")
	doc := delErr(goquery.NewDocumentFromReader(res.Body))
	doc.Find("#list a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		BookLists = append(BookLists, "http://www.biququ.com"+href)
	})

	doc.Find("#info h1").Each(func(i int, s *goquery.Selection) {
		name := "./" + s.Text() + ".txt"
		Download(BookLists, 0, delErr(os.Create(name)))
	})
}

func Download(bookLists []string, idx int, file *os.File) {
	if idx >= len(bookLists) {
		fmt.Print("下载完成 \n")
		return
	}

	res := delErr(http.Get(bookLists[idx]))
	defer res.Body.Close()
	judgeCode(res.StatusCode, "请查看是否下载完成！")
	doc := delErr(goquery.NewDocumentFromReader(res.Body))
	doc.Find("h1").Each(func(i int, s *goquery.Selection) {
		fmt.Println("正在下载", s.Text(), "...")
		file.WriteString(s.Text() + "\n")

	})
	doc.Find("#content p").Each(func(i int, s *goquery.Selection) {
		if s.Find("a").Length() > 0 {
			return
		}
		file.WriteString(s.Text() + "\n")
	})
	idx++
	Download(bookLists, idx, file)
}
