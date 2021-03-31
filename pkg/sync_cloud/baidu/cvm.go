package baidu

import (
	"encoding/json"
	"fiy/app/cmdb/models/resource"
	orm "fiy/common/global"
	"fiy/common/log"
	"fiy/tools"
	"fmt"
	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"gorm.io/gorm/clause"
)

/*
  @Author : lanyulei
*/

type baiDuYun struct {
	SK     string   `json:"sk"`
	AK     string   `json:"ak"`
	Region []string `json:"region"`
}

func NewBaiDuYun(sk, ak string, region []string) *baiDuYun {
	return &baiDuYun{
		SK:     sk,
		AK:     ak,
		Region: region,
	}
}

func (b *baiDuYun) BccList(infoID int) (err error) {
	var (
		result        *api.ListInstanceResult
		instancesList []resource.Data
		bccClient     *bcc.Client
		bccDataList   []api.InstanceModel
	)

	for _, r := range b.Region {
		bccClient, err = bcc.NewClient(
			tools.Strip(b.AK),
			tools.Strip(b.SK),
			tools.Strip(r),
		)
		if err != nil {
			log.Errorf("创建客户端连接失败，%v", err)
			return
		}

		args := &api.ListInstanceArgs{}

		result, err = bccClient.ListInstances(args)
		if err != nil {
			return
		}

		bccDataList = append(bccDataList, result.Instances...)
	}

	// 格式化数据
	for _, v := range bccDataList {
		var d []byte
		d, err = json.Marshal(v)
		if err != nil {
			log.Errorf("序列化服务器数据失败，%v", err)
			return
		}

		tmp := make(map[string]interface{})
		err = json.Unmarshal(d, &tmp)
		if err != nil {
			log.Error("反序列化数据失败，", err)
			return
		}

		tmp["instancesID"] = tmp["id"]
		delete(tmp, "id")
		d, err = json.Marshal(tmp)
		if err != nil {
			log.Errorf("序列化服务器数据失败，%v", err)
			return
		}

		instancesList = append(instancesList, resource.Data{
			Uuid:   fmt.Sprintf("baiduyun-bcc-%s", v.InstanceId),
			InfoId: infoID,
			Status: 0,
			Data:   d,
		})
	}

	// 写入数据
	err = orm.Eloquent.Model(&resource.Data{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.AssignmentColumns([]string{"data"}),
	}).Create(&instancesList).Error

	return
}
