## Test endpoints API using curl
* Create a new task

`Request`

```bash
curl -i --request POST 'http://localhost:8080/v1/tasks' \
--header 'Content-Type: application/json' \
--data-raw '{
    "title": "Task_1"
}'
```

`Response`

```json
{
    "title":"Task_1",
    "created_at":"2024-01-04T10:02:14Z",
    "updated_at":"2024-01-04T10:02:14Z"
}
```

* Update a task

`Request`

```bash
curl -i --request PUT 'http://localhost:8080/v1/tasks' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id":    1,
    "title": "Task_2"
}'
```

* Liste a tasks

`Request`

```bash
curl -i --request GET 'http://localhost:8080/v1/tasks'
```

`Response`

```json
{
    "id":1,
    "title":"Task_1",
},
{
    "id":2,
    "title":"Task_2",
},
{
    "id":3,
    "title":"Task_3",
},
```
