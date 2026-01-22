# Brøkeli

![Logo](https://github.com/user-attachments/assets/84d1a33c-d30c-4727-bb97-c1da0976bd9e)

Grow Your Green, Smartly

---

## About Brøkeli

Brøkeli is your personal finance companion designed to help you **cultivate healthy financial habits** and watch your wealth grow. With an intuitive interface and powerful tools, Brøkeli makes managing your money simple, insightful, and even enjoyable. Whether you're tracking expenses, setting budgets, or planning for future investments, Brøkeli provides the clarity you need to make smart financial decisions.

---

## Technical stuff

This section provides a technical overview of the Brøkeli project, including its architecture, domains, events, and exposed API endpoints.

### Architecture Overview

Brøkeli follows a **CQRS (Command Query Responsibility Segregation)** and **Event Sourcing** inspired architecture.

The data flow is generally as follows:

1. **Write Side (Commands)**:
    * **HTTP Handler**: Receives a request (e.g., "Create Account").
    * **Dispatcher**: Dispatches the command to the appropriate Domain Aggregate.
    * **Aggregate**: Validates the command and produces **Events** (e.g., `AccountCreated`).
    * **Event Store**: Persists the events.
2. **Read Side (Queries)**:
    * **Projections**: Subscribe to the Event Store and build optimized read models (e.g., current account balances) in memory.
    * **HTTP Handler**: Queries the Projection to return data to the user.

### Domains & Events

The system is divided into bounded contexts (Domains), each with its own Aggregate and set of Events.

#### 1. Account Domain

Manages the lifecycle of financial accounts.

* **Events**:
  * `Created`: An account was created.
  * `MoneyDeposited`: Initial balance or direct deposit.
  * `AccountClosed`: The account has been archived/closed.

#### 2. Transaction Domain

Manages the movement of money between accounts or external entities.

* **Events**:
  * `MoneySpent`: An expense was recorded.
  * `MoneyReceived`: Income was recorded.
  * `MoneyTransfered`: Money was moved between two internal accounts.
  * `ReimbursementReceived`: A reimbursement was received for a specific transaction.
  * `ExpectedReimbursementSet`: Marked a transaction as expecting a reimbursement.

### Projections

#### Accounts Projection

Maintains the current state of all accounts, calculating balances by aggregating relevant events from both the **Account** and **Transaction** domains.

### API Endpoints

#### Manage Accounts

Operations related to account management.

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/api/accounts` | List all accounts with current balances. |
| `POST` | `/api/accounts` | Create a new account. |
| `DELETE` | `/api/accounts/{id}` | Close an account. |

#### Manage Transactions

Operations for recording financial transactions.

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/expenses` | Register a new expense (money spent). |
| `POST` | `/api/incomes` | Register a new income (money received). |
| `POST` | `/api/transfers` | Register a transfer between accounts. |
| `POST` | `/api/{transaction_id}/reimbursement` | Record a reimbursement for a transaction. |
| `POST` | `/api/{transaction_id}/expected-reimbursements` | Set expected reimbursement amount. |

