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
	"math/rand"
	"os"
	"testing"
	"time"
)

var (
	appId   = os.Getenv("byte_appId")
	token   = os.Getenv("byte_token")
	cluster = os.Getenv("byte_cluster")
)

func TestEnv(t *testing.T) {
	//fmt.Println(appId, token, cluster)
}

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

	rand.Seed(time.Now().UnixNano())
	randomInt := rand.Intn(1000000)
	fileName := fmt.Sprintf("%d", time.Now().Unix()) + "_" + fmt.Sprintf("%d", randomInt) + "_voice.mp3"
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

func TestLongTextToVoiceCreate(t *testing.T) {
	var params = make(map[string]any)
	params["text"] = "昏黑 hūn hēi 昏黑的夜色中 昏黑的夜色中，月光如水洒落，静谧而神秘，仿佛将世界染上一抹深邃的诗意。"
	params["format"] = "mp3"
	params["voice_type"] = "BV406_V2_streaming"

	tts, err := NewGoTTS(
		context.TODO(),
		WithAppId(appId),
		WithCluster(cluster),
		WithToken(token),
		//WithEmotion(), // 开启情感预测
	)
	if err != nil {
		log.Fatalf("初始化失败，err:%v", err)
	}
	res, err := tts.LongTextToVoiceCreate(params)
	if err != nil {
		log.Fatalf("创建长文本转语音任务失败失败，err:%v", err)
	}

	fmt.Printf("%v \n", res)
}

func TestLongTextToVoiceId(t *testing.T) {
	tts, err := NewGoTTS(
		context.TODO(),
		WithAppId(appId),
		WithCluster(cluster),
		WithToken(token),
		//WithEmotion(), // 开启情感预测
	)
	if err != nil {
		log.Fatalf("初始化失败，err:%v", err)
	}
	res, err := tts.LongTextToVoiceId("f837abed-a3a1-431c-9967-a025b1bea3aa")
	if err != nil {
		log.Fatalf("创建长文本转语音任务失败失败，err:%v", err)
	}

	fileName := time.Now().Format("2006-01-02-15-04-05") + "_file.mp3"
	outFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outFile.Close()

	internal.DownloadToDisk(res.AudioUrl, outFile)

	fmt.Printf("%v \n", res)
}
