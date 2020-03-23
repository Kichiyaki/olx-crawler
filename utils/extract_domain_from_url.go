package utils

import "regexp"

func ExtractDomainFromURL(url string) ([]string, error) {
	re, err := regexp.Compile(`^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/\n]+)`)
	if err != nil {
		return nil, err
	}
	return re.FindAllString(url, -1), nil
}
