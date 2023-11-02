package globals

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

// the openapi key
var OpenaiKey string

// the path of the environment variable
var EnvFile string

// the flag struct
type Flag struct {
	Help  string
	Names []string
}

// enum types
type Model struct {
	Chat bool
	Name string
}

// enums for the models
var (
	Turbo   Model = Model{Name: "gpt-3.5-turbo", Chat: true}
	Davinci Model = Model{Name: "text-davinci-003", Chat: false}
	Curie   Model = Model{Name: "text-curie-001", Chat: false}
	Babbage Model = Model{Name: "text-babbage-001", Chat: false}
	Ada     Model = Model{Name: "text-ada-001", Chat: false}
	GPT4    Model = Model{Name: "gpt-4", Chat: true}
)

var Models = map[string]Model{"turbo": Turbo, "davinci": Davinci, "curie": Curie, "babbage": Babbage, "ada": Ada, "gpt4": GPT4}

// variables for the minimum and maximum ranges
var (
	TempRangeMin float64 = 0
	TempRangeMax float64 = 1
	TopPMin      float64 = 0
	TopPMax      float64 = 1
	PresenceMin  float64 = -2
	PresenceMax  float64 = 2
	FrequenceMin float64 = -2
	FrequenceMax float64 = 2
	BestOfMin    int     = 1
	BestOfMax    int     = 20
)

// the errors to use for wrong formats
var (
	ErrWrongModel          error = errors.New("The model you entered was not correct")
	ErrWrongTempRange      error = fmt.Errorf("The temperature is not in the correct range (%v <= temp <= %v)", TempRangeMin, TempRangeMax)
	ErrWrongToppRange      error = fmt.Errorf("The top_p is not in the correct range (%v <= top_p <= %v)", TopPMin, TopPMax)
	ErrWrongPresenceRange  error = fmt.Errorf("The presence penalty is not in the correct range (%v <= presence <= %v)", PresenceMin, PresenceMax)
	ErrWrongFrequenceRange error = fmt.Errorf("The presence penalty is not in the correct range (%v <= frequence <= %v)", FrequenceMin, FrequenceMax)
	ErrWrongBestOfRange    error = fmt.Errorf("The best of variable is not in the correct range (%v <= best of  <= %v)", BestOfMin, BestOfMax)
	ErrWrongKey            error = errors.New("The API Kry you entered is either wrong or hasn't been set up with a paid account of GPT. You must sign up for a paid account at Openai GPT.")
)

// all the help messages
const (
	inputHelp       string = "The input (git diff file.txt)"
	userStoryHelp   string = "The user story with acceptance criteria"
	jsonHelp        string = "If the output is in JSON"
	verboseHelp     string = "If the output is verbose"
	modelHelp       string = "The model for GPT (see USAGE.md for more details)"
	maxTokensHelp   string = "The length of the max tokens (see USAGE.md for more details)"
	temperatureHelp string = "What sampling temperature to use, between 0 and 2. Higher values like 0.8 will make the output more random, while lower values like 0.2 will make it more focused and deterministic."
	topPHelp        string = "An alternative to sampling with temperature, called nucleus sampling, where the model considers the results of the tokens with top_p probability mass. So 0.1 means only the tokens comprising the top 10% probability mass are considered."
	frequenceHelp   string = "Number between -2.0 and 2.0. Positive values penalize new tokens based on their existing frequency in the text so far, decreasing the model's likelihood to repeat the same line verbatim."
	presenceHelp    string = "Number between -2.0 and 2.0. Positive values penalize new tokens based on whether they appear in the text so far, increasing the model's likelihood to talk about new topics."
	bestOfHelp      string = "Generates best_of completions server-side and returns the 'best' (the one with the highest log probability per token). Results cannot be streamed."
)

// the flag arrays
var (
	inputFlagNames       []string = []string{"input", "i"}
	userStoryFlagNames   []string = []string{"userStory", "us"}
	jsonFlagNames        []string = []string{"json", "j"}
	verboseFlagNames     []string = []string{"verbose", "v"}
	modelFlagNames       []string = []string{"model", "m"}
	maxTokensFlagNames   []string = []string{"max"}
	temperatureFlagNames []string = []string{"temp", "t"}
	toppFlagNames        []string = []string{"topp"}
	frequenceFlagNames   []string = []string{"frequence", "freq", "fr", "f"}
	presenceFlagNames    []string = []string{"pr", "presence", "p", "pres"}
	bestOfFlagNames      []string = []string{"bo", "bestof", "best"}
)

// the flags themselves
var (
	VerboseFlag = Flag{
		Help:  verboseHelp,
		Names: verboseFlagNames,
	}
	JsonFlag = Flag{
		Help:  jsonHelp,
		Names: jsonFlagNames,
	}
	InputFlag = Flag{
		Help:  inputHelp,
		Names: inputFlagNames,
	}
	UserStoryFlag = Flag{
		Help:  userStoryHelp,
		Names: userStoryFlagNames,
	}
	ModelFlag = Flag{
		Help:  modelHelp,
		Names: modelFlagNames,
	}
	MaxTokenFlag = Flag{
		Help:  maxTokensHelp,
		Names: maxTokensFlagNames,
	}
	TemperatureFlag = Flag{
		Help:  temperatureHelp,
		Names: temperatureFlagNames,
	}
	ToppFlag = Flag{
		Help:  topPHelp,
		Names: toppFlagNames,
	}
	FrequenceFlag = Flag{
		Help:  frequenceHelp,
		Names: frequenceFlagNames,
	}
	PresenceFlag = Flag{
		Help:  presenceHelp,
		Names: presenceFlagNames,
	}
	BestOfFlag = Flag{
		Help:  bestOfHelp,
		Names: bestOfFlagNames,
	}
)
var Log zerolog.Logger
var Verbose bool

func Setup(json bool, verbose bool) {
	Verbose = verbose
	if !json {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	Log = log.Logger
	// set the logger
	// set the environment file
	EnvFile = ".rgpt.env"
	home := os.Getenv("HOME")
	// load the environment file
	err := godotenv.Load(fmt.Sprintf("%v/%v", home, EnvFile))
	if err != nil {
		// if the error says the environement file doesn't exist
		if strings.Contains(err.Error(), "no such file") {
			Log.Error().
				Str("Env File", EnvFile).
				Str("Home env var", home).
				Err(err).
				Msg("Env file not found. Did you follow the instructions in the INSTALLATION.md?")
		}
		Log.Error().
			Err(err).
			Msg("Error while loading environment variable")
	}
	// set the openapi key to the environment variable
	OpenaiKey = os.Getenv("OPENAI_KEY")
	if len(OpenaiKey) == 0 {
		Log.Error().
			Str("Env file", EnvFile).
			Msg("Open Ai API Key is empty")
	}
}
