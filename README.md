# Queue Cleaner

連結 RabbitMQ 的 `Management Console` 頁面，使用 API 取得

1. 已經沒有任何 `consumer` 連接的 queue
2. 需刪除但未被刪除的 queue (相關規則定義在 `queue_management/remnant`)

## Env

- RABBITMQ_MANAGEMENT_URL 管理平台網址
- USERNAME 使用者帳號
- PASSWORD 使用者密碼
- DELETE_MODE (`ON`/`OFF`)，若設定為 `ON`，在取得列表後，會將列表內的 queue 刪除
