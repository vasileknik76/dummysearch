# dummysearch

DummySearch is Full Text Search and text comparsion engine.
Its work is based on the TF-IDF metric.
All operations with data are performed via REST API.
Documents in index may have some extra data, but it not uses in search.

You can use any language, but engine uses snowball stemmer (https://github.com/kljensen/snowball), so languages list restricted with:
 - English,
 - Spanish (español),
 - French (le français),
 - Russian (ру́сский язы́к),
 - Swedish (svenska),
 - Norwegian (norsk)

## Contents
- [Build and run](#build-and-run)
- [Operations](#operations)
    - [Creating new index](#creating-new-index)
    - [Remove index](#remove-index)
    - [Add document to index](#add-document-to-index)
    - [Bulk add document to index](#bulk-add-document-to-index)
    - [Get document by id](#get-document-by-id)
    - [Delete document by id](#delete-document-by-id)
    - [Search documents](#search-documents-by-query)
    - [Compare two documents](#compare-two-documents)


### Build and run:

#### Build native:
```shell script
$ go build -o build/dummysearch cmd/dummysearch/main.go
```

#### Build docker:
```shell script
$ docker build -t dummysearch .
```

#### Run native:
```shell script
$ ./build/dummysearch
```

#### Run in docker:
```shell script
$ docker run -p 6745:6745 -it -d dummysearch
```

### Operations:

#### Creating new index:

```shell script
$ curl --location --request POST 'http://localhost:6745/' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "lol",
    "config": {
        "language": "english"
    }
}'
```
Response:
```json
{
  "status": true,
  "payload": {
    "Message": "OK"
  }
}
```

#### Remove index:

```shell script
curl --location --request DELETE 'http://localhost:6745/lol/'
```

Response:

```json
{
  "status": true,
  "payload": {
    "Message": "OK"
  }
}
```

#### Add document to index:

```shell script
curl --location --request POST 'http://localhost:6745/lol/' \
--header 'Content-Type: application/json' \
--data-raw '{
    "content": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
    "meta": {
        "someField": "any value",
        "otherField": 1
    }
}'
```
Response:
```json
{
  "status": true,
  "payload": {
    "Message": "OK",
    "DocumentId": 0
  }
}
```

#### Bulk add document to index:

```shell script
curl --location --request POST 'http://localhost:6745/lol/batch' \
--header 'Content-Type: application/json' \
--data-raw '[
    {
        "content": "some text!",
        "meta": {
            "foo": "bar"
        }
    },
    {
        "content": "london is the capital of great britain",
        "meta": {
            "bar": "baz"
        }
    },
    {
        "content": "any other, text.",
        "meta": {
            "foo": "bar2"
        }
    },
]'
```

Response:

```json
{
  "status": true,
  "payload": {
    "Message": "OK",
    "DocumentIds": [
      1,
      2,
      3
    ]
  }
}
```

#### Get document by id:

Source text content not stored, so you can only receive document meta.

```shell script
curl --location --request GET 'http://localhost:6745/lol/0'
```

Response:

```json
{
  "status": true,
  "payload": {
    "Doc": {
      "Meta": {
        "otherField": 1,
        "someField": "any value"
      }
    }
  }
}
```

#### Delete document by id:

```shell script
curl --location --request DELETE 'http://localhost:6745/lol/0'
```

Response:

```json
{
  "status": true,
  "payload": {
    "Message": "OK"
  }
}
```

#### Search documents by query:

```shell script
curl --location --request GET 'http://localhost:6745/lol/search?query=lorem%20london'
```

Response:

```json
{
  "status": true,
  "payload": [
    {
      "DocId": 2,
      "Meta": {
        "bar": "baz"
      },
      "Score": 0.26726124191242445
    },
    {
      "DocId": 0,
      "Meta": {
        "otherField": 1,
        "someField": "any value"
      },
      "Score": 0.07669649888473704
    }
  ]
}
```

#### Compare two documents:

```shell script
curl --location --request GET 'http://localhost:6745/lol/compare?doc1=1&doc2=3'
```

Response:

```json
{
  "status": true,
  "payload": {
    "score": 0.14907119849998599
  }
}
```
