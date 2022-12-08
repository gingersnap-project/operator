package images

import "os"

const (
	CacheManagerEnvName = "RELATED_IMAGE_CACHE_MANAGER"
	DBSyncerEnvName     = "RELATED_IMAGE_DB_SYNCER"
)

var (
	CacheManager = os.Getenv(CacheManagerEnvName)
	DBSyncer     = os.Getenv(DBSyncerEnvName)
)
