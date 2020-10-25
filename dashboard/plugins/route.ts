export default (context: any) => {
  context.app.router.afterEach((_to: any, _from: any) => {
    if (context.app.store.state.dashboard.error !== null) {
      context.app.store.commit('dashboard/clearError')
    }
  })
}
