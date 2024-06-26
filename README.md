# ScriptService

## ADR:
   Приложение, которое предоставляет REST API для запуска команд. Команда представляет собой bash-скрипт. 
   Приложение позволяет параллельно запускать произвольное количество команд.

## Задача
   Разработать и описать решение для выполнения, управления и поддержки команд.

## Решение
1. Создание новой команды

   При создании новой команды выполняются следующие шаги:

   Команда сохраняется в базу данных для дальнейшего отслеживания.
   Генерируется уникальный идентификатор команды и возвращается пользователю.
   Запускается выполнение команды асинхронно с использованием бэкграунд процесса.
2. Получение списка команд

   Пользователь может получить список всех сохраненных команд из базы данных. Это позволяет пользователям просматривать доступные команды.

3. Получение одной команды по идентификатору

   По запросу идентификатора команды пользователь может получить конкретную команды её результат и статус.

4. Остановка команды

   Для остановки команды используется следующий механизм:

   Используется sync.Map, чтобы связать каждую команду с ее уникальным идентификатором и контекстом выполнения.  
   При запросе на остановку команды по идентификатору, соответствующий контекст извлекается из sync.Map.  
   Вызывается метод отмены контекста, что приводит к остановке выполнения команды.  
5. Поддержка долгих команд  
   Использована функция StdoutPipe() на объекте типа exec.Cmd (команда выполнения).
   С попмощью неё создаётся канал, через который можно считывать вывод стандартного потока вывода (stdout) команды, запущенной с помощью exec.Cmd.

## Результаты
   Эти решения обеспечивают эффективное управление командами, поддерживая их асинхронное выполнение, отслеживание и остановку, 
   а также поддержку долгих команд с сохранением вывода.


## Инструкция по запуску приложения
1. Создайте контейнер и зпустите приложение в Docker 
```ssh
docker-compose up --build script_service
```

2. Дальше вы можете отправлять запросы локально используя порт 8000
   по урлам указанным в [app](./internal/app/app.go)