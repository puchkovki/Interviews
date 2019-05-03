# Приложение-singleton echoloop

Данная задача была предложена в виде вступительного испытания на Межфакультетскую кафедру теоретической и прикладной информатики МФТИ.

## Условие задачи

Напишите простое приложение-singleton ```echoloop```. Само приложение выводит раз в секунду переданные ему параметры.

```> echoloop "test text"```<br>
```test text```<br>
```test text```<br>
```test text```<br>
```...```<br>

При попытке повторного запуска новое приложение должно обнаруживать, что экземпляр уже запущен, отравлять свои параметры в уже запущенное приложение и завершаться, сообщив об этом. Запущенное приложение после этого должно печатать и свои, и полученные параметры.

```> echoloop "test1" &```<br>
```test1```<br>
```test1```<br>
```> echoloop "test2"```<br>
```echoloop for "test2" finished!```<br>
```test1```<br>
```test2```<br>
```test1```<br>
```test2```<br>
```...```<br>

 

## Как начать

Данная программа была написана на языке программирования Golang. Для поддержки Golang терминалом следует скачать бинарные файлы по ссылке:
https://golang.org/dl/

## Идея решения

Изначально был создан mutex файл, который проверяет уникальность процесса.

Используя функцию FcntlFlock, устанавливаем блокировку на запись для mutex файла. Если исполняемый файл был запущен не первым, то его аргументы передаются в fifo файл, сам он завершается.

Создаем и открываем глобальный fifo файл.

Приложение отлавливает сигналы os.Kill и os.Interrupt. В случае получения сигнала приложение закрывает и удаляет все ненужные файлы и завершает работу.

Печать аргументов происходит в отдельной goroutine, связь с которой обеспечивается с помощью канала.

Основная goroutine находится в ожидании поступления аргументов из глобального fifo файла.

## Автор

* **Puchkov Kyryll** — [puchkovki](https://github.com/puchkovki)