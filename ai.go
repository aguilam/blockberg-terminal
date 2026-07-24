package main

import (
	"context"
	"encoding/json"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)
const defaultPrompt = `
You are a text recognition system that extracts Minecraft item offers from input text and returns a JSON object.
- Keep item names in the same language as the input. Do not translate item names.

3.
Example 1:
Input:
Block of quartz 1 stack = 1 diamond Raven_

Output:
{
  "seller_name": "Raven_",
  "item_name": "Block of quartz",
  "quantity": 64,
  "price": 1,
  "type_name": "building blocks"
}

Example 2:
Input:
= Gravel = 64 items - 1 diamond FrostByte =-=-=-=-=-

Output:
{
  "seller_name": "FrostByte",
  "item_name": "Gravel",
  "quantity": 64,
  "price": 1,
  "type_name": "natural blocks"
}


Example 3:
Input:
DEEPSLATE 2 stacks 1 diamond

Output:
{
  "seller_name": "None",
  "item_name": "Deepslate",
  "quantity": 128,
  "price": 1,
  "type_name": "natural blocks"
}


Example 4:
Input:
Eye armor trim smithing template 20 diamonds PixelFox

Output:
{
  "seller_name": "PixelFox",
  "item_name": "Eye armor trim smithing template",
  "quantity": 1,
  "price": 20,
  "type_name": "other"
}


Example 5:
Input:
STONE 2 stacks - 1 diamond block NightWolf

Output:
{
  "seller_name": "NightWolf",
  "item_name": "Stone",
  "quantity": 128,
  "price": 9,
  "type_name": "building blocks"
}


Example 6:
Input:
😀Stone bricks😀😊 3 stacks and 128 items - 3 diamond blocks ShadowCat

Output:
{
  "seller_name": "ShadowCat",
  "item_name": "Stone bricks",
  "quantity": 320,
  "price": 27,
  "type_name": "building blocks"
}


Example 7:
Input:
Golden 🌶CARROT🌶 32 items - 2 diamonds BlueNova

Output:
{
  "seller_name": "BlueNova",
  "item_name": "Golden carrot",
  "quantity": 32,
  "price": 2,
  "type_name": "food & drinks"
}


Example 8:
Input:
Diamond pickaxe with Silk Touch 1 item - 10 diamonds IronMage

Output:
{
  "seller_name": "IronMage",
  "item_name": "Diamond pickaxe with Silk Touch",
  "quantity": 1,
  "price": 10,
  "type_name": "tools & utilities"
}


Example 9:
Input:
Mending book 1 item - 9 diamonds CrystalFox

Output:
{
  "seller_name": "CrystalFox",
  "item_name": "Mending",
  "quantity": 1,
  "price": 9,
  "type_name": "ingredients"
}


Example 10:
Input:
Book 1984 1 item - 5 diamonds VoidWalker

Output:
{
  "seller_name": "VoidWalker",
  "item_name": "Book 1984",
  "quantity": 1,
  "price": 5,
  "type_name": "ingredients"
}`
func GenerateSchema[T any]() map[string]any {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)

	data, _ := json.Marshal(schema)
	var result map[string]any
	json.Unmarshal(data, &result)
	return result
}
var RecognizedBarrelSchema = GenerateSchema[AiRecognizedBarrelItem]()
func aiBarrelRecognition(apiKey string, baseUrl string, model string, message string, prompt string) (*RecognizedBarrelItem, error) {
	ctx := context.Background()
	client := openai.NewClient(option.WithBaseURL(baseUrl),option.WithAPIKey(apiKey))
	resp, err := client.Responses.New(ctx,responses.ResponseNewParams{
			Instructions: openai.String(prompt),
			Input: responses.ResponseNewParamsInputUnion{OfString: openai.String(message)},
			Model: model,
			Text: responses.ResponseTextConfigParam{
				Format: responses.ResponseFormatTextConfigParamOfJSONSchema(
					"barrel_item",
					RecognizedBarrelSchema,
				),
			},
		},
	)
	if err != nil {
		return nil, err
	}
	var recognizedBarrel AiRecognizedBarrelItem
	err = json.Unmarshal([]byte(resp.OutputText()),&recognizedBarrel)
	if err != nil {
		return nil,err
	}
	return &RecognizedBarrelItem{
		SellerName: recognizedBarrel.SellerName,
		ItemName: recognizedBarrel.ItemName,
		BenefitRatio: float64(recognizedBarrel.Quantity) / float64(recognizedBarrel.Price),
		Quantity: float64(recognizedBarrel.Quantity),
		Price: float64(recognizedBarrel.Price),
		TypeName: recognizedBarrel.TypeName,
	},nil
}