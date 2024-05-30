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

func TestTextToJoinVoiceDisk(t *testing.T) {
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
	params["audio"]["emotion"] = "通用" // 情感
	params["audio"]["encoding"] = "mp3"
	params["audio"]["speed_ratio"] = 1.0
	params["audio"]["volume_ratio"] = 1.0
	params["audio"]["pitch_ratio"] = 1.0
	params["request"] = make(map[string]interface{})
	params["request"]["reqid"] = uuid.NewString()
	params["request"]["text"] = "昏黑 hūn hēi 昏黑的夜色中 昏黑的夜色中，月光如水洒落，静谧而神秘，仿佛将世界染上一抹深邃的诗意。1. 描绘场景：此句通过“昏黑的夜色”与“月光如水”形成鲜明对比，细腻地勾勒出一个宁静而富有层次的夜晚。月光如同流水般洒落，给静谧的夜增添了动态美，营造出神秘而深邃的氛围。\\n\\n2. 表达情感：句子透露出一种悠远、沉静的情感，仿佛诗人在这月光之下，感受到世界的诗意与宁静，让人体会到作者内心深处对美好事物的向往与陶醉。\\\"昏黑\\\"通常作为形容词使用，意指天色或光线暗淡，看不清楚，常用来形容夜晚或阴暗的环境。这个词由“昏”和“黑”两个同义词组合而成，强调了一种非常昏暗、缺乏光亮的状态。天色黑暗。唐于鹄《过凌霄洞天谒张先生祠》诗：“断崖昼昏黑，槎臬横隻椽。”《初刻拍案惊奇》卷五：“﹝猛虎﹞擒了德容小姐便走……那时夜已昏黑，虽然聚得些人起来，四目相视，束手无策。”清孙枝蔚《为农》诗之三：“归家已昏黑，浊酒妇须谋。”老舍《龙须沟》第二幕：“黎明之前，满院子还是昏黑的。”比喻社会政治黑暗腐败。茅盾《宿莽·色盲》：“我们讲到国际政治的推移，你又说你只见一片昏黑，你成了精神上的色盲。”"
	params["request"]["text_type"] = "plain"
	params["request"]["operation"] = "query"
	//params["request"]["silence_duration"] = 0 // 句尾静音时长，单位为ms，默认为125

	rand.Seed(time.Now().UnixNano())
	randomInt := rand.Intn(1000000)
	fileName := fmt.Sprintf("%d", time.Now().Unix()) + "_" + fmt.Sprintf("%d", randomInt) + "_voice.mp3"
	outFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outFile.Close()

	err = tts.TextToJoinVoiceDisk(params, outFile)
	if err != nil {
		log.Fatalf("生成失败, err:%v \n", err)
	}
}
