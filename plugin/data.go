package plugin

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/rpc"
	"github.com/jqiris/kungfu/v2/stores"
	"github.com/jqiris/kungfu/v2/utils"
	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type DataItem struct {
	UserId int64 `json:"user_id"`
	MsgId  int32 `json:"msg_id"`
	MsgReq any   `json:"msg_req"`
}

type DataHandler func(item *DataItem)

type ServerData struct {
	dbStore     stores.StoreKeeper
	dataHandler DataHandler
	dealChanArr []chan *DataItem
	maxDbQueue  int64
	dbQueueSize int64
	cancelCtx   context.Context
	cancelFunc  context.CancelFunc
}

func NewServerData(dbStore stores.StoreKeeper, dataHandler DataHandler, maxDbQueue, dbQueueSize int64) *ServerData {
	if maxDbQueue < 1 {
		maxDbQueue = 10
	}
	if dbQueueSize < 1 {
		dbQueueSize = 1000
	}
	s := &ServerData{
		dbStore:     dbStore,
		dataHandler: dataHandler,
		dealChanArr: make([]chan *DataItem, maxDbQueue),
		maxDbQueue:  maxDbQueue,
		dbQueueSize: dbQueueSize,
	}
	s.cancelCtx, s.cancelFunc = context.WithCancel(context.Background())
	return s

}

func (b *ServerData) DealReq(ctx context.Context, s *rpc.ServerBase) {
	for {
		select {
		case <-ctx.Done():
			logger.Infof("DealReq receive cancel signal")
			return
		default:
			var item *DataItem
			if err := b.dbStore.BRPop(s.Server.ServerId, &item); err == nil {
				b.dealChanArr[item.UserId%b.maxDbQueue] <- item
			} else {
				if !stores.IsRedisNull(err) {
					logger.Error(err)
				}
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func (b *ServerData) DealQueue(ctx context.Context, queue chan *DataItem) {
	for {
		select {
		case item := <-queue:
			b.dataHandler(item)
		case <-ctx.Done():
			logger.Infof("DealQueue receive cancel signal")
			return
		}
	}
}

func (b *ServerData) Init(s *rpc.ServerBase) {
}

func (b *ServerData) AfterInit(s *rpc.ServerBase) {
	for k, v := range b.dealChanArr {
		v = make(chan *DataItem, b.dbQueueSize)
		b.dealChanArr[k] = v
		go utils.SafeRun(func() {
			b.DealQueue(b.cancelCtx, v)
		})
	}
	//读取队列
	go utils.SafeRun(func() {
		b.DealReq(b.cancelCtx, s)
	})
}

func (b *ServerData) BeforeShutdown(s *rpc.ServerBase) {
	b.cancelFunc()
	//剩余队列资源处理
	wg := sync.WaitGroup{}
	for dealKey, dealChan := range b.dealChanArr {
		wg.Add(1)
		go func(k int, v chan *DataItem) {
			defer wg.Done()
			logger.Infof("ServerData Shutting down deal chan:%v", k)
			for item := range v {
				b.dataHandler(item)
			}
		}(dealKey, dealChan)
	}
	wg.Wait()
	fmt.Println("ServerData Shutting down deal chan over")
}

func (b *ServerData) Shutdown(s *rpc.ServerBase) {
}
