package parser

import (
	"encoding/json"
	"log"

	"github.com/33cn/plugin/plugin/dapp/evm/executor/abi"
)

type Config struct {
	ParseTopics  []ParseTopic `json:"parseTopics"`
	ListenServer ListenServer `json:"listenServer"`
	Topic        Topic        `json:"topic"`
	Chain33Host  string       `json:"chain33Host"`
}

type ParseTopic struct {
	//合约地址，当合约地址为空的时候，表示不关联合约地址，只进行满足事件ID解析
	ContractAddr string  `json:"contractAddr"`
	Abi          abi.ABI `json:"abi"`
	//需要解析的event事件名称
	EventNames []string `json:"eventNames"`
}

//订阅服务配置
type ListenServer struct {
	ListenAddr string `json:"listenAddr"`
}

//订阅配置
type Topic struct {
	//服务名称
	Name string `json:"name"`
	//ListenServer URL
	URL string `json:"url"`
	//订阅类型 0订阅区块,1订阅区块头,2订阅交易回执, 3订阅交易执行结果,4订阅指定evm合约event事件
	Type int32 `json：“type`
	//编码方式：这里有bug,可以是jrpc这样json编码也可以是填入grpc,proto编码方式
	Encode string `json:"encode"`
	//合约地址,evm合约地址
	Contracts []string `json:"contracts"`
	//推送开始序号
	LastSequence int64 `json:"lastSequence"`
	//推送开始高度
	LastHeight int64 `json:"lastHeight"`
	//推送开启区块哈希
	LastBlockHash string `json:"lastBlockHash"`
}

func ParseConfig(data []byte) (*Config, error) {
	var conf Config
	err := json.Unmarshal(data, &conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return &conf, err
}
