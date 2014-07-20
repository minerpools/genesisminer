// genesisminer
package main

import (
	"fmt"
	"time"

	"github.com/conformal/btcchain"
	"github.com/conformal/btcwire"
)

const (
	// maxNonce is the maximum value a nonce can be in a block header.
	maxNonce = ^uint32(0) // 2^32 - 1
)

// solveBlock attempts to find some combination of a nonce, extra nonce, and
// current timestamp which makes the passed block hash to a value less than the
// target difficulty.  The timestamp is updated periodically and the passed
// block is modified with all tweaks during this process.  This means that
// when the function returns true, the block is ready for submission.
//
// This function will return early with false when conditions that trigger a
// stale block such as a new block showing up or periodically when there are
// new transactions and enough time has elapsed without finding a solution.
func SolveBlock(msgBlock *btcwire.MsgBlock, blockHeight int64) bool {
	// Create a couple of convenience variables.
	header := &msgBlock.Header
	targetDifficulty := btcchain.CompactToBig(header.Bits)

	// Initial state.
	// lastGenerated := time.Now()
	hashesCompleted := uint64(0)

	// Search through the entire nonce range for a solution while
	// periodically checking for early quit and stale block
	// conditions along with updates to the speed monitor.
	for i := uint32(0); i <= maxNonce; i++ {

		// Update the nonce and hash the block header.  Each
		// hash is actually a double sha256 (two hashes), so
		// increment the number of hashes completed for each
		// attempt accordingly.
		header.Nonce = i
		hash, _ := header.BlockSha()
		hashesCompleted += 2

		// The block is solved when the new block hash is less
		// than the target difficulty.  Yay!
		if btcchain.ShaHashToBig(&hash).Cmp(targetDifficulty) <= 0 {
			fmt.Printf("valid nonce found: %v \n", i)
			fmt.Printf("block header hash: %v \n", hash)
			return true
		}

		if i%40000 == 0 {
			fmt.Printf("%v hashes performed\n", i)
		}
	}

	return false
}

func main() {
	genesisBlock := btcwire.MsgBlock{
		Header: btcwire.BlockHeader{
			Version: 1,
			PrevBlock: btcwire.ShaHash([btcwire.HashSize]byte{ // Make go vet happy.
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			}),
			MerkleRoot: btcwire.ShaHash([btcwire.HashSize]byte{ // Make go vet happy.
				0x3b, 0xa3, 0xed, 0xfd, 0x7a, 0x7b, 0x12, 0xb2,
				0x7a, 0xc7, 0x2c, 0x3e, 0x67, 0x76, 0x8f, 0x61,
				0x7f, 0xc8, 0x1b, 0xc3, 0x88, 0x8a, 0x51, 0x32,
				0x3a, 0x9f, 0xb8, 0xaa, 0x4b, 0x1e, 0x5e, 0x4a,
			}),
			Timestamp: time.Unix(0x4966bc61, 0), // 2009-01-08 20:54:25 -0600 CST
			Bits:      0x207fffff,               // 545259519 (regtest)
			Nonce:     0x00000000,
		},
		Transactions: []*btcwire.MsgTx{
			{
				Version: 1,
				TxIn: []*btcwire.TxIn{
					{
						PreviousOutpoint: btcwire.OutPoint{
							Hash:  btcwire.ShaHash{},
							Index: 0xffffffff,
						},
						SignatureScript: []byte{
							0x04, 0xff, 0xff, 0x00, 0x1d, 0x01, 0x04, 0x45, /* |.......E| */
							0x54, 0x68, 0x65, 0x20, 0x54, 0x69, 0x6d, 0x65, /* |The Time| */
							0x73, 0x20, 0x30, 0x33, 0x2f, 0x4a, 0x61, 0x6e, /* |s 03/Jan| */
							0x2f, 0x32, 0x30, 0x30, 0x39, 0x20, 0x43, 0x68, /* |/2009 Ch| */
							0x61, 0x6e, 0x63, 0x65, 0x6c, 0x6c, 0x6f, 0x72, /* |ancellor| */
							0x20, 0x6f, 0x6e, 0x20, 0x62, 0x72, 0x69, 0x6e, /* | on brin| */
							0x6b, 0x20, 0x6f, 0x66, 0x20, 0x73, 0x65, 0x63, /* |k of sec|*/
							0x6f, 0x6e, 0x64, 0x20, 0x62, 0x61, 0x69, 0x6c, /* |ond bail| */
							0x6f, 0x75, 0x74, 0x20, 0x66, 0x6f, 0x72, 0x20, /* |out for |*/
							0x62, 0x61, 0x6e, 0x6b, 0x73, /* |banks| */
						},
						Sequence: 0xffffffff,
					},
				},
				TxOut: []*btcwire.TxOut{
					{
						Value: 0x12a05f200,
						PkScript: []byte{
							0x41, 0x04, 0x67, 0x8a, 0xfd, 0xb0, 0xfe, 0x55, /* |A.g....U| */
							0x48, 0x27, 0x19, 0x67, 0xf1, 0xa6, 0x71, 0x30, /* |H'.g..q0| */
							0xb7, 0x10, 0x5c, 0xd6, 0xa8, 0x28, 0xe0, 0x39, /* |..\..(.9| */
							0x09, 0xa6, 0x79, 0x62, 0xe0, 0xea, 0x1f, 0x61, /* |..yb...a| */
							0xde, 0xb6, 0x49, 0xf6, 0xbc, 0x3f, 0x4c, 0xef, /* |..I..?L.| */
							0x38, 0xc4, 0xf3, 0x55, 0x04, 0xe5, 0x1e, 0xc1, /* |8..U....| */
							0x12, 0xde, 0x5c, 0x38, 0x4d, 0xf7, 0xba, 0x0b, /* |..\8M...| */
							0x8d, 0x57, 0x8a, 0x4c, 0x70, 0x2b, 0x6b, 0xf1, /* |.W.Lp+k.| */
							0x1d, 0x5f, 0xac, /* |._.| */
						},
					},
				},
				LockTime: 0,
			},
		},
	}

	SolveBlock(&genesisBlock, 0)
	fmt.Println("Genesis block", genesisBlock)
}
