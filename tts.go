package go_byte_tts

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zmexing/go-byte-tts/internal"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	apiTts = "https://openspeech.bytedance.com/api/v1/tts"
)

type GoTTS struct {
	ctx     context.Context
	appId   string // 应用标识
	cluster string // 业务集群
	token   string // 应用令牌
}

type Option func(*GoTTS)

func NewGoTTS(ctx context.Context, opts ...Option) (*GoTTS, error) {
	g := &GoTTS{ctx: ctx}
	for _, o := range opts {
		o(g)
	}
	// 参数验证
	if g.appId == "" {
		return nil, errors.New("the parameter appid is defined as")
	}
	if g.cluster == "" {
		return nil, errors.New("the parameter cluster is defined as")
	}
	if g.token == "" {
		return nil, errors.New("the parameter token is defined as")
	}
	return g, nil
}

func WithAppId(speechKey string) Option {
	return func(g *GoTTS) {
		g.appId = speechKey
	}
}

func WithCluster(cluster string) Option {
	return func(g *GoTTS) {
		g.cluster = cluster
	}
}

func WithToken(token string) Option {
	return func(g *GoTTS) {
		g.token = token
	}
}

// TextToVoiceDisk 文本转语音并写入磁盘
func (g *GoTTS) TextToVoiceDisk(params map[string]map[string]any, outFile *os.File) error {
	resp, funcClose, err := g.TextToVoice(params)
	defer funcClose()
	if err != nil {
		return err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	rep := &Rep{}
	err = json.Unmarshal(respBody, rep)
	if err != nil {
		return err
	}

	audio, err := base64.StdEncoding.DecodeString(rep.Data)
	if err != nil {
		return err
	}

	return internal.WriteBytesToDisk(audio, outFile)
}

// TextToVoice 文本转语音
func (g *GoTTS) TextToVoice(params map[string]map[string]any) (*http.Response, func(), error) {
	if params["app"] == nil {
		params["app"] = make(map[string]any)
	}
	params["app"]["appid"] = g.appId
	params["app"]["token"] = "access_token"
	params["app"]["cluster"] = g.cluster

	jsonStr, err := json.Marshal(params)
	if err != nil {
		return nil, func() {}, err
	}

	body := map[string]any{
		"json": string(jsonStr),
	}

	header := map[string]any{
		//"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer;%s", g.token),
	}

	client := internal.NewHTTPClient(
		g.ctx,
		internal.WithTimeout(time.Second*60),
		internal.WithHeader(header),
		internal.WithContentType(internal.HttpJson),
	)
	resp, funcClose, err := client.SendRequest(http.MethodPost, apiTts, body)
	if err != nil {
		return nil, funcClose, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, funcClose, errors.New(resp.Status)
	}

	if resp.ContentLength == 0 {
		return nil, funcClose, errors.New("http response ContentLength=0")
	}

	return resp, funcClose, nil
}
