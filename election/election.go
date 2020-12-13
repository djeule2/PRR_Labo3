package election

import (
	"../client"
	"../config"
	"../utils"
	"fmt"
	"strings"
	"time"
)

//(Pseudo-code (1-Variables)
type Election struct {
	N		int 	// nombre de processus
	moi 	int 	// mon numéro de processus entre 0 et N-1
	monApt	int	 	// aptitude de moi
	T 		time.Duration	 	// temps maximum d'une transmission
	aps[]	int		// tableau des aptitudes des processus
	enCours bool 	//election en cours ou non
	elu 	int 	// numéro de processus elu

	client client.Client

	//Channel pour communication avec le client et le reseau
	fromClient chan string
	toClient chan string
	fromNetwork chan string
	toNetwork chan string
}

const nameElc = "ELC-"

//NewElection permet d'initialiser un nouveau processus Election à partir de son id,
//par ailleurs  on initialise quatre channels pour communiquer avec le client et avec le reseau.
func NewElection(id int, fromClient chan string, toClient chan string,
	fromNetwork chan string, toNetwork chan string, client client.Client) *Election  {
	election :=new(Election)
	election.N = len(config.AllNetwork)
	election.moi = id
	election.T = config.T
	election.elu = election.moi
	election.monApt = config.AllNetwork[id].Apt
	election.aps = make([]int, 0)
	election.enCours = false
	election.client = client

	election.fromClient = fromClient
	election.toClient = toClient
	election.fromNetwork = fromNetwork
	election.toNetwork = toNetwork
	return election

}

//Pseudo-code(2-boucle principale: traitement d'une reception  à la fois)
func (election *Election) resolveMessage()  {
	for  {
		select {
		case newElection := <-election.fromClient:
			utils.PrintMessage(election.moi, nameElc, "message from client" + newElection)
			if strings.HasPrefix(newElection, config.DEMANDE){
				utils.PrintMessage(election.moi, nameElc, "demande election")
				go election.election()
			}
		}

	}

}


//Pseudo-code(3-réception)
func (election *Election) election()  {
	election.demarre()
}

func (election *Election)demarre()  {
	election.enCours =true
	election.aps[election.moi] = election.monApt
	tabMess := []string{fmt.Sprint(election.moi), fmt.Sprint(election.monApt)}
	messageSent := strings.Join(tabMess, config.TIRET)
	for _, i := range election.aps{
		election.sendDemande(messageSent, i)
	}


}
func (election *Election)sendDemande(mess string, i int)  {
  if i != election.moi {
  	tabMessage := []string{mess, fmt.Sprint(i) }
  	election.sendMessageToNetwork(strings.Join(tabMessage, config.SEP))

  }
}

func (election *Election)sendMessageToNetwork(message string)  {
	utils.PrintMessage(election.moi, nameElc, "send to network : "+message)
	election.toNetwork <- message
}
