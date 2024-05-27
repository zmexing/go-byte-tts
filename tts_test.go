package go_byte_tts

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/zmexing/go-byte-tts/internal"
	"io"
	"log"
	"os"
	"testing"
	"time"
)

const (
	appId   = "xxx"
	token   = "xxx"
	cluster = "xxx"
)

func TestTextToVoiceDisk(t *testing.T) {
	tts, err := NewGoTTS(
		context.TODO(),
		WithAppId(appId),
		WithCluster(cluster),
		WithToken(token),
	)
	if err != nil {
		log.Fatalf("初始化失败，err:%v", err)
	}

	// 请求参数定义
	params := make(map[string]map[string]any)
	params["user"] = make(map[string]any)
	//这部分如有需要，可以传递用户真实的ID，方便问题定位
	params["user"]["uid"] = "uid"
	params["audio"] = make(map[string]any)
	//填写选中的音色代号
	params["audio"]["voice_type"] = "BV406_V2_streaming"
	params["audio"]["encoding"] = "mp3"
	params["audio"]["speed_ratio"] = 1.0
	params["audio"]["volume_ratio"] = 1.0
	params["audio"]["pitch_ratio"] = 1.0
	params["request"] = make(map[string]interface{})
	params["request"]["reqid"] = uuid.NewString()
	params["request"]["text"] = "中华兴盛，辛有斌哥。How are you"
	params["request"]["text_type"] = "plain"
	params["request"]["operation"] = "query"

	fileName := time.Now().Format("2006-01-02-15-04-05") + "_file.mp3"
	outFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outFile.Close()

	err = tts.TextToVoiceDisk(params, outFile)
	if err != nil {
		log.Fatalf("生成失败, err:%v \n", err)
	}

}

func TestTextToVoice(t *testing.T) {
	tts, err := NewGoTTS(
		context.TODO(),
		WithAppId(appId),
		WithCluster(cluster),
		WithToken(token),
	)
	if err != nil {
		log.Fatalf("初始化失败，err:%v", err)
	}

	// 请求参数定义
	params := make(map[string]map[string]any)
	params["user"] = make(map[string]any)
	//这部分如有需要，可以传递用户真实的ID，方便问题定位
	params["user"]["uid"] = "uid"
	params["audio"] = make(map[string]any)
	//填写选中的音色代号
	params["audio"]["voice_type"] = "BV406_V2_streaming"
	params["audio"]["encoding"] = "mp3"
	params["audio"]["speed_ratio"] = 1.0
	params["audio"]["volume_ratio"] = 1.0
	params["audio"]["pitch_ratio"] = 1.0
	params["request"] = make(map[string]interface{})
	params["request"]["reqid"] = uuid.NewString()
	params["request"]["text"] = "中华兴盛，辛有斌哥。How are you"
	params["request"]["text_type"] = "plain"
	params["request"]["operation"] = "query"

	resp, funcClose, err := tts.TextToVoice(params)
	if err != nil {
		log.Fatalf("文本转语音失败，err:%v", err)
	}
	defer funcClose()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("读取资源数据失败，err:%v", err)
	}

	rep := &Rep{}
	err = json.Unmarshal(respBody, rep)
	if err != nil {
		log.Fatalf("赋值rep失败，err:%v", err)
	}

	audio, err := base64.StdEncoding.DecodeString(rep.Data)
	if err != nil {
		log.Fatalf("语音转码失败，err:%v", err)
	}

	fileName := time.Now().Format("2006-01-02-15-04-05") + "_file.mp3"
	outFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outFile.Close()

	err = internal.WriteBytesToDisk(audio, outFile)
	if err != nil {
		log.Fatalf("写入磁盘失败，err:%v", err)
	}
}
