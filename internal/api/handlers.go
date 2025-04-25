package api

import (
	"strconv"
	"strings"

	"github.com/WillemCode/AliCloud_Resources/pkg/database"
	"github.com/WillemCode/AliCloud_Resources/pkg/logger"
	"github.com/gin-gonic/gin"
)

// 分页响应结构
type PaginatedResponse struct {
	Data     interface{} `json:"data"`
	Total    int         `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

// 设置路由
func SetupRoutes(router *gin.Engine) {
	// 添加 CORS 中间件
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// API 路由
	router.GET("/ecs", handleECSList)
	router.GET("/rds", handleRDSList)
	router.GET("/slb", handleSLBList)
	router.GET("/redis", handleRedisList)
	router.GET("/polardb", handlePolarDBList)
	router.GET("/search", handleSearch)
}

// 处理 ECS 列表请求
func handleECSList(c *gin.Context) {
	page, pageSize := getPaginationParams(c)

	ecsRecords, err := database.ListECSRecords()
	if err != nil {
		logger.Log.Error("查询 ECS 数据失败: ", err)
		c.JSON(500, gin.H{"error": "failed to query ECS data"})
		return
	}

	// 应用分页
	total := len(ecsRecords)
	paginatedData := applyPagination(ecsRecords, page, pageSize)

	c.JSON(200, PaginatedResponse{
		Data:     paginatedData,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// 处理 RDS 列表请求
func handleRDSList(c *gin.Context) {
	page, pageSize := getPaginationParams(c)

	rdsRecords, err := database.ListRDSRecords()
	if err != nil {
		logger.Log.Error("查询 RDS 数据失败: ", err)
		c.JSON(500, gin.H{"error": "failed to query RDS data"})
		return
	}

	// 应用分页
	total := len(rdsRecords)
	paginatedData := applyPagination(rdsRecords, page, pageSize)

	c.JSON(200, PaginatedResponse{
		Data:     paginatedData,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// 处理 SLB 列表请求
func handleSLBList(c *gin.Context) {
	page, pageSize := getPaginationParams(c)

	slbRecords, err := database.ListSLBRecords()
	if err != nil {
		logger.Log.Error("查询 SLB 数据失败: ", err)
		c.JSON(500, gin.H{"error": "failed to query SLB data"})
		return
	}

	// 应用分页
	total := len(slbRecords)
	paginatedData := applyPagination(slbRecords, page, pageSize)

	c.JSON(200, PaginatedResponse{
		Data:     paginatedData,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// 处理 Tair Redis 列表请求
func handleRedisList(c *gin.Context) {
	page, pageSize := getPaginationParams(c)

	RedisRecords, err := database.ListRedisRecords()
	if err != nil {
		logger.Log.Error("查询 Tair Redis 数据失败: ", err)
		c.JSON(500, gin.H{"error": "failed to query SLB data"})
		return
	}

	// 应用分页
	total := len(RedisRecords)
	paginatedData := applyPagination(RedisRecords, page, pageSize)

	c.JSON(200, PaginatedResponse{
		Data:     paginatedData,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// 处理 PolarDB 列表请求
func handlePolarDBList(c *gin.Context) {
	page, pageSize := getPaginationParams(c)

	polarDBRecords, err := database.ListPolarDBRecords()
	if err != nil {
		logger.Log.Error("查询 PolarDB 数据失败: ", err)
		c.JSON(500, gin.H{"error": "failed to query PolarDB data"})
		return
	}

	// 应用分页
	total := len(polarDBRecords)
	paginatedData := applyPagination(polarDBRecords, page, pageSize)

	c.JSON(200, PaginatedResponse{
		Data:     paginatedData,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// 处理搜索请求
func handleSearch(c *gin.Context) {
	keyword := c.Query("q")
	resourceType := c.DefaultQuery("type", "all")

	if keyword == "" {
		c.JSON(400, gin.H{"error": "search keyword is required"})
		return
	}

	// 转换关键词为小写以进行不区分大小写的搜索
	keyword = strings.ToLower(keyword)

	var results []interface{}

	// 根据资源类型搜索不同的表
	if resourceType == "all" || resourceType == "ecs" {
		ecsRecords, err := database.ListECSRecords()
		if err == nil {
			for _, record := range ecsRecords {
				if containsKeyword(record, keyword) {
					results = append(results, record)
				}
			}
		}
	}

	if resourceType == "all" || resourceType == "rds" {
		rdsRecords, err := database.ListRDSRecords()
		if err == nil {
			for _, record := range rdsRecords {
				if containsKeyword(record, keyword) {
					results = append(results, record)
				}
			}
		}
	}

	if resourceType == "all" || resourceType == "slb" {
		slbRecords, err := database.ListSLBRecords()
		if err == nil {
			for _, record := range slbRecords {
				if containsKeyword(record, keyword) {
					results = append(results, record)
				}
			}
		}
	}

	if resourceType == "all" || resourceType == "polardb" {
		polarDBRecords, err := database.ListPolarDBRecords()
		if err == nil {
			for _, record := range polarDBRecords {
				if containsKeyword(record, keyword) {
					results = append(results, record)
				}
			}
		}
	}

	page, pageSize := getPaginationParams(c)
	total := len(results)
	paginatedResults := applyPagination(results, page, pageSize)

	c.JSON(200, PaginatedResponse{
		Data:     paginatedResults,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// 获取分页参数
func getPaginationParams(c *gin.Context) (int, int) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// 限制最大页面大小
	if pageSize > 50 {
		pageSize = 50
	}

	return page, pageSize
}

// 应用分页到任意切片
func applyPagination(data interface{}, page, pageSize int) interface{} {
	switch v := data.(type) {
	case []database.ECSRecord:
		return paginateSlice(v, page, pageSize)
	case []database.RDSRecord:
		return paginateSlice(v, page, pageSize)
	case []database.SLBRecord:
		return paginateSlice(v, page, pageSize)
	case []database.PolarDBRecord:
		return paginateSlice(v, page, pageSize)
	case []interface{}:
		return paginateSlice(v, page, pageSize)
	default:
		return data
	}
}

// 对任意切片进行分页
func paginateSlice(slice interface{}, page, pageSize int) interface{} {
	switch v := slice.(type) {
	case []database.ECSRecord:
		start, end := calculatePaginationBounds(len(v), page, pageSize)
		if start >= len(v) {
			return []database.ECSRecord{}
		}
		return v[start:end]
	case []database.RDSRecord:
		start, end := calculatePaginationBounds(len(v), page, pageSize)
		if start >= len(v) {
			return []database.RDSRecord{}
		}
		return v[start:end]
	case []database.SLBRecord:
		start, end := calculatePaginationBounds(len(v), page, pageSize)
		if start >= len(v) {
			return []database.SLBRecord{}
		}
		return v[start:end]
	case []database.PolarDBRecord:
		start, end := calculatePaginationBounds(len(v), page, pageSize)
		if start >= len(v) {
			return []database.PolarDBRecord{}
		}
		return v[start:end]
	case []interface{}:
		start, end := calculatePaginationBounds(len(v), page, pageSize)
		if start >= len(v) {
			return []interface{}{}
		}
		return v[start:end]
	default:
		return slice
	}
}

// 计算分页的起始和结束索引
func calculatePaginationBounds(total, page, pageSize int) (int, int) {
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}
	return start, end
}

// 检查记录是否包含关键词
func containsKeyword(record interface{}, keyword string) bool {
	switch v := record.(type) {
	case database.ECSRecord:
		return strings.Contains(strings.ToLower(v.InstanceID), keyword) ||
			strings.Contains(strings.ToLower(v.CloudName), keyword) ||
			strings.Contains(strings.ToLower(v.InstanceName), keyword) ||
			strings.Contains(strings.ToLower(v.PublicIP), keyword) ||
			strings.Contains(strings.ToLower(v.OSName), keyword) ||
			strings.Contains(strings.ToLower(v.PrivateIP), keyword) ||
			strings.Contains(strings.ToLower(v.RegionID), keyword)
	case database.RDSRecord:
		return strings.Contains(strings.ToLower(v.InstanceID), keyword) ||
			strings.Contains(strings.ToLower(v.Engine), keyword) ||
			strings.Contains(strings.ToLower(v.ConnectionString), keyword) ||
			strings.Contains(strings.ToLower(v.RegionID), keyword) ||
			strings.Contains(strings.ToLower(v.CloudName), keyword)
	case database.SLBRecord:
		return strings.Contains(strings.ToLower(v.InstanceID), keyword) ||
			strings.Contains(strings.ToLower(v.LoadBalancerName), keyword) ||
			strings.Contains(strings.ToLower(v.IPAddress), keyword) ||
			strings.Contains(strings.ToLower(v.RegionID), keyword) ||
			strings.Contains(strings.ToLower(v.CloudName), keyword)
	case database.PolarDBRecord:
		return strings.Contains(strings.ToLower(v.InstanceID), keyword) ||
			strings.Contains(strings.ToLower(v.Status), keyword) ||
			strings.Contains(strings.ToLower(v.Engine), keyword) ||
			strings.Contains(strings.ToLower(v.ConnectionString), keyword) ||
			strings.Contains(strings.ToLower(v.RegionID), keyword) ||
			strings.Contains(strings.ToLower(v.CloudName), keyword)
	default:
		return false
	}
}
