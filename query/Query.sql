WITH search_key(keyword) AS (
  SELECT '%查询的内容%'
)

SELECT 
  'ecs' AS source_table,
  instance_id, cloud_name, instance_name, 
  public_ip, private_ip, remarks,
  '匹配字段: ' || 
    CASE 
      WHEN instance_id LIKE search_key.keyword THEN 'instance_id'
      WHEN cloud_name LIKE search_key.keyword THEN 'cloud_name'
      WHEN instance_name LIKE search_key.keyword THEN 'instance_name'
      WHEN public_ip LIKE search_key.keyword THEN 'public_ip'
      WHEN private_ip LIKE search_key.keyword THEN 'private_ip'
      WHEN remarks LIKE search_key.keyword THEN 'remarks'
    END AS match_info
FROM ecs, search_key
WHERE 
  instance_id LIKE search_key.keyword OR
  cloud_name LIKE search_key.keyword OR
  instance_name LIKE search_key.keyword OR
  public_ip LIKE search_key.keyword OR
  private_ip LIKE search_key.keyword OR
  remarks LIKE search_key.keyword

UNION ALL

SELECT 
  'rds',
  instance_id, cloud_name, instance_description, 
  NULL, connection_string, remarks,
  '匹配字段: ' || 
    CASE 
      WHEN instance_id LIKE search_key.keyword THEN 'instance_id'
      WHEN cloud_name LIKE search_key.keyword THEN 'cloud_name'
      WHEN instance_description LIKE search_key.keyword THEN 'instance_description'
      WHEN connection_string LIKE search_key.keyword THEN 'connection_string'
      WHEN remarks LIKE search_key.keyword THEN 'remarks'
    END
FROM rds, search_key
WHERE 
  instance_id LIKE search_key.keyword OR
  cloud_name LIKE search_key.keyword OR
  instance_description LIKE search_key.keyword OR
  connection_string LIKE search_key.keyword OR
  remarks LIKE search_key.keyword

UNION ALL

SELECT 
  'redis',
  instance_id, cloud_name, instance_name, 
  ip_address, connection_string, remarks,
  '匹配字段: ' || 
    CASE 
      WHEN instance_id LIKE search_key.keyword THEN 'instance_id'
      WHEN cloud_name LIKE search_key.keyword THEN 'cloud_name'
      WHEN instance_name LIKE search_key.keyword THEN 'instance_name'
      WHEN ip_address LIKE search_key.keyword THEN 'ip_address'
      WHEN connection_string LIKE search_key.keyword THEN 'connection_string'
      WHEN remarks LIKE search_key.keyword THEN 'remarks'
    END
FROM redis, search_key
WHERE 
  instance_id LIKE search_key.keyword OR
  cloud_name LIKE search_key.keyword OR
  instance_name LIKE search_key.keyword OR
  ip_address LIKE search_key.keyword OR
  connection_string LIKE search_key.keyword OR
  remarks LIKE search_key.keyword

UNION ALL

SELECT 
  'slb',
  lb_id, cloud_name, lb_name, 
  ip_address, NULL, remarks,
  '匹配字段: ' || 
    CASE 
      WHEN lb_id LIKE search_key.keyword THEN 'lb_id'
      WHEN cloud_name LIKE search_key.keyword THEN 'cloud_name'
      WHEN lb_name LIKE search_key.keyword THEN 'lb_name'
      WHEN ip_address LIKE search_key.keyword THEN 'ip_address'
      WHEN remarks LIKE search_key.keyword THEN 'remarks'
    END
FROM slb, search_key
WHERE 
  lb_id LIKE search_key.keyword OR
  cloud_name LIKE search_key.keyword OR
  lb_name LIKE search_key.keyword OR
  ip_address LIKE search_key.keyword OR
  remarks LIKE search_key.keyword

UNION ALL

SELECT 
  'polardb',
  dbcluster_id, cloud_name, dbcluster_description, 
  NULL, connection_string, remarks,
  '匹配字段: ' || 
    CASE 
      WHEN dbcluster_id LIKE search_key.keyword THEN 'dbcluster_id'
      WHEN cloud_name LIKE search_key.keyword THEN 'cloud_name'
      WHEN dbcluster_description LIKE search_key.keyword THEN 'dbcluster_description'
      WHEN connection_string LIKE search_key.keyword THEN 'connection_string'
      WHEN remarks LIKE search_key.keyword THEN 'remarks'
    END
FROM polardb, search_key
WHERE 
  dbcluster_id LIKE search_key.keyword OR
  cloud_name LIKE search_key.keyword OR
  dbcluster_description LIKE search_key.keyword OR
  connection_string LIKE search_key.keyword OR
  remarks LIKE search_key.keyword;

