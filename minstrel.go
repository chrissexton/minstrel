package minstrel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

func defaultConfig() Config {
	return Config{
		Token:     "",
		ProjectID: "",
		// Note that this default was taken from checking the browser example.
		// Google's CURL example code does not function, but this is probably a
		// "bad" value.
		ApiEndpoint: "us-central1-prediction-aiplatform.clients6.google.com",
		ModelID:     "chat-bison@001",
	}
}

func New(token, projectID string) *Minstrel {
	cfg := defaultConfig()
	cfg.Token = token
	cfg.ProjectID = projectID
	return NewWithConfig(cfg)
}

func NewWithConfig(config Config) *Minstrel {
	return &Minstrel{
		config: config,
		inst: &instance{
			Messages: []Message{},
			Examples: []Message{},
		},
	}
}

func (m *Minstrel) SetPrompt(content string) {
	m.inst.Context = content

}
func (m *Minstrel) AddExample(example Message) {
	m.inst.Examples = append(m.inst.Examples, example)
}

func (m *Minstrel) ClearExamples() {
	m.inst.Examples = []Message{}
}

func (m *Minstrel) ClearMessages() {
	m.inst.Messages = []Message{}
}

func (m *Minstrel) Reset() {
	m.ClearExamples()
	m.ClearMessages()
}

func defaultParams() Parameters {
	return Parameters{
		Temperature:     0.2,
		MaxOutputTokens: 256,
		TopP:            0.8,
		TopK:            40,
	}
}

func (m *Minstrel) CompleteCandidates(content string) (Response, error) {
	return m.CompleteCandidatesWithParams(defaultParams(), content)
}

var url = "https://%s/ui/projects/%s/locations/us-central1/publishers/google/models/%s:predict"

func (m *Minstrel) buildURL() string {
	return fmt.Sprintf(url, m.config.ApiEndpoint, m.config.ProjectID, m.config.ModelID)
}

func (m *Minstrel) CompleteCandidatesWithParams(params Parameters, content string) (Response, error) {
	msg := Message{
		Author:  "user",
		Content: content,
	}
	m.inst.Messages = append(m.inst.Messages, msg)
	reqBody, err := json.Marshal(request{
		Instances:  []instance{*m.inst},
		Parameters: params,
	})
	if err != nil {
		return Response{}, err
	}
	buf := bytes.NewBuffer(reqBody)
	req, err := http.NewRequest(http.MethodPost, m.buildURL(), buf)
	if err != nil {
		return Response{}, err
	}
	// If they provide a token, use it. Otherwise, we assume the client is a CGP client
	var bearer = "Bearer " + m.config.Token
	req.Header.Add("Authorization", bearer)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Response{}, err
	}
	if resp.StatusCode >= 400 {
		return Response{}, fmt.Errorf("status %d: %s", resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}
	chatResp := Response{}
	err = json.Unmarshal(body, &chatResp)
	return chatResp, nil
}

func (m *Minstrel) Complete(content string) (string, error) {
	return m.CompleteWithParams(defaultParams(), content)
}

func (m *Minstrel) CompleteWithParams(params Parameters, content string) (string, error) {
	resp, err := m.CompleteCandidatesWithParams(params, content)
	if err != nil {
		return "", err
	}
	body := ""
	if len(resp.Predictions) > 0 && len(resp.Predictions[0].Candidates) > 0 {
		c := rand.Intn(len(resp.Predictions[0].Candidates))
		msg := resp.Predictions[0].Candidates[c]
		body = msg.Content
		m.inst.Messages = append(m.inst.Messages, msg)
	}
	return body, nil
}
