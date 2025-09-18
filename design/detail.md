# System Design

This file contains diagrams in Mermaid format: sequence diagram for user profile update flow and ER diagram for the database schema.

## Sequence Diagram: Update Profile

```mermaid
sequenceDiagram
    participant User
    participant Client
    participant API as "REST API (Fiber)"
    participant Auth as "Auth (JWT)"
    participant DB as "SQLite (GORM)"

    User->>Client: Open Profile Page
    Client->>API: GET /profile (Authorization: Bearer token)
    API->>Auth: validate token
    Auth-->>API: user id
    API->>DB: SELECT * FROM users WHERE id = ?
    DB-->>API: user record
    API-->>Client: 200 { profile }

    User->>Client: Edit fields and Save
    Client->>API: PUT /profile (Authorization: Bearer token, body)
    API->>Auth: validate token
    Auth-->>API: user id
    API->>DB: UPDATE users SET ... WHERE id = ?
    DB-->>API: success
    API-->>Client: 200 { updated profile }

```

## ER Diagram

```mermaid
erDiagram
    USERS {
        INTEGER id PK
        datetime created_at
        datetime updated_at
        datetime deleted_at
        TEXT email "NOT NULL, UNIQUE"
        TEXT password "NOT NULL"
        TEXT first_name
        TEXT last_name
        TEXT phone
        TEXT member_level
        INTEGER points
    }

    %% No other tables in this simple design
```
