package sync_cloud

import (
	"fiy/app/cmdb/models/resource"
	orm "fiy/common/global"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

/*
  @Author : lanyulei
*/

// 执行同步任务
func syncCloud() (err error) {
	type syncStatus struct {
		ID     int  `json:"id"`
		Status bool `json:"status"`
	}

	var (
		taskList []*resource.CloudDiscovery
		ch       chan syncStatus
	)
	// 查询所有的任务列表
	err = orm.Eloquent.Find(&taskList).Error
	if err != nil {
		return
	}

	ch = make(chan syncStatus, len(taskList))
	// 接受云资产同步任务执行结果，并处理
	go func(c <-chan syncStatus) {
		for i := 0; i < len(taskList); i++ {
			r := <-ch
			fmt.Println(r)
		}
		close(ch)
	}(ch)

	// 开启多个goroutine执行云资源任务同步
	for _, task := range taskList {
		go func(t *resource.CloudDiscovery, c chan<- syncStatus) {
			defer func(t1 *resource.CloudDiscovery) {
				if err := recover(); err != nil {
					c <- syncStatus{
						ID:     t1.Id,
						Status: false,
					}
				}
			}(t)
			fmt.Println(t.Id)
			panic(123)
		}(task, ch)
	}

	return
}

// 开始同步数据
func Start() (err error) {
	for range time.Tick(viper.GetDuration(`settings.sync.cloud`) * time.Second) {
		err = syncCloud()
		if err != nil {
			return
		}
	}
	return
}
