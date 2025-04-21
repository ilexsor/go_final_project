# Файлы для итогового задания

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.

# Run taskManager
1. docker build -t taskManager:v1 .
2. docker run -d -p 7540:7540 --name taskManager taskManager:v1