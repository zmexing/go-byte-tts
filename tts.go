package go_byte_tts

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jefferyjob/go-easy-utils/anyUtil"
	"github.com/zmexing/go-byte-tts/internal"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	// 短文本语音合成
	apiTts = "https://openspeech.bytedance.com/api/v1/tts"
	// 创建长文本语音
	apiLongTts        = "https://openspeech.bytedance.com/api/v1/tts_async/submit"
	apiLongEmotionTts = "https://openspeech.bytedance.com/api/v1/tts_async_with_emotion/submit"
	// 查询长文本语音合成结果
	apiLongTtsQuery        = "https://openspeech.bytedance.com/api/v1/tts_async/query"
	apiLongEmotionTtsQuery = "https://openspeech.bytedance.com/api/v1/tts_async_with_emotion/query"
	// 长文本语音资源标识
	apiLongResource        = "volc.tts_async.default"
	apiLongEmotionResource = "volc.tts_async.emotion"
)

type GoTTSInter interface {
	// TextToVoice 文本转语音
	TextToVoice(params map[string]map[string]any) (*http.Response, func(), error)

	// TextToVoiceDisk 文本转语音并写入磁盘
	TextToVoiceDisk(params map[string]map[string]any, outFile *os.File) error

	// TextToJoinVoiceDisk 文本转语音并写入磁盘
	// 方法 [TextToVoiceDisk] 因为超过1024字节提示系统错误，所以建议使用 [TextToJoinVoiceDisk]
	// 该方法会自动将文本按照 1024 字节将文本拆开，最后分片生成后合并成一个语音文件
	TextToJoinVoiceDisk(params map[string]map[string]any, outFile *os.File) error

	// LongTextToVoiceCreate 长文本语音合成 任务创建
	// 创建合成任务的频率限制为10 QPS，请勿一次性提交过多任务。
	LongTextToVoiceCreate(params map[string]any) (*TtsAsyncRep, error)

	// LongTextToVoiceId 长文本语音合成 任务查询
	// 音频URL，有效期为1个小时，请及时下载
	LongTextToVoiceId(id string) (*TtsAsyncQueryRep, error)
}

type GoTTS struct {
	ctx     context.Context
	appId   string // 应用标识
	cluster string // 业务集群
	token   string // 应用令牌
	emotion bool   // 是否启用情感预测
}

type Option func(*GoTTS)

func NewGoTTS(ctx context.Context, opts ...Option) (GoTTSInter, error) {
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

func WithEmotion() Option {
	return func(g *GoTTS) {
		g.emotion = true
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
	if err := internal.CheckParams(params); err != nil {
		return nil, func() {}, fmt.Errorf("invalid parameters: %w", err)
	}

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
		return nil, funcClose, fmt.Errorf("http response code failed: %w", errors.New(resp.Status))
	}

	if resp.ContentLength == 0 {
		return nil, funcClose, errors.New("http response ContentLength=0")
	}

	return resp, funcClose, nil
}

func (g *GoTTS) LongTextToVoiceCreate(params map[string]any) (*TtsAsyncRep, error) {
	params["appid"] = g.appId
	params["reqid"] = uuid.NewString()

	// 是否使用情感预测版本
	url := apiLongTts
	resourceId := apiLongResource
	if g.emotion {
		url = apiLongEmotionTts
		resourceId = apiLongEmotionResource
	}

	jsonStr, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("json markshal error: %w", err)
	}

	body := map[string]any{
		"json": string(jsonStr),
	}

	header := map[string]any{
		"Authorization": fmt.Sprintf("Bearer;%s", g.token),
		"Resource-Id":   resourceId,
	}

	client := internal.NewHTTPClient(
		g.ctx,
		internal.WithTimeout(time.Second*60),
		internal.WithHeader(header),
		internal.WithContentType(internal.HttpJson),
	)

	resp, funcClose, err := client.SendRequest(http.MethodPost, url, body)
	defer funcClose()
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http response code failed: %w", errors.New(resp.Status))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("http io ReadAll error: %w", err)
	}

	var ttsAsyncRep TtsAsyncRep
	err = json.Unmarshal(respBody, &ttsAsyncRep)
	if err != nil {
		return nil, fmt.Errorf("http response body Unmarshal error: %w", err)
	}

	return &ttsAsyncRep, nil
}

func (g *GoTTS) LongTextToVoiceId(id string) (*TtsAsyncQueryRep, error) {
	// 是否使用情感预测版本
	url := apiLongTtsQuery
	resourceId := apiLongResource
	if g.emotion {
		url = apiLongEmotionTtsQuery
		resourceId = apiLongEmotionResource
	}

	var params = make(map[string]any)
	params["appid"] = g.appId
	params["task_id"] = id

	header := map[string]any{
		"Authorization": fmt.Sprintf("Bearer;%s", g.token),
		"Resource-Id":   resourceId,
	}

	client := internal.NewHTTPClient(
		g.ctx,
		internal.WithTimeout(time.Second*60),
		internal.WithHeader(header),
	)

	resp, funcClose, err := client.SendRequest(http.MethodGet, url, params)
	defer funcClose()
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http response code failed: %w", errors.New(resp.Status))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("http io ReadAll error: %w", err)
	}

	var result TtsAsyncQueryRep
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, fmt.Errorf("http response body Unmarshal error: %w", err)
	}
	return &result, nil
}

func (g *GoTTS) TextToJoinVoiceDisk(params map[string]map[string]any, outFile *os.File) error {
	text, _ := params["request"]["text"]
	textList := internal.SplitText(anyUtil.AnyToStr(text), 1024)

	// 协程并行处理多个文本列表
	chWork := make(chan ChanJoinVoice, len(textList))
	defer close(chWork)
	for i, v := range textList {
		newMap := internal.DeepCopyParams(params)
		newMap["request"]["text"] = v

		// 如果文本被拆开，则中间的连接停顿应该减小
		if i != (len(textList) - 1) {
			newMap["request"]["silence_duration"] = 50
		}

		go g.workTextToJoinVoiceDisk(newMap, i, chWork)
	}

	// 按照顺序拼接协程的数据
	resMap := make(map[int][]byte, len(textList))
	for i := 0; i < len(textList); i++ {
		wordRes := <-chWork
		if wordRes.Err != nil {
			return wordRes.Err
		}
		resMap[wordRes.Index] = wordRes.Audio
	}

	// 按照顺序拼接结果
	resAudio := []byte{}
	for i, _ := range textList {
		r, ok := resMap[i]
		if !ok {
			return errors.New("error in sequential splicing")
		}
		resAudio = append(resAudio, r...)
	}
	return internal.WriteBytesToDisk(resAudio, outFile)

}

func (g *GoTTS) workTextToJoinVoiceDisk(params map[string]map[string]any, idx int, ch chan ChanJoinVoice) {
	params["request"]["reqid"] = uuid.NewString()

	resp, funcClose, err := g.TextToVoice(params)
	defer funcClose()
	if err != nil {
		ch <- ChanJoinVoice{Index: idx, Err: fmt.Errorf("TextToVoice error: %w", err)}
		return
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- ChanJoinVoice{Index: idx, Err: fmt.Errorf("ReadAll error: %w", err)}
		return
	}

	var rep Rep
	if err := json.Unmarshal(respBody, &rep); err != nil {
		ch <- ChanJoinVoice{Index: idx, Err: fmt.Errorf("JSON unmarshal error: %w", err)}
		return
	}

	audio, err := base64.StdEncoding.DecodeString(rep.Data)
	if err != nil {
		ch <- ChanJoinVoice{Index: idx, Err: fmt.Errorf("base64 decode error: %w", err)}
		return
	}

	ch <- ChanJoinVoice{
		Index: idx,
		Audio: audio,
	}
}
