package coroutines

import (
	"fmt"
	"testing"

	"github.com/jqiris/kungfu/v2/utils"
)

func TestNumberMap(t *testing.T) {
	data := NewNumberMap[int32, int32]()
	data.Incre(1, 3)
	fmt.Println(data.Load(1))
	data.Decre(1, 1)
	fmt.Println(data.Load(1))
	v, ok := data.LoadOk(1)
	fmt.Println(v, ok)
	v, ok = data.LoadOk(2)
	fmt.Println(v, ok)
	fmt.Println(data)
	a := utils.Stringify(data)
	fmt.Printf("%+v\n", a)
	ndata := NewNumberMap[int32, int32]()
	if err := ndata.UnmarshalJSON([]byte(a)); err != nil {
		fmt.Println(err)
	}
	ndata.Store(2, 5)
	fmt.Println(ndata)
	ndata.Range(func(k, v int32) bool {
		fmt.Println(k, v)
		return true
	})
}
