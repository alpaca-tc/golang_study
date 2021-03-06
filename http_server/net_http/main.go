// メイン関数(実行時に呼ばれる関数)を含むpackageはmainにする
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

var (
	counter = 0
)

const (
	SQUARE_NUM_LIMIT = 100
)

func main() {
	// GET
	http.HandleFunc("/hello", helloHandler)
	// GET 200以外のStatus
	http.HandleFunc("/401", unAuthorizedHandler)
	// GET Headerの読み込み
	http.HandleFunc("/square", squareHandler)
	// POST Bodyの読み込み
	http.HandleFunc("/incr", incrementHandler)

	http.HandleFunc("/alpaca", alpacaHandler)

	// 8080ポートで起動
	http.ListenAndServe(":8080", nil)
}

// レスポンスに`Hello World`を書き込むハンドラー
// 引数をこの形にするのはnet/httpの仕様から決まっている
func helloHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Hello World from Go.")
}

// 200以外のHTTP Statusを返すハンドラー
func unAuthorizedHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprint(w, "UnAuthorized")
}

// Headerから数字を取得して、その二乗を返すハンドラー
func squareHandler(w http.ResponseWriter, req *http.Request) {
	// Headerの読み込み
	numStr := req.Header.Get("num")
	// String -> Intの変換
	num, err := strconv.Atoi(numStr)
	if err != nil {
		// 他のエラーの可能性もあるがサンプルとして纏める
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "num is not integer")
		return
	}

	if num >= SQUARE_NUM_LIMIT {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "num must be less than 100")
	} else {
		// fmt.Sprintfでフォーマットに沿った文字列を生成できる。
		fmt.Fprint(w, fmt.Sprintf("Square of %d is equal to %d", num, num*num))
	}
}

// Bodyから数字を取得してその数字だけCounterをIncrementするハンドラー
// DBがまだないので簡易的なもの
func incrementHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body := req.Body
	// bodyの読み込みに開いたio Readerを最後にCloseする
	defer body.Close()

	buf := new(bytes.Buffer)
	io.Copy(buf, body)

	var incrRequest incrRequest
	// BodyのJSONを構造体に変換する
	json.Unmarshal(buf.Bytes(), &incrRequest)

	counter += incrRequest.Num
	fmt.Fprint(w, fmt.Sprintf("Value of Counter is %d \n", counter))
}

func alpacaHandler(w http.ResponseWriter, req *http.Request) {
	alpacaAa := `
    うるせぇアルパカぶつけるぞ 
    Δ~~~~Δ 
    ξ ･ェ･ ξ 
    ξ　~　ξ 
    ξ　　 ξ 
    ξ　　　“~～~～〇 
    ξ　　　　　　 ξ 
    ξ　ξ　ξ~～~ξ　ξ 
    　ξ_ξξ_ξ　ξ_ξξ_ξ 
    　　ヽ(´･ω･)ﾉ 
    　　　 |　 / 
    　　　 UU"
  `

	fmt.Fprint(w, alpacaAa)
}

type incrRequest struct {
	// jsonタグをつける事でjsonのunmarshalが出来る
	// jsonパッケージに渡すので、Publicである必要がある
	Num int `json:"num"`
}
