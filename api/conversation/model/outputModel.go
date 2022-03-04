package model

type ResponseStartConversation struct {
	Code    int                   `json:"code" example:"200"`
	Message string                `json:"message"`
	Data    DataStartConversation `json:"data"`
}

type DataStartConversation struct {
	OrderId     string `json:"order_id" example:"bc72fg9esvbsjd240v2ihidhv9"`
	CallId      string `json:"call_id" example:" bc72fg9esvbsjd240v2ihidhv9"`
	Token       string `json:"token" example:"data token agora"`
	Channel     string `json:"channel" example:"123456"`
	Uid         string `json:"uid" example:"343f9n39fj23fj40"`
	RingingTime int64  `json:"ringing_time" example:"5"`
}

type ResponseGetCallHistoryUser struct {
	Code    int                      `json:"code" example:"200"`
	Message string                   `json:"message"`
	Data    []DataGetCallHistoryUser `json:"data"`
}

type DataGetCallHistoryUser struct {
	CallId          string   `json:"call_id" bson:"call_id"`
	ListenerId      string   `json:"listener_id" bson:"listener_id"`
	ListenerName    string   `json:"listener_name" bson:"listener_name"`
	StartCall       string   `json:"start_call" bson:"start_call"`
	EndCall         string   `json:"end_call" bson:"end_call"`
	Content         string   `json:"content" bson:"content"`
	ConsultingField []string `json:"consulting_field" bson:"consulting_field"`
}

type ResponseGetCallHistoryListener struct {
	Code    int                          `json:"code" example:"200"`
	Message string                       `json:"message"`
	Data    []DataGetCallHistoryListener `json:"data"`
}

type ResponseGetDetailCallHistoryListener struct {
	Code    int                        `json:"code" example:"200"`
	Message string                     `json:"message"`
	Data    DataGetCallHistoryListener `json:"data"`
}

type DataGetCallHistoryListener struct {
	CallId                string   `json:"call_id" bson:"call_id"`
	UserId                string   `json:"user_id" bson:"user_id"`
	UserName              string   `json:"username" bson:"username"`
	PhoneNumber           string   `json:"phone_number" bson:"phone_number"`
	Status                string   `json:"status" bson:"status"`
	StartCall             string   `json:"start_call" bson:"start_call"`
	EndCall               string   `json:"end_call" bson:"end_call"`
	Content               string   `json:"content" bson:"content"`
	ConsultingField       []string `json:"consulting_field" bson:"consulting_field"`
	ConsultingFeeListener int64    `json:"consulting_fee_listener" bson:"consulting_fee_listener"`
}

type DataConversationsUser struct {
	Id            string `json:"id"`
	EmployeeId    string `json:"employee_id"`
	ListenerName  string `json:"listener_name"`
	ListenerRole  int    `json:"listener_role"`
	ListenerImage string `json:"listener_image"`
	Status        string `json:"status"`
	Date          string `json:"date"`
	TimeSlot      string `json:"time_slot"`
	BookingTime   string `json:"booking_time"`
}

type ResponseGetConversationsUser struct {
	Code    int                     `json:"code" example:"200"`
	Message string                  `json:"message"`
	Data    []DataConversationsUser `json:"data"`
}

type DataDetailConversationUser struct {
	Id              string   `json:"id"`
	ListenerName    string   `json:"listener_name"`
	ListenerRole    int      `json:"listener_role"`
	ListenerImage   string   `json:"listener_image"`
	Date            string   `json:"date"`
	TimeSlot        string   `json:"time_slot"`
	BookingTime     string   `json:"booking_time"`
	DiamondOrder    int64    `json:"diamond_order" `
	DiamondPayment  int64    `json:"diamond_payment"`
	DiamondDiscount int64    `json:"diamond_discount" `
	ConsultingField []string `json:"consulting_field"`
	Tip             int64    `json:"tip" `
	Star            int      `json:"star" `
	NoteCompany     string   `json:"note_company" `
	NoteListener    string   `json:"note_listener" `
}

type ResponseGetDetailConversationUser struct {
	Code    int                        `json:"code" example:"200"`
	Message string                     `json:"message"`
	Data    DataDetailConversationUser `json:"data"`
}
