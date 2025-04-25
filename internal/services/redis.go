package services

import (
	"fmt"
	"strings"

	"github.com/WillemCode/AliCloud_Resources/pkg/database"
	"github.com/WillemCode/AliCloud_Resources/pkg/logger"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/r_kvstore"
)

// SyncECSInfo 同步指定账户和区域的 ECS 实例信息
func SyncRedisInfo(accountName string, redisRegionIds []string, accessKey string, accessSecret string) error {
	for _, regionID := range redisRegionIds {
		if regionID != "nil" && regionID != "" {
			logger.Log.Infof("开始同步信息, 区域=%s, 资源=Tair, 账户=%s", regionID, accountName)

			// 初始化 ECS 客户端
			client, err := r_kvstore.NewClientWithAccessKey(regionID, accessKey, accessSecret)
			if err != nil {
				return fmt.Errorf("tair Redis 客户端初始化失败 (区域=%s, 账户=%s): %w", regionID, accountName, err)
			}

			// 分页请求数据
			pageSize := 10                     // 每页返回的条数，转换为 int64 类型
			pageNumber := 1                    // 从第一页开始，转换为 int64 类型
			totalCount := 0                    // 总条数，初始化为 0，转换为 int64 类型
			var records []database.RedisRecord // 存储 Tair 实例记录
			for {
				// 构造请求
				// 构造请求并获取 Tair Redis 实例列表
				request := r_kvstore.CreateDescribeInstancesRequest()
				// 使用 requests.NewInteger 创建请求中的整数参数
				request.PageSize = requests.NewInteger(pageSize)     // 设置每页最大条数
				request.PageNumber = requests.NewInteger(pageNumber) // 设置当前页数

				response, err := client.DescribeInstances(request)
				if err != nil {
					return fmt.Errorf("tair Redis API 调用失败 (账户=%s, 区域=%s): %w", accountName, regionID, err)
				}

				// 获取总数
				if totalCount == 0 {
					totalCount = int(response.TotalCount) // 获取总条数
					logger.Log.Infof("数据查询完成, 区域=%s, 资源=Tair, 账户=%s, 总数=%d 条", regionID, accountName, totalCount)
				}

				// 将 API 返回的数据转换为本地 RedisRecord 列表
				// var records []database.RedisRecord
				for _, instance := range response.Instances.KVStoreInstance {
					// 为每个 Tair Redis 实例获取其网络连接信息（可能包含多个连接地址）
					tair_instance := r_kvstore.CreateDescribeDBInstanceNetInfoRequest()
					tair_instance.InstanceId = instance.InstanceId
					tairResp, err := client.DescribeDBInstanceNetInfo(tair_instance)
					if err != nil {
						return fmt.Errorf("获取 Tair Redis 网络信息失败 (InstanceID=%s): %w", instance.InstanceId, err)
					}
					// 收集所有连接地址并用逗号拼接
					var addressList []string
					var connectionList []string
					for _, addInfo := range tairResp.NetInfoItems.InstanceNetInfo {
						addressList = append(addressList, addInfo.IPAddress)
					}
					for _, conInfo := range tairResp.NetInfoItems.InstanceNetInfo {
						connectionList = append(connectionList, conInfo.ConnectionString)
					}
					addressStr := strings.Join(addressList, ",")
					connectionStr := strings.Join(connectionList, ",")

					// 构造 ECSRecord
					rec := database.RedisRecord{
						InstanceID:       instance.InstanceId,
						CloudName:        accountName,
						InstanceName:     instance.InstanceName,
						Port:             instance.Port,
						RegionId:         instance.RegionId,
						Capacity:         instance.Capacity,
						InstanceClass:    instance.InstanceClass,
						QPS:              instance.QPS,
						Bandwidth:        instance.Bandwidth,
						Connections:      instance.Connections,
						InstanceType:     instance.InstanceType,
						ConnectionString: connectionStr,
						IPAddress:        addressStr,
					}
					records = append(records, rec)
				}
				// 如果返回的数据条数小于 pageSize，说明已经拉取到最后一页，退出循环
				if len(response.Instances.KVStoreInstance) < pageSize {
					break
				}

				// 请求下一页数据
				pageNumber++
			}

			// 调用数据库包保存 ECS 数据
			if err := database.SaveRedisRecords(records); err != nil {
				return fmt.Errorf("保存 Tair Redis 数据失败 (账户=%s): %w", accountName, err)
			}
			logger.Log.Infof("数据同步完成, 区域=%s, 资源=Tair, 账户=%s, 同步=%d 条", regionID, accountName, len(records))
		} else {
			logger.Log.Warnf("当前阿里账户, 区域=%s, 资源=Tair, 账户=%s, 暂无可用区域。", regionID, accountName)
		}
	}
	return nil
}
