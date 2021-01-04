package main

import (
	"fmt"
	"flag"
)

var listOfTests []testParams
//ListOfTests - Структура списка тестов
type testParams struct {
	NodeKey     string `json:"node_key"`
    URI         string `json:"uri"`
    Params      string `json:"params"`
    Descriprion string `json:"descriprion"`
    RqCode     string  `json:"rq_code"`
    AvailConfig int    `json:"avail_config"`
}


func main() {
	
	


	flagAction := flag.String("a", "help", "action")
	flag.Parse()
	fmt.Println("flagAction:", *flagAction)
	switch *flagAction {
	case "cw"://Прогрев кэша
		fmt.Println("cache warming")
	case "test"://Выполнение тестов по скритам в файлах каталога tests
		fmt.Println("tests")
	}

	
	 /*1. Генерации скриптов (файла в каталоге tests).
	Параметры для генерации находятся в каталоге tests_tpl, в файле cache_warming.json.
	Параметры: 
		- Название файла скрипта, который будет записан в каталог tests.
		- Коды городов вылета и прилета
		- Стартовая дата для вылета
		- Конечная дата для вылета
		- Продолжительность нахождения в городе прилета - массив из чисел
	Процесс генерации запускается командой (cw - cache warming): ./apiavia -a=cw 	
	*/ 
	/*
	  2.1. В конфигурационном файле должны быть реквизиты доступа к БД Монго.
	  2.2. В конфигурационном файле должн быть ключ поиска по кэшу.
	  2.3. В конфигурационном файле должн быть ключ поиска авиа.
	*/
	/*
	  3. Код тестирования CACHE_WARMING 
	  Для этого кода тестирования в файле скрипта должен быть параметр проверки данных в писковом кэше.
	  Если параметр поиска в кэше = true, то делаем сначала поиск в кэше.
	  Если по кэшу поиск неудачный - ищем в авиа (прогреваем кэш).
	*/

	
	
}







