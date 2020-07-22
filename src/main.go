package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"

	"encoding/csv"
	"os"
)

var errRequestFailed = errors.New("Rewquest Failed")
var data = [][]string{}
var keyValueData = map[int][]string{}
var workedInex = 0
var stopGoRoutine bool = false

const (
	maxWorkIndex = 2
)

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func main() {

	dataName := []string{"날짜",
		"종가",
		"전일비",
		"시가",
		"고가",
		"저가",
		"거래량",
	}

	data = append(data, dataName)
	os.Remove("result.csv")
	file, err := os.Create("result.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	errChan := make(chan error)
	mainURL := "https://finance.naver.com/item/sise_day.nhn?code=005930"
	maxPage := 100
	currentIndex := 0

	for i := 1; i < maxPage; i++ {
		newPage := fmt.Sprint(mainURL+"&page=", i)
		fmt.Println(newPage)
		go hitURL(newPage, errChan)
		currentIndex++

		if stopGoRoutine {
			break
		}
	}

	for i := 0; i < currentIndex; i++ {
		fmt.Println("err : ", <-errChan)
	}

	//sort
	keys := make([]int, 0, len(keyValueData))
	for k := range keyValueData {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	//fmt.Println(keys)

	//fmt.Println("data , : ", data)

	for _, k := range keys {

		data = append(data, keyValueData[k])
		//fmt.Println(k, keyValueData[k])
	}

	for _, value := range data {
		err := writer.Write(value)
		checkError("Cannot write to file", err)
	}

	fmt.Scanf("%d")
}

func hitURL(url string, c chan error) {
	workedInex++
	scrapeURL(url, c)
	workedInex--
	fmt.Println("-- worked index")
	c <- nil
	return
}

func scrapeURL(url string, c chan error) {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	utfBody, err := iconv.NewReader(res.Body, "euc-kr", "utf-8")
	if err != nil {
		// handler error
		c <- errRequestFailed
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		log.Fatal(err)
	}

	/*
		<th>날짜</th>
		<th>종가</th>
		<th>전일비</th>
		<th>시가</th>
		<th>고가</th>
		<th>저가</th>
		<th>거래량</th>
	*/

	//findkey: dl.blind
	oneData := make([]string, 0)
	index := 0
	doc.Find("span.tah").Each(func(i int, s *goquery.Selection) {
		filteredDay := s.Filter("span.tah.p10.gray03").Text()

		if filteredDay != "" {
			if len(oneData) > 0 {
				//strDateInt := strings.Replace(filteredDay, ".", "", -1)
				//fmt.Println("strDateInt , ", strDateInt)
				//dateKey, _ := strconv.Atoi(strDateInt)
				//keyValueData[dateKey] = oneData
				//data = append(data, oneData)
			}

			oneData = make([]string, 7)
			index = 0
			oneData[index] = filteredDay
		}

		filteredSome := s.Filter("span.tah.p11").Text()
		if filteredSome != "" {
			filteredSome = strings.TrimSpace(filteredSome)
			filteredSome = strings.Replace(filteredSome, ",", "", -1)
			pasredInt, _ := strconv.Atoi(filteredSome)

			//check is up or down? only check down
			band, ok := s.Attr("class")
			if ok {
				isDown := strings.Contains(band, "nv")
				if isDown {
					pasredInt *= -1
				}
			}
			index++
			oneData[index] = strconv.Itoa(pasredInt)
		}
	})

	c <- nil
}
