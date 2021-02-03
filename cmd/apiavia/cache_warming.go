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
	SleepTime   int    `json:"sleep_time"`
}

//Color подсветка в терминале задаваемая с использованием управляющей последовательности ANSI
// https://en.wikipedia.org/wiki/ANSI_escape_code
type Color string
const (
     //ColorBlack черная подсветка в терминале
    ColorBlack Color = "\u001b[30m"
    //ColorRed красная подсветка в терминале
    ColorRed         = "\u001b[31m"
    //ColorGreen зеленая подсветка в терминале
    ColorGreen       = "\u001b[32m"
    //ColorYellow желтая подсветка в терминале
    ColorYellow      = "\u001b[33m"
    //ColorBlue синяя подсветка в терминале
    ColorBlue        = "\u001b[34m"
    //ColorReset сброс подсветки в терминале
    ColorReset       = "\u001b[0m"
)


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
		nTestTotal := 0//Общее количество запросов.
		nTestCache := 0//Количество успешных запросов к кэш 
		nTestGDS   := 0//Количество запросов к GDS 
		for {
			if (curDate.Unix() > endDate.Unix()) { break }
			//Цикл по кадому элементу из durations
			for _, v := range rule.Durations {
				nTestTotal++ 
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
				//fmt.Println(searchStr)
				fmt.Println("------------------- ", time.Now(), " -------------------------- ")
				fmt.Print(rule.From, "-", rule.To, " вылет с " ,dCurStr, " по 	", dBackStr, " /")
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
				//fmt.Println(datCache)
				fmt.Println("Найдено в кэше рекомендаций:", len(datCache), "Время поиска в КЭШе:",tCache)
				//Если не надено, то ищем в ГДСах
				if len(datCache) == 0 {
					nTestGDS++
					fmt.Print(string(ColorBlue), "В кэше не надено. Ищем в ГДСах / ", string(ColorReset))
					//Формироуем поисковую строку
					searchStr := rule.BaseURI + "search?"
					searchStr = searchStr + "&destinations[0]departure=" + rule.From
					searchStr = searchStr + "&destinations[0]arrival=" + rule.To
					searchStr = searchStr + "&destinations[0]date=" + dCurStr
					searchStr = searchStr + "&destinations[1]departure=" + rule.To
					searchStr = searchStr + "&destinations[1]arrival=" + rule.From
					searchStr = searchStr + "&destinations[1]date=" + dBackStr
					searchStr = searchStr + "&adt=2"
					//fmt.Println(searchStr)
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
					//fmt.Println(datGDS)
					fmt.Print(string(ColorBlue), "Найдено в GDS рекомендаций:", len(datGDS), string(ColorReset))
					fmt.Println(string(ColorBlue), " Время поиска в GDS:", tGDS, string(ColorReset))
				} else {
					nTestCache++	
				}
				fmt.Println("Выполнено всего запросов:", nTestTotal, "Выполнено успешных запросов в КЭШ:", nTestCache, "Выполнено запросов к ГДС:", nTestGDS)
				fmt.Println("Заполнение кэша %:", nTestCache*100/nTestTotal)
				//Пауза перед следующей отправкой запроса
				time.Sleep(time.Duration(rule.SleepTime) * time.Second)
			}

			curDate = curDate.AddDate(0,0,1)//Сдивигаем дату "вылета туда" на один день вперед
		}
		fmt.Println("**************************************************")
		fmt.Println("Выполнено всего запросов:", nTestTotal)
		fmt.Println("Выполнено успешных запросов в КЭШ:", nTestCache)
		fmt.Println("Выполнено запросов к ГДС:", nTestGDS)
		fmt.Println("**************************************************")
	}

}

