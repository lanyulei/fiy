package aliyun

import "github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

/*
  @Author : lanyulei
*/

func getEcsList(sk, ak string, region []string) (err error) {
	var (
		ecsClient *ecs.Client
	)

	for _, r := range region {
		ecsClient, err = ecs.NewClientWithAccessKey(
			r,  // 您的可用区ID
			ak, // 您的AccessKey ID
			sk,
		) // 您的AccessKey Secret
		if err != nil {
			return
		}
	}

	ecsClient = ecsClient

	return
}
