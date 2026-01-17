package centrifugo

import (
	"backend/core/config"
	"backend/core/pkg/errorsx"
	"backend/core/pkg/lifecycle"
	"backend/core/types"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
)

type PublishRequest struct {
	Channel string
	Type    string
	Data    any
}

type Centrifugo struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
	log    *zap.SugaredLogger
	cfg    *config.Config
	client *http.Client
	queue  chan PublishRequest
}

func New(
	ctx context.Context,
	wg *sync.WaitGroup,
	cfg *config.Config,
	log *zap.SugaredLogger,
	bufferSize int,
) *Centrifugo {
	cancelCtx, cancel := lifecycle.Cancel(ctx)

	centrifugo := &Centrifugo{
		ctx:    cancelCtx,
		cancel: cancel,
		wg:     wg,
		log:    log,
		cfg:    cfg,
		client: &http.Client{},
		queue:  make(chan PublishRequest, bufferSize),
	}

	centrifugo.wg.Add(1)
	go centrifugo.worker()

	return centrifugo
}

func (centrifugo *Centrifugo) Publish(channel types.Channel, data any) {
	select {
	case <-centrifugo.ctx.Done():
		centrifugo.log.Infof("Attempted to publish after shutdown: %v", channel)
	case centrifugo.queue <- PublishRequest{
		Channel: centrifugo.getChannel(channel),
		Type:    string(channel.DataType),
		Data:    data,
	}:
	default:
		centrifugo.log.Warn(errorsx.New("Socket publish queue full, dropping message").(*errorsx.Error).ToJSON())
	}
}

func (centrifugo *Centrifugo) Request(ctx context.Context, payload map[string]any) (*http.Request, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		e := errorsx.Wrap(err, "Failed to marshal request payload")
		centrifugo.log.Warn(e.(*errorsx.Error).ToJSON())
		return nil, e
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, centrifugo.apiUrl(), bytes.NewReader(body))
	if err != nil {
		return nil, errorsx.Wrap(err, "Failed to create request")
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "apikey "+centrifugo.cfg.CentrifugoAPIKey)

	return request, nil
}

func (centrifugo *Centrifugo) Send(request *http.Request) (*http.Response, error) {
	response, err := centrifugo.client.Do(request)
	if err != nil {
		return nil, errorsx.Wrap(err, "Failed to send request")
	}

	if response.StatusCode >= http.StatusBadRequest {
		return nil, errorsx.Newf("Centrifugo API error: %s", response.Status)
	}

	return response, nil
}

func (centrifugo *Centrifugo) CloseSend(response *http.Response) {
	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			centrifugo.log.Warn(errorsx.WrapJSON(err, "Failed to close response body"))
		}
	}(response.Body)
}

func (centrifugo *Centrifugo) Close() {
	centrifugo.cancel()
	centrifugo.drainQueue()
	centrifugo.wg.Wait()
}

func (centrifugo *Centrifugo) getChannel(channel types.Channel) string {
	if channel.Personal == nil {
		return string(channel.Name)
	}

	personalID := channel.Personal.Get()
	if personalID == nil {
		return string(channel.Name)
	}

	return string(channel.Name) + "#" + *personalID
}

func (centrifugo *Centrifugo) apiUrl() string {
	return centrifugo.cfg.CentrifugoSchema + "://" +
		centrifugo.cfg.CentrifugoHost + ":" +
		centrifugo.cfg.CentrifugoPort + "/api"
}

func (centrifugo *Centrifugo) worker() {
	defer centrifugo.wg.Done()

	for {
		select {
		case <-centrifugo.ctx.Done():
			centrifugo.drainQueue()
			return
		case req := <-centrifugo.queue:
			centrifugo.publishSync(req)
		}
	}
}

func (centrifugo *Centrifugo) drainQueue() {
	for {
		select {
		case publishRequest := <-centrifugo.queue:
			centrifugo.publishSync(publishRequest)
		default:
			return
		}
	}
}

func (centrifugo *Centrifugo) publishSync(publishRequest PublishRequest) {
	ctx, cancel := lifecycle.Timeout(centrifugo.ctx, 3*time.Second)
	defer cancel()

	payload := map[string]any{
		"method": "publish",
		"params": map[string]any{
			"channel": publishRequest.Channel,
			"data": map[string]any{
				"type": publishRequest.Type,
				"data": publishRequest.Data,
			},
		},
	}

	request, err := centrifugo.Request(ctx, payload)
	if err != nil {
		centrifugo.log.Warn(errorsx.WrapfJSON(err, "Failed to build request for channel %s", publishRequest.Channel))

		return
	}

	response, err := centrifugo.Send(request)
	if err != nil {
		centrifugo.log.Warn(errorsx.WrapfJSON(err, "Failed to send request for channel %s", publishRequest))
		return
	}

	if response.StatusCode < 300 {
		centrifugo.log.Debugf(
			"Message sent to channel %s, type: %v, data: %v",
			publishRequest.Channel,
			publishRequest.Type,
			publishRequest.Data,
		)
	}

	centrifugo.CloseSend(response)
}
