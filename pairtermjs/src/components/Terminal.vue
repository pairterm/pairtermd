
<template>
  <div id="terminal"></div>
</template>

<script>
import { Terminal } from 'xterm'
import * as fit from 'xterm/lib/addons/fit/fit'
import 'xterm/dist/xterm.css'

Terminal.applyAddon(fit)

export default {
  name: 'Terminal',
  mounted: () => {
    var loc = window.location
    var sock = new WebSocket(`ws://${loc.host}/pty`)
    sock.onerror = (e) => {
      console.log('socket error', e)
    }
    // wait for the socket to open before starting the terminal
    // or there will be ordering issues :/
    sock.onopen = (e) => {
      var term = new Terminal()
      term.open(document.getElementById('terminal'))
      term.fit()
      term.focus()
      term.on('title', (title) => { document.title = title })
      // pass data using base64 encoding
      // this is fragile: it will not work with non-ascii text!
      // the Go backend is correctly treating pty IO as opaque
      // byte arrays, while term.js uses javascript strings that
      // are utf16, while the pty is usually utf8.
      // I have some Go code that converts to utf16 before sending but
      // it's ugly and wrong. The right answer is to refactor term.js to use
      // ArrayBuffer with uint8 and convert runes on the fly on the client
      term.on('data', (data) => {
        sock.send(btoa(data))
      })
      sock.onmessage = (msg) => {
        term.write(atob(msg.data))
      }
    }
  }
}
</script>

<style scoped>

#terminal {
  height: 100%;
  width: 100%;
}
</style>
