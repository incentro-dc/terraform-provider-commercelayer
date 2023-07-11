package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceManualTaxCalculator() *schema.Resource {
	return &schema.Resource{
		Description: "Configure the manual tax calculator by creating one or more associated tax rules. " +
			"The rules will apply the related tax rate to the matching orders.",
		ReadContext:   resourceManualTaxCalculatorReadFunc,
		CreateContext: resourceManualTaxCalculatorCreateFunc,
		UpdateContext: resourceManualTaxCalculatorUpdateFunc,
		DeleteContext: resourceManualTaxCalculatorDeleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The manual tax calculator unique identifier",
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
							Description: "The tax calculator's internal name.",
							Type:        schema.TypeString,
							Required:    true,
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
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tax_rule_id": {
							Description: "The associated tax rule id.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func resourceManualTaxCalculatorReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.ManualTaxCalculatorsApi.GETManualTaxCalculatorsManualTaxCalculatorId(ctx, d.Id()).Execute()
	if err != nil {
		return diagErr(err)
	}

	manualTaxCalculator, ok := resp.GetDataOk()
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(manualTaxCalculator.GetId())

	return nil
}

func resourceManualTaxCalculatorCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))
	//relationships := nestedMap(d.Get("relationships"))

	manualTaxCalculatorCreate := commercelayer.ManualTaxCalculatorCreate{
		Data: commercelayer.ManualTaxCalculatorCreateData{
			Type: manualTaxCalculatorsType,
			Attributes: commercelayer.POSTManualTaxCalculators201ResponseDataAttributes{
				Name:            attributes["name"].(string),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
			//Relationships: &commercelayer.ManualTaxCalculatorCreateDataRelationships{
			//	TaxRules: &commercelayer.ManualTaxCalculatorDataRelationshipsTaxRules{
			//		Data: commercelayer.ManualTaxCalculatorDataRelationshipsTaxRulesData{
			//			Type: stringRef(taxRulesType),
			//			Id:   stringRef(relationships["tax_rule_id"].(string)),
			//		}},
			//},
		},
	}

	err := d.Set("type", manualTaxCalculatorsType)
	if err != nil {
		return diagErr(err)
	}

	manualTaxCalculator, _, err := c.ManualTaxCalculatorsApi.POSTManualTaxCalculators(ctx).ManualTaxCalculatorCreate(manualTaxCalculatorCreate).Execute()
	if err != nil {
		return diagErr(err)
	}

	d.SetId(*manualTaxCalculator.Data.Id)

	return nil
}

func resourceManualTaxCalculatorDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.ManualTaxCalculatorsApi.DELETEManualTaxCalculatorsManualTaxCalculatorId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourceManualTaxCalculatorUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))
	//relationships := nestedMap(d.Get("relationships"))

	var manualTaxCalculatorUpdate = commercelayer.ManualTaxCalculatorUpdate{
		Data: commercelayer.ManualTaxCalculatorUpdateData{
			Type: manualTaxCalculatorsType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHManualTaxCalculatorsManualTaxCalculatorId200ResponseDataAttributes{
				Name:            stringRef(attributes["name"]),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
			//Relationships: &commercelayer.ManualTaxCalculatorCreateDataRelationships{
			//	TaxRules: &commercelayer.ManualTaxCalculatorDataRelationshipsTaxRules{
			//		Data: commercelayer.ManualTaxCalculatorDataRelationshipsTaxRulesData{
			//			Type: stringRef(manualTaxCalculatorsType),
			//			Id:   stringRef(relationships["tax_rule_id"].(string)),
			//		},
			//	},
			//},
		},
	}

	_, _, err := c.ManualTaxCalculatorsApi.PATCHManualTaxCalculatorsManualTaxCalculatorId(ctx, d.Id()).
		ManualTaxCalculatorUpdate(manualTaxCalculatorUpdate).Execute()

	return diag.FromErr(err)
}
