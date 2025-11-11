# Horo Platform - C4 Architecture Diagrams

This directory contains C4 architecture diagrams for the Horo platform, a microservices-based consultation and course platform.

## Overview

The Horo platform enables customers to browse courses, purchase consultations from prophets/astrologers, and engage in real-time chat consultations. The system follows a microservices architecture with event-driven communication.

## Diagrams

### Level 1: Context Diagram

- **File**: `01-context-diagram.puml`
- **Image**: [View PNG](images/Horo%20System%20-%20Context%20Diagram.png)
- **Description**: Shows the system boundary and external actors (customers, prophets) and systems (Firebase, databases, message broker)
- **Purpose**: High-level view of the system in its environment

![Context Diagram](images/Horo%20System%20-%20Context%20Diagram.png)

### Level 2: Container Diagram

- **File**: `02-container-diagram.puml`
- **Image**: [View PNG](images/Horo%20System%20-%20Container%20Diagram.png)
- **Description**: Shows all microservices, databases, and the message broker with their relationships
- **Services**:
  - API Gateway (entry point, WebSocket hub)
  - User Management Service (authentication, user profiles)
  - Course Service (course catalog, reviews)
  - Chat Service (real-time messaging)
  - Order Service (order lifecycle)
  - Payment Service (payment processing)

![Container Diagram](images/Horo%20System%20-%20Container%20Diagram.png)

### Level 3: Component Diagrams

#### API Gateway Components

- **File**: `03-component-api-gateway.puml`
- **Image**: [View PNG](images/API%20Gateway%20-%20Component%20Diagram.png)
- **Key Components**: Router, Middleware, HTTP Handlers, WebSocket Handler, Hub, Message Publisher/Consumer

![API Gateway Components](images/API%20Gateway%20-%20Component%20Diagram.png)

#### User Management Service Components

- **File**: `04-component-user-management.puml`
- **Image**: [View PNG](images/User%20Management%20Service%20-%20Component%20Diagram.png)
- **Key Components**: HTTP Handler, gRPC Server, User Management App, Auth App, Firebase Adapter, User Repository

![User Management Service Components](images/User%20Management%20Service%20-%20Component%20Diagram.png)

#### Course Service Components

- **File**: `05-component-course-service.puml`
- **Image**: [View PNG](images/Course%20Service%20-%20Component%20Diagram.png)
- **Key Components**: HTTP Handler, gRPC Server, Course Application, Course Repository, User gRPC Client

![Course Service Components](images/Course%20Service%20-%20Component%20Diagram.png)

#### Chat Service Components

- **File**: `06-component-chat-service.puml`
- **Image**: [View PNG](images/Chat%20Service%20-%20Component%20Diagram.png)
- **Key Components**: HTTP Handler, Chat Application, Message/Room Repositories, Message Consumer/Publisher, gRPC Clients

![Chat Service Components](images/Chat%20Service%20-%20Component%20Diagram.png)

#### Order Service Components

- **File**: `07-component-order-service.puml`
- **Image**: [View PNG](images/Order%20Service%20-%20Component%20Diagram.png)
- **Key Components**: HTTP Handler, Order Application, Order Repository, Event Publisher, Payment Consumer, gRPC Clients

![Order Service Components](images/Order%20Service%20-%20Component%20Diagram.png)

#### Payment Service Components

- **File**: `08-component-payment-service.puml`
- **Image**: [View PNG](images/Payment%20Service%20-%20Component%20Diagram.png)
- **Key Components**: HTTP Handler, Payment Application, Payment Repository, Event Publisher, Order Consumer

![Payment Service Components](images/Payment%20Service%20-%20Component%20Diagram.png)

### Sequence Diagrams

#### Order and Payment Flow

- **File**: `09-sequence-order-payment-flow.puml`
- **Image**: [View PNG](images/Order%20and%20Payment%20Flow%20-%20Sequence%20Diagram.png)
- **Description**: Shows the complete flow from order creation to payment completion and chat room creation

![Order and Payment Flow](images/Order%20and%20Payment%20Flow%20-%20Sequence%20Diagram.png)

#### Real-time Chat Flow

- **File**: `10-sequence-chat-flow.puml`
- **Image**: [View PNG](images/Real-time%20Chat%20Flow%20-%20Sequence%20Diagram.png)
- **Description**: Shows WebSocket-based real-time messaging flow through RabbitMQ

![Real-time Chat Flow](images/Real-time%20Chat%20Flow%20-%20Sequence%20Diagram.png)

## Architecture Patterns

### Communication Patterns

1. **Synchronous**:

   - REST API via API Gateway for client-facing operations
   - gRPC for inter-service communication
   - WebSocket for real-time chat

2. **Asynchronous**:
   - RabbitMQ for event-driven communication
   - Message queues for order/payment/chat events

### Key Design Decisions

1. **API Gateway Pattern**: Single entry point for all client requests, handles authentication and routing
2. **Event-Driven Architecture**: Services communicate via events for loose coupling
3. **Database per Service**: Each service owns its data store (MongoDB or PostgreSQL)
4. **Hexagonal Architecture**: Services follow ports & adapters pattern for clean separation of concerns
5. **Real-time Messaging**: WebSocket connections managed at API Gateway, messages flow through RabbitMQ to Chat Service

## Technology Stack

- **Backend**: Go (Fiber framework)
- **Frontend**: Next.js, React
- **Databases**: MongoDB (users, courses, chat), PostgreSQL (orders, payments)
- **Message Broker**: RabbitMQ
- **Authentication**: Firebase Auth
- **Inter-service Communication**: gRPC
- **Real-time**: WebSocket
- **Orchestration**: Kubernetes (Minikube for local), Tilt

## Viewing the Diagrams

### Online Viewers

1. **PlantUML Web Server**: http://www.plantuml.com/plantuml/uml/
2. **PlantText**: https://www.planttext.com/

### Local Rendering

1. **VS Code Extension**: "PlantUML" by jebbs
2. **IntelliJ Plugin**: PlantUML integration
3. **Command Line**:

   ```bash
   # Install PlantUML
   brew install plantuml  # macOS

   # Generate PNG
   plantuml filename.puml

   # Generate SVG
   plantuml -tsvg filename.puml
   ```

## Event Flow Summary

### Order Creation Flow

1. Customer creates order → Order Service
2. Order Service publishes `OrderCreated` → RabbitMQ
3. Payment Service consumes event → Creates payment record
4. Payment Service publishes `PaymentCreated` → RabbitMQ
5. Order Service binds payment to order

### Payment Completion Flow

1. Customer completes payment → Payment Service
2. Payment Service publishes `PaymentSuccess` → RabbitMQ
3. Order Service consumes event → Updates order status to COMPLETED
4. Order Service publishes `OrderCompleted` → RabbitMQ
5. Chat Service consumes event → Creates chat room for customer + prophet

### Chat Message Flow

1. User sends message via WebSocket → API Gateway
2. API Gateway publishes to `ChatMessageIncoming` queue → RabbitMQ
3. Chat Service consumes message → Validates and saves to database
4. Chat Service publishes to `ChatMessageOutgoing` queue → RabbitMQ
5. API Gateway consumes → Broadcasts via WebSocket to all users in room

## Service Dependencies

```
API Gateway
  └─> All Services (HTTP/REST)

User Management Service
  └─> Firebase (Auth)

Course Service
  └─> User Management Service (gRPC)

Chat Service
  ├─> User Management Service (gRPC)
  └─> Course Service (gRPC)

Order Service
  ├─> Course Service (gRPC)
  └─> Payment Service (gRPC/HTTP)

Payment Service
  (No direct service dependencies)
```

## Message Queue Events

### Queues

- `chat_message_incoming_queue`: Messages from users to be processed
- `chat_message_outgoing_queue`: Messages to be broadcasted to users
- `create_payment_queue`: Order created events
- `update_order_status_queue`: Payment success events
- `update_payment_id_queue`: Payment created events
- `settle_payment_queue`: Settlement events
- `notify_create_payment`: Payment notifications
- `notify_order_completed`: Order completion notifications

## Additional Resources

- [Getting Claims Documentation](../getting_claims.md)
- [Order Service README](../../services/order-service/README.md)
- [Payment Service README](../../services/payment-service/README.md)
- [Main README](../../Readme.md)
