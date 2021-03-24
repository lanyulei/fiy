package aliyun

import (
	"encoding/json"
	"fmt"

	"fiy/app/cmdb/models/resource"
	orm "fiy/common/global"
	"gorm.io/gorm/clause"

	"fiy/common/log"
	"fiy/tools"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

/*
  @Author : lanyulei
*/

type aLiYun struct {
	SK     string   `json:"sk"`
	AK     string   `json:"ak"`
	Region []string `json:"region"`
}

func NewALiYun(sk, ak string, region []string) *aLiYun {
	return &aLiYun{
		SK:     sk,
		AK:     ak,
		Region: region,
	}
}

func (a *aLiYun) EcsList(infoID int) (err error) {
	var (
		response  *ecs.DescribeInstancesResponse
		ecsList   []ecs.Instance
		dataList  []resource.Data
		ecsClient *ecs.Client
	)

	for _, r := range a.Region {
		ecsClient, err = ecs.NewClientWithAccessKey(
			tools.Strip(r),
			tools.Strip(a.AK),
			tools.Strip(a.SK),
		)
		if err != nil {
			log.Errorf("创建客户端连接失败，%v", err)
			return
		}

		request := ecs.CreateDescribeInstancesRequest()
		request.PageSize = "1"

		response, err = ecsClient.DescribeInstances(request)
		if err != nil {
			log.Errorf("查询ECS实例列表失败，%v", err)
			return
		}

		if response.TotalCount > 0 {
			for i := 0; i < response.TotalCount/100+1; i++ {
				request.PageSize = "100"
				r, err := ecsClient.DescribeInstances(request)
				if err != nil {
					log.Errorf("查询ECS实例列表失败，%v", err)
					return err
				}

				ecsList = append(ecsList, r.Instances.Instance...)
			}

		}
	}

	// 格式化数据
	for _, v := range ecsList {
		d, err := json.Marshal(v)
		if err != nil {
			log.Errorf("序列化ecs数据失败，%v", err)
			return err
		}
		dataList = append(dataList, resource.Data{
			Uuid:   fmt.Sprintf("aliyun-ecs-%s", v.InstanceId),
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

// 获取续费相关信息
//func (a *aliyun) getEcsRenewInfo() (result []ecs.InstanceRenewAttribute, err error) {
//	var (
//		response        *ecs.DescribeInstanceAutoRenewAttributeResponse
//		TotalPageNumber int
//	)
//	for _, v := range []string{"AutoRenewal", "Normal", "NotRenewal"} {
//		request := ecs.CreateDescribeInstanceAutoRenewAttributeRequest()
//		//request.InstanceId = ecsInstancesId
//		request.RenewalStatus = v
//		request.PageSize = "1"
//
//		response, err = a.ecsClient.DescribeInstanceAutoRenewAttribute(request)
//		if err != nil {
//			return
//		}
//		if response.TotalCount > 0 {
//			TotalPageNumber = response.TotalCount/100 + 1
//
//			for i := 0; i < TotalPageNumber; i++ {
//				request.PageSize = "100"
//				response, err = a.ecsClient.DescribeInstanceAutoRenewAttribute(request)
//				if err != nil {
//					return
//				}
//				for _, r := range response.InstanceRenewAttributes.InstanceRenewAttribute {
//					result = append(result, r)
//				}
//			}
//		}
//	}
//	return
//}

// 查询实例磁盘信息
//func (a *aliyun) getEcsDisk() (result []ecs.Disk, err error) {
//	var (
//		response        *ecs.DescribeDisksResponse
//		TotalPageNumber int
//	)
//	request := ecs.CreateDescribeDisksRequest()
//	request.PageSize = "1"
//
//	response, err = a.ecsClient.DescribeDisks(request)
//	if err != nil {
//		return
//	}
//	if response.TotalCount > 0 {
//		TotalPageNumber = response.TotalCount/100 + 1
//
//		for i := 0; i < TotalPageNumber; i++ {
//			request.PageSize = "100"
//			response, err = a.ecsClient.DescribeDisks(request)
//			if err != nil {
//				return
//			}
//			for _, r := range response.Disks.Disk {
//				result = append(result, r)
//			}
//		}
//	}
//
//	return
//}
