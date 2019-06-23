// Packge echoV4lambda add EchoV4 support for the aws-severless-go-api library.
// Uses the core package behind the scenes and exposes the New method to
// get a new instance and Proxy method to send request to the echo.Echo
package echoV4adapter

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/labstack/echo/v4"
)

// EchoV4Lambda makes it easy to send API Gateway proxy events to a echo.Echo.
// The library transforms the proxy event into an HTTP request and then
// creates a proxy response object from the http.ResponseWriter
type EchoV4Lambda struct {
	core.RequestAccessor

	Echo *echo.Echo
}

// New creates a new instance of the EchoV4Lambda object.
// Receives an initialized *echo.Echo object - normally created with echo.New().
// It returns the initialized instance of the EchoV4Lambda object.
func New(e *echo.Echo) *EchoV4Lambda {
	return &EchoV4Lambda{Echo: e}
}

// Proxy receives an API Gateway proxy event, transforms it into an http.Request
// object, and sends it to the echo.Echo for routing.
// It returns a proxy response object generated from the http.ResponseWriter.
func (e *EchoV4Lambda) Proxy(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	echoRequest, err := e.ProxyEventToHTTPRequest(req)
	return e.proxyInternal(echoRequest, err)
}

// ProxyWithContext receives context and an API Gateway proxy event,
// transforms them into an http.Request object, and sends it to the echo.Echo for routing.
// It returns a proxy response object generated from the http.ResponseWriter.
func (e *EchoV4Lambda) ProxyWithContext(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	echoRequest, err := e.EventToRequestWithContext(ctx, req)
	return e.proxyInternal(echoRequest, err)
}

func (e *EchoV4Lambda) proxyInternal(req *http.Request, err error) (events.APIGatewayProxyResponse, error) {

	if err != nil {
		return core.GatewayTimeout(), core.NewLoggedError("Could not convert proxy event to request: %v", err)
	}

	respWriter := core.NewProxyResponseWriter()
	e.Echo.ServeHTTP(http.ResponseWriter(respWriter), req)

	proxyResponse, err := respWriter.GetProxyResponse()
	if err != nil {
		return core.GatewayTimeout(), core.NewLoggedError("Error while generating proxy response: %v", err)
	}

	return proxyResponse, nil
}
