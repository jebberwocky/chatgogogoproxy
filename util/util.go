package util

import (
	"chatproxy/models"
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

const (
	HttpSignatureHeader            string = "X-Signature"
	HttpResponseHeaderRecipientKey string = "X-Recipient"
	HttpResponseHeaderRecipient    string = "e9c6bcdd-5a98-45ca-8091-4959ce2f28a8" //use random v4 uuid "hermas"
)

const (
	Environment             string = "ENVIRONMENT"
	SecretKey               string = "SecretKey"
	DomainWhitelistEnable   string = "DOMAIN_WHITELIST_ENABLE"
	DomainWhitelistValues   string = "DOMAIN_WHITELISTS"
	DomainWhiteListAllowAll string = "*"
	ServerPort              string = "PORT"
	UserUploadImgDir        string = "USER_UPLOAD_IMG_DIR"
)

const (
	Production      string = "PRD"
	Development     string = "DEV"
	LogLevel        string = "LOG_LEVEL"
	EnvironmentFile string = ".env"
	LogFilePath     string = "LOG_FILEPATH"
	LogSize         string = "LOG_SIZE"
)

const (
	EchoAppContext string = "APP_ECHO_CONTEXT"
)

func GetQuery(d interface{}, key string) string {
	v := reflect.ValueOf(d)
	t := v.Type()
	var parts []string
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		structField := t.Field(i)
		// Get query tag
		queryTag := structField.Tag.Get("query")
		queryValue := fmt.Sprintf("%v", field.Interface())
		part := fmt.Sprintf("%s=%s", queryTag, url.QueryEscape(queryValue))
		parts = append(parts, part)
	}
	return strings.Join(parts, "&")
}

func FindInChatContent(c models.ChatContent, key string) string {
	for _, content := range c.Content {
		if content.Type == key {
			if key == "image_url" {
				return content.ImgUrl.Url
			} else {
				return content.Text
			}
		}
	}
	return ""
}
