package crossplane

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"log"
	"testing"
)

func TestGetCustomResourceDefinitions(t *testing.T) {
	client := NewClient()
	crds := client.GetCustomResourceDefinitions()
	resources := client.GetKubernetesResourcesFromCustomResourceDefinitions(crds)
	for key, resource := range resources.Data {
		fmt.Printf("%s: %v\n", key, resource.TotalCount)
		for _, node := range resource.Nodes {
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
			fmt.Printf("  %s\n", output.Metadata.Name)
		}
	}
}