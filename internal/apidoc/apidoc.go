package apidoc

type StructureDocs struct {
	Endpoints []EndpointDocs `json:"endpoints"`
}

type EndpointDocs struct {
	Path   string   `json:"path"`
	Params []string `json:"params"`
	Method string   `json:"method"`
}

func GetStructure() StructureDocs {
	return StructureDocs{
		Endpoints: []EndpointDocs{
			EndpointDocs{
				Path:   "/user",
				Params: []string{"username", "password"},
				Method: "POST",
			},
			EndpointDocs{
				Path:   "/vhost",
				Params: []string{"username", "vhost", "ip"},
				Method: "POST",
			},
			EndpointDocs{
				Path:   "/resource",
				Params: []string{"username", "vhost", "resource", "name", "permission"},
				Method: "POST",
			},
			EndpointDocs{
				Path:   "/topic",
				Params: []string{"username", "vhost", "resource", "name", "permission", "routing_key"},
				Method: "POST",
			},
		},
	}
}
