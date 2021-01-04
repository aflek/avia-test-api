package main

import (
	"net/http"
	"errors"
	"time"
	"log"
	"io/ioutil"
)

//*************************************************
//Функция отправки поискового запроса
//-------------------------------------------------
//nodeKey - ключ
//uri - базовая часть запроса: "http://domain.ru/api/v1/avia/search?"
//paramsStr - параметры запроса
//Ответ содержит: body, время получения ответа, ошибку
//*************************************************
func sendRQ(requestURL, nodeKey string) ([]byte, time.Duration, error) {

	//Засекаем начало времени поиска
	startSearch := time.Now()

	client := &http.Client{}
	request, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		log.Fatal(err)
		return nil, 0, errors.New("Err: Ошибка отправки запроса")
	}

	request.Header.Set("Node-Key", nodeKey)

	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return nil, 0, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return nil, 0, err
	}
	//fmt.Println(string(body))

	endSearch := time.Now()               //Засекаем окончание времени поиска
	elapsed := endSearch.Sub(startSearch) //Считаем время поиска
	
	return body, elapsed, nil 
}