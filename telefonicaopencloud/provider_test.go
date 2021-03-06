package telefonicaopencloud

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/helper/pathorcontents"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var (
	OS_DB_ENVIRONMENT         = os.Getenv("OS_DB_ENVIRONMENT")
	OS_DB_DATASTORE_VERSION   = os.Getenv("OS_DB_DATASTORE_VERSION")
	OS_DB_DATASTORE_TYPE      = os.Getenv("OS_DB_DATASTORE_TYPE")
	OS_DEPRECATED_ENVIRONMENT = os.Getenv("OS_DEPRECATED_ENVIRONMENT")
	OS_DNS_ENVIRONMENT        = os.Getenv("OS_DNS_ENVIRONMENT")
	OS_EXTGW_ID               = os.Getenv("OS_EXTGW_ID")
	OS_FLAVOR_ID              = os.Getenv("OS_FLAVOR_ID")
	OS_FLAVOR_NAME            = os.Getenv("OS_FLAVOR_NAME")
	OS_IMAGE_ID               = os.Getenv("OS_IMAGE_ID")
	OS_IMAGE_NAME             = os.Getenv("OS_IMAGE_NAME")
	OS_NETWORK_ID             = os.Getenv("OS_NETWORK_ID")
	OS_VPC_ID                 = os.Getenv("OS_VPC_ID")
	OS_POOL_NAME              = os.Getenv("OS_POOL_NAME")
	OS_REGION_NAME            = os.Getenv("OS_REGION_NAME")
	OS_SWIFT_ENVIRONMENT      = os.Getenv("OS_SWIFT_ENVIRONMENT")
	OS_AVAILABILITY_ZONE      = os.Getenv("OS_AVAILABILITY_ZONE")
	OS_ACCESS_KEY             = os.Getenv("OS_ACCESS_KEY")
	OS_SECRET_KEY             = os.Getenv("OS_SECRET_KEY")
	OS_SRC_ACCESS_KEY         = os.Getenv("OS_SRC_ACCESS_KEY")
	OS_SRC_SECRET_KEY         = os.Getenv("OS_SRC_SECRET_KEY")
	OS_TENANT_ID              = os.Getenv("OS_TENANT_ID")
	OS_BMS_FLAVOR_NAME        = os.Getenv("OS_BMS_FLAVOR_NAME")
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"telefonicaopencloud": testAccProvider,
	}
}

func testAccPreCheckRequiredEnvVars(t *testing.T) {
	v := os.Getenv("OS_AUTH_URL")
	if v == "" {
		t.Fatal("OS_AUTH_URL must be set for acceptance tests")
	}

	if OS_IMAGE_ID == "" && OS_IMAGE_NAME == "" {
		t.Fatal("OS_IMAGE_ID or OS_IMAGE_NAME must be set for acceptance tests")
	}

	if OS_POOL_NAME == "" {
		t.Fatal("OS_POOL_NAME must be set for acceptance tests")
	}

	if OS_FLAVOR_ID == "" && OS_FLAVOR_NAME == "" {
		t.Fatal("OS_FLAVOR_ID or OS_FLAVOR_NAME must be set for acceptance tests")
	}

	if OS_NETWORK_ID == "" {
		t.Fatal("OS_NETWORK_ID must be set for acceptance tests")
	}

	if OS_VPC_ID == "" {
		t.Fatal("OS_VPC_ID must be set for acceptance tests")
	}

	if OS_EXTGW_ID == "" {
		t.Fatal("OS_EXTGW_ID must be set for acceptance tests")
	}
}

func testAccPreCheck(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	// Do not run the test if this is a deprecated testing environment.
	if OS_DEPRECATED_ENVIRONMENT != "" {
		t.Skip("This environment only runs deprecated tests")
	}

	// Do not run the test if this is a standalone DNS environment.
	if OS_DNS_ENVIRONMENT != "" {
		t.Skip("This environment only runs DNS tests")
	}
}

func testAccPreCheckDeprecated(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if OS_DEPRECATED_ENVIRONMENT == "" {
		t.Skip("This environment does not support deprecated tests")
	}
}

func testAccPreCheckDNS(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if OS_DNS_ENVIRONMENT == "" {
		t.Skip("This environment does not support DNS tests")
	}
}

func testAccPreCheckSwift(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if OS_SWIFT_ENVIRONMENT == "" {
		t.Skip("This environment does not support Swift tests")
	}
}

func testAccPreCheckDatabase(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if OS_DB_ENVIRONMENT == "" {
		t.Skip("This environment does not support Database tests")
	}
}

func testAccPreCheckAdminOnly(t *testing.T) {
	v := os.Getenv("OS_USERNAME")
	if v != "admin" {
		t.Skip("Skipping test because it requires the admin user")
	}
}

func testAccPreCheckObs(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if OS_ACCESS_KEY == "" {
		t.Fatal("OS_ACCESS_KEY must be set for obs acceptance tests")
	}

	if OS_SECRET_KEY == "" {
		t.Fatal("OS_SECRET_KEY must be set for obs acceptance tests")
	}
}

func testAccPreCheckELB(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if OS_TENANT_ID == "" {
		t.Skip("This environment does not support ELB tests")
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

// Steps for configuring TelefonicaOpenCloud with SSL validation are here:
// https://github.com/hashicorp/terraform/pull/6279#issuecomment-219020144
func TestAccProvider_caCertFile(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping TelefonicaOpenCloud SSL test.")
	}
	if os.Getenv("OS_CACERT") == "" {
		t.Skip("OS_CACERT is not set; skipping TelefonicaOpenCloud CA test.")
	}

	p := Provider()

	caFile, err := envVarFile("OS_CACERT")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(caFile)

	raw := map[string]interface{}{
		"cacert_file": caFile,
	}
	rawConfig, err := config.NewRawConfig(raw)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = p.Configure(terraform.NewResourceConfig(rawConfig))
	if err != nil {
		t.Fatalf("Unexpected err when specifying TelefonicaOpenCloud CA by file: %s", err)
	}
}

func TestAccProvider_caCertString(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping TelefonicaOpenCloud SSL test.")
	}
	if os.Getenv("OS_CACERT") == "" {
		t.Skip("OS_CACERT is not set; skipping TelefonicaOpenCloud CA test.")
	}

	p := Provider()

	caContents, err := envVarContents("OS_CACERT")
	if err != nil {
		t.Fatal(err)
	}
	raw := map[string]interface{}{
		"cacert_file": caContents,
	}
	rawConfig, err := config.NewRawConfig(raw)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = p.Configure(terraform.NewResourceConfig(rawConfig))
	if err != nil {
		t.Fatalf("Unexpected err when specifying TelefonicaOpenCloud CA by string: %s", err)
	}
}

func TestAccProvider_clientCertFile(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping TelefonicaOpenCloud SSL test.")
	}
	if os.Getenv("OS_CERT") == "" || os.Getenv("OS_KEY") == "" {
		t.Skip("OS_CERT or OS_KEY is not set; skipping TelefonicaOpenCloud client SSL auth test.")
	}

	p := Provider()

	certFile, err := envVarFile("OS_CERT")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(certFile)
	keyFile, err := envVarFile("OS_KEY")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(keyFile)

	raw := map[string]interface{}{
		"cert": certFile,
		"key":  keyFile,
	}
	rawConfig, err := config.NewRawConfig(raw)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = p.Configure(terraform.NewResourceConfig(rawConfig))
	if err != nil {
		t.Fatalf("Unexpected err when specifying TelefonicaOpenCloud Client keypair by file: %s", err)
	}
}

func TestAccProvider_clientCertString(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("OS_SSL_TESTS") == "" {
		t.Skip("TF_ACC or OS_SSL_TESTS not set, skipping TelefonicaOpenCloud SSL test.")
	}
	if os.Getenv("OS_CERT") == "" || os.Getenv("OS_KEY") == "" {
		t.Skip("OS_CERT or OS_KEY is not set; skipping TelefonicaOpenCloud client SSL auth test.")
	}

	p := Provider()

	certContents, err := envVarContents("OS_CERT")
	if err != nil {
		t.Fatal(err)
	}
	keyContents, err := envVarContents("OS_KEY")
	if err != nil {
		t.Fatal(err)
	}

	raw := map[string]interface{}{
		"cert": certContents,
		"key":  keyContents,
	}
	rawConfig, err := config.NewRawConfig(raw)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = p.Configure(terraform.NewResourceConfig(rawConfig))
	if err != nil {
		t.Fatalf("Unexpected err when specifying TelefonicaOpenCloud Client keypair by contents: %s", err)
	}
}

func envVarContents(varName string) (string, error) {
	contents, _, err := pathorcontents.Read(os.Getenv(varName))
	if err != nil {
		return "", fmt.Errorf("Error reading %s: %s", varName, err)
	}
	return contents, nil
}

func envVarFile(varName string) (string, error) {
	contents, err := envVarContents(varName)
	if err != nil {
		return "", err
	}

	tmpFile, err := ioutil.TempFile("", varName)
	if err != nil {
		return "", fmt.Errorf("Error creating temp file: %s", err)
	}
	if _, err := tmpFile.Write([]byte(contents)); err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", fmt.Errorf("Error writing temp file: %s", err)
	}
	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", fmt.Errorf("Error closing temp file: %s", err)
	}
	return tmpFile.Name(), nil
}

func testAccAsConfigPreCheck(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)
	if OS_FLAVOR_ID == "" {
		t.Skip("OS_FLAVOR_ID must be set for acceptance tests")
	}
}

func testAccBmsFlavorPreCheck(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)
	if OS_BMS_FLAVOR_NAME == "" {
		t.Skip("Provide the bms name starting with 'physical'")
	}
}

func testAccPreCheckMaas(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)

	if OS_ACCESS_KEY == "" || OS_SECRET_KEY == "" || OS_SRC_ACCESS_KEY == "" || OS_SRC_SECRET_KEY == "" {
		t.Skip("OS_ACCESS_KEY, OS_SECRET_KEY, OS_SRC_ACCESS_KEY, and OS_SRC_SECRET_KEY  must be set for MAAS acceptance tests")
	}
}
