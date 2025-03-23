package controller

import (
	"fmt"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common"
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
)

// This class acts as a controller for communication and model
type Agency struct {
	betLoader common.BetLoader
	client common.Client
	id string
	stopped chan bool
	state ProtocolState
}

func NewAgency(betLoader common.BetLoader, client common.Client, config common.ClientConfig) *Agency {
	return &Agency{
		betLoader: betLoader,
		client: client,
		id: config.ID,
		stopped: make(chan bool, 1),
	}
}

func (a *Agency) Run() {
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
					isRunning = false
				} else {
					a.state = SendBets
				}
			case SendBets:
				if err := a.sendBets(); err != nil {
					log.Errorf("Could not send bet to server: %v", err)
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

	return a.waitOK()
}

func (a *Agency) sendBets() error {
	// Get Bet
	bet := a.betLoader.GetBet()
	betProtocol := protocol.NewBet(bet)

	buf, bytesAmount := betProtocol.Encode()

	betsInfo := protocol.NewBetInfo(1, int32(bytesAmount))

	if err := a.sendBetInfo(betsInfo); err != nil {
		return err
	}

	log.Debugf("Sending %v to server with len %v", betProtocol, bytesAmount)

	// Send bet as bytes and wait for response
	if err := a.sendBet(buf, bytesAmount); err != nil {
		return err
	}

	log.Infof("action: apuesta_enviada | result: success | dni: %s | numero: %d", bet.Document, bet.Number)

	return nil
}

func (a *Agency) sendBetInfo(betsInfo *protocol.BetInfo) error {
	buf, bytes_amount, err := betsInfo.Encode()

	if err != nil {
		return err
	}

	log.Debugf("Sending bet info: %v", buf)

	a.client.Send(buf, bytes_amount)

	return a.waitOK()
}

func (a *Agency) sendBet(buf []byte, bytesAmount int) error {
	_, err := a.client.Send(buf, bytesAmount)

	if err != nil {
		return err
	}

	return a.waitOK()
}

func (a *Agency) waitOK() error {
	buf, bytesAmount := protocol.NewMessageBuf()

	bytesRecv, err := a.client.Recv(buf, bytesAmount)

	if err != nil {
		log.Errorf("Error reading from server: %v", err)
		return err
	}

	message := protocol.NewServerMessageBuild(buf, bytesRecv)

	if message != protocol.Ok {
		log.Error("Error building message from server: %v", message)
		return fmt.Errorf("error building message from server: %v", buf)
	}

	log.Debugf("Got response from server: %s", message)

	return nil
}

func (a *Agency) Close() error {
	a.client.Stop()

	close(a.stopped)
	
	return nil
}