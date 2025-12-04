package base62snowflake

import (
	"errors"
	"sync"
	"time"
)

// configuration
const (
	epoch     = 1704067200000 // 2024-01-01 00:00:00 UTC
	nodeBits  = 10            // allows 1024 nodes
	stepBits  = 12            // 4096 IDs per millisecond per node
	nodeMax   = -1 ^ (-1 << nodeBits)
	stepMax   = -1 ^ (-1 << stepBits)
	timeShift = nodeBits + stepBits
	nodeShift = stepBits
)

// base62 alphabet
const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// Node holds the state for the snowflake generator
type Node struct {
	mu        sync.Mutex
	timestamp int64
	nodeID    int64
	step      int64
}

var (
	defaultNode *Node
	once        sync.Once
)

// init or lazy-load a default node
func getDefaultNode() *Node {
	once.Do(func() {
		var err error
		defaultNode, err = NewNode(1) // default to node ID 1
		if err != nil {
			panic(err) // shouldn't happen with nodeID=1
		}
	})
	return defaultNode
}

// NewNode creates a generator for a specific machine ID (0-1023)
func NewNode(nodeID int64) (*Node, error) {
	if nodeID < 0 || nodeID > nodeMax {
		return nil, errors.New("node ID out of range")
	}
	return &Node{
		timestamp: 0,
		nodeID:    nodeID,
		step:      0,
	}, nil
}

// Generate creates a base62 encoded snowflake ID
func (n *Node) Generate() string {
	n.mu.Lock()
	defer n.mu.Unlock()

	now := time.Now().UnixMilli()

	if now < n.timestamp {
		now = n.timestamp
	}

	if n.timestamp == now {
		n.step = (n.step + 1) & stepMax
		if n.step == 0 {
			for now <= n.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		n.step = 0
	}

	n.timestamp = now
	id := (now-epoch)<<timeShift | (n.nodeID << nodeShift) | n.step

	return encodeBase62(id)
}

// GetSnowflakeID generates a snowflake ID using the default node
// This is your simple API: GetSnowflakeID()
func GetSnowflakeID() string {
	return getDefaultNode().Generate()
}

// SetDefaultNodeID allows changing the default node ID (must be called before first GetSnowflakeID)
func SetDefaultNodeID(nodeID int64) error {
	if nodeID < 0 || nodeID > nodeMax {
		return errors.New("node ID out of range")
	}
	var err error
	once.Do(func() {
		defaultNode, err = NewNode(nodeID)
	})
	return err
}

// encodeBase62 converts int64 to base62 string
func encodeBase62(id int64) string {
	if id == 0 {
		return "0"
	}

	b := make([]byte, 0, 12)
	val := uint64(id)
	base := uint64(62)

	for val > 0 {
		rem := val % base
		val /= base
		b = append(b, base62Chars[rem])
	}

	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	return string(b)
}
