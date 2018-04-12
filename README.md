
# GO RESTful API

A RESTful API built in go. 

## Setup

Run the MySQL script provided (`schema.sql`) in a MySQL server and configure the database credentials in the application code to run.

## Routes

### _GET /api/v1/auth_

Return the authorization token.

Output:

```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNTIzNDgwODExfQ._S0U3-dTqa3YO1D9x4ZAQ-FIE6A5hGOeO1g-zkU3HIs"
}
```

### _GET /api/v1/users_

Return the list of users.

Output:

```json
[
    {
        "id": 10,
        "name": "Adrastea Reinhold Cecilia",
        "gravatar": "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTwFogmoqorCL-fRpN5i74u5NgaddmgycTiKNHm_SSOnh-vAmjb5Q"
    },
    {
        "id": 9,
        "name": "Iphigeneia Franziska Balbino",
        "gravatar": "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTwFogmoqorCL-fRpN5i74u5NgaddmgycTiKNHm_SSOnh-vAmjb5Q"
    },
    {
        "id": 8,
        "name": "Adrastea Leonard Ligeia",
        "gravatar": "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTwFogmoqorCL-fRpN5i74u5NgaddmgycTiKNHm_SSOnh-vAmjb5Q"
    }
]
```

### _GET /api/v1/users/{user_id}_

Return the user's details.

Output:

```json
{
    "id": 1,
    "name": "Widad Niklas",
    "gravatar": "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcS8RpEBBTnUZhKxq9gHAV_8jVSKGF9E8p6cUzaWaQl8BAII_Elt"
}
```

### _GET /api/v1/widgets_

Return the list of widgets.

Output:

```json
[
    {
        "id": 10,
        "name": "Amulet of Wisdom",
        "color": "Silver",
        "price": 54.2,
        "melts": true,
        "inventory": 4
    },
    {
        "id": 9,
        "name": "Thunder Bow",
        "color": "Brown",
        "price": 65.54,
        "melts": false,
        "inventory": 15
    }
]
```

### _GET /api/v1/widgets/{widget_id}_

Return the widget's details.

Output:

```json
{
    "id": 1,
    "name": "Chest",
    "color": "Red",
    "price": 10.6,
    "melts": false,
    "inventory": 50
}
```

### _POST /api/v1/widgets_

Create a new widget.

Input:

```json
{
  "name": "My new widget",
  "price": 10,
  "inventory": 42,
  "melts": true,
  "color": "Blue"
}
```

Output:

```json
{
    "message": "Widget created successfully"
}
```

### _PUT /api/v1/widgets/{widget_id}_

Update a widget.

Input:

```json
{
  "name": "My widget",
  "price": 10,
  "inventory": 42,
  "melts": true,
  "color": "Blue"
}
```

Output:

```json
{
    "message": "Widget updated successfully"
}
```