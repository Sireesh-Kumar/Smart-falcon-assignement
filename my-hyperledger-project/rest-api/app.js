const express = require('express');
const { Gateway, Wallets } = require('fabric-network');
const path = require('path');
const fs = require('fs');

const app = express();
app.use(express.json());

const ccpPath = path.resolve(__dirname, 'connection-org1.json');
const walletPath = path.join(process.cwd(), 'wallet');

async function main() {
    try {
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: 'appUser', discovery: { enabled: true, asLocalhost: true } });

        const network = await gateway.getNetwork('mychannel');
        const contract = network.getContract('mycontract');

        app.post('/createAccount', async (req, res) => {
            try {
                const { dealerID, msisdn, mpin, balance, status, transAmount, transType, remarks } = req.body;
                await contract.submitTransaction('CreateAccount', dealerID, msisdn, mpin, balance, status, transAmount, transType, remarks);
                res.send('Account created successfully');
            } catch (error) {
                res.status(500).send(`Failed to create account: ${error}`);
            }
        });

        app.listen(3000, () => {
            console.log('Server is running on port 3000');
        });
    } catch (error) {
        console.error(`Failed to connect to gateway: ${error}`);
        process.exit(1);
    }
}

main();
