package commercelayer

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
	"strings"
)

func testAccCheckPaymentMethodDestroy(s *terraform.State) error {
	client := testAccProviderCommercelayer.Meta().(*commercelayer.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "commercelayer_payment_method" {
			_, resp, err := client.PaymentMethodsApi.GETPaymentMethodsPaymentMethodId(context.Background(), rs.Primary.ID).Execute()
			if resp.StatusCode == 404 {
				fmt.Printf("commercelayer_payment_method with id %s has been removed\n", rs.Primary.ID)
				continue
			}
			if err != nil {
				return err
			}

			return fmt.Errorf("received response code with status %d", resp.StatusCode)
		}

		if rs.Type == "commercelayer_market" {
			_, resp, err := client.MarketsApi.GETMarketsMarketId(context.Background(), rs.Primary.ID).Execute()
			if resp.StatusCode == 404 {
				fmt.Printf("commercelayer_market with id %s has been removed\n", rs.Primary.ID)
				continue
			}
			if err != nil {
				return err
			}

			return fmt.Errorf("received response code with status %d", resp.StatusCode)
		}

		if rs.Type == "commercelayer_adyen_gateway" {
			_, resp, err := client.AdyenGatewaysApi.
				GETAdyenGatewaysAdyenGatewayId(context.Background(), rs.Primary.ID).Execute()
			if resp.StatusCode == 404 {
				fmt.Printf("commercelayer_adyen_gateway with id %s has been removed\n", rs.Primary.ID)
				continue
			}
			if err != nil {
				return err
			}

			return fmt.Errorf("received response code with status %d", resp.StatusCode)
		}
	}
	return nil
}

func (s *AcceptanceSuite) TestAccPaymentMethod_basic() {
	resourceName := "commercelayer_payment_method.incentro_payment_method"

	resource.Test(s.T(), resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(s)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckPaymentMethodDestroy,
		Steps: []resource.TestStep{
			{
				Config: strings.Join([]string{testAccAdyenGatewayCreate(resourceName), testAccPaymentMethodCreate(resourceName)}, "\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", paymentMethodType),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.currency_code", "EUR"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.payment_source_type", "AdyenPayment"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.price_amount_cents", "0"),
				),
			},
			{
				Config: strings.Join([]string{testAccAdyenGatewayCreate(resourceName), testAccPaymentMethodUpdate(resourceName)}, "\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "attributes.0.metadata.bar", "foo"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.currency_code", "EUR"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.payment_source_type", "AdyenPayment"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.price_amount_cents", "0"),
				),
			},
		},
	})
}

func testAccPaymentMethodCreate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_payment_method" "incentro_payment_method" {
		  attributes {
      		payment_source_type   = "AdyenPayment"
			currency_code          = "EUR"
			price_amount_cents     = 0
			metadata               = {
			  foo : "bar"
		 	  testName: "{{.testName}}"
			}
		  }

		  relationships {
			payment_gateway_id = commercelayer_adyen_gateway.incentro_adyen_gateway.id
		  }
		}
	`, map[string]any{"testName": testName})
}

func testAccPaymentMethodUpdate(testName string) string {
	return hclTemplate(`
		resource "commercelayer_payment_method" "incentro_payment_method" {
		  attributes {
      		payment_source_type    = "AdyenPayment"
			currency_code          = "EUR"
			price_amount_cents     = 0
			metadata               = {
			  bar : "foo"
		 	  testName: "{{.testName}}"
			}
		  }
  		  relationships {
			payment_gateway_id = commercelayer_adyen_gateway.incentro_adyen_gateway.id
		  }
		}
	`, map[string]any{"testName": testName})
}
