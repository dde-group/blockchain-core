package network

type HubBase struct {
	Name string
	// Registered clients.
	Clients map[*WebsocketListener]struct{}

	// Inbound messages from the clients.
	Broadcast chan []byte

	// Register requests from the clients.
	Register chan *WebsocketListener

	// Unregister requests from clients.
	Unregister chan *WebsocketListener

	UnregisterClose chan *WebsocketListener
}
