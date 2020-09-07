package grpc

import (
	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
	svc "github.com/owncloud/ocis-accounts/pkg/service/v0"
	"github.com/owncloud/ocis-pkg/v2/service/grpc"
)

// Server initializes a new go-micro service ready to run
func Server(opts ...Option) grpc.Service {
	options := newOptions(opts...)

	service := grpc.NewService(
		grpc.Name(options.Config.Server.Name),
		grpc.Context(options.Context),
		grpc.Address(options.Config.GRPC.Addr),
		grpc.Namespace(options.Config.GRPC.Namespace),
		grpc.Logger(options.Logger),
		grpc.Flags(options.Flags...),
	)

	var hdlr *svc.Service
	var err error

	if hdlr, err = svc.New(svc.Logger(options.Logger), svc.Config(options.Config)); err != nil {
		options.Logger.Fatal().Err(err).Msg("could not initialize service handler")
	}
	if err = proto.RegisterAccountsServiceHandler(service.Server(), hdlr); err != nil {
		options.Logger.Fatal().Err(err).Msg("could not register service handler")
	}
	if err = proto.RegisterGroupsServiceHandler(service.Server(), hdlr); err != nil {
		options.Logger.Fatal().Err(err).Msg("could not register groups handler")
	}

	service.Init()
	return service
}
