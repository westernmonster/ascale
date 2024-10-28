package dao

import (
	"ascale/app/api/conf"
	"ascale/pkg/cache/redis"
	"ascale/pkg/database/sqalx"
	"ascale/pkg/log"
	"ascale/pkg/stat/prom"
	"context"
	"fmt"
)

type Dao struct {
	db            sqalx.Node
	c             *conf.Config
	redis         *redis.Pool
	redisExpire   int32
	valcodeExpire int32
}

func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:             c,
		db:            sqalx.NewMySQL(c.DB),
		redis:         redis.NewPool(c.Redis),
		redisExpire:   int32(5 * 60),
		valcodeExpire: int32(2 * 60),
	}

	return
}

func (d *Dao) DB() sqalx.Node {
	return d.db
}

func (d *Dao) Redis() *redis.Pool {
	return d.redis
}

// Ping check db and mc health.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.db.Ping(c); err != nil {
		log.Info(fmt.Sprintf("dao.db.Ping() error(%v)", err))
	}

	// if err = d.pingRedis(c); err != nil {
	// 	return
	// }

	return
}

// Close close connections of mc, redis, db.
func (d *Dao) Close(ctx context.Context) {
	if d.db != nil {
		d.db.Close()
	}

	if d.redis != nil {
		d.redis.Close()
	}
}

// PromError prometheus error count.
func PromError(c context.Context, name, format string, args ...interface{}) {
	prom.BusinessErrCount.Incr(name)
	log.For(c).Error(fmt.Sprintf(format, args...))
}
