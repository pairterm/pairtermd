
<template>
  <div id="terminal"></div>
</template>

<script>
import { Terminal } from 'xterm'
import * as fit from 'xterm/lib/addons/fit/fit'
import * as attach from 'xterm/lib/addons/attach/attach'
import 'xterm/dist/xterm.css'

Terminal.applyAddon(fit)
Terminal.applyAddon(attach)

export default {
  name: 'Terminal',
  mounted: () => {
    var host = window.location.host
    var sock = new WebSocket(`ws://${host}/pty`)
    var term = new Terminal()
    sock.onerror = (e) => { console.log('socket error', e) }
    // wait for the socket to open before starting the terminal
    // or there will be ordering issues :/
    sock.onopen = (e) => {
      term.open(document.getElementById('terminal'))
      term.fit()
      term.focus()
      term.on('title', (title) => { document.title = title })
      term.attach(sock, true, true)
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
