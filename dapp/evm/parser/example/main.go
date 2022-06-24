package main

import (
	"compress/gzip"
	"fmt"
	"github.com/33cn/chain33/rpc/jsonclient"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/33cn/chain33-sdk-go/dapp/evm/parser"
	"github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/abi"
	. "github.com/33cn/plugin/plugin/dapp/evm/executor/vm/common"
)

// 主动关闭服务器
var server *http.Server

//需要解析的parseMap
var parseMap *parser.ParseMap

func main() {
	data, err := ioutil.ReadFile("conf.json")
	if err != nil {
		log.Panic(err)
	}
	cfg, err := parser.ParseConfig(data)
	if err != nil {
		log.Panic(err)
	}
	//初始化并赋值
	parseMap = &parser.ParseMap{
		TopicsContractMap: make(map[string]map[Hash]abi.Event),
		TopicsEventMap:    make(map[Hash]abi.Event),
	}
	for _, parseTopic := range cfg.ParseTopics {
		eventMap := make(map[Hash]abi.Event)
		for _, event := range parseTopic.EventNames {
			parseMap.TopicsEventMap[parseTopic.Abi.Events[event].ID] = parseTopic.Abi.Events[event]
			eventMap[parseTopic.Abi.Events[event].ID] = parseTopic.Abi.Events[event]
		}
		if parseTopic.ContractAddr != "" {
			parseMap.TopicsContractMap[parseTopic.ContractAddr] = eventMap
		}
	}
	bindOrResumePush(cfg)

	// 一个通知退出的chan
	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt)

	mux := http.NewServeMux()
	mux.Handle("/", &Handler{cfg: cfg})
	//mux.HandleFunc("/v1/heath", types.HealthCheck{})

	server = &http.Server{
		Addr:         cfg.ListenServer.ListenAddr,
		WriteTimeout: time.Second * 4,
		Handler:      mux,
	}

	go func() {
		// 接收退出信号
		<-exit
		if err := server.Close(); err != nil {
			log.Fatal("Close server:", err)
		}
	}()

	log.Println("Starting v3 httpserver")
	err = server.ListenAndServe()
	if err != nil {
		// 正常退出
		if err == http.ErrServerClosed {
			log.Fatal("Server closed under request")
		} else {
			log.Fatal("Server closed unexpected", err)
		}
	}
	log.Fatal("Server exited")
}

type Handler struct {
	cfg *parser.Config
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)
	if len(r.Header["Content-Encoding"]) >= 1 && r.Header["Content-Encoding"][0] == "gzip" {
		gr, err := gzip.NewReader(r.Body)
		body, err := ioutil.ReadAll(gr)
		if err != nil {
			log.Fatal("Error while serving JSON request: %v", err)
			return
		}

		err = handlerReq(body, h.cfg)
		if err == nil {
			w.Write([]byte("OK"))
		} else {
			w.Write([]byte(err.Error()))
		}
	}
}

//解析evm订阅
func handlerReq(body []byte, cfg *parser.Config) error {
	//TODO 这里暂时只支持区块订阅类型事件解析
	if cfg.Topic.Type == 0 {
		var reqs types.BlockSeqs
		if cfg.Topic.Encode == "jrpc" {
			err := types.JSONToPB(body, &reqs)
			if err != nil {
				log.Fatal("Decoding JSON body have err: %v", err)
				return err
			}
		} else {
			err := types.Decode(body, &reqs)
			if err != nil {
				log.Fatal("Decoding proto body have err: %v", err)
				return err
			}
		}
		results := parser.ParseBlockReceipts(&reqs, parseMap)
		//TODO 后续处理
		log.Println(results)
		return nil
	}

	if cfg.Topic.Type == 4 {
		var reqs types.EVMTxLogsInBlks
		if cfg.Topic.Encode == "jrpc" {
			err := types.JSONToPB(body, &reqs)
			if err != nil {
				log.Fatal("Decoding JSON body have err: %v", err)
				return err
			}
		} else {
			err := types.Decode(body, &reqs)
			if err != nil {
				log.Fatal("Decoding proto body have err: %v", err)
				return err
			}
		}
		log.Println(reqs)
		results := parser.ParseEVMTxLogs(&reqs, parseMap)
		//TODO 后续处理
		log.Println(results)
		return nil
	}
	return fmt.Errorf("unknown type")
}

func bindOrResumePush(cfg *parser.Config) {
	topic := cfg.Topic
	contract := make(map[string]bool)
	for _, name := range topic.Contracts {
		contract[name] = true
	}

	params := types.PushSubscribeReq{
		Name:          topic.Name,
		URL:           topic.URL,
		Encode:        topic.Encode,
		LastSequence:  topic.LastSequence,
		LastHeight:    topic.LastHeight,
		LastBlockHash: topic.LastBlockHash,
		Type:          topic.Type,
		Contract:      contract,
	}
	var res types.ReplySubscribePush
	ctx := jsonclient.NewRPCCtx(cfg.Chain33Host, "Chain33.AddPushSubscribe", params, &res)
	_, err := ctx.RunResult()
	if err != nil {
		fmt.Println("Failed to AddPushSubscribe to  rpc addr:", cfg.Chain33Host, "ReplySubTxReceipt:", res)
		log.Fatal("bindOrResumePush client failed due to:" + err.Error() + ", cfg.Chain33Host:" + cfg.Chain33Host)
	}
	if !res.IsOk {
		fmt.Println("Failed to AddPushSubscribe to  rpc addr:", cfg.Chain33Host, "ReplySubTxReceipt:", res)
		log.Fatal("bindOrResumePush client failed due to res.Msg:" + res.Msg + ", cfg.Chain33Host:" + cfg.Chain33Host)
	}
	fmt.Println("Succeed to AddPushSubscribe for rpc address:", cfg.Chain33Host, ", contract:", params.Contract)
}
