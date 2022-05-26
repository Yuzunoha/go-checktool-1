package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
)


func main() {
  url := "http://google.co.jp"

  resp, _ := http.Get(url)
  defer resp.Body.Close()

  byteArray, _ := ioutil.ReadAll(resp.Body)
  fmt.Println(string(byteArray)) // htmlをstringで取得
}
