package client

import (
	"bytes"
	"encoding/json"
	"fiy/common/log"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/net"
	"github.com/yumaojun03/dmidecode"
)

/*
  @Author : lanyulei
*/

// 执行命令
func runCommand(timeout int, command string, args ...string) (out string, err error) {
	var (
		stdout io.ReadCloser
	)

	// instantiate new command
	cmd := exec.Command(command, args...)

	// get pipe to standard output
	stdout, err = cmd.StdoutPipe()
	if err != nil {
		return
	}

	// start process via command
	if err = cmd.Start(); err != nil {
		return
	}

	// setup a buffer to capture standard output
	var buf bytes.Buffer

	// create a channel to capture any errors from wait
	done := make(chan error)
	go func() {
		if _, err := buf.ReadFrom(stdout); err != nil {
			panic("buf.Read(stdout) error: " + err.Error())
		}
		done <- cmd.Wait()
	}()

	// block on select, and switch based on actions received
	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		if err = cmd.Process.Kill(); err != nil {
			return
		}
		return
	case err = <-done:
		if err != nil {
			close(done)
			return
		}
		out = buf.String()
		return
	}
}

// 获取uuid
func getUUID() (uuidString string, err error) {
	var (
		content []byte
		file    *os.File
	)

	_, err = os.Stat("./.uuid") //os.Stat获取文件信息
	if err != nil {
		if !os.IsExist(err) {
			// 文件不存在
			file, err = os.Create("./.uuid")
			if err != nil {
				fmt.Println("错误，", err)
				return
			}
			defer file.Close()
			err = ioutil.WriteFile("./.uuid", []byte(uuid.NewV4().String()), 0666)
			if err != nil {
				fmt.Println("错误，", err)
				return
			}
		}
	}

	// err == nil 文件存在
	content, err = ioutil.ReadFile("./.uuid")
	if err != nil {
		fmt.Println("错误，", err)
		return
	}
	uuidString = string(content)
	return
}

// 获取系统相关信息
func getHostInfo() (result map[string]interface{}, err error) {
	var (
		out              []byte
		fanCount         string
		hostInfo         *host.InfoStat
		sshdPort         string
		physicalCPUCount string
		cpuCoreCount     string
		uuidValue        string
	)
	hostInfo, err = host.Info()
	if err != nil {
		log.Info("查询主机信息错误，", err)
		return
	}

	// 获取风扇数量
	cmd := exec.Command("bash", "-c", "ipmitool sensor |grep FAN |wc -l")
	out, err = cmd.CombinedOutput()
	if err != nil {
		log.Info("获取风扇数量失败，错误：", err)
		fanCount = "错误：请确认ipmitool命令执行是否有异常"
		err = nil
	} else {
		fanCount = string(out)
		if len(fanCount) > 5 {
			fanCount = "0"
		}
	}

	// 获取sshd的端口
	cmd = exec.Command("bash", "-c", `netstat -tunlp |grep sshd |grep -v "^tcp6" |head -1 |awk '{print $4}' |awk -F ":" '{print $2}'`)
	out, err = cmd.CombinedOutput()
	if err != nil {
		log.Info("获取sshd端口失败，错误：", err)
		sshdPort = "错误：请确认netstat命令执行是否有异常"
		err = nil
	} else {
		sshdPort = string(out)
	}

	// cpu 物理个数
	cmd = exec.Command("bash", "-c", `cat /proc/cpuinfo | grep 'physical id' | sort | uniq | wc -l`)
	out, err = cmd.CombinedOutput()
	if err != nil {
		log.Info("获取物理CPU个数失败，错误：", err)
		physicalCPUCount = "异常"
		err = nil
	} else {
		physicalCPUCount = string(out)
	}

	// cpu 物理核数
	cmd = exec.Command("bash", "-c", `cat /proc/cpuinfo |grep "cores"|uniq|awk '{print $4}'`)
	out, err = cmd.CombinedOutput()
	if err != nil {
		log.Info("获取CPU核数失败，错误：", err)
		cpuCoreCount = "异常"
		err = nil
	} else {
		cpuCoreCount = string(out)
	}

	uuidValue, err = getUUID()
	if err != nil {
		log.Error("获取UUID失败，", err)
		return
	}

	result = map[string]interface{}{
		"uuid":      uuidValue,
		"info_uuid": "",
		"status":    1,
		"data": map[string]interface{}{
			"built_in_hostname":              hostInfo.Hostname,        // 主机名
			"built_in_host_platform":         hostInfo.Platform,        // 系统
			"built_in_host_platform_version": hostInfo.PlatformVersion, // 系统版本
			"built_in_kernel_version":        hostInfo.KernelVersion,   // 内核版本
			"built_in_kernel_arch":           hostInfo.KernelArch,      // 处理器架构
			"built_in_fan_count":             fanCount,                 // 风扇数量
			"built_in_ssh_port":              sshdPort,                 // sshd端口
			"built_in_cpu_count":             physicalCPUCount,         // CPU物理个数
			"built_in_cpu_core_count":        cpuCoreCount,             // CPU核数
		},
	}

	return
}

// 获取GPU相关信息
func getGPUInfo() (result []map[string]interface{}, err error) {
	var uuidValue string
	uuidValue, err = getUUID()
	if err != nil {
		log.Error("获取UUID失败，", err)
		return
	}

	s, err := runCommand(3, "bash",
		"-c",
		"nvidia-smi --query-gpu=uuid,gpu_name,driver_version,pstate,memory.total,power.draw,power.limit --format=csv")
	if err != nil {
		log.Error("执行 nvidia-smi 命令失败，", err)
		return
	}

	outLineList := strings.Split(s, "\n")
	if len(outLineList) > 0 {
		for idx, outValue := range outLineList[1:] {
			outValueLineList := strings.Split(outValue, ", ")
			if len(outValueLineList) > 0 && outValueLineList[0] != "" {
				result = append(result, map[string]interface{}{
					"uuid":      fmt.Sprintf("%s-gpu-%d", uuidValue, idx),
					"info_uuid": "",
					"status":    1,
					"data": map[string]interface{}{
						"built_in_gpu_uuid":          outValueLineList[0],
						"built_in_gpu_name":          outValueLineList[1],
						"built_in_gpu_version":       outValueLineList[2],
						"built_in_performance_state": outValueLineList[3],
						"built_in_gpu_memory":        outValueLineList[4],
						"built_in_power_waste":       outValueLineList[5],
						"built_in_power_limit":       outValueLineList[6],
					},
				})
			}
		}
	}

	return
}

// 获取内存相关信息
func getMemory() (result []map[string]interface{}, err error) {
	var (
		dmi       *dmidecode.Decoder
		uuidValue string
	)

	uuidValue, err = getUUID()
	if err != nil {
		log.Error("获取UUID失败，", err)
		return
	}

	dmi, err = dmidecode.New()
	if err != nil {
		log.Error("查询内存信息失败，", err)
		return
	}

	memoryDeviceList, err := dmi.MemoryDevice()
	if err != nil {
		log.Error("查询内存信息失败，", err)
		return
	}

	for idx, memoryDevice := range memoryDeviceList {
		if memoryDevice.Size > 0 {
			result = append(result, map[string]interface{}{
				"uuid":      fmt.Sprintf("%s-memory-%d", uuidValue, idx),
				"info_uuid": "",
				"status":    1,
				"data": map[string]interface{}{
					"built_in_memory_sn":           memoryDevice.SerialNumber,
					"built_in_memory_size":         fmt.Sprintf("%d G", memoryDevice.Size/1000),
					"built_in_memory_slot":         memoryDevice.DeviceLocator,
					"built_in_memory_type":         fmt.Sprintf("%v - %v", memoryDevice.Type.String(), memoryDevice.TypeDetail.String()),
					"built_in_memory_manufacturer": memoryDevice.Manufacturer,
				},
			})
		}
	}

	return
}

// 获取CPU相关信息
func getCPUInfo() (result []map[string]interface{}, err error) {
	var (
		cpuInfos  []cpu.InfoStat
		uuidValue string
	)

	uuidValue, err = getUUID()
	if err != nil {
		log.Error("获取UUID失败，", err)
		return
	}

	cpuInfos, err = cpu.Info()
	if err != nil {
		log.Error("查询CPU相关信息失败，", err)
		return
	}
	for idx, info := range cpuInfos {
		result = append(result, map[string]interface{}{
			"uuid":      fmt.Sprintf("%s-cpu-%d", uuidValue, idx),
			"info_uuid": "",
			"status":    1,
			"data": map[string]interface{}{
				"built_in_cpu_vendor_id":   info.VendorID,   // CPU制造商
				"built_in_cpu_family":      info.Family,     // 系列
				"built_in_cpu_model":       info.Model,      // 其系列中的哪一代的代号
				"built_in_cpu_physical_id": info.PhysicalID, // 编号
				"built_in_cpu_cores":       info.Cores,      // 核心数
				"built_in_cpu_model_name":  info.ModelName,  // 型号
				"built_in_cpu_cache_size":  info.CacheSize,  // 缓存大小
			},
		})
	}

	return
}

// 获取磁盘相关信息
func getDiskInfo() (result []map[string]interface{}, err error) {
	var (
		infos     []disk.PartitionStat
		diskInfo  *disk.UsageStat
		uuidValue string
	)

	uuidValue, err = getUUID()
	if err != nil {
		log.Error("获取UUID失败，", err)
		return
	}

	infos, err = disk.Partitions(false)
	if err != nil {
		log.Error("查询磁盘信息失败，", err)
		return
	}
	for idx, info := range infos {
		var (
			diskInfoTotal       uint64
			diskInfoFree        uint64
			diskInfoUsed        uint64
			diskInfoUsedPercent float64
		)
		diskInfo, err = disk.Usage(info.Mountpoint)
		if err != nil {
			log.Error("查询磁盘信息失败，", err)
			return
		}
		if diskInfo != nil {
			diskInfoTotal = diskInfo.Total / 1024 / 1024 / 1024
			diskInfoFree = diskInfo.Free / 1024 / 1024 / 1024
			diskInfoUsed = diskInfo.Used / 1024 / 1024 / 1024
			diskInfoUsedPercent = diskInfo.UsedPercent
		}
		result = append(result, map[string]interface{}{
			"uuid":      fmt.Sprintf("%s-disk-%d", uuidValue, idx),
			"info_uuid": "",
			"status":    1,
			"data": map[string]interface{}{
				"built_in_disk_device_name":  info.Device,                                 // 设备名称
				"built_in_disk_mount_point":  info.Mountpoint,                             // 挂载地址
				"built_in_disk_fstype":       info.Fstype,                                 // 文件系统
				"built_in_disk_total_size":   fmt.Sprintf("%v G", diskInfoTotal),          // 磁盘总大小
				"built_in_disk_free_size":    fmt.Sprintf("%v G", diskInfoFree),           // 剩余磁盘大小
				"built_in_disk_used_size":    fmt.Sprintf("%v G", diskInfoUsed),           // 已使用磁盘大小
				"built_in_disk_used_percent": fmt.Sprintf("%.2f %%", diskInfoUsedPercent), // 已使用磁盘大小占比
			},
		})
	}

	return
}

// 获取网卡信息
func getNetInfo() (result []map[string]interface{}, err error) {
	var (
		interfaces []net.InterfaceStat
		uuidValue  string
	)
	uuidValue, err = getUUID()
	if err != nil {
		log.Error("获取UUID失败，", err)
		return
	}
	interfaces, err = net.Interfaces()
	if err != nil {
		log.Error("查询网卡信息失败，", err)
		return
	}
	for idx, inter := range interfaces {
		ipaddress := make([]string, 0)
		for _, addrs := range inter.Addrs {
			if strings.HasPrefix(addrs.Addr, "1") || strings.HasPrefix(addrs.Addr, "2") {
				ipaddress = append(ipaddress, addrs.Addr)
			}
		}
		ipAddressString := strings.Join(ipaddress, ", ")
		if len(ipaddress) > 0 && !strings.HasPrefix(ipAddressString, "10.") && !strings.HasPrefix(ipAddressString, "127.0.0.1") {
			result = append(result, map[string]interface{}{
				"uuid":      fmt.Sprintf("%s-net-%d", uuidValue, idx),
				"info_uuid": "",
				"status":    0,
				"data": map[string]interface{}{
					"built_in_net_index": inter.Index,        // 索引
					"built_in_net_name":  inter.Name,         // 网卡名称
					"built_in_net_mac":   inter.HardwareAddr, // mac地址
					"built_in_net_ip":    ipAddressString,    // ip地址
				},
			})
		}
	}
	return
}

// 整合数据
func IntegrateData() (dataString string, err error) {
	var dataBytes []byte
	//data := make(map[string]interface{})

	//data["info"], _ = getHostInfo()
	//data["gpu"], _ = getGPUInfo()
	//data["memory"], _ = getMemory()
	//data["cpu"], _ = getCPUInfo()
	//data["disk"], _ = getDiskInfo()
	//data["net"], _ = getNetInfo()

	data := map[string]interface{}{
		"info": map[string]interface{}{
			"uuid":      "d8205faa-a71e-480a-b852-6543c11bbdb6",
			"info_uuid": "built_in_idc_host",
			"status":    1,
			"data": map[string]interface{}{
				"built_in_hostname":              "test1", // 主机名
				"built_in_host_platform":         "test2", // 系统
				"built_in_host_platform_version": "test3", // 系统版本
				"built_in_kernel_version":        "test3", // 内核版本
				"built_in_kernel_arch":           "test3", // 处理器架构
				"built_in_fan_count":             "test3", // 风扇数量
				"built_in_ssh_port":              "test3", // sshd端口
				"built_in_cpu_count":             "test3", // CPU物理个数
				"built_in_cpu_core_count":        "test3", // CPU核数
			},
		},
		"gpu": []map[string]interface{}{{
			"uuid":      "d8205faa-a71e-480a-b852-6543c11bbdb62",
			"info_uuid": "built_in_gpu",
			"status":    1,
			"data": map[string]interface{}{
				"built_in_gpu_uuid":          "test1",
				"built_in_gpu_name":          "test2",
				"built_in_gpu_version":       "test3",
				"built_in_performance_state": "test4",
				"built_in_gpu_memory":        "test5",
				"built_in_power_waste":       "test6",
				"built_in_power_limit":       "test7",
			},
		}},
		"memory": []map[string]interface{}{{
			"uuid":      uuid.NewV4().String(),
			"info_uuid": "built_in_memory",
			"status":    1,
			"data": map[string]interface{}{
				"built_in_memory_sn":           "test1",
				"built_in_memory_size":         "test2",
				"built_in_memory_slot":         "test3",
				"built_in_memory_type":         "test4",
				"built_in_memory_manufacturer": "test5",
			},
		}},
		"cpu": []map[string]interface{}{{
			"uuid":      uuid.NewV4().String(),
			"info_uuid": "built_in_cpu",
			"status":    1,
			"data": map[string]interface{}{
				"built_in_cpu_vendor_id":   "test1", // CPU制造商
				"built_in_cpu_family":      "test2", // 系列
				"built_in_cpu_model":       "test3", // 其系列中的哪一代的代号
				"built_in_cpu_physical_id": "test4", // 编号
				"built_in_cpu_cores":       "test5", // 核心数
				"built_in_cpu_model_name":  "test6", // 型号
				"built_in_cpu_cache_size":  "test7", // 缓存大小
			},
		}},
		"disk": []map[string]interface{}{{
			"uuid":      uuid.NewV4().String(),
			"info_uuid": "built_in_disk",
			"status":    1,
			"data": map[string]interface{}{
				"built_in_disk_device_name":  "test1", // 设备名称
				"built_in_disk_mount_point":  "test2", // 挂载地址
				"built_in_disk_fstype":       "test3", // 文件系统
				"built_in_disk_total_size":   "test4", // 磁盘总大小
				"built_in_disk_free_size":    "test5", // 剩余磁盘大小
				"built_in_disk_used_size":    "test6", // 已使用磁盘大小
				"built_in_disk_used_percent": "test7", // 已使用磁盘大小占比
			},
		}},
		"net": []map[string]interface{}{{
			"uuid":      "d8205faa-a71e-480a-b852-6543c11bbdb63",
			"info_uuid": "built_in_net",
			"status":    1,
			"data": map[string]interface{}{
				"built_in_net_index": "test1", // 索引
				"built_in_net_name":  "test2", // 网卡名称
				"built_in_net_mac":   "test3", // mac地址
				"built_in_net_ip":    "test4", // ip地址
			},
		}, {
			"uuid":      "d8205faa-a71e-480a-b852-6543c11bbdb64",
			"info_uuid": "built_in_net",
			"status":    1,
			"data": map[string]interface{}{
				"built_in_net_index": "test11", // 索引
				"built_in_net_name":  "test22", // 网卡名称
				"built_in_net_mac":   "test33", // mac地址
				"built_in_net_ip":    "test44", // ip地址
			},
		}},
	}

	dataBytes, err = json.Marshal(data)
	if err != nil {
		log.Error("序列化数据失败，", err)
		return
	}

	dataString = string(dataBytes)

	return
}
