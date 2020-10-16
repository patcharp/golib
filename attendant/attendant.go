package attendant

// AttendantEndpoint API Endpoint
const (
	pkgName            = "ATTENDANT"
	DefaultApiEndpoint = "https://attendants.sdi.inet.co.th/3rd"
)

// AttendantAPIResult struct
type APIResult struct {
	Error   interface{} `json:"error"`
	Message interface{} `json:"msg"`
	Data    interface{} `json:"data"`
	Total   int         `json:"total"`
	Count   int         `json:"count"`
}

// Client struct
type Client struct {
	token       string
	tokenType   string
	apiEndpoint string
}

// New Attendant client
func NewClient(token string, tokenType string) (Client, error) {
	var attendant Client
	if tokenType == "" {
		tokenType = "Bearer"
	}
	attendant = Client{
		token:       token,
		tokenType:   tokenType,
		apiEndpoint: DefaultApiEndpoint,
	}
	return attendant, nil
}

func (client *Client) SetApiEndpoint(ep string) {
	client.apiEndpoint = ep
}
