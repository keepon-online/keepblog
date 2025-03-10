package ua

import (
	"encoding/json"
	"fmt"
	"github.com/mssola/user_agent"
)

func UserAgent(ua string) {
	agent := user_agent.New(ua)
	indent, _ := json.MarshalIndent(&agent, "", "  ")
	fmt.Println("ua ", string(indent))
}
