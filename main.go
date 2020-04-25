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
	"google.golang.org/api/option"
)

func main() {

	// アカウント情報JSON生成
	// settingsMap := map[string]interface{}{
	// 	"type":                        os.Getenv("TYPE"),
	// 	"project_id":                  os.Getenv("PROJECT_ID"),
	// 	"private_key_id":              os.Getenv("PRIVATE_KEY_ID"),
	// 	"private_key":                 os.Getenv("PRIVATE_KEY"),
	// 	"client_email":                os.Getenv("CLIENT_EMAIL"),
	// 	"client_id":                   os.Getenv("CLIENT_ID"),
	// 	"auth_uri":                    os.Getenv("AUTH_URI"),
	// 	"token_uri":                   os.Getenv("TOKEN_URI"),
	// 	"auth_provider_x509_cert_url": os.Getenv("AUTH_PROVIDER_X509_CERT_URL"),
	// 	"client_x509_cert_url":        os.Getenv("CLIENT_X509_CERT_URL"),
	// }

	// settings, err := json.Marshal(settingsMap)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// Cloud FireStoreの初期化
	ctx := context.Background()
	sa := option.WithCredentialsFile("path/to/serviceAccount.json")
	// sa := option.WithCredentialsJSON(settings)

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
