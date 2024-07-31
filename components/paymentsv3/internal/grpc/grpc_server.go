package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/formancehq/paymentsv3/internal/grpc/proto"
	"github.com/formancehq/paymentsv3/internal/grpc/proto/services"
	"github.com/formancehq/paymentsv3/internal/plugins/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ services.PluginServer = &GRPCServer{}

type GRPCServer struct {
	services.UnimplementedPluginServer
	// This is the real implementation
	Impl models.Plugin
}

func (s *GRPCServer) Install(ctx context.Context, req *services.InstallRequest) (*services.InstallResponse, error) {
	resp, err := s.Impl.Install(ctx, models.InstallRequest{
		Config: req.GetConfig(),
	})
	if err != nil {
		return nil, translateErrorToGRPC(err)
	}

	var capabilities []proto.Capability
	for _, capability := range resp.Capabilities {
		capabilities = append(capabilities, proto.Capability(capability))
	}

	return &services.InstallResponse{
		Capabilities: capabilities,
		Workflow:     translateWorkflow(resp.Workflow),
	}, nil
}

func (s *GRPCServer) FetchNextAccounts(ctx context.Context, req *services.FetchNextAccountsRequest) (*services.FetchNextAccountsResponse, error) {
	resp, err := s.Impl.FetchNextAccounts(ctx, models.FetchNextAccountsRequest{
		FromPayload: req.GetFromPayload(),
		State:       req.GetState(),
	})
	if err != nil {
		return nil, translateErrorToGRPC(err)
	}

	accounts := make([]*proto.Account, 0, len(resp.Accounts))
	for _, account := range resp.Accounts {
		accounts = append(accounts, translateAccount(account))
	}

	return &services.FetchNextAccountsResponse{
		Accounts: accounts,
		NewState: resp.NewState,
	}, nil
}

func (s *GRPCServer) FetchNextPayments(ctx context.Context, req *services.FetchNextPaymentsRequest) (*services.FetchNextPaymentsResponse, error) {
	resp, err := s.Impl.FetchNextPayments(ctx, models.FetchNextPaymentsRequest{
		FromPayload: req.GetFromPayload(),
		State:       req.GetState(),
	})
	if err != nil {
		return nil, translateErrorToGRPC(err)
	}

	payments := make([]*proto.Payment, 0, len(resp.Payments))
	for _, payment := range resp.Payments {
		payments = append(payments, translatePayment(payment))
	}

	return &services.FetchNextPaymentsResponse{
		Payments: payments,
		NewState: resp.NewState,
	}, nil
}

func (s *GRPCServer) FetchNextOthers(ctx context.Context, req *services.FetchNextOthersRequest) (*services.FetchNextOthersResponse, error) {
	resp, err := s.Impl.FetchNextOthers(ctx, models.FetchNextOthersRequest{
		Name:  req.GetName(),
		State: req.GetState(),
	})
	if err != nil {
		return nil, translateErrorToGRPC(err)
	}

	others := make([][]byte, 0, len(resp.Others))
	for _, other := range resp.Others {
		others = append(others, other)
	}

	return &services.FetchNextOthersResponse{
		Payload:  others,
		NewState: resp.NewState,
	}, nil
}

func translateAccount(account models.Account) *proto.Account {
	return &proto.Account{
		Reference: account.Reference,
		Name: func() *wrapperspb.StringValue {
			if account.Name == nil {
				return nil
			}

			return wrapperspb.String(*account.Name)
		}(),
		CreatedAt: timestamppb.New(account.CreatedAt),
		SyncedAt:  timestamppb.New(time.Now().UTC()),
		DefaultAsset: func() *wrapperspb.StringValue {
			if account.DefaultAsset == nil {
				return nil
			}

			return wrapperspb.String(*account.DefaultAsset)
		}(),
		Metadata: account.Metadata,
		Raw:      account.Raw,
	}
}

func translatePayment(payment models.Payment) *proto.Payment {
	return &proto.Payment{
		Reference:   payment.Reference,
		CreatedAt:   timestamppb.New(payment.CreatedAt),
		SyncedAt:    timestamppb.New(time.Now().UTC()),
		PaymentType: proto.PaymentType(payment.PaymentType),
		Amount: &proto.Monetary{
			Asset:  payment.Asset,
			Amount: []byte(payment.Amount.Text(10)),
		},
		Scheme: proto.PaymentScheme(payment.Scheme),
		Status: proto.PaymentStatus(payment.Status),
		SourceAccountReference: func() *wrapperspb.StringValue {
			if payment.SourceAccountReference == nil {
				return nil
			}

			return wrapperspb.String(*payment.SourceAccountReference)
		}(),
		DestinationAccountReference: func() *wrapperspb.StringValue {
			if payment.DestinationAccountReference == nil {
				return nil
			}

			return wrapperspb.String(*payment.DestinationAccountReference)
		}(),
		Metadata: payment.Metadata,
		Raw:      payment.Raw,
	}
}

func translateTask(taskTree models.TaskTree) *proto.TaskTree {
	res := proto.TaskTree{
		NextTasks: []*proto.TaskTree{},
		Task:      nil,
	}

	switch taskTree.TaskType {
	case models.TASK_FETCH_ACCOUNTS:
		res.Task = &proto.TaskTree_FetchAccounts_{
			FetchAccounts: &proto.TaskTree_FetchAccounts{},
		}
	case models.TASK_FETCH_PAYMENTS:
		res.Task = &proto.TaskTree_FetchPayments_{
			FetchPayments: &proto.TaskTree_FetchPayments{},
		}
	case models.TASK_FETCH_OTHERS:
		res.Task = &proto.TaskTree_FetchOthers_{
			FetchOthers: &proto.TaskTree_FetchOthers{
				Name: taskTree.Name,
			},
		}
	default:
		// TODO(polo): better error handling
		panic("unknown task type")
	}

	for _, nextTask := range taskTree.NextTasks {
		res.NextTasks = append(res.NextTasks, translateTask(nextTask))
	}

	return &res
}

func translateWorkflow(workflows models.Workflow) *proto.Workflow {
	res := proto.Workflow{}

	for _, task := range workflows {
		res.Tasks = append(res.Tasks, translateTask(task))
	}

	return &res
}

func translateErrorToGRPC(err error) error {
	switch {
	case errors.Is(err, models.ErrInvalidConfig):
		return status.Errorf(codes.InvalidArgument, err.Error())
	default:
		return status.Errorf(codes.Internal, err.Error())
	}
}
