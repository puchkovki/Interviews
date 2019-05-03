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
> echoloop "test1" &
test1
test1
> echoloop "test2"
echoloop for "test2" finished!
test1
test2
test1
test2
...

 

## Язык программирования Golang

Данная программа была написана на языке программирования Golang. Для поддержки Golang терминалом следует скачать бинарные файлы по ссылке:
https://golang.org/dl/

## Идея решения

Explain how to run the automated tests for this system

### Break down into end to end tests

Explain what these tests test and why

```
Give an example
```

### And coding style tests

Explain what these tests test and why

```
Give an example
```

## Deployment

Add additional notes about how to deploy this on a live system

## Built With

* [Dropwizard](http://www.dropwizard.io/1.0.2/docs/) - The web framework used
* [Maven](https://maven.apache.org/) - Dependency Management
* [ROME](https://rometools.github.io/rome/) - Used to generate RSS Feeds

## Contributing

Please read [CONTRIBUTING.md](https://gist.github.com/PurpleBooth/b24679402957c63ec426) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/your/project/tags). 

## Authors

* **Billie Thompson** - *Initial work* - [PurpleBooth](https://github.com/PurpleBooth)

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* Hat tip to anyone whose code was used
* Inspiration
* etc
