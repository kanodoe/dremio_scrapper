package entity

// DremioLoginRequest ....
type DremioLoginRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

// DremioLoginResponse ...
type DremioLoginResponse struct {
	Token                     string `json:"token"`
	UserName                  string `json:"userName"`
	FirstName                 string `json:"firstName"`
	LastName                  string `json:"lastName"`
	Expires                   int64  `json:"expires"`
	Email                     string `json:"email"`
	UserID                    string `json:"userId"`
	Admin                     bool   `json:"admin"`
	ClusterID                 string `json:"clusterId"`
	ClusterCreatedAt          int64  `json:"clusterCreatedAt"`
	ShowUserAndUserProperties bool   `json:"showUserAndUserProperties"`
	Version                   string `json:"version"`
	Permissions               struct {
		CanUploadProfiles   bool `json:"canUploadProfiles"`
		CanDownloadProfiles bool `json:"canDownloadProfiles"`
		CanEmailForSupport  bool `json:"canEmailForSupport"`
		CanChatForSupport   bool `json:"canChatForSupport"`
	} `json:"permissions"`
	UserCreatedAt int64 `json:"userCreatedAt"`
}
