package subscription

import (
	"net/http"
	"log"
)

type Manager struct{
	SubMap map[*ItemSearch]map[*Client]bool
	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	Quit chan bool
}

func NewManager() *Manager {
	return &Manager{
		SubMap:   make(map[*ItemSearch]map[*Client]bool),
		register:   make(chan *Client),
		Quit:	make(chan bool),
		unregister: make(chan *Client),

	}
}

func (m *Manager) Run() {
	for {
		select {
		case client := <- m.register:
			log.Println("Trying to register client")
			if _, ok := m.SubMap[client.ItemSearch]; !ok {
				log.Println("New search, creating client map")
				m.SubMap[client.ItemSearch] = make(map[*Client]bool)
			}
			m.SubMap[client.ItemSearch][client] = true
			log.Println(m.SubMap[client.ItemSearch])
		case client := <- m.unregister:
			log.Println("Deleting client")
			delete(m.SubMap[client.ItemSearch], client)
			log.Println(len(m.SubMap[client.ItemSearch]))
			if (!(len(m.SubMap[client.ItemSearch])> 0)){
				delete(m.SubMap, client.ItemSearch)
			}
			log.Println(len(m.SubMap))
		}
	}
}


func (manager *Manager) ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(r.FormValue("league"))
	search := &ItemSearch{Type : r.FormValue("type"), Name : r.FormValue("name"), League : r.FormValue("league")}
	log.Println(search)
	//if search valid
	if (search.League != "" && (search.Type != "" || search.Name != "")){
		client := &Client{manager: manager, conn: conn, Send: make(chan []byte, 256), ItemSearch : search}
		manager.register <- client
		go client.writePump()
	}else{
		log.Println("Client sent invalid search ", search)
	}



	//client.readPump()

}