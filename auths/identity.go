package auths

import (
	"sync"

	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/stores"
)

type IdSourcer func(key string) (int64, error)
type IdStorer func(key string, id int64) error

type IdGenerator struct {
	table    string
	key      string
	idSource IdSourcer
	idStore  IdStorer
	lock     sync.Mutex
	nextId   int64
}

func NewIdGenerator(table string, key string, idSource IdSourcer, idStore IdStorer) *IdGenerator {
	return &IdGenerator{
		table:    table,
		key:      key,
		idSource: idSource,
		idStore:  idStore,
		lock:     sync.Mutex{},
		nextId:   -1,
	}
}

func (g *IdGenerator) check() error {
	if !stores.HExists(g.table, g.key) {
		g.lock.Lock()
		defer g.lock.Unlock()
		id, err := g.idSource(g.key)
		if err != nil {
			logger.Reportf("id generator source err: %v", err)
			return err
		}
		if g.nextId != -1 && g.nextId > id {
			id = g.nextId
		}
		if err = stores.HSetNx(g.table, g.key, id); err != nil {
			logger.Reportf("id generator init err: %v", err)
			return err
		}
	}
	return nil
}

func (g *IdGenerator) NextId() (int64, error) {
	if err := g.check(); err != nil {
		return 0, err
	}
	nextId := stores.HIncrBy(g.table, g.key, 1)
	g.nextId = nextId
	return nextId, nil
}

func (g *IdGenerator) Store() error {
	g.lock.Lock()
	defer g.lock.Unlock()
	var id int64
	if err := stores.HGet(g.table, g.key, &id); err != nil {
		return err
	}
	return g.idStore(g.key, id)
}
