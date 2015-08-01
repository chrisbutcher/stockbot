package main

import (
  "flag"
  "fmt"
  "github.com/chrisbutcher/stockbot/bingclient"
  "github.com/gin-gonic/gin"
  "net/http"
  "net/url"
  "regexp"
)

type SlackCommand struct {
  Token       string  `form:"token" binding:"required"`
  TeamID      string  `form:"team_id" binding:"required"`
  TeamDomain  string  `form:"team_domain" binding:"required"`
  ChannelID   string  `form:"channel_id" binding:"required"`
  ChannelName string  `form:"channel_name" binding:"required"`
  UserID      string  `form:"user_id" binding:"required"`
  Timestamp   float64 `form:"timestamp" binding:"required"`
  UserName    string  `form:"user_name" binding:"required"`
  TriggerWord string  `form:"trigger_word" binding:"required"`
  Text        string  `form:"text" binding:"required"`
}

func buildWebhookPayload(imageUrl, replyChannelName string) string {
  return "{\"channel\": \"#" + replyChannelName + "\", \"username\": \"stockphotos\", \"text\": \"" + imageUrl + "\", \"icon_emoji\": \":necktie:\"}"
}

func prepareSearchTerms(triggerWord, rawSearchTerms string) string {
  re := regexp.MustCompile(triggerWord + " ")
  return re.ReplaceAllString(rawSearchTerms, "")
}

func main() {
  port := flag.String("port", "3000", "HTTP port")
  slackToken := flag.String("slack_token", "debug", "Slack verification token")
  slackWebhookUrl := flag.String("slack_webhook_url", "debug", "Slack response webhook URL")
  bingApiToken := flag.String("bing_token", "debug", "Bing auth token")
  flag.Parse()

  r := gin.New()
  r.Use(gin.Logger())
  r.Use(gin.Recovery())

  client := bingclient.BingClient{ApiToken: *bingApiToken, Market: "en-US", Adult: "Moderate", Format: "json"}

  r.GET("/", func(c *gin.Context) {
    c.String(http.StatusOK, "bot active")
  })

  r.POST("/bot/stockphotos", func(c *gin.Context) {
    var slackCmd SlackCommand
    c.Bind(&slackCmd)

    if slackCmd.Token != *slackToken {
      c.String(http.StatusUnauthorized, "not authorized")
      return
    }

    searchTerms := prepareSearchTerms(slackCmd.TriggerWord, slackCmd.Text)

    imageUrl, _ := client.FetchImageRandomImage(searchTerms)
    fmt.Println(imageUrl)
    c.String(http.StatusFound, "Sending webhook")

    webhookPayload := buildWebhookPayload(imageUrl, slackCmd.ChannelName)

    _, err := http.PostForm(*slackWebhookUrl,
      url.Values{"payload": {webhookPayload}})

    if err != nil {
      panic(err)
    }
  })

  fmt.Println("Listening on port " + *port)
  r.Run(":" + *port)
}
