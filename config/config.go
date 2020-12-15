package config

import (
	"github.com/tkanos/gonfig"
)

type PanneType struct {
	Start 	int
	End 	int
}

type Sites struct {
	Host 	string
	Port 	int
	Apt 	int
	Req []int
	Pannes []PanneType
}

var AllNetwork = make(map[int]Sites)

type NetworkType struct {
	Sites []Sites
}

func SetConfigNetwork()  {
	configNetwork := NetworkType{}

	err := gonfig.GetConf("config/configNetworkBully.json",
		&configNetwork)

	if err != nil {
		panic(err)
	}

	for k, v := range configNetwork.Sites{
		AllNetwork[k] = v
	}

}

