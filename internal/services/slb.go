package services

import (
	"fmt"

	"github.com/WillemCode/AliCloud_Resources/pkg/database"
	"github.com/WillemCode/AliCloud_Resources/pkg/logger"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
)

// SyncSLBInfo 同步指定账户和区域的 SLB 实例信息
func SyncSLBInfo(accountName string, slbRegionIds []string, accessKey string, accessSecret string) error {
	for _, regionID := range slbRegionIds {
		if regionID != "nil" && regionID != "" {
			logger.Log.Infof("开始同步信息, 区域=%s, 资源=CLB, 账户=%s", regionID, accountName)

			// 初始化 SLB 客户端
			client, err := slb.NewClientWithAccessKey(regionID, accessKey, accessSecret)
			if err != nil {
				return fmt.Errorf("SLB 客户端初始化失败 (账户=%s, 区域=%s): %w", accountName, regionID, err)
			}

			// 分页请求数据
			pageSize := 10  // 每页返回的条数，转换为 int64 类型
			pageNumber := 1 // 从第一页开始，转换为 int64 类型
			totalCount := 0 // 总条数，初始化为 0，转换为 int64 类型
			var records []database.SLBRecord

			for {

				// 获取 SLB 实例列表
				request := slb.CreateDescribeLoadBalancersRequest()
				request.PageSize = requests.NewInteger(pageSize)     // 设置每页最大条数
				request.PageNumber = requests.NewInteger(pageNumber) // 设置当前页数

				response, err := client.DescribeLoadBalancers(request)
				if err != nil {
					return fmt.Errorf("SLB API 调用失败 (账户=%s, 区域=%s): %w", accountName, regionID, err)
				}

				// 获取总数
				if totalCount == 0 {
					totalCount = int(response.TotalCount) // 获取总条数
					logger.Log.Infof("数据查询完成, 区域=%s, 资源=CLB, 账户=%s, 总数=%d 条", regionID, accountName, totalCount)
				}

				for _, lb := range response.LoadBalancers.LoadBalancer {
					// 构造 SLBRecord
					rec := database.SLBRecord{
						LoadBalancerID:   lb.LoadBalancerId,
						CloudName:        accountName,
						LoadBalancerName: lb.LoadBalancerName,
						IPAddress:        lb.Address,
						Bandwidth:        int64(lb.Bandwidth),
						NetworkType:      lb.NetworkType,
						RegionID:         lb.RegionId,
						Status:           lb.LoadBalancerStatus,
					}
					records = append(records, rec)
				}
				// 如果返回的数据条数小于 pageSize，说明已经拉取到最后一页，退出循环
				if len(response.LoadBalancers.LoadBalancer) < pageSize {
					break
				}

				// 请求下一页数据
				pageNumber++
			}
			// 保存 SLB 数据
			if err := database.SaveSLBRecords(records); err != nil {
				return fmt.Errorf("保存 SLB 数据失败 (账户=%s): %w", accountName, err)
			}

			logger.Log.Infof("数据同步完成, 区域=%s, 资源=CLB, 账户=%s, 同步=%d 条", regionID, accountName, len(records))
		} else {
			logger.Log.Warnf("当前阿里账户, 区域=%s, 资源=CLB, 账户=%s, 暂无可用区域。", regionID, accountName)
		}
	}
	return nil
}
