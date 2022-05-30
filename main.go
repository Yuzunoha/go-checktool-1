package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

var p = fmt.Println

func main() {
	// コマンドライン引数をパースして取得する
	flag.Parse()
	args := flag.Args()

	// ガード: コマンドライン引数1の存在チェック("check", "submit", ...)
	if len(args) <= 0 {
		p("コマンドライン引数を指定してください")
		p("ex:")
		p("  codegym test")
		p("  codegym submit fizzbuzz")
		return
	}

	// コマンドライン引数1で分岐する
	switch args[0] {
	case "test":
		test()
	case "submit":
		// ガード: コマンドライン引数2の存在チェック("fizzbuzz", "fukuri", ...)
		if len(args) <= 1 {
			p("2つ目のコマンドライン引数を指定してください")
			p("ex:")
			p("  codegym submit fizzbuzz")
			p("  codegym submit fukuri")
			return
		}
		submit(args[1])
	default:
		p("存在しないオプションです: " + args[0])
	}
}

func test() {
	// カレントディレクトリを取得する(goの実行ファイルの場所ではない。実行した場所である)
	currentDir, _ := os.Getwd()
	filePath := currentDir + "/" + "main.go"
	if isExist(filePath) {
		fmt.Println("存在します")
	} else {
		fmt.Println("存在しません")
	}
}

func isExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func submit(taskKey string) {
	// 有効なタスクキー一覧
	taskKeys := []string{"fizzbuzz", "fukuri"}
	// ガード: 有効なタスクキーかどうか
	if false == contains(taskKeys, taskKey) {
		p("存在しないオプションです: " + taskKey)
		return
	}
	// カレントディレクトリを取得する(goの実行ファイルの場所ではない。実行した場所である)
	currentDir, _ := os.Getwd()
	// 送信対象ファイルのフルパス
	filePath := currentDir + "/" + taskKey + ".php"

	// ガード: タスクキーに該当するphpファイルが同じディレクトリに存在すること
	if (true) {
		p("カレントディレクトリに次のファイルが存在しません: " + filePath)
		return
	}
	// p(filename)
	return
	// 送信先のurl
	// var toUrlStr = "https://markdown.yuzunoha.net/api/go"
	// var toUrlStr = "http://localhost/api/go"
	// カレントディレクトリを取得する
	//path, _ := os.Getwd()
	// カレントディレクトリのtest.txtファイルのフルパス
	// var filePath = path + "/test.txt"
	// リクエスト発行
	//var content = SendPostRequest(toUrlStr, filePath, "file")
	// 結果表示
	//p(string(content))
}

// 配列aに要素eが含まれていればtrueを返す関数
func contains(a []string, e string) bool {
	for _, v := range a {
		if e == v {
			return true
		}
	}
	return false
}

func SendPostRequest(url string, filePath string, fieldName string) []byte {
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
	// postのキーバリューを乗せられる
	writer.WriteField("a", "b")
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
