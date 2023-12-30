package restaurant

import (
	abstract "github.com/Dparty/common/abstract"
	"github.com/Dparty/common/snowflake"
	"github.com/Dparty/dao"
	"gorm.io/gorm"
)

var printerPrinterRepository *PrinterRepository

// GetPrinterRepository returns the printer repository by Lazy bones
func GetPrinterRepository() *PrinterRepository {
	if printerPrinterRepository == nil {
		printerPrinterRepository = NewPrinterRepository(dao.GetDBInstance())
	}
	return printerPrinterRepository
}

type PrinterType string

const (
	BILL    PrinterType = "BILL"
	KITCHEN PrinterType = "KITCHEN"
)

type Printer struct {
	gorm.Model
	RestaurantId uint
	Name         string      `json:"name"`
	Sn           string      `json:"sn"`
	Description  string      `json:"description"`
	Type         PrinterType `json:"type" gorm:"type:VARCHAR(128)"`
	PrinterModel string      `json:"printerModel" gorm:"type:VARCHAR(128);default:58mm"`
}

var printerIdGenerator = snowflake.NewIdGenertor(1)

func (a *Printer) BeforeCreate(tx *gorm.DB) (err error) {
	a.Model.ID = printerIdGenerator.Uint()
	return err
}

func (p Printer) ID() uint {
	return p.Model.ID
}

func (p *Printer) SetOwner(owner abstract.Owner) *Printer {
	p.Model.ID = owner.ID()
	return p
}

func NewPrinterRepository(db *gorm.DB) *PrinterRepository {
	return &PrinterRepository{
		db: db,
	}
}

type PrinterRepository struct {
	db *gorm.DB
}

func (p PrinterRepository) Save(printer *Printer) *Printer {
	p.db.Save(printer)
	return printer
}

func (p PrinterRepository) Find(conds ...any) *Printer {
	var printer Printer
	ctx := p.db.Find(&printer, conds...)
	if ctx.RowsAffected == 0 {
		return nil
	}
	return &printer
}

func (p PrinterRepository) GetById(id uint) *Printer {
	return p.Find(id)
}

func (p PrinterRepository) List(conds ...any) []Printer {
	var printers []Printer
	p.db.Find(&printers, conds...)
	return printers
}

func (p PrinterRepository) Delete(id uint) *gorm.DB {
	return p.db.Delete(&Printer{}, id)
}
