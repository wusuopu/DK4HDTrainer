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
    let data = await this.$store.dispatch('call', {method: 'listPort', data: ''})
    console.log(data);
    this.list = data;
    this.loaded = true
  },
}
