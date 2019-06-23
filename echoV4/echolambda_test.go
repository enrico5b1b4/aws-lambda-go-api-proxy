package echoV4adapter_test

import (
	echoV4adapter "github.com/enrico5b1b4/aws-lambda-go-api-proxy/echoV4"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/labstack/echo/v4"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EchoV4Lambda tests", func() {
	Context("Simple ping request", func() {
		It("Proxies the event correctly", func() {
			log.Println("Starting test")
			e := echo.New()
			e.GET("/ping", func(c echo.Context) error {
				log.Println("Handler!!")
				return c.String(200, "pong")
			})

			adapter := echoV4adapter.New(e)

			req := events.APIGatewayProxyRequest{
				Path:       "/ping",
				HTTPMethod: "GET",
			}

			resp, err := adapter.Proxy(req)

			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
		})
	})
})
