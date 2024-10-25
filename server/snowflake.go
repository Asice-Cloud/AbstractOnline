package server

import (
	"errors"
	"sync"
	"time"
)

//snowflake algorithm
/*
第一个数字是符号位，始终为0，表示正数。
接下来的41位是毫秒时间戳，表示当前时间和固定时间点（称为历元）之间的差。这个41位的时间戳可以使用大约69年。
接下来的10位是工作机ID，可以分为两部分：5位数据中心ID和5位工人ID。这10位工作机ID最多可以支持1024个节点。
最后12位是序列号，表示在同一毫秒内生成的不同ID的序列号。12位序列号支持每个节点每毫秒生成4096个ID序列号（相同的机器，相同的时间戳）。
*/

const (
	epoch             int64 = 1526285084376
	workerIDBits      uint8 = 5
	datacenterIDBits  uint8 = 5
	maxWorkerID       int64 = -1 ^ (-1 << workerIDBits)
	maxDatacenterID   int64 = -1 ^ (-1 << datacenterIDBits)
	sequenceBits      uint8 = 12
	workerIDShift     uint8 = sequenceBits
	datacenterIDShift uint8 = sequenceBits + workerIDBits
	timestampShift    uint8 = sequenceBits + workerIDBits + datacenterIDBits
	sequenceMask      int64 = -1 ^ (-1 << sequenceBits)
)

type Snowflake struct {
	mu            sync.Mutex
	lastTimestamp int64
	workerID      int64
	datacenterID  int64
	sequence      int64
}

// NewSnowflake Work node ID, used to identify different work nodes in the same data center. In a distributed system,
// different workerIDs can be assigned according to the actual situation to distinguish different servers or processes
// Data center ID, used to identify different data centers. In distributed systems,
// different datacenterIDs can be assigned according to actual situations to distinguish different data centers.
func NewSnowflake(workerID, datacenterID int64) (*Snowflake, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, errors.New("invalid worker ID")
	}
	if datacenterID < 0 || datacenterID > maxDatacenterID {
		return nil, errors.New("invalid datacenter ID")
	}
	return &Snowflake{
		lastTimestamp: 0,
		workerID:      workerID,
		datacenterID:  datacenterID,
		sequence:      0,
	}, nil
}

func (s *Snowflake) NextID() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ts := time.Now().UnixNano() / 1e6
	if ts < s.lastTimestamp {
		return 0, errors.New("invalid timestamp")
	}
	if ts == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			for ts <= s.lastTimestamp {
				ts = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		s.sequence = 0
	}
	s.lastTimestamp = ts
	return ((ts - epoch) << timestampShift) | (s.datacenterID << datacenterIDShift) | (s.workerID << workerIDShift) | s.sequence, nil
}
