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

func RateHandler(h http.Handler) http.Handler {
	return rateHandler{h}
}

func (h rateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b := checkRate(r.RemoteAddr)
	if b == false {
		log.Printf("in the false  ")
		h.handler.ServeHTTP(w, r)
	}

}

func checkRate(remotAdd string) bool {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer client.Close()

	c := &conf{}
	c = c.getConf()

	key := remotAdd + "_" + strconv.FormatInt(time.Now().Unix(), 10)

	//if exsit
	ex, err := client.Exists(key).Result()
	if err != nil {
		log.Printf("client.Get err   #%v ", err)
	}

	if ex {

		i := client.Incr(key)
		log.Printf("in the else   #%v ", i)
		_ = client.Expire(key, time.Millisecond*time.Duration(c.Time))

		if i.Val() > c.Hits {
			log.Printf("out bbbbbbbbbbbbbbbbbbbbbbbbbbb")

			return false
		}
	} else {
		log.Printf("setttttt")
		_ = client.Set(key, "1", time.Minute*time.Duration(c.Time))
		return true
	}

	return true

}

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
