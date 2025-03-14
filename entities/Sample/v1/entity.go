package Sample

// The code in this file is autogenerated, do not modify manually!

import (
	"bytes"
	"encoding/gob"
	"sync"

	. "github.com/rah-0/hyperion/register"
)

const (
	Version    = "v1"
	Name       = "Sample"
	DbFileName = "SampleV1.bin"
)

var Fields = map[string]int{
	"Name":    1,
	"Surname": 2,
}

var (
	_       Model = (*Sample)(nil)
	mu      sync.Mutex
	Buffer  = new(bytes.Buffer)
	Encoder = gob.NewEncoder(Buffer)
	Decoder = gob.NewDecoder(Buffer)
	Mem     []*Sample
)

func init() {
	gob.Register(Sample{})

	// The following process initializes the encoder and decoder by preloading metadata.
	// This prevents metadata from being stored with the first encoded struct.
	// If the metadata were missing or inconsistent, decoding the struct later could fail.
	x := New()
	if err := x.Encode(); err != nil {
		panic("failed to encode type metadata: " + err.Error())
	}
	if err := x.Decode(); err != nil {
		panic("failed to decode type metadata: " + err.Error())
	}
	x.BufferReset()

	Mem = []*Sample{}

	RegisterEntity(&Entity{
		Version:    Version,
		Name:       Name,
		DbFileName: DbFileName,
		Fields:     Fields,
		New:        New,
	})
}

type Sample struct {
	Name    string
	Surname string
	offset  uint64
}

func New() Model {
	return &Sample{}
}

func (s *Sample) SetFieldValue(fieldName string, value any) {
	switch Fields[fieldName] {
	case 1:
		if v, ok := value.(string); ok {
			s.Name = v
		}
	case 2:
		if v, ok := value.(string); ok {
			s.Surname = v
		}
	}
}

func (s *Sample) GetFieldValue(fieldName string) any {
	switch Fields[fieldName] {
	case 1:
		return s.Name
	case 2:
		return s.Surname
	}
	return nil
}

func (s *Sample) SetOffset(offset uint64) {
	s.offset = offset
}

func (s *Sample) GetOffset() uint64 {
	return s.offset
}

func (s *Sample) Encode() error {
	mu.Lock()
	defer mu.Unlock()
	return Encoder.Encode(s)
}

func (s *Sample) Decode() error {
	mu.Lock()
	defer mu.Unlock()
	return Decoder.Decode(s)
}

func (s *Sample) BufferReset() {
	mu.Lock()
	defer mu.Unlock()
	Buffer.Reset()
}

func (s *Sample) GetBuffer() *bytes.Buffer {
	mu.Lock()
	defer mu.Unlock()
	return Buffer
}

func (s *Sample) GetBufferData() []byte {
	mu.Lock()
	defer mu.Unlock()
	return Buffer.Bytes()
}

func (s *Sample) SetBufferData(data []byte) {
	mu.Lock()
	defer mu.Unlock()
	Buffer.Write(data)
}

func (s *Sample) MemoryAdd() {
	mu.Lock()
	defer mu.Unlock()
	Mem = append(Mem, s)
}

func (s *Sample) MemoryRemove() bool {
	mu.Lock()
	defer mu.Unlock()
	for i, instance := range Mem {
		if instance == s {
			lastIndex := len(Mem) - 1
			Mem[i] = Mem[lastIndex]
			Mem = Mem[:lastIndex]
			return true
		}
	}
	return false
}

func (s *Sample) MemoryClear() {
	mu.Lock()
	defer mu.Unlock()
	Mem = []*Sample{}
}

func (s *Sample) MemoryGetAll() []Model {
	mu.Lock()
	defer mu.Unlock()
	instances := make([]Model, len(Mem))
	for i, instance := range Mem {
		instances[i] = instance
	}
	return instances
}

func (s *Sample) MemoryContains(target Model) bool {
	mu.Lock()
	defer mu.Unlock()

	for _, instance := range Mem {
		if instance == target {
			return true
		}
	}
	return false
}
