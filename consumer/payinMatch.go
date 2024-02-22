package consumer

import (
	"encoding/json"
	"ethCrontask/global"
	"fmt"
	"strconv"
	"time"
)

func RunPayInMatch() {
	config := global.ServerConfig
	for {
		r, err := global.Rdb1.BRPop(0, config.EthBlockList).Result()
		if err != nil {
			time.Sleep(time.Second * 5)
			return
		}
		if len(r) < 2 {
			return
		}
		str := r[1]
		var recordList []scanBlockRecord
		err = json.Unmarshal([]byte(str), &recordList)
		if err != nil {
			fmt.Println("err 转化对象格式失败: ", err)
			continue
		}
		for _, strt := range recordList {
			inPool, _ := global.Rdb2.HExists(config.WalletPools, strt.ToHex).Result()
			if inPool {
				result, err := global.Rdb2.HGet(config.WalletPools, strt.ToHex).Result()
				if err != nil {
					fmt.Println("报错：---->", result)
					return
				}
				ethOrder := ETHOrder{}
				err = json.Unmarshal([]byte(result), &ethOrder)
				if err != nil {
					fmt.Println("报错： json---->", result)
					return
				}

				redisVal := RedisOrder{
					FeeType:        ethOrder.FeeType,
					Fee:            ethOrder.Fee,
					AddrTo:         ethOrder.AddrTo,
					Commission:     ethOrder.Commission,
					OrderNo:        ethOrder.OrderNo,
					CoinTo:         ethOrder.CoinTo,
					BlockchainTo:   ethOrder.BlockchainTo,
					AgreementTo:    ethOrder.AgreementTo,
					AmountTo:       ethOrder.AmountTo,
					CoinFrom:       ethOrder.CoinFrom,
					BlockchainFrom: ethOrder.BlockchainFrom,
					AmountFrom:     ethOrder.AmountFrom,
					AddrFrom:       ethOrder.AddrFrom,
					CreateTime:     ethOrder.CreateTime,
					Hash:           ethOrder.Hash,
				}
				val := round(strt.Coins, strt.Value)
				if ethOrder.State != 3 && ethOrder.AmountFrom == val {
					redisVal.State = 2
					ethOrder.State = 2
					ethOrder.SystemNo = strt.Hash
					ethOrder.Collection = val
					ethOrder.UpdateTime = time.Now().Unix()
					err = global.Db.Model(&ETHOrder{}).Where("order_no = ?", ethOrder.OrderNo).Updates(&ethOrder).Error
					if err != nil {
						fmt.Println("收款,更新状态失败:", err)
						return
					}
				}
				if ethOrder.State != 3 && ethOrder.AmountTo != val {
					redisVal.State = 8
					ethOrder.State = 8
					ethOrder.SystemNo = strt.Hash
					ethOrder.Collection = val
					ethOrder.UpdateTime = time.Now().Unix()
					err = global.Db.Model(&ETHOrder{}).Where("order_no = ?", ethOrder.OrderNo).Updates(&ethOrder).Error
					if err != nil {
						fmt.Println("收款数额不对,更新状态失败:", err)
						return
					}
				}
				jsonData, err := json.Marshal(ethOrder)
				if err != nil {
					fmt.Println("转化json失败,jsonData", err)
					return
				}
				err = global.Rdb2.HSet(config.WalletPools, ethOrder.AddrFrom, jsonData).Err()
				if err != nil {
					fmt.Println("err updating Redis2:", err)
					return
				}
				jsonData, err = json.Marshal(redisVal)
				if err != nil {
					fmt.Println("转化json失败,jsonData", err)
					return
				}
				err = global.Rdb0.HSet(config.OrderState, ethOrder.OrderNo, jsonData).Err()
				if err != nil {
					fmt.Println("Error updating Rdb0:", err)
					return
				}
			}
		}
	}
}
func round(coins, value string) float64 {
	var (
		num    int64
		amount float64
	)
	num, _ = strconv.ParseInt(value, 10, 64)
	switch coins {
	case "USDT":
		amount = float64(num) / 1e6
		break
	case "ETH":
		amount = float64(num) / 1e18
		break
	}
	return amount
}

type scanBlockRecord struct {
	Hash    string `json:"hash"`
	Value   string `json:"value"`
	ToHex   string `json:"toHex"`
	FromHex string `json:"fromHex"`
	Coins   string `json:"coins"`
}

type ETHOrder struct {
	Id             int     `gorm:"column:id;primaryKey;AUTO_INCREMENT;comment:主键" json:"id"`
	FeeType        int     `gorm:"column:fee_type;type:tinyint(1);default:0;comment:状态,1-固定手续费,0-浮动手续费;NOT NULL;"  json:"fee_type"`
	Fee            float64 `gorm:"column:fee;type:DECIMAL(20,6);default:0;comment:费率;NOT NULL;" json:"fee"`
	Commission     float64 `gorm:"column:commission;type:DECIMAL(20,6);comment:手续费;NOT NULL;default:0;" json:"commission"`
	OrderNo        string  `gorm:"column:order_no;type:char(150);comment:订单号;NOT NULL;" json:"order_no"`
	CoinTo         string  `gorm:"column:coin_to;type:char(20);comment:客户币;NOT NULL;" json:"coin_to"`
	BlockchainTo   string  `gorm:"column:blockchain_to;type:char(200);comment:客户区块链;NOT NULL" json:"blockchain_to"`
	AgreementTo    string  `gorm:"column:agreement_to;type:char(200);comment:客户协议;NOT NULL;" json:"agreement_to"`
	AmountTo       float64 `gorm:"column:amount_to;type:DECIMAL(20,18);comment:数量;NOT NULL;default:0;" json:"amount_to"`
	AddrTo         string  `gorm:"column:addr_to;type:char(200);comment:客户地址;NOT NULL;" json:"addr_to"`
	PriceTO        float64 `gorm:"column:price_to;type:DECIMAL(20,6);comment:客户单价;NOT NULL;default:0;" json:"price_to"`
	CoinFrom       string  `gorm:"column:coin_from;type:char(20);comment:平台币;NOT NULL" json:"coin_from"`
	BlockchainFrom string  `gorm:"column:blockchain_from;type:char(200);comment:平台区块链;NOT NULL" json:"blockchain_from"`
	AmountFrom     float64 `gorm:"column:amount_from;type:DECIMAL(20,18);comment:数量;NOT NULL;default:0;" json:"amount_from"`
	AddrFrom       string  `gorm:"column:addr_from;type:char(200);comment:平台钱包临时地址;NOT NULL" json:"addr_from"`
	AgreementFrom  string  `gorm:"column:agreement_from;type:char(200);comment:平台协议;" json:"agreement_from"`
	PriceFrom      float64 `gorm:"column:price_from;type:DECIMAL(20,6);comment:平台单价;NOT NULL;default:0;" json:"price_from"`
	State          int     `gorm:"column:state;type:tinyint(1);comment:订单状态：0下单，1待付款，2待接收,3,已完成,8收款数额不对,9超时;default:0;NOT NULL" json:"state"`
	Collection     float64 `gorm:"column:collection;type:DECIMAL(20,18);comment:实际收数额;default:0;NOT NULL" json:"collection"`
	Payment        float64 `gorm:"column:payment;type:DECIMAL(20,18);comment:实际付款;default:0;NOT NULL" json:"payment"`
	Hash           string  `gorm:"column:hash;type:char(200);comment:哈希地址;NOT NULL"  json:"hash"`
	Public         string  `gorm:"column:public;type:char(200);comment:公钥;NOT NULL" json:"public"`
	Private        string  `gorm:"column:private;type:char(200);comment:密钥;NOT NULL" json:"private"`
	SystemNo       string  `gorm:"column:system_no;type:char(200);comment:区块订单付;" json:"system_no"`
	Timeout        int     `gorm:"column:timeout;type:tinyint(1);default:0;NOT NULL;comment:超时状态：0未超时，1已超时;" json:"timeout"`
	CreateTime     int64   `gorm:"column:create_time;type:int(10);comment:创建时间;default:0;NOT NULL" json:"create_time"`
	UpdateTime     int64   `gorm:"column:update_time;type:int(10);comment:更新时间;default:0;NOT NULL" json:"update_time"`
}
type RedisOrder struct {
	AddrFrom       string  `json:"addr_from"`
	AddrTo         string  `json:"addr_to"`
	AgreementTo    string  `json:"agreement_to"`
	AmountFrom     float64 `json:"amount_from"`
	AmountTo       float64 `json:"amount_to"`
	BlockchainFrom string  `json:"blockchain_from"`
	BlockchainTo   string  `json:"blockchain_to"`
	CoinFrom       string  `json:"coin_from"`
	CoinTo         string  `json:"coin_to"`
	CreateTime     int64   `json:"create_time"`
	Commission     float64 `json:"commission"`
	Fee            float64 `json:"fee"`
	FeeType        int     `json:"fee_type"`
	Hash           string  `json:"hash"`
	OrderNo        string  `json:"order_no"`
	State          int     `json:"state"`
	Timeout        int     `json:"timeout"`
}
