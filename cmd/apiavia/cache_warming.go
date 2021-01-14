package main

import (
	"os"
	"fmt"
	"encoding/json"
	"log"
)

var listOfRules []Rules
//Rules ...
type Rules struct {
	Descriprion string `json:"descriprion"`
	From        string `json:"from"`
	To          string `json:"to"`
	StartDate   string `json:"start_date"`
	EndDate	    string `json:"end_date"`
	Durations   []int  `json:"durations"`
}

// cach_warming -...
func cachWarming() {
 /*1. Генерации скриптов (файла в каталоге test).
	Параметры для генерации находятся в каталоге tests_tpl, в файле cache_warming.json.
	Параметры: 
		- descriprion: Название файла скрипта, который будет записан в каталог tests.
		- from:        IATA-Код города вылета
		- to:          IATA-Код города прилета
		- start_date:  Стартовая дата для вылета
		- end_date:    Конечная дата для вылета
		- durations:   Продолжительность нахождения в городе прилета - массив из чисел, например, [4,5,6,7]
		- node_key:    Ключ для прямого поиска
		- cahe_key:    Ключ для поиска по кэша
	*/
	file, err := os.Open("./tests_tpl/cache_warming.json")
	defer file.Close()
	if err != nil {
		fmt.Println("Open file fatal error ", err.Error())
		os.Exit(1)
	}
	//Парсим файл в массив
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&listOfRules)
	//fmt.Println(listOfRules)

	//Цикл по массиву и отправка запросов
	for index, rule := range listOfRules {
		fmt.Println("*************************************************")
		fmt.Println(index+1, rule.Descriprion)
		fmt.Println("*************************************************")
		//Сформировать массив со строками типа "2MOWDXB20012021DXBMOW27012021200" - параметрами запроса к 
		arrSearchParams, err := getSearchParams(rule.From, rule.To, rule.StartDate, rule.EndDate, rule.Durations)
		if err != nil { log.Fatal(err) }
		fmt.Println(arrSearchParams)
		//Поиск по кэшу
		//Если по кэшу ничего не найдено, то выполняем поиск через ГДСы 
		//Запись результатов тестирования в файл
	}

	
	/*
	  3. Код тестирования CACHE_WARMING 
	  Для этого кода тестирования в файле скрипта должен быть параметр проверки данных в писковом кэше.
	  Если параметр поиска в кэше = true, то делаем сначала поиск в кэше.
	  Если по кэшу поиск неудачный - ищем в авиа (прогреваем кэш).
	*/
}

func getSearchParams(from, to, startDate, endDate string, durations []int) ([]string, error) {
	var arrSearchParams []string
	d, _ := strToDate(startDate)
	dEndDate, _ := strToDate(endDate)
	
	//Цикл от startDate до endDate с шагом 1 день - дата вылета туда
	for {
		if (d.Unix() > dEndDate.Unix()) { break }
		
		//Цикл по кадому элементу из durations
		for _, v := range durations {
			//fmt.Println(v)
			//Добавляем к текущей дате "вылета туда" текущее число из число из durations и получаем обратную дату 
			dBack := d.AddDate(0,0,v)
			//Даты вылетов туа и обратно переводим в строки
			dStr := dateToStr(d)
			dBackStr := dateToStr(dBack)
			//Формироуем поисковую строку
			searchStr := "2" + from + to + dStr + to + from + dBackStr
			//fmt.Println(searchStr)
			arrSearchParams = append(arrSearchParams, searchStr)
		}	
		  
		
			
		d = d.AddDate(0,0,1)	
	}

	return arrSearchParams, nil
}

