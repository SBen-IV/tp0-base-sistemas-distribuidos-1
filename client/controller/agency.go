package controller

import (
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

	buf = make([]byte, 2)

	_, err = a.client.Recv(buf, 2)

	if err != nil {
		log.Errorf("Error reading from server: %v", err)
		return
	}

	resp, _ := a.translator.OKtoString(buf)

	log.Debugf("Got response from server: %s", resp)

	// Get Bet
	// bet := a.betLoader.GetBets()

	// Count bets
	// Translate to bytes
	// Count bytes
	// Send bet as bytes
	// Wait for response

}

func (a *Agency) Close() error {
	a.betLoader.Destroy()
	a.client.Stop()
	
	return nil
}