# Users
## Названачения
Микросервис для управления пользователями в проекте [Farm](https://github.com/mongolianpes/Farm)
Используется для регистрации, авторизации и получения информации о пользователях

## Функционал
- Позволяет регистрировать пользователей, проверяя соответствует ли имя пользователя, пароль и интересы требованиям. Полные требования к регистрации пользователя ниже.
- Способствует авторизации пользователей, но не управляет их сессииями. Управление сессиями лежит на основном сервисе
- Возвращает информацию о пользователей по ID
- Возвращает ID пользователя по login
- Позволяет удалить пользователя

Требования для регистрации пользователя:
- Логин не равен "my"
- Пароль больше 9 символов
- Длина интересов от 10 до 250 символов

## RPC запросы:
```
service Users {
  rpc GetUserInfo(GetUserInfoRequest) returns (GetUserInfoResponse);
  rpc Auth(AuthRequest) returns (AuthResponse);
  rpc Register(RegisterRequest) returns (RegisterResponse);
  rpc AddAvatar(AddAvatarRequest) returns (AddAvatarResponse);
  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
  rpc GetUserID(GetUserIDRequest) returns (GetUserIDResponse);
}
```
Подобнее в файле Users/internal/proto/users.proto

## Переменные окужения
Для корректной работы необходимы переменные окружения:
DB_HOST - адрес БД
DB_PORT - порт БД
DB_USER - пользователей БД
DB_PASSWORD - пароль пользователя БД
DB_NAME - название БД
OLLAMA_HOST - адрес нейросети Ollama и порт, разделенные двоеточием