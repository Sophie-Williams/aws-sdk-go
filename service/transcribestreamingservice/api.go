// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package transcribestreamingservice

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/private/protocol"
	"github.com/aws/aws-sdk-go/private/protocol/eventstream"
	"github.com/aws/aws-sdk-go/private/protocol/eventstream/eventstreamapi"
	"github.com/aws/aws-sdk-go/private/protocol/rest"
	"github.com/aws/aws-sdk-go/private/protocol/restjson"
)

const opStartStreamTranscription = "StartStreamTranscription"

// StartStreamTranscriptionRequest generates a "aws/request.Request" representing the
// client's request for the StartStreamTranscription operation. The "output" return
// value will be populated with the request's response once the request completes
// successfully.
//
// Use "Send" method on the returned Request to send the API call to the service.
// the "output" return value is not valid until after Send returns without error.
//
// See StartStreamTranscription for more information on using the StartStreamTranscription
// API call, and error handling.
//
// This method is useful when you want to inject custom logic or configuration
// into the SDK's request lifecycle. Such as custom headers, or retry logic.
//
//
//    // Example sending a request using the StartStreamTranscriptionRequest method.
//    req, resp := client.StartStreamTranscriptionRequest(params)
//
//    err := req.Send()
//    if err == nil { // resp is now filled
//        fmt.Println(resp)
//    }
//
// See also, https://docs.aws.amazon.com/goto/WebAPI/transcribe-streaming-2017-10-26/StartStreamTranscription
func (c *TranscribeStreamingService) StartStreamTranscriptionRequest(input *StartStreamTranscriptionInput) (req *request.Request, output *StartStreamTranscriptionOutput) {
	op := &request.Operation{
		Name:       opStartStreamTranscription,
		HTTPMethod: "POST",
		HTTPPath:   "/stream-transcription",
	}

	if input == nil {
		input = &StartStreamTranscriptionInput{}
	}

	output = &StartStreamTranscriptionOutput{}
	req = c.newRequest(op, input, output)
	req.Handlers.UnmarshalMeta.PushBack(
		protocol.RequireHTTPMinProtocol{Major: 2}.Handler,
	)

	es := newStartStreamTranscriptionEventStream()
	output.eventStream = es

	inputReader, inputWriter := io.Pipe()
	req.SetReaderBody(aws.ReadSeekCloser(inputReader))
	es.inputWriter = inputWriter

	req.Handlers.Build.PushBack(request.WithSetRequestHeaders(map[string]string{
		"Content-Type":         "application/vnd.amazon.eventstream",
		"X-Amz-Content-Sha256": "STREAMING-AWS4-HMAC-SHA256-EVENTS",
	}))
	req.Handlers.Build.Swap(restjson.BuildHandler.Name, rest.BuildHandler)
	req.Handlers.Send.Swap(client.LogHTTPRequestHandler.Name, client.LogHTTPRequestHeaderHandler)
	req.Handlers.Unmarshal.PushBack(es.runInputStream)

	req.Handlers.Send.Swap(client.LogHTTPResponseHandler.Name, client.LogHTTPResponseHeaderHandler)
	req.Handlers.Unmarshal.Swap(restjson.UnmarshalHandler.Name, rest.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBack(es.runOutputStream)
	req.Handlers.Unmarshal.PushBack(es.runOnStreamPartClose)
	return
}

// StartStreamTranscription API operation for Amazon Transcribe Streaming Service.
//
// Starts a bidirectional HTTP2 stream where audio is streamed to Amazon Transcribe
// and the transcription results are streamed to your application.
//
// The following are encoded as HTTP2 headers:
//
//    * x-amzn-transcribe-language-code
//
//    * x-amzn-transcribe-media-encoding
//
//    * x-amzn-transcribe-sample-rate
//
//    * x-amzn-transcribe-session-id
//
// Returns awserr.Error for service API and SDK errors. Use runtime type assertions
// with awserr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the AWS API reference guide for Amazon Transcribe Streaming Service's
// API operation StartStreamTranscription for usage and error information.
//
// Returned Error Codes:
//   * ErrCodeBadRequestException "BadRequestException"
//   One or more arguments to the StartStreamTranscription operation was invalid.
//   For example, MediaEncoding was not set to pcm or LanguageCode was not set
//   to a valid code. Check the parameters and try your request again.
//
//   * ErrCodeLimitExceededException "LimitExceededException"
//   You have exceeded the maximum number of concurrent transcription streams,
//   are starting transcription streams too quickly, or the maximum audio length
//   of 4 hours. Wait until a stream has finished processing, or break your audio
//   stream into smaller chunks and try your request again.
//
//   * ErrCodeInternalFailureException "InternalFailureException"
//   A problem occurred while processing the audio. Amazon Transcribe terminated
//   processing. Try your request again.
//
//   * ErrCodeConflictException "ConflictException"
//   A new stream started with the same session ID. The current stream has been
//   terminated.
//
// See also, https://docs.aws.amazon.com/goto/WebAPI/transcribe-streaming-2017-10-26/StartStreamTranscription
func (c *TranscribeStreamingService) StartStreamTranscription(input *StartStreamTranscriptionInput) (*StartStreamTranscriptionOutput, error) {
	req, out := c.StartStreamTranscriptionRequest(input)
	return out, req.Send()
}

// StartStreamTranscriptionWithContext is the same as StartStreamTranscription with the addition of
// the ability to pass a context and additional request options.
//
// See StartStreamTranscription for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *TranscribeStreamingService) StartStreamTranscriptionWithContext(ctx aws.Context, input *StartStreamTranscriptionInput, opts ...request.Option) (*StartStreamTranscriptionOutput, error) {
	req, out := c.StartStreamTranscriptionRequest(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

// StartStreamTranscriptionEventStream provides the event stream handling for the StartStreamTranscription.
type StartStreamTranscriptionEventStream struct {

	// Writer is the EventStream writer for the AudioStream
	// events. This value is automatically set by the SDK when the API call is made
	// Use this member when unit testing your code with the SDK to mock out the
	// EventStream Writer.
	//
	// Must not be nil.
	Writer AudioStreamWriter

	inputWriter io.WriteCloser

	// Reader is the EventStream reader for the TranscriptResultStream
	// events. This value is automatically set by the SDK when the API call is made
	// Use this member when unit testing your code with the SDK to mock out the
	// EventStream Reader.
	//
	// Must not be nil.
	Reader TranscriptResultStreamReader

	outputReader io.ReadCloser

	done      chan struct{}
	closeOnce sync.Once
	err       *eventstreamapi.OnceError
}

func newStartStreamTranscriptionEventStream() *StartStreamTranscriptionEventStream {
	return &StartStreamTranscriptionEventStream{
		done: make(chan struct{}),
		err:  eventstreamapi.NewOnceError(),
	}
}

func (es *StartStreamTranscriptionEventStream) runOnStreamPartClose(r *request.Request) {
	if es.done == nil {
		return
	}
	go es.waitStreamPartClose()

}

func (es *StartStreamTranscriptionEventStream) waitStreamPartClose() {
	var inputC <-chan struct{}
	if v, ok := es.Writer.(interface{ ErrorSet() <-chan struct{} }); ok {
		inputC = v.ErrorSet()
	}
	var outputC <-chan struct{}
	if v, ok := es.Reader.(interface{ ErrorSet() <-chan struct{} }); ok {
		outputC = v.ErrorSet()
	}

	select {
	case <-es.done:
	case <-inputC:
		es.err.SetError(es.Writer.Err())
		es.Close()
	case <-outputC:
		es.err.SetError(es.Reader.Err())
		es.Close()
	}
}

// Send writes the event to the stream blocking until the event is written.
// Returns an error if the event was not written.
//
// These events are:
//
//     * AudioEvent
func (es *StartStreamTranscriptionEventStream) Send(ctx aws.Context, event AudioStreamEvent) error {
	return es.Writer.Send(ctx, event)
}

func (es *StartStreamTranscriptionEventStream) runInputStream(r *request.Request) {
	var opts []func(*eventstream.Encoder)
	if r.Config.Logger != nil && r.Config.LogLevel.Matches(aws.LogDebugWithEventStreamBody) {
		opts = append(opts, eventstream.EncodeWithLogger(r.Config.Logger))
	}
	var encoder eventstreamapi.Encoder = eventstream.NewEncoder(es.inputWriter, opts...)

	var closer aws.MultiCloser
	sigSeed, err := v4.GetSignedRequestSignature(r.HTTPRequest)
	if err != nil {
		r.Error = awserr.New(request.ErrCodeSerialization,
			"unable to get initial request's signature", err)
		return
	}
	signer := eventstreamapi.NewSignEncoder(
		v4.NewStreamSigner(r.ClientInfo.SigningRegion, r.ClientInfo.SigningName,
			sigSeed, r.Config.Credentials),
		encoder,
	)
	encoder = signer
	closer = append(closer, signer)
	closer = append(closer, es.inputWriter)

	eventWriter := eventstreamapi.NewEventWriter(encoder,
		protocol.HandlerPayloadMarshal{
			Marshalers: r.Handlers.BuildStream,
		},
		eventTypeForAudioStreamEvent,
	)

	es.Writer = &writeAudioStream{
		StreamWriter: eventstreamapi.NewStreamWriter(eventWriter, closer),
	}
}

// Events returns a channel to read events from.
//
// These events are:
//
//     * TranscriptEvent
func (es *StartStreamTranscriptionEventStream) Events() <-chan TranscriptResultStreamEvent {
	return es.Reader.Events()
}

func (es *StartStreamTranscriptionEventStream) runOutputStream(r *request.Request) {
	var opts []func(*eventstream.Decoder)
	if r.Config.Logger != nil && r.Config.LogLevel.Matches(aws.LogDebugWithEventStreamBody) {
		opts = append(opts, eventstream.DecodeWithLogger(r.Config.Logger))
	}

	decoder := eventstream.NewDecoder(r.HTTPResponse.Body, opts...)
	eventReader := eventstreamapi.NewEventReader(decoder,
		protocol.HandlerPayloadUnmarshal{
			Unmarshalers: r.Handlers.UnmarshalStream,
		},
		unmarshalerForTranscriptResultStreamEvent,
	)

	es.outputReader = r.HTTPResponse.Body
	es.Reader = newReadTranscriptResultStream(eventReader)
}

// Close closes the stream. This will also cause the stream to be closed.
// Close must be called when done using the stream API. Not calling Close
// may result in resource leaks.
//
// Will close the underlying EventStream writer, and no more events can be
// sent.
//
// You can use the closing of the Reader's Events channel to terminate your
// application's read from the API's stream.
//
func (es *StartStreamTranscriptionEventStream) Close() (err error) {
	es.closeOnce.Do(es.safeClose)
	return es.Err()
}

func (es *StartStreamTranscriptionEventStream) safeClose() {
	if es.done != nil {
		close(es.done)
	}

	t := time.NewTicker(time.Second)
	defer t.Stop()
	writeCloseDone := make(chan error)
	go func() {
		if err := es.Writer.Close(); err != nil {
			es.err.SetError(err)
		}
		close(writeCloseDone)
	}()
	select {
	case <-t.C:
	case <-writeCloseDone:
	}
	if es.inputWriter != nil {
		es.inputWriter.Close()
	}

	es.Reader.Close()
	if es.outputReader != nil {
		es.outputReader.Close()
	}
}

// Err returns any error that occurred while reading or writing EventStream
// Events from the service API's response. Returns nil if there were no errors.
func (es *StartStreamTranscriptionEventStream) Err() error {
	if err := es.err.Err(); err != nil {
		return err
	}
	if err := es.Writer.Err(); err != nil {
		return err
	}
	if err := es.Reader.Err(); err != nil {
		return err
	}

	return nil
}

// A list of possible transcriptions for the audio.
type Alternative struct {
	_ struct{} `type:"structure"`

	// One or more alternative interpretations of the input audio.
	Items []*Item `type:"list"`

	// The text that was transcribed from the audio.
	Transcript *string `type:"string"`
}

// String returns the string representation
func (s Alternative) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s Alternative) GoString() string {
	return s.String()
}

// SetItems sets the Items field's value.
func (s *Alternative) SetItems(v []*Item) *Alternative {
	s.Items = v
	return s
}

// SetTranscript sets the Transcript field's value.
func (s *Alternative) SetTranscript(v string) *Alternative {
	s.Transcript = &v
	return s
}

// Provides a wrapper for the audio chunks that you are sending.
type AudioEvent struct {
	_ struct{} `type:"structure" payload:"AudioChunk"`

	// An audio blob that contains the next part of the audio that you want to transcribe.
	//
	// AudioChunk is automatically base64 encoded/decoded by the SDK.
	AudioChunk []byte `type:"blob"`
}

// String returns the string representation
func (s AudioEvent) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s AudioEvent) GoString() string {
	return s.String()
}

// SetAudioChunk sets the AudioChunk field's value.
func (s *AudioEvent) SetAudioChunk(v []byte) *AudioEvent {
	s.AudioChunk = v
	return s
}

// The AudioEvent is and event in the AudioStream group of events.
func (s *AudioEvent) eventAudioStream() {}

// UnmarshalEvent unmarshals the EventStream Message into the AudioEvent value.
// This method is only used internally within the SDK's EventStream handling.
func (s *AudioEvent) UnmarshalEvent(
	payloadUnmarshaler protocol.PayloadUnmarshaler,
	msg eventstream.Message,
) error {
	s.AudioChunk = make([]byte, len(msg.Payload))
	copy(s.AudioChunk, msg.Payload)
	return nil
}

func (s *AudioEvent) MarshalEvent(pm protocol.PayloadMarshaler) (msg eventstream.Message, err error) {
	msg.Headers.Set(eventstreamapi.MessageTypeHeader, eventstream.StringValue(eventstreamapi.EventMessageType))
	msg.Headers.Set(":content-type", eventstream.StringValue("application/octet-stream"))
	msg.Payload = s.AudioChunk
	return msg, err
}

// AudioStreamEvent groups together all EventStream
// events writes for AudioStream.
//
// These events are:
//
//     * AudioEvent
type AudioStreamEvent interface {
	eventAudioStream()
	eventstreamapi.Marshaler
	eventstreamapi.Unmarshaler
}

// AudioStreamWriter provides the interface for writing events to the stream.
// The default implementation for this interface will be AudioStream.
//
// The writer's Close method must allow multiple concurrent calls.
//
// These events are:
//
//     * AudioEvent
type AudioStreamWriter interface {
	// Sends writes events to the stream blocking until the event has been
	// written. An error is returned if the write fails.
	Send(aws.Context, AudioStreamEvent) error

	// Close will stop the writer writing to the event stream.
	Close() error

	// Returns any error that has occurred while writing to the event stream.
	Err() error
}

type writeAudioStream struct {
	*eventstreamapi.StreamWriter
}

func (w *writeAudioStream) Send(ctx aws.Context, event AudioStreamEvent) error {
	return w.StreamWriter.Send(ctx, event)
}

func eventTypeForAudioStreamEvent(event eventstreamapi.Marshaler) (string, error) {
	switch event.(type) {
	case *AudioEvent:
		return "AudioEvent", nil
	default:
		return "", awserr.New(
			request.ErrCodeSerialization,
			fmt.Sprintf("unknown event type, %T, for AudioStream", event),
			nil,
		)
	}
}

// One or more arguments to the StartStreamTranscription operation was invalid.
// For example, MediaEncoding was not set to pcm or LanguageCode was not set
// to a valid code. Check the parameters and try your request again.
type BadRequestException struct {
	_ struct{} `type:"structure"`

	Message_ *string `locationName:"Message" type:"string"`
}

// String returns the string representation
func (s BadRequestException) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s BadRequestException) GoString() string {
	return s.String()
}

// The BadRequestException is and event in the TranscriptResultStream group of events.
func (s *BadRequestException) eventTranscriptResultStream() {}

// UnmarshalEvent unmarshals the EventStream Message into the BadRequestException value.
// This method is only used internally within the SDK's EventStream handling.
func (s *BadRequestException) UnmarshalEvent(
	payloadUnmarshaler protocol.PayloadUnmarshaler,
	msg eventstream.Message,
) error {
	if err := payloadUnmarshaler.UnmarshalPayload(
		bytes.NewReader(msg.Payload), s,
	); err != nil {
		return err
	}
	return nil
}

func (s *BadRequestException) MarshalEvent(pm protocol.PayloadMarshaler) (msg eventstream.Message, err error) {
	msg.Headers.Set(eventstreamapi.MessageTypeHeader, eventstream.StringValue(eventstreamapi.ExceptionMessageType))
	var buf bytes.Buffer
	if err = pm.MarshalPayload(&buf, s); err != nil {
		return eventstream.Message{}, err
	}
	msg.Payload = buf.Bytes()
	return msg, err
}

// Code returns the exception type name.
func (s BadRequestException) Code() string {
	return "BadRequestException"
}

// Message returns the exception's message.
func (s BadRequestException) Message() string {
	return *s.Message_
}

// OrigErr always returns nil, satisfies awserr.Error interface.
func (s BadRequestException) OrigErr() error {
	return nil
}

func (s BadRequestException) Error() string {
	return fmt.Sprintf("%s: %s", s.Code(), s.Message())
}

// A new stream started with the same session ID. The current stream has been
// terminated.
type ConflictException struct {
	_ struct{} `type:"structure"`

	Message_ *string `locationName:"Message" type:"string"`
}

// String returns the string representation
func (s ConflictException) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ConflictException) GoString() string {
	return s.String()
}

// The ConflictException is and event in the TranscriptResultStream group of events.
func (s *ConflictException) eventTranscriptResultStream() {}

// UnmarshalEvent unmarshals the EventStream Message into the ConflictException value.
// This method is only used internally within the SDK's EventStream handling.
func (s *ConflictException) UnmarshalEvent(
	payloadUnmarshaler protocol.PayloadUnmarshaler,
	msg eventstream.Message,
) error {
	if err := payloadUnmarshaler.UnmarshalPayload(
		bytes.NewReader(msg.Payload), s,
	); err != nil {
		return err
	}
	return nil
}

func (s *ConflictException) MarshalEvent(pm protocol.PayloadMarshaler) (msg eventstream.Message, err error) {
	msg.Headers.Set(eventstreamapi.MessageTypeHeader, eventstream.StringValue(eventstreamapi.ExceptionMessageType))
	var buf bytes.Buffer
	if err = pm.MarshalPayload(&buf, s); err != nil {
		return eventstream.Message{}, err
	}
	msg.Payload = buf.Bytes()
	return msg, err
}

// Code returns the exception type name.
func (s ConflictException) Code() string {
	return "ConflictException"
}

// Message returns the exception's message.
func (s ConflictException) Message() string {
	return *s.Message_
}

// OrigErr always returns nil, satisfies awserr.Error interface.
func (s ConflictException) OrigErr() error {
	return nil
}

func (s ConflictException) Error() string {
	return fmt.Sprintf("%s: %s", s.Code(), s.Message())
}

// A problem occurred while processing the audio. Amazon Transcribe terminated
// processing. Try your request again.
type InternalFailureException struct {
	_ struct{} `type:"structure"`

	Message_ *string `locationName:"Message" type:"string"`
}

// String returns the string representation
func (s InternalFailureException) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s InternalFailureException) GoString() string {
	return s.String()
}

// The InternalFailureException is and event in the TranscriptResultStream group of events.
func (s *InternalFailureException) eventTranscriptResultStream() {}

// UnmarshalEvent unmarshals the EventStream Message into the InternalFailureException value.
// This method is only used internally within the SDK's EventStream handling.
func (s *InternalFailureException) UnmarshalEvent(
	payloadUnmarshaler protocol.PayloadUnmarshaler,
	msg eventstream.Message,
) error {
	if err := payloadUnmarshaler.UnmarshalPayload(
		bytes.NewReader(msg.Payload), s,
	); err != nil {
		return err
	}
	return nil
}

func (s *InternalFailureException) MarshalEvent(pm protocol.PayloadMarshaler) (msg eventstream.Message, err error) {
	msg.Headers.Set(eventstreamapi.MessageTypeHeader, eventstream.StringValue(eventstreamapi.ExceptionMessageType))
	var buf bytes.Buffer
	if err = pm.MarshalPayload(&buf, s); err != nil {
		return eventstream.Message{}, err
	}
	msg.Payload = buf.Bytes()
	return msg, err
}

// Code returns the exception type name.
func (s InternalFailureException) Code() string {
	return "InternalFailureException"
}

// Message returns the exception's message.
func (s InternalFailureException) Message() string {
	return *s.Message_
}

// OrigErr always returns nil, satisfies awserr.Error interface.
func (s InternalFailureException) OrigErr() error {
	return nil
}

func (s InternalFailureException) Error() string {
	return fmt.Sprintf("%s: %s", s.Code(), s.Message())
}

// A word or phrase transcribed from the input audio.
type Item struct {
	_ struct{} `type:"structure"`

	// The word or punctuation that was recognized in the input audio.
	Content *string `type:"string"`

	// The offset from the beginning of the audio stream to the end of the audio
	// that resulted in the item.
	EndTime *float64 `type:"double"`

	// The offset from the beginning of the audio stream to the beginning of the
	// audio that resulted in the item.
	StartTime *float64 `type:"double"`

	// The type of the item. PRONUNCIATION indicates that the item is a word that
	// was recognized in the input audio. PUNCTUATION indicates that the item was
	// interpreted as a pause in the input audio.
	Type *string `type:"string" enum:"ItemType"`
}

// String returns the string representation
func (s Item) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s Item) GoString() string {
	return s.String()
}

// SetContent sets the Content field's value.
func (s *Item) SetContent(v string) *Item {
	s.Content = &v
	return s
}

// SetEndTime sets the EndTime field's value.
func (s *Item) SetEndTime(v float64) *Item {
	s.EndTime = &v
	return s
}

// SetStartTime sets the StartTime field's value.
func (s *Item) SetStartTime(v float64) *Item {
	s.StartTime = &v
	return s
}

// SetType sets the Type field's value.
func (s *Item) SetType(v string) *Item {
	s.Type = &v
	return s
}

// You have exceeded the maximum number of concurrent transcription streams,
// are starting transcription streams too quickly, or the maximum audio length
// of 4 hours. Wait until a stream has finished processing, or break your audio
// stream into smaller chunks and try your request again.
type LimitExceededException struct {
	_ struct{} `type:"structure"`

	Message_ *string `locationName:"Message" type:"string"`
}

// String returns the string representation
func (s LimitExceededException) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s LimitExceededException) GoString() string {
	return s.String()
}

// The LimitExceededException is and event in the TranscriptResultStream group of events.
func (s *LimitExceededException) eventTranscriptResultStream() {}

// UnmarshalEvent unmarshals the EventStream Message into the LimitExceededException value.
// This method is only used internally within the SDK's EventStream handling.
func (s *LimitExceededException) UnmarshalEvent(
	payloadUnmarshaler protocol.PayloadUnmarshaler,
	msg eventstream.Message,
) error {
	if err := payloadUnmarshaler.UnmarshalPayload(
		bytes.NewReader(msg.Payload), s,
	); err != nil {
		return err
	}
	return nil
}

func (s *LimitExceededException) MarshalEvent(pm protocol.PayloadMarshaler) (msg eventstream.Message, err error) {
	msg.Headers.Set(eventstreamapi.MessageTypeHeader, eventstream.StringValue(eventstreamapi.ExceptionMessageType))
	var buf bytes.Buffer
	if err = pm.MarshalPayload(&buf, s); err != nil {
		return eventstream.Message{}, err
	}
	msg.Payload = buf.Bytes()
	return msg, err
}

// Code returns the exception type name.
func (s LimitExceededException) Code() string {
	return "LimitExceededException"
}

// Message returns the exception's message.
func (s LimitExceededException) Message() string {
	return *s.Message_
}

// OrigErr always returns nil, satisfies awserr.Error interface.
func (s LimitExceededException) OrigErr() error {
	return nil
}

func (s LimitExceededException) Error() string {
	return fmt.Sprintf("%s: %s", s.Code(), s.Message())
}

// The result of transcribing a portion of the input audio stream.
type Result struct {
	_ struct{} `type:"structure"`

	// A list of possible transcriptions for the audio. Each alternative typically
	// contains one item that contains the result of the transcription.
	Alternatives []*Alternative `type:"list"`

	// The offset in milliseconds from the beginning of the audio stream to the
	// end of the result.
	EndTime *float64 `type:"double"`

	// true to indicate that Amazon Transcribe has additional transcription data
	// to send, false to indicate that this is the last transcription result for
	// the audio stream.
	IsPartial *bool `type:"boolean"`

	// A unique identifier for the result.
	ResultId *string `type:"string"`

	// The offset in milliseconds from the beginning of the audio stream to the
	// beginning of the result.
	StartTime *float64 `type:"double"`
}

// String returns the string representation
func (s Result) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s Result) GoString() string {
	return s.String()
}

// SetAlternatives sets the Alternatives field's value.
func (s *Result) SetAlternatives(v []*Alternative) *Result {
	s.Alternatives = v
	return s
}

// SetEndTime sets the EndTime field's value.
func (s *Result) SetEndTime(v float64) *Result {
	s.EndTime = &v
	return s
}

// SetIsPartial sets the IsPartial field's value.
func (s *Result) SetIsPartial(v bool) *Result {
	s.IsPartial = &v
	return s
}

// SetResultId sets the ResultId field's value.
func (s *Result) SetResultId(v string) *Result {
	s.ResultId = &v
	return s
}

// SetStartTime sets the StartTime field's value.
func (s *Result) SetStartTime(v float64) *Result {
	s.StartTime = &v
	return s
}

type StartStreamTranscriptionInput struct {
	_ struct{} `type:"structure" payload:"AudioStream"`

	// Indicates the language used in the input audio stream.
	//
	// LanguageCode is a required field
	LanguageCode *string `location:"header" locationName:"x-amzn-transcribe-language-code" type:"string" required:"true" enum:"LanguageCode"`

	// The encoding used for the input audio.
	//
	// MediaEncoding is a required field
	MediaEncoding *string `location:"header" locationName:"x-amzn-transcribe-media-encoding" type:"string" required:"true" enum:"MediaEncoding"`

	// The sample rate, in Hertz, of the input audio. We suggest that you use 8000
	// Hz for low quality audio and 16000 Hz for high quality audio.
	//
	// MediaSampleRateHertz is a required field
	MediaSampleRateHertz *int64 `location:"header" locationName:"x-amzn-transcribe-sample-rate" min:"8000" type:"integer" required:"true"`

	// A identifier for the transcription session. Use this parameter when you want
	// to retry a session. If you don't provide a session ID, Amazon Transcribe
	// will generate one for you and return it in the response.
	SessionId *string `location:"header" locationName:"x-amzn-transcribe-session-id" type:"string"`

	// The name of the vocabulary to use when processing the transcription job.
	VocabularyName *string `location:"header" locationName:"x-amzn-transcribe-vocabulary-name" min:"1" type:"string"`
}

// String returns the string representation
func (s StartStreamTranscriptionInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s StartStreamTranscriptionInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *StartStreamTranscriptionInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "StartStreamTranscriptionInput"}
	if s.LanguageCode == nil {
		invalidParams.Add(request.NewErrParamRequired("LanguageCode"))
	}
	if s.MediaEncoding == nil {
		invalidParams.Add(request.NewErrParamRequired("MediaEncoding"))
	}
	if s.MediaSampleRateHertz == nil {
		invalidParams.Add(request.NewErrParamRequired("MediaSampleRateHertz"))
	}
	if s.MediaSampleRateHertz != nil && *s.MediaSampleRateHertz < 8000 {
		invalidParams.Add(request.NewErrParamMinValue("MediaSampleRateHertz", 8000))
	}
	if s.VocabularyName != nil && len(*s.VocabularyName) < 1 {
		invalidParams.Add(request.NewErrParamMinLen("VocabularyName", 1))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetLanguageCode sets the LanguageCode field's value.
func (s *StartStreamTranscriptionInput) SetLanguageCode(v string) *StartStreamTranscriptionInput {
	s.LanguageCode = &v
	return s
}

// SetMediaEncoding sets the MediaEncoding field's value.
func (s *StartStreamTranscriptionInput) SetMediaEncoding(v string) *StartStreamTranscriptionInput {
	s.MediaEncoding = &v
	return s
}

// SetMediaSampleRateHertz sets the MediaSampleRateHertz field's value.
func (s *StartStreamTranscriptionInput) SetMediaSampleRateHertz(v int64) *StartStreamTranscriptionInput {
	s.MediaSampleRateHertz = &v
	return s
}

// SetSessionId sets the SessionId field's value.
func (s *StartStreamTranscriptionInput) SetSessionId(v string) *StartStreamTranscriptionInput {
	s.SessionId = &v
	return s
}

// SetVocabularyName sets the VocabularyName field's value.
func (s *StartStreamTranscriptionInput) SetVocabularyName(v string) *StartStreamTranscriptionInput {
	s.VocabularyName = &v
	return s
}

type StartStreamTranscriptionOutput struct {
	_ struct{} `type:"structure" payload:"TranscriptResultStream"`

	eventStream *StartStreamTranscriptionEventStream

	// The language code for the input audio stream.
	LanguageCode *string `location:"header" locationName:"x-amzn-transcribe-language-code" type:"string" enum:"LanguageCode"`

	// The encoding used for the input audio stream.
	MediaEncoding *string `location:"header" locationName:"x-amzn-transcribe-media-encoding" type:"string" enum:"MediaEncoding"`

	// The sample rate for the input audio stream. Use 8000 Hz for low quality audio
	// and 16000 Hz for high quality audio.
	MediaSampleRateHertz *int64 `location:"header" locationName:"x-amzn-transcribe-sample-rate" min:"8000" type:"integer"`

	// An identifier for the streaming transcription.
	RequestId *string `location:"header" locationName:"x-amzn-request-id" type:"string"`

	// An identifier for a specific transcription session.
	SessionId *string `location:"header" locationName:"x-amzn-transcribe-session-id" type:"string"`

	// The name of the vocabulary used when processing the job.
	VocabularyName *string `location:"header" locationName:"x-amzn-transcribe-vocabulary-name" min:"1" type:"string"`
}

// String returns the string representation
func (s StartStreamTranscriptionOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s StartStreamTranscriptionOutput) GoString() string {
	return s.String()
}

// SetLanguageCode sets the LanguageCode field's value.
func (s *StartStreamTranscriptionOutput) SetLanguageCode(v string) *StartStreamTranscriptionOutput {
	s.LanguageCode = &v
	return s
}

// SetMediaEncoding sets the MediaEncoding field's value.
func (s *StartStreamTranscriptionOutput) SetMediaEncoding(v string) *StartStreamTranscriptionOutput {
	s.MediaEncoding = &v
	return s
}

// SetMediaSampleRateHertz sets the MediaSampleRateHertz field's value.
func (s *StartStreamTranscriptionOutput) SetMediaSampleRateHertz(v int64) *StartStreamTranscriptionOutput {
	s.MediaSampleRateHertz = &v
	return s
}

// SetRequestId sets the RequestId field's value.
func (s *StartStreamTranscriptionOutput) SetRequestId(v string) *StartStreamTranscriptionOutput {
	s.RequestId = &v
	return s
}

// SetSessionId sets the SessionId field's value.
func (s *StartStreamTranscriptionOutput) SetSessionId(v string) *StartStreamTranscriptionOutput {
	s.SessionId = &v
	return s
}

// SetVocabularyName sets the VocabularyName field's value.
func (s *StartStreamTranscriptionOutput) SetVocabularyName(v string) *StartStreamTranscriptionOutput {
	s.VocabularyName = &v
	return s
}

// GetStream returns the type to interact with the event stream.
func (s *StartStreamTranscriptionOutput) GetStream() *StartStreamTranscriptionEventStream {
	return s.eventStream
}

// The transcription in a TranscriptEvent.
type Transcript struct {
	_ struct{} `type:"structure"`

	// Result objects that contain the results of transcribing a portion of the
	// input audio stream. The array can be empty.
	Results []*Result `type:"list"`
}

// String returns the string representation
func (s Transcript) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s Transcript) GoString() string {
	return s.String()
}

// SetResults sets the Results field's value.
func (s *Transcript) SetResults(v []*Result) *Transcript {
	s.Results = v
	return s
}

// Represents a set of transcription results from the server to the client.
// It contains one or more segments of the transcription.
type TranscriptEvent struct {
	_ struct{} `type:"structure"`

	// The transcription of the audio stream. The transcription is composed of all
	// of the items in the results list.
	Transcript *Transcript `type:"structure"`
}

// String returns the string representation
func (s TranscriptEvent) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s TranscriptEvent) GoString() string {
	return s.String()
}

// SetTranscript sets the Transcript field's value.
func (s *TranscriptEvent) SetTranscript(v *Transcript) *TranscriptEvent {
	s.Transcript = v
	return s
}

// The TranscriptEvent is and event in the TranscriptResultStream group of events.
func (s *TranscriptEvent) eventTranscriptResultStream() {}

// UnmarshalEvent unmarshals the EventStream Message into the TranscriptEvent value.
// This method is only used internally within the SDK's EventStream handling.
func (s *TranscriptEvent) UnmarshalEvent(
	payloadUnmarshaler protocol.PayloadUnmarshaler,
	msg eventstream.Message,
) error {
	if err := payloadUnmarshaler.UnmarshalPayload(
		bytes.NewReader(msg.Payload), s,
	); err != nil {
		return err
	}
	return nil
}

func (s *TranscriptEvent) MarshalEvent(pm protocol.PayloadMarshaler) (msg eventstream.Message, err error) {
	msg.Headers.Set(eventstreamapi.MessageTypeHeader, eventstream.StringValue(eventstreamapi.EventMessageType))
	var buf bytes.Buffer
	if err = pm.MarshalPayload(&buf, s); err != nil {
		return eventstream.Message{}, err
	}
	msg.Payload = buf.Bytes()
	return msg, err
}

// TranscriptResultStreamEvent groups together all EventStream
// events writes for TranscriptResultStream.
//
// These events are:
//
//     * TranscriptEvent
type TranscriptResultStreamEvent interface {
	eventTranscriptResultStream()
	eventstreamapi.Marshaler
	eventstreamapi.Unmarshaler
}

// TranscriptResultStreamReader provides the interface for reading to the stream. The
// default implementation for this interface will be TranscriptResultStream.
//
// The reader's Close method must allow multiple concurrent calls.
//
// These events are:
//
//     * TranscriptEvent
type TranscriptResultStreamReader interface {
	// Returns a channel of events as they are read from the event stream.
	Events() <-chan TranscriptResultStreamEvent

	// Close will stop the reader reading events from the stream.
	Close() error

	// Returns any error that has occurred while reading from the event stream.
	Err() error
}

type readTranscriptResultStream struct {
	eventReader *eventstreamapi.EventReader
	stream      chan TranscriptResultStreamEvent
	err         *eventstreamapi.OnceError

	done      chan struct{}
	closeOnce sync.Once
}

func newReadTranscriptResultStream(eventReader *eventstreamapi.EventReader) *readTranscriptResultStream {
	r := &readTranscriptResultStream{
		eventReader: eventReader,
		stream:      make(chan TranscriptResultStreamEvent),
		done:        make(chan struct{}),
		err:         eventstreamapi.NewOnceError(),
	}
	go r.readEventStream()

	return r
}

// Close will close the underlying event stream reader.
func (r *readTranscriptResultStream) Close() error {
	r.closeOnce.Do(r.safeClose)
	return r.Err()
}

func (r *readTranscriptResultStream) ErrorSet() <-chan struct{} {
	return r.err.ErrorSet()
}

func (r *readTranscriptResultStream) safeClose() {
	close(r.done)
}

func (r *readTranscriptResultStream) Err() error {
	return r.err.Err()
}

func (r *readTranscriptResultStream) Events() <-chan TranscriptResultStreamEvent {
	return r.stream
}

func (r *readTranscriptResultStream) readEventStream() {
	defer r.Close()
	defer close(r.stream)

	for {
		event, err := r.eventReader.ReadEvent()
		if err != nil {
			if err == io.EOF {
				return
			}
			select {
			case <-r.done:
				// If closed already ignore the error
				return
			default:
			}
			r.err.SetError(err)
			return
		}

		select {
		case r.stream <- event.(TranscriptResultStreamEvent):
		case <-r.done:
			return
		}
	}
}

func unmarshalerForTranscriptResultStreamEvent(eventType string) (eventstreamapi.Unmarshaler, error) {
	switch eventType {
	case "TranscriptEvent":
		return &TranscriptEvent{}, nil
	case "BadRequestException":
		return &BadRequestException{}, nil
	case "ConflictException":
		return &ConflictException{}, nil
	case "InternalFailureException":
		return &InternalFailureException{}, nil
	case "LimitExceededException":
		return &LimitExceededException{}, nil
	default:
		return nil, awserr.New(
			request.ErrCodeSerialization,
			fmt.Sprintf("unknown event type name, %s, for TranscriptResultStream", eventType),
			nil,
		)
	}
}

const (
	// ItemTypePronunciation is a ItemType enum value
	ItemTypePronunciation = "PRONUNCIATION"

	// ItemTypePunctuation is a ItemType enum value
	ItemTypePunctuation = "PUNCTUATION"
)

const (
	// LanguageCodeEnUs is a LanguageCode enum value
	LanguageCodeEnUs = "en-US"

	// LanguageCodeEnGb is a LanguageCode enum value
	LanguageCodeEnGb = "en-GB"

	// LanguageCodeEsUs is a LanguageCode enum value
	LanguageCodeEsUs = "es-US"

	// LanguageCodeFrCa is a LanguageCode enum value
	LanguageCodeFrCa = "fr-CA"

	// LanguageCodeFrFr is a LanguageCode enum value
	LanguageCodeFrFr = "fr-FR"
)

const (
	// MediaEncodingPcm is a MediaEncoding enum value
	MediaEncodingPcm = "pcm"
)
