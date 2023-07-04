package tokens

import (
	"context"
	"log"
	"time"

	"github.com/amleonc/tabula/internal/helpers/keys"
	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func TokenWithClaims(claims map[string]any) (string, error) {
	now := time.Now()
	t, err := jwt.NewBuilder().
		Issuer("tabula").
		Audience([]string{"tabula-user"}).
		IssuedAt(now).
		Expiration(now.Add(time.Hour * 24 * 30)).
		Build()
	if err != nil {
		return "", err
	}
	for k, v := range claims {
		if err = t.Set(k, v); err != nil {
			log.Println("error setting claim")
			return "", err
		}
	}
	tokenBytes, err := jwt.Sign(t, jwt.WithKey(jwa.RS256, keys.PrivateRSAKey()))
	if err != nil {
		return "", err
	}
	return string(tokenBytes), nil
}

func ClaimsFromToken(ctx context.Context, t jwt.Token) (map[string]any, error) {
	claims, err := t.AsMap(ctx)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func UserIDFromToken(ctx context.Context, i any) (uuid.UUID, error) {
	t := i.(jwt.Token)
	claims, err := ClaimsFromToken(ctx, t)
	if err != nil {
		return uuid.UUID{}, err
	}
	s := claims["id"].(string)
	id, err := uuid.FromString(s)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}
