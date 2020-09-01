# Приложение Покерный банк
Подключение к рабочему экземпляру сервера postgres возможно из авторизованной консоли google cloud console командой:
```shell script
gcloud sql connect postgres-instanse --user=postgres --quiet
````
Команда публикации приложения в Google App Engine, параметр --project не нужен, если проект уже выбран:
````
gcloud app deploy --project [YOUR_PROJECT_ID]
````