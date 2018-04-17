
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
  mounted () {
    var term = this.initTerminal()
    window.addEventListener('resize', () => { term.fit(); return true })
  },
  methods: {
    initTerminal: () => {
      let term = new Terminal()
      term.open(document.getElementById('terminal'))
      term.fit()
      term.focus()

      let sock = new WebSocket(`ws://${window.location.host}/pty?cols=${term.cols}&rows=${term.rows}`)
      sock.onerror = (e) => { console.log('socket error', e) }
      sock.onopen = (e) => { term.attach(sock, true, true) }

      term.on('resize', (size) => {
        let payload = {
          type: 'pt_resize',
          payload: {
            cols: size.cols,
            rows: size.rows
          }
        }
        sock.send(JSON.stringify(payload))
      })

      return term
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
