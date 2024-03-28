### Как запустить
Создать в корне проекта директорию .env, скопировать туда файлы из .env.example (сделал так, потому что не секьюрно заливать настоящие конфиги на гит).
Делаем docker-compose up. Для взаимодействия с миграциями и документацией команды есть в Makefile.

Написать приложение, реализующее REST-API условного маркетплейса:
- авторизация пользователей;
- регистрация пользователей;
- размещение нового объявления;
- отображение ленты объявлений.
API должен быть реализован в формате REST + JSON. Реализовать приложение можно на языке PHP или Go на выбор. Допускается использование фреймворков. База данных любая.

Авторизация пользователя:
- авторизация должна происходить посредством отправки в приложение логина и пароля;
- приложение должно проверить корректность полученных данных и вернуть авторизационный токен в случае успеха;
- токен должен передаваться в хедерах запроса, название хедера на ваш выбор;
- токен может быть передан в любой эндпоинт;
- при наличии токена, он должен быть обработан и проверен.

Регистрация пользователей:
- регистрация должна осуществляться посредством отправки в приложение логина и пароля;
- необходимо предусмотреть и реализовать разумные ограничения на формат логина и пароля;
- в успешном ответе вернуть данные добавленного пользователя.

Размещение объявления:
- размещение объявления должно происходить посредством отправки данными в формате JSON: заголовок, текст объявления, адрес изображения, цена;
- размещать объявления могут только авторизованные пользователи;
- необходимо предусмотреть и реализовать разумные ограничения на длину заголовка, текста объявления, цены, размера и формата изображения;
- в успешном ответе вернуть данные добавленного объявления.

Отображение ленты объявлений:
- лента объявлений представляет из себя список объявлений, отсортированный по дате добавления (самые свежие в начале);
- необходимо реализовать постраничную навигацию, возможность изменения типа и направления сортировки (дата создания и цена), возможность фильтрации по цене (мин. и макс. значение);
- для каждого объявления необходимо вернуть: заголовок, текст объявления, адрес изображения, цену, логин автора;
- для авторизованных пользователей необходимо дополнительно возвращать признак принадлежности объявления текущему пользователю.