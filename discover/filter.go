package discover

import "github.com/jqiris/kungfu/v2/treaty"

type Filter func(s *treaty.Server) bool

type Filters []Filter

func (fs Filters) Apply(s *treaty.Server) bool {
	for _, f := range fs {
		if !f(s) {
			return false
		}
	}
	return true
}

func WithFilterMaintained(maintained bool) Filter {
	return func(s *treaty.Server) bool {
		return s.Maintained == maintained
	}
}

func WithFilterVersion(version string) Filter {
	return func(s *treaty.Server) bool {
		return s.Version == version
	}
}
