package data_structures

type Entity interface {
//ogni entity ha un id
	GetId() string
//serve a darmi l'indirizzo del key value store in cui salvare questo oggetto
	ComputeAddress() string
}

//payload di una transazione
type Payload struct {
	Action string;
	Id     string;
    Time   uint;
    Value  int;
}


