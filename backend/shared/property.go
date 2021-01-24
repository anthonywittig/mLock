package shared

type PropertyX struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedBy string `json:"createdBy"`
}

type Property2 struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	ID        string `json:"id"` // Needed while migrating to dynamo only? Can't guarantee uniqueness?
	CreatedBy string `json:"createdBy"`
}
