# posts_service

# Запуск сервиса

`make start` или `docker-compose up -d`

# Выбор типа хранилища

В docker-compose файле уставноить переменную env: 

REPO_TYPE: "postgres" - хранение в БД \
REPO_TYPE: "in-memory" - хранение в памяти

# Пример запросов

## Создание поста

Запрос:

```
mutation {
  createPost(
     title: "Test Post"
     content: "Content"
     allowComments: true
   )
}
```

Ответ:

```
"2011b617-f29f-49b3-9451-0952197de080" (Пример)
```

## Добавление комментария

```
mutation {
  createComment(postId: "2011bb17-f29f-69b3-9451-0952197de080", parentId: null, content: "Comment") {
    id
    desc
  }
}
```

```
"2bf1b157-f29f-49b3-9451-0952197de087" (Пример)
```

## Добавление ответа комментария 

```
mutation {
  createComment(postId: "2011bb17-f29f-69b3-9451-0952197de080", parentId: "2bf1b157-f29f-49b3-9451-0952197de087", content: "Comment") {
    id
    desc
  }
}
```

## Подписка на пост

```
subscription {
  commentAdded(postId: "2011bb17-f29f-69b3-9451-0952197de080") {
    id
    postId
    parentId
    content
    createdAt
  }
}
```

## Просмотр постов

```
query {
  posts {
    id
    title
    content
    allowComments
  }
}
```
