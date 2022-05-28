package main

import (
  "fmt"
  "net/http"
  "net/url"
  "io/ioutil"
  "log"
  "io"
  "os"
)

func main() {
	// カレントディレクトリを取得する
	p, _ := os.Getwd()
    status, content := SendPostRequest("https://markdown.yuzunoha.net/api/go", p + "/test.txt")
    fmt.Println(status)
    fmt.Println(string(content))
}

func sub1() {
	url := "https://markdown.yuzunoha.net/go"

	resp, _ := http.Get(url)
	defer resp.Body.Close()
  
	byteArray, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byteArray)) // htmlをstringで取得
}

func sub2(argUrlStr string) {
    // url.Values{}でPOSTで送信する入れ物を準備
    ps := url.Values{}

    // Add()でPOSTで送信するデータを作成
    ps.Add("id", "1")
    ps.Add("name", "もりぴ")

    // 特殊文字や日本語をエンコード
    fmt.Println(ps.Encode())

    // http.PostForm()でPOSTメソッドを発行
    res, err := http.PostForm(argUrlStr, ps)

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

func SendPostRequest(url string, filename string) (string, []byte) {
    client := &http.Client{}
    data, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    req, err := http.NewRequest("POST", url, data)
    if err != nil {
        log.Fatal(err)
    }
    resp, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    content, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
    return resp.Status, content
}

// func main() {
//    status, content := SendPostRequest("https://api.example.com/upload", "test.jpg")
//    fmt.Println(status)
//    fmt.Println(string(content))
// }
