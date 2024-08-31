{
  components: {
  },
  data() {
    return {
      loaded: false,
      list: [],
    }
  },
  methods: {
    async handleChangeMoney(orgId) {
      try {
        await this.$store.dispatch('call', {method: 'minusOrgMoney', data: {id: orgId}})
        ElementPlus.ElNotification({ title: '成功', message: '修改成功', type: 'success', })
        await this.fetchList()
      } catch (error) {
        ElementPlus.ElNotification({ title: '错误', message: '修改失败', type: 'error', })
      }
    },
    async fetchList() {
      this.loaded = false
      let data = await this.$store.dispatch('call', {method: 'listOrganization', data: ''})
      _.each(data, (item) => {
        item.TotalAreaValue = 0
        _.each(item.AreaValues, (v, index) => {
          item[`AreaValue${index}`] = v;
          item.TotalAreaValue += v;
        });
      });
      console.log(data)
      this.list = data
      this.loaded = true
    }
  },
  async mounted() {
    await this.fetchList()
  },
}
