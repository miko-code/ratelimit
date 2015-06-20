package rate

import (
	"gopkg.in/redis.v3"
	"testing"
	"time"
)

func TestCheckRate(t *testing.T) {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer client.Close()
	u := &userMD{}
	u.meta = []string{}
	//case of over  a normal call
	b := checkRateOrBlock(client, "10.15.10", u.meta)

	if b == false {
		t.Log("you are over the rate limit ")

	} else {
		t.Log("not over the limit")
	}

	//case of over the limit
	for b != false {
		b = checkRateOrBlock(client, "10.15.10", u.meta)
	}
	if b == false {
		t.Log("you are over the rate limit ")

	} else {
		t.Log("not over the limit")
	}
}

func TestCheckMd(t *testing.T) {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer client.Close()
	client.Set("miki", "miki", time.Minute*1)
	u := &userMD{}
	u.meta = []string{"miki", "10.10.10.5"}
	b := checkRateOrBlock(client, "10.15.10", u.meta)

	if b == false {
		t.Log("key blocked")

	} else {
		t.Log("key not blocked")
	}
}
