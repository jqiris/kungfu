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
}

func NewFilter(options ...FilterOption) *Filter {
	filter := &Filter{
		maintainType: MaintainTypeFalse,
		version:      0,
	}
	for _, option := range options {
		option(filter)
	}
	return filter
}

func (f *Filter) apply(s *treaty.Server) bool {
	if f.maintainType > MaintainTypeAll {
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

func WithFilterMaintained(maintained MaintainType) FilterOption {
	return func(f *Filter) {
		f.maintainType = maintained
	}
}

func WithFilterVersion(version int64) FilterOption {
	return func(f *Filter) {
		f.version = version
	}
}
