package main

import (
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/logging"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/jessevdk/go-flags"
	easy "github.com/vittoriocapocasale/easy_tp/processor"
	"os"
	"syscall"
)

type Opts struct {
	Verbose []bool `short:"v" long:"verbose" description:"Increase verbosity"`
	Connect string `short:"C" long:"connect" description:"Validator component endpoint to connect to" default:"tcp://localhost:4004"`
	Queue   uint   `long:"max-queue-size" description:"Set the maximum queue size before rejecting process requests" default:"100"`
	Threads uint   `long:"worker-thread-count" description:"Set the number of worker threads to use for processing requests in parallel" default:"0"`
}

func main() {
//copiato e incollato dai transaction processor di esempio fino al prossimo commento
	var opts Opts

	logger := logging.Get()

	parser := flags.NewParser(&opts, flags.Default)
	remaining, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			logger.Errorf("Failed to parse args: %v", err)
			os.Exit(2)
		}
	}

	if len(remaining) > 0 {
		fmt.Printf("Error: Unrecognized arguments passed: %v\n", remaining)
		os.Exit(2)
	}

	endpoint := opts.Connect
//creo un handler del tipo che ho dichiarato  nel package processor
	easyHandler:=  easy.NewEasyHandler() 
//copiato e incollato dai transaction processor di esempio fino al prossimo commento
	processor := processor.NewTransactionProcessor(endpoint)
	processor.SetMaxQueueSize(opts.Queue)
	processor.SetThreadCount(1)
	if opts.Threads > 0 {
		processor.SetThreadCount(opts.Threads)
	}
//aggiungo il mio handler. Penso sia conveniente definire un solo handler per tutte le transazioni
	processor.AddHandler(easyHandler)
//copiato e incollato dai transaction processor di esempio fino al prossimo commento
	processor.ShutdownOnSignal(syscall.SIGINT, syscall.SIGTERM)
	err = processor.Start()//avvio il tansaction processor
	if err != nil {
		logger.Error("Processor stopped: ", err)
	}
}

