/*
 * +----------------------------------------------------------------------
 *  | kungfu [ A FAST GAME FRAMEWORK ]
 *  +----------------------------------------------------------------------
 *  | Copyright (c) 2023-2029 All rights reserved.
 *  +----------------------------------------------------------------------
 *  | Licensed ( http:www.apache.org/licenses/LICENSE-2.0 )
 *  +----------------------------------------------------------------------
 *  | Author: jqiris <1920624985@qq.com>
 *  +----------------------------------------------------------------------
 */

package probs

import (
	"errors"
	"math/rand"
	"sort"
)

type WgItem struct {
	Element interface{}
	Weight  int
}

type WgItems []WgItem

func (w WgItems) Len() int           { return len(w) }
func (w WgItems) Swap(i, j int)      { w[i], w[j] = w[j], w[i] }
func (w WgItems) Less(i, j int) bool { return w[i].Weight > w[j].Weight }

type WgRand struct {
	RandItems         WgItems
	TotalWeight       int
	calibratedWeights bool
	precision         int
	calibrateValue    int
	choicePop         bool
}

func (wrc *WgRand) AddElement(element interface{}, weight int) {
	weight *= wrc.calibrateValue
	wrc.RandItems = append(wrc.RandItems, WgItem{Element: element, Weight: weight})
	sort.Sort(wrc.RandItems)
	wrc.TotalWeight += weight
}

func (wrc *WgRand) AddElements(elements map[interface{}]int) {
	for element, weight := range elements {
		wrc.AddElement(element, weight)
	}
}

func (wrc *WgRand) GetRandomChoice() (interface{}, error) {
	if !wrc.calibratedWeights {
		wrc.calibrateWeights()
	}
	if wrc.TotalWeight < 1 {
		return nil, errors.New("权重不正确")
	}
	value := rand.Intn(wrc.TotalWeight)
	for key, item := range wrc.RandItems {
		value -= item.Weight
		if value <= 0 {
			if wrc.choicePop {
				wrc.RandItems = append(wrc.RandItems[:key], wrc.RandItems[key+1:]...)
				wrc.TotalWeight -= item.Weight
			}
			return item.Element, nil
		}
	}
	return nil, errors.New("not found")
}

func (wrc *WgRand) calibrateWeights() {
	if wrc.TotalWeight/wrc.precision < 1 {
		wrc.calibrateValue = wrc.precision / wrc.TotalWeight
		wrc.TotalWeight = 0
		for key, item := range wrc.RandItems {
			newWeight := item.Weight * wrc.calibrateValue
			wrc.RandItems[key].Weight = newWeight
			wrc.TotalWeight += newWeight
		}
		sort.Sort(wrc.RandItems)
		wrc.calibratedWeights = true
	}
}

func NewWgRand(choicePop bool, arguments ...int) WgRand {
	var precision = 1000
	if len(arguments) > 0 {
		precision = arguments[0]
	}
	return WgRand{
		precision:         precision,
		calibratedWeights: false,
		calibrateValue:    1,
		choicePop:         choicePop,
	}
}
