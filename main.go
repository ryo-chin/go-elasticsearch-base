package main

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var svc *sqs.SQS

func main() {
	//
	// Deque ID from SQS
	//
	sess, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String("http://localhost:9324"),
		Credentials: credentials.AnonymousCredentials,
		Region:      aws.String("elasticmq")},
	)
	svc = sqs.New(sess)
	qURL := "http://localhost:9324/queue/user"

	targetUserIdMap := map[string]string{}
	defer RecoverId(qURL, targetUserIdMap)

	result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &qURL,
		MaxNumberOfMessages: aws.Int64(10),
		VisibilityTimeout:   aws.Int64(60), // 60 seconds
		WaitTimeSeconds:     aws.Int64(0),
	})
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	if len(result.Messages) == 0 {
		fmt.Println("Received no messages")
		return
	}
	fmt.Printf("Success: %+v\n", result.Messages)

	for _, message := range result.Messages {
		targetUserIdMap[*message.MessageId] = *message.Body
	}

	//
	// Select from MySQL
	//
	dbUser := "go_elastic_user"
	password := "password"
	protocol := "tcp(localhost:3309)"
	dbName := "go_elastic_db"
	db, err := gorm.Open("mysql", dbUser+":"+password+"@"+protocol+"/"+dbName)
	defer CloseDB(*db, err)

	var users []User
	var targetUserIdList []int
	for _, v := range targetUserIdMap {
		i, _ := strconv.Atoi(v)
		targetUserIdList = append(targetUserIdList, i)
	}
	db.Where("id in (?)", targetUserIdList).Find(&users)
	fmt.Printf("Success: %+v\n", users)
}

func CloseDB(db gorm.DB, err error) {
	fmt.Println("Error", err)
	db.Close()
}

func RecoverId(qURL string, targetUserIdMap map[string]string) {
	var batchRequestEntryList []*sqs.SendMessageBatchRequestEntry
	for k, v := range targetUserIdMap {
		entry := sqs.SendMessageBatchRequestEntry{Id: &k, MessageBody: &v}
		batchRequestEntryList = append(batchRequestEntryList, &entry)
	}
	svc.SendMessageBatchRequest(&sqs.SendMessageBatchInput{
		Entries:  batchRequestEntryList,
		QueueUrl: &qURL,
	})
}

type User struct {
	Id          int
	LastName    string
	FirstName   string
	Gender      string
	PhoneNumber string
	Email       string
	Password    string
}
