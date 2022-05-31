```
#生成seed，注意：生成的seed要记住，如果不小心误删钱包，通过seed可以找回钱包
 ./ycc-cli seed generate -l 0

#保存seed，-p后跟的是钱包密码。 注意：密码可以自定义，并且牢牢记住，后面解锁钱包时会用到。密码需满足8位字母数字组合的要求。
 ./ycc-cli seed save -s [上一步生成的seed值] -p tech1234
    
#解锁钱包
./ycc-cli wallet unlock -p tech1234 -t 0 
   
#导入创始账户地址，创始账户在配置文件中配的
./ycc-cli account import_key -k 3990969DF92A5914F7B71EEB9A4E58D6E255F32BF042FEA5318FC8B3D50EE6E8  -l genesis
    
#新创建一个标签为testA的账户
./ycc-cli account create -l testA
    
#创始账户向账户A中发起转账
./ycc-cli send coins transfer -a 1000 -t 1FbuCnz6Bnw3EiECrv9QYKcnYhP19JL5gb -n "test for transfer bty" -k 3990969DF92A5914F7B71EEB9A4E58D6E255F32BF042FEA5318FC8B3D50EE6E8
    
#查询交易
./ycc-cli tx query_hash -s [上一步生成的交易hash]

#查询最新区块
./ycc-cli block last_header

#查询区块详情
./ycc-cli block view -s [上一步生成的区块hash]

#查询节点是否同步到最新区块
./ycc-cli net is_sync
```