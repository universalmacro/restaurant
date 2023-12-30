package repositories

import (
	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/restaurant/dao/models/restaurant"
)

var restaurantRepositorySingleton = singleton.NewSingleton[RestaurantRepository](newRestaurantRepository, singleton.Eager)

func GetRestaurantRepository() *RestaurantRepository {
	return restaurantRepositorySingleton.Get()
}

type RestaurantRepository struct {
	*dao.Repository[restaurant.Restaurant]
}

func newRestaurantRepository() *RestaurantRepository {
	return &RestaurantRepository{
		Repository: dao.NewRepository[restaurant.Restaurant](),
	}
}
