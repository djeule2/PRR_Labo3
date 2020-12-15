package main

import (
	"./config"
	"./client"
	"./network"
	"./election"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

var waitAllGroupe = &sync.WaitGroup{}

type demandeReqType []int

var demandesreq demandeReqType

//String() Formate la valeur du flag.

func (request *demandeReqType) String() string  {
	return fmt.Sprint(*request)
}

func (r *demandeReqType)Set(s string) error {
	split := strings.Split(s, config.TIRET)

	for _, v := range split {
		i, err := strconv.ParseInt(v, 10, 64)

		if err == nil {
			*r = append(*r, int(i))
		} else {
			return err
		}
	}

	return nil

}

func main()  {
	config.SetConfigNetwork()
	fmt.Println(config.AllNetwork)
	for id, _ := range config.AllNetwork {
		clientElection := make(chan string)
		electionClient := make(chan string)
		electionNetwork := make(chan string)
		networkElection := make(chan string)

		clt := client.NewClient(id, electionClient, clientElection)

		el := election.NewElection(id, clientElection, electionClient,
			networkElection, electionNetwork, *clt)

		net := network.NewNetwork(id, electionNetwork, networkElection, *el)

		waitAllGroupe.Add(1)
		go net.Exec()
	}

	waitAllGroupe.Wait()

}
