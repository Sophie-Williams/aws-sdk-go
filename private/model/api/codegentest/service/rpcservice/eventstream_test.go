// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

// +build go1.6

package rpcservice

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/corehandlers"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/awstesting/unit"
	"github.com/aws/aws-sdk-go/private/protocol"
	"github.com/aws/aws-sdk-go/private/protocol/eventstream"
	"github.com/aws/aws-sdk-go/private/protocol/eventstream/eventstreamapi"
	"github.com/aws/aws-sdk-go/private/protocol/eventstream/eventstreamtest"
	"github.com/aws/aws-sdk-go/private/protocol/jsonrpc"
)

var _ time.Time
var _ awserr.Error

func TestEmptyStream_Read(t *testing.T) {
	expectEvents, eventMsgs := mockEmptyStreamReadEvents()
	sess, cleanupFn, err := eventstreamtest.SetupEventStreamSession(t,
		eventstreamtest.ServeEventStream{
			T:      t,
			Events: eventMsgs,
		},
		true,
	)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := New(sess)
	resp, err := svc.EmptyStream(nil)
	if err != nil {
		t.Fatalf("expect no error got, %v", err)
	}
	defer resp.GetStream().Close()
	// Trim off response output type pseudo event so only event messages remain.
	expectEvents = expectEvents[1:]

	var i int
	for event := range resp.GetStream().Events() {
		if event == nil {
			t.Errorf("%d, expect event, got nil", i)
		}
		if e, a := expectEvents[i], event; !reflect.DeepEqual(e, a) {
			t.Errorf("%d, expect %T %v, got %T %v", i, e, e, a, a)
		}
		i++
	}

	if err := resp.GetStream().Err(); err != nil {
		t.Errorf("expect no error, %v", err)
	}
}

func TestEmptyStream_ReadClose(t *testing.T) {
	_, eventMsgs := mockEmptyStreamReadEvents()
	sess, cleanupFn, err := eventstreamtest.SetupEventStreamSession(t,
		eventstreamtest.ServeEventStream{
			T:      t,
			Events: eventMsgs,
		},
		true,
	)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := New(sess)
	resp, err := svc.EmptyStream(nil)
	if err != nil {
		t.Fatalf("expect no error got, %v", err)
	}

	resp.GetStream().Close()
	<-resp.GetStream().Events()

	if err := resp.GetStream().Err(); err != nil {
		t.Errorf("expect no error, %v", err)
	}
}

func BenchmarkEmptyStream_Read(b *testing.B) {
	_, eventMsgs := mockEmptyStreamReadEvents()
	var buf bytes.Buffer
	encoder := eventstream.NewEncoder(&buf)
	for _, msg := range eventMsgs {
		if err := encoder.Encode(msg); err != nil {
			b.Fatalf("failed to encode message, %v", err)
		}
	}
	stream := &loopReader{source: bytes.NewReader(buf.Bytes())}

	sess := unit.Session
	svc := New(sess, &aws.Config{
		Endpoint:               aws.String("https://example.com"),
		DisableParamValidation: aws.Bool(true),
	})
	svc.Handlers.Send.Swap(corehandlers.SendHandler.Name,
		request.NamedHandler{Name: "mockSend",
			Fn: func(r *request.Request) {
				r.HTTPResponse = &http.Response{
					Status:     "200 OK",
					StatusCode: 200,
					Header:     http.Header{},
					Body:       ioutil.NopCloser(stream),
				}
			},
		},
	)

	resp, err := svc.EmptyStream(nil)
	if err != nil {
		b.Fatalf("failed to create request, %v", err)
	}
	defer resp.GetStream().Close()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err = resp.GetStream().Err(); err != nil {
			b.Fatalf("expect no error, got %v", err)
		}
		event := <-resp.GetStream().Events()
		if event == nil {
			b.Fatalf("expect event, got nil, %v, %d", resp.GetStream().Err(), i)
		}
	}
}

func mockEmptyStreamReadEvents() (
	[]EmptyEventStreamEvent,
	[]eventstream.Message,
) {
	expectEvents := []EmptyEventStreamEvent{
		&EmptyStreamOutput{},
	}

	var marshalers request.HandlerList
	marshalers.PushBackNamed(jsonrpc.BuildHandler)
	payloadMarshaler := protocol.HandlerPayloadMarshal{
		Marshalers: marshalers,
	}
	_ = payloadMarshaler

	eventMsgs := []eventstream.Message{
		{
			Headers: eventstream.Headers{
				eventstreamtest.EventMessageTypeHeader,
				{
					Name:  eventstreamapi.EventTypeHeader,
					Value: eventstream.StringValue("initial-response"),
				},
			},
			Payload: eventstreamtest.MarshalEventPayload(payloadMarshaler, expectEvents[0]),
		},
	}

	return expectEvents, eventMsgs
}

func TestGetEventStream_Read(t *testing.T) {
	expectEvents, eventMsgs := mockGetEventStreamReadEvents()
	sess, cleanupFn, err := eventstreamtest.SetupEventStreamSession(t,
		eventstreamtest.ServeEventStream{
			T:      t,
			Events: eventMsgs,
		},
		true,
	)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := New(sess)
	resp, err := svc.GetEventStream(nil)
	if err != nil {
		t.Fatalf("expect no error got, %v", err)
	}
	defer resp.GetStream().Close()
	expectResp := expectEvents[0].(*GetEventStreamOutput)
	if e, a := expectResp.IntVal, resp.IntVal; !reflect.DeepEqual(e, a) {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := expectResp.StrVal, resp.StrVal; !reflect.DeepEqual(e, a) {
		t.Errorf("expect %v, got %v", e, a)
	}
	// Trim off response output type pseudo event so only event messages remain.
	expectEvents = expectEvents[1:]

	var i int
	for event := range resp.GetStream().Events() {
		if event == nil {
			t.Errorf("%d, expect event, got nil", i)
		}
		if e, a := expectEvents[i], event; !reflect.DeepEqual(e, a) {
			t.Errorf("%d, expect %T %v, got %T %v", i, e, e, a, a)
		}
		i++
	}

	if err := resp.GetStream().Err(); err != nil {
		t.Errorf("expect no error, %v", err)
	}
}

func TestGetEventStream_ReadClose(t *testing.T) {
	_, eventMsgs := mockGetEventStreamReadEvents()
	sess, cleanupFn, err := eventstreamtest.SetupEventStreamSession(t,
		eventstreamtest.ServeEventStream{
			T:      t,
			Events: eventMsgs,
		},
		true,
	)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := New(sess)
	resp, err := svc.GetEventStream(nil)
	if err != nil {
		t.Fatalf("expect no error got, %v", err)
	}

	// Assert calling Err before close does not close the stream.
	resp.GetStream().Err()
	select {
	case _, ok := <-resp.GetStream().Events():
		if !ok {
			t.Fatalf("expect stream not to be closed, but was")
		}
	default:
	}

	resp.GetStream().Close()
	<-resp.GetStream().Events()

	if err := resp.GetStream().Err(); err != nil {
		t.Errorf("expect no error, %v", err)
	}
}

func BenchmarkGetEventStream_Read(b *testing.B) {
	_, eventMsgs := mockGetEventStreamReadEvents()
	var buf bytes.Buffer
	encoder := eventstream.NewEncoder(&buf)
	for _, msg := range eventMsgs {
		if err := encoder.Encode(msg); err != nil {
			b.Fatalf("failed to encode message, %v", err)
		}
	}
	stream := &loopReader{source: bytes.NewReader(buf.Bytes())}

	sess := unit.Session
	svc := New(sess, &aws.Config{
		Endpoint:               aws.String("https://example.com"),
		DisableParamValidation: aws.Bool(true),
	})
	svc.Handlers.Send.Swap(corehandlers.SendHandler.Name,
		request.NamedHandler{Name: "mockSend",
			Fn: func(r *request.Request) {
				r.HTTPResponse = &http.Response{
					Status:     "200 OK",
					StatusCode: 200,
					Header:     http.Header{},
					Body:       ioutil.NopCloser(stream),
				}
			},
		},
	)

	resp, err := svc.GetEventStream(nil)
	if err != nil {
		b.Fatalf("failed to create request, %v", err)
	}
	defer resp.GetStream().Close()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err = resp.GetStream().Err(); err != nil {
			b.Fatalf("expect no error, got %v", err)
		}
		event := <-resp.GetStream().Events()
		if event == nil {
			b.Fatalf("expect event, got nil, %v, %d", resp.GetStream().Err(), i)
		}
	}
}

func mockGetEventStreamReadEvents() (
	[]EventStreamEvent,
	[]eventstream.Message,
) {
	expectEvents := []EventStreamEvent{
		&GetEventStreamOutput{
			IntVal: aws.Int64(123),
			StrVal: aws.String("string value goes here"),
		},
		&EmptyEvent{},
		&ExplicitPayloadEvent{
			LongVal: aws.Int64(1234),
			NestedVal: &NestedShape{
				IntVal: aws.Int64(123),
				StrVal: aws.String("string value goes here"),
			},
			StringVal: aws.String("string value goes here"),
		},
		&HeaderOnlyEvent{
			BlobVal:    []byte("blob value goes here"),
			BoolVal:    aws.Bool(true),
			ByteVal:    aws.Int64(1),
			IntegerVal: aws.Int64(123),
			LongVal:    aws.Int64(1234),
			ShortVal:   aws.Int64(12),
			StringVal:  aws.String("string value goes here"),
			TimeVal:    aws.Time(time.Unix(1396594860, 0).UTC()),
		},
		&ImplicitPayloadEvent{
			ByteVal:    aws.Int64(1),
			IntegerVal: aws.Int64(123),
			ShortVal:   aws.Int64(12),
		},
		&PayloadOnlyEvent{
			NestedVal: &NestedShape{
				IntVal: aws.Int64(123),
				StrVal: aws.String("string value goes here"),
			},
		},
		&PayloadOnlyBlobEvent{
			BlobPayload: []byte("blob value goes here"),
		},
		&PayloadOnlyStringEvent{
			StringPayload: aws.String("string value goes here"),
		},
	}

	var marshalers request.HandlerList
	marshalers.PushBackNamed(jsonrpc.BuildHandler)
	payloadMarshaler := protocol.HandlerPayloadMarshal{
		Marshalers: marshalers,
	}
	_ = payloadMarshaler

	eventMsgs := []eventstream.Message{
		{
			Headers: eventstream.Headers{
				eventstreamtest.EventMessageTypeHeader,
				{
					Name:  eventstreamapi.EventTypeHeader,
					Value: eventstream.StringValue("initial-response"),
				},
			},
			Payload: eventstreamtest.MarshalEventPayload(payloadMarshaler, expectEvents[0]),
		},
		{
			Headers: eventstream.Headers{
				eventstreamtest.EventMessageTypeHeader,
				{
					Name:  eventstreamapi.EventTypeHeader,
					Value: eventstream.StringValue("Empty"),
				},
			},
		},
		{
			Headers: eventstream.Headers{
				eventstreamtest.EventMessageTypeHeader,
				{
					Name:  eventstreamapi.EventTypeHeader,
					Value: eventstream.StringValue("ExplicitPayload"),
				},
				{
					Name:  "LongVal",
					Value: eventstream.Int64Value(*expectEvents[2].(*ExplicitPayloadEvent).LongVal),
				},
				{
					Name:  "StringVal",
					Value: eventstream.StringValue(*expectEvents[2].(*ExplicitPayloadEvent).StringVal),
				},
			},
			Payload: eventstreamtest.MarshalEventPayload(payloadMarshaler, expectEvents[2]),
		},
		{
			Headers: eventstream.Headers{
				eventstreamtest.EventMessageTypeHeader,
				{
					Name:  eventstreamapi.EventTypeHeader,
					Value: eventstream.StringValue("Headers"),
				},
				{
					Name:  "BlobVal",
					Value: eventstream.BytesValue(expectEvents[3].(*HeaderOnlyEvent).BlobVal),
				},
				{
					Name:  "BoolVal",
					Value: eventstream.BoolValue(*expectEvents[3].(*HeaderOnlyEvent).BoolVal),
				},
				{
					Name:  "ByteVal",
					Value: eventstream.Int8Value(int8(*expectEvents[3].(*HeaderOnlyEvent).ByteVal)),
				},
				{
					Name:  "IntegerVal",
					Value: eventstream.Int32Value(int32(*expectEvents[3].(*HeaderOnlyEvent).IntegerVal)),
				},
				{
					Name:  "LongVal",
					Value: eventstream.Int64Value(*expectEvents[3].(*HeaderOnlyEvent).LongVal),
				},
				{
					Name:  "ShortVal",
					Value: eventstream.Int16Value(int16(*expectEvents[3].(*HeaderOnlyEvent).ShortVal)),
				},
				{
					Name:  "StringVal",
					Value: eventstream.StringValue(*expectEvents[3].(*HeaderOnlyEvent).StringVal),
				},
				{
					Name:  "TimeVal",
					Value: eventstream.TimestampValue(*expectEvents[3].(*HeaderOnlyEvent).TimeVal),
				},
			},
		},
		{
			Headers: eventstream.Headers{
				eventstreamtest.EventMessageTypeHeader,
				{
					Name:  eventstreamapi.EventTypeHeader,
					Value: eventstream.StringValue("ImplicitPayload"),
				},
				{
					Name:  "ByteVal",
					Value: eventstream.Int8Value(int8(*expectEvents[4].(*ImplicitPayloadEvent).ByteVal)),
				},
			},
			Payload: eventstreamtest.MarshalEventPayload(payloadMarshaler, expectEvents[4]),
		},
		{
			Headers: eventstream.Headers{
				eventstreamtest.EventMessageTypeHeader,
				{
					Name:  eventstreamapi.EventTypeHeader,
					Value: eventstream.StringValue("PayloadOnly"),
				},
			},
			Payload: eventstreamtest.MarshalEventPayload(payloadMarshaler, expectEvents[5]),
		},
		{
			Headers: eventstream.Headers{
				eventstreamtest.EventMessageTypeHeader,
				{
					Name:  eventstreamapi.EventTypeHeader,
					Value: eventstream.StringValue("PayloadOnlyBlob"),
				},
			},
			Payload: expectEvents[6].(*PayloadOnlyBlobEvent).BlobPayload,
		},
		{
			Headers: eventstream.Headers{
				eventstreamtest.EventMessageTypeHeader,
				{
					Name:  eventstreamapi.EventTypeHeader,
					Value: eventstream.StringValue("PayloadOnlyString"),
				},
			},
			Payload: []byte(*expectEvents[7].(*PayloadOnlyStringEvent).StringPayload),
		},
	}

	return expectEvents, eventMsgs
}
func TestGetEventStream_ReadException(t *testing.T) {
	expectEvents := []EventStreamEvent{
		&GetEventStreamOutput{
			IntVal: aws.Int64(123),
			StrVal: aws.String("string value goes here"),
		},
		&ExceptionEvent{
			IntVal:   aws.Int64(123),
			Message_: aws.String("string value goes here"),
		},
	}

	var marshalers request.HandlerList
	marshalers.PushBackNamed(jsonrpc.BuildHandler)
	payloadMarshaler := protocol.HandlerPayloadMarshal{
		Marshalers: marshalers,
	}

	eventMsgs := []eventstream.Message{
		{
			Headers: eventstream.Headers{
				eventstreamtest.EventMessageTypeHeader,
				{
					Name:  eventstreamapi.EventTypeHeader,
					Value: eventstream.StringValue("initial-response"),
				},
			},
			Payload: eventstreamtest.MarshalEventPayload(payloadMarshaler, expectEvents[0]),
		},
		{
			Headers: eventstream.Headers{
				eventstreamtest.EventExceptionTypeHeader,
				{
					Name:  eventstreamapi.ExceptionTypeHeader,
					Value: eventstream.StringValue("Exception"),
				},
			},
			Payload: eventstreamtest.MarshalEventPayload(payloadMarshaler, expectEvents[1]),
		},
	}

	sess, cleanupFn, err := eventstreamtest.SetupEventStreamSession(t,
		eventstreamtest.ServeEventStream{
			T:      t,
			Events: eventMsgs,
		},
		true,
	)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := New(sess)
	resp, err := svc.GetEventStream(nil)
	if err != nil {
		t.Fatalf("expect no error got, %v", err)
	}

	defer resp.GetStream().Close()

	<-resp.GetStream().Events()

	err = resp.GetStream().Err()
	if err == nil {
		t.Fatalf("expect err, got none")
	}

	expectErr := &ExceptionEvent{
		IntVal:   aws.Int64(123),
		Message_: aws.String("string value goes here"),
	}
	aerr, ok := err.(awserr.Error)
	if !ok {
		t.Errorf("expect exception, got %T, %#v", err, err)
	}
	if e, a := expectErr.Code(), aerr.Code(); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := expectErr.Message(), aerr.Message(); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

	if e, a := expectErr, aerr; !reflect.DeepEqual(e, a) {
		t.Errorf("expect %#v, got %#v", e, a)
	}
}

var _ awserr.Error = (*ExceptionEvent)(nil)
var _ awserr.Error = (*ExceptionEvent2)(nil)

type loopReader struct {
	source *bytes.Reader
}

func (c *loopReader) Read(p []byte) (int, error) {
	if c.source.Len() == 0 {
		c.source.Seek(0, 0)
	}

	return c.source.Read(p)
}
