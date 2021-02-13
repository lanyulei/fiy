package business

import "fiy/common/models"

/*
  @Author : lanyulei
  @Desc : 服务模板进程
*/

type ServiceTemplateProcess struct {
	Id          int    `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	Name        string `gorm:"column:name; type:varchar(128);" json:"name" binding:"required"`   // 名称
	Alias       string `gorm:"column:alias; type:varchar(128);" json:"alias" binding:"required"` // 别名
	StartParams string `gorm:"column:start_params; type:varchar(512);" json:"start_params"`      // 启动参数

	BindIP   string `gorm:"column:bind_ip; type:varchar(45);" json:"bind_ip"`    // 绑定IP
	Port     string `gorm:"column:port; type:varchar(45);" json:"port"`          // 端口
	Protocol string `gorm:"column:protocol; type:varchar(128);" json:"protocol"` // 协议

	WorkPath         string `gorm:"column:work_path; type:varchar(512);" json:"work_path"`                   // 工作路径
	StartUser        string `gorm:"column:start_user; type:varchar(128);" json:"start_user"`                 // 启动用户
	StartNumber      int    `gorm:"column:start_number; type:int(11);" json:"start_number"`                  // 启动数量
	StartPriority    int    `gorm:"column:start_priority; type:int(11);" json:"start_priority"`              // 启动优先级
	StartTimeout     int    `gorm:"column:start_timeout; type:int(11);" json:"start_timeout"`                // 启动超时时长（秒）
	StartCommand     string `gorm:"column:start_command; type:varchar(128);" json:"start_command"`           // 启动命令
	StopCommand      string `gorm:"column:stop_command; type:varchar(128);" json:"stop_command"`             // 停止命令
	RestartCommand   string `gorm:"column:restart_command; type:varchar(128);" json:"restart_command"`       // 重启命令
	ForceStopCommand string `gorm:"column:force_stop_command; type:varchar(128);" json:"force_stop_command"` // 强制停止命令
	ReloadCommand    string `gorm:"column:reload_command; type:varchar(128);" json:"reload_command"`         // 重载命令
	PidPath          string `gorm:"column:pid_path; type:varchar(128);" json:"pid_path"`                     // PID文件
	RestartInterval  int    `gorm:"column:restart_interval; type:varchar(128);" json:"restart_interval"`     // 重启间隔时间（秒）

	GatewayIP       string `gorm:"column:gateway_ip; type:varchar(128);" json:"gateway_ip"`            // 网关IP
	GatewayPort     string `gorm:"column:gateway_port; type:varchar(50);" json:"gateway_port"`         // 网关端口
	GatewayProtocol string `gorm:"column:gateway_protocol; type:varchar(50);" json:"gateway_protocol"` // 网关协议

	SvcTpl int `gorm:"column:svc_tpl; type:int(11);" json:"svc_tpl"` // 绑定服务模板

	Remark string `gorm:"column:remark; type:varchar(1024);" json:"remark"` // 备注
	models.BaseModel
}

func (ServiceTemplateProcess) TableName() string {
	return "cmdb_business_svc_tpl_process"
}
