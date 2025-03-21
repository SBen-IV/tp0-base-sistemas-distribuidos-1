package controller

import (
	"encoding/binary"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

type Agency struct {
	betLoader common.BetLoader
	client common.Client
	id string
	translator *common.ProtocolTranslator
}

func NewAgency(config common.ClientConfig) *Agency {
	return &Agency{
		betLoader: common.CreateBetLoader(),
		client: common.CreateClient(config),
		id: config.ID,
		translator: common.NewProtocolTranslator(),
	}
}

func (a *Agency) Run() {
	a.betLoader.Init()
	err := a.client.Connect()

	if err != nil {
		log.Errorf("Could not connect to server: %v", err)
		return
	}

	a.identifyToLotery()

	// Get Bet
	bet := a.betLoader.GetBet()
	bet_in_bytes := a.translator.BetToBytes(bet)

	a.sendBetInfo(bet_in_bytes)

	log.Debugf("Sending %v to server with len %v", bet_in_bytes, len(bet_in_bytes))

	// Send bet as bytes
	// Wait for response
	a.sendBet(bet_in_bytes)

	log.Infof("action: apuesta_enviada | result: success | dni: %s | numero: %s", bet.Document, bet.Number)
}

func (a *Agency) identifyToLotery() {
	// Send CLI_ID and wait for response
	buf, err := a.translator.IDtoBytes(a.id)

	if err != nil {
		log.Errorf("Error translating ID: %v", err)
		return
	}

	bytes_sent, err := a.client.Send(buf, len(buf))

	log.Debugf("Bytes sent: %d bytes", bytes_sent)
	log.Error("Error %v", err)

	// Wait for response

	a.waitOK()
}

func (a *Agency) sendBetInfo(bet_in_bytes []byte) {

	// Count bets
	// Translate to bytes
	// Count bytes

	buf := make([]byte, 8)

	binary.BigEndian.PutUint32(buf[0:4], 1) // Send only 1 bet
	binary.BigEndian.PutUint32(buf[4:8], uint32(len(bet_in_bytes))) // Send bytes amount

	log.Debugf("Sending bet info: %v", buf)

	a.client.Send(buf, 8)

	a.waitOK()
}

func (a *Agency) sendBet(bet []byte) {
	a.client.Send(bet, len(bet))

	a.waitOK()
}

func (a *Agency) waitOK() {
	buf := make([]byte, 2)

	_, err := a.client.Recv(buf, 2)

	if err != nil {
		log.Errorf("Error reading from server: %v", err)
		return
	}

	resp, _ := a.translator.OKtoString(buf)

	log.Debugf("Got response from server: %s", resp)
}

func (a *Agency) Close() error {
	a.betLoader.Destroy()
	a.client.Stop()
	
	return nil
}