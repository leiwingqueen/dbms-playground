package buffer

const NBuffers = 1000
const PageSize = 8 * 1024

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
	tag      BufferTag
	bufId    int
	freeNext int
	// To simplify the following descriptions, three descriptor states are defined:
	//Empty: When the corresponding buffer pool slot does not store a page (i.e. refcount and usage_count are 0), the state of this descriptor is empty.
	//Pinned: When the corresponding buffer pool slot stores a page and any PostgreSQL processes are accessing the page (i.e. refcount and usage_count are greater than or equal to 1), the state of this buffer descriptor is pinned.
	//Unpinned: When the corresponding buffer pool slot stores a page but no PostgreSQL processes are accessing the page (i.e. usage_count is greater than or equal to 1, but refcount is 0), the state of this buffer descriptor is unpinned.
	pinCount   int
	usageCount int
	dirty      bool
}

type BufferPool struct {
	// mapping from buffer tag to buffer id
	bufferMap         map[BufferTag]int
	bufferDescriptors []BufferDesc
	bufferBlocks      []Page
	freeListHead      int
	freeListTail      int
}

// init buffer pool
func constructBufferPool() *BufferPool {
	pool := BufferPool{
		bufferMap:         make(map[BufferTag]int),
		bufferDescriptors: make([]BufferDesc, NBuffers),
		bufferBlocks:      make([]Page, NBuffers),
		freeListHead:      0,
		freeListTail:      NBuffers - 1,
	}
	for i := 0; i < NBuffers; i++ {
		pool.bufferDescriptors[i] = BufferDesc{
			tag:        BufferTag{},
			bufId:      i,
			freeNext:   i + 1,
			pinCount:   0,
			usageCount: 0,
			dirty:      false,
		}
	}
	pool.bufferDescriptors[NBuffers-1].freeNext = -1
	return &pool
}

func (bufferPool *BufferPool) read(tag BufferTag) *Page {
	// already in buffer pool, return it
	if bufferId, exist := bufferPool.bufferMap[tag]; exist {
		bufferPool.bufferDescriptors[bufferId].pinCount++
		bufferPool.bufferDescriptors[bufferId].usageCount++
		return &bufferPool.bufferBlocks[bufferId]
	}
	// not in buffer pool, get a free buffer
	if bufferPool.freeListHead >= 0 {
		bufferId := bufferPool.freeListHead
		bufferDesc := bufferPool.bufferDescriptors[bufferId]
		bufferPool.freeListHead = bufferDesc.freeNext
		bufferDesc.freeNext = -1
		bufferDesc.tag.blockNum = tag.blockNum
		bufferDesc.tag.dbNode = tag.dbNode
		bufferDesc.pinCount = 1
		bufferDesc.usageCount = 1
		bufferPool.bufferMap[tag] = bufferId
		// load data from disk
		page := loadPageFromDisk(tag)
		bufferPool.bufferBlocks[bufferId] = page
		return &bufferPool.bufferBlocks[bufferId]
	}
	// TODO: page replacement strategy
	evictBufferId := 0
	evictDesc := bufferPool.bufferDescriptors[evictBufferId]
	// update the buffer map
	bufferPool.bufferMap[tag] = evictBufferId
	delete(bufferPool.bufferMap, evictDesc.tag)

	evictDesc.tag.blockNum = tag.blockNum
	evictDesc.tag.dbNode = tag.dbNode
	evictDesc.pinCount = 1
	evictDesc.usageCount = 1
	// load data from disk
	page := loadPageFromDisk(tag)
	bufferPool.bufferBlocks[evictBufferId] = page
	return &bufferPool.bufferBlocks[evictBufferId]
}

func loadPageFromDisk(tag BufferTag) Page {
	// mock function
	data := make([]byte, PageSize)
	return data
}
