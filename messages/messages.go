package messages

var (
	// NotAdminError is an error handler
	NotAdminError = "You are not an admin"
	// NotDeveloperError is an error handler
	NotDeveloperError = "You are not registered as a developer"
	// NotGuarantorError is an error handler
	NotGuarantorError = "You are not registered as a guarantor"
	// NotInvestorError is an error handler
	NotInvestorError = "You are not registered as a investor"
	// NotRecipientError is an error handler
	NotRecipientError = "You are not registered as a recipient"
	// URLEmptyError is an error handler
	URLEmptyError = "Requst URL is empty"
	// EmptyValueError is an error handler
	EmptyValueError = "Empty value passed in request"
	// TokenError is an error handler
	TokenError = "Length of token not 32 characters"
	// TooLongError is an error handler
	TooLongError = "Request value too long"
	// RelayError is an error handler
	RelayError = "Error in relaying request to openx"
	// NotUserError is an error handler
	NotUserError = "You are not registered as a user"
	// TickerError is an error handler
	TickerError = "Unable to fetch exchange rate from ticker"
	// NotEntityError is an error handler
	NotEntityError = "You are not registered as the specified entity"
	// ConversionError is an error handler
	ConversionError = "Type Conversion Error"
)

// ParamError returns a sdtring containing the ParamError
func ParamError(param string) string {
	return "Required param: " + param + " not found"
}
