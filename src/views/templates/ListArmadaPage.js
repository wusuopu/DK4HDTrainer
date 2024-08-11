{
  components: {
  },
  data() {
    return {
      list: [],
    }
  },
  methods: {
  },
  async mounted() {
    let data = await this.$store.dispatch('call', {method: 'listArmada', data: ''})
    _.each(data, (item) => {
      if (item.Longitude) {
        let pos = ""
        if (item.Longitude > 0) {
          pos = "东"
        } else {
          pos = "西"
        }
        item.Longitude = `${pos}${Math.abs(item.Longitude.toFixed(2))}`
      }
      if (item.Latitude) {
        let pos = ""
        if (item.Latitude > 0) {
          pos = "南"
        } else {
          pos = "北"
        }
        item.Latitude = `${pos}${Math.abs(item.Latitude.toFixed(2))}`
      }
    });
    console.log(data);
    data = _.filter(data, "Name")
    this.list = data;
  },
}
