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
			ethOrder := consumer.ETHOrder{}
			err := json.Unmarshal([]byte(s2), &ethOrder)
			if err != nil {
				return
			}
			if ethOrder.State != 3 && ethOrder.FeeType == 0 && (ethOrder.CreateTime+(10*60)) < time.Now().Unix() {
				ethOrder.Timeout = 1
				ethOrder.UpdateTime = time.Now().Unix()
				err := global.Db.Model(&consumer.ETHOrder{}).Where("order_no = ?", ethOrder.OrderNo).Updates(&ethOrder).Error
				if err != nil {
					fmt.Println("mysql更新失败")
					return
				}
				redisVal := consumer.RedisOrder{
					FeeType:    ethOrder.FeeType,
					Fee:        ethOrder.Fee,
					AddrTo:     ethOrder.AddrTo,
					Commission: ethOrder.Commission,
					OrderNo:    ethOrder.OrderNo,
					//--------------
					CoinTo:       ethOrder.CoinTo,
					BlockchainTo: ethOrder.BlockchainTo,
					AgreementTo:  ethOrder.AgreementTo,
					AmountTo:     ethOrder.AmountTo,
					//---------------
					CoinFrom:       ethOrder.CoinFrom,
					BlockchainFrom: ethOrder.BlockchainFrom,
					AmountFrom:     ethOrder.AmountFrom,
					AddrFrom:       ethOrder.AddrFrom,
					//订单状态
					State:      ethOrder.State,
					CreateTime: ethOrder.CreateTime,
					Hash:       ethOrder.Hash,
					Timeout:    1,
				}
				jsonData, err := json.Marshal(redisVal)
				if err != nil {
					fmt.Println("转化json失败,jsonData", err)
					return
				}
				err = global.Rdb0.HSet(global.ServerConfig.OrderState, ethOrder.OrderNo, jsonData).Err()
				if err != nil {
					fmt.Println("redis0 trx 更新失败", err)
					return
				}
				//err = global.Rdb2.HDel(global.ServerConfig.WalletPools, btcOrder.Hash).Err()
				//if err != nil {
				//	fmt.Println("redis3 trx 更新失败", err)
				//	return
				//}
			}

			//if btcOrder.State == 9 {
			//	err = global.Rdb2.HDel(global.ServerConfig.WalletPools, btcOrder.Hash).Err()
			//	if err != nil {
			//		fmt.Println("redis3 trx 更新失败", err)
			//		return
			//	}
			//}

		}
		time.Sleep(time.Second * 3)
	}
}
