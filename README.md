# gokeystat

[![Join the chat at https://gitter.im/AlekseyLobanov/gokeystat](https://badges.gitter.im/AlekseyLobanov/gokeystat.svg)](https://gitter.im/AlekseyLobanov/gokeystat?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Travis status](https://travis-ci.org/AlekseyLobanov/gokeystat.svg)](https://travis-ci.org/AlekseyLobanov/gokeystat)
[![Go Report Card](https://goreportcard.com/badge/github.com/alekseylobanov/gokeystat)](https://goreportcard.com/report/github.com/alekseylobanov/gokeystat)

**gokeystat** позволяет собирать статистику использования клавиатуры, подсчитывая и сохраняя каждые 15 минут количество нажатий каждой клавиши.

Для запуска необходимо выполнить:

`./gokeystat -id <keyboard_id>`

Где `keyboard_id` можно получить командой:

`xinput list`

Собранные данные можно экспортировать в форматы:

* `csv` и `csv.gz`
* `json` и `json.gz`
* `jsl` и `jsl.gz`
 
Для экспорта необходимо запустить с ключом `-o`:

`./gokeystat -o example.csv`

Чтобы вывести информацию по каждой клавише, следует добавить ключ `-full`.