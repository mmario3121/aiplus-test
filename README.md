# AIPlus Test

## Installation

1. Clone the repository from GitHub:
    ```bash
    git clone https://github.com/mmario3121/aiplus-test.git
    ```
2. Navigate into the project directory:
    ```bash
    cd aiplus-test
    ```
3. Copy the `.env.example` file to `.env` and modify it with your own values:
    ```bash
    cp .env.example .env
    ```
   Open the `.env` file and replace the placeholders with your actual values.
4. Build and start the Docker containers using `docker-compose`:
    ```bash
    docker-compose up --build
    ```

## APIs

### Employee

#### Create a new employee
- **Endpoint:** `POST /employee`
- **Description:** Create a new employee. The request body should contain the employee information in JSON format.
- **Example request body:**
    ```json
    {
        "name": "Test",
        "phone": "87077777777",
        "city_id": 1
    }
    ```
- **Example response:**
    ```json
    HTTP/1.1 201 Created
    Content-Type: application/json

    {
        "id": 1,
        "name": "Test",
        "phone": "87077777777",
        "city_id": 1
    }
    ```

#### Retrieve a list of all employees
- **Endpoint:** `GET /employee`
- **Description:** Retrieve a list of all employees.
- **Example response:**
    ```json
    HTTP/1.1 200 OK
    Content-Type: application/json

    [
        {
            "id": 1,
            "name": "Test",
            "phone": "87077777777",
            "city_id": 1
        },
        {
            "id": 2,
            "name": "Test 2",
            "phone": "87077777778",
            "city_id": 2
        }
    ]
    ```

#### Retrieve information about an employee by their ID
- **Endpoint:** `GET /employee/{id}`
- **Description:** Retrieve information about an employee by their ID.
- **Example request:**
    ```bash
    GET /employee/1
    ```
- **Example response:**
    ```json
    HTTP/1.1 200 OK
    Content-Type: application/json

    {
        "id": 1,
        "name": "Test",
        "phone": "87077777777",
        "city_id": 1
    }
    ```

#### Delete an employee by their ID
- **Endpoint:** `DELETE /employee/{id}`
- **Description:** Delete an employee by their ID.
- **Example request:**
    ```bash
    DELETE /employee/1
    ```
- **Example response:**
    ```json
    HTTP/1.1 200 OK
    Content-Type: application/json

    {
        "status": "deleted"
    }
    ```
## Run Tests

To run tests, use the following command:

```bash
    go test
```