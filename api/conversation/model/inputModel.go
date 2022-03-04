package model

type InputStartCall struct {
	OrderId         string   `json:"order_id"`
}

type InputEndCall struct {
	CallId          string   `json:"call_id"`
}


type InputJoinCall struct {
	CallId          string   `json:"call_id"`
}

type InputSubmitInfoConversation struct {
	CallId          string   `json:"call_id" bson:"call_id"`
	StartCall       string   `json:"start_call" bson:"start_call"`
	EndCall         string   `json:"end_call" bson:"end_call"`
	Content         string   `json:"content" bson:"content"`
	ConsultingField []string `json:"consulting_field" bson:"consulting_field"`
}

type InputSubmitCallEvaluation struct {
	CallId       string `json:"call_id"`
	NoteCompany  string `json:"note_company" `
	NoteListener string `json:"note_listener" `
	Star         int    `json:"star"`
	Tip          int64  `json:"tip"`
}

type InputSwitchListener struct {
	OrderId         string   `json:"order_id"`
}

/*Input Request Extend Call*/
type InputExtendCall struct {
	CallId string `json:"call_id"`
}
