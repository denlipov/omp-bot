package request

import (
	"context"
	"errors"
	"time"

	pb "github.com/denlipov/com-request-api/pkg/com-request-api"
	"github.com/denlipov/omp-bot/internal/config"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type RequestService interface {
	Describe(requestID uint64) (*pb.Request, error)
	List(limit uint64, offset uint64) ([]*pb.Request, error)
	Create(req pb.Request) (uint64, error)
	Update(requestID uint64, request pb.Request) error
	Remove(requestID uint64) (bool, error)
}

type DummyRequestService struct {
	connOpts []grpc.DialOption
}

func NewDummyRequestService() *DummyRequestService {
	return &DummyRequestService{
		connOpts: []grpc.DialOption{
			grpc.WithInsecure(),
		},
	}
}

func (s *DummyRequestService) newConnection() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(config.GetConfig().Grpc.ComRequestApiURI, s.connOpts...)
	if err != nil {
		log.Error().Msgf("fail to dial to communication request service: %+v", err)
		return nil, err
	}
	return conn, nil
}

func getTimeoutCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(),
		time.Duration(config.GetConfig().Grpc.ServiceConnTimeout)*time.Second)
}

func (s *DummyRequestService) Describe(requestID uint64) (*pb.Request, error) {
	conn, err := s.newConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewComRequestApiServiceClient(conn)
	ctx, cancel := getTimeoutCtx()
	defer cancel()

	resp, err := client.DescribeRequestV1(ctx, &pb.DescribeRequestV1Request{RequestId: requestID})
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

func (s *DummyRequestService) List(limit uint64, offset uint64) ([]*pb.Request, error) {
	conn, err := s.newConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewComRequestApiServiceClient(conn)
	ctx, cancel := getTimeoutCtx()
	defer cancel()

	resp, err := client.ListRequestV1(ctx, &pb.ListRequestV1Request{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	return resp.Request, nil
}

func (s *DummyRequestService) Create(req pb.Request) (uint64, error) {
	conn, err := s.newConnection()
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	client := pb.NewComRequestApiServiceClient(conn)
	ctx, cancel := getTimeoutCtx()
	defer cancel()

	resp, err := client.CreateRequestV1(ctx, &pb.CreateRequestV1Request{
		Request: &req,
	})
	if err != nil {
		return 0, err
	}

	return resp.RequestId, nil
}

func (s *DummyRequestService) Update(requestID uint64, request pb.Request) error {
	conn, err := s.newConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewComRequestApiServiceClient(conn)
	ctx, cancel := getTimeoutCtx()
	defer cancel()

	resp, err := client.UpdateRequestV1(ctx, &pb.UpdateRequestV1Request{
		RequestId: requestID,
		Body: &pb.UpdateRequestBody{
			Service: request.Service,
			User:    request.User,
			Text:    request.Text,
		},
	})
	if err != nil {
		return err
	}
	if !resp.Status {
		return errors.New("Serv returned status 'Not updated'")
	}
	return nil
}

func (s *DummyRequestService) Remove(requestID uint64) (bool, error) {
	conn, err := s.newConnection()
	if err != nil {
		return false, err
	}
	defer conn.Close()

	client := pb.NewComRequestApiServiceClient(conn)
	ctx, cancel := getTimeoutCtx()
	defer cancel()

	resp, err := client.RemoveRequestV1(ctx, &pb.RemoveRequestV1Request{
		RequestId: requestID,
	})
	if err != nil {
		return false, err
	}
	return resp.Status, nil
}
