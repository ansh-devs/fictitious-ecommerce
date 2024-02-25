
# **Ecommerce Api with NATS**


## Structure

![Project Structure](https://raw.githubusercontent.com/ansh-devs/ecomm-poc/main/assets/project_scaffolding.png)


## API Reference

### Login Api

```http
  POST /auth/login
```

| Parameter | Type     | Description                             |
| :-------- | :------- | :---------------------------------------|
| `username`    | `string`   | **Note**. Username|
| `password`    | `string`   | **Note**. Password|

### Orders Api

#### Get item

```http
  GET /orders/get-order/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of item to fetch |

### Products Api


```http
  GET /products/get-product/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of item to fetch |


```http
  GET /products/get-all/
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of item to fetch |



