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
	demandeTimes int
}

//NewClient initialise nouveau client à partir de son id.

func NewClient(id int, fromElection chan string, toElection chan string) *Client {
	client := new(Client)
	client.id = id
	client.toElection = toElection
	client.fromElection = fromElection
	client.demandeTimes = config.AllNetwork[id].Req
	return client
}

func (client *Client) Exec()  {
	go client.getElu()
	go client.demande()
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
	time.Sleep(time.Duration(2)*time.Second)
	for i := 0; i < client.demandeTimes; i++ {
		//time.Sleep(time.Duration(v)*time.Second)
		utils.PrintMessage(client.id, clientName, "Demande election "+strconv.Itoa(i))
		client.toElection <- config.ELECTION
	}
}
