package middleware

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetCookieSettings() (string, bool, bool, error) {
	env := os.Getenv("GO_ENV")

	var domain string
	var secure, httpOnly bool
	var err error

	if env == "development" {
		domain = os.Getenv("DEV_DOMAIN")
		domain = extractDomain(domain)
		secure, err = strconv.ParseBool(os.Getenv("DEV_SECURE_COOKIE"))
		if err != nil {
			return "", false, false, err
		}
		httpOnly, err = strconv.ParseBool(os.Getenv("DEV_HTTP_ONLY_COOKIE"))
		if err != nil {
			return "", false, false, err
		}
	} else {
		domain = os.Getenv("PROD_DOMAIN")
		domain = extractDomain(domain)
		secure, err = strconv.ParseBool(os.Getenv("PROD_SECURE_COOKIE"))
		if err != nil {
			return "", false, false, err
		}
		httpOnly, err = strconv.ParseBool(os.Getenv("PROD_HTTP_ONLY_COOKIE"))
		if err != nil {
			return "", false, false, err
		}
	}

	return domain, secure, httpOnly, nil
}

func GetLogoutCookieSettings() (string, error) {
	env := os.Getenv("ENV")
	if env == "" {
		return "", fmt.Errorf("environment variable 'ENV' is not set")
	}

	var originEnvVar string
	if env == "development" {
		originEnvVar = os.Getenv("DEV_DOMAIN")
		if originEnvVar == "" {
			return "", fmt.Errorf("environment variable 'DEV_DOMAIN' is not set")
		}
	} else {
		originEnvVar = os.Getenv("PROD_DOMAIN")
		if originEnvVar == "" {
			return "", fmt.Errorf("environment variable 'PROD_DOMAIN' is not set")
		}
	}

	domain := extractDomain(originEnvVar)
	fmt.Println("Extracted domain:", domain)
	if domain == "" {
		return "", fmt.Errorf("failed to extract domain from origin: %s", originEnvVar)
	}

	return domain, nil
}

func extractDomain(fullOrigin string) string {
	if strings.Contains(fullOrigin, "//") {
		parts := strings.Split(fullOrigin, "//")
		fullOrigin = parts[1]
	}
	if strings.Contains(fullOrigin, ":") {
		parts := strings.Split(fullOrigin, ":")
		fullOrigin = parts[0]
	}
	return fullOrigin
}
