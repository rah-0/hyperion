package util

const (
	Size2KB   = 2 * 1024
	Size4KB   = 4 * 1024
	Size8KB   = 8 * 1024
	Size16KB  = 16 * 1024
	Size32KB  = 32 * 1024
	Size64KB  = 64 * 1024
	Size128KB = 128 * 1024
	Size256KB = 256 * 1024
	Size512KB = 512 * 1024
	Size1MB   = 1 * 1024 * 1024
	Size10MB  = 10 * 1024 * 1024
	Size100MB = 100 * 1024 * 1024
	Size1GB   = 1 * 1024 * 1024 * 1024
)

// StructField represents a field in a struct
type StructField struct {
	Name          string
	Type          string
	Tag           string
	IndexSortable bool
}

// StructDef represents the whole struct definition
type StructDef struct {
	Name   string
	Fields []StructField
}
