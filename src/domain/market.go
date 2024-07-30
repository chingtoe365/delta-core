package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SignalType string

const (
	// add more signal types here
	Change SignalType = "change"
)

const (
	CollectionMarketSignal = "markert_signals"
)

type MarketSignalMeta struct {
	Key    string
	Type   SignalType
	Config map[string]interface{}
}

type IMarketSignaler interface {
	Roll(SingalId primitive.ObjectID)
	Notify(string)
	Remove()
}

type MarketSignalDto struct {
	Id         primitive.ObjectID `bson:"_id" json:"-"`
	UserId     primitive.ObjectID `bson:"userId" json:"-"`
	SignalMeta MarketSignalMeta
}
