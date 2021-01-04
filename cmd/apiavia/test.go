package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"encoding/json"
)

func test() {
	//Перебор файлов (сценариев тестирования) в папке tests
	//Тестировани запускается командой ./apiavia -a=test
	files, err := ioutil.ReadDir("./tests")
	if err != nil { log.Fatal(err) }

	for _, file := range files {
		fileListOfTests, err := os.Open("./tests/" + file.Name())
		defer fileListOfTests.Close()
		if err != nil {
			fmt.Println("Open file fatal error ", err.Error())
			os.Exit(1)
		}
		//Парсим файл в массив
		jsonParser := json.NewDecoder(fileListOfTests)
		err = jsonParser.Decode(&listOfTests) //пишем данные в массив

		//Цикл по массиву и отправка запросов
		for index, element := range listOfTests {
			fmt.Println(" ------------------------------------------------------------------- ")
			fmt.Print(index+1); fmt.Println(". " + element.Descriprion)
			runTest(element)
		}
	}
}

//Функция запуска теста
func runTest(params testParams) {
	fmt.Println(params.RqCode)
	switch params.RqCode {
	case "AVIA_API_SEARCH":
		aviaAPISearch(params)
	default:
		fmt.Println("Неизветный ключ теста", params.RqCode)
	}
}