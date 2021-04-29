package cluster

import "fiy/common/models"

/*
  @Author : lanyulei
*/

// 集群规格，即集群配置信息
type Spec struct {
	Id                      int    `gorm:"column:id; primary_key;AUTO_INCREMENT" json:"id"`
	Version                 string `gorm:"column:version;" json:"version"`
	Provider                string `gorm:"column:provider;" json:"provider"`
	NetworkType             string `gorm:"column:network_type;" json:"network_type"`
	FlannelBackend          string `gorm:"column:flannel_backend;" json:"flannel_backend"`
	CalicoIpv4PoolIp        string `gorm:"column:calico_ipv4pool_ipip;" json:"calico_ipv4_pool_ip"`
	RuntimeType             string `gorm:"column:runtime_type;" json:"runtime_type"`
	DockerStorageDir        string `gorm:"column:docker_storage_dir;" json:"docker_storage_dir"`
	ContainerStorageDir     string `gorm:"column:containerd_storage_dir;" json:"container_storage_dir"`
	LbKubeApiServerIp       string `gorm:"column:lb_kube_apiserver_ip;" json:"lb_kube_api_server_ip"`
	KubeApiServerPort       int32  `gorm:"column:kube_api_server_port;" json:"kube_api_server_port"`
	KubeRouter              string `gorm:"column:kube_router;" json:"kube_router"`
	KubePodSubnet           string `gorm:"column:kube_pod_subnet;" json:"kube_pod_subnet"`
	KubeServiceSubnet       string `gorm:"column:kube_service_subnet;" json:"kube_service_subnet"`
	WorkerAmount            int32  `gorm:"column:worker_amount;" json:"worker_amount"`
	KubeMaxPods             int32  `gorm:"column:kube_max_pods;" json:"kube_max_pods"`
	KubeProxyMode           string `gorm:"column:kube_proxy_mode;" json:"kube_proxy_mode"`
	IngressControllerType   string `gorm:"column:ingress_controller_type;" json:"ingress_controller_type"`
	Architectures           string `gorm:"column:architectures;" json:"architectures"`
	KubernetesAudit         string `gorm:"column:kubernetes_audit;" json:"kubernetes_audit"`
	DockerSubnet            string `gorm:"column:docker_subnet;" json:"docker_subnet"`
	UpgradeVersion          string `gorm:"column:upgrade_version;" json:"upgrade_version"`
	HelmVersion             string `gorm:"column:helm_version;" json:"helm_version"`
	NetworkInterface        string `gorm:"column:network_interface;" json:"network_interface"`
	SupportGpu              string `gorm:"column:support_gpu;" json:"support_gpu"`
	YumOperate              string `gorm:"column:yum_operate;" json:"yum_operate"`
	KubeNetworkNodePrefix   int32  `gorm:"column:kube_network_node_prefix;default:16" json:"kube_network_node_prefix"`
	EnableDnsCache          string `gorm:"column:enable_dns_cache;" json:"enable_dns_cache"`
	DnsCacheVersion         string `gorm:"column:dns_cache_version;" json:"dns_cache_version"`
	CiliumVersion           string `gorm:"column:cilium_version;" json:"cilium_version"`
	CiliumTunnelMode        string `gorm:"column:cilium_tunnel_mode;" json:"cilium_tunnel_mode"`
	CiliumNativeRoutingCidr string `gorm:"column:cilium_native_routing_cidr;" json:"cilium_native_routing_cidr"`
	models.BaseModel
}

func (Spec) TableName() string {
	return "k8s_cluster_spec"
}
