package category

import (
	"done-hub/common/requester"
	"done-hub/providers/base"
	"done-hub/types"
	"errors"
	"net/http"
	"strings"
)

type Category struct {
	Category                  string
	ChatComplete              ChatCompletionConvert
	ResponseChatComplete      ChatCompletionResponse
	ResponseChatCompleteStrem ChatCompletionStreamResponse
	ErrorHandler              requester.HttpErrorHandler
	GetModelName              func(string) string
	GetOtherUrl               func(bool) string
}

var CategoryMap = map[string]*Category{}

func GetCategory(modelName string) (*Category, error) {

	category := ""

	if strings.HasPrefix(modelName, "gemini") || strings.HasPrefix(modelName, "imagen") {
		category = "gemini"
	} else if strings.HasPrefix(modelName, "claude") {
		category = "claude"
	}

	if category == "" {
		return nil, errors.New("category_not_found")
	}

	return CategoryMap[category], nil

}

type ChatCompletionConvert func(*types.ChatCompletionRequest) (any, *types.OpenAIErrorWithStatusCode)
type ChatCompletionResponse func(base.ProviderInterface, *http.Response, *types.ChatCompletionRequest) (*types.ChatCompletionResponse, *types.OpenAIErrorWithStatusCode)

type ChatCompletionStreamResponse func(base.ProviderInterface, *types.ChatCompletionRequest) requester.HandlerPrefix[string]
