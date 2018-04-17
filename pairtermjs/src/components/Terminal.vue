
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
    var sock = new WebSocket(`ws://${window.location.host}/pty?cols=${term.cols}&rows=${term.rows}`)

    window.addEventListener('resize', () => { term.fit(); return true })

    sock.onerror = (e) => { console.log('socket error', e) }
    sock.onopen = (e) => { term.attach(sock, true, true) }
  },
  methods: {
    initTerminal: () => {
      var term = new Terminal()
      term.open(document.getElementById('terminal'))
      term.fit()
      term.focus()
      term.on('title', (title) => { document.title = title })
      term.on('resize', (size) => {
        console.log('resized to ', size.cols, size.rows)
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
