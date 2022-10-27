package commercelayer

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"text/template"
	"time"
)

var testAccProviderCommercelayer *schema.Provider
var testAccProviderFactories = map[string]func() (*schema.Provider, error){}

func TestMain(m *testing.M) {
	tokenFile, err := ioutil.TempFile("", "token")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("using token file %s\n", tokenFile.Name())

	testAccProviderCommercelayer = Provider(WithTokenCacheFile(tokenFile))()
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"commercelayer": func() (*schema.Provider, error) {
			return testAccProviderCommercelayer, nil
		},
	}

	retCode := m.Run()

	err = tokenFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = os.Remove(tokenFile.Name())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("removed token file %s\n", tokenFile.Name())

	os.Exit(retCode)
}

func TestProvider(t *testing.T) {
	provider := Provider()()
	if err := provider.InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	requiredEnvs := []string{
		"COMMERCELAYER_CLIENT_ID",
		"COMMERCELAYER_CLIENT_SECRET",
		"COMMERCELAYER_API_ENDPOINT",
		"COMMERCELAYER_AUTH_ENDPOINT",
	}
	for _, val := range requiredEnvs {
		if os.Getenv(val) == "" {
			t.Fatalf("%v must be set for acceptance tests", val)
		}
	}
}

func hclTemplate(data string, params map[string]any) string {
	var out bytes.Buffer
	tmpl := template.Must(template.New("hcl").Parse(data))
	err := tmpl.Execute(&out, params)
	if err != nil {
		log.Fatal(err)
	}
	return out.String()
}

func retryRemoval(times int, callable func() (*http.Response, error)) error {
	for retries := 1; retries < times; retries++ {
		resp, err := callable()
		if resp.StatusCode == 404 {
			return nil
		}
		if err != nil {
			return err
		}

		if resp.StatusCode == 200 {
			log.Println("retrying removal")
			time.Sleep(time.Second)
			continue
		}

		return fmt.Errorf("received response code with status %d", resp.StatusCode)
	}

	return nil
}
