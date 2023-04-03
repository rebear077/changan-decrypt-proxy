package proxy

type ParamList struct {
	IntParam  int
	StrParam  string
	BoolParam bool
}

type JsonCommand struct {
	JsonRpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}
type Signature struct {
	Index     string `json:"index"`
	Signature string `json:"signature"`
}

type TxSignature struct {
}

type Transaction struct {
	BlockHash        string      `json:"blockHash"`
	BlockLimit       string      `json:"blockLimit"`
	BlockNumber      string      `json:"blockNumber"`
	ChainId          string      `json:"chainId"`
	ExtraData        string      `json:"extraData"`
	From             string      `json:"from"`
	Gas              string      `json:"gas"`
	GasPrice         string      `json:"gasPrice"`
	GroupId          string      `json:"groupId"`
	Hash             string      `json:"hash"`
	Input            string      `json:"input"`
	Nonce            string      `json:"nonce"`
	signature        TxSignature `json:"signature"`
	To               string      `json:"to"`
	TransactionIndex string      `json:"transactionIndex"`
	Value            string      `json:"value"`
}

type Result struct {
	DbHash           string        `json:"dbHash"`
	ExtraData        string        `json:"extraData"`
	GasLimit         string        `json:"gasLimit"`
	GasUsed          string        `json:"gasUsed"`
	Hash             string        `json:"hash"`
	LogsBloom        string        `json:"logsBloom"`
	Number           string        `json:"number"`
	ParentHash       string        `json:"parentHash"`
	ReceiptsRoot     string        `json:"receiptsRoot"`
	Sealer           string        `json:"sealer"`
	SealerList       []string      `json:"sealerList"`
	SignatureList    []Signature   `json:"signatureList"`
	StateRoot        string        `json:"stateRoot"`
	Timestamp        string        `json:"timestamp"`
	Transactions     []Transaction `json:"transactions"`
	TransactionsRoot string        `json:"transactionsRoot"`
}

// 区块
type JsonBlockByNumber struct {
	ID      int    `json:"id"`
	JsonRpc string `json:"jsonrpc"`
	Result  Result `json:"result"`
}

// 区块数量(十六进制字符串)
type JsonBlockNumber struct {
	ID      int    `json:"id"`
	JsonRpc string `json:"jsonrpc"`
	Result  string `json:"result"`
}

type TransactionArray struct {
	TxCount      int           `json:"txCount"`
	Transactions []interface{} `json:"transactions"`
}

type BlockArray struct {
	BlockCount int64    `json:"blockCount"`
	Blocks     []Result `json:"blocks"`
}
