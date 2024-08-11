{
  components: [],
  data() {
    return {
      timerId: null,
      menuIndex: '',
      menus: [
        {name: "势力信息", path: "/org",},
        {name: "舰队信息", path: "/armada",},
        {name: "海员信息 ", path: "/seaman",},
        {name: "船只信息", path: "/ship",},
        {name: "港口信息", path: "/port",},
        {name: "霸者之证", path: "/map",},
      ],
      lockFoodFlag: false,
      lockFatigueFlag: false,
      lockShipFlag: false,
    }
  },
  computed: {
    language() {
      return {
        sc: '简体中文',
        tc: '繁體中文',
        en: 'English',
      }[_.get(this.$store.state.game, 'lang')];
    },
  },
  methods: {
    async refreshStatus() {
      await this.$store.dispatch('refresh');
    },
    async checkStatus() {
      await this.$store.dispatch('getStatus');
    },
    async getPlayerInfo() {
      await this.$store.dispatch("getPlayerInfo");
    },
    handleMenuSelect(key) {
      this.$router.push(key);
    },
    async addMoney() {
      try {
        await this.$store.dispatch('call', {method: 'addMoney', data: ''})
        ElementPlus.ElNotification({ title: '成功', message: '修改成功', type: 'success', })
      } catch (error) {
        ElementPlus.ElNotification({ title: '错误', message: '修改失败', type: 'error', })
      }
    },
    async turnOnAllPorts() {
      try {
        await this.$store.dispatch('call', {method: 'turnOnAllPorts', data: ''})
        ElementPlus.ElNotification({ title: '成功', message: '修改成功', type: 'success', })
      } catch (error) {
        ElementPlus.ElNotification({ title: '错误', message: '修改失败', type: 'error', })
      }
    },
    async enhanceSeaman() {
      try {
        await this.$store.dispatch('call', {method: 'enhanceSeaman', data: ''})
        ElementPlus.ElNotification({ title: '成功', message: '修改成功', type: 'success', })
      } catch (error) {
        ElementPlus.ElNotification({ title: '错误', message: '修改失败', type: 'error', })
      }
    },
    async enhanceShip() {
      try {
        await this.$store.dispatch('call', {method: 'enhanceShip', data: ''})
        ElementPlus.ElNotification({ title: '成功', message: '修改成功', type: 'success', })
      } catch (error) {
        ElementPlus.ElNotification({ title: '错误', message: '修改失败', type: 'error', })
      }
    },
    async toggleLockFlag(key, value) {
      try {
        await this.$store.dispatch('call', {method: 'toggleLockFlag', data: {key: key, value: value}})
        ElementPlus.ElNotification({ title: '成功', message: '修改成功', type: 'success', })
      } catch (error) {
        ElementPlus.ElNotification({ title: '错误', message: '修改失败', type: 'error', })
      }
    },
  },
  async mounted() {
    await this.checkStatus()

    this.timerId = setInterval(async () => {
      if (!this.$store.state.game.processId) {
        await this.refreshStatus()
      }
      if (this.$store.state.game.processId) {
        await this.getPlayerInfo()
      }
    }, 2000);
  },
  destroyed() {
    this.timerId && clearInterval(this.timerId);
  },
}
