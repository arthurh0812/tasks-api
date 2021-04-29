package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/arthurh0812/tasks-api/pkg/api"
	"github.com/kataras/iris/v12/context"
	"net/http"
	"strings"
)

var SetAccessControlHeaders context.Handler = func(ctx *context.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "POST,GET,OPTION")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
	ctx.Next()
}

func ExtractBearerToken(ctx *context.Context, authAPIService string) (string, error) {
	auth := ctx.Request().Header.Get("Authorization")
	if auth == "" {
		return "", fmt.Errorf("no authorization header provided")
	}
	var token string
	if parts := strings.Split(auth, " "); parts[0] != "Bearer" {
		return "", fmt.Errorf("authorization requires %q key", "Bearer")
	} else if len(parts) < 2 {
		return "", fmt.Errorf("no Bearer token provided")
	} else {
		token = parts[1]
	}
	res, err := http.Get(fmt.Sprintf("http://%s/verify-token/%s", authAPIService, token))
	if err != nil {
		return "", fmt.Errorf("request to authorization service failed: %v", err)
	}
	dec := json.NewDecoder(res.Body)
	defer res.Body.Close()
	apiRes := api.Response{}
	err = dec.Decode(&apiRes)
	if err != nil {
		return "", fmt.Errorf("failed to decode JSON API response from the authorization service: %v", err)
	}
	if apiRes.Status % 100 > 2 {
		return "", fmt.Errorf("request to authorization service returned error: %s", apiRes.Message)
	}
	if uid := apiRes.Data["uid"]; uid == nil {
		return "", fmt.Errorf("uid not provided")
	} else {
		return uid.(string), nil
	}
}