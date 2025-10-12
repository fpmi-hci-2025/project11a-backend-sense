# project11a-backend-sense

## Структура базы данных

| Сущность | Таблица БД | Поле | Смысл поля | Тип данных |
|----------|------------|------|------------|------------|
| **Пользователь** | `users` | `id` | уникальный идентификатор пользователя (PK) | UUID |
| | | `username` | логин/отображаемое имя (уникальное) | TEXT |
| | | `icon_url` | ссылка на аватар | TEXT |
| | | `email` | электронная почта (уникальная) | TEXT |
| | | `phone` | номер телефона (уникальный) | TEXT |
| | | `registered_at` | дата/время регистрации | TIMESTAMPTZ |
| | | `description` | описание профиля | TEXT |
| | | `role` | роль пользователя | ENUM user_role |
| | | `statistic` | произвольные метрики профиля | JSONB |
| | | `followers_count` | количество подписчиков | INTEGER |
| | | `following_count` | количество подписок | INTEGER |
| **Публикация** | `publications` | `id` | уникальный идентификатор публикации (PK) | UUID |
| | | `author_id` | автор публикации (FK → users.id) | UUID |
| | | `type` | тип публикации | ENUM publication_type |
| | | `content` | текст/контент | TEXT |
| | | `source` | источник (для цитаты) | TEXT |
| | | `publication_date` | дата/время публикации | TIMESTAMPTZ |
| | | `visibility` | видимость публикации | ENUM visibility_type |
| | | `likes_count` | счетчик лайков (агрегат) | INTEGER |
| | | `comments_count` | счетчик комментариев (агрегат) | INTEGER |
| | | `saved_count` | счетчик сохранений (агрегат) | INTEGER |
| **Медиафайл** | `media_assets` | `id` | уникальный идентификатор медиа (PK) | UUID |
| | | `owner_id` | владелец файла (FK → users.id) | UUID |
| | | `url` | ссылка на файл | TEXT |
| | | `mime` | MIME-тип | TEXT |
| | | `width` | ширина пикселей (может быть NULL) | INTEGER |
| | | `height` | высота пикселей (может быть NULL) | INTEGER |
| | | `exif` | метаданные EXIF | JSONB |
| | | `created_at` | когда загружено | TIMESTAMPTZ |
| **Связь публикации–медиа** | `publication_media` | `publication_id` | публикация (FK → publications.id) | UUID |
| | | `media_id` | медиафайл (FK → media_assets.id) | UUID |
| | | `ord` | порядок медиа внутри публикации | INTEGER |
| **Комментарий** | `comments` | `id` | уникальный идентификатор комментария (PK) | UUID |
| | | `publication_id` | ссылка на публикацию (FK) | UUID |
| | | `parent_id` | родительский комментарий (для вложенности, FK, NULL допускается) | UUID |
| | | `author_id` | автор комментария (FK → users.id) | UUID |
| | | `text` | текст комментария | TEXT |
| | | `created_at` | дата/время создания | TIMESTAMPTZ |
| | | `likes_count` | счетчик лайков (агрегат) | INTEGER |
| **Лайк публикации** | `publication_likes` | `id` | уникальный идентификатор лайка (PK) | UUID |
| | | `user_id` | поставивший лайк (FK → users.id) | UUID |
| | | `publication_id` | лайкнутая публикация (FK → publications.id) | UUID |
| | | `created_at` | дата/время лайка | TIMESTAMPTZ |
| **Лайк комментария** | `comment_likes` | `id` | уникальный идентификатор лайка (PK) | UUID |
| | | `user_id` | поставивший лайк (FK → users.id) | UUID |
| | | `comment_id` | лайкнутый комментарий (FK → comments.id) | UUID |
| | | `created_at` | дата/время лайка | TIMESTAMPTZ |
| **Сохранённое** | `saved_items` | `id` | уникальный идентификатор сохранения (PK) | UUID |
| | | `user_id` | кто сохранил (FK → users.id) | UUID |
| | | `publication_id` | что сохранено (FK → publications.id) | UUID |
| | | `added_at` | когда сохранено | TIMESTAMPTZ |
| | | `note` | заметка пользователя | TEXT |
| **Рекомендация** | `recommendations` | `id` | уникальный идентификатор записи (PK) | UUID |
| | | `user_id` | получатель рекомендации (FK → users.id) | UUID |
| | | `publication_id` | рекомендуемая публикация (FK) | UUID |
| | | `algorithm` | название/тип алгоритма | TEXT |
| | | `reason` | обоснование рекомендации | TEXT |
| | | `rank` | позиция/рейтинг рекомендации | INTEGER |
| | | `created_at` | когда сформировано | TIMESTAMPTZ |
| | | `hidden` | скрыта ли рекомендация пользователем | BOOLEAN |
| **Просмотр публикации** | `publications_views` | `view_uuid` | уникальный идентификатор просмотра (PK) | UUID |
| | | `user_id` | кто смотрел (FK → users.id) | UUID |
| | | `publication_id` | что смотрел (FK → publications.id) | UUID |
| | | `viewed_at` | когда смотрел | TIMESTAMPTZ |
| **Подписка** | `user_follows` | `id` | уникальный идентификатор подписки (PK) | UUID |
| | | `follower_id` | кто подписывается (FK → users.id) | UUID |
| | | `following_id` | на кого подписывается (FK → users.id) | UUID |
| | | `created_at` | дата/время подписки | TIMESTAMPTZ |
| **Тег** | `tags` | `id` | уникальный идентификатор тега (PK) | UUID |
| | | `name` | название тега (уникальное) | TEXT |
| | | `description` | описание тега | TEXT |
| | | `usage_count` | счетчик использования | INTEGER |
| | | `created_at` | дата/время создания | TIMESTAMPTZ |
| **Связь публикации–тег** | `publication_tags` | `publication_id` | публикация (FK → publications.id) | UUID |
| | | `tag_id` | тег (FK → tags.id) | UUID |
| | | `created_at` | дата/время связи | TIMESTAMPTZ |
| **Уведомление** | `notifications` | `id` | уникальный идентификатор уведомления (PK) | UUID |
| | | `user_id` | получатель уведомления (FK → users.id) | UUID |
| | | `type` | тип уведомления | TEXT |
| | | `title` | заголовок уведомления | TEXT |
| | | `message` | текст уведомления | TEXT |
| | | `data` | дополнительные данные | JSONB |
| | | `is_read` | прочитано ли уведомление | BOOLEAN |
| | | `created_at` | дата/время создания | TIMESTAMPTZ |
| **Сессия** | `user_sessions` | `id` | уникальный идентификатор сессии (PK) | UUID |
| | | `user_id` | пользователь (FK → users.id) | UUID |
| | | `token_hash` | хеш JWT токена (уникальный) | TEXT |
| | | `refresh_token_hash` | хеш refresh токена | TEXT |
| | | `expires_at` | дата/время истечения | TIMESTAMPTZ |
| | | `created_at` | дата/время создания | TIMESTAMPTZ |
| | | `last_used_at` | дата/время последнего использования | TIMESTAMPTZ |
| | | `user_agent` | информация о браузере | TEXT |
| | | `ip_address` | IP адрес | INET |
| **Упоминание** | `mentions` | `id` | уникальный идентификатор упоминания (PK) | UUID |
| | | `publication_id` | публикация с упоминанием (FK → publications.id) | UUID |
| | | `comment_id` | комментарий с упоминанием (FK → comments.id) | UUID |
| | | `mentioned_user_id` | упомянутый пользователь (FK → users.id) | UUID |
| | | `mentioned_by_user_id` | кто упомянул (FK → users.id) | UUID |
| | | `created_at` | дата/время упоминания | TIMESTAMPTZ |
| **Жалоба** | `reports` | `id` | уникальный идентификатор жалобы (PK) | UUID |
| | | `reporter_id` | кто пожаловался (FK → users.id) | UUID |
| | | `publication_id` | на какую публикацию (FK → publications.id) | UUID |
| | | `comment_id` | на какой комментарий (FK → comments.id) | UUID |
| | | `reported_user_id` | на какого пользователя (FK → users.id) | UUID |
| | | `reason` | причина жалобы | TEXT |
| | | `description` | описание жалобы | TEXT |
| | | `status` | статус обработки | TEXT |
| | | `reviewed_by` | кто рассмотрел (FK → users.id) | UUID |
| | | `reviewed_at` | дата/время рассмотрения | TIMESTAMPTZ |
| | | `created_at` | дата/время жалобы | TIMESTAMPTZ |
| **Коллекция** | `collections` | `id` | уникальный идентификатор коллекции (PK) | UUID |
| | | `user_id` | владелец коллекции (FK → users.id) | UUID |
| | | `name` | название коллекции | TEXT |
| | | `description` | описание коллекции | TEXT |
| | | `is_public` | публичная ли коллекция | BOOLEAN |
| | | `created_at` | дата/время создания | TIMESTAMPTZ |
| | | `updated_at` | дата/время обновления | TIMESTAMPTZ |
| **Элемент коллекции** | `collection_items` | `collection_id` | коллекция (FK → collections.id) | UUID |
| | | `publication_id` | публикация (FK → publications.id) | UUID |
| | | `added_at` | дата/время добавления | TIMESTAMPTZ |
| | | `note` | заметка к элементу | TEXT |
| **Аналитика публикации** | `publication_analytics` | `id` | уникальный идентификатор записи (PK) | UUID |
| | | `publication_id` | публикация (FK → publications.id) | UUID |
| | | `date` | дата аналитики | DATE |
| | | `views_count` | количество просмотров | INTEGER |
| | | `likes_count` | количество лайков | INTEGER |
| | | `comments_count` | количество комментариев | INTEGER |
| | | `shares_count` | количество репостов | INTEGER |
| | | `created_at` | дата/время создания записи | TIMESTAMPTZ |
| | | `updated_at` | дата/время обновления записи | TIMESTAMPTZ |

### Примечания по enum-типам:
- **user_role**: `reader` | `user` | `creator` | `expert` | `super`
- **publication_type**: `quote` | `post` | `article`
- **visibility_type**: `public` | `community` | `private`

## API Эндпоинты

| Актор | Use Case | Маршрут | HTTP-запрос | Аутентификация |
|-------|----------|---------|-------------|----------------|
| **Пользователь** | UC 0 Войти в систему | `/auth/login` | POST | нет |
| **Пользователь** | UC 0.1 Проверить токен | `/auth/check` | GET | да |
| **Пользователь** | UC 0.2 Зарегистрироваться | `/auth/register` | POST | нет |
| **Пользователь** | UC 0.3 Выйти из системы | `/auth/logout` | POST | да |
| **Пользователь** | UC 1.1 Создать публикацию | `/publication/create` | POST | да |
| **Пользователь** | UC 1.2 Получить публикацию | `/publication/{id}` | GET | да |
| **Пользователь** | UC 1.3 Редактировать публикацию | `/publication/{id}` | PUT | да |
| **Пользователь** | UC 1.4 Удалить публикацию | `/publication/{id}` | DELETE | да |
| **Пользователь** | UC 1.5 Поставить лайк | `/publication/{id}/like` | POST | да |
| **Пользователь** | UC 1.6 Получить лайки | `/publication/{id}/likes` | GET | да |
| **Пользователь** | UC 1.7 Сохранить публикацию | `/publication/{id}/save` | POST | да |
| **Пользователь** | UC 1.8 Убрать из сохраненных | `/publication/{id}/save` | DELETE | да |
| **Пользователь** | UC 2.1 Получить комментарии | `/publication/{id}/comments` | GET | да |
| **Пользователь** | UC 2.2 Создать комментарий | `/publication/{id}/comments` | POST | да |
| **Пользователь** | UC 2.3 Получить комментарий | `/comment/{id}` | GET | да |
| **Пользователь** | UC 2.4 Редактировать комментарий | `/comment/{id}` | PUT | да |
| **Пользователь** | UC 2.5 Удалить комментарий | `/comment/{id}` | DELETE | да |
| **Пользователь** | UC 2.6 Ответить на комментарий | `/comment/{id}/reply` | POST | да |
| **Пользователь** | UC 2.7 Лайкнуть комментарий | `/comment/{id}/like` | POST | да |
| **Пользователь** | UC 3.1 Получить ленту | `/feed` | GET | да |
| **Пользователь** | UC 3.2 Мои публикации | `/feed/me` | GET | да |
| **Пользователь** | UC 3.3 Сохраненные публикации | `/feed/me/saved` | GET | да |
| **Пользователь** | UC 3.4 Публикации пользователя | `/feed/user/{id}` | GET | да |
| **Пользователь** | UC 4.1 Мой профиль | `/profile/me` | GET | да |
| **Пользователь** | UC 4.2 Редактировать профиль | `/profile/me` | POST | да |
| **Пользователь** | UC 4.3 Профиль пользователя | `/profile/{id}` | GET | да |
| **Пользователь** | UC 4.4 Статистика пользователя | `/profile/{id}/stats` | GET | да |
| **Пользователь** | UC 4.5 Подписаться на пользователя | `/follow/{id}` | POST | да |
| **Пользователь** | UC 4.6 Отписаться от пользователя | `/follow/{id}` | DELETE | да |
| **Пользователь** | UC 4.7 Получить уведомления | `/notifications` | GET | да |
| **Пользователь** | UC 5.1 Поиск публикаций | `/search` | GET | да |
| **Пользователь** | UC 5.2 Поиск пользователей | `/search/users` | GET | да |
| **Пользователь** | UC 5.3 Прогрев поискового индекса | `/search/warmup` | POST | да |
| **Пользователь** | UC 5.4 Получить популярные теги | `/tags` | GET | да |
| **Пользователь** | UC 6.1 Загрузить медиа-файл | `/media/upload` | POST | да |
| **Пользователь** | UC 6.2 Получить медиа-файл | `/media/{id}` | GET | да |
| **Пользователь** | UC 6.3 Удалить медиа-файл | `/media/{id}` | DELETE | да |
| **Пользователь** | UC 7.1 Получить рекомендации | `/recommendations` | POST | да |
| **Пользователь** | UC 7.2 Лента рекомендаций | `/recommendations/feed` | GET | да |
| **Пользователь** | UC 7.3 Скрыть рекомендацию | `/recommendations/{id}/hide` | POST | да |
| **Пользователь** | UC 7.4 Очистить текст | `/purify` | POST | да |
