package wire

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/parallelcointeam/parallelcoin/pkg/chain/fork"
	chainhash "github.com/parallelcointeam/parallelcoin/pkg/chain/hash"
	"github.com/parallelcointeam/parallelcoin/pkg/util/cl"
)

// MaxBlockHeaderPayload is the maximum number of bytes a block header can be. Version 4 bytes + Timestamp 4 bytes + Bits 4 bytes + Nonce 4 bytes + PrevBlock and MerkleRoot hashes.
const MaxBlockHeaderPayload = 16 + (chainhash.HashSize * 2)

// BlockHeader defines information about a block and is used in the bitcoin block (MsgBlock) and headers (MsgHeaders) messages.
type BlockHeader struct {
	// Version of the block.  This is not the same as the protocol version.
	Version int32
	// Hash of the previous block header in the block chain.
	PrevBlock chainhash.Hash
	// MerkleRoot is the Merkle tree reference to hash of all transactions for the block.
	MerkleRoot chainhash.Hash
	// Time the block was created.  This is, unfortunately, encoded as a uint32 on the wire and therefore is limited to 2106.
	Timestamp time.Time
	// Difficulty target for the block.
	Bits uint32
	// Nonce used to generate the block.
	Nonce uint32
}

// blockHeaderLen is a constant that represents the number of bytes for a block header.
const blockHeaderLen = 80

// BlockHash computes the block identifier hash for the given block header.
func (h *BlockHeader) BlockHash() (out chainhash.Hash) {
	// Encode the header and double sha256 everything prior to the number of transactions.  Ignore the error returns since there is no way the encode could fail except being out of memory which would cause a run-time panic.
	buf := bytes.NewBuffer(make([]byte, 0, MaxBlockHeaderPayload))
	_ = writeBlockHeader(buf, 0, h)
	out = chainhash.DoubleHashH(buf.Bytes())
	return
}

// BlockHashWithAlgos computes the block identifier hash for the given block header. This function is additional because the sync manager and the parallelcoin protocol only use SHA256D hashes for inventories and calculating the scrypt (or other) hash for these blocks when requested via that route causes an 'unrequested block' error.
func (h *BlockHeader) BlockHashWithAlgos(height int32) (out chainhash.Hash) {
	// Encode the header and double sha256 everything prior to the number of transactions.  Ignore the error returns since there is no way the encode could fail except being out of memory which would cause a run-time panic.
	buf := bytes.NewBuffer(make([]byte, 0, MaxBlockHeaderPayload))
	err := writeBlockHeader(buf, 0, h)
	if err != nil {
		fmt.Println("error writing block header to buffer", err, cl.Ine())
	}
	vers := h.Version
	algo := fork.GetAlgoName(vers, height)
	out = fork.Hash(buf.Bytes(), algo, height)
	// fmt.Printf("BlockHashWithAlgos %d %s %s %s\n", vers, algo, out, cl.Ine())
	return
}

// BtcDecode decodes r using the bitcoin protocol encoding into the receiver. This is part of the Message interface implementation. See Deserialize for decoding block headers stored to disk, such as in a database, as opposed to decoding block headers from the wire.
func (h *BlockHeader) BtcDecode(r io.Reader, pver uint32, enc MessageEncoding) error {
	return readBlockHeader(r, pver, h)
}

// BtcEncode encodes the receiver to w using the bitcoin protocol encoding. This is part of the Message interface implementation. See Serialize for encoding block headers to be stored to disk, such as in a database, as opposed to encoding block headers for the wire.
func (h *BlockHeader) BtcEncode(w io.Writer, pver uint32, enc MessageEncoding) error {
	return writeBlockHeader(w, pver, h)
}

// Deserialize decodes a block header from r into the receiver using a format that is suitable for long-term storage such as a database while respecting the Version field.
func (h *BlockHeader) Deserialize(r io.Reader) error {
	// At the current time, there is no difference between the wire encoding at protocol version 0 and the stable long-term storage format.  As a result, make use of readBlockHeader.
	return readBlockHeader(r, 0, h)
}

// Serialize encodes a block header from r into the receiver using a format that is suitable for long-term storage such as a database while respecting the Version field.
func (h *BlockHeader) Serialize(w io.Writer) error {
	// At the current time, there is no difference between the wire encoding at protocol version 0 and the stable long-term storage format.  As a result, make use of writeBlockHeader.
	return writeBlockHeader(w, 0, h)
}

// NewBlockHeader returns a new BlockHeader using the provided version, previous block hash, merkle root hash, difficulty bits, and nonce used to generate the block with defaults for the remaining fields.
func NewBlockHeader(version int32, prevHash, merkleRootHash *chainhash.Hash,
	bits uint32, nonce uint32) *BlockHeader {
	// Limit the timestamp to one second precision since the protocol doesn't support better.
	return &BlockHeader{
		Version:    version,
		PrevBlock:  *prevHash,
		MerkleRoot: *merkleRootHash,
		Timestamp:  time.Unix(time.Now().Unix(), 0),
		Bits:       bits,
		Nonce:      nonce,
	}
}

// readBlockHeader reads a bitcoin block header from r.  See Deserialize for decoding block headers stored to disk, such as in a database, as opposed to decoding from the wire.
func readBlockHeader(r io.Reader, pver uint32, bh *BlockHeader) error {
	return readElements(r, &bh.Version, &bh.PrevBlock, &bh.MerkleRoot,
		(*uint32Time)(&bh.Timestamp), &bh.Bits, &bh.Nonce)
}

// writeBlockHeader writes a bitcoin block header to w.  See Serialize for encoding block headers to be stored to disk, such as in a database, as opposed to encoding for the wire.
func writeBlockHeader(w io.Writer, pver uint32, bh *BlockHeader) error {
	sec := uint32(bh.Timestamp.Unix())
	return writeElements(w, bh.Version, &bh.PrevBlock, &bh.MerkleRoot,
		sec, bh.Bits, bh.Nonce)
}
