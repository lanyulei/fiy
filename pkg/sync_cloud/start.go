package sync_cloud

import (
	"encoding/json"
	"fiy/common/log"
	"fmt"
	"time"

	"fiy/pkg/sync_cloud/aliyun"

	"fiy/app/cmdb/models/resource"
	orm "fiy/common/global"

	"github.com/spf13/viper"
)

/*
  @Author : lanyulei
*/

type syncStatus struct {
	ID     int  `json:"id"`
	Status bool `json:"status"`
}

type cloudInfo struct {
	resource.CloudDiscovery
	AccountName   string `json:"account_name"`
	AccountType   string `json:"account_type"`
	AccountStatus bool   `json:"account_status"`
	AccountSecret string `json:"account_secret"`
	AccountKey    string `json:"account_key"`
}

// 执行同步任务
func syncCloud() (err error) {

	var (
		ch                 chan syncStatus
		cloudDiscoveryList []*cloudInfo
	)
	// 查询所有的任务列表
	err = orm.Eloquent.Model(&resource.CloudDiscovery{}).
		Joins("left join cmdb_resource_cloud_account as crca on crca.id = cmdb_resource_cloud_discovery.cloud_account").
		Select("cmdb_resource_cloud_discovery.*, crca.name as account_name, crca.type as account_type, crca.status as account_status, crca.secret as account_secret, crca.key as account_key").
		Where("crca.status = ? and cmdb_resource_cloud_discovery.status = ?", true, true).
		Find(&cloudDiscoveryList).Error
	if err != nil {
		return
	}

	ch = make(chan syncStatus, len(cloudDiscoveryList))
	// 接受云资产同步任务执行结果，并处理
	go func(c <-chan syncStatus) {
		for i := 0; i < len(cloudDiscoveryList); i++ {
			r := <-ch
			// todo 更新同步状态与时间
			fmt.Println(r)
		}
		close(ch)
	}(ch)

	// 开启多个goroutine执行云资源任务同步
	for _, task := range cloudDiscoveryList {
		go func(t *cloudInfo, c chan<- syncStatus) {
			defer func(t1 *cloudInfo) {
				if err := recover(); err != nil {
					c <- syncStatus{
						ID:     t1.Id,
						Status: false,
					}
				}
			}(t)

			var err error

			if t.AccountType == "aliyun" {
				regionList := make([]string, 0)
				err = json.Unmarshal(t.Region, &regionList)

				aliyunClient := aliyun.NewAliyun(t.AccountSecret, t.AccountKey, regionList)
				if t.ResourceType == 1 { // 查询云主机资产
					err = aliyunClient.GetEcsList()
				}
			}

			if err != nil {
				errValue := fmt.Sprintf("同步云资源失败，%v", err)
				log.Error(errValue)
				panic(errValue)
			}
		}(task, ch)
	}

	return
}

// 开始同步数据
func Start() (err error) {
	if viper.GetInt(`settings.sync.cloud`) > 0 {
		t := time.NewTicker(viper.GetDuration(`settings.sync.cloud`) * time.Second)
		defer t.Stop()
		for {
			<-t.C
			err = syncCloud()
			if err != nil {
				log.Fatalf("同步云资产数据失败，%v", err)
				return
			}
			t.Reset(viper.GetDuration(`settings.sync.cloud`) * time.Second)
		}
	}
	return
}
