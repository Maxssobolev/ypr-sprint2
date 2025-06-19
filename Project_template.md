# Задание 1

#### To-Be:

- Домен "Онлайн-кинотеатр - агрегатор":

	- Поддомен "Управление пользователями"

		- Контекст: регистрация и аутентификация пользователей в системе, восстановление доступа
		- Контекст: авторизация пользователей

	- Поддомен "Управление контентом"

		- Контекст: интеграция с кинотеатрами
		- Контекст: каталог кинотеатров-партнеров

  - Поддомен "Метаданные о фильмах"

		- Контекст: хранение детальной информации о фильмах
		- Контекст: информация о жанрах, актерах, режиссерах
		- Контекст: управление медиа-контентом

	- Поддомен "Избранное и оценки"
	
		- Контекст: оценка фильмов
		- Контекст: управление избранным
		- Контекст: история просмотров

	- Поддомен "Управление подписками и платежами"
		
		- Контекст: интеграция с платежными системами
		- Контекст: просмотр и управление подписками

	- Поддомен "Поддержка и сопровождение"

		- Контекст: система тикетов, взаимодействие с пользователями экосистемы, помощь в подключении устройств
		- Контекст: получение обратной связи от пользователей экосистемы

	- Поддомен "Аналитика для маркетинга"

		- Контекст: сбор статистических данных (просмотры, клики, избранное)
		- Контекст: визуализация собранных данных с целью выявления наиболее популярных кинотеатров и контента среди пользователей


  
[Диаграмма контейнеров](./docs/architecture/container/container.puml)
![Диаграмма контейнеров](./docs-out/docs/architecture/container/container/container.png)

# Задание 2

### 1. Proxy

Реализовано.


### 2. Kafka

Тесты 
![Тесты](./screenshots/tests.png)
Топики
![Топики](./screenshots/topics.png)

# Задание 3

Вывод фильмов
![Вывод фильмов](./screenshots/api-movies.png)
Тесты 
![Тесты](./screenshots/k8s-tests.png)


# Задание 4
Для простоты дальнейшего обновления и развертывания вам как архитектуру необходимо так же реализовать helm-чарты для прокси-сервиса и проверить работу 

Для этого:
1. Перейдите в директорию helm и отредактируйте файл values.yaml

```yaml
# Proxy service configuration
proxyService:
  enabled: true
  image:
    repository: ghcr.io/db-exp/cinemaabysstest/proxy-service
    tag: latest
    pullPolicy: Always
  replicas: 1
  resources:
    limits:
      cpu: 300m
      memory: 256Mi
    requests:
      cpu: 100m
      memory: 128Mi
  service:
    port: 80
    targetPort: 8000
    type: ClusterIP
```

- Вместо ghcr.io/db-exp/cinemaabysstest/proxy-service напишите свой путь до образа для всех сервисов
- для imagePullSecret проставьте свое значение (скопируйте из конфигурации kubernetes)
  ```yaml
  imagePullSecrets:
      dockerconfigjson: ewoJImF1dGhzIjogewoJCSJnaGNyLmlvIjogewoJCQkiYXV0aCI6ICJaR0l0Wlhod09tZG9jRjl2UTJocVZIa3dhMWhKVDIxWmFVZHJOV2hRUW10aFVXbFZSbTVaTjJRMFNYUjRZMWM9IgoJCX0KCX0sCgkiY3JlZHNTdG9yZSI6ICJkZXNrdG9wIiwKCSJjdXJyZW50Q29udGV4dCI6ICJkZXNrdG9wLWxpbnV4IiwKCSJwbHVnaW5zIjogewoJCSIteC1jbGktaGludHMiOiB7CgkJCSJlbmFibGVkIjogInRydWUiCgkJfQoJfSwKCSJmZWF0dXJlcyI6IHsKCQkiaG9va3MiOiAidHJ1ZSIKCX0KfQ==
  ```

2. В папке ./templates/services заполните шаблоны для proxy-service.yaml и events-service.yaml (опирайтесь на свою kubernetes конфигурацию - смысл helm'а сделать шаблоны для быстрого обновления и установки)

```yaml
template:
    metadata:
      labels:
        app: proxy-service
    spec:
      containers:
       Тут ваша конфигурация
```

3. Проверьте установку
Сначала удалим установку руками

```bash
kubectl delete all --all -n cinemaabyss
kubectl delete  namespace cinemaabyss
```
Запустите 
```bash
helm install cinemaabyss .\src\kubernetes\helm --namespace cinemaabyss --create-namespace
```
Если в процессе будет ошибка
```code
[2025-04-08 21:43:38,780] ERROR Fatal error during KafkaServer startup. Prepare to shutdown (kafka.server.KafkaServer)
kafka.common.InconsistentClusterIdException: The Cluster ID OkOjGPrdRimp8nkFohYkCw doesn't match stored clusterId Some(sbkcoiSiQV2h_mQpwy05zQ) in meta.properties. The broker is trying to join the wrong cluster. Configured zookeeper.connect may be wrong.
```

Проверьте развертывание:
```bash
kubectl get pods -n cinemaabyss
minikube tunnel
```

Потом вызовите 
https://cinemaabyss.example.com/api/movies
и приложите скриншот развертывания helm и вывода https://cinemaabyss.example.com/api/movies

## Удаляем все

```bash
kubectl delete all --all -n cinemaabyss
kubectl delete namespace cinemaabyss
```
