package ecc

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type TxFetcher struct {
	cache map[int]*Tx
}

func (Tf *TxFetcher) getUrl(testnet bool) string {
	if bool(testnet) {
		return "https://blockstream.info/testnet/api/"
	} else {
		return "https://blockstream.info/api/"
	}
}

func (Tf *TxFetcher) fetch(txId int, testnet bool, fresh bool) *Tx {
	testnet = false
	fresh = false

	if bool(fresh) || func() int {
		for i, v := range Tf.cache {
			if v == txId {
				return i
			}
		}
		return -1
	}() == -1 {
		url := fmt.Sprintf("%s/tx/%x/hex", Tf.getUrl(testnet), txId)
		var raw []byte
		response, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()
		func() {
			defer func() {
				if r := recover(); r != nil {
					if err, ok := r.(error); ok {
						if strings.HasPrefix(err.Error(), "ValueError") {
							panic(
								fmt.Errorf(
									"ValueError: %v",
									"unexpected response: %s", func() string {
										body, err := ioutil.ReadAll(response.Body)
										if err != nil {
											panic(err)
										}
										return string(body)
									}(),
								),
							)
							return
						}
					}
					panic(r)
				}
			}()
			raw = []byte(strings.TrimSpace(func() string {
				body, err := ioutil.ReadAll(response.Body)
				if err != nil {
					panic(err)
				}
				return string(body)
			}()))
		}()
		/*data, err := hex.DecodeString(raw)
		if err != nil {
			panic(err)
		}
		b := fmt.Sprintf("% x", data) //from hex string to byte array
		raw = []byte(b)*/
		if raw[4] == 0 {
			raw := append(raw[:4], raw[6:]...)
			tx := new(Tx)
			tx = tx.parse(raw, testnet)
			tx.locktime = littleEndianToInt(raw[len(raw)-4:])
		} else {
			tx := new(Tx)
			tx = tx.parse(raw, testnet)
		}
		tx := new(Tx)
		if !reflect.DeepEqual(tx.id(), txId) {
			panic(fmt.Errorf("ValueError: %v", "not the same id: %s vs %d", tx.id(), txId))
		}
		Tf.cache[txId] = tx
	}
	Tf.cache[txId].testnet = testnet
	return Tf.cache[txId]
}

func (Tf *TxFetcher) dumpCache(filename string) {
	f := func() *os.File {
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0o777)
		if err != nil {
			panic(err)
		}
		return f
	}()
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	to_dump := func() (d map[interface{}]interface{}) {
		d = make(map[interface{}]interface{})
		for k, tx := range Tf.cache {
			d[k] = hex.EncodeToString(tx.serialize())
		}
		return
	}()
	s := json.dumps(to_dump, true, 4)
	func() int {
		n, err := f.WriteString(s)
		if err != nil {
			panic(err)
		}
		return n
	}()
}

func (cls *TxFetcher) loadCache(filename string) {
	//need to convert this python line to go
	//disk_cache = json.loads(open(filename, 'r').read())
	for k, rawHex := range diskCache {
		raw := []byte(rawHex)
		if raw[4] == 0 {
			raw = raw[:4] + raw[6:]
			tx := new(Tx)
			tx = tx.parse(raw)
			tx.locktime = littleEndianToInt(raw[-4:])
		} else {
			tx := new(Tx)
			tx = tx.parse(raw)
		}
		tx := new(Tx)
		cls.cache[k] = tx
	}
}

type Tx struct {
	version  int64
	txIns    []*TxIn
	txOuts   []*TxOut
	locktime int64
	testnet  bool
}

func NewTx(version int64, txIns []*TxIn, txOuts []*TxOut, locktime int64, testnet bool) (T *Tx) {
	T = new(Tx)
	T.version = version
	T.txIns = txIns
	T.txOuts = txOuts
	T.locktime = locktime
	T.testnet = testnet
	return
}

func (T *Tx) Repr() (string, string, string, string, string, int64) {
	txIns := ""
	for _, txIn := range T.txIns {
		txIns += fmt.Sprintf("%#v", txIn) + "\n"
	}
	txOuts := ""
	for _, txOut := range T.txOuts {
		txOuts += fmt.Sprintf("%#v", txOut) + "\n"
	}
	return "tx: %s\nversion: %s\ntxIns: %s\ntxOuts: %s\nlocktime: %d",
		T.id(),
		string(T.version),
		txIns,
		txOuts,
		T.locktime
}

func (T *Tx) id() string {
	//Human-readable hexadecimal of the transaction hash"
	return hex.EncodeToString([]byte(T.hash())) //coverts string to hex string
}

func (T *Tx) hash() string {
	//Binary hash of the legacy serialization"
	return hash256(string(T.serialize()))
}

func (T *Tx) parse(s []byte, testnet bool) *Tx {
	var byt bytes.Buffer
	byt.Write(s)
	x, _ := byt.ReadBytes(4)
	version := littleEndianToInt(x)
	numInputs := readVarint(s)
	var inputs []*TxIn
	for i := 0; i < int(numInputs); i++ {
		Ti := new(TxIn)
		inputs = append(inputs, Ti.parse(s))
	}
	numOutputs := readVarint(s)
	var outputs []*TxOut
	for i := 0; i < int(numOutputs); i++ {
		To := new(TxOut)
		outputs = append(outputs, To.parse(s))
	}
	testnet = false
	y, _ := byt.ReadBytes(4)
	locktime := littleEndianToInt(y)
	return NewTx(version, inputs, outputs, locktime, testnet)
}

func (T *Tx) serialize() []byte {
	//"Returns the byte serialization of the transaction"
	result := intToLittleEndian(int(T.version), 4)
	result = append(result, encodeVarint(len(T.txIns))...)
	for _, txIn := range T.txIns {
		result = append(result, txIn.serialize()...)
	}
	result = append(result, encodeVarint(len(T.txOuts))...)
	for _, txOut := range T.txOuts {
		result = append(result, txOut.serialize()...)
	}
	result = append(result, intToLittleEndian(int(T.locktime), 4)...)
	return result
}

func (T *Tx) fee(testnet bool) int {
	testnet = false
	inputSum, outputSum := 0, 0
	for _, txIn := range T.txIns {
		inputSum += txIn.value(testnet)
	}
	for _, txOut := range T.txOuts {
		outputSum += int(txOut.amount)
	}
	return inputSum - outputSum
}

type TxIn struct {
	prevTx    []byte
	prevIndex int64
	scriptSig *Script
	sequence  int64
}

func NewTxIn(prevTx []byte, prevIndex int64, scriptSig *Script, sequence int64) (TI *TxIn) {
	scriptSig = nil
	sequence = 4294967295
	Ti := new(TxIn)
	Ti.prevTx = prevTx
	Ti.prevIndex = prevIndex
	if &scriptSig == nil {
		Ti.scriptSig = new(Script)
	} else {
		Ti.scriptSig = scriptSig
	}
	Ti.sequence = sequence
	return
}

func (Ti *TxIn) Repr() string {
	return fmt.Sprintf("%s:%d", hex.EncodeToString([]byte(Ti.prevTx)), Ti.prevIndex)
}

func (Ti *TxIn) parse(s []byte) *TxIn {
	//"Takes a byte stream and parses the tx_input at the start.
	//Returns a TxIn object.
	var byt bytes.Buffer
	byt.Write(s)
	S := new(Script)
	prevTx, _ := byt.ReadBytes(32)
	x, _ := byt.ReadBytes(4)
	prevIndex := littleEndianToInt(x)
	scriptSig := S.parse(s)
	z, _ := byt.ReadBytes(4)
	sequence := littleEndianToInt(z)
	return NewTxIn(prevTx, prevIndex, scriptSig, sequence)
}

func (s *TxIn) serialize() []byte {
	//"Returns the byte serialization of the transaction input"
	result := s.prevTx[:]
	result = append(result, intToLittleEndian(int(s.prevIndex), 4)...)
	result = append(result, s.scriptSig.serialize()...)
	result = append(result, intToLittleEndian(int(s.sequence), 4)...)
	return result
}

func (Ti *TxIn) fetchTx(testnet bool) *Tx {
	testnet = false
	txf := new(TxFetcher)
	intVar, _ := strconv.Atoi(hex.EncodeToString(Ti.prevTx))
	return txf.fetch(intVar, testnet, false)
}

func (Ti *TxIn) value(testnet bool) int {
	//"Get the output value by looking up the tx hash.Returns the amount in satoshi.
	testnet = false
	tx := Ti.fetchTx(false)
	return int(tx.txOuts[Ti.prevIndex].amount)
}

func (Ti *TxIn) scriptPubkey(testnet bool) *Script {
	//"Get the ScriptPubKey by looking up the tx hash.Returns a Script object.
	testnet = false
	tx := Ti.fetchTx(false)
	return tx.txOuts[Ti.prevIndex].scriptPubkey
}

type TxOut struct {
	amount       int64
	scriptPubkey *Script
}

func NewTxOut(amount int64, scriptPubkey *Script) (To *TxOut) {
	To = new(TxOut)
	To.amount = amount
	To.scriptPubkey = scriptPubkey
	return
}

func (To *TxOut) Repr() string {
	return fmt.Sprintf("%d:%s", To.amount, To.scriptPubkey)
}

func (To *TxOut) parse(s []byte) *TxOut {
	//"Takes a byte stream and parses the tx_output at the start.Returns a TxOut object.
	var byt bytes.Buffer
	byt.Write(s)
	x, _ := byt.ReadBytes(8)
	amount := littleEndianToInt(x)
	S := new(Script)
	scriptPubkey := S.parse(s)
	return NewTxOut(amount, scriptPubkey)
}

func (To *TxOut) serialize() []byte {
	//"Returns the byte serialization of the transaction output"
	result := intToLittleEndian(int(To.amount), 8)
	result = append(result, To.scriptPubkey.serialize()...)
	return result
}
