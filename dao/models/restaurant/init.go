package restaurant

import (
	"github.com/Dparty/dao"
)

var db = dao.GetDBInstance()

func init() {
	db.AutoMigrate(&Restaurant{}, &Table{}, &Printer{}, &Item{}, &Bill{}, &Discount{})
}
