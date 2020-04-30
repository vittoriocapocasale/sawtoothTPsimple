package cbor

import (
	"bytes"
	"github.com/brianolson/cbor_go"
	"github.com/ugorji/go/codec"
	"sync"
)

var cborHandle *codec.CborHandle
var once sync.Once

//singleton, serve solo per avere una handle con cui chiamare encode e decode
func GetCBORHandleInstance() *codec.CborHandle {
	once.Do(func() {
		cborHandle = new(codec.CborHandle)
	})
	return cborHandle
}


//funzioni per codificare e decodificare, nulla di particolarmente rilevante
func EncodeCbor (input interface{}, output *bytes.Buffer) error {
	e:=cbor.NewEncoder(output)
	return e.Encode(input)

}

func DecodeCbor (input []byte , output interface{}) error {
	h:=GetCBORHandleInstance()
	var dec *codec.Decoder = codec.NewDecoderBytes(input , h)
	var err error = dec.Decode(output)
	return err
}
