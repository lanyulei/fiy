package cluster

import "fiy/common/models"

/*
  @Author : lanyulei
*/

// 集群规格，即集群配置信息
type Spec struct {
	Id                      int    `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	Version                 string `gorm:"column:version;type:varchar(255)" json:"version"`
	Provider                string `gorm:"column:provider;type:varchar(255)" json:"provider"`
	NetworkType             string `gorm:"column:network_type;type:varchar(255)" json:"network_type"`
	FlannelBackend          string `gorm:"column:flannel_backend;type:varchar(255)" json:"flannel_backend"`
	CalicoIpv4PoolIpip      string `gorm:"column:calico_ipv4pool_ipip;type:varchar(255)" json:"calico_ipv4pool_ipip"`
	RuntimeType             string `gorm:"column:runtime_type;type:varchar(255)" json:"runtime_type"`
	DockerStorageDir        string `gorm:"column:docker_storage_dir;type:varchar(255)" json:"docker_storage_dir"`
	ContainerdStorageDir    string `gorm:"column:containerd_storage_dir;type:varchar(255)" json:"containerd_storage_dir"`
	LbKubeApiserverIp       string `gorm:"column:lb_kube_apiserver_ip;type:varchar(255)" json:"lb_kube_apiserver_ip"`
	KubeApiServerPort       int    `gorm:"column:kube_api_server_port;type:int(11)" json:"kube_api_server_port"`
	KubeRouter              string `gorm:"column:kube_router;type:varchar(255)" json:"kube_router"`
	KubePodSubnet           string `gorm:"column:kube_pod_subnet;type:varchar(255)" json:"kube_pod_subnet"`
	KubeServiceSubnet       string `gorm:"column:kube_service_subnet;type:varchar(255)" json:"kube_service_subnet"`
	WorkerAmount            int    `gorm:"column:worker_amount;type:int(11)" json:"worker_amount"`
	KubeMaxPods             int    `gorm:"column:kube_max_pods;type:int(11)" json:"kube_max_pods"`
	KubeProxyMode           string `gorm:"column:kube_proxy_mode;type:varchar(255)" json:"kube_proxy_mode"`
	IngressControllerType   string `gorm:"column:ingress_controller_type;type:varchar(255)" json:"ingress_controller_type"`
	Architectures           string `gorm:"column:architectures;type:varchar(255)" json:"architectures"`
	KubernetesAudit         string `gorm:"column:kubernetes_audit;type:varchar(255)" json:"kubernetes_audit"`
	DockerSubnet            string `gorm:"column:docker_subnet;type:varchar(255)" json:"docker_subnet"`
	UpgradeVersion          string `gorm:"column:upgrade_version;type:varchar(255)" json:"upgrade_version"`
	HelmVersion             string `gorm:"column:helm_version;type:varchar(255)" json:"helm_version"`
	NetworkInterface        string `gorm:"column:network_interface;type:varchar(255)" json:"network_interface"`
	SupportGpu              string `gorm:"column:support_gpu;type:varchar(255)" json:"support_gpu"`
	YumOperate              string `gorm:"column:yum_operate;type:varchar(255)" json:"yum_operate"`
	KubeNetworkNodePrefix   int    `gorm:"column:kube_network_node_prefix;type:int(11);default:16" json:"kube_network_node_prefix"`
	EnableDnsCache          string `gorm:"column:enable_dns_cache;type:varchar(255)" json:"enable_dns_cache"`
	DnsCacheVersion         string `gorm:"column:dns_cache_version;type:varchar(255)" json:"dns_cache_version"`
	CiliumVersion           string `gorm:"column:cilium_version;type:varchar(255)" json:"cilium_version"`
	CiliumTunnelMode        string `gorm:"column:cilium_tunnel_mode;type:varchar(255)" json:"cilium_tunnel_mode"`
	CiliumNativeRoutingCidr string `gorm:"column:cilium_native_routing_cidr;type:varchar(255)" json:"cilium_native_routing_cidr"`
	models.BaseModel
}

func (Spec) TableName() string {
	return "k8s_cluster_spec"
}
