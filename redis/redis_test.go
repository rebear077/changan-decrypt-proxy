package redis

import (
	"context"
	"fmt"
	"testing"

	types "github.com/FISCO-BCOS/go-sdk/type"
	"github.com/redis/go-redis/v9"
)

func TestRedis(t *testing.T) {

	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "123456", // no password set
		DB:       0,        // use default DB
	})

	err := rdb.Set(ctx, "2023:1022", "dyy", 0).Err()
	if err != nil {
		panic(err)
	}
	err = rdb.Set(ctx, "2023:1023", "clyy", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "2023*").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key:", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}

func TestScan(t *testing.T) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "123456", // no password set
		DB:       0,        // use default DB
	})
	var cursor uint64
	var n int
	for {
		var keys []string
		var err error
		//*扫描所有key，每次20条
		keys, cursor, err = client.Scan(ctx, cursor, "*", 20).Result()
		if err != nil {
			panic(err)
		}
		n += len(keys)

		fmt.Printf("\nfound %d keys\n", n)
		var value []string
		for _, key := range keys {
			value, err = client.HKeys(ctx, key).Result()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("%v %v\n", key, value)
		}
		if cursor == 0 {
			break
		}
	}
}
func packInvoice() types.RawInvoiceInformation {
	title := new(types.RawInvoiceInformation)
	title.Certificateid = "2018210071"
	title.Customerid = "9527"
	title.Corpname = "米其林轮胎"
	title.Certificatetype = "中华人民共和国居民身份证"
	title.Intercustomerid = "YU225007"
	invoice := new(types.Invoiceinfos)
	invoice.Invoicenotaxamt = "960"
	invoice.Invoiceccy = "人民币"
	invoice.Sellername = "米其林轮胎"
	invoice.Invoicetype = "普通发票"
	invoice.Buyername = "长安汽车"
	invoice.Buyerusccode = "140c737693"
	invoice.Invoicedate = "2022-09-22"
	invoice.Sellerusccode = "150c78387"
	invoice.Invoicecode = "1289789309503-006"
	invoice.Invoicenum = "8234701951089023476"
	invoice.Checkcode = "082976"
	invoice.Invoiceamt = "1000"
	// invoice2 := new(receive.Invoiceinfos)
	// invoice2 = invoice
	title.Invoiceinfos = append(title.Invoiceinfos, *invoice)
	// title.Invoiceinfos = append(title.Invoiceinfos, *invoice2)
	return *title
}
func fliter() (string, map[string]interface{}) {
	invoices := packInvoice()
	key := invoices.Customerid + ":" + invoices.Invoiceinfos[0].Checkcode
	values := make(map[string]interface{})
	values["certificateId"] = invoices.Certificateid
	values["customerId"] = invoices.Certificateid
	values["corpName"] = invoices.Corpname
	values["certificateType"] = invoices.Certificatetype
	values["interCustomerId"] = invoices.Intercustomerid
	return key, values
}
func TestHset(t *testing.T) {
	key, values := fliter()
	ctx := context.Background()
	dber := NewRedisOperator(0)
	err := dber.MultipleSet(ctx, key, values)
	if err != nil {
		fmt.Println(err)
	}
}
func TestHGet(t *testing.T) {
	dber := NewRedisOperator(0)
	ctx := context.Background()
	res, err := dber.GetAll(ctx, "invoice")
	if err != nil {
		fmt.Println(err, "xxxxxx")
	}
	fmt.Println(res)
}
func TestDel(t *testing.T) {
	dber := NewRedisOperator(1)
	ctx := context.Background()
	res, _ := dber.rdb.Del(ctx, "001").Result()
	fmt.Println(res)
}
func TestFlush(t *testing.T) {
	dber := NewRedisOperator(0)
	ctx := context.Background()
	res, err := dber.rdb.FlushAll(ctx).Result()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}
func TestHMSet(t *testing.T) {
	dber := NewRedisOperator(1)
	ctx := context.Background()
	values := make(map[string]interface{})
	values["name"] = "dyy"
	values["age"] = "13"
	values["local"] = "1"
	err := dber.MultipleSet(ctx, "001", values)
	if err != nil {
		fmt.Println(err)
	}
}
