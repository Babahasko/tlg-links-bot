package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	getUpdatesMethod = "getUpdates"
	sendMessageMethod = "sendMessage"	
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host string, token string) *Client {
	return &Client{
		host: host,
		basePath: newBasePath(token),
		client: http.Client{},
	}
}

func newBasePath(token string) (string){
	return "bot"+token
}


func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)
	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return fmt.Errorf("message fail: %w", err)
	}
	return nil
}

// Updates. Get the updates by telegramm api. Offset and limit are the parametes
// for selecting, which updates do we need. https://core.telegram.org/bots/api#getupdates
func (c *Client) Updates(offset int, limit int) ([]Update, error){
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, fmt.Errorf("updates method fail: %w", err)
	}
	var res UpdatesResponse

	if err:=json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("unmarshal fail: %w", err)
	}

	return res.Result, nil
}

// doRequest make a request with given method on defiened query.
// Very common method for making querys and returns body or error
func (c *Client) doRequest(method string, query url.Values) ([]byte, error){
	u := url.URL{
		Scheme: "https",
		Host: c.host,
		Path: path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("doRequst fail: %w", err)
	}

	req.URL.RawQuery = query.Encode()
	
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("doRequst fail: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("doRequest fail: %w", err)
	}

	return body, nil
}
