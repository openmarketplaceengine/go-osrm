package types

// ResponseStatus represent OSRM API response
type ResponseStatus struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	DataVersion string `json:"data_version"`
}

// ErrCode returns error code from OSRM response
func (r ResponseStatus) ErrCode() string {
	return r.Code
}

func (r ResponseStatus) Error() string {
	return r.Code + " - " + r.Message
}

func (r ResponseStatus) ApiError() error {
	if r.Code != errorCodeOK {
		return r
	}
	return nil
}
