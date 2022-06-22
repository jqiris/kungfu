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
)

type DataItem struct {
	UserId int64 `json:"user_id"`
	MsgId  int32 `json:"msg_id"`
	MsgReq any   `json:"msg_req"`
}

type DataHandler func(item *DataItem)

type ServerData struct {
	dataHandler DataHandler
	dealChanArr []chan *DataItem
	doneChan    chan int
	maxDbQueue  int64
	dbQueueSize int64
	cancelCtx   context.Context
	cancelFunc  context.CancelFunc
	serverId    string
}

func NewServerData(dataHandler DataHandler, maxDbQueue, dbQueueSize int64) *ServerData {
	if maxDbQueue < 1 {
		maxDbQueue = 10
	}
	if dbQueueSize < 1 {
		dbQueueSize = 1000
	}
	s := &ServerData{
		dataHandler: dataHandler,
		dealChanArr: make([]chan *DataItem, maxDbQueue),
		doneChan:    make(chan int, maxDbQueue),
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
			logger.Infof("%v DealReq receive cancel signal", b.serverId)
			return
		default:
			var item *DataItem
			if err := stores.BRPop(s.Server.ServerId, &item); err == nil {
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

func (b *ServerData) DealQueue(ctx context.Context, num int, queue chan *DataItem) {
	for {
		select {
		case item := <-queue:
			b.dataHandler(item)
		case <-ctx.Done():
			logger.Infof("%v DealQueue receive cancel signal:%v", b.serverId, num)
			b.doneChan <- num
			return
		}
	}
}

func (b *ServerData) Init(s *rpc.ServerBase) {
}

func (b *ServerData) AfterInit(s *rpc.ServerBase) {
	b.serverId = s.Server.ServerId
	//初始化队列
	for i := 0; i < len(b.dealChanArr); i++ {
		b.dealChanArr[i] = make(chan *DataItem, b.dbQueueSize)
	}
	//读取队列
	go utils.SafeRun(func() {
		b.DealReq(b.cancelCtx, s)
	})
	//处理数据
	for i := 0; i < len(b.dealChanArr); i++ {
		k, v := i, b.dealChanArr[i]
		go utils.SafeRun(func() {
			b.DealQueue(b.cancelCtx, k, v)
		})
	}

}

func (b *ServerData) BeforeShutdown(s *rpc.ServerBase) {
	b.cancelFunc()
	//剩余队列资源处理
	wg := sync.WaitGroup{}
	for i := 0; i < len(b.dealChanArr); i++ {
		k := <-b.doneChan
		v := b.dealChanArr[k]
		close(v)
		wg.Add(1)
		go func(num int, queue chan *DataItem) {
			defer wg.Done()
			logger.Infof("%v ServerData Shutting down deal begin:%v", b.serverId, num)
			for item := range queue {
				b.dataHandler(item)
			}
			logger.Infof("%v ServerData Shutting down deal end:%v", b.serverId, num)
		}(k, v)
	}
	wg.Wait()
	fmt.Println(b.serverId + " ServerData Shutting down deal chan over")
}

func (b *ServerData) Shutdown(s *rpc.ServerBase) {
}
