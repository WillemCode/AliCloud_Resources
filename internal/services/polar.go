package services

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/WillemCode/AliCloud_Resources/pkg/database"
	"github.com/WillemCode/AliCloud_Resources/pkg/logger"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
)

// SyncPolarDBInfo 同步指定账户和区域的 PolarDB 实例信息
func SyncPolarDBInfo(accountName string, polarRegionIds []string, accessKey string, accessSecret string) error {
	for _, regionID := range polarRegionIds {
		if regionID != "nil" && regionID != "" {
			logger.Log.Infof("开始同步信息, 区域=%s, 资源=polar, 账户=%s", regionID, accountName)

			// 初始化 PolarDB 客户端
			client, err := polardb.NewClientWithAccessKey(regionID, accessKey, accessSecret)
			if err != nil {
				return fmt.Errorf("PolarDB 客户端初始化失败 (账户=%s, 区域=%s): %w", accountName, regionID, err)
			}

			// 分页请求数据
			pageSize := 30  // 每页返回的条数，转换为 int64 类型
			pageNumber := 1 // 从第一页开始，转换为 int64 类型
			totalCount := 0 // 总条数，初始化为 0，转换为 int64 类型
			var records []database.PolarDBRecord
			for {

				// 获取 PolarDB 集群列表
				request := polardb.CreateDescribeDBClustersRequest()
				request.PageSize = requests.NewInteger(pageSize)     // 设置每页最大条数
				request.PageNumber = requests.NewInteger(pageNumber) // 设置当前页数

				response, err := client.DescribeDBClusters(request)
				if err != nil {
					return fmt.Errorf("PolarDB API 调用失败 (账户=%s, 区域=%s): %w", accountName, regionID, err)
				}

				// 获取总数
				if totalCount == 0 {
					totalCount = int(response.TotalRecordCount) // 获取总条数
					logger.Log.Infof("数据查询完成, 区域=%s, 资源=polar, 账户=%s, 总数=%d 条", regionID, accountName, totalCount)
				}

				for _, cluster := range response.Items.DBCluster {
					// 获取每个集群的连接地址列表
					epReq := polardb.CreateDescribeDBClusterEndpointsRequest()
					epReq.DBClusterId = cluster.DBClusterId
					epResp, err := client.DescribeDBClusterEndpoints(epReq)
					if err != nil {
						return fmt.Errorf("获取 PolarDB 连接信息失败 (DBClusterID=%s): %w", cluster.DBClusterId, err)
					}
					// 收集所有连接地址并用逗号拼接
					var addrList []string
					for _, ep := range epResp.Items {
						for _, addr := range ep.AddressItems {
							addrList = append(addrList, addr.ConnectionString)
						}
					}
					connectionStr := strings.Join(addrList, ",")

					// 将 MemorySize 从 string 转换为 int64
					memorySize, err := strconv.ParseInt(cluster.MemorySize, 10, 64)
					if err != nil {
						return fmt.Errorf("解析 MemorySize 失败 (DBClusterID=%s): %w", cluster.DBClusterId, err)
					}

					// 构造 PolarDBRecord
					rec := database.PolarDBRecord{
						InstanceID:       cluster.DBClusterId,
						CloudName:        accountName,
						Engine:           cluster.Engine,
						RegionID:         cluster.RegionId,
						Status:           cluster.DBClusterStatus,
						DBNodeCount:      int64(cluster.DBNodeNumber),
						Description:      cluster.DBClusterDescription,
						MemorySize:       memorySize,
						ConnectionString: connectionStr,
					}
					records = append(records, rec)
				}
				// 如果返回的数据条数小于 pageSize，说明已经拉取到最后一页，退出循环
				if len(response.Items.DBCluster) < pageSize {
					break
				}

				// 请求下一页数据
				pageNumber++
			}
			// 保存 PolarDB 数据
			if err := database.SavePolarDBRecords(records); err != nil {
				return fmt.Errorf("保存 PolarDB 数据失败 (账户=%s): %w", accountName, err)
			}

			logger.Log.Infof("数据同步完成, 区域=%s, 资源=polar, 账户=%s, 同步=%d 条", regionID, accountName, len(records))
		} else {
			logger.Log.Warnf("当前阿里账户, 区域=%s, 资源=polar, 账户=%s, 暂无可用区域。", regionID, accountName)
		}
	}
	return nil
}
