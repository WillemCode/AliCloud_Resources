package services

import (
	"fmt"
	"strings"

	"github.com/WillemCode/AliCloud_Resources/pkg/database"
	"github.com/WillemCode/AliCloud_Resources/pkg/logger"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

// SyncECSInfo 同步指定账户和区域的 ECS 实例信息
func SyncECSInfo(accountName string, ecsRegionIds []string, accessKey string, accessSecret string) error {
	for _, regionID := range ecsRegionIds {
		if regionID != "nil" && regionID != "" {
			logger.Log.Infof("开始同步信息, 区域=%s, 资源=ECS, 账户=%s", regionID, accountName)

			// 初始化 ECS 客户端
			client, err := ecs.NewClientWithAccessKey(regionID, accessKey, accessSecret)
			if err != nil {
				return fmt.Errorf("ECS客户端初始化失败 (区域=%s, 账户=%s): %w", regionID, accountName, err)
			}

			// 分页请求数据
			pageSize := 10                   // 每页返回的条数，转换为 int64 类型
			pageNumber := 1                  // 从第一页开始，转换为 int64 类型
			totalCount := 0                  // 总条数，初始化为 0，转换为 int64 类型
			var records []database.ECSRecord // 存储 ECS 实例记录
			for {
				// 构造请求
				// 构造请求并获取 ECS 实例列表
				request := ecs.CreateDescribeInstancesRequest()
				// 使用 requests.NewInteger 创建请求中的整数参数
				request.PageSize = requests.NewInteger(pageSize)     // 设置每页最大条数
				request.PageNumber = requests.NewInteger(pageNumber) // 设置当前页数

				response, err := client.DescribeInstances(request)
				if err != nil {
					return fmt.Errorf("ECS API 调用失败 (账户=%s, 区域=%s): %w", accountName, regionID, err)
				}

				// 获取总数
				if totalCount == 0 {
					totalCount = int(response.TotalCount) // 获取总条数
					logger.Log.Infof("数据查询完成, 区域=%s, 资源=ECS, 账户=%s, 总数=%d 条", regionID, accountName, totalCount)
				}

				// 将 API 返回的数据转换为本地 ECSRecord 列表
				// var records []database.ECSRecord
				for _, instance := range response.Instances.Instance {
					// 收集所有公网 IP（自带公网 IP、EIP、网卡级 EIP）
					var publicIPList []string
					// 1) 实例自带公网 IP（可能有多个）
					if len(instance.PublicIpAddress.IpAddress) > 0 {
						publicIPList = append(publicIPList, instance.PublicIpAddress.IpAddress...)
					}
					// 2) 实例主网卡上的 EIP
					if instance.EipAddress.IpAddress != "" {
						publicIPList = append(publicIPList, instance.EipAddress.IpAddress)
					}
					// // 3) 遍历所有弹性网卡，检查是否有 EIP
					// for _, eni := range instance.NetworkInterfaces.NetworkInterface {
					// 	if eni.EipAddress.IpAddress != "" {
					// 		publicIPList = append(publicIPList, eni.EipAddress.IpAddress)
					// 	}
					// }
					// 收集私网 IP（如果有多张网卡，这里只取第一个主网卡）
					privateIP := ""
					if len(instance.NetworkInterfaces.NetworkInterface) > 0 {
						privateIP = instance.NetworkInterfaces.NetworkInterface[0].PrimaryIpAddress
					}
					// 用逗号拼接所有收集到的公网 IP
					publicIPs := strings.Join(publicIPList, ",")
					// // 提取私网 IP（如果存在多个，仅取第一个）
					// privateIP := ""
					// if len(instance.NetworkInterfaces.NetworkInterface) > 0 {
					// 	privateIP = instance.NetworkInterfaces.NetworkInterface[0].PrimaryIpAddress
					// }
					// // 提取公网 IP 列表并用逗号拼接
					// publicIPs := ""
					// if len(instance.PublicIpAddress.IpAddress) > 0 {
					// 	publicIPs = strings.Join(instance.PublicIpAddress.IpAddress, ",")
					// }

					// 构造 ECSRecord
					rec := database.ECSRecord{
						InstanceID:   instance.InstanceId,
						CloudName:    accountName,
						InstanceName: instance.InstanceName,
						Status:       instance.Status,
						RegionID:     instance.RegionId,
						OSName:       instance.OSName,
						InstanceType: instance.InstanceType,
						CPU:          int64(instance.Cpu),
						Memory:       int64(instance.Memory),
						PublicIP:     publicIPs,
						PrivateIP:    privateIP,
					}
					records = append(records, rec)
				}
				// 如果返回的数据条数小于 pageSize，说明已经拉取到最后一页，退出循环
				if len(response.Instances.Instance) < pageSize {
					break
				}

				// 请求下一页数据
				pageNumber++
			}

			// 调用数据库包保存 ECS 数据
			if err := database.SaveECSRecords(records); err != nil {
				return fmt.Errorf("保存 ECS 数据失败 (账户=%s): %w", accountName, err)
			}
			logger.Log.Infof("数据同步完成, 区域=%s, 资源=ECS, 账户=%s, 同步=%d 条", regionID, accountName, len(records))
		} else {
			logger.Log.Warnf("当前阿里账户, 区域=%s, 资源=ECS, 账户=%s, 暂无可用区域。", regionID, accountName)
		}
	}
	return nil
}
