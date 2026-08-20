# MySQL 统计实战：从基础到性能优化

## 一、统计的基石：核心函数
MySQL 提供丰富的聚合函数来满足统计需求：

1. ​**​COUNT​**​：统计行数
    SELECT COUNT(*) FROM orders; -- 统计总订单量

2. ​**​SUM/AVG​**​：数值计算
    SELECT SUM(amount) AS total_sales, AVG(amount) AS avg_price 
    FROM orders WHERE create_date > '2024-01-01';

3. ​**​MAX/MIN​**​：极值查询
    SELECT MAX(temperature), MIN(humidity) FROM sensor_data;

## 二、进阶统计：分组与过滤
通过 GROUP BY 实现多维统计：

    -- 按日期统计销售额
    SELECT DATE(create_time) AS day, 
           SUM(amount) AS daily_sales,
           COUNT(DISTINCT user_id) AS active_users 
    FROM orders 
    GROUP BY day 
    HAVING daily_sales > 10000;

注意：HAVING 用于分组后过滤，WHERE 用于分组前过滤

## 三、性能优化三板斧
### 1. 索引策略
- 为 WHERE/GROUP BY/ORDER BY 涉及的列创建复合索引
- 优先选择区分度高的字段作为索引前导列

### 2. 统计信息管理
    ANALYZE TABLE orders; -- 手动更新统计信息
    SHOW TABLE STATUS LIKE 'orders'; -- 查看估算值

建议在低峰期执行统计信息更新

## 四、典型应用场景
1. ​**​用户行为分析​**​  
       -- 统计7日留存率
       SELECT reg_date,
              COUNT(DISTINCT user_id) AS reg_users,
              COUNT(DISTINCT CASE WHEN login_date = reg_date + INTERVAL 7 DAY THEN user_id END)/COUNT(DISTINCT user_id) AS retention_rate
       FROM user_events 
       GROUP BY reg_date;


**引用说明**
: 基础聚合函数与应用场景
: 分组统计与函数详解
: 统计信息管理与性能优化
: 典型业务场景示例
: 高性能统计实现方案
