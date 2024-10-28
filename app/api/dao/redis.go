package dao

import (
	"ascale/pkg/cache/redis"
	"ascale/pkg/log"
	"context"
)

func (d *Dao) pingRedis(c context.Context) (err error) {
	var conn redis.Conn
	conn, err = d.redis.GetContext(c)
	if err != nil {
		return
	}
	defer conn.Close()
	_, err = conn.Do("SET", "PING", "PONG")
	return
}

func (p *Dao) FlushRedis(ctx context.Context) (err error) {
	var conn redis.Conn
	if conn, err = p.redis.GetContext(ctx); err != nil {
		log.For(ctx).Errorf("dao.FlushDB(), err(%+v)", err)
		return
	}

	defer conn.Close()

	if err = conn.Send("FLUSHDB"); err != nil {
		log.For(ctx).Errorf("dao.FlushDB(), err(%+v)", err)
		return
	}

	return
}
