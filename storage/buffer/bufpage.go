package buffer

import "unsafe"

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
	Tuples       []byte
	SpecialSpace []byte
}

type PageXLogRecPtr uint64

type LocationIndex uint16

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
	Lower LocationIndex
	// offset to end of free space
	Upper LocationIndex
	// offset to start of special space
	Special LocationIndex
	// page size and page version
	PageSize uint16
}

func (p *Page) PageInit(pageSize uint16, specialSize uint16) {
	// 先不考虑内存对齐的场景
	p.Header.Lower = LocationIndex(unsafe.Sizeof(PageHeaderData{}))
	p.Header.Upper = LocationIndex(pageSize - specialSize)
	p.Header.Special = LocationIndex(pageSize - specialSize)
	p.Header.PageSize = pageSize
}

func (p *Page) PageAddItemExtended(item []byte, size int32, offsetNumber int32, flags int32) int32 {
	return 0
}

func (p *Page) PageGetFreeSpace() int32 {
	// TODO
	return 0
}
