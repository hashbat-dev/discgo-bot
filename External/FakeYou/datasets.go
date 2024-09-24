package fakeyou

const (
	FAKEYOU_SLEEP_BETWEEN_MODELS_MS = 50
)

var (
	FakeYouURLModelList    string = "https://api.fakeyou.com/tts/list"
	FakeYouURLRequestTTS   string = "https://api.fakeyou.com/tts/inference"
	FakeYouURLCheckRequest string = "https://api.fakeyou.com/tts/job/"
	FakeYouTTSAudioBaseUrl string = "https://storage.googleapis.com/vocodes-public"
)
