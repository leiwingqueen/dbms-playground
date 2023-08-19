package buffer

type Page []byte

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
	pd_lsn PageXLogRecPtr
	// checksum
	pd_checksum uint16
	// flag bits, see below
	pd_flags uint16
	// offset to start of free space
	pd_lower LocationIndex
	// offset to end of free space
	pd_upper LocationIndex
	// offset to start of special space
	pd_special LocationIndex
	// page size and page version
	pd_pagesize_version uint16
	// oldest prunable XID, or zero if none
	pd_prune_xid TransactionId
	// line pointer array
	pd_linp []ItemIdData
}
