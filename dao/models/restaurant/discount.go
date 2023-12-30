package restaurant

import (
	"github.com/Dparty/common/singleton"
	"github.com/Dparty/common/snowflake"
	"github.com/Dparty/dao"
	"gorm.io/gorm"
)

type Discount struct {
	gorm.Model
	RestaurantId uint
	Label        string
	Offset       int64
}

var discountIdGenerator = snowflake.NewIdGenertor(1)

func (a *Discount) BeforeCreate(tx *gorm.DB) (err error) {
	a.Model.ID = discountIdGenerator.Uint()
	return err
}

type DiscountRepository struct {
	db *gorm.DB
}

var discountRepository = singleton.NewSingleton[DiscountRepository](newDiscountRepository, singleton.Eager)

func GetDiscountRepository() *DiscountRepository {
	return discountRepository.Get()
}

func newDiscountRepository() *DiscountRepository {
	return &DiscountRepository{
		db: dao.GetDBInstance(),
	}
}

func (r DiscountRepository) Find(conds ...any) *Discount {
	var discount Discount
	ctx := r.db.Find(&discount, conds...)
	if ctx.RowsAffected == 0 {
		return nil
	}
	return &discount
}

func (r DiscountRepository) FindById(id uint) *Discount {
	return r.Find(id)
}

func (r DiscountRepository) List(conds ...any) []Discount {
	var discounts []Discount
	r.db.Find(&discounts, conds...)
	return discounts
}

func (r DiscountRepository) ListBy(restaurantId uint) []Discount {
	var discounts []Discount
	r.db.Where("restaurant_id = ?", restaurantId).Find(&discounts)
	return discounts
}

func (r DiscountRepository) Save(discount *Discount) {
	r.db.Save(discount)
}

func (r DiscountRepository) Delete(discount *Discount) {
	r.db.Delete(discount)
}
