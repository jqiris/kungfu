var protobuf = require('protocol-buffers')
var fs = require('fs')

// protobuf.toJS() takes the same arguments as protobuf()
var js = protobuf.toJS(fs.readFileSync('treaty.proto'))
fs.writeFileSync('treaty.js', js)