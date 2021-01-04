package main

//*************************************************
//Функция проверки: Символ есть число
//*************************************************
func charIsNum(munChar string) bool {
	
	if (munChar == "0") || (munChar == "1") || (munChar == "2") || (munChar == "3") || (munChar == "4") || (munChar == "5") || (munChar == "6") || (munChar == "7") || (munChar == "8") || (munChar == "9") {
		return true
	}
	return false
}