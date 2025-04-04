package repository

import (
	"book-management/internal/repository/cache"
	"book-management/internal/repository/dao"
	"book-management/internal/repository/repo"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(cache.ProviderSet, dao.ProviderSet, repo.ProviderSet)
