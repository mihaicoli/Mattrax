export default async function (context: any) {
  if (context.store.state.authentication.user === null) {
    await context.store.dispatch('authentication/populateUserInfomation')
  }

  const authenticated: boolean = await context.store.dispatch(
    'authentication/isAuthenticated'
  )
  if (!authenticated) {
    let params = ''
    if (context.route.fullPath !== '/') {
      params += '?redirect_to=' + encodeURIComponent(context.route.fullPath)
    }
    context.app.router.push('/login' + params)
  }
}
