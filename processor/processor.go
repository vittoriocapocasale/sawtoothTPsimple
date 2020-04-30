package processor

import (
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/processor_pb2"
	"github.com/vittoriocapocasale/easy_tp/cbor"
	"github.com/vittoriocapocasale/easy_tp/data_structures"
	"github.com/vittoriocapocasale/easy_tp/manager"
)

//nome a piacere, sempre una struct vuota
type EasyHandler struct{}

//chiamata in main. Serve solo a poterci chiamare funzioni 
func NewEasyHandler() *EasyHandler {
	return &EasyHandler{}
}

//modo di go di definire i "metodi" queste funzioni sono standard e ci intteressano poco 
func (self *EasyHandler) FamilyName() string {
	return "easy"
}

func (self *EasyHandler) FamilyVersions() []string {
	return []string{"1.0"}
}

//se usiamo un solo handler, probabilmente ci appropriamo di tutto il namespace
func (self *EasyHandler) Namespaces() []string {
	return []string{"000001","000002"};
}

//qui parte il codice "vero"
func (self *EasyHandler) Apply(request *processor_pb2.TpProcessRequest, context *processor.Context) error {

	payloadData := request.GetPayload()
	if payloadData == nil {
		return &processor.InvalidTransactionError{Msg: "Must contain payload"}
	}
	var payload data_structures.Payload
        //decodifico il payload in una struttura di tipo payload
	err := cbor.DecodeCbor(payloadData, &payload)
	if err != nil {
		return &processor.InternalError{Msg: fmt.Sprint("Failed to decode payload: ", err)}
	}
        //mi copio i campi che mi servono, giusto per sprecare memoria
	payloadAction := payload.Action
	payloadId := payload.Id
	identity := request.Header.GetSignerPublicKey()
	payloadValue := payload.Value
	payloadTime := payload.Time
	fmt.Println(payloadAction, payloadId, payloadTime, payloadValue)
        //logica delle transazioni, non sono sicuro sia giusta, potrebbe essere una versione vecchia e comunque non ho affrontato praticamente nessuno dei problemi di cui abbiamo discusso
	switch payloadAction {
	case "create":
                //in pratica è una new
		entity:=&data_structures.IntStruct{payloadId, 0, payloadTime, identity}
                //preparo la struttura in cui fare la load
		collisionMap:=map[string]*data_structures.IntStruct{}
                //faccio la load. è importante che tutte le strutture in un indirizzo siano uniformi, cioè tutte *data_structures.IntStruct{}
		err:= manager.Load(entity.ComputeAddress(), collisionMap ,context)
		if err!= nil {
			return err;
		}
                //aggiungo entity a collisionMap se non esiste già
		err= data_structures.CreateIntStruct(entity, collisionMap)
		if err!= nil {
			return err;
		}
                //salvo le modifiche
		err = manager.Store(entity.ComputeAddress(), collisionMap, context)
		if err!= nil {
			return err;
		}
		return nil
	case "delete":
		entity:=&data_structures.IntStruct{payloadId, 0, payloadTime, identity}
		collisionMap:=map[string]*data_structures.IntStruct{}
		err:= manager.Load(entity.ComputeAddress(), collisionMap ,context)
		if err!= nil {
			return err;
		}
		entity, err= data_structures.DeleteIntStruct(payloadId, collisionMap, payloadTime, identity)
		if err!= nil {
			return err;
		}
		err = manager.Store(entity.ComputeAddress(), collisionMap, context)
		if err!= nil {
			return err;
		}
		return nil
	case "update":
		entity:=&data_structures.IntStruct{payloadId, 0, payloadTime, identity}
		collisionMap:=map[string]*data_structures.IntStruct{}
		err:= manager.Load(entity.ComputeAddress(), collisionMap ,context)
		if err!= nil {
			return err;
		}
		fmt.Println(collisionMap)
		entity, err= data_structures.UpdateIntStruct(payloadId, collisionMap, payloadValue, payloadTime, identity)
		if err!= nil {
			return err;
		}
		err = manager.Store(entity.ComputeAddress(), collisionMap, context)
		if err!= nil {
			return err;
		}
		return nil
	case "exchange":
		entity:=&data_structures.IntStruct{payloadId, 0, payloadTime, identity}
		entitiesMap:=map[string]*data_structures.IntStruct{}
		err:= manager.Load(entity.ComputeAddress(), entitiesMap ,context)
		if err!= nil {
			return err;
		}
		entity, exists := entitiesMap[entity.Id]
		if !exists {
			return &processor.InvalidTransactionError{Msg: "Entity does not exists"}
		}
		if entity.Time>payloadTime {
			return &processor.InvalidTransactionError{Msg: "Too old to update"}
		}
		if entity.Admin==identity {
			return &processor.InvalidTransactionError{Msg: "Entity already owned"}
		}
		proposal:=&data_structures.Proposal{identity+entity.Id, identity, payloadTime, entity.Id, []string{}, []string{identity, entity.Admin}}
		proposalsMap:= map[string]*data_structures.Proposal{}
		err= manager.Load(proposal.ComputeAddress(), entitiesMap ,context)
		if err!= nil {
			return err;
		}
		err=data_structures.CreateProposal(proposal, proposalsMap)
		if err!= nil {
			return err;
		}
		err = manager.Store(proposal.ComputeAddress(), proposalsMap, context)
		if err!= nil {
			return err;
		}

		return nil
	case "vote":
		proposal:=&data_structures.Proposal{payloadId, identity, payloadTime, payloadId, []string{}, []string{}}
		proposalsMap:=map[string]*data_structures.Proposal{}
		err:= manager.Load(proposal.ComputeAddress(), proposalsMap ,context)
		if err!= nil {
			return err;
		}
		proposal, exists := proposalsMap[proposal.Id]
		if !exists {
			return &processor.InvalidTransactionError{Msg: "Proposal does not exists"}
		}
		if proposal.Time>payloadTime {
			return &processor.InvalidTransactionError{Msg: "Transaction in too old"}
		}
		entity := & data_structures.IntStruct{proposal.EntityId, 0, payloadTime, identity}
		entitiesMap:=map[string]*data_structures.IntStruct{}
		err= manager.Load(entity.ComputeAddress(), entitiesMap ,context)
		if err!= nil {
			return err;
		}
		entity, exists = entitiesMap[entity.Id]
		if !exists {
			return &processor.InvalidTransactionError{Msg: "Entity does not exists"}
		}
		if entity.Time>proposal.Time {
			delete(proposalsMap, proposal.Id)
			err=manager.Store(proposal.ComputeAddress(), proposalsMap, context)
			if err != nil {
				return err
			}
			return nil
		}
		for k:=range proposal.Voters {
			if proposal.Voters[k]==identity {
				return &processor.InvalidTransactionError{Msg: "Double vote attempt"}
			}
		}
		isRequired:=false;
		for k:=range proposal.RequiredVoters {
			if proposal.RequiredVoters[k]==identity {
				isRequired=true
			}
		}
		if !isRequired {
			return &processor.InvalidTransactionError{Msg: "Vote not allowed"}
		}
		proposal.Voters=append(proposal.Voters, identity)
		if len(proposal.Voters)==len(proposal.RequiredVoters) {
			entity.Admin=proposal.Value
			entity.Time=proposal.Time
			delete(proposalsMap, proposal.Id)
			err=manager.Store(entity.ComputeAddress(), entitiesMap, context)
			if err != nil {
				return err
			}
		}
		err=manager.Store(proposal.ComputeAddress(), proposalsMap,  context)
		if err != nil {
			return err
		}
		return nil
	default:
		return &processor.InvalidTransactionError{Msg:"Unknown Action"}
	}
}
