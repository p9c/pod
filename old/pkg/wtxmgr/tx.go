package wtxmgr

import (
	"bytes"
	amount2 "github.com/p9c/pod/pkg/amt"
	"github.com/p9c/pod/pkg/chaincfg"
	"time"
	
	"github.com/p9c/pod/pkg/blockchain"
	"github.com/p9c/pod/pkg/chainhash"
	"github.com/p9c/pod/pkg/walletdb"
	"github.com/p9c/pod/pkg/wire"
)

type (
	// Block contains the minimum amount of data to uniquely identify any block on either the best or side chain.
	Block struct {
		Hash   chainhash.Hash
		Height int32
	}
	// BlockMeta contains the unique identification for a block and any metadata pertaining to the block. At the moment,
	// this additional metadata only includes the block time from the block header.
	BlockMeta struct {
		Block
		Time time.Time
	}
	// blockRecord is an in-memory representation of the block record saved in the database.
	blockRecord struct {
		Block
		Time         time.Time
		transactions []chainhash.Hash
	}
	// incidence records the block hash and blockchain height of a mined transaction. Since a transaction hash alone is
	// not enough to uniquely identify a mined transaction (duplicate transaction hashes are allowed), the incidence is
	// used instead.
	incidence struct {
		txHash chainhash.Hash
		block  Block
	}
	// indexedIncidence records the transaction incidence and an input or output index.
	indexedIncidence struct {
		incidence
		index uint32
	}
	// debit records the debits a transaction record makes from previous wallet transaction credits.
	debit struct {
		txHash chainhash.Hash
		index  uint32
		amount amount2.Amount
		spends indexedIncidence
	}
	// credit describes a transaction output which was or is spendable by
	// wallet.
	credit struct {
		outPoint wire.OutPoint
		block    Block
		amount   amount2.Amount
		change   bool
		spentBy  indexedIncidence // Index == ^uint32(0) if unspent
	}
	// TxRecord represents a transaction managed by the Store.
	TxRecord struct {
		MsgTx        wire.MsgTx
		Hash         chainhash.Hash
		Received     time.Time
		SerializedTx []byte // Optional: may be nil
	}
	// Credit is the type representing a transaction output which was spent or is still spendable by wallet. A UTXO is
	// an unspent Credit, but not all Credits are UTXOs.
	Credit struct {
		wire.OutPoint
		BlockMeta
		Amount       amount2.Amount
		PkScript     []byte
		Received     time.Time
		FromCoinBase bool
	}
	// Store implements a transaction store for storing and managing wallet transactions.
	Store struct {
		chainParams *chaincfg.Params
		// Event callbacks. These execute in the same goroutine as the wtxmgr caller.
		NotifyUnspent func(hash *chainhash.Hash, index uint32)
	}
)

// NewTxRecord creates a new transaction record that may be inserted into the store. It uses memoization to save the
// transaction hash and the serialized transaction.
func NewTxRecord(serializedTx []byte, received time.Time) (*TxRecord, error) {
	rec := &TxRecord{
		Received:     received,
		SerializedTx: serializedTx,
	}
	e := rec.MsgTx.Deserialize(bytes.NewReader(serializedTx))
	if e != nil {
		str := "failed to deserialize transaction"
		return nil, storeError(ErrInput, str, e)
	}
	copy(rec.Hash[:], chainhash.DoubleHashB(serializedTx))
	return rec, nil
}

// NewTxRecordFromMsgTx creates a new transaction record that may be inserted into the store.
func NewTxRecordFromMsgTx(msgTx *wire.MsgTx, received time.Time) (*TxRecord, error) {
	buf := bytes.NewBuffer(make([]byte, 0, msgTx.SerializeSize()))
	e := msgTx.Serialize(buf)
	if e != nil {
		str := "failed to serialize transaction"
		return nil, storeError(ErrInput, str, e)
	}
	rec := &TxRecord{
		MsgTx:        *msgTx,
		Received:     received,
		SerializedTx: buf.Bytes(),
		Hash:         msgTx.TxHash(),
	}
	return rec, nil
}

// DoUpgrades performs any necessary upgrades to the transaction history contained in the wallet database, namespaced by
// the top level bucket key namespaceKey.
func DoUpgrades(db walletdb.DB, namespaceKey []byte) (e error) {
	// No upgrades
	return nil
}

// Open opens the wallet transaction store from a walletdb namespace.
// If the store does not exist, ErrNoExist is returned.
func Open(ns walletdb.ReadBucket, chainParams *chaincfg.Params) (*Store, error) {
	// Open the store.
	e := openStore(ns)
	if e != nil {
		return nil, e
	}
	s := &Store{chainParams, nil} // TODO: set callbacks
	return s, nil
}

// Create creates a new persistent transaction store in the walletdb namespace. Creating the store when one already
// exists in this namespace will error with ErrAlreadyExists.
func Create(ns walletdb.ReadWriteBucket) (e error) {
	return createStore(ns)
}

// updateMinedBalance updates the mined balance within the store, if changed, after processing the given transaction
// record.
func (s *Store) updateMinedBalance(
	ns walletdb.ReadWriteBucket, rec *TxRecord,
	block *BlockMeta,
) (e error) {
	// Fetch the mined balance in case we need to update it.
	minedBalance, e := fetchMinedBalance(ns)
	if e != nil {
		return e
	}
	// Add a debit record for each unspent credit spent by this transaction. The index is set in each iteration below.
	spender := indexedIncidence{
		incidence: incidence{
			txHash: rec.Hash,
			block:  block.Block,
		},
	}
	newMinedBalance := minedBalance
	for i, input := range rec.MsgTx.TxIn {
		unspentKey, credKey := existsUnspent(ns, &input.PreviousOutPoint)
		if credKey == nil {
			// Debits for unmined transactions are not explicitly tracked. Instead, all previous outputs spent by any
			// unmined transaction are added to a map for quick lookups when it must be checked whether a mined output
			// is unspent or not.
			//
			// Tracking individual debits for unmined transactions could be added later to simplify (and increase
			// performance of) determining some details that need the previous outputs (e.g. determining a fee), but at
			// the moment that is not done (and a db lookup is used for those cases instead).
			//
			// There is also a good chance that all unmined transaction handling will move entirely to the db rather
			// than being handled in memory for atomicity reasons, so the simplist implementation is currently used.
			continue
		}
		// If this output is relevant to us, we'll mark the it as spent and remove its amount from the store.
		spender.index = uint32(i)
		amt, e := spendCredit(ns, credKey, &spender)
		if e != nil {
			return e
		}
		e = putDebit(
			ns, &rec.Hash, uint32(i), amt, &block.Block, credKey,
		)
		if e != nil {
			return e
		}
		if e := deleteRawUnspent(ns, unspentKey); E.Chk(e) {
			return e
		}
		newMinedBalance -= amt
	}
	// For each output of the record that is marked as a credit, if the output is marked as a credit by the unconfirmed
	// store, remove the marker and mark the output as a credit in the db.
	//
	// Moved credits are added as unspents, even if there is another unconfirmed transaction which spends them.
	cred := credit{
		outPoint: wire.OutPoint{Hash: rec.Hash},
		block:    block.Block,
		spentBy:  indexedIncidence{index: ^uint32(0)},
	}
	it := makeUnminedCreditIterator(ns, &rec.Hash)
	for it.next() {
		// TODO: This should use the raw apis. The credit value (it.cv) can be moved from unmined directly to the
		//  credits bucket. The key needs a modification to include the block height/hash.
		index, e := fetchRawUnminedCreditIndex(it.ck)
		if e != nil {
			return e
		}
		amount, change, e := fetchRawUnminedCreditAmountChange(it.cv)
		if e != nil {
			return e
		}
		cred.outPoint.Index = index
		cred.amount = amount
		cred.change = change
		if e := putUnspentCredit(ns, &cred); E.Chk(e) {
			return e
		}
		e = putUnspent(ns, &cred.outPoint, &block.Block)
		if e != nil {
			return e
		}
		newMinedBalance += amount
	}
	if it.err != nil {
		return it.err
	}
	// Update the balance if it has changed.
	if newMinedBalance != minedBalance {
		return putMinedBalance(ns, newMinedBalance)
	}
	return nil
}

// deleteUnminedTx deletes an unmined transaction from the store.
//
// NOTE: This should only be used once the transaction has been mined.
func (s *Store) deleteUnminedTx(ns walletdb.ReadWriteBucket, rec *TxRecord) (e error) {
	for i := range rec.MsgTx.TxOut {
		k := canonicalOutPoint(&rec.Hash, uint32(i))
		if e := deleteRawUnminedCredit(ns, k); E.Chk(e) {
			return e
		}
	}
	return deleteRawUnmined(ns, rec.Hash[:])
}

// InsertTx records a transaction as belonging to a wallet's transaction history. If block is nil, the transaction is
// considered unspent, and the transaction's index must be unset.
func (s *Store) InsertTx(ns walletdb.ReadWriteBucket, rec *TxRecord, block *BlockMeta) (e error) {
	if block == nil {
		return s.insertMemPoolTx(ns, rec)
	}
	return s.insertMinedTx(ns, rec, block)
}

// RemoveUnminedTx attempts to remove an unmined transaction from the transaction store. This is to be used in the
// scenario that a transaction that we attempt to rebroadcast, turns out to double spend one of our existing inputs.
// This function we remove the conflicting transaction identified by the tx record, and also recursively remove all
// transactions that depend on it.
func (s *Store) RemoveUnminedTx(ns walletdb.ReadWriteBucket, rec *TxRecord) (e error) {
	// As we already have a tx record, we can directly call the RemoveConflict method. This will do the job of
	// recursively removing this unmined transaction, and any transactions that depend on it.
	return RemoveConflict(ns, rec)
}

// insertMinedTx inserts a new transaction record for a mined transaction into the database under the confirmed bucket.
// It guarantees that, if the tranasction was previously unconfirmed, then it will take care of cleaning up the
// unconfirmed state. All other unconfirmed double spend attempts will be removed as well.
func (s *Store) insertMinedTx(
	ns walletdb.ReadWriteBucket, rec *TxRecord,
	block *BlockMeta,
) (e error) {
	// If a transaction record for this hash and block already exists, we can exit early.
	if _, v := existsTxRecord(ns, &rec.Hash, &block.Block); v != nil {
		return nil
	}
	// If a block record does not yet exist for any transactions from this block, insert a block record first.
	// Otherwise, update it by adding the transaction hash to the set of transactions from this block.
	blockKey, blockValue := existsBlockRecord(ns, block.Height)
	if blockValue == nil {
		e = putBlockRecord(ns, block, &rec.Hash)
	} else {
		blockValue, e = appendRawBlockRecord(blockValue, &rec.Hash)
		if e != nil {
			return e
		}
		e = putRawBlockRecord(ns, blockKey, blockValue)
	}
	if e != nil {
		return e
	}
	if e := putTxRecord(ns, rec, &block.Block); E.Chk(e) {
		return e
	}
	// Determine if this transaction has affected our balance, and if so, update it.
	if e := s.updateMinedBalance(ns, rec, block); E.Chk(e) {
		return e
	}
	// If this transaction previously existed within the store as unmined, we'll need to remove it from the unmined
	// bucket.
	if v := existsRawUnmined(ns, rec.Hash[:]); v != nil {
		I.F("marking unconfirmed transaction %v mined in block %d", &rec.Hash, block.Height)
		if e := s.deleteUnminedTx(ns, rec); E.Chk(e) {
			return e
		}
	}
	// As there may be unconfirmed transactions that are invalidated by this transaction (either being duplicates, or
	// double spends), remove them from the unconfirmed set. This also handles removing unconfirmed transaction spend
	// chains if any other unconfirmed transactions spend outputs of the removed double spend.
	return s.removeDoubleSpends(ns, rec)
}

// AddCredit marks a transaction record as containing a transaction output spendable by wallet. The output is added
// unspent, and is marked spent when a new transaction spending the output is inserted into the store.
//
// TODO(jrick): This should not be necessary. Instead, pass the indexes that are known to contain credits when a
//  transaction or merkleblock is inserted into the store.
func (s *Store) AddCredit(
	ns walletdb.ReadWriteBucket,
	rec *TxRecord,
	block *BlockMeta,
	index uint32,
	change bool,
) (e error) {
	if int(index) >= len(rec.MsgTx.TxOut) {
		str := "transaction output does not exist"
		return storeError(ErrInput, str, nil)
	}
	isNew, e := s.addCredit(ns, rec, block, index, change)
	if e == nil && isNew && s.NotifyUnspent != nil {
		s.NotifyUnspent(&rec.Hash, index)
	}
	return e
}

// addCredit is an AddCredit helper that runs in an update transaction. The bool return specifies whether the unspent
// output is newly added ( true) or a duplicate (false).
func (s *Store) addCredit(
	ns walletdb.ReadWriteBucket,
	rec *TxRecord,
	block *BlockMeta,
	index uint32,
	change bool,
) (bool, error) {
	if block == nil {
		// If the outpoint that we should mark as credit already exists within the store, either as unconfirmed or
		// confirmed, then we have nothing left to do and can exit.
		k := canonicalOutPoint(&rec.Hash, index)
		if existsRawUnminedCredit(ns, k) != nil {
			return false, nil
		}
		if existsRawUnspent(ns, k) != nil {
			return false, nil
		}
		v := valueUnminedCredit(amount2.Amount(rec.MsgTx.TxOut[index].Value), change)
		return true, putRawUnminedCredit(ns, k, v)
	}
	k, v := existsCredit(ns, &rec.Hash, index, &block.Block)
	if v != nil {
		return false, nil
	}
	txOutAmt := amount2.Amount(rec.MsgTx.TxOut[index].Value)
	T.F(
		"marking transaction %v output %d (%v) spendable",
		rec.Hash, index, txOutAmt,
	)
	cred := credit{
		outPoint: wire.OutPoint{
			Hash:  rec.Hash,
			Index: index,
		},
		block:   block.Block,
		amount:  txOutAmt,
		change:  change,
		spentBy: indexedIncidence{index: ^uint32(0)},
	}
	v = valueUnspentCredit(&cred)
	e := putRawCredit(ns, k, v)
	if e != nil {
		return false, e
	}
	minedBalance, e := fetchMinedBalance(ns)
	if e != nil {
		return false, e
	}
	e = putMinedBalance(ns, minedBalance+txOutAmt)
	if e != nil {
		return false, e
	}
	return true, putUnspent(ns, &cred.outPoint, &block.Block)
}

// Rollback removes all blocks at height onwards, moving any transactions within each block to the unconfirmed pool.
func (s *Store) Rollback(ns walletdb.ReadWriteBucket, height int32) (e error) {
	return s.rollback(ns, height)
}
func (s *Store) rollback(ns walletdb.ReadWriteBucket, height int32) (e error) {
	minedBalance, e := fetchMinedBalance(ns)
	if e != nil {
		return e
	}
	// Keep track of all credits that were removed from coinbase transactions. After detaching all blocks, if any
	// transaction record exists in unmined that spends these outputs, remove them and their spend chains.
	//
	// It is necessary to keep these in memory and fix the unmined transactions later since blocks are removed in
	// increasing order.
	var coinBaseCredits []wire.OutPoint
	var heightsToRemove []int32
	it := makeReverseBlockIterator(ns)
	for it.prev() {
		b := &it.elem
		if it.elem.Height < height {
			break
		}
		heightsToRemove = append(heightsToRemove, it.elem.Height)
		T.F("rolling back %d transactions from block %v height %d", len(b.transactions), b.Hash, b.Height)
		for i := range b.transactions {
			txHash := &b.transactions[i]
			recKey := keyTxRecord(txHash, &b.Block)
			recVal := existsRawTxRecord(ns, recKey)
			var rec TxRecord
			e = readRawTxRecord(txHash, recVal, &rec)
			if e != nil {
				return e
			}
			e = deleteTxRecord(ns, txHash, &b.Block)
			if e != nil {
				return e
			}
			// Handle coinbase transactions specially since they are not moved to the unconfirmed store. A coinbase
			// cannot contain any debits, but all credits should be removed and the mined balance decremented.
			if blockchain.IsCoinBaseTx(&rec.MsgTx) {
				op := wire.OutPoint{Hash: rec.Hash}
				for i, output := range rec.MsgTx.TxOut {
					k, v := existsCredit(
						ns, &rec.Hash,
						uint32(i), &b.Block,
					)
					if v == nil {
						continue
					}
					op.Index = uint32(i)
					coinBaseCredits = append(coinBaseCredits, op)
					unspentKey, credKey := existsUnspent(ns, &op)
					if credKey != nil {
						minedBalance -= amount2.Amount(output.Value)
						e = deleteRawUnspent(ns, unspentKey)
						if e != nil {
							return e
						}
					}
					e = deleteRawCredit(ns, k)
					if e != nil {
						return e
					}
				}
				continue
			}
			e = putRawUnmined(ns, txHash[:], recVal)
			if e != nil {
				return e
			}
			// For each debit recorded for this transaction, mark the credit it spends as unspent (as long as it still
			// exists) and delete the debit. The previous output is recorded in the unconfirmed store for every previous
			// output, not just debits.
			for i, input := range rec.MsgTx.TxIn {
				prevOut := &input.PreviousOutPoint
				prevOutKey := canonicalOutPoint(
					&prevOut.Hash,
					prevOut.Index,
				)
				e = putRawUnminedInput(ns, prevOutKey, rec.Hash[:])
				if e != nil {
					return e
				}
				// If this input is a debit, remove the debit record and mark the credit that it spent as unspent,
				// incrementing the mined balance.
				debKey, credKey, e := existsDebit(
					ns,
					&rec.Hash, uint32(i), &b.Block,
				)
				if e != nil {
					return e
				}
				if debKey == nil {
					continue
				}
				// unspendRawCredit does not error in case the no credit exists for this key, but this behavior is
				// correct. Since blocks are removed in increasing order, this credit may have already been removed from
				// a previously removed transaction record in this rollback.
				var amt amount2.Amount
				amt, e = unspendRawCredit(ns, credKey)
				if e != nil {
					return e
				}
				e = deleteRawDebit(ns, debKey)
				if e != nil {
					return e
				}
				// If the credit was previously removed in the rollback, the credit amount is zero. Only mark the
				// previously spent credit as unspent if it still exists.
				if amt == 0 {
					continue
				}
				unspentVal, e := fetchRawCreditUnspentValue(credKey)
				if e != nil {
					return e
				}
				minedBalance += amt
				e = putRawUnspent(ns, prevOutKey, unspentVal)
				if e != nil {
					return e
				}
			}
			// For each detached non-coinbase credit, move the credit output to unmined. If the credit is marked
			// unspent, it is removed from the utxo set and the mined balance is decremented.
			//
			// TODO: use a credit iterator
			for i, output := range rec.MsgTx.TxOut {
				k, v := existsCredit(
					ns, &rec.Hash, uint32(i),
					&b.Block,
				)
				if v == nil {
					continue
				}
				amt, change, e := fetchRawCreditAmountChange(v)
				if e != nil {
					return e
				}
				outPointKey := canonicalOutPoint(&rec.Hash, uint32(i))
				unminedCredVal := valueUnminedCredit(amt, change)
				e = putRawUnminedCredit(ns, outPointKey, unminedCredVal)
				if e != nil {
					return e
				}
				e = deleteRawCredit(ns, k)
				if e != nil {
					return e
				}
				credKey := existsRawUnspent(ns, outPointKey)
				if credKey != nil {
					minedBalance -= amount2.Amount(output.Value)
					e = deleteRawUnspent(ns, outPointKey)
					if e != nil {
						return e
					}
				}
			}
		}
		// reposition cursor before deleting this k/v pair and advancing to the previous.
		it.reposition(it.elem.Height)
		// Avoid cursor deletion until bolt issue #620 is resolved.
		//
		// e = it.delete() if e != nil  {
		// 	return e
		// }
	}
	if it.err != nil {
		return it.err
	}
	// Delete the block records outside of the iteration since cursor deletion is broken.
	for _, h := range heightsToRemove {
		e = deleteBlockRecord(ns, h)
		if e != nil {
			return e
		}
	}
	for _, op := range coinBaseCredits {
		opKey := canonicalOutPoint(&op.Hash, op.Index)
		unminedSpendTxHashKeys := fetchUnminedInputSpendTxHashes(ns, opKey)
		for _, unminedSpendTxHashKey := range unminedSpendTxHashKeys {
			unminedVal := existsRawUnmined(ns, unminedSpendTxHashKey[:])
			// If the spending transaction spends multiple outputs
			// from the same transaction, we'll find duplicate
			// entries within the store, so it's possible we're
			// unable to find it if the conflicts have already been
			// removed in a previous iteration.
			if unminedVal == nil {
				continue
			}
			var unminedRec TxRecord
			unminedRec.Hash = unminedSpendTxHashKey
			e = readRawTxRecord(&unminedRec.Hash, unminedVal, &unminedRec)
			if e != nil {
				return e
			}
			D.F(
				"transaction %v spends a removed coinbase output -- removing as well %s",
				unminedRec.Hash,
			)
			e = RemoveConflict(ns, &unminedRec)
			if e != nil {
				return e
			}
		}
	}
	return putMinedBalance(ns, minedBalance)
}

func // UnspentOutputs returns all unspent received transaction outputs.
// The order is undefined.
(s *Store) UnspentOutputs(ns walletdb.ReadBucket) ([]Credit, error) {
	var unspent []Credit
	var op wire.OutPoint
	var block Block
	e := ns.NestedReadBucket(bucketUnspent).ForEach(
		func(k, v []byte) (e error) {
			e = readCanonicalOutPoint(k, &op)
			if e != nil {
				return e
			}
			if existsRawUnminedInput(ns, k) != nil {
				// Output is spent by an unmined transaction.
				// Skip this k/v pair.
				return nil
			}
			e = readUnspentBlock(v, &block)
			if e != nil {
				return e
			}
			blockTime, e := fetchBlockTime(ns, block.Height)
			if e != nil {
				return e
			}
			// TODO(jrick): reading the entire transaction should
			// be avoidable.  Creating the credit only requires the
			// output amount and pkScript.
			rec, e := fetchTxRecord(ns, &op.Hash, &block)
			if e != nil {
				return e
			}
			txOut := rec.MsgTx.TxOut[op.Index]
			cred := Credit{
				OutPoint: op,
				BlockMeta: BlockMeta{
					Block: block,
					Time:  blockTime,
				},
				Amount:       amount2.Amount(txOut.Value),
				PkScript:     txOut.PkScript,
				Received:     rec.Received,
				FromCoinBase: blockchain.IsCoinBaseTx(&rec.MsgTx),
			}
			unspent = append(unspent, cred)
			return nil
		},
	)
	if e != nil {
		if _, ok := e.(TxMgrError); ok {
			return nil, e
		}
		str := "failed iterating unspent bucket"
		return nil, storeError(ErrDatabase, str, e)
	}
	e = ns.NestedReadBucket(bucketUnminedCredits).ForEach(
		func(k, v []byte) (e error) {
			if existsRawUnminedInput(ns, k) != nil {
				// Output is spent by an unmined transaction.
				// Skip to next unmined credit.
				return nil
			}
			e = readCanonicalOutPoint(k, &op)
			if e != nil {
				return e
			}
			// TODO(jrick): Reading/parsing the entire transaction record
			// just for the output amount and script can be avoided.
			recVal := existsRawUnmined(ns, op.Hash[:])
			var rec TxRecord
			e = readRawTxRecord(&op.Hash, recVal, &rec)
			if e != nil {
				return e
			}
			txOut := rec.MsgTx.TxOut[op.Index]
			cred := Credit{
				OutPoint: op,
				BlockMeta: BlockMeta{
					Block: Block{Height: -1},
				},
				Amount:       amount2.Amount(txOut.Value),
				PkScript:     txOut.PkScript,
				Received:     rec.Received,
				FromCoinBase: blockchain.IsCoinBaseTx(&rec.MsgTx),
			}
			unspent = append(unspent, cred)
			return nil
		},
	)
	if e != nil {
		if _, ok := e.(TxMgrError); ok {
			return nil, e
		}
		str := "failed iterating unmined credits bucket"
		return nil, storeError(ErrDatabase, str, e)
	}
	return unspent, nil
}

func // Balance returns the spendable wallet balance (total value of all unspent
// transaction outputs) given a minimum of minConf confirmations, calculated
// at a current chain height of curHeight.  Coinbase outputs are only included
// in the balance if maturity has been reached.
//
// Balance may return unexpected results if syncHeight is lower than the block
// height of the most recent mined transaction in the store.
(s *Store) Balance(ns walletdb.ReadBucket, minConf int32, syncHeight int32) (amount2.Amount, error) {
	bal, e := fetchMinedBalance(ns)
	if e != nil {
		return 0, e
	}
	// Subtract the balance for each credit that is spent by an unmined
	// transaction.
	var op wire.OutPoint
	var block Block
	e = ns.NestedReadBucket(bucketUnspent).ForEach(
		func(k, v []byte) (e error) {
			e = readCanonicalOutPoint(k, &op)
			if e != nil {
				return e
			}
			e = readUnspentBlock(v, &block)
			if e != nil {
				return e
			}
			if existsRawUnminedInput(ns, k) != nil {
				_, v := existsCredit(ns, &op.Hash, op.Index, &block)
				amt, e := fetchRawCreditAmount(v)
				if e != nil {
					return e
				}
				bal -= amt
			}
			return nil
		},
	)
	if e != nil {
		if _, ok := e.(TxMgrError); ok {
			return 0, e
		}
		str := "failed iterating unspent outputs"
		return 0, storeError(ErrDatabase, str, e)
	}
	// Decrement the balance for any unspent credit with less than
	// minConf confirmations and any (unspent) immature coinbase credit.
	coinbaseMaturity := int32(s.chainParams.CoinbaseMaturity)
	stopConf := minConf
	if coinbaseMaturity > stopConf {
		stopConf = coinbaseMaturity
	}
	lastHeight := syncHeight - stopConf
	blockIt := makeReadReverseBlockIterator(ns)
	for blockIt.prev() {
		block := &blockIt.elem
		if block.Height < lastHeight {
			break
		}
		for i := range block.transactions {
			txHash := &block.transactions[i]
			rec, e := fetchTxRecord(ns, txHash, &block.Block)
			if e != nil {
				return 0, e
			}
			numOuts := uint32(len(rec.MsgTx.TxOut))
			for i := uint32(0); i < numOuts; i++ {
				// Avoid double decrementing the credit amount
				// if it was already removed for being spent by
				// an unmined tx.
				opKey := canonicalOutPoint(txHash, i)
				if existsRawUnminedInput(ns, opKey) != nil {
					continue
				}
				_, v := existsCredit(ns, txHash, i, &block.Block)
				if v == nil {
					continue
				}
				amt, spent, e := fetchRawCreditAmountSpent(v)
				if e != nil {
					return 0, e
				}
				if spent {
					continue
				}
				confs := syncHeight - block.Height + 1
				if confs < minConf || (blockchain.IsCoinBaseTx(&rec.MsgTx) &&
					confs < coinbaseMaturity) {
					bal -= amt
				}
			}
		}
	}
	if blockIt.err != nil {
		return 0, blockIt.err
	}
	// If unmined outputs are included, increment the balance for each
	// output that is unspent.
	if minConf == 0 {
		e = ns.NestedReadBucket(bucketUnminedCredits).ForEach(
			func(k, v []byte) (e error) {
				if existsRawUnminedInput(ns, k) != nil {
					// Output is spent by an unmined transaction.
					// Skip to next unmined credit.
					return nil
				}
				amount, e := fetchRawUnminedCreditAmount(v)
				if e != nil {
					return e
				}
				bal += amount
				return nil
			},
		)
		if e != nil {
			if _, ok := e.(TxMgrError); ok {
				return 0, e
			}
			str := "failed to iterate over unmined credits bucket"
			return 0, storeError(ErrDatabase, str, e)
		}
	}
	return bal, nil
}