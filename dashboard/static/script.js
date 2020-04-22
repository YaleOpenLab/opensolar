(() => {

  'use-strict'

  const themeSwiter = {

    init: function() {
      this.wrapper = document.getElementById('theme-switcher-wrapper')
      this.button = document.getElementById('theme-switcher-button')
      this.slider = this.wrapper.querySelector('[data-slider]')
      this.mode = this.wrapper.querySelectorAll('[data-mode]')
      this.modes = ['mode-light', 'mode-dark']
      this.theme = this.wrapper.querySelectorAll('[data-theme]')
      this.themes = ['theme-orange', 'theme-purple', 'theme-green', 'theme-blue', 'theme-yellow', 'theme-red', 'theme-teal', 'theme-pink']
      this.events()
      this.start()
    },
    
    events: function() {
      this.button.addEventListener('click', this.displayed.bind(this), false)
      this.mode.forEach(mode => mode.addEventListener('click', this.modeed.bind(this), false))
      this.theme.forEach(theme => theme.addEventListener('click', this.themed.bind(this), false))
    },

    start: function() {
      let mode = this.modes[Math.floor(Math.random() * this.modes.length)]
      let theme = this.themes[Math.floor(Math.random() * this.themes.length)]
      document.body.classList.add('mode-dark', theme)
    },

    displayed: function() {
      (this.wrapper.classList.contains('is-open'))
        ? this.wrapper.classList.remove('is-open')
        : this.wrapper.classList.add('is-open')
    },

    modeed: function(e) {
      this.slider.classList.toggle('is-change')
      this.modes.forEach(mode => {
        if(document.body.classList.contains(mode))
          document.body.classList.remove(mode)
      })
      return document.body.classList.add(`mode-${e.currentTarget.dataset.mode}`)
    },

    themed: function(e) {
      this.themes.forEach(theme => {
        if(document.body.classList.contains(theme))
          document.body.classList.remove(theme)
      })
      return document.body.classList.add(`theme-${e.currentTarget.dataset.theme}`)
    }

  }

  themeSwiter.init()

})()