package core

import (
	"log"

	edb "github.com/Varunram/essentials/database"
	"github.com/boltdb/bolt"

	consts "github.com/YaleOpenLab/opensolar/consts"
)

// InvestorBucket is the investor bucket
var InvestorBucket = []byte("Investors")

// RecipientBucket is the recipient bucket
var RecipientBucket = []byte("Recipients")

// ProjectsBucket is the project bucket
var ProjectsBucket = []byte("Projects")

// ContractorBucket is the contractor bucket
var ContractorBucket = []byte("Contractors")

// CreateHomeDir creates a home directory at $HOME. If the user does not have permissions
// to write to home, execution is stopped.
func CreateHomeDir() {
	edb.CreateDirs(consts.HomeDir, consts.DbDir, consts.OpenSolarIssuerDir)
	log.Println("creating db at: ", consts.DbDir+consts.DbName)
	db, err := edb.CreateDB(consts.DbDir+consts.DbName, ProjectsBucket, InvestorBucket, RecipientBucket, ContractorBucket)
	if err != nil {
		log.Fatal(err)
	}
	db.Close()
}

// OpenDB opens the database at the home directory.
func OpenDB() (*bolt.DB, error) {
	return edb.OpenDB(consts.DbDir + consts.DbName)
}

// DeleteKeyFromBucket deletes a given key from the bucket
func DeleteKeyFromBucket(key int, bucketName []byte) error {
	return edb.DeleteKeyFromBucket(consts.DbDir+consts.DbName, key, bucketName)
}
