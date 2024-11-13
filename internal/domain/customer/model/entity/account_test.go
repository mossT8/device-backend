package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"mossT8.github.com/device-backend/internal/application/logger"
)

func TestNewAccount(t *testing.T) {
	timestamp := time.Now()
	account := NewAccount("test@example.com", "Test User", timestamp)

	assert.Equal(t, "test@example.com", account.GetEmail())
	assert.Equal(t, "Test User", account.GetName())
	assert.Equal(t, timestamp.UTC(), account.GetCreatedAt().UTC())
	assert.Equal(t, timestamp.UTC(), account.GetModifiedAt().UTC())
}

func TestAccount_Getters(t *testing.T) {
	timestamp := time.Now()
	account := Account{
		ID:              mysqlRecordId(123),
		Email:           mysqlText("test@example.com"),
		PasswordHash:    mysqlText("hashedpassword"),
		Salt:            mysqlText("salt123"),
		Name:            mysqlText("Test User"),
		Verified:        mysqlBool(true),
		ReceivesUpdates: mysqlBool(false),
		CreatedAt:       mysqlDate(timestamp),
		ModifiedAt:      mysqlDate(timestamp),
	}

	tests := []struct {
		name     string
		getter   interface{}
		expected interface{}
	}{
		{"GetID", account.GetID(), int64(123)},
		{"GetEmail", account.GetEmail(), "test@example.com"},
		{"GetPasswordHash", account.GetPasswordHash(), "hashedpassword"},
		{"GetSalt", account.GetSalt(), "salt123"},
		{"GetName", account.GetName(), "Test User"},
		{"GetVerified", account.GetVerified(), true},
		{"GetReceivesUpdates", account.GetReceivesUpdates(), false},
		{"GetCreatedAt", account.GetCreatedAt().UTC(), timestamp.UTC()},
		{"GetModifiedAt", account.GetModifiedAt().UTC(), timestamp.UTC()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.getter)
		})
	}
}

func TestAccount_Setters(t *testing.T) {
	account := NewAccount("initial@example.com", "Initial Name", time.Now())
	initialModifiedAt := account.GetModifiedAt()
	time.Sleep(time.Millisecond) // Ensure time difference for ModifiedAt checks

	t.Run("SetID", func(t *testing.T) {
		account.SetID(456)
		assert.Equal(t, int64(456), account.GetID())
	})

	t.Run("SetEmail", func(t *testing.T) {
		account.SetEmail("new@example.com")
		assert.Equal(t, "new@example.com", account.GetEmail())
		assert.True(t, account.GetModifiedAt().After(initialModifiedAt))
	})

	t.Run("SetPassword", func(t *testing.T) {
		password := "mypassword"
		salt := "mysalt"
		err := account.SetPassword(password, salt)
		assert.NoError(t, err)
		assert.Equal(t, salt, account.GetSalt())

		// Verify password hash
		err = bcrypt.CompareHashAndPassword([]byte(account.GetPasswordHash()), []byte(password+salt))
		assert.NoError(t, err)
		assert.True(t, account.GetModifiedAt().After(initialModifiedAt))
	})

	t.Run("SetPasswordHash", func(t *testing.T) {
		account.SetPasswordHash("newhash")
		assert.Equal(t, "newhash", account.GetPasswordHash())
		assert.True(t, account.GetModifiedAt().After(initialModifiedAt))
	})

	t.Run("SetSalt", func(t *testing.T) {
		account.SetSalt("newsalt")
		assert.Equal(t, "newsalt", account.GetSalt())
		assert.True(t, account.GetModifiedAt().After(initialModifiedAt))
	})

	t.Run("SetName", func(t *testing.T) {
		account.SetName("New Name")
		assert.Equal(t, "New Name", account.GetName())
		assert.True(t, account.GetModifiedAt().After(initialModifiedAt))
	})

	t.Run("SetVerified", func(t *testing.T) {
		account.SetVerified(true)
		assert.True(t, account.GetVerified())
		assert.True(t, account.GetModifiedAt().After(initialModifiedAt))
	})

	t.Run("SetReceivesUpdates", func(t *testing.T) {
		account.SetReceivesUpdates(true)
		assert.True(t, account.GetReceivesUpdates())
		assert.True(t, account.GetModifiedAt().After(initialModifiedAt))
	})

	t.Run("SetCreatedAt", func(t *testing.T) {
		newTime := time.Now()
		account.SetCreatedAt(newTime)
		assert.Equal(t, newTime.UTC(), account.GetCreatedAt().UTC())
	})

	t.Run("SetModifiedAt", func(t *testing.T) {
		newTime := time.Now()
		account.SetModifiedAt(newTime)
		assert.Equal(t, newTime.UTC(), account.GetModifiedAt().UTC())
	})
}

func TestAccount_SetPassword_Pass(t *testing.T) {
	initialModifiedAt := time.Now()
	account := NewAccount("test@example.com", "Test User", initialModifiedAt)

	password := "123456"
	salt := "salt"

	err := account.SetPassword(password, salt)
	assert.NoError(t, err)

	// The hashed password is already stored as a string in the account
	hashedPassword := account.GetPasswordHash()
	logger.Infof("TEST", "hashed password '%s' and salt '%s'", hashedPassword, salt)

	// Verify the hash works
	err = bcrypt.CompareHashAndPassword([]byte(account.GetPasswordHash()), []byte(password+salt))
	assert.NoError(t, err)
	assert.True(t, account.GetModifiedAt().After(initialModifiedAt))
}

func TestAccount_SetPassword_Error(t *testing.T) {
	account := NewAccount("test@example.com", "Test User", time.Now())

	// Mock an extremely long password that might cause bcrypt to fail
	longPassword := string(make([]byte, 73)) // bcrypt has a maximum length
	err := account.SetPassword(longPassword, "salt")
	assert.Error(t, err)
}
