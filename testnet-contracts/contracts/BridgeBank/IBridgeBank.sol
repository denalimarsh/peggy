pragma solidity ^0.5.0;

/*
    Bridge Bank interface
*/
interface IBridgeBank {

    function createNewBridgeToken(
        string calldata _symbol
    ) external returns(address);

     function mintBridgeTokens(
        bytes calldata _cosmosSender,
        address payable _intendedRecipient,
        address _bridgeTokenAddress,
        string calldata _symbol,
        uint256 _amount
    ) external;

    function lock(
        bytes calldata _recipient,
        address _token,
        uint256 _amount
    ) external payable;

     function unlock(
        address payable _recipient,
        address _token,
        string calldata _symbol,
        uint256 _amount
    ) external;

    function getCosmosDepositStatus(
        bytes32 _id
    ) external view returns(bool);

    function viewCosmosDeposit(
        bytes32 _id
    ) external view returns(bytes memory, address payable, address, uint256);

    function() external payable;
}