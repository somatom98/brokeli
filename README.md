# Brøkeli

<p align="left">
  <img src="web/src/assets/BrokeliLogo.png" height="100" alt="Brøkeli Logo">
  <img src="web/src/assets/BrokeliText.png" height="55" alt="Brøkeli Text">
</p>

Grow Your Green, Smartly

---

## About Brøkeli

Brøkeli is your personal finance companion designed to help you **cultivate healthy financial habits** and watch your wealth grow. With an intuitive interface and powerful tools, Brøkeli makes managing your money simple, insightful, and even enjoyable. Whether you're tracking expenses, setting budgets, or planning for future investments, Brøkeli provides the clarity you need to make smart financial decisions.

---

## Technical stuff

This section provides a technical overview of the Brøkeli project, including its architecture, domains, events, and exposed API endpoints.

### Architecture Overview

Brøkeli follows standard Go project layout conventions, employing Clean Architecture, Domain-Driven Design (DDD), Command Query Responsibility Segregation (CQRS), and Event Sourcing.

- **`cmd/`**: Contains the application entry points.
  - `main.go`: The main executable that wires up and starts the application.
  - `migrate/`: Contains database migration tools/scripts.
- **`internal/`**: Contains private application and library code.
  - `domain/`: Core business logic, separated into bounded contexts (aggregates).
    - `account/`: Account aggregate and related events (e.g., `Opened`, `MoneyDeposited`).
    - `transaction/`: Transaction aggregate and related events (e.g., `MoneySpent`, `MoneyTransfered`).
    - `projections/`: Read models built from the event store (e.g., `accounts` projection).
    - `values/`: Value objects used across domains (e.g., `Currency`, `Entry`).
  - `features/`: Vertical slices containing HTTP handlers/endpoints.
    - `manage_accounts/`: Handlers for account management.
    - `manage_transactions/`: Handlers for transaction recording.
    - `import_transactions/`: Handlers for importing transactions from external sources.
  - `setup/`: Application initialization, dependency injection, and routing wiring.
- **`pkg/`**: Public library code.
  - `event_store/`: Abstractions and implementations (e.g., PostgreSQL) for persisting and subscribing to domain events.
  - `database/`: Database utilities and migration logic.
- **`tests/`**: Integration tests and end-to-end tests.

### Domains & Events

The system is divided into bounded contexts (Domains), each with its own Aggregate and set of Events.

#### 1. Account Domain

Manages the lifecycle and core properties of financial accounts.

- **Events**:
  - `AccountOpened`: A new account was created.
  - `AccountNameUpdated`: An account name was changed.
  - `MoneyDeposited`: Money was added to an account balance.
  - `MoneyWithdrawn`: Money was removed from an account balance.

#### 2. Transaction Domain

Manages the movement of money between accounts or external entities.

- **Events**:
  - `MoneySpent`: An expense was recorded.
  - `MoneyReceived`: Income was recorded.
  - `MoneyTransfered`: Money was moved between two internal accounts.
  - `ReimbursementReceived`: A reimbursement was received for a specific transaction.
  - `ExpectedReimbursementSet`: Marked a transaction as expecting a reimbursement.

### Projections

#### Accounts Projection

Maintains the current state of all accounts, calculating balances by aggregating relevant events from both **Account** and **Transaction** domains.

### API Endpoints

#### Manage Accounts

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/api/accounts` | List all accounts with current balances. |
| `GET` | `/api/accounts/{id}/balances` | Get balances for a specific account. |
| `GET` | `/api/balances` | Get all account balances. |
| `POST` | `/api/accounts` | Create a new account. |
| `POST` | `/api/accounts/{id}/deposits` | Record a deposit into an account. |
| `POST` | `/api/accounts/{id}/withdrawals` | Record a withdrawal from an account. |

#### Manage Transactions

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/expenses` | Register a new expense (money spent). |
| `POST` | `/api/incomes` | Register a new income (money received). |
| `POST` | `/api/transfers` | Register a transfer between accounts. |
| `POST` | `/api/{transaction_id}/reimbursement` | Record a reimbursement for a transaction. |
| `POST` | `/api/{transaction_id}/expected-reimbursements` | Set expected reimbursement amount. |

#### Import Transactions

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/import-transactions` | Import transactions from an external source (e.g., CSV). |
