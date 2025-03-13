package services

import (
	"fmt"
	"strings"

	"github.com/WillemCode/AliCloud_Resources/pkg/database"
	"github.com/WillemCode/AliCloud_Resources/pkg/logger"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
)

// SyncRDSInfo 同步指定账户和区域的 RDS 实例信息
func SyncRDSInfo(accountName string, rdsRegionIds []string, accessKey string, accessSecret string) error {
	for _, regionID := range rdsRegionIds {
		if regionID != "nil" && regionID != "" {
			logger.Log.Infof("开始同步信息, 区域=%s, 资源=RDS, 账户=%s", regionID, accountName)

			// 初始化 RDS 客户端
			client, err := rds.NewClientWithAccessKey(regionID, accessKey, accessSecret)
			if err != nil {
				return fmt.Errorf("RDS 客户端初始化失败 (账户=%s, 区域=%s): %w", accountName, regionID, err)
			}

			// 分页请求数据
			pageSize := 10  // 每页返回的条数，转换为 int64 类型
			pageNumber := 1 // 从第一页开始，转换为 int64 类型
			totalCount := 0 // 总条数，初始化为 0，转换为 int64 类型
			var records []database.RDSRecord
			for {
				// 构造请求
				// 构造请求并获取 RDS 实例列表

				// 获取 RDS 实例列表
				request := rds.CreateDescribeDBInstancesRequest()
				request.PageSize = requests.NewInteger(pageSize)     // 设置每页最大条数
				request.PageNumber = requests.NewInteger(pageNumber) // 设置当前页数

				response, err := client.DescribeDBInstances(request)
				if err != nil {
					return fmt.Errorf("RDS API 调用失败 (账户=%s, 区域=%s): %w", accountName, regionID, err)
				}

				// 获取总数
				if totalCount == 0 {
					totalCount = int(response.TotalRecordCount) // 获取总条数
					logger.Log.Infof("数据查询完成, 区域=%s, 资源=RDS, 账户=%s, 总数=%d 条", regionID, accountName, totalCount)
				}

				for _, instance := range response.Items.DBInstance {
					// 为每个 RDS 实例获取其网络连接信息（可能包含多个连接地址）
					netReq := rds.CreateDescribeDBInstanceNetInfoRequest()
					netReq.DBInstanceId = instance.DBInstanceId
					netResp, err := client.DescribeDBInstanceNetInfo(netReq)
					if err != nil {
						return fmt.Errorf("获取 RDS 网络信息失败 (InstanceID=%s): %w", instance.DBInstanceId, err)
					}
					// 收集所有连接地址并用逗号拼接
					var addressList []string
					for _, netInfo := range netResp.DBInstanceNetInfos.DBInstanceNetInfo {
						addressList = append(addressList, netInfo.ConnectionString)
					}
					connectionStr := strings.Join(addressList, ",")

					// 构造 RDSRecord
					rec := database.RDSRecord{
						InstanceID:       instance.DBInstanceId,
						CloudName:        accountName,
						Engine:           instance.Engine,
						RegionID:         instance.RegionId,
						Status:           instance.DBInstanceStatus,
						Memory:           int64(instance.DBInstanceMemory),
						Description:      instance.DBInstanceDescription,
						ConnectionString: connectionStr,
					}
					records = append(records, rec)
				}
				// 如果返回的数据条数小于 pageSize，说明已经拉取到最后一页，退出循环
				if len(response.Items.DBInstance) < pageSize {
					break
				}

				// 请求下一页数据
				pageNumber++
			}
			// 保存 RDS 数据
			if err := database.SaveRDSRecords(records); err != nil {
				return fmt.Errorf("保存 RDS 数据失败 (账户=%s): %w", accountName, err)
			}

			logger.Log.Infof("数据同步完成, 区域=%s, 资源=RDS, 账户=%s, 同步=%d 条", regionID, accountName, len(records))
		} else {
			logger.Log.Warnf("当前阿里账户, 区域=%s, 资源=RDS, 账户=%s, 暂无可用区域。", regionID, accountName)
		}
	}
	return nil
}
