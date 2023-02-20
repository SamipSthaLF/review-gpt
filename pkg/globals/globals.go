package globals

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)
// the openapi key
var OpenaiKey string
// the path of the environment variable
var EnvFile string
// the flag struct
type Flag struct{
    Help string
    Names []string 
}
// all the help messages
const (
    inputHelp string = "The input (git diff file.txt)"
    modelHelp string= "The model for GPT (see USAGE.md for more details)"
    maxTokensHelp string= "The length of the max tokens (see USAGE.md for more details)"
    temperatureHelp string= "What sampling temperature to use, between 0 and 2. Higher values like 0.8 will make the output more random, while lower values like 0.2 will make it more focused and deterministic."
    topPHelp string= "An alternative to sampling with temperature, called nucleus sampling, where the model considers the results of the tokens with top_p probability mass. So 0.1 means only the tokens comprising the top 10% probability mass are considered."
    frequenceHelp string= "Number between -2.0 and 2.0. Positive values penalize new tokens based on their existing frequency in the text so far, decreasing the model's likelihood to repeat the same line verbatim."
    presenceHelp string= "Number between -2.0 and 2.0. Positive values penalize new tokens based on whether they appear in the text so far, increasing the model's likelihood to talk about new topics." 
    bestOfHelp string= "Generates best_of completions server-side and returns the 'best' (the one with the highest log probability per token). Results cannot be streamed."
)
// the flag arrays
var (
    inputFlagNames []string = []string{"input", "i"}
    modelFlagNames []string = []string{"model", "m"}
    maxTokensFlagNames []string = []string{"max"}
    temperatureFlagNames []string = []string{"temp", "t"}
    toppFlagNames []string = []string{"topp"}
    frequenceFlagNames []string= []string{"frequence", "freq", "fr", "f"}
    presenceFlagNames []string = []string{"pr", "presence", "p", "pres"}
    bestOfFlagNames []string= []string{"bo", "bestof", "best"}
)
// the flags themselves 
var (
    InputFlag = Flag{
        Help: inputHelp,
        Names: inputFlagNames,  
    }
    ModelFlag = Flag{
        Help: modelHelp,
        Names: modelFlagNames,  
    }
    MaxTokenFlag = Flag{
        Help: maxTokensHelp,
        Names: maxTokensFlagNames,
    }
    TemperatureFlag = Flag{
        Help: temperatureHelp,
        Names: temperatureFlagNames,
    }
    ToppFlag= Flag{
        Help: topPHelp,
        Names: toppFlagNames,
    }
    FrequenceFlag = Flag{
        Help: frequenceHelp,
        Names: frequenceFlagNames,
    }
    PresenceFlag = Flag{
        Help: presenceHelp,
        Names: presenceFlagNames,
    }
    BestOfFlag = Flag{
        Help: bestOfHelp,
        Names: bestOfFlagNames,
    }
)
func Setup() {
    // set the environment file
	EnvFile = ".rgpt.env"
	home := os.Getenv("HOME")
    // load the environment file
	err := godotenv.Load(fmt.Sprintf("%v/%v", home, EnvFile))
	if err != nil {
        // if the error says the environement file doesn't exist
		if strings.Contains(err.Error(), "no such file") {
			log.Fatalln(errors.New(".rgpt.env not found. Did you follow the instructions in the INSTALLATION.md?"))
		}
		log.Fatalln(err)
	}
    // set the openapi key to the environment variable
	OpenaiKey = os.Getenv("OPENAI_KEY")
}
