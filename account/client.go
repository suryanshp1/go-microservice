package account

import (
	"context"

	"github.com/suryanshp1/go-microservice/account/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.AccountServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	service := pb.NewAccountServiceClient(conn)
	return &Client{conn: conn, service: service}, nil
}

func (c *Client) Close() error {
	c.conn.Close()
	return nil
}

func (c *Client) PostAccount(ctx context.Context, name string) (*Account, error) {
	req := &pb.PostAccountRequest{Name: name}
	resp, err := c.service.PostAccount(ctx, req)
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:   resp.Account.Id,
		Name: resp.Account.Name,
	}, nil
}

func (c *Client) GetAccount(ctx context.Context, id string) (*Account, error) {
	req := &pb.GetAccountRequest{Id: id}
	resp, err := c.service.GetAccount(ctx, req)
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:   resp.Account.Id,
		Name: resp.Account.Name,
	}, nil
}

func (c *Client) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]*Account, error) {
	req := &pb.GetAccountsRequest{Skip: skip, Take: take}
	resp, err := c.service.GetAccounts(ctx, req)
	if err != nil {
		return nil, err
	}
	accounts := make([]*Account, len(resp.Accounts))
	for i, acc := range resp.Accounts {
		accounts[i] = &Account{
			ID:   acc.Id,
			Name: acc.Name,
		}
	}
	return accounts, nil
}
