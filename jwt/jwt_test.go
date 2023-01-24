package jwt_test

import (
	"context"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/infiniteloopcloud/go/jwt"
	pjwt "github.com/pascaldekloe/jwt"
)

var testIssuer = "test issuer"

//nolint:cyclop // cyclomatic complexity was too large because of the asserts
func TestJwtECDSA_CreateAndVerify(t *testing.T) {
	ctx := context.Background()
	token := Token{
		Token:     uuid.New().String(),
		UserID:    "test_1234",
		ExpiresAt: time.Now().UTC().Add(10 * time.Hour),
		Type:      0,
	}

	user := User{
		ID:       "test_1234",
		Name:     "John Doe",
		Password: "secret password",
		Phone:    "123456789",
		Email:    "test_email@gmail.com",
		Status:   1<<63 | 1<<62 | 1<<2 | 1<<1 | 1,
	}
	result, err := jwt.Create(ctx, jwt.Metadata{
		Issuer:     testIssuer,
		ClientHost: "test_client_host",
	}, user, token)
	if err != nil {
		t.Error(err)
	}
	expected := "eyJhbGciOiJFUzI1NiJ9.eyJhdWQiOlsidGVzdF9jbGllbnRfaG9zdCJdLCJl"
	if !strings.Contains(string(result), expected) {
		t.Errorf("Token should contain %s, rather like to be %s", expected, string(result))
	}

	claims, err := jwt.Verify(ctx, jwt.Metadata{
		Issuer:     testIssuer,
		ClientHost: "test_client_host",
	}, result)
	if err != nil {
		t.Error(err)
	}

	if len(claims.Audiences) != 1 || claims.Audiences[0] != "test_client_host" {
		t.Errorf("Audience should be [%s] instead of %v", "test_client_host", claims.Audiences)
	}

	if claims.Registered.Issuer != testIssuer {
		t.Errorf("issuer should be %s instead of %s", testIssuer, claims.Registered.Issuer)
	}

	if claims.Registered.Subject != user.ID {
		t.Errorf("subject should be %s instead of %s", user.ID, claims.Registered.Issuer)
	}

	if jwtExpires := pjwt.NewNumericTime(token.ExpiresAt.Round(time.Second)); claims.Registered.Expires.String() != jwtExpires.String() {
		t.Errorf("expires should be %s instead of %s", jwtExpires, claims.Registered.Expires)
	}

	if claims.ID != token.Token {
		t.Errorf("id should be %s instead of %s", token.Token, claims.Registered.ID)
	}

	if name, ok := claims.Set["name"]; !ok {
		t.Errorf("missing name")
	} else if name != user.Name {
		t.Errorf("name should be %s instead of %s", user.Name, name)
	}

	if email, ok := claims.Set["email"]; !ok {
		t.Errorf("missing email")
	} else if email != user.Email {
		t.Errorf("email should be %s instead of %s", user.Email, email)
	}

	if phone, ok := claims.Set["phone"]; !ok {
		t.Errorf("missing phone")
	} else if phone != user.Phone {
		t.Errorf("phone should be %s instead of %s", user.Phone, phone)
	}

	if untypedStatus, ok := claims.Set["status"]; !ok {
		t.Errorf("missing status")
	} else if statusStr, ok := untypedStatus.(string); ok {
		status, err := strconv.ParseUint(statusStr, 10, 64)
		if err != nil {
			t.Error(err)
		}
		if user.Status != status {
			t.Errorf("status should be %d instead of %d", user.Status, status)
		}
	}
}

type Token struct {
	Token     string    `json:"token" model:"token"`
	ExpiresAt time.Time `json:"expires_at" model:"expires_at"`
	UserID    string    `json:"user_id" model:"user_id"`
	Type      uint8     `json:"type" model:"type"`
	CreatedAt time.Time `json:"created_at" model:"created_at"`
	UpdatedAt time.Time `json:"updated_at" model:"updated_at"`
}

func (t Token) ClaimsParse() map[string]interface{} {
	var claims = make(map[string]interface{})
	claims["token"] = t.Token
	claims["token_expires_at"] = t.ExpiresAt
	return claims
}

func (t Token) Expired() bool {
	return !t.ExpiresAt.After(time.Now().UTC())
}

type User struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Password           string    `json:"-"`
	PreviousPasswords  string    `json:"-"`
	Salt               string    `json:"-"`
	PasswordCryptoType uint8     `json:"-"`
	Phone              string    `json:"phone"`
	Email              string    `json:"email"`
	Status             uint64    `json:"status"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	DeletedAt          time.Time `json:"deleted_at"`
}

func (u User) ClaimsParse() map[string]interface{} {
	var claims = make(map[string]interface{})
	claims["user_id"] = u.ID
	claims["name"] = u.Name
	claims["email"] = u.Email
	claims["phone"] = u.Phone
	claims["status"] = strconv.FormatUint(u.Status, 10)
	return claims
}
