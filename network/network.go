package network

import (
	"../config"
	"../election"
	"../utils"
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)
const nameNet = " net "

type Network struct {
	id 	int
	enPanne bool
	enCours bool

	fromElection chan string
	toElection chan string
	fromNetwork chan string
	toNetwork chan string

	election election.Election
	address map[int]string
	conns map[int]net.Conn

}

func NewNetwork(id int, fromElection chan string, toElection chan string, election election.Election) *Network {
	network := new(Network)
	network.id = id
	network.fromElection = fromElection
	network.toElection = toElection
	network.fromNetwork = make(chan string)
	network.toNetwork = make(chan string)
	network.election = election
	network.conns = make(map[int]net.Conn)

	return network
}

func (network *Network)Exec()  {

	network.connect()
	go network.election.Exec()
	go network.getMessage()
	go network.listenConn()
}

// la fonction connect lance le processus de connection d'un processus à tous les autres processus définit dans le systeme.
func (network *Network) connect() {
	utils.PrintMessage(network.id, nameNet, "Network connecting to others processus")
	for k, v := range config.AllNetwork{
		if k != network.id{
			address := v.Host + ":" + fmt.Sprint(v.Port)
			network.conns[k] = network.connectTo(address)
		}
	}

}

//la fonction connecTo permet de connecter un processus à un autre spécique défini par l'adresse passé en paramètre
func (network *Network)connectTo(adresse string) net.Conn {
	conn, err := net.Dial("udp", adresse)

	utils.PrintMessage(network.id, nameNet, " Connected to: "+adresse)

	if err != nil {
		log.Fatal(err)
	}

	return conn
}

func (network *Network) getMessage() {
	for{
		select {
			//Gère les arrivé des message du processus election
			//les message sont ensuite transmis aux autre site

			case msg := <- network.fromElection:
				utils.PrintMessage(network.id, nameNet, "message received from mutex : "+msg)
				utils.PrintMessage(network.id, nameNet, "Sending message to other pocessus")
				msgsplit := strings.Split(msg, config.TIRET)
				id, err := strconv.Atoi(msgsplit[0])
				if err == nil{
					for i := 0; i< len(config.AllNetwork); i++{
						if i != id{
							network.sendMessage(network.conns[i], msg)
						}
					}
				}
			case msg := <- network.fromNetwork:
				utils.PrintMessage(network.id, nameNet, "Sending message to processus Election")
				go func() {network.toElection <- msg}()
		}

	}

}

func (network *Network)sendMessage(conn net.Conn, msg string)  {
	_, err := conn.Write([]byte(msg))

	utils.PrintMessage(network.id, nameNet, "Sending message to "+ fmt.Sprint(conn.RemoteAddr())+" : "+msg)

	if err != nil {
		log.Fatal(err)
	}
}

//la fonction ListenConn permetra écouter le message qui arrive à notre adresse.
// pour chaque message on l'envois au processus election
func (network *Network) listenConn() {
	me := config.AllNetwork[network.id]
	addr := me.Host + ":" + fmt.Sprint(me.Port)


	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	utils.PrintMessage(network.id, nameNet,
		"Listenting to "+fmt.Sprint(conn.LocalAddr()))

	buffer := make([]byte, 4096)

	for {
		n, dest, err := conn.ReadFrom(buffer)
		if err != nil {
			log.Fatal(err)
		}

		 time.Sleep(config.T * time.Millisecond)
		s := bufio.NewScanner(bytes.NewReader(buffer[0:n]))

		for s.Scan() {
			utils.PrintMessage(network.id, nameNet, "Received message : "+
				s.Text()+" from "+dest.String())
			network.fromNetwork <- s.Text()

		}
	}
}

