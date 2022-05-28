package main

import (
  "fmt"
  "net/http"
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
	var content = SendPostRequest(toUrlStr, filePath, "file");
	// 結果表示
    fmt.Println(string(content))
}

func SendPostRequest (url string, filePath string, fieldName string) []byte {
	// フルパス指定で送信するファイルを開く
    file, err := os.Open(filePath)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
	// 送信のためのバッファ領域を確保する
	body := &bytes.Buffer{}
	// MultipartWriterを作る。boundaryは自動で設定される。
    writer := multipart.NewWriter(body)
	// 送信するフォームを作る
    part, err := writer.CreateFormFile(fieldName, filepath.Base(file.Name()))
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
