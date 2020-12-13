package config

type PanneType struct {
	Start 	int
	End 	int
}

type ElectionType struct {
	Host 	string
	Port 	int
	Apt 	int
	tabReq []int
	Pannes []PanneType
}

var AllNetwork = make(map[int]ElectionType)

type NetworkType struct {
	Elections []ElectionType
}

