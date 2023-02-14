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

package discover

import "github.com/jqiris/kungfu/v2/treaty"

type MaintainType int

const (
	MaintainTypeAll MaintainType = iota
	MaintainTypeFalse
	MaintainTypeTrue
)

type Filter struct {
	maintainType MaintainType
	version      int64
	ignore       bool //忽略具体状态检查
}

func NewFilter(options ...FilterOption) *Filter {
	filter := &Filter{
		maintainType: MaintainTypeFalse,
		version:      0,
		ignore:       false,
	}
	for _, option := range options {
		option(filter)
	}
	return filter
}

func (f *Filter) apply(s *treaty.Server) bool {
	if f.maintainType > MaintainTypeAll {
		if f.ignore {
			return true
		}
		if f.maintainType == MaintainTypeFalse && s.Maintained {
			return false
		}
		if f.maintainType == MaintainTypeTrue && !s.Maintained {
			return false
		}
	}
	if f.version > 0 && f.version != s.Version {
		return false
	}
	return true
}

type FilterOption func(f *Filter)

func FilterMaintained(maintained MaintainType) FilterOption {
	return func(f *Filter) {
		f.maintainType = maintained
	}
}

func FilterVersion(version int64) FilterOption {
	return func(f *Filter) {
		f.version = version
	}
}

func FilterIgnore(ignore bool) FilterOption {
	return func(f *Filter) {
		f.ignore = ignore
	}
}
