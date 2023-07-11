package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceTaxRule() *schema.Resource {
	return &schema.Resource{
		Description: "Create a tax rule to be associated with a manual tax calculator. " +
			"Use the available regular expressions to match the order shipping method and apply the " +
			"specified tax_rate (default is 0). You can optionally define the tax rule taxable preferences by " +
			"setting the freight_taxable, payment_method_taxable, gift_card_taxable, " +
			"and adjustment_taxable attributes to true (all false by default).",
		ReadContext:   resourceTaxRuleReadFunc,
		CreateContext: resourceTaxRuleCreateFunc,
		UpdateContext: resourceTaxRuleUpdateFunc,
		DeleteContext: resourceTaxRuleDeleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The tax rule unique identifier",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"type": {
				Description: "The resource type",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"attributes": {
				Description: "Resource attributes",
				Type:        schema.TypeList,
				MaxItems:    1,
				MinItems:    1,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The tax rule internal name.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"tax_rate": {
							Description: "The tax rate for this rule.",
							Type:        schema.TypeFloat,
							Optional:    true,
							Default:     0.0,
						},
						"country_code_regex": {
							Description: "The regex that will be evaluated to match the shipping address country code.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"not_country_code_regex": {
							Description: "The regex that will be evaluated as negative match for the shipping " +
								"address country code.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"state_code_regex": {
							Description: "The regex that will be evaluated to match the shipping address state code.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"not_state_code_regex": {
							Description: "The regex that will be evaluated as negative match for the shipping " +
								"address state code.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"zip_code_regex": {
							Description: "The regex that will be evaluated to match the shipping address zip code.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"not_zip_code_regex": {
							Description: "The regex that will be evaluated as negative match for the shipping zip country code.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"freight_taxable": {
							Description: "Indicates if the freight is taxable.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"payment_method_taxable": {
							Description: "Indicates if the payment method is taxable.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"gift_card_taxable": {
							Description: "Indicates if gift cards are taxable.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"adjustment_taxable": {
							Description: "Indicates if adjustments are taxable.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"reference": {
							Description: "A string that you can use to add any external identifier to the resource. This " +
								"can be useful for integrating the resource to an external system, like an ERP, a " +
								"marketing tool, a CRM, or whatever.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"reference_origin": {
							Description: "Any identifier of the third party system that defines the reference code",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"metadata": {
							Description: "Set of key-value pairs that you can attach to the resource. This can be useful " +
								"for storing additional information about the resource in a structured format",
							Type: schema.TypeMap,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
					},
				},
			},
			"relationships": {
				Description: "Resource relationships",
				Type:        schema.TypeList,
				MaxItems:    1,
				MinItems:    1,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"manual_tax_calculator_id": {
							Description: "The associated manual tax calculator id.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func resourceTaxRuleReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.TaxRulesApi.GETTaxRulesTaxRuleId(ctx, d.Id()).Execute()
	if err != nil {
		return diagErr(err)
	}

	taxRule, ok := resp.GetDataOk()
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(taxRule.GetId())

	return nil
}

func resourceTaxRuleCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))
	relationships := nestedMap(d.Get("relationships"))

	taxRuleCreate := commercelayer.TaxRuleCreate{
		Data: commercelayer.TaxRuleCreateData{
			Type: taxRulesType,
			Attributes: commercelayer.POSTTaxRules201ResponseDataAttributes{
				Name:                 attributes["name"].(string),
				TaxRate:              float64ToFloat32Ref(attributes["tax_rate"]),
				CountryCodeRegex:     stringRef(attributes["country_code_regex"]),
				NotCountryCodeRegex:  stringRef(attributes["not_country_code_regex"]),
				StateCodeRegex:       stringRef(attributes["state_code_regex"]),
				NotStateCodeRegex:    stringRef(attributes["not_state_code_regex"]),
				ZipCodeRegex:         stringRef(attributes["zip_code_regex"]),
				NotZipCodeRegex:      stringRef(attributes["not_zip_code_regex"]),
				FreightTaxable:       boolRef(attributes["freight_taxable"]),
				PaymentMethodTaxable: boolRef(attributes["payment_method_taxable"]),
				GiftCardTaxable:      boolRef(attributes["gift_card_taxable"]),
				AdjustmentTaxable:    boolRef(attributes["adjustment_taxable"]),
				Reference:            stringRef(attributes["reference"]),
				ReferenceOrigin:      stringRef(attributes["reference_origin"]),
				Metadata:             keyValueRef(attributes["metadata"]),
			},
			Relationships: &commercelayer.TaxRuleCreateDataRelationships{
				ManualTaxCalculator: commercelayer.TaxRuleCreateDataRelationshipsManualTaxCalculator{
					Data: commercelayer.TaxRuleDataRelationshipsManualTaxCalculatorData{
						Type: stringRef(manualTaxCalculatorsType),
						Id:   stringRef(relationships["manual_tax_calculator_id"].(string)),
					},
				},
			},
		},
	}

	err := d.Set("type", taxRulesType)
	if err != nil {
		return diagErr(err)
	}

	taxRule, _, err := c.TaxRulesApi.POSTTaxRules(ctx).TaxRuleCreate(taxRuleCreate).Execute()
	if err != nil {
		return diagErr(err)
	}

	d.SetId(*taxRule.Data.Id)

	return nil
}

func resourceTaxRuleDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.TaxRulesApi.DELETETaxRulesTaxRuleId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourceTaxRuleUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))
	relationships := nestedMap(d.Get("relationships"))

	var TaxRuleUpdate = commercelayer.TaxRuleUpdate{
		Data: commercelayer.TaxRuleUpdateData{
			Type: taxRulesType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHTaxRulesTaxRuleId200ResponseDataAttributes{
				Name:                 stringRef(attributes["name"].(string)),
				TaxRate:              float64ToFloat32Ref(attributes["tax_rate"]),
				CountryCodeRegex:     stringRef(attributes["country_code_regex"]),
				NotCountryCodeRegex:  stringRef(attributes["not_country_code_regex"]),
				StateCodeRegex:       stringRef(attributes["state_code_regex"]),
				NotStateCodeRegex:    stringRef(attributes["not_state_code_regex"]),
				ZipCodeRegex:         stringRef(attributes["zip_code_regex"]),
				NotZipCodeRegex:      stringRef(attributes["not_zip_code_regex"]),
				FreightTaxable:       boolRef(attributes["freight_taxable"]),
				PaymentMethodTaxable: boolRef(attributes["payment_method_taxable"]),
				GiftCardTaxable:      boolRef(attributes["gift_card_taxable"]),
				AdjustmentTaxable:    boolRef(attributes["adjustment_taxable"]),
				Reference:            stringRef(attributes["reference"]),
				ReferenceOrigin:      stringRef(attributes["reference_origin"]),
				Metadata:             keyValueRef(attributes["metadata"]),
			},
			Relationships: &commercelayer.TaxRuleUpdateDataRelationships{
				ManualTaxCalculator: &commercelayer.TaxRuleCreateDataRelationshipsManualTaxCalculator{
					Data: commercelayer.TaxRuleDataRelationshipsManualTaxCalculatorData{
						Type: stringRef(manualTaxCalculatorsType),
						Id:   stringRef(relationships["manual_tax_calculator_id"].(string)),
					},
				},
			},
		},
	}

	_, _, err := c.TaxRulesApi.PATCHTaxRulesTaxRuleId(ctx, d.Id()).TaxRuleUpdate(TaxRuleUpdate).Execute()

	return diag.FromErr(err)
}
