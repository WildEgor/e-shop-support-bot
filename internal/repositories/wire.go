package repositories

import (
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/db"
	"github.com/google/wire"
)

var RepositoriesSet = wire.NewSet(
	db.DbSet,
	wire.Bind(new(IUserStateRepository), new(*UserStateRepository)),
	NewUserStateRepository,
	wire.Bind(new(ITopicsRepository), new(*TopicRepository)),
	NewTopicsRepository,
	wire.Bind(new(IGroupRepository), new(*GroupRepository)),
	NewGroupRepository,
)
