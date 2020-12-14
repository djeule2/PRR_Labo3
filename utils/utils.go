package utils

import "fmt"

func PrintMessage(id int, partName string, message string )  {
	fmt.Println("[Pr >>: "+ fmt.Sprint(id) + partName + "] " + message)
}


//Maximum nous retourne le max entre 2 numéro
func Maximum(valX, valY int) int {
	if valX > valY {
		return valX
	}
	return valY
}
