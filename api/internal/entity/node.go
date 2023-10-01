package entity

type Node struct {
	Id                 string `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	SenderEmail        string `json:"senderEmail"`
	SenderMacAddress   string `json:"senderMacAddress"`
	ReceiverEmail      string `json:"receiverEmail"`
	ReceiverMacAddress string `json:"ReceiverMacAddress"`
}

// Sender(Up server) - Backend(notification) - Receiver
// Receiver send request to get file - Sender
