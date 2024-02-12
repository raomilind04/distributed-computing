package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type piApiResponse struct {
	Pi  int64  `json:"pi"`
	Err string `json:"error"`
}

func loadBalancingFunction(x float64) int {
    fmt.Println(os.Getenv("SERVER_COUNT"))
    serverCount, err := strconv.ParseInt(os.Getenv("SERVER_COUNT"), 10, 64)
    if err != nil {
        fmt.Println("SERVER_COUNT env variable can not be parsed")
        panic("SERVER_COUNT is invalid")
    }
    var serverNumber int
    if serverCount == 1 {
        serverNumber = 1
    } else if serverCount == 2 {
        if x >= 1 && x < 45215900000000 {
            serverNumber = 1
        } else {
            serverNumber = 2
        }
    } else if serverCount == 4 {
        if x >= 1 && x <= 383027000000 {
            serverNumber = 1
        } else if x >= 383028000000 && x <= 45215900000000 {
            serverNumber = 2
        } else if x >= 45216000000000 && x <= 9412060000000000 {
            serverNumber = 3
        } else {
            serverNumber = 4
        }
    } else if serverCount == 8 {
        if x >= 1 && x <= 34846100000 {
            serverNumber = 1
        } else if x >= 34846200000 && x <= 383027000000 {
            serverNumber = 2
        } else if x >= 383028000000 && x <= 4175930000000 {
            serverNumber = 3
        } else if x >= 4175940000000 && x <= 45215900000000 {
            serverNumber = 4
        } else if x >= 45216000000000 && x <= 906640000000000 {
            serverNumber = 5
        } else if x >= 906641000000000 && x <= 9412060000000000 {
            serverNumber = 6
        } else if x >= 9412070000000000 && x <= 97577200000000000 {
            serverNumber = 7
        } else {
            serverNumber = 8
        }
    } else {
        panic("Invalid server Count: Should be 1, 2, 4, 8")
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
            return
      }
      err = htmlTemplate.Execute(c.Writer, gin.H{
            "Error": piResult.Err,
        })
			fmt.Println("Failure:", piResult.Err)
			fmt.Printf("%d : Something went wrong", res.StatusCode)
		} else {
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
    }
})

	r.Run(":3000")
}
