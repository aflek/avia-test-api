package main

import (
	//"fmt"
	"flag"
)

func main() {
	flagAction := flag.String("a", "help", "action")
	
	flag.Parse()
	//fmt.Println("flagAction:", *flagAction)
	switch *flagAction {
	case "cw":
		cachWarming()//Прогрев кэша
	case "test":
		test()//Выполнение тестов по скритам в файлах каталога tests
	}
}







