Project allows to use API from zakaz.UA (delivery_schedule/plan) and send message to Telegram if slot is available.

Config.json should be placed in root folder.

```  
 BotAPIKey  string `mapstructure:"bot_api_key"`
 TargetURL  string `mapstructure:"target_url"`
 ChatID     int64  `mapstructure:"chat_id"`
 MessageURL string `mapstructure:"message_url"
 ```