package messages

var (
	NotAdminError     = "You are not an admin on opensolar"
	NotDeveloperError = "You are not registered as a developer"
	NotGuarantorError = "You are not registered as a guarantor"
	NotInvestorError  = "You are not registered as a investor"
	NotRecipientError = "You are not registered as a recipient"
	URLEmptyError     = "Requst URL is Empty"
	EmptyValueError   = "Value passed in request is empty"
	TokenError        = "Length of token not 32 characters"
	TooLongError      = "Value passed in request too long"
	RelayError        = "Error in realying request to openx"
	NotUserError      = "You are not registered as a user"
	TickerError       = "Unable to fetch exchange rate from ticker"
	NotEntityError    = "You are not registered as the specified entity"
	ConversionError   = "Error while converting between types"
)

// ParamError returns a sdtring containing the ParamError
func ParamError(param string) string {
	return "Required param: " + param + " not found"
}
