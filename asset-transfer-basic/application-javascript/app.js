/*
 * Copyright IBM Corp. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const FabricCAServices = require("fabric-ca-client");
const { Gateway, Wallets } = require('fabric-network');
const path = require('path');
const { buildCAClient, registerAndEnrollUser, enrollAdmin } = require('../../test-application/javascript/CAUtil.js');
const { buildCCPOrg1, buildCCPOrg2, buildWallet } = require('../../test-application/javascript/AppUtil.js');

const readline = require("node:readline")

const channelName = 'mychannel';
const chaincodeName = 'basic';
const mspOrg1 = 'Org1MSP';
const mspOrg2 = 'Org2MSP';
const walletPath = path.join(__dirname, 'wallet');
const org1UserId = 'appUser';

function prettyJSONString(inputString) {
	return JSON.stringify(JSON.parse(inputString), null, 2);
}

function readInput(prompt) {
	const rl = readline.createInterface({
		input: process.stdin,
		output: process.stdout,
	});

	return new Promise((resolve) => {
		rl.question(prompt, (input) => {
			rl.close();
			resolve(input);
		})
	})
}

function getUnixFormatTime(date = new Date()) {
	const weekdays = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
	const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun',
		'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

	const weekday = weekdays[date.getDay()];
	const month = months[date.getMonth()];
	const day = String(date.getDate()).padStart(2, ' '); // space pad like Go
	const hours = String(date.getHours()).padStart(2, '0');
	const minutes = String(date.getMinutes()).padStart(2, '0');
	const seconds = String(date.getSeconds()).padStart(2, '0');
	const year = date.getFullYear();

	// Extract short timezone abbreviation from toTimeString, e.g. "(CEST)"
	// const tzMatch = date.toTimeString().match(/\(([^)]+)\)/);
	// const tzAbbr = tzMatch ? tzMatch[1] : 'UTC';
	const tzAbbr = "CEST";

	return `${weekday} ${month} ${day} ${hours}:${minutes}:${seconds} ${tzAbbr} ${year}`;
}

async function main() {
	let menuInput = "";
	let ccp, caClient;
	let wallet;
	let gateway;
	do {
		menuInput = await readInput(`Choose option:
1. Login
2. Query chaincode
3. Invoke chaincode
0. Quit
Choice: `);

		switch (menuInput) {
			case "1":
				const loginInput = await readInput(`Choose account:
1. Org1
2. Org2
Choice: `);

				switch (loginInput) {
					case "1":
						try {
							ccp = buildCCPOrg1();
							caClient = buildCAClient(FabricCAServices, ccp, 'ca.org1.example.com');
							wallet = await buildWallet(Wallets, walletPath);
							await enrollAdmin(caClient, wallet, mspOrg1);
							await registerAndEnrollUser(caClient, wallet, mspOrg1, org1UserId, 'org1.department1');
						} catch (error) {
							console.error(`******** FAILED to run the application: ${error}`);
						}
						break;
					case "2":
						try {
							ccp = buildCCPOrg2();
							caClient = buildCAClient(FabricCAServices, ccp, 'ca.org2.example.com');
							wallet = await buildWallet(Wallets, walletPath);
							await enrollAdmin(caClient, wallet, mspOrg2);
							await registerAndEnrollUser(caClient, wallet, mspOrg2, org1UserId, 'org1.department1');
						} catch (error) {
							console.error(`******** FAILED to run the application: ${error}`);
						}
						break;
				}
				break;
			case "2":
				const id = await readInput(`Choose query parameters:
id: `);
				const name = await readInput(`name: `);
				const merchantType = await readInput(`merchant type: `);
				const price = await readInput(`price: `);

				gateway = new Gateway();
				try {
					await gateway.connect(ccp, {
						wallet,
						identity: org1UserId,
						discovery: { enabled: true, asLocalhost: true }
					});

					const network = await gateway.getNetwork(channelName);

					const contract = network.getContract(chaincodeName);

					console.log('\n--> Evaluate Transaction: FindProduct')
					let result = await contract.evaluateTransaction('FindProduct', id, name, merchantType, price);
					console.log(`*** Result: ${prettyJSONString(result.toString())}`);
				} finally {
					gateway.disconnect();
				}
				break;
			case "3":
				gateway = new Gateway();
				try {
					await gateway.connect(ccp, {
						wallet,
						identity: org1UserId,
						discovery: { enabled: true, asLocalhost: true }
					});

					const network = await gateway.getNetwork(channelName);

					const contract = network.getContract(chaincodeName);

					const invokeInput = await readInput(`Choose a function to invoke:
1. Create merchant
2. Add product to merchant
3. Create users
4. Buy product
5. Add funds
Choice: `);
					switch (invokeInput) {
						case "1":
							const id = await readInput(`\tid: `);
							const pib = await readInput(`\tpib: `);
							const type = await readInput(`\ttype (Wholesale, Supermarket, Service): `);
							const balance = await readInput(`\tbalance: `);

							console.log('\n--> Submit Transaction: CreateMerchant')
							await contract.submitTransaction('CreateMerchant', id, pib, type, balance);
							break;

						case "2":
							const merchantId = await readInput(`\t merchantId: `);
							const productId = await readInput(`\t productId: `);
							console.log('\n --> Submit Transaction: AddProductToMerchant');
							await contract.submitTransaction('AddProductToMerchant', merchantId, productId);
							break;

						case "3":
							const users = [];
							while (true) {
								const usersInput = await readInput(`Add more users (y/n)? `);
								if (usersInput !== "y") {
									break;
								}

								const user = {
									"ID": await readInput(`\tid: `),
									"Name": await readInput(`\tname: `),
									"Surname": await readInput(`\tsurname: `),
									"Email": await readInput(`\temail: `),
									"Balance": Number(await readInput(`\tbalance: `)),
									"Bills": [],
								}

								users.push(user);
							}

							console.log('\n--> Submit Transaction: CreateUsers');
							await contract.submitTransaction("CreateUsers", JSON.stringify(users))
							break;

						case "4":
							const userId = await readInput(`\t userId: `);
							const productId1 = await readInput(`\t productId: `);

							console.log('\n--> Submit Transaction: BuyProduct');
							contract.submitTransaction("BuyProduct", userId, productId1, getUnixFormatTime());
							break;

						case "5":
							const userOrMerchantId = await readInput(`\tuser or merchant id: `);
							const amount = await readInput(`\tamount: `);

							console.log('\n--> Submit Transaction: AddFunds');
							contract.submitTransaction("AddFunds", userOrMerchantId, amount);
							break;
					}

				} finally {
					gateway.disconnect();
				}
				break;
		}

	} while (menuInput !== "0")
}

main();
