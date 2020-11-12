package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

type Antrian struct {
	Id     string `json:"id"`
	Status bool   `json:"status"`
}

var data []Antrian

var client *db.Client
var ctx context.Context

func init() {
	ctx = context.Background()
	conf := &firebase.Config{
		DatabaseURL: "https://{URL_FIREBASE}.firebaseio.com/",
	}

	opt := option.WithCredentialsFile("firebase-admin-sdk.json")

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln("Error initializing app:", err)
	}

	client, err = app.Database(ctx)
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
	}
}

func main() {
	router.LoadHTMLGlob("views/*")

	router := gin.Default()
	router.GET("/", getSomething)
	router.POST("/api/v1/antrian", AddAntrianHandler)
	router.GET("/api/v1/antrian/status", GetAntrianHandler)
	router.PUT("/api/v1/antrian/id/:idAntrian", UpdateAntrianHandler)
	router.DELETE("/api/v1/antrian/id/:idAntrian/delete", DeleteAntrianHandler)
	router.Run(":8080")
}

func getSomething(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"body1": "GET.Something Success",
	})
	return
}

func AddAntrianHandler(c *gin.Context) {
	flag, err := addAntrian()
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status": "failed",
		})
		return
	}

	if flag {
		c.JSON(http.StatusOK, map[string]interface{}{
			"status": "success",
		})
	} else {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status": "failed",
		})
		return
	}
}

func addAntrian() (bool, error) {
	_, _, dataAntrian := getAntrian()
	var Id string
	var antrianRef *db.Ref
	ref := client.NewRef("antrian")

	if dataAntrian == nil {
		Id = fmt.Sprintf("B-0")
		antrianRef = ref.Child("0")
	} else {
		Id = fmt.Sprintf("B-%d", len(dataAntrian))
		antrianRef = ref.Child(fmt.Sprintf("%d", len(dataAntrian)))
	}
	antrian := Antrian{
		Id:     Id,
		Status: false,
	}
	if err := antrianRef.Set(ctx, antrian); err != nil {
		log.Fatal(err)
		return false, err
	}
	return true, nil
}

func GetAntrianHandler(c *gin.Context) {
	flag, err, resp := getAntrian()
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status": "failed",
		})
		return
	}
	if flag {
		c.JSON(http.StatusOK, map[string]interface{}{
			"status": "success",
			"data":   resp,
		})
	} else {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "failed",
			"message": "unknown",
		})
		return
	}
}

func getAntrian() (bool, error, []map[string]interface{}) {
	var data []map[string]interface{}
	ref := client.NewRef("antrian")
	if err := ref.Get(ctx, &data); err != nil {
		log.Fatalln("Error reading from database:", err)
		return false, err, nil
	}
	return true, nil, data
}

func UpdateAntrianHandler(c *gin.Context) {
	idAntrian := c.Param("idAntrian")
	err := updateAntrian(idAntrian)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
	})

}

func updateAntrian(idAntrian string) error {
	ref := client.NewRef("antrian")
	id := strings.Split(idAntrian, "-")
	childRef := ref.Child(id[1])
	antrian := Antrian{
		Id:     idAntrian,
		Status: true,
	}
	if err := childRef.Set(ctx, antrian); err != nil {
		log.Fatal(err)
		return false, err
	}
	return true, nil
}

func DeleteAntrianHandler(c *gin.Context) {
	idAntrian := c.Param("idAntrian")
	err := deleteAntrian(idAntrian)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
	})
}

func deleteAntrian(idAntrian string) (bool, error) {

	ref := client.NewRef("antrian")
	id := strings.Split(idAntrian, "-")
	childRef := ref.Child(id[1])
	if err := childRef.Delete(ctx); err != nil {
		log.Fatal(err)
		return false, err
	}
	return true, nil
}
