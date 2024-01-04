## Test endpoints API using curl
* Creating new task

`Request`

```bash
curl -i --request POST 'http://localhost:8080/v1/tasks' \
--header 'Content-Type: application/json' \
--data-raw '{
    "title": "Test Task"
}'
```

```json
{
    "title":"test task",
    "created_at":"2024-01-04T09:41:30Z",
    "updated_at":"2024-01-04T09:41:30Z"
}
```
