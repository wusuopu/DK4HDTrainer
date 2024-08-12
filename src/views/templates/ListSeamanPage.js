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
    let data = await this.$store.dispatch('call', {method: 'listSeaman', data: ''})
    _.each(data, (item) => {
      _.each(item.Metrics, (v, index) => {
        item[`Metrics${index}`] = v;
      });
      _.each(item.Exps, (v, index) => {
        item[`Exps${index}`] = v;
      });
    });
    console.log(data);
    this.list = data;
    this.loaded = true
  },
}
