package cohere

import (
	"done-hub/common/requester"
	"done-hub/model"
	"done-hub/providers/base"
	"done-hub/types"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type CohereProviderFactory struct{}

// 创建 CohereProvider
func (f CohereProviderFactory) Create(channel *model.Channel) base.ProviderInterface {
	return &CohereProvider{
		BaseProvider: base.BaseProvider{
			Config:    getConfig(),
			Channel:   channel,
			Requester: requester.NewHTTPRequester(*channel.Proxy, requestErrorHandle),
		},
	}
}

type CohereProvider struct {
	base.BaseProvider
}

func getConfig() base.ProviderConfig {
	return base.ProviderConfig{
		BaseURL:         "https://api.cohere.ai",
		ChatCompletions: "/v2/chat",
		ModelList:       "/v1/models",
		Rerank:          "/v1/rerank",
	}
}

// 请求错误处理
func requestErrorHandle(resp *http.Response) *types.OpenAIError {
	CohereError := &CohereError{}
	err := json.NewDecoder(resp.Body).Decode(CohereError)
	if err != nil {
		return nil
	}

	return errorHandle(CohereError)
}

// 错误处理
func errorHandle(CohereError *CohereError) *types.OpenAIError {
	if CohereError.Message == "" {
		return nil
	}
	return &types.OpenAIError{
		Message: CohereError.Message,
		Type:    "Cohere error",
	}
}

// 获取请求头
func (p *CohereProvider) GetRequestHeaders() (headers map[string]string) {
	headers = make(map[string]string)
	p.CommonRequestHeaders(headers)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", p.Channel.Key)

	return headers
}

func (p *CohereProvider) GetFullRequestURL(requestURL string) string {
	baseURL := strings.TrimSuffix(p.GetBaseURL(), "/")

	return fmt.Sprintf("%s%s", baseURL, requestURL)
}

func convertFinishReason(finishReason string) string {
	switch finishReason {
	case "COMPLETE", "STOP_SEQUENCE":
		return types.FinishReasonStop
	case "MAX_TOKENS":
		return types.FinishReasonLength
	case "TOOL_CALL":
		return types.FinishReasonToolCalls
	case "ERROR":
		return types.FinishReasonContentFilter
	default:
		return types.FinishReasonNull
	}
}
