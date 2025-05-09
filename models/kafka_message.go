package models

type MessageMetadata struct {
	IssuedDate         string   `json:"issuedDate"`
	ConfidenceLevel    string   `json:"confidenceLevel"`
	Process            string   `json:"process"`
	Attachments        []string `json:"attachments"`
	DocID              string   `json:"docId"`
	Subject            string   `json:"subject"`
	InforSign          string   `json:"inforSign"`
	TypeName           string   `json:"typeName"`
	Format             string   `json:"format"`
	Description        string   `json:"description"`
	Language           []string `json:"language"`
	Autograph          string   `json:"autograph"`
	RiskRecoveryStatus string   `json:"riskRecoveryStatus"`
	CodeNumber         string   `json:"codeNumber"`
	NumberOfPage       string   `json:"numberOfPage"`
	Mode               string   `json:"mode"`
	OrganName          string   `json:"organName"`
	CodeNotation       string   `json:"codeNotation"`
	ArcDocCode         string   `json:"arcDocCode"`
	SchemaID           string   `json:"schemaId"`
	FileExtension      string   `json:"fileExtension"`
	Keyword            string   `json:"keyword"`
	Maintenance        string   `json:"maintenance"`
	RiskRecovery       string   `json:"riskRecovery"`
}

type ReceivedMessage struct {
	Metadata  MessageMetadata `json:"metadata"`
	PartyCode string          `json:"partyCode"`
	PID       string          `json:"pid"`
	ID        string          `json:"id"`
	Source    string          `json:"source"`
	Type      string          `json:"type"`
	FondCode  string          `json:"fondCode"`
	Content   string          `json:"content"`
}
