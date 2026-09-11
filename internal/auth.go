package internal
import (
	"errors"
	"strings"
	"time"
)
const jwtSecret = "super-secret-key-12345"
type Claims struct {
	Username string
	Expiry time.Time
}
func GenerateToken(username string) (string, error) {
	if username == "" {
		return "", errors.New("username cannot be empty")
	}
	expiry := time.Now().Add(24 * time.Hour)
	token := username + "." + expiry.Format(time.RFC3339) + "." + jwtSecret
	return token, nil
}
func ValidateToken(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", errors.New("invalid security token format")
	}
	if parts[2] != jwtSecret {
		return "", errors.New("token signature mismatch")
	}
	return parts[0], nil
}