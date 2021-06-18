require('ts-node/register');
const main = require('./src/');


let test = main.parseToPinusProtobuf('./testTreaty');
console.log('server result', JSON.stringify(test, null, 4));