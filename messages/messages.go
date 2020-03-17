package messages

var (
	NotAdminError     = "You are not an admin on opensolar"
	NotDeveloperError = "You are not registered as the specified entity"
	NotEntityError    = "You are not registered as the specified entity"
	ConversionError   = "Error while converting between types"
)

// ParamError returns a sdtring containing the ParamError
func ParamError(param string) string {
	return "Required param: " + param + " not found"
}
