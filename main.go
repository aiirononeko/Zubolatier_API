package main

import (
	"context"
	"encoding/json"
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

type Settings struct {
	Type                    string `json="type"`
	ProjectId               string `json="project_id"`
	PrivateKeyId            string `json="private_key_id"`
	PrivateKey              string `json="private_key"`
	ClientEmail             string `json="client_email"`
	ClientId                string `json="client_id"`
	AuthUri                 string `json="auth_uri"`
	TokenUri                string `json="token_uri"`
	AuthProviderX509CertUrl string `json="auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json="client_x509_cert_url"`
}

func main() {

	s := Settings{
		Type:                    os.Getenv("TYPE"),
		ProjectId:               "zubolatier",
		PrivateKeyId:            os.Getenv("PRIVATE_KEY_ID"),
		PrivateKey:              os.Getenv("PRIVATE_KEY"),
		ClientEmail:             os.Getenv("CLIENT_EMAIL"),
		ClientId:                os.Getenv("CLIENT_ID"),
		AuthUri:                 os.Getenv("AUTH_URI"),
		TokenUri:                os.Getenv("TOKEN_URI"),
		AuthProviderX509CertUrl: os.Getenv("AUTH_PROVIDER_X590_CERT_URL"),
		ClientX509CertUrl:       os.Getenv("CLIENT_X509_CERT_URL"),
	}

	// アカウント情報JSON生成
	jsonBytes, err := json.Marshal(s)
	if err != nil {
		log.Fatalln(err)
	}

	// Cloud FireStoreの初期化
	ctx := context.Background()
	// sa := option.WithCredentialsFile("path/to/serviceAccount.json")
	sa := option.WithCredentialsJSON(jsonBytes)

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
