package mail

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sandexcare_backend/db"
)

const (
	CollectionEmailsTemplate = "emails_template"
)

/*Model Email*/
type EmailTemplateTable struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	EmailType string             `json:"email_type" bson:"email_type"`
	Template  string             `json:"template" bson:"template"`
}

/*Resource Email*/
type resourceEmailTemplateInterface interface {
	GetEmailTemplate(typeTemplate string) EmailTemplateTable
}

/*Resource Interface of project Ad*/
func NewResource() resourceEmailTemplateInterface {
	return &resourceMail{}
}

/*Init parent struct*/
type resourceMail struct {
}

func (r *resourceMail) GetEmailTemplate(typeTemplate string) EmailTemplateTable {
	var email EmailTemplateTable
	_ = collection(CollectionEmailsTemplate).FindOne(db.GetContext(), bson.M{"email_type": typeTemplate}).Decode(&email)
	return email
}
