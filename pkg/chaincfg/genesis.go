package chaincfg

import (
	"time"
	
	"github.com/p9c/pod/pkg/wire"
	"github.com/p9c/pod/pkg/chainhash"
	"github.com/p9c/pod/pkg/fork"
)

// genesisCoinbaseTx is the coinbase transaction for the genesis blocks for the main network, regression test network,
// and test network (version 3).
var genesisCoinbaseTx = wire.MsgTx{
	Version: 2,
	TxIn: []*wire.TxIn{
		{
			PreviousOutPoint: wire.OutPoint{
				Hash:  chainhash.Hash{},
				Index: 0xffffffff,
			},
			SignatureScript: []byte{
				0x04, 0xff, 0xff, 0x00, 0x1d, 0x01, 0x04, 0x32,
				0x4e, 0x59, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x20,
				0x32, 0x30, 0x31, 0x34, 0x2d, 0x30, 0x37, 0x2d,
				0x31, 0x39, 0x20, 0x2d, 0x20, 0x44, 0x65, 0x6c,
				0x6c, 0x20, 0x42, 0x65, 0x67, 0x69, 0x6e, 0x73,
				0x20, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x69,
				0x6e, 0x67, 0x20, 0x42, 0x69, 0x74, 0x63, 0x6f,
				0x69, 0x6e,
			},
			Sequence: 0xffffffff,
		},
	},
	TxOut: []*wire.TxOut{
		{
			Value: 0x174876E800,
			PkScript: []byte{
				0x41, 0x04, 0xe0, 0xd2, 0x71, 0x72, 0x51, 0x0c,
				0x68, 0x06, 0x88, 0x97, 0x40, 0xed, 0xaf, 0xe6,
				0xe6, 0x3e, 0xb2, 0x3f, 0xca, 0x32, 0x78, 0x6f,
				0xcc, 0xfd, 0xb2, 0x82, 0xbb, 0x28, 0x76, 0xa9,
				0xf4, 0x3b, 0x22, 0x82, 0x45, 0xdf, 0x05, 0x76,
				0x61, 0xff, 0x94, 0x3f, 0x61, 0x50, 0x71, 0x6a,
				0x20, 0xea, 0x18, 0x51, 0xe8, 0xa7, 0xe9, 0xf5,
				0x4e, 0x62, 0x02, 0x97, 0x66, 0x46, 0x18, 0x43,
				0x8d, 0xae, 0xac,
			},
		},
	},
	LockTime: 0,
}

// genesisHash is the hash of the first block in the block chain for the main network (genesis block).
var genesisHash = chainhash.Hash(
	[chainhash.HashSize]byte{
		0xc7, 0xcc, 0x40, 0xc7, 0xc5, 0x4f, 0xd1, 0x39,
		0x1d, 0xdf, 0x3a, 0xe7, 0xcf, 0x98, 0xf2, 0x8b,
		0x23, 0xcf, 0xfd, 0x0c, 0x66, 0xd3, 0x04, 0xc9,
		0xaa, 0xd3, 0xba, 0xfc, 0xf0, 0x09, 0x00, 0x00,
		// 0x00, 0x00, 0x09, 0xf0, 0xfc, 0xba, 0xd3, 0xaa,
		// 0xc9, 0x04, 0xd3, 0x66, 0x0c, 0xfd, 0xcf, 0x23,
		// 0x8b, 0xf2, 0x98, 0xcf, 0xe7, 0x3a, 0xdf, 0x1d,
		// 0x39, 0xd1, 0x4f, 0xc5, 0xc7, 0x40, 0xcc, 0xc7,
	},
)

// genesisMerkleRoot is the hash of the first transaction in the genesis block for the main network.
var genesisMerkleRoot = chainhash.Hash(
	[chainhash.HashSize]byte{
		// 0xc8, 0x43, 0xea, 0xe4, 0x65, 0x8e, 0x3a, 0x51,
		// 0xd2, 0xf2, 0x80, 0xc3, 0x63, 0x76, 0xce, 0x56,
		// 0xdc, 0x71, 0xa6, 0xc7, 0x0e, 0x4b, 0x1c, 0x5a,
		// 0xd2, 0xd7, 0xa9, 0x31, 0x6f, 0x9b, 0x9a, 0xb7,
		0xb7, 0x9a, 0x9b, 0x6f, 0x31, 0xa9, 0xd7, 0xd2,
		0x5a, 0x1c, 0x4b, 0x0e, 0xc7, 0xa6, 0x71, 0xdc,
		0x56, 0xce, 0x76, 0x63, 0xc3, 0x80, 0xf2, 0xd2,
		0x51, 0x3a, 0x8e, 0x65, 0xe4, 0xea, 0x43, 0xc8,
	},
)

// genesisBlock defines the genesis block of the block chain which serves as the public transaction ledger for the main
// network.
var genesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    2,
		PrevBlock:  chainhash.Hash{},
		MerkleRoot: genesisMerkleRoot,
		Timestamp:  time.Unix(0x53c9ecdc, 0), // 2014-07-19 03:58:20 +0000 UTC
		Bits:       0x1e0fffff,               // 4294905630[00000fffff000000000000000000000000000000000000000000000000000000]
		Nonce:      0x10281,                  // 2164392192
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}

// regTestGenesisHash is the hash of the first block in the block chain for the regression test network (genesis block).
var regTestGenesisHash = chainhash.Hash(
	[chainhash.HashSize]byte{
		0x81, 0x91, 0x37, 0x60, 0xab, 0x59, 0x85, 0x57,
		0x7b, 0x23, 0x4d, 0xf6, 0xe2, 0x65, 0xba, 0x6b,
		0x48, 0x7e, 0x66, 0x25, 0xc8, 0x52, 0x2a, 0xdc,
		0x83, 0xa1, 0x0e, 0x22, 0x9e, 0xb7, 0xe9, 0x69,
	},
)

// regTestGenesisMerkleRoot is the hash of the first transaction in the genesis block for the regression test network.
// It is the same as the merkle root for the main network.
var regTestGenesisMerkleRoot = genesisMerkleRoot

// regTestGenesisBlock defines the genesis block of the block chain which serves as the public transaction ledger for
// the regression test network.
var regTestGenesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    2,
		PrevBlock:  chainhash.Hash{}, // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: regTestGenesisMerkleRoot,
		Timestamp:  time.Unix(0x53c9ea84, 0),
		Bits:       0x207fffff, // 4294934304 [7fffff0000000000000000000000000000000000000000000000000000000000]
		Nonce:      0x00000001,
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}

// testNet3GenesisMerkleRoot is the hash of the first transaction in the genesis block for the test network (version 3).
// It is the same as the merkle root for the main network.
var testNet3GenesisMerkleRoot = genesisMerkleRoot

// testNet3GenesisBlock defines the genesis block of the block chain which serves as the public transaction ledger for
// the test network (version 3).
var testNet3GenesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    2,
		PrevBlock:  chainhash.Hash{}, // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: testNet3GenesisMerkleRoot,
		Timestamp:  time.Unix(0x5dda3362, 0),
		Bits:       fork.SecondPowLimitBits, // 0x1e00f1ea, //testnetBits, // 0x1e0fffff,
		// 486604799 [00000000ffff0000000000000000000000000000000000000000000000000000]
		Nonce: 0x001adf18, // 417274368
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}

// testNet3GenesisHash is the hash of the first block in the block chain for the test network (version 3).
// var testNet3GenesisHash = chainhash.Hash([chainhash.HashSize]byte{
// 	0xdf, 0x0c, 0xb3, 0x5f, 0x69, 0x72, 0x75, 0xe1,
// 	0x8f, 0x66, 0xa2, 0x7d, 0xc8, 0xbb, 0x12, 0xfa,
// 	0x85, 0x4d, 0xed, 0x22, 0x2c, 0x0c, 0x1b, 0xf9,
// 	0x5e, 0xa3, 0xba, 0xec, 0x41, 0x0e, 0x00, 0x00,
// })
var testNet3GenesisHash = testNet3GenesisBlock.Header.BlockHash()

// simNetGenesisHash is the hash of the first block in the block chain for the simulation test network.
var simNetGenesisHash = chainhash.Hash(
	[chainhash.HashSize]byte{
		0xdf, 0x0c, 0xb3, 0x5f, 0x69, 0x72, 0x75, 0xe1,
		0x8f, 0x66, 0xa2, 0x7d, 0xc8, 0xbb, 0x12, 0xfa,
		0x85, 0x4d, 0xed, 0x22, 0x2c, 0x0c, 0x1b, 0xf9,
		0x5e, 0xa3, 0xba, 0xec, 0x41, 0x0e, 0x00, 0x00,
	},
)

// simNetGenesisMerkleRoot is the hash of the first transaction in the genesis block for the simulation test network. It
// is the same as the merkle root for the main network.
var simNetGenesisMerkleRoot = genesisMerkleRoot

// simNetGenesisBlock defines the genesis block of the block chain which serves as the public transaction ledger for the
// simulation test network.
var simNetGenesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    2,
		PrevBlock:  chainhash.Hash{}, // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: simNetGenesisMerkleRoot,
		Timestamp:  time.Unix(0x53c9ea84, 0),
		Bits:       0x207fffff, // 545259519 [7fffff0000000000000000000000000000000000000000000000000000000000]
		Nonce:      2,
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}
