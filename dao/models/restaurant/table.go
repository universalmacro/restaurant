package restaurant

import (
	abstract "github.com/Dparty/common/abstract"
	"github.com/Dparty/common/snowflake"
	"github.com/Dparty/dao"
	"gorm.io/gorm"
)

var tableRepository *TableRepository

func GetTableRepository() *TableRepository {
	if tableRepository == nil {
		tableRepository = NewTableRepository()
	}
	return tableRepository
}

var tableIdGenerator = snowflake.NewIdGenertor(10)

type Table struct {
	gorm.Model
	RestaurantId uint
	Label        string `json:"label"`
	X            int64  `json:"x"`
	Y            int64  `json:"y"`
}

func (a *Table) BeforeCreate(tx *gorm.DB) (err error) {
	a.Model.ID = tableIdGenerator.Uint()
	return err
}

func (t Table) ID() uint {
	return t.Model.ID
}

func (t *Table) SetOwner(owner abstract.Owner) *Table {
	t.Model.ID = owner.ID()
	return t
}

func (t Table) Bills(status *string) []Bill {
	var bills []Bill
	ctx := db.Where("table_id = ?", t.ID())
	if status != nil {
		ctx = db.Where("status = ?", *status)
	}
	ctx.Find(&bills)
	return bills
}

type TableRepository struct {
	db             *gorm.DB
	billRepository *BillRepository
}

func NewTableRepository() *TableRepository {
	return &TableRepository{
		db:             dao.GetDBInstance(),
		billRepository: GetBillRepository(),
	}
}

func (t TableRepository) Find(conds ...any) *Table {
	var table Table
	ctx := t.db.Find(&table, conds)
	if ctx.RowsAffected == 0 {
		return nil
	}
	return &table
}

func (t TableRepository) List(conds ...any) []Table {
	var tables []Table
	t.db.Find(&tables, conds...)
	return tables
}

func (t TableRepository) Save(table *Table) *gorm.DB {
	return t.db.Save(table)
}

func (t TableRepository) Delete(table *Table) *gorm.DB {
	return t.db.Delete(&table)
}
