## API Documentation

###Post Structure

```json
{
  "id": "string - AUTO GENERATED - UNIQUE",
  "name": "string - REQUIRED", 
  "body": "string - REQUIRED"
}
```


#### GET `/posts`

returns array of all posts

#### GET `/posts/:id`

returns single post<br />
returns error if post not found<br />

#### POST `/posts`

create new post<br />
returns created post<br />
expects following JSON structure<br />

```json
{
  "name": "string - REQUIRED", 
  "body": "string - REQUIRED"
}
```

#### PUT `/posts/:id`

update existing post<br />
returns updated post<br />
expects following JSON structure

```json
{
  "name": "string - REQUIRED", 
  "body": "string - REQUIRED"
}
```

#### DELETE `/posts/:id`

returns error false when successfully deleting post<br />