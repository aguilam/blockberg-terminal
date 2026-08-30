# Blockberg Terminal

[Client](https://github.com/aguilam/blockberg-terminal-mod) | Server(This)

This project realize GO server for collecting information about the trading zone on the server (e.g. SMP) and subsequently searching best offers.

## Features

- **Authentication:** Set an API key if you want to allow adding new barrels by only users with this api key, or leave it unset to allow all users to post new barrels.
- **Search by sign text:** Search barrels with best offers by parsed sign text
- **Search by saved items:** Search need you items by last snapshoted items in barrel
- **Get barrel content:** Get actual state of items in barrel by snapshot
- **AI Sign Parsing:** If you add ai-key, ai-url and ai-prompt all barrel signs not contain formatted data will be parsed by selected LLM model
- **Zone Boundaries:** Limit zone for save new barrels by min and max coord in flags

### Example usage with argument
```
./blockberg-terminal --port 8081 
```

### Arguments

| Argument    | Base value | Description                                                            |
| ------------| ---------- | -----------                                                            |
| --api-key   | ""         | Key used for add new storages. If not set, anyone can add new storages |
| --port      | 8080       | Server port                                                            |
| --ai-key    | ""         | AI provider API key                                                    |
| --ai-url    | ""         | AI provider base url                                                   |
| --ai-prompt | ""         | Replace base prompt for ai recognition                                 |
| --ai-model  | ""         | Selected AI model name                                                 |
| --x-min     | UnsetCoord | Set minimum X coord for get new data by barrel post endpoint           |
| --x-max     | UnsetCoord | Set maximum X coord for get new data by barrel post endpoint           |
| --y-min     | UnsetCoord | Set minimum Y coord for get new data by barrel post endpoint           |
| --y-max     | UnsetCoord | Set maximum Y coord for get new data by barrel post endpoint           |
| --z-min     | UnsetCoord | Set minimum Z coord for get new data by barrel post endpoint           |
| --z-max     | UnsetCoord | Set maximum Z coord for get new data by barrel post endpoint           |