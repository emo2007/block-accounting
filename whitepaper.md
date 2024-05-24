### Chainlink Multisig Transaction Contract
- Overview:
- This contract is used to manage transactions within an organization that requires multiple confirmations. It leverages a multisig mechanism to ensure that actions within the contract are approved by multiple parties before they are executed.

### Key Features:
- Owners Array: Initial deployment includes an array of addresses (owners[]), representing individuals linked to the organization. And the number of needed confirmations (uint <= owners.length);
- Confirmation Tracking: The contract tracks the number of confirmations required and received for each transaction.
- Transaction States: Each transaction can be in one of three states:
- Submitted: The initial state when a transaction is proposed.
- Confirmed: After the transaction receives the required number of confirmations.
- Executed: The final state after the transaction is executed.
- Deployment: Deployed using a main wallet which is considered the primary interface for administrative interactions.
### Workflow:
- Submission: A user proposes a transaction, which is recorded in the 'Submitted' state.
- Confirmation: As required confirmations are collected, the transaction transitions to the 'Confirmed' state.
- Execution: Once confirmed, the transaction can be executed by anyone, transitioning to the 'Executed' state.

### Chainlink Payroll Contract
### Overview:
This contract manages the payroll system, allowing salaries to be set in USD and then paid out in ETH based on the current exchange rate provided via Chainlink oracles.

### Key Features:
- Authorized Wallet: Only a specified wallet can execute payouts, set during contract deployment.
- Salary Management: Salaries are set in USD for each employee, needing confirmation through the multisig mechanism before execution.
- Currency Conversion: Utilizes Chainlink to fetch real-time ETH/USD prices to calculate the payout amount in ETH.
### Workflow:
- Set Salary: Propose salaries in USD which are confirmed and executed via multisig.
- Payout: On payroll day, the contract calculates the equivalent ETH for each employee's USD salary and transfers it. You have to deposit funds before calling this function.

### Chainlink Licensing Contract
### Overview:
- This contract manages licensing agreements by distributing funds based on predefined shares, after fetching and storing relevant data from Chainlink data feeds.

### Key Features:
- Data Request: Requests data, like total payout amounts, from a Chainlink data feed.
- Share and Owner Management: Stores shares and owner addresses, setting how distributions are handled.
- Multi-layer Confirmation: Deployment and critical functions require multisig confirmation. CREATE2 OPCODE is used on the multisig side.
### Workflow:
- Data Fetching: Requests data from Chainlink Custom Data Feed and stores it.
- Payout Setup: Before executing payouts, set the payroll contract address.
- Distribution: Distributes funds according to shares among the owners.


### Chainlink Agreement Contract
### Overview:
- Similar in functionality to the licensing contract, this contract fetches and evaluates boolean data points to determine outcomes of agreements.

### Key Features:
- Boolean Data Handling: Manages agreement validations based on true/false responses from Chainlink data feeds.
### Workflow:
- Data Request: Fetches a boolean value determining the agreement's state.
- Outcome Execution: Executes actions based on the true/false outcome, similar to licensing contract operations.


These contracts collectively form a robust framework using Chainlink oracles and Ethereum blockchain technology to ensure secure, transparent, and decentralized transaction management within an organization. Each contract's deployment and operations are safeguarded by multisig processes to maintain organizational control and integrity.

During development, we utilized the Polygon Amoy network because it simplifies handling multi calls and offers lower costs and faster transactions. This made Amoy an optimal choice for implementing these contracts.

