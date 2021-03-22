package aliyun

import (
	"fmt"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

/*
  @Author : lanyulei
*/

type aliyun struct {
	SK        string   `json:"sk"`
	AK        string   `json:"ak"`
	Region    []string `json:"region"`
	ecsClient *ecs.Client
}

func NewAliyun(sk, ak string, region []string) *aliyun {
	return &aliyun{
		SK:     sk,
		AK:     ak,
		Region: region,
	}
}

func (a *aliyun) GetEcsList() (err error) {
	var (
		result []map[string]interface{}
	)

	for _, r := range a.Region {
		r = strings.Trim(r, " ")
		r = strings.Trim(r, "\t")
		r = strings.Trim(r, "\n")
		r = strings.Trim(r, "\r")
		a.ecsClient, err = ecs.NewClientWithAccessKey(
			r,
			a.AK,
			a.SK,
		)
		if err != nil {
			return
		}

		result, err = a.getAliyunECSList()
		for _, r := range result {
			fmt.Println(r)
		}
	}

	return
}

func (a *aliyun) getAliyunECSList() (result []map[string]interface{}, err error) {

	request := ecs.CreateDescribeInstancesRequest()
	request.PageSize = "1"

	response, err := a.ecsClient.DescribeInstances(request)
	if err != nil {
		return
	}

	response = response

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
