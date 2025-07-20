package store

import "time"

type InputRef struct {
	TxID string
	VOut uint32
}

type UTXO struct {
	TxID     string
	VOut     int
	Address  string
	Amount   float64
	Created  time.Time
	Modified time.Time
}

type Tracker struct {
	ID        int
	Address   string
	CreatedAt time.Time
}
