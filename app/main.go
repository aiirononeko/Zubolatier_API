package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	firebase "firebase.google.com/go"
	"github.com/labstack/echo"
)

func main() {

	// Cloud FireStoreの初期化
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: "Zubolatier"}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	// サーバーのインスタンスを作成
	e := echo.New()

	// ドキュメントをランダムに1件取得
	e.GET("/get", func(c echo.Context) error {

		// 乱数生成
		rand.Seed(time.Now().UnixNano())
		randNum := rand.Intn(7) + 1 // 1~7
		randStr := strconv.Itoa(randNum)

		// 生成した乱数をIDに持つドキュメントを1件取得
		recipe, err := client.Collection("recipes").Doc(randStr).Get(ctx)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		// ドキュメントのnilチェック
		if recipe == nil {
			return c.JSON(http.StatusNotFound, "ドキュメントが見つかりませんでした")
		}

		res := recipe.Data()
		return c.JSON(http.StatusOK, res)
	})

	// サーバー起動
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
