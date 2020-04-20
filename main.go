package main

import (
	"context"
	"log"
	"net/http"
	"os"

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

	// ルーティング設定
	e.GET("/all", func(c echo.Context) error {

		// 全レコード取得
		iter := client.Collection("recipes").Documents(ctx)
		var res interface{}
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to iterate: %v", err)
			}
			res = doc.Data()
		}

		return c.JSON(http.StatusOK, res)
	})

	// サーバー起動
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))

}
