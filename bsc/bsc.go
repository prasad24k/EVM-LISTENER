package bsc

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	geth "github.com/ethereum/go-ethereum/ethclient"
)

type Transactions struct {
	Amount         float64 `json:"amount"`
	TxHash         string  `json:"txn_hash"`
	BlockNumber    string  `json:"block_number"`
	DepositAddress string  `json:"deposit_address"`
	FromAddress    string  `json:"from_address"`
	Status         string  `json:"status"`
	AssetId        int     `json:"asset_id"`
}
type EthListener struct {
	Client     *geth.Client
	HeaderChan chan *types.Header
	Contracts  []Contract

	Mutex            sync.Mutex
	WatchedAddresses map[string]bool
}

type Contract struct {
	Address common.Address
	AssetID int
}

func NewEthClient(rpc string, contracts ...Contract) *EthListener {
	client, err := geth.DialContext(context.Background(), rpc)
	if err != nil {
		log.Fatalf("Error connecting to Eth Client: %s", err.Error())
		return nil
	}

	return &EthListener{
		Client:           client,
		HeaderChan:       make(chan *types.Header),
		Contracts:        contracts,
		WatchedAddresses: make(map[string]bool),
	}
}

func (l *EthListener) Start() error {
	sub, err := l.Client.SubscribeNewHead(context.Background(), l.HeaderChan)
	if err != nil {
		return err
	}

	go l.processHeaders(sub)
	return nil
}

func (l *EthListener) processHeaders(sub ethereum.Subscription) {
	var wg sync.WaitGroup
	for {
		select {
		case err := <-sub.Err():
			log.Printf("Subscription error: %v", err)
			return
		case header := <-l.HeaderChan:
			wg.Add(1)
			go func(h *types.Header) {
				defer wg.Done()
				if err := l.processBlock(h.Number); err != nil {
					//log.Printf("Error processing block %v: %v", h.Number, err)
				}
			}(header)
		}
	}
}

func (l *EthListener) processBlock(blockNumber *big.Int) error {

	log.Println(blockNumber)
	block, err := l.Client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		return err
	}

	for _, tx := range block.Transactions() {
		if tx.To() == nil {
			continue
		}

		//log.Println("hi ,", tx.To())
		isContract, _ := l.isContractAddress(tx.To())
		if isContract {
			log.Println("entered")
			if err := l.processContractTransaction(tx); err != nil {
				log.Printf("Error processing contract transaction %s: %v", tx.Hash().Hex(), err)
			}
		} else {
			if err := l.processEthTransaction(tx); err != nil {
				//log.Printf("Error processing ETH transaction %s: %v", tx.Hash().Hex(), err)
			}
		}
	}
	return nil
}

func (l *EthListener) processEthTransaction(tx *types.Transaction) error {
	blockNumber, err := l.Client.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	receipt, err := l.Client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		return err
	}

	status := "failed"
	if receipt.Status == types.ReceiptStatusSuccessful {
		status = "success"
	}

	from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		return err
	}

	toAddress := tx.To().Hex()

	var details Transactions

	if l.IsWatchedAddress(toAddress) {
		amountFloat := new(big.Float).SetInt(tx.Value())

		// Divide by 10^18 if assetID is 3

		// log.Println(amountFloat.Float64(), " < 111 111 111")
		amount, _ := amountFloat.Float64()
		log.Println(amount, " < 0 0 00 ")

		details.Amount = amount / (10.0 * 10.0 * 10.0 * 10.0 * 10.0 * 10.0 * 10.0 * 10.0 * 10.0 * 10.0 * 10.0 * 10.0 * 10.0 * 10.0 * 10.0 * 10.0 * 10.0 * 10.0)

		// details.Amount = amount
		details.TxHash = tx.Hash().Hex()
		details.AssetId = 3
		details.BlockNumber = strconv.FormatUint(blockNumber, 10)
		details.DepositAddress = toAddress
		details.FromAddress = from.Hex()
		details.Status = status

		log.Println("bsc hash: ", details.Amount)
		log.Println("bsc hash: ", details.TxHash)
		log.Println("bsc deposit address: ", details.DepositAddress)
		log.Println("bsc from address: ", details.FromAddress)
		log.Println("bsc deposit status: ", details.Status)
		log.Println("asset id : ", details.AssetId)

		// if err := l.callInsertDepositDetailsAPI(details); err != nil {
		// 	log.Println("Calling callInsertDepositDetailsAPI...")
		// 	log.Printf("Error calling API to insert deposit details: %v", err)
		// }
	}

	return nil
}
func (l *EthListener) processContractTransaction(tx *types.Transaction) error {
	blockNumber, err := l.Client.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	receipt, err := l.Client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		return err
	}

	status := "failed"
	if receipt.Status == types.ReceiptStatusSuccessful {
		status = "success"
	}

	from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		return err
	}

	// Debug: Log transaction basics
	//log.Printf("Processing contract transaction: %s to contract: %s",tx.Hash().Hex(), tx.To().Hex())

	receiver, amountStr := decodeUSDTTransfer(tx.Data())

	// Debug: Log decoded data
	//log.Printf("Decoded transfer - recipient: %s, amount: %s", receiver, amountStr)

	isContract, assetID := l.isContractAddress(tx.To())

	// Debug: Log contract check results
	//log.Printf("Contract check - isContract: %v, assetID: %d", isContract, assetID)

	// Debug: Log watched address check
	_ = l.IsWatchedAddress(receiver)
	//isWatched := l.IsWatchedAddress(receiver)
	//log.Printf("Watched address check - address: %s, isWatched: %v", receiver, isWatched)

	if isContract && l.IsWatchedAddress(receiver) {
		//log.Printf("✅ Transaction eligible for processing - assetID: %d, recipient: %s", assetID, receiver)

		//log.Println(amountStr, " < - --  before")

		amountFloat, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			log.Printf("Error converting amount to float64: %v", err)
			return err
		}

		log.Println(assetID, amountFloat, " < - - - -")

		var details Transactions
		details.TxHash = tx.Hash().Hex()
		details.BlockNumber = strconv.FormatUint(blockNumber, 10)
		details.DepositAddress = receiver
		details.FromAddress = from.Hex()
		details.Status = status
		details.AssetId = assetID

		// Set amount based on assetID
		switch assetID {
		case 1, 2, 3:
			details.Amount = amountFloat / 1e18

		default:
			details.Amount = amountFloat
		}
		// Debug: Log transaction details before API call
		log.Printf("Calling API with details - AssetID: %d, Amount: %f, TxHash: %s, Recipient: %s",
			details.AssetId, details.Amount, details.TxHash, details.DepositAddress)

		// if err := l.callInsertDepositDetailsAPI(details); err != nil {
		// 	//log.Println("Calling callInsertDepositDetailsAPI...")
		// 	log.Printf("Error calling API to insert deposit details: %v", err)
		// } else {
		// 	log.Printf("Successfully called API for assetID: %d", assetID)
		// }
	} else {
		log.Printf("❌ Skipping transaction - isContract: %v, isWatched: %v, assetID: %d, receiver: %s",
			isContract, l.IsWatchedAddress(receiver), assetID, receiver)
	}

	return nil
}

//	func (l *EthListener) isContractAddress(addr *common.Address) (bool, int) {
//		for _, contract := range l.Contracts {
//			if contract.Address == *addr {
//				return true, contract.AssetID
//			}
//		}
//		return false, 0
//	}
func (l *EthListener) isContractAddress(addr *common.Address) (bool, int) {
	addrHex := addr.Hex()
	//log.Printf("Checking if %s is a contract", addrHex)

	for _, contract := range l.Contracts {
		_ = contract.Address.Hex()
		if contract.Address == *addr {
			//log.Printf("✅ Match found: %s is contract with assetID: %d", addrHex, contract.AssetID)
			return true, contract.AssetID
		}
		//log.Printf("Compared with contract: %s (assetID: %d) - no match", contractHex, contract.AssetID)
	}
	log.Printf("❌ No matching contract found for address: %s", addrHex)
	return false, 0
}

func (l *EthListener) IsWatchedAddress(address string) bool {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	_, exists := l.WatchedAddresses[address]
	//log.Printf("Checking if %s is watched: %v", address, exists)

	// Debug: If we're expecting this address to be watched but it's not, dump all watched addresses
	if !exists {
		log.Printf("Address %s not found in watched addresses. Current watched addresses:", address)
		for watchedAddr := range l.WatchedAddresses {
			log.Printf("- %s", watchedAddr)
		}
	}

	return exists
}
func (l *EthListener) AddWatchedAddress(address string) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()
	l.WatchedAddresses[address] = true
}
func decodeUSDTTransfer(data []byte) (recipient string, amount string) {
	log.Printf("Decoding transaction data of length: %d", len(data))

	// Debug: Print first few bytes of data
	if len(data) >= 8 {
		//log.Printf("Method signature (first 4 bytes): %x", data[:4])
	}

	if len(data) != 68 {
		log.Printf("❌ Transaction data length is not 68 bytes, got %d bytes", len(data))
		return "", ""
	}

	recipientBytes := data[16:36]
	recipient = common.BytesToAddress(recipientBytes).Hex()

	amountBytes := data[36:68]
	amountBigInt := new(big.Int).SetBytes(amountBytes)

	//log.Printf("✅ Successfully decoded transfer: recipient=%s, amount=%v", recipient, amountBigInt)
	return recipient, fmt.Sprintf("%v", amountBigInt)
}
