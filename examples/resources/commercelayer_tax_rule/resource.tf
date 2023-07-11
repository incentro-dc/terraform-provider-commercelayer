resource "commercelayer_manual_tax_calculator" "incentro_manual_tax_calculator" {
  attributes {
    name = "Incentro Manual Tax Calculator"
  }
}

resource "commercelayer_tax_rule" "incentro_tax_rule" {
  attributes {
    name = "Incentro Tax Rule"
    metadata = {
      foo : "bar"
      testName : "{{.testName}}"
    }
  }
  relationships {
    manual_tax_calculator_id = commercelayer_manual_tax_calculator.incentro_manual_tax_calculator.id
  }
}