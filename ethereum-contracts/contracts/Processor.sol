pragma solidity ^0.5.0;

import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "openzeppelin-solidity/contracts/token/ERC20/ERC20.sol";

  /*
   *  @title: Processor
   *  @dev: Processes requests for deposit locking and unlocking by
   *        storing a deposit's information then relaying the funds
   *        the original sender.
   */
contract Processor {

    using SafeMath for uint256;

    /*
    * @dev:  Deposit struct to store information.
    */    
    struct Deposit {
        address payable sender;
        bytes recipient;
        address token;
        uint256 amount;
        uint256 nonce;
        bool locked;
    }

    uint256 public nonce;
    mapping(bytes32 => Deposit) private deposits;

    /*
    * @dev: Constructor, initalizes deposit count.
    */
    constructor() 
        public
    {
        nonce = 0;
    }

    modifier onlySender(bytes32 _id) {
        require(
            msg.sender == deposits[_id].sender,
            'Must be the original sender.'
        );
        _;
    }

    modifier canDeliver(bytes32 _id) {
        if(deposits[_id].token == address(0)) {
            require(
                address(this).balance >= deposits[_id].amount,
                'Insufficient ethereum balance for delivery.'
            );
        } else {
            require(
                ERC20(deposits[_id].token).balanceOf(address(this)) >= deposits[_id].amount,
                'Insufficient ERC20 token balance for delivery.'
            );            
        }
        _;
    }
  
    modifier availableNonce() {
        require(
            nonce + 1 > nonce,
            'No available nonces.'
        );
        _;
    }

    /*
    * @dev: Creates an item with a unique id.
    *
    * @param _sender: The sender's ethereum address.
    * @param _recipient: The intended recipient's cosmos address.
    * @param _token: The currency type, either erc20 or ethereum.
    * @param _amount: The amount of erc20 tokens/ ethereum (in wei) to be itemized.
    * @return: The newly created item's unique id.
    */
    function create(
        address payable _sender,
        bytes memory _recipient,
        address _token,
        uint256 _amount
    )
        internal
        returns(bytes32)
    {
        nonce++;

        bytes32 depositKey = keccak256(
            abi.encodePacked(
                _sender,
                _recipient,
                _token,
                _amount,
                nonce
            )
        );
        
        deposits[depositKey] = Deposit(
            _sender,
            _recipient,
            _token,
            _amount,
            nonce,
            true
        );

        return depositKey;
    }

    /*
    * @dev: Completes the deposit by sending the funds to the
    *       original sender and unlocking the deposit.
    *
    * @param _id: The deposit to be completed.
    */
    function complete(
        bytes32 _id
    )
        internal
        canDeliver(_id)
        returns(address payable, address, uint256, uint256)
    {
        require(isLocked(_id));

        //Get locked deposit's attributes for return
        address payable sender = deposits[_id].sender;
        address token = deposits[_id].token;
        uint256 amount = deposits[_id].amount;
        uint256 uniqueNonce = deposits[_id].nonce;

        //Update lock status
        deposits[_id].locked = false;

        //Transfers based on token address type
        if (token == address(0)) {
          sender.transfer(amount);
        } else {
          require(ERC20(token).transfer(sender, amount));
        }       

        return(sender, token, amount, uniqueNonce);
    }

    /*
    * @dev: Checks the current nonce.
    *
    * @return: The current nonce.
    */
    function getNonce()
        internal
        view
        returns(uint256)
    {
        return nonce;
    }

    /*
    * @dev: Checks if an individual deposit exists.
    *
    * @param _id: The unique deposit's id.
    * @return: Boolean indicating if the deposit exists in memory.
    */
    function isLocked(
        bytes32 _id
    )
        internal 
        view
        returns(bool)
    {
        return(deposits[_id].locked);
    }

    /*
    * @dev: Gets an deposit's information
    *
    * @param _Id: The deposit containing the desired information.
    * @return: Sender's address.
    * @return: Recipient's address in bytes.
    * @return: Token address.
    * @return: Amount of ethereum/erc20 in the deposit.
    * @return: Unique nonce of the deposit.
    */
    function getDeposit(
        bytes32 _id
    )
        internal 
        view
        returns(address payable, bytes memory, address, uint256, uint256)
    {
        Deposit memory deposit = deposits[_id];

        return(
            deposit.sender,
            deposit.recipient,
            deposit.token,
            deposit.amount,
            deposit.nonce
        );
    }
}
