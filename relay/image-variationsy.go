package relay

import (
	"done-hub/common"
	providersBase "done-hub/providers/base"
	"done-hub/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type relayImageVariations struct {
	relayBase
	request types.ImageEditRequest
}

func NewRelayImageVariations(c *gin.Context) *relayImageVariations {
	relay := &relayImageVariations{}
	relay.c = c
	return relay
}

func (r *relayImageVariations) setRequest() error {
	if err := common.UnmarshalBodyReusable(r.c, &r.request); err != nil {
		return err
	}

	if r.request.Model == "" {
		r.request.Model = "dall-e-2"
	}

	if r.request.Size == "" {
		r.request.Size = "1024x1024"
	}

	r.setOriginalModel(r.request.Model)

	return nil
}

func (r *relayImageVariations) getPromptTokens() (int, error) {
	// 图像变换通常没有prompt文本，返回最小值
	if r.request.Prompt != "" {
		return common.CountTokenText(r.request.Prompt, r.getOriginalModel()), nil
	}
	return 1, nil // 最小计费单位
}

func (r *relayImageVariations) send() (err *types.OpenAIErrorWithStatusCode, done bool) {
	provider, ok := r.provider.(providersBase.ImageVariationsInterface)
	if !ok {
		err = common.StringErrorWrapperLocal("channel not implemented", "channel_error", http.StatusServiceUnavailable)
		done = true
		return
	}

	r.request.Model = r.modelName

	response, err := provider.CreateImageVariations(&r.request)
	if err != nil {
		return
	}
	err = responseJsonClient(r.c, response)

	if err != nil {
		done = true
	}

	return
}
