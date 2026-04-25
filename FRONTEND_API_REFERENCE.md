# Frontend API Reference

This document tracks all REST API endpoints exposed by the API Gateway for frontend consumption. It will be updated as new microservices are integrated.

---

## 1. Auth Service

### **Register Passenger**
Creates a new user passenger account.
* **Method:** `POST`
* **Endpoint:** `/api/v1/auth/register`
* **Headers:** `Content-Type: application/json`

**Request Body:**
```json
{
  "email": "john.doe@example.com",
  "password": "securepassword123",
  "full_name": "John Doe",
  "phone_number": "+251911234567"
}
```

**Response Body (201 Created):**
```json
{
  "user": {
    "id": "uuid-string",
    "email": "john.doe@example.com",
    "full_name": "John Doe",
    "phone_number": "+251911234567",
    "role": "ROLE_PASSENGER",
    "created_at": "2026-04-25T12:00:00Z",
    "updated_at": "2026-04-25T12:00:00Z"
  }
}
```

### **Login**
Authenticates a user and returns an access token.
* **Method:** `POST`
* **Endpoint:** `/api/v1/auth/login`
* **Headers:** `Content-Type: application/json`

**Request Body:**
```json
{
  "email": "john.doe@example.com",
  "password": "securepassword123"
}
```

**Response Body (200 OK):**
```json
{
  "access_token": "jwt.token.string",
  "user": {
    "id": "uuid-string",
    "email": "john.doe@example.com",
    "full_name": "John Doe",
    "phone_number": "+251911234567",
    "role": "ROLE_PASSENGER",
    "created_at": "2026-04-25T12:00:00Z",
    "updated_at": "2026-04-25T12:00:00Z"
  }
}
```

### **Get Current User (Me)**
Fetches the currently authenticated user's profile.
* **Method:** `GET`
* **Endpoint:** `/api/v1/auth/me`
* **Headers:** `Authorization: Bearer <access_token>`

**Response Body (200 OK):**
```json
{
  "user": {
    "id": "uuid-string",
    "email": "john.doe@example.com",
    "full_name": "John Doe",
    "phone_number": "+251911234567",
    "role": "ROLE_PASSENGER",
    "created_at": "2026-04-25T12:00:00Z",
    "updated_at": "2026-04-25T12:00:00Z"
  }
}
```

---

## 2. Fleet Service

### **Add a New Bus**
Registers a new bus entity to the fleet.
* **Method:** `POST`
* **Endpoint:** `/api/v1/fleet/buses`
* **Headers:** `Content-Type: application/json`

**Request Body:**
```json
{
  "plate_number": "AA-1234A",
  "operator_name": "Selam Bus",
  "capacity": 55
}
```
**Response Body (201 Created):**
```json
{
  "id": "uuid-string",
  "plate_number": "AA-1234A",
  "operator_name": "Selam Bus",
  "capacity": 55
}
```

### **List Buses**
* **Method:** `GET`
* **Endpoint:** `/api/v1/fleet/buses?limit=10&offset=0`
**Response Body (200 OK):**
```json
{
  "buses": [
    {
      "id": "uuid-string",
      "plate_number": "AA-1234A",
      "operator_name": "Selam Bus",
      "capacity": 55
    }
  ]
}
```

### **Add a New Route**
* **Method:** `POST`
* **Endpoint:** `/api/v1/fleet/routes`

**Request Body:**
```json
{
  "origin": "Addis Ababa",
  "destination": "Hawassa",
  "distance_km": 275,
  "estimated_duration_mins": 300
}
```

### **List Routes**
* **Method:** `GET`
* **Endpoint:** `/api/v1/fleet/routes`

### **Create a Schedule**
* **Method:** `POST`
* **Endpoint:** `/api/v1/fleet/schedules`

**Request Body:**
```json
{
  "route_id": "uuid-string",
  "bus_id": "uuid-string",
  "departure_time": "2026-04-26T08:00:00Z",
  "arrival_time": "2026-04-26T13:00:00Z",
  "price": 450.00,
  "status": "SCHEDULED"
}
```

### **List Schedules**
* **Method:** `GET`
* **Endpoint:** `/api/v1/fleet/schedules`
