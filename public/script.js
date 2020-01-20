
console.log('app started')

// const socket = new WebSocket('ws://localhost:3000/sub', 'pub.sp.nanomsg.org')
const subsocket = new WebSocket('ws://127.0.0.1:8081/sub', ['pub.sp.nanomsg.org'])
subsocket.binaryType = 'arraybuffer'

const repsocket = new WebSocket('ws://127.0.0.1:8081/req', ['rep.sp.nanomsg.org'])
repsocket.binaryType = 'arraybuffer'

var enc = new TextEncoder(); // always utf-8
var dec = new TextDecoder(); // always utf-8

function reqDATE() {
  console.log('Sending REP message: DATE')
  makeRequest("DATE")
}

function reqGREET() {
  console.log('Sending REP message: DATE')
  makeRequest("GREET")
}

function makeRequest(msg) {
  const length = msg.length || msg.byteLength;
  const data = new Uint8Array(length + 4)

  const reqIdHeader = new Uint8Array(4);
  window.crypto.getRandomValues(reqIdHeader);
  // the first bit HAS TO BE one, in order to get a response
  reqIdHeader[0] |= 1 << 7;

  data.set(reqIdHeader, 0);

  if (typeof msg === 'string' || msg instanceof String) {
    for (let i = 4; i < msg.length + 4; ++i) {
      data[i] = msg.charCodeAt(i - 4);
    }

  } else {
    data.set(msg, 4);
  }

  repsocket.send(data)

  let dataview = new DataView(reqIdHeader.buffer)
  console.log('REQ #', dataview.getUint32(0), ': ', msg)
}


subsocket.onerror = (event) => {
  console.log('SUB error')
  console.log(event)
}

subsocket.onopen = (event) => {
  console.log('SUB open')
  console.log(event)
}

subsocket.onmessage = (event) => {
  console.log(Date.now(), 'SUB message')
  console.log(event)
  console.log("Message: ", dec.decode(event.data))
}

repsocket.onerror = (event) => {
  console.log('REP error')
  console.log(event)
}

repsocket.onopen = (event) => {
  console.log('REP open')
  console.log(event)
}

repsocket.onmessage = (event) => {
  console.log(Date.now(), 'REP message')
  console.log(event)
  let dataview = new DataView(event.data.slice(0,4))
  console.log("REP #", dataview.getUint32(0), ": ", dec.decode(event.data.slice(4)))
}