# PRR_Labo3

## Claude-André Alves, Olivier Djeulezeck

### Lancement
Pour le lancement les paramètres de connection se trouvent dans le fichier configNetworkBully.json chaque client/site est représenté par un payload json une fois mis au bonnes valeur il suffit de lancer le main avec les commandes suivantes :
go get github.com/tkanos/gonfig
go run main.go

Nous avons utiliser le package gonfig pour pouvoir faire la lecture du fichier json de config.

### /!\ Problème connu /!\
Lorsque l'on lance les clients une fois une élection faites le programme se bloque, nous n'avons pas réussi à debugger celà nous suspectons un problème avec les les timers.


