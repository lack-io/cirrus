package cdiscount

import (
	"errors"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/xingyys/cirrus/config"
)

var (
	ErrInvalidConfig = errors.New("invalid config")
	ErrConnectDB     = errors.New("connect database")
	ErrDBWrite       = errors.New("write to database")
	ErrDBRead        = errors.New("read from database")
)

// 商品信息
type Good struct {
	ID uint64 `gorm:"primaryKey"`

	// UID 唯一ID
	UID string `json:"uid,omitempty"`

	// URL 所在网址
	URL string `json:"url,omitempty"`

	// Comments 评论数
	Comments int `json:"comments,omitempty"`

	// Express 快递信息
	Express string `json:"express,omitempty"`

	// 入库时间
	Timeout int64 `json:"timeout,omitempty"`
}

type Store struct {
	cfg *config.Store

	db *gorm.DB
}

func NewStore(cfg *config.Store) (*Store, error) {
	s := &Store{cfg: cfg}

	switch cfg.DB {
	case config.Sqlite:
		if cfg.Sqlite == nil || cfg.Sqlite.Name == "" {
			return nil, fmt.Errorf("%w: missing cfg sqlite", ErrInvalidConfig)
		}
		var err error
		s.db, err = gorm.Open(sqlite.Open(cfg.Sqlite.Name), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrConnectDB, err)
		}
	}

	err := s.db.AutoMigrate(&Good{})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDBWrite, err)
	}

	return s, nil
}

func (s *Store) GetGoods(page, size int) ([]*Good, error) {
	goods := make([]*Good, 0)

	err := s.db.Limit(size).Offset((page - 1) * size).Find(&goods).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDBRead, err)
	}

	return goods, nil
}

func (s *Store) AddGood(good *Good) error {

	err := s.db.Create(good).Error
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDBWrite, err)
	}

	return nil
}

func (s *Store) Reset() {

}
