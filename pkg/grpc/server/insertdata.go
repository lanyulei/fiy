package server

import (
	"encoding/json"
	"fiy/app/cmdb/models/model"
	"fiy/app/cmdb/models/resource"
	orm "fiy/common/global"
	"fiy/common/log"
	"reflect"
)

/*
  @Author : lanyulei
  @Desc : 存入数据
*/

// 更新map
func updateMap(uuid string, newMap map[string]interface{}) (jsonData []byte, err error) {
	var (
		data    resource.Data
		oldData map[string]interface{}
	)

	// 查询现有数据
	err = orm.Eloquent.Model(&resource.Data{}).Where("uuid = ?", uuid).Find(&data).Error
	if err != nil {
		log.Error("查询资源数据失败，", err)
		return
	}

	// 没有数据则直接返回新上报的数据
	if data.Id == 0 {
		jsonData, _ = json.Marshal(newMap)
		return
	}

	// 序列化旧的字段数据
	err = json.Unmarshal(data.Data, &oldData)
	if err != nil {
		log.Error("反序列化数据失败")
		return
	}

	// 更新新上报的数据
	for k, v := range newMap {
		oldData[k] = v
	}
	jsonData, _ = json.Marshal(oldData)

	return
}

func formatData(data string) (result map[string]interface{}, err error) {
	var (
		modelIDs []struct {
			ID         int    `json:"id"`
			Identifies string `json:"identifies"`
		}
		jsonData []byte
	)

	result = make(map[string]interface{})

	// 查询模型ID
	err = orm.Eloquent.Model(&model.Info{}).Where("identifies in ?", []string{
		"built_in_idc_host",
		"built_in_gpu",
		"built_in_memory",
		"built_in_cpu",
		"built_in_disk",
		"built_in_net",
	}).Find(&modelIDs).Error
	if err != nil {
		log.Info("查询模型ID失败，", err)
		return
	}

	modelIDMaps := make(map[string]int)
	for _, i := range modelIDs {
		modelIDMaps[i.Identifies] = i.ID
	}

	dataMap := map[string]interface{}{}
	err = json.Unmarshal([]byte(data), &dataMap)
	if err != nil {
		log.Info("json反序列化失败，", err)
		return
	}

	for k, v := range dataMap {
		if v != nil {
			if k == "info" {
				jsonData, err = updateMap(v.(map[string]interface{})["uuid"].(string), v.(map[string]interface{})["data"].(map[string]interface{}))
				result[k] = &resource.Data{
					Uuid:   v.(map[string]interface{})["uuid"].(string),
					InfoId: modelIDMaps[v.(map[string]interface{})["info_uuid"].(string)],
					Status: int(v.(map[string]interface{})["status"].(float64)),
					Data:   jsonData,
				}
			} else {
				dataList := make([]resource.Data, 0)
				for _, d := range v.([]interface{}) {
					jsonData, err = updateMap(d.(map[string]interface{})["uuid"].(string), d.(map[string]interface{})["data"].(map[string]interface{}))
					dataList = append(dataList, resource.Data{
						Uuid:   d.(map[string]interface{})["uuid"].(string),
						InfoId: modelIDMaps[d.(map[string]interface{})["info_uuid"].(string)],
						Status: int(d.(map[string]interface{})["status"].(float64)),
						Data:   jsonData,
					})
				}
				result[k] = &dataList
			}
		}
	}

	return
}

func insertData(data string) (err error) {
	var (
		result          map[string]interface{}
		hostInfo        resource.Data
		dataRelatedList []resource.DataRelated
		hostInfoCount   int64
		uuidList        []struct {
			Uuid      string `json:"uuid"`
			UuidCount int    `json:"uuid_count"`
		}
		insertDataList []resource.Data
		updateDataList []resource.Data
	)

	result, err = formatData(data)
	if err != nil {
		log.Error("重组数据失败，", err)
		return
	}

	// 查询数据是否存在
	err = orm.Eloquent.Model(&resource.Data{}).
		Where("uuid = ?", result["info"].(*resource.Data).Uuid).
		Count(&hostInfoCount).Error
	if err != nil {
		log.Error("查询主机信息失败，", err)
		return
	}

	tx := orm.Eloquent.Begin()

	if hostInfoCount > 0 {
		// update
		err = tx.Model(&resource.Data{}).
			Where("uuid = ?", result["info"].(*resource.Data).Uuid).
			Updates(map[string]interface{}{
				"data": result["info"].(*resource.Data).Data,
			}).Error
		if err != nil {
			tx.Rollback()
			log.Error("更新主机基础信息失败，", err)
			return
		}
	} else {
		// insert
		err = tx.Create(result["info"].(*resource.Data)).Error
		if err != nil {
			tx.Rollback()
			log.Error("新建主机基础信息失败，", err)
			return
		}
	}

	// 查询ID
	if result["info"].(*resource.Data).Id == 0 {
		err = orm.Eloquent.Model(&resource.Data{}).
			Where("uuid = ?", result["info"].(*resource.Data).Uuid).
			Find(&hostInfo).Error
		if err != nil {
			log.Error("查询数据ID失败，", err)
			tx.Rollback()
			return
		}
	} else {
		hostInfo = *result["info"].(*resource.Data)
	}

	dataUuids := make([]string, 0)
	for k, d := range result {
		if k != "info" {
			for _, z := range *d.(*[]resource.Data) {
				dataUuids = append(dataUuids, z.Uuid)
			}

			// 验证uuid是否存在
			err = tx.Model(&resource.Data{}).
				Where("uuid in ?", dataUuids).
				Select("uuid, count(uuid) as uuid_count").
				Group("uuid").
				Find(&uuidList).Error
			if err != nil {
				log.Error("UUID数据统计失败，", err)
				tx.Rollback()
				return
			}

			uuidMap := make(map[string]interface{})
			for _, u := range uuidList {
				if u.UuidCount > 0 {
					uuidMap[u.Uuid] = u.UuidCount
				}
			}

			for _, t := range *d.(*[]resource.Data) {
				if _, ok := uuidMap[t.Uuid]; ok {
					// update
					updateDataList = append(updateDataList, t)
				} else {
					// insert
					insertDataList = append(insertDataList, t)
				}
			}
		}
	}

	// insert
	if len(insertDataList) > 0 {
		err = tx.Create(&insertDataList).Error
		if err != nil {
			log.Error("插入数据失败，", err)
			tx.Rollback()
			return
		}
	}

	// update
	for _, d := range updateDataList {
		// 查询数据，判断数据是否有变更
		var tmpData resource.Data
		err = tx.Model(&tmpData).Where("uuid = ?", d.Uuid).Find(&tmpData).Error
		if err != nil {
			log.Error("查询数据失败，", err)
			tx.Rollback()
			return
		}

		var (
			oldMapData map[string]interface{}
			newMapData map[string]interface{}
		)

		_ = json.Unmarshal(tmpData.Data, &oldMapData)
		_ = json.Unmarshal(d.Data, &newMapData)

		if !reflect.DeepEqual(oldMapData, newMapData) {
			err = tx.Model(&resource.Data{}).Where("uuid = ?", d.Uuid).Update("data", d.Data).Error
			if err != nil {
				log.Error("更新数据失败", err)
				tx.Rollback()
				return
			}
		}
	}

	for _, i := range insertDataList {
		dataRelatedList = append(dataRelatedList, resource.DataRelated{
			Source:       hostInfo.Id,
			Target:       i.Id,
			SourceInfoId: hostInfo.InfoId,
			TargetInfoId: i.InfoId,
		})
	}

	if len(dataRelatedList) > 0 {
		err = tx.Create(&dataRelatedList).Error
		if err != nil {
			tx.Rollback()
			log.Error("创建数据关联失败，", err)
			return
		}
	}

	tx.Commit()

	return
}
