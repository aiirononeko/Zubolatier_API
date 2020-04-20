package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo"
)

func main() {
	// サーバーのインスタンスを作成
	e := echo.New()

	// ルーティング設定
	e.GET("/hello", hello)

	// サーバー起動
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "hello!")
}
