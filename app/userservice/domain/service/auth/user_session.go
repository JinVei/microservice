package auth

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jinvei/microservice/base/api/proto/v1/dto"
	"github.com/redis/go-redis/v9"
)

const (
	sessionKeyFormat        = "micro:uid:%s:sid:%s"
	userSessionSetkeyFormat = "micro:uid:%s:sessions"
	emVerifyCodeKeyFormat   = "micro:auth:email:%s:code"
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
	} else if err != nil {
		return err
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

func (s *UserSession) DelUserSession(ctx context.Context, uid, sid string) *dto.Status {
	sk := fmt.Sprintf(sessionKeyFormat, uid, sid)
	ssetK := fmt.Sprintf(userSessionSetkeyFormat, uid)
	if res := s.rdb.Del(ctx, sk); res.Err() != nil {
		//TODO: log
		flog.Warn("rdb.Del", "res.Err()", res.Err())
	}
	if res := s.rdb.LRem(ctx, ssetK, 0, sid); res.Err() != nil {
		//TODO: log
		flog.Warn("rdb.LRem", "res.Err()", res.Err())
	}
	return nil
}

func (s *UserSession) PutVerifyCode(ctx context.Context, email, code string) error {
	key := fmt.Sprintf(emVerifyCodeKeyFormat, email)
	res := s.rdb.Set(ctx, key, code, 60*time.Second)
	return res.Err()
}

func (s *UserSession) GetVerifyCode(ctx context.Context, email string) (string, bool) {
	key := fmt.Sprintf(emVerifyCodeKeyFormat, email)
	res := s.rdb.Get(ctx, key)
	if res.Err() != nil {
		if res.Err() != redis.Nil {
			flog.Error(res.Err(), "rdb.Get()")
		}
		return "", false
	}
	return res.Val(), true
}
