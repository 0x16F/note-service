package responses

const (
	internalErrorMessage = "Internal Error"
)

type Error struct {
	Code      int         `json:"code"`
	Message   interface{} `json:"message,omitempty"`
	Developer interface{} `json:"developer,omitempty"`
}
