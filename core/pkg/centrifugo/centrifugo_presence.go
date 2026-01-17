package centrifugo

import (
	"backend/core/pkg/errorsx"
	"backend/core/pkg/lifecycle"
	"backend/core/types"
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"
)

type PresenceClient struct {
	UserID      string                 `json:"user"`
	ConnInfo    map[string]interface{} `json:"conn_info"`
	ChannelInfo map[string]interface{} `json:"channel_info"`
}

type Presence struct {
	ctx        context.Context
	log        *zap.SugaredLogger
	centrifugo *Centrifugo
}

func NewPresence(ctx context.Context, log *zap.SugaredLogger, centrifugo *Centrifugo) *Presence {
	return &Presence{
		ctx:        ctx,
		log:        log,
		centrifugo: centrifugo,
	}
}

func (s *Presence) GetFrom(channel types.ChannelName) map[string]PresenceClient {
	ctx, cancel := lifecycle.Timeout(s.ctx, 3*time.Second)
	defer cancel()

	payload := map[string]any{
		"method": "presence",
		"params": map[string]any{
			"channel": string(channel),
		},
	}

	request, err := s.centrifugo.Request(ctx, payload)
	if err != nil {
		s.log.Warn(errorsx.WrapfJSON(err, "Failed to create request for channel %s", channel))
		return nil
	}

	response, err := s.centrifugo.Send(request)
	if err != nil {
		s.centrifugo.CloseSend(response)
		s.log.Warn(errorsx.WrapfJSON(err, "Failed to send request for channel %s", channel))
		return nil
	}

	var result struct {
		Result struct {
			Presence map[string]PresenceClient `json:"presence"`
		} `json:"result"`
	}

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		s.centrifugo.CloseSend(response)
		s.log.Warn(errorsx.WrapfJSON(err, "Failed to decode response for channel %s", channel))
		return nil
	}

	s.centrifugo.CloseSend(response)

	return result.Result.Presence
}
