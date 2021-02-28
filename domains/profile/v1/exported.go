package testv1

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/soldatov-s/go-garage-profile/internal/cfg"
	"github.com/soldatov-s/go-garage/domains"
	"github.com/soldatov-s/go-garage/providers/db/pq"
	"github.com/soldatov-s/go-garage/providers/httpsrv/echo"
	"github.com/soldatov-s/go-garage/providers/logger"
)

const (
	DomainName = "profilev1"
)

type empty struct{}

type ProfileV1 struct {
	log zerolog.Logger
	ctx context.Context
	db  *pq.Enity
}

func Registrate(ctx context.Context) (context.Context, error) {
	t := &ProfileV1{
		ctx: ctx,
		log: logger.GetPackageLogger(ctx, empty{}),
	}
	var err error
	if t.db, err = pq.GetEnityTypeCast(ctx, cfg.DBName); err != nil {
		return nil, err
	}

	privateV1, err := echo.GetAPIVersionGroup(ctx, cfg.PrivateHTTP, cfg.V1)
	if err != nil {
		return nil, err
	}

	grProtect := privateV1.Group
	grProtect.Use(echo.HydrationLogger(&t.log))
	grProtect.GET("/profile/:id", echo.Handler(t.profileGetHandler))
	grProtect.POST("/profile", echo.Handler(t.profilePostHandler))
	grProtect.DELETE("/profile/:id", echo.Handler(t.profileDeleteHandler))
	grProtect.PUT("/profiles/:id", echo.Handler(t.profilePutHandler))
	grProtect.POST("/users/search", echo.Handler(t.profileSearchPostHandler))

	return domains.RegistrateByName(ctx, DomainName, t), nil
}

func Get(ctx context.Context) (*ProfileV1, error) {
	if v, ok := domains.GetByName(ctx, DomainName).(*ProfileV1); ok {
		return v, nil
	}
	return nil, domains.ErrInvalidDomainType
}
