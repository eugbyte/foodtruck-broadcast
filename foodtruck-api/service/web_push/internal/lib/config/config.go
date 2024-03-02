package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/samber/lo"
)

func conn() *string {
	value := os.Getenv("LOCALSTACK_HOSTNAME")
	value = strings.TrimSpace(value)

	var isDev = value != ""
	if !isDev {
		return nil
	}

	return lo.ToPtr(fmt.Sprintf("http://%s:4566", value))
}

var DBConn = conn
var QConn = conn

const VAPID_PRIVATE_KEY = "2O1psvs_Q_uUoy5yW24D82tIPI3gNUfdPbp89Ygjguo"
const VAPID_PUBLIC_KEY = "BIvYk7KnTMwTVY9ubq55Gkt0FzMVo2Rsm7Cs43MEhFJjXBwu073JvS6oEH56CMdK7FFcjW6mqZgF_zdsnDgCczA"
const VAPID_EMAIL = "eugenetham1994@gmail.com"
