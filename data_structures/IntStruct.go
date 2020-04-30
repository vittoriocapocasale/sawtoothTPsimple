package data_structures

import (
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/vittoriocapocasale/easy_tp/crypto"
)
const IS_NAMESPACE="000001";
type IntStruct struct {
	Id      string;
	Value   int;
	Time    uint;
	Admin   string;
}

func (self *IntStruct) GetId() string{
	return self.Id
}

func (self *IntStruct) ComputeAddress() string {
	hashedName := crypto.Hexdigest(self.Id)
	return IS_NAMESPACE + hashedName[len(hashedName)-(70-len(IS_NAMESPACE)):]
}

func CreateIntStruct(entity *IntStruct, collisionMap map[string]*IntStruct) (error) {

	_, exists := collisionMap[entity.Id]
	if exists {
		return &processor.InvalidTransactionError{Msg: "Entity already existent"}
	}
	collisionMap[entity.Id] = entity
	return nil
}


func DeleteIntStruct(id string, collisionMap map[string]*IntStruct, time uint, identity string) (*IntStruct, error) {
	entity, exists := collisionMap[id]
	if !exists {
		return nil, &processor.InvalidTransactionError{Msg: "Entity not found"}
	}
	if entity.Time>time {
		return nil, &processor.InvalidTransactionError{Msg: "Too old to update"}
	}
	if entity.Admin!=identity {
		return nil, &processor.InvalidTransactionError{Msg: "Permission denied"}
	}
	delete(collisionMap, id)
	return entity, nil
}

func UpdateIntStruct(id string, collisionMap map[string]*IntStruct, newValue int, time uint, identity string) (*IntStruct, error) {

	entity, exists := collisionMap[id]
	if !exists {
		return nil, &processor.InvalidTransactionError{Msg: "Entity not found"}
	}
	if entity.Admin!=identity {
		return nil, &processor.InvalidTransactionError{Msg: "Permission denied"}
	}
	if entity.Time>time {
		return nil, &processor.InvalidTransactionError{Msg: "Too old to update"}
	}
	entity.Value=newValue
	entity.Time=time
	return entity, nil
}
