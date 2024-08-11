async function dk4(type, payload) {
  let result = await dk4Action(type || "", JSON.stringify(payload || "{}"));
  let data = {};
  try {
    data = JSON.parse(result);
  } catch (e) {
    console.log("dk4 error", e, result);
    data = {code: 500, error: e};
  }
  
  return data;
}

$(async function(){
  async function readTemplate(name) {
    let resp = await dk4("readTemplateFile", "templates/" + name);
    return resp.data;
  }
  function InitStore() {
    const store = new Vuex.Store({
      state: {
        game: {
          processId: null,
          lang: "",
          version: "",
        },
        currentOrg: null,
        leadSeaman: null,
      },
      mutations: {
        setGame(state, game) {
          state.game = game
        },
        setPlayer(state, game) {
          state.currentOrg = game.currentOrg
          state.leadSeaman = game.leadSeaman
        },
        error(state, resp) {
          if (resp.code >= 500) {
            state.currentOrg = null
            state.leadSeaman = null
            state.game = {}
          }
          console.log("error", resp)
        },
      },
      actions: {
        async refresh({ commit }) {
          let resp = await dk4("refreshStatus", {});
          if (resp.error) {
            commit("error", resp);
            return;
          }

          commit("setGame", resp.data);
        },
        async getStatus({ commit }) {
          let resp = await dk4("getStatus", {});
          if (resp.error) {
            commit("error", resp);
            return;
          }

          commit("setGame", resp.data);
        },
        async getPlayerInfo({ commit }) {
          let resp = await dk4("getPlayerInfo", {});
          if (resp.error) {
            commit("error", resp);
            return;
          }

          commit("setPlayer", resp.data);
        },
        async call({ commit }, { method, data }) {
          let resp = await dk4(method, data);
          if (resp.error) {
            commit("error", resp);
            throw new Error(resp.error);
          }

          return resp.data;
        }
      },
    });
    return store;
  }
  async function LoadComponent(name) {
    let jsContent = await readTemplate(name + '.js');
    let htmlContent = await readTemplate(name + '.html');
    let component = eval('(' + jsContent + ')');
    component.template = htmlContent;
    return component;
  }

  const components = {
    MainLayout: await LoadComponent('MainLayout'),
    ListOrgPage: await LoadComponent('ListOrgPage'),
    ListArmadaPage: await LoadComponent('ListArmadaPage'),
    ListSeamanPage: await LoadComponent('ListSeamanPage'),
    ListShipPage: await LoadComponent('ListShipPage'),
    ListPortPage: await LoadComponent('ListPortPage'),
    MapPage: await LoadComponent('MapPage'),
  }
  _.each(components, (item) => {
    // 根据名字查找对应的组件
    item.components = _.reduce(item.components, (ret, name) => {
      ret[name] = components[name];
      return ret;
    }, {})
  });


  var router = new VueRouter.createRouter({
    history: VueRouter.createWebHashHistory(),
    routes: [
      {
        path: '/',
        component: components.MainLayout,
        children: [
          {
            path: 'org',
            component: components.ListOrgPage,
            name: 'org',
          },
          {
            path: 'armada',
            component: components.ListArmadaPage,
            name: 'armada',
          },
          {
            path: 'seaman',
            component: components.ListSeamanPage,
            name: 'seaman',
          },
          {
            path: 'ship',
            component: components.ListShipPage,
            name: 'ship',
          },
          {
            path: 'port',
            component: components.ListPortPage,
            name: 'port',
          },
          {
            path: 'map',
            component: components.MapPage,
            name: 'map',
          },
        ],
      }
    ],
    route404: {
      path: '/404',
      component: {
        template: '<div>not found</div>'
      }
    },
  });
  window.store = InitStore();

  window.app = Vue.createApp({});
  app.use(window.store);
  app.use(ElementPlus);
  app.use(router);

  app.mount('#app');
});