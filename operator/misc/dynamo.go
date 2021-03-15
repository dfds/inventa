package misc

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoDbStore struct {
	dynamoClient *dynamodb.DynamoDB
}

type ServiceProxyTable struct {
	ServiceName string `json:"servicename"`
}

type ServiceProxyItem struct {
	ClusterName     string `json:"clustername"`
	ObjectKind      string `json:"objectkind"`
	ObjectNamespace string `json:"objectnamespace"`
	ObjectName      string `json:"objectname"`
	Identifier      string `json:"identifier"`
	Endpoint        string `json:"endpoint"`
}

func SerialiseProxyItem(serviceProxyItem ServiceProxyItem) string {
	return fmt.Sprintf("%s::%s::%s::%s::%s::%s", serviceProxyItem.ClusterName, serviceProxyItem.ObjectKind, serviceProxyItem.ObjectNamespace,
		serviceProxyItem.ObjectName, serviceProxyItem.Identifier, serviceProxyItem.Endpoint)
}

func DeserialiseProxyItem(input string) (ServiceProxyItem, error) {
	split := strings.Split(input, "::")

	item := ServiceProxyItem{}

	if len(split) == 6 {
		fmt.Println("I found something")

		item.ClusterName = split[0]
		item.ObjectKind = split[1]
		item.ObjectNamespace = split[2]
		item.ObjectName = split[3]
		item.Identifier = split[4]
		item.Endpoint = split[5]

		return item, nil
	}

	return item, errors.New("Unable to deserialise")
}

func NewDynamoDbStore() *DynamoDbStore {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
		//		LogLevel: aws.LogLevel(aws.LogDebug),
	})
	if err != nil {
		panic(err)
	}

	return &DynamoDbStore{
		dynamoClient: dynamodb.New(sess),
	}
}

func (s *DynamoDbStore) PutService(entry ServiceProxyItem) {
	dynamoEntry, err := dynamodbattribute.MarshalMap(ServiceProxyTable{ServiceName: SerialiseProxyItem(entry)})
	if err != nil {
		panic(err)
	}

	input := &dynamodb.PutItemInput{
		Item:      dynamoEntry,
		TableName: aws.String("serviceproxy-services"),
	}

	_, err = s.dynamoClient.PutItem(input)
	if err != nil {
		panic(err)
	}

	fmt.Println("Stored object ", entry.ObjectName, " to table ", "serviceproxy-services")
}

func (s *DynamoDbStore) RemoveService(entry ServiceProxyItem) {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"servicename": {
				S: aws.String(SerialiseProxyItem(entry)),
			},
		},
		TableName: aws.String("serviceproxy-services"),
	}

	_, err := s.dynamoClient.DeleteItem(input)
	if err != nil {
		panic(err)
	}

	fmt.Println("Removed object ", entry.ObjectName, " from table ", "serviceproxy-services")
}
