# WB-Product-Info-Fetcher
Консольная утилита, которая по артикулу товара Wildberries (например, 250263101) получает его характеристики (тело ответа JSON, название, цена, цена со скидкой, рейтинг) через WB API и выводит в консоль.
Программа принимает на вход два аргумента — артикул товара на Wildberries и одно из полей: body, name, price, salePrice, rating. По артикулу выполняется поиск товара, а поле служит для сокращения вывода (например, поле price способствует тому, что будет возвращена только цена товара). Второй аргумент (поле) не обязателен: при его отсутствии выводятся все поля, кроме body.

**Установка**:\
[user@nixos:~]$ git clone https://github.com/asquebay/WB-Product-Info-Fetcher.git

[user@nixos:~]$ cd WB-Product-Info-Fetcher

**Пример использования 1**:\
[user@nixos:~]$ go run main.go 250263101

**Вывод 1**:\
Название: Кружки для чая и кофе стеклянные 2 шт "Глория" 320 мл\
Цена: 1270.00 ₽\
Акционная цена: 496.00 ₽\
Рейтинг: 5.0 ★

**Пример использования 2**:\
[user@nixos:~]$ go run main.go 250263101 price

**Вывод 2**:\
1270.00
