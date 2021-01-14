package main

import (
	"time"
	"errors"
	"fmt"
	"log"
	"strconv"
)

//*************************************************
//Функция проверки: Символ есть число
//*************************************************
func charIsNum(munChar string) bool {
	
	if (munChar == "0") || (munChar == "1") || (munChar == "2") || (munChar == "3") || (munChar == "4") || (munChar == "5") || (munChar == "6") || (munChar == "7") || (munChar == "8") || (munChar == "9") {
		return true
	}
	return false
}

//СharIsNum - функция проверки: Символ есть число
func СharIsNum(s byte) bool {
	
	if s> 47 && s < 58 {
		return true
	}
	return false
}

//StrIsNum - функция проверки: Срока состоит только из чисел
func StrIsNum(s string) bool {
    for i:=0; i<len(s); i++ {
		c:= s[i]
        if !СharIsNum(c) {
            return false
        }
    }
    return true
}

//StrToTimeMonth ...
func StrToTimeMonth(mm string) (time.Month, error) {

	var m time.Month
	switch mm {
	case "01":	
		m = time.January
	case "02":	
		m = time.February
	case "03":
		m = time.March
    case "04":
		m = time.April
	case "05":
		m = time.May
	case "06":
		m = time.June
	case "07":
		m = time.July
	case "08":
		m = time.August
	case "09":
		m = time.September
	case "10":
		m = time.October
	case "11":
		m = time.November
	case "12":
		m = time.December
	default:
		return time.January, errors.New("Неверный номер месяца: '" + mm)
	}

	return m, nil
}

//Преобразование даты в строку
func dateToStr(dt time.Time) string {

	d := dt.Day()
	m := int(dt.Month())
	y := dt.Year()
	//return fmt.Sprintf("%v"+"%v"+"%v", d, m, y)
	return intToDateStr(d) + intToDateStr(m) + fmt.Sprintf("%v", y)
}

func intToDateStr(d int) string {
	if d < 10 {
		return fmt.Sprintf("0%v", d)
	}

	return fmt.Sprintf("%v", d)
}

//strToDate(23122020) вернет дату в формате Time
func strToDate(dt string) (time.Time, error) {
	//1.Строка должна иметь длину 10 сиволов
	strLen := len(dt)
	if strLen != 8 {
		return time.Now(), errors.New("Строка с датой должна быть длиной в 8 символов, включая разделитель, например: 23122021. \nДлина " + dt + " составляет " + strconv.Itoa(strLen) + "симоволов")
	}
	//2. Парсим строку
	dd := string(dt[0]) + string(dt[1])//Дата
	mm := string(dt[2]) + string(dt[3])//Месяц
	yyyy := string(dt[4]) + string(dt[5])+ string(dt[6]) + string(dt[7])//Год
	//3.Проверка дня - должны быть числа и их должно быть 2.
	if !StrIsNum(dd) {
		return time.Now(), errors.New("День '" + dd + "' не является числом")
	}
	//4.Проверка месяца - должны быть числа и их должно быть 2.
	if !StrIsNum(mm) {
		return time.Now(), errors.New("Месяц '" + mm + "' не является числом")
	}
	//5.Проверка года - должны быть числа и чисел должно быть 4 шт.
	if !StrIsNum(yyyy) {
		return time.Now(), errors.New("Год '" + yyyy + "' не является числом")
	}

	//6.Формируем дату с типом Time
	d, err := strconv.Atoi(dd); if err != nil { log.Fatal(err) }//День из string в int

	m, err := StrToTimeMonth(mm); if err != nil { log.Fatal(err) }//Месяц из string в time.Month
	
	y, err := strconv.Atoi(yyyy); if err != nil { log.Fatal(err) } //Год из string в int

	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC), nil

}