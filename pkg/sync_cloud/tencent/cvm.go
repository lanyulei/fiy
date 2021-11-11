package tencent

import (
	"encoding/json"
	"fmt"

	"fiy/app/cmdb/models/resource"
	orm "fiy/common/global"

	"gorm.io/gorm/clause"

	"fiy/common/log"
	"fiy/tools"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

/*
  @Author : jiangboyang0930
*/

type tencentCloud struct {
	SK     string   `json:"sk"`
	AK     string   `json:"ak"`
	Region []string `json:"region"`
}

func NewTencentCloud(sk, ak string, region []string) *tencentCloud {
	return &tencentCloud{
		SK:     sk,
		AK:     ak,
		Region: region,
	}
}

func (a *tencentCloud) CvmList(infoID int) (err error) {
	var (
		cvmList   []*cvm.Instance
		dataList  []resource.Data
		cvmClient *cvm.Client
		offset    int64 = 0
		limit     int64 = 100
	)

	for _, r := range a.Region {
		cvmClient, err = cvm.NewClient(
			common.NewCredential(a.SK, a.AK),
			tools.Strip(r),
			profile.NewClientProfile(),
		)
		if err != nil {
			log.Errorf("创建客户端连接失败，%v", err)
			return
		}

		request := cvm.NewDescribeInstancesRequest()
		request.Offset = &offset
		request.Limit = &limit

		r, err := cvmClient.DescribeInstances(request)
		if err != nil {
			log.Errorf("查询CVM实例列表失败，%v", err)
			return err
		}
		// 无实例返回则结束
		if *r.Response.TotalCount == 0 {
			break
		}
		cvmList = append(cvmList, r.Response.InstanceSet...)
	}

	// 格式化数据
	for _, v := range cvmList {
		d, err := json.Marshal(v)
		if err != nil {
			log.Errorf("序列化cvm数据失败，%v", err)
			return err
		}
		dataList = append(dataList, resource.Data{
			Uuid:   fmt.Sprintf("tencent-cvm-%d", v.InstanceId),
			InfoId: infoID,
			Status: 0,
			Data:   d,
		})
	}

	// 写入数据
	err = orm.Eloquent.Model(&resource.Data{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.AssignmentColumns([]string{"data"}),
	}).Create(&dataList).Error

	return
}
