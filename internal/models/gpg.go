package models

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/hibare/DomainHQ/internal/constants"
	"gorm.io/gorm"
)

type GPGUsers struct {
	ID      string `gorm:"primaryKey;autoIncrement" json:"-"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Comment string `json:"comment"`
}

func (GPGUsers) TableName() string {
	return "gpg_users"
}

type GPGPubKeyStore struct {
	KeyID       string     `gorm:"primaryKey;" json:"key_id"`
	KeyIDShort  string     `json:"key_id_short"`
	Fingerprint string     `json:"fingerprint"`
	CreatedAt   time.Time  `json:"created_at"`
	Algorithm   string     `json:"algorithm"`
	Version     int        `json:"version"`
	Revoked     bool       `json:"revoked"`
	Users       []GPGUsers `gorm:"foreignKey:ID;constraint:OnDelete:CASCADE"`
	PublicKey   string     `json:"public_key"`
}

func (GPGPubKeyStore) TableName() string {
	return "gpg_pub_key_stores"
}

func ParsePubKey(keyText string) (GPGPubKeyStore, error) {
	entities, err := openpgp.ReadArmoredKeyRing(bytes.NewBufferString(keyText))
	if err != nil {
		return GPGPubKeyStore{}, err
	}

	// Throw error when there are more than one key
	if len(entities) > 1 {
		return GPGPubKeyStore{}, fmt.Errorf("more than one key found")
	}

	entity := entities[0]
	key := GPGPubKeyStore{
		KeyID:       strings.ToLower(entity.PrimaryKey.KeyIdString()),
		KeyIDShort:  strings.ToLower(entity.PrimaryKey.KeyIdShortString()),
		Fingerprint: strings.ToLower(hex.EncodeToString(entity.PrimaryKey.Fingerprint[:])),
		CreatedAt:   entity.PrimaryKey.CreationTime,
		Algorithm:   string(entity.PrimaryKey.PubKeyAlgo),
		Version:     entity.PrimaryKey.Version,
		Revoked:     entity.Revoked(time.Now()),
		PublicKey:   keyText,
	}

	for _, id := range entity.Identities {
		key.Users = append(key.Users, GPGUsers{
			Name:    id.UserId.Name,
			Email:   strings.ToLower(id.UserId.Email),
			Comment: id.UserId.Comment,
		})
	}

	return key, nil
}

func AddPubKey(db *gorm.DB, key *GPGPubKeyStore) error {
	err := db.First(key).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return db.Create(key).Error
		} else {
			return err
		}
	} else {
		return db.Save(key).Error
	}
}

func LookupPubKey(db *gorm.DB, searchStr string) (*GPGPubKeyStore, error) {
	keys := []GPGPubKeyStore{}
	searchStr = strings.ToLower(searchStr)

	err := db.Preload("Users").Joins("JOIN gpg_users ON gpg_users.id = gpg_pub_key_stores.key_id").
		Where("gpg_pub_key_stores.key_id = ? OR gpg_pub_key_stores.key_id_short = ? OR gpg_pub_key_stores.fingerprint = ? OR gpg_users.email = ?", searchStr, searchStr, strings.TrimPrefix(searchStr, constants.GPGFingerprintPrefix), searchStr).
		Find(&keys).Error
	if err != nil {
		return nil, err
	}

	if len(keys) > 1 {
		return nil, fmt.Errorf("more than one key found")
	}

	if len(keys) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &keys[0], nil
}
