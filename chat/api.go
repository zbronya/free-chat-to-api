package chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/launchdarkly/eventsource"
	"github.com/zbronya/free-chat-to-api/config"
	"github.com/zbronya/free-chat-to-api/httpclient"
	"github.com/zbronya/free-chat-to-api/logger"
	"github.com/zbronya/free-chat-to-api/model"
	"github.com/zbronya/free-chat-to-api/model/request"
	"github.com/zbronya/free-chat-to-api/model/response"
	"github.com/zbronya/free-chat-to-api/proofofwork"
	"github.com/zbronya/free-chat-to-api/utils"
	"net/http"
	"strings"
	"time"
)

func Completions(c *gin.Context) {
	req := &request.ChatRequest{}
	err := c.BindJSON(req)
	if err != nil {
		utils.ErrorResp(c, http.StatusBadRequest, "Invalid parameter", nil)
		return
	}

	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
	client := httpclient.NewReqClient()
	pConfig := proofofwork.GetConfig(ua)
	chatRequirementReq := proofofwork.GetChatRequirementReq(pConfig)
	deviceId := uuid.NewString()

	requirement, err := getChatRequirement(c, client, chatRequirementReq, ua, deviceId)
	if err != nil {
		return
	}

	token := proofofwork.CalcProofToken(pConfig, requirement.Proof.Seed, requirement.Proof.Difficulty)

	doConversation(c, client, req, requirement, ua, deviceId, token)
}

func doConversation(c *gin.Context, client *httpclient.ReqClient, req *request.ChatRequest, requirement *model.ChatRequirementRes, ua string, deviceId string, token string) {
	completionReq := model.ApiReqToChatReq(req)
	url := config.GatewayUrl + "/backend-anon/conversation"

	header := map[string]string{
		"Accept":          "text/event-stream",
		"Accept-Encoding": "gzip, deflate, br",
		"Accept-Language": "en-US,en;q=0.9",
		"Content-Type":    "application/json",
		"Oai-Device-Id":   deviceId,
		"Oai-Language":    "en-US",
		"Openai-Sentinel-Chat-Requirements-Token": requirement.Token,
		"Openai-Sentinel-Proof-Token":             token,
		"Origin":                                  "https://chat.openai.com",
		"Priority":                                "u=1, i",
		"Referer":                                 "https://chat.openai.com/",
		"User-Agent":                              ua,
	}

	j, _ := json.Marshal(completionReq)

	resp, err := client.Post(url, header, j)

	if err != nil {
		utils.ErrorResp(c, http.StatusInternalServerError, "fail to completions", err)
		return
	}

	if resp.StatusCode != 200 {
		utils.ErrorResp(c, http.StatusInternalServerError, "fail to completions", nil)
		return
	}

	if req.Stream {
		conversationStream(c, req, resp)
	} else {
		conversation(c, req, resp)
	}

}

func getChatRequirement(c *gin.Context, client httpclient.HttpClient, req model.ChatRequirementReq, ua string, deviceId string) (*model.ChatRequirementRes, error) {
	url := config.GatewayUrl + "/backend-anon/sentinel/chat-requirements"
	header := map[string]string{
		"Accept":          "*/*",
		"Accept-Language": "en-US,en;q=0.9",
		"Content-Type":    "application/json",
		"Oai-Device-Id":   deviceId,
		"Oai-Language":    "en-US",
		"Origin":          "https://chat.openai.com",
		"Referer":         "https://chat.openai.com/",
		"User-Agent":      ua,
	}

	j, _ := json.Marshal(req)

	resp, err := client.Post(url, header, j)
	if err != nil {
		utils.ErrorResp(c, http.StatusInternalServerError, "fail to get chat requirements", err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		utils.ErrorResp(c, http.StatusInternalServerError, "fail to get chat requirements", nil)
		return nil, errors.New("fail to get chat requirements")
	}

	defer resp.Body.Close()

	var chatRequirementRes model.ChatRequirementRes

	err = json.NewDecoder(resp.Body).Decode(&chatRequirementRes)
	return &chatRequirementRes, err
}

func conversationStream(c *gin.Context, req *request.ChatRequest, resp *http.Response) {
	messageTemp := ""
	decoder := eventsource.NewDecoder(resp.Body)
	defer func(decoder *eventsource.Decoder) {
		_, _ = decoder.Decode()
	}(decoder)
	id := utils.GenerateID(29)
	handlingSigns := false
	for {
		event, err := decoder.Decode()
		if err != nil {
			logger.GetLogger().Error(err.Error())
			utils.ErrorResp(c, http.StatusInternalServerError, "", err)
			break
		}
		name := event.Event()
		data := event.Data()
		if data == "" {
			continue
		}
		if data == "[DONE]" {
			result := &response.Stream{}
			result.ID = id
			result.Created = time.Now().Unix()
			result.Object = "chat.completion.chunk"
			delta := response.StreamDelta{
				Content: "",
			}
			choices := response.StreamChoice{
				Delta:        delta,
				FinishReason: "stop",
			}
			result.Choices = append(result.Choices, choices)
			result.Model = req.Model
			bytes, err := json.Marshal(result)
			if err != nil {
				logger.GetLogger().Error(err.Error())
				continue
			}
			c.SSEvent(name, fmt.Sprint(" ", string(bytes)))
			c.SSEvent(name, " [DONE]")
			break
		}
		chatResp := &model.ChatCompletionResp{}
		err = json.Unmarshal([]byte(data), chatResp)
		if chatResp.Error != nil && !handlingSigns {
			logger.GetLogger().Error(fmt.Sprint(chatResp.Error))
			utils.ErrorResp(c, http.StatusInternalServerError, "", chatResp.Error)
			return
		}
		if err != nil {
			continue
		}

		if chatResp.Message.Author.Role == "assistant" && (chatResp.Message.Status == "in_progress" || handlingSigns) {
			handlingSigns = true
			parts := chatResp.Message.Content.Parts[0]
			content := strings.Replace(parts, messageTemp, "", 1)
			messageTemp = parts
			if content == "" {
				continue
			}
			apiResp := &response.Stream{}
			apiResp.ID = id
			apiResp.Created = time.Now().Unix()
			apiResp.Object = "chat.completion.chunk"
			delta := response.StreamDelta{
				Content: content,
			}
			choices := response.StreamChoice{
				Delta: delta,
			}
			apiResp.Choices = append(apiResp.Choices, choices)
			apiResp.Model = req.Model

			bytes, err := json.Marshal(apiResp)
			if err != nil {
				logger.GetLogger().Error(err.Error())
				continue
			}
			c.SSEvent(name, fmt.Sprint(" ", string(bytes)))
			continue
		}
	}
}

func conversation(c *gin.Context, req *request.ChatRequest, resp *http.Response) {
	content := ""
	decoder := eventsource.NewDecoder(resp.Body)
	defer func(decoder *eventsource.Decoder) {
		_, _ = decoder.Decode()
	}(decoder)

	handlingSigns := false
	for {
		event, err := decoder.Decode()
		if err != nil {
			logger.GetLogger().Error(err.Error())
			utils.ErrorResp(c, http.StatusInternalServerError, "", err)
			break
		}
		data := event.Data()
		if data == "" {
			continue
		}
		if data == "[DONE]" {
			result := &response.ChatResponse{}
			result.ID = utils.GenerateID(29)
			result.Created = time.Now().Unix()
			result.Object = "chat.completion"
			result.Model = req.Model
			usage := response.Usage{
				PromptTokens:     0,
				CompletionTokens: 0,
				TotalTokens:      0,
			}
			result.Usage = usage
			message := response.Message{
				Role:    "assistant",
				Content: content,
			}
			choice := response.Choice{
				Message:      message,
				FinishReason: "stop",
				Index:        0,
			}
			result.Choices = append(result.Choices, choice)
			c.JSON(http.StatusOK, result)
			break
		}

		chatResp := &model.ChatCompletionResp{}
		err = json.Unmarshal([]byte(data), chatResp)
		if chatResp.Error != nil && !handlingSigns {
			logger.GetLogger().Error(fmt.Sprint(chatResp.Error))
			utils.ErrorResp(c, http.StatusInternalServerError, "", chatResp.Error)
			return
		}
		if err != nil {
			continue
		}

		if chatResp.Message.Author.Role == "assistant" && (chatResp.Message.Status == "in_progress" || handlingSigns) {
			handlingSigns = true
			if !strings.Contains(chatResp.Message.Content.Parts[0], content) {
				continue
			}
			content = chatResp.Message.Content.Parts[0]
			if content == "" {
				continue
			}
			continue
		}
	}
}
