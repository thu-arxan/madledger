pragma solidity ^0.4.0;
 
contract Test{
    uint balance;
    
    function update(uint amount) public returns (address, uint){
        balance += amount;
        balance += 1;
        return (msg.sender, balance);
    }
}