package config

type PanneType struct {
	Start 	int
	End 	int
}

type ElectionType struct {
	Host 	string
	Port 	int
	Apt 	int
	Req []int
	Pannes []PanneType
}

var AllNetwork = make(map[int]ElectionType)

type NetworkType struct {
	Elections []ElectionType
}

func SetConfigNetwork()  {
	confiNetwork := NetworkType{}

	for k, v := range confiNetwork.Elections{
		AllNetwork[k] = v
	}

}

