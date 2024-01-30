package crons

import (
	"encoding/json"
	"ethCrontask/consumer"
	"ethCrontask/global"
	"fmt"
	"time"
)

func Init() {
	for {
		result, err := global.Rdb2.HGetAll(global.ServerConfig.WalletPools).Result()
		if err != nil {
			fmt.Println("定时查询订单状态")
			return
		}
		for _, s2 := range result {
			btcOrder := consumer.EthOrder{}
			err := json.Unmarshal([]byte(s2), &btcOrder)
			if err != nil {
				return
			}
			if btcOrder.State == 0 && (btcOrder.CreateTime+(10*60)) < time.Now().Unix() {
				btcOrder.State = 9
				btcOrder.UpdateTime = time.Now().Unix()
				err := global.Db.Model(&consumer.EthOrder{}).Where("order_no = ?", btcOrder.OrderNo).Updates(&btcOrder).Error
				if err != nil {
					fmt.Println("mysql更新失败")
					return
				}
				redisVal := consumer.RedisOrder{
					FeeType:      btcOrder.FeeType,
					Fee:          btcOrder.Fee,
					ExchangeRate: btcOrder.ExchangeRate,
					OrderNo:      btcOrder.OrderNo,
					//--------------
					CoinTo:       btcOrder.CoinTo,
					BlockchainTo: btcOrder.BlockchainTo,
					AgreementTo:  btcOrder.AgreementTo,
					AmountTo:     btcOrder.AmountTo,
					//---------------
					CoinFrom:       btcOrder.CoinFrom,
					BlockchainFrom: btcOrder.BlockchainFrom,
					AmountFrom:     btcOrder.AmountFrom,
					AddrFrom:       btcOrder.AddrFrom,
					//订单状态
					State:      btcOrder.State,
					CreateTime: btcOrder.CreateTime,
					Hash:       btcOrder.Hash,
				}
				jsonData, err := json.Marshal(redisVal)
				if err != nil {
					fmt.Println("转化json失败,jsonData", err)
					return
				}
				err = global.Rdb0.HSet(global.ServerConfig.OrderState, btcOrder.OrderNo, jsonData).Err()
				if err != nil {
					fmt.Println("redis0 trx 更新失败", err)
					return
				}
				err = global.Rdb2.HDel(global.ServerConfig.WalletPools, btcOrder.Hash).Err()
				if err != nil {
					fmt.Println("redis3 trx 更新失败", err)
					return
				}
			}

			if btcOrder.State == 7 || btcOrder.State == 8 {
				err = global.Rdb2.HDel(global.ServerConfig.WalletPools, btcOrder.Hash).Err()
				if err != nil {
					fmt.Println("redis3 trx 更新失败", err)
					return
				}
			}

		}
		time.Sleep(time.Second * 10)
	}
}
