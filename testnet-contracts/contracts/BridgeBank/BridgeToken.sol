pragma solidity ^0.5.0;

import "../OpenZep/ERC20Mintable.sol";
import "../OpenZep/ERC20Detailed.sol";

/**
 * @title BridgeToken
 * @dev Mintable, ERC20 compatible BankToken for use by BridgeBank
 **/

contract BridgeToken is ERC20Mintable, ERC20Detailed {

    constructor(
        string memory _symbol
    )
        public
        ERC20Detailed(
            _symbol,
            _symbol,
            18
        )
    {
        // Intentionally left blank
    }
}