//
// Copyright (c) 2014 MessageBird B.V.
// All rights reserved.
//
// Author: Maurice Nonnekes <maurice@messagebird.com>

// Package messagebird is an official library for interacting with MessageBird.com API.
// The MessageBird API connects your website or application to operators around the world. With our API you can integrate SMS, Chat & Voice.
// More documentation you can find on the MessageBird developers portal: https://developers.messagebird.com/
package messagebird

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"strings"
)

const (
	// ClientVersion is used in User-Agent request header to provide server with API level.
	ClientVersion = "3.0.0"
	// RestEndpoint points to the MessageBird REST API.
	RestEndpoint = "https://rest.messagebird.com"
	// VoiceEndpoint points to the MessageBird Voice API.
	VoiceEndpoint = "https://voice.messagebird.com"
)

const (
	// BalancePath represents the path to the balance resource.
	BalancePath = "balance"
	// HLRPath represents the path to the HLR resource.
	HLRPath = "hlr"
	// MessagePath represents the path to the Message resource.
	MessagePath = "messages"
	// MMSPath represents the path to the MMS resource.
	MMSPath = "mms"
	// VoiceMessagePath represents the path to the VoiceMessage resource.
	VoiceMessagePath = "voicemessages"
	// VerifyPath represents the path to the Verify resource.
	VerifyPath = "verify"
	// LookupPath represents the path to the Lookup resource.
	LookupPath = "lookup"
	// CallFlowPath represents the path to the CallFlow resource.
	CallFlowPath = "call-flows"
	// CallPath represents the path to the Call resource.
	CallPath = "calls"
	// LegPath represents the path to the Leg resource.
	LegPath = "legs"
	// RecordingPath represents the path to the Recording resource.
	RecordingPath = "recordings"
	// TranscriptionPath represents the path to the Transcription resource.
	TranscriptionPath = "transcriptions"
	// WebhookPath represents the path to the Webhook resource.
	WebhookPath = "webhooks"
)

const (
	// Get represents the get resource constant.
	Get = "GET"
	// Post represents the post resource constant.
	Post = "POST"
	// Delete represents the delete resource constant.
	Delete = "Delete"
)

var (
	// ErrResponse is returned when we were able to contact API but request was not successful and contains error details.
	ErrResponse = errors.New("The MessageBird API returned an error")
	// ErrUnexpectedResponse is used when there was an internal server error and nothing can be done at this point.
	ErrUnexpectedResponse = errors.New("The MessageBird API is currently unavailable")
)

// Client is used to access API with a given key.
// Uses standard lib HTTP client internally, so should be reused instead of created as needed and it is safe for concurrent use.
type Client struct {
	AccessKey  string       // The API access key
	HTTPClient *http.Client // The HTTP client to send requests on
	DebugLog   *log.Logger  // Optional logger for debugging purposes
}

// New creates a new MessageBird client object.
func New(AccessKey string) *Client {
	return &Client{AccessKey: AccessKey, HTTPClient: &http.Client{}}
}

func (c *Client) createRequest(method string, endpoint string, path string, params *url.Values) (*http.Request, error) {
	uri, err := url.Parse(endpoint + "/" + path)
	if err != nil {
		return nil, err
	}

	var request *http.Request

	if params != nil {
		body := params.Encode()
		if request, err = http.NewRequest(method, uri.String(), strings.NewReader(body)); err != nil {
			return nil, err
		}

		if c.DebugLog != nil {
			if unescapedBody, queryError := url.QueryUnescape(body); queryError == nil {
				log.Printf("HTTP REQUEST: %s %s %s", method, uri.String(), unescapedBody)
			} else {
				log.Printf("HTTP REQUEST: %s %s %s", method, uri.String(), body)
			}
		}

		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		if request, err = http.NewRequest(method, uri.String(), nil); err != nil {
			return nil, err
		}

		if c.DebugLog != nil {
			log.Printf("HTTP REQUEST: %s %s", method, uri.String())
		}
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", "AccessKey "+c.AccessKey)
	request.Header.Add("User-Agent", "MessageBird/ApiClient/"+ClientVersion+" Go/"+runtime.Version())

	return request, nil
}

func (c *Client) createJSONRequest(method string, endpoint string, path string, params interface{}) (*http.Request, error) {
	uri, err := url.Parse(endpoint + "/" + path)
	if err != nil {
		return nil, err
	}

	jsonEncoded, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, uri.String(), bytes.NewBuffer(jsonEncoded))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", "AccessKey "+c.AccessKey)
	request.Header.Add("User-Agent", "MessageBird/ApiClient/"+ClientVersion+" Go/"+runtime.Version())

	return request, nil
}

func (c *Client) request(v interface{}, request *http.Request) error {
	response, err := c.HTTPClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if c.DebugLog != nil {
		log.Printf("HTTP RESPONSE: %s", string(responseBody))
	}
	fmt.Printf("HTTP RESPONSE: %s\n", string(responseBody))

	// Status code 500 is a server error and means nothing can be done at this
	// point.
	if response.StatusCode == 500 {
		return ErrUnexpectedResponse
	}

	if v != nil {
		if err = json.Unmarshal(responseBody, &v); err != nil {
			return err
		}
	}

	// Status codes 200 and 201 are indicative of being able to convert the
	// response body to the struct that was specified.
	if response.StatusCode == 200 || response.StatusCode == 201 {
		return nil
	}

	// Anything else than a 200/201/500 should be a JSON error.
	return ErrResponse
}

// Balance returns the balance information for the account that is associated
// with the access key.
func (c *Client) Balance() (*Balance, error) {
	request, err := c.createRequest(Get, RestEndpoint, BalancePath, nil)
	if err != nil {
		return nil, err
	}

	balance := &Balance{}
	if err := c.request(balance, request); err != nil {
		if err == ErrResponse {
			return balance, err
		}

		return nil, err
	}

	return balance, nil
}

// HLR looks up an existing HLR object for the specified id that was previously
// created by the NewHLR function.
func (c *Client) HLR(id string) (*HLR, error) {
	request, err := c.createRequest(Get, RestEndpoint, HLRPath+"/"+id, nil)
	if err != nil {
		return nil, err
	}

	hlr := &HLR{}
	if err := c.request(hlr, request); err != nil {
		if err == ErrResponse {
			return hlr, err
		}

		return nil, err
	}

	return hlr, nil
}

// HLRs lists all HLR objects that were previously created by the NewHLR
// function.
func (c *Client) HLRs() (*HLRList, error) {
	request, err := c.createRequest(Get, RestEndpoint, HLRPath, nil)
	if err != nil {
		return nil, err
	}

	hlrList := &HLRList{}
	if err := c.request(hlrList, request); err != nil {
		if err == ErrResponse {
			return hlrList, err
		}

		return nil, err
	}

	return hlrList, nil
}

// NewHLR retrieves the information of an existing HLR.
func (c *Client) NewHLR(msisdn, reference string) (*HLR, error) {
	params := &url.Values{
		"msisdn":    {msisdn},
		"reference": {reference},
	}

	request, err := c.createRequest(Post, RestEndpoint, HLRPath, params)
	if err != nil {
		return nil, err
	}

	hlr := &HLR{}
	if err := c.request(hlr, request); err != nil {
		if err == ErrResponse {
			return hlr, err
		}

		return nil, err
	}

	return hlr, nil
}

// Message retrieves the information of an existing Message.
func (c *Client) Message(id string) (*Message, error) {
	request, err := c.createRequest(Get, RestEndpoint, MessagePath+"/"+id, nil)
	if err != nil {
		return nil, err
	}

	message := &Message{}
	if err := c.request(message, request); err != nil {
		if err == ErrResponse {
			return message, err
		}

		return nil, err
	}

	return message, nil
}

// Messages retrieves all messages of the user represented as a MessageList object.
func (c *Client) Messages() (*MessageList, error) {
	request, err := c.createRequest(Get, RestEndpoint, MessagePath, nil)
	if err != nil {
		return nil, err
	}

	messageList := &MessageList{}
	if err := c.request(messageList, request); err != nil {
		if err == ErrResponse {
			return messageList, err
		}

		return nil, err
	}

	return messageList, nil
}

// NewMessage creates a new message for one or more recipients.
func (c *Client) NewMessage(originator string, recipients []string, body string, msgParams *MessageParams) (*Message, error) {
	params, err := paramsForMessage(msgParams)
	if err != nil {
		return nil, err
	}

	params.Set("originator", originator)
	params.Set("body", body)
	params.Set("recipients", strings.Join(recipients, ","))

	request, err := c.createRequest(Post, RestEndpoint, MessagePath, params)
	if err != nil {
		return nil, err
	}

	message := &Message{}
	if err := c.request(message, request); err != nil {
		if err == ErrResponse {
			return message, err
		}

		return nil, err
	}

	return message, nil
}

// MMSMessage retrieves the information of an existing MmsMessage.
func (c *Client) MMSMessage(id string) (*MMSMessage, error) {
	request, err := c.createRequest(Get, RestEndpoint, MMSPath+"/"+id, nil)
	if err != nil {
		return nil, err
	}

	mmsMessage := &MMSMessage{}
	if err := c.request(mmsMessage, request); err != nil {
		if err == ErrResponse {
			return mmsMessage, err
		}

		return nil, err
	}

	return mmsMessage, nil
}

// NewMMSMessage creates a new MMS message for one or more recipients.
func (c *Client) NewMMSMessage(originator string, recipients []string, msgParams *MMSMessageParams) (*MMSMessage, error) {
	params, err := paramsForMMSMessage(msgParams)
	if err != nil {
		return nil, err
	}

	params.Set("originator", originator)
	params.Set("recipients", strings.Join(recipients, ","))

	request, err := c.createRequest(Post, RestEndpoint, MMSPath, params)
	if err != nil {
		return nil, err
	}

	mmsMessage := &MMSMessage{}
	if err := c.request(mmsMessage, request); err != nil {
		if err == ErrResponse {
			return mmsMessage, err
		}

		return nil, err
	}

	return mmsMessage, nil
}

// VoiceMessage retrieves the information of an existing VoiceMessage.
func (c *Client) VoiceMessage(id string) (*VoiceMessage, error) {
	request, err := c.createRequest(Get, RestEndpoint, VoiceMessagePath+"/"+id, nil)
	if err != nil {
		return nil, err
	}

	message := &VoiceMessage{}
	if err := c.request(message, request); err != nil {
		if err == ErrResponse {
			return message, err
		}

		return nil, err
	}

	return message, nil
}

// VoiceMessages retrieves all VoiceMessages of the user.
func (c *Client) VoiceMessages() (*VoiceMessageList, error) {
	request, err := c.createRequest(Get, RestEndpoint, VoiceMessagePath, nil)
	if err != nil {
		return nil, err
	}

	messageList := &VoiceMessageList{}
	if err := c.request(messageList, request); err != nil {
		if err == ErrResponse {
			return messageList, err
		}

		return nil, err
	}

	return messageList, nil
}

// NewVoiceMessage creates a new voice message for one or more recipients.
func (c *Client) NewVoiceMessage(recipients []string, body string, params *VoiceMessageParams) (*VoiceMessage, error) {
	urlParams := paramsForVoiceMessage(params)
	urlParams.Set("body", body)
	urlParams.Set("recipients", strings.Join(recipients, ","))

	request, err := c.createRequest(Post, RestEndpoint, VoiceMessagePath, urlParams)
	if err != nil {
		return nil, err
	}

	message := &VoiceMessage{}
	if err := c.request(message, request); err != nil {
		if err == ErrResponse {
			return message, err
		}

		return nil, err
	}

	return message, nil
}

// NewVerify generates a new One-Time-Password for one recipient.
func (c *Client) NewVerify(recipient string, params *VerifyParams) (*Verify, error) {
	urlParams := paramsForVerify(params)
	urlParams.Set("recipient", recipient)

	request, err := c.createRequest(Post, RestEndpoint, VerifyPath, urlParams)
	if err != nil {
		return nil, err
	}

	verify := &Verify{}
	if err := c.request(verify, request); err != nil {
		if err == ErrResponse {
			return verify, err
		}

		return nil, err
	}

	return verify, nil
}

// VerifyToken performs token value check against MessageBird API.
func (c *Client) VerifyToken(id, token string) (*Verify, error) {
	params := &url.Values{}
	params.Set("token", token)

	path := VerifyPath + "/" + id + "?" + params.Encode()

	request, err := c.createRequest(Get, RestEndpoint, path, nil)
	if err != nil {
		return nil, err
	}

	verify := &Verify{}
	if err := c.request(verify, request); err != nil {
		if err == ErrResponse {
			return verify, err
		}

		return nil, err
	}

	return verify, nil
}

// Lookup performs a new lookup for the specified number.
func (c *Client) Lookup(phoneNumber string, params *LookupParams) (*Lookup, error) {
	urlParams := paramsForLookup(params)
	path := LookupPath + "/" + phoneNumber + "?" + urlParams.Encode()

	request, err := c.createRequest(Get, RestEndpoint, path, nil)
	if err != nil {
		return nil, err
	}

	lookup := &Lookup{}
	if err := c.request(lookup, request); err != nil {
		if err == ErrResponse {
			return lookup, err
		}

		return nil, err
	}

	return lookup, nil
}

// NewLookupHLR creates a new HLR lookup for the specified number.
func (c *Client) NewLookupHLR(phoneNumber string, params *LookupParams) (*HLR, error) {
	urlParams := paramsForLookup(params)
	path := LookupPath + "/" + phoneNumber + "/hlr"

	request, err := c.createRequest(Post, RestEndpoint, path, urlParams)
	if err != nil {
		return nil, err
	}

	hlr := &HLR{}
	if err := c.request(hlr, request); err != nil {
		if err == ErrResponse {
			return hlr, err
		}

		return nil, err
	}

	return hlr, nil
}

// LookupHLR performs a HLR lookup for the specified number.
func (c *Client) LookupHLR(phoneNumber string, params *LookupParams) (*HLR, error) {
	urlParams := paramsForLookup(params)
	path := LookupPath + "/" + phoneNumber + "/hlr?" + urlParams.Encode()

	request, err := c.createRequest(Get, RestEndpoint, path, nil)
	if err != nil {
		return nil, err
	}

	hlr := &HLR{}
	if err := c.request(hlr, request); err != nil {
		if err == ErrResponse {
			return hlr, err
		}

		return nil, err
	}

	return hlr, nil
}

// CallFlow retrieves the existing CallFlow with the specified id.
func (c *Client) CallFlow(id string) (*CallFlow, error) {
	request, err := c.createRequest(Get, VoiceEndpoint, CallFlowPath+"/"+id, nil)
	if err != nil {
		return nil, err
	}

	callFlowList := &CallFlowList{}
	if err := c.request(callFlowList, request); err != nil {
		if err == ErrResponse {
			return &callFlowList.Data[0], err
		}

		return nil, err
	}
	return &callFlowList.Data[0], nil
}

// CallFlows retrieves all the CallFlows of the user.
func (c *Client) CallFlows() (*CallFlowList, error) {
	request, err := c.createRequest(Get, VoiceEndpoint, CallFlowPath, nil)
	if err != nil {
		return nil, err
	}

	callFlowList := &CallFlowList{}
	if err := c.request(callFlowList, request); err != nil {
		if err == ErrResponse {
			return callFlowList, nil
		}

		return nil, err
	}
	return callFlowList, nil
}

// NewCallFlow creates a new CallFlow.
func (c *Client) NewCallFlow(params *CallFlowParams) (*CallFlow, error) {
	request, err := c.createJSONRequest(Post, VoiceEndpoint, CallFlowPath, params)
	if err != nil {
		return nil, err
	}

	callFlowList := &CallFlowList{}
	if err := c.request(callFlowList, request); err != nil {
		if err == ErrResponse {
			return &callFlowList.Data[0], nil
		}

		return nil, err
	}

	return &callFlowList.Data[0], nil
}

// Call retrieves the existing call with the specified ID.
func (c *Client) Call(id string) (*Call, error) {
	request, err := c.createRequest(Get, VoiceEndpoint, CallPath+"/"+id, nil)
	if err != nil {
		return nil, err
	}

	callList := &CallList{}
	if err := c.request(callList, request); err != nil {
		if err == ErrResponse {
			return &callList.Data[0], nil
		}

		return nil, err
	}

	return &callList.Data[0], nil
}

// Calls retrieves all the Calls of the user.
func (c *Client) Calls() (*CallList, error) {
	request, err := c.createRequest(Get, VoiceEndpoint, CallPath, nil)
	if err != nil {
		return nil, err
	}

	callList := &CallList{}
	if err := c.request(callList, request); err != nil {
		if err == ErrResponse {
			return callList, nil
		}

		return nil, err
	}

	return callList, nil
}

// NewCall creates a new Call resource.
func (c *Client) NewCall(params *CallParams) (*Call, error) {
	request, err := c.createJSONRequest(Post, VoiceEndpoint, CallPath, params)
	if err != nil {
		return nil, err
	}

	callList := &CallList{}
	if err := c.request(callList, request); err != nil {
		if err == ErrResponse {
			return &callList.Data[0], err
		}

		return nil, err
	}

	return &callList.Data[0], err
}

// Leg returns the existing Leg resource with the given legID that belongs to the given Call.
func (c *Client) Leg(callID string, legID string) (*Leg, error) {
	request, err := c.createRequest(Get, VoiceEndpoint, CallPath+"/"+callID+"/"+LegPath+"/"+legID, nil)
	if err != nil {
		return nil, err
	}

	legList := &LegList{}
	if err := c.request(legList, request); err != nil {
		if err == ErrResponse {
			return &legList.Data[0], nil
		}

		return nil, err
	}

	return &legList.Data[0], nil
}

// Legs returns all the Legs belonging to the given Call.
func (c *Client) Legs(callID string) (*LegList, error) {
	request, err := c.createRequest(Get, VoiceEndpoint, CallPath+"/"+callID+"/"+LegPath, nil)
	if err != nil {
		return nil, err
	}

	legList := &LegList{}
	if err := c.request(legList, request); err != nil {
		if err == ErrResponse {
			return legList, nil
		}

		return nil, err
	}

	return legList, nil
}

// Recording returns the existing Recording resource.
func (c *Client) Recording(callID string, legID string, recordingID string) (*Recording, error) {
	var path = CallPath + "/" + callID + "/" + LegPath + "/" + legID + "/" + RecordingPath + "/" + recordingID
	request, err := c.createRequest(Get, VoiceEndpoint, path, nil)
	if err != nil {
		return nil, err
	}

	recordingList := &RecordingList{}
	if err := c.request(recordingList, request); err != nil {
		if err == ErrResponse {
			return &recordingList.Data[0], nil
		}

		return nil, err
	}

	return &recordingList.Data[0], nil
}

// Recordings returns all the recordings of a leg
func (c *Client) Recordings(callID string, legID string) (*RecordingList, error) {
	var path = CallPath + "/" + callID + "/" + LegPath + "/" + legID + "/" + RecordingPath
	request, err := c.createRequest(Get, VoiceEndpoint, path, nil)
	if err != nil {
		return nil, err
	}

	recordingList := &RecordingList{}
	if err := c.request(recordingList, request); err != nil {
		if err == ErrResponse {
			return recordingList, nil
		}

		return nil, err
	}

	return recordingList, nil
}

// Transcription returns the existing Transcription resource.
func (c *Client) Transcription(callID string, legID string, recordingID string, transcriptionID string) (*Transcription, error) {
	var path = CallPath + "/" + callID + "/" + LegPath + "/" + legID + "/" + RecordingPath + "/" + recordingID + "/" + TranscriptionPath + "/" + transcriptionID
	request, err := c.createRequest(Get, VoiceEndpoint, path, nil)
	if err != nil {
		return nil, err
	}

	transcriptionList := &TranscriptionList{}
	if err := c.request(transcriptionList, request); err != nil {
		if err == ErrResponse {
			return &transcriptionList.Data[0], nil
		}

		return nil, err
	}

	return &transcriptionList.Data[0], nil
}

// Transcriptions returns all the Transcriptions of a recording
func (c *Client) Transcriptions(callID string, legID string, recordingID string) (*TranscriptionList, error) {
	var path = CallPath + "/" + callID + "/" + LegPath + "/" + legID + "/" + RecordingPath + "/" + TranscriptionPath
	request, err := c.createRequest(Get, VoiceEndpoint, path, nil)
	if err != nil {
		return nil, err
	}

	transcriptionList := &TranscriptionList{}
	if err := c.request(transcriptionList, request); err != nil {
		if err == ErrResponse {
			return transcriptionList, nil
		}

		return nil, err
	}

	return transcriptionList, nil
}

// NewTranscriptionRequest creates a new Transcription request for the given recording
func (c *Client) NewTranscriptionRequest(callID string, legID string, recordingID string) (*Transcription, error) {
	var path = CallPath + "/" + callID + "/" + LegPath + "/" + legID + "/" + RecordingPath + "/" + TranscriptionPath
	request, err := c.createJSONRequest(Post, VoiceEndpoint, path, nil)
	if err != nil {
		return nil, err
	}

	transcriptionList := &TranscriptionList{}
	if err := c.request(transcriptionList, request); err != nil {
		if err == ErrResponse {
			return &transcriptionList.Data[0], nil
		}

		return nil, err
	}

	return &transcriptionList.Data[0], nil
}

func (c *Client) Webhook(id string) (*Webhook, error) {
	request, err := c.createRequest(Get, VoiceEndpoint, WebhookPath+"/"+id, nil)
	if err != nil {
		return nil, err
	}

	webhookList := &WebhookList{}
	if err = c.request(webhookList, request); err != nil {
		if err == ErrResponse {
			return &webhookList.Data[0], err
		}

		return nil, err
	}

	return &webhookList.Data[0], nil
}

func (c *Client) Webhooks() (*WebhookList, error) {
	request, err := c.createRequest(Get, VoiceEndpoint, WebhookPath, nil)
	if err != nil {
		return nil, err
	}

	webhookList := &WebhookList{}
	if err = c.request(webhookList, request); err != nil {
		if err == ErrResponse {
			return webhookList, err
		}

		return nil, err
	}

	return webhookList, nil
}

func (c *Client) NewWebhook(params *WebhookParams) (*Webhook, error) {
	request, err := c.createJSONRequest(Post, VoiceEndpoint, WebhookPath, params)
	if err != nil {
		return nil, err
	}

	webhookList := &WebhookList{}
	if err = c.request(webhookList, request); err != nil {
		if err == ErrResponse {
			return &webhookList.Data[0], err
		}

		return nil, err
	}

	return &webhookList.Data[0], nil
}

func (c *Client) DeleteWebhook(id string) error {
	request, err := c.createRequest(Delete, VoiceEndpoint, WebhookPath+"/"+id, nil)
	if err != nil {
		return err
	}

	if err = c.request(nil, request); err != nil {
		return err
	}

	return nil
}
