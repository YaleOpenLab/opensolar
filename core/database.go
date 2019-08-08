package core

import (
	"log"

	edb "github.com/Varunram/essentials/database"
	"github.com/boltdb/bolt"

	consts "github.com/YaleOpenLab/opensolar/consts"
)

var InvestorBucket = []byte("Investors")
var RecipientBucket = []byte("Recipients")
var ProjectsBucket = []byte("Projects")
var ContractorBucket = []byte("Contractors")

// CreateHomeDir creates a home directory
func CreateHomeDir() {
	edb.CreateDirs(consts.HomeDir, consts.DbDir, consts.OpenSolarIssuerDir)
	log.Println("creating db at: ", consts.DbDir+consts.DbName)
	db, err := edb.CreateDB(consts.DbDir+consts.DbName, ProjectsBucket, InvestorBucket, RecipientBucket, ContractorBucket)
	if err != nil {
		log.Fatal(err)
	}
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
