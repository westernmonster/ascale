package antispam

import (
	"context"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"

	"ascale/pkg/cache/redis"
	"ascale/pkg/net/http/vin"
	"ascale/pkg/xtime"

	"github.com/stretchr/testify/assert"
)

func TestAntiSpamHandler(t *testing.T) {
	anti := New(
		&Config{
			On:     true,
			Second: 1,
			N:      1,
			Hour:   1,
			M:      1,
			Redis: &redis.Config{
				MaxActive:    10,
				MaxIdle:      10,
				IdleTimeout:  xtime.Duration(time.Second * 60),
				Name:         "test",
				Proto:        "tcp",
				Addr:         "127.0.0.1:6379",
				DialTimeout:  xtime.Duration(time.Second),
				ReadTimeout:  xtime.Duration(time.Second),
				WriteTimeout: xtime.Duration(time.Second),
			},
		},
	)

	engine := vin.New()
	engine.Use(func(c *vin.Context) {
		uid, _ := strconv.ParseInt(c.Request.Form.Get("uid"), 10, 64)
		c.Set("uid", uid)
		c.Next()
	})
	engine.Use(anti.Handler())
	engine.GET("/antispam", func(c *vin.Context) {
		c.String(200, "pass")
	})
	go engine.Run(":18080")

	time.Sleep(time.Millisecond * 50)
	code, content, err := httpGet("http://127.0.0.1:18080/antispam?uid=11")
	if err != nil {
		t.Logf("http get failed, err:=%v", err)
		t.FailNow()
	}
	if code != 200 || string(content) != "pass" {
		t.Logf("request should pass by limiter, but blocked: %d, %v", code, content)
		t.FailNow()
	}

	_, content, err = httpGet("http://127.0.0.1:18080/antispam?uid=11")
	if err != nil {
		t.Logf("http get failed, err:=%v", err)
		t.FailNow()
	}
	if string(content) == "pass" {
		t.Logf("request should block by limiter, but passed")
		t.FailNow()
	}

	engine.Server().Shutdown(context.TODO())
}

func httpGet(url string) (code int, content []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	code = resp.StatusCode
	return
}

func TestConfigValidate(t *testing.T) {
	var conf *Config
	assert.Contains(t, conf.validate().Error(), "empty config")

	conf = &Config{
		Second: 0,
	}
	assert.Contains(t, conf.validate().Error(), "invalid Second")

	conf = &Config{
		Second: 1,
		Hour:   0,
	}
	assert.Contains(t, conf.validate().Error(), "invalid Hour")
}
