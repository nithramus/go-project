const net = require('net')

const client = net.createConnection({host: "wss://chatws29.stream.highwebmedia.com"}, () => {
      // 'connect' listener
  console.log('connected to server!');
  client.write('world!\r\n');
})
