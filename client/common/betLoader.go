package common

import "github.com/7574-sistemas-distribuidos/docker-compose-init/client/model"

type BetLoader interface {
	Init() error
	GetBets() []model.Bet
	Destroy() error
}

type betLoader struct {}

func CreateBetLoader() *betLoader {
	return &betLoader{}
}

func (b *betLoader) Init() error {
	return nil
}

func (b *betLoader) GetBets() []model.Bet {
	return []model.Bet{}
}

func (b *betLoader) Destroy() error {
	return nil
}