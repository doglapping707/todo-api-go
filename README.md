## Test endpoints API using curl
* Creating new task

`Request`

```bash
curl -i --request POST 'http://localhost:8080/v1/tasks' \
--header 'Content-Type: application/json' \
--data-raw '{
    "title": "TestTask 1"
}'
```

`Response`

```json
{
    "title":"Test Task",
    "created_at":"2024-01-04T10:02:14Z",
    "updated_at":"2024-01-04T10:02:14Z"
}
```

* Updating new task

`Request`

```bash
curl -i --request PUT 'http://localhost:8080/v1/tasks' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id":    1,
    "title": "TestTask 2"
}'
```