package consumer

import (
	"encoding/json"
	"ethCrontask/global"
	"fmt"
	"time"
)

func RunPayInMatch() {
	config := global.ServerConfig
	r, err := global.Rdb1.BRPop(0, config.EthBlockList).Result()
	if err != nil {
		time.Sleep(time.Second * 5)
		return
	}
	if len(r) < 2 {
		return
	}
	str := r[1]
	var recordList []ScanBlockRecord
	err = json.Unmarshal([]byte(str), &recordList)

	for _, strt := range recordList {
		inPool, _ := global.Rdb2.HExists(config.WalletPools, strt.ToHex).Result()
		if inPool {
			result, err := global.Rdb2.HGet(config.WalletPools, strt.ToHex).Result()
			if err != nil {
				return
			}
			ethOrder := EthOrder{}
			err = json.Unmarshal([]byte(result), &ethOrder)
			if err != nil {
				return
			}
			redisVal := RedisOrder{
				FeeType:        ethOrder.FeeType,
				Fee:            ethOrder.Fee,
				ExchangeRate:   ethOrder.ExchangeRate,
				OrderNo:        ethOrder.OrderNo,
				CoinTo:         ethOrder.CoinTo,
				BlockchainTo:   ethOrder.BlockchainTo,
				AgreementTo:    ethOrder.AgreementTo,
				AmountTo:       ethOrder.AmountTo,
				CoinFrom:       ethOrder.CoinFrom,
				BlockchainFrom: ethOrder.BlockchainFrom,
				AmountFrom:     ethOrder.AmountFrom,
				AddrFrom:       ethOrder.AddrFrom,
				State:          ethOrder.State,
				CreateTime:     ethOrder.CreateTime,
				Hash:           ethOrder.Hash,
			}
			if ethOrder.State == 0 && ethOrder.AmountFrom == round(strt.Value) {
				redisVal.State = 1
				ethOrder.State = 1
				return
			} else if ethOrder.State == 0 && ethOrder.AmountTo != round(strt.Value) {
				redisVal.State = 8
				ethOrder.State = 8
				ethOrder.Collection = round(strt.Value)
				ethOrder.UpdateTime = time.Now().Unix()
				err = global.Db.Model(&EthOrder{}).Where("order_no = ?", ethOrder.OrderNo).Updates(&ethOrder).Error
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
func round(num int64) float64 {
	return float64(num) / 1e6
}

type ScanBlockRecord struct {
	Hash     string `json:"hash"`
	Value    int64  `json:"value"`
	Gas      int64  `json:"gas"`
	GasPrice int64  `json:"gasPrice"`
	Nonce    int64  `json:"nonce"`
	ToHex    string `json:"toHex"`
	FromHex  string `json:"fromHex"`
}

type EthOrder struct {
	Id             int     `gorm:"column:id;primaryKey;AUTO_INCREMENT;comment:主键" json:"id"`
	FeeType        int     `gorm:"column:fee_type;type:tinyint(1);default:0;comment:状态,1-固定手续费,0-浮动手续费;NOT NULL;"  json:"fee_type"`
	Fee            float64 `gorm:"column:fee;type:DECIMAL(20,6);default:0;comment:费率;NOT NULL;" json:"fee"`
	ExchangeRate   float64 `gorm:"column:exchange_rate;type:DECIMAL(20,6);comment:汇率;NOT NULL;default:0;" json:"exchange_rate"`
	OrderNo        string  `gorm:"column:order_no;type:char(150);comment:订单号;NOT NULL;" json:"order_no"`
	CoinTo         string  `gorm:"column:coin_to;type:char(20);comment:客户币;NOT NULL;" json:"coin_to"`
	BlockchainTo   string  `gorm:"column:blockchain_to;type:char(200);comment:客户区块链;NOT NULL" json:"blockchain_to"`
	AgreementTo    string  `gorm:"column:agreement_to;type:char(200);comment:客户协议;NOT NULL;" json:"agreement_to"`
	AmountTo       float64 `gorm:"column:amount_to;type:DECIMAL(20,6);comment:数量;NOT NULL;default:0;" json:"amount_to"`
	AddrTo         string  `gorm:"column:addr_to;type:char(200);comment:客户地址;NOT NULL;" json:"addr_to"`
	PriceTO        float64 `gorm:"column:price_to;type:DECIMAL(20,6);comment:客户单价;NOT NULL;default:0;" json:"price_to"`
	PriceToTalTo   float64 `gorm:"column:price_total_to;type:DECIMAL(20,6);comment:客户总价;NOT NULL;default:0;" json:"price_total_to"`
	CoinFrom       string  `gorm:"column:coin_from;type:char(20);comment:平台币;NOT NULL" json:"coin_from"`
	BlockchainFrom string  `gorm:"column:blockchain_from;type:char(200);comment:平台区块链;NOT NULL" json:"blockchain_from"`
	AmountFrom     float64 `gorm:"column:amount_from;type:DECIMAL(20,6);comment:数量;NOT NULL;default:0;" json:"amount_from"`
	AddrFrom       string  `gorm:"column:addr_from;type:char(200);comment:平台钱包临时地址;NOT NULL" json:"addr_from"`
	AgreementFrom  string  `gorm:"column:agreement_from;type:char(200);comment:平台协议;" json:"agreement_from"`
	PriceFrom      float64 `gorm:"column:price_from;type:DECIMAL(20,6);comment:平台单价;NOT NULL;default:0;" json:"price_from"`
	PriceTotalFrom float64 `gorm:"column:price_total_from;type:DECIMAL(20,6);comment:平台总价;NOT NULL;default:0;" json:"price_total_from"`
	State          int     `gorm:"column:state;type:tinyint(1);comment:订单状态：用户：0待付款,1待接收,2已完成,8收款数额不对,9超时;default:0;NOT NULL" json:"state"`
	Collection     float64 `gorm:"column:state;type:DECIMAL(20,6);comment:实际收数额;default:0;NOT NULL" json:"collection"`
	Hash           string  `gorm:"column:hash;type:char(200);comment:哈希地址;NOT NULL"  json:"hash"`
	Public         string  `gorm:"column:public;type:char(200);comment:公钥;NOT NULL" json:"public"`
	Private        string  `gorm:"column:private;type:char(200);comment:密钥;NOT NULL" json:"private"`
	CreateTime     int64   `gorm:"column:create_time;type:int(10);comment:创建时间;default:0;NOT NULL" json:"create_time"`
	UpdateTime     int64   `gorm:"column:update_time;type:int(10);comment:更新时间;default:0;NOT NULL" json:"update_time"`
}
type RedisOrder struct {
	AddrFrom       string  `json:"addr_from"`
	AgreementTo    string  `json:"agreement_to"`
	AmountFrom     float64 `json:"amount_from"`
	AmountTo       float64 `json:"amount_to"`
	BlockchainFrom string  `json:"blockchain_from"`
	BlockchainTo   string  `json:"blockchain_to"`
	CoinFrom       string  `json:"coin_from"`
	CoinTo         string  `json:"coin_to"`
	CreateTime     int64   `json:"create_time"`
	ExchangeRate   float64 `json:"exchange_rate"`
	Fee            float64 `json:"fee"`
	FeeType        int     `json:"fee_type"`
	Hash           string  `json:"hash"`
	OrderNo        string  `json:"order_no"`
	State          int     `json:"state"`
}
