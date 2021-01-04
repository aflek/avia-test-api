package main

import (
	"fmt"
	"log"
	"strconv"
	"errors"
	"encoding/json"
)

//Reccomendation - Блок с реккомендацией
type Reccomendation struct {
	RecID      string  `json:"rec_id"`
	ConfigID   int     `json:"config_id"`
	TotalPrice float32 `json:"total_price"`
	Currency   string  `json:"currency"`
}

//RecResponse - Ответ на запрос - массив реккомендаций
type RecResponse struct {
	Collections []Reccomendation
}
//RecResponseError ...
type RecResponseError struct {
	Error string `json:"error"`
}


//Запуск поискового теста
func aviaAPISearch(params testParams) {
	requestURL, err := getURLForSerachRQ(params.URI, params.Params) 
	if err != nil { log.Fatal(err) }
	
	//Запускае поиск
	body, t, err := sendRQ(requestURL, params.NodeKey)
	if err != nil { log.Fatal(err) }
	//fmt.Println(string(body))
	fmt.Println(t)

	//----------------------------------------------------------------------------
	//Парсинг в труктуру
	var dat []Reccomendation

	if err := json.Unmarshal([]byte(body), &dat); err != nil {
		log.Println(string(body))
		panic(err)
		//fmt.Println(err.Error())
	}
	//Анализ структуры
}

//*************************************************
//Функция формирования строки поискового запроса
//*************************************************
func getURLForSerachRQ(uri, paramsStr string) (string, error) {
	var url string
	//var err error

	//1. Проверка что строка не пустая
	if len(paramsStr) < 1 {
		return url, errors.New("Err: Строка с параметрами поиска пустая")
	}
	//2. Проверяем наличие числа, задающего количество роутов (первая позиция в paramsStr)
	nRoutes, _ := strconv.Atoi(paramsStr[:1])
	if nRoutes == 0 {
		return url, errors.New("Err: Неверно указано количество роутов - " + paramsStr[:1] + ". Должно быть число от 1 и более.")
	}
	//3. Должно быть 1 + (14 символов в описании роута)*число роутов + 3
	var realStrSize = 1 + 14*nRoutes + 3
	if realStrSize != len(paramsStr) {
		return url, errors.New("Err: Строка с параметрами поиска для " + paramsStr[:1] + " роутов должна содержать " + strconv.Itoa(realStrSize) + " символов. У вас их " + strconv.Itoa(len(paramsStr)))
	}
	//4.1 Последний символ в строке (кол-во младенцев), должен быть числовой
	sINF := paramsStr[(len(paramsStr) - 1):]
	if !charIsNum(sINF) {
		return url, errors.New("Err: Неверно указано количество INF - " + sINF + ". Последний символ в конце, должн быть числом.")
	}
	//4.2 Второй с конца символ в строке (кол-во детей), должен быть числовой
	sCHD := paramsStr[(len(paramsStr) - 2):(len(paramsStr) - 1)]
	if !charIsNum(sCHD) {
		return url, errors.New("Err: Неверно указано количество CHD - " + sCHD + ". Второй символ с конца, должн быть числом.")
	}
	//4.3 Третий с конца символ в строке (кол-во взрослых), должен быть числовой
	sADT := paramsStr[(len(paramsStr) - 3):(len(paramsStr) - 2)]
	if !charIsNum(sADT) {
		return url, errors.New("Err: Неверно указано количество ADT - " + sADT + ". Третий символ с конца, должн быть числом.")
	}
	//Цикл по роутам
	var routeStr string = ""
	var routeURL string = ""
	for i := 1; i <= nRoutes; i++ {
		routeStr = paramsStr[((i-1)*14 + 1):(i*14 + 1)] //Кусок с роутом и его датой

		routeURL = routeURL + "destinations[" + strconv.Itoa(i-1) + "][departure]=" + routeStr[0:3]
		routeURL = routeURL + "&destinations[" + strconv.Itoa(i-1) + "][arrival]=" + routeStr[3:6]

		routeDate := routeStr[6:8] + "-" + routeStr[8:10] + "-" + routeStr[10:14]
		//TODO добавить проверку что routeDate - дата

		routeURL = routeURL + "&destinations[" + strconv.Itoa(i-1) + "][date]=" + routeDate + "&"
		//fmt.Println(routeUrl)
	}

	url = uri + routeURL + "adt=" + sADT + "&chd=" + sCHD + "&inf=" + sINF

	//fmt.Println(url)

	return url, nil
}