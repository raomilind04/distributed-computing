package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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

func appendToFile(filename, entry string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	_, err = w.WriteString(entry + "\n")
	if err != nil {
		return err
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	return nil
}

func handleRequest(inputNumber float64) {
	serverNumber := loadBalancingFunction(inputNumber)
	url := fmt.Sprintf("http://pi_api_%d:8080/pi?target=%f", serverNumber, inputNumber)

	res, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	var piResult piApiResponse

	if res.StatusCode != http.StatusOK {
		if err := json.NewDecoder(res.Body).Decode(&piResult); err != nil {
			fmt.Println(err)
			return
		}
		appendToFile("output.txt", piResult.Err)
	} else {
		if err := json.NewDecoder(res.Body).Decode(&piResult); err != nil {
			fmt.Println(err)
			return
		}
		result := fmt.Sprintf("%d", piResult.Pi)
		appendToFile("output.txt", result)
	}
}

func main() {

	inputFile, err := os.Open("input.txt")
	if err != nil {
		panic("Unable to open input file")
	}
	scanner := bufio.NewScanner(inputFile)

	t0Master := time.Now()

	counter := 0

	for scanner.Scan() {
		counter++
		line := scanner.Text()
		numbers := strings.Fields(line)

		inputNumber, err := strconv.ParseFloat(numbers[0], 64)
		if err != nil {
			fmt.Print("Could not read the value")
			continue
		}
		handleRequest(inputNumber)
		fmt.Printf("Completed %v reqs", counter)
	}

	t1Master := time.Since(t0Master)
	fmt.Printf("Application took %v to seconds to server the request \n", t1Master)
	fmt.Printf("Average time to serve the requests is %v", t1Master/1000)
}
