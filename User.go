package googledrive

type User struct {
	Kind         string `json:"kind"`
	DisplayName  string `json:"displayName"`
	PhotoLink    string `json:"photoLink"`
	Me           bool   `json:"me"`
	PermissionId string `json:"permissionId"`
	EmailAddress string `json:"emailAddress"`
}
