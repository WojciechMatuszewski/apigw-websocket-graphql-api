package connection

type Connection struct {
	GqlID  string `json:"gqlID"`
	ID     string `json:"id"`
	Params string `json:"params"`
	Event  string `json:"event"`
}

func (c Connection) AsItem() connectionItem {
	return connectionItem{
		Key:          NewKey(c.Event, c.ID),
		ConnectionID: c.ID,
		GqlID:        c.GqlID,
		Params:       c.Params,
	}
}

type connectionItem struct {
	Key
	ConnectionID string `json:"connectionID"`
	GqlID        string `json:"gqlID"`
	Params       string `json:"query"`
}

func (ci connectionItem) AsConnection() Connection {
	return Connection{
		GqlID:  ci.GqlID,
		ID:     ci.Key.AsConnectionID(),
		Params: ci.Params,
		Event:  ci.Key.AsEventName(),
	}
}
