# stockbot
A Slack bot that displays funny stock photos based on your search terms

## Testing on localhost

Using httpie

```
http -f POST http://localhost:3000/bot/stockimage token=debug team_id=T0001 team_domain=example channel_id=1234 channel_name=random user_id=4321 user_name=user text="stockphoto: salad" trigger_word="stockphoto:" timestamp=1234
```
