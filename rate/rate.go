package rate

import (
	"gopkg.in/redis.v3"
	"log"
	"net/http"
	"strconv"
	//	"strings"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

type rateHandler struct {
	handler http.Handler
}
type conf struct {
	Hits int64 `yaml:"hits"`
	Time int64 `yaml:"time"`
}

type userMD struct {
	meta []string
}

func RateHandler(h http.Handler) http.Handler {
	return rateHandler{h}
}

func (h rateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer client.Close()

	//in a real world case ill will probly  use r to get user data like app engine for exemple.
	u := &userMD{}
	u.meta = []string{"aaaa", "10.10.110.5"}

	b := checkRateOrBlock(client, r.RemoteAddr, u.meta)

	if b == false {
		w.WriteHeader(403)
		w.Write([]byte("DONT abuse the server"))
	} else {
		h.handler.ServeHTTP(w, r)

	}

}

func checkRateOrBlock(client *redis.Client, remotAdd string, meta []string) bool {

	c := &conf{}
	c = c.getConf()

	//Limit by parameters
	for _, e := range meta {
		ex, err := client.Exists(e).Result()
		if err != nil {
			log.Printf("client.Exists err   #%v ", err)
		}
		if ex {
			log.Printf("key inside the blacklist  therfor blocked")
			return false
		}
	}

	key := remotAdd + "_" + strconv.FormatInt(time.Now().Unix(), 10)

	//if exsit
	ex, err := client.Exists(key).Result()
	if err != nil {
		log.Printf("client.Exists err   #%v ", err)
	}
	//if key exsit increment and check if over the limit.
	if ex {

		i := client.Incr(key)
		log.Printf("key number   #%v ", i)
		_ = client.Expire(key, time.Millisecond*time.Duration(c.Time))

		if i.Val() > c.Hits {
			log.Printf("over the limit ")

			return false
		}
	} else { //key doesn't exsit insert it to redis
		log.Printf("set new key")
		_ = client.Set(key, "1", time.Minute*time.Duration(c.Time))
		return true
	}

	return true

}

//get the value from the yaml file
func (c *conf) getConf() *conf {

	yamlFile, err := ioutil.ReadFile("./conf.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
