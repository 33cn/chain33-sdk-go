package client


// 发送交易
func (client *JSONClient) SendTransaction(signedTx string) (string, error) {
	var res string
	send := &RawParm{
		Token: "BTY",
		Data:  signedTx,
	}
	err := client.Call("Chain33.SendTransaction", send, &res)
	if err != nil {
		return "", err
	}

	return res, nil
}

// 查询交易
func (client *JSONClient) QueryTransaction(hash string) (*TransactionDetail, error) {
	query := QueryParm{
		Hash: hash,
	}
	var detail TransactionDetail
	err := client.Call("Chain33.QueryTransaction", query, &detail)
	if err != nil {
		return nil, err
	}

	return &detail, nil
}
