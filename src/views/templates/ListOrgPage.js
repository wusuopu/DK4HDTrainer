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
  },
  async mounted() {
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
  },
}
