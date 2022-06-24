package parser

import (
	"fmt"
	"strings"

	"github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/abi"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/common"
)

type ParseMap struct {
	//多个合约订阅事件解析  contractAddr-->eventID--->event
	TopicsContractMap map[string]map[common.Hash]abi.Event
	//接口类型定义解析,合约地址为空的  eventID--->event
	TopicsEventMap map[common.Hash]abi.Event
}

//解析结果
type ParseTxResult struct {
	// contractAddr--->eventID--->paramName--->value
	ParseContractMap map[string]map[common.Hash]map[string]interface{}

	// eventID--->paramName--->value
	ParseEventMap map[common.Hash]map[string]interface{}
}

//多个合约订阅事件解析,全局变量   contractAddr-->eventID--->event
//var TopicsContractsMap = make(map[string]map[common.Hash]abiStr.Event)

func ParseTopics(event abi.Event, topics []string) (map[string]interface{}, error) {
	var hashs []common.Hash
	for _, topic := range topics {
		hashs = append(hashs, common.BytesToHash(common.FromHex(topic)))
	}
	//判断eventID 是否相等,如果不等说明不是该事件
	if len(hashs) == 0 || hashs[0] != event.ID {
		return nil, fmt.Errorf("It's not a listen event!")
	}
	outMap := make(map[string]interface{})
	err := abi.ParseTopicsIntoMap(outMap, event.Inputs, hashs[1:])
	if err != nil {
		return outMap, fmt.Errorf("ParseTopics have a err: %s /n", err.Error())
	}
	return outMap, nil
}

//直接解析evm订阅事件
func ParseEVMTxLogs(blks *types.EVMTxLogsInBlks, parseMap *ParseMap) map[string]*ParseTxResult {
	//txhash--->contractAddr--->eventID--->paramName--->value
	var results = make(map[string]*ParseTxResult)
	for _, blk := range blks.GetLogs4EVMPerBlk() {
		for _, txLog := range blk.GetTxAndLogs() {
			var evmAction types.EVMContractAction4Chain33
			err := types.Decode(txLog.Tx.Payload, &evmAction)
			if nil != err {
				continue
			}
			tx := txLog.Tx
			results[common.Bytes2Hex(tx.Hash())] = &ParseTxResult{
				ParseContractMap: make(map[string]map[common.Hash]map[string]interface{}),
				ParseEventMap:    make(map[common.Hash]map[string]interface{}),
			}
			for _, evmLog := range txLog.GetLogsPerTx().GetLogs() {
				//如果TopicsContractMap中存在该合约
				if topicsEvent, ok := parseMap.TopicsContractMap[evmAction.ContractAddr]; ok {
					//从topicsEvent中匹配相关事件
					if event, ok := topicsEvent[common.BytesToHash(evmLog.GetTopic()[0])]; ok {
						var hashs []common.Hash
						for _, topic := range evmLog.GetTopic() {
							hashs = append(hashs, common.BytesToHash(topic))
						}
						outMap := make(map[string]interface{})
						err := abi.ParseTopicsIntoMap(outMap, event.Inputs, hashs[1:])
						if err != nil {
							continue
						}
						results[common.Bytes2Hex(tx.Hash())].ParseContractMap[evmAction.ContractAddr][event.ID] = outMap
					}
				}
				//如果定义存在订阅事件
				if event, ok := parseMap.TopicsEventMap[common.BytesToHash(evmLog.GetTopic()[0])]; ok {
					var hashs []common.Hash
					for _, topic := range evmLog.GetTopic() {
						hashs = append(hashs, common.BytesToHash(topic))
					}
					outMap := make(map[string]interface{})
					err := abi.ParseTopicsIntoMap(outMap, event.Inputs, hashs[1:])
					if err != nil {
						continue
					}
					results[common.Bytes2Hex(tx.Hash())].ParseEventMap[event.ID] = outMap
				}
			}
		}
	}
	return results
}

// 直接解析block订阅日志,返回数据存储： txhash--->ParseTxResult
func ParseBlockReceipts(reqs *types.BlockSeqs, parseMap *ParseMap) map[string]*ParseTxResult {
	//txhash--->contractAddr--->eventID--->paramName--->value
	var results = make(map[string]*ParseTxResult)
	for _, req := range reqs.GetSeqs() {
		for txIndex, tx := range req.GetDetail().Block.Txs {
			//确认是订阅的交易类型
			if !strings.Contains(string(tx.Execer), "evm") {
				continue
			}
			var evmAction types.EVMContractAction4Chain33
			err := types.Decode(tx.Payload, &evmAction)
			if nil != err {
				continue
			}
			//因为只有交易执行成功时，才会存证log信息，所以需要事先判断
			if types.ExecOk != req.GetDetail().Receipts[txIndex].Ty {
				continue
			}

			results[common.Bytes2Hex(tx.Hash())] = &ParseTxResult{
				ParseContractMap: make(map[string]map[common.Hash]map[string]interface{}),
				ParseEventMap:    make(map[common.Hash]map[string]interface{}),
			}
			for _, log := range req.GetDetail().Receipts[txIndex].Logs {
				//TyLogEVMEventData = 605 这个log类型定义在evm合约内部
				if 605 != log.Ty {
					continue
				}
				var evmLog types.EVMLog
				err := types.Decode(log.Log, &evmLog)
				if nil != err {
					continue
				}

				//如果TopicsContractMap中存在该合约
				if topicsEvent, ok := parseMap.TopicsContractMap[evmAction.ContractAddr]; ok {
					//从topicsEvent中匹配相关事件
					if event, ok := topicsEvent[common.BytesToHash(evmLog.GetTopic()[0])]; ok {
						results[common.Bytes2Hex(tx.Hash())].ParseContractMap[evmAction.ContractAddr] = make(map[common.Hash]map[string]interface{})
						var hashs []common.Hash
						for _, topic := range evmLog.GetTopic() {
							hashs = append(hashs, common.BytesToHash(topic))
						}
						outMap := make(map[string]interface{})
						err := abi.ParseTopicsIntoMap(outMap, event.Inputs, hashs[1:])
						if err != nil {
							continue
						}
						results[common.Bytes2Hex(tx.Hash())].ParseContractMap[evmAction.ContractAddr][event.ID] = outMap
					}
				}
				//如果定义存在订阅事件
				if event, ok := parseMap.TopicsEventMap[common.BytesToHash(evmLog.GetTopic()[0])]; ok {
					var hashs []common.Hash
					for _, topic := range evmLog.GetTopic() {
						hashs = append(hashs, common.BytesToHash(topic))
					}
					outMap := make(map[string]interface{})
					err := abi.ParseTopicsIntoMap(outMap, event.Inputs, hashs[1:])
					if err != nil {
						continue
					}
					results[common.Bytes2Hex(tx.Hash())].ParseEventMap[event.ID] = outMap
				}

			}

		}
	}
	return results
}
