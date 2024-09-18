package crypto

import (
    "testing"
)

func TestCheckHashPasswordMatch(t *testing.T) {
    password := "abc123"
    hash, err := HashPassword(password)
    if err != nil {
        t.Fatalf("Error while hashing password: %v", err)
    }
    if err = CheckPasswordHash(password, hash); err != nil {
        t.Fatalf("Password doesn't match its hash")
    }
}

func TestCheckHashPasswordNoMatch(t *testing.T) {
    hash, err := HashPassword("abc123")
    if err != nil {
        t.Fatalf("Error while hashing password: %v", err)
    }
    if err = CheckPasswordHash("123abc", hash); err == nil {
        t.Fatalf("Hash shouldn't match the wrong password")
    }
}
