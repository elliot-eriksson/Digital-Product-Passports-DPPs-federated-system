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

// ----------------- createPassportHandler api check struct ----------------------
type APICheck struct {
	APIKey string `json: api_key`
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
	remanufacturing_events string `json: remanufacturing_events`
	shipping               string `json: shipping`
	makes                  string `json: makes`
	made_from              string `json: made_from`
}

type QrCodeImage struct {
	Filename string `json: Filename`
	Content  string `json: Content`
}

// ----------------- retriveMutableLog api check struct -----------------------

type MutableLog struct {
	Key string `json: Key`
}

// ----------------- KeyPAIR request api check struct -----------------------
type keyData struct {
	privatekey string `json: privatekey`
	publickey  string `json: publickey`
}

type dataToCA struct {
	cid                    string  `json: cid`
	remanufacturing_events keyData `json: remanufacturing_events`
	shipping               keyData `json: shipping`
	makes                  keyData `json: makes`
	made_from              keyData `json: made_from`
}
