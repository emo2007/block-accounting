// SPDX-License-Identifier: MIT
// 0x74f11486DB0FCAA2dCDE0aEB477e1F37fCAa510A

pragma solidity ^0.8.19;
// The wallet owners can
// submit a transaction
// approve and revoke approval of pending transactions
// anyone can execute a transaction after enough owners has approved it.
contract MultiSigWallet {
    event Deposit(address indexed sender, uint amount, uint balance);
    event SubmitTransaction(
        address indexed owner,
        uint indexed txIndex,
        address indexed to,
        uint value,
        bytes data
    );

    event ConfirmTransaction(address indexed owner, uint indexed txIndex);
    event RevokeConfirmation(address indexed owner, uint indexed txIndex);
    event ExecuteTransaction(address indexed owner, uint indexed txIndex, address indexed to);
    event ExecuteTransactionFailed(address indexed owner, uint indexed txIndex, string reason);
    event ContractDeployed(address indexed contractAddress);


    address[] public owners;

    mapping(address => bool) public isOwner;

    uint public numConfirmationsRequired;

    struct Transaction {
        address to;
        uint value;
        bytes data;
        bool executed;
        uint numConfirmations;
    }

    mapping(uint => mapping(address => bool)) public isConfirmed;

    Transaction[] public transactions;

    modifier onlyOwner() {
        require(isOwner[msg.sender], 'not owner');
        _;
    }

    modifier txExists(uint _txIndex) {
        require(_txIndex < transactions.length, 'tx does not exist');
        _;
    }

    modifier notConfirmed(uint _txIndex) {
        require(!isConfirmed[_txIndex][msg.sender], 'tx already confirmed');
        _;
    }

    modifier notExecuted(uint _txIndex) {
        require(!transactions[_txIndex].executed, 'tx already confirmed');
        _;
    }

    constructor(address[] memory _owners, uint _numConfirmationsRequired) {
        require(_owners.length > 0, 'owners required');
        require(
            _numConfirmationsRequired > 0 &&
                _numConfirmationsRequired <= _owners.length,
            'invalid number of required confirmations'
        );
        for (uint i = 0; i < _owners.length; i++) {
            address owner = _owners[i];
            require(owner != address(0), 'invalid owner');
            require(!isOwner[owner], 'owner not unique');
            isOwner[owner] = true;
            owners.push(owner);
        }
        numConfirmationsRequired = _numConfirmationsRequired;
    }

    receive() external payable {
        emit Deposit(msg.sender, msg.value, address(this).balance);
    }

    function submitTransaction(
        address _to,
        uint _value,
        bytes memory _data
    ) public onlyOwner {
        uint txIndex = transactions.length;
        transactions.push(
            Transaction({
                to: _to,
                value: _value,
                data: _data,
                executed: false,
                numConfirmations: 0
            })
        );
        emit SubmitTransaction(msg.sender, txIndex, _to, _value, _data);
    }

    function confirmTransaction(
        uint _txIndex
    )
        public
        onlyOwner
        txExists(_txIndex)
        notExecuted(_txIndex)
        notConfirmed(_txIndex)
    {
        Transaction storage transaction = transactions[_txIndex];
        transaction.numConfirmations += 1;
        isConfirmed[_txIndex][msg.sender] = true;
        emit ConfirmTransaction(msg.sender, _txIndex);
    }

    function executeTransaction(uint _txIndex)
    public
    onlyOwner
    txExists(_txIndex)
    notExecuted(_txIndex)
    {
        Transaction storage transaction = transactions[_txIndex];
        require(
            transaction.numConfirmations >= numConfirmationsRequired,
            "cannot execute tx"
        );


        (bool success, bytes memory returnData) = transaction.to.call{value: transaction.value}(transaction.data);
        if (success) {
            transaction.executed = true;
            emit ExecuteTransaction(msg.sender, _txIndex, transaction.to);
            removeTransaction(_txIndex);
        } else {
            // Get the revert reason and emit it
            if (returnData.length > 0) {
                // The call reverted with a message
                assembly {
                    let returndata_size := mload(returnData)
                    revert(add(32, returnData), returndata_size)
                }
            } else {
                // The call reverted without a message
                emit ExecuteTransactionFailed(msg.sender, _txIndex, "Transaction failed without a reason");
            }
        }
    }

    function executeDeployTransaction(uint _txIndex, uint256 _salt) public onlyOwner txExists(_txIndex) notExecuted(_txIndex) {
        Transaction storage transaction = transactions[_txIndex];
        require(
            transaction.numConfirmations >= numConfirmationsRequired,
            "cannot execute tx"
        );

        address deployedAddress;

        bytes memory bytecode = transaction.data;

        // Assembly to deploy contract using CREATE2
        assembly {
            deployedAddress :=
            create2(
                callvalue(), // wei sent with current call
                // Actual code starts after skipping the first 32 bytes
                add(bytecode, 0x20),
                mload(bytecode), // Load the size of code contained in the first 32 bytes
                _salt // Salt from function arguments
            )

            if iszero(extcodesize(deployedAddress)) { revert(0, 0) }
        }

        require(deployedAddress != address(0), "Failed to deploy contract");
        transaction.executed = true;
        emit ExecuteTransaction(msg.sender, _txIndex, deployedAddress);
        emit ContractDeployed(deployedAddress);
        removeTransaction(_txIndex);
    }


    function revokeConfirmation(
        uint _txIndex
    ) public onlyOwner txExists(_txIndex) notExecuted(_txIndex) {
        Transaction storage transaction = transactions[_txIndex];
        require(isConfirmed[_txIndex][msg.sender], 'tx not confirmed');
        transaction.numConfirmations -= 1;
        isConfirmed[_txIndex][msg.sender] = false;

        emit RevokeConfirmation(msg.sender, _txIndex);
    }

    function removeTransaction(uint _txIndex) public onlyOwner {
        require(_txIndex < transactions.length, "tx does not exist");
        delete transactions[_txIndex];
    }

    function getOwners() public view returns (address[] memory) {
        return owners;
    }

    function getTransactionCount() public view returns (uint) {
        return transactions.length;
    }

    function getTransaction(
        uint _txIndex
    )
        public
        view
        returns (
            address to,
            uint value,
            bytes memory data,
            bool executed,
            uint numConfirmations
        )
    {
        Transaction storage transaction = transactions[_txIndex];
        return (
            transaction.to,
            transaction.value,
            transaction.data,
            transaction.executed,
            transaction.numConfirmations
        );
    }
}
