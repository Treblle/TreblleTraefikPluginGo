package TreblleTraefikPluginGo

import (
	"fmt"
	"os"
)

func generateFieldsToMask(additionalFieldsToMask []string) map[string]bool {
	defaultFieldsToMask := []string{
		"password",
		"pwd",
		"secret",
		"password_confirmation",
		"passwordConfirmation",
		"cc",
		"card_number",
		"cardNumber",
		"ccv",
		"ssn",
		"credit_score",
		"creditScore",
	}

	fields := append(defaultFieldsToMask, additionalFieldsToMask...)
	fieldsToMask := make(map[string]bool)
	for _, field := range fields {
		fieldsToMask[field] = true
	}

	return fieldsToMask
}

func logError(err error) {
	prefix := "Error in treblle traefik plugin: "

	os.Stdout.WriteString(fmt.Sprintf("%v%v\n", prefix, err.Error()))
}
