const TestProcessor = artifacts.require('TestProcessor');
const TestToken = artifacts.require('TestToken');

const Web3Utils = require('web3-utils');
const EVMRevert = 'revert';
const BigNumber = web3.BigNumber;

require('chai')
  .use(require('chai-as-promised'))
  .use(require('chai-bignumber')(BigNumber))
  .should();

contract('TestProcessor', function (accounts) {

  const userOne = accounts[1];
  const userTwo = accounts[2];
  const userThree = accounts[3];

  describe('Processor contract deployment', function() {

    beforeEach(async function() {
      this.processor = await TestProcessor.new();
    });

    it('should deploy the processor with the correct parameters', async function () {
      this.processor.should.exist;
      
      const nonce = Number(await this.processor.nonce());
      nonce.should.be.bignumber.equal(0);
    });

  });

  describe('Deposit creation', function() {

    beforeEach(async function() {
      this.processor = await TestProcessor.new();
      this.recipient = web3.utils.bytesToHex(['20bytestring'])
      this.amount = 250;

      //Load user account with tokens for testing
      this.token = await TestToken.new();
      await this.token.mint(userOne, 1000, { from: accounts[0] }).should.be.fulfilled;
    });

 
    it('should allow for the creation of deposits', async function () {
      await this.processor.callCreate(userOne, this.recipient, this.token.address, this.amount).should.be.fulfilled;
    });

    it('should generate unique deposit id\'s for a created deposit', async function () {
      //Simulate sha3 hash to get deposit's expected id
      const expectedId = Web3Utils.soliditySha3(
        {t: 'address payable', v: userOne},
        {t: 'bytes', v: this.recipient},
        {t: 'address', v: this.token.address},
        {t: 'int256', v:this.amount},
        {t: 'int256', v:1});

      //Get the deposit's id if it were to be created
      const id = await this.processor.callCreate.call(userOne, this.recipient, this.token.address, this.amount);
      id.should.be.equal(expectedId);
    });

    it('should allow access to an deposit\'s information given it\'s unique id', async function () {
      const id = await this.processor.callCreate.call(userOne, this.recipient, this.token.address, this.amount);
      await this.processor.callCreate(userOne, this.recipient, this.token.address, this.amount);

      //Attempt to get an deposit's information
      await this.processor.callGetDeposit(id).should.be.fulfilled;
    });

    it('should correctly identify the existence of deposits in memory', async function () {
      //Get the deposit's expected id then lock funds
      const id = await this.processor.callCreate.call(userOne, this.recipient, this.token.address, this.amount);
      await this.processor.callCreate(userOne, this.recipient, this.token.address, this.amount).should.be.fulfilled;
      
      //Check if deposit has been created and locked
      const locked = await this.processor.callIsLocked(id);
      locked.should.be.equal(true);
    });

    it('should store deposits with the correct parameters', async function () {
      //Create the deposit and store its id
      const id = await this.processor.callCreate.call(userOne, this.recipient, this.token.address, this.amount);
      await this.processor.callCreate(userOne, this.recipient, this.token.address, this.amount);

      //Get the deposit's information
      const depositInfo = await this.processor.callGetDeposit(id);

      //Parse each attribute
      const sender = depositInfo[0];
      const receiver = depositInfo[1];
      const token = depositInfo[2];
      const amount = Number(depositInfo[3]);
      const nonce = Number(depositInfo[4]);

      //Confirm that each attribute is correct
      sender.should.be.equal(userOne);
      receiver.should.be.equal(this.recipient);
      token.should.be.equal(this.token.address);
      amount.should.be.bignumber.equal(this.amount);
      nonce.should.be.bignumber.equal(1);
    });

  });

  describe('Deposit completion', function() {

    beforeEach(async function() {
      this.processor = await TestProcessor.new();
      this.weiAmount = web3.utils.toWei("0.25", "ether");
      this.recipient = web3.utils.bytesToHex(['20bytestring'])
      this.ethereumToken = '0x0000000000000000000000000000000000000000';

      //Load contract with ethereum so it can complete deposits
      await this.processor.send(web3.utils.toWei("1", "ether"), { from: accounts[0]}).should.be.fulfilled;

      this.depositId = await this.processor.callCreate.call(userOne, this.recipient, this.ethereumToken, this.weiAmount);
      await this.processor.callCreate(userOne, this.recipient, this.ethereumToken, this.weiAmount);
    });

    it('should not allow for the completion of deposits whose value exceeds the contract\'s balance', async function () {
      //Create an deposit with an overlimit amount
      const overlimitAmount = web3.utils.toWei("1.25", "ether");
      const id = await this.processor.callCreate.call(userOne, this.recipient, this.ethereumToken, overlimitAmount);
      await this.processor.callCreate(userOne, this.recipient, this.ethereumToken, overlimitAmount);

      //Attempt to complete the deposit
      await this.processor.callComplete(id).should.be.rejectedWith(EVMRevert);
    });

    it('should not allow for the completion of non-deposits', async function () {
      //Generate a false deposit id
      const fakeId = Web3Utils.soliditySha3(
        {t: 'address payable', v: userOne},
        {t: 'bytes', v: this.recipient},
        {t: 'address', v: this.ethereumToken},
        {t: 'int256', v:12},
        {t: 'int256', v:1});

      await this.processor.callComplete(fakeId).should.be.rejectedWith(EVMRevert);
    
    });

    it('should not allow for the completion of an deposit that has already been completed', async function () {
      //Complete the deposit
      await this.processor.callComplete(this.depositId).should.be.fulfilled;

      //Attempt to complete the deposit again
      await this.processor.callComplete(this.depositId).should.be.rejectedWith(EVMRevert);
    });

    it('should allow for an deposit to be completed', async function () {
      await this.processor.callComplete(this.depositId).should.be.fulfilled;
    });

    it('should update lock status of deposits upon completion', async function () {
      //Confirm that the deposit is active
      const startingLockStatus = await this.processor.callIsLocked(this.depositId);
      startingLockStatus.should.be.equal(true);

      //Complete the deposit
      await this.processor.callComplete(this.depositId).should.be.fulfilled;

      //Check if the deposit still exists
      const completedDeposit = await this.processor.callIsLocked(this.depositId);
      completedDeposit.should.be.equal(false);
    });

    it('should correctly transfer deposited funds to the original sender', async function () {
      //Get prior balances of user and peggy contract
      const beforeUserBalance = Number(await web3.eth.getBalance(userOne));
      const beforeContractBalance = Number(await web3.eth.getBalance(this.processor.address));

      await this.processor.callComplete(this.depositId).should.be.fulfilled;

      //Get balances after completion
      const afterUserBalance = Number(await web3.eth.getBalance(userOne));
      const afterContractBalance = Number(await web3.eth.getBalance(this.processor.address));

      //Expected balances 
      afterUserBalance.should.be.bignumber.equal(beforeUserBalance + Number(this.weiAmount)); 
      afterContractBalance.should.be.bignumber.equal(beforeContractBalance - Number(this.weiAmount));
    });

  });

});