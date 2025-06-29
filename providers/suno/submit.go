package suno

import (
	"done-hub/common"
	"done-hub/types"
	"net/http"
)

func (s *SunoProvider) Submit(action string, request *SunoSubmitReq) (data *types.TaskResponse[string], errWithCode *types.OpenAIErrorWithStatusCode) {
	var submitUri string
	switch action {
	case SunoActionMusic:
		submitUri = s.SubmitMusic
	case SunoActionLyrics:
		submitUri = s.SubmitLyrics
	default:
		return nil, common.StringErrorWrapper("unsupported action: "+action, "invalid_request", http.StatusBadRequest)
	}

	fullRequestURL := s.GetFullRequestURL(submitUri, "")
	headers := s.GetRequestHeaders()

	// 创建请求
	req, err := s.Requester.NewRequest(http.MethodPost, fullRequestURL, s.Requester.WithHeader(headers), s.Requester.WithBody(request))

	if err != nil {
		return nil, common.ErrorWrapper(err, "new_request_failed", http.StatusInternalServerError)
	}

	data = &types.TaskResponse[string]{}
	_, errWithCode = s.Requester.SendRequest(req, data, false)

	return data, errWithCode
}
