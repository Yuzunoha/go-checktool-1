package main

import (
  "fmt"
  "net/http"
  "net/url"
  "io/ioutil"
  "log"
  "io"
)

func main() {
	sub2();
}

func sub1() {
	url := "https://markdown.yuzunoha.net/go"

	resp, _ := http.Get(url)
	defer resp.Body.Close()
  
	byteArray, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byteArray)) // htmlをstringで取得
}

func sub2() {
    // POSTメソッド
    // url.Values{}でPOSTで送信する入れ物を準備
    ps := url.Values{}

    // Add()でPOSTで送信するデータを作成
    ps.Add("id", "1")
    ps.Add("name", "もりぴ")

    // 特殊文字や日本語をエンコード
    fmt.Println(ps.Encode())

    // http.PostForm()でPOSTメソッドを発行
    res, err := http.PostForm("https://markdown.yuzunoha.net/api/go", ps)

    if err != nil {
        log.Fatal(err)
    }

    // deferでクローズ処理
    defer res.Body.Close()
    // Bodyの内容を読み込む
    body, _ := io.ReadAll(res.Body)
    // Bodyの内容を出力する
    fmt.Print(string(body))
}
