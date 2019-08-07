package database

import (
	edb "github.com/Varunram/essentials/database"
	consts "github.com/YaleOpenLab/opensolar/consts"
	"github.com/boltdb/bolt"
)

var InvestorBucket = []byte("Investors")
var RecipientBucket = []byte("Recipients")
var ProjectsBucket = []byte("Projects")
var ContractorBucket = []byte("Contractors")

// CreateHomeDir creates a home directory
func CreateHomeDir() {
	edb.CreateDirs(consts.HomeDir, consts.DbDir, consts.OpenSolarIssuerDir)
	db, _ := edb.CreateDB(consts.DbDir+consts.DbName, ProjectsBucket, InvestorBucket, RecipientBucket, ContractorBucket)
	db.Close()
}

// OpenDB opens the db
func OpenDB() (*bolt.DB, error) {
	return edb.OpenDB(consts.DbDir + consts.DbName)
}

// DeleteKeyFromBucket deletes a given key from the bucket
func DeleteKeyFromBucket(key int, bucketName []byte) error {
	return edb.DeleteKeyFromBucket(consts.DbDir+consts.DbName, key, bucketName)
}
