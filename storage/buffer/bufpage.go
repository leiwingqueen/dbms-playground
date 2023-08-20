package buffer

import "encoding/binary"

/**
* +----------------+---------------------------------+
* | PageHeaderData | linp1 linp2 linp3 ...           |
* +-----------+----+---------------------------------+
* | ... linpN |									  |
* +-----------+--------------------------------------+
* |		   ^ pd_lower							  |
* |												  |
* |			 v pd_upper							  |
* +-------------+------------------------------------+
* |			 | tupleN ...                         |
* +-------------+------------------+-----------------+
* |	   ... tuple3 tuple2 tuple1 | "special space" |
* +--------------------------------+-----------------+
*									^ pd_special
 */
type Page struct {
	Header PageHeaderData
	// line pointer array
	LinePointers []ItemIdData
	Tuples       [][]byte
	SpecialSpace []byte
}

type PageXLogRecPtr uint64

// type LocationIndex uint16

type TransactionId uint32

type ItemIdData struct {
	/* offset to tuple (from start of page) */
	lp_off uint16
	/* state of line pointer, see below */
	lp_flags uint8
	/* byte length of tuple */
	lp_len uint16
}

type PageHeaderData struct {
	// LSN: next byte after last byte of xlog record for last change to this page
	Lsn PageXLogRecPtr
	// checksum
	CheckSum uint16
	// flag bits, see below
	Flags uint16
	// offset to start of free space
	Lower uint16
	// offset to end of free space
	Upper uint16
	// offset to start of special space
	Special uint16
	// page size and page version
	PageSize uint16
}

// PageHeaderData convert to byte array
func (p *PageHeaderData) ToBytes() []byte {
	size := pageHeaderSize()
	bytes := make([]byte, size)
	// LSN
	binary.LittleEndian.PutUint64(bytes[0:8], uint64(p.Lsn))
	// CheckSum
	binary.LittleEndian.PutUint16(bytes[8:10], p.CheckSum)
	// Flags
	binary.LittleEndian.PutUint16(bytes[10:12], p.Flags)
	// Lower
	binary.LittleEndian.PutUint16(bytes[12:14], p.Lower)
	// Upper
	binary.LittleEndian.PutUint16(bytes[14:16], p.Upper)
	// Special
	binary.LittleEndian.PutUint16(bytes[16:18], p.Special)
	// PageSize
	binary.LittleEndian.PutUint16(bytes[18:20], p.PageSize)
	return bytes
}

func pageHeaderSize() uint16 {
	return 8 + 2 + 2 + 2 + 2 + 2 + 2
}

// ItemIdData size
func itemIdDataSize() uint16 {
	return 2 + 1 + 2
}

func (p *Page) PageInit(pageSize uint16, specialSize uint16) {
	// 先不考虑内存对齐的场景
	p.Header.Lower = pageHeaderSize()
	p.Header.Upper = pageSize - specialSize
	p.Header.Special = pageSize - specialSize
	p.Header.PageSize = pageSize
}

func (p *Page) PageAddItemExtended(item []byte, size int32, offsetNumber int32, flags int32) uint16 {
	header := p.Header
	header.Upper -= uint16(size)
	header.Lower += itemIdDataSize()
	p.LinePointers = append(p.LinePointers, ItemIdData{header.Upper, 0, uint16(size)})
	p.Tuples = append(p.Tuples, item)
	return header.Upper
}

func (p *Page) PageGetFreeSpace() int32 {
	// TODO
	return 0
}

func (p *Page) ToByteArray(offsetNumber int32) []byte {
	header := p.Header
	size := header.PageSize
	bytes := make([]byte, size)
	// PageHeaderData
	copy(bytes[0:pageHeaderSize()], header.ToBytes())
	// ItemIdData
	offset := pageHeaderSize()
	for i := 0; i < len(p.LinePointers); i++ {
		copy(bytes[offset:offset+itemIdDataSize()], p.LinePointers[i].ToBytes())
		offset += itemIdDataSize()
	}
	// Tuples
	// TODO
	return nil
}
