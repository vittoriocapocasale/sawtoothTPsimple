package data_structures

import (
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/vittoriocapocasale/easy_tp/crypto"
)

const P_NAMESPACE = "000002"

type Proposal struct {
	Id string;
	Value string;
	Time uint;
	EntityId string;
	Voters []string;
	RequiredVoters []string;
}

func (self *Proposal) GetId() string{
	return self.Id
}

//serve a darmi l'indirizzo del key value store in cui salvare questo oggetto
func (self *Proposal) ComputeAddress() string {
	hashedName := crypto.Hexdigest(self.Id)
	return P_NAMESPACE + hashedName[len(hashedName)-(70-len(P_NAMESPACE)):]
}

//collision Map è quella che viene restituita dalla load
func CreateProposal(entity *Proposal, collisionMap map[string]*Proposal) (error) {

	_, exists := collisionMap[entity.Id]
	if exists {
		return &processor.InvalidTransactionError{Msg: "Entity already existent"}
	}
	collisionMap[entity.Id] = entity
	return nil
}
