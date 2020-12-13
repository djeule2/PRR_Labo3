package utils

import "fmt"

func PrintMessage(id int, partName string, message string )  {
	fmt.Println("[Pr >>: "+ fmt.Sprint(id) + partName + "] " + message)
}
