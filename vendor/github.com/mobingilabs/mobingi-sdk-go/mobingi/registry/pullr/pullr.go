package pullr

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/pkg/errors"
)

type QueueClient struct{}

func (qc *QueueClient) Publish(payload string) error {
	svc := sns.New(session.New())
	params := &sns.PublishInput{
		Message:  aws.String(payload),
		TopicArn: aws.String(os.Getenv("PULLR_SNS_ARN")),
	}

	res, err := svc.Publish(params)
	if err != nil {
		return errors.Wrap(err, "publish failed")
	}

	debug.Info("sns_send:", *res.MessageId)
	return nil
}

type ReadCtx struct {
	MaxNumOfMessages   int64
	VisibilityTimeout  int64
	WaitTimeSeconds    int64
	DontDeleteMessages bool
}

// Read reads message from pullr SQS. First return value is the raw AWS SDK response object. The
// first string slice are the messages, and the next is the receipt handles for easy access.
func (qc *QueueClient) Read(rc *ReadCtx) (*sqs.ReceiveMessageOutput, []string, []string, error) {
	var mxm int64 = 1   // default: 1 message
	var vt int64 = 1800 // default: 30 minutes
	var wt int64 = 10   // default: long polling
	m := make([]string, 0)
	h := make([]string, 0)

	if rc != nil {
		mxm = rc.MaxNumOfMessages
		vt = rc.VisibilityTimeout
		wt = rc.WaitTimeSeconds
	}

	qurl := os.Getenv("PULLR_SQS_URL")
	svc := sqs.New(session.New())
	res, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &qurl,
		MaxNumberOfMessages: aws.Int64(mxm),
		VisibilityTimeout:   aws.Int64(vt),
		WaitTimeSeconds:     aws.Int64(wt),
	})

	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "recv msg failed")
	}

	for _, item := range res.Messages {
		m = append(m, *item.Body)
		h = append(h, *item.ReceiptHandle)
		debug.Info("sqs_recv:", *item.ReceiptHandle)

		// delete message(s) by default
		dontdelm := false
		if rc != nil {
			dontdelm = rc.DontDeleteMessages
		}

		if !dontdelm {
			delres, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl:      &qurl,
				ReceiptHandle: item.ReceiptHandle,
			})

			if err != nil {
				debug.Error(delres, err)
			} else {
				debug.Info("sqs_del:", *item.ReceiptHandle)
			}
		}
	}

	return res, m, h, nil
}
