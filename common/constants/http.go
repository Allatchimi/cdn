package constants

import "fmt"

var HTTP_500_ERROR_MESSAGE = func(message string) error {
	return fmt.Errorf("%s", fmt.Sprintf("Error occurred when trying to %s! Please try again later.", message))
}

var HTTP_401_INVALID_TOKEN_ERROR_MESSAGE = func() error {
	return fmt.Errorf("%s", "Invalid or expired token! Please enter valid information.")
}

var HTTP_401_INVALID_PERMISSION_ERROR_MESSAGE = func() error {
	return fmt.Errorf("%s", "You don't have permission to access this resource! Please enter valid information.")
}

var HTTP_404_ERROR_MESSAGE = func(message string) error {
	return fmt.Errorf("%s", fmt.Sprintf("%s not found! Please enter valid information.", message))
}

var HTTP_400_ERROR_MESSAGE = func(message string) error {
	return fmt.Errorf("%s", fmt.Sprintf("%s! Please enter valid information.", message))
}
