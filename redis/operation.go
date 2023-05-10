package redis

import (
	"context"
	"fmt"

	"github.com/FISCO-BCOS/go-sdk/conf"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

const (
	Increase = "increase"
	Decrease = "decrease"
)

type RedisOperator struct {
	rdb *redis.Client
}
type Rediser interface {
	Set(ctx context.Context, key string, values ...interface{}) error
	Get(ctx context.Context, key, field string) (string, error)
	GetAll(ctx context.Context, key string) (map[string]string, error)
	GetKeys(ctx context.Context, key string) ([]string, error)
	GetHashLen(ctx context.Context, key string) (int64, error)
	MultipleGet(ctx context.Context, key string, fields ...string) ([]interface{}, error)
	MultipleSet(ctx context.Context, key string, values map[string]interface{}) error
	Delete(ctx context.Context, key string, fields ...string) error
	Exists(ctx context.Context, key, field string) (bool, error)
	Scan(ctx context.Context, condition string) int
}

// 初始化一个redis客户端
func NewRedisOperator(db int) *RedisOperator {
	configs, err := conf.ParseConfigFile("./configs/config.toml")
	if err != nil {
		logrus.Fatalln(err)
	}
	config := &configs[0]
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisUrl,
		Password: config.RedisPassword, // no password set
		DB:       db,                   // use default DB
	})
	return &RedisOperator{
		rdb: rdb,
	}
}

// 调用Hset方法，根据key和field字段设置，field字段的值(单个值)
// HSet accepts values in following formats:
//
//   - HSet("myhash", "key1", "value1", "key2", "value2")
//
//   - HSet("myhash", []string{"key1", "value1", "key2", "value2"})
//
//   - HSet("myhash", map[string]interface{}{"key1": "value1", "key2": "value2"})
//
//     Playing struct With "redis" tag.
//     type MyHash struct { Key1 string `redis:"key1"`; Key2 int `redis:"key2"` }
//
//   - HSet("myhash", MyHash{"value1", "value2"})
//
//     For struct, can be a structure pointer type, we only parse the field whose tag is redis.
//     if you don't want the field to be read, you can use the `redis:"-"` flag to ignore it,
//     or you don't need to set the redis tag.
//     For the type of structure field, we only support simple data types:
//     string, int/uint(8,16,32,64), float(32,64), time.Time(to RFC3339Nano), time.Duration(to Nanoseconds ),
//     if you are other more complex or custom data types, please implement the encoding.BinaryMarshaler interface.
//
// Note that it requires Redis v4 for multiple field/value pairs support.
func (operator *RedisOperator) Set(ctx context.Context, key string, values ...interface{}) error {
	err := operator.rdb.HSet(ctx, key, values).Err()

	return err
}

// 调用HGet方法，根据key和field字段，查询field字段的值
func (operator *RedisOperator) Get(ctx context.Context, key, field string) (string, error) {
	data, err := operator.rdb.HGet(ctx, key, field).Result()
	switch {
	case err != nil:
		return "", err
	default:
		return data, nil
	}
}

// HGetAll - 根据key查询所有字段和值
func (operator *RedisOperator) GetAll(ctx context.Context, key string) (map[string]string, error) {
	data, err := operator.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// HKeys-根据key返回所有字段名
func (operator *RedisOperator) GetKeys(ctx context.Context, key string) ([]string, error) {
	data, err := operator.rdb.HKeys(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// HLen-根据key，查询hash的字段数量
func (operator *RedisOperator) GetHashLen(ctx context.Context, key string) (int64, error) {
	len, err := operator.rdb.HLen(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return len, nil
}

// HMGet-根据key和多个字段名，批量查询多个hash字段值
func (operator *RedisOperator) MultipleGet(ctx context.Context, key string, fields ...string) ([]interface{}, error) {
	data, err := operator.rdb.HMGet(ctx, key, fields...).Result()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// HMSet-根据key和多个字段名和字段值，批量设置hash字段值
func (operator *RedisOperator) MultipleSet(ctx context.Context, key string, values map[string]interface{}) error {
	// fmt.Println(values)
	err := operator.rdb.HMSet(ctx, key, values).Err()
	return err
}

// HDel-根据key和字段名，删除hash字段，支持批量删除hash字段
func (operator *RedisOperator) Delete(ctx context.Context, key string, fields ...string) error {
	err := operator.rdb.HDel(ctx, key, fields...).Err()
	return err

}

// HExists-检测hash字段名是否存在。
func (operator *RedisOperator) Exists(ctx context.Context, key, field string) (bool, error) {
	res, err := operator.rdb.HExists(ctx, key, field).Result()
	return res, err
}

// 扫描当前redis数据库中包含的主键数量
func (operator *RedisOperator) Scan(ctx context.Context, condition string) (int, []string) {
	var cursor uint64
	var n int
	var keys []string
	for {
		var key []string
		var err error
		//*扫描所有key，每次20条
		key, cursor, err = operator.rdb.Scan(ctx, cursor, condition, 10).Result()
		if err != nil {
			panic(err)
		}
		n += len(keys)
		keys = append(keys, key...)
		if cursor == 0 {
			break
		}
	}
	fmt.Println(keys)
	return n, keys
}
func (operator *RedisOperator) FlushData(ctx context.Context) {

	res, err := operator.rdb.FlushAll(ctx).Result()
	if err != nil {
		logrus.Errorln(err)
		return
	}
	logrus.Infoln("清除redis所有数据:", res)
}

// Redis Zadd 命令用于将一个或多个成员元素及其分数值加入到有序集当中。
func (operator *RedisOperator) ZAdd(ctx context.Context, key string, member map[string]float64) {
	for ele, score := range member {
		z := redis.Z{
			Score:  score,
			Member: ele,
		}
		_, err := operator.rdb.ZAdd(ctx, key, z).Result()
		if err != nil {
			logrus.Errorln("err", err)
		}
	}

}

// Redis Zcard 命令用于计算集合中元素的数量,当 key 存在且是有序集类型时，返回有序集的基数。 当 key 不存在时，返回 0
// 如果查询出错，则返回-1
func (operator *RedisOperator) Zcard(ctx context.Context, key string) int64 {
	res, err := operator.rdb.ZCard(ctx, key).Result()
	if err != nil {
		logrus.Errorln(err)
		return -1
	}
	return res
}

// Redis Zcount 命令用于计算有序集合中指定分数区间的成员数量。
// 如果出错，则返回-1
func (operator *RedisOperator) Zcount(ctx context.Context, key, min, max string) int64 {
	res, err := operator.rdb.ZCount(ctx, key, min, max).Result()
	if err != nil {
		logrus.Errorln(err)
		return -1
	}
	return res
}

// Redis Zrange 返回有序集中，指定区间内的成员。
// 其中成员的位置按分数值递增(从小到大)来排序。
// 具有相同分数值的成员按字典序(lexicographical order )来排列。
// 如果你需要成员按
// 值递减(从大到小)来排列，请使用 ZREVRANGE 命令。
// 下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。
// 你也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。
func (operator *RedisOperator) ZrangeIncrease(ctx context.Context, key string, start, end int64) []string {
	res, err := operator.rdb.ZRange(ctx, key, start, end).Result()
	if err != nil {
		logrus.Errorln(err)
		return nil
	}
	return res
}
func (operator *RedisOperator) ZrangeDecrease(ctx context.Context, key string, start, end int64) []string {
	res, err := operator.rdb.ZRevRange(ctx, key, start, end).Result()
	if err != nil {
		logrus.Errorln(err)
		return nil
	}
	return res
}
func (operator *RedisOperator) Del(ctx context.Context, keys ...string) {
	operator.rdb.Del(ctx, keys...)
}

// // 返回的每个元素都是一个有序集合元素，一个有序集合元素由一个成员（member）和一个分值（score）组成。
// func (operator *RedisOperator) Zscan(ctx context.Context, key string, cursor uint64, count int64) {
// 	operator.rdb.ZScan(ctx, key, cursor, count)
// }
