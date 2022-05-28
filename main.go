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

var p = fmt.Println

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
	// フォームを作る。右辺の引数はそれぞれ "file", "test.txt"
    part, err := writer.CreateFormFile(fieldName, filepath.Base(file.Name()))
    if err != nil {
        log.Fatal(err)
    }
	// フォームに送信するファイルの中身をコピーする
    io.Copy(part, file)
	// MultipartWriterを閉じる
    writer.Close()
	// リクエストを作成する。ペイロードはバッファ領域へのポインタ「body」を指定することに注意
    request, err := http.NewRequest("POST", url, body)
    if err != nil {
        log.Fatal(err)
    }
	// リクエストヘッダを書き込む。値は "multipart/form-data; boundary=aab16..."
    request.Header.Add("Content-Type", writer.FormDataContentType())
	// クライアントのポインタを取得する
    client := &http.Client{}
	// リクエストを実行してレスポンスを受け取る
    response, err := client.Do(request)
    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()
	// レスポンスからコンテンツを抽出して返却する。おわり
    content, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }
    return content
}
