package openai

import "fmt"

func fullURL(sub string) string {
	return fmt.Sprintf("%s/%s", openaiUrlv1, sub)
}
