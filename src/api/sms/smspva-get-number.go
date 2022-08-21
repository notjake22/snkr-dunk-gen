package sms

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"io/ioutil"
	"log"
	"main/src/settings"
	"net/http"
	"strconv"
	"time"
)

type Options struct {
	BearerHeader string `json:"bearerHeader"`
	ID           string `json:"id"`
	OptionNumber string `json:"optionNumber"`
}

func GetNewNumber() (string, Options, error) {
	apiKey := settings.ReadSettings().SmsAPIKey
	log.Println("Checking carriers")
	option := getOption(apiKey)

	o := Options{
		OptionNumber: option,
	}
	if len(option) < 0 {
		log.Println(option)
		sentry.CaptureMessage("No numbers found at price SmsPva")
		return "", o, errors.New("no numbers found at price")
	}
	number, id, err := genNumber(option, apiKey)
	if err != nil {
		sentry.CaptureMessage("Error getting number SmsPva")
		return "", o, err
	}

	o.ID = id
	return number, o, err
}

func getOption(apiKey string) string {
	for i := 1; i > 0; i++ {
		uri := fmt.Sprintf("https://smspva.com/priemnik.php?metod=get_service_price&country=US&service=opt%s&apikey=%s", strconv.Itoa(i), apiKey)
		req, err := http.NewRequest("GET", uri, nil)
		if err != nil {
			sentry.CaptureException(err)
			log.Fatalln("Error making req")
		}

		req.Header.Set("Accept", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			sentry.CaptureException(err)
			log.Fatalln("Error sending req")
		}
		if res.StatusCode != 200 {
			sentry.CaptureException(err)
			log.Fatalln("Request error")
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			sentry.CaptureException(err)
			log.Fatalln("Error reading body response")
		}

		//fmt.Printf("%+v", string(body))

		var response OptionResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			sentry.CaptureException(err)
			log.Fatalln("Error decoding response")
		}

		respPrice, err := strconv.ParseFloat(response.Price, 64)
		if err != nil {
			sentry.CaptureException(err)
			log.Fatalln("Error parsing price 1")
		}

		maxPrice, err := strconv.ParseFloat("1.0", 64)
		if err != nil {
			sentry.CaptureException(err)
			log.Fatalln("Error parsing price 2")
		}

		if respPrice <= maxPrice {
			return strconv.Itoa(i)
		} else {
			time.Sleep(1500 * time.Millisecond)
			continue
		}
	}
	//options := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25"}
	//for _, optNumber := range options {
	//	uri := fmt.Sprintf("https://smspva.com/priemnik.php?metod=get_service_price&country=US&service=opt%s&apikey=%s", optNumber, apiKey)
	//	req, err := http.NewRequest("GET", uri, nil)
	//	if err != nil {
	//		log.Fatalln("Error making req")
	//	}
	//
	//	req.Header.Set("Accept", "application/json")
	//	res, err := http.DefaultClient.Do(req)
	//	if err != nil {
	//		log.Fatalln("Error sending req")
	//	}
	//	if res.StatusCode != 200 {
	//		log.Fatalln("Request error")
	//	}
	//
	//	body, err := ioutil.ReadAll(res.Body)
	//	if err != nil {
	//		log.Fatalln("Error reading body response")
	//	}
	//
	//	//fmt.Printf("%+v", string(body))
	//
	//	var response OptionResponse
	//	err = json.Unmarshal(body, &response)
	//	if err != nil {
	//		log.Fatalln("Error decoding response")
	//	}
	//
	//	respPrice, err := strconv.ParseFloat(response.Price, 64)
	//	if err != nil {
	//		log.Fatalln("Error parsing price")
	//	}
	//
	//	maxPrice, err := strconv.ParseFloat(settings.ReadSettings().SmsPriceMax, 64)
	//	if err != nil {
	//		log.Fatalln("Error parsing price")
	//	}
	//
	//	if respPrice <= maxPrice {
	//		return optNumber
	//	} else {
	//		time.Sleep(1500 * time.Millisecond)
	//		continue
	//	}
	//}
	return ""
}

func genNumber(optNumber string, apiKey string) (string, string, error) {
	uri := fmt.Sprintf("https://smspva.com/priemnik.php?metod=get_number&country=US&service=opt%s&apikey=%s", optNumber, apiKey)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Println("Error making request")
		sentry.CaptureException(err)
		return "", "", err
	}

	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error sending request")
		sentry.CaptureException(err)
		return "", "", err
	}
	if res.StatusCode != 200 {
		return "", "", errors.New("request error")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading body")
		sentry.CaptureException(err)
		return "", "", err
	}

	var response NumberResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println("Error decoding body")
		fmt.Printf("%+v \n", string(body))
		sentry.CaptureException(err)
		msg := fmt.Sprintf("Error decoding body SmsPva: %+v", string(body))
		sentry.CaptureMessage(msg)
		return "", "", err
	}

	if response.Response != "1" {
		if response.Response == "2" {
			if len(response.Balance) > 1 {
				balance, err := strconv.ParseFloat(response.Balance, 64)
				if err != nil {
					sentry.CaptureException(err)
					log.Fatalln("Error parsing price 3: ", err)
				}
				maxPrice, err := strconv.ParseFloat(settings.ReadSettings().SmsPriceMax, 64)
				if err != nil {
					sentry.CaptureException(err)
					log.Fatalln("Error parsing price 4")
				}
				if balance < maxPrice {
					log.Println("SMS balance too low, balance: ", response.Balance)
					sentry.CaptureMessage("User Sms balance lower than max price SmsPva")
					return "", "", errors.New("SMS balance lower than max price")
				} else {
					log.Println("All numbers taken in this price option")
				}
			}
			log.Println("All numbers taken in this price option")
		}
		sentry.CaptureMessage("Api error getting number (last position) SmsPva")
		return "", "", errors.New("api error getting number")
	}

	return response.Number, strconv.Itoa(response.ID), nil
}

type OptionResponse struct {
	Response string `json:"response"`
	Country  string `json:"country"`
	Service  string `json:"service"`
	Price    string `json:"price"`
}

type NumberResponse struct {
	Response string `json:"response"`
	Number   string `json:"number"`
	Balance  string `json:"balance"`
	ID       int    `json:"id"`
}
