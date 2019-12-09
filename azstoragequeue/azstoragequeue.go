package azstoragequeue

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/Azure/azure-storage-queue-go/azqueue"
)

type AzureQueueOption func(q *AzureQueue)
type AzureQueue struct {
	Name       string
	serviceURL azqueue.ServiceURL
	queueURL   azqueue.QueueURL
	ctx        context.Context
}

var credential *azqueue.SharedKeyCredential
var logger *log.Logger

func CreateQueueManager(opts ...AzureQueueOption) *AzureQueue {
	qm := &AzureQueue{}
	logger = log.New(os.Stdout, "", log.LstdFlags)
	for _, opt := range opts {
		opt(qm)
	}
	return qm
}

func (qm *AzureQueue) Init(accountName string, accountKey string) error {
	credential, err := azqueue.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		logger.Panicln(err.Error())
		return err
	}
	p := azqueue.NewPipeline(credential, azqueue.PipelineOptions{})
	u, _ := url.Parse(fmt.Sprintf("https://%s.queue.core.windows.net", accountName))
	s := azqueue.NewServiceURL(*u, p)

	qm.serviceURL = s
	return nil
}
func (qm *AzureQueue) CreateQueue(name string) error {
	qm.queueURL = qm.serviceURL.NewQueueURL(name)
	qm.ctx = context.TODO()
	_, err := qm.queueURL.Create(qm.ctx, azqueue.Metadata{})
	if err != nil {
		logger.Println(err.Error())
		return err
	}

	return nil
}
func (qm *AzureQueue) Put(messageText string) error {

	messagesURL := qm.queueURL.NewMessagesURL()
	_, err := messagesURL.Enqueue(qm.ctx, messageText, time.Second*0, time.Minute)
	if err != nil {
		logger.Println(err.Error())
		return err
	}

	return nil
}
