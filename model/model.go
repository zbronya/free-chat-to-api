package model

import (
	"github.com/google/uuid"
	"github.com/zbronya/free-chat-to-api/model/request"
	v1 "github.com/zbronya/free-chat-to-api/utils"
)

type ChatRequirementReq struct {
	P string `json:"p"`
}

type ProofWork struct {
	Difficulty string `json:"difficulty"`
	Required   bool   `json:"required"`
	Seed       string `json:"seed"`
}

type ChatRequirementRes struct {
	Token  string    `json:"token"`
	Proof  ProofWork `json:"proofofwork"`
	Arkose struct {
		Required bool   `json:"required"`
		DX       string `json:"dx"`
	} `json:"arkose"`
}

type ChatAuthor struct {
	Role string `json:"role"`
}

type ChatContent struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}

type ChatMessages struct {
	Author  ChatAuthor  `json:"author"`
	Content ChatContent `json:"content"`
}

type ChatConversationMode struct {
	Kind string `json:"kind"`
}

type ChatCompletionRequest struct {
	Action                     string               `json:"action"`
	Messages                   []ChatMessages       `json:"messages"`
	ParentMessageId            string               `json:"parent_message_id"`
	Model                      string               `json:"model"`
	TimeZoneOffsetMin          int                  `json:"timezone_offset_min"`
	Suggestions                []string             `json:"suggestions"`
	HistoryAndTrainingDisabled bool                 `json:"history_and_training_disabled"`
	ConversationMode           ChatConversationMode `json:"conversation_mode"`
	WebsocketRequestId         string               `json:"websocket_request_id"`
}

func ApiReqToChatReq(req *request.ChatRequest) (chatReq *ChatCompletionRequest) {
	messages := make([]ChatMessages, 0)
	for _, apiMessage := range req.Messages {
		chatMessage := ChatMessages{
			Author: ChatAuthor{
				Role: apiMessage.Role,
			},
			Content: ChatContent{
				ContentType: "text",
				Parts:       []string{apiMessage.Content},
			},
		}
		messages = append(messages, chatMessage)
	}

	chatReq = &ChatCompletionRequest{
		Action:                     "next",
		Messages:                   messages,
		ParentMessageId:            uuid.New().String(),
		Model:                      v1.MappingModel(req.Model),
		TimeZoneOffsetMin:          -180,
		Suggestions:                make([]string, 0),
		HistoryAndTrainingDisabled: true,
		ConversationMode: ChatConversationMode{
			Kind: "primary_assistant",
		},
		WebsocketRequestId: uuid.New().String(),
	}
	return chatReq
}

type ChatCompletionResp struct {
	Message struct {
		Id     string `json:"id"`
		Author struct {
			Role     string      `json:"role"`
			Name     interface{} `json:"name"`
			Metadata struct {
			} `json:"metadata"`
		} `json:"author"`
		CreateTime float64     `json:"create_time"`
		UpdateTime interface{} `json:"update_time"`
		Content    struct {
			ContentType string   `json:"content_type"`
			Parts       []string `json:"parts"`
		} `json:"content"`
		Status   string      `json:"status"`
		EndTurn  interface{} `json:"end_turn"`
		Weight   float64     `json:"weight"`
		Metadata struct {
			Citations        []interface{} `json:"citations"`
			GizmoId          interface{}   `json:"gizmo_id"`
			MessageType      string        `json:"message_type"`
			ModelSlug        string        `json:"model_slug"`
			DefaultModelSlug string        `json:"default_model_slug"`
			Pad              string        `json:"pad"`
			ParentId         string        `json:"parent_id"`
		} `json:"metadata"`
		Recipient string `json:"recipient"`
	} `json:"message"`
	ConversationId string      `json:"conversation_id"`
	Error          interface{} `json:"error"`
}
