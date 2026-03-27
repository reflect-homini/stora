package llm

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/itsLeonB/ungerr"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"github.com/reflect-homini/stora/internal/core/config"
	"github.com/reflect-homini/stora/internal/core/otel"
)

type LLMService interface {
	Prompt(ctx context.Context, p Prompt) (string, error)
}

type Prompt struct {
	SystemMessage  string
	UserMessage    string `validate:"required,min=3"`
	ResponseFormat openai.ChatCompletionNewParamsResponseFormatUnion
}

type openAILLMService struct {
	client   openai.Client
	model    string
	validate *validator.Validate
}

func NewLLMService(cfg config.LLM) LLMService {
	client := openai.NewClient(option.WithAPIKey(cfg.ApiKey), option.WithBaseURL(cfg.BaseUrl))
	return &openAILLMService{client, cfg.Model, validator.New()}
}

func (llm *openAILLMService) Prompt(ctx context.Context, p Prompt) (string, error) {
	ctx, span := otel.Tracer.Start(ctx, "openAILLMService.Prompt")
	defer span.End()

	if err := llm.validate.Struct(p); err != nil {
		return "", ungerr.Wrap(err, "failed to validate prompt input")
	}

	msgs := make([]openai.ChatCompletionMessageParamUnion, 0, 2)

	if p.SystemMessage != "" {
		msgs = append(msgs, openai.SystemMessage(p.SystemMessage))
	}

	msgs = append(msgs, openai.UserMessage(p.UserMessage))

	params := openai.ChatCompletionNewParams{
		Model:          llm.model,
		Messages:       msgs,
		ResponseFormat: p.ResponseFormat,
	}

	response, err := llm.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return "", ungerr.Wrap(err, "error getting LLM response")
	}

	if len(response.Choices) < 1 {
		return "", ungerr.Unknown("no response from LLM")
	}

	return response.Choices[0].Message.Content, nil
}
