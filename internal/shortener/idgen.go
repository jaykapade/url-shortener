package shortener

import (
	"log"
	"sync"

	"github.com/sony/sonyflake"
)

var (
	sf   *sonyflake.Sonyflake
	once sync.Once
)

func InitIDGenerator() {
	once.Do(func() {
		sf = sonyflake.NewSonyflake(sonyflake.Settings{})
		if sf == nil {
			log.Fatal("Failed to initialize Sonyflake")
		}
	})
}

func GenerateFlakeID() (uint64, error) {
	return sf.NextID()
}
