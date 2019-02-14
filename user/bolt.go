package user

import (
	"path"
	"time"

	"github.com/arknable/fwdproxy/env"
	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

const (
	// BucketName is the name of user bucket
	BucketName = "Users"
)

// BoltRepository implements Repository using boltdb data source
type BoltRepository struct {
	db *bolt.DB
}

// Initialize implements Repository.Initialize
func (r *BoltRepository) Initialize() error {
	dbPath := path.Join(env.HomePath(), "users.db")
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		return err
	}
	r.db = db
	log.WithField("path", dbPath).Info("Using path as DB path")

	tx, err := r.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bucket, err := tx.CreateBucketIfNotExists([]byte(BucketName))
	if err != nil {
		return err
	}
	val := bucket.Get([]byte(BuiltInUsername))
	if val == nil {
		err = bucket.Put([]byte(BuiltInUsername), []byte(BuiltInUserPwd))
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Close implements Repository.Close
func (r *BoltRepository) Close() error {
	return r.db.Close()
}

// Validate implements Repository.Validate
func (r *BoltRepository) Validate(username, password string) (bool, error) {
	valid := false
	err := r.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		pwd := bucket.Get([]byte(username))
		if pwd != nil {
			valid = (password == string(pwd))
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return valid, nil
}
