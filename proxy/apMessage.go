package proxy

type APChannelInfo struct {
	Method           string  `json:"method"`
	TxHash           string  `json:"txHash"`
	APaddr           Address `json:"-"` // the Address is the 20bytes unique ID, and is computed by the publicKey
	APaddrStr        string  `json:"apAddr"`
	APid             string  `json:"apid"`
	APCoverageRadius float64 `json:"-"`
	APLocation       struct {
		Lat float64 `json:"lat"` // latitude
		Lng float64 `json:"lng"` // longitude
	} `json:"-"`
	APLocationStr         string  `json:"apLocation"`
	APUsingChannelList    []int   `json:"-"`
	APUsingChannelListStr string  `json:"apChannelList"`
	APChannelOccupancy    float64 `json:"-"`
	APPower               float64 `json:"apPower"`
	SpectrumCoinNumber    float64 `json:"-"`
	Square                float64 `json:"-"`
	ApIpinfo              string  `json:"-"`
}

type Seller struct {
	SellerAPaddr    Address `json:"-"` // the address of the seller
	SellerAPaddrStr string  `json:"sellerAPaddr"`
	SellerAPid      string  `json:"sellerAPid"` // the id of the seller
	ChannelPrice    float64 `json:"-"`          // the bid price for a channel
}

type BidingPriceInfo struct {
	Method         string   `json:"method"`
	TxHash         string   `json:"txHash"`
	BuyerAPaddr    Address  `json:"-"` // the address of the buyer AP
	BuyerAPaddrStr string   `json:"buyerAPaddr"`
	BuyerAPid      string   `json:"buyerAPid"` // the id of the buyer AP
	PricesList     []Seller `json:"-"`
	PricesListStr  string   `json:"pricesList"`
}

type ChannelSwitchInfo struct {
	Method          string  `json:"method"`
	TxHash          string  `json:"txHash"`
	BuyerAPaddr     Address `json:"-"`
	BuyerAPaddrStr  string  `json:"buyerAPaddr"`
	BuyerAPid       string  `json:"buyerAPid"` // though is somewhat redundant with the addr, keep it
	SellerAPaddr    Address `json:"-"`
	SellerAPaddrStr string  `json:"sellerAPaddr"`
	SellerAPid      string  `json:"sellerAPid"`     // though is somewhat redundant with the addr, keep it
	SwitchChannels  []int   `json:"switchChannels"` // at this time, the length is always 1 as only one channel is traded at a time.
	TotalPrice      float64 `json:"-"`              // the amount of channel coins should be paid
	blockNum        uint64  `json:"-"`
	TimeStamp       string  `json:"-"` //seller将此ChannelSwitchInfo消息上链的时间戳
}

func (csi *ChannelSwitchInfo) SetBlockNum(num uint64) {
	csi.blockNum = num
}

func (csi *ChannelSwitchInfo) GetBlockNum() uint64 {
	return csi.blockNum
}

type ChannelDealInfo struct {
	Method              string  `json:"method"`
	TxHash              string  `json:"txHash"`
	ReferredTradeTstamp string  `json:"-"` // ++++++++ buyer所选择的ChannelSwitchInfo消息的时间戳
	BuyerAPaddr         Address `json:"-"`
	BuyerAPaddrStr      string  `json:"buyerAPaddr"`
	BuyerAPid           string  `json:"buyerAPid"` // though is somewhat redundant with the addr, keep it
	SellerAPaddr        Address `json:"-"`
	SellerAPaddrStr     string  `json:"sellerAPaddr"`
	SellerAPid          string  `json:"sellerAPid"` // though is somewhat redundant with the addr, keepit
	SwitchChannels      []int   `json:"switchChannels"`
	// the amount of channel coins to send.
	// should equal over the total price in the referred TradeTx,
	// if not, this DealTx cannot achieve agreement through the blockchain consensus.
	SendAmount float64 `json:"-"`
}

type QueryAPChannelInfo struct {
	Method    string `json:"method"`
	TxHash    string `json:"txHash"`
	APaddrStr string `json:"apAddr"`
}
