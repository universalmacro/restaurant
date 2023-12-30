package restaurant

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	abstract "github.com/Dparty/common/abstract"
	"github.com/Dparty/common/data"
	"github.com/Dparty/common/fault"
	"github.com/Dparty/common/snowflake"
	"github.com/Dparty/dao"
	"github.com/Dparty/dao/common"
	"gorm.io/gorm"
)

var itemRepository *ItemRepository

// GetItemRepository returns the item repository by Lazy bones
func GetItemRepository() *ItemRepository {
	if itemRepository == nil {
		itemRepository = NewItemRepository()
	}
	return itemRepository
}

type Item struct {
	gorm.Model
	RestaurantId uint              `json:"restaurantId" gorm:"index:restaurant_index"`
	Name         string            `json:"name"`
	Pricing      int64             `json:"pricing"`
	Attributes   Attributes        `json:"attributes"`
	Images       common.StringList `json:"images" gorm:"type:JSON"`
	Tags         common.StringList `json:"tags"`
	Printers     common.IDList     `json:"printers"`
	Status       string            `json:"status" gorm:"type:VARCHAR(32);default:ACTIVED"`
	Alcohol      bool              `json:"alcohol"`
}

var itemIdGenerator = snowflake.NewIdGenertor(1)

func (a *Item) BeforeCreate(tx *gorm.DB) (err error) {
	a.Model.ID = itemIdGenerator.Uint()
	return err
}

func (i Item) ID() uint {
	return i.Model.ID
}

func (i *Item) SetOwner(owner abstract.Owner) *Item {
	i.RestaurantId = owner.ID()
	return i
}

// func (i Item) Owner() *Restaurant {
// 	return restaurantRepository.GetById(i.RestaurantId)
// }

func (i Item) CreateOrder(specification []data.Pair[string, string]) (Order, error) {
	// TODO: specification verification
	return Order{
		Item:          i,
		Specification: specification,
	}, nil
}

type Attributes []Attribute

func (as Attributes) GetOption(left, right string) (data.Pair[string, string], error) {
	for _, a := range as {
		if left == a.Label {
			for _, option := range a.Options {
				if right == option.Label {
					return data.Pair[string, string]{L: left, R: right}, nil
				}
			}
		}
	}
	return data.Pair[string, string]{}, errors.New("not found")
}

func (Attributes) GormDataType() string {
	return "JSON"
}

func (s *Attributes) Scan(value any) error {
	return json.Unmarshal(value.([]byte), s)
}

func (s Attributes) Value() (driver.Value, error) {
	b, err := json.Marshal(s)
	return b, err
}

type Attribute struct {
	Label   string   `json:"label"`
	Options []Option `json:"options"`
}

type Option struct {
	Label string `json:"label"`
	Extra int64  `json:"extra"`
}

type Options []Option

func (Options) GormDataType() string {
	return "JSON"
}

func (s *Options) Scan(value any) error {
	return json.Unmarshal(value.([]byte), s)
}

func (s Options) Value() (driver.Value, error) {
	b, err := json.Marshal(s)
	return b, err
}

func (Attribute) GormDataType() string {
	return "JSON"
}

func (s *Attribute) Scan(value any) error {
	return json.Unmarshal(value.([]byte), s)
}

func (s Attribute) Value() (driver.Value, error) {
	b, err := json.Marshal(s)
	return b, err
}

type ItemRepository struct {
	db *gorm.DB
}

func NewItemRepository() *ItemRepository {
	return &ItemRepository{
		db: dao.GetDBInstance(),
	}
}

func (i ItemRepository) Get(conds ...any) *Item {
	var item Item
	ctx := i.db.Find(&item, conds...)
	if ctx.RowsAffected == 0 {
		return nil
	}
	return &item
}

func (i ItemRepository) GetById(id uint) *Item {
	return i.Get(id)
}

func (i ItemRepository) Save(item *Item) (*Item, error) {
	var attributesMap map[string]bool = make(map[string]bool)
	for _, attribute := range item.Attributes {
		_, ok := attributesMap[attribute.Label]
		if ok {
			return nil, fault.ErrItemAttributesConflict
		}
		attributesMap[attribute.Label] = true
		var optionMap map[string]bool = make(map[string]bool)
		for _, option := range attribute.Options {
			_, ok := optionMap[option.Label]
			if ok {
				return nil, fault.ErrItemAttributesConflict
			}
			optionMap[option.Label] = true
		}
	}
	i.db.Save(item)
	return item, nil
}

func (i ItemRepository) List(conds ...any) []Item {
	var items []Item
	i.db.Find(&items, conds...)
	return items
}

func (i ItemRepository) Delete(item *Item) *gorm.DB {
	return i.db.Delete(item)
}
