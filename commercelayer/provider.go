package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/incentro-dc/go-commercelayer-sdk/api"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var baseSchema = map[string]*schema.Schema{
	"client_id": {
		Type:        schema.TypeString,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc("COMMERCELAYER_CLIENT_ID", nil),
		Description: "The client id of a Commercelayer store",
		Sensitive:   true,
	},
	"client_secret": {
		Type:        schema.TypeString,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc("COMMERCELAYER_CLIENT_SECRET", nil),
		Description: "The client secret of a Commercelayer store",
		Sensitive:   true,
	},
	"api_endpoint": {
		Type:        schema.TypeString,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc("COMMERCELAYER_API_ENDPOINT", nil),
		Description: "The Commercelayer api endpoint",
	},
	"auth_endpoint": {
		Type:        schema.TypeString,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc("COMMERCELAYER_AUTH_ENDPOINT", nil),
		Description: "The Commercelayer auth endpoint",
	},
}

var baseResourceMap = map[string]*schema.Resource{
	"commercelayer_address":                   resourceAddress(),
	"commercelayer_merchant":                  resourceMerchant(),
	"commercelayer_price_list":                resourcePriceList(),
	"commercelayer_customer_group":            resourceCustomerGroup(),
	"commercelayer_webhook":                   resourceWebhook(),
	"commercelayer_external_gateway":          resourceExternalGateway(),
	"commercelayer_external_tax_calculator":   resourceExternalTaxCalculator(),
	"commercelayer_market":                    resourceMarket(),
	"commercelayer_inventory_model":           resourceInventoryModel(),
	"commercelayer_shipping_method":           resourceShippingMethod(),
	"commercelayer_shipping_zone":             resourceShippingZone(),
	"commercelayer_shipping_category":         resourceShippingCategory(),
	"commercelayer_stock_location":            resourceStockLocation(),
	"commercelayer_inventory_return_location": resourceInventoryReturnLocation(),
	"commercelayer_inventory_stock_location":  resourceInventoryStockLocation(),
	"commercelayer_delivery_lead_time":        resourceDeliveryLeadTime(),
	"commercelayer_manual_gateway":            resourceManualGateway(),
	"commercelayer_adyen_gateway":             resourceAdyenGateway(),
	"commercelayer_paypal_gateway":            resourcePaypalGateway(),
	"commercelayer_klarna_gateway":            resourceKlarnaGateway(),
	"commercelayer_braintree_gateway":         resourceBraintreeGateway(),
	"commercelayer_checkout_com_gateway":      resourceCheckoutComGateway(),
	"commercelayer_google_geocoder":           resourceGoogleGeocoders(),
	"commercelayer_bing_geocoder":             resourceBingGeocoders(),
	"commercelayer_payment_method":            resourcePaymentMethod(),
}

type Configuration struct {
	tokenSource oauth2.TokenSource
}

type ProviderOption func(configuration *Configuration)

func WithTokenSource(tokenSource oauth2.TokenSource) ProviderOption {
	return func(c *Configuration) {
		c.tokenSource = tokenSource
	}
}

func Provider(opts ...ProviderOption) plugin.ProviderFunc {
	c := Configuration{}

	for _, opt := range opts {
		opt(&c)
	}

	return func() *schema.Provider {
		return &schema.Provider{
			Schema:               baseSchema,
			ResourcesMap:         baseResourceMap,
			ConfigureContextFunc: c.configureFunc,
		}
	}
}

func (c *Configuration) configureFunc(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	clientId := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	apiEndpoint := d.Get("api_endpoint").(string)
	authEndpoint := d.Get("auth_endpoint").(string)

	credentials := clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		TokenURL:     authEndpoint,
		Scopes:       []string{},
	}

	newCtx := context.Background()

	var tokenSource = credentials.TokenSource(newCtx)
	if c.tokenSource != nil {
		tokenSource = c.tokenSource
	}

	httpClient := oauth2.NewClient(newCtx, tokenSource)

	commercelayerClient := api.NewAPIClient(&api.Configuration{
		HTTPClient: httpClient,
		Debug:      true,
		Servers: []api.ServerConfiguration{
			{URL: apiEndpoint},
		},
	})

	return commercelayerClient, nil
}
