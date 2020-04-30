package manager

import (
	"bytes"
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/vittoriocapocasale/easy_tp/cbor"
)


//prende un indirizzo e una struttura dati e salva la struttura dati nello stato a quell'indirizzo.
 
func Load(address string, collisionMap interface{}, context *processor.Context) ( error) {

	results, err := context.GetState([]string{address})
	if err != nil {
		 return err
	}
	data, exists := results[address]
	if exists && len(data) > 0 {
		err = cbor.DecodeCbor(data, &collisionMap)
		if err != nil {
			return err
		}
	}
	return  nil
}

//prende un indirizzo e una struttura dati e carica ciò che è a quell'indirizzo nella struttura dati. Essendo un key value store, in caso di key uguali tra due oggetti c'è il rischio sovrascrittura. Quindi per ogni key salvo una mappa di oggetti (per questo collisionMap, idea Sawtooth), così che in caso di conflitto possa salvarli entrambi. Suppongo che se mi fosse capitato realmente un conflitto, il sistema sarebbe diventato inconsistrnte a causa delle hashmap iterate casualmente.

func Store(address string, collisionMap interface{}, context *processor.Context) error {
	buffer := new(bytes.Buffer)
	err := cbor.EncodeCbor(collisionMap, buffer)
	if err != nil {
		return &processor.InternalError{
			Msg: fmt.Sprint("Failed to encode new map:", err),
		}
	}
	addresses, err := context.SetState(map[string][]byte{
		address: buffer.Bytes(),
	})
	if err != nil {
		return err
	}
	if len(addresses) == 0 {
		return &processor.InternalError{Msg: "No addresses in set response"}
	}
	return nil
}

//store su molti oggetti, chiamo una volta setState...abbastanza superflua come funzione
func MultipleStore(collisionMaps map[string]interface{}, context *processor.Context) error {

	e := make(map[string][]byte)
	for k, v := range collisionMaps {
		buffer := new(bytes.Buffer)
		err := cbor.EncodeCbor(v, buffer)
		if err != nil {
			return &processor.InternalError{
				Msg: fmt.Sprint("Failed to encode new map:", err),
			}
		}
		e[k] = buffer.Bytes()
	}
	addresses, err := context.SetState(e)
	if err != nil {
		return err
	}
	if len(addresses) != len(e) {
		return &processor.InternalError{Msg: "No addresses in set response"}
	}
	return nil
}








