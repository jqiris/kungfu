package stores

import "sync"

type StoreMem struct {
	sync.RWMutex // protect data
	data         map[string]any
}

func NewStoreMem() *StoreMem {
	return &StoreMem{
		data: make(map[string]any),
	}
}

// Remove delete data associated with the key from session storage
func (s *StoreMem) Remove(key string) {
	s.Lock()
	defer s.Unlock()
	delete(s.data, key)
}

// Set associates value with the key in session storage
func (s *StoreMem) Set(key string, value any) {
	s.Lock()
	defer s.Unlock()
	s.data[key] = value
}

// HasKey decides whether a key has associated value
func (s *StoreMem) HasKey(key string) bool {
	s.RLock()
	defer s.RUnlock()

	_, has := s.data[key]
	return has
}

// Int returns the value associated with the key as a int.
func (s *StoreMem) Int(key string) int {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(int)
	if !ok {
		return 0
	}
	return value
}

// Int returns the value associated with the key as a bool.
func (s *StoreMem) Bool(key string) bool {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return false
	}

	value, ok := v.(bool)
	if !ok {
		return false
	}
	return value
}

// Int8 returns the value associated with the key as a int8.
func (s *StoreMem) Int8(key string) int8 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(int8)
	if !ok {
		return 0
	}
	return value
}

// Int16 returns the value associated with the key as a int16.
func (s *StoreMem) Int16(key string) int16 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(int16)
	if !ok {
		return 0
	}
	return value
}

// Int32 returns the value associated with the key as a int32.
func (s *StoreMem) Int32(key string) int32 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(int32)
	if !ok {
		return 0
	}
	return value
}

// Int64 returns the value associated with the key as a int64.
func (s *StoreMem) Int64(key string) int64 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(int64)
	if !ok {
		return 0
	}
	return value
}

// Uint returns the value associated with the key as a uint.
func (s *StoreMem) Uint(key string) uint {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(uint)
	if !ok {
		return 0
	}
	return value
}

// Uint8 returns the value associated with the key as a uint8.
func (s *StoreMem) Uint8(key string) uint8 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(uint8)
	if !ok {
		return 0
	}
	return value
}

// Uint16 returns the value associated with the key as a uint16.
func (s *StoreMem) Uint16(key string) uint16 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(uint16)
	if !ok {
		return 0
	}
	return value
}

// Uint32 returns the value associated with the key as a uint32.
func (s *StoreMem) Uint32(key string) uint32 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(uint32)
	if !ok {
		return 0
	}
	return value
}

// Uint64 returns the value associated with the key as a uint64.
func (s *StoreMem) Uint64(key string) uint64 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(uint64)
	if !ok {
		return 0
	}
	return value
}

// Float32 returns the value associated with the key as a float32.
func (s *StoreMem) Float32(key string) float32 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(float32)
	if !ok {
		return 0
	}
	return value
}

// Float64 returns the value associated with the key as a float64.
func (s *StoreMem) Float64(key string) float64 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(float64)
	if !ok {
		return 0
	}
	return value
}

// String returns the value associated with the key as a string.
func (s *StoreMem) String(key string) string {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return ""
	}

	value, ok := v.(string)
	if !ok {
		return ""
	}
	return value
}

// Value returns the value associated with the key as a any.
func (s *StoreMem) Value(key string) any {
	s.RLock()
	defer s.RUnlock()

	return s.data[key]
}

// State returns all session state
func (s *StoreMem) Data() map[string]any {
	s.RLock()
	defer s.RUnlock()
	return s.data
}

// Restore session state after reconnect
func (s *StoreMem) Restore(data map[string]any) {
	s.Lock()
	defer s.Unlock()
	s.data = data
}

// Clear releases all data related to current session
func (s *StoreMem) Clear() {
	s.Lock()
	defer s.Unlock()
	s.data = map[string]any{}
}
