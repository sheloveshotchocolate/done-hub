package stabilityAI

import (
	"bytes"
	"done-hub/common"
	"done-hub/common/config"
	"done-hub/common/storage"
	"done-hub/common/utils"
	"done-hub/types"
	"encoding/base64"
	"net/http"
	"time"
)

func convertModelName(modelName string) string {
	if modelName == "stable-image-core" {
		return "core"
	}

	return "sd3"
}

func (p *StabilityAIProvider) CreateImageGenerations(request *types.ImageRequest) (*types.ImageResponse, *types.OpenAIErrorWithStatusCode) {
	url, errWithCode := p.GetSupportedAPIUri(config.RelayModeImagesGenerations)
	if errWithCode != nil {
		return nil, errWithCode
	}

	// 获取请求地址
	fullRequestURL := p.GetFullRequestURL(url, convertModelName(request.Model))
	if fullRequestURL == "" {
		return nil, common.ErrorWrapper(nil, "invalid_stabilityAI_config", http.StatusInternalServerError)
	}

	// 获取请求头
	headers := p.GetRequestHeaders()
	headers["Accept"] = "application/json; type=image/png"

	var formBody bytes.Buffer
	builder := p.Requester.CreateFormBuilder(&formBody)
	builder.WriteField("prompt", request.Prompt)
	builder.WriteField("output_format", "png")
	if request.Model != "stable-image-core" {
		builder.WriteField("model", request.Model)
	}
	builder.Close()

	req, err := p.Requester.NewRequest(
		http.MethodPost,
		fullRequestURL,
		p.Requester.WithBody(&formBody),
		p.Requester.WithHeader(headers),
		p.Requester.WithContentType(builder.FormDataContentType()))
	req.ContentLength = int64(formBody.Len())

	if err != nil {
		return nil, common.ErrorWrapper(err, "new_request_failed", http.StatusInternalServerError)
	}

	stabilityAIResponse := &generateResponse{}

	// 发送请求
	_, errWithCode = p.Requester.SendRequest(req, stabilityAIResponse, false)
	if errWithCode != nil {
		return nil, errWithCode
	}

	openaiResponse := &types.ImageResponse{
		Created: time.Now().Unix(),
	}

	imgUrl := ""
	if request.ResponseFormat == "" || request.ResponseFormat == "url" {
		body, err := base64.StdEncoding.DecodeString(stabilityAIResponse.Image)
		if err == nil {
			imgUrl = storage.Upload(body, utils.GetUUID()+".png")
		}
	}

	if imgUrl == "" {
		openaiResponse.Data = []types.ImageResponseDataInner{{B64JSON: stabilityAIResponse.Image}}
	} else {
		openaiResponse.Data = []types.ImageResponseDataInner{{URL: imgUrl}}
	}

	imageCount := len(openaiResponse.Data)
	// PromptTokens保持之前根据prompt内容计算的值
	// CompletionTokens根据生成的图像数量计算，避免空回复计费问题
	p.Usage.CompletionTokens = imageCount * 258
	p.Usage.TotalTokens = p.Usage.PromptTokens + p.Usage.CompletionTokens

	return openaiResponse, nil
}
