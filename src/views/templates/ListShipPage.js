{
  components: {
  },
  data() {
    return {
      loaded: false,
      list: [],
      gunOptions: ["散弹", "曲射", "加农曲射", "加农", "重加农", "连射"]
    }
  },
  methods: {
  },
  async mounted() {
    let data = await this.$store.dispatch('call', {method: 'listShip', data: ''})
    _.each(data, (item) => {
      item.GunName = this.gunOptions[item.Gun] || "";
    });
    data = _.filter(data, 'Name');
    console.log(data);
    this.list = data;
    this.loaded = true
  },
}
