package client

import (
	"../config"
	"../utils"
	"strconv"
	"time"
)

const clientName = "clt"
type Client struct {
	id int
	fromElection chan string
	toElection chan string
	demandeTimes []int
}

//NewClient initialise nouveau client à partir de son id.

func NewClient(id int, fromElection chan string, toElection chan string) *Client {
	client := new(Client)
	client.id = id
	client.toElection = toElection
	client.fromElection = fromElection
	client.demandeTimes = config.AllNetwork[id]
	return client
}

//getElu est une goroutine qui retourne l'election d'un nouveau processus.
func (client *Client)getElu()  {
	for  {
		eluMessage := <-client.fromElection
		utils.PrintMessage(client.id, clientName, " New elu :"+eluMessage)
	}
}

//demande et une goroutine qui permet de faire  les demande de client pour une nouvelle election
func (client *Client) demande() {
	for _, v := range client.demandeTimes {
		time.Sleep(time.Duration(v)*time.Second)
		utils.PrintMessage(client.id, clientName, "Demande election "+strconv.Itoa(v))
		client.toElection <- config.DEMANDE
	}
}
