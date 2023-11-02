package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zelas91/gofermart/internal/entities"
	errService "github.com/zelas91/gofermart/internal/error"
	"github.com/zelas91/gofermart/internal/logger"
	"github.com/zelas91/gofermart/internal/service"
	"net/http"
	"strconv"
	"time"
)

const defaultRetryAfter = time.Second

type Client struct {
	client  *http.Client
	baseURL string
	service service.Orders
}

func New(baseURL string, service service.Orders) *Client {
	return &Client{
		&http.Client{
			Timeout: time.Second * 1,
		},
		baseURL,
		service,
	}
}

func (c *Client) StartService(ctx context.Context) {
	c.fetchOrder(ctx)
}
func (c *Client) fetchOrder(ctx context.Context) {
	ticker := time.NewTicker(defaultRetryAfter)

	go func() {
		retryErr := false
		for {
			if ctx.Err() != nil {
				return
			}
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if retryErr {
					retryErr = false
					ticker.Reset(defaultRetryAfter)
				}

				orders, err := c.service.GetOrdersWithoutFinalStatuses(ctx)
				if err != nil {
					logger.GetLogger(ctx).Info("get orders  ", err)
					return
				}

				for _, order := range orders {

					oa, err := c.getOrderAccrual(ctx, order)

					if errClient := new(errService.ErrHTTPClient); errors.As(err, &errClient) {
						if errClient.StatusCode == http.StatusTooManyRequests {
							retryErr = true
							ticker.Reset(errClient.RetryTime)
							return
						}

						logger.GetLogger(ctx).Errorf("accrual get request err: %v", err)
						continue

					}
					if err != nil {
						logger.GetLogger(ctx).Errorf("get order accrual err : %v", err)
						return
					}
					if err = c.service.UpdateOrder(ctx, oa); err != nil {
						logger.GetLogger(ctx).Errorf("update order err: %v", err)
					}
				}
			}
		}
	}()

}

func (c *Client) getOrderAccrual(ctx context.Context, order entities.Order) (entities.OrderAccrual, error) {
	URL := fmt.Sprintf("%s/api/orders/%s", c.baseURL, order.Number)
	resp, err := c.client.Get(URL)
	if err != nil {
		return entities.OrderAccrual{}, err
	}
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusTooManyRequests {
			return entities.OrderAccrual{}, &errService.ErrHTTPClient{Msg: "status code not equals 200",
				StatusCode: resp.StatusCode, RetryTime: c.getRetryAfterByResp(ctx, resp)}
		}
		return entities.OrderAccrual{}, &errService.ErrHTTPClient{Msg: "status code not equals 200",
			StatusCode: resp.StatusCode, RetryTime: defaultRetryAfter}
	}

	var oa entities.OrderAccrual

	if err = json.NewDecoder(resp.Body).Decode(&oa); err != nil {
		return oa, err
	}

	return oa, err

}
func (c *Client) getRetryAfterByResp(ctx context.Context, r *http.Response) time.Duration {
	respRetryAfter := r.Header.Get("Retry-After")
	if respRetryAfter != "" {
		retryAfter, err := strconv.Atoi(respRetryAfter)
		if err == nil {
			return time.Duration(retryAfter) * time.Second
		}
		logger.GetLogger(ctx).Error(err)
	}

	return defaultRetryAfter
}
