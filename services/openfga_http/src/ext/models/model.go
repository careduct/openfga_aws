package models

type Model struct {
	Id       string `dynamodbav:"id" json:"id"`
	FieldOne string `dynamodbav:"fieldOne" json:"fieldOne"`
	FieldTwo string `dynamodbav:"fieldTwo" json:"fieldTwo"`
}
