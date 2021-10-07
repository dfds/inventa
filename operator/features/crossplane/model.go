package crossplane

type customResourceDefinitionsRequest struct {
	CustomResourceDefinitions struct {
		TotalCount int
		Nodes []struct {
			Metadata struct {
				Name string
			}
			Unstructured map[string]interface{}
		}
	}
}

type kubernetesResourceRequest struct {
	KubernetesResources struct {
		TotalCount int
		Nodes []struct {
			Id string
			Unstructured map[string]interface{}
		}
	} `goql:"kubernetesResources(apiVersion:$apiVersion<String!>,kind:$kind<String!>)"`
}

type kubernetesResourcesResponse struct {
	Data map[string] KubernetesResourcesResponseData
}

type KubernetesResourcesResponseData struct {
	TotalCount int
	Nodes []struct {
		Id string
		Unstructured map[string]interface{}
	}
}

type customResourceDefinitionsUnstructured struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		Name              string    `json:"name"`
	} `json:"metadata"`
	Spec struct {
		Group string `json:"group"`
		Names struct {
			Plural     string   `json:"plural"`
			Singular   string   `json:"singular"`
			Kind       string   `json:"kind"`
			ListKind   string   `json:"listKind"`
			Categories []string `json:"categories"`
		} `json:"names"`
		Scope    string `json:"scope"`
		Versions []struct {
			Name string `json:"name"`
		} `json:"versions"`
	} `json:"spec"`
}

type kubernetesResourceUnstructured struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		Name              string    `json:"name"`
		Annotations       map[string]string `json:"annotations"`
	} `json:"metadata"`
}