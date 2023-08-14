package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"time"
)

func main() {

	// 连接到Redis
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       0,
	})

	// 检查是否连接成功
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Redis:", pong)

	// 指定时间范围
	timeRanges := []struct {
		StartTime int64
		EndTime   int64
	}{
		{1690732800, 1690819200}, // 时间范围1: 2023-08-29 00:00:00 to 2023-08-30 00:00:00
		{1690819200, 1690905600}, // 时间范围2: 2023-08-30 00:00:00 to 2023-08-31 00:00:00
		{1690905600, 1690992000}, // 时间范围2: 2023-08-30 00:00:00 to 2023-08-31 00:00:00
		{1690992000, 1691078400}, // 时间范围2: 2023-08-30 00:00:00 to 2023-08-31 00:00:00

		{1691078400, 1691164800}, // 时间范围2: 2023-08-30 00:00:00 to 2023-08-31 00:00:00

		{1691164800, 1691251200}, // 时间范围2: 2023-08-30 00:00:00 to 2023-08-31 00:00:00

		{1691251200, 1691337600}, // 时间范围2: 2023-08-30 00:00:00 to 2023-08-31 00:00:00

		{1691337600, 1691424000}, // 时间范围2: 2023-08-30 00:00:00 to 2023-08-31 00:00:00

		{1691424000, 1691510400}, // 时间范围2: 2023-08-30 00:00:00 to 2023-08-31 00:00:00

		{1691510400, 1691596800}, // 时间范围2: 2023-08-30 00:00:00 to 2023-08-31 00:00:00

		{1691596800, 1691683200}, // 时间范围2: 2023-08-30 00:00:00 to 2023-08-31 00:00:00

		{1691683200, 1691769600}, // 时间范围2: 2023-08-30 00:00:00 to 2023-08-31 00:00:00

		{1691769600, 1691856000}, // 时间范围2: 2023-08-30 00:00:00 to 2023-08-31 00:00:00
		{1691856000, 1691942400}, // 时间范围2: 2023-08-30 00:00:00 to 2023-08-31 00:00:00
		{1691942400, 1692028800}, // 时间范围2: 2023-08-30 00:00:00 to 2023-08-31 00:00:00

		// 添加更多时间范围...
	}

	// 遍历每个时间范围
	for _, tr := range timeRanges {
		keyPattern := fmt.Sprintf("rs_*") // 修改为您的键名模式

		// 遍历符合模式的键
		keys, err := client.Keys(ctx, keyPattern).Result()
		if err != nil {
			log.Fatal(err)
		}

		// 统计每个时间范围内满足条件的记录数量
		count := 0
		for _, key := range keys {
			hgetallResult, err := client.HGetAll(ctx, key).Result()
			if err != nil {
				log.Fatal(err)
			}

			created_atStr, exists := hgetallResult["created_at"]
			if exists {
				created_at, _ := strconv.ParseInt(created_atStr, 10, 64)
				if created_at > tr.StartTime && created_at < tr.EndTime {
					count++
				}
			}
		}

		// 输出统计结果
		startTime := time.Unix(tr.StartTime, 0)
		endTime := time.Unix(tr.EndTime, 0)
		fmt.Printf("Records between %s and %s: %d\n", startTime.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05"), count)
	}
}
