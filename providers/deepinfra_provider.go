package providers

import (
	"bytes"
	"fmt"
	"github.com/google/wire"
	"github.com/sashabaranov/go-openai"
	"net/http"
	"steplems-bot/types"
)

type Interceptor struct {
	core http.RoundTripper
}

func (i Interceptor) RoundTrip(r *http.Request) (*http.Response, error) {
	defer func() {
		_ = r.Body.Close()
	}()

	headers := r.Header.Clone()
	var b bytes.Buffer
	if err := headers.Write(&b); err != nil {
		fmt.Println(err)
	}
	fmt.Println(b.String())

	// send the request using the DefaultTransport
	return i.core.RoundTrip(r)
}

func ProvideDeepInfraToken() (types.DeepInfraToken, error) {
	return ProvideEnvironmentVariable[types.DeepInfraToken]("DEEP_INFRA_TOKEN")()
}

func ProvideDeepInfraClient(token types.DeepInfraToken) *types.DeepInfraClient {
	config := openai.DefaultConfig(string(token))
	config.BaseURL = "https://api.deepinfra.com/v1/openai"
	//config.HTTPClient = &http.Client{
	//	Transport: Interceptor{http.DefaultTransport},
	//}
	oConfig := openai.NewClientWithConfig(config)
	client := (*types.DeepInfraClient)(oConfig)

	return client
}

var DeepInfraProviders = wire.NewSet(ProvideDeepInfraToken, ProvideDeepInfraClient)
