export default (context: any) => {
  context.app.router.afterEach((_to: any, _from: any) => {
    context.app.store.commit('dashboard/clearError')
  })
}
