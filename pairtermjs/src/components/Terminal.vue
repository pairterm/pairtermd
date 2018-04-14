
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
    var host = window.location.host
    var sock = new WebSocket(`ws://${host}/pty`)
    sock.onerror = (e) => { console.log('socket error', e) }
    // wait for the socket to open before starting the terminal
    // or there will be ordering issues :/
    sock.onopen = (e) => {
      var term = new Terminal()
      term.open(document.getElementById('terminal'))
      term.fit()
      term.focus()
      term.on('title', (title) => { document.title = title })
      term.on('data', (data) => { sock.send(data) })
      sock.onmessage = (msg) => { term.write(msg.data) }
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
