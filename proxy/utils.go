package proxy

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	HASH_LENGTH    = 256 / 8
	ADDRESS_LENGTH = 160 / 8
)

// 交易阶段
const (
	Stage_InitAPInfo = iota
	Stage_BidingPrice
	Stage_ChannelSwitch
	Stage_ChannelDeal
	Stage_Finish
)

// 交易角色
const (
	Role_Buyer  = "Buyer"
	Role_Seller = "Seller"
	Role_Keeper = "Keeper"
)

const (
	BidingDistanceLimit    = 3
	BalanceFactor          = 1e8
	ChannelOccupancyFactor = 100
)

const (
	CommentDelimiter = "--------------------------------------------------------------"
)

type Hash [HASH_LENGTH]byte

type Address [ADDRESS_LENGTH]byte

var NULLADDR = Address{}

var (
	ROLE_IS_ERROR       = errors.New("The current role should not enter this stage")
	SPECTRUM_NOT_ENOUGH = errors.New("Spectrum resources is not enough")
	COIN_NOT_ENOUGH     = errors.New("Spectrum coins is not enough")
	STAGE_WAIT_TIMEOUT  = errors.New("Current stage waiting time timeout")
	IS_KEEPER           = errors.New("The current AP node is the keeper node")
)

func ParseAddrFromString(str string) Address {
	addrByte, _ := hex.DecodeString(str)
	var add Address
	copy(add[:], addrByte[:ADDRESS_LENGTH])
	return add
}

func FromStringToAddress(str string) Address {
	temp := make([]byte, ADDRESS_LENGTH)
	copy(temp, []byte(str))
	var add Address
	copy(add[:], temp[:ADDRESS_LENGTH])
	return add
}

func ParseHashFromString(str string) Hash {
	hashByte, _ := hex.DecodeString(str)
	var hash Hash
	copy(hash[:], hashByte[:HASH_LENGTH])
	return hash
}

func FromStringToIntSlice(str string) ([]int, error) {
	if !strings.HasPrefix(str, "[") {
		return []int{}, fmt.Errorf("not a list string: %s", str)
	}
	if !strings.HasSuffix(str, "]") {
		return []int{}, fmt.Errorf("not a list string: %s", str)
	}
	res := strings.TrimPrefix(str, "[")
	res = strings.TrimSuffix(res, "]")

	res = strings.TrimLeft(res, "")
	res = strings.TrimLeft(res, " ")

	resList := strings.Split(res, " ")
	intSlice := make([]int, 0)
	for _, ele := range resList {
		if ele == "" || ele == " " {
			continue
		}

		if len(ele) == 0 {
			break
		}
		temp, err := strconv.Atoi(ele)
		if err != nil {
			return []int{}, fmt.Errorf("element cannot convert into int: %v", err)
		}
		intSlice = append(intSlice, temp)
	}
	return intSlice, nil

}

type InitConf struct {
	ContractLinkerURL  string
	ContractLinkerPort string
	DBPath             string
	APNumber           int
	APid               string
}

func ReadConfJson(jsonFile string, result *InitConf) error {
	byteValue, err := os.ReadFile(jsonFile) //读取json文件
	if err != nil {
		fmt.Printf("%s is not exist\n", jsonFile)
		return err
	}
	err = json.Unmarshal(byteValue, &result) //解析json k-v对
	if err != nil {
		return err
	}
	addr := strings.Split(result.ContractLinkerURL, ":")
	result.ContractLinkerPort = addr[1]

	return nil
}

func ReadAPInfoJson(jsonFile string, result *APChannelInfo, APid string) error {
	byteValue, err := os.ReadFile(jsonFile) //读取json文件
	if err != nil {
		return err
	}
	err = json.Unmarshal(byteValue, &result) //解析json k-v对
	if err != nil {
		return err
	}
	result.APaddr = FromStringToAddress(APid)
	//fmt.Println(*result)
	return nil
}

func Address2String(addr Address) string {
	return fmt.Sprintf("%x", addr)
}
