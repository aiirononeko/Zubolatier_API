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
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func main() {

	// Cloud FireStoreの初期化
	ctx := context.Background()
	sa := option.WithCredentialsFile("path/to/serviceAccount.json")
	app, err := firebase.NewApp(ctx, nil, sa)
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

	/**
	 * 一旦全件取得してから生成した乱数と照合するようになっているので、
	 * 全件取得しなくても良いようにリファクタリングする必要あり.
	 */
	// ドキュメントをランダムに1件取得
	e.GET("/get", func(c echo.Context) error {

		// 乱数生成
		rand.Seed(time.Now().UnixNano())
		randNum := rand.Intn(3) + 1 // 1~3
		randStr := strconv.Itoa(randNum)

		// ドキュメント全件取得
		iter := client.Collection("recipes").Documents(ctx)
		var res interface{}
		for {

			// 1件ずつ処理していく
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to iterate: %v", err)
			}

			// 生成された乱数とIDを照合して1件のドキュメントを返却
			if doc.Ref.ID == randStr {
				res = doc.Data()
			}
		}

		// ドキュメントのnilチェック
		if res == nil {
			return c.JSON(http.StatusNotFound, "ドキュメントが見つかりませんでした")
		}

		return c.JSON(http.StatusOK, res)
	})

	// サーバー起動
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
