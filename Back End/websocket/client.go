package websocket

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

/*
 * Client
 *  - ID:	Client ID,
 *  - Conn: Reference to websocket connection
 *  - Pool: Reference to pool
 */
type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

/*
 * Message
 *  - Type: 0 if bytes, 1 if string (I think)
 *  - Body: String body containing content of message
 */
type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

/*
 * MessageContent
 *  - Type: 		String containing type of data (eg. textMsg)
 *  - Content:	Content struct
 */
type MessageContent struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type Song struct {
	Name   string
	Album  string
	Artist string
	ID     string
}

/*
 * MessageToClient
 *  - Type: 		The type of response
 *  - Response: Content of the response (not always there)
 */
type MessageToClient struct {
	Type     string    `json:"type"`
	Response *Response `json:"response"`
}

/*
 * Response
 *  - Pool:	Tells the client where it has been moved (likely lobby or game)
 */
type Response struct {
	Pool string `json:"pool,omitempty"`
}

/*
 * Content
 *  - TextMsg: If is of textMsg type,
 */
type Content struct {
	TextMsg string `json:"textMsg,omitempty"`
	Song    string `json:"songSearch,omitempty"`
	SongID  string `json:"songID,omitempty"`
}

type ContentClient struct {
	Client  *Client
	Content *Content
}

/*
 * CreateConversation
 *  - Participants: A string delimited by | containing a list of participants in the conversation
 *  - Name:					The name of the conversation
 */
type CreateConversation struct {
	Participants string `json:"participants"`
	Name         string `json:"name"`
}

/*
 * GetConversation
 *  - ConversationID: The hash id of the conversation
 *  - Offset:					Integer offset of range of messages you're grabbing
 *  - ClientID:				ID of the client making the get request
 */
type GetConversation struct {
	ConversationID string `json:"conversationID"`
	Offset         int    `json:"offset"`
	ClientID       string `json:"clientID"`
}

// Read function
func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c

		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		message := Message{Type: messageType, Body: string(p)}
		fmt.Println(message)

		messageContent := &MessageContent{}

		err = json.Unmarshal(p, &messageContent)

		fmt.Println("Type:", messageContent.Type)
		fmt.Println("Content:", messageContent.Content)

		if messageContent.Type == "searchSong" {
			messageSong := &Content{}
			messageClient := &ContentClient{}
			messageClient.Client = c
			messageClient.Content = messageSong

			err = json.Unmarshal([]byte(messageContent.Content), &messageSong)
			c.Pool.SearchSong <- messageClient
		} else if messageContent.Type == "addSong" {
			messageSong := &Content{}
			messageClient := &ContentClient{}
			messageClient.Client = c
			messageClient.Content = messageSong

			err = json.Unmarshal([]byte(messageContent.Content), &messageSong)
			c.Pool.AddSong <- messageClient
		}
	}
}
