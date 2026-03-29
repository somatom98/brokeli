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
  - `lunar-converter/`: Tool to convert transactions from Lunar.
- **`internal/`**: Contains private application and library code.
  - `domain/`: Core business logic, separated into bounded contexts (aggregates).
    - `account/`: Account aggregate and related events (e.g., `Opened`, `MoneyDeposited`).
    - `transaction/`: Transaction aggregate and related events (e.g., `MoneySpent`, `MoneyTransfered`).
    - `budget/`: User budgets and limits management.
    - `projections/`: Read models built from the event store (e.g., `accounts`, `transactions`, `balance_updates` projections).
    - `values/`: Value objects used across domains (e.g., `Currency`, `Entry`).
  - `features/`: Vertical slices containing HTTP handlers/endpoints.
    - `manage_accounts/`: Handlers for account management.
    - `manage_transactions/`: Handlers for transaction recording and querying.
    - `manage_budgets/`: Handlers for budget management.
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

#### 3. Budget Domain

Manages user-defined budgets and spending limits based on transaction categories.

### Projections

#### Accounts Projection

Maintains the current state of all accounts, calculating balances by aggregating relevant events from both **Account** and **Transaction** domains.

#### Balance Updates Projection

Maintains a historical series of account balances over time.

#### Transactions Projection

Maintains a queryable read model of all recorded transactions.

### API Endpoints

#### Manage Accounts

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/api/accounts` | List all accounts with current balances. |
| `GET` | `/api/accounts/{id}/balances` | Get balances for a specific account. |
| `GET` | `/api/accounts/{id}/distributions` | Get distributions for a specific account. |
| `GET` | `/api/balances` | Get all account balances. |
| `POST` | `/api/accounts` | Create a new account. |
| `POST` | `/api/accounts/{id}/deposits` | Record a deposit into an account. |
| `POST` | `/api/accounts/{id}/withdrawals` | Record a withdrawal from an account. |

#### Manage Transactions

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/api/transactions` | List and query transactions. |
| `POST` | `/api/expenses` | Register a new expense (money spent). |
| `POST` | `/api/incomes` | Register a new income (money received). |
| `POST` | `/api/transfers` | Register a transfer between accounts. |
| `POST` | `/api/{transaction_id}/reimbursement` | Record a reimbursement for a transaction. |
| `POST` | `/api/{transaction_id}/expected-reimbursements` | Set expected reimbursement amount. |

#### Manage Budgets

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/api/budgets` | List all budgets. |
| `GET` | `/api/budgets/categories` | Get transaction categories for budgeting. |
| `POST` | `/api/budgets` | Save or update a budget. |
| `DELETE` | `/api/budgets/{id}` | Delete a specific budget. |

#### Import Transactions

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/import-transactions` | Import transactions from an external source (e.g., CSV). |

## Future Improvements & Roadmap

Brøkeli is a living project with a long-term vision to become a comprehensive, AI-enhanced financial platform. Below are the key areas targeted for future development.

### 💱 Currency Conversion & Multi-Currency Support

Enhance the current system to handle accounts and transactions in any currency, providing a seamless overview in a primary "Base Currency."

- **Automatic Exchange Rates**: Integration with external APIs (e.g., Fixer.io, Open Exchange Rates) for daily updates.
- **Historical Tracking**: Ability to see balances over time adjusted for historical exchange rates.
- **Unified Reporting**: View consolidated net worth and spending reports across all currencies.

### 📈 Investment & Asset Tracking

Expand beyond liquid accounts to track investments, providing a holistic view of financial health.

- **Ticker Tracking**: Register stocks, ETFs, and crypto holdings with real-time price updates.
- **Portfolio Analytics**: Track profit/loss, dividend yields, and asset allocation across different sectors.
- **Non-Financial Assets**: Support for tracking real estate, vehicles, or other significant assets.

### 🛠️ Refined Transaction Management

Improve the core transaction workflow with more flexibility and better tracking of complex movements.

- **Reimbursement Linkage**: Allow referencing the original expense in a reimbursement, so that net spending for a specific transaction can be updated retroactively, regardless of timing.
- **Granular Editing**: Provide the ability to modify all fields of an existing transaction (e.g., amount, date, category, tags) with a full audit trail maintained in the event store.
- **Splitting Transactions**: Ability to split a single receipt into multiple categories or accounts (e.g., a supermarket trip including both groceries and household items).

### 🏷️ Enhanced Categorization & AI Insights

Leverage AI to move from manual tracking to proactive financial coaching.

- **Auto-Categorization**: Use machine learning to suggest categories for imported transactions based on historical data.
- **Natural Language Querying (MCP/A2A)**: Ask questions like "How much did I spend on groceries in London last month?" directly in the UI.
