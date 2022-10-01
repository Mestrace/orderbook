package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Mestrace/orderbook/domain/dao"
	"github.com/Mestrace/orderbook/domain/model"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/go-redis/redis/v9"
)

const (
	templateOrderBookSymbolList    = `ob:%s:sl`
	templateOrderBookSymbolItemAsk = `ob:%s:si:%s:a`
	templateOrderBookSymbolItemBid = `ob:%s:si:%s:b`
)

func OrderBookWithRedis(orderBookDAO dao.ExchangeOrderBook, client *redis.Client,
	symbolListCacheDuration, symbolCacheDuration time.Duration) dao.ExchangeOrderBook {
	return &orderbookWithRedis{
		redisClient:             client,
		ExchangeOrderBook:       orderBookDAO,
		symbolListCacheDuration: symbolListCacheDuration,
		symbolCacheDuration:     symbolCacheDuration,
	}
}

type orderbookWithRedis struct {
	redisClient *redis.Client
	dao.ExchangeOrderBook
	symbolListCacheDuration time.Duration
	symbolCacheDuration     time.Duration
}

func (o *orderbookWithRedis) GetSymbolList(
	ctx context.Context, param *model.GetSymbolListParams,
) (*model.GetSymbolListData, error) {
	var (
		result = &model.GetSymbolListData{}
	)

	key := fmt.Sprintf(templateOrderBookSymbolList, param.ExchangeName)

	if err := o.get(ctx, key, result); err != nil {
		logger.CtxWarnf(ctx, "redis_get_failed|key=%s|err=%+v", key, err)
	}

	result, err := o.ExchangeOrderBook.GetSymbolList(ctx, param)
	if err != nil {
		return nil, err
	}

	if err := o.set(ctx, key, result, o.symbolListCacheDuration); err != nil {
		logger.CtxWarnf(ctx, "redis_set_failed|key=%s|err=%+v", key, err)
	}

	return result, nil
}

func (o *orderbookWithRedis) GetSymbolPrice(
	ctx context.Context, param *model.GetSymbolPriceParams,
) (*model.GetSymbolPriceData, error) {
	var (
		askKey   = fmt.Sprintf(templateOrderBookSymbolItemAsk, param.ExchangeName, param.Symbol)
		bidKey   = fmt.Sprintf(templateOrderBookSymbolItemBid, param.ExchangeName, param.Symbol)
		askStat  *model.SymbolStat
		bidStat  *model.SymbolStat
		redisErr error
	)

	// get from redis
	if param.OrderType == model.OrderTypeAll || param.OrderType == model.OrderTypeAsks {
		askStat = &model.SymbolStat{}
		redisErr = o.get(ctx, askKey, askStat)
		if redisErr != nil {
			logger.CtxWarnf(ctx, "redis_get_failed|key=%s|err=%+v", askKey, redisErr)
		}
	}

	if redisErr == nil && (param.OrderType == model.OrderTypeAll || param.OrderType == model.OrderTypeBids) {
		bidStat = &model.SymbolStat{}
		redisErr = o.get(ctx, bidKey, bidStat)
		if redisErr != nil {
			logger.CtxWarnf(ctx, "redis_get_failed|key=%s|err=%+v", bidKey, redisErr)
		}
	}

	// if redis found
	if redisErr == nil {
		return &model.GetSymbolPriceData{
			Ask: askStat,
			Bid: bidStat,
		}, nil
	}
	redisErr = nil

	// get from api
	result, err := o.ExchangeOrderBook.GetSymbolPrice(ctx, param)
	if err != nil {
		return nil, err
	}

	// set redis
	if param.OrderType == model.OrderTypeAll || param.OrderType == model.OrderTypeAsks {
		redisErr = o.set(ctx, askKey, result.Ask, o.symbolCacheDuration)
		if redisErr != nil {
			logger.CtxWarnf(ctx, "redis_set_failed|key=%s|err=%+v", askKey, redisErr)
		}
	}

	if redisErr == nil && (param.OrderType == model.OrderTypeAll || param.OrderType == model.OrderTypeBids) {
		redisErr = o.set(ctx, bidKey, result.Bid, o.symbolCacheDuration)
		if redisErr != nil {
			logger.CtxWarnf(ctx, "redis_set_failed|key=%s|err=%+v", askKey, redisErr)
		}
	}

	return result, nil
}

func (o *orderbookWithRedis) get(ctx context.Context, key string, out interface{}) error {
	cmd := o.redisClient.Get(ctx, key)
	if err := cmd.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			logger.CtxWarnf(ctx, "redis_not_found|key=%s", key)
			return err
		}
		logger.CtxWarnf(ctx, "redis_get_key_failed|key=%s|err=%+v", key, err)
		return err
	}

	b, err := cmd.Bytes()
	if err != nil {
		logger.CtxWarnf(ctx, "parse_bytes_failed|key=%s|err=%+v", key, err)
		return err
	}

	if err := json.Unmarshal(b, out); err != nil {
		logger.CtxWarnf(ctx, "unmarshal_failed|key=%s|err=%+v", key, err)
		return err
	}

	return nil
}

func (o *orderbookWithRedis) set(ctx context.Context, key string, in interface{}, exp time.Duration) error {
	b, err := json.Marshal(in)
	if err != nil {
		logger.CtxWarnf(ctx, "marshal_failed|key=%s|err=%+v", key, err)
		return err
	}

	cmd := o.redisClient.SetEx(ctx, key, b, exp)
	if err := cmd.Err(); err != nil {
		logger.CtxWarnf(ctx, "redis_set_ex_failed|key=%s|err=%+v", key, err)
		return err
	}

	return nil
}
