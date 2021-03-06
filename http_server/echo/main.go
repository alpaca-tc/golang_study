// メイン関数(実行時に呼ばれる関数)を含むpackageはmainにする
package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

var (
	counter = 0
)

const (
	SQUARE_NUM_LIMIT = 100
)

func main() {
	e := echo.New()

	// GET
	e.GET("/hello", helloHandler)
	// GET 200以外のStatus
	e.GET("/401", unAuthorizedHandler)
	// GET Headerの読み込み
	e.GET("/square", squareHandler)
	// POST Bodyの読み込み
	e.POST("/incr", incrementHandler)
	e.GET("/alpaca", alpacaHandler)

	// 8080ポートで起動
	e.Logger.Fatal(e.Start(":8080"))
}

// レスポンスに`Hello World`を書き込むハンドラー
// 引数をこの形にするのはechoの仕様から決まっている
func helloHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World from Go.")
}

// 200以外のHTTP Statusを返すハンドラー
func unAuthorizedHandler(c echo.Context) error {
	return echo.NewHTTPError(http.StatusUnauthorized, "UnAuthorized")
}

// Headerから数字を取得して、その二乗を返すハンドラー
func squareHandler(c echo.Context) error {
	// Headerの読み込み
	numStr := c.Request().Header.Get("num")
	// String -> Intの変換
	num, err := strconv.Atoi(numStr)
	if err != nil {
		// 他のエラーの可能性もあるがサンプルとして纏める
		return echo.NewHTTPError(http.StatusBadRequest, "num is not integer")
	}

	if num >= SQUARE_NUM_LIMIT {
		// fmt.Sprintfでフォーマットに沿った文字列を生成できる。
		return c.String(http.StatusBadRequest, fmt.Sprintf("num must be less than %d", SQUARE_NUM_LIMIT))
	} else {
		// fmt.Sprintfでフォーマットに沿った文字列を生成できる。
		return c.String(http.StatusOK, fmt.Sprintf("Square of %d is equal to %d", num, num*num))
	}
}

// Bodyから数字を取得してその数字だけCounterをIncrementするハンドラー
// DBがまだないので簡易的なもの
func incrementHandler(c echo.Context) error {
	incrRequest := incrRequest{}
	if err := c.Bind(&incrRequest); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	counter += incrRequest.Num
	return c.JSON(http.StatusOK, incrResponse{counter})
}

type incrRequest struct {
	// jsonタグをつける事でjsonのunmarshalが出来る
	// jsonパッケージに渡すので、Publicである必要がある
	Num int `json:"num"`
}

type incrResponse struct {
	Counter int `json:"counter"`
}

func alpacaHandler(c echo.Context) error {
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
    　　　 UU
  `

	return c.String(http.StatusOK, alpacaAa)
}
