package buffer

// Oid object id
type Oid = uint32
type BlockNumber = uint32

// BufferTag
// also name the page id. it comprises the
type BufferTag struct {
	dbNode   Oid
	blockNum BlockNumber
}

// BufferDesc
// buffer description
type BufferDesc struct {
	tag        BufferTag
	bufId      int
	freeNext   int
	pinCount   int
	usageCount int
	dirty      bool
}

// Page
// pg page size is 8KB
type Page struct {
	data []byte
}

type BufferPool struct {
	// mapping from buffer tag to buffer id
	bufferMap         map[BufferTag]int
	bufferDescriptors []BufferDesc
	bufferBlocks      []Page
	freeListHead      int
	freeListTail      int
}
