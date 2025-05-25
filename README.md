# Links tlg bot

Это учебный проект, в котором я реализовал телеграмм бота на go. Собственно вот плечи, на которых стоит этот проект -> серия видео [Николая Тузова](https://www.youtube.com/watch?v=PnOrFYtZJUI&list=PLFAQFisfyqlWDwouVTUztKX2wUjYQ4T3l&index=1)

# Основная идея
Создать простенького телеграмм бота, который позволяет сохранять ссылки и выдаёт рандомную.

# Запуск проекта

```cmd
git clone https://github.com/Babahasko/tlg-links-bot.git .
cd tlg-links-bot
go build
./tlg-links-bot -tg-bot-token="<token>"
```
# Что бы я изменил(возможно когда-нибудь)
1. Хранить данные в БД sqlite
2. Хранить конфигурацию проекта в .env
