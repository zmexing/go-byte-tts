# go-byte-tts

利用字节提供的语音服务，可以通过使用该SDK将文本转化为合成语音并获取受支持的声音列表。

## 快速接入

**安装**
```bash
go get -u github.com/zmexing/go-byte-tts@latest
```

**使用**
短文本转语音、生成mp3文件、直接写入到本地磁盘
```go
package main

import (
	"context"
	byteTts "github.com/zmexing/go-byte-tts"
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"time"
)

const (
	appId   = "xxx"
	token   = "xxx"
	cluster = "xxx"
)

func main() {
	tts, err := byteTts.NewGoTTS(
		context.TODO(),
		byteTts.WithAppId(appId),
		byteTts.WithCluster(cluster),
		byteTts.WithToken(token),
	)
	if err != nil {
		panic("初始化报错: " + fmt.Sprintf("%v", err))
	}

	fileName := time.Now().Format("2006-01-02-15-04-05") + "_file.mp3"
	outFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outFile.Close()

	params := setParams()
	err = tts.TextToVoiceDisk(params, outFile)
	if err != nil {
		log.Fatalf("TextToVoiceDisk失败, err:%v \n", err)
	}
}

func setParams() map[string]map[string]any {
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
	return params
}
```

### 接口
```go
// TextToVoiceDisk 文本转语音并写入磁盘
func (g *GoTTS) TextToVoiceDisk(params map[string]map[string]any, outFile *os.File) error

// TextToVoice 文本转语音
func (g *GoTTS) TextToVoice(params map[string]map[string]any) (*http.Response, func(), error) 
```

## 参考文档
- 字节在线语音合成：https://www.volcengine.com/docs/6561/79820
- 音色列表：https://www.volcengine.com/docs/6561/97465