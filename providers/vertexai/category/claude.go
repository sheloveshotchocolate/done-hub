package category

import (
	"done-hub/common"
	"done-hub/common/requester"
	"done-hub/providers/base"
	"done-hub/providers/claude"
	"done-hub/types"
	"encoding/json"
	"net/http"
)

const AnthropicVersion = "vertex-2023-10-16"

type ClaudeRequest struct {
	*claude.ClaudeRequest
	AnthropicVersion string `json:"anthropic_version"`
}

var claudeMap = map[string]string{
	"claude-3-5-sonnet-20240620": "claude-3-5-sonnet@20240620",
	"claude-3-5-sonnet-20241022": "claude-3-5-sonnet-v2@20241022",
	"claude-3-opus-20240229":     "claude-3-opus@20240229",
	"claude-3-sonnet-20240229":   "claude-3-sonnet@20240229",
	"claude-3-haiku-20240307":    "claude-3-haiku@20240307",
	"claude-3-5-haiku-20241022":  "claude-3-5-haiku@20241022",
	"claude-3-7-sonnet-20250219": "claude-3-7-sonnet@20250219",
	"claude-sonnet-4-20250514":   "claude-sonnet-4@20250514",
	"claude-opus-4-20250514":     "claude-opus-4@20250514",
}

func init() {
	CategoryMap["claude"] = &Category{
		Category:                  "claude",
		ChatComplete:              ConvertClaudeFromChatOpenai,
		ResponseChatComplete:      ConvertClaudeToChatOpenai,
		ResponseChatCompleteStrem: ClaudeChatCompleteStrem,
		ErrorHandler:              claude.RequestErrorHandle,
		GetModelName:              GetClaudeModelName,
		GetOtherUrl:               getClaudeOtherUrl,
	}
}

func ConvertClaudeFromChatOpenai(request *types.ChatCompletionRequest) (any, *types.OpenAIErrorWithStatusCode) {
	rawRequest, err := claude.ConvertFromChatOpenai(request)
	if err != nil {
		return nil, err
	}

	claudeRequest := &ClaudeRequest{}
	claudeRequest.ClaudeRequest = rawRequest
	claudeRequest.AnthropicVersion = AnthropicVersion

	// 删除model字段
	claudeRequest.Model = ""

	return claudeRequest, nil
}

func ConvertClaudeToChatOpenai(provider base.ProviderInterface, response *http.Response, request *types.ChatCompletionRequest) (*types.ChatCompletionResponse, *types.OpenAIErrorWithStatusCode) {
	claudeResponse := &claude.ClaudeResponse{}
	err := json.NewDecoder(response.Body).Decode(claudeResponse)
	if err != nil {
		return nil, common.ErrorWrapper(err, "decode_response_failed", http.StatusInternalServerError)
	}

	return claude.ConvertToChatOpenai(provider, claudeResponse, request)
}

func ClaudeChatCompleteStrem(provider base.ProviderInterface, request *types.ChatCompletionRequest) requester.HandlerPrefix[string] {
	chatHandler := &claude.ClaudeStreamHandler{
		Usage:   provider.GetUsage(),
		Request: request,
		Prefix:  `data: {"type"`,
	}

	return chatHandler.HandlerStream
}

func GetClaudeModelName(modelName string) string {
	if value, exists := claudeMap[modelName]; exists {
		modelName = value
	}

	return modelName
}

func getClaudeOtherUrl(stream bool) string {
	if stream {
		return "streamRawPredict"
	}
	return "rawPredict"
}
