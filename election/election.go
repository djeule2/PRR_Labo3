package election

import (
	"../client"
	"../config"
	"../utils"
	"fmt"
	"strconv"
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
	endTime <- chan time.Time
	fromClient chan string
	toClient chan string
	fromNetwork chan string
	toNetwork chan string
}

const nameElc = " elc "

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

func (election *Election)Exec()  {
	go election.resolveMessage()
	go election.client.Exec()
}

//Pseudo-code(2-boucle principale: traitement d'une reception  à la fois)
func (election *Election) resolveMessage()  {
	for  {
		select {
		case newElection := <-election.fromClient:
			utils.PrintMessage(election.moi, nameElc, "message from client" + newElection)
			if strings.HasPrefix(newElection, config.ELECTION){
				utils.PrintMessage(election.moi, nameElc, "demande election")
				for election.enCours {
					//attente si une élection du processus est déjà en cpours.
					//lorsque c'est fini on peut faire la démandde d'une nouvelle election
				}
				go election.election()
			}
		case netMsg := <- election.fromNetwork:
			utils.PrintMessage(election.moi, nameElc, "from network "+netMsg)
			msgSplit := strings.Split(netMsg, config.TIRET)
			id, err1 := strconv.Atoi(msgSplit[0])
			apt, err2 := strconv.Atoi(msgSplit[1])
			if err1 == nil && err2==nil{
				utils.PrintMessage(election.moi, nameElc, "receive demande election from site: "+msgSplit[0])
				election.receiveMessage(id, apt)
			}
		case <- election.endTime:
			utils.PrintMessage(election.moi, nameElc, "TimeOut")
			election.receiveTimeOut(election.aps)
		}

	}

}


//Pseudo-code(3-réception)
func (election *Election) election()  {
	election.demarre()
}

func (election *Election)demarre()  {
	election.enCours =true

	//initialisation du tableau à 0 avant le démarrage de l'election
	for i := 0; i <election.N; i++{
		election.aps = append(election.aps, 0)
	}
	election.aps[election.moi] = election.monApt
	tabMess := []string{fmt.Sprint(election.moi), fmt.Sprint(election.monApt)}
	messageSent := strings.Join(tabMess, config.TIRET)

	election.sendMessageToNetwork(messageSent)
	timeDelay := 2 * config.T*time.Millisecond
	election.endTime = time.After(timeDelay)
}

func (election *Election)receiveMessage(i int, apt int)  {
	if !election.enCours{
		election.demarre()
	}
	election.aps[i] = apt

}

//receiveTimeOut permet de traiter la fin d'un processus
func (election *Election)receiveTimeOut(tabApt []int)  {
aptMax := election.monApt

 for i, apt := range tabApt {
 	if aptMax == apt{
 		election.elu = utils.Maximum(election.elu, i)
	} else {
		aptMax = utils.Maximum(apt, aptMax)
		if aptMax == apt{
			election.elu = i
		}
	}
 }
 utils.PrintMessage(election.moi, nameElc, "tour effectuer sur les aptitudes obtenus")
 election.enCours = false
 election.sendResults(election.elu, aptMax)

}

func (election *Election)sendResults(idElu int, apt int)  {
	stringElu := []string{fmt.Sprint(idElu), fmt.Sprint(apt)}
	election.toClient <- strings.Join(stringElu, config.TIRET)
}


func (election *Election)sendMessageToNetwork(message string)  {
	utils.PrintMessage(election.moi, nameElc, "send to network : "+message)
	election.toNetwork <- message
}
