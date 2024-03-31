package system

import (
	"flag"
	"os"

	"dario.cat/mergo"
	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/config.yaml", "the config file")

var Config SystemConfig

type SystemConfig struct {
	Host          string `json:",optional,default=0.0.0.0"`
	Port          int    `json:",optional,default=6789"`
	SystemName    string `json:",optional,default=Marino"`
	SessionSecret string `json:",optional,default=4047a33a5e9ec91151230c26cfef1959"`
	Secret        string `json:",optional,default=5163c66c78eec9e4"`
	SystemPrompt  string `json:",optional,default="`

	// jwt 配置
	Auth struct {
		AccessSecret string `json:",optional,default=13450cd8841c0f0"`
		AccessExpire int64  `json:",optional,default=25920000"`
	}

	Initialled struct {
		Db        bool `json:",optional,default="`
		Embedding bool `json:",optional,default="`
		ActGpt    bool `json:",optional,default="`
		Baidu     bool `json:",optional,default="`
		OpenAi    bool `json:",optional,default="`
		Mail      bool `json:",optional,default="`
		Redis     bool `json:",optional,default="`
	}

	Mail struct {
		SMTPFrom  string `json:",optional"`
		SMTPToken string `json:",optional"`
	}

	// 是否内容检测
	Moderation struct {
		Api          string `json:",optional,default=http://localhost:8080/wordscheck"`
		CheckContent bool   `json:",optional,default=false"`
	}

	Organization struct {
		Name    string `json:",optional,default=FlyOnTheWay"`
		Contact string `json:",optional"`
		Phone   string `json:",optional"`
	}

	// sql 配置
	Db struct {
		DataSource string `json:",optional,default="`
		Dimension  int    `json:",optional,default=768"`
		// 高召回率，费内存。https://mp.weixin.qq.com/s/ue9_IixxgKnZ6vZl-4OkLQ
		IndexType string `json:",optional,default=hnsw"`
	}

	Redis struct {
		DataSource string `json:",optional,default="`
	}

	// openai 配置
	OpenAi struct {
		Host      string `json:",optional,default=http://api.openai.com"`
		AccessKey string `json:",optional,default="`
		// azure or openai
		Type string `json:",optional,default=openai"`
		// for Azure
		APIVersion string `json:",optional,default=2023-05-15"`
	}

	// act-gpt
	ActGpt struct {
		Host string `json:",optional,default=https://maas.act-gpt.com"`
		//
		AccessKey string `json:",optional,default="`
		// 模型
		Model string `json:",optional,default=act-gpt-001"`
	}

	Embedding struct {
		Host      string `json:",optional,default=https://maas.act-gpt.com"`
		Api       string `json:",optional,default=/v1/embeddings"`
		Model     string `json:",optional,default=act-embed-001"`
		AccessKey string `json:",optional,default="`
	}

	Reranker struct {
		Host      string `json:",optional,default=http://0.0.0.0:8000"`
		Api       string `json:",optional,default=/v1/reranker"`
		Model     string `json:",optional,default=act-reranker-001"`
		AccessKey string `json:",optional,default="`
	}

	Baidu struct {
		ClientId     string `json:",optional,default="`
		ClientSecret string `json:",optional,default="`
	}

	Parser struct {
		Host               string `json:",optional,default=https://parser.act-gpt.com"`
		TextApi            string `json:",optional,default=/v1/html2text"`
		DocumentApi        string `json:",optional,default=/v1/extract"`
		ChunkTokenNum      int    `json:",optional,default=500"`
		MinChunkCharNum    int    `json:",optional,default=400"`
		MaxChunkNum        int    `json:",optional,default=550"`
		MinChunkLenToEmbed int    `json:",optional,default=10"`
		ChunkOverlap       int    `json:",optional,default=150"`
	}
}

func Merge(config SystemConfig) {
	var c SystemConfig
	conf.MustLoad(*configFile, &c)
	mergo.Merge(&config, c)
	Config = c
}

func InitNeedSave(conf SystemConfig) bool {
	if os.Getenv("MERGE_CONFIG") != "true" {
		Config = conf
		return false
	}
	Merge(conf)
	return true
}

func WithDefault() SystemConfig {
	var c SystemConfig
	conf.LoadFromJsonBytes([]byte("{}"), &c)
	return c
}

func InitFormFile() {
	flag.Parse()
	var c SystemConfig
	conf.MustLoad(*configFile, &c)
	Config = c
}
