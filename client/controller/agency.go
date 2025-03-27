package controller

import (
	"fmt"
	"time"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/model"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/protocol"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// State of the communication with server
type ProtocolState int

const (
	ConnectToNationalLottery ProtocolState = iota
	IdentifyToNationalLottery
	SendBets
	NoMoreBets
	EndConnection
	GetWinners
)

const EMPTY_BETS = 0

// A bet contains
// FIRST_NAME;LAST_NAME;DOCUMENT;BIRTHDAY;NUMBER,
// Suppose the max character amount for each element is:
// FIRST_NAME = 15
// LAST_NAME = 15
// DOCUMENT = 8
// BIRTHDAY = 10
// NUMBER = 4
// SEPARATORS (; y ,) PER BET = 5
// TOTAL = 57 Bytes per bet
// MAX BYTES PER BATCH = 8 kB = 8000 bytes
// 8000 / 57 ~= 140.35 => 140
const MAX_BATCH_AMOUNT = 140

// This class acts as a controller for communication and model
type Agency struct {
	betLoader common.BetLoader
	client common.Client
	id string
	stopped chan bool
	state ProtocolState
	batchMaxAmount int
	loopPeriod time.Duration
}

func NewAgency(betLoader common.BetLoader, client common.Client, config common.ClientConfig) *Agency {
	return &Agency{
		betLoader: betLoader,
		client: client,
		id: config.ID,
		stopped: make(chan bool, 1),
		batchMaxAmount: common.Min(config.BatchMaxAmount, MAX_BATCH_AMOUNT),
		loopPeriod: config.LoopPeriod,
	}
}

func (a *Agency) Run() {
	if err := a.betLoader.Init(); err != nil {
		log.Errorf("Error initializing bet loader: %v", err)
		return 
	}

	var isRunning bool = true

	for isRunning {
		select {
		case <-a.stopped:
			log.Info("Stop received")
 			isRunning = false
		default:
			switch a.state {
			case ConnectToNationalLottery:
				if err := a.client.Connect(); err != nil {
					log.Errorf("Could not connect to server: %v", err)
					isRunning = false
				} else {
					a.state = IdentifyToNationalLottery
				}
			case IdentifyToNationalLottery:
				if err := a.identifyToNationalLottery(); err != nil {
					log.Errorf("Could not identify to server: %v", err)
					isRunning = false
				} else {
					a.state = SendBets
				}
			case SendBets:
				if err := a.manageBets(); err != nil {
					switch err.(type) {
					case *common.NoMoreBets:
						log.Debugf("%v", err)
						a.state = NoMoreBets
					default:
						log.Errorf("Could not manage bets: %v", err)
						isRunning = false
					}
				}
			case NoMoreBets:
				if err := a.sendOperation(model.NoMoreBets); err != nil {
					log.Errorf("Could not send no more bets: %v", err)
					isRunning = false
				} else {
					a.state = GetWinners
				}
			case GetWinners:
				if err := a.manageDraw(); err != nil {
					switch err.(type) {
					case *common.WinnersNotAvailableYet:
						// Wait some time to ask again
						time.Sleep(a.loopPeriod)
					default:
						log.Errorf("Could not manage winners: %v", err)
						isRunning = false
					}
				} else {
					a.state = EndConnection
				}
			case EndConnection:
				if err := a.sendOperation(model.FinOp); err != nil {
					log.Errorf("Could not end connection: %v", err)
				}

				isRunning = false
			}
		}
	}
}

func (a *Agency) identifyToNationalLottery() error {
	// Send CLI_ID and wait for response
	buf, bytesAmount, err := protocol.NewAgencyID(a.id).Encode()

	if err != nil {
		log.Errorf("Error creating AgencyID: %v", err)
		return err
	}

	bytesSent, err := a.client.Send(buf, bytesAmount)

	log.Debugf("Bytes sent: %d bytes", bytesSent)
	if err != nil {
		log.Errorf("Error sending bytes: %v", err)
		return err
	}

	// Wait for response

	return a.waitServerResponse()
}

func (a *Agency) manageBets() error {
	bets, err := a.betLoader.GetBets(a.batchMaxAmount)
	
	if err != nil {
		return err
	}
	
	betsAmount := len(bets)

	if betsAmount == EMPTY_BETS {
		return &common.NoMoreBets{}
	}

	if err := a.sendOperation(model.BetOp); err != nil {
		return err
	}

	betsProtocol := protocol.NewBets(bets)
	buf, bytesAmount := betsProtocol.Encode()
	betsInfo := protocol.NewBetInfo(int32(betsAmount), int32(bytesAmount))

	if err := a.sendBetsInfo(betsInfo); err != nil {
		return err
	}
	
	log.Infof("Sending %v bets to server with len %v", betsAmount, bytesAmount)
	
	if err := a.sendBets(buf, bytesAmount); err != nil {
		return err
	}
	
	log.Infof("action: apuesta_enviada | result: success | cantidad: %v", betsAmount)

	return nil
}

func (a *Agency) manageDraw() error {
	log.Debug("Asking for winners")

	if err := a.sendOperation(model.DrawOp); err != nil {
		return err
	}

	winners, err := a.recvWinners()

	if err != nil {
		return err
	}

	log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %v", len(winners))

	return nil
}

func (a *Agency) sendOperation(op model.Operation) error {
	buf, bytesAmount := protocol.NewOperation(op).Encode()

	bytesSent, err := a.client.Send(buf, bytesAmount)

	log.Debugf("Sent operation %v with len %v", op, bytesSent)

	if err != nil {
		log.Errorf("Error sending bytes: %v", err)
		return err
	}

	// Wait for response

	return a.waitServerResponse()
}

func (a *Agency) sendBetsInfo(betsInfo *protocol.BetInfo) error {
	buf, bytes_amount, err := betsInfo.Encode()

	if err != nil {
		return err
	}

	log.Debugf("Sending bet info: %v", buf)

	a.client.Send(buf, bytes_amount)

	return a.waitServerResponse()
}

func (a *Agency) sendBets(buf []byte, bytesAmount int) error {
	_, err := a.client.Send(buf, bytesAmount)

	if err != nil {
		return err
	}

	return a.waitServerResponse()
}

func (a *Agency) recvWinners() ([]model.Winner, error) {
	buf, bytesAmount := protocol.NewWinnersBytesAmountBuf()

	_, err := a.client.Recv(buf, bytesAmount)

	if err != nil {
		return nil, err
	}

	winnersBytesAmount := protocol.NewWinnersBytesAmountBuild(buf)

	if err := a.sendOk(); err != nil {
		return nil, err
	}

	buf, bytesAmount = protocol.NewWinnersBuf(winnersBytesAmount)

	_, err = a.client.Recv(buf, bytesAmount)

	if err != nil {
		return nil, err
	}

	winners, err := protocol.NewWinnersBuild(buf, bytesAmount)

	if err != nil {
		return nil, err
	}

	if err := a.sendOk(); err != nil {
		return nil, err
	}

	return winners, err
}

func (a *Agency) waitServerResponse() error {
	buf, bytesAmount := protocol.NewMessageBuf()

	bytesRecv, err := a.client.Recv(buf, bytesAmount)

	if err != nil {
		log.Errorf("Error reading from server: %v", err)
		return err
	}

	message := protocol.NewServerMessageBuild(buf, bytesRecv)

	if message == protocol.WinnersNotAvailable {
		return &common.WinnersNotAvailableYet{}
	}
	
	if message != protocol.Ok {
		return fmt.Errorf("error building message from server: %v", message)
	}

	log.Debugf("Got response from server: %s", message)

	return nil
}

func (a *Agency) sendOk() error {
	okMessage := protocol.NewOkMessageBuf()
	buf, bytesAmount := okMessage.Encode()

	_, err := a.client.Send(buf, bytesAmount)

	return err
}

func (a *Agency) Close() error {
	a.client.Stop()

	close(a.stopped)
	
	return nil
}