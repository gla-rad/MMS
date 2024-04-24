package consumer

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/maritimeconnectivity/MMS/mmtp"
	"github.com/maritimeconnectivity/MMS/utils/errors"
	"github.com/maritimeconnectivity/MMS/utils/rw"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"sync"
	"time"
)

type Consumer struct {
	Mrn            string                       // the MRN of the Consumer
	Interests      []string                     // the Interests that the Consumer wants to subscribe to
	Messages       map[string]*mmtp.MmtpMessage // the incoming messages for this Consumer
	MsgMu          *sync.RWMutex                // RWMutex for locking the Messages map
	ReconnectToken string                       // token for reconnecting to a previous session
	Notifications  map[string]*mmtp.MmtpMessage // Map containing pointers to messages, which the Consumer should be notified about
	NotifyMu       *sync.RWMutex                // a Mutex for Notifications map
}

func (c *Consumer) QueueMessage(mmtpMessage *mmtp.MmtpMessage) error {
	if c != nil {
		uUid := mmtpMessage.GetUuid()
		if uUid == "" {
			return fmt.Errorf("the message does not contain a UUID")
		}
		c.MsgMu.Lock()
		c.Messages[uUid] = mmtpMessage
		c.MsgMu.Unlock()
		c.NotifyMu.Lock()
		c.Notifications[uUid] = mmtpMessage
		c.NotifyMu.Unlock()
	} else {
		return fmt.Errorf("consumer resolved to nil while trying to queue message")
	}
	return nil
}

func (c *Consumer) BulkQueueMessages(mmtpMessages []*mmtp.MmtpMessage) {
	if c != nil {
		c.MsgMu.Lock()
		for _, message := range mmtpMessages {
			c.Messages[message.Uuid] = message
		}
		c.MsgMu.Unlock()
	}
}

func (c *Consumer) notify(ctx context.Context, conn *websocket.Conn) error {
	notifications := make([]*mmtp.MessageMetadata, 0, len(c.Notifications))
	for msgUuid, mmtpMsg := range c.Notifications {
		msgMetadata := &mmtp.MessageMetadata{
			Uuid:   mmtpMsg.GetUuid(),
			Header: mmtpMsg.GetProtocolMessage().GetSendMessage().GetApplicationMessage().GetHeader(),
		}
		notifications = append(notifications, msgMetadata)
		delete(c.Notifications, msgUuid)
	}

	notifyMsg := &mmtp.MmtpMessage{
		MsgType: mmtp.MsgType_PROTOCOL_MESSAGE,
		Uuid:    uuid.NewString(),
		Body: &mmtp.MmtpMessage_ProtocolMessage{
			ProtocolMessage: &mmtp.ProtocolMessage{
				ProtocolMsgType: mmtp.ProtocolMessageType_NOTIFY_MESSAGE,
				Body: &mmtp.ProtocolMessage_NotifyMessage{
					NotifyMessage: &mmtp.Notify{
						MessageMetadata: notifications,
					},
				},
			},
		},
	}
	err := rw.WriteMessage(ctx, conn, notifyMsg)
	if err != nil {
		return fmt.Errorf("could not send Notify to Producer: %w", err)
	}
	return nil
}

// CheckNewMessages Checks if there are messages the Agent has not been notified about and notifies about these
func (c *Consumer) CheckNewMessages(ctx context.Context, conn *websocket.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(5 * time.Second):
			c.NotifyMu.Lock()
			if len(c.Notifications) > 0 {
				if err := c.notify(ctx, conn); err != nil {
					log.Println("Failed Notifying Agent:", err)
				}
			}
			c.NotifyMu.Unlock()
			continue
		}
	}
}

// HandleReceive handles request from consumer to receive messages, i.e. lookups buffered messages for the consumer and
// sends these messages to that consumer
func (c *Consumer) HandleReceive(mmtpMessage *mmtp.MmtpMessage, request *http.Request, conn *websocket.Conn) error {
	if receive := mmtpMessage.GetProtocolMessage().GetReceiveMessage(); receive != nil {
		if msgUuids := receive.GetFilter().GetMessageUuids(); msgUuids != nil {
			msgsLen := len(msgUuids)
			mmtpMessages := make([]*mmtp.MmtpMessage, 0, msgsLen)
			appMsgs := make([]*mmtp.ApplicationMessage, 0, msgsLen)
			c.MsgMu.Lock()
			c.NotifyMu.Lock()
			for _, msgUuid := range msgUuids {
				mmtpMsg, exists := c.Messages[msgUuid]
				if exists {
					msg := mmtpMsg.GetProtocolMessage().GetSendMessage().GetApplicationMessage()
					mmtpMessages = append(mmtpMessages, mmtpMsg)
					appMsgs = append(appMsgs, msg)
					delete(c.Messages, msgUuid)
					delete(c.Notifications, msgUuid) //Delete upcoming notification
				}
			}
			c.NotifyMu.Unlock()
			resp := &mmtp.MmtpMessage{
				MsgType: mmtp.MsgType_RESPONSE_MESSAGE,
				Uuid:    uuid.NewString(),
				Body: &mmtp.MmtpMessage_ResponseMessage{ResponseMessage: &mmtp.ResponseMessage{
					ResponseToUuid:      mmtpMessage.GetUuid(),
					Response:            mmtp.ResponseEnum_GOOD,
					ApplicationMessages: appMsgs,
				}},
			}
			err := rw.WriteMessage(request.Context(), conn, resp)
			c.MsgMu.Unlock()
			if err != nil {
				c.BulkQueueMessages(mmtpMessages)
				return fmt.Errorf("could not send messages to Consumer: %w", err)
			}
		} else { // Receive all messages
			c.MsgMu.Lock()
			msgsLen := len(c.Messages)
			appMsgs := make([]*mmtp.ApplicationMessage, 0, msgsLen)
			now := time.Now().UnixMilli()
			c.NotifyMu.Lock()
			for msgUuid, mmtpMsg := range c.Messages {
				msg := mmtpMsg.GetProtocolMessage().GetSendMessage().GetApplicationMessage()
				if now <= msg.Header.Expires {
					appMsgs = append(appMsgs, msg)
					delete(c.Notifications, msgUuid) //Delete upcoming notification
				}
			}
			c.NotifyMu.Unlock()
			resp := &mmtp.MmtpMessage{
				MsgType: mmtp.MsgType_RESPONSE_MESSAGE,
				Uuid:    uuid.NewString(),
				Body: &mmtp.MmtpMessage_ResponseMessage{ResponseMessage: &mmtp.ResponseMessage{
					ResponseToUuid:      mmtpMessage.GetUuid(),
					Response:            mmtp.ResponseEnum_GOOD,
					ApplicationMessages: appMsgs,
				}},
			}

			defer c.MsgMu.Unlock()

			err := rw.WriteMessage(request.Context(), conn, resp)
			if err != nil {
				return fmt.Errorf("could not send messages to Consumer: %w", err)
			} else {
				clear(c.Messages)
			}
		}
	}
	return nil
}

// HandleDisconnect handles a request from a consumer to disconnect, by responding to the consumer and closing the socket
func (c *Consumer) HandleDisconnect(mmtpMessage *mmtp.MmtpMessage, request *http.Request, conn *websocket.Conn) error {
	if disconnect := mmtpMessage.GetProtocolMessage().GetDisconnectMessage(); disconnect != nil {
		resp := &mmtp.MmtpMessage{
			MsgType: mmtp.MsgType_RESPONSE_MESSAGE,
			Uuid:    uuid.NewString(),
			Body: &mmtp.MmtpMessage_ResponseMessage{
				ResponseMessage: &mmtp.ResponseMessage{
					ResponseToUuid: mmtpMessage.GetUuid(),
					Response:       mmtp.ResponseEnum_GOOD,
				}},
		}
		if err := rw.WriteMessage(request.Context(), conn, resp); err != nil {
			return fmt.Errorf("could not send disconnect response to Agent: %w", err)
		}

		if err := conn.Close(websocket.StatusNormalClosure, "Closed connection after receiving Disconnect message"); err != nil {
			return fmt.Errorf("websocket could not be closed cleanly: %w", err)
		}
		return nil
	}
	errors.SendErrorMessage(mmtpMessage.GetUuid(), "Mismatch between protocol message type and message body", request.Context(), conn)
	return fmt.Errorf("message did not contain a Disconnect message in the body")
}
