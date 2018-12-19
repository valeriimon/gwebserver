package mainservice

type Application struct {
	ID            string `json:"ID omitempty"`
	Database      string `json:"Database omitempty"`
	UserType      string `json:"UserType omitempty"`
	ProductType   string `json:"ProductType omitempty"`
	jwtPayload    map[string]interface{}
	AllowedRoutes []ApplicationRoute `json:"AllowedRoutes omitempty"`
}

type ApplicationRoute struct {
	Name        string
	Description string
}
