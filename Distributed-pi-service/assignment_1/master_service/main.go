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
	} else if serverCount == 16 {
		if x >= 1 && x <= 4052740000 {
			serverNumber = 1
		} else if x >= 4052750000 && x <= 34846100000 {
			serverNumber = 2
		} else if x >= 34846200000 && x <= 81574500000 {
			serverNumber = 3
		} else if x >= 81574600000 && x <= 383028000000 {
			serverNumber = 4
		} else if x >= 383029000000 && x <= 850312000000 {
			serverNumber = 5
		} else if x >= 850313000000 && x <= 4175950000000 {
			serverNumber = 6
		} else if x >= 4175960000000 && x <= 8848790000000 {
			serverNumber = 7
		} else if x >= 8848800000000 && x <= 45216200000000 {
			serverNumber = 8
		} else if x >= 45216300000000 && x <= 439360000000000 {
			serverNumber = 9
		} else if x >= 439361000000000 && x <= 906643000000000 {
			serverNumber = 10
		} else if x >= 906644000000000 && x <= 4739250000000000 {
			serverNumber = 11
		} else if x >= 4739260000000000 && x <= 9412080000000000 {
			serverNumber = 12
		} else if x >= 9412090000000000 && x <= 50849000000000000 {
			serverNumber = 13
		} else if x >= 50849100000000000 && x <= 97577300000000000 {
			serverNumber = 14
		} else if x >= 97577400000000000 && x <= 543055000000000000 {
			serverNumber = 15
		} else {
			serverNumber = 16
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
		result := fmt.Sprintf("%v                                             %v", inputNumber, piResult.Pi)
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
		fmt.Printf("Completed %v reqs \n", counter)
	}

	t1Master := time.Since(t0Master)
	fmt.Printf("Application took %v to seconds to server the request \n", t1Master)
	fmt.Printf("Average time to serve the requests is %v", t1Master/1000)

	time.Sleep(30 * time.Minute)
}
