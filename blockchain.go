package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

//blockchain data
type Block struct {
	data         map[string]interface{}
	hash         string
	previousHash string
	timestamp    time.Time
	pow          int
}

//blockchain type
//The genesisBlock property represents the first block added to the blockchain. In contrast, the difficulty property defines the minimum effort miners have to undertake to mine a block and include it in the blockchain.
type Blockchain struct {
	genesisBlock Block
	chain        []Block
	difficulty   int
}

//method for our Block type that implements this functionality
func (b Block) calculateHash() string {
	data, _ := json.Marshal(b.data)
	blockData := b.previousHash + string(data) + b.timestamp.String() + strconv.Itoa(b.pow)
	blockHash := sha256.Sum256([]byte(blockData))
	return fmt.Sprintf("%x", blockHash)
}

//create a mine() method for our Block type that will keep incrementing the PoW value and calculating the block hash until we get a valid hash.
func (b *Block) mine(difficulty int) {
	for !strings.HasPrefix(b.hash, strings.Repeat("0", difficulty)) {
		b.pow++
		b.hash = b.calculateHash()
	}
}

//function that creates a genesis block for our blockchain and returns a new instance of the Blockchain type
func CreateBlockchain(difficulty int) Blockchain {
	genesisBlock := Block{
		hash:      "0",
		timestamp: time.Now(),
	}
	return Blockchain{
		genesisBlock,
		[]Block{genesisBlock},
		difficulty,
	}
}

//new blocks into a blockchain.
func (b *Blockchain) addBlock(from, to string, amount float64) {
	blockData := map[string]interface{}{
		"from":   from,
		"to":     to,
		"amount": amount,
	}
	lastBlock := b.chain[len(b.chain)-1]
	newBlock := Block{
		data:         blockData,
		previousHash: lastBlock.hash,
		timestamp:    time.Now(),
	}
	newBlock.mine(b.difficulty)
	b.chain = append(b.chain, newBlock)
}

//checking the validity of the blockchain
func (b Blockchain) isValid() bool {
	for i := range b.chain[1:] {
		previousBlock := b.chain[i]
		currentBlock := b.chain[i+1]
		if currentBlock.hash != currentBlock.calculateHash() || currentBlock.previousHash != previousBlock.hash {
			return false
		}
	}
	return true
}

//using the blockchain to make transactions
func main() {
	// create a new blockchain instance with a mining difficulty of 2
	blockchain := CreateBlockchain(2)

	// record transactions on the blockchain for Alice, Bob, and John
	blockchain.addBlock("Alice", "Bob", 5)
	blockchain.addBlock("John", "Bob", 2)
	blockchain.addBlock("Daniel", "Bob", 1112)
	fmt.Println(blockchain)
	// check if the blockchain is valid; expecting true
	fmt.Println(blockchain.isValid())
}
