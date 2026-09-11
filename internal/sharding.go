package internal
import (
	"crypto/fnv"
	"errors"
)
type ShardManager struct {
	shards []string
}
func NewShardManager(shards []string) (*ShardManager, error) {
	if len(shards) == 0 {
		return nil, errors.New("shard list cannot be empty")
	}
	return &ShardManager{shards: shards},nil
}
func (sm *ShardManager) GetShard(key string) string {
	hasher := fnv.New32a()
	hasher.Write([]byte(key))
	hashValue := hasher.Sum32()
	index := hashValue % uint32(len(sm.shards))
	return sm.shards[index]
}