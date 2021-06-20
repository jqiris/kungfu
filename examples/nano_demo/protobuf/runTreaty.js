const main = require('pinus-parse-interface');

// const test = main.parseToPinusProtobuf('./testTreaty');
// console.log('result',JSON.stringify(test,null,4));
main.parseAndWrite('./testTreaty', './protos.json');

