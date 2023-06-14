package concurrency

type TxIdGenerator struct {
	id int64
}

func constructTxIdGenerator() *TxIdGenerator {
	// TODO: generate txId
	return nil
}

func (generator *TxIdGenerator) getTxId() int64 {
	// TODO
	return 0
}

type BufferPool struct {
}

func (bufferPool *BufferPool) read() {

}

type TimestampProtocol struct {
}

func (ts *TimestampProtocol) read(txId int64) {

}

func (ts *TimestampProtocol) write(txId int64) {

}
