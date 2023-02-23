package images

import "os"

const (
	CacheManagerMSSQLEnvName    = "RELATED_IMAGE_CACHE_MANAGER_MSSQL"
	CacheManagerMySQLEnvName    = "RELATED_IMAGE_CACHE_MANAGER_MYSQL"
	CacheManagerPostgresEnvName = "RELATED_IMAGE_CACHE_MANAGER_POSTGRES"
	DBSyncerEnvName             = "RELATED_IMAGE_DB_SYNCER"
	IndexEnvName                = "RELATED_IMAGE_INDEX"
)

var (
	CacheManagerMSSQL    = os.Getenv(CacheManagerMSSQLEnvName)
	CacheManagerMySQL    = os.Getenv(CacheManagerMySQLEnvName)
	CacheManagerPostgres = os.Getenv(CacheManagerPostgresEnvName)
	DBSyncer             = os.Getenv(DBSyncerEnvName)
	Index                = os.Getenv(IndexEnvName)
)
