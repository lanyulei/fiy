package business

import (
	"fiy/app/cmdb/models/model"
	"fiy/app/cmdb/models/resource"
	orm "fiy/common/global"
	"fiy/tools/app"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

// 业务树
func BusinessTree(c *gin.Context) {
	var (
		err      error
		biz      model.Info
		treeData []*struct {
			resource.Data
			Children []interface{} `json:"children"`
		}
	)

	// 查询业务模型数据
	err = orm.Eloquent.Model(&model.Info{}).
		Where("identifies = 'built_in_biz'").
		Find(&biz).Error
	if err != nil {
		app.Error(c, -1, err, "查询模型数据失败")
		return
	}
	// 所有业务ID
	bizIdList := make([]int, 0)
	err = orm.Eloquent.Model(&resource.Data{}).Where("info_id = ?", biz.Id).Pluck("id", &bizIdList).Error
	if err != nil {
		app.Error(c, -1, err, "查询业务数据列表失败")
		return
	}
	if len(bizIdList) > 0 {
		// 查询业务数据列表
		err = orm.Eloquent.Model(&resource.Data{}).Find(&treeData, bizIdList).Error
		if err != nil {
			app.Error(c, -1, err, "查询模块列表失败")
			return
		}
	}

	// 所有业务对应的集群ID
	setIdList := make([]int, 0)
	setList := make([]*struct {
		resource.Data
		Biz      int           `json:"biz"`
		Children []interface{} `json:"children"`
	}, 0)
	if len(bizIdList) > 0 {
		err = orm.Eloquent.Model(&resource.Data{}).
			Joins("left join cmdb_resource_data_related as dr on dr.target = cmdb_resource_data.id").
			Where("dr.source in ?", bizIdList).
			Select("cmdb_resource_data.*, dr.source as biz").
			Find(&setList).Error
		if err != nil {
			app.Error(c, -1, err, "查询集群数据失败")
			return
		}

		for _, s := range setList {
			setIdList = append(setIdList, s.Id)
		}
	}

	// 所有业务对应的集群ID
	moduleList := make([]*struct {
		resource.Data
		Cluster int `json:"cluster"`
	}, 0)
	if len(setIdList) > 0 {
		err = orm.Eloquent.Model(&resource.Data{}).
			Joins("left join cmdb_resource_data_related as dr on dr.target = cmdb_resource_data.id").
			Where("dr.source in ?", setIdList).
			Select("cmdb_resource_data.*, dr.source as cluster").
			Find(&moduleList).Error
		if err != nil {
			app.Error(c, -1, err, "查询模块数据失败")
			return
		}
	}

	for _, t := range treeData {
		for _, s := range setList {
			for _, m := range moduleList {
				if m.Cluster == s.Id {
					s.Children = append(s.Children, m)
					continue
				}
			}
			if s.Biz == t.Id {
				t.Children = append(t.Children, s)
				continue
			}
		}
	}

	app.OK(c, treeData, "")
}
