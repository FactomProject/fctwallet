package Wallet

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/FactomProject/FactomCode/common"
	"github.com/FactomProject/factoid"
	"github.com/FactomProject/factoid/block"
	"github.com/FactomProject/factom"
	"log"
	"time"
)

type DataStatusStruct struct {
	DBlockHeight   int
	LastKnownBlock string
}

func TimestampToString(timestamp uint64) string {
	blockTime := time.Unix(int64(timestamp), 0)
	return blockTime.Format("2006-01-02 15:04:05")
}
func ByteSliceToDecodedString(b []byte) DecodedString {
	var ds DecodedString
	ds.Encoded = fmt.Sprintf("%x", b)
	ds.Decoded = string(b)
	return ds
}

var DataStatus *DataStatusStruct

const DataStatusBucket string = "DataStatus"

var BucketList []string = []string{DataStatusBucket}

type Common struct {
	ChainID   string
	Timestamp string

	JSONString   string
	SpewString   string
	BinaryString string
}

func (e *Common) JSON() (string, error) {
	return common.EncodeJSONString(e)
}

func (e *Common) Spew() string {
	return common.Spew(e)
}

type Block struct {
	Common

	FullHash    string //KeyMR
	PartialHash string

	PrevBlockHash string
	NextBlockHash string

	EntryCount int

	EntryList []*Entry

	IsAdminBlock       bool
	IsFactoidBlock     bool
	IsEntryCreditBlock bool
	IsEntryBlock       bool
}
type ListEntry struct {
	ChainID string
	KeyMR   string
}

type DBlock struct {
	DBHash string

	PrevBlockKeyMR string
	NextBlockKeyMR string
	TimeStamp      uint64
	SequenceNumber int

	EntryBlockList   []ListEntry
	AdminBlock       ListEntry
	FactoidBlock     ListEntry
	EntryCreditBlock ListEntry

	BlockTimeStr string
	KeyMR        string

	Blocks int

	AdminEntries       int
	EntryCreditEntries int
	FactoidEntries     int
	EntryEntries       int
}

func (e *DBlock) JSON() (string, error) {
	return common.EncodeJSONString(e)
}

func (e *DBlock) Spew() string {
	return common.Spew(e)
}

type DecodedString struct {
	Encoded string
	Decoded string
}

type Entry struct {
	Common

	ExternalIDs []DecodedString
	Content     DecodedString

	//Marshallable blocks
	Hash string
}

func (e *Entry) JSON() (string, error) {
	return common.EncodeJSONString(e)
}

func (e *Entry) Spew() string {
	return common.Spew(e)
}

func SaveDataStatus(ds *DataStatusStruct) error {
	err := SaveData(DataStatusBucket, DataStatusBucket, ds)
	if err != nil {
		return err
	}
	DataStatus = ds
	return nil
}

func LoadDataStatus() *DataStatusStruct {
	if DataStatus != nil {
		return DataStatus
	}
	ds := new(DataStatusStruct)
	var err error
	ds2, err := LoadData(DataStatusBucket, DataStatusBucket, ds)
	if err != nil {
		panic(err)
	}
	if ds2 == nil {
		ds = new(DataStatusStruct)
		ds.LastKnownBlock = "0000000000000000000000000000000000000000000000000000000000000000"
	}
	DataStatus = ds
	log.Printf("LoadDataStatus DS - %v, %v", ds, ds2)
	return ds
}

func EncodeJSONString(data interface{}) (string, error) {
	encoded, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(encoded), err
}

func Synchronize() error {
	log.Println("Synchronize()")
	head, err := factom.GetDBlockHead()
	if err != nil {
		return err
	}
	previousKeyMR := head.KeyMR
	dataStatus := LoadDataStatus()
	maxHeight := dataStatus.DBlockHeight
	for {
		body, err := GetDBlockFromFactom(previousKeyMR)
		if err != nil {
			return err
		}

		log.Printf("\n\nProcessing dblock number %v\n", body.SequenceNumber)

		str, err := EncodeJSONString(body)
		if err != nil {
			return err
		}
		log.Printf("%v", str)

		for _, v := range body.EntryBlockList {
			if v.ChainID == "000000000000000000000000000000000000000000000000000000000000000a" {
				continue
			}
			fetchedBlock, err := FetchBlock(v.ChainID, v.KeyMR, body.BlockTimeStr)
			if err != nil {
				return err
			}
			fmt.Printf("\nfetchedBlock - %v\n\n", fetchedBlock)
		}

		if maxHeight < body.SequenceNumber {
			maxHeight = body.SequenceNumber
		}
		previousKeyMR = body.PrevBlockKeyMR
		if previousKeyMR == "0000000000000000000000000000000000000000000000000000000000000000" {
			dataStatus.LastKnownBlock = head.KeyMR
			dataStatus.DBlockHeight = maxHeight
			break
		}

	}
	err = SaveDataStatus(dataStatus)
	if err != nil {
		return err
	}
	return nil
}

func FetchBlock(chainID, hash, blockTime string) (*Block, error) {
	block := new(Block)

	raw, err := factom.GetRaw(hash)
	if err != nil {
		return nil, err
	}
	switch chainID {
	case "000000000000000000000000000000000000000000000000000000000000000c":
		block, err = ParseEntryCreditBlock(chainID, hash, raw, blockTime)
		if err != nil {
			return nil, err
		}
		break
	case "000000000000000000000000000000000000000000000000000000000000000f":
		block, err = ParseFactoidBlock(chainID, hash, raw, blockTime)
		if err != nil {
			return nil, err
		}
		break
	default:
		block, err = ParseEntryBlock(chainID, hash, raw, blockTime)
		if err != nil {
			return nil, err
		}
		break
	}

	return block, nil
}

func ParseEntryCreditBlock(chainID, hash string, rawBlock []byte, blockTime string) (*Block, error) {
	answer := new(Block)

	ecBlock := common.NewECBlock()
	_, err := ecBlock.UnmarshalBinaryData(rawBlock)
	if err != nil {
		return nil, err
	}

	answer.ChainID = chainID
	h, err := ecBlock.Hash()
	if err != nil {
		return nil, err
	}
	answer.FullHash = h.String()

	h, err = ecBlock.HeaderHash()
	if err != nil {
		return nil, err
	}
	answer.PartialHash = h.String()

	answer.PrevBlockHash = ecBlock.Header.PrevLedgerKeyMR.String()

	answer.EntryCount = len(ecBlock.Body.Entries)
	answer.EntryList = make([]*Entry, answer.EntryCount)

	answer.BinaryString = fmt.Sprintf("%x", rawBlock)

	for i, v := range ecBlock.Body.Entries {
		entry := new(Entry)

		marshalled, err := v.MarshalBinary()
		if err != nil {
			return nil, err
		}
		entry.BinaryString = fmt.Sprintf("%x", marshalled)
		entry.Timestamp = blockTime
		entry.ChainID = chainID

		entry.Hash = fmt.Sprintf("%x", v.ECID())

		entry.JSONString, err = v.JSONString()
		if err != nil {
			return nil, err
		}
		entry.SpewString = v.Spew()

		answer.EntryList[i] = entry
	}

	answer.JSONString, err = ecBlock.JSONString()
	if err != nil {
		return nil, err
	}
	answer.SpewString = ecBlock.Spew()
	answer.IsEntryCreditBlock = true

	return answer, nil
}

func ParseFactoidBlock(chainID, hash string, rawBlock []byte, blockTime string) (*Block, error) {
	answer := new(Block)

	fBlock := new(block.FBlock)
	_, err := fBlock.UnmarshalBinaryData(rawBlock)
	if err != nil {
		return nil, err
	}

	answer.ChainID = chainID
	answer.PartialHash = fBlock.GetHash().String()
	answer.FullHash = fBlock.GetLedgerKeyMR().String()
	answer.PrevBlockHash = fmt.Sprintf("%x", fBlock.PrevKeyMR.Bytes())

	transactions := fBlock.GetTransactions()
	answer.EntryCount = len(transactions)
	answer.EntryList = make([]*Entry, answer.EntryCount)
	answer.BinaryString = fmt.Sprintf("%x", rawBlock)
	for i, v := range transactions {
		entry := new(Entry)
		bin, err := v.MarshalBinary()

		if err != nil {
			return nil, err
		}

		entry.BinaryString = fmt.Sprintf("%x", bin)
		entry.Timestamp = TimestampToString(v.GetMilliTimestamp() / 1000)
		entry.Hash = v.GetHash().String()
		entry.ChainID = chainID

		entry.JSONString, err = v.JSONString()
		if err != nil {
			return nil, err
		}
		entry.SpewString = v.Spew()

		answer.EntryList[i] = entry
	}
	answer.JSONString, err = fBlock.JSONString()
	if err != nil {
		return nil, err
	}
	answer.SpewString = fBlock.Spew()
	answer.IsFactoidBlock = true

	return answer, nil
}

func ParseEntryBlock(chainID, hash string, rawBlock []byte, blockTime string) (*Block, error) {
	answer := new(Block)

	eBlock := common.NewEBlock()
	_, err := eBlock.UnmarshalBinaryData(rawBlock)
	if err != nil {
		return nil, err
	}

	answer.ChainID = chainID
	h, err := eBlock.KeyMR()
	if err != nil {
		return nil, err
	}
	answer.PartialHash = h.String()
	if err != nil {
		return nil, err
	}
	h, err = eBlock.Hash()
	if err != nil {
		return nil, err
	}
	answer.FullHash = h.String()

	answer.PrevBlockHash = eBlock.Header.PrevKeyMR.String()

	answer.EntryCount = len(eBlock.Body.EBEntries)
	answer.EntryList = make([]*Entry, answer.EntryCount)
	answer.BinaryString = fmt.Sprintf("%x", rawBlock)

	for i, v := range eBlock.Body.EBEntries {
		entry, err := FetchAndParseEntry(v.String(), blockTime)
		if err != nil {
			return nil, err
		}

		answer.EntryList[i] = entry
	}
	answer.JSONString, err = eBlock.JSONString()
	if err != nil {
		return nil, err
	}
	answer.SpewString = eBlock.Spew()

	answer.IsEntryBlock = true

	return answer, nil
}

func FetchAndParseEntry(hash, blockTime string) (*Entry, error) {
	e := new(Entry)
	raw, err := factom.GetRaw(hash)
	if err != nil {
		return nil, err
	}

	entry := new(common.Entry)
	_, err = entry.UnmarshalBinaryData(raw)
	if err != nil {
		return nil, err
	}

	e.ChainID = entry.ChainID.String()
	e.Hash = hash
	str, err := entry.JSONString()
	if err != nil {
		return nil, err
	}
	e.JSONString = str
	e.SpewString = entry.Spew()
	e.BinaryString = fmt.Sprintf("%x", raw)
	e.Timestamp = blockTime

	e.Content = ByteSliceToDecodedString(entry.Content)
	e.ExternalIDs = make([]DecodedString, len(entry.ExtIDs))
	for i, v := range entry.ExtIDs {
		e.ExternalIDs[i] = ByteSliceToDecodedString(v)
	}

	return e, nil
}

func GetDBlockFromFactom(keyMR string) (*DBlock, error) {
	answer := new(DBlock)

	body, err := factom.GetDBlock(keyMR)
	if err != nil {
		return answer, err
	}

	answer = new(DBlock)
	answer.DBHash = body.DBHash
	answer.PrevBlockKeyMR = body.Header.PrevBlockKeyMR
	answer.TimeStamp = body.Header.Timestamp
	answer.SequenceNumber = body.Header.SequenceNumber
	answer.EntryBlockList = make([]ListEntry, len(body.EntryBlockList))
	for i, v := range body.EntryBlockList {
		answer.EntryBlockList[i].ChainID = v.ChainID
		answer.EntryBlockList[i].KeyMR = v.KeyMR
	}
	//answer.DBlock = *body
	answer.BlockTimeStr = TimestampToString(body.Header.Timestamp)
	answer.KeyMR = keyMR

	return answer, nil
}

func Init() {
	/*for _, v := range BucketList {
		err := factoidState.GetDB().Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(v))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
	}*/
}

type ByteData []byte

var _ factoid.IBlock = (*ByteData)(nil)

/*
	encoding.BinaryMarshaler   // Easy to support this, just drop the slice.
	encoding.BinaryUnmarshaler // And once in Binary, it must come back.
	//encoding.TextMarshaler     // Using this mostly for debugging
	CustomMarshalText() ([]byte, error)

	// We need the progress through the slice, so we really can't use the stock spec
	// for the UnmarshalBinary() method from encode.  We define our own method that
	// makes the code easier to read and way more efficent.
	UnmarshalBinaryData(data []byte) ([]byte, error)
	String() string // Makes debugging, logging, and error reporting easier

	IsEqual(IBlock) []IBlock // Check if this block is the same as itself.
	//   Returns nil, or the path to the first difference.

	GetDBHash() IHash       // Identifies the class of the object
	GetHash() IHash         // Returns the hash of the object
	GetNewInstance() IBlock // Get a new instance of this object*/

func (bd ByteData) MarshalBinary() (data []byte, err error) {
	return []byte(bd), nil
}

func (bd ByteData) UnmarshalBinary(data []byte) error {
	bd = data
	return nil
}

func (bd ByteData) CustomMarshalText() ([]byte, error) {
	return []byte(fmt.Sprint("%x", bd)), nil
}

func (bd ByteData) UnmarshalBinaryData(data []byte) ([]byte, error) {
	bd = data
	return nil, nil
}

func (bd ByteData) String() string {
	return fmt.Sprint("%x", bd)
}

func (bd ByteData) IsEqual(factoid.IBlock) []factoid.IBlock {
	return nil
}

func (bd ByteData) GetDBHash() factoid.IHash {
	return factoid.Sha([]byte("ByteData"))
}

func (bd ByteData) GetHash() factoid.IHash {
	return factoid.Sha([]byte(bd))
}

func (bd ByteData) GetNewInstance() factoid.IBlock {
	return new(ByteData)
}

func LoadData(bucket, key string, dst interface{}) (interface{}, error) {
	fmt.Printf("\nLoadData - %v, %v\n\n", bucket, key)
	v := factoidState.GetDB().GetRaw([]byte(bucket), []byte(key))

	if v == nil {
		return nil, nil
	}

	bd := v.(*ByteData)

	dec := gob.NewDecoder(bytes.NewBuffer((*bd)[:]))
	err := dec.Decode(dst)
	if err != nil {
		log.Printf("Error decoding %v of %v", bucket, key)
		return nil, err
	}

	return dst, nil
}

func SaveData(bucket, key string, toStore interface{}) error {
	var data bytes.Buffer

	enc := gob.NewEncoder(&data)

	err := enc.Encode(toStore)
	if err != nil {
		return err
	}

	factoidState.GetDB().PutRaw([]byte(bucket), []byte(key), ByteData(data.Bytes()))

	return nil
}
