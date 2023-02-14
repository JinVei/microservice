package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	sessionKeyFormat        = "micro:uid:%s:sid:%s"
	userSessionSetkeyFormat = "micro:uid:%s:sessions"
)

type UserSession struct {
	rdb         *redis.Client
	sessDura    time.Duration
	maxSesssion int64
}

func NewUserSession(rbd *redis.Client, sessDura time.Duration, maxSesssion int64) *UserSession {
	return &UserSession{
		rdb:         rbd,
		sessDura:    sessDura,
		maxSesssion: maxSesssion,
	}
}

func (s *UserSession) AddSession(ctx context.Context, userID, sid string) error {
	sk := fmt.Sprintf(sessionKeyFormat, userID, sid)
	ssetK := fmt.Sprintf(userSessionSetkeyFormat, userID)

	// limit Session number to maxSesssion
	lres := s.rdb.LLen(ctx, ssetK)
	if len, err := lres.Result(); err == nil && s.maxSesssion <= len {
		dkey := s.rdb.LPop(ctx, ssetK)
		if dkey.Val() != "" {
			k := fmt.Sprintf(sessionKeyFormat, userID, dkey.Val())
			s.rdb.Del(ctx, k)
		}
	}

	// TODO: store entity.Session to redis
	// ss := entity.Session{
	// 	UserID:     userID,
	// 	SessionId:  sessionid,
	// 	LastUpdate: strconv.Itoa(int(time.Now().Unix())),
	// }
	shm := s.rdb.SetEx(ctx, sk, strconv.Itoa(int(time.Now().Unix())), s.sessDura)
	if shm.Err() != nil {
		return shm.Err()
	}

	s.rdb.RPush(ctx, ssetK, sid)
	s.rdb.Expire(ctx, ssetK, s.sessDura)

	return nil
}

func (s *UserSession) DelUserSession(ctx context.Context, uid, sid string) error {
	sk := fmt.Sprintf(sessionKeyFormat, uid, sid)
	ssetK := fmt.Sprintf(userSessionSetkeyFormat, uid)
	if res := s.rdb.Del(ctx, sk); res.Err() != nil {
		//TODO: log
	}
	if res := s.rdb.LRem(ctx, ssetK, 0, sid); res.Err() != nil {
		//TODO: log
	}
	return nil
}
