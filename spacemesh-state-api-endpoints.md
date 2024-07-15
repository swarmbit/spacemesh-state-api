# API

## Requests

### **GET** - /network/info

#### CURL

```sh
curl -X GET "https://spacemesh-api-v2.swarmbit.io/network/info" \
    -H "x-api-key: <api-key>"
```

#### Header Parameters

- **x-api-key** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "<api-key>"
  ],
  "default": "<api-key>"
}
```

### **GET** - /account

#### CURL

```sh
curl -X GET "https://spacemesh-api-v2.swarmbit.io/account\
?offset=0&limit=20&sort=desc" \
    -H "x-api-key: <api-key>"
```

#### Query Parameters

- **offset** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "0"
  ],
  "default": "0"
}
```
- **limit** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "20"
  ],
  "default": "20"
}
```
- **sort** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "desc"
  ],
  "default": "desc"
}
```

#### Header Parameters

- **x-api-key** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "<api-key>"
  ],
  "default": "<api-key>"
}
```

### **GET** - /account/post/epoch/13

#### CURL

```sh
curl -X GET "https://spacemesh-api-v2.swarmbit.io/account/post/epoch/13\
?offset=0&limit=20&sort=desc" \
    -H "x-api-key: <api-key>"
```

#### Query Parameters

- **offset** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "0"
  ],
  "default": "0"
}
```
- **limit** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "20"
  ],
  "default": "20"
}
```
- **sort** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "desc"
  ],
  "default": "desc"
}
```

#### Header Parameters

- **x-api-key** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "<api-key>"
  ],
  "default": "<api-key>"
}
```

### **GET** - /account/sm1qqqqqqpzvpdcm0c09aac3fvzywmt7v0dyqvpygq55xla6

#### CURL

```sh
curl -X GET "https://spacemesh-api-v2.swarmbit.io/account/sm1qqqqqqpzvpdcm0c09aac3fvzywmt7v0dyqvpygq55xla6" \
    -H "x-api-key: <api-key>"
```

#### Header Parameters

- **x-api-key** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "<api-key>"
  ],
  "default": "<api-key>"
}
```

### **GET** - /account/sm1qqqqqqpzvpdcm0c09aac3fvzywmt7v0dyqvpygq55xla6/rewards

#### CURL

```sh
curl -X GET "https://spacemesh-api-v2.swarmbit.io/account/sm1qqqqqqpzvpdcm0c09aac3fvzywmt7v0dyqvpygq55xla6/rewards\
?offset=0&limit=20&sort=desc" \
    -H "x-api-key: <api-key>"
```

#### Query Parameters

- **offset** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "0"
  ],
  "default": "0"
}
```
- **limit** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "20"
  ],
  "default": "20"
}
```
- **sort** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "desc"
  ],
  "default": "desc"
}
```

#### Header Parameters

- **x-api-key** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "<api-key>"
  ],
  "default": "<api-key>"
}
```

### **GET** - /account/sm1qqqqqqpzvpdcm0c09aac3fvzywmt7v0dyqvpygq55xla6/rewards/details

#### CURL

```sh
curl -X GET "https://spacemesh-api-v2.swarmbit.io/account/sm1qqqqqqpzvpdcm0c09aac3fvzywmt7v0dyqvpygq55xla6/rewards/details" \
    -H "x-api-key: <api-key>"
```

#### Header Parameters

- **x-api-key** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "<api-key>"
  ],
  "default": "<api-key>"
}
```

### **GET** - /account/sm1qqqqqqpzvpdcm0c09aac3fvzywmt7v0dyqvpygq55xla6/rewards/details/13

#### CURL

```sh
curl -X GET "https://spacemesh-api-v2.swarmbit.io/account/sm1qqqqqqpzvpdcm0c09aac3fvzywmt7v0dyqvpygq55xla6/rewards/details/13" \
    -H "x-api-key: <api-key>"
```

#### Header Parameters

- **x-api-key** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "<api-key>"
  ],
  "default": "<api-key>"
}
```

### **GET** - /account/sm1qqqqqqpzvpdcm0c09aac3fvzywmt7v0dyqvpygq55xla6/transactions

#### CURL

```sh
curl -X GET "https://spacemesh-api-v2.swarmbit.io/account/sm1qqqqqqpzvpdcm0c09aac3fvzywmt7v0dyqvpygq55xla6/transactions\
?offset=0&limit=20&sort=desc&complete=true&method=spawn&minAmount=-1" \
    -H "x-api-key: <api-key>"
```

#### Query Parameters

- **offset** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "0"
  ],
  "default": "0"
}
```
- **limit** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "20"
  ],
  "default": "20"
}
```
- **sort** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "desc"
  ],
  "default": "desc"
}
```
- **complete** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "true"
  ],
  "default": "true"
}
```

#### Header Parameters

- **x-api-key** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "<api-key>"
  ],
  "default": "<api-key>"
}
```

### **POST** - /account/group

#### CURL

```sh
curl -X POST "https://spacemesh-api-v2.swarmbit.io/account/group" \
    -H "x-api-key: <api-key>" \
    -H "Content-Type: application/json; charset=utf-8" \
    --data-raw "$body"
```

#### Header Parameters

- **x-api-key** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "<api-key>"
  ],
  "default": "<api-key>"
}
```
- **Content-Type** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "application/json; charset=utf-8"
  ],
  "default": "application/json; charset=utf-8"
}
```

#### Body Parameters

- **body** should respect the following schema:

```
{
  "type": "string",
  "default": "{\"accounts\":[\"sm1qqqqqq82d3yv8m632dsn237wg6sa3frsy2eruysclwgvm\",\"sm1qqqqqqpy8svgfxfhh2w42ujhrynsxgwsrrzgplss9t0lv\"]}"
}
```

### **GET** - /account/sm1qqqqqqpzvpdcm0c09aac3fvzywmt7v0dyqvpygq55xla6/atx/13

#### CURL

```sh
curl -X GET "https://spacemesh-api-v2.swarmbit.io/account/sm1qqqqqqpzvpdcm0c09aac3fvzywmt7v0dyqvpygq55xla6/atx/13\
?offset=0&limit=20&sort=desc" \
    -H "x-api-key: <api-key>"
```

#### Query Parameters

- **offset** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "0"
  ],
  "default": "0"
}
```
- **limit** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "20"
  ],
  "default": "20"
}
```
- **sort** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "desc"
  ],
  "default": "desc"
}
```

#### Header Parameters

- **x-api-key** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "<api-key>"
  ],
  "default": "<api-key>"
}
```

### **POST** - /account/sm1qqqqqqpzvpdcm0c09aac3fvzywmt7v0dyqvpygq55xla6/atx/13/filter-active-nodes

#### CURL

```sh
curl -X POST "https://spacemesh-api-v2.swarmbit.io/account/sm1qqqqqqpzvpdcm0c09aac3fvzywmt7v0dyqvpygq55xla6/atx/13/filter-active-nodes" \
    -H "x-api-key: <api-key>" \
    -H "Content-Type: application/json; charset=utf-8" \
    --data-raw "$body"
```

#### Header Parameters

- **x-api-key** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "<api-key>"
  ],
  "default": "<api-key>"
}
```
- **Content-Type** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "application/json; charset=utf-8"
  ],
  "default": "application/json; charset=utf-8"
}
```

#### Body Parameters

- **body** should respect the following schema:

```
{
  "type": "string",
  "default": "{\"nodes\":[\"eb8d43fc76cfb2db703a3b6f620e437c33962a45bc179b7e057558206d22033d\",\"eb8d43fc76cfb2db703a3b6f620e437c33962a45bc179b7e057558206d220222\"]}"
}
```

### **GET** - /nodes/0694caac231c6fe64de0c8f6b9169cbc99a0e9d202894ab26b23260c40e6387c

#### CURL

```sh
curl -X GET "https://spacemesh-api-v2.swarmbit.io/nodes/0694caac231c6fe64de0c8f6b9169cbc99a0e9d202894ab26b23260c40e6387c\
?offset=0&limit=20&sort=desc" \
    -H "x-api-key: <api-key>"
```

#### Query Parameters

- **offset** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "0"
  ],
  "default": "0"
}
```
- **limit** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "20"
  ],
  "default": "20"
}
```
- **sort** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "desc"
  ],
  "default": "desc"
}
```

#### Header Parameters

- **x-api-key** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "<api-key>"
  ],
  "default": "<api-key>"
}
```

### **GET** - /nodes/

#### CURL

```sh
curl -X GET "https://spacemesh-api-v2.swarmbit.io/nodes/" \
    -H "x-api-key: <api-key>"
```

#### Header Parameters

- **x-api-key** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "<api-key>"
  ],
  "default": "<api-key>"
}
```

### **GET** - /nodes/0694caac231c6fe64de0c8f6b9169cbc99a0e9d202894ab26b23260c40e6387c/rewards

#### CURL

```sh
curl -X GET "https://spacemesh-api-v2.swarmbit.io/nodes/0694caac231c6fe64de0c8f6b9169cbc99a0e9d202894ab26b23260c40e6387c/rewards" \
    -H "x-api-key: <api-key>"
```

#### Header Parameters

- **x-api-key** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "<api-key>"
  ],
  "default": "<api-key>"
}
```

### **GET** - /nodes/0694caac231c6fe64de0c8f6b9169cbc99a0e9d202894ab26b23260c40e6387c/rewards/details

#### CURL

```sh
curl -X GET "https://spacemesh-api-v2.swarmbit.io/nodes/0694caac231c6fe64de0c8f6b9169cbc99a0e9d202894ab26b23260c40e6387c/rewards/details" \
    -H "x-api-key: <api-key>"
```

#### Header Parameters

- **x-api-key** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "<api-key>"
  ],
  "default": "<api-key>"
}
```

### **GET** - /nodes/0694caac231c6fe64de0c8f6b9169cbc99a0e9d202894ab26b23260c40e6387c/rewards/eligibility

#### CURL

```sh
curl -X GET "https://spacemesh-api-v2.swarmbit.io/nodes/0694caac231c6fe64de0c8f6b9169cbc99a0e9d202894ab26b23260c40e6387c/rewards/eligibility" \
    -H "x-api-key: <api-key>"
```

#### Header Parameters

- **x-api-key** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "<api-key>"
  ],
  "default": "<api-key>"
}
```

### **GET** - /epochs/13

#### CURL

```sh
curl -X GET "https://spacemesh-api-v2.swarmbit.io/epochs/13" \
    -H "x-api-key: <api-key>"
```

#### Header Parameters

- **x-api-key** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "<api-key>"
  ],
  "default": "<api-key>"
}
```

### **GET** - /epochs/13/atx

#### CURL

```sh
curl -X GET "https://spacemesh-api-v2.swarmbit.io/epochs/13/atx\
?offset=0&limit=20&sort=desc" \
    -H "x-api-key: <api-key>"
```

#### Query Parameters

- **offset** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "0"
  ],
  "default": "0"
}
```
- **limit** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "20"
  ],
  "default": "20"
}
```
- **sort** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "desc"
  ],
  "default": "desc"
}
```

#### Header Parameters

- **x-api-key** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "<api-key>"
  ],
  "default": "<api-key>"
}
```

### **GET** - /layers

#### CURL

```sh
curl -X GET "https://spacemesh-api-v2.swarmbit.io/layers\
?offset=0&limit=20&sort=desc" \
    -H "x-api-key: <api-key>"
```

#### Query Parameters

- **offset** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "0"
  ],
  "default": "0"
}
```
- **limit** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "20"
  ],
  "default": "20"
}
```
- **sort** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "desc"
  ],
  "default": "desc"
}
```

#### Header Parameters

- **x-api-key** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "<api-key>"
  ],
  "default": "<api-key>"
}
```

### **GET** - /layers/52785/rewards

#### CURL

```sh
curl -X GET "https://spacemesh-api-v2.swarmbit.io/layers/52785/rewards" \
    -H "x-api-key: <api-key>"
```

#### Header Parameters

- **x-api-key** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "<api-key>"
  ],
  "default": "<api-key>"
}
```

### **GET** - /layers/52785/transactions

#### CURL

```sh
curl -X GET "https://spacemesh-api-v2.swarmbit.io/layers/52785/transactions\
?complete=true" \
    -H "x-api-key: <api-key>"
```

#### Query Parameters

- **complete** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "true"
  ],
  "default": "true"
}
```

#### Header Parameters

- **x-api-key** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "<api-key>"
  ],
  "default": "<api-key>"
}
```

### **GET** - /transactions

#### CURL

```sh
curl -X GET "https://spacemesh-api-v2.swarmbit.io/transactions\
?offset=0&limit=20&sort=desc&complete=true" \
    -H "x-api-key: <api-key>"
```

#### Query Parameters

- **offset** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "0"
  ],
  "default": "0"
}
```
- **limit** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "20"
  ],
  "default": "20"
}
```
- **sort** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "desc"
  ],
  "default": "desc"
}
```
- **complete** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "true"
  ],
  "default": "true"
}
```

#### Header Parameters

- **x-api-key** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "<api-key>"
  ],
  "default": "<api-key>"
}
```

### **GET** - /transactions/92d44e654c

#### CURL

```sh
curl -X GET "https://spacemesh-api-v2.swarmbit.io/transactions/92d44e654c" \
    -H "x-api-key: <api-key>"
```

#### Header Parameters

- **x-api-key** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "<api-key>"
  ],
  "default": "<api-key>"
}
```

### **GET** - /poets

#### CURL

```sh
curl -X GET "https://spacemesh-api-v2.swarmbit.io/poets\
?offset=0&limit=20&sort=desc&complete=true" \
    -H "x-api-key: <api-key>"
```

#### Query Parameters

- **offset** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "0"
  ],
  "default": "0"
}
```
- **limit** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "20"
  ],
  "default": "20"
}
```
- **sort** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "desc"
  ],
  "default": "desc"
}
```
- **complete** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "true"
  ],
  "default": "true"
}
```

#### Header Parameters

- **x-api-key** should respect the following schema:

```
{
  "type": "string",
  "enum": [
    "<api-key>"
  ],
  "default": "<api-key>"
}
```

## References

