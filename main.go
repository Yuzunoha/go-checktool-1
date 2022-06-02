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
	"time"

	"github.com/briandowns/spinner"
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

// 環境変数 CODEGYM_TOKEN が無ければ終了させる関数
func getEnvCodegymToken() string {
	s := os.Getenv("CODEGYM_TOKEN")
	if "" == s {
		p("環境変数 CODEGYM_TOKEN を設定してください。次のようなコマンドで設定します")
		p("  export CODEGYM_TOKEN=\"8|M9MkLyUztaW0EgaWPwaymOOS1UuJO4wXTlzGPMOZ\"")
		p("  右辺の値はポータルサイトから取得した値に置き換えてください")
		p("  先頭の数字部分 8| も必要なことに注意してください")
		os.Exit(0)
	}
	return s
}

func submit(taskKey string) {
	// ガード: 環境変数を取得する。無ければ強制終了する
	token := getEnvCodegymToken()
	// 有効なタスクキー一覧
	taskKeys := []string{"fizzbuzz", "kuku", "fukuri", "shiharai"}
	// ガード: 有効なタスクキーかどうか
	if false == contains(taskKeys, taskKey) {
		p("存在しないオプションです: " + taskKey)
		return
	}
	// カレントディレクトリを取得する(goの実行ファイルの場所ではない。実行した場所である)
	currentDir, _ := os.Getwd()
	// 送信対象ファイル名(ファイル名はタスクキーと同一)
	fileName := taskKey + ".php"
	// 送信対象ファイルのフルパス
	filePath := currentDir + "/" + fileName
	// ガード: タスクキーに該当するphpファイルが同じディレクトリに存在すること
	if false == isExist(filePath) {
		p("カレントディレクトリに次のファイルが存在しません: " + fileName)
		return
	}
	// 送信先のurl
	// クラウド
	toUrlStr := "https://markdown.yuzunoha.net/api/go/unittest"
	// ローカル
	// toUrlStr := "http://localhost/api/go/unittest"
	// リクエスト発行
	content := SendPostRequest(toUrlStr, token, filePath, taskKey)
	// 結果表示
	p(string(content))
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

func createSpinner() *spinner.Spinner {
	s := spinner.New(spinner.CharSets[1], 100*time.Millisecond) // Build our new spinner
	s.Prefix = "確認中です。ピー、ガー、ヒョロロロー... "                         // Prefix text before the spinner
	return s
}

func SendPostRequest(url string, token string, filePath string, taskKey string) []byte {
	// フルパス指定で送信するファイルを開く
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// 送信のためのバッファ領域を確保する
	body := &bytes.Buffer{}
	// MultipartWriterを作る。boundaryは自動で設定される
	writer := multipart.NewWriter(body)
	// フォームを作る。右辺の引数はそれぞれ "file", "test.txt"
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		log.Fatal(err)
	}
	// フォームに送信するファイルの中身をコピーする
	io.Copy(part, file)
	// postのキーバリュー。"fizzbuzz", "warikan", ...
	writer.WriteField("taskKey", taskKey)
	// MultipartWriterを閉じる
	writer.Close()
	// リクエストを作成する。ペイロードはバッファ領域へのポインタ「body」を指定することに注意
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Fatal(err)
	}
	// "multipart/form-data; boundary=aab16..."
	request.Header.Add("Content-Type", writer.FormDataContentType())
	// "Bearer 8|M9MkLyUztaW0EgaWPwaymOOS1UuJO4wXTlzGPMOZ"
	request.Header.Add("Authorization", "Bearer "+token)
	// クライアントのポインタを取得する
	client := &http.Client{}

	// スピナー開始
	spinner := createSpinner()
	spinner.Start()

	// リクエストを実行してレスポンスを受け取る
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// スピナー終了
	spinner.Stop()

	// レスポンスからコンテンツを抽出して返却する。おわり
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	return content
}
