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

func NewRateHandler(handler http.Handler) *rateHandler {
	return &rateHandler{handler: handler}
}

type conf struct {
	hits int64 `yaml:"hits"`
	time int64 `yaml:"time"`
}

func (s *rateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	b := checkRate(r.RemoteAddr)
	if b == false {
		log.Printf("in the false   #%v ")
		w.WriteHeader(200)
	}
}

func checkRate(remotAdd string) bool {

	//shold i generate  new clinte?
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	c := &conf{5, 5}
	//	c = c.getConf()
	log.Printf("ccccccccccc   #%v ", c.hits)

	key := remotAdd + "_" + strconv.FormatInt(time.Now().Unix(), 10)

	currnet, err := client.Get(key).Result()
	if err != nil {
		log.Printf("client.Get err   #%v ", err)
	}

	if currnet == "" {
		//		log.Printf("in the if   #%v ")
		_ = client.Set(key, "1", time.Minute*10)
		//		log.Printf("in the if   #%v ", b)
		return false
	}

	i := client.Incr(key)
	log.Printf("in the else   #%v ", i)
	_ = client.Expire(key, time.Minute*10)

	if i.Val() > c.hits {
		log.Printf("out bbbbbbbbbbbbbbbbbbbbbbbbbbb")

		return false
	}

	return true

}

func (c *conf) getConf() *conf {

	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
