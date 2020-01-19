
console.log('app started')

// const socket = new WebSocket('ws://localhost:3000/sub', 'pub.sp.nanomsg.org')
const socket = new WebSocket('ws://localhost:3000/sub', ['rep.sp.nanomsg.org', 'pub.sp.nanomsg.org'])
socket.binaryType = 'arraybuffer'

socket.onerror = (event) => {
  console.log('error')
  console.log(event)
}

socket.onopen = (event) => {
  console.log('open')
  console.log(event)
}

socket.onmessage = (event) => {
  console.log(Date.now(), 'message')
  console.log(event)
}
