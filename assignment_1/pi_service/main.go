package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func parseFile(target float64) (int64, error) {
    file, err := os.Open("master_file.txt")
    if err != nil {
        fmt.Println(err)
        fmt.Println("Error in error in opening the file")
        return -1, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    scanner.Scan()
    scanner.Scan()
    
    for scanner.Scan() {
        line := scanner.Text()
        cols := strings.Fields(line)

        x, err := strconv.ParseFloat(cols[0], 64)
        if err != nil {
            fmt.Println("Unable to parse the value of x")
            return -1, err
        }

        if x == target {
            pi, err := strconv.ParseInt(cols[1], 10, 64)
            if err != nil {
                fmt.Println("Unable to parse the value of pi")
                return -1, err
            }
            fmt.Printf("pi value of %v: %v \n", target, pi)
            return pi, nil
        }

        if x > target {
            return -1, errors.New("Value not Found")
        }
    }

    return -1, errors.New("Value not Found")
}

func main() {
    r := gin.Default()
    r.GET("/ping", func (c *gin.Context) {
        c.JSON(200, gin.H {
            "message": "Server is up and running",
        })
    })
    r.GET("/pi", func (c *gin.Context) {
        targetStr := c.Query("target")
        if targetStr == "" {
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H {
                "error": "Missing target paramter", 
            })
            return
        }
        target, err := strconv.ParseFloat(targetStr, 64)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H {
                "error": err.Error(),
            })
            return 
        }
        result, err := parseFile(target)
        if err != nil {
            if err.Error() == "Value not Found" {
                c.AbortWithStatusJSON(http.StatusBadRequest, gin.H {
                    "error": err.Error(),
                })
            } else {
                c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H {
                    "error": "Internal Service Exception",
                })
            }
            return
        }

        c.JSON(http.StatusOK, gin.H {
            "pi": result,
        })
    })
    r.Run()
}
