package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
	"strconv"
	"time"
)

const DefaultSlackTimeout = 5 * time.Second

type SlackClient struct {
	WebHookUrl string
	UserName   string
	Channel    string
	TimeOut    time.Duration
}

type SimpleSlackRequest struct {
	Text      string
	IconEmoji string
}

type SlackJobNotification struct {
	Color     string
	IconEmoji string
	Details   string
	Text      string
}

type SlackMessage struct {
	Username    string       `json:"username,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	Text        string       `json:"text,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Color         string `json:"color,omitempty"`
	Fallback      string `json:"fallback,omitempty"`
	CallbackID    string `json:"callback_id,omitempty"`
	ID            int    `json:"id,omitempty"`
	AuthorID      string `json:"author_id,omitempty"`
	AuthorName    string `json:"author_name,omitempty"`
	AuthorSubname string `json:"author_subname,omitempty"`
	AuthorLink    string `json:"author_link,omitempty"`
	AuthorIcon    string `json:"author_icon,omitempty"`
	Title         string `json:"title,omitempty"`
	TitleLink     string `json:"title_link,omitempty"`
	Pretext       string `json:"pretext,omitempty"`
	Text          string `json:"text,omitempty"`
	ImageURL      string `json:"image_url,omitempty"`
	ThumbURL      string `json:"thumb_url,omitempty"`
	// Fields and actions are not defined.
	MarkdownIn []string    `json:"mrkdwn_in,omitempty"`
	Ts         json.Number `json:"ts,omitempty"`
}

// SendSlackNotification will post to an 'Incoming Webook' url setup in Slack Apps. It accepts
// some text and the slack channel is saved within Slack.
func (sc SlackClient) SendSlackNotification(sr SimpleSlackRequest) error {
	slackRequest := SlackMessage{
		Text:      sr.Text,
		Username:  sc.UserName,
		IconEmoji: sr.IconEmoji,
		Channel:   sc.Channel,
	}
	return sc.sendHttpRequest(slackRequest)
}

func (sc SlackClient) SendJobNotification(job SlackJobNotification) error {
	attachment := Attachment{
		Color: job.Color,
		Text:  job.Details,
		Ts:    json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	slackRequest := SlackMessage{
		Text:        job.Text,
		Username:    sc.UserName,
		IconEmoji:   job.IconEmoji,
		Channel:     sc.Channel,
		Attachments: []Attachment{attachment},
	}
	return sc.sendHttpRequest(slackRequest)
}

func (sc SlackClient) SendError(message string, options ...string) (err error) {
	return sc.funcName("danger", message, options)
}

func (sc SlackClient) SendInfo(message string, options ...string) (err error) {
	return sc.funcName("good", message, options)
}

func (sc SlackClient) SendWarning(message string, options ...string) (err error) {
	return sc.funcName("warning", message, options)
}

func (sc SlackClient) funcName(color string, message string, options []string) error {
	emoji := ":hammer_and_wrench"
	if len(options) > 0 {
		emoji = options[0]
	}
	sjn := SlackJobNotification{
		Color:     color,
		IconEmoji: emoji,
		Details:   message,
	}
	return sc.SendJobNotification(sjn)
}
func (sc SlackClient) sendHttpRequest(slackRequest SlackMessage) error {
	slackBody, _ := json.Marshal(slackRequest)
	req, err := http.NewRequest(http.MethodPost, sc.WebHookUrl, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	if sc.TimeOut == 0 {
		sc.TimeOut = DefaultSlackTimeout
	}
	client := &http.Client{Timeout: sc.TimeOut}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	if buf.String() != "ok" {
		return errors.New("Non-ok response returned from Slack")
	}
	return nil
}

// AWS Lambda
type SimpleType struct {
	Version    string     `json:"version"`
	ID         string     `json:"id"`
	DetailType string     `json:"detail-type"`
	Source     string     `json:"source"`
	Time       string     `json:"time"`
	Region     string     `json:"region"`
	Resources  []string   `json:"resources"`
	Account    string     `json:"account"`
	Detail     DetailType `json:"detail"`
}

type DetailType struct {
	ScanStatus            string                    `json:"scan-status"`
	RepositoryName        string                    `json:"repository-name"`
	FindingSeverityCounts FindingSeverityCountsType `json:"finding-severity-counts"`
	ImageDigest           string                    `json:"image-digest"`
	ImageTags             []string                  `json:"image-tags"`
}

type FindingSeverityCountsType struct {
	Critical int64 `json:"CRITICAL"`
	High     int64 `json:"HIGH"`
	Medium   int64 `json:"MEDIUM"`
}

func HandleRequest(ctx context.Context, event SimpleType) (events.APIGatewayProxyResponse, error) {
	sc := SlackClient{
		WebHookUrl: "https://hooks.slack.com/services/XXXXX/XXXXXXX/XXXXXXX",
		UserName:   "incoming-webhook",
		Channel:    "test-monitoring",
	}

	log.Print(fmt.Sprintf("imageTag:[%s] ", event.Detail.ImageTags[0]))
	log.Print(fmt.Sprintf("ImageDigest:[%s] ", event.Detail.ImageDigest))
	log.Print(fmt.Sprintf("RepositoryName:[%s] ", event.Detail.RepositoryName))
	//log.Print(fmt.Sprintf("Account:[%s] ", event.Account))

	c := event.Detail.FindingSeverityCounts.Critical
	h := event.Detail.FindingSeverityCounts.High
	m := event.Detail.FindingSeverityCounts.Medium

	log.Print(fmt.Sprintf("Critical:[%d] ", c))
	log.Print(fmt.Sprintf("High:[%d] ", h))
	log.Print(fmt.Sprintf("Medium:[%d] ", m))

	var color string
	if h == 0 && m == 0 && c == 0 {
		color = "good"
	} else if h > 0 || c > 0 {
		color = "danger"
	} else if c == 0 && h == 0 && m > 0 {
		color = "warning"
	}

	detail := "CRITICAL: " + strconv.FormatInt(c, 10) + "\n" + "HIGH: " + strconv.FormatInt(h, 10) + "\n" + "MEDIUM:" + strconv.FormatInt(m, 10)

	//To send a notification with status (slack attachments)
	sr := SlackJobNotification{
		Color:     color,
		IconEmoji: ":hammer_and_wrench",
		Details:   detail,
		Text:      event.Detail.RepositoryName + ":" + event.Detail.ImageTags[0],
	}

	err := sc.SendJobNotification(sr)
	if err != nil {
		log.Fatal(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "OK",
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
