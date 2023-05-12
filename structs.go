package minstrel

type Minstrel struct {
	config Config
	inst   *instance
}

type Config struct {
	Token       string
	ProjectID   string
	ApiEndpoint string
	ModelID     string
}

type Message struct {
	Author  string `json:"author"`
	Content string `json:"content"`
}

type request struct {
	Instances  []instance `json:"instances"`
	Parameters Parameters `json:"parameters"`
}

type instance struct {
	Context  string    `json:"context"`
	Examples []Message `json:"examples"`
	Messages []Message `json:"messages"`
}

type Parameters struct {
	Temperature     float64 `json:"temperature"`
	MaxOutputTokens int     `json:"tokenLimits"`
	TopP            float64 `json:"topP"`
	TopK            float64 `json:"topK"`
}

type Response struct {
	Predictions []struct {
		SafetyAttributes struct {
			Blocked    bool          `json:"blocked"`
			Scores     []interface{} `json:"scores"`
			Categories []interface{} `json:"categories"`
		} `json:"safetyAttributes"`
		Candidates []Message `json:"candidates"`
	} `json:"predictions"`
}
