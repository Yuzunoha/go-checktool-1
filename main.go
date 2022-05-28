package main

import (
  "fmt"
  "net/http"
  "net/url"
  "io/ioutil"
  "log"
  "io"
  "os"
  "bytes"
  "mime/multipart"
  "path/filepath"    
)

func main() {
	// 送信先のurl
	var toUrlStr = "https://markdown.yuzunoha.net/api/go"
	// カレントディレクトリを取得する
	p, _ := os.Getwd()
	// カレントディレクトリのtest.txtファイルのフルパス
	var filePath = p + "/test.txt"
	// リクエスト発行
	var content = SendPostRequest(toUrlStr, filePath, "multipart/form-data");
	// 結果表示
    fmt.Println(string(content))
}

func SendPostRequest (url string, filename string, filetype string) []byte {
    file, err := os.Open(filename)

    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()


    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    part, err := writer.CreateFormFile(filetype, filepath.Base(file.Name()))

    if err != nil {
        log.Fatal(err)
    }

    io.Copy(part, file)
    writer.Close()
    request, err := http.NewRequest("POST", url, body)

    if err != nil {
        log.Fatal(err)
    }

    request.Header.Add("Content-Type", writer.FormDataContentType())
    client := &http.Client{}

    response, err := client.Do(request)

    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()

    content, err := ioutil.ReadAll(response.Body)

    if err != nil {
        log.Fatal(err)
    }

    return content
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

func SendPostRequestOld(url string, filename string) (string, []byte) {
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
