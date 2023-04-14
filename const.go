package openai

const (
	openaiUrlv1         = "https://api.openai.com/v1"
	chatCompletion      = "chat/completions"
	completions         = "completions"
	edits               = "edits"
	audioTranscriptions = "audio/transcriptions"
	audioTranslations   = "audio/translations"
	fineTunes           = "fine-tunes"
	embeddings          = "embeddings"
	moderations         = "moderations"
	models              = "models"
)

// model list from https://platform.openai.com/docs/models/model-endpoint-compatibility
const (
	GPT4        = "gpt-4"
	GPT40314    = "gpt-4-0314"
	GPT432K     = "gpt-4-32k"
	GPT432K0314 = "gpt-4-32k-0314"

	GPT3Dot5Turbo          = "gpt-3.5-turbo"
	GPT3Dot5Turbo0301      = "gpt-3.5-turbo-0301"
	GPT3Dot5TextDavinci003 = "text-davinci-003"
	GPT3Dot5TextDavinci002 = "text-davinci-002"
	GPT3Dot5CodeDavinci002 = "code-davinci-002"

	GPT3TextCurie001       = "text-curie-001"
	GPT3TextBabbage001     = "text-babbage-001"
	GPT3TextAda001         = "text-ada-001"
	GPT3TextDavincEdit001  = "text-davinci-edit-001"
	GPT3CodeDavinciEdit001 = "code-davinci-edit-001"
	GPT3Whisper1           = "whisper-1"
	GPT3Davinci            = "davinci"
	GPT3Curie              = "curie"
	GPT3Babbage            = "babbage"
	GPT3Ada                = "ada"
	TextEmbeddingAda002    = "text-embedding-ada-002"
	TextSearchAdaDoc001    = "text-search-ada-doc-001"
	TextModerationStable   = "text-moderation-stable"
	TextModerationLatest   = "text-moderation-latest"
)

// chat role defined by OpenAI
const (
	ChatMessageRoleUser      = "user"
	ChatMessageRoleSystem    = "system"
	ChatMessageRoleAssistant = "assistant"
)
