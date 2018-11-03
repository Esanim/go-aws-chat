package chatsess

import (
	"html"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Chat represents base chat model.
type Chat struct {
	DateID   string
	Time     time.Time
	Username string
	Text     string
}

// NewChat creates a new chat instance.
func NewChat(Username, Text string) Chat {
	return Chat{
		DateID:   time.Now().Format(DateFmt),
		Time:     time.Now(),
		Username: Username,
		Text:     html.EscapeString(Text),
	}
}

// ChatFromItem creates a new chat instance based on input's parameters.
func ChatFromItem(item map[string]*dynamodb.AttributeValue) Chat {
	dateav := item["DateID"]
	timeav := item["Tmstp"]
	unameav := item["Username"]
	txav := item["Text"]

	return Chat{
		DateID:   *dateav.S,
		Time:     DBtoTime(timeav.N),
		Username: *unameav.S,
		Text:     *txav.S,
	}
}

// Put saves the chat to the database.
func (c Chat) Put(sess *session.Session) error {
	dbc := dynamodb.New(sess)
	_, err := dbc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("ch_chats"),
		Item: map[string]*dynamodb.AttributeValue{
			"DateID":   {S: aws.String(c.DateID)},
			"Tmstp":    {N: TimetoDB(c.Time)},
			"Username": {S: aws.String(c.Username)},
			"Text":     {S: aws.String(c.Text)},
		},
	})
	return err
}

// GetChat returns the current instance of the chat.
func GetChat(sess *session.Session) ([]Chat, error) {
	dbc := dynamodb.New(sess)
	dbres, err := dbc.Query(&dynamodb.QueryInput{
		TableName:              aws.String("ch_chats"),
		KeyConditionExpression: aws.String("DateID = :a"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":a": {S: aws.String(time.Now().Format(DateFmt))},
		},
	})

	if err != nil {
		return []Chat{}, err
	}

	res := []Chat{}
	for _, v := range dbres.Items {
		res = append(res, ChatFromItem(v))
	}

	return res, nil
}

// GetChatAfter returns the current instance of the chat after specified time.
func GetChatAfter(dateID string, t time.Time, sess *session.Session) ([]Chat, error) {
	dbc := dynamodb.New(sess)
	dbres, err := dbc.Query(&dynamodb.QueryInput{
		TableName:              aws.String("ch_chats"),
		KeyConditionExpression: aws.String("DateID = :a"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":a": {S: aws.String(time.Now().Format(DateFmt))},
		},
		ExclusiveStartKey: map[string]*dynamodb.AttributeValue{
			"DateID": {S: aws.String(dateID)},
			"Tmstp":  {N: TimetoDB(t)},
		},
	})

	if err != nil {
		return []Chat{}, err
	}

	res := []Chat{}
	for _, v := range dbres.Items {
		res = append(res, ChatFromItem(v))
	}

	return res, nil
}
