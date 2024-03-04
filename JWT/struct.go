package main

// -----------------Get passport struct ----------------------
type tmpStringClaim struct {
	CID string `json: CID`
}

// -----------------Add MutableData struct ----------------------
type httpData struct {
	Key       string `json: Key`
	Eventtype string `json: Eventtype`
	Datetime  string `json: Datetime`
	Data      string `json: Data`
}

type ledgerData struct {
	Eventtype string `json: Eventtype`
	Data      string `json: Data`
	Datetime  string `json: Datetime`
}

type appendEntry struct {
	CID       string `json: CID`
	Eventtype string `json: Eventtype`
	// Name string `json: Name`
	Datetime string `json: Datetime`
}

// ----------------- retrive Event struct ----------------------
type chooseEvent struct {
	Key  string `json: Key`
	Type string `json: Type`
	CID  string `json: CID`
}
type getEvent struct {
	CID string `json:"CID"`
}

// ----------------- add Mutable Product struct ----------------------
type httpDataProduct struct {
	Key          string `json: Key`
	CID          string `json: CID`
	ProductType  string `json: ProductType`
	Datetime     string `json: Datetime`
	CIDToReplace string `json: CIDToReplace`
}

type ledgerDataProduct struct {
	ProductType string `json: ProductType`
	CID         string `json: CID`
	Datetime    string `json: Datetime`
}

type appendEntryProduct struct {
	CID         string `json: CID`
	ProductType string `json: ProductType`
	Datetime    string `json: Datetime`
}

// ----------------- Generate QR Code struct ----------------------

type QrCode struct {
	CID                    string `json: CID`
	Data                   string `json: Data`
	MaterialId             string `json: MaterialId`
	OrderId                string `json: OrderId`
	Dimensions             string `json: Dimensions`
	Plant                  string `json: Plant`
	Entrydate              string `json: Entrydate`
	Remanufacturing_events string `json:"remanufacturing_events"`
	Shipping               string `json:"shipping"`
	Makes                  string `json:"makes"`
	Made_from              string `json:"made_from"`
}

type QrCodeImage struct {
	Filename string `json: Filename`
	Content  string `json: Content`
}

// ----------------- retriveMutableLog api check struct -----------------------

type MutableLog struct {
	Key string `json: Key`
}

// ----------------- KeyPAIR request api struct -----------------------
type keyData struct {
	Privatekey string `json:"privatekey"`
	Publickey  string `json:"publickey"`
}

type dataToCA struct {
	Cid                    string  `json:"cid"`
	Remanufacturing_events keyData `json:"remanufacturing_events"`
	Shipping               keyData `json:"shipping"`
	Makes                  keyData `json:"makes"`
	Made_from              keyData `json:"made_from"`
}

type dataFromCa struct {
	PrivateKey string `json:"privateKey"`
	Success    string `json:"success"`
	Message    string `json:"message"`
}
