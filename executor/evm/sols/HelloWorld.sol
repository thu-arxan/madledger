pragma solidity ^0.4.0;
 
contract HelloWorld{
    uint public balance;
    
    function update(uint amount) public returns (address, uint){
        balance += amount;
        return (msg.sender, balance);
    }
}