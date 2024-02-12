package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type piApiResponse struct {
	Pi  int64  `json:"pi"`
	Err string `json:"error"`
}

func loadBalancingFunction(inputNumber float64) int {
	var serverNumber int
	if inputNumber < 100 {
		serverNumber = 1
	} else {
		serverNumber = 2
	}
	return serverNumber
}

func main() {

	htmlTemplate, err := template.ParseFiles("index.html")
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		err := htmlTemplate.Execute(c.Writer, gin.H{})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	})

	r.POST("/", func(c *gin.Context) {

		inputNumber, err := strconv.ParseFloat(c.Request.PostFormValue("inputNumber"), 64)
		if err != nil {
			fmt.Println("Error in parsing input: ", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid number",
			})
			return
		}

		serverNumber := loadBalancingFunction(inputNumber)
		url := fmt.Sprintf("http://pi_api_%d:8080/pi?target=%f", serverNumber, inputNumber)

		res, err := http.Get(url)

		if err != nil {
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer res.Body.Close()

		var piResult piApiResponse

		if res.StatusCode != http.StatusOK {
			if err := json.NewDecoder(res.Body).Decode(&piResult); err != nil {
				fmt.Println(err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
			}
			fmt.Println("Failure:", piResult.Err)
			fmt.Printf("%d : Something went wrong", res.StatusCode)
			return
		}

		if err := json.NewDecoder(res.Body).Decode(&piResult); err != nil {
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = htmlTemplate.Execute(c.Writer, gin.H{
			"Result": piResult.Pi,
		})
	})

	r.Run(":3000")
}
