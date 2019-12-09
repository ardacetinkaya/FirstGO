package azstoragequeue

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/Azure/azure-storage-queue-go/azqueue"
)

type QueueData struct {
	Value      string
	Date       time.Time
	serviceURL azqueue.ServiceURL
	queueURL   azqueue.QueueURL
	ctx        context.Context
}

var credential *azqueue.SharedKeyCredential

func CreateQueueManager() *QueueData {
	return &QueueData{}
}

func (qm *QueueData) Init(accountName string, accountKey string) error {
	credential, err := azqueue.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	p := azqueue.NewPipeline(credential, azqueue.PipelineOptions{})
	u, _ := url.Parse(fmt.Sprintf("https://%s.queue.core.windows.net", accountName))
	s := azqueue.NewServiceURL(*u, p)

	qm.serviceURL = s
	return nil
}
func (qm *QueueData) CreateQueue(name string) error {
	qm.queueURL = qm.serviceURL.NewQueueURL(name)
	qm.ctx = context.TODO()
	_, err := qm.queueURL.Create(qm.ctx, azqueue.Metadata{})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
func (qm *QueueData) Put(messageText string) error {

	messagesURL := qm.queueURL.NewMessagesURL()
	_, err := messagesURL.Enqueue(qm.ctx, messageText, time.Second*0, time.Minute)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
