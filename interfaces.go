package jsonrpc

type (

	// Method used as RPC method handler.
	Method interface {
		// GetName returns method name in string representation.
		GetName() string

		// GetParamsType says to Router a structure (or nil) which must be used for params parsing (fields, etc.).
		GetParamsType() interface{}

		// Handle will be called by Router when method with current name will be requested.
		Handle(params interface{}) (interface{}, Error)
	}

	// Router is used for methods registration and invoking.
	Router interface {
		// RegisterMethod make a method registration for later invoking.
		RegisterMethod(method Method) error

		// MethodIsRegistered returns `true` only if passed method is registered.
		MethodIsRegistered(methodName string) bool

		// Invoke accepts method name and invoke registered method with same name.
		Invoke(methodName string, params interface{}) (interface{}, Error)
	}

	// Error is general RPC error.
	Error interface {
		error

		// GetCode returns error code.
		GetCode() int

		// GetMessage returns error message.
		GetMessage() string

		// GetData returns error extra-data.
		GetData() interface{}
	}

	// Validator allows to validate different structures, like method params (but not only).
	Validator interface {
		// IsValid returns `error` only if structure has INCORRECT state or properties.
		Validate() error
	}
)
