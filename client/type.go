package client

import "encoding/json"

// RawParm defines raw parameter command
type RawParm struct {
	Token string `json:"token"`
	Data  string `json:"data"`
}

// CreateTxIn create tx input
type CreateTxIn struct {
	Execer     string          `json:"execer"`
	ActionName string          `json:"actionName"`
	Payload    json.RawMessage `json:"payload"`
}

// QueryParm Query parameter
type QueryParm struct {
	Hash string `json:"hash"`
}

// Signature parameter
type Signature struct {
	Ty        int32  `json:"ty"`
	Pubkey    string `json:"pubkey"`
	Signature string `json:"signature"`
}

// Transaction parameter
type Transaction struct {
	Execer     string          `json:"execer"`
	Payload    json.RawMessage `json:"payload"`
	RawPayload string          `json:"rawPayload"`
	Signature  *Signature      `json:"signature"`
	Fee        int64           `json:"fee"`
	FeeFmt     string          `json:"feefmt"`
	Expire     int64           `json:"expire"`
	Nonce      int64           `json:"nonce"`
	From       string          `json:"from,omitempty"`
	To         string          `json:"to"`
	Amount     int64           `json:"amount,omitempty"`
	AmountFmt  string          `json:"amountfmt,omitempty"`
	GroupCount int32           `json:"groupCount,omitempty"`
	Header     string          `json:"header,omitempty"`
	Next       string          `json:"next,omitempty"`
	Hash       string          `json:"hash,omitempty"`
}

// ReceiptDataResult receipt data result
type ReceiptDataResult struct {
	Ty     int32               `json:"ty"`
	TyName string              `json:"tyName"`
	Logs   []*ReceiptLogResult `json:"logs"`
}

// ReceiptLogResult receipt log result
type ReceiptLogResult struct {
	Ty     int32           `json:"ty"`
	TyName string          `json:"tyName"`
	Log    json.RawMessage `json:"log"`
	RawLog string          `json:"rawLog"`
}

// Asset asset
type Asset struct {
	Exec   string `json:"exec"`
	Symbol string `json:"symbol"`
	Amount int64  `json:"amount"`
}

// TxProof :
type TxProof struct {
	Proofs   []string `json:"proofs"`
	Index    uint32   `json:"index"`
	RootHash string   `json:"rootHash"`
}

// TransactionDetail transaction detail
type TransactionDetail struct {
	Tx         *Transaction       `json:"tx"`
	Receipt    *ReceiptDataResult `json:"receipt"`
	Proofs     []string           `json:"proofs"`
	Height     int64              `json:"height"`
	Index      int64              `json:"index"`
	Blocktime  int64              `json:"blockTime"`
	Amount     int64              `json:"amount"`
	Fromaddr   string             `json:"fromAddr"`
	ActionName string             `json:"actionName"`
	Assets     []*Asset           `json:"assets"`
	TxProofs   []*TxProof         `json:"txProofs"`
	FullHash   string             `json:"fullHash"`
}

// Query4Jrpc query jrpc
type Query4Jrpc struct {
	Execer   string          `json:"execer"`
	FuncName string          `json:"funcName"`
	Payload  json.RawMessage `json:"payload"`
}
