# C4 Architecture Diagrams - Visual Index

This page provides a quick visual index of all architecture diagrams for the Horo platform.

---

## 1. System Context Diagram

Shows the high-level system boundary with external actors and systems.

![Context Diagram](images/Horo%20System%20-%20Context%20Diagram.png)

---

## 2. Container Diagram

Shows all microservices, databases, and message broker with their relationships.

![Container Diagram](images/Horo%20System%20-%20Container%20Diagram.png)

---

## 3. Component Diagrams

### 3.1 API Gateway Components

![API Gateway](images/API%20Gateway%20-%20Component%20Diagram.png)

### 3.2 User Management Service Components

![User Management Service](images/User%20Management%20Service%20-%20Component%20Diagram.png)

### 3.3 Course Service Components

![Course Service](images/Course%20Service%20-%20Component%20Diagram.png)

### 3.4 Chat Service Components

![Chat Service](images/Chat%20Service%20-%20Component%20Diagram.png)

### 3.5 Order Service Components

![Order Service](images/Order%20Service%20-%20Component%20Diagram.png)

### 3.6 Payment Service Components

![Payment Service](images/Payment%20Service%20-%20Component%20Diagram.png)

---

## 4. Sequence Diagrams

### 4.1 Order and Payment Flow

Complete flow from order creation to payment completion and chat room creation.

![Order and Payment Flow](images/Order%20and%20Payment%20Flow%20-%20Sequence%20Diagram.png)

### 4.2 Real-time Chat Flow

WebSocket-based real-time messaging flow through RabbitMQ.

![Real-time Chat Flow](images/Real-time%20Chat%20Flow%20-%20Sequence%20Diagram.png)

---

## Quick Navigation

- [Main Documentation README](README.md)
- [PlantUML Source Files](.)
- [Generated Images](images/)

## Regenerating Images

To regenerate all PNG images from the PlantUML source files:

```bash
cd docs/c4-diagrams
for file in *.puml; do plantuml -o images "$file"; done
```

Note: Requires PlantUML to be installed (`brew install plantuml`)

