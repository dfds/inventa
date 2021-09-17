package crossplane

import (
	"bytes"
	"context"
	"fmt"
	"github.com/dfds/inventa/operator/misc"
	"github.com/getoutreach/goql"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/json"
	"log"
	"net/http"
	"strings"
	"text/template"
)

type Client struct {
	goqlClient *goql.Client
	httpClient *http.Client
	queryEndpoint string
}

func (c *Client) GetCustomResourceDefinitions() []customResourceDefinitionsUnstructured {
	var resp []customResourceDefinitionsUnstructured
	query := &customResourceDefinitionsRequest{}
	err := c.goqlClient.Query(context.Background(), &goql.Operation{OperationType: query})
	if err != nil {
		log.Fatal(err)
	}

	for _, node := range query.CustomResourceDefinitions.Nodes {
		var output customResourceDefinitionsUnstructured
		decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			Metadata: nil,
			Result: &output,
			TagName: "json",
		})
		err := decoder.Decode(node.Unstructured)
		if err != nil {
			log.Fatal(err)
		}
		if strings.Contains(output.Spec.Group, "crossplane") {
			resp = append(resp, output)
		}
	}

	return resp
}

func (c *Client) GetKubernetesResources(apiVersion string, kind string, version string) []kubernetesResourceUnstructured {
	var resp []kubernetesResourceUnstructured
	query := &kubernetesResourceRequest{}
	vars := make(map[string]interface{})
	vars["apiVersion"] = fmt.Sprintf("%s/%s", apiVersion, version)
	vars["kind"] = kind

	err := c.goqlClient.Query(context.Background(), &goql.Operation{OperationType: query, Variables: vars})
	if err != nil {
		log.Fatal(err)
	}

	for _, node := range query.KubernetesResources.Nodes {
		var output kubernetesResourceUnstructured
		decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			Metadata: nil,
			Result: &output,
			TagName: "json",
		})
		err := decoder.Decode(node.Unstructured)
		if err != nil {
			log.Fatal(err)
		}
		resp = append(resp, output)
	}

	return resp
}

func (c * Client) GetKubernetesResourcesFromCustomResourceDefinitions(crds []customResourceDefinitionsUnstructured) kubernetesResourcesResponse {
	req := NewKubernetesResourceCombinedQuery(crds)
	var graphQlRequest struct {
		Query string `json:"query"`
	}
	graphQlRequest.Query = string(req)
	payload, err := json.Marshal(graphQlRequest)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := c.httpClient.Post(c.queryEndpoint, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var kubernetesResources kubernetesResourcesResponse
	err = json.Unmarshal(body, &kubernetesResources)
	if err != nil {
		log.Fatal(err)
	}

	return kubernetesResources
}

var kubernetesResourceQuery string = `
   {{.Alias}}: kubernetesResources(apiVersion: "{{.ApiVersion}}", kind: "{{.Kind}}")
{
    totalCount
    nodes {
      id
      unstructured
    }
  }`
var KubernetesResourceQueryTemplate *template.Template = SetupQueryTemplate()

func SetupQueryTemplate() *template.Template{
	tmpl, err := template.New("KubernetesResourceQuery").Parse(kubernetesResourceQuery)
	if err != nil {
		log.Fatal(err)
	}

	return tmpl
}

func NewKubernetesResourceCombinedQuery(crds []customResourceDefinitionsUnstructured) []byte {
	var buf bytes.Buffer
	buf.WriteString("{")
	for _, crd := range crds {
		groupName := strings.ReplaceAll(crd.Spec.Group, ".", "")
		bytes := newKubernetesResourceQueryObject(fmt.Sprintf("%s%s", groupName,crd.Spec.Names.Kind), fmt.Sprintf("%s/%s", crd.Spec.Group, crd.Spec.Versions[0].Name), crd.Spec.Names.Kind)
		buf.Write(bytes)
	}
	buf.WriteString("\n}")
	return buf.Bytes()
}

func newKubernetesResourceQueryObject(alias string, apiVersion string, kind string) []byte{
	var buf bytes.Buffer
	vars := make(map[string]string)
	vars["Alias"] = alias
	vars["ApiVersion"] = apiVersion
	vars["Kind"] = kind
	err := KubernetesResourceQueryTemplate.Execute(&buf, vars)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}

func NewClient() *Client {
	endpoint := misc.GetEnvValue("INVENTA_CROSSPLANE_ENDPOINT", "http://crossplane-host-placeholder:8080/query")
	tokenResult, err := misc.GetInClusterK8sToken()
	if err != nil {
		log.Println("Unable to find Pod service account token, defaulting to value of INVENTA_CROSSPLANE_TOKEN")
	}
	token := misc.GetEnvValue("INVENTA_CROSSPLANE_TOKEN", tokenResult)
	httpClient := &http.Client{
		Transport: &addAuthTransport{T: http.DefaultTransport, token: token},
	}

	goqlClient := goql.NewClient(endpoint, goql.ClientOptions{HTTPClient: httpClient})

	client := &Client{
		goqlClient: goqlClient,
		httpClient: httpClient,
		queryEndpoint: endpoint,
	}
	return client
}

type addAuthTransport struct {
	T http.RoundTripper
	token string
}

func (adt *addAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", adt.token))
	return adt.T.RoundTrip(req)
}