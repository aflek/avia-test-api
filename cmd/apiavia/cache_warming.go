package main

import (
	"os"
	"fmt"
	"encoding/json"
	"log"
	"time"
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
	NodeKey     string `json:"node_key"`
	CaheKey     string `json:"cahe_key"`
	BaseURI     string `json:"base_uri"`
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

	//Цикл по массиву правил в json файле
	for index, rule := range listOfRules {
		fmt.Println("*************************************************")
		fmt.Println(index+1, rule.Descriprion)
		fmt.Println("*************************************************")
		//Цикл по датам вылета туда и обратно
		curDate, _ := strToDate(rule.StartDate); if err != nil { log.Fatal(err) }
		endDate, _ := strToDate(rule.EndDate); if err != nil { log.Fatal(err) }
		for {
			if (curDate.Unix() > endDate.Unix()) { break }

			//Цикл по кадому элементу из durations
			for _, v := range rule.Durations { 
				//Добавляем к текущей дате "вылета туда" текущее число из число из durations и получаем обратную дату 
				dBack := curDate.AddDate(0,0,v)
				//Даты вылетов туа и обратно переводим в строки
				dCurStr  := dateToStr(curDate, "-")
				dBackStr := dateToStr(dBack, "-")
				//Формироуем поисковую строку
				searchStr := rule.BaseURI + "cache_search?"
				searchStr = searchStr + "&destinations[0]departure=" + rule.From
				searchStr = searchStr + "&destinations[0]arrival=" + rule.To
				searchStr = searchStr + "&destinations[0]date=" + dCurStr
				searchStr = searchStr + "&destinations[1]departure=" + rule.To
				searchStr = searchStr + "&destinations[1]arrival=" + rule.From
				searchStr = searchStr + "&destinations[1]date=" + dBackStr
				searchStr = searchStr + "&adt=2"
				fmt.Println(searchStr)
				fmt.Println("------------------------")
				//Отправка запроса в кэш
				bodyCache, tCache, err := sendRQ(searchStr, rule.CaheKey); if err != nil { log.Fatal(err) }
				//fmt.Println(cacheRS)
				//Парсинг в труктуру
				var datCache []Reccomendation

				if err := json.Unmarshal([]byte(bodyCache), &datCache); err != nil {
					log.Println(string(bodyCache))
					panic(err)
					//fmt.Println(err.Error())
				}
				//Анализ структуры
				fmt.Println(datCache)
				fmt.Println(len(datCache))
				fmt.Println(tCache)
				fmt.Println("------------------------")
				//Если не надено, то ищем в ГДСах
				if len(datCache) == 0 {
					fmt.Println("В кэше не надено. Ищем в ГДСах")
					//Формироуем поисковую строку
					searchStr := rule.BaseURI + "search?"
					searchStr = searchStr + "&destinations[0]departure=" + rule.From
					searchStr = searchStr + "&destinations[0]arrival=" + rule.To
					searchStr = searchStr + "&destinations[0]date=" + dCurStr
					searchStr = searchStr + "&destinations[1]departure=" + rule.To
					searchStr = searchStr + "&destinations[1]arrival=" + rule.From
					searchStr = searchStr + "&destinations[1]date=" + dBackStr
					searchStr = searchStr + "&adt=2"
					fmt.Println(searchStr)
					//Отправка запроса в ГДС
					bodyGDS, tGDS, err := sendRQ(searchStr, rule.CaheKey); if err != nil { log.Fatal(err) }
					//Парсинг в труктуру
					var datGDS []Reccomendation
					if err := json.Unmarshal([]byte(bodyGDS), &datGDS); err != nil {
						log.Println(string(bodyGDS))
						panic(err)
						//fmt.Println(err.Error())
					}
					//Анализ структуры
					fmt.Println(datGDS)
					fmt.Println(len(datGDS))
					fmt.Println(tGDS)
				}
				//Пауза перед следующей отправкой запроса
				time.Sleep(5 * time.Second)
			}

			curDate = curDate.AddDate(0,0,1)//Сдивигаем дату "вылета туда" на один день вперед
		}

	/*	
		//Сформировать массив со строками типа "2MOWDXB20012021DXBMOW27012021200" - параметрами запроса к 
		arrSearchParams, err := getSearchParams(rule.From, rule.To, rule.StartDate, rule.EndDate, rule.Durations)
		if err != nil { log.Fatal(err) }
		//fmt.Println(arrSearchParams)
		//Поиск по кэшу
		uri := rule.BaseURI + "cache_search?"
		fmt.Println(uri)
		for _, v := range arrSearchParams {

			fmt.Println(v)
		}
		//params := arrSearchParams
	*/	
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

/*
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
			searchStr := "2" + from + to + dStr + to + from + dBackStr + "200"
			//fmt.Println(searchStr)
			arrSearchParams = append(arrSearchParams, searchStr)
		}	
		  
		
			
		d = d.AddDate(0,0,1)	
	}

	return arrSearchParams, nil
}

*/